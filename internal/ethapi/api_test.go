package ethapi

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/tomochain/tomochain/accounts"
	"github.com/tomochain/tomochain/accounts/abi/bind"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/common/hexutil"
	"github.com/tomochain/tomochain/consensus"
	"github.com/tomochain/tomochain/consensus/ethash"
	"github.com/tomochain/tomochain/contracts/trc21issuer"
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
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
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
		genBlocks      = 10
		signer         = types.HomesteadSigner{}
		randomAccounts = newAccounts(2)
	)
	api := NewPublicBlockChainAPI(newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
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
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: nil,
			want:      21000,
		},
		// simple transfer with insufficient funds on latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: CallArgs{
				From:  &randomAccounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: core.ErrInsufficientFundsForTransfer,
			want:      21000,
		},
		// empty create
		{
			blockNumber: rpc.LatestBlockNumber,
			call:        CallArgs{},
			expectErr:   nil,
			want:        53000,
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call: CallArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			expectErr: core.ErrInsufficientFundsForTransfer,
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

func TestTRC21(t *testing.T) {
	// Initialize test accounts
	testPriKey, _ := crypto.HexToECDSA("0d782c534042ab93092d1baaf188e041ae429ca27d28d1a0d2ded2d3dd04c717")
	testAddr := crypto.PubkeyToAddress(testPriKey.PublicKey)
	fmt.Println("Public key: ", testAddr.String()) // 0x5C845F19EB923eEE213b620c12cc6D1d4E6E3506

	client, err := ethclient.Dial("http://127.0.0.1:8547")
	if err != nil {
		t.Fatal("Can't connect to RPC server: %", err)
	}

	nonce, _ := client.NonceAt(context.Background(), testAddr, nil)
	fmt.Println("Nonce", nonce)

	// Setup transactOpts
	auth := bind.NewKeyedTransactor(testPriKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei

	// Deploy TRC21
	trc21Addr, trc21Instance, err := trc21issuer.DeployTRC21(auth, client, "Viction", "VIC", 18, big.NewInt(1000000000000000000), big.NewInt(0))
	if err != nil {
		t.Fatal("Can't deploy TRC21: ", err)
	}
	fmt.Println("TRC21 address: ", trc21Addr.String())
	time.Sleep(10 * time.Second)

	// Get TRC21 name
	name, err := trc21Instance.Name()
	if err != nil {
		t.Fatal("Can't get name of TRC21: ", err)
	}
	fmt.Println("TRC21 name: ", name)

	// Attach TRC21Issuer to TRC21Issuer address
	trc21issuerAddr := common.TRC21IssuerSMC
	trc21issuerInstance, _ := trc21issuer.NewTRC21Issuer(auth, trc21issuerAddr, client)
	trc21IssuerMincap, err := trc21issuerInstance.MinCap()
	if err != nil {
		t.Fatal("Can't get min cap of trc21 issuer smart contract:", err)
	}
	fmt.Println("TRC21 Issuer min cap: ", trc21IssuerMincap)

	// Apply TRC21 issuer
	trc21issuerInstance.TransactOpts.Nonce = big.NewInt(int64(nonce + 1))
	trc21issuerInstance.TransactOpts.Value = new(big.Int).SetUint64(10000000000000000000)
	applyTx, err := trc21issuerInstance.Apply(trc21Addr)
	if err != nil {
		t.Fatal("Can't Apply free gas for token: ", err)
	}
	fmt.Println("Apply TRC21Issuer transaction: ", applyTx.Hash().Hex())
	time.Sleep(10 * time.Second)
	applyTxReceipt, err := client.TransactionReceipt(context.Background(), applyTx.Hash())
	if err != nil {
		t.Fatal("Can't get transaction receipt: ", err)
	}
	fmt.Println("Transaction receipt: ", applyTxReceipt)

	// Get balance token
	balanceBefore, err := trc21issuerInstance.GetTokenCapacity(trc21Addr)
	if err != nil {
		t.Fatal("Can't get token capacity of trc21 issuer smart contract:", err)
	}
	fmt.Println("TRC21 Issuer token capacity: ", balanceBefore)

	// Get test account balance
	testAccountBalanceBefore, err := client.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal("Can't get balance of test account: ", err)
	}

	// Transfer token to another address
	trc21Instance.TransactOpts.Nonce = big.NewInt(int64(nonce + 2))
	transferTx, err := trc21Instance.Transfer(common.HexToAddress("0x8A244cfdd4777E44bedEDCD478e62AC311EC30Dc"), big.NewInt(1000000000000000000))
	if err != nil {
		t.Fatal("Can't transfer token: ", err)
	}
	fmt.Println("Transfer token transaction: ", transferTx.Hash().Hex())
	time.Sleep(10 * time.Second)
	transferTxReceipt, err := client.TransactionReceipt(context.Background(), transferTx.Hash())
	if err != nil {
		t.Fatal("Can't get transaction receipt: ", err)
	}
	fmt.Println("Transaction receipt: ", transferTxReceipt)

	// Get test account balance after transfer
	testAccountBalanceAfter, err := client.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal("Can't get balance of test account: ", err)
	}

	if testAccountBalanceBefore.Cmp(testAccountBalanceAfter) != 0 {
		fmt.Println("Test failed: Test account balance before and after transfer is not equal")
	}

	// Get balance token
	balanceAfter, err := trc21issuerInstance.GetTokenCapacity(trc21Addr)
	if err != nil {
		t.Fatal("Can't get token capacity of trc21 issuer smart contract:", err)
	}
	fmt.Println("TRC21 Issuer token capacity: ", balanceAfter)

	if balanceBefore.Cmp(balanceAfter) <= 0 {
		t.Fatal("Test failed: Token balance fee before and after transfer is not correct")
	}
}
