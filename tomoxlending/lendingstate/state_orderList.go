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
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/rlp"
)

// stateObject represents an Ethereum orderId which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// exchangeObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type stateOrderList struct {
	Interest  common.Hash
	orderBook common.Hash
	orderType string
	data      orderList
	db        *TomoXStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by TomoXStateDB.Commit.
	dbErr error

	// Write caches.
	trie Trie // storage trie, which becomes non-nil on first access

	cachedStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to disk

	onDirty func(Interest common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *stateOrderList) empty() bool {
	return s.data.Volume == nil || s.data.Volume.Cmp(Zero) == 0
}

// newObject creates a state object.
func newStateOrderList(db *TomoXStateDB, orderType string, orderBook common.Hash, Interest common.Hash, data orderList, onDirty func(Interest common.Hash)) *stateOrderList {
	return &stateOrderList{
		db:            db,
		orderType:     orderType,
		orderBook:     orderBook,
		Interest:      Interest,
		data:          data,
		cachedStorage: make(map[common.Hash]common.Hash),
		dirtyStorage:  make(map[common.Hash]common.Hash),
		onDirty:       onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateOrderList) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateOrderList) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *stateOrderList) getTrie(db Database) Trie {
	if c.trie == nil {
		var err error
		c.trie, err = db.OpenStorageTrie(c.Interest, c.data.Root)
		if err != nil {
			c.trie, _ = db.OpenStorageTrie(c.Interest, EmptyHash)
			c.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return c.trie
}

// GetState returns a value in orderId storage.
func (self *stateOrderList) GetOrderAmount(db Database, orderId common.Hash) common.Hash {
	amount, exists := self.cachedStorage[orderId]
	if exists {
		return amount
	}
	// Load from DB in case it is missing.
	enc, err := self.getTrie(db).TryGet(orderId[:])
	if err != nil {
		self.setError(err)
		return EmptyHash
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			self.setError(err)
		}
		amount.SetBytes(content)
	}
	self.cachedStorage[orderId] = amount
	return amount
}

// SetState updates a value in orderId storage.
func (self *stateOrderList) insertLendingItem(db Database, orderId common.Hash, amount common.Hash) {
	self.setLendingItem(orderId, amount)
	self.setError(self.getTrie(db).TryUpdate(orderId[:], amount[:]))
}

// SetState updates a value in orderId storage.
func (self *stateOrderList) removeLendingItem(db Database, orderId common.Hash) {
	tr := self.getTrie(db)
	self.setError(tr.TryDelete(orderId[:]))
	self.setLendingItem(orderId, EmptyHash)
}

func (self *stateOrderList) setLendingItem(orderId common.Hash, amount common.Hash) {
	self.cachedStorage[orderId] = amount
	self.dirtyStorage[orderId] = amount

	if self.onDirty != nil {
		self.onDirty(self.Interest)
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *stateOrderList) updateTrie(db Database) Trie {
	tr := self.getTrie(db)
	for orderId, amount := range self.dirtyStorage {
		delete(self.dirtyStorage, orderId)
		if amount == EmptyHash {
			self.setError(tr.TryDelete(orderId[:]))
			continue
		}
		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(amount[:], "\x00"))
		self.setError(tr.TryUpdate(orderId[:], v))
	}
	return tr
}

// UpdateRoot sets the trie root to the current root orderId of
func (self *stateOrderList) updateRoot(db Database) error {
	self.updateTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.trie.Commit(nil)
	if err == nil {
		self.data.Root = root
	}
	return err
}

func (self *stateOrderList) deepCopy(db *TomoXStateDB, onDirty func(Interest common.Hash)) *stateOrderList {
	stateOrderList := newStateOrderList(db, self.orderType, self.orderBook, self.Interest, self.data, onDirty)
	if self.trie != nil {
		stateOrderList.trie = db.db.CopyTrie(self.trie)
	}
	for orderId, amount := range self.dirtyStorage {
		stateOrderList.dirtyStorage[orderId] = amount
	}
	for orderId, amount := range self.cachedStorage {
		stateOrderList.cachedStorage[orderId] = amount
	}
	return stateOrderList
}

// AddVolume removes amount from c's balance.
// It is used to add funds to the destination exchanges of a transfer.
func (c *stateOrderList) AddVolume(amount *big.Int) {
	c.setVolume(new(big.Int).Add(c.data.Volume, amount))
}

// AddVolume removes amount from c's balance.
// It is used to add funds to the destination exchanges of a transfer.
func (c *stateOrderList) subVolume(amount *big.Int) {
	c.setVolume(new(big.Int).Sub(c.data.Volume, amount))
}

func (self *stateOrderList) setVolume(volume *big.Int) {
	self.data.Volume = volume
	if self.onDirty != nil {
		self.onDirty(self.Interest)
		self.onDirty = nil
	}
}

func (self *stateOrderList) Volume() *big.Int {
	return self.data.Volume
}
