package ethapi

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tomochain/tomochain/accounts"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/consensus/ethash"
	"github.com/tomochain/tomochain/core"
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
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

type testBackend struct {
	db      ethdb.Database
	chain   *core.BlockChain
	pending *types.Block
	TomoX   *tomox.TomoX
}

func (t testBackend) Downloader() *downloader.Downloader {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) ProtocolVersion() int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) ChainDb() ethdb.Database {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) EventMux() *event.TypeMux {
	//TODO implement me
	panic("implement me")
}

func (b testBackend) AccountManager() *accounts.Manager {
	return &accounts.Manager{}
}

func (b testBackend) TomoxService() *tomox.TomoX {
	return b.TomoX
}

func (t testBackend) LendingService() *tomoxlending.Lending {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SetHead(number uint64) {
	//TODO implement me
	panic("implement me")
}

func (b testBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	if blockNr == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock().Header(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	if number == rpc.PendingBlockNumber {
		return b.pending, nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b testBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.chain.StateAt(header.Root)
	return stateDb, header, err
}

func (t testBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(t.db, blockHash, core.GetBlockNumber(t.db, blockHash)), nil
}

func (t testBackend) GetTd(blockHash common.Hash) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (b testBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, tomoxState *tradingstate.TradingStateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.chain, nil)
	return vm.NewEVM(context, state, tomoxState, b.chain.Config(), vmCfg), vmError, nil
}

func (t testBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetPoolTransactions() (types.Transactions, error) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) Stats() (pending int, queued int) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SubscribeTxPreEvent(events chan<- core.TxPreEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SendOrderTx(ctx context.Context, signedTx *types.OrderTransaction) error {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) OrderTxPoolContent() (map[common.Address]types.OrderTransactions, map[common.Address]types.OrderTransactions) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) OrderStats() (pending int, queued int) {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) SendLendingTx(ctx context.Context, signedTx *types.LendingTransaction) error {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) ChainConfig() *params.ChainConfig {
	return t.chain.Config()
}

func (t testBackend) CurrentBlock() *types.Block {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetIPCClient() (*ethclient.Client, error) {
	//TODO implement me
	panic("implement me")
}

func (b testBackend) GetEngine() consensus.Engine {
	return b.chain.Engine()
}

func (t testBackend) GetRewardByHash(hash common.Hash) map[string]map[string]map[string]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetVotersRewards(address common.Address) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetVotersCap(checkpoint *big.Int, masterAddr common.Address, voters []common.Address) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetEpochDuration() *big.Int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetMasternodesCap(checkpoint uint64) map[common.Address]*big.Int {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetBlocksHashCache(blockNr uint64) []common.Hash {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) AreTwoBlockSamePath(newBlock common.Hash, oldBlock common.Hash) bool {
	//TODO implement me
	panic("implement me")
}

func (t testBackend) GetOrderNonce(address common.Hash) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func newTestBackend(t *testing.T, n int, gspec *core.Genesis, generator func(i int, b *core.BlockGen)) *testBackend {
	var (
		engine      = ethash.NewFaker()
		cacheConfig = &core.CacheConfig{
			TrieTimeLimit: 5 * time.Minute,
			TrieNodeLimit: 256 * 1024 * 1024,
		}
	)
	// Generate blocks for testing
	db, blocks, _ := core.GenerateChainWithGenesis(gspec, engine, n, generator)
	chain, err := core.NewBlockChain(db, cacheConfig, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	tomo := tomox.New(&tomox.DefaultConfig)

	backend := &testBackend{db: db, chain: chain, TomoX: tomo}
	return backend
}

type Account struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

func newAccounts(n int) (accounts []Account) {
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	slices.SortFunc(accounts, func(a, b Account) int {
		return a.addr.Cmp(b.addr)
	})
	return accounts
}

func TestEstimateGas(t *testing.T) {
	t.Parallel()
	// Initialize test accounts
	var (
		accounts = newAccounts(2)
		genesis  = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc: core.GenesisAlloc{
				accounts[0].addr: {Balance: big.NewInt(params.Ether)},
				accounts[1].addr: {Balance: big.NewInt(params.Ether)},
			},
		}
		genBlocks = 10
		signer    = types.HomesteadSigner{}
	)
	api := NewPublicBlockChainAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		//tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &accounts[1].addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), signer, accounts[0].key)
		tx, err := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, nil, nil), signer, accounts[0].key)
		if err != nil {
			panic(err)
		}
		b.AddTx(tx)
	}))
	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        CallArgs
		expectErr   error
		want        uint64
	}{
		// simple transfer on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: CallArgs{
				From:  accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (hexutil.Big)(*big.NewInt(1000)),
			},
			expectErr: nil,
			want:      21000,
		},
		// empty create
		{
			blockNumber: rpc.LatestBlockNumber,
			call:        CallArgs{},
			expectErr:   nil,
			want:        53000,
		},
	}
	for i, tc := range testSuite {
		result, err := api.EstimateGas(context.Background(), tc.call, &tc.blockNumber)
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tc.expectErr)
				continue
			}
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, have %v", i, tc.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		if uint64(result) != tc.want {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, uint64(result), tc.want)
		}
	}
}

