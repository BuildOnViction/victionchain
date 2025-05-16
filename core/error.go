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

package core

import "errors"

var (
	// ErrKnownBlock is returned when a block to import is already known locally.
	ErrKnownBlock = errors.New("block already known")

	// ErrGasLimitReached is returned by the gas pool if the amount of gas required
	// by a transaction is higher than what's left in the block.
	ErrGasLimitReached = errors.New("gas limit reached")

	// ErrBlacklistedHash is returned if a block to import is on the blacklist.
	ErrBlacklistedHash = errors.New("blacklisted hash")

	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// next one expected based on the local chain.
	ErrNonceTooHigh = errors.New("nonce too high")

	ErrNotPoSV = errors.New("Posv not found in config")

	ErrNotFoundM1 = errors.New("list M1 not found ")

	ErrStopPreparingBlock = errors.New("stop calculating a block not verified by M2")

	// ErrInsufficientPayerFunds is returned if the gas fee cost of executing a transaction
	// is higher than the balance of the payer's account.
	ErrInsufficientPayerFunds = errors.New("insufficient payer funds for gas * price")

	// ErrInsufficientSenderFunds is returned if the value in transaction
	// is higher than the balance of the user's account.
	ErrInsufficientSenderFunds = errors.New("insufficient sender funds for value")
)
