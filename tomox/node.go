// Copyright (c) 2019, Agiletech Viet Nam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tomox

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

const (
	black, red bool = true, false
)

type KeyMeta struct {
	Left   []byte
	Right  []byte
	Parent []byte
}

type KeyMetaBSON struct {
	Left   string
	Right  string
	Parent string
}

func (keys *KeyMeta) String(tree *Tree) string {
	return fmt.Sprintf("L: %v, P: %v, R: %v", tree.FormatBytes(keys.Left),
		tree.FormatBytes(keys.Parent), tree.FormatBytes(keys.Right))
}

type Item struct {
	Keys  *KeyMeta
	Value []byte
	// Deleted bool
	Color bool
}

type ItemBSON struct {
	Keys  *KeyMetaBSON
	Value string
	// Deleted bool
	Color bool
}

type ItemRecord struct {
	Key   string
	Value *Item
}

type ItemRecordBSON struct {
	Key   string
	Value *ItemBSON
}

// Node is a single element within the tree
type Node struct {
	Key  []byte
	Item *Item
}

func (node *Node) String(tree *Tree) string {
	if node == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v -> %x, (%v)\n", tree.FormatBytes(node.Key), node.Value(), node.Item.Keys.String(tree))
}

func (node *Node) maximumNode(tree *Tree) *Node {
	newNode := node
	if newNode == nil {
		return newNode
	}
	for !tree.IsEmptyKey(node.RightKey()) {
		newNode = newNode.Right(tree)
	}
	return newNode
}

func (node *Node) LeftKey(keys ...[]byte) []byte {
	if node == nil || node.Item == nil || node.Item.Keys == nil {
		return nil
	}
	if len(keys) == 1 {
		node.Item.Keys.Left = keys[0]
	}

	return node.Item.Keys.Left
}

func (node *Node) RightKey(keys ...[]byte) []byte {
	if node == nil || node.Item == nil || node.Item.Keys == nil {
		return nil
	}
	if len(keys) == 1 {
		// if string(node.Key) == "1" {
		// 	fmt.Printf("Update right key: %s\n", string(keys[0]))
		// if string(keys[0]) == "3" {
		// 	panic("should stops")
		// }
		// }
		node.Item.Keys.Right = keys[0]
	}

	return node.Item.Keys.Right
}

func (node *Node) ParentKey(keys ...[]byte) []byte {
	if node == nil || node.Item == nil || node.Item.Keys == nil {
		return nil
	}
	if len(keys) == 1 {
		node.Item.Keys.Parent = keys[0]
	}

	return node.Item.Keys.Parent
}

func (node *Node) Left(tree *Tree) *Node {
	key := node.LeftKey()

	newNode, err := tree.GetNode(key)
	if err != nil {
		log.Error("Error at left", "err", err)
	}

	return newNode
}

func (node *Node) Right(tree *Tree) *Node {
	key := node.RightKey()
	newNode, err := tree.GetNode(key)
	if err != nil {
		fmt.Println(err)
	}
	return newNode
}

func (node *Node) Parent(tree *Tree) *Node {
	key := node.ParentKey()
	newNode, err := tree.GetNode(key)
	if err != nil {
		log.Error("Error at parent", "err", err)
	}
	return newNode
}

func (node *Node) Value() []byte {
	return node.Item.Value
}

func (node *Node) grandparent(tree *Tree) *Node {
	if node != nil && !tree.IsEmptyKey(node.ParentKey()) {
		return node.Parent(tree).Parent(tree)
	}
	return nil
}

func (node *Node) uncle(tree *Tree) *Node {
	if node == nil || tree.IsEmptyKey(node.ParentKey()) {
		return nil
	}
	parent := node.Parent(tree)
	// if tree.IsEmptyKey(parent.ParentKey()) {
	// 	return nil
	// }

	return parent.sibling(tree)
}

func (node *Node) sibling(tree *Tree) *Node {
	if node == nil || tree.IsEmptyKey(node.ParentKey()) {
		return nil
	}
	parent := node.Parent(tree)
	// fmt.Printf("Parent: %s\n", parent)
	if tree.Comparator(node.Key, parent.LeftKey()) == 0 {
		return parent.Right(tree)
	}
	return parent.Left(tree)
}
