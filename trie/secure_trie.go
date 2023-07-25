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

package trie

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/rlp"
)

// SecureTrie wraps a trie with key hashing. In a secure trie, all
// access operations hash the key using keccak256. This prevents
// calling code from creating long chains of nodes that
// increase the access time.
//
// Contrary to a regular trie, a SecureTrie can only be created with
// New and must have an attached database. The database also stores
// the Preimage of each key.
//
// SecureTrie is not safe for concurrent use.
type SecureTrie struct {
	trie             Trie
	hashKeyBuf       [common.HashLength]byte
	secKeyCache      map[string][]byte
	secKeyCacheOwner *SecureTrie // Pointer to self, replace the key Cache on mismatch
}

// NewSecure creates a trie with an existing root Node from a backing database
// and optional intermediate in-memory Node pool.
//
// If root is the zero hash or the sha3 hash of an empty string, the
// trie is initially empty. Otherwise, New will panic if Db is nil
// and returns MissingNodeError if the root Node cannot be found.
//
// Accessing the trie loads nodes from the database or Node pool on demand.
// Loaded nodes are kept around until their 'Cache generation' expires.
// A new Cache generation is created by each call to Commit.
// cache limit sets the number of past Cache generations to keep.
func NewSecure(root common.Hash, db *Database) (*SecureTrie, error) {
	if db == nil {
		panic("trie.NewSecure called without a database")
	}
	trie, err := New(root, db)
	if err != nil {
		return nil, err
	}
	return &SecureTrie{trie: *trie}, nil
}

// MustGet returns the value for key stored in the trie.
// The value bytes must not be modified by the caller.
//
// This function will omit any encountered error but just
// print out an error message.
func (t *SecureTrie) MustGet(key []byte) []byte {
	return t.trie.MustGet(t.hashKey(key))
}

// GetStorage attempts to retrieve a storage slot with provided account address
// and slot key. The value bytes must not be modified by the caller.
// If the specified storage slot is not in the trie, nil will be returned.
// If a trie node is not found in the database, a MissingNodeError is returned.
func (t *SecureTrie) GetStorage(_ common.Address, key []byte) ([]byte, error) {
	enc, err := t.trie.Get(t.hashKey(key))
	if err != nil || len(enc) == 0 {
		return nil, err
	}
	_, content, _, err := rlp.Split(enc)
	return content, err
}

// GetAccount attempts to retrieve an account with provided account address.
// If the specified account is not in the trie, nil will be returned.
// If a trie node is not found in the database, a MissingNodeError is returned.
func (t *SecureTrie) GetAccount(address common.Address) (*types.StateAccount, error) {
	res, err := t.trie.Get(t.hashKey(address.Bytes()))
	if res == nil || err != nil {
		return nil, err
	}
	ret := new(types.StateAccount)
	err = rlp.DecodeBytes(res, ret)
	return ret, err
}

// GetAccountByHash does the same thing as GetAccount, however it expects an
// account hash that is the hash of address. This constitutes an abstraction
// leak, since the client code needs to know the key format.
func (t *SecureTrie) GetAccountByHash(addrHash common.Hash) (*types.StateAccount, error) {
	res, err := t.trie.Get(addrHash.Bytes())
	if res == nil || err != nil {
		return nil, err
	}
	ret := new(types.StateAccount)
	err = rlp.DecodeBytes(res, ret)
	return ret, err
}

// GetNode attempts to retrieve a trie node by compact-encoded path. It is not
// possible to use keybyte-encoding as the path might contain odd nibbles.
// If the specified trie node is not in the trie, nil will be returned.
// If a trie node is not found in the database, a MissingNodeError is returned.
func (t *SecureTrie) GetNode(path []byte) ([]byte, int, error) {
	return t.trie.GetNode(path)
}

// MustUpdate associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
//
// This function will omit any encountered error but just print out an
// error message.
func (t *SecureTrie) MustUpdate(key, value []byte) {
	hk := t.hashKey(key)
	t.trie.MustUpdate(hk, value)
	t.getSecKeyCache()[string(hk)] = common.CopyBytes(key)
}

