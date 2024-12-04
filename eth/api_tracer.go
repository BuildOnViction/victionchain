// Copyright 2017 The go-ethereum Authors
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

package eth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/tomochain/tomochain/tomox/tradingstate"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/eth/tracers"
	"github.com/tomochain/tomochain/internal/ethapi"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/rpc"
	"github.com/tomochain/tomochain/trie"
)

const (
	// defaultTraceTimeout is the amount of time a single transaction can execute
	// by default before being forcefully aborted.
	defaultTraceTimeout = 5 * time.Second

	// defaultTraceReexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultTraceReexec = uint64(128)
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*vm.LogConfig
	Tracer  *string
	Timeout *string
	Reexec  *uint64
}

// txTraceResult is the result of a single transaction trace.
type txTraceResult struct {
	Result interface{} `json:"result,omitempty"` // Trace results produced by the tracer
	Error  string      `json:"error,omitempty"`  // Trace failure produced by the tracer
}

// blockTraceTask represents a single block trace task when an entire chain is
// being traced.
type blockTraceTask struct {
	statedb *state.StateDB   // Intermediate state prepped for tracing
	block   *types.Block     // Block to trace the transactions from
	rootref common.Hash      // Trie root reference held for this task
	results []*txTraceResult // Trace results procudes by the task
}

// blockTraceResult represets the results of tracing a single block when an entire
// chain is being traced.
type blockTraceResult struct {
	Block  hexutil.Uint64   `json:"block"`  // Block number corresponding to this trace
	Hash   common.Hash      `json:"hash"`   // Block hash corresponding to this trace
	Traces []*txTraceResult `json:"traces"` // Trace results produced by the task
}

// txTraceTask represents a single transaction trace task when an entire block
// is being traced.
type txTraceTask struct {
	statedb *state.StateDB // Intermediate state prepped for tracing
	index   int            // Transaction offset in the block
}

// TraceChain returns the structured logs created during the execution of EVM
// between two blocks (excluding start) and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceChain(ctx context.Context, start, end rpc.BlockNumber, config *TraceConfig) (*rpc.Subscription, error) {
	// Fetch the block interval that we want to trace
	var from, to *types.Block

	switch start {
	case rpc.PendingBlockNumber:
		from = api.eth.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		from = api.eth.blockchain.CurrentBlock()
	default:
		from = api.eth.blockchain.GetBlockByNumber(uint64(start))
	}
	switch end {
	case rpc.PendingBlockNumber:
		to = api.eth.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		to = api.eth.blockchain.CurrentBlock()
	default:
		to = api.eth.blockchain.GetBlockByNumber(uint64(end))
	}
	// Trace the chain if we've found all our blocks
	if from == nil {
		return nil, fmt.Errorf("starting block #%d not found", start)
	}
	if to == nil {
		return nil, fmt.Errorf("end block #%d not found", end)
	}
	return api.traceChain(ctx, from, to, config)
}

