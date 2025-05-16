// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/rlp"
)

const (
	LegacyTxType    = 0x00
	SponsoredTxType = 100
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

var (
	ErrInvalidSig                 = errors.New("invalid transaction v, r, s values")
	ErrUnexpectedProtection       = errors.New("transaction type does not supported EIP-155 protected signatures")
	ErrShortTypedTx               = errors.New("typed transaction too short")
	ErrTxTypeNotSupported         = errors.New("transaction type not supported")
	ErrSamePayerSenderSponsoredTx = errors.New("payer = sender in sponsored transaction")
	ErrInvalidTxType              = errors.New("transaction type not valid in this context")
	errEmptyTypedTx               = errors.New("empty typed transaction bytes")

	// ErrExpiredSponsoredTx is returned if the sponsored transaction is expired.
	ErrExpiredSponsoredTx = errors.New("sponsored transaction is expired")

	skipNonceDestinationAddress = map[string]bool{
		common.TomoXAddr:                         true,
		common.TradingStateAddr:                  true,
		common.TomoXLendingAddress:               true,
		common.TomoXLendingFinalizedTradeAddress: true,
	}
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(tx *Transaction) Signer {
	if tx.Type() == SponsoredTxType {
		if tx.ChainId() != nil {
			return NewMikoSigner(tx.ChainId())
		}
		return HomesteadSigner{}
	}

	v, _, _ := tx.RawSignatureValues()
	if v != nil && v.Sign() != 0 && isProtectedV(v) {
		return NewEIP155Signer(deriveChainId(v))
	}
	return HomesteadSigner{}
}

type Transaction struct {
	inner TxData    // Consensus contents of a transaction
	time  time.Time // Time first seen locally (spam avoidance)
	// caches
	hash  atomic.Value
	size  atomic.Value
	from  atomic.Value
	payer atomic.Value
}

// TxData is the underlying data of a transaction.
// This is implemented by LegacyTx and AccessListTx.
type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	value() *big.Int
	nonce() uint64
	to() *common.Address
	accessList() AccessList

	expiredTime() uint64
	rawPayerSignatureValues() (v, r, s *big.Int)

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)

	encode(*bytes.Buffer) error
	decode([]byte) error
}

// NewTx creates a new transaction.
func NewTx(inner TxData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction) setDecoded(inner TxData, size uint64) {
	tx.inner = inner
	tx.time = time.Now()
	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func sanityCheckSignature(v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return ErrUnexpectedProtection
	}
	var plainV byte
	if isProtectedV(v) {
		chainID := deriveChainId(v).Uint64()
		plainV = byte(v.Uint64() - 35 - 2*chainID)
	} else if maybeProtected {
		plainV = byte(v.Uint64() - 27)
	} else {
		plainV = byte(v.Uint64())
	}
	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}
	return nil
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return tx.inner.chainID()
}

// Type returns type of transaction
func (tx *Transaction) Type() uint8 {
	return tx.inner.txType()
}

// Data returns data of transaction
func (tx *Transaction) Data() []byte {
	return tx.inner.data()
}

// Gas returns gas limit of transaciton
func (tx *Transaction) Gas() uint64 {
	return tx.inner.gas()
}

// GasPrice returns gas price in wei of transaction
func (tx *Transaction) GasPrice() *big.Int {
	return tx.inner.gasPrice()
}

// Value returns amount native coin of transaction
func (tx *Transaction) Value() *big.Int {
	return tx.inner.value()
}

// Nonce returns nonce of send on transaction
func (tx *Transaction) Nonce() uint64 {
	return tx.inner.nonce()
}
func (tx *Transaction) AccessList() AccessList { return nil }

// From returns creator of transaction
func (tx *Transaction) From() *common.Address {
	v, _, _ := tx.RawSignatureValues()
	if v != nil {
		signer := deriveSigner(tx)
		if f, err := Sender(signer, tx); err != nil {
			return nil
		} else {
			return &f
		}
	} else {
		return nil
	}
}

// To returns recipient of transaction
func (tx *Transaction) To() *common.Address {
	ito := tx.inner.to()
	if ito == nil {
		return nil
	}
	cpy := *ito
	return &cpy
}

// ExpiredTime returns the expired time of the sponsored transaction
func (tx *Transaction) ExpiredTime() uint64 {
	return tx.inner.expiredTime()
}

// Cost returns gas * gasPrice + value (amount)
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
	total.Add(total, tx.Value())
	return total
}

// RawSignatureValues returns the V, R, S payer signature values of the transaction.
// The return values should not be modified by the caller.
func (tx *Transaction) RawPayerSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawPayerSignatureValues()
}

