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

package trading_state

import (
	"github.com/tomochain/tomochain/common"
	"math/big"
)

type journalEntry interface {
	undo(db *TradingStateDB)
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

func (ch insertOrder) undo(s *TradingStateDB) {
	s.CancelOrder(ch.orderBook, ch.order)
}
func (ch cancelOrder) undo(s *TradingStateDB) {
	s.InsertOrderItem(ch.orderBook, ch.orderId, ch.order)
}
func (ch subAmountOrder) undo(s *TradingStateDB) {
	priceHash := common.BigToHash(ch.order.Price)
	stateOrderBook := s.getStateExchangeObject(ch.orderBook)
	var stateOrderList *orderListState
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
}
func (ch nonceChange) undo(s *TradingStateDB) {
	s.SetNonce(ch.hash, ch.prev)
}
func (ch priceChange) undo(s *TradingStateDB) {
	s.SetPrice(ch.hash, ch.prev)
}
