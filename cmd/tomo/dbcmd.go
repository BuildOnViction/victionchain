// Copyright (c) 2020 Victionchain
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// this program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/tomochain/tomochain/cmd/utils"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus/posv"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/tomox"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending"
	"github.com/tomochain/tomochain/trie"
	"gopkg.in/urfave/cli.v1"
)

var (
	dbCommand = cli.Command{
		Name:      "db",
		Usage:     "Low level database operations",
		ArgsUsage: "",
		Category:  "DATABASE COMMANDS",
		Subcommands: []cli.Command{
			dbRepairSnapshotCmd,
		},
	}
	dbRepairSnapshotCmd = cli.Command{
		Action: utils.MigrateFlags(dbRepairSnapshot),
		Name:   "repair-snapshot",
		Usage:  "Repair a corrupted snapshot",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.ReexecFlag,
		},
		Description: `This command sets new list signer of snapshot from contract.`,
	}
)

// setupChain initializes the blockchain and its components based on the provided CLI context.
// It sets up the chain database, genesis block, consensus mechanism, and blockchain instance.
//
// Parameters:
// - ctx: The CLI context containing the command-line arguments and flags.
//
// Returns:
// - ethdb.Database: The initialized chain database.
// - tomoConfig: The configuration for the TomoChain node.
// - *params.ChainConfig: The chain configuration parameters.
// - *core.BlockChain: The initialized blockchain instance.
// - error: An error if the setup process fails.
func setupChain(ctx *cli.Context) (ethdb.Database, tomoConfig, *params.ChainConfig, *core.BlockChain, error) {
	stack, nodeConfig := makeConfigNode(ctx)
	chainDB := utils.MakeChainDatabase(ctx, stack)

	chainConfig, _, _ := core.SetupGenesisBlock(chainDB, nodeConfig.Eth.Genesis)

	var (
		vmConfig       = vm.Config{}
		cacheConfig    = &core.CacheConfig{}
		tomoXService   *tomox.TomoX
		lendingService *tomoxlending.Lending
	)

	posvConsensus := posv.New(chainConfig.Posv, chainDB)
	posvConsensus.GetTomoXService = func() posv.TradingService {
		return tomoXService
	}
	posvConsensus.GetLendingService = func() posv.LendingService {
		return lendingService
	}

	blockchain, err := core.NewBlockChain(chainDB, cacheConfig, chainConfig, posvConsensus, vmConfig)
	if err != nil {
		return nil, tomoConfig{}, nil, nil, err
	}

	return chainDB, nodeConfig, chainConfig, blockchain, nil
}

