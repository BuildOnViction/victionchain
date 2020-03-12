package chancoinx

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	LimitThresholdOrderNonceInQueue = 100
)

// List of errors
var (
	ErrNoTopics          = errors.New("missing topic(s)")
	ErrOrderNonceTooLow  = errors.New("OrderNonce too low")
	ErrOrderNonceTooHigh = errors.New("OrderNonce too high")
)

// PublicChancoinXAPI provides the chancoinX RPC service that can be
// use publicly without security implications.
type PublicChancoinXAPI struct {
	t        *ChancoinX
	mu       sync.Mutex
	lastUsed map[string]time.Time // keeps track when a filter was polled for the last time.

}

// NewPublicChancoinXAPI create a new RPC chancoinX service.
func NewPublicChancoinXAPI(t *ChancoinX) *PublicChancoinXAPI {
	api := &PublicChancoinXAPI{
		t:        t,
		lastUsed: make(map[string]time.Time),
	}
	return api
}

// Version returns the ChancoinX sub-protocol version.
func (api *PublicChancoinXAPI) Version(ctx context.Context) string {
	return ProtocolVersionStr
}
