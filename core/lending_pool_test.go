package core

import (
	"context"
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/crypto/sha3"
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
	AutoTopUp       bool           `json:"autoTopUp,omitempty"`
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
	err = rpcClient.Call(&result, "tomox_getLendingOrderCount", userAddress)
	if err != nil {
		return 0, err
	}
	s := result.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	return uint64(n), nil
}


func (l *LendingMsg) computeHash() common.Hash {
	sha := sha3.NewKeccak256()
	if l.Status == lendingstate.LendingStatusCancelled {
		sha := sha3.NewKeccak256()
		sha.Write(l.Hash.Bytes())
		sha.Write(common.BigToHash(big.NewInt(int64(l.AccountNonce))).Bytes())
		sha.Write(l.UserAddress.Bytes())
		sha.Write(common.BigToHash(big.NewInt(int64(l.LendingId))).Bytes())
		sha.Write([]byte(l.Status))
		sha.Write(l.RelayerAddress.Bytes())
	} else {
		sha.Write(l.RelayerAddress.Bytes())
		sha.Write(l.UserAddress.Bytes())
		sha.Write(l.CollateralToken.Bytes())
		sha.Write(l.LendingToken.Bytes())
		sha.Write(common.BigToHash(l.Quantity).Bytes())
		sha.Write(common.BigToHash(big.NewInt(int64(l.Term))).Bytes())
		if l.Type == lendingstate.Limit {
			sha.Write(common.BigToHash(big.NewInt(int64(l.Interest))).Bytes())
		}
		sha.Write([]byte(l.Side))
		sha.Write([]byte(l.Status))
		sha.Write([]byte(l.Type))
		sha.Write(common.BigToHash(big.NewInt(int64(l.AccountNonce))).Bytes())
		sha.Write(common.BigToHash(big.NewInt(int64(l.LendingTradeId))).Bytes())
		if l.Status == types.LendingSideBorrow {
			autoTopUp := int64(0)
			if l.AutoTopUp {
				autoTopUp = int64(1)
			}
			sha.Write(common.BigToHash(big.NewInt(autoTopUp)).Bytes())
		}
	}
	return common.BytesToHash(sha.Sum(nil))

}
func testSendLending(t *testing.T, amount *big.Int, interest uint64, side string, status string, autoTopUp bool, lendingId, tradeId uint64, cancelledHash common.Hash, extraData string) {

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
		Term:            60,
		AutoTopUp:       autoTopUp,
		Interest:        interest,
		LendingId:       lendingId,
		LendingTradeId:  tradeId,
		ExtraData:       extraData,
	}
	if cancelledHash != (common.Hash{}) {
		msg.Hash = cancelledHash
	} else {
		msg.Hash = msg.computeHash()
	}

	tx := types.NewLendingTransaction(nonce, msg.Quantity, msg.Interest, msg.Term, msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.AutoTopUp, msg.Status, msg.Side, msg.Type, msg.Hash, lendingId, tradeId, msg.ExtraData)
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
	// 10%
	interestRate := 10 * common.BaseLendingInterest.Uint64()
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Investing, lendingstate.LendingStatusNew, true,0, 0, common.Hash{}, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Investing, lendingstate.LendingStatusNew, true,0, 0,  common.Hash{},"")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Investing, lendingstate.LendingStatusNew, true,0, 0,  common.Hash{},"")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Investing, lendingstate.LendingStatusNew, true,0, 0,  common.Hash{},"")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Borrowing, lendingstate.LendingStatusNew, true,0, 1, common.Hash{}, "")
	time.Sleep(2 * time.Second)
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Borrowing, lendingstate.LendingStatusNew, true,0, 0, common.Hash{}, "")
}

func TestCancelLending(t *testing.T) {
	// 10%
	interestRate := 10 * common.BaseLendingInterest.Uint64()
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Borrowing, lendingstate.LendingStatusNew, true,0, 0, common.Hash{}, "")
	time.Sleep(2 * time.Second)
	//TODO: run the above testcase first, then updating lendingId, Hash
	testSendLending(t, new(big.Int).Mul(_1E8, big.NewInt(1000)), interestRate, lendingstate.Investing, lendingstate.LendingStatusCancelled, true,1, 0, common.HexToHash("0x3da4e24b9c0f60e04cdb4c4494de37203c6e1a354907cbd6d9bbbe2e52aecaab"), "")
}