// reexecState re-executes blocks to find the latest valid state.
//
// Parameters:
// - reexec: The number of blocks to re-execute.
// - db: The database containing the blockchain data.
// - nodeConfig: The chain configuration containing the POSV parameters.
// - block: The current block from which to start re-executing.
// - blockchain: The blockchain instance.
//
// Returns:
// - *state.StateDB: The state database after re-executing the blocks.
// - *tradingstate.TradingStateDB: The trading state database after re-executing the blocks.
// - state.Database: The state database instance.
// - error: An error if the re-execution process fails.
func reexecState(
	reexec uint64,
	db ethdb.Database,
	nodeConfig tomoConfig,
	block *types.Block,
	blockchain *core.BlockChain,
) (*state.StateDB, *tradingstate.TradingStateDB, state.Database, error) {
	database := state.NewDatabase(db)

	var (
		stateDB      *state.StateDB
		tomoxStateDB *tradingstate.TradingStateDB
		err          error
	)

	// Re-execute to find the latest valid state
	for i := uint64(0); i < reexec; i++ {
		block = blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		if block == nil {
			break
		}

		if stateDB, err = state.New(block.Root(), database); err == nil {
			if block.NumberU64() >= common.TIPTomoXBlock.Uint64() {
				tomoxDB := tradingstate.NewDatabase(tomox.NewLDBEngine(&nodeConfig.TomoX))
				tomoxStateDB, err = tradingstate.New(block.Root(), tomoxDB)
			}
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		if _, ok := err.(*trie.MissingNodeError); ok {
			return nil, nil, nil, errors.New("required historical state unavailable")
		}
		return nil, nil, nil, err
	}
	return stateDB, tomoxStateDB, database, nil
}

// regenerateState regenerates the state from the given block to the target block.
//
// Parameters:
// - database: The state database instance.
// - block: The current block from which to start regenerating the state.
// - targetBlock: The target block to which the state needs to be regenerated.
// - stateDB: The state database to be updated during regeneration.
// - blockchain: The blockchain instance.
// - tomoxStateDB: The trading state database.
//
// Returns:
// - *state.StateDB: The updated state database after regeneration.
// - *tradingstate.TradingStateDB: The updated trading state database after regeneration.
// - error: An error if the regeneration process fails.
func regenerateState(
	database state.Database,
	block *types.Block,
	targetBlock *types.Block,
	stateDB *state.StateDB,
	blockchain *core.BlockChain,
	tomoxStateDB *tradingstate.TradingStateDB,
) (*state.StateDB, *tradingstate.TradingStateDB, error) {
	var (
		start  = time.Now()
		logged time.Time
		proot  common.Hash
	)
	for block.NumberU64() < targetBlock.NumberU64() {
		if time.Since(logged) > 8*time.Second {
			log.Info("Regenerating historical state", "block", block.NumberU64()+1, "target", targetBlock.NumberU64(), "elapsed", time.Since(start))
			logged = time.Now()
		}
		// Retrieve the next block to regenerate and process it
		if block = blockchain.GetBlockByNumber(block.NumberU64() + 1); block == nil {
			return nil, nil, fmt.Errorf("block #%d not found", block.NumberU64()+1)
		}
		feeCapacity := state.GetTRC21FeeCapacityFromState(stateDB)
		_, _, _, err := blockchain.Processor().Process(block, stateDB, tomoxStateDB, vm.Config{}, feeCapacity)
		if err != nil {
			return nil, nil, err
		}
		root := stateDB.IntermediateRoot(true)
		if root != block.Root() {
			return nil, nil, fmt.Errorf("invalid merkle root (number :%d  got : %x expect: %x)", block.NumberU64(), root.Hex(), block.Root())
		}
		// Finalize the state so any modifications are written to the trie
		root, err = stateDB.Commit(true)
		if err != nil {
			return nil, nil, err
		}
		if err := stateDB.Reset(root); err != nil {
			return nil, nil, err
		}
		database.TrieDB().Reference(root, common.Hash{})
		if proot != (common.Hash{}) {
			database.TrieDB().Dereference(proot)
		}
		proot = root
	}
	size, _ := database.TrieDB().Size()
	log.Info("Historical state regenerated", "block", block.NumberU64(), "elapsed", time.Since(start), "size", size)
	return stateDB, tomoxStateDB, nil
}

// finaliseBlock processes the transactions in the given block (usually gap block) and finalizes the state.
//
// Parameters:
// - targetBlock: The target block containing the transactions to be processed.
// - stateDB: The state database to be updated with the transaction results.
// - blockchain: The blockchain instance.
// - chainConfig: The chain configuration parameters.
// - tomoxStateDB: The trading state database.
//
// Returns:
// - *state.StateDB: The updated state database after processing the transactions.
// - error: An error if the transaction processing or state finalization fails.
func finaliseBlock(
	targetBlock *types.Block,
	stateDB *state.StateDB,
	blockchain *core.BlockChain,
	chainConfig *params.ChainConfig,
	tomoxStateDB *tradingstate.TradingStateDB,
) (*state.StateDB, error) {
	var (
		signer = types.MakeSigner(chainConfig, targetBlock.Number())
		txs    = targetBlock.Transactions()
	)

	feeCapacity := state.GetTRC21FeeCapacityFromState(stateDB)
	for _, tx := range txs {
		var balance *big.Int
		if tx.To() != nil {
			if value, ok := feeCapacity[*tx.To()]; ok {
				balance = value
			}
		}
		msg, _ := tx.AsMessage(signer, balance, targetBlock.Number(), false, chainConfig.IsAtlas(targetBlock.Number()))
		vmctx := core.NewEVMContext(msg, targetBlock.Header(), blockchain, nil)
		vmenv := vm.NewEVM(vmctx, stateDB, tomoxStateDB, chainConfig, vm.Config{})
		owner := common.Address{}
		if _, _, _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()), owner); err != nil {
			return nil, err
		}
		stateDB.Finalise(true)
	}
	return stateDB, nil
}