// TRC21Cost returns amount + gasprice * gaslimit.
func (tx *Transaction) TRC21Cost() *big.Int {
	total := new(big.Int).Mul(common.TRC21GasPrice, new(big.Int).SetUint64(tx.inner.copy().gas()))
	total.Add(total, tx.Value())
	return total
}

// return R, S, V signature values of transaction
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawSignatureValues()
}

// return the true RLP encoded storage size of transaction, either by encoding and returning it, or returning a previously cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.inner)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	switch tx := tx.inner.(type) {
	case *LegacyTx:
		return tx.V != nil && isProtectedV(tx.V)
	default:
		return true
	}
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 are considered unprotected
	return true
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return tx.inner.encode(w)
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	if tx.Type() == LegacyTxType {
		return rlp.Encode(w, &tx.inner)
	}
	buffer := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buffer)
	buffer.Reset()
	if err := tx.encodeTyped(buffer); err != nil {
		return err
	}
	return rlp.Encode(w, buffer.Bytes())
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	kind, size, err := s.Kind()
	switch {
	case err != nil:
		// fmt.Println(">>>>>>> DecodeRLP", err)
		return err
	case kind == rlp.List:
		// fmt.Println(">>>>>>> DecodeRLP:kind == rlp.List", kind, size)
		// It's a legacy transaction.
		var inner LegacyTx
		err := s.Decode(&inner)
		// fmt.Println(">>>>>>> DecodeRLP:kind == rlp.List::s.Decode(&inner)", err)
		if err == nil {
			// fmt.Println(">>>>>>> DecodeRLP:kind == rlp.List::s.Decode(&inner) setDecoded")
			tx.setDecoded(&inner, uint64(rlp.ListSize(size)))
		}
		return err
	case kind == rlp.String:
		// fmt.Println(">>>>>>> DecodeRLP:kind == rlp.String", kind, size)
		// It's an EIP-2718 typed TX envelope.
		var b []byte
		if b, err = s.Bytes(); err != nil {
			return err
		}
		inner, err := tx.decodeTyped(b)
		if err == nil {
			tx.setDecoded(inner, uint64(len(b)))
		}
		return err
	default:
		return rlp.ErrExpectedList
	}
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction) decodeTyped(b []byte) (TxData, error) {
	if len(b) == 0 {
		return nil, errEmptyTypedTx
	}
	var inner TxData
	switch b[0] {

	case SponsoredTxType:
		inner = new(SponsoredTx)

	default:
		return nil, ErrTxTypeNotSupported
	}
	err := inner.decode(b[1:])
	return inner, err
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *Transaction) CacheHash() {
	v := rlpHash(tx)
	tx.hash.Store(v)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *Transaction) AsMessage(s Signer, balanceFee *big.Int, number *big.Int, checkNonce bool) (Message, error) {
	msg := Message{
		nonce:           tx.Nonce(),
		gasLimit:        tx.Gas(),
		gasPrice:        new(big.Int).Set(tx.GasPrice()),
		to:              tx.To(),
		amount:          tx.Value(),
		data:            tx.Data(),
		checkNonce:      checkNonce,
		balanceTokenFee: balanceFee,
		expiredTime:     tx.ExpiredTime(),
	}
	var err error
	msg.from, err = Sender(s, tx)
	if err != nil {
		return Message{}, err
	}
	if tx.Type() == SponsoredTxType {
		msg.payer, err = Payer(s, tx)
		if err != nil {
			return Message{}, err
		}
		fmt.Println("=---= payer", msg.payer.String(), "sender", msg.from.String())
		if msg.payer == msg.from {
			// Reject sponsored transaction with identical payer and sender
			return Message{}, ErrSamePayerSenderSponsoredTx
		}
		return msg, nil
	} else {
		msg.payer = msg.from
	}
	if balanceFee != nil {
		if number.Cmp(common.TIPTRC21FeeBlock) > 0 {
			msg.gasPrice = common.TRC21GasPrice
		} else {
			msg.gasPrice = common.TRC21GasPriceBefore
		}
	}
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	if r == nil || s == nil || v == nil {
		return nil, fmt.Errorf("%w: r: %s, s: %s, v: %s", ErrInvalidSig, r, s, v)
	}
	cpy := tx.inner.copy()
	cpy.setSignatureValues(signer.ChainID(), v, r, s)
	return &Transaction{inner: cpy, time: tx.time}, nil
}

func (tx *Transaction) IsSpecialTransaction() bool {
	if tx.To() == nil {
		return false
	}
	return tx.To().String() == common.RandomizeSMC || tx.To().String() == common.BlockSigners
}

func (tx *Transaction) IsTradingTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.TomoXAddr {
		return false
	}

	return true
}

