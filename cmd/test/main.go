package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tomochain/tomochain/cmd/utils"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/consensus/posv"
	"github.com/tomochain/tomochain/contracts"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/eth"
	"github.com/tomochain/tomochain/eth/tracers"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/internal/ethapi"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/trie"
	"math/big"
	"os"
	"time"
)

var (
	tx                  = common.HexToHash("0x1ad76c9a8b8db44d68aacc9063a5962c56e8fe802aefcd7f68136162197b24be")
)

func main() {
	lddb, _ := ethdb.NewLDBDatabase("/home/tamnb/_projects/tomochain/src/github.com/tomochain/backup/chaindata", eth.DefaultConfig.DatabaseCache, utils.MakeDatabaseHandles())
	tx, blockHash, _, index := core.GetTransaction(lddb, tx)
	if tx == nil {
		fmt.Println("transaction not found", tx.Hash().Hex())
		return
	}
	chainConfig, _, err := core.SetupGenesisBlock(lddb, nil)
	if err != nil {
		fmt.Println("Get Chain Config fail ", err)
		return
	}
	engine := posv.New(chainConfig.Posv, lddb)
	engine.HookReward = func(chain consensus.ChainReader, stateBlock *state.StateDB, parentState *state.StateDB, header *types.Header) (error, map[string]interface{}) {
		number := header.Number.Uint64()
		rCheckpoint := chain.Config().Posv.RewardCheckpoint
		foundationWalletAddr := chain.Config().Posv.FoudationWalletAddr
		if foundationWalletAddr == (common.Address{}) {
			log.Error("Foundation Wallet Address is empty", "error", foundationWalletAddr)
			return errors.New("Foundation Wallet Address is empty"), nil
		}
		rewards := make(map[string]interface{})
		if number > 0 && number-rCheckpoint > 0 && foundationWalletAddr != (common.Address{}) {
			start := time.Now()
			// Get signers in blockSigner smartcontract.
			// Get reward inflation.
			chainReward := new(big.Int).Mul(new(big.Int).SetUint64(chain.Config().Posv.Reward), new(big.Int).SetUint64(params.Ether))
			chainReward = rewardInflation(chainReward, number, common.BlocksPerYear)

			totalSigner := new(uint64)
			signers, err := contracts.GetRewardForCheckpoint(engine, chain, header, rCheckpoint, totalSigner)

			log.Debug("Time Get Signers", "block", header.Number.Uint64(), "time", common.PrettyDuration(time.Since(start)))
			if err != nil {
				log.Crit("Fail to get signers for reward checkpoint", "error", err)
			}
			rewards["signers"] = signers
			rewardSigners, err := contracts.CalculateRewardForSigner(chainReward, signers, *totalSigner)
			if err != nil {
				log.Crit("Fail to calculate reward for signers", "error", err)
			}
			// Add reward for coin holders.
			voterResults := make(map[common.Address]interface{})
			if len(signers) > 0 {
				for signer, calcReward := range rewardSigners {
					err, rewards := contracts.CalculateRewardForHolders(foundationWalletAddr, parentState, signer, calcReward, number)
					if err != nil {
						log.Crit("Fail to calculate reward for holders.", "error", err)
					}
					if len(rewards) > 0 {
						for holder, reward := range rewards {
							stateBlock.AddBalance(holder, reward)
						}
					}
					voterResults[signer] = rewards
				}
			}
			rewards["rewards"] = voterResults
			log.Debug("Time Calculated HookReward ", "block", header.Number.Uint64(), "time", common.PrettyDuration(time.Since(start)))
		}
		return nil, rewards
	}
	blockchain, err := core.NewBlockChain(lddb, &core.CacheConfig{Disabled: false, TrieNodeLimit: eth.DefaultConfig.TrieCache, TrieTimeLimit: eth.DefaultConfig.TrieTimeout}, chainConfig, engine, vm.Config{EnablePreimageRecording: false})

	msg, vmctx, statedb, err := computeTxEnv(blockchain, lddb, chainConfig, blockHash, int(index))
	if err != nil {
		fmt.Println("computeTxEnv ", err)
		return
	}
	fmt.Println("test")
	f, err := os.Create("out.txt")
	defer f.Close()
	tracer:="callTracer"
	result, err := traceInternalTx(chainConfig, msg, vmctx, statedb, &eth.TraceConfig{Tracer:&tracer})
	if err != nil {
		fmt.Println("traceInternalTx", err)
	}
	if data,err := result.(json.RawMessage); !err{
		fmt.Println("RawMessage",len(data))
		f.Write(data)
	} else {
		data, err := json.Marshal(result)
		fmt.Println("Marshal",len(data))
		if err != nil {
			fmt.Println("Marshal", err)
		}
		f.Write(data)
	}
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func traceInternalTx(chainConfig *params.ChainConfig, message core.Message, vmctx vm.Context, statedb *state.StateDB, config *eth.TraceConfig) (interface{}, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer vm.Tracer
		err    error
	)
	switch {
	case config != nil && config.Tracer != nil:
		// Constuct the JavaScript tracer to execute with
		if tracer, err = tracers.New(*config.Tracer); err != nil {
			return nil, err
		}

	case config == nil:
		tracer = vm.NewStructLogger(nil)

	default:
		tracer = vm.NewStructLogger(config.LogConfig)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, statedb, chainConfig, vm.Config{Debug: true, Tracer: tracer})

	owner := common.Address{}
	ret, gas, failed, err := core.ApplyMessage(vmenv, message, new(core.GasPool).AddGas(message.Gas()), owner)
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %v", err)
	}
	// Depending on the tracer type, format and return the output
	switch tracer := tracer.(type) {
	case *vm.StructLogger:
		return &ethapi.ExecutionResult{
			Gas:         gas,
			Failed:      failed,
			ReturnValue: fmt.Sprintf("%x", ret),
			StructLogs:  ethapi.FormatLogs(tracer.StructLogs()),
		}, nil

	case *tracers.Tracer:
		return tracer.GetResult()

	default:
		panic(fmt.Sprintf("bad tracer type %T", tracer))
	}
}

