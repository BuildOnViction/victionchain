// Copyright (c) 2019, Agiletech Viet Nam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// References: http://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package tomox

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/log"
	"encoding/hex"
)

// Tree holds elements of the red-black tree
type Tree struct {
	db          OrderDao
	rootKey     []byte
	size        uint64
	Comparator  Comparator
	FormatBytes FormatBytes
}

// NewWith instantiates a red-black tree with the custom comparator.
func NewWith(comparator Comparator, db OrderDao) *Tree {
	tree := &Tree{
		Comparator: comparator,
		db:         db,
	}

	return tree
}

func NewWithBytesComparator(db OrderDao) *Tree {
	return NewWith(
		bytes.Compare,
		db,
	)
}

func (tree *Tree) Root(dryrun bool, blockHash common.Hash) *Node {
	root, err := tree.GetNode(tree.rootKey, dryrun, blockHash)
	if err != nil {
		log.Error("Can't get tree.Root", "rootKey", hex.EncodeToString(tree.rootKey), "err", err)
		return nil
	}
	return root
}

func (tree *Tree) IsEmptyKey(key []byte) bool {
	return tree.db.IsEmptyKey(key)
}

func (tree *Tree) SetRootKey(key []byte, size uint64) {
	tree.rootKey = key
	tree.size = size
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Put(key []byte, value []byte, dryrun bool, blockHash common.Hash) error {
	var insertedNode *Node
	if tree.IsEmptyKey(tree.rootKey) {
		// Assert key is of comparator's type for initial tree
		item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
		log.Debug("RootKey hasn't been set. Set it now", "rootkey", hex.EncodeToString(key), "all keys", tree.KeysinString(true, blockHash))
		tree.rootKey = key
		insertedNode = &Node{Key: key, Item: item}
	} else {
		node := tree.Root(dryrun, blockHash)
		if node == nil {
			return fmt.Errorf("Error on inserting node into the tree. tree.Root() is nil")
		}
		loop := true

		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				node.Item.Value = value
				tree.Save(node, dryrun, blockHash)
				return nil
			case compare < 0:
				if tree.IsEmptyKey(node.LeftKey()) {
					node.LeftKey(key)
					tree.Save(node, dryrun, blockHash)
					item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
					nodeLeft := &Node{Key: key, Item: item}
					insertedNode = nodeLeft
					loop = false
				} else {
					node = node.Left(tree, dryrun, blockHash)
				}
			case compare > 0:
				if tree.IsEmptyKey(node.RightKey()) {
					node.RightKey(key)
					tree.Save(node, dryrun, blockHash)
					item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
					nodeRight := &Node{Key: key, Item: item}
					insertedNode = nodeRight
					loop = false
				} else {
					node = node.Right(tree, dryrun, blockHash)
				}

			}
		}

		insertedNode.ParentKey(node.Key)
		tree.Save(insertedNode, dryrun, blockHash)
	}

	tree.insertCase1(insertedNode, dryrun, blockHash)
	log.Debug("Put node", "insertedNode key", hex.EncodeToString(insertedNode.Key), "all keys", tree.KeysinString(true, blockHash))
	tree.Save(insertedNode, dryrun, blockHash)

	tree.size++
	return nil
}

func (tree *Tree) GetNode(key []byte, dryrun bool, blockHash common.Hash) (*Node, error) {

	item := &Item{}

	val, err := tree.db.Get(key, item, dryrun, blockHash)

	if err != nil || val == nil {
		return nil, err
	}
	return &Node{Key: key, Item: val.(*Item)}, err
}

