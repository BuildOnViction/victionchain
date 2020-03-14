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
const LendingABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"COLLATERALS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"term\",\"type\":\"uint256\"}],\"name\":\"addTerm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"LENDINGRELAYER_LIST\",\"outputs\":[{\"name\":\"_tradeFee\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"baseTokens\",\"type\":\"address[]\"},{\"name\":\"terms\",\"type\":\"uint256[]\"},{\"name\":\"collaterals\",\"type\":\"address[]\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"depositRate\",\"type\":\"uint256\"},{\"name\":\"liquidationRate\",\"type\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"addILOCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"BASES\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"setCollateralPrice\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"COLLATERAL_LIST\",\"outputs\":[{\"name\":\"_depositRate\",\"type\":\"uint256\"},{\"name\":\"_liquidationRate\",\"type\":\"uint256\"},{\"name\":\"_price\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"}],\"name\":\"addBaseToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"relayer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ALL_COLLATERALS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"depositRate\",\"type\":\"uint256\"},{\"name\":\"liquidationRate\",\"type\":\"uint256\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"addCollateral\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_OWNER\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"getLendingRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"r\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// LendingBin is the compiled bytecode used for deploying new contracts.
const LendingBin = `0x608060405234801561001057600080fd5b50604051602080611c92833981016040525160068054600160a060020a03909216600160a060020a03199283161790556007805490911633179055611c388061005a6000396000f3006080604052600436106100b65763ffffffff60e060020a6000350416630811f05a81146100bb5780630c655955146100ef5780630faf292c146101095780632ddada4c146101415780633b8748271461021e5780636d1dc42a14610248578063757ff0e314610260578063822507011461028457806383e280d9146102c35780638406c079146102e4578063c4691d36146102f9578063e5eecf6814610311578063fd301c491461033b578063fe82470014610350575b600080fd5b3480156100c757600080fd5b506100d360043561045e565b60408051600160a060020a039092168252519081900360200190f35b3480156100fb57600080fd5b50610107600435610486565b005b34801561011557600080fd5b5061012a600160a060020a03600435166105da565b6040805161ffff9092168252519081900360200190f35b34801561014d57600080fd5b506040805160206004604435818101358381028086018501909652808552610107958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506105f09650505050505050565b34801561022a57600080fd5b50610107600160a060020a0360043516602435604435606435610d18565b34801561025457600080fd5b506100d360043561105c565b34801561026c57600080fd5b50610107600160a060020a036004351660243561106a565b34801561029057600080fd5b506102a5600160a060020a0360043516611369565b60408051938452602084019290925282820152519081900360600190f35b3480156102cf57600080fd5b50610107600160a060020a0360043516611389565b3480156102f057600080fd5b506100d36115a0565b34801561030557600080fd5b506100d36004356115af565b34801561031d57600080fd5b50610107600160a060020a03600435166024356044356064356115bd565b34801561034757600080fd5b506100d36118d5565b34801561035c57600080fd5b50610371600160a060020a03600435166118e4565b604051808561ffff1661ffff168152602001806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b838110156103c75781810151838201526020016103af565b50505050905001848103835286818151815260200191508051906020019060200280838360005b838110156104065781810151838201526020016103ee565b50505050905001848103825285818151815260200191508051906020019060200280838360005b8381101561044557818101518382015260200161042d565b5050505090500197505050505050505060405180910390f35b600280548290811061046c57fe5b600091825260209091200154600160a060020a0316905081565b600754600160a060020a031633146104e8576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b603c811015610541576040805160e560020a62461bcd02815260206004820152600c60248201527f496e76616c6964207465726d0000000000000000000000000000000000000000604482015290519081900360640190fd5b61059b600480548060200260200160405190810160405280929190818152602001828054801561059057602002820191906000526020600020905b81548152602001906001019080831161057c575b505050505082611a2d565b15156105d757600480546001810182556000919091527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b018190555b50565b60006020819052908152604090205461ffff1681565b600654604080517f540105c7000000000000000000000000000000000000000000000000000000008152600160a060020a03888116600483015291516000938493849391169163540105c791602480820192869290919082900301818387803b15801561065c57600080fd5b505af1158015610670573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405260c081101561069957600080fd5b8151602083015160408401516060850151608086018051949693959294919392830192916401000000008111156106cf57600080fd5b820160208101848111156106e257600080fd5b81518560208202830111640100000000821117156106ff57600080fd5b5050929190602001805164010000000081111561071b57600080fd5b8201602081018481111561072e57600080fd5b815185602082028301116401000000008211171561074b57600080fd5b50979b505050600160a060020a038a16331496506107ba95505050505050576040805160e560020a62461bcd02815260206004820152601660248201527f52656c61796572206f776e657220726571756972656400000000000000000000604482015290519081900360640190fd5b600654604080517f500f99f7000000000000000000000000000000000000000000000000000000008152600160a060020a038b811660048301529151919092169163500f99f79160248083019260209291908290030181600087803b15801561082257600080fd5b505af1158015610836573d6000803e3d6000fd5b505050506040513d602081101561084c57600080fd5b5051156108a3576040805160e560020a62461bcd02815260206004820152601960248201527f52656c6179657220726571756972656420746f20636c6f736500000000000000604482015290519081900360640190fd5b60018761ffff16101580156108bd57506103e88761ffff16105b1515610913576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c696420747261646520466565000000000000000000000000000000604482015290519081900360640190fd5b845186511461096c576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f66207465726d7300000000000000604482015290519081900360640190fd5b83518651146109c5576040805160e560020a62461bcd02815260206004820152601f60248201527f4e6f742076616c6964206e756d626572206f6620636f6c6c61746572616c7300604482015290519081900360640190fd5b5060009050805b8551811015610ab457610a506003805480602002602001604051908101604052809291908181526020018280548015610a2e57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610a10575b50505050508783815181101515610a4157fe5b90602001906020020151611a76565b9150600182151514610aac576040805160e560020a62461bcd02815260206004820152601560248201527f496e76616c6964206c656e64696e6720746f6b656e0000000000000000000000604482015290519081900360640190fd5b6001016109cc565b5060005b8451811015610b9657610b326004805480602002602001604051908101604052809291908181526020018280548015610b1057602002820191906000526020600020905b815481526020019060010190808311610afc575b50505050508683815181101515610b2357fe5b90602001906020020151611a2d565b9150600182151514610b8e576040805160e560020a62461bcd02815260206004820152600c60248201527f496e76616c6964207465726d0000000000000000000000000000000000000000604482015290519081900360640190fd5b600101610ab8565b5060005b8351811015610c6d57610c0f6005805480602002602001604051908101604052809291908181526020018280548015610bfc57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610bde575b50505050508583815181101515610a4157fe5b1515610c65576040805160e560020a62461bcd02815260206004820152601f60248201527f4e6f742076616c6964206e756d626572206f6620636f6c6c61746572616c7300604482015290519081900360640190fd5b600101610b9a565b6040805160808101825261ffff898116825260208083018a81528385018a905260608401899052600160a060020a038d166000908152808352949094208351815461ffff191693169290921782559251805192939192610cd39260018501920190611ac5565b5060408201518051610cef916002840191602090910190611b37565b5060608201518051610d0b916003840191602090910190611ac5565b5050505050505050505050565b60008060648510158015610d2c5750606484115b1515610d70576040805160e560020a62461bcd02815260206004820152600d6024820152600080516020611bcd833981519152604482015290519081900360640190fd5b838511610db5576040805160e560020a62461bcd02815260206004820152600d6024820152600080516020611bcd833981519152604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a0389811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b158015610e0757600080fd5b505af1158015610e1b573d6000803e3d6000fd5b505050506040513d6020811015610e3157600080fd5b50519150811515610e7a576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611bed833981519152604482015290519081900360640190fd5b85905033600160a060020a031681600160a060020a0316631d1438486040518163ffffffff1660e060020a028152600401602060405180830381600087803b158015610ec557600080fd5b505af1158015610ed9573d6000803e3d6000fd5b505050506040513d6020811015610eef57600080fd5b5051600160a060020a031614610f4f576040805160e560020a62461bcd02815260206004820152601560248201527f526571756972656420746f6b656e206973737565720000000000000000000000604482015290519081900360640190fd5b604080516060810182528681526020808201878152828401878152600160a060020a038b166000908152600180855290869020945185559151918401919091555160029092019190915560058054835181840281018401909452808452610ff59392830182828015610fea57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610fcc575b505050505087611a76565b151561105457600580546001810182556000919091527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db001805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0388161790555b505050505050565b600380548290811061046c57fe5b6008546040805160e060020a63a3ff31b5028152600160a060020a03858116600483015291516000938493169163a3ff31b591602480830192602092919082900301818787803b1580156110bd57600080fd5b505af11580156110d1573d6000803e3d6000fd5b505050506040513d60208110156110e757600080fd5b5051806110fd5750600160a060020a0384166001145b9150811515611144576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611bed833981519152604482015290519081900360640190fd5b600160a060020a038416600090815260016020526040902054606411156111a3576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611bed833981519152604482015290519081900360640190fd5b61120760028054806020026020016040519081016040528092919081815260200182805480156111fc57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116111de575b505050505085611a76565b1561127357600754600160a060020a0316331461126e576040805160e560020a62461bcd02815260206004820152601760248201527f436f6e7472616374206f776e6572207265717569726564000000000000000000604482015290519081900360640190fd5b611348565b83905033600160a060020a031681600160a060020a0316631d1438486040518163ffffffff1660e060020a028152600401602060405180830381600087803b1580156112be57600080fd5b505af11580156112d2573d6000803e3d6000fd5b505050506040513d60208110156112e857600080fd5b5051600160a060020a031614611348576040805160e560020a62461bcd02815260206004820152601560248201527f526571756972656420746f6b656e206973737565720000000000000000000000604482015290519081900360640190fd5b5050600160a060020a03909116600090815260016020526040902060020155565b600160208190526000918252604090912080549181015460029091015483565b600754600090600160a060020a031633146113ee576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a0385811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b15801561144057600080fd5b505af1158015611454573d6000803e3d6000fd5b505050506040513d602081101561146a57600080fd5b5051806114805750600160a060020a0382166001145b90508015156114d9576040805160e560020a62461bcd02815260206004820152601260248201527f496e76616c6964206261736520746f6b656e0000000000000000000000000000604482015290519081900360640190fd5b61153d600380548060200260200160405190810160405280929190818152602001828054801561153257602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611514575b505050505083611a76565b151561159c57600380546001810182556000919091527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b01805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0384161790555b5050565b600654600160a060020a031681565b600580548290811061046c57fe5b600754600090600160a060020a03163314611622576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b606484101580156116335750606483115b1515611677576040805160e560020a62461bcd02815260206004820152600d6024820152600080516020611bcd833981519152604482015290519081900360640190fd5b8284116116bc576040805160e560020a62461bcd02815260206004820152600d6024820152600080516020611bcd833981519152604482015290519081900360640190fd5b6008546040805160e060020a63a3ff31b5028152600160a060020a0388811660048301529151919092169163a3ff31b59160248083019260209291908290030181600087803b15801561170e57600080fd5b505af1158015611722573d6000803e3d6000fd5b505050506040513d602081101561173857600080fd5b50518061174e5750600160a060020a0385166001145b9050801515611795576040805160e560020a62461bcd0281526020600482015260126024820152600080516020611bed833981519152604482015290519081900360640190fd5b604080516060810182528581526020808201868152828401868152600160a060020a038a1660009081526001808552908690209451855591519184019190915551600292830155815483518183028101830190945280845261183793929183018282801561182c57602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161180e575b505050505086611a76565b15156118ce576002805460018181019092557f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054600160a060020a03881673ffffffffffffffffffffffffffffffffffffffff1991821681179092556005805493840181556000527f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db090920180549092161790555b5050505050565b600754600160a060020a031681565b600160a060020a03811660009081526020818152604080832080546001820180548451818702810187019095528085526060958695869561ffff90951694600281019360039091019285919083018282801561196957602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161194b575b50505050509250818054806020026020016040519081016040528092919081815260200182805480156119bb57602002820191906000526020600020905b8154815260200190600101908083116119a7575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015611a1757602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116119f9575b5050505050905093509350935093509193509193565b6000805b8351811015611a6a57828482815181101515611a4957fe5b906020019060200201511415611a625760019150611a6f565b600101611a31565b600091505b5092915050565b6000805b8351811015611a6a5782600160a060020a03168482815181101515611a9b57fe5b90602001906020020151600160a060020a03161415611abd5760019150611a6f565b600101611a7a565b828054828255906000526020600020908101928215611b27579160200282015b82811115611b27578251825473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03909116178255602090920191600190910190611ae5565b50611b33929150611b7e565b5090565b828054828255906000526020600020908101928215611b72579160200282015b82811115611b72578251825591602001919060010190611b57565b50611b33929150611bb2565b611baf91905b80821115611b3357805473ffffffffffffffffffffffffffffffffffffffff19168155600101611b84565b90565b611baf91905b80821115611b335760008155600101611bb85600496e76616c696420726174657300000000000000000000000000000000000000496e76616c696420636f6c6c61746572616c0000000000000000000000000000a165627a7a7230582016452d9cde49a89380b8cc9dddfa746784e7ebff15def5d464a5b90162fd70230029`

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

