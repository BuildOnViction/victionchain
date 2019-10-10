// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// RelayerRegistrationABI is the input ABI used to generate the binding from.
const RelayerRegistrationABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"MaximumRelayers\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"depositMore\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"RELAYER_COINBASES\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"getRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"maxRelayer\",\"type\":\"uint256\"},{\"name\":\"maxToken\",\"type\":\"uint256\"},{\"name\":\"minDeposit\",\"type\":\"uint256\"}],\"name\":\"reconfigure\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"cancelSelling\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"sellRelayer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"RelayerCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"RELAYER_ON_SALE_LIST\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"resign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"new_owner\",\"type\":\"address\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MinimumDeposit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"register\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MaximumTokenList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"buyRelayer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_OWNER\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"maxRelayers\",\"type\":\"uint256\"},{\"name\":\"maxTokenList\",\"type\":\"uint256\"},{\"name\":\"minDeposit\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"max_relayer\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"max_token\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"min_deposit\",\"type\":\"uint256\"}],\"name\":\"ConfigEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"RegisterEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"UpdateEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"TransferEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit_release_time\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deposit_amount\",\"type\":\"uint256\"}],\"name\":\"ResignEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"remaining_time\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deposit_amount\",\"type\":\"uint256\"}],\"name\":\"RefundEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"is_on_sale\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"coinbase\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"SellEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"coinbase\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"BuyEvent\",\"type\":\"event\"}]"

