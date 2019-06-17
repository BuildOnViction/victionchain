package tomox

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"
)

func buildOrder(nonce *big.Int) *OrderItem {
	rand.Seed(time.Now().UTC().UnixNano())
	v := []byte(string(rand.Intn(999)))
	lstBuySell := []string{"BUY", "SELL"}
	order := &OrderItem{
		Quantity:        new(big.Int).SetUint64(uint64(rand.Intn(10)) * 1000000000000000000),
		Price:           new(big.Int).SetUint64(uint64(rand.Intn(10)) * 100000000000000000),
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            lstBuySell[rand.Int()%len(lstBuySell)],
		Type:            "LO",
		PairName:        "0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd" + "::" + "0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd",
		//Hash:            common.StringToHash("0xdc842ea4a239d1a4e56f1e7ba31aab5a307cb643a9f5b89f972f2f5f0d1e7587"),
		Hash: common.StringToHash(nonce.String()),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        nonce,
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	}
	return order
}

func testCreateOrder(t *testing.T, nonce *big.Int) {
	order := buildOrder(nonce)
	topic := order.BaseToken.Hex() + "::" + order.QuoteToken.Hex()
	encodedTopic := fmt.Sprintf("0x%s", hex.EncodeToString([]byte(topic)))
	fmt.Println("topic: ", encodedTopic)

	ipaddress := "0.0.0.0"
	url := fmt.Sprintf("http://%s:8501", ipaddress)

	//create topic
	rpcClient, err := rpc.DialHTTP(url)
	defer rpcClient.Close()
	if err != nil {
		t.Error("rpc.DialHTTP failed", "err", err)
	}
	var result interface{}
	params := make(map[string]interface{})
	params["topic"] = encodedTopic
	err = rpcClient.Call(&result, "tomoX_newTopic", params)
	if err != nil {
		t.Error("rpcClient.Call tomoX_newTopic failed", "err", err)
	}

	//create new order
	params["payload"], err = json.Marshal(order)
	if err != nil {
		t.Error("json.Marshal failed", "err", err)
	}

	err = rpcClient.Call(&result, "tomoX_createOrder", params)
	if err != nil {
		t.Error("rpcClient.Call tomoX_createOrder failed", "err", err)
	}
}

func TestCreate10Orders(t *testing.T) {
	for i := 1; i <= 20; i++ {
		testCreateOrder(t, new(big.Int).SetUint64(uint64(i)))
		time.Sleep(1 * time.Second)
	}
}

func TestCancelOrder(t *testing.T) {
	order := buildOrder(new(big.Int).SetInt64(1))
	topic := order.BaseToken.Hex() + "::" + order.QuoteToken.Hex()
	encodedTopic := fmt.Sprintf("0x%s", hex.EncodeToString([]byte(topic)))
	fmt.Println("topic: ", encodedTopic)

	ipaddress := "0.0.0.0"
	url := fmt.Sprintf("http://%s:8501", ipaddress)

	//cancel order
	rpcClient, err := rpc.DialHTTP(url)
	defer rpcClient.Close()
	if err != nil {
		t.Error("rpc.DialHTTP failed", "err", err)
	}
	var result interface{}
	params := make(map[string]interface{})
	params["topic"] = encodedTopic
	params["payload"], err = json.Marshal(order)
	if err != nil {
		t.Error("json.Marshal failed", "err", err)
	}

	err = rpcClient.Call(&result, "tomoX_cancelOrder", params)
	if err != nil {
		t.Error("rpcClient.Call tomoX_createOrder failed", "err", err)
	}
}

