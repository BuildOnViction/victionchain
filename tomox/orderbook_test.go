package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"
)
var testPairName = "aaa/tomo"
var sampleTestOrder = &OrderItem{
	ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
	UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
	BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
	QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
	Status:          "New",
	Side:            Ask,
	Type:            Limit,
	PairName:        testPairName,
	Hash:            common.StringToHash(string(rand.Intn(1000))),
	Signature: &Signature{
		V: 1,
		R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
		S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
	},
	FilledAmount: new(big.Int).SetUint64(0),
	Nonce:        new(big.Int).SetUint64(1),
	CreatedAt:    time.Now(),
	UpdatedAt:    time.Now(),
}
var testingBlockHash = common.HexToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0")
func initTestOrderBook(testDir, pairName string) *OrderBook {
	tomox := &TomoX{
		Orderbooks: map[string]*OrderBook{},
		activePairs: map[string]bool{},
		db: NewLDBEngine(&Config{
			DataDir:  testDir,
			DBEngine: "leveldb",
		}),
	}
	ob, _ := tomox.GetOrderBook(pairName, true, testingBlockHash)
	return ob
}

//  this order doesn't match any order existing in orderTree
// inserting it to bid/ask Tree accordingly
func TestOrderBook_ProcessLimitOrder_InsertToOrderTree(t *testing.T) {
	testDir := "TestOrderBook_ProcessLimitOrder_InsertToOrderTree"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Price = big.NewInt(1000) // ask order, price = 1000
	order1.Quantity = big.NewInt(1000)
	trades, _, err := ob.processLimitOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// process one more ask order
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Price = big.NewInt(2000) // ask order, price = 2000
	order2.Quantity = big.NewInt(1000)
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order2)
	}

	// process one bid order
	order3 := &OrderItem{}
	*order3 = *sampleTestOrder
	order3.Side = Bid
	order3.Price = big.NewInt(500) // bid order, price = 500
	order3.Quantity = big.NewInt(1000)
	trades, _, err = ob.ProcessOrder(order3, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order3)
	}

	// check size of orderTree
	if ob.Bids.PriceTree.Size() != 1 || ob.Asks.PriceTree.Size() != 2 {
		t.Error("Wrong priceTree size")
		t.Log("Expected Bid tree, size = 1", "actual", ob.Bids.PriceTree.Size())
		t.Log("Expected Ask tree, size = 2", "actual", ob.Asks.PriceTree.Size())
	}
}

// this order matches one order of orderTree, order.quantity < orderList.headOrder.Item.Quantity
// as a result, after matching, quantityToTrade = 0, update quantity of headOrder
// refer this case tomox/orderbook.go:357:
// IsStrictlySmallerThan(quantityToTrade, headOrder.Item.Quantity)
func TestOrderBook_ProcessLimitOrder_OneToOneMatching_FullMatching_Case1(t *testing.T) {
	testDir := "TestOrderBook_ProcessLimitOrder_OneToOneMatching_FullMatching_Case1"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(100) // ask order, price = 100
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// process one bid order to match the above order
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Bid
	order2.Quantity = big.NewInt(900)
	order2.Price = big.NewInt(101) // bid order, price = 101
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	// as a result of this matching process, we expect txMatch.Trades has one item
	// quantityToTrade is zero and order1 of askTree (tomox/orderbook_test.go:113) should remain quantity = 100
	if len(trades) != 1 {
		t.Error("Wrong number of trades, expected: 1", "actual", len(trades))
	}

	// check tradedQuantity
	if trades[0]["quantity"] != "900" {
		t.Error("Wrong tradedQuantity, expected: 900", "actual", trades[0]["quantity"])
	}
	remainingOrder := ob.Asks.GetOrder(GetKeyFromBig(big.NewInt(1)), big.NewInt(100), true, testingBlockHash)
	if remainingOrder.Item.Quantity.Cmp(big.NewInt(100)) != 0 {
		t.Error("Wrong remaining quantity")
	}
}


