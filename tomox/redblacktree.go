// Copyright (c) 2019, Agiletech Viet Nam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// References: http://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package tomox

import (
	"bytes"
	"fmt"

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

func (tree *Tree) Root(dryrun bool) *Node {
	root, err := tree.GetNode(tree.rootKey, dryrun)
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
func (tree *Tree) Put(key []byte, value []byte, dryrun bool) error {
	var insertedNode *Node
	if tree.IsEmptyKey(tree.rootKey) {
		// Assert key is of comparator's type for initial tree
		item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
		tree.rootKey = key
		insertedNode = &Node{Key: key, Item: item}
	} else {
		node := tree.Root(dryrun)
		if node == nil {
			return fmt.Errorf("Error on inserting node into the tree. tree.Root() is nil")
		}
		loop := true

		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:

				node.Item.Value = value
				tree.Save(node, dryrun)
				return nil
			case compare < 0:
				if tree.IsEmptyKey(node.LeftKey()) {
					node.LeftKey(key)
					tree.Save(node, dryrun)
					item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
					nodeLeft := &Node{Key: key, Item: item}
					insertedNode = nodeLeft
					loop = false
				} else {
					node = node.Left(tree, dryrun)
				}
			case compare > 0:

				if tree.IsEmptyKey(node.RightKey()) {
					node.RightKey(key)
					tree.Save(node, dryrun)
					item := &Item{Value: value, Color: red, Keys: &KeyMeta{}}
					nodeRight := &Node{Key: key, Item: item}
					insertedNode = nodeRight
					loop = false
				} else {
					node = node.Right(tree, dryrun)
				}

			}
		}

		insertedNode.ParentKey(node.Key)
		tree.Save(insertedNode, dryrun)
	}

	tree.insertCase1(insertedNode, dryrun)
	tree.Save(insertedNode, dryrun)

	tree.size++
	return nil
}

func (tree *Tree) GetNode(key []byte, dryrun bool) (*Node, error) {

	item := &Item{}

	val, err := tree.db.Get(key, item, dryrun)

	if err != nil || val == nil {
		return nil, err
	}
	return &Node{Key: key, Item: val.(*Item)}, err
}

