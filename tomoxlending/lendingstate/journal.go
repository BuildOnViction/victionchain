// Copyright 2016 The go-ethereum Authors
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
	"github.com/tomochain/tomochain/common"
	"math/big"
)

type journalEntry interface {
	undo(db *TomoXStateDB)
}

type journal []journalEntry

type (
	// Changes to the account trie.
	insertOrder struct {
		orderBook common.Hash
		orderId   common.Hash
		order     *LendingItem
	}
	cancelOrder struct {
		orderBook common.Hash
		orderId   common.Hash
		order     LendingItem
	}
	subAmountOrder struct {
		orderBook common.Hash
		orderId   common.Hash
		order     LendingItem
		amount    *big.Int
	}
	nonceChange struct {
		hash common.Hash
		prev uint64
	}
	InterestChange struct {
		hash common.Hash
		prev *big.Int
	}
)

func (ch insertOrder) undo(s *TomoXStateDB) {
	s.CancelOrder(ch.orderBook, ch.order)
}
func (ch cancelOrder) undo(s *TomoXStateDB) {
	s.InsertLendingItem(ch.orderBook, ch.orderId, ch.order)
}
func (ch subAmountOrder) undo(s *TomoXStateDB) {
	InterestHash := common.BigToHash(ch.order.Interest)
	stateOrderBook := s.getStateExchangeObject(ch.orderBook)
	var stateOrderList *stateOrderList
	switch ch.order.Side {
	case Ask:
		stateOrderList = stateOrderBook.getStateOrderListAskObject(s.db, InterestHash)
	case Bid:
		stateOrderList = stateOrderBook.getStateBidOrderListObject(s.db, InterestHash)
	default:
		return
	}
	stateLendingItem := stateOrderBook.getStateOrderObject(s.db, ch.orderId)
	newAmount := new(big.Int).Add(stateLendingItem.Quantity(), ch.amount)
	stateLendingItem.setVolume(newAmount)
	stateOrderList.insertLendingItem(s.db, ch.orderId, common.BigToHash(newAmount))
	stateOrderList.AddVolume(ch.amount)
}
func (ch nonceChange) undo(s *TomoXStateDB) {
	s.SetNonce(ch.hash, ch.prev)
}
func (ch InterestChange) undo(s *TomoXStateDB) {
	s.SetInterest(ch.hash, ch.prev)
}