func (tx *Transaction) IsLendingTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.TomoXLendingAddress {
		return false
	}
	return true
}

func (tx *Transaction) IsLendingFinalizedTradeTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.TomoXLendingFinalizedTradeAddress {
		return false
	}
	return true
}

func (tx *Transaction) IsSkipNonceTransaction() bool {
	if tx.To() == nil {
		return false
	}
	if skip := skipNonceDestinationAddress[tx.To().String()]; skip {
		return true
	}
	return false
}

func (tx *Transaction) IsSigningTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.BlockSigners {
		return false
	}

	method := common.ToHex(tx.Data()[0:4])

	if method != common.SignMethod {
		return false
	}

	if len(tx.Data()) != (32*2 + 4) {
		return false
	}

	return true
}

func (tx *Transaction) IsVotingTransaction() (bool, *common.Address) {
	if tx.To() == nil {
		return false, nil
	}
	b := (tx.To().String() == common.MasternodeVotingSMC)

	if !b {
		return b, nil
	}

	method := common.ToHex(tx.Data()[0:4])
	if b = (method == common.VoteMethod); b {
		addr := tx.Data()[len(tx.Data())-20:]
		m := common.BytesToAddress(addr)
		return b, &m
	}

	if b = (method == common.UnvoteMethod); b {
		addr := tx.Data()[len(tx.Data())-32-20 : len(tx.Data())-32]
		m := common.BytesToAddress(addr)
		return b, &m
	}

	if b = (method == common.ProposeMethod); b {
		addr := tx.Data()[len(tx.Data())-20:]
		m := common.BytesToAddress(addr)
		return b, &m
	}

	if b = (method == common.ResignMethod); b {
		addr := tx.Data()[len(tx.Data())-20:]
		m := common.BytesToAddress(addr)
		return b, &m
	}

	return b, nil
}

func (tx *Transaction) IsTomoXApplyTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.TomoXListingSMC.String() {
		return false
	}

	method := common.ToHex(tx.Data()[0:4])

	if method != common.TomoXApplyMethod {
		return false
	}

	// 4 bytes for function name
	// 32 bytes for 1 parameter
	if len(tx.Data()) != (32 + 4) {
		return false
	}
	return true
}

func (tx *Transaction) IsTomoZApplyTransaction() bool {
	if tx.To() == nil {
		return false
	}

	if tx.To().String() != common.TRC21IssuerSMC.String() {
		return false
	}

	method := common.ToHex(tx.Data()[0:4])
	if method != common.TomoZApplyMethod {
		return false
	}

	// 4 bytes for function name
	// 32 bytes for 1 parameter
	if len(tx.Data()) != (32 + 4) {
		return false
	}

	return true
}

