package operator

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/streadway/amqp"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Broker           *rabbitmq.Connection
	AccountService   interfaces.AccountService
	TokenService     interfaces.TokenService
}

type TxQueueOrder struct {
	userAddress common.Address
	baseToken   common.Address
	quoteToken  common.Address
	amount      *big.Int
	pricepoint  *big.Int
	side        *big.Int
	salt        *big.Int
	feeMake     *big.Int
	feeTake     *big.Int
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr interfaces.TradeService,
	p interfaces.EthereumProvider,
	o interfaces.OrderService,
	w *types.Wallet,
	rabbitConn *rabbitmq.Connection,
	accountService interfaces.AccountService,
	tokenService interfaces.TokenService,
) (*TxQueue, error) {
	txq := &TxQueue{
		Name:             n,
		TradeService:     tr,
		OrderService:     o,
		EthereumProvider: p,
		Wallet:           w,
		Broker:           rabbitConn,
		AccountService:   accountService,
		TokenService:     tokenService,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	err = txq.Broker.ConsumeQueuedTrades(ch, &q, txq.ExecuteTrade)
	if err != nil {
		logger.Error(err)
	}

	return txq, nil
}

func (txq *TxQueue) GetChannel() *amqp.Channel {
	name := "TX_QUEUES" + txq.Name
	return txq.Broker.GetChannel(name)
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	return q.Messages
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(m *types.Matches, tag uint64) error {
	logger.Infof("Executing trades: %+v", m)

	makerOrders := m.MakerOrders
	trades := m.Trades
	takerOrder := m.TakerOrder

	orderValues := [][10]*big.Int{}
	orderAddresses := [][4]common.Address{}
	vValues := [][2]uint8{}
	rsValues := [][4][32]byte{}
	amounts := []*big.Int{}

	for i := range makerOrders {
		mo := makerOrders[i]
		to := takerOrder
		t := trades[i]

		orderValues = append(orderValues, [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee})
		orderAddresses = append(orderAddresses, [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, to.QuoteToken})
		vValues = append(vValues, [2]uint8{mo.Signature.V, to.Signature.V})
		rsValues = append(rsValues, [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S})
		amounts = append(amounts, t.Amount)
	}

	for i := range orderAddresses {
		mOrder := TxQueueOrder{
			userAddress: orderAddresses[i][0],
			baseToken:   orderAddresses[i][2],
			quoteToken:  orderAddresses[i][3],
			amount:      orderValues[i][0],
			pricepoint:  orderValues[i][1],
			side:        orderValues[i][2],
			salt:        orderValues[i][3],
			feeMake:     orderValues[i][8],
			feeTake:     orderValues[i][9],
		}

		tOrder := TxQueueOrder{
			userAddress: orderAddresses[i][1],
			baseToken:   orderAddresses[i][2],
			quoteToken:  orderAddresses[i][3],
			amount:      orderValues[i][4],
			pricepoint:  orderValues[i][5],
			side:        orderValues[i][6],
			salt:        orderValues[i][7],
			feeMake:     orderValues[i][8],
			feeTake:     orderValues[i][9],
		}

		baseToken, err := txq.TokenService.GetByAddress(orderAddresses[i][2])

		if err != nil {
			logger.Errorf("Base token address %s not found", orderAddresses[i][2])
			continue
		}

		baseTokenAmount := amounts[i]
		quoteTokenAmount := math.Div(math.Div(math.Mul(amounts[i], mOrder.pricepoint), math.Exp(big.NewInt(10), big.NewInt(int64(baseToken.Decimals)))), big.NewInt(1e18))

		if math.IsEqual(mOrder.side, big.NewInt(0)) {
			err := txq.AccountService.Transfer(mOrder.quoteToken, mOrder.userAddress, tOrder.userAddress, quoteTokenAmount)
			if err != nil {
				logger.Error(err)
			}

			err = txq.AccountService.Transfer(tOrder.baseToken, tOrder.userAddress, mOrder.userAddress, baseTokenAmount)
			if err != nil {
				logger.Error(err)
			}
		} else {
			err := txq.AccountService.Transfer(mOrder.baseToken, mOrder.userAddress, tOrder.userAddress, baseTokenAmount)
			if err != nil {
				logger.Error(err)
			}

			err = txq.AccountService.Transfer(tOrder.quoteToken, tOrder.userAddress, mOrder.userAddress, quoteTokenAmount)
			if err != nil {
				logger.Error(err)
			}
		}
	}

	//err := txq.Broker.PublishTradeSentMessage(m)
	//if err != nil {
	//	logger.Error(err)
	//	return errors.New("Could not update")
	//}

	err := txq.HandleTxSuccess(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	txq.triggerStopOrders(m.Trades)

	return nil
}

func (txq *TxQueue) triggerStopOrders(trades []*types.Trade) {
	for _, trade := range trades {
		stopOrders, err := txq.OrderService.GetTriggeredStopOrders(trade.BaseToken, trade.QuoteToken, trade.PricePoint)

		if err != nil {
			logger.Error(err)
			continue
		}

		for _, stopOrder := range stopOrders {
			err := txq.handleStopOrder(stopOrder)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func (txq *TxQueue) handleStopOrder(so *types.StopOrder) error {
	o, err := so.ToOrder()

	if err != nil {
		logger.Error(err)
		return err
	}

	err = txq.OrderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	so.Status = types.StopOrderStatusDone
	err = txq.OrderService.UpdateStopOrder(so.Hash, so)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) HandleTradeInvalid(m *types.Matches) error {
	logger.Errorf("Trade invalid: %v", m)

	err := txq.Broker.PublishTradeInvalidMessage(m)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) HandleTxError(m *types.Matches) error {
	logger.Errorf("Transaction Error: %v", m)

	errType := "Transaction error"
	err := txq.Broker.PublishTxErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) HandleTxSuccess(m *types.Matches) error {
	logger.Infof("Transaction success: %v", m)

	err := txq.Broker.PublishTradeSuccessMessage(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) HandleError(m *types.Matches) error {
	logger.Errorf("Operator Error: %v", m)

	errType := "Server error"
	err := txq.Broker.PublishErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) PublishPendingTrades(m *types.Matches) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q := txq.Broker.GetQueue(ch, name)

	b, err := json.Marshal(m)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = txq.Broker.Publish(ch, q, b)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PurgePendingTrades() error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	err := txq.Broker.Purge(ch, name)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
