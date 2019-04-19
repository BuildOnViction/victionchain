package tomox

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomodex/utils/math"
)

func EncodeNodeItem(item *Item) (interface{}, error) {
	n := ItemBSON{
		Keys: &KeyMetaBSON{
			Left:   string(item.Keys.Left),
			Right:  string(item.Keys.Right),
			Parent: string(item.Keys.Parent),
		},
		Value: string(item.Value),
		Color: item.Color,
	}

	return n, nil
}

func DecodeNodeItem(raw bson.Raw, item *Item) error {
	decoded := new(struct {
		Keys  *KeyMetaBSON
		Value string
		Color bool
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	item.Keys = &KeyMeta{
		Left:   []byte(decoded.Keys.Left),
		Right:  []byte(decoded.Keys.Right),
		Parent: []byte(decoded.Keys.Parent),
	}
	item.Value = []byte(decoded.Value)
	item.Color = decoded.Color

	return nil
}

func EncodeOrderItem(o *OrderItem) (interface{}, error) {
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

func DecodeOrderItem(raw bson.Raw, o *OrderItem) error {
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
		PricePoint      string           `json:"pricepoint" bson:"pricepoint"`
		Amount          string           `json:"amount" bson:"amount"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		MakeFee         string           `json:"makeFee" bson:"makeFee"`
		TakeFee         string           `json:"takeFee" bson:"takeFee"`
		Signature       *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
		UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
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
	o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	o.Nonce = math.ToBigInt(decoded.Nonce)
	o.MakeFee = math.ToBigInt(decoded.MakeFee)
	o.TakeFee = math.ToBigInt(decoded.TakeFee)
	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Type = decoded.Type
	o.Hash = common.HexToHash(decoded.Hash)

	if decoded.Amount != "" {
		o.Quantity = math.ToBigInt(decoded.Amount)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.PricePoint != "" {
		o.Price = math.ToBigInt(decoded.PricePoint)
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

	return nil
}

func EncodeOrderListItem(item *OrderListItem) (interface{}, error) {

	return nil, nil
}

func DecodeOrderListItem(data bson.Raw, item *OrderListItem) error {

	return nil
}

func EncodeOrderTreeItem(oti *OrderTreeItem) (interface{}, error) {
	otib := OrderTreeItemBSON{
		Volume:        oti.Volume.String(),
		NumOrders:     strconv.FormatUint(oti.NumOrders, 10),
		PriceTreeKey:  string(oti.PriceTreeKey),
		PriceTreeSize: strconv.FormatUint(oti.PriceTreeSize, 10),
	}

	return otib, nil
}

func DecodeOrderTreeItem(raw bson.Raw, oti *OrderTreeItem) error {

	decoded := new(struct {
		Volume        string `json:"volume" bson:"volume"`
		NumOrders     string `json:"numOrders" bson:"numOrders"`
		PriceTreeKey  string `json:"priceTreeKey"`
		PriceTreeSize string `json:"priceTreeSize" bson:"priceTreeSize"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	oti.Volume = math.ToBigInt(decoded.Volume)
	numOrders, err := strconv.ParseInt(decoded.NumOrders, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", numOrders, numOrders)
	}
	oti.NumOrders = uint64(numOrders)
	oti.PriceTreeKey = []byte(decoded.PriceTreeKey)

	priceTreeSize, err := strconv.ParseInt(decoded.PriceTreeSize, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", priceTreeSize, priceTreeSize)
	}
	oti.PriceTreeSize = uint64(priceTreeSize)

	return nil
}

func EncodeOrderBookItem(item *OrderBookItem) (interface{}, error) {

	return nil, nil
}

func DecodeOrderBookItem(data bson.Raw, item *OrderBookItem) error {

	return nil
}

func EncodeItem(val interface{}) (interface{}, error) {
	switch val.(type) {
	case *Item:
		return EncodeNodeItem(val.(*Item))
	case *OrderItem:
		return EncodeOrderItem(val.(*OrderItem))
	case *OrderListItem:
		return EncodeOrderListItem(val.(*OrderListItem))
	case *OrderTreeItem:
		return EncodeOrderTreeItem(val.(*OrderTreeItem))
	case *OrderBookItem:
		return EncodeOrderBookItem(val.(*OrderBookItem))
	default:
		return nil, nil
	}
}

func DecodeItem(data bson.Raw, val interface{}) error {
	switch val.(type) {
	case *Item:
		return DecodeNodeItem(data, val.(*Item))
	case *OrderItem:
		return DecodeOrderItem(data, val.(*OrderItem))
	case *OrderListItem:
		return DecodeOrderListItem(data, val.(*OrderListItem))
	case *OrderTreeItem:
		return DecodeOrderTreeItem(data, val.(*OrderTreeItem))
	case *OrderBookItem:
		return DecodeOrderBookItem(data, val.(*OrderBookItem))
	default:
		return nil
	}
}
