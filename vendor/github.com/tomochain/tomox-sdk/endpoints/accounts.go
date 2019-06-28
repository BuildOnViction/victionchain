package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/middlewares"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type AccountEndpoint struct {
	AccountService interfaces.AccountService
}

func ServeAccountResource(
	r *mux.Router,
	accountService interfaces.AccountService,
) {

	e := &AccountEndpoint{AccountService: accountService}

	r.Handle(
		"/account/create", http.HandlerFunc(e.handleCreateAccount),
	).Methods("POST")

	r.Handle(
		"/account/favorite/{address}",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleGetFavoriteTokens)),
	).Methods("GET")

	r.Handle(
		"/account/favorite/add",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleAddFavoriteToken)),
	).Methods("POST")

	r.Handle(
		"/account/favorite/remove",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleRemoveFavoriteToken)),
	).Methods("POST")

	r.Handle(
		"/account/{address}", http.HandlerFunc(e.handleGetAccount),
	).Methods("GET")

	r.Handle(
		"/account/{address}/{token}", http.HandlerFunc(e.handleGetAccountTokenBalance),
	).Methods("GET")
}

func (e *AccountEndpoint) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)
	existingAccount, err := e.AccountService.GetByAddress(a)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if existingAccount != nil {
		httputils.WriteJSON(w, http.StatusOK, "Account already exists")
		return
	}

	acc := &types.Account{Address: a}
	err = e.AccountService.Create(acc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, acc)
}

func (e *AccountEndpoint) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr := vars["address"]
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)
	a, err := e.AccountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *AccountEndpoint) handleGetAccountTokenBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	t := vars["token"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	b, err := e.AccountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, b)
}

func (e *AccountEndpoint) handleGetFavoriteTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr := vars["address"]
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)

	publicKeyBytes := common.Hex2Bytes(r.Header["Pubkey"][0])
	publicAddress := utils.GetAddressFromPublicKey(publicKeyBytes)

	if address != publicAddress {
		httputils.WriteError(w, http.StatusUnauthorized, "Request is not sent from address's owner")
		return
	}

	a, err := e.AccountService.GetFavoriteTokens(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *AccountEndpoint) handleAddFavoriteToken(w http.ResponseWriter, r *http.Request) {
	var tr *types.FavoriteTokenRequest
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&tr)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if !common.IsHexAddress(tr.Address) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	if !common.IsHexAddress(tr.Token) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
		return
	}

	address := common.HexToAddress(tr.Address)
	tokenAddr := common.HexToAddress(tr.Token)

	publicKeyBytes := common.Hex2Bytes(r.Header["Pubkey"][0])
	publicAddress := utils.GetAddressFromPublicKey(publicKeyBytes)

	if address != publicAddress {
		httputils.WriteError(w, http.StatusUnauthorized, "Request is not sent from address's owner")
		return
	}

	err = e.AccountService.AddFavoriteToken(address, tokenAddr)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, tokenAddr)
}

func (e *AccountEndpoint) handleRemoveFavoriteToken(w http.ResponseWriter, r *http.Request) {
	var tr *types.FavoriteTokenRequest
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&tr)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if !common.IsHexAddress(tr.Address) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	if !common.IsHexAddress(tr.Token) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
		return
	}

	address := common.HexToAddress(tr.Address)
	tokenAddr := common.HexToAddress(tr.Token)

	publicKeyBytes := common.Hex2Bytes(r.Header["Pubkey"][0])
	publicAddress := utils.GetAddressFromPublicKey(publicKeyBytes)

	if address != publicAddress {
		httputils.WriteError(w, http.StatusUnauthorized, "Request is not sent from address's owner")
		return
	}

	err = e.AccountService.DeleteFavoriteToken(address, tokenAddr)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, tokenAddr)
}
