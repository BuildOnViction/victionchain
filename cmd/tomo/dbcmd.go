// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/accounts/abi"
	"github.com/tomochain/tomochain/cmd/utils"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus/posv"
	"github.com/tomochain/tomochain/contracts/validator/contract"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
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

// callmsg implements core.Message to allow passing it as a transaction simulator.
// Copied from accounts/abi/bind/backends/simulated.go
type callmsg struct {
	tomochain.CallMsg
}

func (m callmsg) From() common.Address      { return m.CallMsg.From }
func (m callmsg) Nonce() uint64             { return 0 }
func (m callmsg) CheckNonce() bool          { return false }
func (m callmsg) To() *common.Address       { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int        { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64               { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int           { return m.CallMsg.Value }
func (m callmsg) Data() []byte              { return m.CallMsg.Data }
func (m callmsg) BalanceTokenFee() *big.Int { return m.CallMsg.BalanceTokenFee }

// getNearestGap calculates the nearest gap block number and its hash based on the current head block.
// It uses the provided database and chain configuration to determine the nearest epoch and gap blocks.
//
// Parameters:
// - db: The database containing the blockchain data.
// - config: The chain configuration containing the POSV parameters.
//
// Returns:
// - blockNumber: The nearest gap block number.
// - blockHash: The hash of the nearest gap block.
// - err: An error if the nearest epoch is genesis.
func getNearestGap(db ethdb.Database, config *params.ChainConfig) (blockNumber uint64, blockHash common.Hash, err error) {
	// Get the hash of the current head block.
	headHash := core.GetHeadHeaderHash(db)
	// Get the block number of the current head block.
	headBlockNumber := core.GetBlockNumber(db, headHash)
	// Calculate the nearest epoch block number.
	nearestEpochBlockNumber := headBlockNumber - (headBlockNumber % config.Posv.Epoch)
	// If the nearest epoch block is the genesis block, return an error.
	if nearestEpochBlockNumber == 0 {
		return 0, common.Hash{}, fmt.Errorf("got genesis block")
	}
	// Calculate the nearest gap block number.
	nearestGapBlockNumber := nearestEpochBlockNumber - config.Posv.Gap
	// Get the hash of the nearest gap block.
	gapBlockHash := core.GetCanonicalHash(db, nearestGapBlockNumber)
	// Return the nearest gap block number and its hash.
	return nearestGapBlockNumber, gapBlockHash, nil
}

// getCandidateCap retrieves the candidate's stake from the smart contract.
//
// Parameters:
// - stateDB: The copy state database containing the blockchain state.
// - header: The block header containing the current block information.
// - config: The chain configuration containing the POSV parameters.
// - candidate: The address of the candidate whose stake is being retrieved.
// - contractAddress: The address of the smart contract.
// - abiContract: The ABI of the smart contract.
//
// Returns:
// - *big.Int: The stake amount of the candidate.
// - error: An error if the retrieval or unpacking fails.
func getCandidateCap(stateDB *state.StateDB, header *types.Header, config *params.ChainConfig, candidate common.Address, contractAddress common.Address, abiContract abi.ABI) (*big.Int, error) {
	// Define the method name to be called on the smart contract.
	method := "getCandidateCap"

	// Pack the input parameters for the smart contract call.
	input, err := abiContract.Pack(method, candidate)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("failed to pack input: %v", err)
	}

	// Define a fake caller address.
	fakeCaller := common.HexToAddress("0x0000000000000000000000000000000000000001")
	// Retrieve the TRC21 fee capacity from the state database.
	feeCapacity := state.GetTRC21FeeCapacityFromState(stateDB)

	// Create the call message for the smart contract call.
	callMsg := tomochain.CallMsg{
		To:       &contractAddress,
		Data:     input,
		From:     fakeCaller,
		GasPrice: big.NewInt(0),
		Gas:      1000000,
		Value:    new(big.Int),
	}

	// Set the balance token fee for the call message.
	callMsg.BalanceTokenFee = feeCapacity[*callMsg.To]

	// Wrap the call message in a new callmsg type.
	newCallMsg := callmsg{callMsg}

	// Create the EVM context for the smart contract call.
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     nil,
		Origin:      newCallMsg.From(),
		Coinbase:    header.Coinbase,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		GasLimit:    header.GasLimit,
		GasPrice:    new(big.Int).Set(newCallMsg.GasPrice()),
	}

	// Create a new EVM instance with the provided context, state database, and configuration.
	evm := vm.NewEVM(evmContext, stateDB, nil, config, vm.Config{})

	// Create a new gas pool and add the gas from the call message.
	gasPool := new(core.GasPool).AddGas(callMsg.Gas)
	owner := common.Address{}
	// Execute the state transition.
	returnValue, _, _, err := core.NewStateTransition(evm, newCallMsg, gasPool).TransitionDb(owner)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("state transition error: %v", err)
	}

	// Unpack the result from the call.
	var candidateStake *big.Int
	if err := abiContract.Unpack(&candidateStake, method, returnValue); err != nil {
		return big.NewInt(0), fmt.Errorf("failed to unpack result: %v", err)
	}
	return candidateStake, nil
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

	// Retrieve the gap block and its header.
	gapBlock := core.GetBlock(chainDB, gapBlockHash, gapBlockNumber)
	gapHeader := core.GetHeader(chainDB, gapBlockHash, gapBlockNumber)

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

	contractAddress := common.HexToAddress(common.MasternodeVotingSMC)
	abiContract, err := abi.JSON(strings.NewReader(contract.TomoValidatorABI))
	if err != nil {
		return err
	}

	// Create a copy of the state database to avoid modifying the original.
	stateDBCopy := stateDB.Copy()

	// Iterate over the candidate addresses and retrieve their stakes.
	for _, candidate := range candidateAddresses {
		candidateStake, err := getCandidateCap(stateDBCopy, gapHeader, chainConfig, candidate, contractAddress, abiContract)
		if err != nil {
			return err
		}
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

	newCandidates := make(map[common.Address]struct{})
	log.Info("New masternodes list ordered by stake")
	for _, masternode := range masternodes {
		log.Info("Masternode", "Address", masternode.Address.Hex(), "Stake:", masternode.Stake.String())
		newCandidates[masternode.Address] = struct{}{}
	}

	return updateSnapshot(chainDB, gapBlockHash, newCandidates)
}
