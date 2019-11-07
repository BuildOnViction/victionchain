package core

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/ethclient"
	"github.com/tomochain/tomochain/rpc"
)

type OrderMsg struct {
	AccountNonce    uint64         `json:"nonce"    gencodec:"required"`
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	OrderID         uint64         `json:"orderid,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

func getNonce(t *testing.T, userAddress common.Address) (uint64, error) {
	rpcClient, err := rpc.DialHTTP("http://127.0.0.1:8501")
	defer rpcClient.Close()
	if err != nil {
		return 0, err
	}
	var result interface{}
	if err != nil {

		return 0, err
	}
	err = rpcClient.Call(&result, "tomox_getOrderCount", userAddress)
	if err != nil {
		return 0, err
	}
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	return uint64(n), nil
}
func testSendOrder(t *testing.T, amount, price *big.Int, side string, status string, orderID uint64) {

	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		log.Print(err)
	}

	privateKey, err := crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
	if err != nil {
		log.Print(err)
	}
	msg := &OrderMsg{
		Quantity:        amount,
		Price:           price,
		ExchangeAddress: common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e"),
		UserAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		BaseToken:       common.HexToAddress("0x4d7eA2cE949216D6b120f3AA10164173615A2b6C"),
		QuoteToken:      common.HexToAddress("0x0000000000000000000000000000000000000001"),
		Status:          status,
		Side:            side,
		Type:            "LO",
		PairName:        "BTC/TOMO",
	}
	nonce, _ := getNonce(t, msg.UserAddress)
	tx := types.NewOrderTransaction(nonce, msg.Quantity, msg.Price, msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.PairName, common.Hash{}, orderID)
	signedTx, err := types.OrderSignTx(tx, types.OrderTxSigner{}, privateKey)
	if err != nil {
		log.Print(err)
	}

	err = client.SendOrderTransaction(context.Background(), signedTx)
	if err != nil {
		log.Print(err)
	}
}

func testSendOrderETHBTC(t *testing.T, amount, price *big.Int, side string, status string, orderID uint64) {

	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		log.Print(err)
	}

	privateKey, err := crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
	if err != nil {
		log.Print(err)
	}
	msg := &OrderMsg{
		Quantity:        amount,
		Price:           price,
		ExchangeAddress: common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e"),
		UserAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		BaseToken:       common.HexToAddress("0xC2fa1BA90b15E3612E0067A0020192938784D9C5"),
		QuoteToken:      common.HexToAddress("0x4d7eA2cE949216D6b120f3AA10164173615A2b6C"),
		Status:          status,
		Side:            side,
		Type:            "LO",
		PairName:        "ETH/BTC",
	}
	nonce, _ := getNonce(t, msg.UserAddress)
	tx := types.NewOrderTransaction(nonce, msg.Quantity, msg.Price, msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.PairName, common.Hash{}, orderID)
	signedTx, err := types.OrderSignTx(tx, types.OrderTxSigner{}, privateKey)
	if err != nil {
		log.Print(err)
	}

	err = client.SendOrderTransaction(context.Background(), signedTx)
	if err != nil {
		log.Print(err)
	}
}

func TestSendBuyOrder(t *testing.T) {
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "BUY", "NEW", 0)
}

func TestSendSellOrder(t *testing.T) {
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "SELL", "NEW", 0)
}
func TestFilled(t *testing.T) {
	//BTC/TOMO
	price := new(big.Int).Mul(big.NewInt(1000000000000000000), big.NewInt(5000))
	testSendOrder(t, new(big.Int).SetUint64(10000000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(10000000000000), price, "SELL", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(100000000000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(100000000000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(200000000000000000), price, "SELL", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(10000000000000), price, "SELL", "NEW", 0)
	//ETH/BTC
	price = new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(20))
	testSendOrderETHBTC(t, new(big.Int).SetUint64(1000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrderETHBTC(t, new(big.Int).SetUint64(1000000000), price, "SELL", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrderETHBTC(t, new(big.Int).SetUint64(100000000000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrderETHBTC(t, new(big.Int).SetUint64(100000000000000000), price, "BUY", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrderETHBTC(t, new(big.Int).SetUint64(200000000000000000), price, "SELL", "NEW", 0)
	time.Sleep(5 * time.Second)
	testSendOrderETHBTC(t, new(big.Int).SetUint64(100000000), price, "SELL", "NEW", 0)
}
func TestPartialFilled(t *testing.T) {

}
func TestNoMatch(t *testing.T) {

}

func TestCancelOrder(t *testing.T) {
	//testSendOrder(t, new(big.Int).SetUint64(48), new(big.Int).SetUint64(15), "BUY", "NEW", 0)
	//time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(48), new(big.Int).SetUint64(15), "BUY", "CANCELLED", 3)
	//time.Sleep(5 * time.Second)
	//testSendOrder(t, new(big.Int).SetUint64(48), new(big.Int).SetUint64(15), "SELL", "NEW", 0)
}
