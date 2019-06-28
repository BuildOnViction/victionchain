package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type NotificationEndpoint struct {
	NotificationService interfaces.NotificationService
}

// ServeNotificationResource sets up the routing of notification endpoints and the corresponding handlers.
func ServeNotificationResource(
	r *mux.Router,
	notificationService interfaces.NotificationService,
) {
	e := &NotificationEndpoint{notificationService}

	r.HandleFunc("/notifications", e.HandleGetNotifications).Methods("GET")
	r.HandleFunc("/notifications/{id}", e.HandleUpdateNotification).Methods("PUT")

	ws.RegisterChannel(ws.NotificationChannel, e.handleNotificationWebSocket)
}

func (e *NotificationEndpoint) HandleGetNotifications(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	page := v.Get("page")
	perPage := v.Get("perPage")
	userAddress := v.Get("userAddress")

	p, err := strconv.Atoi(page)

	if err != nil {
		err = errors.New(fmt.Sprintf("%s is not an integer.", page))
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if p <= 0 {
		p = 1
	}

	pp, err := strconv.Atoi(perPage)

	if err != nil {
		err = errors.New(fmt.Sprintf("%s is not an integer.", perPage))
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if pp <= 0 || pp%10 != 0 || pp > 50 {
		pp = 10
	}

	if !common.IsHexAddress(userAddress) {
		err = errors.New("Invalid user address")
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	a := common.HexToAddress(userAddress)

	notifications, err := e.NotificationService.GetByUserAddress(a, pp, (p-1)*pp) // limit = perPage, offset = (page-1)*perPage

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, notifications)
}

func (e *NotificationEndpoint) HandleUpdateNotification(w http.ResponseWriter, r *http.Request) {
	var n types.Notification
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&n)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	defer r.Body.Close()

	n.Status = types.StatusRead
	updated, err := e.NotificationService.Update(&n)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, updated)
}

func (e *NotificationEndpoint) handleNotificationWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	if ev.Type != types.SUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		ws.SendNotificationErrorMessage(c, err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		b, _ = json.Marshal(ev.Payload)
		var addr string

		err = json.Unmarshal(b, &addr)
		if err != nil {
			logger.Error(err)
			ws.SendNotificationErrorMessage(c, err)
			return
		}

		if !common.IsHexAddress(addr) {
			err := map[string]string{"Message": "Invalid address"}
			ws.SendNotificationErrorMessage(c, err)
			return
		}

		a := common.HexToAddress(addr)

		ws.RegisterNotificationConnection(a, c)
		notifications, err := e.NotificationService.GetByUserAddress(a, 0, 0)

		if err != nil {
			logger.Error(err)
			ws.SendNotificationErrorMessage(c, err)
			return
		}

		ws.SendNotificationMessage(types.INIT, a, notifications)
	}
}
