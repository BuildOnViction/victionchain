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
	"fmt"
	"sort"

	"github.com/tomochain/tomochain/cmd/utils"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus/posv"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
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
		},
		Description: `This command sets new list signer of snapshot from contract.`,
	}
)

// getNearestGap calculates the nearest gap block number and its hash based on the current head block.
// It uses the provided database and chain configuration to determine the nearest gap blocks.
//
// Parameters:
// - db: The database containing the blockchain data.
// - config: The chain configuration containing the POSV parameters.
//
// Returns:
// - blockNumber: The nearest gap block number.
// - blockHash: The hash of the nearest gap block.
// - err: An error if sync block not checkpoint.
func getNearestGap(db ethdb.Database, config *params.ChainConfig) (blockNumber uint64, blockHash common.Hash, err error) {
	// Get the hash of the current head block.
	headHash := core.GetHeadHeaderHash(db)
	// Get the block number of the current head block.
	headBlockNumber := core.GetBlockNumber(db, headHash)
	// Get current sync block number.
	syncBlockNumber := headBlockNumber + 1

	if syncBlockNumber%config.Posv.Epoch != 0 {
		return 0, common.Hash{}, fmt.Errorf("mismatched signer only appears at a checkpoint")
	}
	// Calculate the nearest gap block number.
	nearestGapBlockNumber := syncBlockNumber - config.Posv.Gap
	// Get the hash of the nearest gap block.
	gapBlockHash := core.GetCanonicalHash(db, nearestGapBlockNumber)
	// Return the nearest gap block number and its hash.
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
	log.Info("Old masternodes list in order")
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
	// Check if --datadir is provided
	dataDir := ctx.GlobalIsSet(utils.DataDirFlag.Name)
	if !dataDir {
		return fmt.Errorf("missing required flag: --datadir")
	}

	// Create the chain database and configuration.
	stack, config := makeConfigNode(ctx)
	// Open the chain database.
	chainDB := utils.MakeChainDatabase(ctx, stack)
	defer chainDB.Close()

	// Setup the genesis block and retrieve the chain configuration.
	chainConfig, _, _ := core.SetupGenesisBlock(chainDB, config.Eth.Genesis)
	log.Info("Current chain configuration", "epoch", chainConfig.Posv.Epoch, "gap", chainConfig.Posv.Gap)

	gapBlockNumber, gapBlockHash, err := getNearestGap(chainDB, chainConfig)
	if err != nil {
		return err
	}

	// Retrieve the gap block.
	gapBlock := core.GetBlock(chainDB, gapBlockHash, gapBlockNumber)

	var (
		candidateAddresses []common.Address
		masternodes        []posv.Masternode
	)

	// Create a new state database from the gap block.
	stateDB, err := state.New(gapBlock.Root(), state.NewDatabase(chainDB))
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
	log.Info("New masternodes list ordered by stake")
	for _, masternode := range masternodes {
		log.Info("Masternode", "Address", masternode.Address.Hex(), "Stake:", masternode.Stake.String())
		newCandidates[masternode.Address] = struct{}{}
	}

	return updateSnapshot(chainDB, gapBlockHash, newCandidates)
}
