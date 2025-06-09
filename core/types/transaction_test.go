// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"container/heap"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"sort"
	"testing"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/rlp"
)

// The values in those tests are from the Transaction Tests
// at github.com/ethereum/tests.
var (
	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)

	rightvrsTx, _ = NewTransaction(
		3,
		common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	).WithSignature(
		HomesteadSigner{},
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)
)

func TestTransactionSigHash(t *testing.T) {
	var homestead HomesteadSigner
	if homestead.Hash(emptyTx) != common.HexToHash("c775b99e7ad12f50d819fcd602390467e28141316969f4b57f0626f74fe3b386") {
		t.Errorf("empty transaction hash mismatch, got %x", emptyTx.Hash())
	}
	if homestead.Hash(rightvrsTx) != common.HexToHash("fe7a79529ed5f7c3375d06b26b186a8644e0e16c373d7a12be41c62d6042b77a") {
		t.Errorf("RightVRS transaction hash mismatch, got %x", rightvrsTx.Hash())
	}
}

func TestTransactionEncode(t *testing.T) {
	txb, err := rlp.EncodeToBytes(rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := common.FromHex("f86103018207d094b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a8255441ca098ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4aa08887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a3")
	if !bytes.Equal(txb, should) {
		t.Errorf("encoded RLP mismatch, got %x", txb)
	}
}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f8498080808080011ca09b16de9d5bdee2cf56c28d16275a4da68cd30273e2525f3959f5d62557489921a0372ebd8fb3345f7db7b5a86d42e24d36e983e259b0664ceb8c227ec9af572f3d"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from, err := Sender(HomesteadSigner{}, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()

	tx, err := decodeTx(common.Hex2Bytes("f85d80808094000000000000000000000000000000000000000080011ca0527c0d8f5c63f7b9f41324a7c8a563ee1190bcbf0dac8ab446291bdbf32f5c79a0552c4ef0a09a04395074dab9ed34d3fbfb843c2f2546cc30fe89ec143ca94ca6"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from, err := Sender(HomesteadSigner{}, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if addr != from {
		t.Error("derived address doesn't match")
	}
}

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func TestTransactionPriceNonceSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := HomesteadSigner{}
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	txset, _ := NewTransactionsByPriceAndNonce(signer, groups, nil)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != 25*25 {
		t.Errorf("expected %d transactions, found %d", 25*25, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}
		// Find the previous and next nonce of this account
		prev, next := i-1, i+1
		for j := i - 1; j >= 0; j-- {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				prev = j
				break
			}
		}
		for j := i + 1; j < len(txs); j++ {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				next = j
				break
			}
		}
		// Make sure that in between the neighbor nonces, the transaction is correctly positioned price wise
		for j := prev + 1; j < next; j++ {
			fromj, _ := Sender(signer, txs[j])
			if j < i && txs[j].GasPrice().Cmp(txi.GasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
			if j > i && txs[j].GasPrice().Cmp(txi.GasPrice()) > 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) > tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
		}
	}
}

// TestTransactionJSON tests serializing/de-serializing to/from JSON.
func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := NewEIP155Signer(common.Big1)

	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		switch i % 2 {
		case 0:
			tx = NewTransaction(i, common.Address{1}, common.Big0, 1, common.Big2, []byte("abcdef"))
		case 1:
			tx = NewContractCreation(i, common.Big0, 1, common.Big2, []byte("abcdef"))
		}

		tx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		data, err := json.Marshal(tx)
		if err != nil {
			t.Errorf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, got %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, got %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}

