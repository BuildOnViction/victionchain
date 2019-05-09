package tomox

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBytesItem(val interface{}) ([]byte, error) {

	switch val.(type) {
	case *Item:
		return rlp.EncodeToBytes(val.(*Item))
	case *OrderItem:
		return rlp.EncodeToBytes(val.(*OrderItem))
	case *OrderListItem:
		return rlp.EncodeToBytes(val.(*OrderListItem))
	case *OrderTreeItem:
		return rlp.EncodeToBytes(val.(*OrderTreeItem))
	case *OrderBookItem:
		return rlp.EncodeToBytes(val.(*OrderBookItem))
	default:
		return rlp.EncodeToBytes(val)
	}
}

func DecodeBytesItem(bytes []byte, val interface{}) error {

	switch val.(type) {
	case *Item:
		return rlp.DecodeBytes(bytes, val.(*Item))
	case *OrderItem:
		return rlp.DecodeBytes(bytes, val.(*OrderItem))
	case *OrderListItem:
		return rlp.DecodeBytes(bytes, val.(*OrderListItem))
	case *OrderTreeItem:
		return rlp.DecodeBytes(bytes, val.(*OrderTreeItem))
	case *OrderBookItem:
		return rlp.DecodeBytes(bytes, val.(*OrderBookItem))
	default:
		return rlp.DecodeBytes(bytes, val)
	}

}
