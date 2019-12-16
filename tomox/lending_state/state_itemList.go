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

package lending_state

import (
	"bytes"
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/tomox/database"
	"io"
)

// stateObject represents an Ethereum orderId which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// lendingObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type itemListState struct {
	lendingBook common.Hash
	lendingType string
	price       common.Hash
	data        itemList
	db          *LendingStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by LendingStateDB.Commit.
	dbErr error

	// Write caches.
	trie database.Trie // storage trie, which becomes non-nil on first access

	cachedStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to disk

	onDirty func(price common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *itemListState) empty() bool {
	return s.data.Volume == nil || s.data.Volume.Sign() == 0
}

// newObject creates a state object.
func newItemListState(db *LendingStateDB, lendingType string, lendingBook common.Hash, price common.Hash, data itemList, onDirty func(price common.Hash)) *itemListState {
	return &itemListState{
		db:            db,
		lendingType:   lendingType,
		lendingBook:   lendingBook,
		price:         price,
		data:          data,
		cachedStorage: make(map[common.Hash]common.Hash),
		dirtyStorage:  make(map[common.Hash]common.Hash),
		onDirty:       onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *itemListState) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *itemListState) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *itemListState) getTrie(db database.Database) database.Trie {
	if c.trie == nil {
		var err error
		c.trie, err = db.OpenStorageTrie(c.price, c.data.Root)
		if err != nil {
			c.trie, _ = db.OpenStorageTrie(c.price, EmptyHash)
			c.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return c.trie
}

// GetState returns a value in orderId storage.
func (self *itemListState) GetOrderAmount(db database.Database, orderId common.Hash) common.Hash {
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
func (self *itemListState) insertLendingItem(db database.Database, orderId common.Hash, amount common.Hash) {
	self.setOrderItem(orderId, amount)
	self.setError(self.getTrie(db).TryUpdate(orderId[:], amount[:]))
}

// SetState updates a value in orderId storage.
func (self *itemListState) removeOrderItem(db database.Database, orderId common.Hash) {
	tr := self.getTrie(db)
	self.setError(tr.TryDelete(orderId[:]))
	self.setOrderItem(orderId, EmptyHash)
}

func (self *itemListState) setOrderItem(orderId common.Hash, amount common.Hash) {
	self.cachedStorage[orderId] = amount
	self.dirtyStorage[orderId] = amount

	if self.onDirty != nil {
		self.onDirty(self.lendingBook)
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *itemListState) updateTrie(db database.Database) database.Trie {
	tr := self.getTrie(db)
	for orderId, amount := range self.dirtyStorage {
		delete(self.dirtyStorage, orderId)
		if (amount == EmptyHash) {
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
func (self *itemListState) updateRoot(db database.Database) error {
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

func (self *itemListState) deepCopy(db *LendingStateDB, onDirty func(price common.Hash)) *itemListState {
	stateOrderList := newItemListState(db, self.lendingType, self.lendingBook, self.price, self.data, onDirty)
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
