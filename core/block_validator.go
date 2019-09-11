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
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/posv"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/tomox"
	"sort"
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

	engine, _ := v.engine.(*posv.Posv)
	tomoXService := engine.GetTomoXService()

	currentState, err := v.bc.State()
	if err != nil {
		return err
	}

	// validate matchedOrder txs
	processedHashes := []common.Hash{}

	// clear the previous dry-run cache
	if tomoXService != nil {
		tomoXService.GetDB().InitDryRunMode()
	}
	txMatchBatchData, err := ExtractMatchingTransactions(block.Transactions())
	if err != nil {
		return err
	}
	for _, txMatchBatch := range txMatchBatchData {
		if tomoXService == nil {
			log.Error("tomox not found")
			return tomox.ErrTomoXServiceNotFound
		}
		log.Debug("Verify matching transaction", "txHash", txMatchBatch.TxHash)
		hashes, err := v.validateMatchingOrder(tomoXService, currentState, txMatchBatch)
		if err != nil {
			return err
		}
		processedHashes = append(processedHashes, hashes...)
	}
	hashNoValidator := block.HashNoValidator()
	if _, ok := v.bc.processedOrderHashes.Get(hashNoValidator); !ok && len(processedHashes) > 0 {
		v.bc.processedOrderHashes.Add(hashNoValidator, processedHashes)
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

func (v *BlockValidator) validateMatchingOrder(tomoXService *tomox.TomoX, currentState *state.StateDB, txMatchBatch tomox.TxMatchBatch) ([]common.Hash, error) {
	processedHashes := []common.Hash{}
	log.Debug("verify matching transaction found a TxMatches Batch", "numTxMatches", len(txMatchBatch.Data))

	for _, txMatch := range txMatchBatch.Data {
		// verify orderItem
		order, err := txMatch.DecodeOrder()
		if err != nil {
			return []common.Hash{}, fmt.Errorf("transaction match is corrupted. Failed decode order. Error: %s ", err)
		}
		if tomoXService.ExistProcessedOrderHash(order.Hash) {
			log.Debug("This order has been processed", "hash", hex.EncodeToString(order.Hash.Bytes()))
			continue
		}
		log.Debug("process tx match", "order", order)

		processedHashes = append(processedHashes, order.Hash)


		// Remove order from db pending.
		if err := tomoXService.RemovePendingHash(order.Hash); err != nil {
			log.Debug("Fail to remove pending hash", "err", err)
		}
		if err := tomoXService.RemoveOrderPending(order.Hash); err != nil {
			log.Debug("Fail to remove order pending", "err", err)
		}

		// SDK node doesn't need to run ME
		if tomoXService.IsSDKNode() {
			log.Debug("SDK node ignore running matching engine")
			continue
		}
		if err := order.VerifyMatchedOrder(currentState); err != nil {
			return []common.Hash{}, err
		}

		ob, err := tomoXService.GetOrderBook(order.PairName, true)
		// if orderbook of this pairName has been updated by previous tx in this block, use it

		if err != nil {
			return []common.Hash{}, err
		}

		// verify old state: orderbook hash, bidTree hash, askTree hash
		if err := txMatch.VerifyOldTomoXState(ob); err != nil {
			return []common.Hash{}, err
		}

		// process Matching Engine
		if _, _, err := ob.ProcessOrder(order, true, true); err != nil {
			return []common.Hash{}, err
		}

		// verify new state
		if err := txMatch.VerifyNewTomoXState(ob); err != nil {
			return []common.Hash{}, err
		}
	}

	return processedHashes, nil
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

func ExtractMatchingTransactions(transactions types.Transactions) ([]tomox.TxMatchBatch, error) {
	txMatchBatchData := []tomox.TxMatchBatch{}
	for _, tx := range transactions {
		if tx.IsMatchingTransaction() {
			txMatchBatch, err := tomox.DecodeTxMatchesBatch(tx.Data())
			if err != nil {
				return []tomox.TxMatchBatch{}, fmt.Errorf("transaction match is corrupted. Failed to decode txMatchBatch. Error: %s", err)
			}
			txMatchBatch.TxHash = tx.Hash()
			txMatchBatchData = append(txMatchBatchData, txMatchBatch)
		}
	}
	sort.Slice(txMatchBatchData, func(i, j int) bool {
		return txMatchBatchData[i].Timestamp < txMatchBatchData[j].Timestamp
	})
	return txMatchBatchData, nil
}
