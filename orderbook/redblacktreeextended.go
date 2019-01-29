// Copyright (c) 2019, Agiletech Viet Nam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orderbook

import (
	"math/big"
)

// RedBlackTreeExtended to demonstrate how to extend a RedBlackTree to include new functions
type RedBlackTreeExtended struct {
	*Tree
}

func NewRedBlackTreeExtended(obdb *BatchDatabase) *RedBlackTreeExtended {
	tree := &RedBlackTreeExtended{NewWith(CmpBigInt, obdb)}

	tree.FormatBytes = func(key []byte) string {
		if tree.IsEmptyKey(key) {
			return "<nil>"
		}
		return new(big.Int).SetBytes(key).String()
	}

	return tree
}

// GetMin gets the min value and flag if found
func (tree *RedBlackTreeExtended) GetMin() (value []byte, found bool) {
	node, found := tree.getMinFromNode(tree.Root())
	if node != nil {
		return node.Value(), found
	}
	return nil, false
}

// GetMax gets the max value and flag if found
func (tree *RedBlackTreeExtended) GetMax() (value []byte, found bool) {
	node, found := tree.getMaxFromNode(tree.Root())
	if node != nil {
		return node.Value(), found
	}
	return nil, false
}

// RemoveMin removes the min value and flag if found
func (tree *RedBlackTreeExtended) RemoveMin() (value []byte, deleted bool) {
	node, found := tree.getMinFromNode(tree.Root())
	// fmt.Println("found min", node)
	if found {
		tree.Remove(node.Key)
		// fmt.Printf("%x\n", node.Key)
		return node.Value(), found
	}
	return nil, false
}

// RemoveMax removes the max value and flag if found
func (tree *RedBlackTreeExtended) RemoveMax() (value []byte, deleted bool) {
	// fmt.Println("found max with root", tree.Root())
	node, found := tree.getMaxFromNode(tree.Root())
	// fmt.Println("found max", node)
	if found {
		tree.Remove(node.Key)
		return node.Value(), found
	}
	return nil, false
}

func (tree *RedBlackTreeExtended) getMinFromNode(node *Node) (foundNode *Node, found bool) {
	if node == nil {
		return nil, false
	}
	nodeLeft := node.Left(tree.Tree)
	if nodeLeft == nil {
		return node, true
	}
	return tree.getMinFromNode(nodeLeft)
}

func (tree *RedBlackTreeExtended) getMaxFromNode(node *Node) (foundNode *Node, found bool) {
	if node == nil {
		return nil, false
	}
	nodeRight := node.Right(tree.Tree)
	if nodeRight == nil {
		return node, true
	}
	return tree.getMaxFromNode(nodeRight)
}
