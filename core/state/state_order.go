package state

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

// OrderState order state
type OrderState struct {
	orderNonce map[common.Address]*big.Int
}

// NewOrderState new state order
func NewOrderState(orderNonce map[common.Address]*big.Int) *OrderState {
	return &OrderState{
		orderNonce: orderNonce,
	}
}

// GetNonce get order nonce from order state
func (os *OrderState) GetNonce(addr common.Address) uint64 {
	if orderNonce, ok := os.orderNonce[addr]; ok {
		bigstr := orderNonce.String()
		n, err := strconv.ParseInt(bigstr, 10, 64)
		if err != nil {
			return uint64(n)
		}
	}
	return 0
}
