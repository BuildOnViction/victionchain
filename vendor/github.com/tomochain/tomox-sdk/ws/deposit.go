package ws

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/types"
)

// DepositConn is websocket deposit connection struct
// It holds the reference to connection and the channel of type DepositMessage

// Send update directly to client based on wallet address
type DepositConnection []*Client

var depositConnections map[string]DepositConnection

// GetDepositConn returns the connection associated with an deposit ID
func GetDepositConnections(a common.Address) DepositConnection {
	c := depositConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return depositConnections[a.Hex()]
}

func DepositSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		logger.Info("In unsubscription handler")
		depositConnection := depositConnections[a.Hex()]
		if depositConnection == nil {
			logger.Info("No subscriptions")
		}

		if depositConnection != nil {
			logger.Info("%v connections before unsubscription", len(depositConnections[a.Hex()]))
			for i, c := range depositConnection {
				if client == c {
					depositConnection = append(depositConnection[:i], depositConnection[i+1:]...)
				}
			}

		}

		depositConnections[a.Hex()] = depositConnection
		logger.Info("%v connections after unsubscription", len(depositConnections[a.Hex()]))
	}
}

// RegisterDepositConnection registers a connection with and depositID.
// It is called whenever a message is recieved over deposit channel
func RegisterDepositConnection(a common.Address, c *Client) {
	logger.Info("Registering new deposit connection")

	if depositConnections == nil {
		depositConnections = make(map[string]DepositConnection)
	}

	if depositConnections[a.Hex()] == nil {
		logger.Info("Registering a new deposit connection")
		depositConnections[a.Hex()] = DepositConnection{c}
		RegisterConnectionUnsubscribeHandler(c, DepositSocketUnsubscribeHandler(a))
		logger.Info("Number of connections for this address: %v", len(depositConnections))
	}

	if depositConnections[a.Hex()] != nil {
		if !isClientConnected(depositConnections[a.Hex()], c) {
			logger.Info("Registering a new deposit connection")
			depositConnections[a.Hex()] = append(depositConnections[a.Hex()], c)
			RegisterConnectionUnsubscribeHandler(c, DepositSocketUnsubscribeHandler(a))
			logger.Info("Number of connections for this address: %v", len(depositConnections))
		}
	}
}

func SendDepositMessage(msgType types.SubscriptionEvent, a common.Address, payload interface{}) {
	conn := GetDepositConnections(a)
	if conn == nil {
		logger.Infof("No connection found")
		return
	}

	for _, c := range conn {
		c.SendMessage(DepositChannel, msgType, payload)
	}
}
