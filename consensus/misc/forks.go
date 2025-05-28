// Copyright 2017 The go-ethereum Authors
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

package misc

import (
	"fmt"
	"math/big"

	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/contracts/vrc25issuer"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/params"
)

// VerifyForkHashes verifies that blocks conforming to network hard-forks do have
// the correct hashes, to avoid clients going off on different chains. This is an
// optional feature.
func VerifyForkHashes(config *params.ChainConfig, header *types.Header, uncle bool) error {
	// We don't care about uncles
	if uncle {
		return nil
	}
	// If the homestead reprice hash is set, validate it
	if config.EIP150Block != nil && config.EIP150Block.Cmp(header.Number) == 0 {
		if config.EIP150Hash != (common.Hash{}) && config.EIP150Hash != header.Hash() {
			return fmt.Errorf("homestead gas reprice fork: have 0x%x, want 0x%x", header.Hash(), config.EIP150Hash)
		}
	}
	// All ok, return
	return nil
}

// ApplySaigonHardFork mint additional token to EcoSystem Multisig preiodly for 4 years
func ApplySaigonHardFork(statedb *state.StateDB, saigonBlock *big.Int, headBlock *big.Int) {
	endBlock := new(big.Int).Add(saigonBlock, new(big.Int).SetUint64(common.SaigonEcoSystemFundInterval*(common.SaigonEcoSystemFundTotalRepeat-1))) // additional token will be minted at block 0 of each interval 4 intervals
	if headBlock.Cmp(saigonBlock) < 0 || headBlock.Cmp(endBlock) > 0 {
		return
	}
	blockOfInterval := new(big.Int).Mod(new(big.Int).Sub(headBlock, saigonBlock), new(big.Int).SetUint64(common.SaigonEcoSystemFundInterval))
	if blockOfInterval.Cmp(big.NewInt(0)) == 0 {
		ecoSystemFund := new(big.Int).Mul(common.SaigonEcoSystemFundAnnually, new(big.Int).SetUint64(params.Ether))
		statedb.AddBalance(common.SaigonEcoSystemFundAddress, ecoSystemFund)
	}
}

// ApplySaigonHardFork mint additional token to EcoSystem Multisig once. For testnet only.
func ApplySaigonHardForkTestnet(statedb *state.StateDB, saigonBlock *big.Int, headBlock *big.Int, posv *params.PosvConfig) {
	if headBlock.Cmp(saigonBlock) == 0 && posv != nil {
		ecoSystemFund := new(big.Int).Mul(new(big.Int).Mul(common.SaigonEcoSystemFundAnnually, new(big.Int).SetUint64(params.Ether)), new(big.Int).SetUint64(common.SaigonEcoSystemFundTotalRepeat))
		statedb.AddBalance(posv.FoudationWalletAddr, ecoSystemFund)
	}
}

func ApplyVIPVRC25Upgarde(statedb *state.StateDB, vipVRC25Block *big.Int, headBlock *big.Int) {
	if headBlock.Cmp(vipVRC25Block) == 0 {
		minCapLoc := state.GetLocSimpleVariable(common.TRC21IssuerMinCapSlot)
		statedb.SetState(common.TRC21IssuerSMC, minCapLoc, common.BigToHash(common.VRC25IssuerMinCap))
		statedb.SetCode(common.TRC21IssuerSMC, []byte(vrc25issuer.Vrc25issuerBin))
	}
}
