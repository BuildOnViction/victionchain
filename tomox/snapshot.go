package tomox

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	snapshotPrefix    = "tomox-snap-"
	latestSnapshotKey = "tomox-snap-latest"
)

var (
	// errors
	errOrderBooKSnapshotNotFound = errors.New("orderBook snapshot not found")
	errInvalidTreeHash           = errors.New("orderTree: invalid hash")
	errInvalidOrderBookHash      = errors.New("orderBook: invalid hash")
	errInvalidOrderListHash      = errors.New("orderList: invalid hash")
)

// tomox snapshot
type Snapshot struct {
	OrderBooks map[string]*OrderBookSnapshot
	Hash       common.Hash
}

// snapshot of OrderBook
type OrderBookSnapshot struct {
	Bids          *OrderTreeSnapshot
	Asks          *OrderTreeSnapshot
	OrderBookItem []byte
}

// snapshot of OrderTree
type OrderTreeSnapshot struct {
	OrderTreeItem []byte
	OrderList     map[common.Hash]*OrderListSnapshot // common.BytesToHash(getKeyFromPrice(price)) => orderlist
}

// snapshot of OrderList
type OrderListSnapshot struct {
	OrderListItem []byte
	OrderItem     [][]byte // slice of orderItems, encode each orderItem to []byte
}

// put tomox snapshot to db
func (s *Snapshot) store(db OrderDao) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return db.Put(append([]byte(snapshotPrefix), s.Hash[:]...), &blob, false, common.Hash{})
}

// take a snapshot of data of tomox
func newSnapshot(tomox *TomoX, blockHash common.Hash) (*Snapshot, error) {
	var (
		ob                       *OrderBook
		bidTreeSnap, askTreeSnap *OrderTreeSnapshot
		encodedBytes             []byte
		err                      error
	)
	snap := new(Snapshot)
	snap.Hash = blockHash
	snap.OrderBooks = make(map[string]*OrderBookSnapshot)
	for _, pair := range tomox.listTokenPairs() {
		ob, err = tomox.GetOrderBook(pair, false, common.Hash{})
		if err != nil {
			return nil, err
		}
		obSnap := new(OrderBookSnapshot)
		encodedBytes, err = EncodeBytesItem(ob.Item)
		if err != nil {
			return snap, err
		}
		obSnap.OrderBookItem = encodedBytes
		if bidTreeSnap, err = prepareOrderTreeData(ob.Bids); err != nil {
			return snap, err
		}
		if askTreeSnap, err = prepareOrderTreeData(ob.Asks); err != nil {
			return snap, err
		}
		obSnap.Bids = bidTreeSnap
		obSnap.Asks = askTreeSnap
		snap.OrderBooks[pair] = obSnap
	}
	return snap, nil
}

// take a snapshot of orderTree
func prepareOrderTreeData(tree *OrderTree) (*OrderTreeSnapshot, error) {
	var (
		serializedTree []byte
		err            error
	)
	snap := new(OrderTreeSnapshot)
	serializedTree, err = EncodeBytesItem(tree.Item)
	if err != nil {
		return nil, err
	}
	snap.OrderTreeItem = serializedTree

	snap.OrderList = make(map[common.Hash]*OrderListSnapshot)
	if !tree.NotEmpty() || tree.PriceTree.Size() == 0 {
		return snap, nil
	}

	// foreach each price, snapshot its orderlist
	for _, key := range tree.PriceTree.Keys(false, common.Hash{}) {
		priceKeyHash := common.BytesToHash(key)
		bytes, found := tree.PriceTree.Get(key, false, common.Hash{})
		if found {
			var ol *OrderList
			ol, err = tree.decodeOrderList(bytes)
			if err != nil {
				return nil, err
			}
			if ol.Item.Length == 0 {
				return snap, nil
			}
			// snapshot orderItems
			var (
				items    [][]byte
				byteItem []byte
			)
			order := ol.GetOrder(ol.Item.HeadOrder, false, common.Hash{})
			for order != nil {
				if byteItem, err = EncodeBytesItem(order.Item); err != nil {
					return nil, err
				}
				items = append(items, byteItem)
				order = order.GetNextOrder(ol, false, common.Hash{})
			}
			snap.OrderList[priceKeyHash] = &OrderListSnapshot{
				OrderListItem: bytes,
				OrderItem:     items,
			}
		}
	}
	return snap, nil
}

