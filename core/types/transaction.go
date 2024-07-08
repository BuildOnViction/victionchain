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

	"github.com/tomochain/tomochain/common"
	cmath "github.com/tomochain/tomochain/common/math"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/rlp"
)

// Transaction types.
const (
	LegacyTxType     = 0x00
	DynamicFeeTxType = 0x02
)

var (
	ErrInvalidSig           = errors.New("invalid transaction v, r, s values")
	ErrUnexpectedProtection = errors.New("transaction type does not supported EIP-155 protected signatures")
	ErrTxTypeNotSupported   = errors.New("transaction type not supported")
	ErrGasFeeCapTooLow      = errors.New("fee cap less than base fee")
	errShortTypedTx         = errors.New("typed transaction too short")
	errInvalidYParity       = errors.New("'yParity' field must be 0 or 1")
	errVYParityMismatch     = errors.New("'v' and 'yParity' fields do not match")
	errVYParityMissing      = errors.New("missing 'yParity' or 'v' field in transaction")
)

var (
	skipNonceDestinationAddress = map[string]bool{
		common.TomoXAddr:                         true,
		common.TradingStateAddr:                  true,
		common.TomoXLendingAddress:               true,
		common.TomoXLendingFinalizedTradeAddress: true,
	}
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(V *big.Int) Signer {
	if V.Sign() != 0 && isProtectedV(V) {
		return NewEIP155Signer(deriveChainId(V))
	} else {
		return HomesteadSigner{}
	}
}

type Transaction struct {
	inner TxData // Consensus contents of a transaction
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

// NewTx creates a new transaction.
func NewTx(inner TxData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// TxData is the underlying inner of a transaction.
//
// This is implemented by LegacyTx.
type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	gasTipCap() *big.Int
	gasFeeCap() *big.Int
	value() *big.Int
	nonce() uint64
	to() *common.Address

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)

	// effectiveGasPrice computes the gas price paid by the transaction, given
	// the inclusion block baseFee.
	//
	// Unlike other TxData methods, the returned *big.Int should be an independent
	// copy of the computed value, i.e. callers are allowed to mutate the result.
	// Method implementations can use 'dst' to store the result.
	effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int

	encode(*bytes.Buffer) error // used to RLP encode the typed transaction contents
	decode([]byte) error        // used to RLP decode the typed transaction contents
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	if tx.Type() == LegacyTxType {
		return rlp.Encode(w, tx.inner)
	}
	// It's an EIP-2718 typed TX envelope. Shouldn't reach here for now
	// Getting buffer from a sync pool
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buf)
	buf.Reset()
	if err := tx.encodeTyped(buf); err != nil {
		return err
	}
	return rlp.Encode(w, buf.Bytes())
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction) encodeTyped(w *bytes.Buffer) error {
	// Encode and write the transaction type to buffer first
	w.WriteByte(tx.Type())
	// Then concatenate the RLP-encoded transaction content later
	return tx.inner.encode(w)
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

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	kind, size, err := s.Kind()
	switch {
	case err != nil:
		return err
	case kind == rlp.List:
		// It's a legacy transaction.
		var inner LegacyTx
		err := s.Decode(&inner)
		if err == nil {
			tx.setDecoded(&inner, int(rlp.ListSize(size)))
		}
		return err
	case kind == rlp.String:
		// It's an EIP-2718 typed TX envelope.
		var b []byte
		if b, err = s.Bytes(); err != nil {
			return err
		}
		inner, err := tx.decodeTyped(b)
		if err == nil {
			tx.setDecoded(inner, len(b))
		}
		return err
	default:
		return rlp.ErrExpectedList
	}
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP-2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return err
		}
		tx.setDecoded(&data, len(b))
		return nil
	}
	// It's an EIP-2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner, len(b))
	return nil
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction) decodeTyped(b []byte) (TxData, error) {
	if len(b) <= 1 {
		return nil, errShortTypedTx
	}
	var inner TxData
	switch b[0] {
	default:
		return nil, ErrTxTypeNotSupported
	}
	err := inner.decode(b[1:])
	return inner, err
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction) setDecoded(inner TxData, size int) {
	tx.inner = inner
	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered unprotected
	return true
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
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = byte(v.Uint64() - 27)
	} else {
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = byte(v.Uint64())
	}
	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}

	return nil
}

