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

package lending_state

import (
	"github.com/tomochain/tomochain/common"
	"math/big"
)

type journalEntry interface {
	undo(db *LendingStateDB)
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
	priceChange struct {
		hash common.Hash
		prev *big.Int
	}
)

func (ch insertOrder) undo(s *LendingStateDB) {
	s.CancelLendingOrder(ch.orderBook, ch.order)
}
func (ch cancelOrder) undo(s *LendingStateDB) {
	s.InsertLendingItem(ch.orderBook, ch.orderId, ch.order)
}
func (ch subAmountOrder) undo(s *LendingStateDB) {
	priceHash := common.BigToHash(ch.order.Price)
	stateOrderBook := s.getLendingExchange(ch.orderBook)
	var stateOrderList *itemListState
	switch ch.order.Side {
	case INVESTING:
		stateOrderList = stateOrderBook.getInvestingOrderList(s.db, priceHash)
	case BORROWING:
		stateOrderList = stateOrderBook.getBorrowingOrderList(s.db, priceHash)
	default:
		return
	}
	stateOrderItem := stateOrderBook.getLendingItem(s.db, ch.orderId)
	newAmount := new(big.Int).Add(stateOrderItem.Quantity(), ch.amount)
	stateOrderItem.setVolume(newAmount)
	stateOrderList.insertLendingItem(s.db, ch.orderId, common.BigToHash(newAmount))
}
func (ch nonceChange) undo(s *LendingStateDB) {
	s.SetNonce(ch.hash, ch.prev)
}
