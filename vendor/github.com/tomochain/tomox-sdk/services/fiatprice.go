package services

import (
	"fmt"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type FiatPriceService struct {
	TokenDao     interfaces.TokenDao
	FiatPriceDao interfaces.FiatPriceDao
}

// NewTradeService returns a new instance of TradeService
func NewFiatPriceService(
	tokenDao interfaces.TokenDao,
	fiatPriceDao interfaces.FiatPriceDao,
) *FiatPriceService {
	return &FiatPriceService{
		TokenDao:     tokenDao,
		FiatPriceDao: fiatPriceDao,
	}
}

// InitFiatPrice will query Coingecko API and stores fiat price data in the last 1 day after booting up server
func (s *FiatPriceService) InitFiatPrice() {
	// Fix ids with 4 coins
	symbols := []string{"bitcoin", "ethereum", "ripple", "tomochain"}
	// Fix fiat currency with USD
	vsCurrency := "usd"

	for _, symbol := range symbols {
		data, err := s.FiatPriceDao.GetCoinMarketChart(symbol, vsCurrency, "2")

		if err != nil {
			logger.Error(err)
			continue
		}

		items := data.Prices

		for _, item := range items {
			fiatPriceItem := &types.FiatPriceItem{
				Symbol:       symbol,
				Timestamp:    fmt.Sprintf("%d", int64(item[0])), // Convert timestamp from float64 to int64
				Price:        fmt.Sprintf("%f", item[1]),
				FiatCurrency: vsCurrency,
			}

			_, err := s.FiatPriceDao.FindAndModify(
				fiatPriceItem.Symbol,
				fiatPriceItem.FiatCurrency,
				fiatPriceItem.Timestamp,
				fiatPriceItem,
			)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

// UpdateFiatPrice will query Coingecko API and stores fiat price data in the last 30 minutes
func (s *FiatPriceService) UpdateFiatPrice() {
	// Fix ids with 4 coins
	symbols := []string{"bitcoin", "ethereum", "ripple", "tomochain"}
	// Fix fiat currency with USD
	vsCurrency := "usd"

	for _, symbol := range symbols {
		data, err := s.FiatPriceDao.GetCoinMarketChart(symbol, vsCurrency, "2")

		if err != nil {
			logger.Error(err)
			continue
		}

		items := data.Prices

		for _, item := range items {
			fiatPriceItem := &types.FiatPriceItem{
				Symbol:       symbol,
				Timestamp:    fmt.Sprintf("%d", int64(item[0])), // Convert timestamp from float64 to int64
				Price:        fmt.Sprintf("%f", item[1]),
				FiatCurrency: vsCurrency,
			}

			_, err := s.FiatPriceDao.FindAndModify(
				fiatPriceItem.Symbol,
				fiatPriceItem.FiatCurrency,
				fiatPriceItem.Timestamp,
				fiatPriceItem,
			)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func (s *FiatPriceService) SyncFiatPrice() error {
	prices, err := s.FiatPriceDao.GetLatestQuotes()

	if err != nil {
		logger.Error(err)
		return err
	}

	for k, v := range prices {
		err := s.TokenDao.UpdateFiatPriceBySymbol(k, v)

		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *FiatPriceService) GetFiatPriceChart() (map[string][]*types.FiatPriceItem, error) {
	result := make(map[string][]*types.FiatPriceItem)

	// Fix ids with 4 coins
	symbols := []string{"bitcoin", "ethereum", "ripple", "tomochain"}

	for _, symbol := range symbols {
		data, err := s.FiatPriceDao.Get24hChart(symbol, "usd")

		if err != nil {
			logger.Error(err)
			continue
		}

		result[symbol] = data
	}

	return result, nil
}
