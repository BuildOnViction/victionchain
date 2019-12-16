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

package trading_state

import (
	"fmt"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/log"
	"github.com/tomochain/tomochain/rlp"
	"github.com/tomochain/tomochain/tomox/database"
	"io"
	"math/big"
)

// stateObject represents an Ethereum orderId which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// tradingObject values can be accessed and modified through the object.
// Finally, call CommitAskTrie to write the modified storage trie into a database.
type tradingState struct {
	orderBookHash common.Hash // orderbookHashprice of ethereum address of the orderId
	data          tradingObject
	db            *TradingStateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by TradingStateDB.Commit.
	dbErr error

	// Write caches.
	asksTrie             database.Trie // storage trie, which becomes non-nil on first access
	bidsTrie             database.Trie // storage trie, which becomes non-nil on first access
	ordersTrie           database.Trie // storage trie, which becomes non-nil on first access
	liquidationPriceTrie database.Trie // storage trie, which becomes non-nil on first access

	askStates      map[common.Hash]*orderListState
	askStatesDirty map[common.Hash]struct{}

	bidStates      map[common.Hash]*orderListState
	bidStatesDirty map[common.Hash]struct{}

	orderStates      map[common.Hash]*orderItemState
	orderStatesDirty map[common.Hash]struct{}

	liquidationPriceStates      map[common.Hash]*liquidationPriceState
	liquidationPriceStatesDirty map[common.Hash]struct{}
	onDirty                     func(hash common.Hash) // Callback method to mark a state object newly dirty
}

// empty returns whether the orderId is considered empty.
func (s *tradingState) empty() bool {
	return s.data.Nonce == 0 && common.EmptyHash(s.data.AskRoot) && common.EmptyHash(s.data.BidRoot) && common.EmptyHash(s.data.LiquidationPriceRoot)
}

// newObject creates a state object.
func newStateExchanges(db *TradingStateDB, hash common.Hash, data tradingObject, onDirty func(addr common.Hash)) *tradingState {
	return &tradingState{
		db:               db,
		orderBookHash:    hash,
		data:             data,
		askStates:        make(map[common.Hash]*orderListState),
		bidStates:        make(map[common.Hash]*orderListState),
		orderStates:      make(map[common.Hash]*orderItemState),
		askStatesDirty:   make(map[common.Hash]struct{}),
		bidStatesDirty:   make(map[common.Hash]struct{}),
		orderStatesDirty: make(map[common.Hash]struct{}),
		onDirty:          onDirty,
	}
}

// EncodeRLP implements rlp.Encoder.
func (self *tradingState) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, self.data)
}

// setError remembers the first non-nil error it is called with.
func (self *tradingState) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *tradingState) getAsksTrie(db database.Database) database.Trie {
	if self.asksTrie == nil {
		var err error
		self.asksTrie, err = db.OpenStorageTrie(self.orderBookHash, self.data.AskRoot)
		if err != nil {
			self.asksTrie, _ = db.OpenStorageTrie(self.orderBookHash, EmptyHash)
			self.setError(fmt.Errorf("can't create asks trie: %v", err))
		}
	}
	return self.asksTrie
}

func (self *tradingState) getOrdersTrie(db database.Database) database.Trie {
	if self.ordersTrie == nil {
		var err error
		self.ordersTrie, err = db.OpenStorageTrie(self.orderBookHash, self.data.OrderRoot)
		if err != nil {
			self.ordersTrie, _ = db.OpenStorageTrie(self.orderBookHash, EmptyHash)
			self.setError(fmt.Errorf("can't create asks trie: %v", err))
		}
	}
	return self.ordersTrie
}

