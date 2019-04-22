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

func EncodeOrderListItem(item *OrderListItem) (interface{}, error) {

	return nil, nil
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

func EncodeOrderBookItem(item *OrderBookItem) (interface{}, error) {

	return nil, nil
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

func (nir *ItemRecord) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Key   string
		Value ItemBSON
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	nir.Key = decoded.Key
	nir.Value = &Item{
		Keys: &KeyMeta{
			Left:   []byte(decoded.Value.Keys.Left),
			Right:  []byte(decoded.Value.Keys.Right),
			Parent: []byte(decoded.Value.Keys.Parent),
		},
		Value: []byte(decoded.Value.Value),
		Color: decoded.Value.Color,
	}

	return nil
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

func (otir *OrderTreeItemRecord) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Key   string
		Value OrderTreeItemBSON
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	otir.Key = decoded.Key

	otir.Value.Volume = math.ToBigInt(decoded.Value.Volume)
	numOrders, err := strconv.ParseInt(decoded.Value.NumOrders, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", numOrders, numOrders)
	}
	otir.Value.NumOrders = uint64(numOrders)
	otir.Value.PriceTreeKey = []byte(decoded.Value.PriceTreeKey)

	priceTreeSize, err := strconv.ParseInt(decoded.Value.PriceTreeSize, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", priceTreeSize, priceTreeSize)
	}
	otir.Value.PriceTreeSize = uint64(priceTreeSize)

	return nil
}

func (olir *OrderListItemRecord) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Key   string
		Value OrderListItemBSON
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	olir.Key = decoded.Key

	olir.Value.HeadOrder = []byte(decoded.Value.HeadOrder)
	olir.Value.TailOrder = []byte(decoded.Value.TailOrder)

	length, err := strconv.ParseInt(decoded.Value.Length, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", length, length)
	}
	olir.Value.Length = uint64(length)

	olir.Value.Volume = math.ToBigInt(decoded.Value.Volume)
	olir.Value.Price = math.ToBigInt(decoded.Value.Price)

	return nil
}

func (obir *OrderBookItemRecord) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Key   string
		Value OrderBookItemBSON
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	obir.Key = decoded.Key

	timestamp, err := strconv.ParseInt(decoded.Value.Timestamp, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", timestamp, timestamp)
	}
	obir.Value.Timestamp = uint64(timestamp)

	nextOrderID, err := strconv.ParseInt(decoded.Value.NextOrderID, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", nextOrderID, nextOrderID)
	}
	obir.Value.NextOrderID = uint64(nextOrderID)

	maxPricePoint, err := strconv.ParseInt(decoded.Value.MaxPricePoint, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", maxPricePoint, maxPricePoint)
	}
	obir.Value.MaxPricePoint = uint64(maxPricePoint)

	obir.Value.Name = decoded.Value.Name

	return nil
}
