// Copyright 2016 The go-ethereum Authors
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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/params"
)

var (
	ErrInvalidChainId    = errors.New("invalid chain id for signer")
	errMissingPayerField = errors.New("transaction has no payer field")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) Signer {
	var signer Signer
	switch {
	case config.IsMiko(blockNumber):
		signer = NewMikoSigner(config.ChainId)
	case config.IsEIP155(blockNumber):
		signer = NewEIP155Signer(config.ChainId)
	case config.IsHomestead(blockNumber):
		signer = HomesteadSigner{}
	default:
		signer = FrontierSigner{}
	}
	return signer
}

// LatestSigner returns the 'most permissive' Signer available for the given chain
// configuration. Specifically, this enables support of EIP-155 replay protection and
// EIP-2930 access list transactions when their respective forks are scheduled to occur at
// any block number in the chain config.
//
// Use this in transaction-handling code where the current block number is unknown. If you
// have the current block number available, use MakeSigner instead.
func LatestSigner(config *params.ChainConfig) Signer {
	if config.ChainId != nil {
		if config.MikoBlock != nil {
			return NewMikoSigner(config.ChainId)
		}
		if config.EIP155Block != nil {
			return NewEIP155Signer(config.ChainId)
		}
	}
	return HomesteadSigner{}
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// SignNewTx creates a transaction and signs it.
func SignNewTx(prv *ecdsa.PrivateKey, s Signer, txdata TxData) (*Transaction, error) {
	tx := NewTx(txdata)
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
	// ChainID
	ChainID() *big.Int

	// Payer returns the payer address of sponsored transaction
	Payer(tx *Transaction) (common.Address, error)
}

// EIP155Transaction implements Signer using the EIP155 rules.
type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func (s EIP155Signer) ChainID() *big.Int {
	return s.chainId
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s EIP155Signer) Equal(s2 Signer) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	if !tx.Protected() {
		return HomesteadSigner{}.Sender(tx)
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlain(s.Hash(tx), R, S, V, true)
}

// WithSignature returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = HomesteadSigner{}.SignatureValues(tx, sig)
	if err != nil {
		return nil, nil, nil, err
	}
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		s.chainId, uint(0), uint(0),
	})
}

func (s EIP155Signer) Payer(tx *Transaction) (common.Address, error) {
	return common.Address{}, ErrInvalidTxType
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) Equal(s2 Signer) bool {
	_, ok := s2.(HomesteadSigner)
	return ok
}

func (hs HomesteadSigner) ChainID() *big.Int {
	return nil
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()

	return recoverPlain(hs.Hash(tx), r, s, v, true)
}

type FrontierSigner struct{}

func (fs FrontierSigner) ChainID() *big.Int {
	return nil
}

func (s FrontierSigner) Equal(s2 Signer) bool {
	_, ok := s2.(FrontierSigner)
	return ok
}

func (fs FrontierSigner) Payer(tx *Transaction) (common.Address, error) {
	return common.Address{}, ErrInvalidTxType
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
	})
}

func (fs FrontierSigner) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlain(fs.Hash(tx), r, s, v, false)
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}

func CacheSigner(signer Signer, tx *Transaction) {
	if tx == nil {
		return
	}
	addr, err := signer.Sender(tx)
	if err != nil {
		return
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
}

type MikoSigner struct {
	EIP155Signer
}

func NewMikoSigner(chainId *big.Int) MikoSigner {
	return MikoSigner{NewEIP155Signer(chainId)}
}

func (s MikoSigner) Equal(s2 Signer) bool {
	miko, ok := s2.(MikoSigner)
	return ok && miko.chainId.Cmp(s.chainId) == 0
}

func (s MikoSigner) Sender(tx *Transaction) (common.Address, error) {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.Sender(tx)
	case SponsoredTxType:
		if tx.ChainId().Cmp(s.chainId) != 0 {
			return common.Address{}, ErrInvalidChainId
		}
		// V in sponsored signature is {0, 1}, but the recoverPlain expects
		// {0, 1} + 27, so we need to add 27 to V
		V, R, S := tx.RawSignatureValues()
		V = new(big.Int).Add(V, big.NewInt(27))
		return recoverPlain(s.Hash(tx), R, S, V, true)
	default:
		return common.Address{}, ErrTxTypeNotSupported
	}
}

func (s MikoSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.SignatureValues(tx, sig)
	case SponsoredTxType:
		// V in sponsored signature is {0, 1}, get it directly from raw signature
		// because decodeSignature returns {0, 1} + 27
		R, S, _ := decodeSignature(sig)
		V := big.NewInt(int64(sig[64]))
		return R, S, V, nil
	default:
		return nil, nil, nil, ErrTxTypeNotSupported
	}
}

func (s MikoSigner) Hash(tx *Transaction) common.Hash {
	switch tx.Type() {
	case LegacyTxType:
		return s.EIP155Signer.Hash(tx)
	case SponsoredTxType:
		payerV, payerR, payerS := tx.RawPayerSignatureValues()
		return prefixedRlpHash(
			tx.Type(),
			[]interface{}{
				s.chainId,
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.ExpiredTime(),
				payerV, payerR, payerS,
			},
		)
	default:
		return common.Hash{}
	}
}

func payerInternal(s Signer, tx *Transaction) (common.Address, error) {
	if tx.Type() != SponsoredTxType {
		return common.Address{}, ErrInvalidTxType
	}

	sender, err := Sender(s, tx)
	if err != nil {
		return common.Address{}, err
	}

	payerV, payerR, payerS := tx.RawPayerSignatureValues()
	payerHash := rlpHash([]interface{}{
		tx.ChainId(), // The chainId is checked in Sender already
		sender,
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		tx.ExpiredTime(),
	})

	// V in payer signature is {0, 1}, but the recoverPlain expects
	// {0, 1} + 27, so we need to add 27 to V
	payerV = new(big.Int).Add(payerV, big.NewInt(27))
	return recoverPlain(payerHash, payerR, payerS, payerV, true)
}

func (s MikoSigner) Payer(tx *Transaction) (common.Address, error) {
	return payerInternal(s, tx)
}

// Payer returns the address derived from payer's signature in sponsored
// transaction or nil in other transaction types.
//
// Payer may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Payer(signer Signer, tx *Transaction) (common.Address, error) {
	if tx.Type() != SponsoredTxType {
		return common.Address{}, errMissingPayerField
	}

	if sc := tx.payer.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Payer(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.payer.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// PayerHash returns the hash to be signed by the payer
func (s MikoSigner) PayerHash(tx *SponsoredTx) common.Hash {
	return rlpHash([]interface{}{
		tx.ChainID,
		tx.Nonce,
		tx.GasPrice,
		tx.Gas,
		tx.To,
		tx.Value,
		tx.Data,
		tx.ExpiredTime,
	})
}

func PayerSign(prv *ecdsa.PrivateKey, signer Signer, sender common.Address, txdata TxData) (r, s, v *big.Int, err error) {
	payerHash := rlpHash([]interface{}{
		signer.ChainID(),
		txdata.nonce(),
		txdata.gasPrice(),
		txdata.gas(),
		txdata.to(),
		txdata.value(),
		txdata.data(),
		txdata.expiredTime(),
	})

	sig, err := crypto.Sign(payerHash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	r, s, _ = decodeSignature(sig)
	v = big.NewInt(int64(sig[64] + 27))
	return r, s, v, nil
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}
