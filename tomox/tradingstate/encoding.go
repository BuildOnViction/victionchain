package tradingstate

import (
	"github.com/tomochain/tomochain/rlp"
)

func EncodeBytesItem(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func DecodeBytesItem(bytes []byte, val interface{}) error {
	return rlp.DecodeBytes(bytes, val)

}
