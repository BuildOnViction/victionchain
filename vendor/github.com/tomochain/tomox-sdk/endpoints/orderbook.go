package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type OrderBookEndpoint struct {
	orderBookService interfaces.OrderBookService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeOrderBookResource(
	r *mux.Router,
	orderBookService interfaces.OrderBookService,
) {
	e := &OrderBookEndpoint{orderBookService}
	r.HandleFunc("/orderbook/raw", e.handleGetRawOrderBook)
	r.HandleFunc("/orderbook", e.handleGetOrderBook)
	ws.RegisterChannel(ws.OrderBookChannel, e.orderBookWebSocket)
	ws.RegisterChannel(ws.RawOrderBookChannel, e.rawOrderBookWebSocket)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetOrderBook(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")

	if bt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter missing")
		return
	}

	if qt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter missing")
		return
	}

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetRawOrderBook(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")

	if bt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter missing")
		return
	}

	if qt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter missing")
		return
	}

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetRawOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) rawOrderBookWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetRawOrderBookSocket()

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload

	err = json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.orderBookService.UnsubscribeRawOrderBook(c)
		return
	}

	if (p.BaseToken == common.Address{}) {
		msg := map[string]string{"Message": "Invalid Base Token"}
		socket.SendErrorMessage(c, msg)
		return
	}

	if (p.QuoteToken == common.Address{}) {
		msg := map[string]string{"Message": "Invalid Quote Token"}
		socket.SendErrorMessage(c, msg)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		e.orderBookService.SubscribeRawOrderBook(c, p.BaseToken, p.QuoteToken)
	}
}

func (e *OrderBookEndpoint) orderBookWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOrderBookSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload

	err = json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		msg := map[string]string{"Message": "Internal server error"}
		socket.SendErrorMessage(c, msg)
	}

	if ev.Type == types.SUBSCRIBE {

		if (p.BaseToken == common.Address{}) {
			msg := map[string]string{"Message": "Invalid base token"}
			socket.SendErrorMessage(c, msg)
			return
		}

		if (p.QuoteToken == common.Address{}) {
			msg := map[string]string{"Message": "Invalid quote token"}
			socket.SendErrorMessage(c, msg)
			return
		}

		e.orderBookService.SubscribeOrderBook(c, p.BaseToken, p.QuoteToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.orderBookService.UnsubscribeOrderBook(c)
			return
		}

		e.orderBookService.UnsubscribeOrderBookChannel(c, p.BaseToken, p.QuoteToken)
	}
}
