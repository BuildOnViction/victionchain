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

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RelayerRegistrationABI is the input ABI used to generate the binding from.
const RelayerRegistrationABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"MaximumRelayers\",\"outputs\":[{\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"maxRelayer\",\"type\":\"uint16\"},{\"name\":\"maxToken\",\"type\":\"uint256\"},{\"name\":\"minDeposit\",\"type\":\"uint256\"}],\"name\":\"reconfigure\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"depositMore\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"getRelayerByCoinbase\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint16\"},{\"name\":\"\",\"type\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"update\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"new_owner\",\"type\":\"address\"},{\"name\":\"new_coinbase\",\"type\":\"address\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"RelayerCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"resign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getRelayerByOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MinimumDeposit\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"},{\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"register\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MaximumTokenList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"coinbase\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CONTRACT_OWNER\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"maxRelayers\",\"type\":\"uint16\"},{\"name\":\"maxTokenList\",\"type\":\"uint256\"},{\"name\":\"minDeposit\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"max_relayer\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"max_token\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"min_deposit\",\"type\":\"uint256\"}],\"name\":\"ConfigEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"RegisterEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"UpdateEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tradeFee\",\"type\":\"uint16\"},{\"indexed\":false,\"name\":\"fromTokens\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"toTokens\",\"type\":\"address[]\"}],\"name\":\"TransferEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"deposit_release_time\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deposit_amount\",\"type\":\"uint256\"}],\"name\":\"ResignEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"remaining_time\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deposit_amount\",\"type\":\"uint256\"}],\"name\":\"RefundEvent\",\"type\":\"event\"}]"

// RelayerRegistrationFuncSigs maps the 4-byte function signature to its string representation.
var RelayerRegistrationFuncSigs = map[string]string{
	"fd301c49": "CONTRACT_OWNER()",
	"0e5c0fee": "MaximumRelayers()",
	"cfaece12": "MaximumTokenList()",
	"c635a9f2": "MinimumDeposit()",
	"87d340ab": "RelayerCount()",
	"4ce69bf5": "depositMore(address)",
	"540105c7": "getRelayerByCoinbase(address)",
	"bf4d79bd": "getRelayerByOwner(address)",
	"2704a032": "reconfigure(uint16,uint256,uint256)",
	"fa89401a": "refund(address)",
	"c6c71aed": "register(address,uint16,address[],address[])",
	"ae6e43f5": "resign(address)",
	"7138bc92": "transfer(address,address,address)",
	"56246b68": "update(address,uint16,address[],address[])",
}

// RelayerRegistrationBin is the compiled bytecode used for deploying new contracts.
var RelayerRegistrationBin = "0x608060405234801561001057600080fd5b50604051606080612265833981016040908152815160208301519190920151600060068190558054600193909355670de0b6b3a7640000909102600755600160a060020a031961ffff909316740100000000000000000000000000000000000000000260a060020a61ffff02199092169190911791909116331790556121ca8061009b6000396000f3006080604052600436106100cf5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630e5c0fee81146100d45780632704a032146101005780634ce69bf514610124578063540105c71461013857806356246b68146102215780637138bc92146102c557806387d340ab146102f2578063ae6e43f514610319578063bf4d79bd1461033a578063c635a9f2146103ab578063c6c71aed146103c0578063cfaece1214610457578063fa89401a1461046c578063fd301c491461048d575b600080fd5b3480156100e057600080fd5b506100e96104be565b6040805161ffff9092168252519081900360200190f35b34801561010c57600080fd5b5061012261ffff600435166024356044356104e0565b005b610122600160a060020a0360043516610629565b34801561014457600080fd5b50610159600160a060020a0360043516610915565b6040518086600160a060020a0316600160a060020a031681526020018581526020018461ffff1661ffff1681526020018060200180602001838103835285818151815260200191508051906020019060200280838360005b838110156101c95781810151838201526020016101b1565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156102085781810151838201526020016101f0565b5050505090500197505050505050505060405180910390f35b34801561022d57600080fd5b506040805160206004604435818101358381028086018501909652808552610122958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a998901989297509082019550935083925085019084908082843750949750610a249650505050505050565b3480156102d157600080fd5b50610122600160a060020a0360043581169060243581169060443516610dc7565b3480156102fe57600080fd5b506103076114c5565b60408051918252519081900360200190f35b34801561032557600080fd5b50610122600160a060020a03600435166114cb565b34801561034657600080fd5b5061035b600160a060020a0360043516611698565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561039757818101518382015260200161037f565b505050509050019250505060405180910390f35b3480156103b757600080fd5b5061030761170e565b6040805160206004604435818101358381028086018501909652808552610122958335600160a060020a0316956024803561ffff1696369695606495939492019291829185019084908082843750506040805187358901803560208181028481018201909552818452989b9a9989019892975090820195509350839250850190849080828437509497506117149650505050505050565b34801561046357600080fd5b50610307611d7e565b34801561047857600080fd5b50610122600160a060020a0360043516611d84565b34801561049957600080fd5b506104a261207b565b60408051600160a060020a039092168252519081900360200190f35b60005474010000000000000000000000000000000000000000900461ffff1681565b60008054600160a060020a03163314610543576040805160e560020a62461bcd02815260206004820152601460248201527f436f6e7472616374204f776e6572204f6e6c792e000000000000000000000000604482015290519081900360640190fd5b60065461ffff85161161055557600080fd5b60048311801561056657506103e983105b151561057157600080fd5b612710821161057f57600080fd5b506000805461ffff8086167401000000000000000000000000000000000000000090810275ffff00000000000000000000000000000000000000001990931692909217928390556001859055670de0b6b3a7640000848102600781905560408051949095049092168352602083018690528284019190915291517f8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d939239181900360600190a150505050565b600160a060020a03818116600090815260036020526040902054829116331461068a576040805160e560020a62461bcd028152602060048201526013602482015260008051602061217f833981519152604482015290519081900360640190fd5b600160a060020a038216600090815260056020526040902054829015610720576040805160e560020a62461bcd02815260206004820152602860248201527f5468652072656c6179657220686173206265656e20726571756573746564207460448201527f6f20636c6f73652e000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60003411610778576040805160e560020a62461bcd02815260206004820152601a60248201527f5472616e736665722076616c7565206d757374206265203e2030000000000000604482015290519081900360640190fd5b670de0b6b3a76400003410156107fe576040805160e560020a62461bcd02815260206004820152603160248201527f4174206c65617374203120544f4d4f20697320726571756972656420666f722060448201527f61206465706f7369742072657175657374000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a038316600090815260026020818152604092839020805434018082556001820154855182815261ffff909116938101849052608095810186815294830180549682018790527fcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f969295909360030192606083019060a0840190869080156108b557602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610897575b505083810382528481815481526020019150805480156108fe57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116108e0575b5050965050505050505060405180910390a1505050565b600160a060020a03808216600090815260036020818152604080842054600280845282862080546001820154928201805486518189028101890190975280875298998a998a996060998a999890931697949661ffff16959394930192918491908301828280156109ae57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610990575b5050505050915080805480602002602001604051908101604052809291908181526020018280548015610a0a57602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116109ec575b505050505090509450945094509450945091939590929450565b600160a060020a038481166000908152600360205260409020548591163314610a85576040805160e560020a62461bcd028152602060048201526013602482015260008051602061217f833981519152604482015290519081900360640190fd5b600160a060020a038516600090815260056020526040902054859015610b1b576040805160e560020a62461bcd02815260206004820152602860248201527f5468652072656c6179657220686173206265656e20726571756573746564207460448201527f6f20636c6f73652e000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60018561ffff1610158015610b3557506127108561ffff16105b1515610b8b576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c6964204d616b657220466565000000000000000000000000000000604482015290519081900360640190fd5b60015484511115610be6576040805160e560020a62461bcd02815260206004820152601f60248201527f457863656564696e67206e756d626572206f6620747261646520706169727300604482015290519081900360640190fd5b8351835114610c3f576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f6620506169727300000000000000604482015290519081900360640190fd5b600160a060020a038616600090815260026020818152604090922060018101805461ffff191661ffff8a161790558651610c819391909201919087019061208a565b50600160a060020a03861660009081526002602090815260409091208451610cb19260039092019186019061208a565b50600160a060020a03861660009081526002602081815260409283902080546001820154855182815261ffff909116938101849052608095810186815294830180549682018790527fcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f969295909360030192606083019060a084019086908015610d6457602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610d46575b50508381038252848181548152602001915080548015610dad57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610d8f575b5050965050505050505060405180910390a1505050505050565b6000610dd16120ef565b600160a060020a038581166000908152600360205260409020548691163314610e32576040805160e560020a62461bcd028152602060048201526013602482015260008051602061217f833981519152604482015290519081900360640190fd5b600160a060020a038616600090815260056020526040902054869015610ec8576040805160e560020a62461bcd02815260206004820152602860248201527f5468652072656c6179657220686173206265656e20726571756573746564207460448201527f6f20636c6f73652e000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a03861615801590610ee95750600160a060020a0386163314155b1515610ef457600080fd5b600160a060020a03861660009081526002602052604090206001015461ffff1615610f8f576040805160e560020a62461bcd02815260206004820152603c60248201527f4f776e65722061646472657373206d757374206e6f742062652063757272656e60448201527f746c7920757365642061732072656c617965722d636f696e6261736500000000606482015290519081900360840190fd5b600160a060020a0385161515610fa457600080fd5b600054600160a060020a0386811691161415610fbf57600080fd5b600160a060020a038581169088161461110157600160a060020a03851660009081526002602052604090206001015461ffff161561106d576040805160e560020a62461bcd02815260206004820152602360248201527f546865206e657720636f696e6261736520697320616c726561647920696e207560448201527f7365640000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a03851660009081526004602052604090205415611101576040805160e560020a62461bcd02815260206004820152602b60248201527f546865206e657720636f696e626173652069732075736564206173206120526560448201527f6c617965722d6f776e6572000000000000000000000000000000000000000000606482015290519081900360840190fd5b600093505b336000908152600460205260409020548410156114bc573360009081526004602052604090208054600160a060020a03891691908690811061114457fe5b600091825260209091200154600160a060020a031614156114b157600160a060020a038716600090815260026020818152604092839020835160808101855281548152600182015461ffff16818401529281018054855181850281018501875281815294959294928601938301828280156111e857602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116111ca575b505050505081526020016003820180548060200260200160405190810160405280929190818152602001828054801561124a57602002820191906000526020600020905b8154600160a060020a0316815260019091019060200180831161122c575b505050919092525050336000908152600460205260409020805492955091869150811061127357fe5b60009182526020909120018054600160a060020a0319169055600160a060020a038581169088161461135a57600160a060020a0387166000908152600260208190526040822082815560018101805461ffff1916905591906112d79083018261211c565b6112e560038301600061211c565b5050600160a060020a038516600090815260026020818152604092839020865181558682015160018201805461ffff191661ffff90921691909117905592860151805187949361133a9390850192019061208a565b506060820151805161135691600384019160209091019061208a565b5050505b600160a060020a0380861660008181526003602081815260408084208054968d16600160a060020a031997881681179091558452600482528084208054600181810183559186528386200180549097168617909655938352600280825292849020805495810154855187815261ffff909116928101839052608095810186815294820180549682018790527fccbab7f516e3706c2a308cdde87979595b62b4b2079cdf7141926bc256467e3b9796939590949290920192909190606083019060a08401908690801561145557602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611437575b5050838103825284818154815260200191508054801561149e57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611480575b5050965050505050505060405180910390a15b600190930192611106565b50505050505050565b60065481565b600160a060020a03818116600090815260036020526040902054829116331461152c576040805160e560020a62461bcd028152602060048201526013602482015260008051602061217f833981519152604482015290519081900360640190fd5b600160a060020a038216600090815260026020526040812054116115c0576040805160e560020a62461bcd02815260206004820152602760248201527f4e6f2072656c61796572206173736f636961746564207769746820746869732060448201527f6164647265737300000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a0382166000908152600560205260409020541561162e576040805160e560020a62461bcd02815260206004820152601860248201527f5265717565737420616c72656164792072656365697665640000000000000000604482015290519081900360640190fd5b600160a060020a03821660009081526005602090815260408083206224ea0042019081905560028352928190205481519384529183019190915280517f2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae59281900390910190a15050565b600160a060020a03811660009081526004602090815260409182902080548351818402810184019094528084526060939283018282801561170257602002820191906000526020600020905b8154600160a060020a031681526001909101906020018083116116e4575b50505050509050919050565b60075481565b61171c6120ef565b600054600160a060020a03163314156117a5576040805160e560020a62461bcd02815260206004820152602f60248201527f436f6e7472616374204f776e657220697320666f7262696464656e20746f206360448201527f726561746520612052656c617965720000000000000000000000000000000000606482015290519081900360840190fd5b33600160a060020a038616141561182c576040805160e560020a62461bcd02815260206004820152603660248201527f436f696e6261736520616e642052656c617965724f776e65722061646472657360448201527f73206d757374206e6f74206265207468652073616d6500000000000000000000606482015290519081900360840190fd5b600054600160a060020a03868116911614156118b8576040805160e560020a62461bcd02815260206004820152602b60248201527f436f696e62617365206d757374206e6f742062652073616d6520617320434f4e60448201527f54524143545f4f574e4552000000000000000000000000000000000000000000606482015290519081900360840190fd5b600754341015611912576040805160e560020a62461bcd02815260206004820152601e60248201527f4d696e696d756d206465706f736974206e6f74207361746973666965642e0000604482015290519081900360640190fd5b60018461ffff161015801561192c57506127108461ffff16105b1515611982576040805160e560020a62461bcd02815260206004820152601160248201527f496e76616c6964204d616b657220466565000000000000000000000000000000604482015290519081900360640190fd5b600154835111156119dd576040805160e560020a62461bcd02815260206004820152601f60248201527f457863656564696e67206e756d626572206f6620747261646520706169727300604482015290519081900360640190fd5b8251825114611a36576040805160e560020a62461bcd02815260206004820152601960248201527f4e6f742076616c6964206e756d626572206f6620506169727300000000000000604482015290519081900360640190fd5b600160a060020a03851660009081526002602052604090205415611aa4576040805160e560020a62461bcd02815260206004820152601c60248201527f436f696e6261736520616c726561647920726567697374657265642e00000000604482015290519081900360640190fd5b600160a060020a038581166000908152600360205260409020541615611b14576040805160e560020a62461bcd02815260206004820152601b60248201527f436f696e6261736520616c726561647920726567697374657265640000000000604482015290519081900360640190fd5b6000546006547401000000000000000000000000000000000000000090910461ffff1611611b8c576040805160e560020a62461bcd02815260206004820152601b60248201527f4d6178696d756d2072656c617965727320726567697374657265640000000000604482015290519081900360640190fd5b506040805160808101825234815261ffff858116602080840191825283850187815260608501879052600160a060020a038a166000908152600280845296902085518155925160018401805461ffff1916919095161790935591518051939485949293611bfe9391850192019061208a565b5060608201518051611c1a91600384019160209091019061208a565b505050600160a060020a0385166000818152600360208181526040808420805433600160a060020a031991821681179092559085526004835281852080546001818101835591875284872001805490921687179091556006805482019055948452600280835293819020805495810154825187815261ffff909116938101849052608092810183815295820180549382018490527fcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b979694959094920192909190606083019060a084019086908015611d1c57602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611cfe575b50508381038252848181548152602001915080548015611d6557602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311611d47575b5050965050505050505060405180910390a15050505050565b60015481565b600160a060020a03818116600090815260036020526040812054909182918491163314611de9576040805160e560020a62461bcd028152602060048201526013602482015260008051602061217f833981519152604482015290519081900360640190fd5b600160a060020a03841660009081526005602052604081205411611e57576040805160e560020a62461bcd02815260206004820152601160248201527f52657175657374206e6f7420666f756e64000000000000000000000000000000604482015290519081900360640190fd5b600160a060020a03841660009081526002602090815260408083205460059092529091205490935042111561201857600091505b33600090815260046020526040902054821015612013573360009081526004602052604090208054600160a060020a038616919084908110611ec957fe5b600091825260209091200154600160a060020a0316141561200857336000908152600460205260409020805483908110611eff57fe5b600091825260208083209091018054600160a060020a0319169055600160a060020a03861682526002908190526040822082815560018101805461ffff191690559190611f4e9083018261211c565b611f5c60038301600061211c565b5050600160a060020a03841660009081526003602090815260408083208054600160a060020a031916905560059091528082208290556006805460001901905551339185156108fc02918691818181858888f19350505050158015611fc5573d6000803e3d6000fd5b5060408051600181526000602082015280820185905290517ffaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d302509181900360600190a15b600190910190611e8b565b612075565b600160a060020a038416600090815260056020908152604080832054815193845242900391830191909152818101859052517ffaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d302509181900360600190a15b50505050565b600054600160a060020a031681565b8280548282559060005260206000209081019282156120df579160200282015b828111156120df5782518254600160a060020a031916600160a060020a039091161782556020909201916001909101906120aa565b506120eb92915061213d565b5090565b60806040519081016040528060008152602001600061ffff16815260200160608152602001606081525090565b508054600082559060005260206000209081019061213a9190612164565b50565b61216191905b808211156120eb578054600160a060020a0319168155600101612143565b90565b61216191905b808211156120eb576000815560010161216a560052656c61796572204f776e6572204f6e6c792e00000000000000000000000000a165627a7a723058201be4d6ef450bff4a436da0267e998f6f2fd66c93bcb66c3887186ec9dde674890029"

// DeployRelayerRegistration deploys a new Ethereum contract, binding an instance of RelayerRegistration to it.
func DeployRelayerRegistration(auth *bind.TransactOpts, backend bind.ContractBackend, maxRelayers uint16, maxTokenList *big.Int, minDeposit *big.Int) (common.Address, *types.Transaction, *RelayerRegistration, error) {
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
// Solidity: function MaximumRelayers() constant returns(uint16)
func (_RelayerRegistration *RelayerRegistrationCaller) MaximumRelayers(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "MaximumRelayers")
	return *ret0, err
}

// MaximumRelayers is a free data retrieval call binding the contract method 0x0e5c0fee.
//
// Solidity: function MaximumRelayers() constant returns(uint16)
func (_RelayerRegistration *RelayerRegistrationSession) MaximumRelayers() (uint16, error) {
	return _RelayerRegistration.Contract.MaximumRelayers(&_RelayerRegistration.CallOpts)
}

// MaximumRelayers is a free data retrieval call binding the contract method 0x0e5c0fee.
//
// Solidity: function MaximumRelayers() constant returns(uint16)
func (_RelayerRegistration *RelayerRegistrationCallerSession) MaximumRelayers() (uint16, error) {
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
// Solidity: function getRelayerByCoinbase(address coinbase) constant returns(address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationCaller) GetRelayerByCoinbase(opts *bind.CallOpts, coinbase common.Address) (common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(uint16)
		ret3 = new([]common.Address)
		ret4 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _RelayerRegistration.contract.Call(opts, out, "getRelayerByCoinbase", coinbase)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase(address coinbase) constant returns(address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationSession) GetRelayerByCoinbase(coinbase common.Address) (common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByCoinbase(&_RelayerRegistration.CallOpts, coinbase)
}

// GetRelayerByCoinbase is a free data retrieval call binding the contract method 0x540105c7.
//
// Solidity: function getRelayerByCoinbase(address coinbase) constant returns(address, uint256, uint16, address[], address[])
func (_RelayerRegistration *RelayerRegistrationCallerSession) GetRelayerByCoinbase(coinbase common.Address) (common.Address, *big.Int, uint16, []common.Address, []common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByCoinbase(&_RelayerRegistration.CallOpts, coinbase)
}

// GetRelayerByOwner is a free data retrieval call binding the contract method 0xbf4d79bd.
//
// Solidity: function getRelayerByOwner(address owner) constant returns(address[])
func (_RelayerRegistration *RelayerRegistrationCaller) GetRelayerByOwner(opts *bind.CallOpts, owner common.Address) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _RelayerRegistration.contract.Call(opts, out, "getRelayerByOwner", owner)
	return *ret0, err
}

// GetRelayerByOwner is a free data retrieval call binding the contract method 0xbf4d79bd.
//
// Solidity: function getRelayerByOwner(address owner) constant returns(address[])
func (_RelayerRegistration *RelayerRegistrationSession) GetRelayerByOwner(owner common.Address) ([]common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByOwner(&_RelayerRegistration.CallOpts, owner)
}

// GetRelayerByOwner is a free data retrieval call binding the contract method 0xbf4d79bd.
//
// Solidity: function getRelayerByOwner(address owner) constant returns(address[])
func (_RelayerRegistration *RelayerRegistrationCallerSession) GetRelayerByOwner(owner common.Address) ([]common.Address, error) {
	return _RelayerRegistration.Contract.GetRelayerByOwner(&_RelayerRegistration.CallOpts, owner)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) DepositMore(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "depositMore", coinbase)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationSession) DepositMore(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.DepositMore(&_RelayerRegistration.TransactOpts, coinbase)
}