func (tree *Tree) Has(key []byte, dryrun bool) (bool, error) {
	return tree.db.Has(key, dryrun)
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key []byte, dryrun bool) (value []byte, found bool) {

	node, err := tree.GetNode(key, dryrun)
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
func (tree *Tree) Remove(key []byte, dryrun bool) {
	var child *Node
	node, err := tree.GetNode(key, dryrun)
	if err != nil || node == nil {
		return
	}
	log.Debug("Get node", "node", node, "dryrun", dryrun)

	var left, right *Node = nil, nil
	if !tree.IsEmptyKey(node.LeftKey()) {
		left = node.Left(tree, dryrun)
	}
	if !tree.IsEmptyKey(node.RightKey()) {
		right = node.Right(tree, dryrun)
	}

	if left != nil && right != nil {
		node = left.maximumNode(tree, dryrun)
	}

	if left == nil || right == nil {
		if right == nil {
			child = left
		} else {
			child = right
		}
		if child == nil {
			tree.deleteNode(node, dryrun)
		} else {
			if node.Item.Color == black {
				node.Item.Color = nodeColor(child)
				tree.Save(node, dryrun)

				tree.deleteCase1(node, dryrun)
			}

			tree.replaceNode(node, child, dryrun)

			if tree.IsEmptyKey(node.ParentKey()) && child != nil {
				child.Item.Color = black
				tree.Save(child, dryrun)
			}
		}
	}

	tree.size--
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
func (tree *Tree) Keys(dryrun bool) [][]byte {
	keys := make([][]byte, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(dryrun); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree) Values(dryrun bool) [][]byte {
	values := make([][]byte, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(dryrun); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left(dryrun bool) *Node {
	var parent *Node
	current := tree.Root(dryrun)
	for current != nil {
		parent = current
		current = current.Left(tree, dryrun)
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right(dryrun bool) *Node {
	var parent *Node
	current := tree.Root(dryrun)
	for current != nil {
		parent = current
		current = current.Right(tree, dryrun)
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
func (tree *Tree) Floor(key []byte, dryrun bool) (floor *Node, found bool) {
	found = false
	node := tree.Root(dryrun)
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left(tree, dryrun)
		case compare > 0:
			floor, found = node, true
			node = node.Right(tree, dryrun)
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
func (tree *Tree) Ceiling(key []byte, dryrun bool) (ceiling *Node, found bool) {
	found = false
	node := tree.Root(dryrun)
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left(tree, dryrun)
		case compare > 0:
			node = node.Right(tree, dryrun)
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
	tree.rootKey = EmptyKey()
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String(dryrun bool) string {
	str := fmt.Sprintf("RedBlackTree, size: %d\n", tree.size)

	// if !tree.Empty() {
	output(tree, tree.Root(dryrun), "", true, &str, dryrun)
	// }
	return str
}

func output(tree *Tree, node *Node, prefix string, isTail bool, str *string, dryrun bool) {
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
		output(tree, node.Right(tree, dryrun), newPrefix, false, str, dryrun)
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
		output(tree, node.Left(tree, dryrun), newPrefix, true, str, dryrun)
	}
}

func (tree *Tree) rotateLeft(node *Node, dryrun bool) {
	right := node.Right(tree, dryrun)
	tree.replaceNode(node, right, dryrun)
	node.RightKey(right.LeftKey())
	if !tree.IsEmptyKey(right.LeftKey()) {
		rightLeft := right.Left(tree, dryrun)
		rightLeft.ParentKey(node.Key)
		tree.Save(rightLeft, dryrun)
	}
	right.LeftKey(node.Key)
	node.ParentKey(right.Key)
	tree.Save(node, dryrun)
	tree.Save(right, dryrun)
}

func (tree *Tree) rotateRight(node *Node, dryrun bool) {
	left := node.Left(tree, dryrun)
	tree.replaceNode(node, left, dryrun)
	node.LeftKey(left.RightKey())
	if !tree.IsEmptyKey(left.RightKey()) {
		leftRight := left.Right(tree, dryrun)
		leftRight.ParentKey(node.Key)
		tree.Save(leftRight, dryrun)
	}
	left.RightKey(node.Key)
	node.ParentKey(left.Key)
	tree.Save(node, dryrun)
	tree.Save(left, dryrun)
}

func (tree *Tree) replaceNode(old *Node, new *Node, dryrun bool) {
	log.Debug("Replace node", "old", old, "new", new, "dryrun", dryrun)

	// we do not change any byte of Key so we can copy the reference to save directly to db
	var newKey []byte
	if new == nil {
		newKey = EmptyKey()
	} else {
		newKey = new.Key
	}

	if tree.IsEmptyKey(old.ParentKey()) {
		tree.rootKey = newKey
	} else {
		// update left and right for oldParent
		oldParent := old.Parent(tree, dryrun)
		if tree.Comparator(old.Key, oldParent.LeftKey()) == 0 {
			oldParent.LeftKey(newKey)
		} else {
			// remove oldParent right
			oldParent.RightKey(newKey)
		}
		// we can have case like: remove a node, then add it again
		if oldParent != nil {
			tree.Save(oldParent, dryrun)
		}
	}
	if new != nil {
		// here is the swap, not update key
		// new.Parent = old.Parent
		new.ParentKey(old.ParentKey())
		tree.Save(new, dryrun)
	}

}

func (tree *Tree) insertCase1(node *Node, dryrun bool) {

	if tree.IsEmptyKey(node.ParentKey()) {
		node.Item.Color = black
	} else {
		tree.insertCase2(node, dryrun)
	}
}

func (tree *Tree) insertCase2(node *Node, dryrun bool) {
	parent := node.Parent(tree, dryrun)
	if nodeColor(parent) == black {
		return
	}

	tree.insertCase3(node, dryrun)
}

func (tree *Tree) insertCase3(node *Node, dryrun bool) {
	parent := node.Parent(tree, dryrun)
	uncle := node.uncle(tree, dryrun)
	if nodeColor(uncle) == red {
		parent.Item.Color = black
		uncle.Item.Color = black
		tree.Save(uncle, dryrun)
		tree.Save(parent, dryrun)
		grandparent := parent.Parent(tree, dryrun)
		tree.assertNotNull(grandparent, "grant parent")

		grandparent.Item.Color = red
		tree.insertCase1(grandparent, dryrun)
		tree.Save(grandparent, dryrun)
	} else {
		tree.insertCase4(node, dryrun)
	}
}

func (tree *Tree) insertCase4(node *Node, dryrun bool) {
	parent := node.Parent(tree, dryrun)
	grandparent := parent.Parent(tree, dryrun)
	tree.assertNotNull(grandparent, "grant parent")
	if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.LeftKey()) == 0 {
		tree.rotateLeft(parent, dryrun)
		node = node.Left(tree, dryrun)
	} else if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.RightKey()) == 0 {
		tree.rotateRight(parent, dryrun)
		node = node.Right(tree, dryrun)
	}

	tree.insertCase5(node, dryrun)
}

func (tree *Tree) assertNotNull(node *Node, name string) {
	if node == nil {
		panic(fmt.Sprintf("%s is nil\n", name))
	}
}

func (tree *Tree) insertCase5(node *Node, dryrun bool) {
	parent := node.Parent(tree, dryrun)
	parent.Item.Color = black
	tree.Save(parent, dryrun)

	grandparent := parent.Parent(tree, dryrun)
	tree.assertNotNull(grandparent, "grant parent")
	grandparent.Item.Color = red
	tree.Save(grandparent, dryrun)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.LeftKey()) == 0 {
		tree.rotateRight(grandparent, dryrun)
	} else if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		tree.Comparator(parent.Key, grandparent.RightKey()) == 0 {
		tree.rotateLeft(grandparent, dryrun)
	}

}

func (tree *Tree) Save(node *Node, dryrun bool) error {
	log.Debug("Save node", "node", node, "dryrun", dryrun)
	return tree.db.Put(node.Key, node.Item, dryrun)
}

func (tree *Tree) deleteCase1(node *Node, dryrun bool) {
	log.Debug("delete case 1", "node value", hex.EncodeToString(node.Value()))
	if tree.IsEmptyKey(node.ParentKey()) {
		tree.deleteNode(node, dryrun)
		return
	}

	tree.deleteCase2(node, dryrun)
}

func (tree *Tree) deleteCase2(node *Node, dryrun bool) {
	log.Debug("delete case 2", "node value", hex.EncodeToString(node.Value()))
	parent := node.Parent(tree, dryrun)
	sibling := node.sibling(tree, dryrun)

	if nodeColor(sibling) == red {
		parent.Item.Color = red
		sibling.Item.Color = black
		tree.Save(parent, dryrun)
		tree.Save(sibling, dryrun)
		if tree.Comparator(node.Key, parent.LeftKey()) == 0 {
			tree.rotateLeft(parent, dryrun)
		} else {
			tree.rotateRight(parent, dryrun)
		}
	}

	tree.deleteCase3(node, dryrun)
}

func (tree *Tree) deleteCase3(node *Node, dryrun bool) {
	log.Debug("delete case 3", "node value", hex.EncodeToString(node.Value()))
	parent := node.Parent(tree, dryrun)
	sibling := node.sibling(tree, dryrun)
	siblingLeft := sibling.Left(tree, dryrun)
	siblingRight := sibling.Right(tree, dryrun)

	if nodeColor(parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == black &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		tree.Save(sibling, dryrun)
		tree.deleteCase1(parent, dryrun)
		tree.deleteNode(node, dryrun)
	} else {
		tree.deleteCase4(node, dryrun)
	}

}

func (tree *Tree) deleteCase4(node *Node, dryrun bool) {
	log.Debug("delete case 4", "node value", hex.EncodeToString(node.Value()))
	parent := node.Parent(tree, dryrun)
	sibling := node.sibling(tree, dryrun)
	siblingLeft := sibling.Left(tree, dryrun)
	siblingRight := sibling.Right(tree, dryrun)

	if nodeColor(parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == black &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		parent.Item.Color = black
		tree.Save(sibling, dryrun)
		tree.Save(parent, dryrun)
	} else {
		tree.deleteCase5(node, dryrun)
	}
}

func (tree *Tree) deleteCase5(node *Node, dryrun bool) {
	log.Debug("delete case 5", "node value", hex.EncodeToString(node.Value()))
	parent := node.Parent(tree, dryrun)
	sibling := node.sibling(tree, dryrun)
	siblingLeft := sibling.Left(tree, dryrun)
	siblingRight := sibling.Right(tree, dryrun)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 &&
		nodeColor(sibling) == black &&
		nodeColor(siblingLeft) == red &&
		nodeColor(siblingRight) == black {
		sibling.Item.Color = red
		siblingLeft.Item.Color = black

		tree.Save(sibling, dryrun)
		tree.Save(siblingLeft, dryrun)

		tree.rotateRight(sibling, dryrun)

	} else if tree.Comparator(node.Key, parent.RightKey()) == 0 &&
		nodeColor(sibling) == black &&
		nodeColor(siblingRight) == red &&
		nodeColor(siblingLeft) == black {
		sibling.Item.Color = red
		siblingRight.Item.Color = black

		tree.Save(sibling, dryrun)
		tree.Save(siblingRight, dryrun)

		tree.rotateLeft(sibling, dryrun)

	}

	tree.deleteCase6(node, dryrun)
}

func (tree *Tree) deleteCase6(node *Node, dryrun bool) {
	log.Debug("delete case 6", "node value", hex.EncodeToString(node.Value()))
	parent := node.Parent(tree, dryrun)
	sibling := node.sibling(tree, dryrun)
	siblingLeft := sibling.Left(tree, dryrun)
	siblingRight := sibling.Right(tree, dryrun)

	sibling.Item.Color = nodeColor(parent)
	parent.Item.Color = black

	tree.Save(sibling, dryrun)
	tree.Save(parent, dryrun)

	if tree.Comparator(node.Key, parent.LeftKey()) == 0 && nodeColor(siblingRight) == red {
		siblingRight.Item.Color = black
		tree.Save(siblingRight, dryrun)
		tree.rotateLeft(parent, dryrun)
	} else if nodeColor(siblingLeft) == red {
		siblingLeft.Item.Color = black
		tree.Save(siblingLeft, dryrun)
		tree.rotateRight(parent, dryrun)
	}

	// update the parent meta then delete the current node from db
	tree.deleteNode(node, dryrun)
}

func nodeColor(node *Node) bool {
	if node == nil {
		return black
	}
	return node.Item.Color
}

func (tree *Tree) deleteNode(node *Node, dryrun bool) {
	log.Debug("Delete node", "node value", hex.EncodeToString(node.Value()), "dryrun", dryrun)
	tree.db.Delete(node.Key, dryrun)
}
