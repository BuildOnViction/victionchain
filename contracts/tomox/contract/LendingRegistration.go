// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	"github.com/tomochain/tomochain/accounts/abi"
	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
)

// LAbstractRegistrationABI is the input ABI used to generate the binding from.
const LAbstractRegistrationABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"RESIGN_REQUESTS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"getRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LAbstractRegistrationBin is the compiled bytecode used for deploying new contracts.
const LAbstractRegistrationBin = `0x`

// DeployLAbstractRegistration deploys a new Ethereum contract, binding an instance of LAbstractRegistration to it.
func DeployLAbstractRegistration(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LAbstractRegistration, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractRegistrationABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LAbstractRegistrationBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LAbstractRegistration{LAbstractRegistrationCaller: LAbstractRegistrationCaller{contract: contract}, LAbstractRegistrationTransactor: LAbstractRegistrationTransactor{contract: contract}, LAbstractRegistrationFilterer: LAbstractRegistrationFilterer{contract: contract}}, nil
}

// LAbstractRegistration is an auto generated Go binding around an Ethereum contract.
type LAbstractRegistration struct {
	LAbstractRegistrationCaller     // Read-only binding to the contract
	LAbstractRegistrationTransactor // Write-only binding to the contract
	LAbstractRegistrationFilterer   // Log filterer for contract events
}

// LAbstractRegistrationCaller is an auto generated read-only Go binding around an Ethereum contract.
type LAbstractRegistrationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractRegistrationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LAbstractRegistrationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractRegistrationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LAbstractRegistrationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractRegistrationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LAbstractRegistrationSession struct {
	Contract     *LAbstractRegistration // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LAbstractRegistrationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LAbstractRegistrationCallerSession struct {
	Contract *LAbstractRegistrationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// LAbstractRegistrationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LAbstractRegistrationTransactorSession struct {
	Contract     *LAbstractRegistrationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// LAbstractRegistrationRaw is an auto generated low-level Go binding around an Ethereum contract.
type LAbstractRegistrationRaw struct {
	Contract *LAbstractRegistration // Generic contract binding to access the raw methods on
}

// LAbstractRegistrationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LAbstractRegistrationCallerRaw struct {
	Contract *LAbstractRegistrationCaller // Generic read-only contract binding to access the raw methods on
}

// LAbstractRegistrationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LAbstractRegistrationTransactorRaw struct {
	Contract *LAbstractRegistrationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLAbstractRegistration creates a new instance of LAbstractRegistration, bound to a specific deployed contract.
