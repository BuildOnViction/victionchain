package endpoints

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

type MarketsEndpoint struct {
	MarketsService interfaces.MarketsService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeMarketsResource(
	r *mux.Router,
	marketsService interfaces.MarketsService,
) {
	e := &MarketsEndpoint{marketsService}

	ws.RegisterChannel(ws.MarketsChannel, e.handleMarketsWebSocket)
}

func (e *MarketsEndpoint) handleMarketsWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetMarketSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		e.MarketsService.Subscribe(c)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.MarketsService.Unsubscribe(c)
	}
}
