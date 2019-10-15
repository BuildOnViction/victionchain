package tomox_state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
)

func GetLocMappingAtKey(key common.Hash, slot uint64) *big.Int {
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	retByte := crypto.Keccak256(key.Bytes(), slotHash.Bytes())
	ret := new(big.Int)
	ret.SetBytes(retByte)
	return ret
}

// getCoinbaseCount get number of coinbase
func getCoinbaseCount(statedb *state.StateDB) *big.Int {
	slot := RelayerMappingSlot["RelayerCount"]
	slotHash := common.BigToHash(big.NewInt(int64(slot)))
	return statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), slotHash).Big()
}

// GetCoinbaseList get all coinbase
func getCoinbaseList(statedb *state.StateDB) []common.Address {
	var listCoinBase []common.Address
	numberCoinbase := getCoinbaseCount(statedb)
	log.Info("getCoinbaseList", "total", numberCoinbase)
	slot := RelayerMappingSlot["RELAYER_COINBASES"]
	for i := big.NewInt(0); i.Cmp(numberCoinbase) == -1; i.Add(i, big.NewInt(1)) {
		locBig := GetLocMappingAtKey(common.BigToHash(i), slot)
		locHash := common.BigToHash(locBig)
		coinbase := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Big()
		listCoinBase = append(listCoinBase, common.BytesToAddress(coinbase.Bytes()))
	}
	return listCoinBase
}

//GetCoinbaseFeeList get add coinbase fee
func GetCoinbaseFeeList(statedb *state.StateDB) map[common.Address]*big.Int {
	log.Info("GetCoinbaseFeeList start...")
	relayerFeeList := make(map[common.Address]*big.Int)
	coinbaseList := getCoinbaseList(statedb)
	for _, cb := range coinbaseList {
		tradeFee := GetExRelayerFee(cb, statedb)
		log.Info("GetCoinbaseFeeList", "coinbase", cb, "fee", tradeFee)
		relayerFeeList[cb] = tradeFee
	}
	return relayerFeeList
}

func GetExRelayerFee(relayer common.Address, statedb *state.StateDB) *big.Int {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_fee"])
	locHash := common.BigToHash(locBig)
	return statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Big()
}
func GetRelayerOwner(relayer common.Address, statedb *state.StateDB) common.Address {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	log.Debug("GetRelayerOwner", "relayer", relayer.Hex(), "slot", slot, "locBig", locBig)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_owner"])
	locHash := common.BigToHash(locBig)
	return common.BytesToAddress(statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Bytes())
}

func GetBaseTokenLength(relayer common.Address, statedb *state.StateDB) uint64 {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_fromTokens"])
	locHash := common.BigToHash(locBig)
	return statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Big().Uint64()
}

func GetBaseTokenAtIndex(relayer common.Address, statedb *state.StateDB, index uint64) common.Address {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_fromTokens"])
	locHash := common.BigToHash(locBig)
	loc := state.GetLocDynamicArrAtElement(locHash, index, 1)
	return common.BytesToAddress(statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), loc).Bytes())
}

func GetQuoteTokenLength(relayer common.Address, statedb *state.StateDB) uint64 {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_toTokens"])
	locHash := common.BigToHash(locBig)
	return statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHash).Big().Uint64()
}

func GetQuoteTokenAtIndex(relayer common.Address, statedb *state.StateDB, index uint64) common.Address {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)
	locBig = locBig.Add(locBig, RelayerStructMappingSlot["_toTokens"])
	locHash := common.BigToHash(locBig)
	loc := state.GetLocDynamicArrAtElement(locHash, index, 1)
	return common.BytesToAddress(statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), loc).Bytes())
}

func SubRelayerFee(relayer common.Address, fee *big.Int, statedb *state.StateDB) error {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)

	locBigDeposit := new(big.Int).SetUint64(uint64(0)).Add(locBig, RelayerStructMappingSlot["_deposit"])
	locHashDeposit := common.BigToHash(locBigDeposit)
	balance := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit).Big()
	log.Debug("ApplyTomoXMatchedTransaction settle balance: SubRelayerFee BEFORE", "relayer", relayer.String(), "balance", balance)
	if balance.Cmp(fee) < 0 {
		return errors.Errorf("relayer %s isn't enough tomo fee", relayer.String())
	} else {
		balance = balance.Sub(balance, fee)
		statedb.SetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit, common.BigToHash(balance))
		statedb.SubBalance(common.HexToAddress(common.RelayerRegistrationSMC), fee)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SubRelayerFee AFTER", "relayer", relayer.String(), "balance", balance)
		return nil
	}
}

func AddTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB) error {
	// TOMO native
	if token.String() == common.TomoNativeAddress {
		balance := statedb.GetBalance(addr)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: ADD TOKEN TOMO NATIVE BEFORE", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		balance = balance.Add(balance, value)
		statedb.SetBalance(addr, balance)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: ADD TOMO NATIVE BALANCE AFTER", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)

		return nil
	}

	// TRC tokens
	if statedb.Exist(token) {
		slot := TokenMappingSlot["balances"]
		locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
		balance := statedb.GetState(token, locHash).Big()
		log.Debug("ApplyTomoXMatchedTransaction settle balance: ADD TOKEN BALANCE BEFORE", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		balance = balance.Add(balance, value)
		statedb.SetState(token, locHash, common.BigToHash(balance))
		log.Debug("ApplyTomoXMatchedTransaction settle balance: ADD TOKEN BALANCE AFTER", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		return nil
	} else {
		return errors.Errorf("token %s isn't exist", token.String())
	}
}

func SubTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB) error {
	// TOMO native
	if token.String() == common.TomoNativeAddress {

		balance := statedb.GetBalance(addr)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SUB TOMO NATIVE BALANCE BEFORE", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		if balance.Cmp(value) < 0 {
			return errors.Errorf("value %s in token %s not enough , have : %s , want : %s  ", addr.String(), token.String(), balance, value)
		}
		balance = balance.Sub(balance, value)
		statedb.SetBalance(addr, balance)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SUB TOMO NATIVE BALANCE AFTER", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		return nil
	}

	// TRC tokens
	if statedb.Exist(token) {
		slot := TokenMappingSlot["balances"]
		locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
		balance := statedb.GetState(token, locHash).Big()
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SUB TOKEN BALANCE BEFORE", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		if balance.Cmp(value) < 0 {
			return errors.Errorf("value %s in token %s not enough , have : %s , want : %s  ", addr.String(), token.String(), balance, value)
		}
		balance = balance.Sub(balance, value)
		statedb.SetState(token, locHash, common.BigToHash(balance))
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SUB TOKEN BALANCE AFTER", "token", token.String(), "address", addr.String(), "balance", balance, "orderValue", value)
		return nil
	} else {
		return errors.Errorf("token %s isn't exist", token.String())
	}
}