func TestDBPending(t *testing.T) {
	testDir := "TestDBPending"

	tomox := &TomoX{
		Orderbooks:  map[string]*OrderBook{},
		activePairs: make(map[string]bool),
		db: NewLDBEngine(&Config{
			DataDir:  testDir,
			DBEngine: "leveldb",
		}),
	}
	defer os.RemoveAll(testDir)

	if pHashes := tomox.getPendingHashes(); len(pHashes) != 0 {
		t.Error("Expected: no pending hash", "Actual:", len(pHashes))
	}

	var hash common.Hash
	hash = common.StringToHash("0x0000000000000000000000000000000000000000")
	tomox.addPendingHash(hash)
	hash = common.StringToHash("0x0000000000000000000000000000000000000001")
	tomox.addPendingHash(hash)
	hash = common.StringToHash("0x0000000000000000000000000000000000000002")
	tomox.addPendingHash(hash)
	if pHashes := tomox.getPendingHashes(); len(pHashes) != 3 {
		t.Error("Expected: 3 pending hash", "Actual:", len(pHashes))
	}

	// Test remove hash
	hash = common.StringToHash("0x0000000000000000000000000000000000000002")
	tomox.removePendingHash(hash)
	if pHashes := tomox.getPendingHashes(); len(pHashes) != 2 {
		t.Error("Expected: 2 pending hash", "Actual:", len(pHashes))
	}

	order := buildOrder(new(big.Int).SetInt64(1))
	tomox.addOrderPending(order)
	od := tomox.getOrderPending(order.Hash)
	if order.Hash.String() != od.Hash.String() {
		t.Error("Fail to add order pending", "orderOld", order, "orderNew", od)
	}
}

func TestTomoX_GetActivePairs(t *testing.T) {
	testDir := "TestTomoX_GetActivePairs"

	tomox := &TomoX{
		Orderbooks:  map[string]*OrderBook{},
		activePairs: make(map[string]bool),
		db: NewLDBEngine(&Config{
			DataDir:  testDir,
			DBEngine: "leveldb",
		}),
	}
	defer os.RemoveAll(testDir)

	if pairs := tomox.listTokenPairs(); len(pairs) != 0 {
		t.Error("Expected: no active pair", "Actual:", len(pairs))
	}

	savedPairs := map[string]bool{}
	savedPairs["xxx/tomo"] = true
	savedPairs["aaa/tomo"] = true
	if err := tomox.updatePairs(savedPairs); err != nil {
		t.Error("Failed to save active pairs", err)
	}

	// a node has just been restarted, haven't inserted any order yet
	// in memory: there is no activePairsKey
	// in db: there are 2 active pairs
	// expected: tomox.listTokenPairs return 2
	tomox.activePairs = map[string]bool{} // reset tomox.activePairs to simulate the case: a node was restarted
	if pairs := tomox.listTokenPairs(); len(pairs) != 2 {
		t.Error("Expected: 2 active pairs", "Actual:", len(pairs))
	}

	// a node has just been restarted, then insert an order of "aaa/tomo"
	// in db: there are 2 active pairs
	// expected: tomox.listTokenPairs return 2
	tomox.activePairs = map[string]bool{} // reset tomox.activePairsKey to simulate the case: a node was restarted
	tomox.GetOrderBook("aaa/tomo")
	if pairs := tomox.listTokenPairs(); len(pairs) != 2 {
		t.Error("Expected: 2 active pairs", "Actual:", len(pairs))
	}

	// insert an order of existing pair: xxx/tomo
	// expected: tomox.listTokenPairs return 2 pairs
	tomox.GetOrderBook("xxx/tomo")
	if pairs := tomox.listTokenPairs(); len(pairs) != 2 {
		t.Error("Expected: 2 active pairs", "Actual:", len(pairs))
	}

	// now, activePairsKey in tomox.activePairsKey and db are same
	// try to add one more pair to orderbook
	tomox.GetOrderBook("xxx/tomo")
	tomox.GetOrderBook("yyy/tomo")

	if pairs := tomox.listTokenPairs(); len(pairs) != 3 {
		t.Error("Expected: 3 active pairs", "Actual:", len(pairs))
	}
}

func TestEncodeDecodeTXMatch(t *testing.T) {
	var trades []map[string]string
	var txMatches map[common.Hash]TxDataMatch
	var decodeMatches map[common.Hash]TxDataMatch

	transactionRecord := make(map[string]string)
	transactionRecord["price"] = new(big.Int).SetUint64(uint64(25) * 100000000000000000).String()
	transactionRecord["quantity"] = new(big.Int).SetUint64(uint64(12) * 1000000000000000000).String()
	trades = append(trades, transactionRecord)

	transactionRecord = make(map[string]string)
	transactionRecord["price"] = new(big.Int).SetUint64(uint64(14) * 1000000000000000000).String()
	transactionRecord["quantity"] = new(big.Int).SetUint64(uint64(15) * 1000000000000000000).String()
	trades = append(trades, transactionRecord)

	order := buildOrder(new(big.Int).SetInt64(1))
	value, err := EncodeBytesItem(order)
	if err != nil {
		t.Error("Can't encode", "order", order, "err", err)
	}
	txMatches = make(map[common.Hash]TxDataMatch)
	txMatches[order.Hash] = TxDataMatch{
		order:  value,
		trades: trades,
	}
	encode, err := json.Marshal(txMatches)
	if err != nil {
		t.Error("Fail to marshal txMatches", "err", err)
	}

	err = json.Unmarshal(encode, &decodeMatches)
	if err != nil {
		t.Error("Fail to unmarshal txMatches", "err", err)
	}

	if _, ok := decodeMatches[order.Hash]; !ok {
		t.Error("marshal and unmarshal txMatches not valid", "mashal", txMatches[order.Hash], "unmarshal", decodeMatches[order.Hash])
	}
}

