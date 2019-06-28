package ws

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/types"
)

// NotificationConnection is websocket notification connection struct
// It holds the reference to connection and the channel of type NotificationChannel

type NotificationConnection []*Client

var notificationConnections map[string]NotificationConnection

// GetNotificationConnections returns the connection associated with an user address
func GetNotificationConnections(a common.Address) NotificationConnection {
	c := notificationConnections[a.Hex()]
	if c == nil {
		logger.Warning("No connection found")
		return nil
	}

	return notificationConnections[a.Hex()]
}

func NotificationSocketUnsubscribeHandler(a common.Address) func(client *Client) {
	return func(client *Client) {
		logger.Info("In unsubscription handler")
		notificationConnection := notificationConnections[a.Hex()]
		if notificationConnection == nil {
			logger.Info("No subscriptions")
		}

		if notificationConnection != nil {
			logger.Info("%v connections before unsubscription", len(notificationConnections[a.Hex()]))
			for i, c := range notificationConnection {
				if client == c {
					notificationConnection = append(notificationConnection[:i], notificationConnection[i+1:]...)
				}
			}

		}

		notificationConnections[a.Hex()] = notificationConnection
		logger.Info("%v connections after unsubscription", len(notificationConnections[a.Hex()]))
	}
}

// RegisterNotificationConnection registers a connection with an user address
// It is called whenever a message is received over notification channel
func RegisterNotificationConnection(a common.Address, c *Client) {
	logger.Info("Registering new notification connection")

	if notificationConnections == nil {
		notificationConnections = make(map[string]NotificationConnection)
	}

	if notificationConnections[a.Hex()] == nil {
		logger.Info("Registering a new notification connection")
		notificationConnections[a.Hex()] = NotificationConnection{c}
		RegisterConnectionUnsubscribeHandler(c, NotificationSocketUnsubscribeHandler(a))
		logger.Info("Number of connections for this address: %v", len(notificationConnections))
	}

	if notificationConnections[a.Hex()] != nil {

		if !isClientConnected(notificationConnections[a.Hex()], c) {
			logger.Info("Registering a new notification connection")
			notificationConnections[a.Hex()] = append(notificationConnections[a.Hex()], c)
			RegisterConnectionUnsubscribeHandler(c, NotificationSocketUnsubscribeHandler(a))
			logger.Info("Number of connections for this address: %v", len(notificationConnections))
		}
	}
}

func SendNotificationMessage(msgType types.SubscriptionEvent, a common.Address, payload interface{}) {
	conn := GetNotificationConnections(a)
	if conn == nil {
		return
	}

	for _, c := range conn {
		c.SendMessage(NotificationChannel, msgType, payload)
	}
}

// SendNotificationErrorMessage sends error message on markets channel
func SendNotificationErrorMessage(c *Client, data interface{}) {
	c.SendMessage(NotificationChannel, types.ERROR, data)
}
