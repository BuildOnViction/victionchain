package protocol

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/tomox"
)

// remember that API structs to be offered MUST be exported
type OrderbookAPI struct {
	V      int
	Engine *tomox.Engine
}

// Version : return version
func (api *OrderbookAPI) Version() (int, error) {
	return api.V, nil
}

func NewOrderbookAPI(v int, orderbookEngine *tomox.Engine) *OrderbookAPI {
	return &OrderbookAPI{
		V:      v,
		Engine: orderbookEngine,
	}
}

func (api *OrderbookAPI) getRecordFromOrder(order *tomox.Order, ob *tomox.OrderBook) map[string]string {
	record := make(map[string]string)
	record["timestamp"] = strconv.FormatUint(order.Item.Timestamp, 10)
	record["price"] = order.Item.Price.String()
	record["quantity"] = order.Item.Quantity.String()
	// retrieve the input order_id, by default it is set when retrieving from orderbook
	record["order_id"] = new(big.Int).SetBytes(order.Key).String()
	record["trade_id"] = order.Item.TradeID
	return record
}

func (api *OrderbookAPI) GetBestAskList(pairName string) []map[string]string {
	ob, _ := api.Engine.GetOrderBook(pairName)
	if ob == nil {
		return nil
	}
	orderList := ob.Asks.MaxPriceList()
	if orderList == nil {
		return nil
	}

	// t.Logf("Best ask List : %s", orderList.String(0))
	cursor := orderList.Head()
	// we have length
	results := make([]map[string]string, orderList.Item.Length)
	for cursor != nil {
		record := api.getRecordFromOrder(cursor, ob)
		results = append(results, record)
		cursor = cursor.GetNextOrder(orderList)
	}
	return results
}

func (api *OrderbookAPI) GetBestBidList(pairName string) []map[string]string {
	ob, _ := api.Engine.GetOrderBook(pairName)
	if ob == nil {
		return nil
	}
	orderList := ob.Bids.MinPriceList()
	// t.Logf("Best ask List : %s", orderList.String(0))
	if orderList == nil {
		return nil
	}
	cursor := orderList.Tail()
	// we have length
	results := make([]map[string]string, orderList.Item.Length)
	for cursor != nil {
		record := api.getRecordFromOrder(cursor, ob)
		results = append(results, record)
		cursor = cursor.GetPrevOrder(orderList)
	}
	return results

}

func (api *OrderbookAPI) GetOrder(pairName, orderID string) map[string]string {
	var result map[string]string
	ob, _ := api.Engine.GetOrderBook(pairName)
	if ob == nil {
		return nil
	}
	key := tomox.GetKeyFromString(orderID)
	order := ob.GetOrder(key)
	if order != nil {
		result = api.getRecordFromOrder(order, ob)
	}
	return result
}
