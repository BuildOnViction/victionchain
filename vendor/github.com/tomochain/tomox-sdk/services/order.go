package services

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/utils/math"
	"github.com/tomochain/tomox-sdk/ws"
)

// OrderService
type OrderService struct {
	orderDao        interfaces.OrderDao
	stopOrderDao    interfaces.StopOrderDao
	pairDao         interfaces.PairDao
	accountDao      interfaces.AccountDao
	tradeDao        interfaces.TradeDao
	notificationDao interfaces.NotificationDao
	engine          interfaces.Engine
	validator       interfaces.ValidatorService
	broker          *rabbitmq.Connection
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	stopOrderDao interfaces.StopOrderDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	notificationDao interfaces.NotificationDao,
	engine interfaces.Engine,
	validator interfaces.ValidatorService,
	broker *rabbitmq.Connection,
) *OrderService {

	return &OrderService{
		orderDao,
		stopOrderDao,
		pairDao,
		accountDao,
		tradeDao,
		notificationDao,
		engine,
		validator,
		broker,
	}
}

// GetOrderCountByUserAddress get the total number of orders created by a user
func (s *OrderService) GetOrderCountByUserAddress(addr common.Address) (int, error) {
	return s.orderDao.GetOrderCountByUserAddress(addr)
}

// GetByID fetches the details of an order using order's mongo ID
func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}

// GetByUserAddress fetches all the orders placed by passed user address
func (s *OrderService) GetByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(a, bt, qt, from, to, limit...)
}

// GetByHash fetches all trades corresponding to a trade hash
func (s *OrderService) GetByHash(hash common.Hash) (*types.Order, error) {
	return s.orderDao.GetByHash(hash)
}

func (s *OrderService) GetByHashes(hashes []common.Hash) ([]*types.Order, error) {
	return s.orderDao.GetByHashes(hashes)
}

