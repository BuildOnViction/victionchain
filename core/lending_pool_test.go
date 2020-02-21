package core

import (
	"context"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/crypto/sha3"
	"github.com/tomochain/tomochain/ethclient"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"
	"log"
	"math/big"
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
	}
	return common.BytesToHash(sha.Sum(nil))

}
func testSendLending(t *testing.T, nonce uint64, amount *big.Int, interest uint64, side string, status string, lendingId, tradeId uint64, cancelledHash common.Hash, extraData string) {

	client, err := ethclient.Dial("http://127.0.0.1:1545")
	if err != nil {
		log.Print(err)
	}
	privateKey, err := crypto.HexToECDSA("65ec4d4dfbcac594a14c36baa462d6f73cd86134840f6cf7b80a1e1cd33473e2")
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
		Status:          status,
		Side:            side,
		Type:            "LO",
		Term:            60,
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

	tx := types.NewLendingTransaction(nonce, msg.Quantity, msg.Interest, msg.Term, msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.Status, msg.Side, msg.Type, msg.Hash, lendingId, tradeId, msg.ExtraData)
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
	//test matching FULL FILLED
	testSendLending(t, 0, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, common.Hash{}, "")
	time.Sleep(2000)
	testSendLending(t, 1, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusNew, 0, 0, common.Hash{}, "")
	time.Sleep(2000)
	//test pay the above loan
	testSendLending(t, 2, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Borrowing, lendingstate.Payment, 0, 1, common.Hash{}, "")
	time.Sleep(2000)

	// test matching PARTIAL FILLED
	testSendLending(t, 3, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, common.Hash{}, "")
	time.Sleep(2000)
	testSendLending(t, 4, new(big.Int).SetUint64(500000000000000000), 10, lendingstate.Borrowing, lendingstate.LendingStatusNew, 0, 0, common.Hash{}, "")
	time.Sleep(2000)

	//// test cancel
	testSendLending(t, 5, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusNew, 0, 0, common.Hash{}, "")
	time.Sleep(2000)
	////TODO: update cancelled hash
	testSendLending(t, 6, new(big.Int).SetUint64(1000000000000000000), 10, lendingstate.Investing, lendingstate.LendingStatusCancelled, 3, 0, common.HexToHash("0xdf2efbe1970dddb9ced42d6f5cf4a3618522bc57704515c5b47027d85b43fec1"), "")

}