// UpdateStorage associates key with value in the trie. Subsequent calls to
// Get will return value. If value has length zero, any existing value
// is deleted from the trie and calls to Get will return nil.
//
// The value bytes must not be modified by the caller while they are
// stored in the trie.
//
// If a node is not found in the database, a MissingNodeError is returned.
func (t *SecureTrie) UpdateStorage(_ common.Address, key, value []byte) error {
	hk := t.hashKey(key)
	v, _ := rlp.EncodeToBytes(value)
	err := t.trie.Update(hk, v)
	if err != nil {
		return err
	}
	t.getSecKeyCache()[string(hk)] = common.CopyBytes(key)
	return nil
}

// UpdateAccount will abstract the write of an account to the secure trie.

func (t *SecureTrie) UpdateAccount(address common.Address, acc *types.StateAccount) error {
	hk := t.hashKey(address.Bytes())
	data, err := rlp.EncodeToBytes(acc)
	if err != nil {
		return err
	}
	if err := t.trie.Update(hk, data); err != nil {
		return err
	}
	t.getSecKeyCache()[string(hk)] = address.Bytes()
	return nil
}

func (t *SecureTrie) UpdateContractCode(_ common.Address, _ common.Hash, _ []byte) error {
	return nil
}

// MustDelete removes any existing value for key from the trie. This function
// will omit any encountered error but just print out an error message.
func (t *SecureTrie) MustDelete(key []byte) {
	hk := t.hashKey(key)
	delete(t.getSecKeyCache(), string(hk))
	t.trie.MustDelete(hk)
}

// DeleteStorage removes any existing storage slot from the trie.
// If the specified trie node is not in the trie, nothing will be changed.
// If a node is not found in the database, a MissingNodeError is returned.
func (t *SecureTrie) DeleteStorage(_ common.Address, key []byte) error {
	hk := t.hashKey(key)
	delete(t.getSecKeyCache(), string(hk))
	return t.trie.Delete(hk)
}

// DeleteAccount abstracts an account deletion from the trie.
func (t *SecureTrie) DeleteAccount(address common.Address) error {
	hk := t.hashKey(address.Bytes())
	delete(t.getSecKeyCache(), string(hk))
	return t.trie.Delete(hk)
}

// GetKey returns the sha3 Preimage of a hashed key that was
// previously used to store a value.
func (t *SecureTrie) GetKey(shaKey []byte) []byte {
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
func (t *SecureTrie) Commit(onleaf LeafCallback) (root common.Hash, err error) {
	// Write all the pre-images to the actual disk database
	if len(t.getSecKeyCache()) > 0 {
		t.trie.Db.Lock.Lock()
		for hk, key := range t.secKeyCache {
			t.trie.Db.InsertPreimage(common.BytesToHash([]byte(hk)), key)
		}
		t.trie.Db.Lock.Unlock()

		t.secKeyCache = make(map[string][]byte)
	}
	// Commit the trie to its intermediate Node database
	return t.trie.Commit(onleaf)
}

// Hash returns the root hash of SecureTrie. It does not write to the
// database and can be used even if the trie doesn't have one.
func (t *SecureTrie) Hash() common.Hash {
	return t.trie.Hash()
}

// Copy returns a copy of SecureTrie.
func (t *SecureTrie) Copy() *SecureTrie {
	cpy := *t
	return &cpy
}

// NodeIterator returns an iterator that returns nodes of the underlying trie. Iteration
// starts at the key after the given start key.
func (t *SecureTrie) NodeIterator(start []byte) NodeIterator {
	return t.trie.NodeIterator(start)
}

// hashKey returns the hash of key as an ephemeral buffer.
// The caller must not hold onto the return value because it will become
// invalid on the next call to hashKey or secKey.
func (t *SecureTrie) hashKey(key []byte) []byte {
	h := newHasher(false)
	h.sha.Reset()
	h.sha.Write(key)
	buf := h.sha.Sum(t.hashKeyBuf[:0])
	returnHasherToPool(h)
	return buf
}

// getSecKeyCache returns the current secure key Cache, creating a new one if
// ownership changed (i.e. the current secure trie is a copy of another owning
// the actual Cache).
func (t *SecureTrie) getSecKeyCache() map[string][]byte {
	if t != t.secKeyCacheOwner {
		t.secKeyCacheOwner = t
		t.secKeyCache = make(map[string][]byte)
	}
	return t.secKeyCache
}
