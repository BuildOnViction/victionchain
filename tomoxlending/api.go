package tomoxlending

import (
	"context"
	"errors"
	"sync"
	"time"
)


// List of errors
var (
	ErrOrderNonceTooLow  = errors.New("OrderNonce too low")
	ErrOrderNonceTooHigh = errors.New("OrderNonce too high")
)

// PublicTomoXLendingAPI provides the tomoX RPC service that can be
// use publicly without security implications.
type PublicTomoXLendingAPI struct {
	t        *Lending
	mu       sync.Mutex
	lastUsed map[string]time.Time // keeps track when a filter was polled for the last time.

}

// NewPublicTomoXLendingAPI create a new RPC tomoX service.
func NewPublicTomoXLendingAPI(t *Lending) *PublicTomoXLendingAPI {
	api := &PublicTomoXLendingAPI{
		t:        t,
		lastUsed: make(map[string]time.Time),
	}
	return api
}

// Version returns the Lending sub-protocol version.
func (api *PublicTomoXLendingAPI) Version(ctx context.Context) string {
	return ProtocolVersionStr
}
