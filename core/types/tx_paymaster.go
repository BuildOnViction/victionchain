package types

import (
	"bytes"
	"math/big"

	"github.com/tomochain/tomochain/rlp"

	"github.com/tomochain/tomochain/common"
)

// PaymasterTx indicates the transactions which can be sponsored gas by a custom paymaster contract.
type PaymasterTx struct {
	ChainID   *big.Int        // destination chain ID
	Nonce     uint64          // nonce of sender account
	GasPrice  *big.Int        // wei per gas
	Gas       uint64          // gas limit
	To        *common.Address `rlp:"nil"` // nil means contract creation
	Value     *big.Int        // wei amount
	Data      []byte          // contract invocation input inner
	PmPayload []byte          // Payload for calling paymaster contracts = PM contract address (required) + custom payload (if any)
	V, R, S   *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *PaymasterTx) copy() TxData {
	cpy := &PaymasterTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are initialized below.
		Value:    new(big.Int),
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *PaymasterTx) txType() byte        { return PaymasterTxType }
func (tx *PaymasterTx) chainID() *big.Int   { return tx.ChainID }
func (tx *PaymasterTx) data() []byte        { return tx.Data }
func (tx *PaymasterTx) gas() uint64         { return tx.Gas }
func (tx *PaymasterTx) gasPrice() *big.Int  { return tx.GasPrice }
func (tx *PaymasterTx) value() *big.Int     { return tx.Value }
func (tx *PaymasterTx) nonce() uint64       { return tx.Nonce }
func (tx *PaymasterTx) to() *common.Address { return tx.To }
func (tx *PaymasterTx) pmPayload() []byte   { return tx.PmPayload }

func (tx *PaymasterTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *PaymasterTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.V, tx.R, tx.S = v, r, s
}

func (tx *PaymasterTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *PaymasterTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
