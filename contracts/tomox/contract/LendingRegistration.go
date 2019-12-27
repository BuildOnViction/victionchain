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

// AbstractRegistrationABI is the input ABI used to generate the binding from.
const AbstractRegistrationABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"RESIGN_REQUESTS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"getRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AbstractRegistrationBin is the compiled bytecode used for deploying new contracts.
const AbstractRegistrationBin = `0x`

// DeployAbstractRegistration deploys a new Ethereum contract, binding an instance of AbstractRegistration to it.
func DeployAbstractRegistration(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AbstractRegistration, error) {
	parsed, err := abi.JSON(strings.NewReader(AbstractRegistrationABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AbstractRegistrationBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AbstractRegistration{AbstractRegistrationCaller: AbstractRegistrationCaller{contract: contract}, AbstractRegistrationTransactor: AbstractRegistrationTransactor{contract: contract}, AbstractRegistrationFilterer: AbstractRegistrationFilterer{contract: contract}}, nil
}

// AbstractRegistration is an auto generated Go binding around an Ethereum contract.
type AbstractRegistration struct {
	AbstractRegistrationCaller     // Read-only binding to the contract
	AbstractRegistrationTransactor // Write-only binding to the contract
	AbstractRegistrationFilterer   // Log filterer for contract events
}

// AbstractRegistrationCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbstractRegistrationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbstractRegistrationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbstractRegistrationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbstractRegistrationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbstractRegistrationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbstractRegistrationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbstractRegistrationSession struct {
	Contract     *AbstractRegistration // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AbstractRegistrationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbstractRegistrationCallerSession struct {
	Contract *AbstractRegistrationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// AbstractRegistrationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbstractRegistrationTransactorSession struct {
	Contract     *AbstractRegistrationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// AbstractRegistrationRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbstractRegistrationRaw struct {
	Contract *AbstractRegistration // Generic contract binding to access the raw methods on
}

// AbstractRegistrationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbstractRegistrationCallerRaw struct {
	Contract *AbstractRegistrationCaller // Generic read-only contract binding to access the raw methods on
}

// AbstractRegistrationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbstractRegistrationTransactorRaw struct {
	Contract *AbstractRegistrationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbstractRegistration creates a new instance of AbstractRegistration, bound to a specific deployed contract.
