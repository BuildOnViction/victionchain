package ws

import (
	"github.com/op/go-logging"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
)

type WebsocketLogger struct {
	*logging.Logger
	websocketMessageLogger *logging.Logger
}

func NewWebsocketLogger() *WebsocketLogger {
	return &WebsocketLogger{
		utils.StdoutLogger,
		utils.Logger,
	}
}

func (l *WebsocketLogger) LogMessageIn(msg *types.WebsocketMessage) {
	l.websocketMessageLogger.Infof("Receiving %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
}

func (l *WebsocketLogger) LogMessageOut(msg *types.WebsocketMessage) {
	l.websocketMessageLogger.Infof("Sending %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
}
