package core

import (
	"fmt"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/params"
)

// ValidationOptions define certain differences between transaction validation
// across the different pools without having to duplicate those checks.
type ValidationOptions struct {
	Config *params.ChainConfig // Chain configuration to selectively validate based on current fork rules
	Accept uint8               // Bitmap of transaction types that should be accepted for the calling pool
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx *types.Transaction, local bool, opts *ValidationOptions) error {
	// Ensure transactions not implemented by the calling pool are rejected
	if opts.Accept&(1<<tx.Type()) == 0 {
		return fmt.Errorf("%w: tx type %v not supported by this pool", ErrTxTypeNotSupported, tx.Type())
	}
	// check if sender is in black list
	if tx.From() != nil && common.Blacklist[*tx.From()] {
		return fmt.Errorf("Reject transaction with sender in black-list: %v", tx.From().Hex())
	}
	// check if receiver is in black list
	if tx.To() != nil && common.Blacklist[*tx.To()] {
		return fmt.Errorf("Reject transaction with receiver in black-list: %v", tx.To().Hex())
	}

	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	if tx.Size() > 32*1024 {
		return ErrOversizedData
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return ErrNegativeValue
	}
	// Ensure the transaction doesn't exceed the current block limit gas.
	if pool.currentMaxGas < tx.Gas() {
		return ErrGasLimit
	}
	// Make sure the transaction is signed properly
	from, err := types.Sender(pool.signer, tx)
	if err != nil {
		return ErrInvalidSender
	}
	// Drop non-local transactions under our own minimal accepted gas price
	local = local || pool.locals.contains(from) // account may be local even if the transaction arrived from the network
	if !local && pool.gasPrice.Cmp(tx.GasPrice()) > 0 {
		if !tx.IsSpecialTransaction() || (pool.IsSigner != nil && !pool.IsSigner(from)) {
			return ErrUnderpriced
		}
	}
	// Ensure the transaction adheres to nonce ordering
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}
	if pool.pendingState.GetNonce(from)+common.LimitThresholdNonceInQueue < tx.Nonce() {
		return ErrNonceTooHigh
	}
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	balance := pool.currentState.GetBalance(from)
	cost := tx.Cost()
	minGasPrice := common.MinGasPrice
	feeCapacity := big.NewInt(0)

	if tx.To() != nil {
		if value, ok := pool.trc21FeeCapacity[*tx.To()]; ok {
			feeCapacity = value
			if !state.ValidateTRC21Tx(pool.pendingState.StateDB, from, *tx.To(), tx.Data()) {
				return ErrInsufficientFunds
			}
			cost = tx.TRC21Cost()
			minGasPrice = common.TRC21GasPrice
		}
	}
	if new(big.Int).Add(balance, feeCapacity).Cmp(cost) < 0 {
		return ErrInsufficientFunds
	}

	if tx.To() == nil || (tx.To() != nil && !tx.IsSpecialTransaction()) {
		intrGas, err := IntrinsicGas(tx.Data(), tx.To() == nil, pool.homestead)
		if err != nil {
			return err
		}
		// Exclude check smart contract sign address.
		if tx.Gas() < intrGas {
			return ErrIntrinsicGas
		}

		// Check zero gas price.
		if tx.GasPrice().Cmp(new(big.Int).SetInt64(0)) == 0 {
			return ErrZeroGasPrice
		}

		// under min gas price
		if tx.GasPrice().Cmp(minGasPrice) < 0 {
			return ErrUnderMinGasPrice
		}
	}

	/*
		minGasDeploySMC := new(big.Int).Mul(new(big.Int).SetUint64(10), new(big.Int).SetUint64(params.Ether))
		if tx.To() == nil && (tx.Cost().Cmp(minGasDeploySMC) < 0 || tx.GasPrice().Cmp(new(big.Int).SetUint64(10000*params.Shannon)) < 0) {
			return ErrMinDeploySMC
		}
	*/

	// validate minFee slot for TomoZ
	if tx.IsTomoZApplyTransaction() {
		copyState := pool.currentState.Copy()
		return ValidateTomoZApplyTransaction(pool.chain, nil, copyState, common.BytesToAddress(tx.Data()[4:]))
	}

	// validate balance slot, token decimal for TomoX
	if tx.IsTomoXApplyTransaction() {
		copyState := pool.currentState.Copy()
		return ValidateTomoXApplyTransaction(pool.chain, nil, copyState, common.BytesToAddress(tx.Data()[4:]))
	}

	// validate the length of paymaster payload
	if tx.Type() == types.PaymasterTxType && len(tx.PmPayload()) < 20 {
		return ErrPmPayloadTooShort
	}
	return nil
}