func NewLAbstractRegistration(address common.Address, backend bind.ContractBackend) (*LAbstractRegistration, error) {
	contract, err := bindLAbstractRegistration(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LAbstractRegistration{LAbstractRegistrationCaller: LAbstractRegistrationCaller{contract: contract}, LAbstractRegistrationTransactor: LAbstractRegistrationTransactor{contract: contract}, LAbstractRegistrationFilterer: LAbstractRegistrationFilterer{contract: contract}}, nil
}

// NewLAbstractRegistrationCaller creates a new read-only instance of LAbstractRegistration, bound to a specific deployed contract.
func NewLAbstractRegistrationCaller(address common.Address, caller bind.ContractCaller) (*LAbstractRegistrationCaller, error) {
	contract, err := bindLAbstractRegistration(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractRegistrationCaller{contract: contract}, nil
}

// NewLAbstractRegistrationTransactor creates a new write-only instance of LAbstractRegistration, bound to a specific deployed contract.
func NewLAbstractRegistrationTransactor(address common.Address, transactor bind.ContractTransactor) (*LAbstractRegistrationTransactor, error) {
	contract, err := bindLAbstractRegistration(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractRegistrationTransactor{contract: contract}, nil
}

// NewLAbstractRegistrationFilterer creates a new log filterer instance of LAbstractRegistration, bound to a specific deployed contract.
func NewLAbstractRegistrationFilterer(address common.Address, filterer bind.ContractFilterer) (*LAbstractRegistrationFilterer, error) {
	contract, err := bindLAbstractRegistration(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LAbstractRegistrationFilterer{contract: contract}, nil
}

// bindLAbstractRegistration binds a generic wrapper to an already deployed contract.
func bindLAbstractRegistration(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractRegistrationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractRegistration *LAbstractRegistrationRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractRegistration.Contract.LAbstractRegistrationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractRegistration *LAbstractRegistrationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractRegistration.Contract.LAbstractRegistrationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractRegistration *LAbstractRegistrationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractRegistration.Contract.LAbstractRegistrationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractRegistration *LAbstractRegistrationCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractRegistration.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractRegistration *LAbstractRegistrationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractRegistration.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractRegistration *LAbstractRegistrationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractRegistration.Contract.contract.Transact(opts, method, params...)
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_LAbstractRegistration *LAbstractRegistrationCaller) RESIGNREQUESTS(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LAbstractRegistration.contract.Call(opts, out, "RESIGN_REQUESTS", arg0)
	return *ret0, err
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_LAbstractRegistration *LAbstractRegistrationSession) RESIGNREQUESTS(arg0 common.Address) (*big.Int, error) {
	return _LAbstractRegistration.Contract.RESIGNREQUESTS(&_LAbstractRegistration.CallOpts, arg0)
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_LAbstractRegistration *LAbstractRegistrationCallerSession) RESIGNREQUESTS(arg0 common.Address) (*big.Int, error) {
	return _LAbstractRegistration.Contract.RESIGNREQUESTS(&_LAbstractRegistration.CallOpts, arg0)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_LAbstractRegistration *LAbstractRegistrationCaller) GetRelayerByCoinbase(opts *bind.CallOpts, arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(common.Address)
		ret2 = new(*big.Int)
		ret3 = new(uint16)
		ret4 = new([]common.Address)
		ret5 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
	}
	err := _LAbstractRegistration.contract.Call(opts, out, "getRelayerByCoinbase", arg0)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, err
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_LAbstractRegistration *LAbstractRegistrationSession) GetRelayerByCoinbase(arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _LAbstractRegistration.Contract.GetRelayerByCoinbase(&_LAbstractRegistration.CallOpts, arg0)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_LAbstractRegistration *LAbstractRegistrationCallerSession) GetRelayerByCoinbase(arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _LAbstractRegistration.Contract.GetRelayerByCoinbase(&_LAbstractRegistration.CallOpts, arg0)
}

// LAbstractTOMOXListingABI is the input ABI used to generate the binding from.
const LAbstractTOMOXListingABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"getTokenStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LAbstractTOMOXListingBin is the compiled bytecode used for deploying new contracts.
const LAbstractTOMOXListingBin = `0x`

// DeployLAbstractTOMOXListing deploys a new Ethereum contract, binding an instance of LAbstractTOMOXListing to it.
func DeployLAbstractTOMOXListing(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LAbstractTOMOXListing, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractTOMOXListingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LAbstractTOMOXListingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LAbstractTOMOXListing{LAbstractTOMOXListingCaller: LAbstractTOMOXListingCaller{contract: contract}, LAbstractTOMOXListingTransactor: LAbstractTOMOXListingTransactor{contract: contract}, LAbstractTOMOXListingFilterer: LAbstractTOMOXListingFilterer{contract: contract}}, nil
}

// LAbstractTOMOXListing is an auto generated Go binding around an Ethereum contract.
type LAbstractTOMOXListing struct {
	LAbstractTOMOXListingCaller     // Read-only binding to the contract
	LAbstractTOMOXListingTransactor // Write-only binding to the contract
	LAbstractTOMOXListingFilterer   // Log filterer for contract events
}

// LAbstractTOMOXListingCaller is an auto generated read-only Go binding around an Ethereum contract.
type LAbstractTOMOXListingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTOMOXListingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LAbstractTOMOXListingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTOMOXListingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LAbstractTOMOXListingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTOMOXListingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LAbstractTOMOXListingSession struct {
	Contract     *LAbstractTOMOXListing // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LAbstractTOMOXListingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LAbstractTOMOXListingCallerSession struct {
	Contract *LAbstractTOMOXListingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// LAbstractTOMOXListingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LAbstractTOMOXListingTransactorSession struct {
	Contract     *LAbstractTOMOXListingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// LAbstractTOMOXListingRaw is an auto generated low-level Go binding around an Ethereum contract.
type LAbstractTOMOXListingRaw struct {
	Contract *LAbstractTOMOXListing // Generic contract binding to access the raw methods on
}

// LAbstractTOMOXListingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LAbstractTOMOXListingCallerRaw struct {
	Contract *LAbstractTOMOXListingCaller // Generic read-only contract binding to access the raw methods on
}

// LAbstractTOMOXListingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LAbstractTOMOXListingTransactorRaw struct {
	Contract *LAbstractTOMOXListingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLAbstractTOMOXListing creates a new instance of LAbstractTOMOXListing, bound to a specific deployed contract.
func NewLAbstractTOMOXListing(address common.Address, backend bind.ContractBackend) (*LAbstractTOMOXListing, error) {
	contract, err := bindLAbstractTOMOXListing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LAbstractTOMOXListing{LAbstractTOMOXListingCaller: LAbstractTOMOXListingCaller{contract: contract}, LAbstractTOMOXListingTransactor: LAbstractTOMOXListingTransactor{contract: contract}, LAbstractTOMOXListingFilterer: LAbstractTOMOXListingFilterer{contract: contract}}, nil
}

// NewLAbstractTOMOXListingCaller creates a new read-only instance of LAbstractTOMOXListing, bound to a specific deployed contract.
func NewLAbstractTOMOXListingCaller(address common.Address, caller bind.ContractCaller) (*LAbstractTOMOXListingCaller, error) {
	contract, err := bindLAbstractTOMOXListing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractTOMOXListingCaller{contract: contract}, nil
}

// NewLAbstractTOMOXListingTransactor creates a new write-only instance of LAbstractTOMOXListing, bound to a specific deployed contract.
func NewLAbstractTOMOXListingTransactor(address common.Address, transactor bind.ContractTransactor) (*LAbstractTOMOXListingTransactor, error) {
	contract, err := bindLAbstractTOMOXListing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractTOMOXListingTransactor{contract: contract}, nil
}

// NewLAbstractTOMOXListingFilterer creates a new log filterer instance of LAbstractTOMOXListing, bound to a specific deployed contract.
func NewLAbstractTOMOXListingFilterer(address common.Address, filterer bind.ContractFilterer) (*LAbstractTOMOXListingFilterer, error) {
	contract, err := bindLAbstractTOMOXListing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LAbstractTOMOXListingFilterer{contract: contract}, nil
}

// bindLAbstractTOMOXListing binds a generic wrapper to an already deployed contract.
func bindLAbstractTOMOXListing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractTOMOXListingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractTOMOXListing.Contract.LAbstractTOMOXListingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractTOMOXListing.Contract.LAbstractTOMOXListingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractTOMOXListing.Contract.LAbstractTOMOXListingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractTOMOXListing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractTOMOXListing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractTOMOXListing *LAbstractTOMOXListingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractTOMOXListing.Contract.contract.Transact(opts, method, params...)
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus( address) constant returns(bool)
func (_LAbstractTOMOXListing *LAbstractTOMOXListingCaller) GetTokenStatus(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LAbstractTOMOXListing.contract.Call(opts, out, "getTokenStatus", arg0)
	return *ret0, err
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus( address) constant returns(bool)
func (_LAbstractTOMOXListing *LAbstractTOMOXListingSession) GetTokenStatus(arg0 common.Address) (bool, error) {
	return _LAbstractTOMOXListing.Contract.GetTokenStatus(&_LAbstractTOMOXListing.CallOpts, arg0)
}

// GetTokenStatus is a free data retrieval call binding the contract method 0xa3ff31b5.
//
// Solidity: function getTokenStatus( address) constant returns(bool)
func (_LAbstractTOMOXListing *LAbstractTOMOXListingCallerSession) GetTokenStatus(arg0 common.Address) (bool, error) {
	return _LAbstractTOMOXListing.Contract.GetTokenStatus(&_LAbstractTOMOXListing.CallOpts, arg0)
}

// LAbstractTokenTRC21ABI is the input ABI used to generate the binding from.
const LAbstractTokenTRC21ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"issuer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LAbstractTokenTRC21Bin is the compiled bytecode used for deploying new contracts.
const LAbstractTokenTRC21Bin = `0x`

// DeployLAbstractTokenTRC21 deploys a new Ethereum contract, binding an instance of LAbstractTokenTRC21 to it.
func DeployLAbstractTokenTRC21(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LAbstractTokenTRC21, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractTokenTRC21ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LAbstractTokenTRC21Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LAbstractTokenTRC21{LAbstractTokenTRC21Caller: LAbstractTokenTRC21Caller{contract: contract}, LAbstractTokenTRC21Transactor: LAbstractTokenTRC21Transactor{contract: contract}, LAbstractTokenTRC21Filterer: LAbstractTokenTRC21Filterer{contract: contract}}, nil
}

// LAbstractTokenTRC21 is an auto generated Go binding around an Ethereum contract.
type LAbstractTokenTRC21 struct {
	LAbstractTokenTRC21Caller     // Read-only binding to the contract
	LAbstractTokenTRC21Transactor // Write-only binding to the contract
	LAbstractTokenTRC21Filterer   // Log filterer for contract events
}

// LAbstractTokenTRC21Caller is an auto generated read-only Go binding around an Ethereum contract.
type LAbstractTokenTRC21Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTokenTRC21Transactor is an auto generated write-only Go binding around an Ethereum contract.
type LAbstractTokenTRC21Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTokenTRC21Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LAbstractTokenTRC21Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LAbstractTokenTRC21Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LAbstractTokenTRC21Session struct {
	Contract     *LAbstractTokenTRC21 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// LAbstractTokenTRC21CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LAbstractTokenTRC21CallerSession struct {
	Contract *LAbstractTokenTRC21Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// LAbstractTokenTRC21TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LAbstractTokenTRC21TransactorSession struct {
	Contract     *LAbstractTokenTRC21Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// LAbstractTokenTRC21Raw is an auto generated low-level Go binding around an Ethereum contract.
type LAbstractTokenTRC21Raw struct {
	Contract *LAbstractTokenTRC21 // Generic contract binding to access the raw methods on
}

// LAbstractTokenTRC21CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LAbstractTokenTRC21CallerRaw struct {
	Contract *LAbstractTokenTRC21Caller // Generic read-only contract binding to access the raw methods on
}

// LAbstractTokenTRC21TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LAbstractTokenTRC21TransactorRaw struct {
	Contract *LAbstractTokenTRC21Transactor // Generic write-only contract binding to access the raw methods on
}

// NewLAbstractTokenTRC21 creates a new instance of LAbstractTokenTRC21, bound to a specific deployed contract.
func NewLAbstractTokenTRC21(address common.Address, backend bind.ContractBackend) (*LAbstractTokenTRC21, error) {
	contract, err := bindLAbstractTokenTRC21(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LAbstractTokenTRC21{LAbstractTokenTRC21Caller: LAbstractTokenTRC21Caller{contract: contract}, LAbstractTokenTRC21Transactor: LAbstractTokenTRC21Transactor{contract: contract}, LAbstractTokenTRC21Filterer: LAbstractTokenTRC21Filterer{contract: contract}}, nil
}

// NewLAbstractTokenTRC21Caller creates a new read-only instance of LAbstractTokenTRC21, bound to a specific deployed contract.
func NewLAbstractTokenTRC21Caller(address common.Address, caller bind.ContractCaller) (*LAbstractTokenTRC21Caller, error) {
	contract, err := bindLAbstractTokenTRC21(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractTokenTRC21Caller{contract: contract}, nil
}

// NewLAbstractTokenTRC21Transactor creates a new write-only instance of LAbstractTokenTRC21, bound to a specific deployed contract.
func NewLAbstractTokenTRC21Transactor(address common.Address, transactor bind.ContractTransactor) (*LAbstractTokenTRC21Transactor, error) {
	contract, err := bindLAbstractTokenTRC21(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LAbstractTokenTRC21Transactor{contract: contract}, nil
}

// NewLAbstractTokenTRC21Filterer creates a new log filterer instance of LAbstractTokenTRC21, bound to a specific deployed contract.
func NewLAbstractTokenTRC21Filterer(address common.Address, filterer bind.ContractFilterer) (*LAbstractTokenTRC21Filterer, error) {
	contract, err := bindLAbstractTokenTRC21(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LAbstractTokenTRC21Filterer{contract: contract}, nil
}

// bindLAbstractTokenTRC21 binds a generic wrapper to an already deployed contract.
func bindLAbstractTokenTRC21(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LAbstractTokenTRC21ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractTokenTRC21.Contract.LAbstractTokenTRC21Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractTokenTRC21.Contract.LAbstractTokenTRC21Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractTokenTRC21.Contract.LAbstractTokenTRC21Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LAbstractTokenTRC21.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LAbstractTokenTRC21.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LAbstractTokenTRC21.Contract.contract.Transact(opts, method, params...)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21Caller) Issuer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LAbstractTokenTRC21.contract.Call(opts, out, "issuer")
	return *ret0, err
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21Session) Issuer() (common.Address, error) {
	return _LAbstractTokenTRC21.Contract.Issuer(&_LAbstractTokenTRC21.CallOpts)
}

// Issuer is a free data retrieval call binding the contract method 0x1d143848.
//
// Solidity: function issuer() constant returns(address)
func (_LAbstractTokenTRC21 *LAbstractTokenTRC21CallerSession) Issuer() (common.Address, error) {
	return _LAbstractTokenTRC21.Contract.Issuer(&_LAbstractTokenTRC21.CallOpts)
}

// LendingABI is the input ABI used to generate the binding from.
const LendingABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"depositRate\",\"type\":\"uint256\"},{\"name\":\"liquidationRate\",\"type\":\"uint256\"},{\"name\":\"recallRate\",\"type\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"addCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"COLLATERALS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"term\",\"type\":\"uint256\"}],\"name\":\"addTerm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"LENDINGRELAYER_LIST\",\"outputs\":[{\"name\":\"_tradeFee\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"Relayer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TomoXListing\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"baseTokens\",\"type\":\"address[]\"},{\"name\":\"terms\",\"type\":\"uint256[]\"},{\"name\":\"collaterals\",\"type\":\"address[]\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"depositRate\",\"type\":\"uint256\"},{\"name\":\"liquidationRate\",\"type\":\"uint256\"},{\"name\":\"recallRate\",\"type\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"addILOCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"TERMS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"BASES\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"setCollateralPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"COLLATERAL_LIST\",\"outputs\":[{\"name\":\"_depositRate\",\"type\":\"uint256\"},{\"name\":\"_liquidationRate\",\"type\":\"uint256\"},{\"name\":\"_recallRate\",\"type\":\"uint256\"},{\"name\":\"_price\",\"type\":\"uint256\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"addBaseToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ILO_COLLATERALS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_OWNER\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"getLendingRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"r\",\"type\":\"address\"},{\"name\":\"t\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// LendingBin is the compiled bytecode used for deploying new contracts.
const LendingBin = `0x608060405234801561001057600080fd5b50604051604080611ea183398101604052805160209091015160068054600160a060020a03938416600160a060020a03199182161790915560088054939092169281169290921790556007805490911633179055611e2e806100736000396000f3006080604052600436106100cb5763ffffffff60e060020a60003504166271127581146100d05780630811f05a146100ff5780630c655955146101335780630faf292c1461014b578063264949d81461018357806329a4ddec146101985780632ddada4c146101ad578063502e28f41461028a57806356327f57146102b75780636d1dc42a146102e1578063757ff0e3146102f9578063822507011461031d57806383e280d914610369578063b8687ec41461038a578063fd301c49146103a2578063fe824700146103b7575b600080fd5b3480156100dc57600080fd5b506100fd600160a060020a03600435166024356044356064356084356104c5565b005b34801561010b57600080fd5b5061011760043561086e565b60408051600160a060020a039092168252519081900360200190f35b34801561013f57600080fd5b506100fd600435610896565b34801561015757600080fd5b5061016c600160a060020a03600435166109ea565b6040805161ffff9092168252519081900360200190f35b34801561018f57600080fd5b50610117610a00565b3480156101a457600080fd5b50610117610a0f565b3480156101b957600080fd5b5060408051602060046044358181013583810280860185019096528085526100fd958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610a1e9650505050505050565b34801561029657600080fd5b506100fd600160a060020a036004351660243560443560643560843561115d565b3480156102c357600080fd5b506102cf600435611561565b60408051918252519081900360200190f35b3480156102ed57600080fd5b50610117600435611580565b34801561030557600080fd5b506100fd600160a060020a036004351660243561158e565b34801561032957600080fd5b5061033e600160a060020a0360043516611898565b6040805195865260208601949094528484019290925260608401526080830152519081900360a00190f35b34801561037557600080fd5b506100fd600160a060020a03600435166118c6565b34801561039657600080fd5b50610117600435611add565b3480156103ae57600080fd5b50610117611aeb565b3480156103c357600080fd5b506103d8600160a060020a0360043516611afa565b604051808561ffff1661ffff168152602001806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b8381101561042e578181015183820152602001610416565b50505050905001848103835286818151815260200191508051906020019060200280838360005b8381101561046d578181015183820152602001610455565b50505050905001848103825285818151815260200191508051906020019060200280838360005b838110156104ac578181015183820152602001610494565b5050505090500197505050505050505060405180910390f35b6007546000908190600160a060020a0316331461052c576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b6064861015801561053d5750606485115b1515610593576040805160e560020a62461bcd02815260206004820152600d60248201527f496e76616c696420726174657300000000000000000000000000000000000000604482015290519081900360640190fd5b8486116105ea576040805160e560020a62461bcd02815260206004820152601560248201527f496e76616c6964206465706f7369742072617465730000000000000000000000604482015290519081900360640190fd5b858411610641576040805160e560020a62461bcd02815260206004820152601460248201527f496e76616c696420726563616c6c207261746573000000000000000000000000604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a038a811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b15801561069357600080fd5b505af11580156106a7573d6000803e3d6000fd5b505050506040513d60208110156106bd57600080fd5b5051806106d35750600160a060020a0387166001145b915081151561071a576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611de3833981519152604482015290519081900360640190fd5b50600160a060020a03861660009081526001602052604090206004015482156107405750435b6040805160a08101825287815260208082018881528284018881526060840188815260808501878152600160a060020a038e16600090815260018087529088902096518755935193860193909355905160028086019190915590516003850155905160049093019290925581548351818302810183019094528084526108069392918301828280156107fb57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116107dd575b505050505088611c43565b151561086557600280546001810182556000919091527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace01805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0389161790555b50505050505050565b600280548290811061087c57fe5b600091825260209091200154600160a060020a0316905081565b600754600160a060020a031633146108f8576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b603c811015610951576040805160e560020a62461bcd02815260206004820152600c60248201527f496e76616c6964207465726d0000000000000000000000000000000000000000604482015290519081900360640190fd5b6109ab60048054806020026020016040519081016040528092919081815260200182805480156109a057602002820191906000526020600020905b81548152602001906001019080831161098c575b505050505082611c9e565b15156109e757600480546001810182556000919091527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b018190555b50565b60006020819052908152604090205461ffff1681565b600654600160a060020a031681565b600854600160a060020a031681565b600654604080517f540105c7000000000000000000000000000000000000000000000000000000008152600160a060020a03888116600483015291516000938493849391169163540105c791602480820192869290919082900301818387803b158015610a8a57600080fd5b505af1158015610a9e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405260c0811015610ac757600080fd5b815160208301516040840151606085015160808601805194969395929491939283019291640100000000811115610afd57600080fd5b82016020810184811115610b1057600080fd5b8151856020820283011164010000000082111715610b2d57600080fd5b50509291906020018051640100000000811115610b4957600080fd5b82016020810184811115610b5c57600080fd5b8151856020820283011164010000000082111715610b7957600080fd5b50979b505050600160a060020a038a1633149650610be895505050505050576040805160e560020a62461bcd02815260206004820152601660248201527f52656c61796572206f776e657220726571756972656400000000000000000000604482015290519081900360640190fd5b600654604080517f500f99f7000000000000000000000000000000000000000000000000000000008152600160a060020a038b811660048301529151919092169163500f99f79160248083019260209291908290030181600087803b158015610c5057600080fd5b505af1158015610c64573d6000803e3d6000fd5b505050506040513d6020811015610c7a57600080fd5b505115610cd1576040805160e560020a62461bcd02815260206004820152601960248201527f52656c6179657220726571756972656420746f20636c6f736500000000000000604482015290519081900360640190fd5b60008761ffff1610158015610ceb57506103e88761ffff16105b1515610d41576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c696420747261646520466565000000000000000000000000000000604482015290519081900360640190fd5b8451865114610d9a576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f66207465726d7300000000000000604482015290519081900360640190fd5b8351865114610df3576040805160e560020a62461bcd02815260206004820152601f60248201527f4e6f742076616c6964206e756d626572206f6620636f6c6c61746572616c7300604482015290519081900360640190fd5b5060009050805b8551811015610ee257610e7e6003805480602002602001604051908101604052809291908181526020018280548015610e5c57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610e3e575b50505050508783815181101515610e6f57fe5b90602001906020020151611c43565b9150600182151514610eda576040805160e560020a62461bcd02815260206004820152601560248201527f496e76616c6964206c656e64696e6720746f6b656e0000000000000000000000604482015290519081900360640190fd5b600101610dfa565b5060005b8451811015610fc457610f606004805480602002602001604051908101604052809291908181526020018280548015610f3e57602002820191906000526020600020905b815481526020019060010190808311610f2a575b50505050508683815181101515610f5157fe5b90602001906020020151611c9e565b9150600182151514610fbc576040805160e560020a62461bcd02815260206004820152600c60248201527f496e76616c6964207465726d0000000000000000000000000000000000000000604482015290519081900360640190fd5b600101610ee6565b5060005b83518110156110b2578351600090859083908110610fe257fe5b60209081029091010151600160a060020a0316146110aa57611066600580548060200260200160405190810160405280929190818152602001828054801561105357602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611035575b50505050508583815181101515610e6f57fe5b15156110aa576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611de3833981519152604482015290519081900360640190fd5b600101610fc8565b6040805160808101825261ffff898116825260208083018a81528385018a905260608401899052600160a060020a038d166000908152808352949094208351815461ffff1916931692909217825592518051929391926111189260018501920190611cdb565b5060408201518051611134916002840191602090910190611d4d565b5060608201518051611150916003840191602090910190611cdb565b5050505050505050505050565b6000806000606487101580156111735750606486115b15156111c9576040805160e560020a62461bcd02815260206004820152600d60248201527f496e76616c696420726174657300000000000000000000000000000000000000604482015290519081900360640190fd5b858711611220576040805160e560020a62461bcd02815260206004820152601560248201527f496e76616c6964206465706f7369742072617465730000000000000000000000604482015290519081900360640190fd5b868511611277576040805160e560020a62461bcd02815260206004820152601460248201527f496e76616c696420726563616c6c207261746573000000000000000000000000604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a038b811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b1580156112c957600080fd5b505af11580156112dd573d6000803e3d6000fd5b505050506040513d60208110156112f357600080fd5b5051925082151561133c576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611de3833981519152604482015290519081900360640190fd5b87915033600160a060020a031682600160a060020a0316631d1438486040518163ffffffff1660e060020a028152600401602060405180830381600087803b15801561138757600080fd5b505af115801561139b573d6000803e3d6000fd5b505050506040513d60208110156113b157600080fd5b5051600160a060020a031614611411576040805160e560020a62461bcd02815260206004820152601560248201527f526571756972656420746f6b656e206973737565720000000000000000000000604482015290519081900360640190fd5b50600160a060020a03871660009081526001602052604090206004015483156114375750435b6040805160a08101825288815260208082018981528284018981526060840189815260808501878152600160a060020a038f166000908152600180875290889020965187559351938601939093559051600285015551600384015551600490920191909155600580548351818402810184019094528084526114f893928301828280156114ed57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116114cf575b505050505089611c43565b151561155757600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a038a161790555b5050505050505050565b600480548290811061156f57fe5b600091825260209091200154905081565b600380548290811061087c57fe5b6008546040805160e060020a63a3ff31b5028152600160a060020a03858116600483015291516000938493169163a3ff31b591602480830192602092919082900301818787803b1580156115e157600080fd5b505af11580156115f5573d6000803e3d6000fd5b505050506040513d602081101561160b57600080fd5b5051806116215750600160a060020a0384166001145b9150811515611668576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611de3833981519152604482015290519081900360640190fd5b600160a060020a038416600090815260016020526040902054606411156116c7576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611de3833981519152604482015290519081900360640190fd5b61172b600280548060200260200160405190810160405280929190818152602001828054801561172057602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611702575b505050505085611c43565b1561179757600754600160a060020a03163314611792576040805160e560020a62461bcd02815260206004820152601760248201527f436f6e7472616374206f776e6572207265717569726564000000000000000000604482015290519081900360640190fd5b61186c565b83905033600160a060020a031681600160a060020a0316631d1438486040518163ffffffff1660e060020a028152600401602060405180830381600087803b1580156117e257600080fd5b505af11580156117f6573d6000803e3d6000fd5b505050506040513d602081101561180c57600080fd5b5051600160a060020a03161461186c576040805160e560020a62461bcd02815260206004820152601560248201527f526571756972656420746f6b656e206973737565720000000000000000000000604482015290519081900360640190fd5b5050600160a060020a039091166000908152600160205260409020600381019190915543600490910155565b6001602081905260009182526040909120805491810154600282015460038301546004909301549192909185565b600754600090600160a060020a0316331461192b576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a0385811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b15801561197d57600080fd5b505af1158015611991573d6000803e3d6000fd5b505050506040513d60208110156119a757600080fd5b5051806119bd5750600160a060020a0382166001145b9050801515611a16576040805160e560020a62461bcd02815260206004820152601260248201527f496e76616c6964206261736520746f6b656e0000000000000000000000000000604482015290519081900360640190fd5b611a7a6003805480602002602001604051908101604052809291908181526020018280548015611a6f57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611a51575b505050505083611c43565b1515611ad957600380546001810182556000919091527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b01805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0384161790555b5050565b600580548290811061087c57fe5b600754600160a060020a031681565b600160a060020a03811660009081526020818152604080832080546001820180548451818702810187019095528085526060958695869561ffff909516946002810193600390910192859190830182828015611b7f57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611b61575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015611bd157602002820191906000526020600020905b815481526020019060010190808311611bbd575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015611c2d57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611c0f575b5050505050905093509350935093509193509193565b6000805b8351811015611c925782600160a060020a03168482815181101515611c6857fe5b90602001906020020151600160a060020a03161415611c8a5760019150611c97565b600101611c47565b600091505b5092915050565b6000805b8351811015611c9257828482815181101515611cba57fe5b906020019060200201511415611cd35760019150611c97565b600101611ca2565b828054828255906000526020600020908101928215611d3d579160200282015b82811115611d3d578251825473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909116178255602090920191600190910190611cfb565b50611d49929150611d94565b5090565b828054828255906000526020600020908101928215611d88579160200282015b82811115611d88578251825591602001919060010190611d6d565b50611d49929150611dc8565b611dc591905b80821115611d4957805473ffffffffffffffffffffffffffffffffffffffff19168155600101611d9a565b90565b611dc591905b80821115611d495760008155600101611dce5600496e76616c696420636f6c6c61746572616c0000000000000000000000000000a165627a7a723058206da634fec04d831052ef479f43cade9482ccaa6b2c8c5f04d4850d54a353d2750029`

// DeployLending deploys a new Ethereum contract, binding an instance of Lending to it.
func DeployLending(auth *bind.TransactOpts, backend bind.ContractBackend, r common.Address, t common.Address) (common.Address, *types.Transaction, *Lending, error) {
	parsed, err := abi.JSON(strings.NewReader(LendingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LendingBin), backend, r, t)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Lending{LendingCaller: LendingCaller{contract: contract}, LendingTransactor: LendingTransactor{contract: contract}, LendingFilterer: LendingFilterer{contract: contract}}, nil
}

// Lending is an auto generated Go binding around an Ethereum contract.
type Lending struct {
	LendingCaller     // Read-only binding to the contract
	LendingTransactor // Write-only binding to the contract
	LendingFilterer   // Log filterer for contract events
}

// LendingCaller is an auto generated read-only Go binding around an Ethereum contract.
type LendingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LendingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LendingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LendingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LendingSession struct {
	Contract     *Lending          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LendingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LendingCallerSession struct {
	Contract *LendingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// LendingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LendingTransactorSession struct {
	Contract     *LendingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LendingRaw is an auto generated low-level Go binding around an Ethereum contract.
type LendingRaw struct {
	Contract *Lending // Generic contract binding to access the raw methods on
}

// LendingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LendingCallerRaw struct {
	Contract *LendingCaller // Generic read-only contract binding to access the raw methods on
}

// LendingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LendingTransactorRaw struct {
	Contract *LendingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLending creates a new instance of Lending, bound to a specific deployed contract.
func NewLending(address common.Address, backend bind.ContractBackend) (*Lending, error) {
	contract, err := bindLending(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Lending{LendingCaller: LendingCaller{contract: contract}, LendingTransactor: LendingTransactor{contract: contract}, LendingFilterer: LendingFilterer{contract: contract}}, nil
}

// NewLendingCaller creates a new read-only instance of Lending, bound to a specific deployed contract.
func NewLendingCaller(address common.Address, caller bind.ContractCaller) (*LendingCaller, error) {
	contract, err := bindLending(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LendingCaller{contract: contract}, nil
}

// NewLendingTransactor creates a new write-only instance of Lending, bound to a specific deployed contract.
func NewLendingTransactor(address common.Address, transactor bind.ContractTransactor) (*LendingTransactor, error) {
	contract, err := bindLending(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LendingTransactor{contract: contract}, nil
}

// NewLendingFilterer creates a new log filterer instance of Lending, bound to a specific deployed contract.
func NewLendingFilterer(address common.Address, filterer bind.ContractFilterer) (*LendingFilterer, error) {
	contract, err := bindLending(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LendingFilterer{contract: contract}, nil
}

// bindLending binds a generic wrapper to an already deployed contract.
func bindLending(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LendingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Lending *LendingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Lending.Contract.LendingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Lending *LendingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Lending.Contract.LendingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Lending *LendingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Lending.Contract.LendingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Lending *LendingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Lending.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Lending *LendingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Lending.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Lending *LendingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Lending.Contract.contract.Transact(opts, method, params...)
}

// BASES is a free data retrieval call binding the contract method 0x6d1dc42a.
//
// Solidity: function BASES( uint256) constant returns(address)
func (_Lending *LendingCaller) BASES(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "BASES", arg0)
	return *ret0, err
}

// BASES is a free data retrieval call binding the contract method 0x6d1dc42a.
//
// Solidity: function BASES( uint256) constant returns(address)
func (_Lending *LendingSession) BASES(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.BASES(&_Lending.CallOpts, arg0)
}

// BASES is a free data retrieval call binding the contract method 0x6d1dc42a.
//
// Solidity: function BASES( uint256) constant returns(address)
func (_Lending *LendingCallerSession) BASES(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.BASES(&_Lending.CallOpts, arg0)
}

// COLLATERALS is a free data retrieval call binding the contract method 0x0811f05a.
//
// Solidity: function COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCaller) COLLATERALS(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "COLLATERALS", arg0)
	return *ret0, err
}

// COLLATERALS is a free data retrieval call binding the contract method 0x0811f05a.
//
// Solidity: function COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingSession) COLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.COLLATERALS(&_Lending.CallOpts, arg0)
}

// COLLATERALS is a free data retrieval call binding the contract method 0x0811f05a.
//
// Solidity: function COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCallerSession) COLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.COLLATERALS(&_Lending.CallOpts, arg0)
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _recallRate uint256, _price uint256, _blockNumber uint256)
func (_Lending *LendingCaller) COLLATERALLIST(opts *bind.CallOpts, arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	RecallRate      *big.Int
	Price           *big.Int
	BlockNumber     *big.Int
}, error) {
	ret := new(struct {
		DepositRate     *big.Int
		LiquidationRate *big.Int
		RecallRate      *big.Int
		Price           *big.Int
		BlockNumber     *big.Int
	})
	out := ret
	err := _Lending.contract.Call(opts, out, "COLLATERAL_LIST", arg0)
	return *ret, err
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _recallRate uint256, _price uint256, _blockNumber uint256)
func (_Lending *LendingSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	RecallRate      *big.Int
	Price           *big.Int
	BlockNumber     *big.Int
}, error) {
	return _Lending.Contract.COLLATERALLIST(&_Lending.CallOpts, arg0)
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _recallRate uint256, _price uint256, _blockNumber uint256)
func (_Lending *LendingCallerSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	RecallRate      *big.Int
	Price           *big.Int
	BlockNumber     *big.Int
}, error) {
	return _Lending.Contract.COLLATERALLIST(&_Lending.CallOpts, arg0)
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_Lending *LendingCaller) CONTRACTOWNER(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "CONTRACT_OWNER")
	return *ret0, err
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_Lending *LendingSession) CONTRACTOWNER() (common.Address, error) {
	return _Lending.Contract.CONTRACTOWNER(&_Lending.CallOpts)
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_Lending *LendingCallerSession) CONTRACTOWNER() (common.Address, error) {
	return _Lending.Contract.CONTRACTOWNER(&_Lending.CallOpts)
}

// ILOCOLLATERALS is a free data retrieval call binding the contract method 0xb8687ec4.
//
// Solidity: function ILO_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCaller) ILOCOLLATERALS(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "ILO_COLLATERALS", arg0)
	return *ret0, err
}

// ILOCOLLATERALS is a free data retrieval call binding the contract method 0xb8687ec4.
//
// Solidity: function ILO_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingSession) ILOCOLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.ILOCOLLATERALS(&_Lending.CallOpts, arg0)
}

// ILOCOLLATERALS is a free data retrieval call binding the contract method 0xb8687ec4.
//
// Solidity: function ILO_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCallerSession) ILOCOLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.ILOCOLLATERALS(&_Lending.CallOpts, arg0)
}