// RelayerRegistrationBin is the compiled bytecode used for deploying new contracts.
const RelayerRegistrationBin = `0x608060405234801561001057600080fd5b506040516060806125c8833981016040908152815160208301519190920151600060078190556001939093556002919091556008558054600160a060020a03191633179055612564806100646000396000f3006080604052600436106100fb5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630e5c0fee81146101005780634ce69bf5146101275780634fa339271461013d578063540105c71461017157806356246b681461026157806357ea3c41146103055780635b673b1f1461032357806387c6bbcd1461034457806387d340ab14610368578063885b71371461037d578063ae6e43f51461039e578063ba45b0b8146103bf578063c635a9f2146103e6578063c6c71aed146103fb578063cfaece1214610492578063e699df0e146104a7578063fa89401a146104bb578063fd301c49146104dc575b600080fd5b34801561010c57600080fd5b506101156104f1565b60408051918252519081900360200190f35b61013b600160a060020a03600435166104f7565b005b34801561014957600080fd5b50610155600435610859565b60408051600160a060020a039092168252519081900360200190f35b34801561017d57600080fd5b50610192600160a060020a0360043516610874565b6040518087815260200186600160a060020a0316600160a060020a031681526020018581526020018461ffff1661ffff1681526020018060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156102085781810151838201526020016101f0565b50505050905001838103825284818151815260200191508051906020019060200280838360005b8381101561024757818101518382015260200161022f565b505050509050019850505050505050505060405180910390f35b34801561026d57600080fd5b50604080516020600460443581810135838102808601850190965280855261013b958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506109879650505050505050565b34801561031157600080fd5b5061013b600435602435604435610d7b565b34801561032f57600080fd5b5061013b600160a060020a0360043516610e6a565b34801561035057600080fd5b5061013b600160a060020a0360043516602435611036565b34801561037457600080fd5b506101156111ed565b34801561038957600080fd5b50610115600160a060020a03600435166111f3565b3480156103aa57600080fd5b5061013b600160a060020a0360043516611205565b3480156103cb57600080fd5b5061013b600160a060020a0360043581169060243516611448565b3480156103f257600080fd5b50610115611792565b604080516020600460443581810135838102808601850190965280855261013b958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506117989650505050505050565b34801561049e57600080fd5b50610115611d89565b61013b600160a060020a0360043516611d8f565b3480156104c757600080fd5b5061013b600160a060020a0360043516612048565b3480156104e857600080fd5b50610155612365565b60015481565b600160a060020a03818116600090815260036020526040902060050154829116331461055b576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a0382166000908152600560205260409020548290156105cd576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b600160a060020a03831660009081526006602052604090205483901561063f576040805160e560020a62461bcd02815260206004820152602a60248201526000805160206124d98339815191526044820152600080516020612519833981519152606482015290519081900360840190fd5b60003411610697576040805160e560020a62461bcd02815260206004820152601a60248201527f5472616e736665722076616c7565206d757374206265203e2030000000000000604482015290519081900360640190fd5b670de0b6b3a764000034101561071d576040805160e560020a62461bcd02815260206004820152603160248201527f4174206c65617374203120544f4d4f20697320726571756972656420666f722060448201527f61206465706f7369742072657175657374000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a038416600090815260036020526040902054610746903463ffffffff61237416565b50600160a060020a03841660009081526003602081815260409283902080546001820154855182815261ffff90911693810184905260809581018681526002840180549783018890527fcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f97939690940192606083019060a0840190869080156107f857602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116107da575b5050838103825284818154815260200191508054801561084157602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610823575b5050965050505050505060405180910390a150505050565b600460205260009081526040902054600160a060020a031681565b600160a060020a03808216600090815260036020818152604080842060048101546005820154825460018401546002850180548751818a0281018a01909852808852999a8b9a8b9a8b9a60609a8b9a989094169761ffff909616959091019291849183018282801561090f57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116108f1575b505050505091508080548060200260200160405190810160405280929190818152602001828054801561096b57602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161094d575b5050505050905095509550955095509550955091939550919395565b600160a060020a0384811660009081526003602052604090206005015485911633146109eb576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a038516600090815260056020526040902054859015610a5d576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b600160a060020a038616600090815260066020526040902054869015610acf576040805160e560020a62461bcd02815260206004820152602a60248201526000805160206124d98339815191526044820152600080516020612519833981519152606482015290519081900360840190fd5b60018661ffff1610158015610ae957506127108661ffff16105b1515610b3f576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c6964204d616b657220466565000000000000000000000000000000604482015290519081900360640190fd5b60025485511115610b9a576040805160e560020a62461bcd02815260206004820152601f60248201527f457863656564696e67206e756d626572206f6620747261646520706169727300604482015290519081900360640190fd5b8451845114610bf3576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f6620506169727300000000000000604482015290519081900360640190fd5b600160a060020a038716600090815260036020908152604090912060018101805461ffff191661ffff8a161790558651610c359260029092019188019061238d565b50600160a060020a03871660009081526003602081815260409092208651610c659391909201919087019061238d565b50600160a060020a03871660009081526003602081815260409283902080546001820154855182815261ffff90911693810184905260809581018681526002840180549783018890527fcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f97939690940192606083019060a084019086908015610d1757602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610cf9575b50508381038252848181548152602001915080548015610d6057602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610d42575b5050965050505050505060405180910390a150505050505050565b600054600160a060020a03163314610ddd576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b600754831015610dec57600080fd5b600482118015610dfd57506103e982105b1515610e0857600080fd5b6127108111610e1657600080fd5b600183905560028290556008819055604080518481526020810184905280820183905290517f8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d939239181900360600190a1505050565b600160a060020a038181166000908152600360205260409020600501548291163314610ece576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a038216600090815260056020526040902054829015610f40576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b600160a060020a03831660009081526006602052604081205411610fd4576040805160e560020a62461bcd02815260206004820152602160248201527f52656c61796572206973206e6f742063757272656e746c7920666f722073616c60448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a03831660008181526006602090815260408083208390558051838152918201939093528083019190915290517fdb3d5e65fcde89731529c01d62b87bab1c64471cffdd528fc1adbc1712b5d0829181900360600190a1505050565b600160a060020a03828116600090815260036020526040902060050154839116331461109a576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a03831660009081526005602052604090205483901561110c576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b6000831161118a576040805160e560020a62461bcd02815260206004820152602860248201527f507269636520746167206d75737420626520646966666572656e74207468616e60448201527f205a65726f283029000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a03841660008181526006602090815260409182902086905581516001815290810192909252818101859052517fdb3d5e65fcde89731529c01d62b87bab1c64471cffdd528fc1adbc1712b5d0829181900360600190a150505050565b60075481565b60066020526000908152604090205481565b600160a060020a038181166000908152600360205260409020600501548291163314611269576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a0382166000908152600660205260409020548290156112db576040805160e560020a62461bcd02815260206004820152602a60248201526000805160206124d98339815191526044820152600080516020612519833981519152606482015290519081900360840190fd5b600160a060020a0383166000908152600360205260408120541161136f576040805160e560020a62461bcd02815260206004820152602760248201527f4e6f2072656c61796572206173736f636961746564207769746820746869732060448201527f6164647265737300000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a038316600090815260056020526040902054156113dd576040805160e560020a62461bcd02815260206004820152601860248201527f5265717565737420616c72656164792072656365697665640000000000000000604482015290519081900360640190fd5b600160a060020a03831660009081526005602090815260408083206224ea0042019081905560038352928190205481519384529183019190915280517f2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae59281900390910190a1505050565b600160a060020a0382811660009081526003602052604090206005015483911633146114ac576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a03831660009081526005602052604090205483901561151e576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b600160a060020a038416600090815260066020526040902054849015611590576040805160e560020a62461bcd02815260206004820152602a60248201526000805160206124d98339815191526044820152600080516020612519833981519152606482015290519081900360840190fd5b600160a060020a038416158015906115b15750600160a060020a0384163314155b15156115bc57600080fd5b600160a060020a03841660009081526003602052604090206001015461ffff1615611657576040805160e560020a62461bcd02815260206004820152603c60248201527f4f776e65722061646472657373206d757374206e6f742062652063757272656e60448201527f746c7920757365642061732072656c617965722d636f696e6261736500000000606482015290519081900360840190fd5b600160a060020a03858116600090815260036020818152604092839020600581018054600160a060020a0319168a871617908190558154600183015486519290971680835293820181905261ffff90961694810185905260a0606082018181526002840180549284018390527fc13ab794f75ba420a1f52192a8e35a2cf2c74ae31ed94f53f47ce7712011b66298959795969094019291608083019060c08401908690801561172f57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611711575b5050838103825284818154815260200191508054801561177857602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161175a575b505097505050505050505060405180910390a15050505050565b60085481565b6117a06123f2565b600054600160a060020a0316331415611829576040805160e560020a62461bcd02815260206004820152602f60248201527f436f6e7472616374204f776e657220697320666f7262696464656e20746f206360448201527f726561746520612052656c617965720000000000000000000000000000000000606482015290519081900360840190fd5b33600160a060020a03861614156118b0576040805160e560020a62461bcd02815260206004820152603660248201527f436f696e6261736520616e642052656c617965724f776e65722061646472657360448201527f73206d757374206e6f74206265207468652073616d6500000000000000000000606482015290519081900360840190fd5b600054600160a060020a038681169116141561193c576040805160e560020a62461bcd02815260206004820152602b60248201527f436f696e62617365206d757374206e6f742062652073616d6520617320434f4e60448201527f54524143545f4f574e4552000000000000000000000000000000000000000000606482015290519081900360840190fd5b600854341015611996576040805160e560020a62461bcd02815260206004820152601e60248201527f4d696e696d756d206465706f736974206e6f74207361746973666965642e0000604482015290519081900360640190fd5b60018461ffff16101580156119b057506127108461ffff16105b1515611a06576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c6964204d616b657220466565000000000000000000000000000000604482015290519081900360640190fd5b60025483511115611a61576040805160e560020a62461bcd02815260206004820152601f60248201527f457863656564696e67206e756d626572206f6620747261646520706169727300604482015290519081900360640190fd5b8251825114611aba576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f6620506169727300000000000000604482015290519081900360640190fd5b600160a060020a03851660009081526003602052604090205415611b28576040805160e560020a62461bcd02815260206004820152601c60248201527f436f696e6261736520616c726561647920726567697374657265642e00000000604482015290519081900360640190fd5b60015460075410611b83576040805160e560020a62461bcd02815260206004820152601b60248201527f4d6178696d756d2072656c617965727320726567697374657265640000000000604482015290519081900360640190fd5b506040805160c08101825234815261ffff858116602080840191825283850187815260608501879052600754608086018190523360a08701526000908152600483528681208054600160a060020a031916600160a060020a038d169081179091558152600383529590952084518155915160018301805461ffff1916919094161790925592518051929384939092611c2292600285019291019061238d565b5060608201518051611c3e91600384019160209091019061238d565b50608082810151600483015560a09283015160059092018054600160a060020a031916600160a060020a03938416179055600780546001908101909155918816600090815260036020818152604092839020805495810154845187815261ffff9091169281018390529384018581526002820180549686018790527fcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b9893969095929094019391929091606084019184019086908015611d2757602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611d09575b50508381038252848181548152602001915080548015611d7057602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611d52575b5050965050505050505060405180910390a15050505050565b60025481565b600160a060020a0381166000908152600560205260408120548190839015611e03576040805160e560020a62461bcd02815260206004820152602860248201526000805160206124b983398151915260448201526000805160206124f9833981519152606482015290519081900360840190fd5b600160a060020a03841660009081526006602052604081205493508311611e9a576040805160e560020a62461bcd02815260206004820152602160248201527f52656c61796572206973206e6f742063757272656e746c7920666f722073616c60448201527f6500000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b348314611ef1576040805160e560020a62461bcd02815260206004820152601960248201527f50726963652d746167206d757374206265206d61746368656400000000000000604482015290519081900360640190fd5b600160a060020a038085166000908152600360205260409020600501541691503315801590611f29575033600160a060020a03831614155b8015611f3d5750600160a060020a03821615155b1515611f93576040805160e560020a62461bcd02815260206004820152601160248201527f41646472657373206e6f742076616c6964000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a0380851660009081526003602090815260408083206005018054600160a060020a031916331790556006909152808220829055519184169185156108fc0291869190818181858888f19350505050158015611ff9573d6000803e3d6000fd5b506040805160018152600160a060020a0386166020820152348183015290517f07e248a3b3d2184a9491c3b45089a6e15aac742b9d974e691e7beb0f6e7c58c69181900360600190a150505050565b600160a060020a0381811660009081526003602052604081206005015490918291829185911633146120b2576040805160e560020a62461bcd0281526020600482015260136024820152600080516020612499833981519152604482015290519081900360640190fd5b600160a060020a038516600090815260066020526040902054859015612124576040805160e560020a62461bcd02815260206004820152602a60248201526000805160206124d98339815191526044820152600080516020612519833981519152606482015290519081900360840190fd5b600160a060020a03861660009081526005602052604081205411612192576040805160e560020a62461bcd02815260206004820152601160248201527f52657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a038616600090815260036020908152604080832080546004909101546005909352922054919650945042111561230057600160a060020a038616600090815260036020526040812081815560018101805461ffff19169055906121ff6002830182612436565b61220d600383016000612436565b506000600482810182905560059283018054600160a060020a0319908116909155600160a060020a038a811684526020948552604080852085905560078054600019908101875285885282872080548087169091558c88528388208054919095169516851790935583865260039096528085209093018990558454019093555191945033916108fc88150291889190818181858888f193505050501580156122b9573d6000803e3d6000fd5b5060408051600181526000602082015280820187905290517ffaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d302509181900360600190a161235d565b600160a060020a038616600090815260056020908152604080832054815193845242900391830191909152818101879052517ffaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d302509181900360600190a15b505050505050565b600054600160a060020a031681565b60008282018381101561238657600080fd5b9392505050565b8280548282559060005260206000209081019282156123e2579160200282015b828111156123e25782518254600160a060020a031916600160a060020a039091161782556020909201916001909101906123ad565b506123ee929150612457565b5090565b60c06040519081016040528060008152602001600061ffff1681526020016060815260200160608152602001600081526020016000600160a060020a031681525090565b5080546000825590600052602060002090810190612454919061247e565b50565b61247b91905b808211156123ee578054600160a060020a031916815560010161245d565b90565b61247b91905b808211156123ee5760008155600101612484560052656c61796572204f776e6572204f6e6c792e000000000000000000000000005468652072656c6179657220686173206265656e2072657175657374656420745468652072656c61796572206d757374206265206e6f742063757272656e746c6f20636c6f73652e0000000000000000000000000000000000000000000000007920666f722053616c6500000000000000000000000000000000000000000000a165627a7a72305820eb227ccc6810f0152c3830c58c0811f6334019fc0ad0d2634f0791c6867bc3400029`

