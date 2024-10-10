package backends

import (
	"context"
	"errors"
	ethereum "github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/params"
	"math/big"
	"testing"
)

func TestSimulatedBackend_EstimateGasWithPrice(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	sim := NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(params.Ether*2 + 2e17)}})
	defer sim.Close()

	receipant := common.HexToAddress("deadbeef")
	var cases = []struct {
		name        string
		message     ethereum.CallMsg
		expect      uint64
		expectError error
	}{
		{"EstimateWithoutPrice", ethereum.CallMsg{
			From:     addr,
			To:       &receipant,
			Gas:      0,
			GasPrice: big.NewInt(0),
			Value:    big.NewInt(1000),
			Data:     nil,
		}, 21000, nil},

		{"EstimateWithPrice", ethereum.CallMsg{
			From:     addr,
			To:       &receipant,
			Gas:      0,
			GasPrice: big.NewInt(1000),
			Value:    big.NewInt(1000),
			Data:     nil,
		}, 21000, nil},

		{"EstimateWithVeryHighPrice", ethereum.CallMsg{
			From:     addr,
			To:       &receipant,
			Gas:      0,
			GasPrice: big.NewInt(1e14), // gascost = 2.1ether
			Value:    big.NewInt(1e17), // the remaining balance for fee is 2.1ether
			Data:     nil,
		}, 21000, nil},

		{"EstimateWithSuperhighPrice", ethereum.CallMsg{
			From:     addr,
			To:       &receipant,
			Gas:      0,
			GasPrice: big.NewInt(2e14), // gascost = 4.2ether
			Value:    big.NewInt(1000),
			Data:     nil,
		}, 21000, errors.New("gas required exceeds allowance (10999)")}, // 10999=(2.2ether-1000wei)/(2e14)
	}
	for _, c := range cases {
		got, err := sim.EstimateGas(context.Background(), c.message)
		if c.expectError != nil {
			if err == nil {
				t.Fatalf("Expect error, got nil")
			}
			if c.expectError.Error() != err.Error() {
				t.Fatalf("Expect error, want %v, got %v", c.expectError, err)
			}
			continue
		}
		if got != c.expect {
			t.Fatalf("Gas estimation mismatch, want %d, got %d", c.expect, got)
		}
	}
}
