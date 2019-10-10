// Copyright 2015 The go-ethereum Authors
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

package state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// ManagedOrderState mange virual nonce for order
type ManagedOrderState struct {
	*OrderState

	mu sync.RWMutex

	accounts map[common.Address]uint64
}

// NewManagedOrderState returns a new managed state with the statedb as it's backing layer
func NewManagedOrderState(orderstate *OrderState) *ManagedOrderState {
	return &ManagedOrderState{
		OrderState: orderstate,
		accounts:   make(map[common.Address]uint64),
	}
}

// SetState sets the backing layer of the managed state
func (ms *ManagedOrderState) SetState(statedb *OrderState) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.OrderState = statedb
}

// GetNonce returns the canonical nonce for the managed or unmanaged account.
//
// Because GetNonce mutates the DB, we must take a write lock.
func (ms *ManagedOrderState) GetNonce(addr common.Address) uint64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if val, ok := ms.accounts[addr]; ok {
		return val
	} else {
		return ms.OrderState.GetNonce(addr)
	}
}

// SetNonce sets the new canonical nonce for the managed state
func (ms *ManagedOrderState) SetNonce(addr common.Address, nonce uint64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.accounts[addr] = nonce
}