// DeployRelayerRegistration deploys a new Ethereum contract, binding an instance of RelayerRegistration to it.
func DeployRelayerRegistration(auth *bind.TransactOpts, backend bind.ContractBackend, maxRelayers *big.Int, maxTokenList *big.Int, minDeposit *big.Int) (common.Address, *types.Transaction, *RelayerRegistration, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayerRegistrationABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RelayerRegistrationBin), backend, maxRelayers, maxTokenList, minDeposit)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RelayerRegistration{RelayerRegistrationCaller: RelayerRegistrationCaller{contract: contract}, RelayerRegistrationTransactor: RelayerRegistrationTransactor{contract: contract}, RelayerRegistrationFilterer: RelayerRegistrationFilterer{contract: contract}}, nil
}

// RelayerRegistration is an auto generated Go binding around an Ethereum contract.
type RelayerRegistration struct {
	RelayerRegistrationCaller     // Read-only binding to the contract
	RelayerRegistrationTransactor // Write-only binding to the contract
	RelayerRegistrationFilterer   // Log filterer for contract events
}

// RelayerRegistrationCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayerRegistrationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerRegistrationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayerRegistrationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerRegistrationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayerRegistrationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayerRegistrationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelayerRegistrationSession struct {
	Contract     *RelayerRegistration // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RelayerRegistrationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayerRegistrationCallerSession struct {
	Contract *RelayerRegistrationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// RelayerRegistrationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayerRegistrationTransactorSession struct {
	Contract     *RelayerRegistrationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// RelayerRegistrationRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayerRegistrationRaw struct {
	Contract *RelayerRegistration // Generic contract binding to access the raw methods on
}

// RelayerRegistrationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayerRegistrationCallerRaw struct {
	Contract *RelayerRegistrationCaller // Generic read-only contract binding to access the raw methods on
}

// RelayerRegistrationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayerRegistrationTransactorRaw struct {
	Contract *RelayerRegistrationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelayerRegistration creates a new instance of RelayerRegistration, bound to a specific deployed contract.
func NewRelayerRegistration(address common.Address, backend bind.ContractBackend) (*RelayerRegistration, error) {
	contract, err := bindRelayerRegistration(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RelayerRegistration{RelayerRegistrationCaller: RelayerRegistrationCaller{contract: contract}, RelayerRegistrationTransactor: RelayerRegistrationTransactor{contract: contract}, RelayerRegistrationFilterer: RelayerRegistrationFilterer{contract: contract}}, nil
}

// NewRelayerRegistrationCaller creates a new read-only instance of RelayerRegistration, bound to a specific deployed contract.
func NewRelayerRegistrationCaller(address common.Address, caller bind.ContractCaller) (*RelayerRegistrationCaller, error) {
	contract, err := bindRelayerRegistration(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationCaller{contract: contract}, nil
}

// NewRelayerRegistrationTransactor creates a new write-only instance of RelayerRegistration, bound to a specific deployed contract.
func NewRelayerRegistrationTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayerRegistrationTransactor, error) {
	contract, err := bindRelayerRegistration(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationTransactor{contract: contract}, nil
}

// NewRelayerRegistrationFilterer creates a new log filterer instance of RelayerRegistration, bound to a specific deployed contract.
func NewRelayerRegistrationFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayerRegistrationFilterer, error) {
	contract, err := bindRelayerRegistration(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationFilterer{contract: contract}, nil
}

// bindRelayerRegistration binds a generic wrapper to an already deployed contract.
func bindRelayerRegistration(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayerRegistrationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerRegistration *RelayerRegistrationRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RelayerRegistration.Contract.RelayerRegistrationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerRegistration *RelayerRegistrationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.RelayerRegistrationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerRegistration *RelayerRegistrationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.RelayerRegistrationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RelayerRegistration *RelayerRegistrationCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RelayerRegistration.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RelayerRegistration *RelayerRegistrationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RelayerRegistration *RelayerRegistrationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.contract.Transact(opts, method, params...)
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_RelayerRegistration *RelayerRegistrationCaller) CONTRACTOWNER(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "CONTRACT_OWNER")
	return *ret0, err
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_RelayerRegistration *RelayerRegistrationSession) CONTRACTOWNER() (common.Address, error) {
	return _RelayerRegistration.Contract.CONTRACTOWNER(&_RelayerRegistration.CallOpts)
}

// CONTRACTOWNER is a free data retrieval call binding the contract method 0xfd301c49.
//
// Solidity: function CONTRACT_OWNER() constant returns(address)
func (_RelayerRegistration *RelayerRegistrationCallerSession) CONTRACTOWNER() (common.Address, error) {
	return _RelayerRegistration.Contract.CONTRACTOWNER(&_RelayerRegistration.CallOpts)
}

// MaximumRelayers is a free data retrieval call binding the contract method 0x0e5c0fee.
//
// Solidity: function MaximumRelayers() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCaller) MaximumRelayers(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "MaximumRelayers")
	return *ret0, err
}

// MaximumRelayers is a free data retrieval call binding the contract method 0x0e5c0fee.
//
// Solidity: function MaximumRelayers() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationSession) MaximumRelayers() (*big.Int, error) {
	return _RelayerRegistration.Contract.MaximumRelayers(&_RelayerRegistration.CallOpts)
}