// Type returns the transaction type.
func (tx *Transaction) Type() uint8 {
	return tx.inner.txType()
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return tx.inner.chainID()
}

// Protected says whether the transaction is replay-protected.
func (tx *Transaction) Protected() bool {
	switch tx := tx.inner.(type) {
	case *LegacyTx:
		return tx.V != nil && isProtectedV(tx.V)
	default:
		return true
	}
}

func (tx *Transaction) Data() []byte       { return common.CopyBytes(tx.inner.data()) }
func (tx *Transaction) Gas() uint64        { return tx.inner.gas() }
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }
func (tx *Transaction) Value() *big.Int    { return new(big.Int).Set(tx.inner.value()) }
func (tx *Transaction) Nonce() uint64      { return tx.inner.nonce() }
func (tx *Transaction) CheckNonce() bool   { return true }

// GasTipCap returns the gasTipCap per gas of the transaction.
func (tx *Transaction) GasTipCap() *big.Int { return new(big.Int).Set(tx.inner.gasTipCap()) }

// GasFeeCap returns the fee cap per gas of the transaction.
func (tx *Transaction) GasFeeCap() *big.Int { return new(big.Int).Set(tx.inner.gasFeeCap()) }

// GasFeeCapCmp compares the fee cap of two transactions.
func (tx *Transaction) GasFeeCapCmp(other *Transaction) int {
	return tx.inner.gasFeeCap().Cmp(other.inner.gasFeeCap())
}

// GasFeeCapIntCmp compares the fee cap of the transaction against the given fee cap.
func (tx *Transaction) GasFeeCapIntCmp(other *big.Int) int {
	return tx.inner.gasFeeCap().Cmp(other)
}

// GasTipCapCmp compares the gasTipCap of two transactions.
func (tx *Transaction) GasTipCapCmp(other *Transaction) int {
	return tx.inner.gasTipCap().Cmp(other.inner.gasTipCap())
}

// GasTipCapIntCmp compares the gasTipCap of the transaction against the given gasTipCap.
func (tx *Transaction) GasTipCapIntCmp(other *big.Int) int {
	return tx.inner.gasTipCap().Cmp(other)
}

// EffectiveGasTip returns the effective miner gasTipCap for the given base fee.
// Note: if the effective gasTipCap is negative, this method returns both error
// the actual negative value, _and_ ErrGasFeeCapTooLow
func (tx *Transaction) EffectiveGasTip(baseFee *big.Int) (*big.Int, error) {
	if baseFee == nil {
		return tx.GasTipCap(), nil
	}
	var err error
	gasFeeCap := tx.GasFeeCap()
	if gasFeeCap.Cmp(baseFee) == -1 {
		err = ErrGasFeeCapTooLow
	}
	return cmath.BigMin(tx.GasTipCap(), gasFeeCap.Sub(gasFeeCap, baseFee)), err
}

// EffectiveGasTipValue is identical to EffectiveGasTip, but does not return an
// error in case the effective gasTipCap is negative
func (tx *Transaction) EffectiveGasTipValue(baseFee *big.Int) *big.Int {
	effectiveTip, _ := tx.EffectiveGasTip(baseFee)
	return effectiveTip
}

// EffectiveGasTipCmp compares the effective gasTipCap of two transactions assuming the given base fee.
func (tx *Transaction) EffectiveGasTipCmp(other *Transaction, baseFee *big.Int) int {
	if baseFee == nil {
		return tx.GasTipCapCmp(other)
	}
	return tx.EffectiveGasTipValue(baseFee).Cmp(other.EffectiveGasTipValue(baseFee))
}

// EffectiveGasTipIntCmp compares the effective gasTipCap of a transaction to the given gasTipCap.
func (tx *Transaction) EffectiveGasTipIntCmp(other *big.Int, baseFee *big.Int) int {
	if baseFee == nil {
		return tx.GasTipCapIntCmp(other)
	}
	return tx.EffectiveGasTipValue(baseFee).Cmp(other)
}

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *common.Address {
	if tx.inner.to() == nil {
		return nil
	}
	to := *tx.inner.to()
	return &to
}