func rewardInflation(chainReward *big.Int, number uint64, blockPerYear uint64) *big.Int {
	if blockPerYear*2 <= number && number < blockPerYear*6 {
		chainReward.Div(chainReward, new(big.Int).SetUint64(2))
	}
	if blockPerYear*6 <= number {
		chainReward.Div(chainReward, new(big.Int).SetUint64(4))
	}

	return chainReward
}

// computeTxEnv returns the execution environment of a certain transaction.
func computeTxEnv(blockchain *core.BlockChain, ethdb ethdb.Database, chainConfig *params.ChainConfig, blockHash common.Hash, txIndex int) (core.Message, vm.Context, *state.StateDB, error) {
	// Create the parent state database
	block := blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("block %x not found", blockHash)
	}
	parent := blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("parent %x not found", block.ParentHash())
	}
	statedb, err := computeStateDB(blockchain, ethdb, parent)
	fmt.Println("computeStateDB", err)
	if err != nil {
		return nil, vm.Context{}, nil, err
	}
	statedb.Commit(true)
	blockchain.StateCache.TrieDB().Commit(parent.Root(),false)
	config := chainConfig
	header := block.Header()
	usedGas := new(uint64)
	// Recompute transactions up to the target index.
	signer := types.MakeSigner(chainConfig, block.Number())
	feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
	for idx, tx := range block.Transactions() {
		var balacne *big.Int
		if tx.To() != nil {
			if value, ok := feeCapacity[*tx.To()]; ok {
				balacne = value
			}
		}
		// Assemble the transaction call message and return if the requested offset
		msg, _ := tx.AsMessage(signer, balacne, block.Number())
		context := core.NewEVMContext(msg, block.Header(), blockchain, nil)
		if idx == txIndex {
			return msg, context, statedb, nil
		}
		_, _, err, _ := core.ApplyTransaction(config, feeCapacity, blockchain, nil, new(core.GasPool), statedb, header, tx, usedGas, vm.Config{})
		if err != nil {
			return nil, vm.Context{}, nil, fmt.Errorf("tx %x failed: %v", tx.Hash(), err)
		}
		statedb.DeleteSuicides()
	}
	return nil, vm.Context{}, nil, fmt.Errorf("tx index %d out of range for block %x", txIndex, blockHash)
}

// computeStateDB retrieves the state database associated with a certain block.
// If no state is locally available for the given block, a number of blocks are
// attempted to be reexecuted to generate the desired state.
func computeStateDB(blockchain *core.BlockChain, chaindb ethdb.Database, block *types.Block) (*state.StateDB, error) {
	// If we have the state fully available, use that
	statedb, err := blockchain.StateAt(block.Root())
	if err == nil {
		return statedb, nil
	}
	// Otherwise try to reexec blocks until we find a state or reach our limit
	origin := block.NumberU64()
	database := state.NewDatabase(chaindb)

	for block.NumberU64() > 0 {
		block = blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		if block == nil {
			break
		}
		if statedb, err = state.New(block.Root(), database); err == nil {
			break
		}
	}
	fmt.Println("find state at block ", block.NumberU64())
	if err != nil {
		switch err.(type) {
		case *trie.MissingNodeError:
			return nil, errors.New("required historical state unavailable")
		default:
			return nil, err
		}
	}
	// State was available at historical point, regenerate
	var (
		start = time.Now()
		proot common.Hash
	)
	for block.NumberU64() < origin {
		// Retrieve the next block to regenerate and process it
		if block = blockchain.GetBlockByNumber(block.NumberU64() + 1); block == nil {
			return nil, fmt.Errorf("block #%d not found", block.NumberU64()+1)
		}
		feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
		_, _, index, err := blockchain.Processor().Process(block, statedb, vm.Config{}, feeCapacity)
		if err != nil {
			fmt.Println("blockchain.Processor().Process ", index, "err", err)
			return nil, err
		}
		// Finalize the state so any modifications are written to the trie
		root, err := statedb.Commit(true)
		if err != nil {
			return nil, err
		}
		if root != block.Root() {
			err = fmt.Errorf("invalid tomox merke trie got : %s , expect : %s ", root.Hex(), block.Root().Hex())
			fmt.Println("block number", block.NumberU64(), "hash", block.Hash().Hex(), "root", root.Hex())
			return nil, err
		}
		if err := statedb.Reset(root); err != nil {
			return nil, err
		}
		database.TrieDB().Reference(root, common.Hash{})
		database.TrieDB().Dereference(proot, common.Hash{})
		proot = root
	}
	log.Info("Historical state regenerated", "block", block.NumberU64(), "elapsed", time.Since(start), "size", database.TrieDB().Size())
	return statedb, nil
}
