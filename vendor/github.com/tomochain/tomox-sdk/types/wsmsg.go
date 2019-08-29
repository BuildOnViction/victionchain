package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "SUBSCRIBE"
	UNSUBSCRIBE SubscriptionEvent = "UNSUBSCRIBE"
	Fetch       SubscriptionEvent = "fetch"

	UPDATE        SubscriptionEvent = "UPDATE"
	ERROR         SubscriptionEvent = "ERROR"
	SUCCESS_EVENT SubscriptionEvent = "SUCCESS"
	INIT          SubscriptionEvent = "INIT"
	CANCEL        SubscriptionEvent = "CANCEL"

	// status

	ORDER_ADDED            = "ORDER_ADDED"
	ORDER_FILLED           = "ORDER_FILLED"
	ORDER_PARTIALLY_FILLED = "ORDER_PARTIALLY_FILLED"
	ORDER_CANCELLED        = "ORDER_CANCELLED"
	ERROR_STATUS           = "ERROR"

	TradeAdded   = "TRADE_ADDED"
	TradeUpdated = "TRADE_UPDATED"
	// channel
	TradeChannel     = "trades"
	OrderbookChannel = "orderbook"
	OrderChannel     = "orders"
	OHLCVChannel     = "ohlcv"
)

type WebsocketMessage struct {
	Channel string         `json:"channel"`
	Event   WebsocketEvent `json:"event"`
}

func (ev *WebsocketMessage) String() string {
	return fmt.Sprintf("%v/%v", ev.Channel, ev.Event.String())
}

type WebsocketEvent struct {
	Type    SubscriptionEvent `json:"type"`
	Hash    string            `json:"hash,omitempty"`
	Payload interface{}       `json:"payload"`
}

func (ev *WebsocketEvent) String() string {
	return fmt.Sprintf("%v", ev.Type)
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	PairID   string `json:"pair"`
}

type OrderPendingPayload struct {
	Matches *Matches `json:"matches"`
}

type OrderSuccessPayload struct {
	Matches *Matches `json:"matches"`
}

type OrderMatchedPayload struct {
	Matches *Matches `json:"matches"`
}

type SubscriptionPayload struct {
	PairName   string         `json:"pairName,omitempty"`
	QuoteToken common.Address `json:"quoteToken,omitempty"`
	BaseToken  common.Address `json:"baseToken,omitempty"`
	From       int64          `json"from"`
	To         int64          `json:"to"`
	Duration   int64          `json:"duration"`
	Units      string         `json:"units"`
}

func NewOrderWebsocketMessage(o *Order) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "NEW_ORDER",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderAddedWebsocketMessage(o *Order, p *Pair, filled int64) *WebsocketMessage {
	o.Process(p)
	o.FilledAmount = big.NewInt(filled)
	o.Status = "OPEN"
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "ORDER_ADDED",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderCancelWebsocketMessage(oc *OrderCancel) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "CANCEL_ORDER",
			Hash:    oc.Hash.Hex(),
			Payload: oc,
		},
	}
}