// MaximumRelayers is a free data retrieval call binding the contract method 0x0e5c0fee.
//
// Solidity: function MaximumRelayers() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCallerSession) MaximumRelayers() (*big.Int, error) {
	return _RelayerRegistration.Contract.MaximumRelayers(&_RelayerRegistration.CallOpts)
}

// MaximumTokenList is a free data retrieval call binding the contract method 0xcfaece12.
//
// Solidity: function MaximumTokenList() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCaller) MaximumTokenList(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "MaximumTokenList")
	return *ret0, err
}

// MaximumTokenList is a free data retrieval call binding the contract method 0xcfaece12.
//
// Solidity: function MaximumTokenList() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationSession) MaximumTokenList() (*big.Int, error) {
	return _RelayerRegistration.Contract.MaximumTokenList(&_RelayerRegistration.CallOpts)
}

// MaximumTokenList is a free data retrieval call binding the contract method 0xcfaece12.
//
// Solidity: function MaximumTokenList() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCallerSession) MaximumTokenList() (*big.Int, error) {
	return _RelayerRegistration.Contract.MaximumTokenList(&_RelayerRegistration.CallOpts)
}

// MinimumDeposit is a free data retrieval call binding the contract method 0xc635a9f2.
//
// Solidity: function MinimumDeposit() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCaller) MinimumDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "MinimumDeposit")
	return *ret0, err
}

// MinimumDeposit is a free data retrieval call binding the contract method 0xc635a9f2.
//
// Solidity: function MinimumDeposit() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationSession) MinimumDeposit() (*big.Int, error) {
	return _RelayerRegistration.Contract.MinimumDeposit(&_RelayerRegistration.CallOpts)
}

// MinimumDeposit is a free data retrieval call binding the contract method 0xc635a9f2.
//
// Solidity: function MinimumDeposit() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCallerSession) MinimumDeposit() (*big.Int, error) {
	return _RelayerRegistration.Contract.MinimumDeposit(&_RelayerRegistration.CallOpts)
}

// RELAYERCOINBASES is a free data retrieval call binding the contract method 0x4fa33927.
//
// Solidity: function RELAYER_COINBASES( uint256) constant returns(address)
func (_RelayerRegistration *RelayerRegistrationCaller) RELAYERCOINBASES(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "RELAYER_COINBASES", arg0)
	return *ret0, err
}

// RELAYERCOINBASES is a free data retrieval call binding the contract method 0x4fa33927.
//
// Solidity: function RELAYER_COINBASES( uint256) constant returns(address)
func (_RelayerRegistration *RelayerRegistrationSession) RELAYERCOINBASES(arg0 *big.Int) (common.Address, error) {
	return _RelayerRegistration.Contract.RELAYERCOINBASES(&_RelayerRegistration.CallOpts, arg0)
}

// RELAYERCOINBASES is a free data retrieval call binding the contract method 0x4fa33927.
//
// Solidity: function RELAYER_COINBASES( uint256) constant returns(address)
func (_RelayerRegistration *RelayerRegistrationCallerSession) RELAYERCOINBASES(arg0 *big.Int) (common.Address, error) {
	return _RelayerRegistration.Contract.RELAYERCOINBASES(&_RelayerRegistration.CallOpts, arg0)
}

// RELAYERONSALELIST is a free data retrieval call binding the contract method 0x885b7137.
//
// Solidity: function RELAYER_ON_SALE_LIST( address) constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCaller) RELAYERONSALELIST(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "RELAYER_ON_SALE_LIST", arg0)
	return *ret0, err
}

// RELAYERONSALELIST is a free data retrieval call binding the contract method 0x885b7137.
//
// Solidity: function RELAYER_ON_SALE_LIST( address) constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationSession) RELAYERONSALELIST(arg0 common.Address) (*big.Int, error) {
	return _RelayerRegistration.Contract.RELAYERONSALELIST(&_RelayerRegistration.CallOpts, arg0)
}

// RELAYERONSALELIST is a free data retrieval call binding the contract method 0x885b7137.
//
// Solidity: function RELAYER_ON_SALE_LIST( address) constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCallerSession) RELAYERONSALELIST(arg0 common.Address) (*big.Int, error) {
	return _RelayerRegistration.Contract.RELAYERONSALELIST(&_RelayerRegistration.CallOpts, arg0)
}

// RelayerCount is a free data retrieval call binding the contract method 0x87d340ab.
//
// Solidity: function RelayerCount() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCaller) RelayerCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "RelayerCount")
	return *ret0, err
}

// RelayerCount is a free data retrieval call binding the contract method 0x87d340ab.
//
// Solidity: function RelayerCount() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationSession) RelayerCount() (*big.Int, error) {
	return _RelayerRegistration.Contract.RelayerCount(&_RelayerRegistration.CallOpts)
}

