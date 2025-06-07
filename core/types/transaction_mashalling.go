package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
)

// txJSON is the JSON representation of transactions.
type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID  *hexutil.Big    `json:"chainId,omitempty"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`

	Value   *hexutil.Big    `json:"value"`
	Input   *hexutil.Bytes  `json:"input"`
	V       *hexutil.Big    `json:"v"`
	R       *hexutil.Big    `json:"r"`
	S       *hexutil.Big    `json:"s"`
	YParity *hexutil.Uint64 `json:"yParity,omitempty"`

	// Access list transaction fields:
	AccessList *AccessList `json:"accessList,omitempty"`

	// Sponsored transaction fields
	ExpiredTime *hexutil.Uint64 `json:"expiredTime,omitempty"`
	PayerV      *hexutil.Big    `json:"payerV,omitempty"`
	PayerR      *hexutil.Big    `json:"payerR,omitempty"`
	PayerS      *hexutil.Big    `json:"payerS,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	var enc txJSON
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())

	switch itx := tx.inner.(type) {
	case *LegacyTx:
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		if tx.Protected() {
			enc.ChainID = (*hexutil.Big)(tx.ChainId())
		}
	case *SponsoredTx:
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.GasPrice = (*hexutil.Big)(itx.GasPrice)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		enc.PayerV = (*hexutil.Big)(itx.PayerV)
		enc.PayerR = (*hexutil.Big)(itx.PayerR)
		enc.PayerS = (*hexutil.Big)(itx.PayerS)
		enc.ExpiredTime = (*hexutil.Uint64)(&itx.ExpiredTime)
	}
	return json.Marshal(&enc)
}

func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txJSON
	err := json.Unmarshal(input, &dec)
	if err != nil {
		return err
	}

	// Decode / verify fields according to transaction type.
	var inner TxData
	switch dec.Type {
	case LegacyTxType:
		var itx LegacyTx
		inner = &itx
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input

		// R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// V
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return err
			}
		}
	case SponsoredTxType:
		var itx SponsoredTx
		inner = &itx
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input

		// R
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// S
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// V
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return err
			}
		}
		// Sponsored transaction fields
		itx.ExpiredTime = uint64(*dec.ExpiredTime)
		if dec.PayerV == nil {
			return errors.New("missing required field 'payerV' in transaction")
		}
		itx.PayerV = (*big.Int)(dec.PayerV)
		if dec.PayerR == nil {
			return errors.New("missing required field 'payerR' in transaction")
		}
		itx.PayerR = (*big.Int)(dec.PayerR)
		if dec.PayerS == nil {
			return errors.New("missing required field 'payerS' in transaction")
		}
		itx.PayerS = (*big.Int)(dec.PayerS)
		if err := sanityCheckSignature(itx.PayerV, itx.PayerR, itx.PayerS, false); err != nil {
			return err
		}

	default:
		return ErrTxTypeNotSupported
	}
	tx.setDecoded(inner, 0)

	return nil
}
