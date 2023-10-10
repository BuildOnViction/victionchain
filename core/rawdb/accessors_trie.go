// Copyright 2022 The go-ethereum Authors
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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>

package rawdb

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/log"
)

// HashScheme is the legacy hash-based state scheme with which trie nodes are
// stored in the disk with node hash as the database key. The advantage of this
// scheme is that different versions of trie nodes can be stored in disk, which
// is very beneficial for constructing archive nodes. The drawback is it will
// store different trie nodes on the same path to different locations on the disk
// with no data locality, and it's unfriendly for designing state pruning.
//
// Now this scheme is still kept for backward compatibility, and it will be used
// for archive node and some other tries(e.g. light trie).
const HashScheme = "hashScheme"

// ReadLegacyTrieNode retrieves the legacy trie node with the given
// associated node hash.
func ReadLegacyTrieNode(db ethdb.KeyValueReader, hash common.Hash) []byte {
	data, err := db.Get(hash.Bytes())
	if err != nil {
		return nil
	}
	return data
}

// HasLegacyTrieNode checks if the trie node with the provided hash is present in db.
func HasLegacyTrieNode(db ethdb.KeyValueReader, hash common.Hash) bool {
	ok, _ := db.Has(hash.Bytes())
	return ok
}

// WriteLegacyTrieNode writes the provided legacy trie node to database.
func WriteLegacyTrieNode(db ethdb.KeyValueWriter, hash common.Hash, node []byte) {
	if err := db.Put(hash.Bytes(), node); err != nil {
		log.Crit("Failed to store legacy trie node", "err", err)
	}
}

// DeleteLegacyTrieNode deletes the specified legacy trie node from database.
func DeleteLegacyTrieNode(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Delete(hash.Bytes()); err != nil {
		log.Crit("Failed to delete legacy trie node", "err", err)
	}
}
