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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/tomox"
	"math/big"
)
import (
	"fmt"
	"runtime"
	"sync"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}
type CalculatedBlock struct {
	block *types.Block
	stop  bool
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config, balanceFee map[common.Address]*big.Int) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	if common.TIPSigning.Cmp(header.Number) == 0 {
		statedb.DeleteAddress(common.HexToAddress(common.BlockSigners))
	}
	InitSignerInTransactions(p.config, header, block.Transactions())
	balanceUpdated := map[common.Address]*big.Int{}
	totalFeeUsed := uint64(0)
	for i, tx := range block.Transactions() {
		// check black-list txs after hf
		if block.Number().Uint64() >= common.BlackListHFNumber {
			// check if sender is in black list
			if tx.From() != nil && common.Blacklist[*tx.From()] {
				return nil, nil, 0, fmt.Errorf("Block contains transaction with sender in black-list: %v", tx.From().Hex())
			}
			// check if receiver is in black list
			if tx.To() != nil && common.Blacklist[*tx.To()] {
				return nil, nil, 0, fmt.Errorf("Block contains transaction with receiver in black-list: %v", tx.To().Hex())
			}
		}
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, gas, err, tokenFeeUsed := ApplyTransaction(p.config, balanceFee, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
		if tokenFeeUsed {
			balanceFee[*tx.To()] = new(big.Int).Sub(balanceFee[*tx.To()], new(big.Int).SetUint64(gas))
			balanceUpdated[*tx.To()] = balanceFee[*tx.To()]
			totalFeeUsed = totalFeeUsed + gas
		}
	}
	state.UpdateTRC21Fee(statedb, balanceUpdated, totalFeeUsed)
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)
	return receipts, allLogs, *usedGas, nil
}

func (p *StateProcessor) ProcessBlockNoValidator(cBlock *CalculatedBlock, statedb *state.StateDB, cfg vm.Config, balanceFee map[common.Address]*big.Int) (types.Receipts, []*types.Log, uint64, error) {
	block := cBlock.block
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	if common.TIPSigning.Cmp(header.Number) == 0 {
		statedb.DeleteAddress(common.HexToAddress(common.BlockSigners))
	}
	if cBlock.stop {
		return nil, nil, 0, ErrStopPreparingBlock
	}
	InitSignerInTransactions(p.config, header, block.Transactions())
	balanceUpdated := map[common.Address]*big.Int{}
	totalFeeUsed := uint64(0)

	if cBlock.stop {
		return nil, nil, 0, ErrStopPreparingBlock
	}
	// Iterate over and process the individual transactions
	receipts = make([]*types.Receipt, block.Transactions().Len())
	for i, tx := range block.Transactions() {
		// check black-list txs after hf
		if block.Number().Uint64() >= common.BlackListHFNumber {
			// check if sender is in black list
			if tx.From() != nil && common.Blacklist[*tx.From()] {
				return nil, nil, 0, fmt.Errorf("Block contains transaction with sender in black-list: %v", tx.From().Hex())
			}
			// check if receiver is in black list
			if tx.To() != nil && common.Blacklist[*tx.To()] {
				return nil, nil, 0, fmt.Errorf("Block contains transaction with receiver in black-list: %v", tx.To().Hex())
			}
		}
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, gas, err, tokenFeeUsed := ApplyTransaction(p.config, balanceFee, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		if cBlock.stop {
			return nil, nil, 0, ErrStopPreparingBlock
		}
		receipts[i] = receipt
		allLogs = append(allLogs, receipt.Logs...)
		if tokenFeeUsed {
			balanceFee[*tx.To()] = new(big.Int).Sub(balanceFee[*tx.To()], new(big.Int).SetUint64(gas))
			balanceUpdated[*tx.To()] = balanceFee[*tx.To()]
			totalFeeUsed = totalFeeUsed + gas
		}
	}
	state.UpdateTRC21Fee(statedb, balanceUpdated, totalFeeUsed)
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)
	return receipts, allLogs, *usedGas, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, tokensFee map[common.Address]*big.Int, bc *BlockChain, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, uint64, error, bool) {
	if tx.To() != nil && tx.To().String() == common.BlockSigners && config.IsTIPSigning(header.Number) {
		return ApplySignTransaction(config, statedb, header, tx, usedGas)
	}
	if tx.IsMatchingTransaction() && config.IsTIPTomoX(header.Number) {
		return ApplyTomoXMatchedTransaction(config, statedb, header, tx, usedGas)
	}
	var balanceFee *big.Int
	if tx.To() != nil {
		if value, ok := tokensFee[*tx.To()]; ok {
			balanceFee = value
		}
	}
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number), balanceFee)
	if err != nil {
		return nil, 0, err, false
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, 0, err, false
	}
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	receipt := types.NewReceipt(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	if balanceFee != nil && failed {
		state.PayFeeWithTRC21TxFail(statedb, msg.From(), *tx.To())
	}
	return receipt, gas, err, balanceFee != nil
}