func (self *tradingState) getBestPriceAsksTrie(db database.Database) common.Hash {
	trie := self.getAsksTrie(db)
	encKey, encValue, err := trie.TryGetBestLeftKeyAndValue()
	if err != nil {
		log.Error("Failed find best liquidationPrice ask trie ", "orderbook", self.orderBookHash.Hex())
		return EmptyHash
	}
	if len(encKey) == 0 || len(encValue) == 0 {
		log.Debug("Not found get best ask trie", "encKey", encKey, "encValue", encValue)
		return EmptyHash
	}
	var data itemList
	if err := rlp.DecodeBytes(encValue, &data); err != nil {
		log.Error("Failed to decode state get best ask trie", "err", err)
		return EmptyHash
	}
	return common.BytesToHash(encKey)
}

func (self *tradingState) getBestBidsTrie(db database.Database) common.Hash {
	trie := self.getBidsTrie(db)
	encKey, encValue, err := trie.TryGetBestRightKeyAndValue()
	if err != nil {
		log.Error("Failed find best liquidationPrice bid trie ", "orderbook", self.orderBookHash.Hex())
		return EmptyHash
	}
	if len(encKey) == 0 || len(encValue) == 0 {
		log.Debug("Not found get best bid trie", "encKey", encKey, "encValue", encValue)
		return EmptyHash
	}
	var data itemList
	if err := rlp.DecodeBytes(encValue, &data); err != nil {
		log.Error("Failed to decode state get best bid trie", "err", err)
		return EmptyHash
	}
	return common.BytesToHash(encKey)
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *tradingState) updateAsksTrie(db database.Database) database.Trie {
	tr := self.getAsksTrie(db)
	for price, orderList := range self.askStates {
		if _, isDirty := self.askStatesDirty[price]; isDirty {
			delete(self.askStatesDirty, price)
			if (orderList.empty()) {
				self.setError(tr.TryDelete(price[:]))
				continue
			}
			orderList.updateRoot(db)
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(orderList)
			self.setError(tr.TryUpdate(price[:], v))
		}
	}

	return tr
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *tradingState) updateAsksRoot(db database.Database) error {
	self.updateAsksTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	self.data.AskRoot = self.asksTrie.Hash()
	return nil
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *tradingState) CommitAsksTrie(db database.Database) error {
	self.updateAsksTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.asksTrie.Commit(func(leaf []byte, parent common.Hash) error {
		var orderList itemList
		if err := rlp.DecodeBytes(leaf, &orderList); err != nil {
			return nil
		}
		if orderList.Root != EmptyRoot {
			db.TrieDB().Reference(orderList.Root, parent)
		}
		return nil
	})
	if err == nil {
		self.data.AskRoot = root
	}
	return err
}

func (self *tradingState) getBidsTrie(db database.Database) database.Trie {
	if self.bidsTrie == nil {
		var err error
		self.bidsTrie, err = db.OpenStorageTrie(self.orderBookHash, self.data.BidRoot)
		if err != nil {
			self.bidsTrie, _ = db.OpenStorageTrie(self.orderBookHash, EmptyHash)
			self.setError(fmt.Errorf("can't create bids trie: %v", err))
		}
	}
	return self.bidsTrie
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *tradingState) updateBidsTrie(db database.Database) database.Trie {
	tr := self.getBidsTrie(db)
	for price, orderList := range self.bidStates {
		if _, isDirty := self.bidStatesDirty[price]; isDirty {
			delete(self.bidStatesDirty, price)
			if (orderList.empty()) {
				self.setError(tr.TryDelete(price[:]))
				continue
			}
			orderList.updateRoot(db)
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(orderList)
			self.setError(tr.TryUpdate(price[:], v))
		}
	}
	return tr
}

func (self *tradingState) updateBidsRoot(db database.Database) {
	self.updateBidsTrie(db)
	self.data.BidRoot = self.bidsTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *tradingState) CommitBidsTrie(db database.Database) error {
	self.updateBidsTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.bidsTrie.Commit(func(leaf []byte, parent common.Hash) error {
		var orderList itemList
		if err := rlp.DecodeBytes(leaf, &orderList); err != nil {
			return nil
		}
		if orderList.Root != EmptyRoot {
			db.TrieDB().Reference(orderList.Root, parent)
		}
		return nil
	})
	if err == nil {
		self.data.BidRoot = root
	}
	return err
}

