package types

// Pair struct is used to model the pair data in the system and DB
type PriceBoardData struct {
	Ticks          []*Tick `json:"ticks" bson:"ticks"`
	PriceUSD       string  `json:"usd" bson:"usd"`
	LastTradePrice string  `json:"last_trade_price" bson:"last_trade_price"`
}
