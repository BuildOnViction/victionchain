package types

type OrderBook struct {
	PairName string              `json:"pairName"`
	Asks     []map[string]string `json:"asks"`
	Bids     []map[string]string `json:"bids"`
}

type RawOrderBook struct {
	PairName string   `json:"pairName"`
	Orders   []*Order `json:"orders"`
}
