package tomox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/tomox/tomox_state"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	filterTimeout                   = 300 // filters are considered timeout out after filterTimeout seconds
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

// Info contains diagnostic information.
type Info struct {
	Memory         int    `json:"memory"`         // Memory size of the floating messages in bytes.
	Messages       int    `json:"messages"`       // Number of floating messages.
	MaxMessageSize uint32 `json:"maxMessageSize"` // Maximum accepted message size
}

// Info returns diagnostic information about the tomoX node.
func (api *PublicTomoXAPI) Info(ctx context.Context) Info {
	return Info{
		Messages: len(api.t.messageQueue) + len(api.t.p2pMsgQueue),
	}
}

// MarkTrustedPeer marks a peer trusted, which will allow it to send historic (expired) messages.
// Note: This function is not adding new nodes, the node needs to exists as a peer.
func (api *PublicTomoXAPI) MarkTrustedPeer(ctx context.Context, enode string) (bool, error) {
	n, err := discover.ParseNode(enode)
	if err != nil {
		return false, err
	}
	return true, api.t.AllowP2PMessagesFromPeer(n.ID[:])
}

// MakeLightClient turns the node into light client, which does not forward
// any incoming messages, and sends only messages originated in this node.
func (api *PublicTomoXAPI) MakeLightClient(ctx context.Context) bool {
	api.t.lightClient = true
	return api.t.lightClient
}

// CancelLightClient cancels light client mode.
func (api *PublicTomoXAPI) CancelLightClient(ctx context.Context) bool {
	api.t.lightClient = false
	return !api.t.lightClient
}

//go:generate gencodec -type NewMessage -field-override newMessageOverride -out gen_newmessage_json.go

// NewMessage represents a new tomoX message that is posted through the RPC.
type NewMessage struct {
	TTL        uint32    `json:"ttl"`
	Topic      TopicType `json:"topic"`
	Payload    []byte    `json:"payload"`
	Padding    []byte    `json:"padding"`
	PowTime    uint32    `json:"powTime"`
	PowTarget  float64   `json:"powTarget"`
	TargetPeer string    `json:"targetPeer"`
}

type newMessageOverride struct {
	Payload hexutil.Bytes
	Padding hexutil.Bytes
}

// Post a message on the TomoX network.
func (api *PublicTomoXAPI) CreateOrder(ctx context.Context, req NewMessage) (bool, error) {
	var (
		err error
	)

	params := &MessageParams{
		TTL:      req.TTL,
		Payload:  req.Payload,
		Padding:  req.Padding,
		WorkTime: req.PowTime,
		Topic:    req.Topic,
	}

	if params.Topic == (TopicType{}) {
		log.Error("Missing topic(s)", "params.Topic", params.Topic)
		return false, ErrNoTopics
	}

	// encrypt and sent message
	tomoXMsg, err := NewSentMessage(params)
	if err != nil {
		return false, err
	}

	env, err := tomoXMsg.Wrap(params)
	if err != nil {
		return false, err
	}

	// send to specific node
	if len(req.TargetPeer) > 0 {
		n, err := discover.ParseNode(req.TargetPeer)
		if err != nil {
			return false, fmt.Errorf("failed to parse target peer: %s", err)
		}
		return true, api.t.SendP2PMessage(n.ID[:], env)
	}

	return true, api.t.Send(env)
}