// RelayerCount is a free data retrieval call binding the contract method 0x87d340ab.
//
// Solidity: function RelayerCount() constant returns(uint256)
func (_RelayerRegistration *RelayerRegistrationCallerSession) RelayerCount() (*big.Int, error) {
	return _RelayerRegistration.Contract.RelayerCount(&_RelayerRegistration.CallOpts)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase(coinbase address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationCaller) GetRelayerByCoinbase(opts *bind.CallOpts, coinbase common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
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
	err := _RelayerRegistration.contract.Call(opts, out, "getRelayerByCoinbase", coinbase)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, err
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase(coinbase address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationSession) GetRelayerByCoinbase(coinbase common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByCoinbase(&_RelayerRegistration.CallOpts, coinbase)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase(coinbase address) constant returns(uint256, address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationCallerSession) GetRelayerByCoinbase(coinbase common.Address) (*big.Int, common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByCoinbase(&_RelayerRegistration.CallOpts, coinbase)
}

// BuyRelayer is a paid mutator transaction binding the contract method 0xe699df0e.
//
// Solidity: function buyRelayer(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) BuyRelayer(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "buyRelayer", coinbase)
}

// BuyRelayer is a paid mutator transaction binding the contract method 0xe699df0e.
//
// Solidity: function buyRelayer(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) BuyRelayer(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.BuyRelayer(&_RelayerRegistration.TransactOpts, coinbase)
}

// BuyRelayer is a paid mutator transaction binding the contract method 0xe699df0e.
//
// Solidity: function buyRelayer(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) BuyRelayer(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.BuyRelayer(&_RelayerRegistration.TransactOpts, coinbase)
}

// CancelSelling is a paid mutator transaction binding the contract method 0x5b673b1f.
//
// Solidity: function cancelSelling(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) CancelSelling(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "cancelSelling", coinbase)
}

// CancelSelling is a paid mutator transaction binding the contract method 0x5b673b1f.
//
// Solidity: function cancelSelling(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) CancelSelling(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.CancelSelling(&_RelayerRegistration.TransactOpts, coinbase)
}

// CancelSelling is a paid mutator transaction binding the contract method 0x5b673b1f.
//
// Solidity: function cancelSelling(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) CancelSelling(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.CancelSelling(&_RelayerRegistration.TransactOpts, coinbase)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) DepositMore(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "depositMore", coinbase)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) DepositMore(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.DepositMore(&_RelayerRegistration.TransactOpts, coinbase)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) DepositMore(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.DepositMore(&_RelayerRegistration.TransactOpts, coinbase)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x57ea3c41.
//
// Solidity: function reconfigure(maxRelayer uint256, maxToken uint256, minDeposit uint256) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Reconfigure(opts *bind.TransactOpts, maxRelayer *big.Int, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "reconfigure", maxRelayer, maxToken, minDeposit)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x57ea3c41.
//
// Solidity: function reconfigure(maxRelayer uint256, maxToken uint256, minDeposit uint256) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Reconfigure(maxRelayer *big.Int, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Reconfigure(&_RelayerRegistration.TransactOpts, maxRelayer, maxToken, minDeposit)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x57ea3c41.
//
// Solidity: function reconfigure(maxRelayer uint256, maxToken uint256, minDeposit uint256) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Reconfigure(maxRelayer *big.Int, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Reconfigure(&_RelayerRegistration.TransactOpts, maxRelayer, maxToken, minDeposit)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Refund(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "refund", coinbase)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Refund(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Refund(&_RelayerRegistration.TransactOpts, coinbase)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Refund(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Refund(&_RelayerRegistration.TransactOpts, coinbase)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Register(opts *bind.TransactOpts, coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "register", coinbase, tradeFee, fromTokens, toTokens)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Register(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Register(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Register(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Register(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Resign(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "resign", coinbase)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Resign(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Resign(&_RelayerRegistration.TransactOpts, coinbase)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(coinbase address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Resign(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Resign(&_RelayerRegistration.TransactOpts, coinbase)
}

// SellRelayer is a paid mutator transaction binding the contract method 0x87c6bbcd.
//
// Solidity: function sellRelayer(coinbase address, price uint256) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) SellRelayer(opts *bind.TransactOpts, coinbase common.Address, price *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "sellRelayer", coinbase, price)
}

// SellRelayer is a paid mutator transaction binding the contract method 0x87c6bbcd.
//
// Solidity: function sellRelayer(coinbase address, price uint256) returns()
func (_RelayerRegistration *RelayerRegistrationSession) SellRelayer(coinbase common.Address, price *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.SellRelayer(&_RelayerRegistration.TransactOpts, coinbase, price)
}

// SellRelayer is a paid mutator transaction binding the contract method 0x87c6bbcd.
//
// Solidity: function sellRelayer(coinbase address, price uint256) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) SellRelayer(coinbase common.Address, price *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.SellRelayer(&_RelayerRegistration.TransactOpts, coinbase, price)
}

// Transfer is a paid mutator transaction binding the contract method 0xba45b0b8.
//
// Solidity: function transfer(coinbase address, new_owner address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Transfer(opts *bind.TransactOpts, coinbase common.Address, new_owner common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "transfer", coinbase, new_owner)
}

// Transfer is a paid mutator transaction binding the contract method 0xba45b0b8.
//
// Solidity: function transfer(coinbase address, new_owner address) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Transfer(coinbase common.Address, new_owner common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Transfer(&_RelayerRegistration.TransactOpts, coinbase, new_owner)
}

