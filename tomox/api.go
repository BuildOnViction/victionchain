package tomox

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

// GetBestBid returns the bestBid price of the given pair
func (api *PublicTomoXAPI) GetBestBid(pairName string) (*big.Int, error) {
	//TODO: get BestBid from tomox state trie
	return big.NewInt(0), nil
}

// GetBestAsk returns the bestAsk price of the given pair
func (api *PublicTomoXAPI) GetBestAsk(pairName string) (*big.Int, error) {
	//TODO: get BestAsk from tomox state trie
	return big.NewInt(0), nil
}

// GetBidTree returns the bidTreeItem of the given pair
func (api *PublicTomoXAPI) GetBidTree(pairName string) (interface{}, error) {
	//TODO: get BidTree from tomox state trie
	return nil, nil
}

// GetAskTree returns the askTreeItem of the given pair
func (api *PublicTomoXAPI) GetAskTree(pairName string) (interface{}, error) {
	//TODO: get AskTree from tomox state trie
	return nil, nil
}

// GetPendingOrders returns pending orders of the given pair
func (api *PublicTomoXAPI) GetPendingOrders(pairName string) ([]*tomox_state.OrderItem, error) {
	result := []*tomox_state.OrderItem{}
	//TODO: get pending orders from orderpool
	return result, nil
}