func TestRPCGetBlockReceipts(t *testing.T) {
	t.Parallel()

	var (
		genBlocks  = 3
		backend, _ = setupReceiptBackend(t, genBlocks)
		api        = NewPublicBlockChainAPI(backend)
	)
	blockHashes := make([]common.Hash, genBlocks+1)
	ctx := context.Background()
	for i := 0; i <= genBlocks; i++ {
		header, err := backend.HeaderByNumber(ctx, rpc.BlockNumber(i))
		if err != nil {
			t.Errorf("failed to get block: %d err: %v", i, err)
		}
		blockHashes[i] = header.Hash()
	}

	var testSuite = []struct {
		test rpc.BlockNumber
		want string
	}{
		// 1. block without any txs(number)
		{
			test: rpc.BlockNumber(0),
			want: `[]`,
		},
		// 2. earliest tag
		{
			test: rpc.EarliestBlockNumber,
			want: `[]`,
		},
		// 3. latest tag
		{
			test: rpc.LatestBlockNumber,
			want: `[{"blockHash":"0x7b30611be396a2b3135482fb49975fa1641b9703da2bb9e8ddef4dd5ab0c36e8", "blockNumber":"0x3", "contractAddress":null, "cumulativeGasUsed":"0xea60", "from":"0x703c4b2bd70c169f5717101caee543299fc946c7", "gasUsed":"0xea60", "logs":[], "logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "status":"0x0", "to":"0x0000000000000000000000000000000000031ec7", "transactionHash":"0x0fa8c0c52f331c690c832c11c9cdc6c9e635bc5b055729230b1eb2b35c53419f", "transactionIndex":"0x0"}]`,
		},
		// 5. block with contract create tx(number)
		{
			test: rpc.BlockNumber(2),
			want: `[{"blockHash":"0xa56b19f6ed7acd69a6b17ab17388cca59de28fe8c49ae62be68752476386b39d","blockNumber":"0x2","contractAddress":null,"cumulativeGasUsed":"0x5318","from":"0x703c4b2bd70c169f5717101caee543299fc946c7","gasUsed":"0x5318","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0x0000000000000000000000000000000000000000","transactionHash":"0x537c16d5b0f04d33a2a40bc879f892c2a8e5866a3a7db99eeb78165b003d3d55","transactionIndex":"0x0"}]`,
		},
		// 10. block is not found
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
		signer   = types.HomesteadSigner{}
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
			//tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &acc2Addr, Value: big.NewInt(1000), Gas: params.TxGas, GasPrice: b.BaseFee(), Data: nil}), types.HomesteadSigner{}, acc1Key)
			tx, err = types.SignTx(types.NewTransaction(uint64(i), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
		case 1:
			// create contract
			//tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: nil, Gas: 53100, GasPrice: b.BaseFee(), Data: common.FromHex("0x60806040")}), signer, acc1Key)
			tx, err = types.SignTx(types.NewTransaction(uint64(i), common.Address{}, nil, 53100, nil, common.FromHex("0x60806040")), signer, acc1Key)
		case 2:
			// with logs
			// transfer(address to, uint256 value)
			data := fmt.Sprintf("0xa9059cbb%s%s", common.HexToHash(common.BigToAddress(big.NewInt(int64(i + 1))).Hex()).String()[2:], common.BytesToHash([]byte{byte(i + 11)}).String()[2:])
			//tx, err = types.SignTx(types.NewTx(&types.LegacyTx{Nonce: uint64(i), To: &contract, Gas: 60000, GasPrice: b.BaseFee(), Data: common.FromHex(data)}), signer, acc1Key)
			tx, err = types.SignTx(types.NewTransaction(uint64(i), contract, nil, 60000, nil, common.FromHex(data)), signer, acc1Key)
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

func testRPCResponseWithFile(t *testing.T, testid int, result interface{}, rpc string, file string) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Errorf("test %d: json marshal error", testid)
		return
	}
	outputFile := filepath.Join("testdata", fmt.Sprintf("%s-%s.json", rpc, file))
	fmt.Println("outputFile: ", outputFile)
	if os.Getenv("WRITE_TEST_FILES") != "" {
		os.WriteFile(outputFile, data, 0644)
	}
	want, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("error reading expected test file: %s output: %v", outputFile, err)
	}
	require.JSONEqf(t, string(want), string(data), "test %d: json not match, want: %s, have: %s", testid, string(want), string(data))
}
