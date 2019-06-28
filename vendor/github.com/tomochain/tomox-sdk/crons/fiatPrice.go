package crons

import (
	"github.com/robfig/cron"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) getFiatPriceCron(c *cron.Cron) {
	c.AddFunc("0 */30 * * * *", s.updateFiatPrice())
	c.AddFunc("0 0 * * * *", s.syncFiatPrice())
}

// updateFiatPrice function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) updateFiatPrice() func() {
	return func() {
		s.FiatPriceService.UpdateFiatPrice()
	}
}

// syncFiatPrice function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) syncFiatPrice() func() {
	return func() {
		s.FiatPriceService.SyncFiatPrice()
	}
}
