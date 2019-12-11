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

package lendingstate

import (
	"fmt"
	"github.com/tomochain/tomochain/rlp"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/trie"
)

type DumpOrderList struct {
	Volume *big.Int
	Orders map[*big.Int]*big.Int
}

func (self *TomoXStateDB) DumpAskTrie(orderBook common.Hash) (map[*big.Int]DumpOrderList, error) {
	exhangeObject := self.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("Order book not found orderBook : %v ", orderBook.Hex())
	}
	result := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getAsksTrie(self.db).NodeIterator(nil))
	for it.Next() {
		InterestByte := self.trie.GetKey(it.Key)
		Interest := new(big.Int).SetBytes(InterestByte)
		var data orderList
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			return nil, fmt.Errorf("Fail when decode order iist orderBook : %v ,Interest :%v ", orderBook.Hex(), Interest)
		}
		orderList := newStateOrderList(self, Ask, orderBook, common.BytesToHash(InterestByte), data, nil)
		dumpOrderList := DumpOrderList{Volume: data.Volume, Orders: map[*big.Int]*big.Int{}}
		orderListIt := trie.NewIterator(orderList.getTrie(self.db).NodeIterator(nil))
		for orderListIt.Next() {
			dumpOrderList.Orders[new(big.Int).SetBytes(self.trie.GetKey(orderListIt.Key))] = new(big.Int).SetBytes(orderListIt.Value)
		}
		result[Interest] = dumpOrderList
	}
	return result, nil
}

func (self *TomoXStateDB) DumpBidTrie(orderBook common.Hash) (map[*big.Int]DumpOrderList, error) {
	exhangeObject := self.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("Order book not found orderBook : %v ", orderBook.Hex())
	}
	result := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getBidsTrie(self.db).NodeIterator(nil))
	for it.Next() {
		InterestByte := self.trie.GetKey(it.Key)
		Interest := new(big.Int).SetBytes(InterestByte)
		var data orderList
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			return nil, fmt.Errorf("Fail when decode order iist orderBook : %v ,Interest :%v ", orderBook.Hex(), Interest)
		}
		orderList := newStateOrderList(self, Bid, orderBook, common.BytesToHash(InterestByte), data, nil)
		dumpOrderList := DumpOrderList{Volume: data.Volume, Orders: map[*big.Int]*big.Int{}}
		orderListIt := trie.NewIterator(orderList.getTrie(self.db).NodeIterator(nil))
		for orderListIt.Next() {
			dumpOrderList.Orders[new(big.Int).SetBytes(self.trie.GetKey(orderListIt.Key))] = new(big.Int).SetBytes(orderListIt.Value)
		}
		result[Interest] = dumpOrderList
	}
	return result, nil
}
