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
	"github.com/ethereum/go-ethereum/crypto"
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
	orderCountKey        = "ORDER_COUNT"
	activePairsKey       = "ACTIVE_PAIRS"
	pendingHash          = "PENDING_HASH"
	pendingPrefix        = "XP"
	processedHash        = "PROCESSED_HASH"
	orderProcessLimit    = 5
)

type Config struct {
	DataDir       string `toml:",omitempty"`
	DBEngine      string `toml:",omitempty"`
	ConnectionUrl string `toml:",omitempty"`
}

type TxDataMatch struct {
	order  []byte
	trades []map[string]string
	obOld  common.Hash
	obNew  common.Hash
	askOld common.Hash
	askNew common.Hash
	bidOld common.Hash
	bidNew common.Hash
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: node.DefaultDataDir(),
}

type TomoX struct {
	// Order related
	Orderbooks map[string]*OrderBook
	db         OrderDao
	orderCount map[common.Address]*big.Int

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

	activePairs map[string]bool // hold active pairs
}

func NewLDBEngine(cfg *Config) *BatchDatabase {
	datadir := cfg.DataDir
	batchDB := NewBatchDatabaseWithEncode(datadir, 0, 0)
	return batchDB
}

func NewMongoDBEngine(cfg *Config) *MongoDatabase {
	mongoDB, err := NewMongoDatabase(nil, cfg.ConnectionUrl)

	if err != nil {
		log.Error(err.Error())
		return &MongoDatabase{}
	}

	return mongoDB
}

