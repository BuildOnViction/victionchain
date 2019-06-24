package tomox

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"
)

var ether = big.NewInt(1000000000000000000)

func prepareOrderbookData(pair string, db OrderDao) (*OrderBook, error) {
	var (
		ob  *OrderBook
		err error
	)
	v := []byte(string(rand.Intn(999)))

	ob = NewOrderBook(pair, db)

	// insert order to bid tree: price 99
	price := CloneBigInt(ether)
	err = ob.Bids.InsertOrder(&OrderItem{
		OrderID:         uint64(1),
		Quantity:        big.NewInt(100),
		Price:           price.Mul(price, big.NewInt(99)),
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            Bid,
		Type:            "LO",
		PairName:        "aaa/tomo",
		Hash:            common.StringToHash(string(rand.Intn(1000))),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        new(big.Int).SetUint64(1),
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	})

	// insert order to bid tree: price 98
	price = CloneBigInt(ether)
	err = ob.Bids.InsertOrder(&OrderItem{
		OrderID:         uint64(2),
		Quantity:        big.NewInt(50),
		Price:           price.Mul(price, big.NewInt(98)),
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            Bid,
		Type:            "LO",
		PairName:        "aaa/tomo",
		Hash:            common.StringToHash(string(rand.Intn(1000))),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        new(big.Int).SetUint64(1),
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	})
	if err != nil {
		return ob, err
	}

	// insert order to ask tree: price 101
	price = CloneBigInt(ether)
	err = ob.Asks.InsertOrder(&OrderItem{
		OrderID:         uint64(3),
		Quantity:        big.NewInt(200),
		Price:           price.Mul(price, big.NewInt(101)),
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            Ask,
		Type:            "LO",
		PairName:        "aaa/tomo",
		Hash:            common.StringToHash(string(rand.Intn(1000))),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        new(big.Int).SetUint64(1),
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	})

	// insert order to ask tree: price 102
	price = CloneBigInt(ether)
	err = ob.Asks.InsertOrder(&OrderItem{
		OrderID:         uint64(4),
		Quantity:        big.NewInt(300),
		Price:           price.Mul(price, big.NewInt(102)),
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            Ask,
		Type:            "LO",
		PairName:        "aaa/tomo",
		Hash:            common.StringToHash(string(rand.Intn(1000))),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        new(big.Int).SetUint64(1),
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	})
	return ob, nil
}

/**
*	Test scenario:
*	- create a new snapshot
*	- store snapshot
*	- load snapshot from disk and compare
**/
func TestTomoX_Snapshot(t *testing.T) {
	testDir := "TestTomoX_Snapshot"

	tomox := &TomoX{
		Orderbooks: map[string]*OrderBook{},
		db: NewLDBEngine(&Config{
			DataDir:  testDir,
			DBEngine: "leveldb",
		}),
	}
	defer os.RemoveAll(testDir)
	blockHash := common.StringToHash("aaa")
	pair := "aaa/tomo"

	// orderbooks["aaa/tomo"] has 2 bids orders and 2 asks orders
	//- order 1: buy 100 aaa, price 99 tomo
	//- order 2: buy 50 aaa, price 98 tomo
	//- order 3: sell 200 aaa, price 101 tomo
	//- order 4: sell 300 aaa, price 102 tomo
	ob, err := prepareOrderbookData(pair, tomox.db)
	if err != nil {
		t.Error("Failed to create orderbook", err)
	}
	tomox.Orderbooks[pair] = ob

	if err := tomox.Snapshot(blockHash); err != nil {
		t.Error("Failed to store snapshot", "err", err, "blockHash", blockHash)
	}

	// load snapshot with invalid hash
	newSnap, err := loadSnapshot(tomox.db, common.StringToHash("xxx"))
	if err == nil {
		t.Error("Expected an error due to wrong hash")
	}

	newSnap, err = loadSnapshot(tomox.db, blockHash)
	if err != nil {
		t.Error("Failed to load snapshot", err)
	}

	// verify snapshot hash
	if newSnap.Hash != blockHash {
		t.Error("Wrong snapshot hash", "expected", blockHash, "actual", newSnap.Hash)
	}

	// load orderbook of an invalid pair
	if _, ok := newSnap.OrderBooks["btc/tomo"]; ok {
		t.Error("Expected an error due to wrong pair")
	}

	// verify orderbook hash
	hash, err := ob.Hash()
	if err != nil {
		t.Error(err)
	}
	var newOb *OrderBook
	newOb, err = newSnap.RestoreOrderBookFromSnapshot(tomox.db, pair)
	if err != nil {
		t.Error("Failed to restore orderbook from snapshot", err)
	}

	newHash, err := newOb.Hash()
	if err != nil {
		t.Error(err)
	}
	if err != nil || newHash != hash {
		t.Error("Wrong orderbook hash", "expected", hash, "actual", newHash)
	}

	var (
		bidTree, askTree *OrderTree
		treeHash         common.Hash
	)

	// verify bid tree
	hash, err = ob.Bids.Hash()

	bidTree = newOb.Bids
	treeHash, err = bidTree.Hash()
	if err != nil || treeHash != hash {
		t.Error("Wrong bid tree hash", "expected", hash, "actual", treeHash)
	}

	// verify ask tree
	hash, err = ob.Asks.Hash()
	askTree = newOb.Asks
	treeHash, err = askTree.Hash()
	if err != nil || treeHash != hash {
		t.Error("Wrong bid tree hash", "expected", hash, "actual", treeHash)
	}

	// verify bid order, orderId = 1, price = 99
	price := CloneBigInt(ether)
	price = price.Mul(price, big.NewInt(99))
	order := bidTree.GetOrder(GetKeyFromBig(big.NewInt(1)), price)
	if order.Item.Quantity.Cmp(big.NewInt(100)) != 0 {
		t.Error("Wrong order item", "expected quantity", 100, "actual quantity", order.Item.Quantity.Uint64())
	}

	// verify bid  order, orderId = 2, price = 98
	price = CloneBigInt(ether)
	price = price.Mul(price, big.NewInt(98))
	order = bidTree.GetOrder(GetKeyFromBig(big.NewInt(2)), price)
	if order.Item.Quantity.Cmp(big.NewInt(50)) != 0 {
		t.Error("Wrong order item", "expected quantity", 50, "actual quantity", order.Item.Quantity.Uint64())
	}

	// verify ask order, orderId =3, price = 101
	price = CloneBigInt(ether)
	price = price.Mul(price, big.NewInt(101))
	order = askTree.GetOrder(GetKeyFromBig(big.NewInt(3)), price)
	if order.Item.Quantity.Cmp(big.NewInt(200)) != 0 {
		t.Error("Wrong order item", "expected quantity", 200, "actual quantity", order.Item.Quantity.Uint64())
	}

	// verify ask order, orderId = 4, price = 102
	price = CloneBigInt(ether)
	price = price.Mul(price, big.NewInt(102))
	order = askTree.GetOrder(GetKeyFromBig(big.NewInt(4)), price)
	if order.Item.Quantity.Cmp(big.NewInt(300)) != 0 {
		t.Error("Wrong order item", "expected quantity", 300, "actual quantity", order.Item.Quantity.Uint64())
	}
}
