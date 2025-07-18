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

package core

import (
	"errors"
	"math"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
)

var (
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==

	4a) Attempt to run transaction data
	4b) If valid, use result as code for the new state object

== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *GasPool
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	//FromFrontier() (common.Address, error)
	To() *common.Address

	GasPrice() *big.Int
	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte
	BalanceTokenFee() *big.Int
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation, homestead bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation && homestead {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGas < nz {
			return 0, vm.ErrOutOfGas
		}
		gas += nz * params.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, vm.ErrOutOfGas
		}
		gas += z * params.TxDataZeroGas
	}
	return gas, nil
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool, owner common.Address) ([]byte, uint64, bool, error) {
	return NewStateTransition(evm, msg, gp).TransitionDb(owner)
}

func (st *StateTransition) from() vm.AccountRef {
	f := st.msg.From()
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return vm.AccountRef(f)
}

func (st *StateTransition) balanceTokenFee() *big.Int {
	return st.msg.BalanceTokenFee()
}

func (st *StateTransition) to() vm.AccountRef {
	if st.msg == nil {
		return vm.AccountRef{}
	}
	to := st.msg.To()
	if to == nil {
		return vm.AccountRef{} // contract creation
	}

	reference := vm.AccountRef(*to)
	if !st.state.Exist(*to) {
		st.state.CreateAccount(*to)
	}
	return reference
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return vm.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas(useAtlasRule bool) (bool, error) {
	var (
		state           = st.state
		balanceTokenFee = st.balanceTokenFee()
		from            = st.from()
	)
	// Check balance based on hard fork status
	if err := st.checkBalance(balanceTokenFee, state, from, useAtlasRule); err != nil {
		return false, err
	}

	// Update gas tracking
	if err := st.gp.SubGas(st.msg.Gas()); err != nil {
		return false, err
	}
	st.gas += st.msg.Gas()
	st.initialGas = st.msg.Gas()

	// Subtract balance based on hard fork status
	isUsedTokenFee := st.subtractBalance(balanceTokenFee, state, from, useAtlasRule)
	return isUsedTokenFee, nil
}

func (st *StateTransition) checkBalance(balanceTokenFee *big.Int, state vm.StateDB, from vm.AccountRef, useAtlasRule bool) error {
	vrc25val := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), common.TRC21GasPrice)
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)
	if useAtlasRule {
		// After Atlas HF: Check balance if no token fee or insufficient token fee
		if balanceTokenFee == nil || balanceTokenFee.Cmp(vrc25val) <= 0 {
			if state.GetBalance(from.Address()).Cmp(mgval) < 0 {
				return errInsufficientBalanceForGas
			}
		}
	} else {
		// Before Atlas HF: Check balance for regular tx or insufficient token fee
		if balanceTokenFee == nil {
			if state.GetBalance(from.Address()).Cmp(mgval) < 0 {
				return errInsufficientBalanceForGas
			}
		} else if balanceTokenFee.Cmp(mgval) < 0 {
			return errInsufficientBalanceForGas
		}
	}
	return nil
}

func (st *StateTransition) subtractBalance(balanceTokenFee *big.Int, stateDB vm.StateDB, from vm.AccountRef, useAtlasRule bool) bool {
	vrc25val := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), common.TRC21GasPrice)
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)
	if useAtlasRule {
		// After Atlas HF: Subtract balance if no token fee or insufficient token fee
		if balanceTokenFee == nil || balanceTokenFee.Cmp(vrc25val) <= 0 {
			stateDB.SubBalance(from.Address(), mgval)
			return false
		} else {
			st.vrc25PayGas(*st.msg.To(), vrc25val)
			return true
		}
	} else {
		// Before Atlas HF: Subtract balance only for regular transactions
		if balanceTokenFee == nil {
			stateDB.SubBalance(from.Address(), mgval)
		}
	}
	return false
}

