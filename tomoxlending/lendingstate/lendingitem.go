package lendingstate

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/crypto/sha3"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"math/big"
	"strconv"
	"time"
)

const (
	Investing                  = "INVESTING"
	Borrowing                  = "BORROWING"
	Deposit                    = "DEPOSIT"
	Payment                    = "PAYMENT"
	LendingStatusNew           = "NEW"
	LendingStatusOpen          = "OPEN"
	LendingStatusReject        = "REJECTED"
	LendingStatusFilled        = "FILLED"
	LendingStatusPartialFilled = "PARTIAL_FILLED"
	LendingStatusCancelled     = "CANCELLED"
	Market                     = "MO"
	Limit                      = "LO"
)

var ValidInputLendingStatus = map[string]bool{
	Deposit:                true,
	Payment:                true,
	LendingStatusNew:       true,
	LendingStatusCancelled: true,
}

// Signature struct
type Signature struct {
	V byte        `bson:"v" json:"v"`
	R common.Hash `bson:"r" json:"r"`
	S common.Hash `bson:"s" json:"s"`
}

type SignatureRecord struct {
	V byte   `bson:"v" json:"v"`
	R string `bson:"r" json:"r"`
	S string `bson:"s" json:"s"`
}

type LendingItem struct {
	Quantity        *big.Int       `bson:"quantity" json:"quantity"`
	Interest        *big.Int       `bson:"interest" json:"interest"`
	Side            string         `bson:"side" json:"side"` // INVESTING/BORROWING
	Type            string         `bson:"type" json:"type"` // LIMIT/MARKET
	LendingToken    common.Address `bson:"lendingToken" json:"lendingToken"`
	CollateralToken common.Address `bson:"collateralToken" json:"collateralToken"`
	FilledAmount    *big.Int       `bson:"filledAmount" json:"filledAmount"`
	Status          string         `bson:"status" json:"status"`
	Relayer         common.Address `bson:"relayer" json:"relayer"`
	Term            uint64         `bson:"term" json:"term"`
	UserAddress     common.Address `bson:"userAddress" json:"userAddress"`
	Signature       *Signature     `bson:"signature" json:"signature"`
	Hash            common.Hash    `bson:"hash" json:"hash"`
	TxHash          common.Hash    `bson:"txHash" json:"txHash"`
	Nonce           *big.Int       `bson:"nonce" json:"nonce"`
	CreatedAt       time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time      `bson:"updatedAt" json:"updatedAt"`
	LendingId       uint64         `bson:"tradeId" json:"tradeId"`
	ExtraData       string         `bson:"extraData" json:"extraData"`
}

type LendingItemBSON struct {
	Quantity        string           `bson:"quantity" json:"quantity"`
	Interest        string           `bson:"interest" json:"interest"`
	Side            string           `bson:"side" json:"side"` // INVESTING/BORROWING
	Type            string           `bson:"type" json:"type"` // LIMIT/MARKET
	LendingToken    string           `bson:"lendingToken" json:"lendingToken"`
	CollateralToken string           `bson:"collateralToken" json:"collateralToken"`
	FilledAmount    string           `bson:"filledAmount" json:"filledAmount"`
	Status          string           `bson:"status" json:"status"`
	Relayer         string           `bson:"relayer" json:"relayer"`
	Term            string           `bson:"term" json:"term"`
	UserAddress     string           `bson:"userAddress" json:"userAddress"`
	Signature       *SignatureRecord `bson:"signature" json:"signature"`
	Hash            string           `bson:"hash" json:"hash"`
	TxHash          string           `bson:"txHash" json:"txHash"`
	Nonce           string           `bson:"nonce" json:"nonce"`
	CreatedAt       time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `bson:"updatedAt" json:"updatedAt"`
	LendingId       string           `bson:"tradeId" json:"tradeId"`
	ExtraData       string           `bson:"extraData" json:"extraData"`
}

