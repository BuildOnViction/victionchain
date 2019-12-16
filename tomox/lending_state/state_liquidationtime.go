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
	"github.com/tomochain/tomochain/trie"
	"io"
)

type liquidationTimeState struct {
	time        common.Hash
	lendingBook common.Hash
	data        itemList
	db          *LendingStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by TradingStateDB.Commit.
	dbErr error

	// Write caches.
	trie database.Trie // storage trie, which becomes non-nil on first access

	cachedStorage map[common.Hash]common.Hash
	dirtyStorage  map[common.Hash]common.Hash

	onDirty func(time common.Hash) // Callback method to mark a state object newly dirty
}

func (s *liquidationTimeState) empty() bool {
	return common.EmptyHash(s.data.Root)
}

func newLiquidationTimeState(db *LendingStateDB, time common.Hash, lendingBook common.Hash, data itemList, onDirty func(time common.Hash)) *liquidationTimeState {
	return &liquidationTimeState{
		db:            db,
		lendingBook:   lendingBook,
		time:          time,
		data:          data,
		cachedStorage: make(map[common.Hash]common.Hash),
		dirtyStorage:  make(map[common.Hash]common.Hash),
		onDirty:       onDirty,
	}
}

func (self *liquidationTimeState) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, self.data)
}

func (self *liquidationTimeState) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *liquidationTimeState) getTrie(db database.Database) database.Trie {
	if self.trie == nil {
		var err error
		self.trie, err = db.OpenStorageTrie(self.lendingBook, self.data.Root)
		if err != nil {
			self.trie, _ = db.OpenStorageTrie(self.time, EmptyHash)
			self.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return self.trie
}

func (self *liquidationTimeState) getAllLendingIds(db database.Database) []common.Hash {
	lendingIds := []common.Hash{}
	lendingBookTrie := self.getTrie(db)
	if lendingBookTrie == nil {
		return lendingIds
	}
	for id, _ := range self.cachedStorage {
		lendingIds = append(lendingIds, id)
	}
	orderListIt := trie.NewIterator(lendingBookTrie.NodeIterator(nil))
	for orderListIt.Next() {
		id := common.BytesToHash(orderListIt.Value)
		if _, exist := self.cachedStorage[id]; exist {
			continue
		}
		lendingIds = append(lendingIds, id)
	}
	return lendingIds
}

func (self *liquidationTimeState) insertLendingId(db database.Database, lendingId common.Hash) {
	self.setLendingId(lendingId, lendingId)
	self.setError(self.getTrie(db).TryUpdate(lendingId[:], lendingId[:]))
}

func (self *liquidationTimeState) removeLendingId(db database.Database, lendingId common.Hash) {
	tr := self.getTrie(db)
	self.setError(tr.TryDelete(lendingId[:]))
	self.setLendingId(lendingId, EmptyHash)
}

func (self *liquidationTimeState) setLendingId(lendingId common.Hash, value common.Hash) {
	self.cachedStorage[lendingId] = value
	self.dirtyStorage[lendingId] = value

	if self.onDirty != nil {
		self.onDirty(self.lendingBook)
		self.onDirty = nil
	}
}

func (self *liquidationTimeState) updateTrie(db database.Database) database.Trie {
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

func (self *liquidationTimeState) updateRoot(db database.Database) error {
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

func (self *liquidationTimeState) deepCopy(db *LendingStateDB, onDirty func(time common.Hash)) *liquidationTimeState {
	stateLendingBook := newLiquidationTimeState(db, self.lendingBook, self.time, self.data, onDirty)
	if self.trie != nil {
		stateLendingBook.trie = db.db.CopyTrie(self.trie)
	}
	for key, value := range self.dirtyStorage {
		stateLendingBook.dirtyStorage[key] = value
	}
	for key, value := range self.cachedStorage {
		stateLendingBook.cachedStorage[key] = value
	}
	return stateLendingBook
}