func (self *tradingState) deepCopy(db *TradingStateDB, onDirty func(hash common.Hash)) *tradingState {
	stateExchanges := newStateExchanges(db, self.orderBookHash, self.data, onDirty)
	if self.asksTrie != nil {
		stateExchanges.asksTrie = db.db.CopyTrie(self.asksTrie)
	}
	if self.bidsTrie != nil {
		stateExchanges.bidsTrie = db.db.CopyTrie(self.bidsTrie)
	}
	if self.ordersTrie != nil {
		stateExchanges.ordersTrie = db.db.CopyTrie(self.ordersTrie)
	}
	for price, bidObject := range self.bidStates {
		stateExchanges.bidStates[price] = bidObject.deepCopy(db, self.MarkStateBidObjectDirty)
	}
	for price, _ := range self.bidStatesDirty {
		stateExchanges.bidStatesDirty[price] = struct{}{}
	}
	for price, askObject := range self.askStates {
		stateExchanges.askStates[price] = askObject.deepCopy(db, self.MarkStateAskObjectDirty)
	}
	for price, _ := range self.askStatesDirty {
		stateExchanges.askStatesDirty[price] = struct{}{}
	}
	for orderId, orderItem := range self.orderStates {
		stateExchanges.orderStates[orderId] = orderItem.deepCopy(self.MarkStateOrderObjectDirty)
	}
	for orderId, _ := range self.orderStatesDirty {
		stateExchanges.orderStatesDirty[orderId] = struct{}{}
	}
	return stateExchanges
}

// Returns the address of the contract/orderId
func (self *tradingState) Hash() common.Hash {
	return self.orderBookHash
}

func (self *tradingState) SetNonce(nonce uint64) {
	self.setNonce(nonce)
}

func (self *tradingState) setNonce(nonce uint64) {
	self.data.Nonce = nonce
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *tradingState) Nonce() uint64 {
	return self.data.Nonce
}

func (self *tradingState) setPrice(price *big.Int) {
	self.data.Price = price
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *tradingState) Price() *big.Int {
	return self.data.Price
}

// updateStateExchangeObject writes the given object to the trie.
func (self *tradingState) removeStateOrderListAskObject(db database.Database, stateOrderList *orderListState) {
	self.setError(self.asksTrie.TryDelete(stateOrderList.price[:]))
}

