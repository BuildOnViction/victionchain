// Copyright 2016 The go-ethereum Authors
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

package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/event"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/params"
	"github.com/tomochain/tomochain/rlp"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
)

// txPermanent is the number of mined blocks after a mined transaction is
// considered permanent and no rollback is expected
var txPermanent = uint64(500)

// TxPool implements the transaction pool for light clients, which keeps track
// of the status of locally created transactions, detecting if they are included
// in a block (mined) or rolled back. There are no queued transactions since we
// always receive all locally signed transactions in the same order as they are
// created.
type TxPool struct {
	config       *params.ChainConfig
	signer       types.Signer
	quit         chan bool
	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	mu           sync.RWMutex
	chain        *LightChain
	odr          OdrBackend
	chainDb      ethdb.Database
	relay        TxRelayBackend
	head         common.Hash
	nonce        map[common.Address]uint64            // "pending" nonce
	pending      map[common.Hash]*types.Transaction   // pending transactions by tx hash
	mined        map[common.Hash][]*types.Transaction // mined transactions by block hash
	clearIdx     uint64                               // earliest block nr that can contain mined tx info

	homestead bool
}

// TxRelayBackend provides an interface to the mechanism that forwards transacions
// to the ETH network. The implementations of the functions should be non-blocking.
//
// Send instructs backend to forward new transactions
// NewHead notifies backend about a new head after processed by the tx pool,
//
//	including  mined and rolled back transactions since the last event
//
// Discard notifies backend about transactions that should be discarded either
//
//	because they have been replaced by a re-send or because they have been mined
//	long ago and no rollback is expected
type TxRelayBackend interface {
	Send(txs types.Transactions)
	NewHead(head common.Hash, mined []common.Hash, rollback []common.Hash)
	Discard(hashes []common.Hash)
}

// NewTxPool creates a new light transaction pool
func NewTxPool(config *params.ChainConfig, chain *LightChain, relay TxRelayBackend) *TxPool {
	pool := &TxPool{
		config:      config,
		signer:      types.NewEIP155Signer(config.ChainId),
		nonce:       make(map[common.Address]uint64),
		pending:     make(map[common.Hash]*types.Transaction),
		mined:       make(map[common.Hash][]*types.Transaction),
		quit:        make(chan bool),
		chainHeadCh: make(chan core.ChainHeadEvent, chainHeadChanSize),
		chain:       chain,
		relay:       relay,
		odr:         chain.Odr(),
		chainDb:     chain.Odr().Database(),
		head:        chain.CurrentHeader().Hash(),
		clearIdx:    chain.CurrentHeader().Number.Uint64(),
	}
	// Subscribe events from blockchain
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)
	go pool.eventLoop()

	return pool
}

// currentState returns the light state of the current head header
func (pool *TxPool) currentState(ctx context.Context) *state.StateDB {
	return NewState(ctx, pool.chain.CurrentHeader(), pool.odr)
}

// GetNonce returns the "pending" nonce of a given address. It always queries
// the nonce belonging to the latest header too in order to detect if another
// client using the same key sent a transaction.
func (pool *TxPool) GetNonce(ctx context.Context, addr common.Address) (uint64, error) {
	state := pool.currentState(ctx)
	nonce := state.GetNonce(addr)
	if state.Error() != nil {
		return 0, state.Error()
	}
	sn, ok := pool.nonce[addr]
	if ok && sn > nonce {
		nonce = sn
	}
	if !ok || sn < nonce {
		pool.nonce[addr] = nonce
	}
	return nonce, nil
}

// txStateChanges stores the recent changes between pending/mined states of
// transactions. True means mined, false means rolled back, no entry means no change
type txStateChanges map[common.Hash]bool

// setState sets the status of a tx to either recently mined or recently rolled back
func (txc txStateChanges) setState(txHash common.Hash, mined bool) {
	val, ent := txc[txHash]
	if ent && (val != mined) {
		delete(txc, txHash)
	} else {
		txc[txHash] = mined
	}
}

// getLists creates lists of mined and rolled back tx hashes
func (txc txStateChanges) getLists() (mined []common.Hash, rollback []common.Hash) {
	for hash, val := range txc {
		if val {
			mined = append(mined, hash)
		} else {
			rollback = append(rollback, hash)
		}
	}
	return
}

