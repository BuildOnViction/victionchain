// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/accounts/abi"
	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/event"
)

// ITRC21ABI is the input ABI used to generate the binding from.
const ITRC21ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"estimateFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Fee\",\"type\":\"event\"}]"

// ITRC21Bin is the compiled bytecode used for deploying new contracts.
const ITRC21Bin = `0x`

// DeployITRC21 deploys a new Ethereum contract, binding an instance of ITRC21 to it.
func DeployITRC21(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ITRC21, error) {
	parsed, err := abi.JSON(strings.NewReader(ITRC21ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ITRC21Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ITRC21{ITRC21Caller: ITRC21Caller{contract: contract}, ITRC21Transactor: ITRC21Transactor{contract: contract}, ITRC21Filterer: ITRC21Filterer{contract: contract}}, nil
}

// ITRC21 is an auto generated Go binding around an Ethereum contract.
type ITRC21 struct {
	ITRC21Caller     // Read-only binding to the contract
	ITRC21Transactor // Write-only binding to the contract
	ITRC21Filterer   // Log filterer for contract events
}

// ITRC21Caller is an auto generated read-only Go binding around an Ethereum contract.
type ITRC21Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITRC21Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ITRC21Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITRC21Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ITRC21Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITRC21Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ITRC21Session struct {
	Contract     *ITRC21           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ITRC21CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ITRC21CallerSession struct {
	Contract *ITRC21Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ITRC21TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ITRC21TransactorSession struct {
	Contract     *ITRC21Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ITRC21Raw is an auto generated low-level Go binding around an Ethereum contract.
type ITRC21Raw struct {
	Contract *ITRC21 // Generic contract binding to access the raw methods on
}

// ITRC21CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ITRC21CallerRaw struct {
	Contract *ITRC21Caller // Generic read-only contract binding to access the raw methods on
}

// ITRC21TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ITRC21TransactorRaw struct {
	Contract *ITRC21Transactor // Generic write-only contract binding to access the raw methods on
}

// NewITRC21 creates a new instance of ITRC21, bound to a specific deployed contract.
func NewITRC21(address common.Address, backend bind.ContractBackend) (*ITRC21, error) {
	contract, err := bindITRC21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ITRC21{ITRC21Caller: ITRC21Caller{contract: contract}, ITRC21Transactor: ITRC21Transactor{contract: contract}, ITRC21Filterer: ITRC21Filterer{contract: contract}}, nil
}

// NewITRC21Caller creates a new read-only instance of ITRC21, bound to a specific deployed contract.
func NewITRC21Caller(address common.Address, caller bind.ContractCaller) (*ITRC21Caller, error) {
	contract, err := bindITRC21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ITRC21Caller{contract: contract}, nil
}

// NewITRC21Transactor creates a new write-only instance of ITRC21, bound to a specific deployed contract.
func NewITRC21Transactor(address common.Address, transactor bind.ContractTransactor) (*ITRC21Transactor, error) {
	contract, err := bindITRC21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ITRC21Transactor{contract: contract}, nil
}

// NewITRC21Filterer creates a new log filterer instance of ITRC21, bound to a specific deployed contract.
func NewITRC21Filterer(address common.Address, filterer bind.ContractFilterer) (*ITRC21Filterer, error) {
	contract, err := bindITRC21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ITRC21Filterer{contract: contract}, nil
}

// bindITRC21 binds a generic wrapper to an already deployed contract.
func bindITRC21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ITRC21ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITRC21 *ITRC21Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ITRC21.Contract.ITRC21Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITRC21 *ITRC21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITRC21.Contract.ITRC21Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITRC21 *ITRC21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITRC21.Contract.ITRC21Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITRC21 *ITRC21CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ITRC21.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITRC21 *ITRC21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITRC21.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITRC21 *ITRC21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITRC21.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ITRC21 *ITRC21Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ITRC21 *ITRC21Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ITRC21.Contract.Allowance(&_ITRC21.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_ITRC21 *ITRC21CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ITRC21.Contract.Allowance(&_ITRC21.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ITRC21 *ITRC21Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ITRC21 *ITRC21Session) BalanceOf(who common.Address) (*big.Int, error) {
	return _ITRC21.Contract.BalanceOf(&_ITRC21.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_ITRC21 *ITRC21CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _ITRC21.Contract.BalanceOf(&_ITRC21.CallOpts, who)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ITRC21 *ITRC21Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ITRC21 *ITRC21Session) Decimals() (uint8, error) {
	return _ITRC21.Contract.Decimals(&_ITRC21.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ITRC21 *ITRC21CallerSession) Decimals() (uint8, error) {
	return _ITRC21.Contract.Decimals(&_ITRC21.CallOpts)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_ITRC21 *ITRC21Caller) EstimateFee(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "estimateFee", value)
	return *ret0, err
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_ITRC21 *ITRC21Session) EstimateFee(value *big.Int) (*big.Int, error) {
	return _ITRC21.Contract.EstimateFee(&_ITRC21.CallOpts, value)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_ITRC21 *ITRC21CallerSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _ITRC21.Contract.EstimateFee(&_ITRC21.CallOpts, value)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_ITRC21 *ITRC21Caller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "issuer")
	return *ret0, err
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_ITRC21 *ITRC21Session) Issuer() (common.Address, error) {
	return _ITRC21.Contract.Issuer(&_ITRC21.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_ITRC21 *ITRC21CallerSession) Issuer() (common.Address, error) {
	return _ITRC21.Contract.Issuer(&_ITRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ITRC21 *ITRC21Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ITRC21.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ITRC21 *ITRC21Session) TotalSupply() (*big.Int, error) {
	return _ITRC21.Contract.TotalSupply(&_ITRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ITRC21 *ITRC21CallerSession) TotalSupply() (*big.Int, error) {
	return _ITRC21.Contract.TotalSupply(&_ITRC21.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.Approve(&_ITRC21.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_ITRC21 *ITRC21TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.Approve(&_ITRC21.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.Transfer(&_ITRC21.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.Transfer(&_ITRC21.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.TransferFrom(&_ITRC21.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_ITRC21 *ITRC21TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ITRC21.Contract.TransferFrom(&_ITRC21.TransactOpts, from, to, value)
}

// ITRC21ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ITRC21 contract.
type ITRC21ApprovalIterator struct {
	Event *ITRC21Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITRC21ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITRC21Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ITRC21Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ITRC21ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITRC21ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITRC21Approval represents a Approval event raised by the ITRC21 contract.
type ITRC21Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ITRC21ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ITRC21.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ITRC21ApprovalIterator{contract: _ITRC21.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ITRC21Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ITRC21.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITRC21Approval)
				if err := _ITRC21.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ITRC21FeeIterator is returned from FilterFee and is used to iterate over the raw logs and unpacked data for Fee events raised by the ITRC21 contract.
type ITRC21FeeIterator struct {
	Event *ITRC21Fee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITRC21FeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITRC21Fee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ITRC21Fee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ITRC21FeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITRC21FeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITRC21Fee represents a Fee event raised by the ITRC21 contract.
type ITRC21Fee struct {
	From   common.Address
	To     common.Address
	Issuer common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFee is a free log retrieval operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) FilterFee(opts *bind.FilterOpts, from []common.Address, to []common.Address, issuer []common.Address) (*ITRC21FeeIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _ITRC21.contract.FilterLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return &ITRC21FeeIterator{contract: _ITRC21.contract, event: "Fee", logs: logs, sub: sub}, nil
}

// WatchFee is a free log subscription operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) WatchFee(opts *bind.WatchOpts, sink chan<- *ITRC21Fee, from []common.Address, to []common.Address, issuer []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _ITRC21.contract.WatchLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITRC21Fee)
				if err := _ITRC21.contract.UnpackLog(event, "Fee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ITRC21TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ITRC21 contract.
type ITRC21TransferIterator struct {
	Event *ITRC21Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITRC21TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITRC21Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ITRC21Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ITRC21TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITRC21TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITRC21Transfer represents a Transfer event raised by the ITRC21 contract.
type ITRC21Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ITRC21TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ITRC21.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ITRC21TransferIterator{contract: _ITRC21.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_ITRC21 *ITRC21Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ITRC21Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ITRC21.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITRC21Transfer)
				if err := _ITRC21.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21ABI is the input ABI used to generate the binding from.
const MyTRC21ABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"owners\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"estimateFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"revokeConfirmation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setMinFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"confirmations\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"depositFee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"pending\",\"type\":\"bool\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"name\":\"getTransactionCount\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"isConfirmed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"getConfirmationCount\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"isMintingTx\",\"type\":\"bool\"},{\"name\":\"destination\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"WITHDRAW_FEE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"from\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"uint256\"},{\"name\":\"pending\",\"type\":\"bool\"},{\"name\":\"executed\",\"type\":\"bool\"}],\"name\":\"getTransactionIds\",\"outputs\":[{\"name\":\"_transactionIds\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"getConfirmations\",\"outputs\":[{\"name\":\"_confirmations\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"withdrawFee\",\"type\":\"uint256\"}],\"name\":\"setWithdrawFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"transactionCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_required\",\"type\":\"uint256\"}],\"name\":\"changeRequirement\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"confirmTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferContractIssuer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"destination\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"submitTransaction\",\"outputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_OWNER_COUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"required\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"DEPOSIT_FEE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"replaceOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"executeTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"transactionId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owners\",\"type\":\"address[]\"},{\"name\":\"_required\",\"type\":\"uint256\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_decimals\",\"type\":\"uint8\"},{\"name\":\"cap\",\"type\":\"uint256\"},{\"name\":\"minFee\",\"type\":\"uint256\"},{\"name\":\"depositFee\",\"type\":\"uint256\"},{\"name\":\"withdrawFee\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Confirmation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Revocation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"Submission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"isMintingTx\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Execution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"ExecutionFailure\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerAddition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerRemoval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"required\",\"type\":\"uint256\"}],\"name\":\"RequirementChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Fee\",\"type\":\"event\"}]"

// MyTRC21Bin is the compiled bytecode used for deploying new contracts.
const MyTRC21Bin = `0x`

// DeployMyTRC21 deploys a new Ethereum contract, binding an instance of MyTRC21 to it.
func DeployMyTRC21(auth *bind.TransactOpts, backend bind.ContractBackend, _owners []common.Address, _required *big.Int, _name string, _symbol string, _decimals uint8, cap *big.Int, minFee *big.Int, depositFee *big.Int, withdrawFee *big.Int) (common.Address, *types.Transaction, *MyTRC21, error) {
	parsed, err := abi.JSON(strings.NewReader(MyTRC21ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MyTRC21Bin), backend, _owners, _required, _name, _symbol, _decimals, cap, minFee, depositFee, withdrawFee)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MyTRC21{MyTRC21Caller: MyTRC21Caller{contract: contract}, MyTRC21Transactor: MyTRC21Transactor{contract: contract}, MyTRC21Filterer: MyTRC21Filterer{contract: contract}}, nil
}

// MyTRC21 is an auto generated Go binding around an Ethereum contract.
type MyTRC21 struct {
	MyTRC21Caller     // Read-only binding to the contract
	MyTRC21Transactor // Write-only binding to the contract
	MyTRC21Filterer   // Log filterer for contract events
}

// MyTRC21Caller is an auto generated read-only Go binding around an Ethereum contract.
type MyTRC21Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyTRC21Transactor is an auto generated write-only Go binding around an Ethereum contract.
type MyTRC21Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyTRC21Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MyTRC21Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MyTRC21Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MyTRC21Session struct {
	Contract     *MyTRC21          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MyTRC21CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MyTRC21CallerSession struct {
	Contract *MyTRC21Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// MyTRC21TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MyTRC21TransactorSession struct {
	Contract     *MyTRC21Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MyTRC21Raw is an auto generated low-level Go binding around an Ethereum contract.
type MyTRC21Raw struct {
	Contract *MyTRC21 // Generic contract binding to access the raw methods on
}

// MyTRC21CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MyTRC21CallerRaw struct {
	Contract *MyTRC21Caller // Generic read-only contract binding to access the raw methods on
}

// MyTRC21TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MyTRC21TransactorRaw struct {
	Contract *MyTRC21Transactor // Generic write-only contract binding to access the raw methods on
}

// NewMyTRC21 creates a new instance of MyTRC21, bound to a specific deployed contract.
func NewMyTRC21(address common.Address, backend bind.ContractBackend) (*MyTRC21, error) {
	contract, err := bindMyTRC21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MyTRC21{MyTRC21Caller: MyTRC21Caller{contract: contract}, MyTRC21Transactor: MyTRC21Transactor{contract: contract}, MyTRC21Filterer: MyTRC21Filterer{contract: contract}}, nil
}

// NewMyTRC21Caller creates a new read-only instance of MyTRC21, bound to a specific deployed contract.
func NewMyTRC21Caller(address common.Address, caller bind.ContractCaller) (*MyTRC21Caller, error) {
	contract, err := bindMyTRC21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MyTRC21Caller{contract: contract}, nil
}

// NewMyTRC21Transactor creates a new write-only instance of MyTRC21, bound to a specific deployed contract.
func NewMyTRC21Transactor(address common.Address, transactor bind.ContractTransactor) (*MyTRC21Transactor, error) {
	contract, err := bindMyTRC21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MyTRC21Transactor{contract: contract}, nil
}

// NewMyTRC21Filterer creates a new log filterer instance of MyTRC21, bound to a specific deployed contract.
func NewMyTRC21Filterer(address common.Address, filterer bind.ContractFilterer) (*MyTRC21Filterer, error) {
	contract, err := bindMyTRC21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MyTRC21Filterer{contract: contract}, nil
}

// bindMyTRC21 binds a generic wrapper to an already deployed contract.
func bindMyTRC21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MyTRC21ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MyTRC21 *MyTRC21Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MyTRC21.Contract.MyTRC21Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MyTRC21 *MyTRC21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MyTRC21.Contract.MyTRC21Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MyTRC21 *MyTRC21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MyTRC21.Contract.MyTRC21Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MyTRC21 *MyTRC21CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MyTRC21.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MyTRC21 *MyTRC21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MyTRC21.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MyTRC21 *MyTRC21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MyTRC21.Contract.contract.Transact(opts, method, params...)
}

// DEPOSITFEE is a free data retrieval call binding the contract method 0xde363e65.
//
// Solidity: function DEPOSIT_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) DEPOSITFEE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "DEPOSIT_FEE")
	return *ret0, err
}

// DEPOSITFEE is a free data retrieval call binding the contract method 0xde363e65.
//
// Solidity: function DEPOSIT_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) DEPOSITFEE() (*big.Int, error) {
	return _MyTRC21.Contract.DEPOSITFEE(&_MyTRC21.CallOpts)
}

// DEPOSITFEE is a free data retrieval call binding the contract method 0xde363e65.
//
// Solidity: function DEPOSIT_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) DEPOSITFEE() (*big.Int, error) {
	return _MyTRC21.Contract.DEPOSITFEE(&_MyTRC21.CallOpts)
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) MAXOWNERCOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "MAX_OWNER_COUNT")
	return *ret0, err
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) MAXOWNERCOUNT() (*big.Int, error) {
	return _MyTRC21.Contract.MAXOWNERCOUNT(&_MyTRC21.CallOpts)
}

// MAXOWNERCOUNT is a free data retrieval call binding the contract method 0xd74f8edd.
//
// Solidity: function MAX_OWNER_COUNT() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) MAXOWNERCOUNT() (*big.Int, error) {
	return _MyTRC21.Contract.MAXOWNERCOUNT(&_MyTRC21.CallOpts)
}

// WITHDRAWFEE is a free data retrieval call binding the contract method 0x9bff5ddb.
//
// Solidity: function WITHDRAW_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) WITHDRAWFEE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "WITHDRAW_FEE")
	return *ret0, err
}

// WITHDRAWFEE is a free data retrieval call binding the contract method 0x9bff5ddb.
//
// Solidity: function WITHDRAW_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) WITHDRAWFEE() (*big.Int, error) {
	return _MyTRC21.Contract.WITHDRAWFEE(&_MyTRC21.CallOpts)
}

// WITHDRAWFEE is a free data retrieval call binding the contract method 0x9bff5ddb.
//
// Solidity: function WITHDRAW_FEE() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) WITHDRAWFEE() (*big.Int, error) {
	return _MyTRC21.Contract.WITHDRAWFEE(&_MyTRC21.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MyTRC21.Contract.Allowance(&_MyTRC21.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MyTRC21.Contract.Allowance(&_MyTRC21.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _MyTRC21.Contract.BalanceOf(&_MyTRC21.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _MyTRC21.Contract.BalanceOf(&_MyTRC21.CallOpts, owner)
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_MyTRC21 *MyTRC21Caller) Confirmations(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "confirmations", arg0, arg1)
	return *ret0, err
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_MyTRC21 *MyTRC21Session) Confirmations(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MyTRC21.Contract.Confirmations(&_MyTRC21.CallOpts, arg0, arg1)
}

// Confirmations is a free data retrieval call binding the contract method 0x3411c81c.
//
// Solidity: function confirmations( uint256,  address) constant returns(bool)
func (_MyTRC21 *MyTRC21CallerSession) Confirmations(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _MyTRC21.Contract.Confirmations(&_MyTRC21.CallOpts, arg0, arg1)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_MyTRC21 *MyTRC21Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_MyTRC21 *MyTRC21Session) Decimals() (uint8, error) {
	return _MyTRC21.Contract.Decimals(&_MyTRC21.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_MyTRC21 *MyTRC21CallerSession) Decimals() (uint8, error) {
	return _MyTRC21.Contract.Decimals(&_MyTRC21.CallOpts)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) EstimateFee(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "estimateFee", value)
	return *ret0, err
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) EstimateFee(value *big.Int) (*big.Int, error) {
	return _MyTRC21.Contract.EstimateFee(&_MyTRC21.CallOpts, value)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _MyTRC21.Contract.EstimateFee(&_MyTRC21.CallOpts, value)
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_MyTRC21 *MyTRC21Caller) GetConfirmationCount(opts *bind.CallOpts, transactionId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "getConfirmationCount", transactionId)
	return *ret0, err
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_MyTRC21 *MyTRC21Session) GetConfirmationCount(transactionId *big.Int) (*big.Int, error) {
	return _MyTRC21.Contract.GetConfirmationCount(&_MyTRC21.CallOpts, transactionId)
}

// GetConfirmationCount is a free data retrieval call binding the contract method 0x8b51d13f.
//
// Solidity: function getConfirmationCount(transactionId uint256) constant returns(count uint256)
func (_MyTRC21 *MyTRC21CallerSession) GetConfirmationCount(transactionId *big.Int) (*big.Int, error) {
	return _MyTRC21.Contract.GetConfirmationCount(&_MyTRC21.CallOpts, transactionId)
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_MyTRC21 *MyTRC21Caller) GetConfirmations(opts *bind.CallOpts, transactionId *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "getConfirmations", transactionId)
	return *ret0, err
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_MyTRC21 *MyTRC21Session) GetConfirmations(transactionId *big.Int) ([]common.Address, error) {
	return _MyTRC21.Contract.GetConfirmations(&_MyTRC21.CallOpts, transactionId)
}

// GetConfirmations is a free data retrieval call binding the contract method 0xb5dc40c3.
//
// Solidity: function getConfirmations(transactionId uint256) constant returns(_confirmations address[])
func (_MyTRC21 *MyTRC21CallerSession) GetConfirmations(transactionId *big.Int) ([]common.Address, error) {
	return _MyTRC21.Contract.GetConfirmations(&_MyTRC21.CallOpts, transactionId)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_MyTRC21 *MyTRC21Caller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "getOwners")
	return *ret0, err
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_MyTRC21 *MyTRC21Session) GetOwners() ([]common.Address, error) {
	return _MyTRC21.Contract.GetOwners(&_MyTRC21.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() constant returns(address[])
func (_MyTRC21 *MyTRC21CallerSession) GetOwners() ([]common.Address, error) {
	return _MyTRC21.Contract.GetOwners(&_MyTRC21.CallOpts)
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_MyTRC21 *MyTRC21Caller) GetTransactionCount(opts *bind.CallOpts, pending bool, executed bool) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "getTransactionCount", pending, executed)
	return *ret0, err
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_MyTRC21 *MyTRC21Session) GetTransactionCount(pending bool, executed bool) (*big.Int, error) {
	return _MyTRC21.Contract.GetTransactionCount(&_MyTRC21.CallOpts, pending, executed)
}

// GetTransactionCount is a free data retrieval call binding the contract method 0x54741525.
//
// Solidity: function getTransactionCount(pending bool, executed bool) constant returns(count uint256)
func (_MyTRC21 *MyTRC21CallerSession) GetTransactionCount(pending bool, executed bool) (*big.Int, error) {
	return _MyTRC21.Contract.GetTransactionCount(&_MyTRC21.CallOpts, pending, executed)
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_MyTRC21 *MyTRC21Caller) GetTransactionIds(opts *bind.CallOpts, from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "getTransactionIds", from, to, pending, executed)
	return *ret0, err
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_MyTRC21 *MyTRC21Session) GetTransactionIds(from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	return _MyTRC21.Contract.GetTransactionIds(&_MyTRC21.CallOpts, from, to, pending, executed)
}

// GetTransactionIds is a free data retrieval call binding the contract method 0xa8abe69a.
//
// Solidity: function getTransactionIds(from uint256, to uint256, pending bool, executed bool) constant returns(_transactionIds uint256[])
func (_MyTRC21 *MyTRC21CallerSession) GetTransactionIds(from *big.Int, to *big.Int, pending bool, executed bool) ([]*big.Int, error) {
	return _MyTRC21.Contract.GetTransactionIds(&_MyTRC21.CallOpts, from, to, pending, executed)
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_MyTRC21 *MyTRC21Caller) IsConfirmed(opts *bind.CallOpts, transactionId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "isConfirmed", transactionId)
	return *ret0, err
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_MyTRC21 *MyTRC21Session) IsConfirmed(transactionId *big.Int) (bool, error) {
	return _MyTRC21.Contract.IsConfirmed(&_MyTRC21.CallOpts, transactionId)
}

// IsConfirmed is a free data retrieval call binding the contract method 0x784547a7.
//
// Solidity: function isConfirmed(transactionId uint256) constant returns(bool)
func (_MyTRC21 *MyTRC21CallerSession) IsConfirmed(transactionId *big.Int) (bool, error) {
	return _MyTRC21.Contract.IsConfirmed(&_MyTRC21.CallOpts, transactionId)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_MyTRC21 *MyTRC21Caller) IsOwner(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "isOwner", arg0)
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_MyTRC21 *MyTRC21Session) IsOwner(arg0 common.Address) (bool, error) {
	return _MyTRC21.Contract.IsOwner(&_MyTRC21.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner( address) constant returns(bool)
func (_MyTRC21 *MyTRC21CallerSession) IsOwner(arg0 common.Address) (bool, error) {
	return _MyTRC21.Contract.IsOwner(&_MyTRC21.CallOpts, arg0)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_MyTRC21 *MyTRC21Caller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "issuer")
	return *ret0, err
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_MyTRC21 *MyTRC21Session) Issuer() (common.Address, error) {
	return _MyTRC21.Contract.Issuer(&_MyTRC21.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_MyTRC21 *MyTRC21CallerSession) Issuer() (common.Address, error) {
	return _MyTRC21.Contract.Issuer(&_MyTRC21.CallOpts)
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) MinFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "minFee")
	return *ret0, err
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) MinFee() (*big.Int, error) {
	return _MyTRC21.Contract.MinFee(&_MyTRC21.CallOpts)
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) MinFee() (*big.Int, error) {
	return _MyTRC21.Contract.MinFee(&_MyTRC21.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_MyTRC21 *MyTRC21Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_MyTRC21 *MyTRC21Session) Name() (string, error) {
	return _MyTRC21.Contract.Name(&_MyTRC21.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_MyTRC21 *MyTRC21CallerSession) Name() (string, error) {
	return _MyTRC21.Contract.Name(&_MyTRC21.CallOpts)
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_MyTRC21 *MyTRC21Caller) Owners(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "owners", arg0)
	return *ret0, err
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_MyTRC21 *MyTRC21Session) Owners(arg0 *big.Int) (common.Address, error) {
	return _MyTRC21.Contract.Owners(&_MyTRC21.CallOpts, arg0)
}

// Owners is a free data retrieval call binding the contract method 0x025e7c27.
//
// Solidity: function owners( uint256) constant returns(address)
func (_MyTRC21 *MyTRC21CallerSession) Owners(arg0 *big.Int) (common.Address, error) {
	return _MyTRC21.Contract.Owners(&_MyTRC21.CallOpts, arg0)
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) Required(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "required")
	return *ret0, err
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) Required() (*big.Int, error) {
	return _MyTRC21.Contract.Required(&_MyTRC21.CallOpts)
}

// Required is a free data retrieval call binding the contract method 0xdc8452cd.
//
// Solidity: function required() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) Required() (*big.Int, error) {
	return _MyTRC21.Contract.Required(&_MyTRC21.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_MyTRC21 *MyTRC21Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_MyTRC21 *MyTRC21Session) Symbol() (string, error) {
	return _MyTRC21.Contract.Symbol(&_MyTRC21.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_MyTRC21 *MyTRC21CallerSession) Symbol() (string, error) {
	return _MyTRC21.Contract.Symbol(&_MyTRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) TotalSupply() (*big.Int, error) {
	return _MyTRC21.Contract.TotalSupply(&_MyTRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) TotalSupply() (*big.Int, error) {
	return _MyTRC21.Contract.TotalSupply(&_MyTRC21.CallOpts)
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_MyTRC21 *MyTRC21Caller) TransactionCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MyTRC21.contract.Call(opts, out, "transactionCount")
	return *ret0, err
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_MyTRC21 *MyTRC21Session) TransactionCount() (*big.Int, error) {
	return _MyTRC21.Contract.TransactionCount(&_MyTRC21.CallOpts)
}

// TransactionCount is a free data retrieval call binding the contract method 0xb77bf600.
//
// Solidity: function transactionCount() constant returns(uint256)
func (_MyTRC21 *MyTRC21CallerSession) TransactionCount() (*big.Int, error) {
	return _MyTRC21.Contract.TransactionCount(&_MyTRC21.CallOpts)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(isMintingTx bool, destination address, value uint256, data bytes, executed bool)
func (_MyTRC21 *MyTRC21Caller) Transactions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	IsMintingTx bool
	Destination common.Address
	Value       *big.Int
	Data        []byte
	Executed    bool
}, error) {
	ret := new(struct {
		IsMintingTx bool
		Destination common.Address
		Value       *big.Int
		Data        []byte
		Executed    bool
	})
	out := ret
	err := _MyTRC21.contract.Call(opts, out, "transactions", arg0)
	return *ret, err
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(isMintingTx bool, destination address, value uint256, data bytes, executed bool)
func (_MyTRC21 *MyTRC21Session) Transactions(arg0 *big.Int) (struct {
	IsMintingTx bool
	Destination common.Address
	Value       *big.Int
	Data        []byte
	Executed    bool
}, error) {
	return _MyTRC21.Contract.Transactions(&_MyTRC21.CallOpts, arg0)
}

// Transactions is a free data retrieval call binding the contract method 0x9ace38c2.
//
// Solidity: function transactions( uint256) constant returns(isMintingTx bool, destination address, value uint256, data bytes, executed bool)
func (_MyTRC21 *MyTRC21CallerSession) Transactions(arg0 *big.Int) (struct {
	IsMintingTx bool
	Destination common.Address
	Value       *big.Int
	Data        []byte
	Executed    bool
}, error) {
	return _MyTRC21.Contract.Transactions(&_MyTRC21.CallOpts, arg0)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_MyTRC21 *MyTRC21Transactor) AddOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "addOwner", owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_MyTRC21 *MyTRC21Session) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.AddOwner(&_MyTRC21.TransactOpts, owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(owner address) returns()
func (_MyTRC21 *MyTRC21TransactorSession) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.AddOwner(&_MyTRC21.TransactOpts, owner)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.Approve(&_MyTRC21.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.Approve(&_MyTRC21.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21Transactor) Burn(opts *bind.TransactOpts, value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "burn", value, data)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21Session) Burn(value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.Contract.Burn(&_MyTRC21.TransactOpts, value, data)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21TransactorSession) Burn(value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.Contract.Burn(&_MyTRC21.TransactOpts, value, data)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) ChangeRequirement(opts *bind.TransactOpts, _required *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "changeRequirement", _required)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_MyTRC21 *MyTRC21Session) ChangeRequirement(_required *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ChangeRequirement(&_MyTRC21.TransactOpts, _required)
}

// ChangeRequirement is a paid mutator transaction binding the contract method 0xba51a6df.
//
// Solidity: function changeRequirement(_required uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) ChangeRequirement(_required *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ChangeRequirement(&_MyTRC21.TransactOpts, _required)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) ConfirmTransaction(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "confirmTransaction", transactionId)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Session) ConfirmTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ConfirmTransaction(&_MyTRC21.TransactOpts, transactionId)
}

// ConfirmTransaction is a paid mutator transaction binding the contract method 0xc01a8c84.
//
// Solidity: function confirmTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) ConfirmTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ConfirmTransaction(&_MyTRC21.TransactOpts, transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) ExecuteTransaction(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "executeTransaction", transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Session) ExecuteTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ExecuteTransaction(&_MyTRC21.TransactOpts, transactionId)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xee22610b.
//
// Solidity: function executeTransaction(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) ExecuteTransaction(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.ExecuteTransaction(&_MyTRC21.TransactOpts, transactionId)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_MyTRC21 *MyTRC21Transactor) RemoveOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "removeOwner", owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_MyTRC21 *MyTRC21Session) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.RemoveOwner(&_MyTRC21.TransactOpts, owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(owner address) returns()
func (_MyTRC21 *MyTRC21TransactorSession) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.RemoveOwner(&_MyTRC21.TransactOpts, owner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_MyTRC21 *MyTRC21Transactor) ReplaceOwner(opts *bind.TransactOpts, owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "replaceOwner", owner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_MyTRC21 *MyTRC21Session) ReplaceOwner(owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.ReplaceOwner(&_MyTRC21.TransactOpts, owner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(owner address, newOwner address) returns()
func (_MyTRC21 *MyTRC21TransactorSession) ReplaceOwner(owner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.ReplaceOwner(&_MyTRC21.TransactOpts, owner, newOwner)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) RevokeConfirmation(opts *bind.TransactOpts, transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "revokeConfirmation", transactionId)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21Session) RevokeConfirmation(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.RevokeConfirmation(&_MyTRC21.TransactOpts, transactionId)
}

// RevokeConfirmation is a paid mutator transaction binding the contract method 0x20ea8d86.
//
// Solidity: function revokeConfirmation(transactionId uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) RevokeConfirmation(transactionId *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.RevokeConfirmation(&_MyTRC21.TransactOpts, transactionId)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(depositFee uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) SetDepositFee(opts *bind.TransactOpts, depositFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "setDepositFee", depositFee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(depositFee uint256) returns()
func (_MyTRC21 *MyTRC21Session) SetDepositFee(depositFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetDepositFee(&_MyTRC21.TransactOpts, depositFee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(depositFee uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) SetDepositFee(depositFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetDepositFee(&_MyTRC21.TransactOpts, depositFee)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) SetMinFee(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "setMinFee", value)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_MyTRC21 *MyTRC21Session) SetMinFee(value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetMinFee(&_MyTRC21.TransactOpts, value)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) SetMinFee(value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetMinFee(&_MyTRC21.TransactOpts, value)
}

// SetWithdrawFee is a paid mutator transaction binding the contract method 0xb6ac642a.
//
// Solidity: function setWithdrawFee(withdrawFee uint256) returns()
func (_MyTRC21 *MyTRC21Transactor) SetWithdrawFee(opts *bind.TransactOpts, withdrawFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "setWithdrawFee", withdrawFee)
}

// SetWithdrawFee is a paid mutator transaction binding the contract method 0xb6ac642a.
//
// Solidity: function setWithdrawFee(withdrawFee uint256) returns()
func (_MyTRC21 *MyTRC21Session) SetWithdrawFee(withdrawFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetWithdrawFee(&_MyTRC21.TransactOpts, withdrawFee)
}

// SetWithdrawFee is a paid mutator transaction binding the contract method 0xb6ac642a.
//
// Solidity: function setWithdrawFee(withdrawFee uint256) returns()
func (_MyTRC21 *MyTRC21TransactorSession) SetWithdrawFee(withdrawFee *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.SetWithdrawFee(&_MyTRC21.TransactOpts, withdrawFee)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc6427474.
//
// Solidity: function submitTransaction(destination address, value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21Transactor) SubmitTransaction(opts *bind.TransactOpts, destination common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "submitTransaction", destination, value, data)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc6427474.
//
// Solidity: function submitTransaction(destination address, value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21Session) SubmitTransaction(destination common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.Contract.SubmitTransaction(&_MyTRC21.TransactOpts, destination, value, data)
}

// SubmitTransaction is a paid mutator transaction binding the contract method 0xc6427474.
//
// Solidity: function submitTransaction(destination address, value uint256, data bytes) returns(transactionId uint256)
func (_MyTRC21 *MyTRC21TransactorSession) SubmitTransaction(destination common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _MyTRC21.Contract.SubmitTransaction(&_MyTRC21.TransactOpts, destination, value, data)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.Transfer(&_MyTRC21.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.Transfer(&_MyTRC21.TransactOpts, to, value)
}

// TransferContractIssuer is a paid mutator transaction binding the contract method 0xc28ce6ff.
//
// Solidity: function transferContractIssuer(newOwner address) returns()
func (_MyTRC21 *MyTRC21Transactor) TransferContractIssuer(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "transferContractIssuer", newOwner)
}

// TransferContractIssuer is a paid mutator transaction binding the contract method 0xc28ce6ff.
//
// Solidity: function transferContractIssuer(newOwner address) returns()
func (_MyTRC21 *MyTRC21Session) TransferContractIssuer(newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.TransferContractIssuer(&_MyTRC21.TransactOpts, newOwner)
}

// TransferContractIssuer is a paid mutator transaction binding the contract method 0xc28ce6ff.
//
// Solidity: function transferContractIssuer(newOwner address) returns()
func (_MyTRC21 *MyTRC21TransactorSession) TransferContractIssuer(newOwner common.Address) (*types.Transaction, error) {
	return _MyTRC21.Contract.TransferContractIssuer(&_MyTRC21.TransactOpts, newOwner)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.TransferFrom(&_MyTRC21.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_MyTRC21 *MyTRC21TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MyTRC21.Contract.TransferFrom(&_MyTRC21.TransactOpts, from, to, value)
}

// MyTRC21ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MyTRC21 contract.
type MyTRC21ApprovalIterator struct {
	Event *MyTRC21Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Approval represents a Approval event raised by the MyTRC21 contract.
type MyTRC21Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MyTRC21ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21ApprovalIterator{contract: _MyTRC21.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MyTRC21Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Approval)
				if err := _MyTRC21.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21ConfirmationIterator is returned from FilterConfirmation and is used to iterate over the raw logs and unpacked data for Confirmation events raised by the MyTRC21 contract.
type MyTRC21ConfirmationIterator struct {
	Event *MyTRC21Confirmation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21ConfirmationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Confirmation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Confirmation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21ConfirmationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21ConfirmationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Confirmation represents a Confirmation event raised by the MyTRC21 contract.
type MyTRC21Confirmation struct {
	Sender        common.Address
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfirmation is a free log retrieval operation binding the contract event 0x4a504a94899432a9846e1aa406dceb1bcfd538bb839071d49d1e5e23f5be30ef.
//
// Solidity: event Confirmation(sender indexed address, transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterConfirmation(opts *bind.FilterOpts, sender []common.Address, transactionId []*big.Int) (*MyTRC21ConfirmationIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Confirmation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21ConfirmationIterator{contract: _MyTRC21.contract, event: "Confirmation", logs: logs, sub: sub}, nil
}

// WatchConfirmation is a free log subscription operation binding the contract event 0x4a504a94899432a9846e1aa406dceb1bcfd538bb839071d49d1e5e23f5be30ef.
//
// Solidity: event Confirmation(sender indexed address, transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchConfirmation(opts *bind.WatchOpts, sink chan<- *MyTRC21Confirmation, sender []common.Address, transactionId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Confirmation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Confirmation)
				if err := _MyTRC21.contract.UnpackLog(event, "Confirmation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21ExecutionIterator is returned from FilterExecution and is used to iterate over the raw logs and unpacked data for Execution events raised by the MyTRC21 contract.
type MyTRC21ExecutionIterator struct {
	Event *MyTRC21Execution // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21ExecutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Execution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Execution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21ExecutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21ExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Execution represents a Execution event raised by the MyTRC21 contract.
type MyTRC21Execution struct {
	TransactionId *big.Int
	Sender        common.Address
	IsMintingTx   bool
	Value         *big.Int
	Data          []byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterExecution is a free log retrieval operation binding the contract event 0x7ae9762dba14e21841a58a1cd988fb97578e27069c9bcc32d65682127141e730.
//
// Solidity: event Execution(transactionId indexed uint256, sender indexed address, isMintingTx bool, value uint256, data bytes)
func (_MyTRC21 *MyTRC21Filterer) FilterExecution(opts *bind.FilterOpts, transactionId []*big.Int, sender []common.Address) (*MyTRC21ExecutionIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Execution", transactionIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21ExecutionIterator{contract: _MyTRC21.contract, event: "Execution", logs: logs, sub: sub}, nil
}

// WatchExecution is a free log subscription operation binding the contract event 0x7ae9762dba14e21841a58a1cd988fb97578e27069c9bcc32d65682127141e730.
//
// Solidity: event Execution(transactionId indexed uint256, sender indexed address, isMintingTx bool, value uint256, data bytes)
func (_MyTRC21 *MyTRC21Filterer) WatchExecution(opts *bind.WatchOpts, sink chan<- *MyTRC21Execution, transactionId []*big.Int, sender []common.Address) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Execution", transactionIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Execution)
				if err := _MyTRC21.contract.UnpackLog(event, "Execution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21ExecutionFailureIterator is returned from FilterExecutionFailure and is used to iterate over the raw logs and unpacked data for ExecutionFailure events raised by the MyTRC21 contract.
type MyTRC21ExecutionFailureIterator struct {
	Event *MyTRC21ExecutionFailure // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21ExecutionFailureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21ExecutionFailure)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21ExecutionFailure)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21ExecutionFailureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21ExecutionFailureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21ExecutionFailure represents a ExecutionFailure event raised by the MyTRC21 contract.
type MyTRC21ExecutionFailure struct {
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterExecutionFailure is a free log retrieval operation binding the contract event 0x526441bb6c1aba3c9a4a6ca1d6545da9c2333c8c48343ef398eb858d72b79236.
//
// Solidity: event ExecutionFailure(transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterExecutionFailure(opts *bind.FilterOpts, transactionId []*big.Int) (*MyTRC21ExecutionFailureIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "ExecutionFailure", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21ExecutionFailureIterator{contract: _MyTRC21.contract, event: "ExecutionFailure", logs: logs, sub: sub}, nil
}

// WatchExecutionFailure is a free log subscription operation binding the contract event 0x526441bb6c1aba3c9a4a6ca1d6545da9c2333c8c48343ef398eb858d72b79236.
//
// Solidity: event ExecutionFailure(transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchExecutionFailure(opts *bind.WatchOpts, sink chan<- *MyTRC21ExecutionFailure, transactionId []*big.Int) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "ExecutionFailure", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21ExecutionFailure)
				if err := _MyTRC21.contract.UnpackLog(event, "ExecutionFailure", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21FeeIterator is returned from FilterFee and is used to iterate over the raw logs and unpacked data for Fee events raised by the MyTRC21 contract.
type MyTRC21FeeIterator struct {
	Event *MyTRC21Fee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21FeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Fee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Fee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21FeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21FeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Fee represents a Fee event raised by the MyTRC21 contract.
type MyTRC21Fee struct {
	From   common.Address
	To     common.Address
	Issuer common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFee is a free log retrieval operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterFee(opts *bind.FilterOpts, from []common.Address, to []common.Address, issuer []common.Address) (*MyTRC21FeeIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21FeeIterator{contract: _MyTRC21.contract, event: "Fee", logs: logs, sub: sub}, nil
}

// WatchFee is a free log subscription operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchFee(opts *bind.WatchOpts, sink chan<- *MyTRC21Fee, from []common.Address, to []common.Address, issuer []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Fee)
				if err := _MyTRC21.contract.UnpackLog(event, "Fee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21OwnerAdditionIterator is returned from FilterOwnerAddition and is used to iterate over the raw logs and unpacked data for OwnerAddition events raised by the MyTRC21 contract.
type MyTRC21OwnerAdditionIterator struct {
	Event *MyTRC21OwnerAddition // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21OwnerAdditionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21OwnerAddition)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21OwnerAddition)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21OwnerAdditionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21OwnerAdditionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21OwnerAddition represents a OwnerAddition event raised by the MyTRC21 contract.
type MyTRC21OwnerAddition struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerAddition is a free log retrieval operation binding the contract event 0xf39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d.
//
// Solidity: event OwnerAddition(owner indexed address)
func (_MyTRC21 *MyTRC21Filterer) FilterOwnerAddition(opts *bind.FilterOpts, owner []common.Address) (*MyTRC21OwnerAdditionIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "OwnerAddition", ownerRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21OwnerAdditionIterator{contract: _MyTRC21.contract, event: "OwnerAddition", logs: logs, sub: sub}, nil
}

// WatchOwnerAddition is a free log subscription operation binding the contract event 0xf39e6e1eb0edcf53c221607b54b00cd28f3196fed0a24994dc308b8f611b682d.
//
// Solidity: event OwnerAddition(owner indexed address)
func (_MyTRC21 *MyTRC21Filterer) WatchOwnerAddition(opts *bind.WatchOpts, sink chan<- *MyTRC21OwnerAddition, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "OwnerAddition", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21OwnerAddition)
				if err := _MyTRC21.contract.UnpackLog(event, "OwnerAddition", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21OwnerRemovalIterator is returned from FilterOwnerRemoval and is used to iterate over the raw logs and unpacked data for OwnerRemoval events raised by the MyTRC21 contract.
type MyTRC21OwnerRemovalIterator struct {
	Event *MyTRC21OwnerRemoval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21OwnerRemovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21OwnerRemoval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21OwnerRemoval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21OwnerRemovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21OwnerRemovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21OwnerRemoval represents a OwnerRemoval event raised by the MyTRC21 contract.
type MyTRC21OwnerRemoval struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerRemoval is a free log retrieval operation binding the contract event 0x8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b90.
//
// Solidity: event OwnerRemoval(owner indexed address)
func (_MyTRC21 *MyTRC21Filterer) FilterOwnerRemoval(opts *bind.FilterOpts, owner []common.Address) (*MyTRC21OwnerRemovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "OwnerRemoval", ownerRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21OwnerRemovalIterator{contract: _MyTRC21.contract, event: "OwnerRemoval", logs: logs, sub: sub}, nil
}

// WatchOwnerRemoval is a free log subscription operation binding the contract event 0x8001553a916ef2f495d26a907cc54d96ed840d7bda71e73194bf5a9df7a76b90.
//
// Solidity: event OwnerRemoval(owner indexed address)
func (_MyTRC21 *MyTRC21Filterer) WatchOwnerRemoval(opts *bind.WatchOpts, sink chan<- *MyTRC21OwnerRemoval, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "OwnerRemoval", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21OwnerRemoval)
				if err := _MyTRC21.contract.UnpackLog(event, "OwnerRemoval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21RequirementChangeIterator is returned from FilterRequirementChange and is used to iterate over the raw logs and unpacked data for RequirementChange events raised by the MyTRC21 contract.
type MyTRC21RequirementChangeIterator struct {
	Event *MyTRC21RequirementChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21RequirementChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21RequirementChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21RequirementChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21RequirementChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21RequirementChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21RequirementChange represents a RequirementChange event raised by the MyTRC21 contract.
type MyTRC21RequirementChange struct {
	Required *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRequirementChange is a free log retrieval operation binding the contract event 0xa3f1ee9126a074d9326c682f561767f710e927faa811f7a99829d49dc421797a.
//
// Solidity: event RequirementChange(required uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterRequirementChange(opts *bind.FilterOpts) (*MyTRC21RequirementChangeIterator, error) {

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "RequirementChange")
	if err != nil {
		return nil, err
	}
	return &MyTRC21RequirementChangeIterator{contract: _MyTRC21.contract, event: "RequirementChange", logs: logs, sub: sub}, nil
}

// WatchRequirementChange is a free log subscription operation binding the contract event 0xa3f1ee9126a074d9326c682f561767f710e927faa811f7a99829d49dc421797a.
//
// Solidity: event RequirementChange(required uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchRequirementChange(opts *bind.WatchOpts, sink chan<- *MyTRC21RequirementChange) (event.Subscription, error) {

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "RequirementChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21RequirementChange)
				if err := _MyTRC21.contract.UnpackLog(event, "RequirementChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21RevocationIterator is returned from FilterRevocation and is used to iterate over the raw logs and unpacked data for Revocation events raised by the MyTRC21 contract.
type MyTRC21RevocationIterator struct {
	Event *MyTRC21Revocation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21RevocationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Revocation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Revocation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21RevocationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21RevocationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Revocation represents a Revocation event raised by the MyTRC21 contract.
type MyTRC21Revocation struct {
	Sender        common.Address
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRevocation is a free log retrieval operation binding the contract event 0xf6a317157440607f36269043eb55f1287a5a19ba2216afeab88cd46cbcfb88e9.
//
// Solidity: event Revocation(sender indexed address, transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterRevocation(opts *bind.FilterOpts, sender []common.Address, transactionId []*big.Int) (*MyTRC21RevocationIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Revocation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21RevocationIterator{contract: _MyTRC21.contract, event: "Revocation", logs: logs, sub: sub}, nil
}

// WatchRevocation is a free log subscription operation binding the contract event 0xf6a317157440607f36269043eb55f1287a5a19ba2216afeab88cd46cbcfb88e9.
//
// Solidity: event Revocation(sender indexed address, transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchRevocation(opts *bind.WatchOpts, sink chan<- *MyTRC21Revocation, sender []common.Address, transactionId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Revocation", senderRule, transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Revocation)
				if err := _MyTRC21.contract.UnpackLog(event, "Revocation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21SubmissionIterator is returned from FilterSubmission and is used to iterate over the raw logs and unpacked data for Submission events raised by the MyTRC21 contract.
type MyTRC21SubmissionIterator struct {
	Event *MyTRC21Submission // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21SubmissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Submission)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Submission)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21SubmissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21SubmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Submission represents a Submission event raised by the MyTRC21 contract.
type MyTRC21Submission struct {
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSubmission is a free log retrieval operation binding the contract event 0xc0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e51.
//
// Solidity: event Submission(transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterSubmission(opts *bind.FilterOpts, transactionId []*big.Int) (*MyTRC21SubmissionIterator, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Submission", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21SubmissionIterator{contract: _MyTRC21.contract, event: "Submission", logs: logs, sub: sub}, nil
}

// WatchSubmission is a free log subscription operation binding the contract event 0xc0ba8fe4b176c1714197d43b9cc6bcf797a4a7461c5fe8d0ef6e184ae7601e51.
//
// Solidity: event Submission(transactionId indexed uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchSubmission(opts *bind.WatchOpts, sink chan<- *MyTRC21Submission, transactionId []*big.Int) (event.Subscription, error) {

	var transactionIdRule []interface{}
	for _, transactionIdItem := range transactionId {
		transactionIdRule = append(transactionIdRule, transactionIdItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Submission", transactionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Submission)
				if err := _MyTRC21.contract.UnpackLog(event, "Submission", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// MyTRC21TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MyTRC21 contract.
type MyTRC21TransferIterator struct {
	Event *MyTRC21Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MyTRC21TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MyTRC21Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MyTRC21Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MyTRC21TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MyTRC21TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MyTRC21Transfer represents a Transfer event raised by the MyTRC21 contract.
type MyTRC21Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MyTRC21TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MyTRC21.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MyTRC21TransferIterator{contract: _MyTRC21.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_MyTRC21 *MyTRC21Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MyTRC21Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MyTRC21.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MyTRC21Transfer)
				if err := _MyTRC21.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TRC21ABI is the input ABI used to generate the binding from.
const TRC21ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"estimateFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setMinFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Fee\",\"type\":\"event\"}]"

// TRC21Bin is the compiled bytecode used for deploying new contracts.
const TRC21Bin = `0x`

// DeployTRC21 deploys a new Ethereum contract, binding an instance of TRC21 to it.
func DeployTRC21(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TRC21, error) {
	parsed, err := abi.JSON(strings.NewReader(TRC21ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TRC21Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TRC21{TRC21Caller: TRC21Caller{contract: contract}, TRC21Transactor: TRC21Transactor{contract: contract}, TRC21Filterer: TRC21Filterer{contract: contract}}, nil
}

// TRC21 is an auto generated Go binding around an Ethereum contract.
type TRC21 struct {
	TRC21Caller     // Read-only binding to the contract
	TRC21Transactor // Write-only binding to the contract
	TRC21Filterer   // Log filterer for contract events
}

// TRC21Caller is an auto generated read-only Go binding around an Ethereum contract.
type TRC21Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TRC21Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TRC21Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TRC21Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TRC21Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TRC21Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TRC21Session struct {
	Contract     *TRC21            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TRC21CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TRC21CallerSession struct {
	Contract *TRC21Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TRC21TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TRC21TransactorSession struct {
	Contract     *TRC21Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TRC21Raw is an auto generated low-level Go binding around an Ethereum contract.
type TRC21Raw struct {
	Contract *TRC21 // Generic contract binding to access the raw methods on
}

// TRC21CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TRC21CallerRaw struct {
	Contract *TRC21Caller // Generic read-only contract binding to access the raw methods on
}

// TRC21TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TRC21TransactorRaw struct {
	Contract *TRC21Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTRC21 creates a new instance of TRC21, bound to a specific deployed contract.
func NewTRC21(address common.Address, backend bind.ContractBackend) (*TRC21, error) {
	contract, err := bindTRC21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TRC21{TRC21Caller: TRC21Caller{contract: contract}, TRC21Transactor: TRC21Transactor{contract: contract}, TRC21Filterer: TRC21Filterer{contract: contract}}, nil
}

// NewTRC21Caller creates a new read-only instance of TRC21, bound to a specific deployed contract.
func NewTRC21Caller(address common.Address, caller bind.ContractCaller) (*TRC21Caller, error) {
	contract, err := bindTRC21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TRC21Caller{contract: contract}, nil
}

// NewTRC21Transactor creates a new write-only instance of TRC21, bound to a specific deployed contract.
func NewTRC21Transactor(address common.Address, transactor bind.ContractTransactor) (*TRC21Transactor, error) {
	contract, err := bindTRC21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TRC21Transactor{contract: contract}, nil
}

// NewTRC21Filterer creates a new log filterer instance of TRC21, bound to a specific deployed contract.
func NewTRC21Filterer(address common.Address, filterer bind.ContractFilterer) (*TRC21Filterer, error) {
	contract, err := bindTRC21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TRC21Filterer{contract: contract}, nil
}

// bindTRC21 binds a generic wrapper to an already deployed contract.
func bindTRC21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TRC21ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TRC21 *TRC21Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TRC21.Contract.TRC21Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TRC21 *TRC21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TRC21.Contract.TRC21Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TRC21 *TRC21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TRC21.Contract.TRC21Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TRC21 *TRC21CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TRC21.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TRC21 *TRC21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TRC21.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TRC21 *TRC21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TRC21.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_TRC21 *TRC21Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_TRC21 *TRC21Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TRC21.Contract.Allowance(&_TRC21.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_TRC21 *TRC21CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TRC21.Contract.Allowance(&_TRC21.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_TRC21 *TRC21Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "balanceOf", owner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_TRC21 *TRC21Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TRC21.Contract.BalanceOf(&_TRC21.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(owner address) constant returns(uint256)
func (_TRC21 *TRC21CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TRC21.Contract.BalanceOf(&_TRC21.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TRC21 *TRC21Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TRC21 *TRC21Session) Decimals() (uint8, error) {
	return _TRC21.Contract.Decimals(&_TRC21.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_TRC21 *TRC21CallerSession) Decimals() (uint8, error) {
	return _TRC21.Contract.Decimals(&_TRC21.CallOpts)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_TRC21 *TRC21Caller) EstimateFee(opts *bind.CallOpts, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "estimateFee", value)
	return *ret0, err
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_TRC21 *TRC21Session) EstimateFee(value *big.Int) (*big.Int, error) {
	return _TRC21.Contract.EstimateFee(&_TRC21.CallOpts, value)
}

// EstimateFee is a free data retrieval call binding the contract method 0x127e8e4d.
//
// Solidity: function estimateFee(value uint256) constant returns(uint256)
func (_TRC21 *TRC21CallerSession) EstimateFee(value *big.Int) (*big.Int, error) {
	return _TRC21.Contract.EstimateFee(&_TRC21.CallOpts, value)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_TRC21 *TRC21Caller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "issuer")
	return *ret0, err
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_TRC21 *TRC21Session) Issuer() (common.Address, error) {
	return _TRC21.Contract.Issuer(&_TRC21.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_TRC21 *TRC21CallerSession) Issuer() (common.Address, error) {
	return _TRC21.Contract.Issuer(&_TRC21.CallOpts)
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_TRC21 *TRC21Caller) MinFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "minFee")
	return *ret0, err
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_TRC21 *TRC21Session) MinFee() (*big.Int, error) {
	return _TRC21.Contract.MinFee(&_TRC21.CallOpts)
}

// MinFee is a free data retrieval call binding the contract method 0x24ec7590.
//
// Solidity: function minFee() constant returns(uint256)
func (_TRC21 *TRC21CallerSession) MinFee() (*big.Int, error) {
	return _TRC21.Contract.MinFee(&_TRC21.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TRC21 *TRC21Caller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TRC21 *TRC21Session) Name() (string, error) {
	return _TRC21.Contract.Name(&_TRC21.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_TRC21 *TRC21CallerSession) Name() (string, error) {
	return _TRC21.Contract.Name(&_TRC21.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TRC21 *TRC21Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TRC21 *TRC21Session) Symbol() (string, error) {
	return _TRC21.Contract.Symbol(&_TRC21.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_TRC21 *TRC21CallerSession) Symbol() (string, error) {
	return _TRC21.Contract.Symbol(&_TRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TRC21 *TRC21Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _TRC21.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TRC21 *TRC21Session) TotalSupply() (*big.Int, error) {
	return _TRC21.Contract.TotalSupply(&_TRC21.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_TRC21 *TRC21CallerSession) TotalSupply() (*big.Int, error) {
	return _TRC21.Contract.TotalSupply(&_TRC21.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_TRC21 *TRC21Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_TRC21 *TRC21Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.Approve(&_TRC21.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns(bool)
func (_TRC21 *TRC21TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.Approve(&_TRC21.TransactOpts, spender, value)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_TRC21 *TRC21Transactor) SetMinFee(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _TRC21.contract.Transact(opts, "setMinFee", value)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_TRC21 *TRC21Session) SetMinFee(value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.SetMinFee(&_TRC21.TransactOpts, value)
}

// SetMinFee is a paid mutator transaction binding the contract method 0x31ac9920.
//
// Solidity: function setMinFee(value uint256) returns()
func (_TRC21 *TRC21TransactorSession) SetMinFee(value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.SetMinFee(&_TRC21.TransactOpts, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_TRC21 *TRC21Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_TRC21 *TRC21Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.Transfer(&_TRC21.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns(bool)
func (_TRC21 *TRC21TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.Transfer(&_TRC21.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_TRC21 *TRC21Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_TRC21 *TRC21Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.TransferFrom(&_TRC21.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns(bool)
func (_TRC21 *TRC21TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TRC21.Contract.TransferFrom(&_TRC21.TransactOpts, from, to, value)
}

// TRC21ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TRC21 contract.
type TRC21ApprovalIterator struct {
	Event *TRC21Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TRC21ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TRC21Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TRC21Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TRC21ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TRC21ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TRC21Approval represents a Approval event raised by the TRC21 contract.
type TRC21Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_TRC21 *TRC21Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TRC21ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TRC21.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TRC21ApprovalIterator{contract: _TRC21.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(owner indexed address, spender indexed address, value uint256)
func (_TRC21 *TRC21Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TRC21Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TRC21.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TRC21Approval)
				if err := _TRC21.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TRC21FeeIterator is returned from FilterFee and is used to iterate over the raw logs and unpacked data for Fee events raised by the TRC21 contract.
type TRC21FeeIterator struct {
	Event *TRC21Fee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TRC21FeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TRC21Fee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TRC21Fee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TRC21FeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TRC21FeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TRC21Fee represents a Fee event raised by the TRC21 contract.
type TRC21Fee struct {
	From   common.Address
	To     common.Address
	Issuer common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFee is a free log retrieval operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_TRC21 *TRC21Filterer) FilterFee(opts *bind.FilterOpts, from []common.Address, to []common.Address, issuer []common.Address) (*TRC21FeeIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _TRC21.contract.FilterLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return &TRC21FeeIterator{contract: _TRC21.contract, event: "Fee", logs: logs, sub: sub}, nil
}

// WatchFee is a free log subscription operation binding the contract event 0xfcf5b3276434181e3c48bd3fe569b8939808f11e845d4b963693b9d7dbd6dd99.
//
// Solidity: event Fee(from indexed address, to indexed address, issuer indexed address, value uint256)
func (_TRC21 *TRC21Filterer) WatchFee(opts *bind.WatchOpts, sink chan<- *TRC21Fee, from []common.Address, to []common.Address, issuer []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _TRC21.contract.WatchLogs(opts, "Fee", fromRule, toRule, issuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TRC21Fee)
				if err := _TRC21.contract.UnpackLog(event, "Fee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// TRC21TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TRC21 contract.
type TRC21TransferIterator struct {
	Event *TRC21Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TRC21TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TRC21Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TRC21Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TRC21TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TRC21TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TRC21Transfer represents a Transfer event raised by the TRC21 contract.
type TRC21Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_TRC21 *TRC21Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TRC21TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TRC21.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TRC21TransferIterator{contract: _TRC21.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(from indexed address, to indexed address, value uint256)
func (_TRC21 *TRC21Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TRC21Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TRC21.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TRC21Transfer)
				if err := _TRC21.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
