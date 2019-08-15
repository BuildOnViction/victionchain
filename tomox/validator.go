package tomox

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
	"bytes"
	"encoding/hex"
)

var (
	// errors
	errFutureOrder             = errors.New("verify matched order: future order")
	errNoTimestamp             = errors.New("verify matched order: no timestamp")
	errWrongHash               = errors.New("verify matched order: wrong hash")
	errInvalidSignature        = errors.New("verify matched order: invalid signature")
	errNotEnoughBalance        = errors.New("verify matched order: not enough balance")
	errInvalidPrice            = errors.New("verify matched order: invalid price")
	errInvalidQuantity         = errors.New("verify matched order: invalid quantity")
	errInvalidRelayer          = errors.New("verify matched order: invalid relayer")
	errInvalidOrderType        = errors.New("verify matched order: unsupported order type")
	errInvalidOrderSide        = errors.New("verify matched order: invalid order side")
	errOrderBookHashNotMatch   = errors.New("verify matched order: orderbook hash not match")
	errOrderTreeHashNotMatch   = errors.New("verify matched order: ordertree hash not match")

	// supported order types
	MatchingOrderType = map[string]bool{
		Market: true,
		Limit:  true,
	}
)

// verify orderItem
func (o *OrderItem) VerifyMatchedOrder(state *state.StateDB) error {
	if err := o.verifyTimestamp(); err != nil {
		return err
	}
	if err := o.verifyOrderSide(); err != nil {
		return err
	}
	if err := o.verifyOrderType(); err != nil {
		return err
	}
	// TODO: for testing without relayer, ignore the rest
	// TODO: remove it
	return nil

	if err := o.verifyRelayer(state); err != nil {
		return err
	}
	if err := o.verifyBalance(state); err != nil {
		return err
	}
	if err := o.verifySignature(); err != nil {
		return err
	}
	return nil
}

// verify token balance make sure that user has enough balance
func (o *OrderItem) verifyBalance(state *state.StateDB) error {
	orderValueByQuoteToken := Zero()
	balance := Zero()
	tokenAddr := common.Address{}

	if o.Price == nil || o.Price.Cmp(Zero()) <= 0 {
		return errInvalidPrice
	}
	if o.Quantity == nil || o.Quantity.Cmp(Zero()) <= 0 {
		return errInvalidQuantity
	}
	if o.Side == Bid {
		tokenAddr = o.QuoteToken
		orderValueByQuoteToken = orderValueByQuoteToken.Mul(o.Quantity, o.Price)
	} else {
		tokenAddr = o.BaseToken
		orderValueByQuoteToken = o.Quantity
	}
	if tokenAddr == (common.Address{}) {
		// native TOMO
		balance = state.GetBalance(o.UserAddress)
	} else {
		// query balance from tokenContract
		balance = GetTokenBalance(state, o.UserAddress, tokenAddr)
	}
	if balance.Cmp(orderValueByQuoteToken) < 0 {
		return errNotEnoughBalance
	}
	return nil
}

// verify whether the exchange applies to become relayer
func (o *OrderItem) verifyRelayer(state *state.StateDB) error {
	if !IsValidRelayer(state, o.ExchangeAddress) {
		return errInvalidRelayer
	}
	return nil
}

