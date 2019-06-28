package crons

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startPriceBoardCron(c *cron.Cron) {
	pairs, err := s.PairService.GetAll()

	if err != nil {
		log.Println(err.Error())
	}

	for _, p := range pairs {
		c.AddFunc("*/3 * * * * *", s.getPriceBoardData(p.BaseTokenAddress, p.QuoteTokenAddress))
	}
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getPriceBoardData(bt, qt common.Address) func() {
	return func() {
		p := make([]types.PairAddresses, 0)
		p = []types.PairAddresses{{
			BaseToken:  bt,
			QuoteToken: qt,
		}}

		// Fix the value at 1 day because we only care about 24h change
		duration := int64(1)
		unit := "day"

		ticks, err := s.PriceBoardService.GetPriceBoardData(p, duration, unit)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		quoteToken, err := s.PriceBoardService.TokenDao.GetByAddress(qt)

		if err != nil {
			log.Printf("%s", err)
			return
		}

		var lastTradePrice string
		lastTrade, err := s.PriceBoardService.TradeDao.GetLatestTrade(bt, qt)
		if lastTrade == nil {
			lastTradePrice = "?"
		} else {
			lastTradePrice = lastTrade.PricePoint.String()
		}

		if err != nil {
			log.Printf("%s", err)
			return
		}

		id := utils.GetPriceBoardChannelID(bt, qt)

		result := types.PriceBoardData{
			Ticks:          ticks,
			PriceUSD:       quoteToken.USD,
			LastTradePrice: lastTradePrice,
		}

		ws.GetPriceBoardSocket().BroadcastMessage(id, result)
	}
}
