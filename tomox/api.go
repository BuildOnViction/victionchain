package tomox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	filterTimeout = 300 // filters are considered timeout out after filterTimeout seconds
)

// List of errors
var (
	ErrTooLowPoW            = errors.New("message rejected, PoW too low")
	ErrNoTopics             = errors.New("missing topic(s)")
)

// PublicTomoXAPI provides the TomoX RPC service that can be
// use publicly without security implications.
type PublicTomoXAPI struct {
	w *TomoX

	mu       sync.Mutex
	lastUsed map[string]time.Time // keeps track when a filter was polled for the last time.
}

// NewPublicTomoXAPI create a new RPC TomoX service.
func NewPublicTomoXAPI(w *TomoX) *PublicTomoXAPI {
	api := &PublicTomoXAPI{
		w:        w,
		lastUsed: make(map[string]time.Time),
	}
	return api
}

// Version returns the TomoX sub-protocol version.
func (api *PublicTomoXAPI) Version(ctx context.Context) string {
	return ProtocolVersionStr
}

// Info contains diagnostic information.
type Info struct {
	Memory         int     `json:"memory"`         // Memory size of the floating messages in bytes.
	Messages       int     `json:"messages"`       // Number of floating messages.
	MinPow         float64 `json:"minPow"`         // Minimal accepted PoW
	MaxMessageSize uint32  `json:"maxMessageSize"` // Maximum accepted message size
}

// Info returns diagnostic information about the TomoX node.
func (api *PublicTomoXAPI) Info(ctx context.Context) Info {
	return Info{
		Messages:       len(api.w.messageQueue) + len(api.w.p2pMsgQueue),
	}
}

// MarkTrustedPeer marks a peer trusted, which will allow it to send historic (expired) messages.
// Note: This function is not adding new nodes, the node needs to exists as a peer.
func (api *PublicTomoXAPI) MarkTrustedPeer(ctx context.Context, enode string) (bool, error) {
	n, err := discover.ParseNode(enode)
	if err != nil {
		return false, err
	}
	return true, api.w.AllowP2PMessagesFromPeer(n.ID[:])
}

// MakeLightClient turns the node into light client, which does not forward
// any incoming messages, and sends only messages originated in this node.
func (api *PublicTomoXAPI) MakeLightClient(ctx context.Context) bool {
	api.w.lightClient = true
	return api.w.lightClient
}

// CancelLightClient cancels light client mode.
func (api *PublicTomoXAPI) CancelLightClient(ctx context.Context) bool {
	api.w.lightClient = false
	return !api.w.lightClient
}

//go:generate gencodec -type NewMessage -field-override newMessageOverride -out gen_newmessage_json.go

// NewMessage represents a new TomoX message that is posted through the RPC.
type NewMessage struct {
	TTL        uint32    `json:"ttl"`
	Payload    []byte    `json:"payload"`
	Padding    []byte    `json:"padding"`
	PowTime    uint32    `json:"powTime"`
	PowTarget  float64   `json:"powTarget"`
	TargetPeer string    `json:"targetPeer"`
}

type newMessageOverride struct {
	Payload   hexutil.Bytes
	Padding   hexutil.Bytes
}

// Post a message on the TomoX network.
func (api *PublicTomoXAPI) Post(ctx context.Context, req NewMessage) (bool, error) {
	var (
		err         error
	)

	params := &MessageParams{
		TTL:      req.TTL,
		Payload:  req.Payload,
		Padding:  req.Padding,
		WorkTime: req.PowTime,
	}

	// encrypt and sent message
	whisperMsg, err := NewSentMessage(params)
	if err != nil {
		return false, err
	}

	env, err := whisperMsg.Wrap(params)
	if err != nil {
		return false, err
	}

	// send to specific node (skip PoW check)
	if len(req.TargetPeer) > 0 {
		n, err := discover.ParseNode(req.TargetPeer)
		if err != nil {
			return false, fmt.Errorf("failed to parse target peer: %s", err)
		}
		return true, api.w.SendP2PMessage(n.ID[:], env)
	}

	return true, api.w.Send(env)
}

//go:generate gencodec -type Criteria -field-override criteriaOverride -out gen_criteria_json.go

// Criteria holds various filter options for inbound messages.
type Criteria struct {
	MinPow       float64     `json:"minPow"`
	AllowP2P     bool        `json:"allowP2P"`
}