func (api *PublicTomoXAPI) CancelOrder(ctx context.Context, req NewMessage) (bool, error) {
	params := &MessageParams{
		TTL:      req.TTL,
		Payload:  req.Payload,
		Padding:  req.Padding,
		WorkTime: req.PowTime,
		Topic:    req.Topic,
	}
	payload := &tomox_state.OrderItem{}
	err := json.Unmarshal(params.Payload, &payload)
	if err != nil {
		log.Error("Wrong order payload format", "err", err)
		return false, err
	}
	//set cancel signature to the order payload
	payload.Status = OrderStatusCancelled
	//then encode it again
	params.Payload, err = json.Marshal(payload)
	if err != nil {
		log.Error("Can't encode order payload", "err", err)
		return false, err
	}
	if params.Topic == (TopicType{}) {
		log.Error("Missing topic(s)", "params.Topic", params.Topic)
		return false, ErrNoTopics
	}

	// encrypt and sent message
	tomoXMsg, err := NewSentMessage(params)
	if err != nil {
		return false, err
	}

	env, err := tomoXMsg.Wrap(params)
	if err != nil {
		return false, err
	}

	return true, api.t.Send(env)
}

//go:generate gencodec -type Criteria -field-override criteriaOverride -out gen_criteria_json.go

// Criteria holds various filter options for inbound messages.
type Criteria struct {
	Topics   []TopicType `json:"topics"`
	AllowP2P bool        `json:"allowP2P"`
}

type criteriaOverride struct {
	Sig hexutil.Bytes
}

// Messages set up a subscription that fires events when messages arrive that match
// the given set of criteria.
func (api *PublicTomoXAPI) Messages(ctx context.Context, crit Criteria) (*rpc.Subscription, error) {
	var (
		err error
	)

	// ensure that the RPC connection supports subscriptions
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return nil, rpc.ErrNotificationsUnsupported
	}

	filter := Filter{
		Messages: make(map[common.Hash]*ReceivedMessage),
		AllowP2P: crit.AllowP2P,
	}

	for i, bt := range crit.Topics {
		if len(bt) == 0 || len(bt) > 4 {
			return nil, fmt.Errorf("subscribe: topic %d has wrong size: %d", i, len(bt))
		}
		filter.Topics = append(filter.Topics, bt[:])
	}

	// listen for message that are encrypted with the given symmetric key
	if len(filter.Topics) == 0 {
		return nil, ErrNoTopics
	}

	id, err := api.t.Subscribe(&filter)
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
				if filter := api.t.GetFilter(id); filter != nil {
					for _, rpcMessage := range toMessage(filter.Retrieve()) {
						if err := notifier.Notify(rpcSub.ID, rpcMessage); err != nil {
							log.Error("Failed to send notification", "err", err)
						}
					}
				}
			case <-rpcSub.Err():
				api.t.Unsubscribe(id)
				return
			case <-notifier.Closed():
				api.t.Unsubscribe(id)
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
	Topic     TopicType `json:"topic"`
	Payload   []byte    `json:"payload"`
	Padding   []byte    `json:"padding"`
	Hash      []byte    `json:"hash"`
}

type messageOverride struct {
	Payload hexutil.Bytes
	Padding hexutil.Bytes
	Hash    hexutil.Bytes
}