// checkMinedTxs checks newly added blocks for the currently pending transactions
// and marks them as mined if necessary. It also stores block position in the db
// and adds them to the received txStateChanges map.
func (pool *TxPool) checkMinedTxs(ctx context.Context, hash common.Hash, number uint64, txc txStateChanges) error {
	// If no transactions are pending, we don't care about anything
	if len(pool.pending) == 0 {
		return nil
	}
	block, err := GetBlock(ctx, pool.odr, hash, number)
	if err != nil {
		return err
	}
	// Gather all the local transaction mined in this block
	list := pool.mined[hash]
	for _, tx := range block.Transactions() {
		if _, ok := pool.pending[tx.Hash()]; ok {
			list = append(list, tx)
		}
	}
	// If some transactions have been mined, write the needed data to disk and update
	if list != nil {
		// Retrieve all the receipts belonging to this block and write the loopup table
		if _, err := GetBlockReceipts(ctx, pool.odr, hash, number); err != nil { // ODR caches, ignore results
			return err
		}
		if err := core.WriteTxLookupEntries(pool.chainDb, block); err != nil {
			return err
		}
		// Update the transaction pool's state
		for _, tx := range list {
			delete(pool.pending, tx.Hash())
			txc.setState(tx.Hash(), true)
		}
		pool.mined[hash] = list
	}
	return nil
}

// rollbackTxs marks the transactions contained in recently rolled back blocks
// as rolled back. It also removes any positional lookup entries.
func (pool *TxPool) rollbackTxs(hash common.Hash, txc txStateChanges) {
	if list, ok := pool.mined[hash]; ok {
		for _, tx := range list {
			txHash := tx.Hash()
			core.DeleteTxLookupEntry(pool.chainDb, txHash)
			pool.pending[txHash] = tx
			txc.setState(txHash, false)
		}
		delete(pool.mined, hash)
	}
}

// reorgOnNewHead sets a new head header, processing (and rolling back if necessary)
// the blocks since the last known head and returns a txStateChanges map containing
// the recently mined and rolled back transaction hashes. If an error (context
// timeout) occurs during checking new blocks, it leaves the locally known head
// at the latest checked block and still returns a valid txStateChanges, making it
// possible to continue checking the missing blocks at the next chain head event
func (pool *TxPool) reorgOnNewHead(ctx context.Context, newHeader *types.Header) (txStateChanges, error) {
	txc := make(txStateChanges)
	oldh := pool.chain.GetHeaderByHash(pool.head)
	newh := newHeader
	// find common ancestor, create list of rolled back and new block hashes
	var oldHashes, newHashes []common.Hash
	for oldh.Hash() != newh.Hash() {
		if oldh.Number.Uint64() >= newh.Number.Uint64() {
			oldHashes = append(oldHashes, oldh.Hash())
			oldh = pool.chain.GetHeader(oldh.ParentHash, oldh.Number.Uint64()-1)
		}
		if oldh.Number.Uint64() < newh.Number.Uint64() {
			newHashes = append(newHashes, newh.Hash())
			newh = pool.chain.GetHeader(newh.ParentHash, newh.Number.Uint64()-1)
			if newh == nil {
				// happens when CHT syncing, nothing to do
				newh = oldh
			}
		}
	}
	if oldh.Number.Uint64() < pool.clearIdx {
		pool.clearIdx = oldh.Number.Uint64()
	}
	// roll back old blocks
	for _, hash := range oldHashes {
		pool.rollbackTxs(hash, txc)
	}
	pool.head = oldh.Hash()
	// check mined txs of new blocks (array is in reversed order)
	for i := len(newHashes) - 1; i >= 0; i-- {
		hash := newHashes[i]
		if err := pool.checkMinedTxs(ctx, hash, newHeader.Number.Uint64()-uint64(i), txc); err != nil {
			return txc, err
		}
		pool.head = hash
	}

	// clear old mined tx entries of old blocks
	if idx := newHeader.Number.Uint64(); idx > pool.clearIdx+txPermanent {
		idx2 := idx - txPermanent
		if len(pool.mined) > 0 {
			for i := pool.clearIdx; i < idx2; i++ {
				hash := core.GetCanonicalHash(pool.chainDb, i)
				if list, ok := pool.mined[hash]; ok {
					hashes := make([]common.Hash, len(list))
					for i, tx := range list {
						hashes[i] = tx.Hash()
					}
					pool.relay.Discard(hashes)
					delete(pool.mined, hash)
				}
			}
		}
		pool.clearIdx = idx2
	}

	return txc, nil
}

