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

type OHLCVEndpoint struct {
	ohlcvService interfaces.OHLCVService
}

func ServeOHLCVResource(
	r *mux.Router,
	ohlcvService interfaces.OHLCVService,
) {
	e := &OHLCVEndpoint{ohlcvService}
	r.HandleFunc("/ohlcv", e.handleGetOHLCV).Methods("GET")
	ws.RegisterChannel(ws.OHLCVChannel, e.ohlcvWebSocket)
}

func processTimeInterval(i string) (string, int) {
	var unit string
	var duration int

	switch i {
	case "1m":
		unit = "min"
		duration = 1
		break
	case "3m":
		unit = "min"
		duration = 3
		break
	case "5m":
		unit = "min"
		duration = 5
		break
	case "15m":
		unit = "min"
		duration = 15
		break
	case "30m":
		unit = "min"
		duration = 30
		break
	case "1h":
		unit = "hour"
		duration = 1
		break
	case "2h":
		unit = "hour"
		duration = 2
		break
	case "4h":
		unit = "hour"
		duration = 4
		break
	case "6h":
		unit = "hour"
		duration = 6
		break
	case "8h":
		unit = "hour"
		duration = 8
		break
	case "12h":
		unit = "hour"
		duration = 12
		break
	case "1d":
		unit = "day"
		duration = 1
		break
	case "1w":
		unit = "week"
		duration = 1
		break
	case "1mo":
		unit = "month"
		duration = 1
		break
	default:
		unit = "hour"
		duration = 1
		break
	}

	return unit, duration
}

func (e *OHLCVEndpoint) handleGetOHLCV(w http.ResponseWriter, r *http.Request) {
	var p types.OHLCVParams

	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")
	from := v.Get("from")
	to := v.Get("to")
	timeInterval := v.Get("timeInterval")

	if timeInterval == "" {
		httputils.WriteError(w, http.StatusBadRequest, "timeInterval Parameter is missing")
		return
	}

	unit, duration := processTimeInterval(timeInterval)

	p.Units = unit
	p.Duration = int64(duration)

	now := time.Now()

	if to == "" {
		p.To = now.Unix()
	} else {
		t, _ := strconv.Atoi(to)
		p.To = int64(t)
	}

	if from == "" {
		p.From = now.AddDate(-1, 0, 0).Unix()
	} else {
		f, _ := strconv.Atoi(from)
		p.From = int64(f)
	}

	if bt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter is missing")
		return
	}

	if qt == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter is missing")
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

	p.Pair = []types.PairAddresses{{
		BaseToken:  common.HexToAddress(bt),
		QuoteToken: common.HexToAddress(qt),
	}}

	res, err := e.ohlcvService.GetOHLCV(p.Pair, p.Duration, p.Units, p.From, p.To)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []*types.Tick{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *OHLCVEndpoint) ohlcvWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOHLCVSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		socket.SendErrorMessage(c, "Invalid payload")
		return
	}

	if ev.Type == types.SUBSCRIBE {
		b, _ = json.Marshal(ev.Payload)
		var p *types.SubscriptionPayload

		err = json.Unmarshal(b, &p)
		if err != nil {
			logger.Error(err)
		}

		if (p.BaseToken == common.Address{}) {
			socket.SendErrorMessage(c, "Invalid base token")
			return
		}

		if (p.QuoteToken == common.Address{}) {
			socket.SendErrorMessage(c, "Invalid Quote Token")
			return
		}

		now := time.Now()

		if p.From == 0 {
			p.From = now.AddDate(-1, 0, 0).Unix()
		}

		if p.To == 0 {
			p.To = now.Unix()
		}

		if p.Duration == 0 {
			p.Duration = 24
		}

		if p.Units == "" {
			p.Units = "hour"
		}

		e.ohlcvService.Subscribe(c, p)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.ohlcvService.Unsubscribe(c)
	}
}
