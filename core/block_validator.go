// Copyright 2015 The go-ethereum Authors
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

package core

import (
	"fmt"
	"github.com/tomochain/tomochain/consensus/posv"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending/lendingstate"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
)

// BlockValidator is responsible for validating block headers, uncles and
// processed state.
//
// BlockValidator implements Validator.
type BlockValidator struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody validates the given block's uncles and verifies the the block
// header's transaction and uncle roots. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(block *types.Block) error {
	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		return ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			return consensus.ErrUnknownAncestor
		}
		return consensus.ErrPrunedAncestor
	}
	// Header validity is known at this point, check the uncles and transactions
	header := block.Header()
	if err := v.engine.VerifyUncles(v.bc, block); err != nil {
		return err
	}
	if hash := types.CalcUncleHash(block.Uncles()); hash != header.UncleHash {
		return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	}
	if hash := types.DeriveSha(block.Transactions()); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}
	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as amount of used gas, the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(block, parent *types.Block, statedb *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	header := block.Header()
	if block.GasUsed() != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", block.GasUsed(), usedGas)
	}
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveSha(receipts)
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(v.config.IsEIP158(header.Number)); header.Root != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	return nil
}

func (v *BlockValidator) ValidateTradingOrder(statedb *state.StateDB, tomoxStatedb *tradingstate.TradingStateDB, txMatchBatch tradingstate.TxMatchBatch, coinbase common.Address) error {
	posvEngine, ok := v.bc.Engine().(*posv.Posv)
	if posvEngine == nil || !ok {
		return ErrNotPoSV
	}
	tomoXService := posvEngine.GetTomoXService()
	if tomoXService == nil {
		return fmt.Errorf("tomox not found")
	}
	log.Debug("verify matching transaction found a TxMatches Batch", "numTxMatches", len(txMatchBatch.Data))
	tradingResult := map[common.Hash]tradingstate.MatchingResult{}
	for _, txMatch := range txMatchBatch.Data {
		// verify orderItem
		order, err := txMatch.DecodeOrder()
		if err != nil {
			return fmt.Errorf("transaction match is corrupted. Failed decode order. Error: %s ", err)
		}

		log.Debug("process tx match", "order", order)
		if err := order.VerifyOrder(statedb); err != nil {
			return fmt.Errorf("invalid order . Error: %v", err)
		}
		// process Matching Engine
		newTrades, newRejectedOrders, err := tomoXService.ApplyOrder(coinbase, v.bc, statedb, tomoxStatedb, tradingstate.GetTradingOrderBookHash(order.BaseToken, order.QuoteToken), order)
		if err != nil {
			return err
		}
		tradingResult[order.Hash] = tradingstate.MatchingResult{
			Trades:  newTrades,
			Rejects: newRejectedOrders,
		}
	}
	if tomoXService.IsSDKNode() {
		v.bc.AddMatchingResult(txMatchBatch.TxHash, tradingResult)
	}
	return nil
}

func (v *BlockValidator) ValidateLendingOrder(statedb *state.StateDB, lendingStateDb *lendingstate.LendingStateDB, tomoxStatedb *tradingstate.TradingStateDB, batch lendingstate.TxLendingBatch, coinbase common.Address) error {
	posvEngine, ok := v.bc.Engine().(*posv.Posv)
	if posvEngine == nil || !ok {
		return ErrNotPoSV
	}
	tomoXService := posvEngine.GetTomoXService()
	if tomoXService == nil {
		return fmt.Errorf("tomox not found")
	}
	lendingService := posvEngine.GetLendingService()
	if lendingService == nil {
		return fmt.Errorf("lendingService not found")
	}
	log.Debug("verify lendingItem ", "numItems", len(batch.Data))
	lendingResult := map[common.Hash]lendingstate.MatchingResult{}
	for _, l := range batch.Data {
		// verify lendingItem

		log.Debug("process lending tx", "lendingItem", lendingstate.ToJSON(l))
		if err := l.VerifyLendingItem(statedb); err != nil {
			return fmt.Errorf("invalid lendingItem . Error: %v", err)
		}
		// process Matching Engine
		newTrades, newRejectedOrders, err := lendingService.ApplyOrder(uint64(batch.Timestamp), coinbase, v.bc, statedb, lendingStateDb, tomoxStatedb, lendingstate.GetLendingOrderBookHash(l.LendingToken, l.Term), l)
		if err != nil {
			return err
		}
		lendingResult[l.Hash] = lendingstate.MatchingResult{
			Trades:  newTrades,
			Rejects: newRejectedOrders,
		}
	}
	if tomoXService.IsSDKNode() {
		v.bc.AddLendingResult(batch.TxHash, lendingResult)
	}
	return nil
}

