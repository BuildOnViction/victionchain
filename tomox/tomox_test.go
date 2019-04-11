package tomox

import (
	"math/big"
	"testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"time"
	"encoding/json"
	"fmt"
	"encoding/hex"
)

func TestCreateOrder(t *testing.T) {
	price, ok := new(big.Int).SetString("250000000000000000000000000000000000000", 10)
	if !ok {
		t.Error("bad value", "price", "250000000000000000000000000000000000000")
	}
	v := []byte("123")
	order := &OrderItem{
		Quantity:        new(big.Int).SetUint64(1000000000000000000),
		Price:           price,
		ExchangeAddress: common.StringToAddress("0x0000000000000000000000000000000000000000"),
		UserAddress:     common.StringToAddress("0xf069080f7acb9a6705b4a51f84d9adc67b921bdf"),
		BaseToken:       common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		QuoteToken:      common.StringToAddress("0x9a8531c62d02af08cf237eb8aecae9dbcb69b6fd"),
		Status:          "New",
		Side:            "BUY",
		Type:            "LO",
		Hash:            common.StringToHash("0xdc842ea4a239d1a4e56f1e7ba31aab5a307cb643a9f5b89f972f2f5f0d1e7587"),
		Signature: &Signature{
			V: v[0],
			R: common.StringToHash("0xe386313e32a83eec20ecd52a5a0bd6bb34840416080303cecda556263a9270d0"),
			S: common.StringToHash("0x05cd5304c5ead37b6fac574062b150db57a306fa591c84fc4c006c4155ebda2a"),
		},
		FilledAmount: new(big.Int).SetUint64(0),
		Nonce:        new(big.Int).SetUint64(7761321472406488),
		MakeFee:      new(big.Int).SetUint64(4000000000000000),
		TakeFee:      new(big.Int).SetUint64(4000000000000000),
		CreatedAt:    uint64(time.Now().Unix()),
		UpdatedAt:    uint64(time.Now().Unix()),
	}

	topic := order.BaseToken.Hex() + "::" + order.QuoteToken.Hex()
	encodedTopic := fmt.Sprintf("0x%s", hex.EncodeToString([]byte(topic)))
	fmt.Println("topic: ", encodedTopic)

	ipaddress := "178.128.53.170"
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