// updateStateExchangeObject writes the given object to the trie.
func (self *tradingState) removeStateOrderListBidObject(db database.Database, stateOrderList *orderListState) {
	self.setError(self.bidsTrie.TryDelete(stateOrderList.price[:]))
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *tradingState) getStateOrderListAskObject(db database.Database, price common.Hash) (stateOrderList *orderListState) {
	// Prefer 'live' objects.
	if obj := self.askStates[price]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getAsksTrie(db).TryGet(price[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data itemList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "orderId", price, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newOrderListState(self.db, Bid, self.orderBookHash, price, data, self.MarkStateAskObjectDirty)
	self.askStates[price] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *tradingState) MarkStateAskObjectDirty(price common.Hash) {
	self.askStatesDirty[price] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *tradingState) createStateOrderListAskObject(db database.Database, price common.Hash) (newobj *orderListState) {
	newobj = newOrderListState(self.db, Ask, self.orderBookHash, price, itemList{}, self.MarkStateAskObjectDirty)
	self.askStates[price] = newobj
	self.askStatesDirty[price] = struct{}{}
	data, err := rlp.EncodeToBytes(newobj)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", price[:], err))
	}
	self.setError(self.asksTrie.TryUpdate(price[:], data))
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
	return newobj
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *tradingState) getStateBidOrderListObject(db database.Database, price common.Hash) (stateOrderList *orderListState) {
	// Prefer 'live' objects.
	if obj := self.bidStates[price]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getBidsTrie(db).TryGet(price[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data itemList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order list object", "orderId", price, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newOrderListState(self.db, Bid, self.orderBookHash, price, data, self.MarkStateBidObjectDirty)
	self.bidStates[price] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *tradingState) MarkStateBidObjectDirty(price common.Hash) {
	self.bidStatesDirty[price] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *tradingState) createStateBidOrderListObject(db database.Database, price common.Hash) (newobj *orderListState) {
	newobj = newOrderListState(self.db, Bid, self.orderBookHash, price, itemList{}, self.MarkStateBidObjectDirty)
	self.bidStates[price] = newobj
	self.bidStatesDirty[price] = struct{}{}
	data, err := rlp.EncodeToBytes(newobj)
	if err != nil {
		panic(fmt.Errorf("can't encode order list object at %x: %v", price[:], err))
	}
	self.setError(self.bidsTrie.TryUpdate(price[:], data))
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
	return newobj
}

// Retrieve a state object given my the address. Returns nil if not found.
func (self *tradingState) getStateOrderObject(db database.Database, orderId common.Hash) (stateOrderItem *orderItemState) {
	// Prefer 'live' objects.
	if obj := self.orderStates[orderId]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getOrdersTrie(db).TryGet(orderId[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data OrderItem
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state order object", "orderId", orderId, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newOrderItemState(self.orderBookHash, orderId, data, self.MarkStateOrderObjectDirty)
	self.orderStates[orderId] = obj
	return obj
}

// MarkStateAskObjectDirty adds the specified object to the dirty map to avoid costly
// state object cache iteration to find a handful of modified ones.
func (self *tradingState) MarkStateOrderObjectDirty(orderId common.Hash) {
	self.orderStatesDirty[orderId] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

// createStateOrderListObject creates a new state object. If there is an existing orderId with
// the given address, it is overwritten and returned as the second return value.
func (self *tradingState) createStateOrderObject(db database.Database, orderId common.Hash, order OrderItem) (newobj *orderItemState) {
	newobj = newOrderItemState(self.orderBookHash, orderId, order, self.MarkStateOrderObjectDirty)
	orderIdHash := common.BigToHash(new(big.Int).SetUint64(order.OrderID))
	self.orderStates[orderIdHash] = newobj
	self.orderStatesDirty[orderIdHash] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.orderBookHash)
		self.onDirty = nil
	}
	return newobj
}

// updateAskTrie writes cached storage modifications into the object's storage trie.
func (self *tradingState) updateOrdersTrie(db database.Database) database.Trie {
	tr := self.getOrdersTrie(db)
	for orderId, orderItem := range self.orderStates {
		if _, isDirty := self.orderStatesDirty[orderId]; isDirty {
			delete(self.orderStatesDirty, orderId)
			if (orderItem.empty()) {
				self.setError(tr.TryDelete(orderId[:]))
				continue
			}
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(orderItem)
			self.setError(tr.TryUpdate(orderId[:], v))
		}
	}
	return tr
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *tradingState) updateOrdersRoot(db database.Database) {
	self.updateOrdersTrie(db)
	self.data.OrderRoot = self.ordersTrie.Hash()
}

// CommitAskTrie the storage trie of the object to dwb.
// This updates the trie root.
func (self *tradingState) CommitOrdersTrie(db database.Database) error {
	self.updateOrdersTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.ordersTrie.Commit(nil)
	if err == nil {
		self.data.OrderRoot = root
	}
	return err
}

func (self *tradingState) MarkStateLiquidationPriceDirty(price common.Hash) {
	self.liquidationPriceStatesDirty[price] = struct{}{}
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
}

func (self *tradingState) createStateLiquidationPrice(db database.Database, liquidationPrice common.Hash) (newobj *liquidationPriceState) {
	newobj = newLiquidationPriceState(self.db, self.orderBookHash, liquidationPrice, itemList{}, self.MarkStateLiquidationPriceDirty)
	self.liquidationPriceStates[liquidationPrice] = newobj
	self.liquidationPriceStatesDirty[liquidationPrice] = struct{}{}
	data, err := rlp.EncodeToBytes(newobj)
	if err != nil {
		panic(fmt.Errorf("can't encode liquidation price object at %x: %v", liquidationPrice[:], err))
	}
	self.setError(self.bidsTrie.TryUpdate(liquidationPrice[:], data))
	if self.onDirty != nil {
		self.onDirty(self.Hash())
		self.onDirty = nil
	}
	return newobj
}

func (self *tradingState) getLiquidationPriceTrie(db database.Database) database.Trie {
	if self.liquidationPriceTrie == nil {
		var err error
		self.liquidationPriceTrie, err = db.OpenStorageTrie(self.orderBookHash, self.data.LiquidationPriceRoot)
		if err != nil {
			self.liquidationPriceTrie, _ = db.OpenStorageTrie(self.orderBookHash, EmptyHash)
			self.setError(fmt.Errorf("can't create liquidation liquidationPrice trie: %v", err))
		}
	}
	return self.liquidationPriceTrie
}

func (self *tradingState) getStateLiquidationPrice(db database.Database, price common.Hash) (stateObject *liquidationPriceState) {
	// Prefer 'live' objects.
	if obj := self.liquidationPriceStates[price]; obj != nil {
		return obj
	}

	// Load the object from the database.
	enc, err := self.getLiquidationPriceTrie(db).TryGet(price[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data itemList
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state liquidation liquidationPrice", "liquidationPrice", price, "err", err)
		return nil
	}
	// Insert into the live set.
	obj := newLiquidationPriceState(self.db, self.orderBookHash, price, data, self.MarkStateLiquidationPriceDirty)
	self.liquidationPriceStates[price] = obj
	return obj
}

func (self *tradingState) getLowestLiquidationPrice(db database.Database) (common.Hash, *liquidationPriceState) {
	trie := self.getLiquidationPriceTrie(db)
	encKey, encValue, err := trie.TryGetBestLeftKeyAndValue()
	if err != nil {
		log.Error("Failed find best liquidationPrice ask trie ", "orderbook", self.orderBookHash.Hex())
		return EmptyHash, nil
	}
	if len(encKey) == 0 || len(encValue) == 0 {
		log.Debug("Not found get best ask trie", "encKey", encKey, "encValue", encValue)
		return EmptyHash, nil
	}
	var data itemList
	if err := rlp.DecodeBytes(encValue, &data); err != nil {
		log.Error("Failed to decode state get best ask trie", "err", err)
		return EmptyHash, nil
	}
	price := common.BytesToHash(encKey)
	obj := newLiquidationPriceState(self.db, self.orderBookHash, price, data, self.MarkStateLiquidationPriceDirty)
	self.liquidationPriceStates[price] = obj
	return price, obj
}

func (self *tradingState) updateLiquidationPriceTrie(db database.Database) database.Trie {
	tr := self.getLiquidationPriceTrie(db)
	for price, stateObject := range self.liquidationPriceStates {
		if _, isDirty := self.liquidationPriceStatesDirty[price]; isDirty {
			delete(self.orderStatesDirty, price)
			if (stateObject.empty()) {
				self.setError(tr.TryDelete(price[:]))
				continue
			}
			stateObject.updateRoot(db)
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(stateObject)
			self.setError(tr.TryUpdate(price[:], v))
		}
	}
	return tr
}

func (self *tradingState) updateLiquidationPriceRoot(db database.Database) {
	self.updateLiquidationPriceTrie(db)
	self.data.LiquidationPriceRoot = self.liquidationPriceTrie.Hash()
}

func (self *tradingState) CommitLiquidationPriceTrie(db database.Database) error {
	self.updateLiquidationPriceTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	root, err := self.liquidationPriceTrie.Commit(nil)
	if err == nil {
		self.data.LiquidationPriceRoot = root
	}
	return err
}
