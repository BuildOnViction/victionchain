package tomox

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
)

func EncodeBytesItem(val interface{}) ([]byte, error) {

	switch val.(type) {
	case *Item:
		return rlp.EncodeToBytes(val.(*Item))
	case *tomox_state.OrderItem:
		return rlp.EncodeToBytes(val.(*tomox_state.OrderItem))
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
	case *tomox_state.OrderItem:
		return rlp.DecodeBytes(bytes, val.(*tomox_state.OrderItem))
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
