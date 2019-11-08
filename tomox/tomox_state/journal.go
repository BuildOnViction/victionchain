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

package tomox_state

import (
	"github.com/ethereum/go-ethereum/common"
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
		order     *OrderItem
	}
	cancelOrder struct {
		orderBook common.Hash
		orderId   common.Hash
		order     OrderItem
	}
	subAmountOrder struct {
		orderBook common.Hash
		orderId   common.Hash
		order     OrderItem
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

func (ch insertOrder) undo(s *TomoXStateDB) {
	s.CancelOrder(ch.orderBook, ch.order)
}
func (ch cancelOrder) undo(s *TomoXStateDB) {
	s.InsertOrderItem(ch.orderBook, ch.orderId, ch.order)
}
func (ch subAmountOrder) undo(s *TomoXStateDB) {
	priceHash := common.BigToHash(ch.order.Price)
	stateOrderBook := s.getStateExchangeObject(ch.orderBook)
	var stateOrderList *stateOrderList
	switch ch.order.Side {
	case Ask:
		stateOrderList = stateOrderBook.getStateOrderListAskObject(s.db, priceHash)
	case Bid:
		stateOrderList = stateOrderBook.getStateBidOrderListObject(s.db, priceHash)
	default:
		return
	}
	stateOrderItem := stateOrderBook.getStateOrderObject(s.db, ch.orderId)
	newAmount := new(big.Int).Add(stateOrderItem.Quantity(), ch.amount)
	stateOrderItem.setVolume(newAmount)
	stateOrderList.insertOrderItem(s.db, ch.orderId, common.BigToHash(newAmount))
	stateOrderList.AddVolume(ch.amount)
}
func (ch nonceChange) undo(s *TomoXStateDB) {
	s.SetNonce(ch.hash, ch.prev)
}
func (ch priceChange) undo(s *TomoXStateDB) {
	s.SetPrice(ch.hash, ch.prev)
}
