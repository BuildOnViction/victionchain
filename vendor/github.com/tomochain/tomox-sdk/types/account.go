package types

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/go-ozzo/ozzo-validation"
)

// Account corresponds to a single Ethereum address. It contains a list of token balances for that address
type Account struct {
	ID             bson.ObjectId                    `json:"-" bson:"_id"`
	Address        common.Address                   `json:"address" bson:"address"`
	TokenBalances  map[common.Address]*TokenBalance `json:"tokenBalances" bson:"tokenBalances"`
	FavoriteTokens map[common.Address]bool          `json:"favoriteTokens" bson:"favoriteTokens"`
	IsBlocked      bool                             `json:"isBlocked" bson:"isBlocked"`
	CreatedAt      time.Time                        `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time                        `json:"updatedAt" bson:"updatedAt"`
}

// GetBSON implements bson.Getter
func (a *Account) GetBSON() (interface{}, error) {
	ar := AccountRecord{
		IsBlocked: a.IsBlocked,
		Address:   a.Address.Hex(),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	tokenBalances := make(map[string]TokenBalanceRecord)

	for key, value := range a.TokenBalances {
		tokenBalances[key.Hex()] = TokenBalanceRecord{
			Address:          value.Address.Hex(),
			Symbol:           value.Symbol,
			Balance:          value.Balance.String(),
			InOrderBalance:   value.InOrderBalance.String(),
			AvailableBalance: value.AvailableBalance.String(),
		}
	}

	ar.TokenBalances = tokenBalances

	favoriteTokens := make(map[string]bool)

	for key, value := range a.FavoriteTokens {
		favoriteTokens[key.Hex()] = value
	}

	ar.FavoriteTokens = favoriteTokens

	if a.ID.Hex() == "" {
		ar.ID = bson.NewObjectId()
	} else {
		ar.ID = a.ID
	}

	return ar, nil
}

// SetBSON implemenets bson.Setter
func (a *Account) SetBSON(raw bson.Raw) error {
	decoded := &AccountRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	a.TokenBalances = make(map[common.Address]*TokenBalance)
	for key, value := range decoded.TokenBalances {

		balance := new(big.Int)
		balance, _ = balance.SetString(value.Balance, 10)
		inOrderBalance := new(big.Int)
		inOrderBalance, _ = inOrderBalance.SetString(value.InOrderBalance, 10)
		availableBalance := new(big.Int)
		availableBalance, _ = availableBalance.SetString(value.AvailableBalance, 10)

		a.TokenBalances[common.HexToAddress(key)] = &TokenBalance{
			Address:          common.HexToAddress(value.Address),
			Symbol:           value.Symbol,
			Balance:          balance,
			InOrderBalance:   inOrderBalance,
			AvailableBalance: availableBalance,
		}
	}

	a.FavoriteTokens = make(map[common.Address]bool)
	for key, value := range decoded.FavoriteTokens {
		a.FavoriteTokens[common.HexToAddress(key)] = value
	}

	a.Address = common.HexToAddress(decoded.Address)
	a.ID = decoded.ID
	a.IsBlocked = decoded.IsBlocked
	a.CreatedAt = decoded.CreatedAt
	a.UpdatedAt = decoded.UpdatedAt

	return nil
}

// JSON Marshal/Unmarshal interface

// MarshalJSON implements the json.Marshal interface
func (a *Account) MarshalJSON() ([]byte, error) {
	logger.Debug(a.FavoriteTokens)

	account := map[string]interface{}{
		"id":        a.ID,
		"address":   a.Address,
		"isBlocked": a.IsBlocked,
		"createdAt": a.CreatedAt.String(),
		"updatedAt": a.UpdatedAt.String(),
	}

	tokenBalance := make(map[string]interface{})

	for address, balance := range a.TokenBalances {
		tokenBalance[address.Hex()] = map[string]interface{}{
			"address":          balance.Address.Hex(),
			"symbol":           balance.Symbol,
			"balance":          balance.Balance.String(),
			"inOrderBalance":   balance.InOrderBalance.String(),
			"availableBalance": balance.AvailableBalance.String(),
		}
	}

	account["tokenBalances"] = tokenBalance

	favoriteTokens := make(map[string]bool)

	for address, isFavorite := range a.FavoriteTokens {
		favoriteTokens[address.Hex()] = isFavorite
	}

	account["favoriteTokens"] = favoriteTokens

	return json.Marshal(account)
}

func (a *Account) UnmarshalJSON(b []byte) error {
	account := map[string]interface{}{}
	err := json.Unmarshal(b, &account)
	if err != nil {
		return err
	}

	if account["id"] != nil && bson.IsObjectIdHex(account["id"].(string)) {
		a.ID = bson.ObjectIdHex(account["id"].(string))
	}

	if account["address"] != nil {
		a.Address = common.HexToAddress(account["address"].(string))
	}

	if account["tokenBalances"] != nil {
		tokenBalances := account["tokenBalances"].(map[string]interface{})
		a.TokenBalances = make(map[common.Address]*TokenBalance)
		for address, balance := range tokenBalances {
			if !common.IsHexAddress(address) {
				continue
			}

			tokenBalance := balance.(map[string]interface{})
			tb := &TokenBalance{}

			if tokenBalance["address"] != nil && common.IsHexAddress(tokenBalance["address"].(string)) {
				tb.Address = common.HexToAddress(tokenBalance["address"].(string))
			}

			if tokenBalance["symbol"] != nil {
				tb.Symbol = tokenBalance["symbol"].(string)
			}

			tb.Balance = new(big.Int)
			tb.InOrderBalance = new(big.Int)
			tb.AvailableBalance = new(big.Int)

			if tokenBalance["balance"] != nil {
				tb.Balance.UnmarshalJSON([]byte(tokenBalance["balance"].(string)))
			}

			if tokenBalance["inOrderBalance"] != nil {
				tb.InOrderBalance.UnmarshalJSON([]byte(tokenBalance["inOrderBalance"].(string)))
			}

			if tokenBalance["availableBalance"] != nil {
				tb.AvailableBalance.UnmarshalJSON([]byte(tokenBalance["availableBalance"].(string)))
			}

			a.TokenBalances[common.HexToAddress(address)] = tb
		}
	}

	if account["favoriteTokens"] != nil {
		favoriteTokens := account["favoriteTokens"].(map[string]interface{})
		for address, isFavorite := range favoriteTokens {
			if !common.IsHexAddress(address) {
				continue
			}

			a.FavoriteTokens[common.HexToAddress(address)] = isFavorite.(bool)
		}
	}

	return nil
}

// Validate enforces the account model
func (a Account) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address, validation.Required),
	)
}

// AccountRecord corresponds to what is stored in the DB. big.Ints are encoded as strings
type AccountRecord struct {
	ID             bson.ObjectId                 `json:"id" bson:"_id"`
	Address        string                        `json:"address" bson:"address"`
	TokenBalances  map[string]TokenBalanceRecord `json:"tokenBalances" bson:"tokenBalances"`
	FavoriteTokens map[string]bool               `json:"favoriteTokens" bson:"favoriteTokens"`
	IsBlocked      bool                          `json:"isBlocked" bson:"isBlocked"`
	CreatedAt      time.Time                     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time                     `json:"updatedAt" bson:"updatedAt"`
}

type AccountBSONUpdate struct {
	*Account
}

func (a *AccountBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()
	tokenBalances := make(map[string]TokenBalanceRecord)

	//TODO validate this. All the fields have to be set
	for key, value := range a.TokenBalances {
		tokenBalances[key.Hex()] = TokenBalanceRecord{
			Address:          value.Address.Hex(),
			Symbol:           value.Symbol,
			Balance:          value.Balance.String(),
			InOrderBalance:   value.InOrderBalance.String(),
			AvailableBalance: value.AvailableBalance.String(),
		}
	}

	set := bson.M{
		"updatedAt": now,
		"address":   a.Address,
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}

// TokenBalance holds the Balance and the Locked balance values for a single Ethereum token
// Balance and Locked Balance are stored as big.Int as they represent uint256 values
type TokenBalance struct {
	Address          common.Address `json:"address" bson:"address"`
	Symbol           string         `json:"symbol" bson:"symbol"`
	Balance          *big.Int       `json:"balance" bson:"balance"`
	AvailableBalance *big.Int       `json:"availableBalance" bson:"availableBalance"`
	InOrderBalance   *big.Int       `json:"inOrderBalance" bson:"inOrderBalance"`
}

// MarshalJSON implements the json.Marshal interface
func (t *TokenBalance) MarshalJSON() ([]byte, error) {
	tb := map[string]interface{}{
		"address":          t.Address.Hex(),
		"symbol":           t.Symbol,
		"balance":          t.Balance.String(),
		"inOrderBalance":   t.InOrderBalance.String(),
		"availableBalance": t.AvailableBalance.String(),
	}

	return json.Marshal(tb)
}

func (t *TokenBalance) UnmarshalJSON(b []byte) error {
	tb := map[string]interface{}{}
	err := json.Unmarshal(b, &tb)
	if err != nil {
		return err
	}

	if tb["address"] != nil {
		t.Address = common.HexToAddress(tb["address"].(string))
	}

	if tb["symbol"] != nil {
		t.Symbol = tb["symbol"].(string)
	}

	t.Balance = new(big.Int)
	t.InOrderBalance = new(big.Int)
	t.AvailableBalance = new(big.Int)

	if tb["balance"] != nil {
		t.Balance.UnmarshalJSON([]byte(tb["balance"].(string)))
	}

	if tb["inOrderBalance"] != nil {
		t.InOrderBalance.UnmarshalJSON([]byte(tb["inOrderBalance"].(string)))
	}

	if tb["availableBalance"] != nil {
		t.AvailableBalance.UnmarshalJSON([]byte(tb["availableBalance"].(string)))
	}

	return nil
}

// TokenBalanceRecord corresponds to a TokenBalance struct that is stored in the DB. big.Ints are encoded as strings
type TokenBalanceRecord struct {
	Address          string `json:"address" bson:"address"`
	Symbol           string `json:"symbol" bson:"symbol"`
	Balance          string `json:"balance" bson:"balance"`
	AvailableBalance string `json:"availableBalance" base:"availableBalance"`
	InOrderBalance   string `json:"inOrderBalance" bson:"inOrderBalance"`
}

type FavoriteTokenRequest struct {
	Address string `json:"address" bson:"address"`
	Token   string `json:"token" bson:"token"`
}