// Transfer is a paid mutator transaction binding the contract method 0xba45b0b8.
//
// Solidity: function transfer(coinbase address, new_owner address) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Transfer(coinbase common.Address, new_owner common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Transfer(&_RelayerRegistration.TransactOpts, coinbase, new_owner)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Update(opts *bind.TransactOpts, coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "update", coinbase, tradeFee, fromTokens, toTokens)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Update(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Update(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(coinbase address, tradeFee uint16, fromTokens address[], toTokens address[]) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Update(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Update(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// RelayerRegistrationBuyEventIterator is returned from FilterBuyEvent and is used to iterate over the raw logs and unpacked data for BuyEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationBuyEventIterator struct {
	Event *RelayerRegistrationBuyEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationBuyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationBuyEvent)
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
		it.Event = new(RelayerRegistrationBuyEvent)
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
func (it *RelayerRegistrationBuyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationBuyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationBuyEvent represents a BuyEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationBuyEvent struct {
	Success  bool
	Coinbase common.Address
	Price    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBuyEvent is a free log retrieval operation binding the contract event 0x07e248a3b3d2184a9491c3b45089a6e15aac742b9d974e691e7beb0f6e7c58c6.
//
// Solidity: event BuyEvent(success bool, coinbase address, price uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterBuyEvent(opts *bind.FilterOpts) (*RelayerRegistrationBuyEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "BuyEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationBuyEventIterator{contract: _RelayerRegistration.contract, event: "BuyEvent", logs: logs, sub: sub}, nil
}

// WatchBuyEvent is a free log subscription operation binding the contract event 0x07e248a3b3d2184a9491c3b45089a6e15aac742b9d974e691e7beb0f6e7c58c6.
//
// Solidity: event BuyEvent(success bool, coinbase address, price uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchBuyEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationBuyEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "BuyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationBuyEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "BuyEvent", log); err != nil {
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

// RelayerRegistrationConfigEventIterator is returned from FilterConfigEvent and is used to iterate over the raw logs and unpacked data for ConfigEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationConfigEventIterator struct {
	Event *RelayerRegistrationConfigEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationConfigEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationConfigEvent)
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
		it.Event = new(RelayerRegistrationConfigEvent)
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
func (it *RelayerRegistrationConfigEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationConfigEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationConfigEvent represents a ConfigEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationConfigEvent struct {
	MaxRelayer *big.Int
	MaxToken   *big.Int
	MinDeposit *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterConfigEvent is a free log retrieval operation binding the contract event 0x8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d93923.
//
// Solidity: event ConfigEvent(max_relayer uint256, max_token uint256, min_deposit uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterConfigEvent(opts *bind.FilterOpts) (*RelayerRegistrationConfigEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "ConfigEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationConfigEventIterator{contract: _RelayerRegistration.contract, event: "ConfigEvent", logs: logs, sub: sub}, nil
}

// WatchConfigEvent is a free log subscription operation binding the contract event 0x8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d93923.
//
// Solidity: event ConfigEvent(max_relayer uint256, max_token uint256, min_deposit uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchConfigEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationConfigEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "ConfigEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationConfigEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "ConfigEvent", log); err != nil {
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

// RelayerRegistrationRefundEventIterator is returned from FilterRefundEvent and is used to iterate over the raw logs and unpacked data for RefundEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationRefundEventIterator struct {
	Event *RelayerRegistrationRefundEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationRefundEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationRefundEvent)
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
		it.Event = new(RelayerRegistrationRefundEvent)
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
func (it *RelayerRegistrationRefundEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationRefundEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationRefundEvent represents a RefundEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationRefundEvent struct {
	Success       bool
	RemainingTime *big.Int
	DepositAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRefundEvent is a free log retrieval operation binding the contract event 0xfaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d30250.
//
// Solidity: event RefundEvent(success bool, remaining_time uint256, deposit_amount uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterRefundEvent(opts *bind.FilterOpts) (*RelayerRegistrationRefundEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "RefundEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationRefundEventIterator{contract: _RelayerRegistration.contract, event: "RefundEvent", logs: logs, sub: sub}, nil
}

// WatchRefundEvent is a free log subscription operation binding the contract event 0xfaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d30250.
//
// Solidity: event RefundEvent(success bool, remaining_time uint256, deposit_amount uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchRefundEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationRefundEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "RefundEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationRefundEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "RefundEvent", log); err != nil {
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

// RelayerRegistrationRegisterEventIterator is returned from FilterRegisterEvent and is used to iterate over the raw logs and unpacked data for RegisterEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationRegisterEventIterator struct {
	Event *RelayerRegistrationRegisterEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationRegisterEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationRegisterEvent)
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
		it.Event = new(RelayerRegistrationRegisterEvent)
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
func (it *RelayerRegistrationRegisterEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationRegisterEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationRegisterEvent represents a RegisterEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationRegisterEvent struct {
	Deposit    *big.Int
	TradeFee   uint16
	FromTokens []common.Address
	ToTokens   []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegisterEvent is a free log retrieval operation binding the contract event 0xcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b.
//
// Solidity: event RegisterEvent(deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterRegisterEvent(opts *bind.FilterOpts) (*RelayerRegistrationRegisterEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "RegisterEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationRegisterEventIterator{contract: _RelayerRegistration.contract, event: "RegisterEvent", logs: logs, sub: sub}, nil
}

// WatchRegisterEvent is a free log subscription operation binding the contract event 0xcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b.
//
// Solidity: event RegisterEvent(deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchRegisterEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationRegisterEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "RegisterEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationRegisterEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "RegisterEvent", log); err != nil {
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

// RelayerRegistrationResignEventIterator is returned from FilterResignEvent and is used to iterate over the raw logs and unpacked data for ResignEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationResignEventIterator struct {
	Event *RelayerRegistrationResignEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationResignEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationResignEvent)
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
		it.Event = new(RelayerRegistrationResignEvent)
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
func (it *RelayerRegistrationResignEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationResignEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationResignEvent represents a ResignEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationResignEvent struct {
	DepositReleaseTime *big.Int
	DepositAmount      *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterResignEvent is a free log retrieval operation binding the contract event 0x2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae5.
//
// Solidity: event ResignEvent(deposit_release_time uint256, deposit_amount uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterResignEvent(opts *bind.FilterOpts) (*RelayerRegistrationResignEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "ResignEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationResignEventIterator{contract: _RelayerRegistration.contract, event: "ResignEvent", logs: logs, sub: sub}, nil
}

// WatchResignEvent is a free log subscription operation binding the contract event 0x2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae5.
//
// Solidity: event ResignEvent(deposit_release_time uint256, deposit_amount uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchResignEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationResignEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "ResignEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationResignEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "ResignEvent", log); err != nil {
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

// RelayerRegistrationSellEventIterator is returned from FilterSellEvent and is used to iterate over the raw logs and unpacked data for SellEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationSellEventIterator struct {
	Event *RelayerRegistrationSellEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationSellEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationSellEvent)
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
		it.Event = new(RelayerRegistrationSellEvent)
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
func (it *RelayerRegistrationSellEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationSellEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationSellEvent represents a SellEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationSellEvent struct {
	IsOnSale bool
	Coinbase common.Address
	Price    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSellEvent is a free log retrieval operation binding the contract event 0xdb3d5e65fcde89731529c01d62b87bab1c64471cffdd528fc1adbc1712b5d082.
//
// Solidity: event SellEvent(is_on_sale bool, coinbase address, price uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterSellEvent(opts *bind.FilterOpts) (*RelayerRegistrationSellEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "SellEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationSellEventIterator{contract: _RelayerRegistration.contract, event: "SellEvent", logs: logs, sub: sub}, nil
}

// WatchSellEvent is a free log subscription operation binding the contract event 0xdb3d5e65fcde89731529c01d62b87bab1c64471cffdd528fc1adbc1712b5d082.
//
// Solidity: event SellEvent(is_on_sale bool, coinbase address, price uint256)
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchSellEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationSellEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "SellEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationSellEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "SellEvent", log); err != nil {
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

// RelayerRegistrationTransferEventIterator is returned from FilterTransferEvent and is used to iterate over the raw logs and unpacked data for TransferEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationTransferEventIterator struct {
	Event *RelayerRegistrationTransferEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationTransferEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationTransferEvent)
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
		it.Event = new(RelayerRegistrationTransferEvent)
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
func (it *RelayerRegistrationTransferEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationTransferEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationTransferEvent represents a TransferEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationTransferEvent struct {
	Owner      common.Address
	Deposit    *big.Int
	TradeFee   uint16
	FromTokens []common.Address
	ToTokens   []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferEvent is a free log retrieval operation binding the contract event 0xc13ab794f75ba420a1f52192a8e35a2cf2c74ae31ed94f53f47ce7712011b662.
//
// Solidity: event TransferEvent(owner address, deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterTransferEvent(opts *bind.FilterOpts) (*RelayerRegistrationTransferEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "TransferEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationTransferEventIterator{contract: _RelayerRegistration.contract, event: "TransferEvent", logs: logs, sub: sub}, nil
}

// WatchTransferEvent is a free log subscription operation binding the contract event 0xc13ab794f75ba420a1f52192a8e35a2cf2c74ae31ed94f53f47ce7712011b662.
//
// Solidity: event TransferEvent(owner address, deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchTransferEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationTransferEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "TransferEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationTransferEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "TransferEvent", log); err != nil {
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

// RelayerRegistrationUpdateEventIterator is returned from FilterUpdateEvent and is used to iterate over the raw logs and unpacked data for UpdateEvent events raised by the RelayerRegistration contract.
type RelayerRegistrationUpdateEventIterator struct {
	Event *RelayerRegistrationUpdateEvent // Event containing the contract specifics and raw log

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
func (it *RelayerRegistrationUpdateEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayerRegistrationUpdateEvent)
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
		it.Event = new(RelayerRegistrationUpdateEvent)
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
func (it *RelayerRegistrationUpdateEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayerRegistrationUpdateEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayerRegistrationUpdateEvent represents a UpdateEvent event raised by the RelayerRegistration contract.
type RelayerRegistrationUpdateEvent struct {
	Deposit    *big.Int
	TradeFee   uint16
	FromTokens []common.Address
	ToTokens   []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdateEvent is a free log retrieval operation binding the contract event 0xcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f.
//
// Solidity: event UpdateEvent(deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterUpdateEvent(opts *bind.FilterOpts) (*RelayerRegistrationUpdateEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "UpdateEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationUpdateEventIterator{contract: _RelayerRegistration.contract, event: "UpdateEvent", logs: logs, sub: sub}, nil
}

// WatchUpdateEvent is a free log subscription operation binding the contract event 0xcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f.
//
// Solidity: event UpdateEvent(deposit uint256, tradeFee uint16, fromTokens address[], toTokens address[])
func (_RelayerRegistration *RelayerRegistrationFilterer) WatchUpdateEvent(opts *bind.WatchOpts, sink chan<- *RelayerRegistrationUpdateEvent) (event.Subscription, error) {

	logs, sub, err := _RelayerRegistration.contract.WatchLogs(opts, "UpdateEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayerRegistrationUpdateEvent)
				if err := _RelayerRegistration.contract.UnpackLog(event, "UpdateEvent", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820181fdc8a1e3a308513e71780ec652fd7655798259507ccfb5b18a1bb0e3c880e0029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