func (l *LendingItem) GetBSON() (interface{}, error) {
	lr := LendingItemBSON{
		Quantity:        l.Quantity.String(),
		Interest:        l.Interest.String(),
		Side:            l.Side,
		Type:            l.Type,
		LendingToken:    l.LendingToken.Hex(),
		CollateralToken: l.CollateralToken.Hex(),
		Status:          l.Status,
		Relayer:         l.Relayer.Hex(),
		Term:            strconv.FormatUint(l.Term, 10),
		UserAddress:     l.UserAddress.Hex(),
		Hash:            l.Hash.Hex(),
		TxHash:          l.TxHash.Hex(),
		Nonce:           l.Nonce.String(),
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
		LendingId:       strconv.FormatUint(l.LendingId, 10),
		ExtraData:       l.ExtraData,
	}

	if l.FilledAmount != nil {
		lr.FilledAmount = l.FilledAmount.String()
	}

	if l.Signature != nil {
		lr.Signature = &SignatureRecord{
			V: l.Signature.V,
			R: l.Signature.R.Hex(),
			S: l.Signature.S.Hex(),
		}
	}

	return lr, nil
}

func (l *LendingItem) SetBSON(raw bson.Raw) error {
	decoded := new(LendingItemBSON)

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}
	if decoded.Quantity != "" {
		l.Quantity = ToBigInt(decoded.Quantity)
	}
	l.Interest = ToBigInt(decoded.Interest)
	l.Side = decoded.Side
	l.Type = decoded.Type
	l.LendingToken = common.HexToAddress(decoded.LendingToken)
	l.CollateralToken = common.HexToAddress(decoded.CollateralToken)
	l.FilledAmount = ToBigInt(decoded.FilledAmount)
	l.Status = decoded.Status
	l.Relayer = common.HexToAddress(decoded.Relayer)
	term, err := strconv.ParseInt(decoded.Term, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse lendingItem.term. Err: %v", err)
	}
	l.Term = uint64(term)
	l.UserAddress = common.HexToAddress(decoded.UserAddress)

	if decoded.Signature != nil {
		l.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	l.Hash = common.HexToHash(decoded.Hash)
	l.TxHash = common.HexToHash(decoded.TxHash)
	l.Nonce = ToBigInt(decoded.Nonce)

	l.CreatedAt = decoded.CreatedAt
	l.UpdatedAt = decoded.UpdatedAt
	lendingId, err := strconv.ParseInt(decoded.LendingId, 10, 64)
	if err != nil {
		return err
	}
	l.LendingId = uint64(lendingId)
	l.ExtraData = decoded.ExtraData
	return nil
}

func (l *LendingItem) VerifyLendingItem(state *state.StateDB) error {
	if err := l.VerifyLendingStatus(); err != nil {
		return err
	}
	if l.Status == LendingStatusNew {
		if err := l.VerifyLendingType(); err != nil {
			return err
		}
		if err := l.VerifyLendingQuantity(); err != nil {
			return err
		}
		if err := l.VerifyLendingSide(); err != nil {
			return err
		}
		if l.Type == Limit {
			if err := l.VerifyLendingInterest(); err != nil {
				return err
			}
		}

	}
	if !IsValidRelayer(state, l.Relayer) {
		return fmt.Errorf("VerifyLendingItem: invalid relayer. address: %s", l.Relayer.Hex())
	}
	if err := l.VerifyLendingSignature(); err != nil {
		return err
	}
	return nil
}

func (l *LendingItem) VerifyLendingSide() error {
	if l.Side != Borrowing && l.Side != Investing {
		return fmt.Errorf("VerifyLendingSide: invalid side . Side: %s", l.Side)
	}
	return nil
}

func (l *LendingItem) VerifyLendingInterest() error {
	if l.Interest == nil || l.Interest.Sign() <= 0 {
		return fmt.Errorf("VerifyLendingInterest: invalid interest. Interest: %v", l.Interest)
	}
	return nil
}

func (l *LendingItem) VerifyLendingQuantity() error {
	if l.Quantity == nil || l.Quantity.Sign() <= 0 {
		return fmt.Errorf("VerifyLendingQuantity: invalid quantity. Quantity: %v", l.Quantity)
	}
	return nil
}

func (l *LendingItem) VerifyLendingType() error {
	if l.Type != Market && l.Type != Limit {
		return fmt.Errorf("VerifyLendingType: invalid lending type. Type: %s", l.Type)
	}
	return nil
}

