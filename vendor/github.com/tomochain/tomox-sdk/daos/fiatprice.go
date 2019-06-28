package daos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tidwall/gjson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

type FiatPriceDao struct {
	collectionName string
	dbName         string
}

// NewFiatPriceDao returns a new instance of FiatPriceDao.
func NewFiatPriceDao() *FiatPriceDao {
	dbName := app.Config.DBName
	collection := "fiat_price"

	return &FiatPriceDao{
		collectionName: collection,
		dbName:         dbName,
	}
}

func (dao *FiatPriceDao) GetLatestQuotes() (map[string]float64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/cryptocurrency/quotes/latest?symbol=%s&convert=USD", app.Config.CoinmarketcapAPIUrl, app.Config.SupportedCurrencies)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("X-CMC_PRO_API_KEY", app.Config.CoinmarketcapAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	status := gjson.Get(string(body), "status")
	statusErrorCode := status.Get("error_code")
	statusErrorMessage := status.Get("error_message")

	if statusErrorCode.Int() != 0 {
		logger.Error(statusErrorMessage.String())
		return nil, errors.New(statusErrorMessage.String())
	}

	data := gjson.Get(string(body), "data")
	result := make(map[string]float64)
	data.ForEach(func(key, value gjson.Result) bool {
		result[key.String()] = value.Get("quote.USD.price").Float()
		return true // keep iterating
	})

	return result, nil
}

func (dao *FiatPriceDao) GetCoinMarketChart(id string, vsCurrency string, days string) (*types.CoinsIDMarketChart, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=%s&days=%s", app.Config.CoingeckoAPIUrl, id, vsCurrency, days)

	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var data *types.CoinsIDMarketChart

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (dao *FiatPriceDao) GetCoinMarketChartRange(id string, vsCurrency string, from int64, to int64) (*types.CoinsIDMarketChart, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/coins/%s/market_chart/range?vs_currency=%s&from=%d&to=%d", app.Config.CoingeckoAPIUrl, id, vsCurrency, from, to)

	req, err := http.NewRequest("GET", url, nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var data *types.CoinsIDMarketChart

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Get24hChart gets price chart of symbol by fiatCurrency in 24h
// It's not guaranteed in exact 24h because we are using data from Coingecko
func (dao *FiatPriceDao) Get24hChart(symbol, fiatCurrency string) ([]*types.FiatPriceItem, error) {
	var res []*types.FiatPriceItem
	q := bson.M{"symbol": symbol, "fiatCurrency": fiatCurrency}

	limit := 24

	err := db.GetAndSort(dao.dbName, dao.collectionName, q, []string{"-timestamp"}, 0, limit, &res)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return []*types.FiatPriceItem{}, nil
	}

	return res, nil
}

// Create function performs the DB insertion task for fiat_price collection
// It accepts 1 or more fiat price items as input.
// All the fiat price items are inserted in one query itself.
func (dao *FiatPriceDao) Create(items ...*types.FiatPriceItem) error {
	y := make([]interface{}, len(items))

	for _, item := range items {
		y = append(y, item)
	}

	err := db.Create(dao.dbName, dao.collectionName, y...)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (dao *FiatPriceDao) FindAndModify(symbol, fiatCurrency, timestamp string, i *types.FiatPriceItem) (*types.FiatPriceItem, error) {
	query := bson.M{
		"symbol":       symbol,
		"fiatCurrency": fiatCurrency,
		"timestamp":    timestamp,
	}
	updated := &types.FiatPriceItem{}
	change := mgo.Change{
		Update:    types.FiatPriceItemBSONUpdate{FiatPriceItem: i},
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	err := db.FindAndModify(dao.dbName, dao.collectionName, query, change, &updated)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

func (dao *FiatPriceDao) Upsert(symbol, fiatCurrency, timestamp string, i *types.FiatPriceItem) error {
	_, err := db.Upsert(dao.dbName, dao.collectionName, bson.M{
		"symbol":       symbol,
		"fiatCurrency": fiatCurrency,
		"timestamp":    timestamp,
	}, i)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Aggregate function calls the aggregate pipeline of mongodb
func (dao *FiatPriceDao) Aggregate(q []bson.M) ([]*types.FiatPriceItem, error) {
	var res []*types.FiatPriceItem

	err := db.Aggregate(dao.dbName, dao.collectionName, q, &res)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

// Drop drops all the order documents in the current database
func (dao *FiatPriceDao) Drop() {
	db.DropCollection(dao.dbName, dao.collectionName)
}
