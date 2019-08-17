package tomox

import "testing"

// SCENARIO:
// pair BTC/USDT
// userA: sell 1 BTC, price 9000 USDT
// userB: buy 1 BTC, price 9000 USDT
// FEE: 0.1% = 1/1000, fee is calculated in quoteToken (USDT)
// EXPECTED:
// userA: received 9000 - 9 = 8991 USDT, send out 1 BTC
// userB received 1 BTC, send out 9000 + 9 = 9009 USDT

// A is taker
func TestSettleBalance_TakerSell(t *testing.T) {

}

// A is maker
func TestSettleBalance_TakerBuy(t *testing.T) {

}

