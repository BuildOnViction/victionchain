package crons

import (
	"log"

	"github.com/robfig/cron"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startMarketsCron(c *cron.Cron) {
	c.AddFunc("*/3 * * * * *", s.getMarketsData())
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getMarketsData() func() {
	return func() {
		pairData, err := s.PairService.GetAllTokenPairData()

		if err != nil {
			log.Printf("%s", err)
			return
		}

		smallChartsDataResult, err := s.FiatPriceService.GetFiatPriceChart()

		res := &types.MarketData{
			PairData:        pairData,
			SmallChartsData: smallChartsDataResult,
		}

		id := utils.GetMarketsChannelID(ws.MarketsChannel)

		ws.GetMarketSocket().BroadcastMessage(id, res)
	}
}
