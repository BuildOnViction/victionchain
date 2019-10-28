package tomox_state

import (
	"fmt"
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
		balance = new(big.Int).Sub(balance, fee)
		statedb.SetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit, common.BigToHash(balance))
		statedb.SubBalance(common.HexToAddress(common.RelayerRegistrationSMC), fee)
		log.Debug("ApplyTomoXMatchedTransaction settle balance: SubRelayerFee AFTER", "relayer", relayer.String(), "balance", balance)
		return nil
	}
}

func CheckRelayerFee(relayer common.Address, fee *big.Int, statedb *state.StateDB) error {
	slot := RelayerMappingSlot["RELAYER_LIST"]
	locBig := GetLocMappingAtKey(relayer.Hash(), slot)

	locBigDeposit := new(big.Int).SetUint64(uint64(0)).Add(locBig, RelayerStructMappingSlot["_deposit"])
	locHashDeposit := common.BigToHash(locBigDeposit)
	balance := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit).Big()
	if balance.Cmp(fee) < 0 {
		return errors.Errorf("relayer %s isn't enough tomo fee : balance %d , fee : %d ", relayer.Hex(), balance.Uint64(), fee.Uint64())
	}
	return nil
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

func CheckSubTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB, mapBalances map[common.Address]map[common.Address]*big.Int) (*big.Int, error) {
	// TOMO native
	if token.String() == common.TomoNativeAddress {
		var balance *big.Int
		if value := mapBalances[token][addr]; value != nil {
			balance = value
		} else {
			balance = statedb.GetBalance(addr)
		}
		if balance.Cmp(value) < 0 {
			return nil, errors.Errorf("value %s in token %s not enough , have : %s , want : %s  ", addr.String(), token.String(), balance, value)
		}
		newBalance := new(big.Int).Sub(balance, value)
		log.Debug("CheckSubTokenBalance settle balance: SUB TOMO NATIVE BALANCE ", "token", token.String(), "address", addr.String(), "balance", balance, "value", value, "newBalance", newBalance)
		return newBalance, nil
	}
	// TRC tokens
	if statedb.Exist(token) {
		var balance *big.Int
		if value := mapBalances[token][addr]; value != nil {
			balance = value
		} else {
			slot := TokenMappingSlot["balances"]
			locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
			balance = statedb.GetState(token, locHash).Big()
		}
		if balance.Cmp(value) < 0 {
			return nil, errors.Errorf("value %s in token %s not enough , have : %s , want : %s  ", addr.String(), token.String(), balance, value)
		}
		newBalance := new(big.Int).Sub(balance, value)
		log.Debug("CheckSubTokenBalance settle balance: SUB TOKEN BALANCE ", "token", token.String(), "address", addr.String(), "balance", balance, "value", value, "newBalance", newBalance)
		return newBalance, nil
	} else {
		return nil, errors.Errorf("token %s isn't exist", token.String())
	}
}

func CheckAddTokenBalance(addr common.Address, value *big.Int, token common.Address, statedb *state.StateDB, mapBalances map[common.Address]map[common.Address]*big.Int) (*big.Int, error) {
	// TOMO native
	if token.String() == common.TomoNativeAddress {
		var balance *big.Int
		if value := mapBalances[token][addr]; value != nil {
			balance = value
		} else {
			balance = statedb.GetBalance(addr)
		}
		newBalance := new(big.Int).Add(balance, value)
		log.Debug("CheckAddTokenBalance settle balance: ADD TOMO NATIVE BALANCE ", "token", token.String(), "address", addr.String(), "balance", balance, "value", value, "newBalance", newBalance)
		return newBalance, nil
	}
	// TRC tokens
	if statedb.Exist(token) {
		var balance *big.Int
		if value := mapBalances[token][addr]; value != nil {
			balance = value
		} else {
			slot := TokenMappingSlot["balances"]
			locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
			balance = statedb.GetState(token, locHash).Big()
		}
		newBalance := new(big.Int).Add(balance, value)
		log.Debug("CheckAddTokenBalance settle balance: ADD TOKEN BALANCE ", "token", token.String(), "address", addr.String(), "balance", balance, "value", value, "newBalance", newBalance)
		if common.BigToHash(newBalance).Big().Cmp(newBalance) != 0 {
			return nil, fmt.Errorf("Overflow when try add token balance , max is 2^256 , balance : %v , value:%v ", balance, value)
		} else {
			return newBalance, nil
		}
	} else {
		return nil, errors.Errorf("token %s isn't exist", token.String())
	}
}
func GetTokenBalance(addr common.Address, token common.Address, statedb *state.StateDB) *big.Int {
	// TOMO native
	if token.String() == common.TomoNativeAddress {
		return statedb.GetBalance(addr)
	}
	// TRC tokens
	if statedb.Exist(token) {
		slot := TokenMappingSlot["balances"]
		locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
		return statedb.GetState(token, locHash).Big()
	} else {
		return common.Big0
	}
}

func SetTokenBalance(addr common.Address, balance *big.Int, token common.Address, statedb *state.StateDB) error {
	// TOMO native
	if token.String() == common.TomoNativeAddress {
		statedb.SetBalance(addr, balance)
		return nil
	}

	// TRC tokens
	if statedb.Exist(token) {
		slot := TokenMappingSlot["balances"]
		locHash := common.BigToHash(GetLocMappingAtKey(addr.Hash(), slot))
		statedb.SetState(token, locHash, common.BigToHash(balance))
		return nil
	} else {
		return errors.Errorf("token %s isn't exist", token.String())
	}
}