// // GetByAddress fetches the detailed document of a token using its contract address
// func (s *OrderService) GetTokenByAddress(addr common.Address) (*types.Token, error) {
// 	return s.tokenDao.GetByAddress(addr)
// }

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetCurrentByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetCurrentByUserAddress(addr, limit...)
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetHistoryByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetHistoryByUserAddress(addr, bt, qt, from, to, limit...)
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewOrder(o *types.Order) error {
	if err := o.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := o.VerifySignature()
	if err != nil {
		logger.Error(err)
	}

	if !ok {
		return errors.New("Invalid Signature")
	}

	p, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	if math.IsStrictlySmallerThan(o.QuoteAmount(p), p.MinQuoteAmount()) {
		return errors.New("Order amount too low")
	}

	// Fill token and pair data
	err = o.Process(p)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.validator.ValidateAvailableBalance(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.broker.PublishNewOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewStopOrder(so *types.StopOrder) error {
	if err := so.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := so.VerifySignature()
	if err != nil {
		logger.Error(err)
	}

	if !ok {
		return errors.New("Invalid Signature")
	}

	p, err := s.pairDao.GetByTokenAddress(so.BaseToken, so.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	if math.IsStrictlySmallerThan(so.QuoteAmount(p), p.MinQuoteAmount()) {
		return errors.New("Order amount too low")
	}

	// Fill token and pair data
	err = so.Process(p)
	if err != nil {
		logger.Error(err)
		return err
	}

	//err = s.validator.ValidateAvailableBalance(so)
	//if err != nil {
	//	logger.Error(err)
	//	return err
	//}

	err = s.broker.PublishNewStopOrderMessage(so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(oc *types.OrderCancel) error {
	o, err := s.orderDao.GetByHash(oc.OrderHash)
	if err != nil {
		logger.Error(err)
		return err
	}

	if o == nil {
		return errors.New("No order with corresponding hash")
	}

	if o.Status == types.FILLED || o.Status == types.ERROR_STATUS || o.Status == types.CANCELLED {
		return fmt.Errorf("Cannot cancel order. Status is %v", o.Status)
	}

	err = s.broker.PublishCancelOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelAllOrder(a common.Address) error {
	orders, err := s.orderDao.GetOpenOrdersByUserAddress(a)

	if err != nil {
		logger.Error(err)
		return err
	}

	if len(orders) == 0 {
		return nil
	}

	for _, o := range orders {
		err = s.broker.PublishCancelOrderMessage(o)

		if err != nil {
			logger.Error(err)
			continue
		}
	}

	return nil
}

// CancelStopOrder handles the cancellation stop order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelStopOrder(oc *types.OrderCancel) error {
	o, err := s.stopOrderDao.GetByHash(oc.OrderHash)
	if err != nil {
		logger.Error(err)
		return err
	}

	if o == nil {
		return errors.New("No stop order with corresponding hash")
	}

	if o.Status == types.FILLED || o.Status == types.ERROR_STATUS || o.Status == types.CANCELLED {
		return fmt.Errorf("cannot cancel order. Status is %v", o.Status)
	}

	err = s.broker.PublishCancelStopOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *types.EngineResponse) error {
	switch res.Status {
	case types.ORDER_ADDED:
		s.handleEngineOrderAdded(res)
	case types.ORDER_FILLED:
		s.handleEngineOrderMatched(res)
	case types.ORDER_PARTIALLY_FILLED:
		s.handleEngineOrderMatched(res)
	case types.ORDER_CANCELLED:
		s.handleOrderCancelled(res)
	case types.TRADES_CANCELLED:
		s.handleOrdersInvalidated(res)
	case types.ERROR_STATUS:
		s.handleEngineError(res)
	default:
		s.handleEngineUnknownMessage(res)
	}

	return nil
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	o := res.Order

	// Save notification
	notifications, err := s.notificationDao.Create(&types.Notification{
		Recipient: o.UserAddress,
		Message:   fmt.Sprintf("ORDER_ADDED - Order Hash: %s", o.Hash.Hex()),
		Type:      types.TypeLog,
		Status:    types.StatusUnread,
	})

	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_ADDED", o.UserAddress, o)
	ws.SendNotificationMessage("ORDER_ADDED", o.UserAddress, notifications)

	s.broadcastOrderBookUpdate([]*types.Order{o})
	s.broadcastRawOrderBookUpdate([]*types.Order{o})
}

// handleEngineOrderMatched returns a websocket message informing the client that his order has been added.
// The request signature message also signals the client to sign trades.
func (s *OrderService) handleEngineOrderMatched(res *types.EngineResponse) {
	o := res.Order //res.Order is the "taker" order
	taker := o.UserAddress
	matches := *res.Matches

	orders := []*types.Order{o}
	validMatches := types.Matches{TakerOrder: o}
	invalidMatches := types.Matches{TakerOrder: o}

	//res.Matches is an array of (order, trade) pairs where each order is an "maker" order that is being matched
	for i, _ := range matches.Trades {
		err := s.validator.ValidateBalance(matches.MakerOrders[i])
		if err != nil {
			logger.Error(err)
			invalidMatches.AppendMatch(matches.MakerOrders[i], matches.Trades[i])

		} else {
			validMatches.AppendMatch(matches.MakerOrders[i], matches.Trades[i])
			orders = append(orders, matches.MakerOrders[i])
		}
	}

	// if there are any invalid matches, the maker orders are at cause (since maker orders have been validated in the
	// newOrder() function. We remove the maker orders from the orderbook)
	if invalidMatches.Length() > 0 {
		err := s.broker.PublishInvalidateMakerOrdersMessage(invalidMatches)
		if err != nil {
			logger.Error(err)
		}
	}

	if validMatches.Length() > 0 {
		err := s.tradeDao.Create(validMatches.Trades...)
		if err != nil {
			logger.Error(err)
			ws.SendOrderMessage("ERROR", taker, err)
			return
		}

		err = s.broker.PublishTrades(&validMatches)
		if err != nil {
			logger.Error(err)
			ws.SendOrderMessage("ERROR", taker, err)
			return
		}

		ws.SendOrderMessage("ORDER_MATCHED", taker, types.OrderMatchedPayload{&matches})
	}

	// we only update the orderbook with the current set of orders if there are no invalid matches.
	// If there are invalid matches, the corresponding maker orders will be removed and the taker order
	// amount filled will be updated as a result, and therefore does not represent the current state of the orderbook
	if invalidMatches.Length() == 0 {
		s.broadcastOrderBookUpdate(orders)
		s.broadcastRawOrderBookUpdate(orders)
	}
}

func (s *OrderService) handleOrderCancelled(res *types.EngineResponse) {
	o := res.Order

	// Save notification
	notifications, err := s.notificationDao.Create(&types.Notification{
		Recipient: o.UserAddress,
		Message:   fmt.Sprintf("ORDER_CANCELLED - Order Hash: %s", o.Hash.Hex()),
		Type:      types.TypeLog,
		Status:    types.StatusUnread,
	})

	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_CANCELLED", o.UserAddress, o)
	ws.SendNotificationMessage("ORDER_CANCELLED", o.UserAddress, notifications)

	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	s.broadcastRawOrderBookUpdate([]*types.Order{res.Order})
}

func (s *OrderService) handleOrdersInvalidated(res *types.EngineResponse) error {
	orders := res.InvalidatedOrders
	trades := res.CancelledTrades

	for _, o := range *orders {
		ws.SendOrderMessage("ORDER_INVALIDATED", o.UserAddress, o)
	}

	if orders != nil && len(*orders) != 0 {
		s.broadcastOrderBookUpdate(*orders)
	}

	if orders != nil && len(*orders) != 0 {
		s.broadcastRawOrderBookUpdate(*orders)
	}

	if trades != nil && len(*trades) != 0 {
		s.broadcastTradeUpdate(*trades)
	}

	return nil
}

// handleEngineError returns an websocket error message to the client and recovers orders on the
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ERROR", o.UserAddress, nil)
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *OrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
	log.Print("Receiving unknown engine message")
	utils.PrintJSON(res)
}

func (s *OrderService) HandleOperatorMessages(msg *types.OperatorMessage) error {
	switch msg.MessageType {
	case types.TRADE_ERROR:
		s.handleOperatorTradeError(msg)
	case types.TRADE_TX_PENDING:
		s.handleOperatorTradeTxPending(msg)
	case types.TRADE_TX_SUCCESS:
		s.handleOperatorTradeTxSuccess(msg)
	case types.TRADE_TX_ERROR:
		s.handleOperatorTradeTxError(msg)
	case types.TRADE_INVALID:
		s.handleOperatorTradeInvalid(msg)
	default:
		s.handleOperatorUnknownMessage(msg)
	}

	return nil
}

func (s *OrderService) handleOperatorTradeTxPending(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades
	orders := matches.MakerOrders

	taker := trades[0].Taker
	ws.SendOrderMessage("ORDER_PENDING", taker, types.OrderPendingPayload{matches})

	for _, o := range orders {
		maker := o.UserAddress
		ws.SendOrderMessage("ORDER_PENDING", maker, types.OrderPendingPayload{matches})
	}

	s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeSuccess handles successfull trade messages from the orderbook. It updates
// the trade status in the database and
func (s *OrderService) handleOperatorTradeTxSuccess(msg *types.OperatorMessage) {
	matches := msg.Matches
	hashes := []common.Hash{}
	trades := matches.Trades

	for _, t := range trades {
		hashes = append(hashes, t.Hash)
	}

	if len(hashes) == 0 {
		return
	}

	trades, err := s.tradeDao.UpdateTradeStatuses(types.SUCCESS, hashes...)
	if err != nil {
		logger.Error(err)
	}

	// Send ORDER_SUCCESS message to order takers
	//taker := trades[0].Taker
	//ws.SendOrderMessage("ORDER_SUCCESS", taker, types.OrderSuccessPayload{matches})
	//
	//// Send ORDER_SUCCESS message to order makers
	//for i, _ := range trades {
	//	match := matches.NthMatch(i)
	//	maker := match.MakerOrders[0].UserAddress
	//	ws.SendOrderMessage("ORDER_SUCCESS", maker, types.OrderSuccessPayload{match})
	//}
	//
	//s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeTxError handles cases where a blockchain transaction is reverted
func (s *OrderService) handleOperatorTradeTxError(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades
	//orders := matches.MakerOrders

	errType := msg.ErrorType
	if errType != "" {
		logger.Error("")
	}

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "ERROR"
	}

	//taker := trades[0].Taker
	//ws.SendOrderMessage("ORDER_ERROR", taker, matches)
	//
	//for _, o := range orders {
	//	maker := o.UserAddress
	//	ws.SendOrderMessage("ORDER_ERROR", maker, o)
	//}
	//
	//s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeError handles error messages from the operator (case where the blockchain tx was made
// but ended up failing. It updates the trade status in the db.
// orderbook.
func (s *OrderService) handleOperatorTradeError(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades
	//orders := matches.MakerOrders

	errType := msg.ErrorType
	if errType != "" {
		logger.Error("")
	}

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "ERROR"
	}

	//taker := trades[0].Taker
	//ws.SendOrderMessage("ORDER_ERROR", taker, matches)
	//
	//for _, o := range orders {
	//	maker := o.UserAddress
	//	ws.SendOrderMessage("ORDER_ERROR", maker, o)
	//}
	//
	//s.broadcastTradeUpdate(trades)
}

// handleOperatorTradeInvalid handles the case where one of the two orders is invalid
// which can be the case for example if one of the account addresses does suddendly
// not have enough tokens to satisfy the order. Ultimately, the goal would be to
// reinclude the non-invalid orders in the orderbook
func (s *OrderService) handleOperatorTradeInvalid(msg *types.OperatorMessage) {
	matches := msg.Matches
	trades := matches.Trades
	//orders := matches.MakerOrders

	errType := msg.ErrorType
	if errType != "" {
		logger.Error("")
	}

	for _, t := range trades {
		err := s.tradeDao.UpdateTradeStatus(t.Hash, "ERROR")
		if err != nil {
			logger.Error(err)
		}

		t.Status = "ERROR"
	}

	//taker := trades[0].Taker
	//ws.SendOrderMessage("ORDER_ERROR", taker, matches)
	//
	//for _, o := range orders {
	//	maker := o.UserAddress
	//	ws.SendOrderMessage("ORDER_ERROR", maker, o)
	//}
	//
	//s.broadcastTradeUpdate(trades)
}

func (s *OrderService) handleOperatorUnknownMessage(msg *types.OperatorMessage) {
	log.Print("Receiving unknown message")
	utils.PrintJSON(msg)
}

func (s *OrderService) broadcastOrderBookUpdate(orders []*types.Order) {
	bids := []map[string]string{}
	asks := []map[string]string{}

	p, err := orders[0].Pair()
	if err != nil {
		logger.Error()
		return
	}

	for _, o := range orders {
		pp := o.PricePoint
		side := o.Side

		amount, err := s.orderDao.GetOrderBookPricePoint(p, pp, side)
		if err != nil {
			logger.Error(err)
		}

		// case where the amount at the pricepoint is equal to 0
		if amount == nil {
			amount = big.NewInt(0)
		}

		update := map[string]string{
			"pricepoint": pp.String(),
			"amount":     amount.String(),
		}

		if side == "BUY" {
			bids = append(bids, update)
		} else {
			asks = append(asks, update)
		}
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetOrderBookSocket().BroadcastMessage(id, &types.OrderBook{
		PairName: orders[0].PairName,
		Bids:     bids,
		Asks:     asks,
	})
}

func (s *OrderService) broadcastRawOrderBookUpdate(orders []*types.Order) {
	p, err := orders[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetRawOrderBookSocket().BroadcastMessage(id, orders)
}

func (s *OrderService) broadcastTradeUpdate(trades []*types.Trade) {
	p, err := trades[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetTradeChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetTradeSocket().BroadcastMessage(id, trades)
}

func (s *OrderService) WatchChanges() {
	pipeline := []bson.M{}

	ct, err := s.orderDao.GetCollection().Watch(pipeline, mgo.ChangeStreamOptions{FullDocument: mgo.UpdateLookup})

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()

	// Watch the event again in case there is error and function returned
	defer s.WatchChanges()

	ctx := context.Background()

	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done(): // if parent context was cancelled
			err := ct.Close() // close the stream
			if err != nil {
				logger.Error("Change stream closed")
			}
			return //exiting from the func
		default:
			ev := types.OrderChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					//if err is not nil, it means something bad happened, let's finish our func
					logger.Error(err)
					return
				}
			}

			//if item from the stream un-marshaled successfully, do something with it
			if ok {
				logger.Debugf("Operation Type: %s", ev.OperationType)
				s.HandleDocumentType(ev)
			}
		}
	}
}

func (s *OrderService) HandleDocumentType(ev types.OrderChangeEvent) error {
	res := &types.EngineResponse{}

	switch ev.OperationType {
	case types.OPERATION_TYPE_INSERT:
		if ev.FullDocument.Status == "OPEN" {
			res.Status = types.ORDER_ADDED
			res.Order = ev.FullDocument
		}
		break
	case types.OPERATION_TYPE_UPDATE:
		if ev.FullDocument.Status == "CANCELLED" {
			res.Status = types.ORDER_CANCELLED
			res.Order = ev.FullDocument
		}
		break
	case types.OPERATION_TYPE_REPLACE:
		if ev.FullDocument.Status == "CANCELLED" {
			res.Status = types.ORDER_CANCELLED
			res.Order = ev.FullDocument
		}
		break
	default:
		break
	}

	if res.Status != "" {
		err := s.broker.PublishEngineResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (s *OrderService) GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error) {
	return s.stopOrderDao.GetTriggeredStopOrders(baseToken, quoteToken, lastPrice)
}

func (s *OrderService) UpdateStopOrder(h common.Hash, so *types.StopOrder) error {
	return s.stopOrderDao.UpdateByHash(h, so)
}