// blockCheckTimeout is the time limit for checking new blocks for mined
// transactions. Checking resumes at the next chain head event if timed out.
const blockCheckTimeout = time.Second * 3

// eventLoop processes chain head events and also notifies the tx relay backend
// about the new head hash and tx state changes
func (pool *TxPool) eventLoop() {
	for {
		select {
		case ev := <-pool.chainHeadCh:
			pool.setNewHead(ev.Block.Header())
			// hack in order to avoid hogging the lock; this part will
			// be replaced by a subsequent PR.
			time.Sleep(time.Millisecond)

		// System stopped
		case <-pool.chainHeadSub.Err():
			return
		}
	}
}

func (pool *TxPool) setNewHead(head *types.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), blockCheckTimeout)
	defer cancel()

	txc, _ := pool.reorgOnNewHead(ctx, head)
	m, r := txc.getLists()
	pool.relay.NewHead(pool.head, m, r)
	pool.homestead = pool.config.IsHomestead(head.Number)
	pool.signer = types.MakeSigner(pool.config, head.Number)
}

// Stop stops the light transaction pool
func (pool *TxPool) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	pool.scope.Close()
	// Unsubscribe subscriptions registered from blockchain
	pool.chainHeadSub.Unsubscribe()
	close(pool.quit)
	log.Info("Transaction pool stopped")
}

// SubscribeTxPreEvent registers a subscription of core.TxPreEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

// Stats returns the number of currently pending (locally created) transactions
func (pool *TxPool) Stats() (pending int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	pending = len(pool.pending)
	return
}

// validateTx checks whether a transaction is valid according to the consensus rules.
func (pool *TxPool) validateTx(ctx context.Context, tx *types.Transaction) error {
	// Validate sender
	var (
		from common.Address
		err  error
	)

	// check if sender is in black list
	if tx.From() != nil && common.Blacklist[*tx.From()] {
		return fmt.Errorf("Reject transaction with sender in black-list: %v", tx.From().Hex())
	}
	// check if receiver is in black list
	if tx.To() != nil && common.Blacklist[*tx.To()] {
		return fmt.Errorf("Reject transaction with receiver in black-list: %v", tx.To().Hex())
	}

	header := pool.chain.GetHeaderByHash(pool.head)
	// validate balance slot, minFee slot for TomoZ
	if pool.config.IsTomoZEnabled(header.Number) && tx.IsTomoZApplyTransaction() {
		copyState := pool.currentState(ctx).Copy()
		if err := core.ValidateTomoZApplyTransaction(pool.chain, copyState, common.BytesToAddress(tx.Data()[4:])); err != nil {
			return err
		}
	}
	// validate balance slot, token decimal for TomoX
	if pool.config.IsTomoXEnabled(header.Number) && tx.IsTomoXApplyTransaction() {
		copyState := pool.currentState(ctx).Copy()
		if err := core.ValidateTomoXApplyTransaction(pool.chain, copyState, common.BytesToAddress(tx.Data()[4:])); err != nil {
			return err
		}
	}

	// Validate the transaction sender and it's sig. Throw
	// if the from fields is invalid.
	if from, err = types.Sender(pool.signer, tx); err != nil {
		return core.ErrInvalidSender
	}
	// Last but not least check for nonce errors
	currentState := pool.currentState(ctx)
	if n := currentState.GetNonce(from); n > tx.Nonce() {
		return core.ErrNonceTooLow
	}

	// Check the transaction doesn't exceed the current
	// block limit gas.
	if header.GasLimit < tx.Gas() {
		return core.ErrGasLimit
	}

	// Transactions can't be negative. This may never happen
	// using RLP decoded transactions but may occur if you create
	// a transaction using the RPC for example.
	if tx.Value().Sign() < 0 {
		return core.ErrNegativeValue
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if b := currentState.GetBalance(from); b.Cmp(tx.Cost()) < 0 {
		return core.ErrInsufficientFunds
	}

	// Should supply enough intrinsic gas
	gas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, pool.homestead)
	if err != nil {
		return err
	}
	if tx.Gas() < gas {
		return core.ErrIntrinsicGas
	}
	return currentState.Error()
}

