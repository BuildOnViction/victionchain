package ws

import (
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var marketsSocket *MarketsSocket

// MarketsSocket holds the map of subscriptions subscribed to markets channels
// corresponding to the key/event they have subscribed to.
type MarketsSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

func NewMarketsSocket() *MarketsSocket {
	return &MarketsSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetMarketSocket return singleton instance of MarketsSocket type struct
func GetMarketSocket() *MarketsSocket {
	if marketsSocket == nil {
		marketsSocket = NewMarketsSocket()
	}

	return marketsSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
func (s *MarketsSocket) Subscribe(channelID string, c *Client) error {
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

// UnsubscribeHandler unsubscribes a connection from a certain markets channel id
func (s *MarketsSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *MarketsSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// Unsubscribe removes a websocket connection from the markets channel updates
func (s *MarketsSocket) UnsubscribeChannel(channelID string, c *Client) {
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *MarketsSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscriptions subscribed to the pair
func (s *MarketsSocket) BroadcastMessage(channelID string, p interface{}) error {

	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the markets channel
func (s *MarketsSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(MarketsChannel, msgType, p)
}

// SendInitMessage sends INIT message on markets channel on subscription event
func (s *MarketsSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(MarketsChannel, types.INIT, data)
}

// SendUpdateMessage sends UPDATE message on markets channel as new data is created
func (s *MarketsSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(MarketsChannel, types.UPDATE, data)
}

// SendErrorMessage sends error message on markets channel
func (s *MarketsSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(MarketsChannel, types.ERROR, data)
}
