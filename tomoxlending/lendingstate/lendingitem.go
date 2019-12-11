package lendingstate

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomochain/common"
	"math/big"
	"strconv"
	"time"
)

const (
	Investing                  = "INVESTING"
	Borrowing                  = "BORROWING"
	LendingStatusNew           = "NEW"
	LendingStatusOpen          = "OPEN"
	LendingStatusReject        = "REJECTED"
	LendingStatusFilled        = "FILLED"
	LendingStatusPartialFilled = "PARTIAL_FILLED"
	LendingStatusCancelled     = "CANCELLED"
	Ask                        = "SELL"
	Bid                        = "BUY"
	Market                     = "MO"
	Limit                      = "LO"
)

var (
	BaseInterest  = big.NewInt(10000)
	SupportedTerm = map[uint64]bool{
		1:   true,
		7:   true,
		30:  true,
		60:  true,
		90:  true,
		120: true,
		180: true,
		365: true,
	}
)

// Signature struct
type Signature struct {
	V byte        `bson:"v" json:"v"`
	R common.Hash `bson:"r" json:"r"`
	S common.Hash `bson:"s" json:"s"`
}

type SignatureRecord struct {
	V byte   `bson:"v" json:"v"`
	R string `bson:"r" json:"r"`
	S string `bson:"s" json:"s"`
}

type LendingItem struct {
	Quantity *big.Int `bson:"quantity" json:"quantity"`
	Interest *big.Int `bson:"interest" json:"interest"`
	// INVESTING/BORROWING
	Side string `bson:"side" json:"side"`
	// LIMIT/MARKET
	Type            string         `bson:"type" json:"type"`
	LendingToken    common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken common.Address `bson:"collateralToken" json:"collateralToken"`
	FilledAmount    *big.Int       `bson:"filledAmount" json:"filledAmount"`
	Status          string         `bson:"status" json:"status"`
	Relayer         common.Address `bson:"relayer" json:"relayer"`
	Term            uint64         `bson:"term" json:"term"`
	UserAddress     common.Address `bson:"userAddress" json:"userAddress"`
	Signature       *Signature     `bson:"signature" json:"signature"`
	Hash            common.Hash    `bson:"hash" json:"hash"`
	TxHash          common.Hash    `bson:"txHash" json:"txHash"`
	Nonce           *big.Int       `bson:"nonce" json:"nonce"`
	CreatedAt       time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time      `bson:"updatedAt" json:"updatedAt"`
	LendingId       uint64         `bson:"lendingId" json:"lendingId"`
	ExtraData       string         `bson:"extraData" json:"extraData"`
}

type LendingItemBSON struct {
	Quantity string `bson:"quantity" json:"quantity"`
	Interest string `bson:"interest" json:"interest"`
	// INVESTING/BORROWING
	Side string `bson:"side" json:"side"`
	// LIMIT/MARKET
	Type            string           `bson:"type" json:"type"`
	LendingToken    string           `bson:"lendingToken" json:"lendingToken"`
	CollateralToken string           `bson:"collateralToken" json:"collateralToken"`
	FilledAmount    string           `bson:"filledAmount" json:"filledAmount"`
	Status          string           `bson:"status" json:"status"`
	Relayer         string           `bson:"relayer" json:"relayer"`
	Term            string           `bson:"term" json:"term"`
	UserAddress     string           `bson:"userAddress" json:"userAddress"`
	Signature       *SignatureRecord `bson:"signature" json:"signature"`
	Hash            string           `bson:"hash" json:"hash"`
	TxHash          string           `bson:"txHash" json:"txHash"`
	Nonce           string           `bson:"nonce" json:"nonce"`
	CreatedAt       time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `bson:"updatedAt" json:"updatedAt"`
	LendingId       string           `bson:"lendingId" json:"lendingId"`
	ExtraData       string           `bson:"extraData" json:"extraData"`
}

func (l *LendingItem) GetBSON() (interface{}, error) {
	lr := LendingItemBSON{
		Quantity:        l.Quantity.String(),
		Interest:        l.Interest.String(),
		Side:            l.Side,
		Type:            l.Type,
		LendingToken:    l.LendingToken.Hex(),
		CollateralToken: l.CollateralToken.Hex(),
		Status:          l.Status,
		Relayer:         l.Relayer.Hex(),
		Term:            strconv.FormatUint(l.Term, 10),
		UserAddress:     l.UserAddress.Hex(),

		Hash:      l.Hash.Hex(),
		TxHash:    l.TxHash.Hex(),
		Nonce:     l.Nonce.String(),
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		LendingId: strconv.FormatUint(l.LendingId, 10),
		ExtraData: l.ExtraData,
	}

	if l.FilledAmount != nil {
		lr.FilledAmount = l.FilledAmount.String()
	}

	if l.Signature != nil {
		lr.Signature = &SignatureRecord{
			V: l.Signature.V,
			R: l.Signature.R.Hex(),
			S: l.Signature.S.Hex(),
		}
	}

	return lr, nil
}

func (l *LendingItem) SetBSON(raw bson.Raw) error {
	decoded := new(LendingItemBSON)

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	if decoded.Quantity != "" {
		l.Quantity = ToBigInt(decoded.Quantity)
	}
	l.Interest = ToBigInt(decoded.Interest)
	l.Side = decoded.Side
	l.Type = decoded.Type
	l.LendingToken = common.HexToAddress(decoded.LendingToken)
	l.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	l.FilledAmount = ToBigInt(decoded.FilledAmount)
	l.Status = decoded.Status
	l.Relayer = common.HexToAddress(decoded.Relayer)
	term, err := strconv.ParseInt(decoded.Term, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.term. Err: %v", err)
	}
	l.Term = uint64(term)
	l.UserAddress = common.HexToAddress(decoded.UserAddress)

	if decoded.Signature != nil {
		l.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	l.Hash = common.HexToHash(decoded.Hash)
	l.TxHash = common.HexToHash(decoded.TxHash)
	l.Nonce = ToBigInt(decoded.Nonce)

	l.CreatedAt = decoded.CreatedAt
	l.UpdatedAt = decoded.UpdatedAt
	lendingId, err := strconv.ParseInt(decoded.LendingId, 10, 64)
	if err != nil {
		return err
	}
	l.LendingId = uint64(lendingId)
	l.ExtraData = decoded.ExtraData
	return nil
}
