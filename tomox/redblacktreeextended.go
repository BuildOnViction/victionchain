// Copyright (c) 2019, Agiletech Viet Nam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tomox

import (
	"math/big"
)

// RedBlackTreeExtended to demonstrate how to extend a RedBlackTree to include new functions
type RedBlackTreeExtended struct {
	*Tree
}

func NewRedBlackTreeExtended(obdb OrderDao) *RedBlackTreeExtended {
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
func (tree *RedBlackTreeExtended) GetMin(dryrun bool) (value []byte, found bool) {
	node, found := tree.getMinFromNode(tree.Root(dryrun), dryrun)
	if node != nil {
		return node.Value(), found
	}
	return nil, false
}

// GetMax gets the max value and flag if found
func (tree *RedBlackTreeExtended) GetMax(dryrun bool) (value []byte, found bool) {
	node, found := tree.getMaxFromNode(tree.Root(dryrun), dryrun)
	if node != nil {
		return node.Value(), found
	}
	return nil, false
}

// RemoveMin removes the min value and flag if found
func (tree *RedBlackTreeExtended) RemoveMin(dryrun bool) (value []byte, deleted bool) {
	node, found := tree.getMinFromNode(tree.Root(dryrun), dryrun)
	// fmt.Println("found min", node)
	if found {
		tree.Remove(node.Key, false)
		// fmt.Printf("%x\n", node.Key)
		return node.Value(), found
	}
	return nil, false
}

// RemoveMax removes the max value and flag if found
func (tree *RedBlackTreeExtended) RemoveMax(dryrun bool) (value []byte, deleted bool) {
	// fmt.Println("found max with root", tree.Root())
	node, found := tree.getMaxFromNode(tree.Root(dryrun), dryrun)
	// fmt.Println("found max", node)
	if found {
		tree.Remove(node.Key, false)
		return node.Value(), found
	}
	return nil, false
}

func (tree *RedBlackTreeExtended) getMinFromNode(node *Node, dryrun bool) (foundNode *Node, found bool) {
	if node == nil {
		return nil, false
	}
	nodeLeft := node.Left(tree.Tree, dryrun)
	if nodeLeft == nil {
		return node, true
	}
	return tree.getMinFromNode(nodeLeft, dryrun)
}

func (tree *RedBlackTreeExtended) getMaxFromNode(node *Node, dryrun bool) (foundNode *Node, found bool) {
	if node == nil {
		return nil, false
	}
	nodeRight := node.Right(tree.Tree, dryrun)
	if nodeRight == nil {
		return node, true
	}
	return tree.getMaxFromNode(nodeRight, dryrun)
}
