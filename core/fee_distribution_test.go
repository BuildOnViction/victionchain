// Copyright 2024 The VictionChain Authors
// This file is part of the VictionChain library.
//
// The VictionChain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The VictionChain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the VictionChain library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/rawdb"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/params"
)

var (
	// VRC25GasPrice is the gas price for VRC25 token fee payment (0.25 Gwei)
	VRC25GasPrice = big.NewInt(250000000)
)

// TestSingleTransactionFeeDistribution tests fee distribution for individual transactions
// including standard ETH transactions and miner fee crediting
func TestNormalTransactionFeeDistribution(t *testing.T) {
	tests := []struct {
		name                         string
		gasPrice                     *big.Int
		gasLimit                     uint64
		gasUsed                      uint64
		transferValue                *big.Int
		minerBalance                 *big.Int
		senderBalance                *big.Int
		expectSuccess                bool
		experientialHardFork         bool
		expectedMinerBalanceAfterTx  *big.Int
		expectedSenderBalanceAfterTx *big.Int
	}{
		{
			name:                         "Before Atlas - Standard ETH transaction fee distribution",
			gasPrice:                     big.NewInt(1000000000), // 1 Gwei
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000000000000000000), // 1 ETH
			minerBalance:                 big.NewInt(0),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			expectSuccess:                true,
			experientialHardFork:         false,
			expectedMinerBalanceAfterTx:  big.NewInt(21000000000000),     // 0 + (21000 * 1 Gwei)
			expectedSenderBalanceAfterTx: big.NewInt(999979000000000000), // 2 ETH - 1 ETH - 0.000021 ETH
		},
		{
			name:                         "Before Atlas - High gas price fee distribution",
			gasPrice:                     big.NewInt(100000000000), // 100 Gwei
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(100000000000000000), // 0.1 ETH
			minerBalance:                 big.NewInt(0),
			senderBalance:                big.NewInt(5000000000000000000), // 5 ETH
			expectSuccess:                true,
			experientialHardFork:         false,
			expectedMinerBalanceAfterTx:  big.NewInt(2100000000000000),    // 0 + (21000 * 100 Gwei)
			expectedSenderBalanceAfterTx: big.NewInt(4897900000000000000), // 5 ETH - 0.1 ETH - 0.0021 ETH
		},
		{
			name:                         "After Experimental - Standard ETH transaction fee distribution",
			gasPrice:                     big.NewInt(2000000000), // 2 Gwei
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000000000000000000), // 1 ETH
			minerBalance:                 big.NewInt(0),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			expectSuccess:                true,
			experientialHardFork:         true,
			expectedMinerBalanceAfterTx:  big.NewInt(42000000000000),     // 0 + (21000 * 2000000000)
			expectedSenderBalanceAfterTx: big.NewInt(999958000000000000), // 2 ETH - 1 ETH - 0.000042 ETH
		},
		{
			name:                         "After Experimental - Gas refund scenario - partial gas usage",
			gasPrice:                     big.NewInt(1000000000), // 1 Gwei
			gasLimit:                     50000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000000000000000000), // 1 ETH
			minerBalance:                 big.NewInt(0),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			expectSuccess:                true,
			experientialHardFork:         true,
			expectedMinerBalanceAfterTx:  big.NewInt(21000000000000),     // 0 + (21000 * 1 Gwei)
			expectedSenderBalanceAfterTx: big.NewInt(999979000000000000), // 2 ETH - 1 ETH - 0.000021 ETH
		},

		{
			name:                         "Before Atlas - Gas refund scenario - partial gas usage",
			gasPrice:                     big.NewInt(1000000000),          // 1 Gwei
			gasLimit:                     50000,                           // Higher limit
			gasUsed:                      21000,                           // Only basic transfer gas used
			transferValue:                big.NewInt(1000000000000000000), // 1 ETH
			minerBalance:                 big.NewInt(0),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			expectSuccess:                true,
			experientialHardFork:         false,
			expectedMinerBalanceAfterTx:  big.NewInt(21000000000000),     // 0 + (21000 * 1 Gwei) - only actual gas used
			expectedSenderBalanceAfterTx: big.NewInt(999979000000000000), // 2 ETH - 1 ETH - 0.000021 ETH (refund 29000 gas)
		},
		{
			name:                         "After Experimental - Gas refund with high gas limit",
			gasPrice:                     VRC25GasPrice,                   // VRC25 gas price = 0.25 Gwei
			gasLimit:                     100000,                          // Much higher limit
			gasUsed:                      21000,                           // Only basic transfer gas used
			transferValue:                big.NewInt(1000000000000000000), // 1 ETH
			minerBalance:                 big.NewInt(0),                   // 1 ETH initial
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			expectSuccess:                true,
			experientialHardFork:         true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),      //  + (21000 * 250000000)
			expectedSenderBalanceAfterTx: big.NewInt(999994750000000000), // 2 ETH - 1 ETH - 0.00000525 ETH (refund 79000 gas = 19750000000000 wei )
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			key, _ := crypto.GenerateKey()
			sender := crypto.PubkeyToAddress(key.PublicKey)
			recipient := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			miner := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

			db := rawdb.NewMemoryDatabase()
			statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))

			// Set initial balances
			statedb.AddBalance(sender, tt.senderBalance)
			statedb.AddBalance(miner, tt.minerBalance)

			// Create chain config
			config := &params.ChainConfig{
				ChainId:           big.NewInt(89),
				HomesteadBlock:    big.NewInt(0),
				EIP150Block:       big.NewInt(0),
				EIP155Block:       big.NewInt(0),
				EIP158Block:       big.NewInt(0),
				ByzantiumBlock:    big.NewInt(0),
				AtlasBlock:        big.NewInt(13523400),
				ExperientialBlock: big.NewInt(13523410), // Fix refund issue at Atlas
			}

			var blockNumber *big.Int
			if tt.experientialHardFork {
				blockNumber = big.NewInt(13523411) // After Experiential
			} else {
				blockNumber = big.NewInt(13523399) // Before Experiential
			}

			// Create EVM context
			context := vm.Context{
				CanTransfer: CanTransfer,
				Transfer:    Transfer,
				GetHash:     func(uint64) common.Hash { return common.Hash{} },
				Origin:      sender,
				Coinbase:    miner,
				BlockNumber: blockNumber,
				Time:        big.NewInt(0),
				Difficulty:  big.NewInt(0),
				GasLimit:    8000000,
				GasPrice:    tt.gasPrice,
			}

			// Record initial balances
			initialSenderBalance := statedb.GetBalance(sender)
			initialMinerBalance := statedb.GetBalance(miner)
			initialRecipientBalance := statedb.GetBalance(recipient)

			evm := vm.NewEVM(context, statedb, nil, config, vm.Config{})
			// Create message
			msg := types.NewMessage(sender, &recipient, 0, tt.transferValue, tt.gasLimit, tt.gasPrice, nil, false, nil)

			// Create state transition and execute
			gp := new(GasPool).AddGas(tt.gasLimit)
			st := NewStateTransition(evm, msg, gp)

			// Execute transaction
			_, gasUsed, failed, err := st.TransitionDb(miner)

			if tt.expectSuccess {
				if err != nil {
					t.Fatalf("Transaction failed unexpectedly: %v", err)
				}
				if failed {
					t.Fatal("Transaction marked as failed but should have succeeded")
				}

				// Verify gas usage
				if gasUsed != tt.gasUsed {
					t.Errorf("Gas usage mismatch: got %d, want %d", gasUsed, tt.gasUsed)
				}

				// Calculate expected fee
				expectedFee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), tt.gasPrice)
				fmt.Println("Expected Fee:", expectedFee, "Gas Used:", gasUsed, "Gas Price:", tt.gasPrice)

				// Verify final balances
				finalSenderBalance := statedb.GetBalance(sender)
				finalMinerBalance := statedb.GetBalance(miner)
				finalRecipientBalance := statedb.GetBalance(recipient)

				// Check sender balance against expected value
				if finalSenderBalance.Cmp(tt.expectedSenderBalanceAfterTx) != 0 {
					t.Errorf("Sender balance incorrect: got %v, want %v", finalSenderBalance, tt.expectedSenderBalanceAfterTx)
					t.Errorf("  Initial: %v, Transfer: %v, Fee: %v", initialSenderBalance, tt.transferValue, expectedFee)
				}

				// Check miner balance against expected value
				if finalMinerBalance.Cmp(tt.expectedMinerBalanceAfterTx) != 0 {
					t.Errorf("Miner balance incorrect: got %v, want %v", finalMinerBalance, tt.expectedMinerBalanceAfterTx)
					t.Errorf("  Initial: %v, Fee received: %v", initialMinerBalance, expectedFee)
				}

				// Verify the calculation logic matches expected values
				calculatedSenderBalance := new(big.Int).Sub(initialSenderBalance, tt.transferValue)
				calculatedSenderBalance.Sub(calculatedSenderBalance, expectedFee)
				if calculatedSenderBalance.Cmp(tt.expectedSenderBalanceAfterTx) != 0 {
					t.Errorf("Expected sender balance calculation mismatch: calculated %v, expected %v", calculatedSenderBalance, tt.expectedSenderBalanceAfterTx)
				}

				calculatedMinerBalance := new(big.Int).Add(initialMinerBalance, expectedFee)
				if calculatedMinerBalance.Cmp(tt.expectedMinerBalanceAfterTx) != 0 {
					t.Errorf("Expected miner balance calculation mismatch: calculated %v, expected %v", calculatedMinerBalance, tt.expectedMinerBalanceAfterTx)
				}

				// For gas refund scenarios, verify that only actual gas used is charged
				if tt.gasUsed < tt.gasLimit {
					refundedGas := tt.gasLimit - tt.gasUsed
					refundAmount := new(big.Int).Mul(new(big.Int).SetUint64(refundedGas), tt.gasPrice)
					t.Logf("Gas refund: %d gas units = %v wei refunded", refundedGas, refundAmount)
				}

				// Check recipient balance (should receive the transfer value)
				expectedRecipientBalance := new(big.Int).Add(initialRecipientBalance, tt.transferValue)
				if finalRecipientBalance.Cmp(expectedRecipientBalance) != 0 {
					t.Errorf("Recipient balance incorrect: got %v, want %v", finalRecipientBalance, expectedRecipientBalance)
				}

				// Verify fee conservation
				totalInitial := new(big.Int).Add(initialSenderBalance, initialMinerBalance)
				totalInitial.Add(totalInitial, initialRecipientBalance)
				totalFinal := new(big.Int).Add(finalSenderBalance, finalMinerBalance)
				totalFinal.Add(totalFinal, finalRecipientBalance)

				if totalInitial.Cmp(totalFinal) != 0 {
					t.Errorf("Total balance conservation failed: initial %v, final %v", totalInitial, totalFinal)
				}
			} else {
				if err == nil && !failed {
					t.Error("Expected transaction to fail but it succeeded")
				}
			}
		})

		fmt.Println("___________________________________________________________________________________")
	}
}

