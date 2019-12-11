package lendingstate

import (
	"fmt"
	"github.com/tomochain/tomochain/tomox/tomox_state"
	"math/big"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomochain/common"
)

type LendingTrade struct {
	Borrower            common.Address `bson:"borrower" json:"borrower"`
	Investor            common.Address `bson:"investor" json:"investor"`
	LendingToken        common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken     common.Address `bson:"collateralToken" json:"collateralToken"`
	TakerOrderHash      common.Hash    `bson:"takerOrderHash" json:"takerOrderHash"`
	MakerOrderHash      common.Hash    `bson:"makerOrderHash" json:"makerOrderHash"`
	BorrowingRelayer    common.Address `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer    common.Address `bson:"investingRelayer" json:"investingRelayer"`
	Term                uint64         `bson:"term" json:"term"`
	Interest            uint64         `bson:"interest" json:"interest"`
	CollateralInterest  *big.Int       `bson:"collateralInterest" json:"collateralInterest"`
	LiquidationInterest *big.Int       `bson:"liquidationInterest" json:"liquidationInterest"`
	Amount              *big.Int       `bson:"amount" json:"amount"`
	BorrowingFee        *big.Int       `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee        *big.Int       `bson:"investingFee" json:"investingFee"`
	Status              string         `bson:"status" json:"status"`
	TakerOrderSide      string         `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType      string         `bson:"takerOrderType" json:"takerOrderType"`
	MakerOrderType      string         `bson:"makerOrderType" json:"makerOrderType"`
	TradeId             uint64         `bson:"tradeId" json:"tradeId"`
	Hash                common.Hash    `bson:"hash" json:"hash"`
	TxHash              common.Hash    `bson:"txHash" json:"txHash"`
	ExtraData           string         `bson:"extraData" json:"extraData"`
	CreatedAt           time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time      `bson:"updatedAt" json:"updatedAt"`
}

type LendingTradeBSON struct {
	Borrower            string    `bson:"borrower" json:"borrower"`
	Investor            string    `bson:"investor" json:"investor"`
	LendingToken        string    `bson:"lendingToken" json:"lendingToken"`
	CollateralToken     string    `bson:"collateralToken" json:"collateralToken"`
	TakerOrderHash      string    `bson:"takerOrderHash" json:"takerOrderHash"`
	MakerOrderHash      string    `bson:"makerOrderHash" json:"makerOrderHash"`
	BorrowingRelayer    string    `bson:"borrowingRelayer" json:"borrowingRelayer"`
	InvestingRelayer    string    `bson:"investingRelayer" json:"investingRelayer"`
	Term                string    `bson:"term" json:"term"`
	Interest            string    `bson:"interest" json:"interest"`
	CollateralInterest  string    `bson:"collateralInterest" json:"collateralInterest"`
	LiquidationInterest string    `bson:"liquidationInterest" json:"liquidationInterest"`
	Amount              string    `bson:"amount" json:"amount"`
	BorrowingFee        string    `bson:"borrowingFee" json:"borrowingFee"`
	InvestingFee        string    `bson:"investingFee" json:"investingFee"`
	Status              string    `bson:"status" json:"status"`
	TakerOrderSide      string    `bson:"takerOrderSide" json:"takerOrderSide"`
	TakerOrderType      string    `bson:"takerOrderType" json:"takerOrderType"`
	MakerOrderType      string    `bson:"makerOrderType" json:"makerOrderType"`
	TradeId             string    `bson:"tradeId" json:"tradeId"`
	Hash                string    `bson:"hash" json:"hash"`
	TxHash              string    `bson:"txHash" json:"txHash"`
	ExtraData           string    `bson:"extraData" json:"extraData"`
	CreatedAt           time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt           time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (t *LendingTrade) GetBSON() (interface{}, error) {
	tr := LendingTradeBSON{
		Borrower:            t.Borrower.Hex(),
		Investor:            t.Investor.Hex(),
		LendingToken:        t.LendingToken.Hex(),
		CollateralToken:     t.CollateralToken.Hex(),
		TakerOrderHash:      t.TakerOrderHash.Hex(),
		MakerOrderHash:      t.MakerOrderHash.Hex(),
		BorrowingRelayer:    t.BorrowingRelayer.Hex(),
		InvestingRelayer:    t.InvestingRelayer.Hex(),
		Term:                strconv.FormatUint(t.Term, 10),
		Interest:            strconv.FormatUint(t.Interest, 10),
		CollateralInterest:  t.CollateralInterest.String(),
		LiquidationInterest: t.LiquidationInterest.String(),
		Amount:              t.Amount.String(),
		BorrowingFee:        t.BorrowingFee.String(),
		InvestingFee:        t.InvestingFee.String(),
		Status:              t.Status,
		TakerOrderSide:      t.TakerOrderSide,
		TakerOrderType:      t.TakerOrderType,
		MakerOrderType:      t.MakerOrderType,
		TradeId:             strconv.FormatUint(t.TradeId, 10),
		Hash:                t.Hash.Hex(),
		TxHash:              t.TxHash.Hex(),
		ExtraData:           t.ExtraData,
		CreatedAt:           t.CreatedAt,
		UpdatedAt:           t.UpdatedAt,
	}

	return tr, nil
}

func (t *LendingTrade) SetBSON(raw bson.Raw) error {
	decoded := new(LendingTradeBSON)

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	t.Borrower = common.HexToAddress(decoded.Borrower)
	t.Investor = common.HexToAddress(decoded.Investor)
	t.LendingToken = common.HexToAddress(decoded.LendingToken)
	t.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	t.TakerOrderHash = common.HexToHash(decoded.TakerOrderHash)
	t.MakerOrderHash = common.HexToHash(decoded.MakerOrderHash)
	t.BorrowingRelayer = common.HexToAddress(decoded.BorrowingRelayer)
	t.InvestingRelayer = common.HexToAddress(decoded.InvestingRelayer)
	term, err := strconv.ParseInt(decoded.Term, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.term. Err: %v", err)
	}
	t.Term = uint64(term)
	interest, err := strconv.ParseInt(decoded.Interest, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.interest. Err: %v", err)
	}
	t.Interest = uint64(interest)
	t.CollateralInterest = ToBigInt(decoded.CollateralInterest)
	t.LiquidationInterest = ToBigInt(decoded.LiquidationInterest)
	t.Amount = tomox_state.ToBigInt(decoded.Amount)
	t.BorrowingFee = tomox_state.ToBigInt(decoded.BorrowingFee)
	t.InvestingFee = tomox_state.ToBigInt(decoded.InvestingFee)
	t.Status = decoded.Status
	t.TakerOrderSide = decoded.TakerOrderSide
	t.TakerOrderType = decoded.TakerOrderType
	t.MakerOrderType = decoded.MakerOrderType
	t.ExtraData = decoded.ExtraData
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt

	return nil
}
