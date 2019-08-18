package tomox

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

// SCENARIO:
// pair BTC/USDT
// userA: sell 1 BTC, price 9000 USDT
// userB: buy 1 BTC, price 9000 USDT
// FEE: 0.1% = 1/1000, fee is calculated in quoteToken (USDT)
// EXPECTED:
// userA: received 9000 - 9 = 8991 USDT, send out 1 BTC
// userB received 1 BTC, send out 9000 + 9 = 9009 USDT

var (
	// parameters
	userA      = common.HexToAddress("0x000000000000000000000000000000000000000a")
	userB      = common.HexToAddress("0x000000000000000000000000000000000000000b")
	baseToken  = common.HexToAddress("0x000000000000000000000000000000000000000x")
	quoteToken = common.HexToAddress("0x000000000000000000000000000000000000000y")
	feeRate    = big.NewInt(0).Mul(big.NewInt(1), common.BasePrice)
	quantity   = big.NewInt(0).Mul(big.NewInt(1), common.BasePrice)
	price      = big.NewInt(0).Mul(big.NewInt(9000), common.BasePrice) // 9000 * 10^18

	// expected
	userAReceived = big.NewInt(0).Mul(big.NewInt(8991), common.BasePrice)
	userASend     = big.NewInt(0).Mul(big.NewInt(1), common.BasePrice)
	userBReceived = big.NewInt(0).Mul(big.NewInt(1), common.BasePrice)
	userBSend     = big.NewInt(0).Mul(big.NewInt(9009), common.BasePrice)
	expectedFee   = big.NewInt(0).Mul(big.NewInt(9), common.BasePrice)
)

// A is taker
func TestSettleBalance_TakerSell(t *testing.T) {
	result := SettleBalance(
		userB,
		userA,
		baseToken,
		quoteToken,
		false,
		feeRate,
		feeRate,
		common.TomoXBaseFee,
		quantity,
		price)

	// taker
	takerInToken := result[userA][InToken].(common.Address)
	takerInTotal := result[userA][InTotal].(*big.Int)
	takerOutToken := result[userA][OutToken].(common.Address)
	takerOutTotal := result[userA][OutTotal].(*big.Int)

	// verify token type
	if !bytes.Equal(takerInToken.Bytes(), quoteToken.Bytes()) || !bytes.Equal(takerOutToken.Bytes(), baseToken.Bytes()) {
		t.Error("Wrong token type of taker",
			"Expected inToken: ", quoteToken, "Actual inToken: ", takerInToken,
			"Expected outToken: ", baseToken, "Actual outToken: ", takerOutToken)
	}

	// verify quantity
	if takerInTotal.Cmp(userAReceived) != 0 {
		t.Error("Taker received wrong quantity", "Expected received:", userAReceived, "Actual received:", takerInTotal)
	}
	if takerOutTotal.Cmp(userASend) != 0 {
		t.Error("Taker sends wrong quantity", "Expected send:", userASend, "Actual send:", takerOutTotal)
	}

	// maker
	makerInToken := result[userB][InToken].(common.Address)
	makerInTotal := result[userB][InTotal].(*big.Int)
	makerOutToken := result[userB][OutToken].(common.Address)
	makerOutTotal := result[userB][OutTotal].(*big.Int)

	// verify token type
	if !bytes.Equal(makerInToken.Bytes(), baseToken.Bytes()) || !bytes.Equal(makerOutToken.Bytes(), quoteToken.Bytes()) {
		t.Error("Wrong token type of maker",
			"Expected inToken: ", baseToken, "Actual inToken: ", makerInToken,
			"Expected outToken: ", quoteToken, "Actual outToken: ", makerOutToken)
	}

	// verify quantity
	if makerInTotal.Cmp(userBReceived) != 0 {
		t.Error("Maker received wrong quantity", "Expected received:", userBReceived, "Actual received:", makerInTotal)
	}
	if makerOutTotal.Cmp(userBSend) != 0 {
		t.Error("Maker sends wrong quantity", "Expected send:", userBSend, "Actual send:", makerOutTotal)
	}

	// fee
	// taker fee
	takerFee := result[userA][Fee].(*big.Int)
	if takerFee.Cmp(expectedFee) != 0 {
		t.Error("Wrong taker fee amount", "Expected: ", 9, "Actual: ", takerFee)
	}
	// maker fee
	makerFee := result[userB][Fee].(*big.Int)
	if takerFee.Cmp(expectedFee) != 0 {
		t.Error("Wrong makerFee fee amount", "Expected: ", 9, "Actual: ", makerFee)
	}
}

// A is maker
func TestSettleBalance_TakerBuy(t *testing.T) {
	result := SettleBalance(
		userB,
		userA,
		baseToken,
		quoteToken,
		false,
		feeRate,
		feeRate,
		common.TomoXBaseFee,
		quantity,
		price)

	// taker
	takerInToken := result[userB][InToken].(common.Address)
	takerInTotal := result[userB][InTotal].(*big.Int)
	takerOutToken := result[userB][OutToken].(common.Address)
	takerOutTotal := result[userB][OutTotal].(*big.Int)

	// verify token type
	if !bytes.Equal(takerInToken.Bytes(), baseToken.Bytes()) || !bytes.Equal(takerOutToken.Bytes(), quoteToken.Bytes()) {
		t.Error("Wrong token type of taker",
			"Expected inToken: ", baseToken, "Actual inToken: ", takerInToken,
			"Expected outToken: ", quoteToken, "Actual outToken: ", takerOutToken)
	}

	// verify quantity
	if takerInTotal.Cmp(userBReceived) != 0 {
		t.Error("Taker received wrong quantity", "Expected received:", userBReceived, "Actual received:", takerInTotal)
	}
	if takerOutTotal.Cmp(userBSend) != 0 {
		t.Error("Taker sends wrong quantity", "Expected send:", userBSend, "Actual send:", takerOutTotal)
	}

	// maker
	makerInToken := result[userA][InToken].(common.Address)
	makerInTotal := result[userA][InTotal].(*big.Int)
	makerOutToken := result[userA][OutToken].(common.Address)
	makerOutTotal := result[userA][OutTotal].(*big.Int)

	// verify token type
	if !bytes.Equal(makerInToken.Bytes(), quoteToken.Bytes()) || !bytes.Equal(makerOutToken.Bytes(), baseToken.Bytes()) {
		t.Error("Wrong token type of maker",
			"Expected inToken: ", quoteToken, "Actual inToken: ", makerInToken,
			"Expected outToken: ", baseToken, "Actual outToken: ", makerOutToken)
	}

	// verify quantity
	if makerInTotal.Cmp(userAReceived) != 0 {
		t.Error("Maker received wrong quantity", "Expected received:", userAReceived, "Actual received:", makerInTotal)
	}
	if makerOutTotal.Cmp(userASend) != 0 {
		t.Error("Maker sends wrong quantity", "Expected send:", userASend, "Actual send:", makerOutTotal)
	}

	// fee
	// taker fee
	takerFee := result[userB][Fee].(*big.Int)
	if takerFee.Cmp(expectedFee) != 0 {
		t.Error("Wrong taker fee amount", "Expected: ", 9, "Actual: ", takerFee)
	}
	// maker fee
	makerFee := result[userA][Fee].(*big.Int)
	if takerFee.Cmp(expectedFee) != 0 {
		t.Error("Wrong makerFee fee amount", "Expected: ", 9, "Actual: ", makerFee)
	}
}