// LENDINGRELAYERLIST is a free data retrieval call binding the contract method 0x0faf292c.
//
// Solidity: function LENDINGRELAYER_LIST( address) constant returns(_tradeFee uint16)
func (_Lending *LendingCaller) LENDINGRELAYERLIST(opts *bind.CallOpts, arg0 common.Address) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "LENDINGRELAYER_LIST", arg0)
	return *ret0, err
}

// LENDINGRELAYERLIST is a free data retrieval call binding the contract method 0x0faf292c.
//
// Solidity: function LENDINGRELAYER_LIST( address) constant returns(_tradeFee uint16)
func (_Lending *LendingSession) LENDINGRELAYERLIST(arg0 common.Address) (uint16, error) {
	return _Lending.Contract.LENDINGRELAYERLIST(&_Lending.CallOpts, arg0)
}

// LENDINGRELAYERLIST is a free data retrieval call binding the contract method 0x0faf292c.
//
// Solidity: function LENDINGRELAYER_LIST( address) constant returns(_tradeFee uint16)
func (_Lending *LendingCallerSession) LENDINGRELAYERLIST(arg0 common.Address) (uint16, error) {
	return _Lending.Contract.LENDINGRELAYERLIST(&_Lending.CallOpts, arg0)
}

// Relayer is a free data retrieval call binding the contract method 0x264949d8.
//
// Solidity: function Relayer() constant returns(address)
func (_Lending *LendingCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "Relayer")
	return *ret0, err
}