func ApplySignTransaction(config *params.ChainConfig, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64) (*types.Receipt, uint64, error, bool) {
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	from, err := types.Sender(types.MakeSigner(config, header.Number), tx)
	if err != nil {
		return nil, 0, err, false
	}
	nonce := statedb.GetNonce(from)
	if nonce < tx.Nonce() {
		return nil, 0, ErrNonceTooHigh, false
	} else if nonce > tx.Nonce() {
		return nil, 0, ErrNonceTooLow, false
	}
	statedb.SetNonce(from, nonce+1)
	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	receipt := types.NewReceipt(root, false, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = 0
	// if the transaction created a contract, store the creation address in the receipt.
	// Set the receipt logs and create a bloom for filtering
	log := &types.Log{}
	log.Address = common.HexToAddress(common.BlockSigners)
	log.BlockNumber = header.Number.Uint64()
	statedb.AddLog(log)
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt, 0, nil, false
}

func ApplyTomoXMatchedTransaction(config *params.ChainConfig, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64) (*types.Receipt, uint64, error, bool) {
	// Update the state with pending changes
	from, err := types.Sender(types.MakeSigner(config, header.Number), tx)
	if err != nil {
		return nil, 0, err, false
	}
	nonce := statedb.GetNonce(from)

	txMatchBatches, err := tomox.DecodeTxMatchesBatch(tx.Data())
	if err != nil || len(txMatchBatches.Data) == 0 {
		return nil, 0, err, false
	}
	gasUsed := big.NewInt(0)
	for _, txMatch := range txMatchBatches.Data {
		orderItem, err := txMatch.DecodeOrder()
		if err != nil {
			return nil, 0, err, false
		}
		takerAddr := orderItem.UserAddress
		takerExAddr := orderItem.ExchangeAddress
		takerExOwner := tomox.GetRelayerOwner(orderItem.ExchangeAddress, statedb)
		baseToken := orderItem.BaseToken
		quoteToken := orderItem.QuoteToken
		takerExfee := tomox.GetExRelayerFee(orderItem.ExchangeAddress, statedb)
		baseFee := common.TomoXBaseFee

		for i := 0; i < len(txMatch.Trades); i++ {
			price := tomox.ToBigInt(txMatch.Trades[i][tomox.TradedPrice])
			quantityString := txMatch.Trades[i][tomox.TradedQuantity]
			quantity := tomox.ToBigInt(quantityString)
			if price.Cmp(big.NewInt(0)) <= 0 || quantity.Cmp(big.NewInt(0)) <= 0 {
				return nil, 0, fmt.Errorf("trade misses important information. tradedPrice %v, tradedQuantity %v", price, quantity), false
			}
			makerExAddr := common.HexToAddress(txMatch.Trades[i][tomox.TradedMakerExchangeAddress])
			makerExfee := tomox.GetExRelayerFee(makerExAddr, statedb)
			makerExOwner := tomox.GetRelayerOwner(makerExAddr, statedb)
			makerAddr := common.HexToAddress(txMatch.Trades[i][tomox.TradedMaker])
			log.Debug("ApplyTomoXMatchedTransaction : trades quantityString", "i", i, "trade", txMatch.Trades[i], "price", price)
			if makerExAddr != (common.Address{}) && makerAddr != (common.Address{}) {
				// take relayer fee
				err := tomox.SubRelayerFee(takerExAddr, common.RelayerFee, statedb)
				if err != nil {
					return nil, 0, err, false
				}
				err = tomox.SubRelayerFee(makerExAddr, common.RelayerFee, statedb)
				if err != nil {
					return nil, 0, err, false
				}

				// masternodes charges fee of both 2 relayers. If maker and taker are on same relayer, that relayer is charged fee twice
				gasUsed = gasUsed.Add(gasUsed, common.RelayerFee)
				gasUsed = gasUsed.Add(gasUsed, common.RelayerFee)

				//log.Debug("ApplyTomoXMatchedTransaction quantity check", "i", i, "trade", txMatch.Trades[i], "price", price, "quantity", quantity)

				isTakerBuy := orderItem.Side == tomox.Bid
				settleBalanceResult := tomox.SettleBalance(
					makerAddr,
					takerAddr,
					baseToken,
					quoteToken,
					isTakerBuy,
					makerExfee,
					takerExfee,
					baseFee,
					quantity,
					price)

				// TAKER
				//log.Debug("ApplyTomoXMatchedTransaction settle balance for taker",
				//	"taker", takerAddr,
				//	"inToken", settleBalanceResult[takerAddr][tomox.InToken].(common.Address), "inQuantity", settleBalanceResult[takerAddr][tomox.InQuantity].(*big.Int),
				//	"inTotal", settleBalanceResult[takerAddr][tomox.InTotal].(*big.Int),
				//	"outToken", settleBalanceResult[takerAddr][tomox.OutToken].(common.Address), "outQuantity", settleBalanceResult[takerAddr][tomox.OutQuantity].(*big.Int),
				//	"outTotal", settleBalanceResult[takerAddr][tomox.OutTotal].(*big.Int))
				err = tomox.AddTokenBalance(takerAddr, settleBalanceResult[takerAddr][tomox.InTotal].(*big.Int), settleBalanceResult[takerAddr][tomox.InToken].(common.Address), statedb)
				if err != nil {
					return nil, 0, err, false
				}
				err = tomox.SubTokenBalance(takerAddr, settleBalanceResult[takerAddr][tomox.OutTotal].(*big.Int), settleBalanceResult[takerAddr][tomox.OutToken].(common.Address), statedb)
				if err != nil {
					return nil, 0, err, false
				}

				// MAKER
				//log.Debug("ApplyTomoXMatchedTransaction settle balance for maker",
				//	"maker", makerAddr,
				//	"inToken", settleBalanceResult[makerAddr][tomox.InToken].(common.Address), "inQuantity", settleBalanceResult[makerAddr][tomox.InQuantity].(*big.Int),
				//	"inTotal", settleBalanceResult[makerAddr][tomox.InTotal].(*big.Int),
				//	"outToken", settleBalanceResult[makerAddr][tomox.OutToken].(common.Address), "outQuantity", settleBalanceResult[makerAddr][tomox.OutQuantity].(*big.Int),
				//	"outTotal", settleBalanceResult[makerAddr][tomox.OutTotal].(*big.Int))
				err = tomox.AddTokenBalance(makerAddr, settleBalanceResult[makerAddr][tomox.InTotal].(*big.Int), settleBalanceResult[makerAddr][tomox.InToken].(common.Address), statedb)
				if err != nil {
					return nil, 0, err, false
				}
				err = tomox.SubTokenBalance(makerAddr, settleBalanceResult[makerAddr][tomox.OutTotal].(*big.Int), settleBalanceResult[makerAddr][tomox.OutToken].(common.Address), statedb)
				if err != nil {
					return nil, 0, err, false
				}

				// add balance for relayers
				//log.Debug("ApplyTomoXMatchedTransaction settle fee for relayers",
				//	"takerRelayerOwner", takerExOwner,
				//	"takerFeeToken", quoteToken, "takerFee", settleBalanceResult[takerAddr][tomox.Fee].(*big.Int),
				//	"makerRelayerOwner", makerExOwner,
				//	"makerFeeToken", quoteToken, "makerFee", settleBalanceResult[makerAddr][tomox.Fee].(*big.Int))
				// takerFee
				err = tomox.AddTokenBalance(takerExOwner, settleBalanceResult[takerAddr][tomox.Fee].(*big.Int), quoteToken, statedb)
				if err != nil {
					return nil, 0, err, false
				}
				// makerFee
				err = tomox.AddTokenBalance(makerExOwner, settleBalanceResult[makerAddr][tomox.Fee].(*big.Int), quoteToken, statedb)
				if err != nil {
					return nil, 0, err, false
				}
			}
		}
	}
	statedb.SetNonce(from, nonce+1)
	statedb.AddBalance(header.Coinbase, gasUsed)
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	receipt := types.NewReceipt(root, false, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gasUsed.Uint64()
	// if the transaction created a contract, store the creation address in the receipt.
	// Set the receipt logs and create a bloom for filtering
	log := &types.Log{}
	log.Address = common.HexToAddress(common.BlockSigners)
	log.BlockNumber = header.Number.Uint64()
	statedb.AddLog(log)
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt, 0, nil, false
}

func InitSignerInTransactions(config *params.ChainConfig, header *types.Header, txs types.Transactions) {
	nWorker := runtime.NumCPU()
	signer := types.MakeSigner(config, header.Number)
	chunkSize := txs.Len() / nWorker
	if txs.Len()%nWorker != 0 {
		chunkSize++
	}
	wg := sync.WaitGroup{}
	wg.Add(nWorker)
	for i := 0; i < nWorker; i++ {
		from := i * chunkSize
		to := from + chunkSize
		if to > txs.Len() {
			to = txs.Len()
		}
		go func(from int, to int) {
			for j := from; j < to; j++ {
				types.CacheSigner(signer, txs[j])
				txs[j].CacheHash()
			}
			wg.Done()
		}(from, to)
	}
	wg.Wait()
}
