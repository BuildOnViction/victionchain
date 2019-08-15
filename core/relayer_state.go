package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/tomox"
	"github.com/pkg/errors"
	"math/big"
)

func getLocMappingAtKey(key common.Hash, slot uint64) *big.Int {
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	retByte := crypto.Keccak256(key.Bytes(), slotHash.Bytes())
	ret := new(big.Int)
	ret.SetBytes(retByte)
	return ret
}

func GetExRelayerFee(relayer common.Address, statedb *state.StateDB) *big.Int {
	slot := tomox.RelayerMappingSlot["RELAYER_LIST"]
	locBig := getLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, tomox.RelayerStructMappingSlot["_fee"])
	locHash := common.BigToHash(locBig)
	return statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Big()

}
func GetRelayerOwner(relayer common.Address, statedb *state.StateDB) common.Address {
	slot := tomox.RelayerMappingSlot["OWNER_LIST"]
	locBig := getLocMappingAtKey(relayer.Hash(), slot)
	locHash := common.BigToHash(locBig)
	return common.BytesToAddress(statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Bytes())

}
func SubRelayerFee(relayer common.Address, fee *big.Int, statedb *state.StateDB) error {
	slot := tomox.RelayerMappingSlot["RELAYER_LIST"]
	locBig := getLocMappingAtKey(relayer.Hash(), slot)

	locBigDeposit := new(big.Int).SetUint64(uint64(0)).Add(locBig, tomox.RelayerStructMappingSlot["_deposit"])
	locHashDeposit := common.BigToHash(locBigDeposit)
	balance := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit).Big()

	if balance.Cmp(fee) < 0 {
		return errors.Errorf("relayer %s isn't enough tomo fee", relayer.String())
	} else {
		balance = balance.Sub(balance, fee)
		statedb.SetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit, common.BigToHash(balance))
		statedb.SubBalance(common.HexToAddress(common.RelayerRegistrationSMC), fee)
		return nil
	}
}

func AddTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB) error {
	if statedb.Exist(addr) {
		slot := tomox.TokenMappingSlot["balances"]
		locHash := common.BigToHash(getLocMappingAtKey(addr.Hash(), slot))
		balance := statedb.GetState(token, locHash).Big()
		balance = balance.Add(balance, value)
		statedb.SetState(token, locHash, common.BigToHash(balance))
		return nil
	} else {
		return errors.Errorf("token %s isn't exist", addr)
	}
}

func SubTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB) error {
	if statedb.Exist(addr) {
		slot := tomox.TokenMappingSlot["balances"]
		locHash := common.BigToHash(getLocMappingAtKey(addr.Hash(), slot))
		balance := statedb.GetState(token, locHash).Big()
		if balance.Cmp(value) < 0 {
			return errors.Errorf("value %s in token %s not enough , have : %s , want : %s  ", addr, token, balance, value)
		}
		balance = balance.Sub(balance, value)
		statedb.SetState(token, locHash, common.BigToHash(balance))
		return nil
	} else {
		return errors.Errorf("token %s isn't exist", addr)
	}
}
