package lendingstate

import (
	"math/big"
	"testing"
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
