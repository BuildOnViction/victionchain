// Copyright 2014 The go-ethereum Authors
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

// Package state provides a caching layer atop the Ethereum state trie.
package lending_state

import (
	"fmt"
	"github.com/tomochain/tomochain/tomox/database"
	"math/big"
	"sort"
	"sync"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/trie"
)

type revision struct {
	id           int
	journalIndex int
}

// StateDBs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type LendingStateDB struct {
	db   database.Database
	trie database.Trie

	// This map holds 'live' objects, which will get modified while processing a state transition.
	lendingExchangeStates      map[common.Hash]*lendingExchangeState
	lendingExchangeStatesDirty map[common.Hash]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by LendingStateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        journal
	validRevisions []revision
	nextRevisionId int

	lock sync.Mutex
}

// Create a new state from a given trie.
func New(root common.Hash, db database.Database) (*LendingStateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return &LendingStateDB{
		db:                         db,
		trie:                       tr,
		lendingExchangeStates:      make(map[common.Hash]*lendingExchangeState),
		lendingExchangeStatesDirty: make(map[common.Hash]struct{}),
	}, nil
}

// setError remembers the first non-nil error it is called with.
func (self *LendingStateDB) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *LendingStateDB) Error() error {
	return self.dbErr
}

// Exist reports whether the given orderId address exists in the state.
// Notably this also returns true for suicided lenddinges.
func (self *LendingStateDB) Exist(addr common.Hash) bool {
	return self.getLendingExchange(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (self *LendingStateDB) Empty(addr common.Hash) bool {
	so := self.getLendingExchange(addr)
	return so == nil || so.empty()
}

func (self *LendingStateDB) GetNonce(addr common.Hash) uint64 {
	stateObject := self.getLendingExchange(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}
	return 0
}

// Database retrieves the low level database supporting the lower level trie ops.
func (self *LendingStateDB) Database() database.Database {
	return self.db
}

func (self *LendingStateDB) SetNonce(addr common.Hash, nonce uint64) {
	stateObject := self.GetOrNewLendingExchangeObject(addr)
	if stateObject != nil {
		self.journal = append(self.journal, nonceChange{
			hash: addr,
			prev: self.GetNonce(addr),
		})
		stateObject.SetNonce(nonce)
	}
}

func (self *LendingStateDB) InsertLendingItem(orderBook common.Hash, orderId common.Hash, order LendingItem) {
	priceHash := common.BigToHash(order.Price)
	stateExchange := self.getLendingExchange(orderBook)
	if stateExchange == nil {
		stateExchange = self.createLendingExchangeObject(orderBook)
	}
	var stateOrderList *itemListState
	switch order.Side {
	case INVESTING:
		stateOrderList = stateExchange.getInvestingOrderList(self.db, priceHash)
		if stateOrderList == nil {
			stateOrderList = stateExchange.createInvestingOrderList(self.db, priceHash)
		}
	case BORROWING:
		stateOrderList = stateExchange.getBorrowingOrderList(self.db, priceHash)
		if stateOrderList == nil {
			stateOrderList = stateExchange.createBorrowingOrderList(self.db, priceHash)
		}
	default:
		return
	}
	self.journal = append(self.journal, insertOrder{
		orderBook: orderBook,
		orderId:   orderId,
		order:     &order,
	})
	stateExchange.createLendingItem(self.db, orderId, order)
	stateOrderList.insertLendingItem(self.db, orderId, common.BigToHash(order.Quantity))
}

func (self *LendingStateDB) GetLendingOrder(orderBook common.Hash, orderId common.Hash) LendingItem {
	stateObject := self.GetOrNewLendingExchangeObject(orderBook)
	if stateObject == nil {
		return EmptyOrder
	}
	stateOrderItem := stateObject.getLendingItem(self.db, orderId)
	if stateOrderItem == nil {
		return EmptyOrder
	}
	return stateOrderItem.data
}

func (self *LendingStateDB) SubAmountLendingItem(orderBook common.Hash, orderId common.Hash, price *big.Int, amount *big.Int, side string) error {
	priceHash := common.BigToHash(price)
	lendingExchange := self.GetOrNewLendingExchangeObject(orderBook)
	if lendingExchange == nil {
		return fmt.Errorf("Order book not found : %s ", orderBook.Hex())
	}
	var itelmList *itemListState
	switch side {
	case INVESTING:
		itelmList = lendingExchange.getInvestingOrderList(self.db, priceHash)
	case BORROWING:
		itelmList = lendingExchange.getBorrowingOrderList(self.db, priceHash)
	default:
		return fmt.Errorf("Order type not found : %s ", side)
	}
	if itelmList == nil || itelmList.empty() {
		return fmt.Errorf("Order list empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderId.Hex(), priceHash.Hex())
	}
	lendingItem := lendingExchange.getLendingItem(self.db, orderId)
	if lendingItem == nil || lendingItem.empty() {
		return fmt.Errorf("Order item empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderId.Hex(), priceHash.Hex())
	}
	currentAmount := new(big.Int).SetBytes(itelmList.GetOrderAmount(self.db, orderId).Bytes()[:])
	if currentAmount.Cmp(amount) < 0 {
		return fmt.Errorf("Order amount not enough : %s , have : %d , want : %d ", orderId.Hex(), currentAmount, amount)
	}
	self.journal = append(self.journal, subAmountOrder{
		orderBook: orderBook,
		orderId:   orderId,
		order:     self.GetLendingOrder(orderBook, orderId),
		amount:    amount,
	})
	newAmount := new(big.Int).Sub(currentAmount, amount)
	lendingItem.setVolume(newAmount)
	if newAmount.Sign() == 0 {
		itelmList.removeOrderItem(self.db, orderId)
	} else {
		itelmList.setOrderItem(orderId, common.BigToHash(newAmount))
	}
	if itelmList.empty() {
		switch side {
		case INVESTING:
			lendingExchange.removeInvestingOrderList(self.db, itelmList)
		case BORROWING:
			lendingExchange.removeBorrowingOrderList(self.db, itelmList)
		default:
		}
	}
	return nil
}

func (self *LendingStateDB) CancelLendingOrder(orderBook common.Hash, order *LendingItem) error {
	priceHash := common.BigToHash(order.Price)
	orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
	stateObject := self.GetOrNewLendingExchangeObject(orderBook)
	if stateObject == nil {
		return fmt.Errorf("Order book not found : %s ", orderBook.Hex())
	}
	var lendingList *itemListState
	switch order.Side {
	case INVESTING:
		lendingList = stateObject.getInvestingOrderList(self.db, priceHash)
	case BORROWING:
		lendingList = stateObject.getBorrowingOrderList(self.db, priceHash)
	default:
		return fmt.Errorf("Order side not found : %s ", order.Side)
	}
	if lendingList == nil || lendingList.empty() {
		return fmt.Errorf("Order list empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderIdHash.Hex(), priceHash.Hex())
	}
	lendingItem := stateObject.getLendingItem(self.db, orderIdHash)
	if lendingItem == nil || lendingItem.empty() {
		return fmt.Errorf("Order item empty  order book : %s , order id  : %s , price  : %s ", orderBook, orderIdHash.Hex(), priceHash.Hex())
	}
	if lendingItem.data.UserAddress != order.UserAddress {
		return fmt.Errorf("Error Order User Address mismatch when cancel order book : %s , order id  : %s , got : %s , expect : %s ", orderBook, orderIdHash.Hex(), lendingItem.data.UserAddress.Hex(), order.UserAddress.Hex())
	}
	if lendingItem.data.ExchangeAddress != order.ExchangeAddress {
		return fmt.Errorf("Exchange Address mismatch when cancel. order book : %s , order id  : %s , got : %s , expect : %s ", orderBook, orderIdHash.Hex(), order.ExchangeAddress.Hex(), lendingItem.data.ExchangeAddress.Hex())
	}
	self.journal = append(self.journal, cancelOrder{
		orderBook: orderBook,
		orderId:   orderIdHash,
		order:     self.GetLendingOrder(orderBook, orderIdHash),
	})
	lendingItem.setVolume(big.NewInt(0))
	lendingList.removeOrderItem(self.db, orderIdHash)
	if lendingList.empty() {
		switch order.Side {
		case INVESTING:
			stateObject.removeInvestingOrderList(self.db, lendingList)
		case BORROWING:
			stateObject.removeBorrowingOrderList(self.db, lendingList)
		default:
		}
	}
	return nil
}

func (self *LendingStateDB) GetBestInvestingRate(orderBook common.Hash) *big.Int {
	stateObject := self.getLendingExchange(orderBook)
	if stateObject != nil {
		priceHash := stateObject.getBestInvestingRate(self.db)
		if common.EmptyHash(priceHash) {
			return Zero
		}
		orderList := stateObject.getInvestingOrderList(self.db, priceHash)
		if orderList == nil {
			log.Error("order list ask not found", "price", priceHash.Hex())
			return Zero
		}
		return new(big.Int).SetBytes(priceHash.Bytes())
	}
	return Zero
}

func (self *LendingStateDB) GetBestBorrowRate(orderBook common.Hash) *big.Int {
	stateObject := self.getLendingExchange(orderBook)
	if stateObject != nil {
		priceHash := stateObject.getBestBorrowingRate(self.db)
		if common.EmptyHash(priceHash) {
			return Zero
		}
		return new(big.Int).SetBytes(priceHash.Bytes())
	}
	return Zero
}

func (self *LendingStateDB) GetBestLendingIdAndAmount(orderBook common.Hash, price *big.Int, side string) (common.Hash, *big.Int, error) {
	stateObject := self.GetOrNewLendingExchangeObject(orderBook)
	if stateObject != nil {
		var stateOrderList *itemListState
		switch side {
		case INVESTING:
			stateOrderList = stateObject.getInvestingOrderList(self.db, common.BigToHash(price))
		case BORROWING:
			stateOrderList = stateObject.getBorrowingOrderList(self.db, common.BigToHash(price))
		default:
			return EmptyHash, Zero, fmt.Errorf("not found side :%s ", side)
		}
		if stateOrderList != nil {
			key, _, err := stateOrderList.getTrie(self.db).TryGetBestLeftKeyAndValue()
			if err != nil {
				return EmptyHash, Zero, err
			}
			orderId := common.BytesToHash(key)
			amount := stateOrderList.GetOrderAmount(self.db, orderId)
			return orderId, new(big.Int).SetBytes(amount.Bytes()), nil
		}
		return EmptyHash, Zero, fmt.Errorf("not found order list with lendingBook : %s , price : %d , side :%s ", orderBook.Hex(), price, side)
	}
	return EmptyHash, Zero, fmt.Errorf("not found lendingBook : %s ", orderBook.Hex())
}

// updateLendingExchange writes the given object to the trie.
func (self *LendingStateDB) updateLendingExchange(stateObject *lendingExchangeState) {
	addr := stateObject.Hash()
	data, err := rlp.EncodeToBytes(stateObject)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
	}
	self.setError(self.trie.TryUpdate(addr[:], data))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *LendingStateDB) getLendingExchange(addr common.Hash) (stateObject *lendingExchangeState) {
	// Prefer 'live' objects.
	if obj := self.lendingExchangeStates[addr]; obj != nil {
		return obj
	}
	// Load the object from the database.
	enc, err := self.trie.TryGet(addr[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data lendingObject
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateExchanges(self, addr, data, self.MarkLendingExchangeObjectDirty)
	self.lendingExchangeStates[addr] = obj
	return obj
}

func (self *LendingStateDB) setLendingExchangeObject(object *lendingExchangeState) {
	self.lendingExchangeStates[object.Hash()] = object
	self.lendingExchangeStatesDirty[object.Hash()] = struct{}{}
}

// Retrieve a state object or create a new state object if nil.
func (self *LendingStateDB) GetOrNewLendingExchangeObject(addr common.Hash) *lendingExchangeState {
	stateExchangeObject := self.getLendingExchange(addr)
	if stateExchangeObject == nil {
		stateExchangeObject = self.createLendingExchangeObject(addr)
	}
	return stateExchangeObject
}

// MarkStateLendObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *LendingStateDB) MarkLendingExchangeObjectDirty(addr common.Hash) {
	self.lendingExchangeStatesDirty[addr] = struct{}{}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *LendingStateDB) createLendingExchangeObject(hash common.Hash) (newobj *lendingExchangeState) {
	newobj = newStateExchanges(self, hash, lendingObject{}, self.MarkLendingExchangeObjectDirty)
	newobj.setNonce(0) // sets the object to dirty
	self.setLendingExchangeObject(newobj)
	return newobj
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (self *LendingStateDB) Copy() *LendingStateDB {
	self.lock.Lock()
	defer self.lock.Unlock()

	// Copy all the basic fields, initialize the memory ones
	state := &LendingStateDB{
		db:                         self.db,
		trie:                       self.db.CopyTrie(self.trie),
		lendingExchangeStates:      make(map[common.Hash]*lendingExchangeState, len(self.lendingExchangeStatesDirty)),
		lendingExchangeStatesDirty: make(map[common.Hash]struct{}, len(self.lendingExchangeStatesDirty)),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range self.lendingExchangeStatesDirty {
		state.lendingExchangeStatesDirty[addr] = struct{}{}
	}
	for addr, exchangeObject := range self.lendingExchangeStates {
		state.lendingExchangeStates[addr] = exchangeObject.deepCopy(state, state.MarkLendingExchangeObjectDirty)
	}

	return state
}

func (s *LendingStateDB) clearJournalAndRefund() {
	s.journal = nil
	s.validRevisions = s.validRevisions[:0]
}

// Snapshot returns an identifier for the current revision of the state.
func (self *LendingStateDB) Snapshot() int {
	id := self.nextRevisionId
	self.nextRevisionId++
	self.validRevisions = append(self.validRevisions, revision{id, len(self.journal)})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (self *LendingStateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(self.validRevisions), func(i int) bool {
		return self.validRevisions[i].id >= revid
	})
	if idx == len(self.validRevisions) || self.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := self.validRevisions[idx].journalIndex

	// Replay the journal to undo changes.
	for i := len(self.journal) - 1; i >= snapshot; i-- {
		self.journal[i].undo(self)
	}
	self.journal = self.journal[:snapshot]

	// Remove invalidated snapshots from the stack.
	self.validRevisions = self.validRevisions[:idx]
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (s *LendingStateDB) Finalise() {
	// Commit objects to the trie.
	for addr, stateObject := range s.lendingExchangeStates {
		if _, isDirty := s.lendingExchangeStatesDirty[addr]; isDirty {
			// Write any storage changes in the state object to its storage trie.
			stateObject.updateInvestingRoot(s.db)
			stateObject.updateBorrowingRoot(s.db)
			stateObject.updateOrderRoot(s.db)
			stateObject.updateLiquidationTimeRoot(s.db)
			// Update the object in the main orderId trie.
			s.updateLendingExchange(stateObject)
			//delete(s.investingStatesDirty, addr)
		}
	}
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root lendingBook of the state trie.
// It is called in between transactions to get the root lendingBook that
// goes into transaction receipts.
func (s *LendingStateDB) IntermediateRoot() common.Hash {
	s.Finalise()
	return s.trie.Hash()
}

// Commit writes the state to the underlying in-memory trie database.
func (s *LendingStateDB) Commit() (root common.Hash, err error) {
	defer s.clearJournalAndRefund()
	// Commit objects to the trie.
	for addr, stateObject := range s.lendingExchangeStates {
		if _, isDirty := s.lendingExchangeStatesDirty[addr]; isDirty {
			// Write any storage changes in the state object to its storage trie.
			if err := stateObject.CommitInvestingTrie(s.db); err != nil {
				return EmptyHash, err
			}
			if err := stateObject.CommitBorrowingTrie(s.db); err != nil {
				return EmptyHash, err
			}
			if err := stateObject.CommitLendingItemTrie(s.db); err != nil {
				return EmptyHash, err
			}
			if err := stateObject.CommitLiquidationTimeTrie(s.db); err != nil {
				return EmptyHash, err
			}
			// Update the object in the main orderId trie.
			s.updateLendingExchange(stateObject)
			delete(s.lendingExchangeStatesDirty, addr)
		}
	}
	// Write trie changes.
	root, err = s.trie.Commit(func(leaf []byte, parent common.Hash) error {
		var exchange lendingObject
		if err := rlp.DecodeBytes(leaf, &exchange); err != nil {
			return nil
		}
		if exchange.InvestingRoot != EmptyRoot {
			s.db.TrieDB().Reference(exchange.InvestingRoot, parent)
		}
		if exchange.BorrowingRoot != EmptyRoot {
			s.db.TrieDB().Reference(exchange.BorrowingRoot, parent)
		}
		if exchange.LendingItemRoot != EmptyRoot {
			s.db.TrieDB().Reference(exchange.LendingItemRoot, parent)
		}
		if exchange.LiquidationTimeRoot != EmptyRoot {
			s.db.TrieDB().Reference(exchange.LiquidationTimeRoot, parent)
		}
		return nil
	})
	log.Debug("Lending Trie cache stats after commit", "misses", trie.CacheMisses(), "unloads", trie.CacheUnloads(), "root", root.Hex())
	return root, err
}

func (self *LendingStateDB) InsertLiquidationData(lendingBook common.Hash, time uint64, lendingId common.Hash) {
	timeHash := common.BigToHash(new(big.Int).SetUint64(time))
	lendingExchangeState := self.getLendingExchange(lendingBook)
	if lendingExchangeState == nil {
		lendingExchangeState = self.createLendingExchangeObject(lendingBook)
	}
	liquidationPriceState := lendingExchangeState.getLiquidationTimeOrderList(self.db, timeHash)
	if liquidationPriceState == nil {
		liquidationPriceState = lendingExchangeState.createLiquidationTime(self.db, timeHash)
	}
	liquidationPriceState.insertLendingId(self.db, lendingId)
}

func (self *LendingStateDB) RemoveLiquidationData(lendingBook common.Hash, time uint64, lendingId common.Hash) {
	timeHash := common.BigToHash(new(big.Int).SetUint64(time))
	lendingExchangeState := self.getLendingExchange(lendingBook)
	if lendingExchangeState == nil {
		lendingExchangeState = self.createLendingExchangeObject(lendingBook)
	}
	liquidationPriceState := lendingExchangeState.getLiquidationTimeOrderList(self.db, timeHash)
	if liquidationPriceState == nil {
		liquidationPriceState = lendingExchangeState.createLiquidationTime(self.db, timeHash)
	}
	liquidationPriceState.removeLendingId(self.db, lendingId)
}

func (self *LendingStateDB) GetAllLiquidationData(lendingBook common.Hash, time uint64) []common.Hash {
	timeBig := new(big.Int).SetUint64(time)
	liquidationData := []common.Hash{}
	lendingExchangeState := self.getLendingExchange(lendingBook)
	if lendingExchangeState == nil {
		lendingExchangeState = self.createLendingExchangeObject(lendingBook)
	}
	lowestPriceHash, liquidationState := lendingExchangeState.getLowestLiquidationTime(self.db)
	lowestTime := new(big.Int).SetBytes(lowestPriceHash[:])
	for liquidationState != nil && lowestTime.Sign() > 0 && lowestTime.Cmp(timeBig) <= 0 {
		data := liquidationState.getAllLendingIds(self.db)
		if len(data) == 0 {
			continue
		}
		liquidationData = append(liquidationData, data...)
	}
	return liquidationData
}