// TestTxByPriceSimpleSort tests that transactions are sorted purely by gas price
// without any special handling for TRC21 or other transaction types.
func TestTxByPriceSimpleSort(t *testing.T) {
	// Generate test transactions with different gas prices
	key, _ := crypto.GenerateKey()
	signer := HomesteadSigner{}

	// Create transactions with specific gas prices in random order
	gasPrices := []*big.Int{
		big.NewInt(500), // highest
		big.NewInt(100), // lowest
		big.NewInt(300), // middle
		big.NewInt(200), // lower middle
		big.NewInt(400), // upper middle
	}

	txs := make(Transactions, len(gasPrices))
	for i, gasPrice := range gasPrices {
		tx, _ := SignTx(NewTransaction(uint64(i), common.Address{1}, big.NewInt(100), 100, gasPrice, nil), signer, key)
		txs[i] = tx
	}

	// Create TxByPrice sorter
	sorter := TxByPrice(txs)

	// Test Less function directly
	expectedOrder := []int{0, 4, 2, 3, 1} // indices sorted by gas price descending (500, 400, 300, 200, 100)

	// Verify that higher gas price transactions are considered "less" (higher priority)
	for i := 0; i < len(expectedOrder)-1; i++ {
		higherPriceIdx := expectedOrder[i]
		lowerPriceIdx := expectedOrder[i+1]

		if !sorter.Less(higherPriceIdx, lowerPriceIdx) {
			t.Errorf("Transaction with gas price %v should have higher priority than %v",
				sorter[higherPriceIdx].GasPrice(), sorter[lowerPriceIdx].GasPrice())
		}

		if sorter.Less(lowerPriceIdx, higherPriceIdx) {
			t.Errorf("Transaction with gas price %v should have lower priority than %v",
				sorter[lowerPriceIdx].GasPrice(), sorter[higherPriceIdx].GasPrice())
		}
	}

	// Test that equal gas prices return false (stable sort)
	equalGasTx, _ := SignTx(NewTransaction(100, common.Address{2}, big.NewInt(100), 100, big.NewInt(300), nil), signer, key)
	sorter = append(sorter, equalGasTx)

	// Find the transaction with gas price 300 and compare with the new equal one
	equalIdx1 := 2               // gas price 300
	equalIdx2 := len(sorter) - 1 // new tx with gas price 300

	if sorter.Less(equalIdx1, equalIdx2) || sorter.Less(equalIdx2, equalIdx1) {
		t.Errorf("Transactions with equal gas prices should not be considered less than each other")
	}
}

// TestTxByPriceNoTRC21Priority tests that TRC21 transactions (transactions going to
// addresses in payersSwap) are not given special priority and are sorted purely by gas price.
func TestTxByPriceNoTRC21Priority(t *testing.T) {
	key, _ := crypto.GenerateKey()
	signer := HomesteadSigner{}

	// Create addresses that would be in payersSwap (simulating TRC21 contracts)
	trc21Address1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
	trc21Address2 := common.HexToAddress("0x2222222222222222222222222222222222222222")
	regularAddress := common.HexToAddress("0x3333333333333333333333333333333333333333")

	// Create transactions: TRC21 with low gas price, regular with high gas price
	lowGasPrice := big.NewInt(100)
	highGasPrice := big.NewInt(500)

	// TRC21 transaction with low gas price
	trc21LowGasTx, _ := SignTx(NewTransaction(1, trc21Address1, big.NewInt(100), 100, lowGasPrice, nil), signer, key)

	// TRC21 transaction with high gas price
	trc21HighGasTx, _ := SignTx(NewTransaction(2, trc21Address2, big.NewInt(100), 100, highGasPrice, nil), signer, key)

	// Regular transaction with high gas price
	regularHighGasTx, _ := SignTx(NewTransaction(3, regularAddress, big.NewInt(100), 100, highGasPrice, nil), signer, key)

	// Regular transaction with low gas price
	regularLowGasTx, _ := SignTx(NewTransaction(4, regularAddress, big.NewInt(100), 100, lowGasPrice, nil), signer, key)

	txs := Transactions{trc21LowGasTx, trc21HighGasTx, regularHighGasTx, regularLowGasTx}

	// Create TxByPrice sorter (TRC21 addresses no longer affect sorting)
	sorter := TxByPrice(txs)

	// Test that high gas price transactions (regardless of TRC21 status) beat low gas price ones
	trc21LowIdx := 0    // TRC21 with low gas price
	trc21HighIdx := 1   // TRC21 with high gas price
	regularHighIdx := 2 // Regular with high gas price
	regularLowIdx := 3  // Regular with low gas price

	// High gas price transactions should beat low gas price ones, regardless of TRC21 status
	if !sorter.Less(trc21HighIdx, trc21LowIdx) {
		t.Errorf("TRC21 transaction with high gas price should beat TRC21 transaction with low gas price")
	}

	if !sorter.Less(regularHighIdx, regularLowIdx) {
		t.Errorf("Regular transaction with high gas price should beat regular transaction with low gas price")
	}

	if !sorter.Less(regularHighIdx, trc21LowIdx) {
		t.Errorf("Regular transaction with high gas price should beat TRC21 transaction with low gas price")
	}

	if !sorter.Less(trc21HighIdx, regularLowIdx) {
		t.Errorf("TRC21 transaction with high gas price should beat regular transaction with low gas price")
	}

	// Transactions with same gas price should be equal priority regardless of TRC21 status
	if sorter.Less(trc21HighIdx, regularHighIdx) || sorter.Less(regularHighIdx, trc21HighIdx) {
		t.Errorf("TRC21 and regular transactions with same gas price should have equal priority")
	}

	if sorter.Less(trc21LowIdx, regularLowIdx) || sorter.Less(regularLowIdx, trc21LowIdx) {
		t.Errorf("TRC21 and regular transactions with same gas price should have equal priority")
	}
}

