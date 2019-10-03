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
	"errors"
	"io"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

var (
	// ErrInvalidOrderSig invalidate signer
	ErrInvalidOrderSig = errors.New("invalid transaction v, r, s values")
	errNoSignerOrder   = errors.New("missing signing methods")
)

// OrderTransaction order transaction
type OrderTransaction struct {
	data ordertxdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type ordertxdata struct {
	AccountNonce    uint64         `json:"nonce"    gencodec:"required"`
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

// EncodeRLP implements rlp.Encoder
func (tx *OrderTransaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

// DecodeRLP implements rlp.Decoder
func (tx *OrderTransaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

// Nonce return nonce of account
func (tx *OrderTransaction) Nonce() uint64                   { return tx.data.AccountNonce }
func (tx *OrderTransaction) Quantity() *big.Int              { return tx.data.Quantity }
func (tx *OrderTransaction) Price() *big.Int                 { return tx.data.Price }
func (tx *OrderTransaction) ExchangeAddress() common.Address { return tx.data.ExchangeAddress }
func (tx *OrderTransaction) UserAddress() common.Address     { return tx.data.UserAddress }
func (tx *OrderTransaction) BaseToken() common.Address       { return tx.data.BaseToken }
func (tx *OrderTransaction) QuoteToken() common.Address      { return tx.data.QuoteToken }
func (tx *OrderTransaction) Status() string                  { return tx.data.Status }
func (tx *OrderTransaction) Side() string                    { return tx.data.Side }
func (tx *OrderTransaction) Type() string                    { return tx.data.Type }
func (tx *OrderTransaction) PairName() string                { return tx.data.PairName }
func (tx *OrderTransaction) Signature() (V, R, S *big.Int)   { return tx.data.V, tx.data.R, tx.data.S }

// From get transaction from
func (tx *OrderTransaction) From() *common.Address {
	if tx.data.V != nil {
		signer := OrderTxSigner{}
		if f, err := OrderSender(signer, tx); err != nil {
			return nil
		} else {
			return &f
		}
	} else {
		return nil
	}
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *OrderTransaction) WithSignature(signer OrderSigner, sig []byte) (*OrderTransaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &OrderTransaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *OrderTransaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

// CacheHash cache hash
func (tx *OrderTransaction) CacheHash() {
	v := rlpHash(tx)
	tx.hash.Store(v)
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *OrderTransaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// NewOrderTransaction init order from value
func NewOrderTransaction(nonce uint64, quantity, price *big.Int, ex, ua, b, q common.Address, status, side, t, pair string) *OrderTransaction {
	return newOrderTransaction(nonce, quantity, price, ex, ua, b, q, status, side, t, pair)
}

func newOrderTransaction(nonce uint64, quantity, price *big.Int, ex, ua, b, q common.Address, status, side, t, pair string) *OrderTransaction {
	d := ordertxdata{
		AccountNonce:    nonce,
		Quantity:        new(big.Int),
		Price:           new(big.Int),
		ExchangeAddress: ex,
		UserAddress:     ua,
		BaseToken:       b,
		QuoteToken:      q,
		Status:          status,
		Side:            side,
		Type:            t,
		PairName:        pair,
		V:               new(big.Int),
		R:               new(big.Int),
		S:               new(big.Int),
	}
	if quantity != nil {
		d.Quantity.Set(quantity)
	}
	if price != nil {
		d.Price.Set(price)
	}

	return &OrderTransaction{data: d}
}

// OrderTransactions is a Transaction slice type for basic sorting.
type OrderTransactions []*OrderTransaction

// Len returns the length of s.
func (s OrderTransactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s OrderTransactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s OrderTransactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// OrderTxDifference returns a new set t which is the difference between a to b.
func OrderTxDifference(a, b OrderTransactions) (keep OrderTransactions) {
	keep = make(OrderTransactions, 0, len(a))

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

// OrderTxByNonce sorted order by nonce defined
type OrderTxByNonce OrderTransactions

func (s OrderTxByNonce) Len() int           { return len(s) }
func (s OrderTxByNonce) Less(i, j int) bool { return s[i].data.AccountNonce < s[j].data.AccountNonce }
func (s OrderTxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
