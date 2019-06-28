package types

import (
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/utils/math"
)

const (
	TypeStopMarketOrder = "SMO"
	TypeStopLimitOrder  = "SLO"

	StopOrderStatusOpen      = "OPEN"
	StopOrderStatusDone      = "DONE"
	StopOrderStatusCancelled = "CANCELLED"
)

type StopOrder struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	UserAddress     common.Address `json:"userAddress" bson:"userAddress"`
	ExchangeAddress common.Address `json:"exchangeAddress" bson:"exchangeAddress"`
	BaseToken       common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken      common.Address `json:"quoteToken" bson:"quoteToken"`
	Status          string         `json:"status" bson:"status"`
	Side            string         `json:"side" bson:"side"`
	Type            string         `json:"type" bson:"type"`
	Hash            common.Hash    `json:"hash" bson:"hash"`
	Signature       *Signature     `json:"signature,omitempty" bson:"signature"`
	StopPrice       *big.Int       `json:"stopPrice" bson:"stopPrice"`
	LimitPrice      *big.Int       `json:"limitPrice" bson:"limitPrice"`
	Direction       int            `json:"direction" bson:"direction"`
	Amount          *big.Int       `json:"amount" bson:"amount"`
	FilledAmount    *big.Int       `json:"filledAmount" bson:"filledAmount"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	MakeFee         *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee         *big.Int       `json:"takeFee" bson:"takeFee"`
	PairName        string         `json:"pairName" bson:"pairName"`
	CreatedAt       time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// MarshalJSON implements the json.Marshal interface
func (so *StopOrder) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"exchangeAddress": so.ExchangeAddress,
		"userAddress":     so.UserAddress,
		"baseToken":       so.BaseToken,
		"quoteToken":      so.QuoteToken,
		"side":            so.Side,
		"type":            so.Type,
		"status":          so.Status,
		"pairName":        so.PairName,
		"amount":          so.Amount.String(),
		"stopPrice":       so.StopPrice.String(),
		"limitPrice":      so.LimitPrice.String(),
		"direction":       strconv.Itoa(so.Direction),
		"makeFee":         so.MakeFee.String(),
		"takeFee":         so.TakeFee.String(),
		"createdAt":       so.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":       so.UpdatedAt.Format(time.RFC3339Nano),
	}

	if so.FilledAmount != nil {
		order["filledAmount"] = so.FilledAmount.String()
	}

	if so.Hash.Hex() != "" {
		order["hash"] = so.Hash.Hex()
	}

	if so.Nonce != nil {
		order["nonce"] = so.Nonce.String()
	}

	if so.Signature != nil {
		order["signature"] = map[string]interface{}{
			"V": so.Signature.V,
			"R": so.Signature.R,
			"S": so.Signature.S,
		}
	}

	return json.Marshal(order)
}

// UnmarshalJSON : write custom logic to unmarshal bytes to StopOrder
func (so *StopOrder) UnmarshalJSON(b []byte) error {
	order := map[string]interface{}{}

	err := json.Unmarshal(b, &order)
	if err != nil {
		return err
	}

	if order["id"] != nil && bson.IsObjectIdHex(order["id"].(string)) {
		so.ID = bson.ObjectIdHex(order["id"].(string))
	}

	if order["pairName"] != nil {
		so.PairName = order["pairName"].(string)
	}

	if order["exchangeAddress"] != nil {
		so.ExchangeAddress = common.HexToAddress(order["exchangeAddress"].(string))
	}

	if order["userAddress"] != nil {
		so.UserAddress = common.HexToAddress(order["userAddress"].(string))
	}

	if order["baseToken"] != nil {
		so.BaseToken = common.HexToAddress(order["baseToken"].(string))
	}

	if order["quoteToken"] != nil {
		so.QuoteToken = common.HexToAddress(order["quoteToken"].(string))
	}

	if order["stopPrice"] != nil {
		so.StopPrice = math.ToBigInt(order["stopPrice"].(string))
	}

	if order["limitPrice"] != nil {
		so.LimitPrice = math.ToBigInt(order["limitPrice"].(string))
	}

	if order["direction"] != nil {

		if direction, err := strconv.Atoi(order["direction"].(string)); err == nil {
			so.Direction = direction
		} else {
			return errors.New("Direction parameter is not an integer.")
		}

	}

	if order["amount"] != nil {
		so.Amount = math.ToBigInt(order["amount"].(string))
	}

	if order["filledAmount"] != nil {
		so.FilledAmount = math.ToBigInt(order["filledAmount"].(string))
	}

	if order["nonce"] != nil {
		so.Nonce = math.ToBigInt(order["nonce"].(string))
	}

	if order["makeFee"] != nil {
		so.MakeFee = math.ToBigInt(order["makeFee"].(string))
	}

	if order["takeFee"] != nil {
		so.TakeFee = math.ToBigInt(order["takeFee"].(string))
	}

	if order["hash"] != nil {
		so.Hash = common.HexToHash(order["hash"].(string))
	}

	if order["side"] != nil {
		so.Side = order["side"].(string)
	}

	if order["type"] != nil {
		so.Type = order["type"].(string)
	}

	if order["status"] != nil {
		so.Status = order["status"].(string)
	}

	if order["signature"] != nil {
		signature := order["signature"].(map[string]interface{})
		so.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	if order["createdAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["createdAt"].(string))
		so.CreatedAt = t
	}

	if order["updatedAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["updatedAt"].(string))
		so.UpdatedAt = t
	}

	return nil
}

func (so *StopOrder) GetBSON() (interface{}, error) {
	or := StopOrderRecord{
		PairName:        so.PairName,
		ExchangeAddress: so.ExchangeAddress.Hex(),
		UserAddress:     so.UserAddress.Hex(),
		BaseToken:       so.BaseToken.Hex(),
		QuoteToken:      so.QuoteToken.Hex(),
		Status:          so.Status,
		Side:            so.Side,
		Type:            so.Type,
		Hash:            so.Hash.Hex(),
		Amount:          so.Amount.String(),
		StopPrice:       so.StopPrice.String(),
		LimitPrice:      so.LimitPrice.String(),
		Direction:       so.Direction,
		Nonce:           so.Nonce.String(),
		MakeFee:         so.MakeFee.String(),
		TakeFee:         so.TakeFee.String(),
		CreatedAt:       so.CreatedAt,
		UpdatedAt:       so.UpdatedAt,
	}

	if so.ID.Hex() == "" {
		or.ID = bson.NewObjectId()
	} else {
		or.ID = so.ID
	}

	if so.FilledAmount != nil {
		or.FilledAmount = so.FilledAmount.String()
	}

	if so.Signature != nil {
		or.Signature = &SignatureRecord{
			V: so.Signature.V,
			R: so.Signature.R.Hex(),
			S: so.Signature.S.Hex(),
		}
	}

	return or, nil
}

func (so *StopOrder) SetBSON(raw bson.Raw) error {
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
		StopPrice       string           `json:"stopPrice" bson:"stopPrice"`
		LimitPrice      string           `json:"limitPrice" bson:"limitPrice"`
		Direction       int              `json:"direction" bson:"direction"`
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
		logger.Error(err)
		return err
	}

	so.ID = decoded.ID
	so.PairName = decoded.PairName
	so.ExchangeAddress = common.HexToAddress(decoded.ExchangeAddress)
	so.UserAddress = common.HexToAddress(decoded.UserAddress)
	so.BaseToken = common.HexToAddress(decoded.BaseToken)
	so.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	so.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	so.Nonce = math.ToBigInt(decoded.Nonce)
	so.MakeFee = math.ToBigInt(decoded.MakeFee)
	so.TakeFee = math.ToBigInt(decoded.TakeFee)
	so.Status = decoded.Status
	so.Side = decoded.Side
	so.Type = decoded.Type
	so.Hash = common.HexToHash(decoded.Hash)

	if decoded.Amount != "" {
		so.Amount = math.ToBigInt(decoded.Amount)
	}

	if decoded.FilledAmount != "" {
		so.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.StopPrice != "" {
		so.StopPrice = math.ToBigInt(decoded.StopPrice)
	}

	if decoded.LimitPrice != "" {
		so.LimitPrice = math.ToBigInt(decoded.LimitPrice)
	}

	so.Direction = decoded.Direction

	if decoded.Signature != nil {
		so.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	so.CreatedAt = decoded.CreatedAt
	so.UpdatedAt = decoded.UpdatedAt

	return nil
}

// ToOrder converts a stop order to a real order that will be pushed to TomoX
func (so *StopOrder) ToOrder() (*Order, error) {
	var o *Order

	switch so.Type {
	case TypeStopMarketOrder:
		o = &Order{
			UserAddress:     so.UserAddress,
			ExchangeAddress: so.ExchangeAddress,
			BaseToken:       so.BaseToken,
			QuoteToken:      so.QuoteToken,
			Status:          OrderStatusOpen,
			Side:            so.Side,
			Type:            TypeMarketOrder,
			Hash:            so.Hash,
			Signature:       so.Signature,
			PricePoint:      so.StopPrice,
			Amount:          so.Amount,
			FilledAmount:    big.NewInt(0),
			Nonce:           so.Nonce,
			MakeFee:         so.MakeFee,
			TakeFee:         so.TakeFee,
			PairName:        so.PairName,
		}

		break
	case TypeStopLimitOrder:
		o = &Order{
			UserAddress:     so.UserAddress,
			ExchangeAddress: so.ExchangeAddress,
			BaseToken:       so.BaseToken,
			QuoteToken:      so.QuoteToken,
			Status:          OrderStatusOpen,
			Side:            so.Side,
			Type:            TypeLimitOrder,
			Hash:            so.Hash,
			Signature:       so.Signature,
			PricePoint:      so.LimitPrice,
			Amount:          so.Amount,
			FilledAmount:    big.NewInt(0),
			Nonce:           so.Nonce,
			MakeFee:         so.MakeFee,
			TakeFee:         so.TakeFee,
			PairName:        so.PairName,
		}

		break
	default:
		return nil, errors.New("Unknown stop order type")
	}

	return o, nil
}

// TODO: Verify userAddress, baseToken, quoteToken, etc. conditions are working
func (so *StopOrder) Validate() error {
	if so.ExchangeAddress != common.HexToAddress(app.Config.Ethereum["exchange_address"]) {
		return errors.New("Order 'exchangeAddress' parameter is incorrect")
	}

	if (so.UserAddress == common.Address{}) {
		return errors.New("Order 'userAddress' parameter is required")
	}

	if so.Nonce == nil {
		return errors.New("Order 'nonce' parameter is required")
	}

	if (so.BaseToken == common.Address{}) {
		return errors.New("Order 'baseToken' parameter is required")
	}

	if (so.QuoteToken == common.Address{}) {
		return errors.New("Order 'quoteToken' parameter is required")
	}

	if so.MakeFee == nil {
		return errors.New("Order 'makeFee' parameter is required")
	}

	if so.TakeFee == nil {
		return errors.New("Order 'takeFee' parameter is required")
	}

	if so.Amount == nil {
		return errors.New("Order 'amount' parameter is required")
	}

	if so.StopPrice == nil {
		return errors.New("Order 'stopPrice' parameter is required")
	}

	if so.Side != BUY && so.Side != SELL {
		return errors.New("Order 'side' should be 'SELL' or 'BUY', but got: '" + so.Side + "'")
	}

	if so.Signature == nil {
		return errors.New("Order 'signature' parameter is required")
	}

	if math.IsSmallerThan(so.Nonce, big.NewInt(0)) {
		return errors.New("Order 'nonce' parameter should be positive")
	}

	if math.IsEqualOrSmallerThan(so.Amount, big.NewInt(0)) {
		return errors.New("Order 'amount' parameter should be strictly positive")
	}

	if math.IsEqualOrSmallerThan(so.StopPrice, big.NewInt(0)) {
		return errors.New("Order 'stopPrice' parameter should be strictly positive")
	}

	valid, err := so.VerifySignature()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("Order 'signature' parameter is invalid")
	}

	return nil
}

// ComputeHash calculates the orderRequest hash
func (so *StopOrder) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(so.ExchangeAddress.Bytes())
	sha.Write(so.UserAddress.Bytes())
	sha.Write(so.BaseToken.Bytes())
	sha.Write(so.QuoteToken.Bytes())
	sha.Write(common.BigToHash(so.Amount).Bytes())
	sha.Write(common.BigToHash(so.StopPrice).Bytes())
	sha.Write(common.BigToHash(so.EncodedSide()).Bytes())
	sha.Write(common.BigToHash(so.Nonce).Bytes())
	sha.Write(common.BigToHash(so.MakeFee).Bytes())
	sha.Write(common.BigToHash(so.TakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (so *StopOrder) VerifySignature() (bool, error) {
	so.Hash = so.ComputeHash()

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		so.Hash.Bytes(),
	)

	address, err := so.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != so.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

func (so *StopOrder) Process(p *Pair) error {
	if so.FilledAmount == nil {
		so.FilledAmount = big.NewInt(0)
	}

	// TODO: Handle this in Validate function
	if so.Type != TypeStopMarketOrder && so.Type != TypeStopLimitOrder {
		so.Type = TypeStopLimitOrder
	}

	if !math.IsEqual(so.MakeFee, p.MakeFee) {
		return errors.New("Invalid MakeFee")
	}

	if !math.IsEqual(so.TakeFee, p.TakeFee) {
		return errors.New("Invalid TakeFee")
	}

	so.PairName = p.Name()
	so.CreatedAt = time.Now()
	so.UpdatedAt = time.Now()
	return nil
}

func (so *StopOrder) QuoteAmount(p *Pair) *big.Int {
	pairMultiplier := p.PairMultiplier()
	return math.Div(math.Mul(so.Amount, so.StopPrice), pairMultiplier)
}

//TODO handle error case ?
func (so *StopOrder) EncodedSide() *big.Int {
	if so.Side == BUY {
		return big.NewInt(0)
	} else {
		return big.NewInt(1)
	}
}

func (so *StopOrder) PairCode() (string, error) {
	if so.PairName == "" {
		return "", errors.New("Pair name is required")
	}

	return so.PairName + "::" + so.BaseToken.Hex() + "::" + so.QuoteToken.Hex(), nil
}

// StopOrderRecord is the object that will be saved in the database
type StopOrderRecord struct {
	ID              bson.ObjectId    `json:"id" bson:"_id"`
	UserAddress     string           `json:"userAddress" bson:"userAddress"`
	ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
	BaseToken       string           `json:"baseToken" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
	Status          string           `json:"status" bson:"status"`
	Side            string           `json:"side" bson:"side"`
	Type            string           `json:"type" bson:"type"`
	Hash            string           `json:"hash" bson:"hash"`
	StopPrice       string           `json:"stopPrice" bson:"stopPrice"`
	LimitPrice      string           `json:"limitPrice" bson:"limitPrice"`
	Direction       int              `json:"direction" bson:"direction"`
	Amount          string           `json:"amount" bson:"amount"`
	FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
	Nonce           string           `json:"nonce" bson:"nonce"`
	MakeFee         string           `json:"makeFee" bson:"makeFee"`
	TakeFee         string           `json:"takeFee" bson:"takeFee"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`

	PairName  string    `json:"pairName" bson:"pairName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type StopOrderBSONUpdate struct {
	*StopOrder
}

func (o StopOrderBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()

	set := bson.M{
		"pairName":        o.PairName,
		"exchangeAddress": o.ExchangeAddress.Hex(),
		"userAddress":     o.UserAddress.Hex(),
		"baseToken":       o.BaseToken.Hex(),
		"quoteToken":      o.QuoteToken.Hex(),
		"status":          o.Status,
		"side":            o.Side,
		"type":            o.Type,
		"stopPrice":       o.StopPrice.String(),
		"limitPrice":      o.LimitPrice.String(),
		"direction":       o.Direction,
		"amount":          o.Amount.String(),
		"nonce":           o.Nonce.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"updatedAt":       now,
	}

	if o.FilledAmount != nil {
		set["filledAmount"] = o.FilledAmount.String()
	}

	if o.Signature != nil {
		set["signature"] = bson.M{
			"V": o.Signature.V,
			"R": o.Signature.R.Hex(),
			"S": o.Signature.S.Hex(),
		}
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"hash":      o.Hash.Hex(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}
