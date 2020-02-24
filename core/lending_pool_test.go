package core

import (
	"context"
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/ethclient"
	"github.com/tomochain/tomochain/rpc"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"log"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

type LendingMsg struct {
	AccountNonce    uint64         `json:"nonce"    gencodec:"required"`
	Quantity        *big.Int       `json:"quantity,omitempty"`
	RelayerAddress  common.Address `json:"relayerAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	CollateralToken common.Address `json:"collateralToken,omitempty"`
	LendingToken    common.Address `json:"lendingToken,omitempty"`
	Term            uint64         `json:"term,omitempty"`
	Interest        uint64         `json:"interest,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	LendingId       uint64         `json:"lendingId,omitempty"`
	LendingTradeId  uint64         `json:"tradeId,omitempty"`
	ExtraData       string         `json:"extraData,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

func getLendingNonce(t *testing.T, userAddress common.Address) (uint64, error) {
	rpcClient, err := rpc.DialHTTP("http://127.0.0.1:8501")
	defer rpcClient.Close()
	if err != nil {
		return 0, err
	}
	var result interface{}
	if err != nil {

		return 0, err
	}
	err = rpcClient.Call(&result, "tomox_getLendingOrderCount", userAddress)
	if err != nil {
		return 0, err
	}
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	return uint64(n), nil
}

func testSendLending(t *testing.T, amount *big.Int, interest uint64, side string, status string, lendingId, tradeId uint64, extraData string) {

	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		log.Print(err)
	}
	privateKey, err := crypto.HexToECDSA("3b43d337ae657c351d2542c7ee837c39f5db83da7ffffb611992ebc2f676743b")
	if err != nil {
		log.Print(err)
	}
	nonce, err := getLendingNonce(t, crypto.PubkeyToAddress(privateKey.PublicKey))
	fmt.Println("nonce", nonce, "err", err)
	msg := &LendingMsg{
		AccountNonce:    nonce,
		Quantity:        amount,
		RelayerAddress:  common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e"),
		UserAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		CollateralToken: BTCAddress,
		LendingToken:    USDAddress,
		Status:          status,
		Side:            side,
		Type:            "LO",
		Term:            86400,
		Interest:        interest,
		LendingId:       lendingId,
		LendingTradeId:  tradeId,
		ExtraData:       extraData,
	}
	tx := types.NewLendingTransaction(nonce, msg.Quantity, msg.Interest, msg.Term, msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.Status, msg.Side, msg.Type, common.Hash{}, lendingId, tradeId, msg.ExtraData)
	signedTx, err := types.LendingSignTx(tx, types.LendingTxSigner{}, privateKey)
	if err != nil {
		log.Print(err)
	}

	err = client.SendLendingTransaction(context.Background(), signedTx)
	if err != nil {
		log.Print(err)
	}
}

func testSendCancelLendingOrder(t *testing.T, nonce uint64, amount *big.Int, interest uint64, side string, lendingId, tradeId uint64, extraData string) {
	client, err := ethclient.Dial("http://127.0.0.1:8501")
	if err != nil {
		log.Print(err)
	}
	privateKey, err := crypto.HexToECDSA("3b43d337ae657c351d2542c7ee837c39f5db83da7ffffb611992ebc2f676743b")
	if err != nil {
		log.Print(err)
	}
	msg := &LendingMsg{
		AccountNonce:    nonce,
		Quantity:        amount,
		RelayerAddress:  common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e"),
		UserAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		CollateralToken: BTCAddress,
		LendingToken:    USDAddress,
		Status:          lendingstate.LendingStatusCancelled,
		Side:            side,
		Type:            "LO",
		Term:            86400,
		Interest:        interest,
		LendingId:       lendingId,
		LendingTradeId:  tradeId,
		ExtraData:       extraData,
	}

	tx := types.NewLendingTransaction(nonce, msg.Quantity, msg.Interest, msg.Term, msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.Status, msg.Side, msg.Type, common.Hash{}, lendingId, msg.LendingTradeId, msg.ExtraData)
	signedTx, err := types.LendingSignTx(tx, types.LendingTxSigner{}, privateKey)
	if err != nil {
		log.Print(err)
	}

	err = client.SendLendingTransaction(context.Background(), signedTx)
	if err != nil {
		log.Print(err)
	}
}

func TestSendLending(t *testing.T) {
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusNew, 0, 0, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusNew, 0, 0, "")
	//time.Sleep(2  * time.Second)
	//testSendLending(t, new(big.Int).Mul(new(big.Int).SetUint64(1000000000000000000), big.NewInt(1005)), 10, lendingstate.Borrowing, lendingstate.Payment, 0, 1, common.Uint64ToHash(1).Hex())
}

func TestCancelLending(t *testing.T) {
	//testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, "")
	//time.Sleep(2  * time.Second)
	testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusCancelled, 1, 0, "")
	time.Sleep(2  * time.Second)
	//testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusNew, 0, 0, "")
	//time.Sleep(2 * time.Second)
	//testSendLending(t, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusCancelled, 2, 0, "")
	//time.Sleep(2  * time.Second)
}