// load snapshot from database when nodes restart
func getSnapshot(db OrderDao, blockHash common.Hash) (*Snapshot, error) {
	var (
		blob interface{}
		err  error
	)
	blob, err = db.Get(append([]byte(snapshotPrefix), blockHash[:]...), &[]byte{}, false, common.Hash{})
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if blob != nil {
		if err = json.Unmarshal(*blob.(*[]byte), snap); err != nil {
			return nil, err
		}
	}
	return snap, nil
}

// import all orderbooks from snapshot
func (s *Snapshot) RestoreOrderBookFromSnapshot(db OrderDao, pairName string) (*OrderBook, error) {
	var (
		obSnap     *OrderBookSnapshot
		bids, asks *OrderTree
		ob         *OrderBook
		err        error
		ok         bool
	)
	if obSnap, ok = s.OrderBooks[pairName]; !ok {
		return &OrderBook{}, errOrderBooKSnapshotNotFound
	}
	orderBookItem := &OrderBookItem{}
	if err = DecodeBytesItem(obSnap.OrderBookItem, orderBookItem); err != nil {
		return &OrderBook{}, err
	}
	key := crypto.Keccak256([]byte(orderBookItem.Name))
	slot := new(big.Int).SetBytes(key)

	ob = &OrderBook{
		Item: orderBookItem,
		Key:  key,
		Slot: slot,
		db:   db,
	}
	ob.Bids = NewOrderTree(db, GetSegmentHash(key, 1, SlotSegment), ob)
	ob.Asks = NewOrderTree(db, GetSegmentHash(key, 2, SlotSegment), ob)

	if bids, err = s.RestoreOrderTree(obSnap.Bids, ob.Bids); err != nil {
		return &OrderBook{}, err
	}
	if asks, err = s.RestoreOrderTree(obSnap.Asks, ob.Asks); err != nil {
		return &OrderBook{}, err
	}
	ob.Bids = bids
	ob.Asks = asks
	// verify hash
	if err = verifyHash(ob, common.BytesToHash(obSnap.OrderBookItem)); err != nil {
		return &OrderBook{}, err
	}
	return ob, nil
}

func verifyHash(o interface{}, hash common.Hash) error {
	var (
		h   common.Hash
		err error
	)
	switch o.(type) {
	case *OrderBook:
		h, err = o.(*OrderBook).Hash()
		if h != hash {
			return errInvalidOrderBookHash
		}
		break
	case *OrderTree:
		h, err = o.(*OrderTree).Hash()
		if h != hash {
			return errInvalidTreeHash
		}
		break
	case *OrderList:
		h, err = o.(*OrderList).Hash()
		if h != hash {
			return errInvalidOrderListHash
		}
		break
	default:
		return nil

	}
	return err
}

// restore orderTree from snapshot
func (s *Snapshot) RestoreOrderTree(treeSnap *OrderTreeSnapshot, tree *OrderTree) (*OrderTree, error) {
	var err error
	orderTreeItem := &OrderTreeItem{}
	orderListItem := &OrderListItem{}

	// restore bids tree
	if err = DecodeBytesItem(treeSnap.OrderTreeItem, orderTreeItem); err != nil {
		return tree, err
	}
	tree.Item = orderTreeItem
	for _, olSnap := range treeSnap.OrderList {
		err = DecodeBytesItem(olSnap.OrderListItem, orderListItem)
		if err != nil {
			return tree, err
		}
		ol := NewOrderListWithItem(orderListItem, tree)
		// orderlist hash from snapshot
		orderListSnapHash := common.BytesToHash(treeSnap.OrderList[common.BytesToHash(ol.Key)].OrderListItem)
		if err = verifyHash(ol, orderListSnapHash); err != nil {
			return tree, err
		}
		if err = tree.SaveOrderList(ol, false, common.Hash{}); err != nil {
			return tree, err
		}

		// try to update order from snapshot to db in case of missing order in db
		for _, item := range olSnap.OrderItem {
			orderItem := &OrderItem{}
			if err = DecodeBytesItem(item, orderItem); err != nil {
				return tree, err
			}
			order := NewOrder(orderItem, GetOrderListCommonKey(ol.Key, tree.orderBook.Item.Name))
			if err = ol.SaveOrder(order, false, common.Hash{}); err != nil {
				return tree, err
			}
		}

	}
	if err = verifyHash(tree, common.BytesToHash(treeSnap.OrderTreeItem)); err != nil {
		return tree, err
	}
	return tree, nil
}