// initializeState initializes the state database and trading state database by either
// creating a new state from the block root or re-executing blocks to find the latest valid state.
//
// Parameters:
// - reexec: The number of blocks to re-execute.
// - db: The database containing the blockchain data.
// - nodeConfig: The chain configuration containing the POSV parameters.
// - block: The current block from which to start re-executing.
// - targetBlock: The target block to which the state needs to be regenerated.
// - blockchain: The blockchain instance.
//
// Returns:
// - *state.StateDB: The state database after initialization or re-execution.
// - error: An error if the initialization or re-execution process fails.
func initializeState(
	reexec uint64,
	db ethdb.Database,
	nodeConfig tomoConfig,
	block *types.Block,
	targetBlock *types.Block,
	blockchain *core.BlockChain,
) (*state.StateDB, error) {
	// Attempt to create a new state database from the block root.
	stateDB, err := state.New(block.Root(), state.NewDatabase(db))
	if err != nil {
		// If creating a new state database fails, re-execute blocks to find the latest valid state.
		stateDB, tomoxStateDB, database, err := reexecState(reexec, db, nodeConfig, block, blockchain)
		if err != nil {
			return nil, err
		}
		// Regenerate the state up to the gap block.
		stateDB, tomoxStateDB, err = regenerateState(database, block, targetBlock, stateDB, blockchain, tomoxStateDB)
		if err != nil {
			return nil, err
		}
		return finaliseBlock(block, stateDB, blockchain, blockchain.Config(), tomoxStateDB)
	}
	// Return the newly created state database.
	return stateDB, nil
}

// getNearestGap calculates the nearest gap block number and its hash based on the current head block.
// It uses the provided database and chain configuration to determine the nearest gap blocks.
//
// Parameters:
// - db: The database containing the blockchain data.
// - nodeConfig: The chain configuration containing the POSV parameters.
//
// Returns:
// - blockNumber: The nearest gap block number.
// - blockHash: The hash of the nearest gap block.
// - err: An error if sync block not checkpoint.
func getNearestGap(db ethdb.Database, nodeConfig *params.ChainConfig) (uint64, common.Hash, error) {
	// Retrieve the current head block's hash and number.
	headHash := core.GetHeadHeaderHash(db)
	headBlockNumber := core.GetBlockNumber(db, headHash)

	modulo := headBlockNumber % nodeConfig.Posv.Epoch
	gapThreshold := nodeConfig.Posv.Epoch - nodeConfig.Posv.Gap

	var nearestGapBlockNumber uint64

	switch {
	case modulo > gapThreshold:
		nearestGapBlockNumber = headBlockNumber - (modulo % gapThreshold)
	case modulo < gapThreshold:
		nearestGapBlockNumber = headBlockNumber - modulo - nodeConfig.Posv.Gap
	default:
		nearestGapBlockNumber = headBlockNumber
	}

	// Retrieve and return the canonical hash of the nearest gap block.
	gapBlockHash := core.GetCanonicalHash(db, nearestGapBlockNumber)
	return nearestGapBlockNumber, gapBlockHash, nil
}

