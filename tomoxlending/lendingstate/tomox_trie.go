// Copyright 2015 The go-ethereum Authors
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
	"github.com/tomochain/tomochain/ethdb"
	"github.com/tomochain/tomochain/trie"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
)

// TomoXTrie wraps a trie with key hashing. In a secure trie, all
// access operations hash the key using keccak256. This prevents
// calling code from creating long chains of nodes that
// increase the access time.
//
// Contrary to a regular trie, a TomoXTrie can only be created with
// New and must have an attached database. The database also stores
// the preimage of each key.
//
// TomoXTrie is not safe for concurrent use.
type TomoXTrie struct {
	trie             trie.Trie
	hashKeyBuf       [common.HashLength]byte
	secKeyCache      map[string][]byte
	secKeyCacheOwner *TomoXTrie // Pointer to self, replace the key cache on mismatch
}

// NewTomoXTrie creates a trie with an existing root node from a backing database
// and optional intermediate in-memory node pool.
//
// If root is the zero hash or the sha3 hash of an empty string, the
// trie is initially empty. Otherwise, New will panic if db is nil
// and returns MissingNodeError if the root node cannot be found.
//
// Accessing the trie loads nodes from the database or node pool on demand.
// Loaded nodes are kept around until their 'cache generation' expires.
// A new cache generation is created by each call to Commit.
// cachelimit sets the number of past cache generations to keep.
func NewTomoXTrie(root common.Hash, db *trie.Database, cachelimit uint16) (*TomoXTrie, error) {
	if db == nil {
		panic("trie.NewTomoXTrie called without a database")
	}
	trie, err := trie.New(root, db)
	if err != nil {
		return nil, err
	}
	trie.SetCacheLimit(cachelimit)
	return &TomoXTrie{trie: *trie}, nil
}

// Get returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
func (t *TomoXTrie) Get(key []byte) []byte {
	res, err := t.TryGet(key)
	if err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
	return res
}

// TryGet returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *TomoXTrie) TryGet(key []byte) ([]byte, error) {
	return t.trie.TryGet(key)
}

// TryGetBestLeftKey returns the value of max left leaf
// If a node was not found in the database, a MissingNodeError is returned.
func (t *TomoXTrie) TryGetBestLeftKeyAndValue() ([]byte, []byte, error) {
	return t.trie.TryGetBestLeftKeyAndValue()
}

// TryGetBestRightKey returns the value of max left leaf
// If a node was not found in the database, a MissingNodeError is returned.
func (t *TomoXTrie) TryGetBestRightKeyAndValue() ([]byte, []byte, error) {
	return t.trie.TryGetBestRightKeyAndValue()
}

// Update associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
func (t *TomoXTrie) Update(key, value []byte) {
	if err := t.TryUpdate(key, value); err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
}

// TryUpdate associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
//
// If a node was not found in the database, a MissingNodeError is returned.
func (t *TomoXTrie) TryUpdate(key, value []byte) error {
	err := t.trie.TryUpdate(key, value)
	if err != nil {
		return err
	}
	t.getSecKeyCache()[string(key)] = common.CopyBytes(key)
	return nil
}

// Delete removes any existing value for key from the trie.
func (t *TomoXTrie) Delete(key []byte) {
	if err := t.TryDelete(key); err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
}

// TryDelete removes any existing value for key from the trie.
// If a node was not found in the database, a MissingNodeError is returned.
func (t *TomoXTrie) TryDelete(key []byte) error {
	delete(t.getSecKeyCache(), string(key))
	return t.trie.TryDelete(key)
}

// GetKey returns the sha3 preimage of a hashed key that was
// previously used to store a value.
func (t *TomoXTrie) GetKey(shaKey []byte) []byte {
	if key, ok := t.getSecKeyCache()[string(shaKey)]; ok {
		return key
	}
	key, _ := t.trie.Db.Preimage(common.BytesToHash(shaKey))
	return key
}

// Commit writes all nodes and the secure hash pre-images to the trie's database.
// Nodes are stored with their sha3 hash as the key.
//
// Committing flushes nodes from memory. Subsequent Get calls will load nodes
// from the database.
func (t *TomoXTrie) Commit(onleaf trie.LeafCallback) (root common.Hash, err error) {
	// Write all the pre-images to the actual disk database
	if len(t.getSecKeyCache()) > 0 {
		t.trie.Db.Lock.Lock()
		for hk, key := range t.secKeyCache {
			t.trie.Db.InsertPreimage(common.BytesToHash([]byte(hk)), key)
		}
		t.trie.Db.Lock.Unlock()

		t.secKeyCache = make(map[string][]byte)
	}
	// Commit the trie to its intermediate node database
	return t.trie.Commit(onleaf)
}

func (t *TomoXTrie) Hash() common.Hash {
	return t.trie.Hash()
}

func (t *TomoXTrie) Root() []byte {
	return t.trie.Root()
}

func (t *TomoXTrie) Copy() *TomoXTrie {
	cpy := *t
	return &cpy
}

// NodeIterator returns an iterator that returns nodes of the underlying trie. Iteration
// starts at the key after the given start key.
func (t *TomoXTrie) NodeIterator(start []byte) trie.NodeIterator {
	return t.trie.NodeIterator(start)
}

// hashKey returns the hash of key as an ephemeral buffer.
// The caller must not hold onto the return value because it will become
// invalid on the next call to hashKey or secKey.
//func (t *TomoXTrie) hashKey(key []byte) []byte {
//	h := newHasher(0, 0, nil)
//	h.sha.Reset()
//	h.sha.Write(key)
//	buf := h.sha.Sum(t.hashKeyBuf[:0])
//	returnHasherToPool(h)
//	return buf
//}

// getSecKeyCache returns the current secure key cache, creating a new one if
// ownership changed (i.e. the current secure trie is a copy of another owning
// the actual cache).
func (t *TomoXTrie) getSecKeyCache() map[string][]byte {
	if t != t.secKeyCacheOwner {
		t.secKeyCacheOwner = t
		t.secKeyCache = make(map[string][]byte)
	}
	return t.secKeyCache
}

// Prove constructs a merkle proof for key. The result contains all encoded nodes
// on the path to the value at key. The value itself is also included in the last
// node and can be retrieved by verifying the proof.
//
// If the trie does not contain a value for key, the returned proof contains all
// nodes of the longest existing prefix of the key (at least the root node), ending
// with the node that proves the absence of the key.
func (t *TomoXTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	return t.trie.Prove(key, fromLevel, proofDb)
}