// following: https://github.com/tomochain/tomox-sdk/blob/master/types/order.go#L125
func (o *OrderItem) computeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.ExchangeAddress.Bytes())
	sha.Write(o.UserAddress.Bytes())
	sha.Write(o.BaseToken.Bytes())
	sha.Write(o.QuoteToken.Bytes())
	sha.Write(common.BigToHash(o.Quantity).Bytes())
	sha.Write(common.BigToHash(o.Price).Bytes())
	sha.Write(common.BigToHash(o.encodedSide()).Bytes())
	sha.Write(common.BigToHash(o.Nonce).Bytes())
	sha.Write(common.BigToHash(o.MakeFee).Bytes())
	sha.Write(common.BigToHash(o.TakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

//verify signatures
func (o *OrderItem) verifySignature() error {
	var (
		hash           common.Hash
		err            error
		signatureBytes []byte
	)
	hash = o.computeHash()
	if hash != o.Hash {
		return errWrongHash
	}
	signatureBytes = append(signatureBytes, o.Signature.R.Bytes()...)
	signatureBytes = append(signatureBytes, o.Signature.S.Bytes()...)
	signatureBytes = append(signatureBytes, o.Signature.V-27)
	pubkey, err := crypto.Ecrecover(hash.Bytes(), signatureBytes)
	if err != nil {
		return err
	}
	var userAddress common.Address
	copy(userAddress[:], crypto.Keccak256(pubkey[1:])[12:])
	if userAddress != o.UserAddress {
		return errInvalidSignature
	}
	return nil
}

// verify order type
func (o *OrderItem) verifyOrderType() error {
	if _, ok := MatchingOrderType[o.Type]; !ok {
		return errInvalidOrderType
	}
	return nil
}

//verify order side
func (o *OrderItem) verifyOrderSide() error {

	if o.Side != Bid && o.Side != Ask {
		return errInvalidOrderSide
	}
	return nil
}

//verify timestamp
func (o *OrderItem) verifyTimestamp() error {
	// check timestamp of buyOrder
	if o.CreatedAt == 0 || o.UpdatedAt == 0 {
		return errNoTimestamp
	}
	if o.CreatedAt > uint64(time.Now().Unix()) || o.UpdatedAt > uint64(time.Now().Unix()) {
		return errFutureOrder
	}
	return nil
}

func (o *OrderItem) encodedSide() *big.Int {
	if o.Side == Bid {
		return big.NewInt(0)
	}
	return big.NewInt(1)
}

func IsValidRelayer(statedb *state.StateDB, address common.Address) bool {
	slotHash := common.BigToHash(new(big.Int).SetUint64(RelayerMappingSlot["RELAYER_LIST"]))
	retByte := crypto.Keccak256(address.Bytes(), slotHash.Bytes())
	locRelayerState := new(big.Int)
	locRelayerState.SetBytes(retByte)

	ret := statedb.GetState(common.HexToAddress(common.RelayerRegistrationSMC), common.BigToHash(locRelayerState))
	if ret.Big().Cmp(new(big.Int).SetUint64(uint64(0))) > 0 {
		return true
	}
	return false
}

func GetTokenBalance(statedb *state.StateDB, address common.Address, contractAddr common.Address) *big.Int {
	slotHash := common.BigToHash(new(big.Int).SetUint64(TokenMappingSlot["balances"]))
	retByte := crypto.Keccak256(address.Bytes(), slotHash.Bytes())
	locBalance := new(big.Int)
	locBalance.SetBytes(retByte)

	ret := statedb.GetState(contractAddr, common.BigToHash(locBalance))
	return ret.Big()
}

// verify orderbook, orderTrees before running matching engine
func (tx TxDataMatch) VerifyOldTomoXState(ob *OrderBook) error {
	// verify orderbook
	if hash, err := ob.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.ObOld.Bytes()) {
		log.Error("wrong old orderbook", "expected", hex.EncodeToString(tx.ObOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderBookHashNotMatch
	}

	// verify order trees
	// bidTree tree
	bidTree := ob.Bids
	if hash, err := bidTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.BidOld.Bytes()) {
		log.Error("wrong old bid tree", "expected",  hex.EncodeToString(tx.BidOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderTreeHashNotMatch
	}
	// askTree tree
	askTree := ob.Asks
	if hash, err := askTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.AskOld.Bytes()) {
		log.Error("wrong old ask tree", "expected", hex.EncodeToString(tx.AskOld.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderTreeHashNotMatch
	}
	return nil
}

// verify orderbook, orderTrees after running matching engine
func (tx TxDataMatch) VerifyNewTomoXState(ob *OrderBook) error {
	// verify orderbook
	if hash, err := ob.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.ObNew.Bytes()) {
		log.Error("wrong new orderbook", "expected", hex.EncodeToString(tx.ObNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderBookHashNotMatch
	}

	// verify order trees
	// bidTree tree
	bidTree := ob.Bids
	if hash, err := bidTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.BidNew.Bytes()) {
		log.Error("wrong new bid tree", "expected", hex.EncodeToString(tx.BidNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderTreeHashNotMatch
	}
	// askTree tree
	askTree := ob.Asks
	if hash, err := askTree.Hash(); err != nil || !bytes.Equal(hash.Bytes(), tx.AskNew.Bytes()) {
		log.Error("wrong new ask tree", "expected", hex.EncodeToString(tx.AskNew.Bytes()), "actual", hex.EncodeToString(hash.Bytes()), "err", err)
		return errOrderTreeHashNotMatch
	}
	return nil
}

func (tx TxDataMatch) DecodeOrder() (*OrderItem, error) {
	order := &OrderItem{}
	if err := DecodeBytesItem(tx.Order, order); err != nil {
		return order, err
	}
	return order, nil
}

func (tx TxDataMatch) GetTrades() []map[string]string {
	return tx.Trades
}

func (tx TxDataMatch) DecodeOrderInBook() (*OrderItem, error) {
	if len(tx.OrderInBook) == 0 {
		return nil, nil
	}
	orderInBook := &OrderItem{}
	if err := DecodeBytesItem(tx.OrderInBook, orderInBook); err != nil {
		return orderInBook, err
	}
	return orderInBook, nil
}