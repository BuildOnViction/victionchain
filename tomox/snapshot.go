package tomox

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	prefix            = "tomox-snap-"
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
	Bids *OrderTreeSnapshot
	Asks *OrderTreeSnapshot
	Data []byte
	Hash common.Hash
}

// snapshot of OrderTree
type OrderTreeSnapshot struct {
	Data      []byte
	Hash      common.Hash
	OrderList map[common.Hash]*OrderListSnapshot // common.BytesToHash(getKeyFromPrice(price)) => orderlist
}

// snapshot of OrderList
type OrderListSnapshot struct {
	Data []byte
}

// put tomox snapshot to db
func (s *Snapshot) store(db OrderDao) error {
	// TODO: need to clean up this kind of data later
	return db.Put(append([]byte(prefix), s.Hash[:]...), s)
}

// take a snapshot of data of tomox
func newSnapshot(tomox *TomoX, blockHash common.Hash) (*Snapshot, error) {
	var (
		orderBookHash            common.Hash
		bidTreeSnap, askTreeSnap *OrderTreeSnapshot
		encodedBytes             []byte
		err                      error
	)
	snap := new(Snapshot)
	snap.Hash = blockHash
	snap.OrderBooks = make(map[string]*OrderBookSnapshot)
	for pair, ob := range tomox.Orderbooks {
		obSnap := new(OrderBookSnapshot)
		encodedBytes, err = EncodeBytesItem(ob.Item)
		if err != nil {
			return snap, err
		}
		obSnap.Data = encodedBytes

		orderBookHash, err = ob.Hash()
		if err != nil {
			return snap, err
		}
		obSnap.Hash = orderBookHash
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
		orderTreeHash  common.Hash
		serializedTree []byte
		err            error
	)
	snap := new(OrderTreeSnapshot)
	if orderTreeHash, err = tree.Hash(); err != nil {
		return &OrderTreeSnapshot{}, err
	}
	snap.Hash = orderTreeHash
	serializedTree, err = EncodeBytesItem(tree.Item)
	if err != nil {
		return &OrderTreeSnapshot{}, err
	}
	snap.Data = serializedTree

	snap.OrderList = make(map[common.Hash]*OrderListSnapshot)
	// foreach each price, snapshot its orderlist
	for _, key := range tree.PriceTree.Keys() {
		priceKeyHash := common.BytesToHash(key)
		bytes, found := tree.PriceTree.Get(key)
		if found {
			snap.OrderList[priceKeyHash] = &OrderListSnapshot{
				Data: bytes,
			}
		}
	}
	return snap, nil
}

// load snapshot from database when nodes restart
func loadSnapshot(db OrderDao, blockHash common.Hash) (*Snapshot, error) {
	var (
		blob interface{}
		err  error
	)
	blob, err = db.Get(append([]byte(prefix), blockHash[:]...), blob)
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	snap = blob.(*Snapshot)
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
	if err = DecodeBytesItem(obSnap.Data, orderBookItem); err != nil {
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
	if err = verifyHash(ob, obSnap.Hash); err != nil {
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
	if err = DecodeBytesItem(treeSnap.Data, orderTreeItem); err != nil {
		return tree, err
	}
	tree.Item = orderTreeItem
	for _, olSnap := range treeSnap.OrderList {
		err = DecodeBytesItem(olSnap.Data, orderListItem)
		if err != nil {
			return tree, err
		}
		ol := NewOrderListWithItem(orderListItem, tree)
		// orderlist hash from snapshot
		orderListSnapHash := common.BytesToHash(treeSnap.OrderList[common.BytesToHash(ol.Key)].Data)
		if err = verifyHash(ol, orderListSnapHash); err != nil {
			return tree, err
		}
		if err = tree.SaveOrderList(ol); err != nil {
			return tree, err
		}
	}
	if err = verifyHash(tree, treeSnap.Hash); err != nil {
		return tree, err
	}
	return tree, nil
}
