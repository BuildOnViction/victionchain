package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type orderEndpoint struct {
	orderService   interfaces.OrderService
	accountService interfaces.AccountService
}

// ServeOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(
	r *mux.Router,
	orderService interfaces.OrderService,
	accountService interfaces.AccountService,
) {
	e := &orderEndpoint{orderService, accountService}

	r.HandleFunc("/orders/count", e.handleGetCountOrder).Methods("GET")
	r.HandleFunc("/orders/history", e.handleGetOrderHistory).Methods("GET")
	r.HandleFunc("/orders/positions", e.handleGetPositions).Methods("GET")
	r.HandleFunc("/orders", e.handleGetOrders).Methods("GET")
	r.HandleFunc("/orders", e.handleNewOrder).Methods("POST")
	r.HandleFunc("/orders/stop", e.handleNewStopOrder).Methods("POST")
	r.HandleFunc("/orders/cancel", e.handleCancelOrder).Methods("POST")
	r.HandleFunc("/orders/cancelAll", e.handleCancelAllOrders).Methods("POST")
	r.HandleFunc("/orders/stop/cancel", e.handleCancelStopOrder).Methods("POST")

	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) handleGetCountOrder(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter Missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)

	total, err := e.orderService.GetOrderCountByUserAddress(a)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, total)
}

func (e *orderEndpoint) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")
	baseToken := v.Get("baseToken")
	quoteToken := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter Missing")
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

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	var baseTokenAddr, quoteTokenAddr common.Address
	if baseToken != "" && quoteToken != "" {
		baseTokenAddr = common.HexToAddress(baseToken)
		quoteTokenAddr = common.HexToAddress(quoteToken)
	} else {
		baseTokenAddr = common.Address{}
		quoteTokenAddr = common.Address{}
	}

	lim := types.DefaultLimit
	if limit != "" {
		lim, _ = strconv.Atoi(limit)
	}

	orders, err = e.orderService.GetByUserAddress(address, baseTokenAddr, quoteTokenAddr, from, to, lim)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetPositions(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	if limit == "" {
		orders, err = e.orderService.GetCurrentByUserAddress(address)
	} else {
		lim, _ := strconv.Atoi(limit)
		orders, err = e.orderService.GetCurrentByUserAddress(address, lim)
	}

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetOrderHistory(w http.ResponseWriter, r *http.Request) {
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

	// Client must provides with both tokens or none of them
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

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	var baseTokenAddr, quoteTokenAddr common.Address
	if baseToken != "" && quoteToken != "" {
		baseTokenAddr = common.HexToAddress(baseToken)
		quoteTokenAddr = common.HexToAddress(quoteToken)
	} else {
		baseTokenAddr = common.Address{}
		quoteTokenAddr = common.Address{}
	}

	lim := types.DefaultLimit
	if limit != "" {
		lim, _ = strconv.Atoi(limit)
	}

	orders, err = e.orderService.GetHistoryByUserAddress(address, baseTokenAddr, quoteTokenAddr, from, to, lim)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleNewOrder(w http.ResponseWriter, r *http.Request) {
	var o *types.Order
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	o.Hash = o.ComputeHash()

	acc, err := e.accountService.FindOrCreate(o.UserAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if acc.IsBlocked {
		httputils.WriteError(w, http.StatusForbidden, "Account is blocked")
		return
	}

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, o)
}

func (e *orderEndpoint) handleNewStopOrder(w http.ResponseWriter, r *http.Request) {
	var so *types.StopOrder
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&so)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	so.Hash = so.ComputeHash()

	acc, err := e.accountService.FindOrCreate(so.UserAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if acc.IsBlocked {
		httputils.WriteError(w, http.StatusForbidden, "Account is blocked")
		return
	}

	err = e.orderService.NewStopOrder(so)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, so)
}

func (e *orderEndpoint) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	oc := &types.OrderCancel{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err = oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.orderService.CancelOrder(oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, oc.Hash)
}

// handleCancelAllOrder cancels all open/partial filled orders of an user address
func (e *orderEndpoint) handleCancelAllOrders(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "Address parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)

	err := e.orderService.CancelAllOrder(a)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *orderEndpoint) handleCancelStopOrder(w http.ResponseWriter, r *http.Request) {
	oc := &types.OrderCancel{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err = oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.orderService.CancelStopOrder(oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, oc.Hash)
}

// ws function handles incoming websocket messages on the order channel
func (e *orderEndpoint) ws(input interface{}, c *ws.Client) {
	msg := &types.WebsocketEvent{}

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		c.SendMessage(ws.OrderChannel, types.ERROR, err.Error())
	}

	switch msg.Type {
	case "NEW_ORDER":
		e.handleWSNewOrder(msg, c)
	case "NEW_STOP_ORDER":
		e.handleWSNewStopOrder(msg, c)
	case "CANCEL_ORDER":
		e.handleWSCancelOrder(msg, c)
	case "CANCEL_STOP_ORDER":
		e.handleWSCancelStopOrder(msg, c)
	default:
		log.Print("Response with error")
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleWSNewOrder(ev *types.WebsocketEvent, c *ws.Client) {
	o := &types.Order{}

	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.OrderChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
		return
	}

	o.Hash = o.ComputeHash()
	ws.RegisterOrderConnection(o.UserAddress, c)

	acc, err := e.accountService.FindOrCreate(o.UserAddress)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
	}

	if acc.IsBlocked {
		c.SendMessage(ws.OrderChannel, types.ERROR, errors.New("Account is blocked"))
	}

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
		return
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleWSNewStopOrder(ev *types.WebsocketEvent, c *ws.Client) {
	so := &types.StopOrder{}

	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.OrderChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, &so)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, so.Hash)
		return
	}

	so.Hash = so.ComputeHash()
	ws.RegisterOrderConnection(so.UserAddress, c)

	acc, err := e.accountService.FindOrCreate(so.UserAddress)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, so.Hash)
	}

	if acc.IsBlocked {
		c.SendMessage(ws.OrderChannel, types.ERROR, errors.New("Account is blocked"))
	}

	err = e.orderService.NewStopOrder(so)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, so.Hash)
		return
	}
}

// handleCancelOrder handles CancelOrder message.
func (e *orderEndpoint) handleWSCancelOrder(ev *types.WebsocketEvent, c *ws.Client) {
	bytes, err := json.Marshal(ev.Payload)
	oc := &types.OrderCancel{}

	err = json.Unmarshal(bytes, &oc)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	addr, err := oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	ws.RegisterOrderConnection(addr, c)

	orderErr := e.orderService.CancelOrder(oc)
	if orderErr != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(orderErr, oc.Hash)
		return
	}
}

// handleWSCancelStopOrder handles CancelStopOrder message.
func (e *orderEndpoint) handleWSCancelStopOrder(ev *types.WebsocketEvent, c *ws.Client) {
	bytes, err := json.Marshal(ev.Payload)
	oc := &types.OrderCancel{}

	err = json.Unmarshal(bytes, &oc)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	addr, err := oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	ws.RegisterOrderConnection(addr, c)

	orderErr := e.orderService.CancelStopOrder(oc)
	if orderErr != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(orderErr, oc.Hash)
		return
	}
}