func (tree *Tree) Has(key []byte, dryrun bool, blockHash common.Hash) (bool, error) {
	return tree.db.Has(key, dryrun, blockHash)
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key []byte, dryrun bool, blockHash common.Hash) (value []byte, found bool) {

	node, err := tree.GetNode(key, dryrun, blockHash)
	if err != nil {
		return nil, false
	}
	if node != nil {
		return node.Item.Value, true
	}
	return nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key []byte, dryrun bool, blockHash common.Hash) {
	var child *Node
	node, err := tree.GetNode(key, dryrun, blockHash)
	if err != nil {
		log.Error("Can't remove node", "key", hex.EncodeToString(key), "err", err)
		return
	}
	if node == nil {
		log.Error("Can't remove node. Node == nil", "key", hex.EncodeToString(key))
		return
	}

	var left, right *Node = nil, nil
	if !tree.IsEmptyKey(node.LeftKey()) {
		left = node.Left(tree, dryrun, blockHash)
	}
	if !tree.IsEmptyKey(node.RightKey()) {
		right = node.Right(tree, dryrun, blockHash)
	}

	if left != nil && right != nil {
		node = left.maximumNode(tree, dryrun, blockHash)
	}

	if left == nil || right == nil {
		if right == nil {
			child = left
		} else {
			child = right
		}
		if child == nil {
			parent := node.Parent(tree, dryrun, blockHash)
			if parent != nil {
				if !tree.IsEmptyKey(parent.LeftKey()) && bytes.Equal(parent.LeftKey(), node.Key) {
					parent.Item.Keys.Left = EmptyKey()
					log.Debug("Set parent.LeftKey() to nil", "parent.Key", parent.Item.Keys)
					tree.Save(parent, dryrun, blockHash)
				} else if !tree.IsEmptyKey(parent.RightKey()) && bytes.Equal(parent.RightKey(), node.Key) {
					parent.Item.Keys.Right = EmptyKey()
					log.Debug("Set parent.RightKey() to nil", "parent.Key", parent.Item.Keys)
					tree.Save(parent, dryrun, blockHash)
				} else {
					log.Error("Parent doesn't have node as child", "node.Key", hex.EncodeToString(node.Key))
					return
				}
			} else {
				// this node is root
				tree.rootKey = EmptyKey()
			}
			tree.deleteNode(node, dryrun, blockHash)
			log.Debug("Removed node with child = nil", "node", hex.EncodeToString(node.Key), "all keys", tree.KeysinString(true, blockHash))
		} else {
			if node.Item.Color == black {
				node.Item.Color = nodeColor(child)
				tree.Save(node, dryrun, blockHash)

				tree.deleteCase1(node, dryrun, blockHash)
			}

			tree.replaceNode(node, child, dryrun, blockHash)

			if tree.IsEmptyKey(node.ParentKey()) && child != nil {
				child.Item.Color = black
				tree.Save(child, dryrun, blockHash)
			}
		}
	}

	tree.size--
	log.Debug("Deleted node", "node key", hex.EncodeToString(node.Key), "all keys", tree.KeysinString(true, blockHash), "tree.size", tree.size)
}

// // Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *Tree) Size() uint64 {
	return tree.size
}