// ALLCOLLATERALS is a free data retrieval call binding the contract method 0xc4691d36.
//
// Solidity: function ALL_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCaller) ALLCOLLATERALS(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Lending.contract.Call(opts, out, "ALL_COLLATERALS", arg0)
	return *ret0, err
}

// ALLCOLLATERALS is a free data retrieval call binding the contract method 0xc4691d36.
//
// Solidity: function ALL_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingSession) ALLCOLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.ALLCOLLATERALS(&_Lending.CallOpts, arg0)
}

// ALLCOLLATERALS is a free data retrieval call binding the contract method 0xc4691d36.
//
// Solidity: function ALL_COLLATERALS( uint256) constant returns(address)
func (_Lending *LendingCallerSession) ALLCOLLATERALS(arg0 *big.Int) (common.Address, error) {
	return _Lending.Contract.ALLCOLLATERALS(&_Lending.CallOpts, arg0)
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
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _price uint256)
func (_Lending *LendingCaller) COLLATERALLIST(opts *bind.CallOpts, arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	Price           *big.Int
}, error) {
	ret := new(struct {
		DepositRate     *big.Int
		LiquidationRate *big.Int
		Price           *big.Int
	})
	out := ret
	err := _Lending.contract.Call(opts, out, "COLLATERAL_LIST", arg0)
	return *ret, err
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _price uint256)
func (_Lending *LendingSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	Price           *big.Int
}, error) {
	return _Lending.Contract.COLLATERALLIST(&_Lending.CallOpts, arg0)
}

