package tomox_state

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	EmptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	Ask       = "SELL"
	Bid       = "BUY"
	Market    = "MO"
	Limit     = "LO"
	Cancel    = "CANCELLED"
)

var EmptyHash = common.Hash{}
var Zero = big.NewInt(0)
var EmptyOrderList = orderList{
	Volume: nil,
	Root:   EmptyHash,
}
var EmptyExchangeOnject = exchangeObject{
	Nonce:   0,
	AskRoot: EmptyHash,
	BidRoot: EmptyHash,
}
var EmptyOrder = OrderItem{
	Quantity: Zero,
}

var (
	ErrWrongHash             = errors.New("verify order: wrong hash")
	ErrInvalidSignature      = errors.New("verify order: invalid signature")
	ErrInvalidPrice          = errors.New("verify order: invalid price")
	ErrInvalidQuantity       = errors.New("verify order: invalid quantity")
	ErrInvalidRelayer        = errors.New("verify order: invalid relayer")
	ErrInvalidOrderType      = errors.New("verify order: unsupported order type")
	ErrInvalidOrderSide      = errors.New("verify order: invalid order side")
	ErrOrderBookHashNotMatch = errors.New("verify order: orderbook hash not match")
	ErrOrderTreeHashNotMatch = errors.New("verify order: ordertree hash not match")

	// supported order types
	MatchingOrderType = map[string]bool{
		Market: true,
		Limit:  true,
	}
)

// exchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type orderList struct {
	Volume *big.Int
	Root   common.Hash // merkle root of the storage trie
}

// exchangeObject is the Ethereum consensus representation of exchanges.
// These objects are stored in the main orderId trie.
type exchangeObject struct {
	Nonce     uint64
	Price     *big.Int    // price in native coin
	AskRoot   common.Hash // merkle root of the storage trie
	BidRoot   common.Hash // merkle root of the storage trie
	OrderRoot common.Hash
}

var (
	TokenMappingSlot = map[string]uint64{
		"balances": 0,
	}
	RelayerMappingSlot = map[string]uint64{
		"CONTRACT_OWNER":       0,
		"MaximumRelayers":      1,
		"MaximumTokenList":     2,
		"RELAYER_LIST":         3,
		"RELAYER_COINBASES":    4,
		"RESIGN_REQUESTS":      5,
		"RELAYER_ON_SALE_LIST": 6,
		"RelayerCount":         7,
		"MinimumDeposit":       8,
	}
	RelayerStructMappingSlot = map[string]*big.Int{
		"_deposit":    big.NewInt(0),
		"_fee":        big.NewInt(1),
		"_fromTokens": big.NewInt(2),
		"_toTokens":   big.NewInt(3),
		"_index":      big.NewInt(4),
		"_owner":      big.NewInt(5),
	}
)