// Keys returns all keys in-order
func (tree *Tree) Keys(dryrun bool, blockHash common.Hash) [][]byte {
	keys := make([][]byte, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(dryrun, blockHash) && i < len(keys); i++ {
		keys[i] = it.Key()
	}
	return keys
}

func (tree *Tree) KeysinString(dryrun bool, blockHash common.Hash) []string {
	keys := make([]string, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(dryrun, blockHash) && i < len(keys); i++ {
		keys[i] = hex.EncodeToString(it.Key())
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree) Values(dryrun bool, blockHash common.Hash) [][]byte {
	values := make([][]byte, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(dryrun, blockHash) && i < len(values); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left(dryrun bool, blockHash common.Hash) *Node {
	var parent *Node
	current := tree.Root(dryrun, blockHash)
	for current != nil {
		parent = current
		current = current.Left(tree, dryrun, blockHash)
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right(dryrun bool, blockHash common.Hash) *Node {
	var parent *Node
	current := tree.Root(dryrun, blockHash)
	for current != nil {
		parent = current
		current = current.Right(tree, dryrun, blockHash)
	}
	return parent
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Floor(key []byte, dryrun bool, blockHash common.Hash) (floor *Node, found bool) {
	found = false
	node := tree.Root(dryrun, blockHash)
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left(tree, dryrun, blockHash)
		case compare > 0:
			floor, found = node, true
			node = node.Right(tree, dryrun, blockHash)
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Ceiling(key []byte, dryrun bool, blockHash common.Hash) (ceiling *Node, found bool) {
	found = false
	node := tree.Root(dryrun, blockHash)
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left(tree, dryrun, blockHash)
		case compare > 0:
			node = node.Right(tree, dryrun, blockHash)
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

// Clear removes all nodes from the tree.
// we do not delete other children, but update them by overriding later
func (tree *Tree) Clear() {
	log.Debug("Remove all nodes from the tree", "tree", tree)
	tree.rootKey = EmptyKey()
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String(dryrun bool, blockHash common.Hash) string {
	str := fmt.Sprintf("RedBlackTree, size: %d\n", tree.size)

	// if !tree.Empty() {
	output(tree, tree.Root(dryrun, blockHash), "", true, &str, dryrun, blockHash)
	// }
	return str
}

func output(tree *Tree, node *Node, prefix string, isTail bool, str *string, dryrun bool, blockHash common.Hash) {
	// fmt.Printf("Node : %v+\n", node)
	if node == nil {
		return
	}

	if !tree.IsEmptyKey(node.RightKey()) {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(tree, node.Right(tree, dryrun, blockHash), newPrefix, false, str, dryrun, blockHash)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}

	if tree.FormatBytes != nil {
		*str += node.String(tree) + "\n"
	} else {
		*str += string(node.Key) + "\n"
	}

	if !tree.IsEmptyKey(node.LeftKey()) {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(tree, node.Left(tree, dryrun, blockHash), newPrefix, true, str, dryrun, blockHash)
	}
}

func (tree *Tree) rotateLeft(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Rotate left - before", "Root key", hex.EncodeToString(tree.Root(dryrun, blockHash).Key), "all keys", tree.KeysinString(true, blockHash))
	right := node.Right(tree, dryrun, blockHash)
	tree.replaceNode(node, right, dryrun, blockHash)
	node.RightKey(right.LeftKey())
	if !tree.IsEmptyKey(right.LeftKey()) {
		rightLeft := right.Left(tree, dryrun, blockHash)
		rightLeft.ParentKey(node.Key)
		tree.Save(rightLeft, dryrun, blockHash)
	}
	right.LeftKey(node.Key)
	node.ParentKey(right.Key)
	tree.Save(node, dryrun, blockHash)
	tree.Save(right, dryrun, blockHash)
	log.Debug("Rotate left - after", "Root key", hex.EncodeToString(tree.Root(dryrun, blockHash).Key), "all keys", tree.KeysinString(true, blockHash))
}

func (tree *Tree) rotateRight(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Rotate right - before", "Root key", hex.EncodeToString(tree.Root(dryrun, blockHash).Key), "all keys", tree.KeysinString(true, blockHash))
	left := node.Left(tree, dryrun, blockHash)
	tree.replaceNode(node, left, dryrun, blockHash)
	node.LeftKey(left.RightKey())
	if !tree.IsEmptyKey(left.RightKey()) {
		leftRight := left.Right(tree, dryrun, blockHash)
		leftRight.ParentKey(node.Key)
		tree.Save(leftRight, dryrun, blockHash)
	}
	left.RightKey(node.Key)
	node.ParentKey(left.Key)
	tree.Save(node, dryrun, blockHash)
	tree.Save(left, dryrun, blockHash)
	log.Debug("Rotate right - after", "Root key", hex.EncodeToString(tree.Root(dryrun, blockHash).Key), "all keys", tree.KeysinString(true, blockHash))
}

func (tree *Tree) replaceNode(old *Node, new *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Replace node", "old", old, "new", new)

	// we do not change any byte of Key so we can copy the reference to save directly to db
	var newKey []byte
	if new == nil {
		newKey = EmptyKey()
	} else {
		newKey = new.Key
	}

	if tree.IsEmptyKey(old.ParentKey()) {
		log.Debug("Set root key - replaceNode", "key", hex.EncodeToString(newKey), "size", tree.size, "all keys", tree.KeysinString(true, blockHash))
		tree.rootKey = newKey
	} else {
		// update left and right for oldParent
		oldParent := old.Parent(tree, dryrun, blockHash)
		if tree.Comparator(old.Key, oldParent.LeftKey()) == 0 {
			oldParent.LeftKey(newKey)
		} else {
			// remove oldParent right
			oldParent.RightKey(newKey)
		}
		// we can have case like: remove a node, then add it again
		if oldParent != nil {
			tree.Save(oldParent, dryrun, blockHash)
		}
	}
	if new != nil {
		// here is the swap, not update key
		// new.Parent = old.Parent
		new.ParentKey(old.ParentKey())
		tree.Save(new, dryrun, blockHash)
	}

}

func (tree *Tree) insertCase1(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Insert case 1", "node key", hex.EncodeToString(node.Key))
	if tree.IsEmptyKey(node.ParentKey()) {
		node.Item.Color = black
	} else {
		tree.insertCase2(node, dryrun, blockHash)
	}
}

func (tree *Tree) insertCase2(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Insert case 2", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	if nodeColor(parent) == black {
		return
	}

	tree.insertCase3(node, dryrun, blockHash)
}

func (tree *Tree) insertCase3(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Insert case 3", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	uncle := node.uncle(tree, dryrun, blockHash)
	if nodeColor(uncle) == red {
		parent.Item.Color = black
		uncle.Item.Color = black
		tree.Save(uncle, dryrun, blockHash)
		tree.Save(parent, dryrun, blockHash)
		grandparent := parent.Parent(tree, dryrun, blockHash)
		tree.assertNotNull(grandparent, "grant parent")

		grandparent.Item.Color = red
		tree.insertCase1(grandparent, dryrun, blockHash)
		tree.Save(grandparent, dryrun, blockHash)
	} else {
		tree.insertCase4(node, dryrun, blockHash)
	}
}

func (tree *Tree) insertCase4(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Insert case 4", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	grandparent := parent.Parent(tree, dryrun, blockHash)
	tree.assertNotNull(grandparent, "grant parent")
	if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.LeftKey()) == 0 {
		tree.rotateLeft(parent, dryrun, blockHash)
		node = node.Left(tree, dryrun, blockHash)
	} else if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.RightKey()) == 0 {
		tree.rotateRight(parent, dryrun, blockHash)
		node = node.Right(tree, dryrun, blockHash)
	}

	tree.insertCase5(node, dryrun, blockHash)
}

func (tree *Tree) assertNotNull(node *Node, name string) {
	if node == nil {
		panic(fmt.Sprintf("%s is nil\n", name))
	}
}

func (tree *Tree) insertCase5(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Insert case 5", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	parent.Item.Color = black
	tree.Save(parent, dryrun, blockHash)

	grandparent := parent.Parent(tree, dryrun, blockHash)
	tree.assertNotNull(grandparent, "grant parent")
	grandparent.Item.Color = red
	tree.Save(grandparent, dryrun, blockHash)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.LeftKey()) == 0 {
		tree.rotateRight(grandparent, dryrun, blockHash)
	} else if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.RightKey()) == 0 {
		tree.rotateLeft(grandparent, dryrun, blockHash)
	}

}

func (tree *Tree) Save(node *Node, dryrun bool, blockHash common.Hash) error {
	log.Debug("Save node", "node", node, "key", hex.EncodeToString(node.Key))
	return tree.db.Put(node.Key, node.Item, dryrun, blockHash)
}

func (tree *Tree) deleteCase1(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 1", "node key", hex.EncodeToString(node.Key))
	if tree.IsEmptyKey(node.ParentKey()) {
		//tree.deleteNode(node, dryrun, blockHash)
		return
	}

	tree.deleteCase2(node, dryrun, blockHash)
}

func (tree *Tree) deleteCase2(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 2", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	sibling := node.sibling(tree, dryrun, blockHash)

	if nodeColor(sibling) == red {
		parent.Item.Color = red
		sibling.Item.Color = black
		tree.Save(parent, dryrun, blockHash)
		tree.Save(sibling, dryrun, blockHash)
		if tree.Comparator(node.Key, parent.LeftKey()) == 0 {
			tree.rotateLeft(parent, dryrun, blockHash)
		} else {
			tree.rotateRight(parent, dryrun, blockHash)
		}
	}

	tree.deleteCase3(node, dryrun, blockHash)
}

func (tree *Tree) deleteCase3(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 3", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	sibling := node.sibling(tree, dryrun, blockHash)
	siblingLeft := sibling.Left(tree, dryrun, blockHash)
	siblingRight := sibling.Right(tree, dryrun, blockHash)

	if nodeColor(parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == black &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		tree.Save(sibling, dryrun, blockHash)
		tree.deleteCase1(parent, dryrun, blockHash)
		//tree.deleteNode(node, dryrun, blockHash)
	} else {
		tree.deleteCase4(node, dryrun, blockHash)
	}

}

func (tree *Tree) deleteCase4(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 4", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	sibling := node.sibling(tree, dryrun, blockHash)
	siblingLeft := sibling.Left(tree, dryrun, blockHash)
	siblingRight := sibling.Right(tree, dryrun, blockHash)

	if nodeColor(parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == black &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		parent.Item.Color = black
		tree.Save(sibling, dryrun, blockHash)
		tree.Save(parent, dryrun, blockHash)
	} else {
		tree.deleteCase5(node, dryrun, blockHash)
	}
}

func (tree *Tree) deleteCase5(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 5", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	sibling := node.sibling(tree, dryrun, blockHash)
	siblingLeft := sibling.Left(tree, dryrun, blockHash)
	siblingRight := sibling.Right(tree, dryrun, blockHash)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == red &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		siblingLeft.Item.Color = black

		tree.Save(sibling, dryrun, blockHash)
		tree.Save(siblingLeft, dryrun, blockHash)

		tree.rotateRight(sibling, dryrun, blockHash)

	} else if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		nodeColor(sibling) == black &&
		nodeColor(siblingRight) == red &&
		nodeColor(siblingLeft) == black {
		sibling.Item.Color = red
		siblingRight.Item.Color = black

		tree.Save(sibling, dryrun, blockHash)
		tree.Save(siblingRight, dryrun, blockHash)

		tree.rotateLeft(sibling, dryrun, blockHash)

	}

	tree.deleteCase6(node, dryrun, blockHash)
}

func (tree *Tree) deleteCase6(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("delete case 6", "node key", hex.EncodeToString(node.Key))
	parent := node.Parent(tree, dryrun, blockHash)
	sibling := node.sibling(tree, dryrun, blockHash)
	siblingLeft := sibling.Left(tree, dryrun, blockHash)
	siblingRight := sibling.Right(tree, dryrun, blockHash)

	sibling.Item.Color = nodeColor(parent)
	parent.Item.Color = black

	tree.Save(sibling, dryrun, blockHash)
	tree.Save(parent, dryrun, blockHash)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 && nodeColor(siblingRight) == red {
		siblingRight.Item.Color = black
		tree.Save(siblingRight, dryrun, blockHash)
		tree.rotateLeft(parent, dryrun, blockHash)
	} else if nodeColor(siblingLeft) == red {
		siblingLeft.Item.Color = black
		tree.Save(siblingLeft, dryrun, blockHash)
		tree.rotateRight(parent, dryrun, blockHash)
	}

	// update the parent meta then delete the current node from db
	//tree.deleteNode(node, dryrun, blockHash)
}

func nodeColor(node *Node) bool {
	if node == nil {
		return black
	}
	return node.Item.Color
}

func (tree *Tree) deleteNode(node *Node, dryrun bool, blockHash common.Hash) {
	log.Debug("Delete node", "node key", hex.EncodeToString(node.Key))
	tree.db.Delete(node.Key, dryrun, blockHash)
}