// traceChain configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *PrivateDebugAPI) traceChain(ctx context.Context, start, end *types.Block, config *TraceConfig) (*rpc.Subscription, error) {
	// Tracing a chain is a **long** operation, only do with subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	sub := notifier.CreateSubscription()

	// Ensure we have a valid starting state before doing any work
	origin := start.NumberU64()
	database := state.NewDatabase(api.eth.ChainDb())

	if number := start.NumberU64(); number > 0 {
		start = api.eth.blockchain.GetBlock(start.ParentHash(), start.NumberU64()-1)
		if start == nil {
			return nil, fmt.Errorf("parent block #%d not found", number-1)
		}
	}
	statedb, err := state.New(start.Root(), database)
	var tomoxState *tradingstate.TradingStateDB
	if err != nil {
		// If the starting state is missing, allow some number of blocks to be reexecuted
		reexec := defaultTraceReexec
		if config != nil && config.Reexec != nil {
			reexec = *config.Reexec
		}
		// Find the most recent block that has the state available
		for i := uint64(0); i < reexec; i++ {
			start = api.eth.blockchain.GetBlock(start.ParentHash(), start.NumberU64()-1)
			if start == nil {
				break
			}
			if statedb, err = state.New(start.Root(), database); err == nil {
				tomoxState, err = tradingstate.New(start.Root(), tradingstate.NewDatabase(api.eth.TomoX.GetLevelDB()))
				if err == nil {
					break
				}
			}
		}
		// If we still don't have the state available, bail out
		if err != nil {
			switch err.(type) {
			case *trie.MissingNodeError:
				return nil, errors.New("required historical state unavailable")
			default:
				return nil, err
			}
		}
	}
	// Execute all the transaction contained within the chain concurrently for each block
	blocks := int(end.NumberU64() - origin)

	threads := runtime.NumCPU()
	if threads > blocks {
		threads = blocks
	}
	var (
		pend    = new(sync.WaitGroup)
		tasks   = make(chan *blockTraceTask, threads)
		results = make(chan *blockTraceTask, threads)
	)
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next block trace tasks
			for task := range tasks {
				signer := types.MakeSigner(api.config, task.block.Number())
				feeCapacity := state.GetTRC21FeeCapacityFromState(task.statedb)
				// Trace all the transactions contained within
				for i, tx := range task.block.Transactions() {
					var balance *big.Int
					if tx.To() != nil {
						if value, ok := feeCapacity[*tx.To()]; ok {
							balance = value
						}
					}
					msg, _ := tx.AsMessage(signer, balance, task.block.Number(), false)
					vmctx := core.NewEVMContext(msg, task.block.Header(), api.eth.blockchain, nil)

					res, err := api.traceTx(ctx, msg, vmctx, task.statedb, config)
					if err != nil {
						task.results[i] = &txTraceResult{Error: err.Error()}
						log.Warn("Tracing failed", "hash", tx.Hash(), "block", task.block.NumberU64(), "err", err)
						break
					}
					task.statedb.DeleteSuicides()
					task.results[i] = &txTraceResult{Result: res}
				}
				// Stream the result back to the user or abort on teardown
				select {
				case results <- task:
				case <-notifier.Closed():
					return
				}
			}
		}()
	}
	// Start a goroutine to feed all the blocks into the tracers
	begin := time.Now()

	go func() {
		var (
			logged time.Time
			number uint64
			traced uint64
			failed error
			proot  common.Hash
		)
		// Ensure everything is properly cleaned up on any exit path
		defer func() {
			close(tasks)
			pend.Wait()

			switch {
			case failed != nil:
				log.Warn("Chain tracing failed", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin), "err", failed)
			case number < end.NumberU64():
				log.Warn("Chain tracing aborted", "start", start.NumberU64(), "end", end.NumberU64(), "abort", number, "transactions", traced, "elapsed", time.Since(begin))
			default:
				log.Info("Chain tracing finished", "start", start.NumberU64(), "end", end.NumberU64(), "transactions", traced, "elapsed", time.Since(begin))
			}
			close(results)
		}()
		// Feed all the blocks both into the tracer, as well as fast process concurrently
		for number = start.NumberU64() + 1; number <= end.NumberU64(); number++ {
			// Stop tracing if interruption was requested
			select {
			case <-notifier.Closed():
				return
			default:
			}
			// Print progress logs if long enough time elapsed
			if time.Since(logged) > 8*time.Second {
				if number > origin {
					memory, _ := database.TrieDB().Size()
					log.Info("Tracing chain segment", "start", origin, "end", end.NumberU64(), "current", number, "transactions", traced, "elapsed", time.Since(begin), "memory", memory)
				} else {
					log.Info("Preparing state for chain trace", "block", number, "start", origin, "elapsed", time.Since(begin))
				}
				logged = time.Now()
			}
			// Retrieve the next block to trace
			block := api.eth.blockchain.GetBlockByNumber(number)
			if block == nil {
				failed = fmt.Errorf("block #%d not found", number)
				break
			}
			// Send the block over to the concurrent tracers (if not in the fast-forward phase)
			if number > origin {
				txs := block.Transactions()

				select {
				case tasks <- &blockTraceTask{statedb: statedb.Copy(), block: block, rootref: proot, results: make([]*txTraceResult, len(txs))}:
				case <-notifier.Closed():
					return
				}
				traced += uint64(len(txs))
			}
			feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
			// Generate the next state snapshot fast without tracing
			_, _, _, err := api.eth.blockchain.Processor().Process(block, statedb, tomoxState, vm.Config{}, feeCapacity)
			if err != nil {
				failed = err
				break
			}
			// Finalize the state so any modifications are written to the trie
			root, err := statedb.Commit(true)
			if err != nil {
				failed = err
				break
			}
			if err := statedb.Reset(root); err != nil {
				failed = err
				break
			}
			// Reference the trie twice, once for us, once for the trancer
			database.TrieDB().Reference(root, common.Hash{})
			if number >= origin {
				database.TrieDB().Reference(root, common.Hash{})
			}
			// Dereference all past tries we ourselves are done working with
			database.TrieDB().Dereference(proot)
			proot = root
		}
	}()

	// Keep reading the trace results and stream the to the user
	go func() {
		var (
			done = make(map[uint64]*blockTraceResult)
			next = origin + 1
		)
		for res := range results {
			// Queue up next received result
			result := &blockTraceResult{
				Block:  hexutil.Uint64(res.block.NumberU64()),
				Hash:   res.block.Hash(),
				Traces: res.results,
			}
			done[uint64(result.Block)] = result

			// Dereference any paret tries held in memory by this task
			database.TrieDB().Dereference(res.rootref)

			// Stream completed traces to the user, aborting on the first error
			for result, ok := done[next]; ok; result, ok = done[next] {
				if len(result.Traces) > 0 || next == end.NumberU64() {
					notifier.Notify(sub.ID, result)
				}
				delete(done, next)
				next++
			}
		}
	}()
	return sub, nil
}