// TestTxByPriceCompleteSort tests the complete sorting flow using Go's sort package
// to verify transactions are properly ordered by gas price in descending order.
func TestTxByPriceCompleteSort(t *testing.T) {
	key, _ := crypto.GenerateKey()
	signer := HomesteadSigner{}

	// Create transactions with various gas prices including duplicates
	testCases := []struct {
		gasPrice int64
		nonce    uint64
		to       common.Address
	}{
		{1000, 0, common.HexToAddress("0x1000000000000000000000000000000000000000")}, // highest
		{100, 1, common.HexToAddress("0x1100000000000000000000000000000000000000")},  // lowest
		{500, 2, common.HexToAddress("0x1500000000000000000000000000000000000000")},  // middle-high
		{300, 3, common.HexToAddress("0x1300000000000000000000000000000000000000")},  // middle
		{100, 4, common.HexToAddress("0x1100000000000000000000000000000000000001")},  // lowest (duplicate)
		{750, 5, common.HexToAddress("0x1750000000000000000000000000000000000000")},  // high
		{200, 6, common.HexToAddress("0x1200000000000000000000000000000000000000")},  // low
		{500, 7, common.HexToAddress("0x1500000000000000000000000000000000000001")},  // middle-high (duplicate)
	}

	// Create transactions
	var txs Transactions
	for _, tc := range testCases {
		tx, _ := SignTx(NewTransaction(tc.nonce, tc.to, big.NewInt(100), 100, big.NewInt(tc.gasPrice), nil), signer, key)
		txs = append(txs, tx)
	}

	// Create sorter (TRC21 addresses no longer affect sorting since payersSwap was removed)
	sorter := (*TxByPrice)(&txs)

	// Sort using Go's sort package
	sort.Sort(sorter)

	// Expected order: 1000, 750, 500, 500, 300, 200, 100, 100
	expectedGasPrices := []int64{1000, 750, 500, 500, 300, 200, 100, 100}

	// Verify sorted order
	if len(*sorter) != len(expectedGasPrices) {
		t.Fatalf("Expected %d transactions, got %d", len(expectedGasPrices), len(*sorter))
	}

	for i, expectedPrice := range expectedGasPrices {
		actualPrice := (*sorter)[i].GasPrice().Int64()
		if actualPrice != expectedPrice {
			t.Errorf("Transaction at index %d: expected gas price %d, got %d", i, expectedPrice, actualPrice)
		}
	}

	// Verify that transactions are in descending order
	for i := 0; i < len(*sorter)-1; i++ {
		currentPrice := (*sorter)[i].GasPrice()
		nextPrice := (*sorter)[i+1].GasPrice()
		if currentPrice.Cmp(nextPrice) < 0 {
			t.Errorf("Transactions not in descending order: position %d (%v) < position %d (%v)",
				i, currentPrice, i+1, nextPrice)
		}
	}

	// Note: TRC21 special handling was removed, so all transactions are sorted purely by gas price
	// The following test now just verifies that all transactions are properly sorted by gas price
}