func (tx *Transaction) From() *common.Address {
	if v, _, _ := tx.RawSignatureValues(); v != nil {
		signer := deriveSigner(v)
		if f, err := Sender(signer, tx); err != nil {
			return nil
		} else {
			return &f
		}
	} else {
		return nil
	}
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

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.inner)
	// For typed transactions later, the encoding also includes the leading type byte.
	if tx.Type() != LegacyTxType {
		c += 1
	}
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *Transaction) AsMessage(s Signer, balanceFee *big.Int, number *big.Int, baseFee *big.Int) (Message, error) {
	msg := Message{
		nonce:           tx.Nonce(),
		gasLimit:        tx.Gas(),
		gasPrice:        tx.GasPrice(),
		gasFeeCap:       tx.GasFeeCap(),
		gasTipCap:       tx.GasTipCap(),
		to:              tx.To(),
		amount:          tx.Value(),
		data:            tx.Data(),
		checkNonce:      tx.CheckNonce(),
		balanceTokenFee: balanceFee,
	}
	var err error
	msg.from, err = Sender(s, tx)
	if balanceFee != nil {
		if number.Cmp(common.TIPTRC21FeeBlock) > 0 {
			msg.gasPrice = common.TRC21GasPrice
		} else {
			msg.gasPrice = common.TRC21GasPriceBefore
		}
	} else if baseFee != nil {
		// If baseFee provided, set gasPrice to effectiveGasPrice.
		msg.gasPrice = cmath.BigMin(msg.gasPrice.Add(msg.gasTipCap, baseFee), msg.gasFeeCap)
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
	return &Transaction{inner: cpy}, nil
}

// Cost returns amount + gasprice * gaslimit.
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
	total.Add(total, tx.Value())
	return total
}