// updateSnapshot updates the snapshot of the blockchain with new candidate signers.
//
// Parameters:
// - db: The database containing the blockchain data.
// - blockHash: The hash of the block for which the snapshot is being updated.
// - newCandidates: A map of new candidate addresses to be included in the snapshot.
//
// Returns:
// - error: An error if the update fails.
func updateSnapshot(db ethdb.Database, blockHash common.Hash, newCandidates map[common.Address]struct{}) error {
	// Retrieve the snapshot from the database.
	blob, err := db.Get(append([]byte("posv-"), blockHash[:]...))
	if err != nil {
		return err
	}
	snap := new(posv.Snapshot)
	// Unmarshal the snapshot from the database.
	if err := json.Unmarshal(blob, snap); err != nil {
		return err
	}

	// Log the old masternodes list.
	log.Info("Old masternodes list")
	for m := range snap.Signers {
		log.Info("Masternode", "Address", m.Hex())
	}

	// Update the masternodes list with the new candidates.
	snap.Signers = newCandidates

	blobWrite, err := json.Marshal(snap)
	if err != nil {
		return err
	}

	// Write the updated snapshot back to the database.
	err = db.Put(append([]byte("posv-"), snap.Hash[:]...), blobWrite)
	if err != nil {
		return err
	}

	// Log the successful repair of the snapshot.
	log.Info("Snapshot repaired successfully at block", "number", snap.Number, "hash", snap.Hash)
	return nil
}

// dbRepairSnapshot repairs the snapshot of the blockchain by updating the list of masternodes.
// It retrieves the nearest gap block, fetches candidate addresses, and updates their stakes.
//
// Parameters:
// - ctx: The CLI context containing the command-line arguments and flags.
//
// Returns:
// - error: An error if the repair process fails.
func dbRepairSnapshot(ctx *cli.Context) error {
	// Check if flag is provided
	dataDir := ctx.GlobalIsSet(utils.DataDirFlag.Name)
	if !dataDir {
		return fmt.Errorf("missing required flag: --datadir")
	}

	reexec := uint64(ctx.GlobalInt(utils.ReexecFlag.Name))

	chainDB, nodeConfig, chainConfig, blockchain, err := setupChain(ctx)
	if err != nil {
		return err
	}

	gapBlockNumber, gapBlockHash, err := getNearestGap(chainDB, chainConfig)
	if err != nil {
		return err
	}

	// Retrieve the gap block.
	gapBlock := core.GetBlock(chainDB, gapBlockHash, gapBlockNumber)

	var (
		candidateAddresses []common.Address
		masternodes        []posv.Masternode
		stateDB            *state.StateDB
		block              *types.Block
	)

	// Create a new state database from the gap block.
	block = gapBlock
	stateDB, err = initializeState(reexec, chainDB, nodeConfig, block, gapBlock, blockchain)
	if err != nil {
		return err
	}

	// Retrieve the candidate addresses from the state database.
	candidateAddresses = state.GetCandidates(stateDB)

	if len(candidateAddresses) == 0 {
		return fmt.Errorf("no candidates found")
	}

	// Iterate over the candidate addresses and retrieve their stakes.
	for _, candidate := range candidateAddresses {
		candidateStake := state.GetCandidateCap(stateDB, candidate)
		// Append the candidate to the list if it is valid.
		if candidate.Cmp(common.Address{}) != 0 {
			masternodes = append(masternodes, posv.Masternode{
				Address: candidate,
				Stake:   candidateStake,
			})
		}
	}

	// Sort the masternodes list by stake.
	sort.Slice(masternodes, func(i, j int) bool {
		return masternodes[i].Stake.Cmp(masternodes[j].Stake) >= 0
	})

	if len(masternodes) > common.MaxMasternodes {
		masternodes = masternodes[:common.MaxMasternodes]
	}

	newCandidates := make(map[common.Address]struct{})
	log.Info("New masternodes list")
	for _, masternode := range masternodes {
		log.Info("Masternode", "Address", masternode.Address.Hex())
		newCandidates[masternode.Address] = struct{}{}
	}

	return updateSnapshot(chainDB, gapBlockHash, newCandidates)
}