// Relayer is a free data retrieval call binding the contract method 0x264949d8.
//
// Solidity: function Relayer() constant returns(address)
func (_Lending *LendingSession) Relayer() (common.Address, error) {
	return _Lending.Contract.Relayer(&_Lending.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x264949d8.
//
// Solidity: function Relayer() constant returns(address)
func (_Lending *LendingCallerSession) Relayer() (common.Address, error) {
	return _Lending.Contract.Relayer(&_Lending.CallOpts)
}

// TERMS is a free data retrieval call binding the contract method 0x56327f57.
//
// Solidity: function TERMS( uint256) constant returns(uint256)
func (_Lending *LendingCaller) TERMS(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "TERMS", arg0)
	return *ret0, err
}

// TERMS is a free data retrieval call binding the contract method 0x56327f57.
//
// Solidity: function TERMS( uint256) constant returns(uint256)
func (_Lending *LendingSession) TERMS(arg0 *big.Int) (*big.Int, error) {
	return _Lending.Contract.TERMS(&_Lending.CallOpts, arg0)
}

// TERMS is a free data retrieval call binding the contract method 0x56327f57.
//
// Solidity: function TERMS( uint256) constant returns(uint256)
func (_Lending *LendingCallerSession) TERMS(arg0 *big.Int) (*big.Int, error) {
	return _Lending.Contract.TERMS(&_Lending.CallOpts, arg0)
}

// TomoXListing is a free data retrieval call binding the contract method 0x29a4ddec.
//
// Solidity: function TomoXListing() constant returns(address)
func (_Lending *LendingCaller) TomoXListing(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "TomoXListing")
	return *ret0, err
}