func New(cfg *Config) *TomoX {
	tomoX := &TomoX{
		Orderbooks:    make(map[string]*OrderBook),
		orderCount:    make(map[common.Address]*big.Int),
		peers:         make(map[*Peer]struct{}),
		quit:          make(chan struct{}),
		envelopes:     make(map[common.Hash]*Envelope),
		syncAllowance: DefaultSyncAllowance,
		expirations:   make(map[uint32]*set.SetNonTS),
		messageQueue:  make(chan *Envelope, messageQueueLimit),
		p2pMsgQueue:   make(chan *Envelope, messageQueueLimit),
		activePairs:   make(map[string]bool),
	}
	switch cfg.DBEngine {
	case "leveldb":
		tomoX.db = NewLDBEngine(cfg)
	case "mongodb":
		tomoX.db = NewMongoDBEngine(cfg)
	default:
		log.Crit("wrong database engine, only accept either leveldb or mongodb")
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

	if order.Status == Cancel {
		if err := tomox.CancelOrder(order); err != nil {
			log.Error("Can't cancel order", "order", order, "err", err)
			return err
		}
		log.Debug("Cancelled order", "order", order)
	} else {
		if err := tomox.InsertOrder(order); err != nil {
			log.Error("Can't insert order", "order", order, "err", err)
			return err
		}
		log.Debug("Inserted order", "order", order)
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
	key := crypto.Keccak256([]byte(name)) //name is already in lower format
	val, err := tomox.db.Get(key, &OrderBookItem{})
	if val == nil {
		if err != nil {
			log.Error("Can't get orderbook in DB", "err", err)
		}
		return false
	}
	if val.(*OrderBookItem) == nil {
		return false
	}
	return true
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
		log.Debug("Create new orderbook", "ob", ob)

		// updating new pairs
		if len(tomox.activePairs) == 0 {
			if pairs, err := tomox.loadPairs(); err == nil {
				tomox.activePairs = pairs
			}
		}

		if _, ok := tomox.activePairs[name]; !ok {
			tomox.activePairs[name] = true
			if err := tomox.updatePairs(tomox.activePairs); err != nil {
				log.Error("Failed to save active pairs", "err", err)
			}
		}

		return ob, nil
	} else {
		ob := NewOrderBook(name, tomox.db)
		if err := ob.Restore(); err != nil {
			log.Debug("Can't restore orderbook", "err", err)
			return nil, err
		}
		return ob, nil
	}
}

func (tomox *TomoX) GetOrder(pairName, orderID string) *Order {
	ob, _ := tomox.getAndCreateIfNotExisted(pairName)
	if ob == nil {
		return nil
	}
	key := GetKeyFromString(orderID)
	return ob.GetOrder(key)
}

func (tomox *TomoX) InsertOrder(order *OrderItem) error {
	ob, err := tomox.getAndCreateIfNotExisted(order.PairName)
	if err != nil {
		return err
	}

	if ob != nil {
		// insert
		if order.OrderID == 0 {
			if err := tomox.verifyOrderNonce(order); err != nil {
				return err
			}
			// Save order into orderbook tree.
			if err := tomox.addPendingHash(order.Hash); err != nil {
				return err
			}
			if err := tomox.addOrderPending(order); err != nil {
				return err
			}

			log.Info("Process saved")
			tomox.orderCount[order.UserAddress] = order.Nonce
			if err := tomox.updateOrderCount(tomox.orderCount); err != nil {
				log.Error("Failed to save orderCount", "err", err)
			}

		} else {
			log.Info("Update order")
			if err := ob.UpdateOrder(order); err != nil {
				log.Error("Update order failed", "order", order, "err", err)
				return err
			}
		}
	}

	return nil
}

func (tomox *TomoX) verifyOrderNonce(order *OrderItem) error {
	var (
		orderCount *big.Int
		ok         bool
	)

	// in case of restarting nodes, data in memory has lost
	// should load from persistent storage
	if len(tomox.orderCount) == 0 {
		if err := tomox.loadOrderCount(); err != nil {
			// if a node has just started, its database doesn't have orderCount information
			// Hence, we should not throw error here
			log.Debug("orderCount is empty in leveldb", "err", err)
		}
	}
	if orderCount, ok = tomox.orderCount[order.UserAddress]; !ok {
		orderCount = big.NewInt(0)
	}

	if order.Nonce.Cmp(orderCount) <= 0 {
		return ErrOrderNonceTooLow
	}
	distance := Sub(order.Nonce, orderCount)
	if distance.Cmp(new(big.Int).SetUint64(LimitThresholdOrderNonceInQueue)) > 0 {
		return ErrOrderNonceTooHigh
	}
	return nil
}

// load orderCount from persistent storage
func (tomox *TomoX) loadOrderCount() error {
	var orderCount map[common.Address]*big.Int
	val, err := tomox.db.Get([]byte(orderCountKey), orderCount)
	if err != nil {
		return err
	}
	tomox.orderCount = val.(map[common.Address]*big.Int)
	return nil
}

// update orderCount to persistent storage
func (tomox *TomoX) updateOrderCount(orderCount map[common.Address]*big.Int) error {
	if err := tomox.db.Put([]byte(orderCountKey), orderCount); err != nil {
		return err
	}
	return nil
}

func (tomox *TomoX) CancelOrder(order *OrderItem) error {
	ob, err := tomox.getAndCreateIfNotExisted(order.PairName)
	if ob != nil && err == nil {
		return ob.CancelOrder(order)
	}

	return err
}

func (tomox *TomoX) GetBidsTree(pairName string) (*OrderTree, error) {
	ob, err := tomox.GetOrderBook(pairName)
	if err != nil {
		return nil, err
	}
	return ob.Bids, nil
}

func (tomox *TomoX) GetAsksTree(pairName string) (*OrderTree, error) {
	ob, err := tomox.GetOrderBook(pairName)
	if err != nil {
		return nil, err
	}
	return ob.Asks, nil
}

func (tomox *TomoX) ProcessOrderPending() map[common.Hash]TxDataMatch {
	txMatches := make(map[common.Hash]TxDataMatch)
	pendingHashes := tomox.getPendingHashes()
	if len(pendingHashes) > 0 {
		for i, orderHash := range pendingHashes {
			if i <= orderProcessLimit {
				order := tomox.getOrderPending(orderHash)
				if order != nil {
					if !tomox.existProcessedOrderHash(orderHash) {
						ob, err := tomox.getAndCreateIfNotExisted(order.PairName)
						if err != nil || ob == nil {
							log.Error("Fail to get/create orderbook", "order.PairName", order.PairName)
							continue
						}

						log.Info("Process order pending", "orderPending", order)
						obOld, err := ob.Hash()
						if err != nil {
							log.Error("Fail to get orderbook hash old", "err", err)
						}
						askOld, err := ob.Asks.Hash()
						if err != nil {
							log.Error("Fail to get ask tree hash old", "err", err)
						}
						bidOld, err := ob.Bids.Hash()
						if err != nil {
							log.Error("Fail to get bid tree hash old", "err", err)
						}

						// Process ME.
						trades, _ := ob.ProcessOrder(order, true)

						obNew, err := ob.Hash()
						if err != nil {
							log.Error("Fail to get orderbook hash new", "err", err)
						}
						askNew, err := ob.Asks.Hash()
						if err != nil {
							log.Error("Fail to get ask tree hash new", "err", err)
						}
						bidNew, err := ob.Bids.Hash()
						if err != nil {
							log.Error("Fail to get bid tree hash new", "err", err)
						}

						value, err := EncodeBytesItem(order)
						if err != nil {
							log.Error("Can't encode", "order", order, "err", err)
							continue
						} else {
							txMatches[order.Hash] = TxDataMatch{
								order:  value,
								trades: trades,
								obOld:  obOld,
								obNew:  obNew,
								askOld: askOld,
								askNew: askNew,
								bidOld: bidOld,
								bidNew: bidNew,
							}
						}
						if err := tomox.addProcessedOrderHash(orderHash, 1000); err != nil {
							log.Error("Fail to save processed order hash", "err", err)
							continue
						}
					}

					// Remove order from db pending.
					if err := tomox.removePendingHash(orderHash); err != nil {
						log.Error("Fail to remove pending hash", "err", err)
						continue
					}
					if err := tomox.removeOrderPending(orderHash); err != nil {
						log.Error("Fail to remove order pending", "err", err)
					}
				} else {
					log.Error("Fail to get order pending from db", "hash", orderHash)
				}
			}
		}
	}

	return txMatches
}

func (tomox *TomoX) getOrderPending(orderHash common.Hash) *OrderItem {
	var (
		val interface{}
		err error
	)
	prefix := []byte(pendingPrefix)
	key := append(prefix, orderHash.Bytes()...)
	if ok, _ := tomox.db.Has(key); ok {
		val, err = tomox.db.Get(key, val)
		if err != nil {
			log.Error("Fail to get order pending", "err", err)

			return nil
		}
	}

	if val == nil {
		return nil
	}

	return val.(*OrderItem)
}

func (tomox *TomoX) addOrderPending(order *OrderItem) error {
	prefix := []byte(pendingPrefix)
	key := append(prefix, order.Hash.Bytes()...)
	// Insert new order pending.
	if err := tomox.db.Put(key, order); err != nil {
		log.Error("Fail to save order pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) removeOrderPending(orderHash common.Hash) error {
	prefix := []byte(pendingPrefix)
	key := append(prefix, orderHash.Bytes()...)

	if err := tomox.db.Delete(key, true); err != nil {
		log.Error("Fail to delete order pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) addPendingHash(orderHash common.Hash) error {
	pendingHashes := tomox.getPendingHashes()
	if pendingHashes == nil {
		return nil
	}
	find := false
	for _, v := range pendingHashes {
		if v == orderHash {
			find = true
		}
	}
	if !find {
		pendingHashes = append(pendingHashes, orderHash)
	}
	// Store pending hash.
	key := []byte(pendingHash)
	if err := tomox.db.Put(key, pendingHashes); err != nil {
		log.Error("Fail to save order hash pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) removePendingHash(orderHash common.Hash) error {
	pendingHashes := tomox.getPendingHashes()
	if pendingHashes == nil {
		return nil
	}
	for i, v := range pendingHashes {
		if v == orderHash {
			pendingHashes = append(pendingHashes[:i], pendingHashes[i+1:]...)
		}
	}
	// Store pending hash.
	if err := tomox.db.Put([]byte(pendingHash), pendingHashes); err != nil {
		log.Error("Fail to delete order hash pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) getPendingHashes() []common.Hash {
	var (
		val interface{}
		err error
	)
	key := []byte(pendingHash)
	if ok, _ := tomox.db.Has(key); ok {
		if val, err = tomox.db.Get(key, val); err != nil {
			log.Error("Fail to get pending hash", "err", err)

			return []common.Hash{}
		}
	}

	if val == nil {
		return []common.Hash{}
	}
	pendingHashes := val.([]common.Hash)

	return pendingHashes
}

func (tomox *TomoX) addProcessedOrderHash(orderHash common.Hash, limit int) error {
	key := []byte(processedHash)

	processedHashes := tomox.getProcessedOrderHash()
	if len(processedHashes) >= limit {
		// Remove first.
		processedHashes = processedHashes[1:]
	}
	processedHashes = append(processedHashes, orderHash)

	if err := tomox.db.Put(key, processedHashes); err != nil {
		log.Error("Fail to save processed order hashes", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) existProcessedOrderHash(orderHash common.Hash) bool {
	processedHashes := tomox.getProcessedOrderHash()

	for _, k := range processedHashes {
		if k == orderHash {
			return true
			break
		}
	}

	return false
}

func (tomox *TomoX) getProcessedOrderHash() []common.Hash {
	var (
		val interface{}
		err error
	)

	key := []byte(processedHash)
	if ok, _ := tomox.db.Has(key); ok {
		if val, err = tomox.db.Get(key, val); err != nil {
			log.Error("Failed to load processed order hashes", "err", err)
			return nil
		}
	}

	var processedHashes []common.Hash
	if val != nil {
		processedHashes = val.([]common.Hash)
	}

	return processedHashes
}

func (tomox *TomoX) updatePairs(pairs map[string]bool) error {
	if err := tomox.db.Put([]byte(activePairsKey), pairs); err != nil {
		return err
	}
	return nil
}

func (tomox *TomoX) loadPairs() (map[string]bool, error) {
	var (
		val interface{}
		err error
	)
	if val, err = tomox.db.Get([]byte(activePairsKey), val); err != nil {
		log.Error("Failed to load active pairs:", "err", err)
		return map[string]bool{}, err
	}
	pairs := val.(map[string]bool)
	activePairs := map[string]bool{}
	for pairName := range pairs {
		if pairs[pairName] {
			activePairs[pairName] = pairs[pairName]
		}
	}
	return activePairs, nil
}

func (tomox *TomoX) listTokenPairs() []string {
	var activePairs []string
	if len(tomox.activePairs) == 0 {
		if pairs, err := tomox.loadPairs(); err == nil {
			tomox.activePairs = pairs
		}
	}
	for p := range tomox.activePairs {
		activePairs = append(activePairs, p)
	}
	return activePairs
}
