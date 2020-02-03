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

package tomox_state

import (
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/trie"
	"math/big"
	"sort"
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
	mapResult := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getAsksTrie(self.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(self.trie.GetKey(it.Key))
		if common.EmptyHash(priceHash) {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateAskObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("Fail when decode order iist orderBook : %v ,price :%v ", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(self, Ask, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.DumpOrderList(self.db)
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateAskObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.DumpOrderList(self.db)
		}
	}
	listPrice := []*big.Int{}
	for price, _ := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]DumpOrderList{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return result, nil
}

func (self *TomoXStateDB) DumpBidTrie(orderBook common.Hash) (map[*big.Int]DumpOrderList, error) {
	exhangeObject := self.getStateExchangeObject(orderBook)
	if exhangeObject == nil {
		return nil, fmt.Errorf("Order book not found orderBook : %v ", orderBook.Hex())
	}
	mapResult := map[*big.Int]DumpOrderList{}
	it := trie.NewIterator(exhangeObject.getBidsTrie(self.db).NodeIterator(nil))
	for it.Next() {
		priceHash := common.BytesToHash(self.trie.GetKey(it.Key))
		if common.EmptyHash(priceHash) {
			continue
		}
		price := new(big.Int).SetBytes(priceHash.Bytes())
		if _, exist := exhangeObject.stateBidObjects[priceHash]; exist {
			continue
		} else {
			var data orderList
			if err := rlp.DecodeBytes(it.Value, &data); err != nil {
				return nil, fmt.Errorf("Fail when decode order iist orderBook : %v ,price :%v ", orderBook.Hex(), price)
			}
			stateOrderList := newStateOrderList(self, Bid, orderBook, priceHash, data, nil)
			mapResult[price] = stateOrderList.DumpOrderList(self.db)
		}
	}
	for priceHash, stateOrderList := range exhangeObject.stateBidObjects {
		if stateOrderList.Volume().Sign() > 0 {
			mapResult[new(big.Int).SetBytes(priceHash.Bytes())] = stateOrderList.DumpOrderList(self.db)
		}
	}
	listPrice := []*big.Int{}
	for price, _ := range mapResult {
		listPrice = append(listPrice, price)
	}
	sort.Slice(listPrice, func(i, j int) bool {
		return listPrice[i].Cmp(listPrice[j]) < 0
	})
	result := map[*big.Int]DumpOrderList{}
	for _, price := range listPrice {
		result[price] = mapResult[price]
	}
	return mapResult, nil
}

func (self *stateOrderList) DumpOrderList(db Database) DumpOrderList {
	mapResult := DumpOrderList{Volume: self.Volume(), Orders: map[*big.Int]*big.Int{}}
	orderListIt := trie.NewIterator(self.getTrie(db).NodeIterator(nil))
	for orderListIt.Next() {
		keyHash := common.BytesToHash(self.trie.GetKey(orderListIt.Key))
		if common.EmptyHash(keyHash) {
			continue
		}
		if _, exist := self.cachedStorage[keyHash]; exist {
			continue
		} else {
			mapResult.Orders[new(big.Int).SetBytes(keyHash.Bytes())] = new(big.Int).SetBytes(orderListIt.Value)
		}
	}
	for key, value := range self.cachedStorage {
		if !common.EmptyHash(value) {
			mapResult.Orders[new(big.Int).SetBytes(key.Bytes())] = new(big.Int).SetBytes(value.Bytes())
		}
	}
	listIds := []*big.Int{}
	for id, _ := range mapResult.Orders {
		listIds = append(listIds, id)
	}
	sort.Slice(listIds, func(i, j int) bool {
		return listIds[i].Cmp(listIds[j]) < 0
	})
	result := DumpOrderList{Volume: self.Volume(), Orders: map[*big.Int]*big.Int{}}
	for _, id := range listIds {
		result.Orders[id] = mapResult.Orders[id]
	}
	return mapResult
}
