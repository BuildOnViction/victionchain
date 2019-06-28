package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/tomochain/tomox-sdk/utils/math"
)

// Token struct is used to model the token data in the system and DB
type Token struct {
	ID              bson.ObjectId  `json:"-" bson:"_id"`
	Name            string         `json:"name" bson:"name"`
	Symbol          string         `json:"symbol" bson:"symbol"`
	Address         common.Address `json:"address" bson:"address"`
	Image           Image          `json:"image" bson:"image"`
	ContractAddress common.Address `json:"contractAddress" bson:"contractAddress"`
	Decimals        int            `json:"decimals" bson:"decimals"`
	Active          bool           `json:"active" bson:"active"`
	Listed          bool           `json:"listed" bson:"listed"`
	Quote           bool           `json:"quote" bson:"quote"`
	MakeFee         *big.Int       `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee         *big.Int       `json:"takeFee,omitempty" bson:"makeFee,omitempty"`
	USD             string         `json:"usd,omitempty" bson:"usd,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// TokenRecord is the struct which is stored in db
type TokenRecord struct {
	ID              bson.ObjectId `json:"-" bson:"_id"`
	Name            string        `json:"name" bson:"name"`
	Symbol          string        `json:"symbol" bson:"symbol"`
	Image           Image         `json:"image" bson:"image"`
	ContractAddress string        `json:"contractAddress" bson:"contractAddress"`
	Decimals        int           `json:"decimals" bson:"decimals"`
	Active          bool          `json:"active" bson:"active"`
	Quote           bool          `json:"quote" bson:"quote"`
	MakeFee         string        `json:"makeFee,omitempty" bson:"makeFee,omitempty"`
	TakeFee         string        `json:"takeFee,omitempty" bson:"takeFee,omitempty"`
	USD             string        `json:"usd,omitempty" bson:"usd,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Image is a sub document used to store data related to images
type Image struct {
	URL  string                 `json:"url" bson:"url"`
	Meta map[string]interface{} `json:"meta" bson:"meta"`
}

type NativeCurrency struct {
	Address  common.Address `json:"address" bson:"address"`
	Symbol   string         `json:"symbol" bson:"symbol"`
	Decimals int            `json:"decimals" bson:"decimals"`
}

func GetNativeCurrency() NativeCurrency {
	return NativeCurrency{
		Address:  common.HexToAddress("0x1"),
		Symbol:   "TOMO",
		Decimals: 18,
	}
}

// DefaultTestBalance returns the default balance
// Only for testing/mock purpose
func DefaultTestBalance() int64 {
	return 50000
}

// DefaultTestBalance returns the default locked balance
// Only for testing/mock purpose
func DefaultTestInOrderBalance() int64 {
	return 0
}

// DefaultTestAvailableBalance returns the default available balance
// Only for testing/mock purpose
func DefaultTestAvailableBalance() int64 {
	return 0
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Symbol, validation.Required),
		validation.Field(&t.ContractAddress, validation.Required),
		validation.Field(&t.Decimals, validation.Required),
	)
}

func (t *Token) MarshalJSON() ([]byte, error) {
	token := map[string]interface{}{
		"id":              t.ID,
		"symbol":          t.Symbol,
		"contractAddress": t.ContractAddress.Hex(),
		"decimals":        t.Decimals,
		"image":           t.Image,
		"active":          t.Active,
		"quote":           t.Quote,
		"usd":             t.USD,
		"createdAt":       t.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":       t.UpdatedAt.Format(time.RFC3339Nano),
	}

	if t.MakeFee != nil {
		token["makeFee"] = t.MakeFee.String()
	}

	if t.TakeFee != nil {
		token["takeFee"] = t.TakeFee.String()
	}

	return json.Marshal(token)
}

func (t *Token) UnmarshalJSON(b []byte) error {
	token := map[string]interface{}{}

	err := json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

	t.ID = bson.ObjectIdHex(token["id"].(string))
	t.Symbol = token["symbol"].(string)
	t.Decimals = token["decimals"].(int)
	t.Active = token["active"].(bool)
	t.Quote = token["quote"].(bool)
	t.USD = token["usd"].(string)

	if token["createdAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, token["createdAt"].(string))
		t.CreatedAt = tm
	}

	if token["updatedAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, token["updatedAt"].(string))
		t.UpdatedAt = tm
	}

	if token["makeFee"] != nil {
		t.MakeFee = math.ToBigInt(token["makeFee"].(string))
	}

	if token["takeFee"] != nil {
		t.TakeFee = math.ToBigInt(token["takeFee"].(string))
	}

	image, ok := token["image"].(map[string]interface{})
	if ok {
		t.Image.URL = image["url"].(string)
		t.Image.Meta = image["meta"].(map[string]interface{})
	}

	return nil
}

// GetBSON implements bson.Getter
func (t *Token) GetBSON() (interface{}, error) {

	tr := TokenRecord{
		ID:              t.ID,
		Name:            t.Name,
		Symbol:          t.Symbol,
		Image:           t.Image,
		ContractAddress: t.ContractAddress.Hex(),
		Decimals:        t.Decimals,
		Active:          t.Active,
		Quote:           t.Quote,
		USD:             t.USD,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}

	if t.MakeFee != nil {
		tr.MakeFee = t.MakeFee.String()
	}

	if t.TakeFee != nil {
		tr.TakeFee = t.TakeFee.String()
	}

	return tr, nil
}

// SetBSON implemenets bson.Setter
func (t *Token) SetBSON(raw bson.Raw) error {
	decoded := &TokenRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	t.ID = decoded.ID
	t.Name = decoded.Name
	t.Symbol = decoded.Symbol
	t.Image = decoded.Image
	if common.IsHexAddress(decoded.ContractAddress) {
		t.ContractAddress = common.HexToAddress(decoded.ContractAddress)
	}
	t.Decimals = decoded.Decimals
	t.Active = decoded.Active
	t.Quote = decoded.Quote
	t.USD = decoded.USD
	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	if decoded.MakeFee != "" {
		t.MakeFee = math.ToBigInt(decoded.MakeFee)
	}

	if decoded.TakeFee != "" {
		t.TakeFee = math.ToBigInt(decoded.TakeFee)
	}
	return nil
}