// COLLATERALLIST is a free data retrieval call binding the contract method 0x82250701.
//
// Solidity: function COLLATERAL_LIST( address) constant returns(_depositRate uint256, _liquidationRate uint256, _price uint256)
func (_Lending *LendingCallerSession) COLLATERALLIST(arg0 common.Address) (struct {
	DepositRate     *big.Int
	LiquidationRate *big.Int
	Price           *big.Int
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

// AddCollateral is a paid mutator transaction binding the contract method 0xe5eecf68.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingTransactor) AddCollateral(opts *bind.TransactOpts, token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addCollateral", token, depositRate, liquidationRate, price)
}

// AddCollateral is a paid mutator transaction binding the contract method 0xe5eecf68.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingSession) AddCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, price)
}

// AddCollateral is a paid mutator transaction binding the contract method 0xe5eecf68.
//
// Solidity: function addCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingTransactorSession) AddCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x3b874827.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingTransactor) AddILOCollateral(opts *bind.TransactOpts, token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.contract.Transact(opts, "addILOCollateral", token, depositRate, liquidationRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x3b874827.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingSession) AddILOCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddILOCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, price)
}

// AddILOCollateral is a paid mutator transaction binding the contract method 0x3b874827.
//
// Solidity: function addILOCollateral(token address, depositRate uint256, liquidationRate uint256, price uint256) returns()
func (_Lending *LendingTransactorSession) AddILOCollateral(token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) (*types.Transaction, error) {
	return _Lending.Contract.AddILOCollateral(&_Lending.TransactOpts, token, depositRate, liquidationRate, price)
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
