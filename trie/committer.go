// Copyright 2019 The go-ethereum Authors
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
	"errors"
	"fmt"
	"sync"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/rlp"
	"golang.org/x/crypto/sha3"
)

// leafChanSize is the size of the leafCh. It's a pretty arbitrary number, to allow
// some parallelism but not incur too much memory overhead.
const leafChanSize = 200

// leaf represents a trie leaf value
type leaf struct {
	size   int         // size of the rlp data (estimate)
	hash   common.Hash // hash of rlp data
	node   Node        // the Node to commit
	vnodes bool        // set to true if the Node (possibly) contains a ValueNode
}

// committer is a type used for the trie Commit operation. A committer has some
// internal preallocated temp space, and also a callback that is invoked when
// leaves are committed. The leafs are passed through the `leafCh`,  to allow
// some level of parallelism.
// By 'some level' of parallelism, it's still the case that all leaves will be
// processed sequentially - onleaf will never be called in parallel or out of order.
type committer struct {
	tmp sliceBuffer
	sha keccakState

	onleaf LeafCallback
	leafCh chan *leaf
}

// committers live in a global sync.Pool
var committerPool = sync.Pool{
	New: func() interface{} {
		return &committer{
			tmp: make(sliceBuffer, 0, 550), // cap is as large as a full FullNode.
			sha: sha3.NewLegacyKeccak256().(keccakState),
		}
	},
}

// newCommitter creates a new committer or picks one from the pool.
func newCommitter() *committer {
	return committerPool.Get().(*committer)
}

func returnCommitterToPool(h *committer) {
	h.onleaf = nil
	h.leafCh = nil
	committerPool.Put(h)
}

// commitNeeded returns 'false' if the given Node is already in sync with Db
func (c *committer) commitNeeded(n Node) bool {
	hash, dirty := n.Cache()
	return hash == nil || dirty
}

// commit collapses a Node down into a hash Node and inserts it into the database
func (c *committer) Commit(n Node, db *Database) (HashNode, int, error) {
	if db == nil {
		return nil, 0, errors.New("no Db provided")
	}
	h, committed, err := c.commit(n, db, true)
	if err != nil {
		return nil, 0, err
	}
	return h.(HashNode), committed, nil
}

// commit collapses a Node down into a hash Node and inserts it into the database
func (c *committer) commit(n Node, db *Database, force bool) (Node, int, error) {
	// if this path is clean, use available cached data
	hash, dirty := n.Cache()
	if hash != nil && !dirty {
		return hash, 0, nil
	}
	// Commit children, then parent, and remove remove the dirty flag.
	switch cn := n.(type) {
	case *ShortNode:
		// Commit child
		collapsed := cn.copy()

		// If the child is fullNode, recursively commit,
		// otherwise it can only be hashNode or valueNode.
		var childCommitted int
		if _, ok := cn.Val.(ValueNode); !ok {
			if childV, committed, err := c.commit(cn.Val, db, false); err != nil {
				return nil, 0, err
			} else {
				collapsed.Val, childCommitted = childV, committed
			}
		}
		// The key needs to be copied, since we're delivering it to database
		collapsed.Key = hexToCompact(cn.Key)
		hashedNode := c.store(collapsed, db, force, true)
		if hn, ok := hashedNode.(HashNode); ok {
			return hn, childCommitted + 1, nil
		} else {
			return collapsed, childCommitted, nil
		}
	case *FullNode:
		hashedKids, hasVnodes, childCommitted, err := c.commitChildren(cn, db, force)
		if err != nil {
			return nil, 0, err
		}
		collapsed := cn.copy()
		collapsed.Children = hashedKids

		hashedNode := c.store(collapsed, db, force, hasVnodes)
		if hn, ok := hashedNode.(HashNode); ok {
			return hn, childCommitted + 1, nil
		} else {
			return collapsed, childCommitted, nil
		}
	case ValueNode:
		return c.store(cn, db, force, false), 0, nil
	// hashnodes aren't stored
	case HashNode:
		return cn, 0, nil
	}
	return hash, 0, nil
}

