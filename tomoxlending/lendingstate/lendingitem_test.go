package lendingstate

import (
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/ethdb"
	"math/big"
	"testing"
	"time"
)

func TestLendingItem_VerifyLendingSide(t *testing.T) {
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{
		{"wrong side", &LendingItem{Side: "GIVE"}, true},
		{"side: borrowing", &LendingItem{Side: Borrowing}, false},
		{"side: investing", &LendingItem{Side: Investing}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LendingItem{
				Side: tt.fields.Side,
			}
			if err := l.VerifyLendingSide(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyLendingSide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLendingItem_VerifyLendingInterest(t *testing.T) {
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{
		{"no interest information", &LendingItem{}, true},
		{"negative interest", &LendingItem{Interest: big.NewInt(-1)}, true},
		{"zero interest", &LendingItem{Interest: Zero}, true},
		{"positive interest", &LendingItem{Interest: big.NewInt(2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LendingItem{
				Interest: tt.fields.Interest,
			}
			if err := l.VerifyLendingInterest(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyLendingSide() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLendingItem_VerifyLendingQuantity(t *testing.T) {
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{
		{"no quantity information", &LendingItem{}, true},
		{"negative quantity", &LendingItem{Quantity: big.NewInt(-1)}, true},
		{"zero quantity", &LendingItem{Quantity: Zero}, true},
		{"positive quantity", &LendingItem{Quantity: big.NewInt(2)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LendingItem{
				Quantity: tt.fields.Quantity,
			}
			if err := l.VerifyLendingQuantity(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyLendingQuantity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLendingItem_VerifyLendingType(t *testing.T) {
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{
		{"type: stop limit", &LendingItem{Type: "stop limit"}, true},
		{"type: take profit", &LendingItem{Type: "take profit"}, true},
		{"type: limit", &LendingItem{Type: Limit}, false},
		{"type: market", &LendingItem{Type: Market}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LendingItem{
				Type: tt.fields.Type,
			}
			if err := l.VerifyLendingType(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyLendingType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLendingItem_VerifyLendingStatus(t *testing.T) {
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{

		{"status: new", &LendingItem{Status: LendingStatusNew}, false},
		{"status: open", &LendingItem{Status: LendingStatusOpen}, true},
		{"status: partial_filled", &LendingItem{Status: LendingStatusPartialFilled}, true},
		{"status: filled", &LendingItem{Status: LendingStatusFilled}, true},
		{"status: cancelled", &LendingItem{Status: LendingStatusCancelled}, false},
		{"status: rejected", &LendingItem{Status: LendingStatusReject}, true},
		{"status: deposit", &LendingItem{Status: Deposit}, false},
		{"status: payment", &LendingItem{Status: Payment}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LendingItem{
				Status: tt.fields.Status,
			}
			if err := l.VerifyLendingStatus(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyLendingStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func SetFee(statedb *state.StateDB, coinbase common.Address, feeRate *big.Int) {
	locRelayerState := state.GetLocMappingAtKey(coinbase.Hash(), LendingRelayerListSlot)
	locHash := common.BytesToHash(new(big.Int).Add(locRelayerState, LendingRelayerStructSlots["fee"]).Bytes())
	statedb.SetState(common.HexToAddress(common.LendingRegistrationSMC), locHash, common.BigToHash(feeRate))
}

func SetCollateralDetail(statedb *state.StateDB, token common.Address, depositRate *big.Int, liquidationRate *big.Int, price *big.Int) {
	collateralState := GetLocMappingAtKey(token.Hash(), CollateralMapSlot)
	locDepositRate := state.GetLocOfStructElement(collateralState, CollateralStructSlots["depositRate"])
	locLiquidationRate := state.GetLocOfStructElement(collateralState, CollateralStructSlots["liquidationRate"])
	locCollateralPrice := state.GetLocOfStructElement(collateralState, CollateralStructSlots["price"])
	statedb.SetState(common.HexToAddress(common.LendingRegistrationSMC), locDepositRate, common.BigToHash(depositRate))
	statedb.SetState(common.HexToAddress(common.LendingRegistrationSMC), locLiquidationRate, common.BigToHash(liquidationRate))
	statedb.SetState(common.HexToAddress(common.LendingRegistrationSMC), locCollateralPrice, common.BigToHash(price))
}

func TestVerifyBalance(t *testing.T) {
	db, _ := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	relayer := common.HexToAddress("0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e")
	uAddr := common.HexToAddress("0xDeE6238780f98c0ca2c2C28453149bEA49a3Abc9")
	lendingToken := common.HexToAddress("0xd9bb01454c85247B2ef35BB5BE57384cC275a8cf")    // USD
	collateralToken := common.HexToAddress("0x4d7eA2cE949216D6b120f3AA10164173615A2b6C") // BTC

	SetFee(statedb, relayer, big.NewInt(100))
	SetCollateralDetail(statedb, collateralToken, big.NewInt(150), big.NewInt(110), big.NewInt(8000)) // BTC price: 8k USD

	// have 10k USD
	statedb.GetOrNewStateObject(lendingToken)
	if err := SetTokenBalance(uAddr, EtherToWei(big.NewInt(10000)), lendingToken, statedb); err != nil {
		t.Error(err.Error())
	}

	// have 2 BTC
	statedb.GetOrNewStateObject(collateralToken)
	if err := SetTokenBalance(uAddr, EtherToWei(big.NewInt(2)), collateralToken, statedb); err != nil {
		t.Error(err.Error())
	}
	lendingdb, _ := ethdb.NewMemDatabase()
	stateCache := NewDatabase(lendingdb)
	lendingstatedb, _ := New(EmptyRoot, stateCache)

	// insert lendingItem1 for testing cancel (side investing)
	lendingItem1 := LendingItem{
		Quantity:        EtherToWei(big.NewInt(11000000000)),
		Interest:        big.NewInt(10),
		Side:            Investing,
		Type:            Limit,
		LendingToken:    lendingToken,
		CollateralToken: collateralToken,
		FilledAmount:    nil,
		Status:          LendingStatusOpen,
		Relayer:         relayer,
		Term:            uint64(30),
		UserAddress:     uAddr,
		Signature:       nil,
		Hash:            common.Hash{},
		TxHash:          common.Hash{},
		Nonce:           nil,
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
		LendingId:       uint64(1),
		ExtraData:       "",
	}
	lendingstatedb.InsertLendingItem(GetLendingOrderBookHash(lendingItem1.LendingToken, lendingItem1.Term), common.BigToHash(new(big.Int).SetUint64(lendingItem1.LendingId)), lendingItem1)

	// insert lendingItem2 for testing cancel (side borrowing)
	lendingItem2 := LendingItem{
		Quantity:        EtherToWei(big.NewInt(8000)),
		Interest:        big.NewInt(10),
		Side:            Borrowing,
		Type:            Limit,
		LendingToken:    lendingToken,
		CollateralToken: collateralToken,
		FilledAmount:    nil,
		Status:          LendingStatusOpen,
		Relayer:         relayer,
		Term:            uint64(30),
		UserAddress:     uAddr,
		Signature:       nil,
		Hash:            common.Hash{},
		TxHash:          common.Hash{},
		Nonce:           nil,
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
		LendingId:       uint64(2),
		ExtraData:       "",
	}
	lendingstatedb.InsertLendingItem(GetLendingOrderBookHash(lendingItem2.LendingToken, lendingItem2.Term), common.BigToHash(new(big.Int).SetUint64(lendingItem2.LendingId)), lendingItem2)

	// insert lendingTrade for testing deposit (side: borrowing)
	lendingstatedb.InsertTradingItem(
		GetLendingOrderBookHash(lendingItem2.LendingToken, lendingItem2.Term),
		uint64(1),
		LendingTrade{
			TradeId:         uint64(1),
			CollateralToken: collateralToken,
			LendingToken:    lendingToken,
			Borrower:        uAddr,
			Amount:          EtherToWei(big.NewInt(8000)),
			LiquidationTime: uint64(time.Now().AddDate(0, 1, 0).UnixNano()),
		},
	)

	// make a big lendingTrade to test case: not enough balance to process payment
	lendingstatedb.InsertTradingItem(
		GetLendingOrderBookHash(lendingItem2.LendingToken, lendingItem2.Term),
		uint64(2),
		LendingTrade{
			TradeId:         uint64(2),
			CollateralToken: collateralToken,
			LendingToken:    lendingToken,
			Borrower:        uAddr,
			Amount:          EtherToWei(big.NewInt(20000)), // user have 10k USD, expect: fail
			LiquidationTime: uint64(time.Now().AddDate(0, 1, 0).UnixNano()),
		},
	)
	tests := []struct {
		name    string
		fields  *LendingItem
		wantErr bool
	}{
		{"Investor doesn't have enough balance. side: investing, quantity 11k USD",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Investing,
				Status:          LendingStatusNew,
				Quantity:        EtherToWei(big.NewInt(11000)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
			},
			true,
		},
		{"Investor has enough balance. side: investing, quantity 10k USD",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Investing,
				Status:          LendingStatusNew,
				Quantity:        EtherToWei(big.NewInt(10000)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
			},
			false,
		},
		{"Investor cancel lendingItem",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Investing,
				Status:          LendingStatusCancelled,
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				Term:            lendingItem1.Term,
				LendingId:       uint64(1),
			},
			true,
		},
		{"Invalid status",
			&LendingItem{
				Side:   Investing,
				Status: Deposit,
			},
			true,
		},
		// have 2BTC = 16k USD => max borrow = 16 / 1.5 = 10.66
		{"Borrower doesn't have enough balance. side: borrowing, quantity 12k USD",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          LendingStatusNew,
				Quantity:        EtherToWei(big.NewInt(12000)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
			},
			true,
		},
		// have 2BTC = 16k USD => max borrow = 16 / 1.5 = 10.66
		{"Borrower has enough balance. side: borrowing, quantity 10k USD",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          LendingStatusNew,
				Quantity:        EtherToWei(big.NewInt(10000)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
			},
			false,
		},
		{"Borrower has enough balance to pay cancel fee. side: borrowing",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          LendingStatusCancelled,
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				Term:            lendingItem2.Term,
				LendingId:       uint64(2),
			},
			false,
		},
		{"Make a deposit to an empty LendingTrade.",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          Deposit,
				Quantity:        EtherToWei(big.NewInt(1)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				ExtraData:       common.BigToAddress(big.NewInt(0)).Hex(),
			},
			true,
		},
		// have 2BTC. make deposit 1 BTC
		{"Borrower has enough balance to make a deposit. side: borrowing",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          Deposit,
				Quantity:        EtherToWei(big.NewInt(1)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				Term:            uint64(30),
				ExtraData:       common.Uint64ToHash(1).Hex(),
			},
			false,
		},
		{"Make a payment to an empty LendingTrade.",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          Payment,
				Quantity:        EtherToWei(big.NewInt(1)),
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				ExtraData:       common.BigToAddress(big.NewInt(0)).Hex(),
			},
			true,
		},
		// have 10k USDT
		{"Borrower has enough balance to make a payment transaction. side: borrowing",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          Payment,
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				Term:            uint64(30),
				ExtraData:       common.Uint64ToHash(1).Hex(),
			},
			false,
		},
		// have 10k USDT
		{"Borrower doesn't haave enough balance to make a payment transaction. side: borrowing",
			&LendingItem{
				UserAddress:     uAddr,
				Relayer:         relayer,
				Side:            Borrowing,
				Status:          Payment,
				LendingToken:    lendingToken,
				CollateralToken: collateralToken,
				Term:            uint64(30),
				ExtraData:       common.Uint64ToHash(2).Hex(),
			},
			true,
		},
		{"Invalid status",
			&LendingItem{
				Side:   Borrowing,
				Status: LendingStatusOpen,
			},
			true,
		},
		{"Invalid side",
			&LendingItem{
				Side: "abc",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyBalance(statedb,
				lendingstatedb,
				tt.fields.Side,
				tt.fields.Status,
				tt.fields.UserAddress,
				tt.fields.Relayer,
				tt.fields.LendingToken,
				tt.fields.CollateralToken,
				tt.fields.Quantity,
				EtherToWei(big.NewInt(1)),
				EtherToWei(big.NewInt(1)),
				EtherToWei(big.NewInt(2)),    // TOMO price: 0.5 USD => USD/TOMO = 2
				EtherToWei(big.NewInt(8000)), // BTC = 8000 USD
				tt.fields.Term,
				tt.fields.LendingId,
				common.HexToHash(tt.fields.ExtraData),
			); (err != nil) != tt.wantErr {
				t.Errorf("VerifyBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
