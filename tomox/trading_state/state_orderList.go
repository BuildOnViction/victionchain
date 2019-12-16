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

package trading_state

import (
	"bytes"
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/tomox/database"
	"io"
)

type orderListState struct {
	price     common.Hash
	orderBook common.Hash
	orderType string
	data      itemList
	db        *TradingStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by TradingStateDB.Commit.
	dbErr error

	// Write caches.
	trie database.Trie // storage trie, which becomes non-nil on first access

	cachedStorage map[common.Hash]common.Hash // Storage entry cache to avoid duplicate reads
	dirtyStorage  map[common.Hash]common.Hash // Storage entries that need to be flushed to disk

	onDirty func(price common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *orderListState) empty() bool {
	return common.EmptyHash(s.data.Root)
}

// newObject creates a state object.
func newOrderListState(db *TradingStateDB, orderType string, orderBook common.Hash, price common.Hash, data itemList, onDirty func(price common.Hash)) *orderListState {
	return &orderListState{
		db:            db,
		orderType:     orderType,
		orderBook:     orderBook,
		price:         price,
		data:          data,
		cachedStorage: make(map[common.Hash]common.Hash),
		dirtyStorage:  make(map[common.Hash]common.Hash),
		onDirty:       onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *orderListState) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *orderListState) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (c *orderListState) getTrie(db database.Database) database.Trie {
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
func (self *orderListState) GetOrderAmount(db database.Database, orderId common.Hash) common.Hash {
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
func (self *orderListState) insertOrderItem(db database.Database, orderId common.Hash, amount common.Hash) {
	self.setOrderItem(orderId, amount)
	self.setError(self.getTrie(db).TryUpdate(orderId[:], amount[:]))
}

// SetState updates a value in orderId storage.
func (self *orderListState) removeOrderItem(db database.Database, orderId common.Hash) {
	tr := self.getTrie(db)
	self.setError(tr.TryDelete(orderId[:]))
	self.setOrderItem(orderId, EmptyHash)
}

func (self *orderListState) setOrderItem(orderId common.Hash, amount common.Hash) {
	self.cachedStorage[orderId] = amount
	self.dirtyStorage[orderId] = amount

	if self.onDirty != nil {
		self.onDirty(self.Price())
		self.onDirty = nil
	}
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *orderListState) updateTrie(db database.Database) database.Trie {
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
func (self *orderListState) updateRoot(db database.Database) error {
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

func (self *orderListState) deepCopy(db *TradingStateDB, onDirty func(price common.Hash)) *orderListState {
	stateObject := newOrderListState(db, self.orderType, self.orderBook, self.price, self.data, onDirty)
	if self.trie != nil {
		stateObject.trie = db.db.CopyTrie(self.trie)
	}
	for orderId, amount := range self.dirtyStorage {
		stateObject.dirtyStorage[orderId] = amount
	}
	for orderId, amount := range self.cachedStorage {
		stateObject.cachedStorage[orderId] = amount
	}
	return stateObject
}

// Returns the address of the contract/orderId
func (c *orderListState) Price() common.Hash {
	return c.price
}
