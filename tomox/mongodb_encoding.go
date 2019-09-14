package tomox

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
)

func (o *OrderItem) GetBSON() (interface{}, error) {
	or := OrderItemBSON{
		PairName:        o.PairName,
		ExchangeAddress: o.ExchangeAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		BaseToken:       o.BaseToken.Hex(),
		QuoteToken:      o.QuoteToken.Hex(),
		Status:          o.Status,
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash.Hex(),
		Quantity:        o.Quantity.String(),
		Price:           o.Price.String(),
		Nonce:           o.Nonce.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       strconv.FormatUint(o.CreatedAt, 10),
		UpdatedAt:       strconv.FormatUint(o.UpdatedAt, 10),
		OrderID:         strconv.FormatUint(o.OrderID, 10),
		Key:             o.Key,
	}

	if o.FilledAmount != nil {
		or.FilledAmount = o.FilledAmount.String()
	}

	if o.Signature != nil {
		or.Signature = &SignatureRecord{
			V: o.Signature.V,
			R: o.Signature.R.Hex(),
			S: o.Signature.S.Hex(),
		}
	}

	return or, nil
}

func (o *OrderItem) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID              bson.ObjectId    `json:"id,omitempty" bson:"_id"`
		PairName        string           `json:"pairName" bson:"pairName"`
		ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
		UserAddress     string           `json:"userAddress" bson:"userAddress"`
		BaseToken       string           `json:"baseToken" bson:"baseToken"`
		QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
		Status          string           `json:"status" bson:"status"`
		Side            string           `json:"side" bson:"side"`
		Type            string           `json:"type" bson:"type"`
		Hash            string           `json:"hash" bson:"hash"`
		Price           string           `json:"price" bson:"price"`
		Quantity        string           `json:"quantity" bson:"quantity"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		MakeFee         string           `json:"makeFee" bson:"makeFee"`
		TakeFee         string           `json:"takeFee" bson:"takeFee"`
		Signature       *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
		UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
		OrderID         string           `json:"orderID" bson:"orderID"`
		Key             string           `json:"key" bson:"key"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	o.PairName = decoded.PairName
	o.ExchangeAddress = common.HexToAddress(decoded.ExchangeAddress)
	o.UserAddress = common.HexToAddress(decoded.UserAddress)
	o.BaseToken = common.HexToAddress(decoded.BaseToken)
	o.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	o.FilledAmount = ToBigInt(decoded.FilledAmount)
	o.Nonce = ToBigInt(decoded.Nonce)
	o.MakeFee = ToBigInt(decoded.MakeFee)
	o.TakeFee = ToBigInt(decoded.TakeFee)
	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Type = decoded.Type
	o.Hash = common.HexToHash(decoded.Hash)

	if decoded.Quantity != "" {
		o.Quantity = ToBigInt(decoded.Quantity)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = ToBigInt(decoded.FilledAmount)
	}

	if decoded.Price != "" {
		o.Price = ToBigInt(decoded.Price)
	}

	if decoded.Signature != nil {
		o.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	o.CreatedAt = uint64(decoded.CreatedAt.Unix())
	o.UpdatedAt = uint64(decoded.UpdatedAt.Unix())
	orderID, err := strconv.ParseInt(decoded.OrderID, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", orderID, orderID)
	}
	o.OrderID = uint64(orderID)
	o.Key = decoded.Key

	return nil
}
