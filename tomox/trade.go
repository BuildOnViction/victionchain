package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/utils/math"
	"math/big"
	"time"
)

const (
	TradeStatusPending = "PENDING"
	TradeStatusSuccess = "SUCCESS"
	TradeStatusError   = "ERROR"
)

type Trade struct {
	ID             bson.ObjectId  `json:"id,omitempty" bson:"_id"`
	Taker          common.Address `json:"taker" bson:"taker"`
	Maker          common.Address `json:"maker" bson:"maker"`
	BaseToken      common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken     common.Address `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash common.Hash    `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash common.Hash    `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           common.Hash    `json:"hash" bson:"hash"`
	TxHash         common.Hash    `json:"txHash" bson:"txHash"`
	PairName       string         `json:"pairName" bson:"pairName"`
	PricePoint     *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Amount         *big.Int       `json:"amount" bson:"amount"`
	MakeFee        *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee        *big.Int       `json:"takeFee" bson:"takeFee"`
	Status         string         `json:"status" bson:"status"`
	CreatedAt      time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt" bson:"updatedAt"`
	TakerOrderSide string         `json:"takerOrderSide" bson:"takerOrderSide"`
}

type TradeBSON struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Taker          string        `json:"taker" bson:"taker"`
	Maker          string        `json:"maker" bson:"maker"`
	BaseToken      string        `json:"baseToken" bson:"baseToken"`
	QuoteToken     string        `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash string        `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash string        `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           string        `json:"hash" bson:"hash"`
	TxHash         string        `json:"txHash" bson:"txHash"`
	PairName       string        `json:"pairName" bson:"pairName"`
	Amount         string        `json:"amount" bson:"amount"`
	MakeFee        string        `json:"makeFee" bson:"makeFee"`
	TakeFee        string        `json:"takeFee" bson:"takeFee"`
	PricePoint     string        `json:"pricepoint" bson:"pricepoint"`
	Status         string        `json:"status" bson:"status"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	TakerOrderSide string        `json:"takerOrderSide" bson:"takerOrderSide"`
}


func (t *Trade) GetBSON() (interface{}, error) {
	tr := TradeBSON{
		ID:             t.ID,
		PairName:       t.PairName,
		Maker:          t.Maker.Hex(),
		Taker:          t.Taker.Hex(),
		BaseToken:      t.BaseToken.Hex(),
		QuoteToken:     t.QuoteToken.Hex(),
		MakerOrderHash: t.MakerOrderHash.Hex(),
		Hash:           t.Hash.Hex(),
		TxHash:         t.TxHash.Hex(),
		TakerOrderHash: t.TakerOrderHash.Hex(),
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		PricePoint:     t.PricePoint.String(),
		Status:         t.Status,
		Amount:         t.Amount.String(),
		MakeFee:        t.MakeFee.String(),
		TakeFee:        t.TakeFee.String(),
		TakerOrderSide: t.TakerOrderSide,
	}

	return tr, nil
}

func (t *Trade) SetBSON(raw bson.Raw) error {
	decoded := &TradeBSON{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.ID = decoded.ID
	t.PairName = decoded.PairName
	t.Taker = common.HexToAddress(decoded.Taker)
	t.Maker = common.HexToAddress(decoded.Maker)
	t.BaseToken = common.HexToAddress(decoded.BaseToken)
	t.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	t.MakerOrderHash = common.HexToHash(decoded.MakerOrderHash)
	t.TakerOrderHash = common.HexToHash(decoded.TakerOrderHash)
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.Status = decoded.Status
	t.Amount = math.ToBigInt(decoded.Amount)
	t.PricePoint = math.ToBigInt(decoded.PricePoint)

	t.MakeFee = math.ToBigInt(decoded.MakeFee)
	t.TakeFee = math.ToBigInt(decoded.TakeFee)

	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	t.TakerOrderSide = decoded.TakerOrderSide
	return nil
}

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(t.MakerOrderHash.Bytes())
	sha.Write(t.TakerOrderHash.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}