// TRC21Cost returns amount + common.TRC21GasPrice * gaslimit.
func (tx *Transaction) TRC21Cost() *big.Int {
	total := new(big.Int).Mul(common.TRC21GasPrice, new(big.Int).SetUint64(tx.Gas()))
	total.Add(total, tx.Value())
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
// The return values may be nil or zero, if the transaction is unsigned.
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawSignatureValues()
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

	addr := common.TomoXListingSMC
	if tx.To().String() != addr.String() {
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

	addr := common.TRC21IssuerSMC
	if tx.To().String() != addr.String() {
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
	var (
		from, to string
		v, r, s  = tx.RawSignatureValues()
	)
	if v != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(v)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.To() == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.To()[:])
	}
	enc, _ := rlp.EncodeToBytes(&tx.inner)
	return fmt.Sprintf(`
	TX(%x)
	Contract:  %v
	From:      %s
	To:        %s
	Nonce:     %v
	GasPrice:  %#x
	GasLimit:  %#x
	GasFeeCap: %#x
	GasTipCap: %#x
	Value:     %#x
	Data:      0x%x
	V:         %#x
	R:         %#x
	S:         %#x
	Hex:       %x
`,
		tx.Hash(),
		tx.To() == nil,
		from,
		to,
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.GasFeeCap(),
		tx.GasTipCap(),
		tx.Value(),
		tx.Data(),
		v, r, s,
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

// txWithMinerFee wraps a transaction with its gas price or effective miner gasTipCap
type txWithMinerFee struct {
	tx   *Transaction
	from common.Address
	fees *big.Int
}

// newTxWithMinerFee creates a wrapped transaction, calculating the effective
// miner gasTipCap if a base fee is provided.
// Returns error in case of a negative effective miner gasTipCap.
func newTxWithMinerFee(tx *Transaction, from common.Address, baseFee *big.Int) (*txWithMinerFee, error) {
	tip := new(big.Int).Set(tx.GasTipCap())
	if baseFee != nil {
		if tx.GasFeeCap().Cmp(baseFee) < 0 {
			return nil, ErrGasFeeCapTooLow
		}
		tip = new(big.Int).Sub(tx.GasFeeCap(), baseFee)
		if tip.Cmp(tx.GasTipCap()) > 0 {
			tip = tx.GasTipCap()
		}
	}
	return &txWithMinerFee{
		tx:   tx,
		from: from,
		fees: tip,
	}, nil
}

// TxByPrice implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByPrice struct {
	txs        []*txWithMinerFee
	payersSwap map[common.Address]*big.Int
}

func (s TxByPrice) Len() int { return len(s.txs) }
func (s TxByPrice) Less(i, j int) bool {
	// Price sorting now based on the effective miner gasTipCap which miner will actually receive, not the gas price
	iPrice := s.txs[i].fees
	if s.txs[i].tx.To() != nil {
		if _, ok := s.payersSwap[*s.txs[i].tx.To()]; ok {
			iPrice = common.TRC21GasPrice
		}
	}

	jPrice := s.txs[j].fees
	if s.txs[j].tx.To() != nil {
		if _, ok := s.payersSwap[*s.txs[j].tx.To()]; ok {
			jPrice = common.TRC21GasPrice
		}
	}
	return iPrice.Cmp(jPrice) > 0
}
func (s TxByPrice) Swap(i, j int) { s.txs[i], s.txs[j] = s.txs[j], s.txs[i] }

func (s *TxByPrice) Push(x interface{}) {
	s.txs = append(s.txs, x.(*txWithMinerFee))
}

func (s *TxByPrice) Pop() interface{} {
	old := s.txs
	n := len(old)
	x := old[n-1]
	// Heap Pop() can cause a memory leak if the element that is returned is a pointer and is not set to nil.
	// This causes the pointer to be held by the underlying array, such that it cannot be garbage collected.
	old[n-1] = nil
	s.txs = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce struct {
	txs     map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads   TxByPrice                       // Next transaction for each unique account (price heap)
	signer  Signer                          // Signer for the set of transactions
	baseFee *big.Int                        // Current base fee
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.

// It also classifies special txs and normal txs
func NewTransactionsByPriceAndNonce(signer Signer, txs map[common.Address]Transactions, signers map[common.Address]struct{},
	payersSwap map[common.Address]*big.Int, baseFee *big.Int) (*TransactionsByPriceAndNonce, Transactions) {
	// Initialize a price based heap with the head transactions
	heads := TxByPrice{}
	heads.payersSwap = payersSwap
	specialTxs := Transactions{}
	for from, accTxs := range txs {
		// Special transactions are always included in the block, whether they satisfy the base fee requirement or not
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
			wrapped, err := newTxWithMinerFee(normalTxs[0], from, baseFee)
			if err != nil {
				// Remove all normal transactions left from this address.
				// Because the transactions have been sorted by nonce already.
				delete(txs, from)
				continue
			}
			heads.txs = append(heads.txs, wrapped)
			txs[from] = normalTxs[1:]
		}
	}
	heap.Init(&heads)
	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce{
		txs:     txs,
		heads:   heads,
		signer:  signer,
		baseFee: baseFee,
	}, specialTxs
}

// Peek returns the next transaction by price.
func (t *TransactionsByPriceAndNonce) Peek() *Transaction {
	if len(t.heads.txs) == 0 {
		return nil
	}
	return t.heads.txs[0].tx
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce) Shift() {
	acc := t.heads.txs[0].from
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		if wrapped, err := newTxWithMinerFee(txs[0], acc, t.baseFee); err == nil {
			t.heads.txs[0], t.txs[acc] = wrapped, txs[1:]
			heap.Fix(&t.heads, 0)
			return
		}
	}
	heap.Pop(&t.heads)
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
	gasFeeCap       *big.Int
	gasTipCap       *big.Int
	data            []byte
	checkNonce      bool
	balanceTokenFee *big.Int
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int,
	gasFeeCap *big.Int, gasTipCap *big.Int, data []byte, checkNonce bool, balanceTokenFee *big.Int) Message {
	if balanceTokenFee != nil {
		gasPrice = common.TRC21GasPrice
	}
	return Message{
		from:            from,
		to:              to,
		nonce:           nonce,
		amount:          amount,
		gasLimit:        gasLimit,
		gasFeeCap:       gasFeeCap,
		gasTipCap:       gasTipCap,
		gasPrice:        gasPrice,
		data:            data,
		checkNonce:      checkNonce,
		balanceTokenFee: balanceTokenFee,
	}
}

func (m Message) From() common.Address      { return m.from }
func (m Message) BalanceTokenFee() *big.Int { return m.balanceTokenFee }
func (m Message) To() *common.Address       { return m.to }
func (m Message) GasPrice() *big.Int        { return m.gasPrice }
func (m Message) GasFeeCap() *big.Int       { return m.gasFeeCap }
func (m Message) GasTipCap() *big.Int       { return m.gasTipCap }
func (m Message) Value() *big.Int           { return m.amount }
func (m Message) Gas() uint64               { return m.gasLimit }
func (m Message) Nonce() uint64             { return m.nonce }
func (m Message) Data() []byte              { return m.data }
func (m Message) CheckNonce() bool          { return m.checkNonce }

// copyAddressPtr copies an address.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}