func (tx *Transaction) String() string {
	var from, to, payer string
	v, r, s := tx.RawSignatureValues()
	if v != nil {
		signer := deriveSigner(tx)
		f, err := Sender(signer, tx)
		if err != nil {
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.To() == nil {
		to = "nil"
	} else {
		to = tx.To().Hex()
	}

	enc, _ := rlp.EncodeToBytes(&tx.inner)
	// Handle sponsored transaction
	if tx.Type() == SponsoredTxType {
		if sponsoredTx, ok := tx.inner.(*SponsoredTx); ok {
			pv, pr, ps := sponsoredTx.rawPayerSignatureValues()
			signer := deriveSigner(tx)
			addr, _ := Payer(signer, tx)

			payer = fmt.Sprintf(`
	Payer:    %s
	PayerV:   %#x
	PayerR:   %#x
	PayerS:   %#x
	Expired:  %d`, addr.String(), pv, pr, ps, sponsoredTx.expiredTime())
		}
	}

	return fmt.Sprintf(`
	TX(%x)
	Type:     %v
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit: %#x
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x%s
	Hex:      %x`,
		tx.Hash(),
		tx.Type(),
		tx.To() == nil,
		from,
		to,
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.Value(),
		tx.Data(),
		v,
		r,
		s,
		payer, // Add payer info if sponsored tx
		enc,
	)
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// TxDifference returns a new set t which is the difference between a to b.
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// TxByPrice implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByPrice struct {
	txs        Transactions
	payersSwap map[common.Address]*big.Int
}

func (s TxByPrice) Len() int { return len(s.txs) }
func (s TxByPrice) Less(i, j int) bool {
	i_price := s.txs[i].inner.gasPrice()
	if s.txs[i].To() != nil {
		if _, ok := s.payersSwap[*s.txs[i].To()]; ok {
			i_price = common.TRC21GasPrice
		}
	}

	j_price := s.txs[j].inner.gasPrice()
	if s.txs[j].To() != nil {
		if _, ok := s.payersSwap[*s.txs[j].To()]; ok {
			j_price = common.TRC21GasPrice
		}
	}
	return i_price.Cmp(j_price) > 0
}
func (s TxByPrice) Swap(i, j int) { s.txs[i], s.txs[j] = s.txs[j], s.txs[i] }

func (s *TxByPrice) Push(x interface{}) {
	s.txs = append(s.txs, x.(*Transaction))
}

func (s *TxByPrice) Pop() interface{} {
	old := s.txs
	n := len(old)
	x := old[n-1]
	s.txs = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce struct {
	txs    map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads  TxByPrice                       // Next transaction for each unique account (price heap)
	signer Signer                          // Signer for the set of transactions
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.

// It also classifies special txs and normal txs
func NewTransactionsByPriceAndNonce(signer Signer, txs map[common.Address]Transactions, signers map[common.Address]struct{}, payersSwap map[common.Address]*big.Int) (*TransactionsByPriceAndNonce, Transactions) {
	// Initialize a price based heap with the head transactions
	heads := TxByPrice{}
	heads.payersSwap = payersSwap
	specialTxs := Transactions{}
	for _, accTxs := range txs {
		from, _ := Sender(signer, accTxs[0])
		var normalTxs Transactions
		lastSpecialTx := -1
		if len(signers) > 0 {
			if _, ok := signers[from]; ok {
				for i, tx := range accTxs {
					if tx.IsSpecialTransaction() {
						lastSpecialTx = i
					}
				}
			}
		}
		if lastSpecialTx >= 0 {
			for i := 0; i <= lastSpecialTx; i++ {
				specialTxs = append(specialTxs, accTxs[i])
			}
			normalTxs = accTxs[lastSpecialTx+1:]
		} else {
			normalTxs = accTxs
		}
		if len(normalTxs) > 0 {
			heads.txs = append(heads.txs, normalTxs[0])
			// Ensure the sender address is from the signer
			txs[from] = normalTxs[1:]
		}
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce{
		txs:    txs,
		heads:  heads,
		signer: signer,
	}, specialTxs
}

// Peek returns the next transaction by price.
func (t *TransactionsByPriceAndNonce) Peek() *Transaction {
	if len(t.heads.txs) == 0 {
		return nil
	}
	return t.heads.txs[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce) Shift() {
	acc, _ := Sender(t.signer, t.heads.txs[0])
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads.txs[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriceAndNonce) Pop() {
	heap.Pop(&t.heads)
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to              *common.Address
	from            common.Address
	nonce           uint64
	amount          *big.Int
	gasLimit        uint64
	gasPrice        *big.Int
	data            []byte
	checkNonce      bool
	balanceTokenFee *big.Int

	payer       common.Address
	expiredTime uint64
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool, balanceTokenFee *big.Int, accessList AccessList,
) Message {
	if balanceTokenFee != nil {
		gasPrice = common.TRC21GasPrice
	}
	return Message{
		from:            from,
		to:              to,
		nonce:           nonce,
		amount:          amount,
		gasLimit:        gasLimit,
		gasPrice:        gasPrice,
		data:            data,
		checkNonce:      checkNonce,
		balanceTokenFee: balanceTokenFee,
		payer:           from,
		expiredTime:     0,
	}
}

func (m Message) From() common.Address      { return m.from }
func (m Message) BalanceTokenFee() *big.Int { return m.balanceTokenFee }
func (m Message) To() *common.Address       { return m.to }
func (m Message) GasPrice() *big.Int        { return m.gasPrice }
func (m Message) Value() *big.Int           { return m.amount }
func (m Message) Gas() uint64               { return m.gasLimit }
func (m Message) Nonce() uint64             { return m.nonce }
func (m Message) Data() []byte              { return m.data }
func (m Message) CheckNonce() bool          { return m.checkNonce }
func (m Message) Payer() common.Address     { return m.payer }
func (m Message) ExpiredTime() uint64       { return m.expiredTime }

// copyAddressPtr copies an address.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	if tx.Type() == LegacyTxType {
		return rlp.EncodeToBytes(tx.inner)
	}
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		fmt.Println("UnmarshalBinary:legacy transaction")
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		fmt.Println("UnmarshalBinary:err", err)
		if err != nil {
			return err
		}
		tx.setDecoded(&data, uint64(len(b)))
		return nil
	}
	fmt.Println("UnmarshalBinary:typed transaction")
	// It's an EIP2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		fmt.Println("UnmarshalBinary:err", err)
		return err
	}
	tx.setDecoded(inner, uint64(len(b)))
	return nil
}