// ToTomoXMessage converts an internal message into an API version.
func ToTomoXMessage(message *ReceivedMessage) *Message {
	msg := Message{
		Payload:   message.Payload,
		Padding:   message.Padding,
		Timestamp: message.Sent,
		TTL:       message.TTL,
		Hash:      message.EnvelopeHash.Bytes(),
		Topic:     message.Topic,
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

// GetOrders returns the orders that match the filter criteria and
// are received between the last poll and now.
func (api *PublicTomoXAPI) GetOrders(id string) ([]*Message, error) {
	api.mu.Lock()
	f := api.t.GetFilter(id)
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

// DeleteTopic deletes a topic.
func (api *PublicTomoXAPI) DeleteTopic(id string) (bool, error) {
	api.mu.Lock()
	defer api.mu.Unlock()

	delete(api.lastUsed, id)
	return true, api.t.Unsubscribe(id)
}

// NewTopic creates a new topic that can be used to poll for
// (new) messages that satisfy the given criteria.
func (api *PublicTomoXAPI) NewTopic(req Criteria) (string, error) {
	var (
		topics [][]byte
		err    error
	)

	if len(req.Topics) > 0 {
		topics = make([][]byte, len(req.Topics))
		for i, topic := range req.Topics {
			topics[i] = make([]byte, TopicLength)
			copy(topics[i], topic[:])
		}
	}

	f := &Filter{
		AllowP2P: req.AllowP2P,
		Topics:   topics,
		Messages: make(map[common.Hash]*ReceivedMessage),
	}

	id, err := api.t.Subscribe(f)
	if err != nil {
		return "", err
	}

	api.mu.Lock()
	api.lastUsed[id] = time.Now()
	api.mu.Unlock()

	return id, nil
}

// TODO: for testing purpose only, remove in production
// PurgePendingOrders remove all pending orders
func (api *PublicTomoXAPI) PurgePendingOrders() error {
	pending := api.t.getPendingOrders()
	for _, p := range pending {
		if err := api.t.RemoveOrderFromPending(p.Hash, p.Cancel); err != nil {
			log.Error("Failed to purge pending hash", "err", err)
			return err
		}
		if err := api.t.RemoveOrderPendingFromDB(p.Hash, p.Cancel); err != nil {
			log.Error("Failed to purge pending orders", "err", err)
			return err
		}
	}
	return nil
}

// GetOrderNonce returns the latest orderNonce of the given address
func (api *PublicTomoXAPI) GetOrderNonce(address common.Address) (*big.Int, error) {
	return api.t.GetOrderNextNonce(address)
}

// GetBestBid returns the bestBid price of the given pair
func (api *PublicTomoXAPI) GetBestBid(pairName string) (*big.Int, error) {
	ob, err := api.t.getAndCreateIfNotExisted(pairName, false, common.Hash{})
	if err != nil {
		return big.NewInt(0), err
	}
	return ob.BestBid(false, common.Hash{}), nil
}

// GetBestAsk returns the bestAsk price of the given pair
func (api *PublicTomoXAPI) GetBestAsk(pairName string) (*big.Int, error) {
	ob, err := api.t.getAndCreateIfNotExisted(pairName, false, common.Hash{})
	if err != nil {
		return big.NewInt(0), err
	}
	return ob.BestAsk(false, common.Hash{}), nil
}

// GetBidTree returns the bidTreeItem of the given pair
func (api *PublicTomoXAPI) GetBidTree(pairName string) (*OrderTreeItem, error) {
	ob, err := api.t.getAndCreateIfNotExisted(pairName, false, common.Hash{})
	if err != nil {
		return nil, err
	}
	return ob.Bids.Item, nil
}

// GetAskTree returns the askTreeItem of the given pair
func (api *PublicTomoXAPI) GetAskTree(pairName string) (*OrderTreeItem, error) {
	ob, err := api.t.getAndCreateIfNotExisted(pairName, false, common.Hash{})
	if err != nil {
		return nil, err
	}
	return ob.Asks.Item, nil
}

// GetPendingOrders returns pending orders of the given pair
func (api *PublicTomoXAPI) GetPendingOrders(pairName string) ([]*tomox_state.OrderItem, error) {
	result := []*tomox_state.OrderItem{}
	pending := api.t.getPendingOrders()
	for _, p := range pending {
		order := api.t.getOrderPendingFromDB(p.Hash, p.Cancel)
		if order != nil && strings.ToLower(order.PairName) == strings.ToLower(pairName) {
			result = append(result, order)
		}
	}
	return result, nil
}

// GetAllPendingHashes returns all pending order hashes
func (api *PublicTomoXAPI) GetAllPendingHashes() ([]OrderPending, error) {
	pending := api.t.getPendingOrders()
	return pending, nil
}

// GetProcessedHashes returns hashes which already processed
func (api *PublicTomoXAPI) GetProcessedHashes() ([]common.Hash, error) {
	result := []common.Hash{}
	for _, val := range api.t.processedOrderCache.Keys() {
		hash, ok := val.(common.Hash)
		if ok {
			result = append(result, hash)
		}
	}
	return result, nil
}