// CalcGasLimit computes the gas limit of the next block after parent.
// This is miner strategy, not consensus protocol.
func CalcGasLimit(parent *types.Block) uint64 {
	// contrib = (parentGasUsed * 3 / 2) / 1024
	contrib := (parent.GasUsed() + parent.GasUsed()/2) / params.GasLimitBoundDivisor

	// decay = parentGasLimit / 1024 -1
	decay := parent.GasLimit()/params.GasLimitBoundDivisor - 1

	/*
		strategy: gasLimit of block-to-mine is set based on parent's
		gasUsed value.  if parentGasUsed > parentGasLimit * (2/3) then we
		increase it, otherwise lower it (or leave it unchanged if it's right
		at that usage) the amount increased/decreased depends on how far away
		from parentGasLimit * (2/3) parentGasUsed is.
	*/
	limit := parent.GasLimit() - decay + contrib
	if limit < params.MinGasLimit {
		limit = params.MinGasLimit
	}
	// however, if we're now below the target (TargetGasLimit) we increase the
	// limit as much as we can (parentGasLimit / 1024 -1)
	if limit < params.TargetGasLimit {
		limit = parent.GasLimit() + decay
		if limit > params.TargetGasLimit {
			limit = params.TargetGasLimit
		}
	}
	return limit
}

func ExtractTradingTransactions(transactions types.Transactions) ([]tradingstate.TxMatchBatch, error) {
	txMatchBatchData := []tradingstate.TxMatchBatch{}
	for _, tx := range transactions {
		if tx.IsTradingTransaction() {
			txMatchBatch, err := tradingstate.DecodeTxMatchesBatch(tx.Data())
			if err != nil {
				return []tradingstate.TxMatchBatch{}, fmt.Errorf("transaction match is corrupted. Failed to decode txMatchBatch. Error: %s", err)
			}
			txMatchBatch.TxHash = tx.Hash()
			txMatchBatchData = append(txMatchBatchData, txMatchBatch)
		}
	}
	return txMatchBatchData, nil
}

func ExtractLendingTransactions(transactions types.Transactions) ([]lendingstate.TxLendingBatch, error) {
	batchData := []lendingstate.TxLendingBatch{}
	for _, tx := range transactions {
		if tx.IsLendingTransaction() {
			txMatchBatch, err := lendingstate.DecodeTxLendingBatch(tx.Data())
			if err != nil {
				return []lendingstate.TxLendingBatch{}, fmt.Errorf("transaction match is corrupted. Failed to decode lendingTransaction. Error: %s", err)
			}
			txMatchBatch.TxHash = tx.Hash()
			batchData = append(batchData, txMatchBatch)
		}
	}
	return batchData, nil
}

func ExtractLendingLiquidatedTradeTransactions(transactions types.Transactions) (lendingstate.LiquidatedResult, error) {
	for _, tx := range transactions {
		if tx.IsLendingLiquidatedTradeTransaction() {
			liquidatedTrades, err := lendingstate.DecodeLiquidatedResult(tx.Data())
			if err != nil {
				return lendingstate.LiquidatedResult{}, fmt.Errorf("transaction is corrupted. Failed to decode LendingClosedTradeTransaction. Error: %s", err)
			}
			liquidatedTrades.TxHash = tx.Hash()
			// each block has only one tx of this type
			return liquidatedTrades, nil
		}
	}
	return lendingstate.LiquidatedResult{}, nil
}
