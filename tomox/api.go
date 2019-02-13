package tomox

import (
	"time"
	"sync"
	"context"
)

// PublicWhisperAPI provides the whisper RPC service that can be
// use publicly without security implications.
type PublicTomoXAPI struct {
	t *TomoX

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

// Version returns the tomoX sub-protocol version.
func (api *PublicTomoXAPI) Version(ctx context.Context) string {
	return ProtocolVersionStr
}