// TraceBlockByNumber returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByNumber(ctx context.Context, number rpc.BlockNumber, config *TraceConfig) ([]*txTraceResult, error) {
	// Fetch the block that we want to trace
	var block *types.Block

	switch number {
	case rpc.PendingBlockNumber:
		block = api.eth.miner.PendingBlock()
	case rpc.LatestBlockNumber:
		block = api.eth.blockchain.CurrentBlock()
	default:
		block = api.eth.blockchain.GetBlockByNumber(uint64(number))
	}
	// Trace the block if it was found
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", number)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockByHash returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockByHash(ctx context.Context, hash common.Hash, config *TraceConfig) ([]*txTraceResult, error) {
	block := api.eth.blockchain.GetBlockByHash(hash)
	if block == nil {
		return nil, fmt.Errorf("block #%x not found", hash)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlock returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlock(ctx context.Context, blob []byte, config *TraceConfig) ([]*txTraceResult, error) {
	block := new(types.Block)
	if err := rlp.Decode(bytes.NewReader(blob), block); err != nil {
		return nil, fmt.Errorf("could not decode block: %v", err)
	}
	return api.traceBlock(ctx, block, config)
}

// TraceBlockFromFile returns the structured logs created during the execution of
// EVM and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceBlockFromFile(ctx context.Context, file string, config *TraceConfig) ([]*txTraceResult, error) {
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	return api.TraceBlock(ctx, blob, config)
}

// traceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requestd tracer.
func (api *PrivateDebugAPI) traceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*txTraceResult, error) {
	// Create the parent state database
	parent := api.eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, fmt.Errorf("parent %x not found", block.ParentHash())
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, tomoxState, err := api.computeStateDB(parent, reexec)
	if err != nil {
		return nil, err
	}
	// Execute all the transaction contained within the block concurrently
	var (
		signer = types.MakeSigner(api.config, block.Number())

		txs     = block.Transactions()
		results = make([]*txTraceResult, len(txs))

		pend = new(sync.WaitGroup)
		jobs = make(chan *txTraceTask, len(txs))
	)
	threads := runtime.NumCPU()
	if threads > len(txs) {
		threads = len(txs)
	}
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()

			// Fetch and execute the next transaction trace tasks
			for task := range jobs {
				feeCapacity := state.GetTRC21FeeCapacityFromState(task.statedb)
				var balance *big.Int
				if txs[task.index].To() != nil {
					if value, ok := feeCapacity[*txs[task.index].To()]; ok {
						balance = value
					}
				}
				msg, _ := txs[task.index].AsMessage(signer, balance, block.Number(), false)
				vmctx := core.NewEVMContext(msg, block.Header(), api.eth.blockchain, nil)

				res, err := api.traceTx(ctx, msg, vmctx, task.statedb, config)
				if err != nil {
					results[task.index] = &txTraceResult{Error: err.Error()}
					continue
				}
				results[task.index] = &txTraceResult{Result: res}
			}
		}()
	}
	// Feed the transactions into the tracers and return
	feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
	var failed error
	for i, tx := range txs {
		// Send the trace task over for execution
		jobs <- &txTraceTask{statedb: statedb.Copy(), index: i}
		var balance *big.Int
		if tx.To() != nil {
			if value, ok := feeCapacity[*tx.To()]; ok {
				balance = value
			}
		}
		// Generate the next state snapshot fast without tracing
		msg, _ := tx.AsMessage(signer, balance, block.Number(), false)
		vmctx := core.NewEVMContext(msg, block.Header(), api.eth.blockchain, nil)

		vmenv := vm.NewEVM(vmctx, statedb, tomoxState, api.config, vm.Config{})
		owner := common.Address{}
		if _, _, _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()), owner); err != nil {
			failed = err
			break
		}

		// Finalize the state so any modifications are written to the trie
		statedb.Finalise(true)
	}
	close(jobs)
	pend.Wait()

	// If execution failed in between, abort
	if failed != nil {
		return nil, failed
	}
	return results, nil
}