func (l *LendingItem) VerifyLendingStatus() error {
	if valid, ok := ValidInputLendingStatus[l.Status]; ok && valid {
		return fmt.Errorf("VerifyLendingStatus: invalid lending status. Status: %s", l.Status)
	}
	return nil
}

func (l *LendingItem) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	if l.Status == LendingStatusNew {
		sha.Write(l.Relayer.Bytes())
		sha.Write(l.UserAddress.Bytes())
		sha.Write(l.LendingToken.Bytes())
		sha.Write(l.CollateralToken.Bytes())
		sha.Write([]byte(strconv.FormatInt(int64(l.Term), 10)))
		sha.Write(common.BigToHash(l.Quantity).Bytes())
		if l.Type == Limit {
			if l.Interest != nil {
				sha.Write(common.BigToHash(l.Interest).Bytes())
			}
		}
		sha.Write(common.BigToHash(l.EncodedSide()).Bytes())
		sha.Write([]byte(l.Status))
		sha.Write([]byte(l.Type))
		sha.Write(common.BigToHash(l.Nonce).Bytes())
	} else if l.Status == LendingStatusCancelled {
		sha.Write(l.Hash.Bytes())
		sha.Write(common.BigToHash(l.Nonce).Bytes())
		sha.Write(l.UserAddress.Bytes())
		sha.Write(common.BigToHash(big.NewInt(int64(l.LendingId))).Bytes())
		sha.Write([]byte(l.Status))
		sha.Write(l.Relayer.Bytes())
		sha.Write(l.LendingToken.Bytes())
		sha.Write(l.CollateralToken.Bytes())
	} else {
		return common.Hash{}
	}

	return common.BytesToHash(sha.Sum(nil))
}

func (l *LendingItem) EncodedSide() *big.Int {
	if l.Side == Borrowing {
		return big.NewInt(0)
	}
	return big.NewInt(1)
}

//verify signatures
func (l *LendingItem) VerifyLendingSignature() error {
	bigstr := l.Nonce.String()
	n, err := strconv.ParseInt(bigstr, 10, 64)
	if err != nil {
		return fmt.Errorf("verify lending item: invalid signature")
	}
	V := big.NewInt(int64(l.Signature.V))
	R := l.Signature.R.Big()
	S := l.Signature.S.Big()

	//(nonce uint64, quantity *big.Int, interest, duration uint64, relayerAddress, userAddress, lendingToken, collateralToken common.Address, status, side, typeLending string, hash common.Hash, id uint64
	tx := types.NewLendingTransaction(uint64(n), l.Quantity, l.Interest.Uint64(), l.Term, l.Relayer, l.UserAddress,
		l.LendingToken, l.CollateralToken, l.Status, l.Side, l.Type, l.Hash, l.LendingId)
	tx.ImportSignature(V, R, S)
	from, _ := types.LendingSender(types.LendingTxSigner{}, tx)
	if from != tx.UserAddress() {
		return fmt.Errorf("verify lending item: invalid signature")
	}
	return nil
}

func VerifyBalance(statedb *state.StateDB, lendingStateDb *LendingStateDB, tradingStateDb *tradingstate.TradingStateDB, l *types.LendingTransaction, lendingTokenDecimal, collateralTokenDecimal *big.Int) error {
	switch l.Side() {
	case Investing:
		if balance := tradingstate.GetTokenBalance(l.UserAddress(), l.LendingToken(), statedb); balance.Cmp(l.Quantity()) < 0 {
			return fmt.Errorf("VerifyBalance: investor doesn't have enough lendingToken. User: %s. Token: %s. Expected: %v. Have: %v", l.UserAddress().Hex(), l.LendingToken().Hex(), l.Quantity(), balance)
		}
		return nil
	case Borrowing:
		//balance >= depositRate * lendingAMount + fee
		//make sure fee > 0.01 TOMO
		switch l.Status() {
		case LendingStatusNew:
			//TODO:@nguyennguyen
			return nil
		case LendingStatusCancelled:
			//TODO:@nguyennguyen
			return nil
		case Deposit:
			//TODO:@nguyennguyen
			return nil
		case Payment:
			//TODO:@nguyennguyen
			return nil
		default:
			return nil
		}
	default:
		return fmt.Errorf("VerifyBalance: unknown lending side")
	}
}
