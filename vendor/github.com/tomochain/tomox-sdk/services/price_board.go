package services

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type PriceBoardService struct {
	TokenDao interfaces.TokenDao
	TradeDao interfaces.TradeDao
}

// NewTradeService returns a new instance of TradeService
func NewPriceBoardService(
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
) *PriceBoardService {
	return &PriceBoardService{
		TokenDao: tokenDao,
		TradeDao: tradeDao,
	}
}

// Subscribe
func (s *PriceBoardService) Subscribe(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	// Fix the value at 1 day because we only care about 24h change
	duration := int64(1)
	unit := "day"

	ticks, err := s.GetPriceBoardData(
		[]types.PairAddresses{{BaseToken: bt, QuoteToken: qt}},
		duration,
		unit,
	)

	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	quoteToken, err := s.TokenDao.GetByAddress(qt)

	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	var lastTradePrice string
	lastTrade, err := s.TradeDao.GetLatestTrade(bt, qt)
	if lastTrade == nil {
		lastTradePrice = "?"
	} else {
		lastTradePrice = lastTrade.PricePoint.String()
	}

	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetPriceBoardChannelID(bt, qt)
	err = socket.Subscribe(id, c)

	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))

	result := types.PriceBoardData{
		Ticks:          ticks,
		PriceUSD:       quoteToken.USD,
		LastTradePrice: lastTradePrice,
	}

	socket.SendInitMessage(c, result)
}

// Unsubscribe
func (s *PriceBoardService) UnsubscribeChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	id := utils.GetPriceBoardChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe
func (s *PriceBoardService) Unsubscribe(c *ws.Client) {
	socket := ws.GetPriceBoardSocket()
	socket.Unsubscribe(c)
}

func (s *PriceBoardService) GetPriceBoardData(pairs []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)

	currentTimestamp := time.Now().Unix()

	_, intervalInSeconds := getModTime(currentTimestamp, duration, unit)

	start := time.Unix(currentTimestamp-intervalInSeconds, 0)
	end := time.Unix(currentTimestamp, 0)

	if len(timeInterval) >= 1 {
		end = time.Unix(timeInterval[1], 0)
		start = time.Unix(timeInterval[0], 0)
	}

	match := make(bson.M)
	match = getMatchQuery(start, end, pairs...)
	match = bson.M{"$match": match}

	group := getGroupBson()
	group = bson.M{"$group": group}

	query := []bson.M{match, group}

	res, err := s.TradeDao.Aggregate(query)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []*types.Tick{}, nil
	}

	return res, nil
}

// query for grouping of the documents into one
func getGroupBson() bson.M {
	var group bson.M

	one, _ := bson.ParseDecimal128("1")
	group = bson.M{
		"count":  bson.M{"$sum": one},
		"high":   bson.M{"$max": "$pricepoint"},
		"low":    bson.M{"$min": "$pricepoint"},
		"open":   bson.M{"$first": "$pricepoint"},
		"close":  bson.M{"$last": "$pricepoint"},
		"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
	}
	groupID := make(bson.M)
	groupID["pairName"] = "$pairName"
	groupID["baseToken"] = "$baseToken"
	groupID["quoteToken"] = "$quoteToken"
	group["_id"] = groupID

	return group
}
