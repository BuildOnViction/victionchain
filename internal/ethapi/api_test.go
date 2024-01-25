// Copyright 2023 The go-ethereum Authors
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

package ethapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tomochain/tomochain"
	"github.com/tomochain/tomochain/accounts"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/math"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/consensus/ethash"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/bloombits"
	"github.com/tomochain/tomochain/core/rawdb"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/crypto"
	"github.com/tomochain/tomochain/eth/downloader"
	"github.com/tomochain/tomochain/ethclient"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/event"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/rpc"
	"github.com/tomochain/tomochain/tomox"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"github.com/tomochain/tomochain/tomoxlending"
)

type testBackend struct {
	db      ethdb.Database
	chain   *core.BlockChain
	pending *types.Block
}

func (b *testBackend) Downloader() *downloader.Downloader {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) ProtocolVersion() int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) EventMux() *event.TypeMux {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) TomoxService() *tomox.TomoX {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) LendingService() *tomoxlending.Lending {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) SendOrderTx(ctx context.Context, signedTx *types.OrderTransaction) error {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) OrderTxPoolContent() (map[common.Address]types.OrderTransactions, map[common.Address]types.OrderTransactions) {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) OrderStats() (pending int, queued int) {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) SendLendingTx(ctx context.Context, signedTx *types.LendingTransaction) error {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetIPCClient() (*ethclient.Client, error) {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetEngine() consensus.Engine {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetRewardByHash(hash common.Hash) map[string]map[string]map[string]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetVotersRewards(address common.Address) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetVotersCap(checkpoint *big.Int, masterAddr common.Address, voters []common.Address) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetEpochDuration() *big.Int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetMasternodesCap(checkpoint uint64) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetBlocksHashCache(blockNr uint64) []common.Hash {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) AreTwoBlockSamePath(newBlock common.Hash, oldBlock common.Hash) bool {
	//TODO implement me
	panic("implement me")
}

func (b *testBackend) GetOrderNonce(address common.Hash) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func newTestBackend(t *testing.T, n int, gspec *core.Genesis, generator func(i int, b *core.BlockGen)) *testBackend {
	var (
		engine      = ethash.NewFaker()
		cacheConfig = &core.CacheConfig{
			TrieTimeLimit: 5 * time.Minute,
		}
		db     = rawdb.NewMemoryDatabase()
		gBlock = gspec.MustCommit(db)
	)
	// Generate blocks for testing
	blocks, _ := core.GenerateChain(gspec.Config, gBlock, engine, db, n, generator)
	chain, err := core.NewBlockChain(db, cacheConfig, gspec.Config, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	backend := &testBackend{db: db, chain: chain}
	return backend
}

func (b *testBackend) setPendingBlock(block *types.Block) {
	b.pending = block
}

func (b testBackend) SyncProgress() tomochain.SyncProgress { return tomochain.SyncProgress{} }
func (b testBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (b testBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	return nil, nil, nil, nil, nil
}
func (b testBackend) ChainDb() ethdb.Database           { return b.db }
func (b testBackend) AccountManager() *accounts.Manager { return nil }
func (b testBackend) ExtRPCEnabled() bool               { return false }
func (b testBackend) RPCGasCap() uint64                 { return 10000000 }
func (b testBackend) RPCEVMTimeout() time.Duration      { return time.Second }
func (b testBackend) RPCTxFeeCap() float64              { return 0 }
func (b testBackend) UnprotectedAllowed() bool          { return false }
func (b testBackend) SetHead(number uint64)             {}
func (b testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock().Header(), nil
	}
	if number == rpc.PendingBlockNumber && b.pending != nil {
		return b.pending.Header(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}
func (b testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.chain.GetHeaderByHash(hash), nil
}
func (b testBackend) CurrentHeader() *types.Header { return b.chain.CurrentBlock().Header() }
func (b testBackend) CurrentBlock() *types.Block   { return b.chain.CurrentBlock() }
func (b testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.LatestBlockNumber {
		head := b.chain.CurrentBlock()
		return b.chain.GetBlock(head.Hash(), head.NumberU64()), nil
	}
	if number == rpc.PendingBlockNumber {
		return b.pending, nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}
func (b testBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.chain.GetBlockByHash(hash), nil
}
func (b testBackend) GetBody(ctx context.Context, hash common.Hash, number rpc.BlockNumber) (*types.Body, error) {
	return b.chain.GetBlock(hash, uint64(number.Int64())).Body(), nil
}
func (b testBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	if number == rpc.PendingBlockNumber {
		panic("pending state not implemented")
	}
	header, err := b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, err := b.chain.StateAt(header.Root)
	return stateDb, header, err
}
func (b testBackend) PendingBlockAndReceipts() (*types.Block, types.Receipts) { panic("implement me") }
func (b testBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	header, err := b.HeaderByHash(ctx, hash)
	if header == nil || err != nil {
		return nil, err
	}
	receipts := rawdb.GetBlockReceipts(b.db, hash, header.Number.Uint64(), b.chain.Config())
	return receipts, nil
}
func (b testBackend) GetTd(hash common.Hash) *big.Int {
	if b.pending != nil && hash == b.pending.Hash() {
		return nil
	}
	return big.NewInt(1)
}
func (b testBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, tomoxState *tradingstate.TradingStateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.chain, nil)
	return vm.NewEVM(context, state, tomoxState, b.chain.Config(), vmCfg), vmError, nil
}
func (b testBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	panic("implement me")
}
func (b testBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.GetTransaction(b.db, txHash)
	return tx, blockHash, blockNumber, index, nil
}
func (b testBackend) GetPoolTransactions() (types.Transactions, error)         { panic("implement me") }
func (b testBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction { panic("implement me") }
func (b testBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	panic("implement me")
}
func (b testBackend) Stats() (pending int, queued int) { panic("implement me") }
func (b testBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	panic("implement me")
}
func (b testBackend) TxPoolContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
	panic("implement me")
}
func (b testBackend) SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) ChainConfig() *params.ChainConfig { return b.chain.Config() }
func (b testBackend) Engine() consensus.Engine         { return b.chain.Engine() }
func (b testBackend) GetLogs(ctx context.Context, blockHash common.Hash, number uint64) ([][]*types.Log, error) {
	panic("implement me")
}
func (b testBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}
func (b testBackend) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}
func (b testBackend) BloomStatus() (uint64, uint64) { panic("implement me") }
func (b testBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	panic("implement me")
}

func setupReceiptBackend(t *testing.T, genBlocks int) (*testBackend, []common.Hash) {
	// Initialize test accounts
	var (
		acc1Key, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		acc2Key, _ = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		acc1Addr   = crypto.PubkeyToAddress(acc1Key.PublicKey)
		acc2Addr   = crypto.PubkeyToAddress(acc2Key.PublicKey)
		contract   = common.HexToAddress("0000000000000000000000000000000000031ec7")
		genesis    = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				acc1Addr: {Balance: big.NewInt(params.Ether)},
				acc2Addr: {Balance: big.NewInt(params.Ether)},
				// // SPDX-License-Identifier: GPL-3.0
				// pragma solidity >=0.7.0 <0.9.0;
				//
				// contract Token {
				//     event Transfer(address indexed from, address indexed to, uint256 value);
				//     function transfer(address to, uint256 value) public returns (bool) {
				//         emit Transfer(msg.sender, to, value);
				//         return true;
				//     }
				// }
				contract: {Balance: big.NewInt(params.Ether), Code: common.FromHex("0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063a9059cbb14610030575b600080fd5b61004a6004803603810190610045919061016a565b610060565b60405161005791906101c5565b60405180910390f35b60008273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516100bf91906101ef565b60405180910390a36001905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610101826100d6565b9050919050565b610111816100f6565b811461011c57600080fd5b50565b60008135905061012e81610108565b92915050565b6000819050919050565b61014781610134565b811461015257600080fd5b50565b6000813590506101648161013e565b92915050565b60008060408385031215610181576101806100d1565b5b600061018f8582860161011f565b92505060206101a085828601610155565b9150509250929050565b60008115159050919050565b6101bf816101aa565b82525050565b60006020820190506101da60008301846101b6565b92915050565b6101e981610134565b82525050565b600060208201905061020460008301846101e0565b9291505056fea2646970667358221220b469033f4b77b9565ee84e0a2f04d496b18160d26034d54f9487e57788fd36d564736f6c63430008120033")},
			},
		}
		signer   = types.LatestSignerForChainID(params.TestChainConfig.ChainId)
		txHashes = make([]common.Hash, genBlocks)
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		var (
			tx  *types.Transaction
			err error
		)
		switch i {
		case 0:
			// transfer 1000wei
			tx, err = types.SignTx(types.NewTransaction(uint64(i), acc2Addr, big.NewInt(1000), params.TxGas, common.MinGasPrice, nil), types.HomesteadSigner{}, acc1Key)
		case 1:
			// create contract
			tx, err = types.SignTx(types.NewContractCreation(uint64(i), big.NewInt(0), 55000, common.MinGasPrice, common.FromHex("0x60806040")), signer, acc1Key)
		case 2:
			// with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			tx, err = types.SignTx(types.NewTransaction(uint64(i), contract, big.NewInt(0), 60000, common.MinGasPrice, common.FromHex(data)), signer, acc1Key)
		}
		if err != nil {
			t.Errorf("failed to sign tx: %v", err)
		}
		if tx != nil {
			b.AddTx(tx)
			txHashes[i] = tx.Hash()
		}
	})
	return backend, txHashes
}

func TestRPCGetBlockReceipts(t *testing.T) {
	t.Parallel()

	var (
		genBlocks  = 3
		backend, _ = setupReceiptBackend(t, genBlocks)
		api        = NewPublicBlockChainAPI(backend)
	)

	var testSuite = []struct {
		test rpc.BlockNumber
		want string
	}{
		// 0. block without any txs(number)
		{
			test: rpc.BlockNumber(0),
			want: `[]`,
		},
		// 1. earliest tag
		{
			test: rpc.EarliestBlockNumber,
			want: `[]`,
		},
		// 2. latest tag with token transfer event emitted
		{
			test: rpc.LatestBlockNumber,
			want: `[{"blockHash":"0x47e3bf718d613e3cb6ad43d24da111c09648e123e356f794b7a57c17836c36c4","blockNumber":"0x3","contractAddress":null,"cumulativeGasUsed":"0x5f60","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0x5f60","logs":[{"address":"0x0000000000000000000000000000000000031ec7","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x000000000000000000000000703c4b2bd70c169f5717101caee543299fc946c7","0x0000000000000000000000000000000000000000000000000000000000000003"],"data":"0x000000000000000000000000000000000000000000000000000000000000000d","blockNumber":"0x3","transactionHash":"0x88ca741c77147cd8caee87a3aaa23d516ea5339fe6773679db85290e4bef3b0d","transactionIndex":"0x0","blockHash":"0x47e3bf718d613e3cb6ad43d24da111c09648e123e356f794b7a57c17836c36c4","logIndex":"0x0","removed":false}],"logsBloom":"0x00000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000800000000000000008000000000000000000000000000000000020000000080000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000400000000002000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000","status":"0x1","to":"0x0000000000000000000000000000000000031ec7","transactionHash":"0x88ca741c77147cd8caee87a3aaa23d516ea5339fe6773679db85290e4bef3b0d","transactionIndex":"0x0"}]`,
		},
		// 3. block with contract create tx(number)
		{
			test: rpc.BlockNumber(2),
			want: `[{"blockHash":"0x451ef5c58313dfe7ba9c31dd82325a3eced23810339ae262638f9ab0f5e4ed9f","blockNumber":"0x2","contractAddress":"0xae9bea628c4ce503dcfd7e305cab4e29e7476592","cumulativeGasUsed":"0xd01e","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0xd01e","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":null,"transactionHash":"0xd5d28884a80b3b467613541503b83712979d06e6d6392e4609ffe6998b11fd1b","transactionIndex":"0x0"}]`,
		},
		// 4. block number is not found
		{
			test: rpc.BlockNumber(genBlocks + 1),
			want: `null`,
		},
	}

	for i, tt := range testSuite {
		var (
			result interface{}
			err    error
		)
		result, err = api.GetBlockReceipts(context.Background(), tt.test)
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		data, err := json.Marshal(result)
		if err != nil {
			t.Errorf("test %d: json marshal error", i)
			continue
		}
		want, have := tt.want, string(data)
		require.JSONEqf(t, want, have, "test %d: json not match, want: %s, have: %s", i, want, have)
	}
}