// this order matches one order of orderTree, order.quantity == orderList.headOrder.Item.Quantity
// as a result, after matching, quantityToTrade = 0, remove headOrder from orderList
// refer this case tomox/orderbook.go:365
// IsEqual(quantityToTrade, headOrder.Item.Quantity)
func TestOrderBook_ProcessLimitOrder_OneToOneMatching_FullMatching_Case2(t *testing.T) {
	testDir := "TestOrderBook_ProcessLimitOrder_OneToOneMatching_FullMatching_Case2"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(100) // ask order, price = 100
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// process one bid order to match the above order
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Bid
	order2.Quantity = big.NewInt(1000) // same as order1's quantity
	order2.Price = big.NewInt(101) // bid order, price = 101
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	// as a result of this matching process, we expect txMatch.Trades has one item
	// quantityToTrade is zero and order1 of askTree (tomox/orderbook_test.go:113) should remain quantity = 100
	if len(trades) != 1 {
		t.Error("Wrong number of trades, expected: 1", "actual", len(trades))
	}

	// check tradedQuantity
	if trades[0]["quantity"] != "1000" {
		t.Error("Wrong tradedQuantity, expected: 1000", "actual", trades[0]["quantity"])
	}
	if remainingOrder := ob.Asks.GetOrder(GetKeyFromBig(big.NewInt(1)), big.NewInt(100), true, testingBlockHash); remainingOrder != nil {
		t.Error("Expected: fully matched with order in the orderTree")
	}
}

// this order matches one order of orderTree, order.quantity > orderList.headOrder.Item.Quantity
// as a result, after matching, quantityToTrade = Sub(quantityToTrade, tradedQuantity), remove headOrder from orderList
// refer this case tomox/orderbook.go:378
func TestOrderBook_ProcessLimitOrder_OneToOneMatching_PartialMatching(t *testing.T) {
	testDir := "TestOrderBook_ProcessLimitOrder_OneToOneMatching_PartialMatching"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(100) // ask order, price = 100
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// process one bid order to match the above order
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Bid
	order2.Quantity = big.NewInt(1200)
	order2.Price = big.NewInt(101) // bid order, price = 101
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	// as a result of this matching process, we expect txMatch.Trades has one item
	// quantityToTrade is zero and order1 of askTree (tomox/orderbook_test.go:113) should remain quantity = 100
	if len(trades) != 1 {
		t.Error("Wrong number of trades, expected: 1", "actual", len(trades))
	}

	// check tradedQuantity
	if trades[0]["quantity"] != "1000" {
		t.Error("Wrong tradedQuantity, expected: 1000", "actual", trades[0]["quantity"])
	}
	remainingOrder := ob.Bids.GetOrder(GetKeyFromBig(big.NewInt(2)), big.NewInt(101), true, testingBlockHash)
	if remainingOrder.Item.Quantity.Cmp(big.NewInt(200)) != 0 {
		t.Error("Wrong remaining quantity")
	}
}

// this order matches many orders existing on the orderTree
// in this case, txMatch.Trades should have multiple items
func TestOrderBook_ProcessLimitOrder_OneToManyMatching(t *testing.T) {
	testDir := "TestOrderBook_ProcessLimitOrder_OneToManyMatching"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(98) // ask order, price = 98
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// one more askOrder
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Ask
	order2.Quantity = big.NewInt(1000)
	order2.Price = big.NewInt(99) // ask order, price = 99
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order2)
	}

	// process a bidOrder which completely matches order1, partially matches order2
	order3 := &OrderItem{}
	*order3 = *sampleTestOrder
	order3.Side = Bid
	order3.Quantity = big.NewInt(1600)
	order3.Price = big.NewInt(100)
	trades, _, err = ob.ProcessOrder(order3, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) != 2 {
		t.Error("Expected: 2 trades, actual", len(trades))
		t.Log(order3)
	}
	// check tradedQuantity
	if trades[0]["quantity"] != "1000" && trades[1]["quantity"] != "600" {
		t.Error("Wrong trade quantity")
		t.Log("Expected: Trade[0][quantity] = 1000", "actual", trades[0]["quantity"])
		t.Log("Expected: Trade[1][quantity] = 600", "actual", trades[1]["quantity"])
	}
	remainingOrder := ob.Asks.GetOrder(GetKeyFromBig(big.NewInt(2)), big.NewInt(99), true, testingBlockHash)
	if remainingOrder.Item.Quantity.Cmp(big.NewInt(400)) != 0 {
		t.Error("Wrong remaining quantity. Expected: 400. Actual:", remainingOrder.Item.Quantity.Uint64())
	}
}