// DepositMore is a paid mutator transaction binding the contract method 0x4ce69bf5.
//
// Solidity: function depositMore(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) DepositMore(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.DepositMore(&_RelayerRegistration.TransactOpts, coinbase)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x2704a032.
//
// Solidity: function reconfigure(uint16 maxRelayer, uint256 maxToken, uint256 minDeposit) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Reconfigure(opts *bind.TransactOpts, maxRelayer uint16, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "reconfigure", maxRelayer, maxToken, minDeposit)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x2704a032.
//
// Solidity: function reconfigure(uint16 maxRelayer, uint256 maxToken, uint256 minDeposit) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Reconfigure(maxRelayer uint16, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Reconfigure(&_RelayerRegistration.TransactOpts, maxRelayer, maxToken, minDeposit)
}

// Reconfigure is a paid mutator transaction binding the contract method 0x2704a032.
//
// Solidity: function reconfigure(uint16 maxRelayer, uint256 maxToken, uint256 minDeposit) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Reconfigure(maxRelayer uint16, maxToken *big.Int, minDeposit *big.Int) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Reconfigure(&_RelayerRegistration.TransactOpts, maxRelayer, maxToken, minDeposit)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Refund(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "refund", coinbase)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Refund(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Refund(&_RelayerRegistration.TransactOpts, coinbase)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Refund(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Refund(&_RelayerRegistration.TransactOpts, coinbase)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Register(opts *bind.TransactOpts, coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "register", coinbase, tradeFee, fromTokens, toTokens)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Register(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Register(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Register is a paid mutator transaction binding the contract method 0xc6c71aed.
//
// Solidity: function register(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Register(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Register(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Resign(opts *bind.TransactOpts, coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "resign", coinbase)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Resign(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Resign(&_RelayerRegistration.TransactOpts, coinbase)
}