func (st *StateTransition) preCheck(useAtlasRule bool) (bool, error) {
	msg := st.msg
	sender := st.from()

	// Make sure this transaction's nonce is correct
	if msg.CheckNonce() {
		nonce := st.state.GetNonce(sender.Address())
		if nonce < msg.Nonce() {
			return false, ErrNonceTooHigh
		} else if nonce > msg.Nonce() {
			return false, ErrNonceTooLow
		}
	}
	return st.buyGas(useAtlasRule)
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the the used gas. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb(owner common.Address) (ret []byte, usedGas uint64, failed bool, err error) {
	currentBlock := st.evm.Context.BlockNumber
	isAtlas := st.evm.ChainConfig().IsAtlas(currentBlock)
	var isUsedTokenFee bool
	if isUsedTokenFee, err = st.preCheck(isAtlas); err != nil {
		return nil, 0, false, err
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	// Pay intrinsic gas
	gas, err := IntrinsicGas(st.data, contractCreation, homestead)
	if err != nil {
		return nil, 0, false, err
	}
	if err = st.useGas(gas); err != nil {
		return nil, 0, false, err
	}

	var (
		evm = st.evm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	// for debugging purpose
	// TODO: clean it after fixing the issue https://github.com/tomochain/tomochain/issues/401
	var contractAction string
	nonce := uint64(1)
	if contractCreation {
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
		contractAction = "contract creation"
	} else {
		// Increment the nonce for the next transaction
		nonce = st.state.GetNonce(sender.Address()) + 1
		st.state.SetNonce(sender.Address(), nonce)
		ret, st.gas, vmerr = evm.Call(sender, st.to().Address(), st.data, st.gas, st.value)
		contractAction = "contract call"
	}
	if vmerr != nil {
		log.Debug("VM returned with error", "action", contractAction, "contract address", st.to().Address(), "gas", st.gas, "gasPrice", st.gasPrice, "nonce", nonce, "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == vm.ErrInsufficientBalance {
			return nil, 0, false, vmerr
		}
	}
	st.refundGas(isUsedTokenFee)

	transactionFee := new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice)

	if isAtlas {
		// need to handle it because after Atlas HF, we're not override gas price
		if isUsedTokenFee {
			transactionFee = new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), common.TRC21GasPrice)
		}
	}
	if st.evm.BlockNumber.Cmp(common.TIPTRC21FeeBlock) > 0 {
		if (owner != common.Address{}) {
			st.state.AddBalance(owner, transactionFee)
		}
	} else {
		st.state.AddBalance(st.evm.Coinbase, transactionFee)
	}

	return ret, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) refundGas(isUsedTokenFee bool) {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	balanceTokenFee := st.balanceTokenFee()
	if st.evm.ChainConfig().IsAtlas(st.evm.BlockNumber) {
		if isUsedTokenFee {
			remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), common.TRC21GasPrice)
			st.vrc25RefundGas(*st.msg.To(), remaining)
		}
	} else {
		if balanceTokenFee == nil {
			from := st.from()
			// Return ETH for remaining gas, exchanged at the original rate.
			remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
			st.state.AddBalance(from.Address(), remaining)
		}
	}

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

func (st *StateTransition) vrc25PayGas(token common.Address, usedFee *big.Int) {
	if usedFee.Cmp(big.NewInt(0)) == 0 {
		return
	}
	slotTokensState := state.SlotTRC21Issuer["tokensState"]
	balanceKey := state.GetLocMappingAtKey(token.Hash(), slotTokensState)
	balanceHash := st.state.GetState(common.TRC21IssuerSMC, common.BigToHash(balanceKey))
	currentBalanceInt := new(big.Int).SetBytes(balanceHash[:])

	// Subtract used amount from current balance
	newBalance := new(big.Int).Sub(currentBalanceInt, usedFee)

	st.state.SetState(common.TRC21IssuerSMC, common.BigToHash(balanceKey), common.BigToHash(newBalance))
	st.state.SubBalance(common.TRC21IssuerSMC, usedFee)
}

func (st *StateTransition) vrc25RefundGas(token common.Address, remaining *big.Int) {
	if remaining.Cmp(big.NewInt(0)) == 0 {
		return
	}
	slotTokensState := state.SlotTRC21Issuer["tokensState"]
	balanceKey := state.GetLocMappingAtKey(token.Hash(), slotTokensState)
	balanceHash := st.state.GetState(common.TRC21IssuerSMC, common.BigToHash(balanceKey))
	currentBalanceInt := new(big.Int).SetBytes(balanceHash[:])

	newBalance := new(big.Int).Add(currentBalanceInt, remaining)
	st.state.SetState(common.TRC21IssuerSMC, common.BigToHash(balanceKey), common.BigToHash(newBalance))
	st.state.AddBalance(common.TRC21IssuerSMC, remaining)
}
