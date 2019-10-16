package tomox

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
)

func IsValidRelayer(statedb *state.StateDB, address common.Address) bool {
	slot := tomox_state.RelayerMappingSlot["RELAYER_LIST"]
	locRelayerState := tomox_state.GetLocMappingAtKey(address.Hash(), slot)

	locBigDeposit := new(big.Int).SetUint64(uint64(0)).Add(locRelayerState, tomox_state.RelayerStructMappingSlot["_deposit"])
	locHashDeposit := common.BigToHash(locBigDeposit)
	balance := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), locHashDeposit).Big()
	if balance.Cmp(new(big.Int).SetUint64(uint64(0))) > 0 {
		return true
	}
	log.Debug("Balance of relayer is not enough", "relayer", address.String(), "balance", balance)
	return false
}

func GetTokenBalance(statedb *state.StateDB, address common.Address, contractAddr common.Address) *big.Int {
	slot := tomox_state.TokenMappingSlot["balances"]
	locBalance := tomox_state.GetLocMappingAtKey(address.Hash(), slot)

	ret := statedb.GetState(contractAddr, common.BigToHash(locBalance))
	return ret.Big()
}

// verify orderbook, orderTrees before running matching engine
func (tx TxDataMatch) VerifyOldTomoXState(ob *OrderBook) error {
	// verify orderbook
	if hash, err := ob.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.ObOld.Bytes()) {
		log.Error("wrong old orderbook", "expected", hex.EncodeToString(tx.ObOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderBookHashNotMatch
	}

	// verify order trees
	// bidTree tree
	bidTree := ob.Bids
	if hash, err := bidTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.BidOld.Bytes()) {
		log.Error("wrong old bid tree", "expected", hex.EncodeToString(tx.BidOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderTreeHashNotMatch
	}
	// askTree tree
	askTree := ob.Asks
	if hash, err := askTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.AskOld.Bytes()) {
		log.Error("wrong old ask tree", "expected", hex.EncodeToString(tx.AskOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderTreeHashNotMatch
	}
	return nil
}

// verify orderbook, orderTrees after running matching engine
func (tx TxDataMatch) VerifyNewTomoXState(ob *OrderBook) error {
	// verify orderbook
	if hash, err := ob.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.ObNew.Bytes()) {
		log.Error("wrong new orderbook", "expected", hex.EncodeToString(tx.ObNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderBookHashNotMatch
	}

	// verify order trees
	// bidTree tree
	bidTree := ob.Bids
	if hash, err := bidTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.BidNew.Bytes()) {
		log.Error("wrong new bid tree", "expected", hex.EncodeToString(tx.BidNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderTreeHashNotMatch
	}
	// askTree tree
	askTree := ob.Asks
	if hash, err := askTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.AskNew.Bytes()) {
		log.Error("wrong new ask tree", "expected", hex.EncodeToString(tx.AskNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return tomox_state.ErrOrderTreeHashNotMatch
	}
	return nil
}

func (tx TxDataMatch) DecodeOrder() (*tomox_state.OrderItem, error) {
	order := &tomox_state.OrderItem{}
	if err := DecodeBytesItem(tx.Order, order); err != nil {
		return order, err
	}
	return order, nil
}

func (tx TxDataMatch) GetTrades() []map[string]string {
	return tx.Trades
}

func (tx TxDataMatch) GetRejectedOrders() []map[string]string {
	return tx.RejectedOders
}
