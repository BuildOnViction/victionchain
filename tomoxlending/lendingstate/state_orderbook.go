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

package lendingstate

import (
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/rlp"
	"io"
	"math/big"
)

// stateObject represents an Ethereum orderId which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// exchangeObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type stateExchanges struct {
	hash common.Hash // orderbookHashInterest of ethereum address of the orderId
	data exchangeObject
	db   *TomoXStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by TomoXStateDB.Commit.
	dbErr error

	// Write caches.
	asksTrie   Trie // storage trie, which becomes non-nil on first access
	bidsTrie   Trie // storage trie, which becomes non-nil on first access
	ordersTrie Trie // storage trie, which becomes non-nil on first access

	stateAskObjects      map[common.Hash]*stateOrderList
	stateAskObjectsDirty map[common.Hash]struct{}

	stateBidObjects      map[common.Hash]*stateOrderList
	stateBidObjectsDirty map[common.Hash]struct{}

	stateOrderObjects      map[common.Hash]*stateLendingItem
	stateOrderObjectsDirty map[common.Hash]struct{}

	onDirty func(hash common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *stateExchanges) empty() bool {
	return s.data.Nonce == 0 && common.EmptyHash(s.data.AskRoot) && common.EmptyHash(s.data.BidRoot)
}

// newObject creates a state object.
func newStateExchanges(db *TomoXStateDB, hash common.Hash, data exchangeObject, onDirty func(addr common.Hash)) *stateExchanges {
	return &stateExchanges{
		db:                     db,
		hash:                   hash,
		data:                   data,
		stateAskObjects:        make(map[common.Hash]*stateOrderList),
		stateBidObjects:        make(map[common.Hash]*stateOrderList),
		stateOrderObjects:      make(map[common.Hash]*stateLendingItem),
		stateAskObjectsDirty:   make(map[common.Hash]struct{}),
		stateBidObjectsDirty:   make(map[common.Hash]struct{}),
		stateOrderObjectsDirty: make(map[common.Hash]struct{}),
		onDirty:                onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateExchanges) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateExchanges) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *stateExchanges) getAsksTrie(db Database) Trie {
	if c.asksTrie == nil {
		var err error
		c.asksTrie, err = db.OpenStorageTrie(c.hash, c.data.AskRoot)
		if err != nil {
			c.asksTrie, _ = db.OpenStorageTrie(c.hash, EmptyHash)
			c.setError(fmt.Errorf("can't create asks trie: %v", err))
		}
	}
	return c.asksTrie
}

func (c *stateExchanges) getOrdersTrie(db Database) Trie {
	if c.ordersTrie == nil {
		var err error
		c.ordersTrie, err = db.OpenStorageTrie(c.hash, c.data.OrderRoot)
		if err != nil {
			c.ordersTrie, _ = db.OpenStorageTrie(c.hash, EmptyHash)
			c.setError(fmt.Errorf("can't create asks trie: %v", err))
		}
	}
	return c.ordersTrie
}

func (c *stateExchanges) getBestInterestAsksTrie(db Database) common.Hash {
	trie := c.getAsksTrie(db)
	encKey, encValue, err := trie.TryGetBestLeftKeyAndValue()
	if err != nil {
		log.Error("Failed find best Interest ask trie ", "orderbook", c.hash.Hex())
		return EmptyHash
	}
	if len(encKey) == 0 || len(encValue) == 0 {
		log.Debug("Not found get best ask trie", "encKey", encKey, "encValue", encValue)
		return EmptyHash
	}
	var data orderList
	if err := rlp.DecodeBytes(encValue, &data); err != nil {
		log.Error("Failed to decode state get best ask trie", "err", err)
		return EmptyHash
	}
	return common.BytesToHash(encKey)
}