// Resign is a paid mutator transaction binding the contract method 0xae6e43f5.
//
// Solidity: function resign(address coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Resign(coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Resign(&_RelayerRegistration.TransactOpts, coinbase)
}

// Transfer is a paid mutator transaction binding the contract method 0x7138bc92.
//
// Solidity: function transfer(address coinbase, address new_owner, address new_coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Transfer(opts *bind.TransactOpts, coinbase common.Address, new_owner common.Address, new_coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "transfer", coinbase, new_owner, new_coinbase)
}

// Transfer is a paid mutator transaction binding the contract method 0x7138bc92.
//
// Solidity: function transfer(address coinbase, address new_owner, address new_coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Transfer(coinbase common.Address, new_owner common.Address, new_coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Transfer(&_RelayerRegistration.TransactOpts, coinbase, new_owner, new_coinbase)
}

// Transfer is a paid mutator transaction binding the contract method 0x7138bc92.
//
// Solidity: function transfer(address coinbase, address new_owner, address new_coinbase) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Transfer(coinbase common.Address, new_owner common.Address, new_coinbase common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Transfer(&_RelayerRegistration.TransactOpts, coinbase, new_owner, new_coinbase)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationTransactor) Update(opts *bind.TransactOpts, coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.contract.Transact(opts, "update", coinbase, tradeFee, fromTokens, toTokens)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationSession) Update(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Update(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
}