type criteriaOverride struct {
	Sig hexutil.Bytes
}

// Messages set up a subscription that fires events when messages arrive that match
// the given set of criteria.
func (api *PublicTomoXAPI) Messages(ctx context.Context, crit Criteria) (*rpc.Subscription, error) {
	var (
		err         error
	)

	// ensure that the RPC connection supports subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	filter := Filter{
		PoW:      crit.MinPow,
		Messages: make(map[common.Hash]*ReceivedMessage),
		AllowP2P: crit.AllowP2P,
	}

	id, err := api.w.Subscribe(&filter)
	if err != nil {
		return nil, err
	}

	// create subscription and start waiting for message events
	rpcSub := notifier.CreateSubscription()
	go func() {
		// for now poll internally, refactor TomoX internal for channel support
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if filter := api.w.GetFilter(id); filter != nil {
					for _, rpcMessage := range toMessage(filter.Retrieve()) {
						if err := notifier.Notify(rpcSub.ID, rpcMessage); err != nil {
							log.Error("Failed to send notification", "err", err)
						}
					}
				}
			case <-rpcSub.Err():
				api.w.Unsubscribe(id)
				return
			case <-notifier.Closed():
				api.w.Unsubscribe(id)
				return
			}
		}
	}()

	return rpcSub, nil
}

//go:generate gencodec -type Message -field-override messageOverride -out gen_message_json.go

// Message is the RPC representation of a TomoX message.
type Message struct {
	TTL       uint32    `json:"ttl"`
	Timestamp uint32    `json:"timestamp"`
	Payload   []byte    `json:"payload"`
	Padding   []byte    `json:"padding"`
	PoW       float64   `json:"pow"`
	Hash      []byte    `json:"hash"`
	Dst       []byte    `json:"recipientPublicKey,omitempty"`
}

type messageOverride struct {
	Payload hexutil.Bytes
	Padding hexutil.Bytes
	Hash    hexutil.Bytes
	Dst     hexutil.Bytes
}

// ToTomoXMessage converts an internal message into an API version.
func ToTomoXMessage(message *ReceivedMessage) *Message {
	msg := Message{
		Payload:   message.Payload,
		Padding:   message.Padding,
		Timestamp: message.Sent,
		TTL:       message.TTL,
		Hash:      message.EnvelopeHash.Bytes(),
	}

	return &msg
}

// toMessage converts a set of messages to its RPC representation.
func toMessage(messages []*ReceivedMessage) []*Message {
	msgs := make([]*Message, len(messages))
	for i, msg := range messages {
		msgs[i] = ToTomoXMessage(msg)
	}
	return msgs
}

// GetFilterMessages returns the messages that match the filter criteria and
// are received between the last poll and now.
func (api *PublicTomoXAPI) GetFilterMessages(id string) ([]*Message, error) {
	api.mu.Lock()
	f := api.w.GetFilter(id)
	if f == nil {
		api.mu.Unlock()
		return nil, fmt.Errorf("filter not found")
	}
	api.lastUsed[id] = time.Now()
	api.mu.Unlock()

	receivedMessages := f.Retrieve()
	messages := make([]*Message, 0, len(receivedMessages))
	for _, msg := range receivedMessages {
		messages = append(messages, ToTomoXMessage(msg))
	}

	return messages, nil
}

// DeleteMessageFilter deletes a filter.
func (api *PublicTomoXAPI) DeleteMessageFilter(id string) (bool, error) {
	api.mu.Lock()
	defer api.mu.Unlock()

	delete(api.lastUsed, id)
	return true, api.w.Unsubscribe(id)
}

// NewMessageFilter creates a new filter that can be used to poll for
// (new) messages that satisfy the given criteria.
func (api *PublicTomoXAPI) NewMessageFilter(req Criteria) (string, error) {
	var (
		err error
	)


	f := &Filter{
		PoW:      req.MinPow,
		AllowP2P: req.AllowP2P,
		Messages: make(map[common.Hash]*ReceivedMessage),
	}

	id, err := api.w.Subscribe(f)
	if err != nil {
		return "", err
	}

	api.mu.Lock()
	api.lastUsed[id] = time.Now()
	api.mu.Unlock()

	return id, nil
}