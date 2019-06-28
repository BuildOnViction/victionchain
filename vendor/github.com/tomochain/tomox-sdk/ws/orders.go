package ws

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/types"
)

// OrderConn is websocket order connection struct
// It holds the reference to connection and the channel of type OrderMessage

type OrderConnection []*Client

var orderConnections map[string]OrderConnection

// GetOrderConn returns the connection associated with an order ID
func GetOrderConnections(a common.Address) OrderConnection {
	c := orderConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return orderConnections[a.Hex()]
}

func OrderSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		logger.Info("In unsubscription handler")
		orderConnection := orderConnections[a.Hex()]
		if orderConnection == nil {
			logger.Info("No subscriptions")
		}

		if orderConnection != nil {
			logger.Info("%v connections before unsubscription", len(orderConnections[a.Hex()]))
			for i, c := range orderConnection {
				if client == c {
					orderConnection = append(orderConnection[:i], orderConnection[i+1:]...)
				}
			}

		}

		orderConnections[a.Hex()] = orderConnection
		logger.Info("%v connections after unsubscription", len(orderConnections[a.Hex()]))
	}
}

// RegisterOrderConnection registers a connection with and orderID.
// It is called whenever a message is recieved over order channel
func RegisterOrderConnection(a common.Address, c *Client) {
	logger.Info("Registering new order connection")

	if orderConnections == nil {
		orderConnections = make(map[string]OrderConnection)
	}

	if orderConnections[a.Hex()] == nil {
		logger.Info("Registering a new order connection")
		orderConnections[a.Hex()] = OrderConnection{c}
		RegisterConnectionUnsubscribeHandler(c, OrderSocketUnsubscribeHandler(a))
		logger.Info("Number of connections for this address: %v", len(orderConnections))
	}

	if orderConnections[a.Hex()] != nil {

		if !isClientConnected(orderConnections[a.Hex()], c) {
			logger.Info("Registering a new order connection")
			orderConnections[a.Hex()] = append(orderConnections[a.Hex()], c)
			RegisterConnectionUnsubscribeHandler(c, OrderSocketUnsubscribeHandler(a))
			logger.Info("Number of connections for this address: %v", len(orderConnections))
		}
	}
}

func SendOrderMessage(msgType types.SubscriptionEvent, a common.Address, payload interface{}) {
	conn := GetOrderConnections(a)
	if conn == nil {
		return
	}

	for _, c := range conn {
		c.SendMessage(OrderChannel, msgType, payload)
	}
}
