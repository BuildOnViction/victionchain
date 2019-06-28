package engine

import (
	"encoding/json"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/ethereum"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
)

// Engine contains daos required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	rabbitMQConn *rabbitmq.Connection
	provider     *ethereum.EthereumProvider
}

var logger = utils.Logger

// NewEngine initializes the engine singleton instance
func NewEngine(
	rabbitMQConn *rabbitmq.Connection,
	orderDao interfaces.OrderDao,
	stopOrderDao interfaces.StopOrderDao,
	tradeDao interfaces.TradeDao,
	pairDao interfaces.PairDao,
	provider *ethereum.EthereumProvider,
) *Engine {
	pairs, err := pairDao.GetAll()

	if err != nil {
		panic(err)
	}

	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := NewOrderBook(rabbitMQConn, orderDao, stopOrderDao, tradeDao, p)

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, rabbitMQConn, provider}
	return engine
}

// Provider : implement engine interface
func (e *Engine) Provider() interfaces.EthereumProvider {
	return e.provider
}

// HandleOrders parses incoming rabbitmq order messages and redirects them to the appropriate
// engine function
func (e *Engine) HandleOrders(msg *rabbitmq.Message) error {
	switch msg.Type {
	case "NEW_ORDER":
		err := e.handleNewOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "NEW_STOP_ORDER":
		err := e.handleNewStopOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "CANCEL_ORDER":
		err := e.handleCancelOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "CANCEL_STOP_ORDER":
		err := e.handleCancelStopOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_MAKER_ORDERS":
		err := e.handleInvalidateMakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_TAKER_ORDERS":
		err := e.handleInvalidateTakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		logger.Error("Unknown message", msg)
	}

	return nil
}

func (e *Engine) handleNewOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.newOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleNewStopOrder(bytes []byte) error {
	so := &types.StopOrder{}
	err := json.Unmarshal(bytes, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := so.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.newStopOrder(so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleCancelOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.cancelOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleCancelStopOrder(bytes []byte) error {
	so := &types.StopOrder{}
	err := json.Unmarshal(bytes, so)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := so.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.cancelStopOrder(so)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateMakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.invalidateMakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateTakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		logger.Error(err)
		return err
	}

	err = ob.invalidateTakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
