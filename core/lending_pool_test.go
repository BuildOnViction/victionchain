package core

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/ethclient"
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
	LendingID       uint64         `json:"lendingId,omitempty"`
	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash common.Hash `json:"hash" rlp:"-"`
}

func testSendLending(t *testing.T, nonce uint64, amount *big.Int, interest uint64, side string, status string, lendingID uint64) {

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
		CollateralToken: common.HexToAddress("0xC2fa1BA90b15E3612E0067A0020192938784D9C5"),
		LendingToken:    common.HexToAddress("0x45c25041b8e6CBD5c963E7943007187C3673C7c9"),
		Status:          status,
		Side:            side,
		Type:            "LO",
		Term:            86400,
		Interest:        interest,
	}

	tx := types.NewLendingTransaction(nonce, msg.Quantity, msg.Interest, msg.Term, msg.RelayerAddress, msg.UserAddress, msg.LendingToken, msg.CollateralToken, msg.Status, msg.Side, msg.Type, common.Hash{}, lendingID)
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
	testSendLending(t, 0, new(big.Int).SetUint64(1000000000000000000), 10, "INVEST", "NEW", 0)
	time.Sleep(2000)
	testSendLending(t, 1, new(big.Int).SetUint64(1000000000000000000), 10, "BORROW", "NEW", 0)
}