// add validates a new transaction and sets its state pending if processable.
// It also updates the locally stored nonce if necessary.
func (self *TxPool) add(ctx context.Context, tx *types.Transaction) error {
	hash := tx.Hash()

	if self.pending[hash] != nil {
		return fmt.Errorf("Known transaction (%x)", hash[:4])
	}
	err := self.validateTx(ctx, tx)
	if err != nil {
		return err
	}

	if _, ok := self.pending[hash]; !ok {
		self.pending[hash] = tx

		nonce := tx.Nonce() + 1

		addr, _ := types.Sender(self.signer, tx)
		if nonce > self.nonce[addr] {
			self.nonce[addr] = nonce
		}

		// Notify the subscribers. This event is posted in a goroutine
		// because it's possible that somewhere during the post "Remove transaction"
		// gets called which will then wait for the global tx pool lock and deadlock.
		go self.txFeed.Send(core.TxPreEvent{Tx: tx})
	}

	// Print a log message if low enough level is set
	log.Debug("Pooled new transaction", "hash", hash, "from", log.Lazy{Fn: func() common.Address { from, _ := types.Sender(self.signer, tx); return from }}, "to", tx.To())
	return nil
}

// Add adds a transaction to the pool if valid and passes it to the tx relay
// backend
func (self *TxPool) Add(ctx context.Context, tx *types.Transaction) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	if err := self.add(ctx, tx); err != nil {
		return err
	}
	//fmt.Println("Send", tx.Hash())
	self.relay.Send(types.Transactions{tx})

	self.chainDb.Put(tx.Hash().Bytes(), data)
	return nil
}

// AddTransactions adds all valid transactions to the pool and passes them to
// the tx relay backend
func (self *TxPool) AddBatch(ctx context.Context, txs []*types.Transaction) {
	self.mu.Lock()
	defer self.mu.Unlock()
	var sendTx types.Transactions

	for _, tx := range txs {
		if err := self.add(ctx, tx); err == nil {
			sendTx = append(sendTx, tx)
		}
	}
	if len(sendTx) > 0 {
		self.relay.Send(sendTx)
	}
}

// GetTransaction returns a transaction if it is contained in the pool
// and nil otherwise.
func (tp *TxPool) GetTransaction(hash common.Hash) *types.Transaction {
	// check the txs first
	if tx, ok := tp.pending[hash]; ok {
		return tx
	}
	return nil
}

// GetTransactions returns all currently processable transactions.
// The returned slice may be modified by the caller.
func (self *TxPool) GetTransactions() (txs types.Transactions, err error) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	txs = make(types.Transactions, len(self.pending))
	i := 0
	for _, tx := range self.pending {
		txs[i] = tx
		i++
	}
	return txs, nil
}

// Content retrieves the data content of the transaction pool, returning all the
// pending as well as queued transactions, grouped by account and nonce.
func (self *TxPool) Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	// Retrieve all the pending transactions and sort by account and by nonce
	pending := make(map[common.Address]types.Transactions)
	for _, tx := range self.pending {
		account, _ := types.Sender(self.signer, tx)
		pending[account] = append(pending[account], tx)
	}
	// There are no queued transactions in a light pool, just return an empty map
	queued := make(map[common.Address]types.Transactions)
	return pending, queued
}

// RemoveTransactions removes all given transactions from the pool.
func (self *TxPool) RemoveTransactions(txs types.Transactions) {
	self.mu.Lock()
	defer self.mu.Unlock()
	var hashes []common.Hash
	for _, tx := range txs {
		//self.RemoveTx(tx.Hash())
		hash := tx.Hash()
		delete(self.pending, hash)
		self.chainDb.Delete(hash[:])
		hashes = append(hashes, hash)
	}
	self.relay.Discard(hashes)
}

// RemoveTx removes the transaction with the given hash from the pool.
func (pool *TxPool) RemoveTx(hash common.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	// delete from pending pool
	delete(pool.pending, hash)
	pool.chainDb.Delete(hash[:])
	pool.relay.Discard([]common.Hash{hash})
}