func (c *stateExchanges) getBestBidsTrie(db Database) common.Hash {
	trie := c.getBidsTrie(db)
	encKey, encValue, err := trie.TryGetBestRightKeyAndValue()
	if err != nil {
		log.Error("Failed find best Interest bid trie ", "orderbook", c.hash.Hex())
		return EmptyHash
	}
	if len(encKey) == 0 || len(encValue) == 0 {
		log.Debug("Not found get best bid trie", "encKey", encKey, "encValue", encValue)
		return EmptyHash
	}
	var data orderList
	if err := rlp.DecodeBytes(encValue, &data); err != nil {
		log.Error("Failed to decode state get best bid trie", "err", err)
		return EmptyHash
	}
	return common.BytesToHash(encKey)
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateAsksTrie(db Database) Trie {
	tr := self.getAsksTrie(db)
	for Interest, orderList := range self.stateAskObjects {
		if _, isDirty := self.stateAskObjectsDirty[Interest]; isDirty {
			delete(self.stateAskObjectsDirty, Interest)
			if orderList.empty() {
				self.setError(tr.TryDelete(Interest[:]))
				continue
			}
			orderList.updateRoot(db)
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(orderList)
			self.setError(tr.TryUpdate(Interest[:], v))
		}
	}

	return tr
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) updateAsksRoot(db Database) error {
	self.updateAsksTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	self.data.AskRoot = self.asksTrie.Hash()
	return nil
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitAsksTrie(db Database) error {
	self.updateAsksTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.asksTrie.Commit(func(leaf []byte, parent common.Hash) error {
		var orderList orderList
		if err := rlp.DecodeBytes(leaf, &orderList); err != nil {
			return nil
		}
		if orderList.Root != EmptyRoot {
			db.TrieDB().Reference(orderList.Root, parent)
		}
		return nil
	})
	if err == nil {
		self.data.AskRoot = root
	}
	return err
}

func (c *stateExchanges) getBidsTrie(db Database) Trie {
	if c.bidsTrie == nil {
		var err error
		c.bidsTrie, err = db.OpenStorageTrie(c.hash, c.data.BidRoot)
		if err != nil {
			c.bidsTrie, _ = db.OpenStorageTrie(c.hash, EmptyHash)
			c.setError(fmt.Errorf("can't create bids trie: %v", err))
		}
	}
	return c.bidsTrie
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateBidsTrie(db Database) Trie {
	tr := self.getBidsTrie(db)
	for Interest, orderList := range self.stateBidObjects {
		if _, isDirty := self.stateBidObjectsDirty[Interest]; isDirty {
			delete(self.stateBidObjectsDirty, Interest)
			if orderList.empty() {
				self.setError(tr.TryDelete(Interest[:]))
				continue
			}
			orderList.updateRoot(db)
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(orderList)
			self.setError(tr.TryUpdate(Interest[:], v))
		}
	}
	return tr
}

func (self *stateExchanges) updateBidsRoot(db Database) {
	self.updateBidsTrie(db)
	self.data.BidRoot = self.bidsTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitBidsTrie(db Database) error {
	self.updateBidsTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.bidsTrie.Commit(func(leaf []byte, parent common.Hash) error {
		var orderList orderList
		if err := rlp.DecodeBytes(leaf, &orderList); err != nil {
			return nil
		}
		if orderList.Root != EmptyRoot {
			db.TrieDB().Reference(orderList.Root, parent)
		}
		return nil
	})
	if err == nil {
		self.data.BidRoot = root
	}
	return err
}

func (self *stateExchanges) deepCopy(db *TomoXStateDB, onDirty func(hash common.Hash)) *stateExchanges {
	stateExchanges := newStateExchanges(db, self.hash, self.data, onDirty)
	if self.asksTrie != nil {
		stateExchanges.asksTrie = db.db.CopyTrie(self.asksTrie)
	}
	if self.bidsTrie != nil {
		stateExchanges.bidsTrie = db.db.CopyTrie(self.bidsTrie)
	}
	if self.ordersTrie != nil {
		stateExchanges.ordersTrie = db.db.CopyTrie(self.ordersTrie)
	}
	for Interest, bidObject := range self.stateBidObjects {
		stateExchanges.stateBidObjects[Interest] = bidObject.deepCopy(db, self.MarkStateBidObjectDirty)
	}
	for Interest, _ := range self.stateBidObjectsDirty {
		stateExchanges.stateBidObjectsDirty[Interest] = struct{}{}
	}
	for Interest, askObject := range self.stateAskObjects {
		stateExchanges.stateAskObjects[Interest] = askObject.deepCopy(db, self.MarkStateAskObjectDirty)
	}
	for Interest, _ := range self.stateAskObjectsDirty {
		stateExchanges.stateAskObjectsDirty[Interest] = struct{}{}
	}
	for orderId, LendingItem := range self.stateOrderObjects {
		stateExchanges.stateOrderObjects[orderId] = LendingItem.deepCopy(self.MarkStateOrderObjectDirty)
	}
	for orderId, _ := range self.stateOrderObjectsDirty {
		stateExchanges.stateOrderObjectsDirty[orderId] = struct{}{}
	}
	return stateExchanges
}

// Returns the address of the contract/orderId
func (c *stateExchanges) Hash() common.Hash {
	return c.hash
}

func (self *stateExchanges) SetNonce(nonce uint64) {
	self.setNonce(nonce)
}

func (self *stateExchanges) setNonce(nonce uint64) {
	self.data.Nonce = nonce
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *stateExchanges) Nonce() uint64 {
	return self.data.Nonce
}

func (self *stateExchanges) setInterest(Interest *big.Int) {
	self.data.Interest = Interest
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *stateExchanges) Interest() *big.Int {
	return self.data.Interest
}

// updateStateExchangeObject writes the given object to the trie.
func (self *stateExchanges) removeStateOrderListAskObject(db Database, stateOrderList *stateOrderList) {
	self.setError(self.asksTrie.TryDelete(stateOrderList.Interest[:]))
}

// updateStateExchangeObject writes the given object to the trie.
func (self *stateExchanges) removeStateOrderListBidObject(db Database, stateOrderList *stateOrderList) {
	self.setError(self.bidsTrie.TryDelete(stateOrderList.Interest[:]))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *stateExchanges) getStateOrderListAskObject(db Database, Interest common.Hash) (stateOrderList *stateOrderList) {
	// Prefer 'live' objects.
	if obj := self.stateAskObjects[Interest]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getAsksTrie(db).TryGet(Interest[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data orderList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "orderId", Interest, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateOrderList(self.db, Bid, self.hash, Interest, data, self.MarkStateAskObjectDirty)
	self.stateAskObjects[Interest] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *stateExchanges) MarkStateAskObjectDirty(Interest common.Hash) {
	self.stateAskObjectsDirty[Interest] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *stateExchanges) createStateOrderListAskObject(db Database, Interest common.Hash) (newobj *stateOrderList) {
	newobj = newStateOrderList(self.db, Ask, self.hash, Interest, orderList{Volume: Zero}, self.MarkStateAskObjectDirty)
	self.stateAskObjects[Interest] = newobj
	self.stateAskObjectsDirty[Interest] = struct{}{}
	data, err := rlp.EncodeToBytes(newobj)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", Interest[:], err))
	}
	self.setError(self.asksTrie.TryUpdate(Interest[:], data))
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
	return newobj
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *stateExchanges) getStateBidOrderListObject(db Database, Interest common.Hash) (stateOrderList *stateOrderList) {
	// Prefer 'live' objects.
	if obj := self.stateBidObjects[Interest]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getBidsTrie(db).TryGet(Interest[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data orderList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "orderId", Interest, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateOrderList(self.db, Bid, self.hash, Interest, data, self.MarkStateBidObjectDirty)
	self.stateBidObjects[Interest] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *stateExchanges) MarkStateBidObjectDirty(Interest common.Hash) {
	self.stateBidObjectsDirty[Interest] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *stateExchanges) createStateBidOrderListObject(db Database, Interest common.Hash) (newobj *stateOrderList) {
	newobj = newStateOrderList(self.db, Bid, self.hash, Interest, orderList{Volume: Zero}, self.MarkStateBidObjectDirty)
	self.stateBidObjects[Interest] = newobj
	self.stateBidObjectsDirty[Interest] = struct{}{}
	data, err := rlp.EncodeToBytes(newobj)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", Interest[:], err))
	}
	self.setError(self.bidsTrie.TryUpdate(Interest[:], data))
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
	return newobj
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *stateExchanges) getStateOrderObject(db Database, orderId common.Hash) (stateLendingItem *stateLendingItem) {
	// Prefer 'live' objects.
	if obj := self.stateOrderObjects[orderId]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getOrdersTrie(db).TryGet(orderId[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data LendingItem
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order object", "orderId", orderId, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newStateLendingItem(self.hash, orderId, data, self.MarkStateOrderObjectDirty)
	self.stateOrderObjects[orderId] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *stateExchanges) MarkStateOrderObjectDirty(orderId common.Hash) {
	self.stateOrderObjectsDirty[orderId] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *stateExchanges) createStateOrderObject(db Database, orderId common.Hash, order LendingItem) (newobj *stateLendingItem) {
	newobj = newStateLendingItem(self.hash, orderId, order, self.MarkStateOrderObjectDirty)
	orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.LendingId))
	self.stateOrderObjects[orderIdHash] = newobj
	self.stateOrderObjectsDirty[orderIdHash] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.hash)
		self.onDirty = nil
	}
	return newobj
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateExchanges) updateOrdersTrie(db Database) Trie {
	tr := self.getOrdersTrie(db)
	for orderId, LendingItem := range self.stateOrderObjects {
		if _, isDirty := self.stateOrderObjectsDirty[orderId]; isDirty {
			delete(self.stateOrderObjectsDirty, orderId)
			if LendingItem.empty() {
				self.setError(tr.TryDelete(orderId[:]))
				continue
			}
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(LendingItem)
			self.setError(tr.TryUpdate(orderId[:], v))
		}
	}
	return tr
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) updateOrdersRoot(db Database) {
	self.updateOrdersTrie(db)
	self.data.OrderRoot = self.ordersTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *stateExchanges) CommitOrdersTrie(db Database) error {
	self.updateOrdersTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.ordersTrie.Commit(nil)
	if err == nil {
		self.data.OrderRoot = root
	}
	return err
}
