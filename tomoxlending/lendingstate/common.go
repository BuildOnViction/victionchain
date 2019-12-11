package lendingstate

import (
	"encoding/json"
	"errors"
	"github.com/tomochain/tomochain/crypto"
	"math/big"
	"time"

	"github.com/tomochain/tomochain/common"
)

const (
	OrderCacheLimit = 10000
)

var (
	EmptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

var EmptyHash = common.Hash{}
var Zero = big.NewInt(0)
var EmptyOrderList = orderList{
	Volume: nil,
	Root:   EmptyHash,
}
var EmptyExchangeObject = exchangeObject{
	Nonce:   0,
	AskRoot: EmptyHash,
	BidRoot: EmptyHash,
}
var EmptyOrder = LendingItem{
	Quantity: Zero,
}

var (
	ErrInvalidSignature = errors.New("verify lending item: invalid signature")
	ErrInvalidInterest  = errors.New("verify lending item: invalid Interest")
	ErrInvalidQuantity  = errors.New("verify lending item: invalid quantity")
	ErrInvalidRelayer   = errors.New("verify lending item: invalid relayer")
	ErrInvalidOrderType = errors.New("verify lending item: unsupported order type")
	ErrInvalidOrderSide = errors.New("verify lending item: invalid order side")
	ErrInvalidStatus    = errors.New("verify lending item: invalid status")

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
	Interest  *big.Int    // Interest in native coin
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

type TxLendingBatch struct {
	Data      []*LendingItem
	Timestamp int64
	TxHash    common.Hash
}

type MatchingResult struct {
	Trades  []*LendingTrade
	Rejects []*LendingItem
}

func EncodeTxLendingBatch(txLendingBatch TxLendingBatch) ([]byte, error) {
	data, err := json.Marshal(txLendingBatch)
	if err != nil || data == nil {
		return []byte{}, err
	}
	return data, nil
}

func DecodeTxLendingBatch(data []byte) (TxLendingBatch, error) {
	txLendingResult := TxLendingBatch{}
	if err := json.Unmarshal(data, &txLendingResult); err != nil {
		return TxLendingBatch{}, err
	}
	return txLendingResult, nil
}

// use orderHash instead of orderId
// because both takerOrders don't have orderId
func GetLendingItemHistoryKey(lendingToken, collateralToken common.Address, lendingItemHash common.Hash) common.Hash {
	return crypto.Keccak256Hash(lendingToken.Bytes(), collateralToken.Bytes(), lendingItemHash.Bytes())
}

type LendingItemHistoryItem struct {
	TxHash       common.Hash
	FilledAmount *big.Int
	Status       string
	UpdatedAt    time.Time
}

// use alloc to prevent reference manipulation
func EmptyKey() []byte {
	key := make([]byte, common.HashLength)
	return key
}

// ToJSON : log json string
func ToJSON(object interface{}, args ...string) string {
	var str []byte
	if len(args) == 2 {
		str, _ = json.MarshalIndent(object, args[0], args[1])
	} else {
		str, _ = json.Marshal(object)
	}
	return string(str)
}

func Mul(x, y *big.Int) *big.Int {
	return big.NewInt(0).Mul(x, y)
}

func Div(x, y *big.Int) *big.Int {
	return big.NewInt(0).Div(x, y)
}

func Add(x, y *big.Int) *big.Int {
	return big.NewInt(0).Add(x, y)
}

func Sub(x, y *big.Int) *big.Int {
	return big.NewInt(0).Sub(x, y)
}

func Neg(x *big.Int) *big.Int {
	return big.NewInt(0).Neg(x)
}

func ToBigInt(s string) *big.Int {
	res := big.NewInt(0)
	res.SetString(s, 10)
	return res
}

func CloneBigInt(bigInt *big.Int) *big.Int {
	res := new(big.Int).SetBytes(bigInt.Bytes())
	return res
}

func Exp(x, y *big.Int) *big.Int {
	return big.NewInt(0).Exp(x, y, nil)
}

func Max(a, b *big.Int) *big.Int {
	if a.Cmp(b) == 1 {
		return a
	} else {
		return b
	}
}

func GetOrderBookHash(baseToken common.Address, quoteToken common.Address) common.Hash {
	return common.BytesToHash(append(baseToken[:16], quoteToken[4:]...))
}
