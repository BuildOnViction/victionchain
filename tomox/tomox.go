package tomox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/sync/syncmap"
	"gopkg.in/fatih/set.v0"
)

const (
	ProtocolName         = "tomoX"
	ProtocolVersion      = uint64(1)
	ProtocolVersionStr   = "1.0"
	expirationCycle      = time.Second
	transmissionCycle    = 300 * time.Millisecond
	statusCode           = 10  // used by TomoX protocol
	messagesCode         = 11  // normal TomoX message
	p2pMessageCode       = 127 // peer-to-peer message (to be consumed by the peer, but not forwarded any further)
	NumberOfMessageCodes = 128
	DefaultTTL           = 50 // seconds
	DefaultSyncAllowance = 10 // seconds
	messageQueueLimit    = 1024
	overflowIdx                // Indicator of message queue overflow
	signatureLength      = 65  // in bytes
	padSizeLimit         = 256 // just an arbitrary number, could be changed without breaking the protocol
	flagsLength          = 1
	SizeMask             = byte(3) // mask used to extract the size of payload size field from the flags
	TopicLength          = 86      // in bytes
	keyIDSize            = 32      // in bytes
)

type Config struct {
	DataDir string `toml:",omitempty"`
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: node.DefaultDataDir(),
}

var AllowedPairs = map[string]*big.Int{
	"TOMO/WETH": big.NewInt(10e9),
}

type TomoX struct {
	// Order related
	Orderbooks map[string]*OrderBook
	db         *BatchDatabase

	// P2P messaging related
	protocol p2p.Protocol
	filters  *Filters // Message filters installed with Subscribe function
	quit     chan struct{}
	peers    map[*Peer]struct{} // Set of currently active peers
	peerMu   sync.RWMutex       // Mutex to sync the active peer set

	messageQueue chan *Envelope // Message queue for normal TomoX messages
	p2pMsgQueue  chan *Envelope // Message queue for peer-to-peer messages (not to be forwarded any further)

	envelopes   map[common.Hash]*Envelope
	expirations map[uint32]*set.SetNonTS // Message expiration pool
	poolMu      sync.RWMutex             // Mutex to sync the message and expiration pools

	syncAllowance int // maximum time in seconds allowed to process the tomoX-related messages

	lightClient bool // indicates is this node is pure light client (does not forward any messages)

	statsMu sync.Mutex // guard stats

	settings syncmap.Map // holds configuration settings that can be dynamically changed

}

func New(cfg *Config) *TomoX {
	datadir := cfg.DataDir
	batchDB := NewBatchDatabaseWithEncode(datadir, 0, 0,
		EncodeBytesItem, DecodeBytesItem)

	tomoX := &TomoX{
		Orderbooks: make(map[string]*OrderBook),
		db:         batchDB,
		peers:         make(map[*Peer]struct{}),
		quit:          make(chan struct{}),
		envelopes:     make(map[common.Hash]*Envelope),
		syncAllowance: DefaultSyncAllowance,
		expirations:   make(map[uint32]*set.SetNonTS),
		messageQueue:  make(chan *Envelope, messageQueueLimit),
		p2pMsgQueue:   make(chan *Envelope, messageQueueLimit),
	}
	tomoX.filters = NewFilters(tomoX)

	tomoX.settings.Store(overflowIdx, false)

	// p2p tomoX sub protocol handler
	tomoX.protocol = p2p.Protocol{
		Name:    ProtocolName,
		Version: uint(ProtocolVersion),
		Length:  NumberOfMessageCodes,
		Run:     tomoX.HandlePeer,
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version": ProtocolVersionStr,
			}
		},
	}

	return tomoX
}

// Overflow returns an indication if the message queue is full.
func (tomox *TomoX) Overflow() bool {
	val, _ := tomox.settings.Load(overflowIdx)
	return val.(bool)
}

// APIs returns the RPC descriptors the TomoX implementation offers
func (tomox *TomoX) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: ProtocolName,
			Version:   ProtocolVersionStr,
			Service:   NewPublicTomoXAPI(tomox),
			Public:    true,
		},
	}
}

// Protocols returns the whisper sub-protocols ran by this particular client.
func (tomox *TomoX) Protocols() []p2p.Protocol {
	return []p2p.Protocol{tomox.protocol}
}

// Version returns the TomoX sub-protocols version number.
func (tomox *TomoX) Version() uint {
	return tomox.protocol.Version
}

func (tomox *TomoX) getPeers() []*Peer {
	arr := make([]*Peer, len(tomox.peers))
	i := 0
	tomox.peerMu.Lock()
	for p := range tomox.peers {
		arr[i] = p
		i++
	}
	tomox.peerMu.Unlock()
	return arr
}