// commitChildren commits the children of the given fullnode
func (c *committer) commitChildren(n *FullNode, db *Database, force bool) ([17]Node, bool, int, error) {
	var (
		committed int
		children  [17]Node
	)

	var hasValueNodeChildren = false
	for i, child := range n.Children {
		if child == nil {
			continue
		}
		hnode, childCommitted, err := c.commit(child, db, false)
		if err != nil {
			return children, false, 0, err
		}
		children[i] = hnode
		if _, ok := hnode.(ValueNode); ok {
			hasValueNodeChildren = true
		}

		committed += childCommitted
	}
	return children, hasValueNodeChildren, committed, nil
}

// store hashes the Node n and if we have a storage layer specified, it writes
// the key/value pair to it and tracks any Node->child references as well as any
// Node->external trie references.
func (c *committer) store(n Node, db *Database, force bool, hasVnodeChildren bool) Node {
	// Larger nodes are replaced by their hash and stored in the database.
	var (
		hash, _ = n.Cache()
		size    int
	)
	if hash == nil {
		if vn, ok := n.(ValueNode); ok {
			c.tmp.Reset()
			if err := rlp.Encode(&c.tmp, vn); err != nil {
				panic("encode error: " + err.Error())
			}
			size = len(c.tmp)
			if size < 32 && !force {
				return n // Nodes smaller than 32 bytes are stored inside their parent
			}
			hash = c.makeHashNode(c.tmp)
		} else {
			// This was not generated - must be a small Node stored in the parent
			// No need to do anything here
			return n
		}
	} else {
		// We have the hash already, estimate the RLP encoding-size of the Node.
		// The size is used for mem tracking, does not need to be exact
		size = estimateSize(n)
	}
	// If we're using channel-based leaf-reporting, send to channel.
	// The leaf channel will be active only when there an active leaf-callback
	if c.leafCh != nil {
		c.leafCh <- &leaf{
			size:   size,
			hash:   common.BytesToHash(hash),
			node:   n,
			vnodes: hasVnodeChildren,
		}
	} else if db != nil {
		// No leaf-callback used, but there's still a database. Do serial
		// insertion
		db.Lock.Lock()
		db.insert(common.BytesToHash(hash), size, n)
		db.Lock.Unlock()
	}
	return hash
}

// commitLoop does the actual insert + leaf callback for nodes
func (c *committer) commitLoop(db *Database) {
	for item := range c.leafCh {
		var (
			hash      = item.hash
			size      = item.size
			n         = item.node
			hasVnodes = item.vnodes
		)
		// We are pooling the trie nodes into an intermediate memory Cache
		db.Lock.Lock()
		db.insert(hash, size, n)
		db.Lock.Unlock()
		if c.onleaf != nil && hasVnodes {
			switch n := n.(type) {
			case *ShortNode:
				if child, ok := n.Val.(ValueNode); ok {
					c.onleaf(child, hash)
				}
			case *FullNode:
				for i := 0; i < 16; i++ {
					if child, ok := n.Children[i].(ValueNode); ok {
						c.onleaf(child, hash)
					}
				}
			}
		}
	}
}

func (c *committer) makeHashNode(data []byte) HashNode {
	n := make(HashNode, c.sha.Size())
	c.sha.Reset()
	c.sha.Write(data)
	c.sha.Read(n)
	return n
}

// estimateSize estimates the size of an rlp-encoded Node, without actually
// rlp-encoding it (zero allocs). This method has been experimentally tried, and with a trie
// with 1000 leafs, the only errors above 1% are on small shortnodes, where this
// method overestimates by 2 or 3 bytes (e.g. 37 instead of 35)
func estimateSize(n Node) int {
	switch n := n.(type) {
	case *ShortNode:
		// A short Node contains a compacted key, and a value.
		return 3 + len(n.Key) + estimateSize(n.Val)
	case *FullNode:
		// A full Node contains up to 16 hashes (some nils), and a key
		s := 3
		for i := 0; i < 16; i++ {
			if child := n.Children[i]; child != nil {
				s += estimateSize(child)
			} else {
				s += 1
			}
		}
		return s
	case ValueNode:
		return 1 + len(n)
	case HashNode:
		return 1 + len(n)
	default:
		panic(fmt.Sprintf("Node type %T", n))

	}
}