// TomoXListing is a free data retrieval call binding the contract method 0x29a4ddec.
//
// Solidity: function TomoXListing() constant returns(address)
func (_Lending *LendingSession) TomoXListing() (common.Address, error) {
	return _Lending.Contract.TomoXListing(&_Lending.CallOpts)
}

// TomoXListing is a free data retrieval call binding the contract method 0x29a4ddec.
//
// Solidity: function TomoXListing() constant returns(address)
func (_Lending *LendingCallerSession) TomoXListing() (common.Address, error) {
	return _Lending.Contract.TomoXListing(&_Lending.CallOpts)
}

// GetLendingRelayerByCoinbase is a free data retrieval call binding the contract method 0xfe824700.
//
// Solidity: function getLendingRelayerByCoinbase(coinbase address) constant returns(uint16, address[], uint256[], address[])
func (_Lending *LendingCaller) GetLendingRelayerByCoinbase(opts *bind.CallOpts, coinbase common.Address) (uint16, []common.Address, []*big.Int, []common.Address, error) {
	var (
		ret0 = new(uint16)
		ret1 = new([]common.Address)
		ret2 = new([]*big.Int)
		ret3 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _Lending.contract.Call(opts, out, "getLendingRelayerByCoinbase", coinbase)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetLendingRelayerByCoinbase is a free data retrieval call binding the contract method 0xfe824700.
//
// Solidity: function getLendingRelayerByCoinbase(coinbase address) constant returns(uint16, address[], uint256[], address[])
func (_Lending *LendingSession) GetLendingRelayerByCoinbase(coinbase common.Address) (uint16, []common.Address, []*big.Int, []common.Address, error) {
	return _Lending.Contract.GetLendingRelayerByCoinbase(&_Lending.CallOpts, coinbase)
}

// GetLendingRelayerByCoinbase is a free data retrieval call binding the contract method 0xfe824700.
//
// Solidity: function getLendingRelayerByCoinbase(coinbase address) constant returns(uint16, address[], uint256[], address[])
func (_Lending *LendingCallerSession) GetLendingRelayerByCoinbase(coinbase common.Address) (uint16, []common.Address, []*big.Int, []common.Address, error) {
	return _Lending.Contract.GetLendingRelayerByCoinbase(&_Lending.CallOpts, coinbase)
}

// AddBaseToken is a paid mutator transaction binding the contract method 0x83e280d9.
//
// Solidity: function addBaseToken(token address) returns()
func (_Lending *LendingTransactor) AddBaseToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addBaseToken", token)
}

// AddBaseToken is a paid mutator transaction binding the contract method 0x83e280d9.
//
// Solidity: function addBaseToken(token address) returns()
func (_Lending *LendingSession) AddBaseToken(token common.Address) (*types.Transaction, error) {
	return _Lending.Contract.AddBaseToken(&_Lending.TransactOpts, token)
}

// AddBaseToken is a paid mutator transaction binding the contract method 0x83e280d9.
//
// Solidity: function addBaseToken(token address) returns()
func (_Lending *LendingTransactorSession) AddBaseToken(token common.Address) (*types.Transaction, error) {
	return _Lending.Contract.AddBaseToken(&_Lending.TransactOpts, token)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x00711275.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingTransactor) AddCollateral(opts *bind.TransactOpts, token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addCollateral", token, depositRate, liquidationRate, recallRate, price)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x00711275.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingSession) AddCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, recallRate, price)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x00711275.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingTransactorSession) AddCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, recallRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x502e28f4.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingTransactor) AddILOCollateral(opts *bind.TransactOpts, token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addILOCollateral", token, depositRate, liquidationRate, recallRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x502e28f4.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingSession) AddILOCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddILOCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, recallRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x502e28f4.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, recallRate uint256, price uint256) returns()