func NewAbstractRegistration(address common.Address, backend bind.ContractBackend) (*AbstractRegistration, error) {
	contract, err := bindAbstractRegistration(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AbstractRegistration{AbstractRegistrationCaller: AbstractRegistrationCaller{contract: contract}, AbstractRegistrationTransactor: AbstractRegistrationTransactor{contract: contract}, AbstractRegistrationFilterer: AbstractRegistrationFilterer{contract: contract}}, nil
}

// NewAbstractRegistrationCaller creates a new read-only instance of AbstractRegistration, bound to a specific deployed contract.
func NewAbstractRegistrationCaller(address common.Address, caller bind.ContractCaller) (*AbstractRegistrationCaller, error) {
	contract, err := bindAbstractRegistration(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbstractRegistrationCaller{contract: contract}, nil
}

// NewAbstractRegistrationTransactor creates a new write-only instance of AbstractRegistration, bound to a specific deployed contract.
func NewAbstractRegistrationTransactor(address common.Address, transactor bind.ContractTransactor) (*AbstractRegistrationTransactor, error) {
	contract, err := bindAbstractRegistration(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbstractRegistrationTransactor{contract: contract}, nil
}

// NewAbstractRegistrationFilterer creates a new log filterer instance of AbstractRegistration, bound to a specific deployed contract.
func NewAbstractRegistrationFilterer(address common.Address, filterer bind.ContractFilterer) (*AbstractRegistrationFilterer, error) {
	contract, err := bindAbstractRegistration(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbstractRegistrationFilterer{contract: contract}, nil
}

// bindAbstractRegistration binds a generic wrapper to an already deployed contract.
func bindAbstractRegistration(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AbstractRegistrationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AbstractRegistration *AbstractRegistrationRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AbstractRegistration.Contract.AbstractRegistrationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AbstractRegistration *AbstractRegistrationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbstractRegistration.Contract.AbstractRegistrationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AbstractRegistration *AbstractRegistrationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbstractRegistration.Contract.AbstractRegistrationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AbstractRegistration *AbstractRegistrationCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AbstractRegistration.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AbstractRegistration *AbstractRegistrationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbstractRegistration.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AbstractRegistration *AbstractRegistrationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbstractRegistration.Contract.contract.Transact(opts, method, params...)
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_AbstractRegistration *AbstractRegistrationCaller) RESIGNREQUESTS(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AbstractRegistration.contract.Call(opts, out, "RESIGN_REQUESTS", arg0)
	return *ret0, err
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_AbstractRegistration *AbstractRegistrationSession) RESIGNREQUESTS(arg0 common.Address) (*big.Int, error) {
	return _AbstractRegistration.Contract.RESIGNREQUESTS(&_AbstractRegistration.CallOpts, arg0)
}

// RESIGNREQUESTS is a free data retrieval call binding the contract method 0x500f99f7.
//
// Solidity: function RESIGN_REQUESTS( address) constant returns(uint256)
func (_AbstractRegistration *AbstractRegistrationCallerSession) RESIGNREQUESTS(arg0 common.Address) (*big.Int, error) {
	return _AbstractRegistration.Contract.RESIGNREQUESTS(&_AbstractRegistration.CallOpts, arg0)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_AbstractRegistration *AbstractRegistrationCaller) GetRelayerByCoinbase(opts *bind.CallOpts, arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
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
	err := _AbstractRegistration.contract.Call(opts, out, "getRelayerByCoinbase", arg0)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, err
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_AbstractRegistration *AbstractRegistrationSession) GetRelayerByCoinbase(arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _AbstractRegistration.Contract.GetRelayerByCoinbase(&_AbstractRegistration.CallOpts, arg0)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase( address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_AbstractRegistration *AbstractRegistrationCallerSession) GetRelayerByCoinbase(arg0 common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _AbstractRegistration.Contract.GetRelayerByCoinbase(&_AbstractRegistration.CallOpts, arg0)
}

// LendingABI is the input ABI used to generate the binding from.
const LendingABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"COLLATERALS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"term\",\"type\":\"uint256\"}],\"name\":\"addTerm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"LENDINGRELAYER_LIST\",\"outputs\":[{\"name\":\"_tradeFee\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"baseTokens\",\"type\":\"address[]\"},{\"name\":\"terms\",\"type\":\"uint256[]\"},{\"name\":\"collaterals\",\"type\":\"address[]\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"BASES\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"depositRate\",\"type\":\"uint256\"}],\"name\":\"addCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"COLLATERAL_LIST\",\"outputs\":[{\"name\":\"_depositRate\",\"type\":\"uint256\"},{\"name\":\"_price\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"addBaseToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"relayer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"getLendingRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"r\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// LendingBin is the compiled bytecode used for deploying new contracts.
const LendingBin = `0x608060405234801561001057600080fd5b50604051602080610c9f833981016040525160058054600160a060020a031916600160a060020a03909216919091179055610c4f806100506000396000f3006080604052600436106100a35763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630811f05a81146100a85780630c655955146100dc5780630faf292c146100f65780632ddada4c1461012e5780636d1dc42a1461020b5780636d75b9ee14610223578063822507011461024757806383e280d9146102815780638406c079146102a2578063fe824700146102b7575b600080fd5b3480156100b457600080fd5b506100c06004356103c5565b60408051600160a060020a039092168252519081900360200190f35b3480156100e857600080fd5b506100f46004356103ed565b005b34801561010257600080fd5b50610117600160a060020a0360043516610422565b6040805161ffff9092168252519081900360200190f35b34801561013a57600080fd5b5060408051602060046044358181013583810280860185019096528085526100f4958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506104389650505050505050565b34801561021757600080fd5b506100c06004356108b4565b34801561022f57600080fd5b506100f4600160a060020a03600435166024356108c2565b34801561025357600080fd5b50610268600160a060020a036004351661094c565b6040805192835260208301919091528051918290030190f35b34801561028d57600080fd5b506100f4600160a060020a0360043516610965565b3480156102ae57600080fd5b506100c06109c4565b3480156102c357600080fd5b506102d8600160a060020a03600435166109d3565b604051808561ffff1661ffff168152602001806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b8381101561032e578181015183820152602001610316565b50505050905001848103835286818151815260200191508051906020019060200280838360005b8381101561036d578181015183820152602001610355565b50505050905001848103825285818151815260200191508051906020019060200280838360005b838110156103ac578181015183820152602001610394565b5050505090500197505050505050505060405180910390f35b60028054829081106103d357fe5b600091825260209091200154600160a060020a0316905081565b600480546001810182556000919091527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b0155565b60006020819052908152604090205461ffff1681565b600554604080517f540105c7000000000000000000000000000000000000000000000000000000008152600160a060020a0388811660048301529151600093929092169163540105c791602480820192869290919082900301818387803b1580156104a257600080fd5b505af11580156104b6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405260c08110156104df57600080fd5b81516020830151604084015160608501516080860180519496939592949193928301929164010000000081111561051557600080fd5b8201602081018481111561052857600080fd5b815185602082028301116401000000008211171561054557600080fd5b5050929190602001805164010000000081111561056157600080fd5b8201602081018481111561057457600080fd5b815185602082028301116401000000008211171561059157600080fd5b509799505050600160a060020a0388163314965061060095505050505050576040805160e560020a62461bcd02815260206004820152601660248201527f72656c61796572206f776e657220726571756972656400000000000000000000604482015290519081900360640190fd5b600554604080517f500f99f7000000000000000000000000000000000000000000000000000000008152600160a060020a0389811660048301529151919092169163500f99f79160248083019260209291908290030181600087803b15801561066857600080fd5b505af115801561067c573d6000803e3d6000fd5b505050506040513d602081101561069257600080fd5b5051156106e9576040805160e560020a62461bcd02815260206004820152601960248201527f72656c6179657220726571756972656420746f20636c6f736500000000000000604482015290519081900360640190fd5b60018561ffff161015801561070357506103e88561ffff16105b1515610759576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c6964204d616b657220466565000000000000000000000000000000604482015290519081900360640190fd5b82518451146107b2576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f66207465726d7300000000000000604482015290519081900360640190fd5b815184511461080b576040805160e560020a62461bcd02815260206004820152601f60248201527f4e6f742076616c6964206e756d626572206f6620636f6c6c61746572616c7300604482015290519081900360640190fd5b6040805160808101825261ffff8781168252602080830188815283850188905260608401879052600160a060020a038b166000908152808352949094208351815461ffff1916931692909217825592518051929391926108719260018501920190610b1c565b506040820151805161088d916002840191602090910190610b8e565b50606082015180516108a9916003840191602090910190610b1c565b505050505050505050565b60038054829081106103d357fe5b60408051808201825291825260006020808401828152600160a060020a03909516808352600191829052928220935184559351928401929092556002805493840181559091527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace909101805473ffffffffffffffffffffffffffffffffffffffff19169091179055565b6001602081905260009182526040909120805491015482565b600380546001810182556000919091527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b01805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b600554600160a060020a031681565b600160a060020a03811660009081526020818152604080832080546001820180548451818702810187019095528085526060958695869561ffff909516946002810193600390910192859190830182828015610a5857602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610a3a575b5050505050925081805480602002602001604051908101604052809291908181526020018280548015610aaa57602002820191906000526020600020905b815481526020019060010190808311610a96575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015610b0657602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610ae8575b5050505050905093509350935093509193509193565b828054828255906000526020600020908101928215610b7e579160200282015b82811115610b7e578251825473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909116178255602090920191600190910190610b3c565b50610b8a929150610bd5565b5090565b828054828255906000526020600020908101928215610bc9579160200282015b82811115610bc9578251825591602001919060010190610bae565b50610b8a929150610c09565b610c0691905b80821115610b8a57805473ffffffffffffffffffffffffffffffffffffffff19168155600101610bdb565b90565b610c0691905b80821115610b8a5760008155600101610c0f5600a165627a7a723058206a2d4823358133b266d2d9854a24b5c170e349e6710ce42317d817cc8fc11a9d0029`

// DeployLending deploys a new Ethereum contract, binding an instance of Lending to it.
func DeployLending(auth *bind.TransactOpts, backend bind.ContractBackend, r common.Address) (common.Address, *types.Transaction, *Lending, error) {
	parsed, err := abi.JSON(strings.NewReader(LendingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LendingBin), backend, r)
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
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _price uint256)
func (_Lending *LendingCaller) COLLATERALLIST(opts *bind.CallOpts, arg0 common.Address) (struct {
	DepositRate *big.Int
	Price       *big.Int
}, error) {
	ret := new(struct {
		DepositRate *big.Int
		Price       *big.Int
	})
	out := ret
	err := _Lending.contract.Call(opts, out, "COLLATERAL_LIST", arg0)
	return *ret, err
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _price uint256)
func (_Lending *LendingSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate *big.Int
	Price       *big.Int
}, error) {
	return _Lending.Contract.COLLATERALLIST(&_Lending.CallOpts, arg0)
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _price uint256)
func (_Lending *LendingCallerSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate *big.Int
	Price       *big.Int
}, error) {
	return _Lending.Contract.COLLATERALLIST(&_Lending.CallOpts, arg0)
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

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() constant returns(address)
func (_Lending *LendingCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "relayer")
	return *ret0, err
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() constant returns(address)
func (_Lending *LendingSession) Relayer() (common.Address, error) {
	return _Lending.Contract.Relayer(&_Lending.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() constant returns(address)
func (_Lending *LendingCallerSession) Relayer() (common.Address, error) {
	return _Lending.Contract.Relayer(&_Lending.CallOpts)
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

// AddCollateral is a paid mutator transaction binding the contract method 0x6d75b9ee.
//
// Solidity: function addCollateral(token address, depositRate uint256) returns()
func (_Lending *LendingTransactor) AddCollateral(opts *bind.TransactOpts, token common.Address, depositRate *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addCollateral", token, depositRate)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x6d75b9ee.
//
// Solidity: function addCollateral(token address, depositRate uint256) returns()
func (_Lending *LendingSession) AddCollateral(token common.Address, depositRate *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x6d75b9ee.
//
// Solidity: function addCollateral(token address, depositRate uint256) returns()
func (_Lending *LendingTransactorSession) AddCollateral(token common.Address, depositRate *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate)
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