func TestProcessedHash(t *testing.T) {
	testDir := "TestProcessedHash"

	tomox := &TomoX{
		Orderbooks:  map[string]*OrderBook{},
		activePairs: make(map[string]bool),
		db: NewLDBEngine(&Config{
			DataDir:  testDir,
			DBEngine: "leveldb",
		}),
	}
	defer os.RemoveAll(testDir)

	if pHashes := tomox.getPendingHashes(); len(pHashes) != 0 {
		t.Error("Expected: no pending hash", "Actual:", len(pHashes))
	}

	var hash common.Hash
	hash = common.StringToHash("0x0000000000000000000000000000000000000000")
	tomox.addProcessedOrderHash(hash, 3)
	hash = common.StringToHash("0x0000000000000000000000000000000000000001")
	tomox.addProcessedOrderHash(hash, 3)
	hash = common.StringToHash("0x0000000000000000000000000000000000000002")
	tomox.addProcessedOrderHash(hash, 3)
	hash = common.StringToHash("0x0000000000000000000000000000000000000003")
	tomox.addProcessedOrderHash(hash, 3)
	hash = common.StringToHash("0x0000000000000000000000000000000000000004")
	tomox.addProcessedOrderHash(hash, 3)
	hash = common.StringToHash("0x0000000000000000000000000000000000000005")
	tomox.addProcessedOrderHash(hash, 3)
	pHashes := tomox.getProcessedOrderHash()
	if len(pHashes) != 3 {
		t.Error("Expected: 3 processed hash", "Actual:", len(pHashes), "pHashes", pHashes)
	}
	if !tomox.existProcessedOrderHash(common.StringToHash("0x0000000000000000000000000000000000000004")) {
		t.Error("Processed hash not exist")
	}
}

func TestTomoX_VerifyOrderNonce(t *testing.T) {
	testDir := "test_VerifyOrderNonce"

	tomox := &TomoX{
		orderCount: make(map[common.Address]*big.Int),
	}
	tomox.db = NewLDBEngine(&Config{
		DataDir:  testDir,
		DBEngine: "leveldb",
	})
	defer os.RemoveAll(testDir)

	// initial: orderCount is empty
	// verifyOrderNonce should PASS
	order := &OrderItem{
		Nonce:       big.NewInt(1),
		UserAddress: common.StringToAddress("0x00011"),
	}
	if err := tomox.verifyOrderNonce(order); err != nil {
		t.Error("Expected: no error")
	}

	storedOrderCountMap := make(map[common.Address]*big.Int)
	storedOrderCountMap[common.StringToAddress("0x00011")] = big.NewInt(5)
	if err := tomox.updateOrderCount(storedOrderCountMap); err != nil {
		t.Error("Failed to save orderCount", "err", err)
	}

	// set duplicated nonce
	order = &OrderItem{
		Nonce:       big.NewInt(5), //duplicated nonce
		UserAddress: common.StringToAddress("0x00011"),
	}
	if err := tomox.verifyOrderNonce(order); err != ErrOrderNonceTooLow {
		t.Error("Expected error: " + ErrOrderNonceTooLow.Error())
	}

	// set nonce too high
	order.Nonce = big.NewInt(110)
	if err := tomox.verifyOrderNonce(order); err != ErrOrderNonceTooHigh {
		t.Error("Expected error: " + ErrOrderNonceTooHigh.Error())
	}

	order.Nonce = big.NewInt(10)
	if err := tomox.verifyOrderNonce(order); err != nil {
		t.Error("Expected: no error")
	}

	// test new account
	order.UserAddress = common.StringToAddress("0x0022")
	if err := tomox.verifyOrderNonce(order); err != nil {
		t.Error("Expected: no error")
	}
}