// computeStateDB retrieves the state database associated with a certain block.
// If no state is locally available for the given block, a number of blocks are
// attempted to be reexecuted to generate the desired state.
func (api *PrivateDebugAPI) computeStateDB(block *types.Block, reexec uint64) (*state.StateDB, *tradingstate.TradingStateDB, error) {
	// If we have the state fully available, use that
	statedb, err := api.eth.blockchain.StateAt(block.Root())
	tomoxState := &tradingstate.TradingStateDB{}
	if err == nil {
		tomoxState, err = api.eth.blockchain.OrderStateAt(block)
		if err == nil {
			return statedb, tomoxState, nil
		}
	}
	// Otherwise try to reexec blocks until we find a state or reach our limit
	origin := block.NumberU64()
	database := state.NewDatabase(api.eth.ChainDb())

	for i := uint64(0); i < reexec; i++ {
		block = api.eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		if block == nil {
			break
		}
		if statedb, err = state.New(block.Root(), database); err == nil {
			tomoxState, err = api.eth.blockchain.OrderStateAt(block)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		switch err.(type) {
		case *trie.MissingNodeError:
			return nil, nil, errors.New("required historical state unavailable")
		default:
			return nil, nil, err
		}
	}
	// State was available at historical point, regenerate
	var (
		start  = time.Now()
		logged time.Time
		proot  common.Hash
	)
	for block.NumberU64() < origin {
		// Print progress logs if long enough time elapsed
		if time.Since(logged) > 8*time.Second {
			log.Info("Regenerating historical state", "block", block.NumberU64()+1, "target", origin, "elapsed", time.Since(start))
			logged = time.Now()
		}
		// Retrieve the next block to regenerate and process it
		if block = api.eth.blockchain.GetBlockByNumber(block.NumberU64() + 1); block == nil {
			return nil, nil, fmt.Errorf("block #%d not found", block.NumberU64()+1)
		}
		feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
		_, _, _, err := api.eth.blockchain.Processor().Process(block, statedb, tomoxState, vm.Config{}, feeCapacity)
		if err != nil {
			return nil, nil, err
		}
		root := statedb.IntermediateRoot(true)
		if root != block.Root() {
			return nil, nil, fmt.Errorf("invalid merkle root (number :%d  got : %x expect: %x)", block.NumberU64(), root.Hex(), block.Root())
		}
		// Finalize the state so any modifications are written to the trie
		root, err = statedb.Commit(true)
		if err != nil {
			return nil, nil, err
		}
		if err := statedb.Reset(root); err != nil {
			return nil, nil, err
		}
		database.TrieDB().Reference(root, common.Hash{})
		database.TrieDB().Dereference(proot)
		proot = root
	}
	size, _ := database.TrieDB().Size()
	log.Info("Historical state regenerated", "block", block.NumberU64(), "elapsed", time.Since(start), "size", size)
	return statedb, tomoxState, nil
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PrivateDebugAPI) TraceTransaction(ctx context.Context, hash common.Hash, config *TraceConfig) (interface{}, error) {
	// Retrieve the transaction and assemble its EVM context
	tx, blockHash, _, index := core.GetTransaction(api.eth.ChainDb(), hash)
	if tx == nil {
		return nil, fmt.Errorf("transaction %x not found", hash)
	}
	reexec := defaultTraceReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	msg, vmctx, statedb, err := api.computeTxEnv(blockHash, int(index), reexec)
	if err != nil {
		return nil, err
	}
	// Trace the transaction and return
	return api.traceTx(ctx, msg, vmctx, statedb, config)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be tracer dependent.
func (api *PrivateDebugAPI) traceTx(ctx context.Context, message core.Message, vmctx vm.Context, statedb *state.StateDB, config *TraceConfig) (interface{}, error) {
	// Assemble the structured logger or the JavaScript tracer
	var (
		tracer vm.Tracer
		err    error
	)
	switch {
	case config != nil && config.Tracer != nil:
		// Define a meaningful timeout of a single transaction trace
		timeout := defaultTraceTimeout
		if config.Timeout != nil {
			if timeout, err = time.ParseDuration(*config.Timeout); err != nil {
				return nil, err
			}
		}
		// Constuct the JavaScript tracer to execute with
		if tracer, err = tracers.New(*config.Tracer); err != nil {
			return nil, err
		}
		// Handle timeouts and RPC cancellations
		deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-deadlineCtx.Done()
			tracer.(*tracers.Tracer).Stop(errors.New("execution timeout"))
		}()
		defer cancel()

	case config == nil:
		tracer = vm.NewStructLogger(nil)

	default:
		tracer = vm.NewStructLogger(config.LogConfig)
	}
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, statedb, nil, api.config, vm.Config{Debug: true, Tracer: tracer})

	owner := common.Address{}
	
	// Bypass blacklist address
	maxBlockNumber := new(big.Int).SetInt64(9147459)
	if vmctx.BlockNumber.Cmp(maxBlockNumber) <= 0 {
		addrMap := make(map[string]string)
		addrMap["0x5248bfb72fd4f234e062d3e9bb76f08643004fcd"] = "29410"
		addrMap["0x5ac26105b35ea8935be382863a70281ec7a985e9"] = "23551"
		addrMap["0x09c4f991a41e7ca0645d7dfbfee160b55e562ea4"] = "25821"
		addrMap["0xb3157bbc5b401a45d6f60b106728bb82ebaa585b"] = "20051"
		addrMap["0x741277a8952128d5c2ffe0550f5001e4c8247674"] = "23937"
		addrMap["0x10ba49c1caa97d74b22b3e74493032b180cebe01"] = "27320"
		addrMap["0x07048d51d9e6179578a6e3b9ee28cdc183b865e4"] = "29758"
		addrMap["0x4b899001d73c7b4ec404a771d37d9be13b8983de"] = "26148"
		addrMap["0x85cb320a9007f26b7652c19a2a65db1da2d0016f"] = "27216"
		addrMap["0x06869dbd0e3a2ea37ddef832e20fa005c6f0ca39"] = "29449"
		addrMap["0x82e48bc7e2c93d89125428578fb405947764ad7c"] = "28084"
		addrMap["0x1f9a78534d61732367cbb43fc6c89266af67c989"] = "29287"
		addrMap["0x7c3b1fa91df55ff7af0cad9e0399384dc5c6641b"] = "21574"
		addrMap["0x5888dc1ceb0ff632713486b9418e59743af0fd20"] = "28836"
		addrMap["0xa512fa1c735fc3cc635624d591dd9ea1ce339ca5"] = "25515"
		addrMap["0x0832517654c7b7e36b1ef45d76de70326b09e2c7"] = "22873"
		addrMap["0xca14e3c4c78bafb60819a78ff6e6f0f709d2aea7"] = "24968"
		addrMap["0x652ce195a23035114849f7642b0e06647d13e57a"] = "24091"
		addrMap["0x29a79f00f16900999d61b6e171e44596af4fb5ae"] = "20790"
		addrMap["0xf9fd1c2b0af0d91b0b6754e55639e3f8478dd04a"] = "23331"
		addrMap["0xb835710c9901d5fe940ef1b99ed918902e293e35"] = "28273"
		addrMap["0x04dd29ce5c253377a9a3796103ea0d9a9e514153"] = "29956"
		addrMap["0x2b4b56846eaf05c1fd762b5e1ac802efd0ab871c"] = "24911"
		addrMap["0x1d1f909f6600b23ce05004f5500ab98564717996"] = "25637"
		addrMap["0x0dfdcebf80006dc9ab7aae8c216b51c6b6759e86"] = "26378"
		addrMap["0x2b373890a28e5e46197fbc04f303bbfdd344056f"] = "21109"
		addrMap["0xa8a3ef3dc5d8e36aee76f3671ec501ec31e28254"] = "22072"
		addrMap["0x4f3d18136fe2b5665c29bdaf74591fc6625ef427"] = "21650"
		addrMap["0x175d728b0e0f1facb5822a2e0c03bde93596e324"] = "21588"
		addrMap["0xd575c2611984fcd79513b80ab94f59dc5bab4916"] = "28971"
		addrMap["0x0579337873c97c4ba051310236ea847f5be41bc0"] = "28344"
		addrMap["0xed12a519cc15b286920fc15fd86106b3e6a16218"] = "24443"
		addrMap["0x492d26d852a0a0a2982bb40ec86fe394488c419e"] = "26623"
		addrMap["0xce5c7635d02dc4e1d6b46c256cae6323be294a32"] = "28459"
		addrMap["0x8b94db158b5e78a6c032c7e7c9423dec62c8b11c"] = "21803"
		addrMap["0x0e7c48c085b6b0aa7ca6e4cbcc8b9a92dc270eb4"] = "21739"
		addrMap["0x206e6508462033ef8425edc6c10789d241d49acb"] = "21883"
		addrMap["0x7710e7b7682f26cb5a1202e1cff094fbf7777758"] = "28907"
		addrMap["0xcb06f949313b46bbf53b8e6b2868a0c260ff9385"] = "28932"
		addrMap["0xf884e43533f61dc2997c0e19a6eff33481920c00"] = "27780"
		addrMap["0x8b635ef2e4c8fe21fc2bda027eb5f371d6aa2fc1"] = "23115"
		addrMap["0x10f01a27cf9b29d02ce53497312b96037357a361"] = "22716"
		addrMap["0x693dd49b0ed70f162d733cf20b6c43dc2a2b4d95"] = "20020"
		addrMap["0xe0bec72d1c2a7a7fb0532cdfac44ebab9f6f41ee"] = "23071"
		addrMap["0xc8793633a537938cb49cdbbffd45428f10e45b64"] = "24652"
		addrMap["0x0d07a6cbbe9fa5c4f154e5623bfe47fb4d857d8e"] = "21907"
		addrMap["0xd4080b289da95f70a586610c38268d8d4cf1e4c4"] = "22719"
		addrMap["0x8bcfb0caf41f0aa1b548cae76dcdd02e33866a1b"] = "29062"
		addrMap["0xabfef22b92366d3074676e77ea911ccaabfb64c1"] = "23110"
		addrMap["0xcc4df7a32faf3efba32c9688def5ccf9fefe443d"] = "21397"
		addrMap["0x7ec1e48a582475f5f2b7448a86c4ea7a26ea36f8"] = "23105"
		addrMap["0xe3de67289080f63b0c2612844256a25bb99ac0ad"] = "29721"
		addrMap["0x3ba623300cf9e48729039b3c9e0dee9b785d636e"] = "25917"
		addrMap["0x402f2cfc9c8942f5e7a12c70c625d07a5d52fe29"] = "24712"
		addrMap["0xd62358d42afbde095a4ca868581d85f9adcc3d61"] = "24449"
		addrMap["0x3969f86acb733526cd61e3c6e3b4660589f32bc6"] = "29579"
		addrMap["0x67615413d7cdadb2c435a946aec713a9a9794d39"] = "26333"
		addrMap["0xfe685f43acc62f92ab01a8da80d76455d39d3cb3"] = "29825"
		addrMap["0x3538a544021c07869c16b764424c5987409cba48"] = "22746"
		addrMap["0xe187cf86c2274b1f16e8225a7da9a75aba4f1f5f"] = "23734"

		blockMap := make(map[int64]string)

		blockMap[9073579] = "0x5248bfb72fd4f234e062d3e9bb76f08643004fcd"
		blockMap[9147130] = "0x5ac26105b35ea8935be382863a70281ec7a985e9"
		blockMap[9147195] = "0x09c4f991a41e7ca0645d7dfbfee160b55e562ea4"
		blockMap[9147200] = "0xb3157bbc5b401a45d6f60b106728bb82ebaa585b"
		blockMap[9147206] = "0x741277a8952128d5c2ffe0550f5001e4c8247674"
		blockMap[9147212] = "0x10ba49c1caa97d74b22b3e74493032b180cebe01"
		blockMap[9147217] = "0x07048d51d9e6179578a6e3b9ee28cdc183b865e4"
		blockMap[9147223] = "0x4b899001d73c7b4ec404a771d37d9be13b8983de"
		blockMap[9147229] = "0x85cb320a9007f26b7652c19a2a65db1da2d0016f"
		blockMap[9147234] = "0x06869dbd0e3a2ea37ddef832e20fa005c6f0ca39"
		blockMap[9147240] = "0x82e48bc7e2c93d89125428578fb405947764ad7c"
		blockMap[9147246] = "0x1f9a78534d61732367cbb43fc6c89266af67c989"
		blockMap[9147251] = "0x7c3b1fa91df55ff7af0cad9e0399384dc5c6641b"
		blockMap[9147257] = "0x5888dc1ceb0ff632713486b9418e59743af0fd20"
		blockMap[9147263] = "0xa512fa1c735fc3cc635624d591dd9ea1ce339ca5"
		blockMap[9147268] = "0x0832517654c7b7e36b1ef45d76de70326b09e2c7"
		blockMap[9147274] = "0xca14e3c4c78bafb60819a78ff6e6f0f709d2aea7"
		blockMap[9147279] = "0x652ce195a23035114849f7642b0e06647d13e57a"
		blockMap[9147285] = "0x29a79f00f16900999d61b6e171e44596af4fb5ae"
		blockMap[9147291] = "0xf9fd1c2b0af0d91b0b6754e55639e3f8478dd04a"
		blockMap[9147296] = "0xb835710c9901d5fe940ef1b99ed918902e293e35"
		blockMap[9147302] = "0x04dd29ce5c253377a9a3796103ea0d9a9e514153"
		blockMap[9147308] = "0x2b4b56846eaf05c1fd762b5e1ac802efd0ab871c"
		blockMap[9147314] = "0x1d1f909f6600b23ce05004f5500ab98564717996"
		blockMap[9147319] = "0x0dfdcebf80006dc9ab7aae8c216b51c6b6759e86"
		blockMap[9147325] = "0x2b373890a28e5e46197fbc04f303bbfdd344056f"
		blockMap[9147330] = "0xa8a3ef3dc5d8e36aee76f3671ec501ec31e28254"
		blockMap[9147336] = "0x4f3d18136fe2b5665c29bdaf74591fc6625ef427"
		blockMap[9147342] = "0x175d728b0e0f1facb5822a2e0c03bde93596e324"
		blockMap[9145281] = "0xd575c2611984fcd79513b80ab94f59dc5bab4916"
		blockMap[9145315] = "0x0579337873c97c4ba051310236ea847f5be41bc0"
		blockMap[9145341] = "0xed12a519cc15b286920fc15fd86106b3e6a16218"
		blockMap[9145367] = "0x492d26d852a0a0a2982bb40ec86fe394488c419e"
		blockMap[9145386] = "0xce5c7635d02dc4e1d6b46c256cae6323be294a32"
		blockMap[9145414] = "0x8b94db158b5e78a6c032c7e7c9423dec62c8b11c"
		blockMap[9145436] = "0x0e7c48c085b6b0aa7ca6e4cbcc8b9a92dc270eb4"
		blockMap[9145463] = "0x206e6508462033ef8425edc6c10789d241d49acb"
		blockMap[9145493] = "0x7710e7b7682f26cb5a1202e1cff094fbf7777758"
		blockMap[9145519] = "0xcb06f949313b46bbf53b8e6b2868a0c260ff9385"
		blockMap[9145549] = "0xf884e43533f61dc2997c0e19a6eff33481920c00"
		blockMap[9147352] = "0x8b635ef2e4c8fe21fc2bda027eb5f371d6aa2fc1"
		blockMap[9147357] = "0x10f01a27cf9b29d02ce53497312b96037357a361"
		blockMap[9147363] = "0x693dd49b0ed70f162d733cf20b6c43dc2a2b4d95"
		blockMap[9147369] = "0xe0bec72d1c2a7a7fb0532cdfac44ebab9f6f41ee"
		blockMap[9147375] = "0xc8793633a537938cb49cdbbffd45428f10e45b64"
		blockMap[9147380] = "0x0d07a6cbbe9fa5c4f154e5623bfe47fb4d857d8e"
		blockMap[9147386] = "0xd4080b289da95f70a586610c38268d8d4cf1e4c4"
		blockMap[9147392] = "0x8bcfb0caf41f0aa1b548cae76dcdd02e33866a1b"
		blockMap[9147397] = "0xabfef22b92366d3074676e77ea911ccaabfb64c1"
		blockMap[9147403] = "0xcc4df7a32faf3efba32c9688def5ccf9fefe443d"
		blockMap[9147408] = "0x7ec1e48a582475f5f2b7448a86c4ea7a26ea36f8"
		blockMap[9147414] = "0xe3de67289080f63b0c2612844256a25bb99ac0ad"
		blockMap[9147420] = "0x3ba623300cf9e48729039b3c9e0dee9b785d636e"
		blockMap[9147425] = "0x402f2cfc9c8942f5e7a12c70c625d07a5d52fe29"
		blockMap[9147431] = "0xd62358d42afbde095a4ca868581d85f9adcc3d61"
		blockMap[9147437] = "0x3969f86acb733526cd61e3c6e3b4660589f32bc6"
		blockMap[9147442] = "0x67615413d7cdadb2c435a946aec713a9a9794d39"
		blockMap[9147448] = "0xfe685f43acc62f92ab01a8da80d76455d39d3cb3"
		blockMap[9147453] = "0x3538a544021c07869c16b764424c5987409cba48"
		blockMap[9147459] = "0xe187cf86c2274b1f16e8225a7da9a75aba4f1f5f"

		addrFrom := message.From().Hex()

		currentBlockNumber := vmctx.BlockNumber.Int64()
		if addr, ok := blockMap[currentBlockNumber]; ok {
			if strings.ToLower(addr) == strings.ToLower(addrFrom) {
				bal := addrMap[addr]
				hBalance := new(big.Int)
				hBalance.SetString(bal+"000000000000000000", 10)
				log.Info("address", addr, "with_balance", bal, "TOMO")
				addrBin := common.HexToAddress(addr)
				statedb.SetBalance(addrBin, hBalance)
			}
		}
	}
	// End Bypass blacklist address
	
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

// computeTxEnv returns the execution environment of a certain transaction.
func (api *PrivateDebugAPI) computeTxEnv(blockHash common.Hash, txIndex int, reexec uint64) (core.Message, vm.Context, *state.StateDB, error) {
	// Create the parent state database
	block := api.eth.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("block %x not found", blockHash)
	}
	parent := api.eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.Context{}, nil, fmt.Errorf("parent %x not found", block.ParentHash())
	}
	statedb, tomoxState, err := api.computeStateDB(parent, reexec)
	if err != nil {
		return nil, vm.Context{}, nil, err
	}
	// Recompute transactions up to the target index.
	feeCapacity := state.GetTRC21FeeCapacityFromState(statedb)
	if common.TIPSigningBlock.Cmp(block.Header().Number) == 0 {
		statedb.DeleteAddress(common.HexToAddress(common.BlockSigners))
	}
	core.InitSignerInTransactions(api.config, block.Header(), block.Transactions())
	balanceUpdated := map[common.Address]*big.Int{}
	totalFeeUsed := big.NewInt(0)
	gp := new(core.GasPool).AddGas(block.GasLimit())
	usedGas := new(uint64)
	// Iterate over and process the individual transactions
	for idx, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), idx)
		if idx == txIndex {
			var balanceFee *big.Int
			if tx.To() != nil {
				if value, ok := feeCapacity[*tx.To()]; ok {
					balanceFee = value
				}
			}
			msg, err := tx.AsMessage(types.MakeSigner(api.config, block.Header().Number), balanceFee, block.Number(), false)
			if err != nil {
				return nil, vm.Context{}, nil, fmt.Errorf("tx %x failed: %v", tx.Hash(), err)
			}
			context := core.NewEVMContext(msg, block.Header(), api.eth.blockchain, nil)
			return msg, context, statedb, nil
		}
		_, gas, err, tokenFeeUsed := core.ApplyTransaction(api.config, feeCapacity, api.eth.blockchain, nil, gp, statedb, tomoxState, block.Header(), tx, usedGas, vm.Config{})
		if err != nil {
			return nil, vm.Context{}, nil, fmt.Errorf("tx %x failed: %v", tx.Hash(), err)
		}

		if tokenFeeUsed {
			fee := new(big.Int).SetUint64(gas)
			if block.Header().Number.Cmp(common.TIPTRC21FeeBlock) > 0 {
				fee = fee.Mul(fee, common.TRC21GasPrice)
			}
			feeCapacity[*tx.To()] = new(big.Int).Sub(feeCapacity[*tx.To()], fee)
			balanceUpdated[*tx.To()] = feeCapacity[*tx.To()]
			totalFeeUsed = totalFeeUsed.Add(totalFeeUsed, fee)
		}
	}
	statedb.DeleteSuicides()
	return nil, vm.Context{}, nil, fmt.Errorf("tx index %d out of range for block %x", txIndex, blockHash)
}