func (_Lending *LendingTransactorSession) AddILOCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, recallRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddILOCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, recallRate, price)
}

// AddTerm is a paid mutator transaction binding the contract method 0x0c655955.
//
// Solidity: function addTerm(term uint256) returns()
func (_Lending *LendingTransactor) AddTerm(opts *bind.TransactOpts, term *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addTerm", term)
}

// AddTerm is a paid mutator transaction binding the contract method 0x0c655955.
//
// Solidity: function addTerm(term uint256) returns()
func (_Lending *LendingSession) AddTerm(term *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddTerm(&_Lending.TransactOpts, term)
}

// AddTerm is a paid mutator transaction binding the contract method 0x0c655955.
//
// Solidity: function addTerm(term uint256) returns()
func (_Lending *LendingTransactorSession) AddTerm(term *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddTerm(&_Lending.TransactOpts, term)
}

// SetCollateralPrice is a paid mutator transaction binding the contract method 0x757ff0e3.
//
// Solidity: function setCollateralPrice(token address, price uint256) returns()
func (_Lending *LendingTransactor) SetCollateralPrice(opts *bind.TransactOpts, token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "setCollateralPrice", token, price)
}

// SetCollateralPrice is a paid mutator transaction binding the contract method 0x757ff0e3.
//
// Solidity: function setCollateralPrice(token address, price uint256) returns()
func (_Lending *LendingSession) SetCollateralPrice(token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.SetCollateralPrice(&_Lending.TransactOpts, token, price)
}