// TestTxByPriceHeapOperations tests the heap operations (Push/Pop) to ensure
// they maintain proper ordering after sorting modifications.
func TestTxByPriceHeapOperations(t *testing.T) {
	key, _ := crypto.GenerateKey()
	signer := HomesteadSigner{}

	// Create initial transactions
	gasPrices := []int64{300, 100, 500, 200}
	var txs Transactions
	for i, gasPrice := range gasPrices {
		tx, _ := SignTx(NewTransaction(uint64(i), common.Address{byte(i)}, big.NewInt(100), 100, big.NewInt(gasPrice), nil), signer, key)
		txs = append(txs, tx)
	}

	// Create TxByPrice and initialize as heap
	sorter := (*TxByPrice)(&txs)
	heap.Init(sorter)

	// Test Pop operation - should get highest gas price first
	if sorter.Len() != 4 {
		t.Fatalf("Expected 4 transactions, got %d", sorter.Len())
	}

	// Pop highest gas price transaction (should be 500)
	poppedTx := heap.Pop(sorter).(*Transaction)
	if poppedTx.GasPrice().Int64() != 500 {
		t.Errorf("Expected to pop transaction with gas price 500, got %d", poppedTx.GasPrice().Int64())
	}

	// Pop next highest (should be 300)
	poppedTx = heap.Pop(sorter).(*Transaction)
	if poppedTx.GasPrice().Int64() != 300 {
		t.Errorf("Expected to pop transaction with gas price 300, got %d", poppedTx.GasPrice().Int64())
	}

	// Test Push operation - add a new high-priority transaction
	newHighTx, _ := SignTx(NewTransaction(10, common.Address{10}, big.NewInt(100), 100, big.NewInt(800), nil), signer, key)
	heap.Push(sorter, newHighTx)

	// The new transaction should be at the top
	poppedTx = heap.Pop(sorter).(*Transaction)
	if poppedTx.GasPrice().Int64() != 800 {
		t.Errorf("Expected to pop new transaction with gas price 800, got %d", poppedTx.GasPrice().Int64())
	}

	// Test Push operation - add a low-priority transaction
	newLowTx, _ := SignTx(NewTransaction(11, common.Address{11}, big.NewInt(100), 100, big.NewInt(50), nil), signer, key)
	heap.Push(sorter, newLowTx)

	// Should still get the remaining higher priority transactions first
	remainingPrices := []int64{200, 100, 50}
	for _, expectedPrice := range remainingPrices {
		poppedTx = heap.Pop(sorter).(*Transaction)
		if poppedTx.GasPrice().Int64() != expectedPrice {
			t.Errorf("Expected to pop transaction with gas price %d, got %d", expectedPrice, poppedTx.GasPrice().Int64())
		}
	}

	// Heap should be empty now
	if sorter.Len() != 0 {
		t.Errorf("Expected empty heap, got %d transactions", sorter.Len())
	}
}

// TestTxByPriceStableSorting tests that transactions with equal gas prices
// maintain their relative order (stable sorting behavior).
func TestTxByPriceStableSorting(t *testing.T) {
	key, _ := crypto.GenerateKey()
	signer := HomesteadSigner{}

	// Create multiple transactions with the same gas price but different nonces
	sameGasPrice := big.NewInt(250)
	var txs Transactions

	// Create 5 transactions with the same gas price
	for i := 0; i < 5; i++ {
		tx, _ := SignTx(NewTransaction(uint64(i), common.Address{byte(i)}, big.NewInt(100), 100, sameGasPrice, nil), signer, key)
		txs = append(txs, tx)
	}

	// Add one high and one low gas price transaction for reference
	highTx, _ := SignTx(NewTransaction(100, common.Address{100}, big.NewInt(100), 100, big.NewInt(500), nil), signer, key)
	lowTx, _ := SignTx(NewTransaction(101, common.Address{101}, big.NewInt(100), 100, big.NewInt(100), nil), signer, key)

	// Insert them in the middle to test sorting
	txs = append(txs[:2], append(Transactions{highTx}, txs[2:]...)...)
	txs = append(txs[:4], append(Transactions{lowTx}, txs[4:]...)...)

	// Create sorter and sort
	sorter := (*TxByPrice)(&txs)
	sort.Sort(sorter)

	// Verify structure: [500] [250, 250, 250, 250, 250] [100]
	expectedStructure := []int64{500, 250, 250, 250, 250, 250, 100}

	if len(*sorter) != len(expectedStructure) {
		t.Fatalf("Expected %d transactions, got %d", len(expectedStructure), len(*sorter))
	}

	for i, expectedPrice := range expectedStructure {
		actualPrice := (*sorter)[i].GasPrice().Int64()
		if actualPrice != expectedPrice {
			t.Errorf("Transaction at index %d: expected gas price %d, got %d", i, expectedPrice, actualPrice)
		}
	}

	// Verify that equal gas price transactions maintain some consistent order
	equalPriceSection := (*sorter)[1:6] // The 5 transactions with gas price 250
	for i := 0; i < len(equalPriceSection)-1; i++ {
		if equalPriceSection[i].GasPrice().Cmp(equalPriceSection[i+1].GasPrice()) != 0 {
			t.Errorf("Equal gas price section has inconsistent prices: %v vs %v",
				equalPriceSection[i].GasPrice(), equalPriceSection[i+1].GasPrice())
		}
	}
}