// as a result of processMarketOrder, either quantityToTrade of this order is zero or the orderTree will be empty
// this testcase verify the case which marketOrder matches all orders of orderTree, orderTree will be empty
func TestOrderBook_ProcessMarketOrder_FullMatching(t *testing.T) {
	testDir := "TestOrderBook_ProcessMarketOrder_FullMatching"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(98) // ask order, price = 98
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// one more askOrder
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Ask
	order2.Quantity = big.NewInt(1000)
	order2.Price = big.NewInt(99) // ask order, price = 99
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order2)
	}

	// process a bidOrder which completely matches order1, partially matches order2
	order3 := &OrderItem{}
	*order3 = *sampleTestOrder
	order3.Side = Bid
	order3.Quantity = big.NewInt(2500)
	order3.Type = Market
	trades, _, err = ob.ProcessOrder(order3, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) != 2 {
		t.Error("Expected: 2 trades, actual", len(trades))
		t.Log(order3)
	}
	// check tradedQuantity
	if trades[0]["quantity"] != "1000" && trades[1]["quantity"] != "1000" {
		t.Error("Wrong trade quantity")
		t.Log("Expected: Trade[0][quantity] = 1000", "actual", trades[0]["quantity"])
		t.Log("Expected: Trade[1][quantity] = 1000", "actual", trades[1]["quantity"])
	}
	// check size of orderTree
	if ob.Bids.PriceTree.Size() != 0 || ob.Asks.PriceTree.Size() != 0 {
		t.Error("Wrong priceTree size. Expected: empty", "actual_bidTree_size", ob.Bids.PriceTree.Size(), "actual_askTree_size", ob.Asks.PriceTree.Size())
	}
}


// as a result of processMarketOrder, either quantityToTrade of this order is zero or the orderTree will be empty
// this testcase verify the case which marketOrder matches some orders of orderTree
// After matching, orderTree is not empty, quantityToTrade = 0
func TestOrderBook_ProcessMarketOrder_PartialMatching(t *testing.T) {
	testDir := "TestOrderBook_ProcessMarketOrder_PartialMatching"
	defer os.RemoveAll(testDir)
	ob := initTestOrderBook(testDir, testPairName)

	order1 := &OrderItem{}
	*order1 = *sampleTestOrder
	order1.Quantity = big.NewInt(1000)
	order1.Price = big.NewInt(98) // ask order, price = 98
	trades, _, err := ob.ProcessOrder(order1, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order1)
	}

	// one more askOrder
	order2 := &OrderItem{}
	*order2 = *sampleTestOrder
	order2.Side = Ask
	order2.Quantity = big.NewInt(1000)
	order2.Price = big.NewInt(99) // ask order, price = 99
	trades, _, err = ob.ProcessOrder(order2, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) > 0 {
		t.Error("Trades should be empty. Order should be inserted to orderTree")
		t.Log(order2)
	}

	// process a bidOrder which completely matches order1, partially matches order2
	order3 := &OrderItem{}
	*order3 = *sampleTestOrder
	order3.Side = Bid
	order3.Quantity = big.NewInt(1300)
	order3.Type = Market
	trades, _, err = ob.ProcessOrder(order3, true, true, testingBlockHash)
	if err != nil {
		t.Error(err.Error())
	}
	if len(trades) != 2 {
		t.Error("Expected: 2 trades, actual", len(trades))
		t.Log(order3)
	}
	// check tradedQuantity
	if trades[0]["quantity"] != "1000" && trades[1]["quantity"] != "300" {
		t.Error("Wrong trade quantity")
		t.Log("Expected: Trade[0][quantity] = 1000", "actual", trades[0]["quantity"])
		t.Log("Expected: Trade[1][quantity] = 300", "actual", trades[1]["quantity"])
	}
	// verify remaining orders
	remainingOrder := ob.Asks.GetOrder(GetKeyFromBig(big.NewInt(2)), big.NewInt(99), true, testingBlockHash)
	if remainingOrder.Item.Quantity.Cmp(big.NewInt(700)) != 0 {
		t.Error("Wrong remaining quantity. Expected: 700. Actual:", remainingOrder.Item.Quantity.Uint64())
	}
}