// getPeer retrieves peer by ID
func (tomox *TomoX) getPeer(peerID []byte) (*Peer, error) {
	tomox.peerMu.Lock()
	defer tomox.peerMu.Unlock()
	for p := range tomox.peers {
		id := p.peer.ID()
		if bytes.Equal(peerID, id[:]) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Could not find peer with ID: %x", peerID)
}

// AllowP2PMessagesFromPeer marks specific peer trusted,
// which will allow it to send historic (expired) messages.
func (tomox *TomoX) AllowP2PMessagesFromPeer(peerID []byte) error {
	p, err := tomox.getPeer(peerID)
	if err != nil {
		return err
	}
	p.trusted = true
	return nil
}

// SendP2PMessage sends a peer-to-peer message to a specific peer.
func (tomox *TomoX) SendP2PMessage(peerID []byte, envelope *Envelope) error {
	p, err := tomox.getPeer(peerID)
	if err != nil {
		return err
	}
	return tomox.SendP2PDirect(p, envelope)
}

// SendP2PDirect sends a peer-to-peer message to a specific peer.
func (tomox *TomoX) SendP2PDirect(peer *Peer, envelope *Envelope) error {
	return p2p.Send(peer.ws, p2pMessageCode, envelope)
}

// Subscribe installs a new message handler used for filtering, decrypting
// and subsequent storing of incoming messages.
func (tomox *TomoX) Subscribe(f *Filter) (string, error) {
	s, err := tomox.filters.Install(f)
	return s, err
}

// GetFilter returns the filter by id.
func (tomox *TomoX) GetFilter(id string) *Filter {
	return tomox.filters.Get(id)
}

// Unsubscribe removes an installed message handler.
func (tomox *TomoX) Unsubscribe(id string) error {
	ok := tomox.filters.Uninstall(id)
	if !ok {
		return fmt.Errorf("Unsubscribe: Invalid ID")
	}
	return nil
}

// Send injects a message into the whisper send queue, to be distributed in the
// network in the coming cycles.
func (tomox *TomoX) Send(envelope *Envelope) error {
	ok, err := tomox.add(envelope, false)
	if err == nil && !ok {
		return fmt.Errorf("failed to add envelope")
	}
	return err
}

// Start implements node.Service, starting the background data propagation thread
// of the TomoX protocol.
func (tomox *TomoX) Start(*p2p.Server) error {
	log.Info("started tomoX v." + ProtocolVersionStr)
	go tomox.update()

	numCPU := runtime.NumCPU()
	for i := 0; i < numCPU; i++ {
		go tomox.processQueue()
	}

	return nil
}

// Stop implements node.Service, stopping the background data propagation thread
// of the TomoX protocol.
func (tomox *TomoX) Stop() error {
	close(tomox.quit)
	log.Info("tomoX stopped")
	return nil
}

// HandlePeer is called by the underlying P2P layer when the TomoX sub-protocol
// connection is negotiated.
func (tomox *TomoX) HandlePeer(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	log.Debug("TomoX handshake start", "peer", peer.Name())
	// Create the new peer and start tracking it
	tomoPeer := newPeer(tomox, peer, rw)

	tomox.peerMu.Lock()
	tomox.peers[tomoPeer] = struct{}{}
	tomox.peerMu.Unlock()

	defer func() {
		tomox.peerMu.Lock()
		delete(tomox.peers, tomoPeer)
		tomox.peerMu.Unlock()
	}()

	// Run the peer handshake and state updates
	if err := tomoPeer.handshake(); err != nil {
		log.Error("TomoX handshake failed", "peer", peer.Name(), "err", err)
		return err
	}
	log.Debug("TomoX handshake success", "peer", peer.Name())
	tomoPeer.start()
	defer tomoPeer.stop()

	return tomox.runMessageLoop(tomoPeer, rw)
}

// runMessageLoop reads and processes inbound messages directly to merge into client-global state.
func (tomox *TomoX) runMessageLoop(p *Peer, rw p2p.MsgReadWriter) error {
	for {
		// fetch the next packet
		packet, err := rw.ReadMsg()
		if err != nil {
			log.Warn("message loop", "peer", p.peer.ID(), "err", err)
			return err
		}

		switch packet.Code {
		case statusCode:
			// this should not happen, but no need to panic; just ignore this message.
			log.Warn("unxepected status message received", "peer", p.peer.ID())
		case messagesCode:
			// decode the contained envelopes
			var envelopes []*Envelope
			if err := packet.Decode(&envelopes); err != nil {
				log.Warn("failed to decode envelopes, peer will be disconnected", "peer", p.peer.ID(), "err", err)
				return errors.New("invalid envelopes")
			}

			trouble := false
			for _, env := range envelopes {
				cached, err := tomox.add(env, tomox.lightClient)
				if err != nil {
					trouble = true
					log.Error("bad envelope received, peer will be disconnected", "peer", p.peer.ID(), "err", err)
				}
				if cached {
					p.mark(env)
				}
			}

			if trouble {
				return errors.New("invalid envelope")
			}
		case p2pMessageCode:
			// peer-to-peer message, sent directly to peer.
			// this message is not supposed to be forwarded to other peers.
			// these messages are only accepted from the trusted peer.
			if p.trusted {
				var envelope Envelope
				if err := packet.Decode(&envelope); err != nil {
					log.Warn("failed to decode direct message, peer will be disconnected", "peer", p.peer.ID(), "err", err)
					return errors.New("invalid direct message")
				}
				err := tomox.postEvent(&envelope, true)
				if err != nil {
					return err
				}
			}
		default:
		}

		packet.Discard()
	}
}

// add inserts a new envelope into the message pool to be distributed within the
// TomoX network. It also inserts the envelope into the expiration pool at the
// appropriate time-stamp. In case of error, connection should be dropped.
// param isP2P indicates whether the message is peer-to-peer (should not be forwarded).
func (tomox *TomoX) add(envelope *Envelope, isP2P bool) (bool, error) {
	now := uint32(time.Now().Unix())
	sent := envelope.Expiry - envelope.TTL

	if sent > now {
		if sent-DefaultSyncAllowance > now {
			return false, fmt.Errorf("envelope created in the future [%x]", envelope.Hash())
		}
	}

	if envelope.Expiry < now {
		if envelope.Expiry+DefaultSyncAllowance*2 < now {
			return false, fmt.Errorf("very old message")
		}
		log.Debug("expired envelope dropped", "hash", envelope.Hash().Hex())
		return false, nil // drop envelope without error
	}

	hash := envelope.Hash()

	tomox.poolMu.Lock()
	_, alreadyCached := tomox.envelopes[hash]
	if !alreadyCached {
		tomox.envelopes[hash] = envelope
		if tomox.expirations[envelope.Expiry] == nil {
			tomox.expirations[envelope.Expiry] = set.NewNonTS()
		}
		if !tomox.expirations[envelope.Expiry].Has(hash) {
			tomox.expirations[envelope.Expiry].Add(hash)
		}
	}
	tomox.poolMu.Unlock()

	if alreadyCached {
		log.Trace("tomoX envelope already cached", "hash", envelope.Hash().Hex())
	} else {
		log.Trace("cached tomoX envelope", "hash", envelope.Hash().Hex())
		tomox.statsMu.Lock()
		tomox.statsMu.Unlock()
		err := tomox.postEvent(envelope, isP2P) // notify the local node about the new message
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// postEvent queues the message for further processing.
func (tomox *TomoX) postEvent(envelope *Envelope, isP2P bool) error {
	log.Debug("Received envelope", "hash", envelope.hash.Hex())
	if isP2P {
		tomox.p2pMsgQueue <- envelope
	} else {
		tomox.checkOverflow()
		tomox.messageQueue <- envelope
	}

	order := &OrderItem{}
	msg := envelope.Open()
	err := json.Unmarshal(msg.Payload, &order)
	if err != nil {
		log.Error("Fail to parse envelope", "err", err)
		return err
	}

	if order.Type == Cancel {
		err := tomox.CancelOrder(order)
		if err != nil {
			log.Error("Can't cancel order", "err", err)
			return err
		}
		log.Debug("Cancelled order", "detail", order)
	} else {
		log.Info("Save order", "detail", order)
		trades, orderInBook, err := tomox.ProcessOrder(order)
		if err != nil {
			log.Error("Can't process order", "err", err)
			return err
		}
		log.Info("Orderbook result", "Trade", trades, "OrderInBook", orderInBook)
	}
	return nil
}

// checkOverflow checks if message queue overflow occurs and reports it if necessary.
func (tomox *TomoX) checkOverflow() {
	queueSize := len(tomox.messageQueue)

	if queueSize == messageQueueLimit {
		if !tomox.Overflow() {
			tomox.settings.Store(overflowIdx, true)
			log.Warn("message queue overflow")
		}
	} else if queueSize <= messageQueueLimit/2 {
		if tomox.Overflow() {
			tomox.settings.Store(overflowIdx, false)
			log.Warn("message queue overflow fixed (back to normal)")
		}
	}
}

// processQueue delivers the messages to the watchers during the lifetime of the whisper node.
func (tomox *TomoX) processQueue() {
	var e *Envelope
	for {
		select {
		case <-tomox.quit:
			return

		case e = <-tomox.messageQueue:
			tomox.filters.NotifyWatchers(e, false)

		case e = <-tomox.p2pMsgQueue:
			tomox.filters.NotifyWatchers(e, true)
		}
	}
}

// update loops until the lifetime of the whisper node, updating its internal
// state by expiring stale messages from the pool.
func (tomox *TomoX) update() {
	// Start a ticker to check for expirations
	expire := time.NewTicker(expirationCycle)

	// Repeat updates until termination is requested
	for {
		select {
		case <-expire.C:
			tomox.expire()

		case <-tomox.quit:
			return
		}
	}
}

// expire iterates over all the expiration timestamps, removing all stale
// messages from the pools.
func (tomox *TomoX) expire() {
	tomox.poolMu.Lock()
	defer tomox.poolMu.Unlock()

	now := uint32(time.Now().Unix())
	for expiry, hashSet := range tomox.expirations {
		if expiry < now {
			// Dump all expired messages and remove timestamp
			hashSet.Each(func(v interface{}) bool {
				delete(tomox.envelopes, v.(common.Hash))
				return true
			})
			tomox.expirations[expiry].Clear()
			delete(tomox.expirations, expiry)
		}
	}
}

// Envelopes retrieves all the messages currently pooled by the node.
func (tomox *TomoX) Envelopes() []*Envelope {
	tomox.poolMu.RLock()
	defer tomox.poolMu.RUnlock()

	all := make([]*Envelope, 0, len(tomox.envelopes))
	for _, envelope := range tomox.envelopes {
		all = append(all, envelope)
	}
	return all
}

// isEnvelopeCached checks if envelope with specific hash has already been received and cached.
func (tomox *TomoX) isEnvelopeCached(hash common.Hash) bool {
	tomox.poolMu.Lock()
	defer tomox.poolMu.Unlock()

	_, exist := tomox.envelopes[hash]
	return exist
}

// validateDataIntegrity returns false if the data have the wrong or contains all zeros,
// which is the simplest and the most common bug.
func validateDataIntegrity(k []byte, expectedSize int) bool {
	if len(k) != expectedSize {
		return false
	}
	if expectedSize > 3 && containsOnlyZeros(k) {
		return false
	}
	return true
}

// containsOnlyZeros checks if the data contain only zeros.
func containsOnlyZeros(data []byte) bool {
	for _, b := range data {
		if b != 0 {
			return false
		}
	}
	return true
}

// bytesToUintLittleEndian converts the slice to 64-bit unsigned integer.
func bytesToUintLittleEndian(b []byte) (res uint64) {
	mul := uint64(1)
	for i := 0; i < len(b); i++ {
		res += uint64(b[i]) * mul
		mul *= 256
	}
	return res
}

// list Orderbook by topic
func (tomox *TomoX) GetOrderBook(pairName string) (*OrderBook, error) {
	return tomox.getAndCreateIfNotExisted(pairName)
}

func (tomox *TomoX) hasOrderBook(name string) bool {
	_, ok := tomox.Orderbooks[name]
	return ok
}

// commit for all orderbooks
func (tomox *TomoX) Commit() error {
	return tomox.db.Commit()
}

func (tomox *TomoX) getAndCreateIfNotExisted(pairName string) (*OrderBook, error) {

	name := strings.ToLower(pairName)

	if !tomox.hasOrderBook(name) {
		// then create one
		ob := NewOrderBook(name, tomox.db)
		if ob != nil {
			ob.Restore()
			tomox.Orderbooks[name] = ob
		}
	}

	// return from map
	return tomox.Orderbooks[name], nil
}

func (tomox *TomoX) GetOrder(pairName, orderID string) *Order {
	ob, _ := tomox.getAndCreateIfNotExisted(pairName)
	if ob == nil {
		return nil
	}
	key := GetKeyFromString(orderID)
	return ob.GetOrder(key)
}

func (tomox *TomoX) ProcessOrder(order *OrderItem) ([]map[string]string, *OrderItem, error) {
	ob, _ := tomox.getAndCreateIfNotExisted(order.PairName)
	var trades []map[string]string
	var orderInBook *OrderItem

	if ob != nil {
		// insert
		if order.OrderID == 0 {
			// Save order into orderbook tree.
			log.Info("Process saved")
			err := ob.SaveOrderPending(order)
			if err != nil {
				log.Error("Error Save Order Pending", "error", err)
			}
			//log.Info("Process order")
			//trades, orderInBook = ob.ProcessOrder(order, true)
		} else {
			log.Info("Update order")
			err := ob.UpdateOrder(order)
			if err != nil {
				log.Error("Update order failed", "order", order, "err", err)
				return trades, orderInBook, err
			}
		}
	}

	return trades, orderInBook, nil
}

func (tomox *TomoX) CancelOrder(order *OrderItem) error {
	ob, err := tomox.getAndCreateIfNotExisted(order.PairName)
	if ob != nil && err == nil {
		return ob.CancelOrder(order.Side, order.OrderID, order.Price)
	}

	return err
}