// SetCollateralPrice is a paid mutator transaction binding the contract method 0x757ff0e3.
//
// Solidity: function setCollateralPrice(token address, price uint256) returns()
func (_Lending *LendingTransactorSession) SetCollateralPrice(token common.Address, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.SetCollateralPrice(&_Lending.TransactOpts, token, price)
}

// Update is a paid mutator transaction binding the contract method 0x2ddada4c.
//
// Solidity: function update(coinbase address, tradeFee uint16, baseTokens address[], terms uint256[], collaterals address[]) returns()
func (_Lending *LendingTransactor) Update(opts *bind.TransactOpts, coinbase common.Address, tradeFee uint16, baseTokens []common.Address, terms []*big.Int, collaterals []common.Address) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "update", coinbase, tradeFee, baseTokens, terms, collaterals)
}

// Update is a paid mutator transaction binding the contract method 0x2ddada4c.
//
// Solidity: function update(coinbase address, tradeFee uint16, baseTokens address[], terms uint256[], collaterals address[]) returns()
func (_Lending *LendingSession) Update(coinbase common.Address, tradeFee uint16, baseTokens []common.Address, terms []*big.Int, collaterals []common.Address) (*types.Transaction, error) {
	return _Lending.Contract.Update(&_Lending.TransactOpts, coinbase, tradeFee, baseTokens, terms, collaterals)
}

// Update is a paid mutator transaction binding the contract method 0x2ddada4c.
//
// Solidity: function update(coinbase address, tradeFee uint16, baseTokens address[], terms uint256[], collaterals address[]) returns()
func (_Lending *LendingTransactorSession) Update(coinbase common.Address, tradeFee uint16, baseTokens []common.Address, terms []*big.Int, collaterals []common.Address) (*types.Transaction, error) {
	return _Lending.Contract.Update(&_Lending.TransactOpts, coinbase, tradeFee, baseTokens, terms, collaterals)
}
