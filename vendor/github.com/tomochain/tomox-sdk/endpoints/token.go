package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/services"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type tokenEndpoint struct {
	tokenService interfaces.TokenService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(
	r *mux.Router,
	tokenService interfaces.TokenService,
) {
	e := &tokenEndpoint{tokenService}
	r.HandleFunc("/tokens/base", e.HandleGetBaseTokens).Methods("GET")
	r.HandleFunc("/tokens/quote", e.HandleGetQuoteTokens).Methods("GET")
	r.HandleFunc("/tokens/{address}", e.HandleGetToken).Methods("GET")
	r.HandleFunc("/tokens", e.HandleGetTokens).Methods("GET")
	r.HandleFunc("/tokens", e.HandleCreateToken).Methods("POST")

	ws.RegisterChannel(ws.TokenChannel, e.ws)
}

func (e *tokenEndpoint) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var t types.Token
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&t)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
	}

	defer r.Body.Close()

	err = e.tokenService.Create(&t)
	if err != nil {
		if err == services.ErrTokenExists {
			httputils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			logger.Error(err)
			httputils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	httputils.WriteJSON(w, http.StatusCreated, t)
}

func (e *tokenEndpoint) HandleGetTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetAll()

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetQuoteTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetQuoteTokens()

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetBaseTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.tokenService.GetBaseTokens()

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tokenEndpoint) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]

	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	tokenAddress := common.HexToAddress(a)
	res, err := e.tokenService.GetByAddress(tokenAddress)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

// ws function handles incoming websocket messages on the order channel
func (e *tokenEndpoint) ws(input interface{}, c *ws.Client) {
	// it means that we can handle not only WebSocketPayload but other Payloads as well
	msg := &types.WebsocketEvent{}
	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		c.SendMessage(ws.TokenChannel, types.ERROR, err.Error())
	}

	switch msg.Type {
	case "GET_TOKENS":
		e.handleGetTokensWS(msg, c)
		log.Printf("Data: %+v", msg)
	default:
		log.Print("Response with error")
	}

}

// handleSubmitSignatures handles NewTrade messages. New trade messages are transmitted to the corresponding order channel
// and received in the handleClientResponse.
func (e *tokenEndpoint) handleGetTokensWS(ev *types.WebsocketEvent, c *ws.Client) {
	res, err := e.tokenService.GetAll()
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.TokenChannel, types.ERROR, err.Error())
		return
	}

	c.SendMessage(ws.TokenChannel, types.UPDATE, res)
}