// TestVRC25SingleTokenFeeDistribution tests VRC25 token fee distribution mechanisms
// including fee capacity deduction and miner payment
func TestVRC25TokenFeeDistribution(t *testing.T) {
	tests := []struct {
		name                         string
		tokenCapacity                *big.Int
		gasLimit                     uint64
		gasUsed                      uint64
		transferValue                *big.Int
		senderBalance                *big.Int
		minerBalance                 *big.Int
		expectTokenFee               bool
		expectSuccess                bool
		expectedMinerBalanceAfterTx  *big.Int
		expectedSenderBalanceAfterTx *big.Int
		expectedTokenCapacityAfterTx *big.Int
	}{
		{
			name:                         "VRC25 sufficient capacity - Post Atlas",
			tokenCapacity:                big.NewInt(10000000000000000), // 0.01 ETH worth (sufficient for 21000 * 250000000)
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               true,
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999999999999999000), // 2 ETH - 1000 wei (no ETH fee, token fee used)
			expectedTokenCapacityAfterTx: big.NewInt(9994750000000000),    // 10000000000000000 - (21000 * 250000000) = 9994750000000000
		},
		{
			name:                         "VRC25 exact capacity match - Post Atlas",
			tokenCapacity:                big.NewInt(5250000000000), // Exactly 21000 * 250000000 = 5250000000000
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               false, // Falls back to ETH when capacity <= required fee
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999994749999999000), // 2 ETH - 1000 wei - 5250000000000 (ETH fee) = 1999994749999999000
			expectedTokenCapacityAfterTx: big.NewInt(5250000000000),       // No change (ETH fallback)
		},
		{
			name:                         "VRC25 high gas limit with sufficient capacity",
			tokenCapacity:                big.NewInt(20000000000000000), // 0.02 ETH worth (sufficient for gas)
			gasLimit:                     100000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               true,
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999999999999999000), // 2 ETH - 1000 wei (no ETH fee)
			expectedTokenCapacityAfterTx: big.NewInt(19994750000000000),   // 20000000000000000 - (21000 * 250000000) = 19994750000000000
		},
		{
			name:                         "VRC25 zero capacity - fallback to ETH",
			tokenCapacity:                big.NewInt(0),
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               false,
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999994749999999000), // 2 ETH - 1000 wei - (21000 * 250000000) = 1999994749999999000
			expectedTokenCapacityAfterTx: big.NewInt(0),                   // Remains zero
		},
		{
			name:                         "VRC25 large gas refund with token fee",
			tokenCapacity:                big.NewInt(25000000000000000), // 0.025 ETH worth (sufficient for gas)
			gasLimit:                     100000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               true,
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999999999999999000), // 2 ETH - 1000 wei (no ETH fee, token fee used)
			expectedTokenCapacityAfterTx: big.NewInt(24994750000000000),   // 25000000000000000 - (21000 * 250000000) = 24994750000000000
		},
		{
			name:                         "VRC25 small gas refund with exact capacity",
			tokenCapacity:                big.NewInt(7500000000000), // Less than required for gas limit (30000 * 250000000 = 7500000000000)
			gasLimit:                     30000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               false, // Falls back to ETH when capacity <= required fee for gas limit
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999994749999999000), // 2 ETH - 1000 wei - 5250000000000 (ETH fee for gas limit) = 1999994749999999000
			expectedTokenCapacityAfterTx: big.NewInt(7500000000000),       // No change (ETH fallback)
		},
		{
			name:                         "VRC25 maximum gas refund scenario",
			tokenCapacity:                big.NewInt(50000000000000000), // 0.05 ETH worth (sufficient for gas)
			gasLimit:                     200000,                        // Very high limit
			gasUsed:                      21000,                         // Minimal usage - maximum refund
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(2000000000000000000), // 2 ETH
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               true,
			expectSuccess:                true,
			expectedMinerBalanceAfterTx:  big.NewInt(5250000000000),       // 0 + (21000 * 250000000) = 5250000000000
			expectedSenderBalanceAfterTx: big.NewInt(1999999999999999000), // 2 ETH - 1000 wei (no ETH fee, token fee used)
			expectedTokenCapacityAfterTx: big.NewInt(49994750000000000),   // 50000000000000000 - (21000 * 250000000) = 49994750000000000
		},
		{
			name:                         "VRC25 insufficient capacity and balance - fails",
			tokenCapacity:                big.NewInt(1000000000000), // Very low capacity (insufficient for 21000 * 250000000 = 5250000000000)
			gasLimit:                     21000,
			gasUsed:                      21000,
			transferValue:                big.NewInt(1000),
			senderBalance:                big.NewInt(3000000000000), // Very low balance (insufficient for ETH fallback fee + transfer)
			minerBalance:                 big.NewInt(0),
			expectTokenFee:               false,                     // Falls back to ETH but insufficient balance
			expectSuccess:                false,                     // Should fail due to insufficient funds
			expectedMinerBalanceAfterTx:  big.NewInt(0),             // No change (transaction fails)
			expectedSenderBalanceAfterTx: big.NewInt(3000000000000), // No change (transaction fails)
			expectedTokenCapacityAfterTx: big.NewInt(1000000000000), // No change (transaction fails)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			key, _ := crypto.GenerateKey()
			sender := crypto.PubkeyToAddress(key.PublicKey)
			tokenAddr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
			miner := common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

			db := rawdb.NewMemoryDatabase()
			statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))

			// Set initial balances
			statedb.AddBalance(sender, tt.senderBalance)
			statedb.AddBalance(miner, tt.minerBalance)

			// Set up VRC25 token fee capacity
			if tt.tokenCapacity.Sign() > 0 {
				slotTokensState := state.SlotTRC21Issuer["tokensState"]
				balanceKey := state.GetLocMappingAtKey(tokenAddr.Hash(), slotTokensState)
				statedb.SetState(common.TRC21IssuerSMC, common.BigToHash(balanceKey), common.BigToHash(tt.tokenCapacity))
				statedb.AddBalance(common.TRC21IssuerSMC, tt.tokenCapacity)
			}

			// Create chain config
			config := &params.ChainConfig{
				ChainId:           big.NewInt(89),
				HomesteadBlock:    big.NewInt(0),
				EIP150Block:       big.NewInt(0),
				EIP155Block:       big.NewInt(0),
				EIP158Block:       big.NewInt(0),
				ByzantiumBlock:    big.NewInt(0),
				AtlasBlock:        big.NewInt(13523400), // VRC25 active after this block
				ExperientialBlock: big.NewInt(13523410), // Fix refund issue at Atlas
			}

			// Create EVM context
			context := vm.Context{
				CanTransfer: CanTransfer,
				Transfer:    Transfer,
				GetHash:     func(uint64) common.Hash { return common.Hash{} },
				Origin:      sender,
				Coinbase:    miner,
				BlockNumber: big.NewInt(13523411), // After ExperientialBlock
				Time:        big.NewInt(0),
				Difficulty:  big.NewInt(0),
				GasLimit:    8000000,
				GasPrice:    VRC25GasPrice,
			}

			// Record initial state
			initialSenderBalance := statedb.GetBalance(sender)
			initialMinerBalance := statedb.GetBalance(miner)
			initialTokenCapacity := state.GetTRC21FeeCapacityFromStateWithToken(statedb, &tokenAddr)

			evm := vm.NewEVM(context, statedb, nil, config, vm.Config{})
			// Create VRC25 message
			msg := types.NewMessage(sender, &tokenAddr, 0, tt.transferValue, tt.gasLimit, VRC25GasPrice, nil, true, tt.tokenCapacity)

			// Create state transition and execute
			gp := new(GasPool).AddGas(tt.gasLimit)
			st := NewStateTransition(evm, msg, gp)

			_, gasUsed, failed, err := st.TransitionDb(miner)

			if tt.expectSuccess {
				if err != nil {
					t.Fatalf("VRC25 transaction failed unexpectedly: %v", err)
				}
				if failed {
					t.Fatal("VRC25 transaction marked as failed but should have succeeded")
				}

				// Verify gas usage
				if gasUsed != tt.gasUsed {
					t.Errorf("Gas usage mismatch: got %d, want %d", gasUsed, tt.gasUsed)
				}

				// Calculate actual fee distribution from balance changes
				finalSenderBalance := statedb.GetBalance(sender)
				finalMinerBalance := statedb.GetBalance(miner)
				finalTokenCapacity := state.GetTRC21FeeCapacityFromStateWithToken(statedb, &tokenAddr)

				// Check sender balance against expected value
				if finalSenderBalance.Cmp(tt.expectedSenderBalanceAfterTx) != 0 {
					t.Errorf("Sender balance incorrect: got %v, want %v", finalSenderBalance, tt.expectedSenderBalanceAfterTx)
					t.Errorf("  Initial: %v, Transfer: %v", initialSenderBalance, tt.transferValue)
				}

				// Check miner balance against expected value
				if finalMinerBalance.Cmp(tt.expectedMinerBalanceAfterTx) != 0 {
					t.Errorf("Miner balance incorrect: got %v, want %v", finalMinerBalance, tt.expectedMinerBalanceAfterTx)
					t.Errorf("  Initial: %v", initialMinerBalance)
				}

				// Check token capacity against expected value
				if finalTokenCapacity == nil && tt.expectedTokenCapacityAfterTx.Sign() > 0 {
					t.Errorf("Token capacity is nil but expected %v", tt.expectedTokenCapacityAfterTx)
				} else if finalTokenCapacity != nil && finalTokenCapacity.Cmp(tt.expectedTokenCapacityAfterTx) != 0 {
					t.Errorf("Token capacity incorrect: got %v, want %v", finalTokenCapacity, tt.expectedTokenCapacityAfterTx)
				}

				// Calculate what actually happened for debugging
				actualSenderDeduction := new(big.Int).Sub(initialSenderBalance, finalSenderBalance)
				actualMinerGain := new(big.Int).Sub(finalMinerBalance, initialMinerBalance)
				actualFeePaid := new(big.Int).Sub(actualSenderDeduction, tt.transferValue)

				// Validate token fee behavior
				if tt.expectTokenFee {
					// Token fee scenario - sender should not pay ETH fee
					expectedFee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), common.TRC21GasPrice)
					if actualFeePaid.Sign() > 0 {
						t.Errorf("Expected no ETH fee for token transaction, but sender paid %v", actualFeePaid)
					}
					// Verify miner receives the fee from token capacity
					if actualMinerGain.Cmp(expectedFee) != 0 {
						t.Errorf("Miner fee mismatch: got %v, want %v", actualMinerGain, expectedFee)
					}
					// Verify token capacity reduction
					if initialTokenCapacity != nil && finalTokenCapacity != nil {
						expectedCapacityReduction := new(big.Int).Sub(initialTokenCapacity, finalTokenCapacity)
						if expectedCapacityReduction.Cmp(expectedFee) != 0 {
							t.Errorf("Token capacity reduction mismatch: got %v, want %v", expectedCapacityReduction, expectedFee)
						}
					}
				} else {
					// ETH fee fallback scenario - sender pays ETH fee
					// For ETH fallback, when token capacity is insufficient, the system may charge for gas limit
					expectedFeeForGasUsed := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), VRC25GasPrice)
					expectedFeeForGasLimit := new(big.Int).Mul(new(big.Int).SetUint64(tt.gasLimit), VRC25GasPrice)

					// Accept either gas used fee or gas limit fee depending on the implementation
					if actualFeePaid.Cmp(expectedFeeForGasUsed) != 0 && actualFeePaid.Cmp(expectedFeeForGasLimit) != 0 {
						t.Errorf("ETH fee mismatch: got %v, want %v (gas used) or %v (gas limit)", actualFeePaid, expectedFeeForGasUsed, expectedFeeForGasLimit)
					}
					// Verify fee conservation (sender fee = miner gain)
					if actualFeePaid.Cmp(actualMinerGain) != 0 {
						t.Errorf("Fee distribution inconsistent: sender fee %v != miner gain %v", actualFeePaid, actualMinerGain)
					}
				}

				// For gas refund scenarios, verify refund calculation
				if tt.gasUsed < tt.gasLimit {
					refundedGas := tt.gasLimit - tt.gasUsed
					t.Logf("Gas refund: %d gas units refunded", refundedGas)
				}

				// Log actual values for debugging
				t.Logf("Test: %s - Sender deduction: %v, Miner gain: %v, Transfer: %v, Fee: %v",
					tt.name, actualSenderDeduction, actualMinerGain, tt.transferValue, actualFeePaid)
				t.Logf("Token capacity: initial %v, final %v", initialTokenCapacity, finalTokenCapacity)
			} else {
				if err == nil && !failed {
					t.Error("Expected VRC25 transaction to fail but it succeeded")
				}
			}
		})
		fmt.Println("___________________________________________________________________________________")
	}
}

// TestComprehensiveFeeDistribution is the main comprehensive test function
func TestComprehensiveFeeDistribution(t *testing.T) {
	t.Run("Normal Transaction Tests", TestNormalTransactionFeeDistribution)
	t.Run("VRC25 Token Tests", TestVRC25TokenFeeDistribution)
}
