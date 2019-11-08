// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

// TOMOXListingABI is the input ABI used to generate the binding from.
const TOMOXListingABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"apply\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// TOMOXListingBin is the compiled bytecode used for deploying new contracts.
const TOMOXListingBin = `0x608060405234801561001057600080fd5b5061027b806100206000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416639d63848a811461005b578063a3ff31b5146100c0578063c6b32f34146100f5575b600080fd5b34801561006757600080fd5b5061007061010b565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156100ac578181015183820152602001610094565b505050509050019250505060405180910390f35b3480156100cc57600080fd5b506100e1600160a060020a036004351661016d565b604080519115158252519081900360200190f35b610109600160a060020a036004351661018b565b005b6060600080548060200260200160405190810160405280929190818152602001828054801561016357602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610145575b5050505050905090565b600160a060020a031660009081526001602052604090205460ff1690565b80600160a060020a03811615156101a157600080fd5b600160a060020a03811660009081526001602081905260409091205460ff16151514156101cd57600080fd5b5060008054600180820183557f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563909101805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039490941693841790556040805160208082018352838252948452919093529190209051815460ff19169015151790555600a165627a7a72305820bbfa9118404af15e7ad54981a507fb7d558745e19d4195847a04f82a410101a30029`

// DeployTOMOXListing deploys a new Ethereum contract, binding an instance of TOMOXListing to it.
func DeployTOMOXListing(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TOMOXListing, error) {
	parsed, err := abi.JSON(strings.NewReader(TOMOXListingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TOMOXListingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TOMOXListing{TOMOXListingCaller: TOMOXListingCaller{contract: contract}, TOMOXListingTransactor: TOMOXListingTransactor{contract: contract}, TOMOXListingFilterer: TOMOXListingFilterer{contract: contract}}, nil
}

// TOMOXListing is an auto generated Go binding around an Ethereum contract.
type TOMOXListing struct {
	TOMOXListingCaller     // Read-only binding to the contract
	TOMOXListingTransactor // Write-only binding to the contract
	TOMOXListingFilterer   // Log filterer for contract events
}

// TOMOXListingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TOMOXListingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TOMOXListingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TOMOXListingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TOMOXListingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TOMOXListingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TOMOXListingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TOMOXListingSession struct {
	Contract     *TOMOXListing     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TOMOXListingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TOMOXListingCallerSession struct {
	Contract *TOMOXListingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TOMOXListingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TOMOXListingTransactorSession struct {
	Contract     *TOMOXListingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TOMOXListingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TOMOXListingRaw struct {
	Contract *TOMOXListing // Generic contract binding to access the raw methods on
}

// TOMOXListingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TOMOXListingCallerRaw struct {
	Contract *TOMOXListingCaller // Generic read-only contract binding to access the raw methods on
}

// TOMOXListingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TOMOXListingTransactorRaw struct {
	Contract *TOMOXListingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTOMOXListing creates a new instance of TOMOXListing, bound to a specific deployed contract.
func NewTOMOXListing(address common.Address, backend bind.ContractBackend) (*TOMOXListing, error) {
	contract, err := bindTOMOXListing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TOMOXListing{TOMOXListingCaller: TOMOXListingCaller{contract: contract}, TOMOXListingTransactor: TOMOXListingTransactor{contract: contract}, TOMOXListingFilterer: TOMOXListingFilterer{contract: contract}}, nil
}

// NewTOMOXListingCaller creates a new read-only instance of TOMOXListing, bound to a specific deployed contract.
func NewTOMOXListingCaller(address common.Address, caller bind.ContractCaller) (*TOMOXListingCaller, error) {
	contract, err := bindTOMOXListing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TOMOXListingCaller{contract: contract}, nil
}

// NewTOMOXListingTransactor creates a new write-only instance of TOMOXListing, bound to a specific deployed contract.
func NewTOMOXListingTransactor(address common.Address, transactor bind.ContractTransactor) (*TOMOXListingTransactor, error) {
	contract, err := bindTOMOXListing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TOMOXListingTransactor{contract: contract}, nil
}

// NewTOMOXListingFilterer creates a new log filterer instance of TOMOXListing, bound to a specific deployed contract.
func NewTOMOXListingFilterer(address common.Address, filterer bind.ContractFilterer) (*TOMOXListingFilterer, error) {
	contract, err := bindTOMOXListing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TOMOXListingFilterer{contract: contract}, nil
}

// bindTOMOXListing binds a generic wrapper to an already deployed contract.
func bindTOMOXListing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TOMOXListingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TOMOXListing *TOMOXListingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TOMOXListing.Contract.TOMOXListingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TOMOXListing *TOMOXListingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TOMOXListing.Contract.TOMOXListingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TOMOXListing *TOMOXListingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TOMOXListing.Contract.TOMOXListingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TOMOXListing *TOMOXListingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _TOMOXListing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TOMOXListing *TOMOXListingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TOMOXListing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TOMOXListing *TOMOXListingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TOMOXListing.Contract.contract.Transact(opts, method, params...)
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus(token address) constant returns(bool)
func (_TOMOXListing *TOMOXListingCaller) GetTokenStatus(opts *bind.CallOpts, token common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _TOMOXListing.contract.Call(opts, out, "getTokenStatus", token)
	return *ret0, err
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus(token address) constant returns(bool)
func (_TOMOXListing *TOMOXListingSession) GetTokenStatus(token common.Address) (bool, error) {
	return _TOMOXListing.Contract.GetTokenStatus(&_TOMOXListing.CallOpts, token)
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus(token address) constant returns(bool)
func (_TOMOXListing *TOMOXListingCallerSession) GetTokenStatus(token common.Address) (bool, error) {
	return _TOMOXListing.Contract.GetTokenStatus(&_TOMOXListing.CallOpts, token)
}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() constant returns(address[])
func (_TOMOXListing *TOMOXListingCaller) Tokens(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _TOMOXListing.contract.Call(opts, out, "tokens")
	return *ret0, err
}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() constant returns(address[])
func (_TOMOXListing *TOMOXListingSession) Tokens() ([]common.Address, error) {
	return _TOMOXListing.Contract.Tokens(&_TOMOXListing.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x9d63848a.
//
// Solidity: function tokens() constant returns(address[])
func (_TOMOXListing *TOMOXListingCallerSession) Tokens() ([]common.Address, error) {
	return _TOMOXListing.Contract.Tokens(&_TOMOXListing.CallOpts)
}

// Apply is a paid mutator transaction binding the contract method 0xc6b32f34.
//
// Solidity: function apply(token address) returns()
func (_TOMOXListing *TOMOXListingTransactor) Apply(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _TOMOXListing.contract.Transact(opts, "apply", token)
}

// Apply is a paid mutator transaction binding the contract method 0xc6b32f34.
//
// Solidity: function apply(token address) returns()
func (_TOMOXListing *TOMOXListingSession) Apply(token common.Address) (*types.Transaction, error) {
	return _TOMOXListing.Contract.Apply(&_TOMOXListing.TransactOpts, token)
}

// Apply is a paid mutator transaction binding the contract method 0xc6b32f34.
//
// Solidity: function apply(token address) returns()
func (_TOMOXListing *TOMOXListingTransactorSession) Apply(token common.Address) (*types.Transaction, error) {
	return _TOMOXListing.Contract.Apply(&_TOMOXListing.TransactOpts, token)
}
