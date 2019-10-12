package core

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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
	err = rpcClient.Call(&result, "tomox_getOrderNonce", userAddress.Hex())
	if err != nil {
		return 0, err
	}
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	return uint64(n), nil
}
func testSendOrder(t *testing.T, amount, price *big.Int, side string) {

	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("3b43d337ae657c351d2542c7ee837c39f5db83da7ffffb611992ebc2f676743b")
	if err != nil {
		log.Fatal(err)
	}

	msg := &OrderMsg{
		Quantity:        amount,
		Price:           price,
		ExchangeAddress: common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e"),
		UserAddress:     common.HexToAddress("0xF7349C253FF7747Df661296E0859c44e974fb52E"),
		BaseToken:       common.HexToAddress("0x4d7eA2cE949216D6b120f3AA10164173615A2b6C"),
		QuoteToken:      common.HexToAddress("0x0000000000000000000000000000000000000001"),
		Status:          "NEW",
		Side:            side,
		Type:            "LO",
		PairName:        "BTC/TOMO",
	}
	nonce, _ := getNonce(t, msg.UserAddress)
	tx := types.NewOrderTransaction(nonce, msg.Quantity, msg.Price, msg.ExchangeAddress, msg.UserAddress, msg.BaseToken, msg.QuoteToken, msg.Status, msg.Side, msg.Type, msg.PairName, common.Hash{}, 0)
	signedTx, err := types.OrderSignTx(tx, types.OrderTxSigner{}, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendOrderTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSendBuyOrder(t *testing.T) {
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "BUY")
}

func TestSendSellOrder(t *testing.T) {
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "SELL")
}
func TestFilled(t *testing.T) {
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "BUY")
	time.Sleep(5 * time.Second)
	testSendOrder(t, new(big.Int).SetUint64(1000000000000000000), new(big.Int).SetUint64(100000000000000000), "SELL")
}
func TestPartialFilled(t *testing.T) {

}
func TestNoMatch(t *testing.T) {

}
