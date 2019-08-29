package types

import (
	"github.com/globalsign/mgo/bson"
)

type MarketData struct {
	PairData        []*PairData                 `json:"pairData" bson:"pairData"`
	SmallChartsData map[string][]*FiatPriceItem `json:"smallChartsData" bson:"smallChartsData"`
}

type ChartItem [2]float64

type CoinsIDMarketChart struct {
	Prices       []*ChartItem `json:"prices"`
	MarketCaps   []*ChartItem `json:"market_caps"`
	TotalVolumes []*ChartItem `json:"total_volumes"`
}

type FiatPriceItem struct {
	Symbol       string `json:"-" bson:"symbol"`
	Price        string `json:"price" bson:"price"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp"`
	FiatCurrency string `json:"fiatCurrency" bson:"fiatCurrency"`
	TotalVolume  string `json:"totalVolume" bson:"totalVolume"`
}

type FiatPriceItemBSONUpdate struct {
	*FiatPriceItem
}

func (i *FiatPriceItem) GetBSON() (interface{}, error) {
	nr := FiatPriceItem{
		Symbol:       i.Symbol,
		Price:        i.Price,
		Timestamp:    i.Timestamp,
		FiatCurrency: i.FiatCurrency,
		TotalVolume:  i.TotalVolume,
	}

	return nr, nil
}

func (i *FiatPriceItem) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Symbol       string `json:"symbol" bson:"symbol"`
		Price        string `json:"price" bson:"price"`
		Timestamp    int64  `json:"timestamp" bson:"timestamp"`
		FiatCurrency string `json:"fiatCurrency" bson:"fiatCurrency"`
		TotalVolume  string `json:"totalVolume" bson:"totalVolume"`
	})

	err := raw.Unmarshal(decoded)

	if err != nil {
		return err
	}

	i.Symbol = decoded.Symbol
	i.Price = decoded.Price
	i.Timestamp = decoded.Timestamp
	i.FiatCurrency = decoded.FiatCurrency
	i.TotalVolume = decoded.TotalVolume

	return nil
}
