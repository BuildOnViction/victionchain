package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type tradeEndpoint struct {
	tradeService interfaces.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeTradeResource(
	r *mux.Router,
	tradeService interfaces.TradeService,
) {
	e := &tradeEndpoint{tradeService}
	r.HandleFunc("/trades", e.HandleGetTrades)
	r.HandleFunc("/trades/history", e.HandleGetTradesHistory)
	ws.RegisterChannel(ws.TradeChannel, e.tradeWebsocket)
}

// HandleGetTrades is responsible for getting pair's trade history requests
func (e *tradeEndpoint) HandleGetTrades(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")
	l := v.Get("limit")
	fromParam := v.Get("from")
	toParam := v.Get("to")

	if bt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter missing")
		return
	}

	if qt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter missing")
		return
	}

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid base token address")
		return
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid quote token address")
		return
	}

	// Client must provides both "from" and "to" or none of them
	if (fromParam != "" && toParam == "") || (toParam != "" && fromParam == "") {
		httputils.WriteError(w, http.StatusBadRequest, "Both \"from\" and \"to\" are required")
		return
	}

	var from, to int64
	now := time.Now()

	if toParam == "" {
		to = now.Unix()
	} else {
		t, _ := strconv.Atoi(toParam)
		to = int64(t)
	}

	if fromParam == "" {
		from = now.AddDate(-1, 0, 0).Unix()
	} else {
		f, _ := strconv.Atoi(fromParam)
		from = int64(f)
	}

	limit := types.DefaultLimit
	if l != "" {
		limit, _ = strconv.Atoi(l)
	}

	baseToken := common.HexToAddress(bt)
	quoteToken := common.HexToAddress(qt)
	res, err := e.tradeService.GetSortedTrades(baseToken, quoteToken, from, to, limit)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

// HandleGetTradesHistory is responsible for handling user's trade history requests
func (e *tradeEndpoint) HandleGetTradesHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")
	baseToken := v.Get("baseToken")
	quoteToken := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	// Client must provides both tokens or none of them
	if (baseToken != "" && quoteToken == "") || (quoteToken != "" && baseToken == "") {
		httputils.WriteError(w, http.StatusBadRequest, "Both token addresses are required")
		return
	}

	if baseToken != "" && !common.IsHexAddress(baseToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if quoteToken != "" && !common.IsHexAddress(quoteToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	// Client must provides both "from" and "to" or none of them
	if (fromParam != "" && toParam == "") || (toParam != "" && fromParam == "") {
		httputils.WriteError(w, http.StatusBadRequest, "Both \"from\" and \"to\" are required")
		return
	}

	var from, to int64
	now := time.Now()

	if toParam == "" {
		to = now.Unix()
	} else {
		t, _ := strconv.Atoi(toParam)
		to = int64(t)
	}

	if fromParam == "" {
		from = now.AddDate(-1, 0, 0).Unix()
	} else {
		f, _ := strconv.Atoi(fromParam)
		from = int64(f)
	}

	lim := types.DefaultLimit
	if limit != "" {
		lim, _ = strconv.Atoi(limit)
	}

	address := common.HexToAddress(addr)

	var baseTokenAddr, quoteTokenAddr common.Address
	if baseToken != "" && quoteToken != "" {
		baseTokenAddr = common.HexToAddress(baseToken)
		quoteTokenAddr = common.HexToAddress(quoteToken)
	} else {
		baseTokenAddr = common.Address{}
		quoteTokenAddr = common.Address{}
	}

	res, err := e.tradeService.GetSortedTradesByUserAddress(address, baseTokenAddr, quoteTokenAddr, from, to, lim)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tradeEndpoint) tradeWebsocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	if err := json.Unmarshal(b, &ev); err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetTradeSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload
	err := json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		if (p.BaseToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid base token"}
			socket.SendErrorMessage(c, err)
			return
		}

		if (p.QuoteToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid quote token"}
			socket.SendErrorMessage(c, err)
			return
		}

		e.tradeService.Subscribe(c, p.BaseToken, p.QuoteToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.tradeService.Unsubscribe(c)
			return
		}

		e.tradeService.UnsubscribeChannel(c, p.BaseToken, p.QuoteToken)
	}
}
