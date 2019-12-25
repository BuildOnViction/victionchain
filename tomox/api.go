package tomox

import (
	"context"
	"errors"
	"github.com/tomochain/tomochain/tomox/tradingstate"
	"math/big"
	"sync"
	"time"

	"github.com/tomochain/tomochain/common"
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

// PublicTomoXAPI provides the tomoX RPC service that can be
// use publicly without security implications.
type PublicTomoXAPI struct {
	t        *TomoX
	mu       sync.Mutex
	lastUsed map[string]time.Time // keeps track when a filter was polled for the last time.

}

// NewPublicTomoXAPI create a new RPC tomoX service.
func NewPublicTomoXAPI(t *TomoX) *PublicTomoXAPI {
	api := &PublicTomoXAPI{
		t:        t,
		lastUsed: make(map[string]time.Time),
	}
	return api
}

// Version returns the TomoX sub-protocol version.
func (api *PublicTomoXAPI) Version(ctx context.Context) string {
	return ProtocolVersionStr
}

// GetOrderNonce returns the latest orderNonce of the given address
func (api *PublicTomoXAPI) GetOrderNonce(address common.Address) (*big.Int, error) {
	//TODO: getOrderNonce from state
	return big.NewInt(0), nil
}

// GetPendingOrders returns pending orders of the given pair
func (api *PublicTomoXAPI) GetPendingOrders(pairName string) ([]*tradingstate.OrderItem, error) {
	result := []*tradingstate.OrderItem{}
	//TODO: get pending orders from orderpool
	return result, nil
}