// Update is a paid mutator transaction binding the contract method 0x56246b68.
//
// Solidity: function update(address coinbase, uint16 tradeFee, address[] fromTokens, address[] toTokens) returns()
func (_RelayerRegistration *RelayerRegistrationTransactorSession) Update(coinbase common.Address, tradeFee uint16, fromTokens []common.Address, toTokens []common.Address) (*types.Transaction, error) {
	return _RelayerRegistration.Contract.Update(&_RelayerRegistration.TransactOpts, coinbase, tradeFee, fromTokens, toTokens)
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
// Solidity: event ConfigEvent(uint256 max_relayer, uint256 max_token, uint256 min_deposit)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterConfigEvent(opts *bind.FilterOpts) (*RelayerRegistrationConfigEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "ConfigEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationConfigEventIterator{contract: _RelayerRegistration.contract, event: "ConfigEvent", logs: logs, sub: sub}, nil
}

// WatchConfigEvent is a free log subscription operation binding the contract event 0x8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d93923.
//
// Solidity: event ConfigEvent(uint256 max_relayer, uint256 max_token, uint256 min_deposit)
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

// ParseConfigEvent is a log parse operation binding the contract event 0x8f6bd709a98381db4e403a67ba106d598972dad177e946f19b54777f54d93923.
//
// Solidity: event ConfigEvent(uint256 max_relayer, uint256 max_token, uint256 min_deposit)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseConfigEvent(log types.Log) (*RelayerRegistrationConfigEvent, error) {
	event := new(RelayerRegistrationConfigEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "ConfigEvent", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event RefundEvent(bool success, uint256 remaining_time, uint256 deposit_amount)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterRefundEvent(opts *bind.FilterOpts) (*RelayerRegistrationRefundEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "RefundEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationRefundEventIterator{contract: _RelayerRegistration.contract, event: "RefundEvent", logs: logs, sub: sub}, nil
}

// WatchRefundEvent is a free log subscription operation binding the contract event 0xfaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d30250.
//
// Solidity: event RefundEvent(bool success, uint256 remaining_time, uint256 deposit_amount)
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

// ParseRefundEvent is a log parse operation binding the contract event 0xfaba1aac53309af4c1c439f38c29500d3828405ee1ca5e7641b0432d17d30250.
//
// Solidity: event RefundEvent(bool success, uint256 remaining_time, uint256 deposit_amount)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseRefundEvent(log types.Log) (*RelayerRegistrationRefundEvent, error) {
	event := new(RelayerRegistrationRefundEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "RefundEvent", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event RegisterEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterRegisterEvent(opts *bind.FilterOpts) (*RelayerRegistrationRegisterEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "RegisterEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationRegisterEventIterator{contract: _RelayerRegistration.contract, event: "RegisterEvent", logs: logs, sub: sub}, nil
}

// WatchRegisterEvent is a free log subscription operation binding the contract event 0xcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b.
//
// Solidity: event RegisterEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
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

// ParseRegisterEvent is a log parse operation binding the contract event 0xcf24380d990b0bb3dd21518926bca48f81495ac131ee92655696db28c43b1b1b.
//
// Solidity: event RegisterEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseRegisterEvent(log types.Log) (*RelayerRegistrationRegisterEvent, error) {
	event := new(RelayerRegistrationRegisterEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "RegisterEvent", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event ResignEvent(uint256 deposit_release_time, uint256 deposit_amount)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterResignEvent(opts *bind.FilterOpts) (*RelayerRegistrationResignEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "ResignEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationResignEventIterator{contract: _RelayerRegistration.contract, event: "ResignEvent", logs: logs, sub: sub}, nil
}

// WatchResignEvent is a free log subscription operation binding the contract event 0x2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae5.
//
// Solidity: event ResignEvent(uint256 deposit_release_time, uint256 deposit_amount)
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

// ParseResignEvent is a log parse operation binding the contract event 0x2e821a4329d6351a6b13fe0c12fd7674cd0f4a2283685a4713e1325f36415ae5.
//
// Solidity: event ResignEvent(uint256 deposit_release_time, uint256 deposit_amount)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseResignEvent(log types.Log) (*RelayerRegistrationResignEvent, error) {
	event := new(RelayerRegistrationResignEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "ResignEvent", log); err != nil {
		return nil, err
	}
	return event, nil
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
	Deposit    *big.Int
	TradeFee   uint16
	FromTokens []common.Address
	ToTokens   []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferEvent is a free log retrieval operation binding the contract event 0xccbab7f516e3706c2a308cdde87979595b62b4b2079cdf7141926bc256467e3b.
//
// Solidity: event TransferEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterTransferEvent(opts *bind.FilterOpts) (*RelayerRegistrationTransferEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "TransferEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationTransferEventIterator{contract: _RelayerRegistration.contract, event: "TransferEvent", logs: logs, sub: sub}, nil
}

// WatchTransferEvent is a free log subscription operation binding the contract event 0xccbab7f516e3706c2a308cdde87979595b62b4b2079cdf7141926bc256467e3b.
//
// Solidity: event TransferEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
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

// ParseTransferEvent is a log parse operation binding the contract event 0xccbab7f516e3706c2a308cdde87979595b62b4b2079cdf7141926bc256467e3b.
//
// Solidity: event TransferEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseTransferEvent(log types.Log) (*RelayerRegistrationTransferEvent, error) {
	event := new(RelayerRegistrationTransferEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "TransferEvent", log); err != nil {
		return nil, err
	}
	return event, nil
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
// Solidity: event UpdateEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) FilterUpdateEvent(opts *bind.FilterOpts) (*RelayerRegistrationUpdateEventIterator, error) {

	logs, sub, err := _RelayerRegistration.contract.FilterLogs(opts, "UpdateEvent")
	if err != nil {
		return nil, err
	}
	return &RelayerRegistrationUpdateEventIterator{contract: _RelayerRegistration.contract, event: "UpdateEvent", logs: logs, sub: sub}, nil
}

// WatchUpdateEvent is a free log subscription operation binding the contract event 0xcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f.
//
// Solidity: event UpdateEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
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

// ParseUpdateEvent is a log parse operation binding the contract event 0xcaa8c94daf6ecfd00518cea95158f5273730574cca907eb0cd47e50732314c4f.
//
// Solidity: event UpdateEvent(uint256 deposit, uint16 tradeFee, address[] fromTokens, address[] toTokens)
func (_RelayerRegistration *RelayerRegistrationFilterer) ParseUpdateEvent(log types.Log) (*RelayerRegistrationUpdateEvent, error) {
	event := new(RelayerRegistrationUpdateEvent)
	if err := _RelayerRegistration.contract.UnpackLog(event, "UpdateEvent", log); err != nil {
		return nil, err
	}
	return event, nil
}
