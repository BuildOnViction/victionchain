package ws

import (
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var orderbookSocket *OrderBookSocket

// OrderBookSocket holds the map of subscriptions subscribed to orderbook channels
// corresponding to the key/event they have subscribed to.
type OrderBookSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewOrderBookSocket() *OrderBookSocket {
	return &OrderBookSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetOrderBookSocket return singleton instance of OrderBookSocket type struct
func GetOrderBookSocket() *OrderBookSocket {
	if orderbookSocket == nil {
		orderbookSocket = NewOrderBookSocket()
	}

	return orderbookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *OrderBookSocket) Subscribe(channelID string, c *Client) error {
	if c == nil {
		return errors.New("No connection found")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Client]bool)
	}

	s.subscriptions[channelID][c] = true

	if s.subscriptionsList[c] == nil {
		s.subscriptionsList[c] = []string{}
	}

	s.subscriptionsList[c] = append(s.subscriptionsList[c], channelID)

	return nil
}

// UnsubscribeHandler unsubscribes a connection from a certain orderbook channel id
func (s *OrderBookSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *OrderBookSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe removes a websocket connection from the orderbook channel updates
func (s *OrderBookSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *OrderBookSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *OrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {

	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the orderbook channel
func (s *OrderBookSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(OrderBookChannel, msgType, p)
}

// SendInitMessage sends INIT message on orderbook channel on subscription event
func (s *OrderBookSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, types.INIT, data)
}

// SendUpdateMessage sends UPDATE message on orderbook channel as new data is created
func (s *OrderBookSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, types.UPDATE, data)
}

// SendErrorMessage sends error message on orderbook channel
func (s *OrderBookSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(OrderBookChannel, types.ERROR, data)
}
