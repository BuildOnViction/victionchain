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

	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/hashicorp/golang-lru"
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
	orderNonceKey        = "ORDER_NONCES"
	activePairsKey       = "ACTIVE_PAIRS"
	pendingOrder         = "PENDING_ORDER"
	pendingPrefix        = "XP"
	pendingCancelPrefix  = "XPCANCEL"
	orderProcessedLimit  = 1000
	orderProcessLimit    = 20
)

type Config struct {
	DataDir        string `toml:",omitempty"`
	DBEngine       string `toml:",omitempty"`
	DBName         string `toml:",omitempty"`
	ConnectionUrl  string `toml:",omitempty"`
	ReplicaSetName string `toml:",omitempty"`
}

type OrderPending struct {
	Hash   common.Hash
	Cancel bool
}

type OrderProcessed struct {
	Hash   common.Hash
	Cancel bool
}

type TxDataMatch struct {
	Order       []byte // serialized data of order has been processed in this tx
	Trades      []map[string]string
	OrderInBook []byte // serialized data of order which remaining after matching
	ObOld       common.Hash
	ObNew       common.Hash
	AskOld      common.Hash
	AskNew      common.Hash
	BidOld      common.Hash
	BidNew      common.Hash
}

type TxMatchBatch struct {
	Data      []TxDataMatch
	Timestamp uint64
	TxHash    common.Hash
}

// DefaultConfig represents (shocker!) the default configuration.
var DefaultConfig = Config{
	DataDir: node.DefaultDataDir(),
}

type TomoX struct {
	// Order related
	Orderbooks map[string]*OrderBook
	db         OrderDao
	orderNonce map[common.Address]*big.Int

	// P2P messaging related
	protocol p2p.Protocol
	filters  *Filters // Message filters installed with Subscribe function
	quit     chan struct{}
	peers    map[*Peer]struct{} // Set of currently active peers
	peerMu   sync.RWMutex       // Mutex to sync the active peer set

	messageQueue chan *Envelope // Message queue for normal TomoX messages
	p2pMsgQueue  chan *Envelope // Message queue for peer-to-peer messages (not to be forwarded any further)

	processedOrderCache *lru.Cache // Caching processed orders

	envelopes   map[common.Hash]*Envelope
	expirations map[uint32]*set.SetNonTS // Message expiration pool
	poolMu      sync.RWMutex             // Mutex to sync the message and expiration pools

	syncAllowance int // maximum time in seconds allowed to process the tomoX-related messages

	lightClient bool // indicates is this node is pure light client (does not forward any messages)
	sdkNode     bool

	settings syncmap.Map // holds configuration settings that can be dynamically changed

	activePairs map[string]bool // hold active pairs

	tokenDecimalCache *lru.Cache
}

func NewLDBEngine(cfg *Config) *BatchDatabase {
	datadir := cfg.DataDir
	batchDB := NewBatchDatabaseWithEncode(datadir, 0)
	return batchDB
}

func NewMongoDBEngine(cfg *Config) *MongoDatabase {
	mongoDB, err := NewMongoDatabase(nil, cfg.DBName, cfg.ConnectionUrl, cfg.ReplicaSetName, 0)

	if err != nil {
		log.Crit("Failed to init mongodb engine", "err", err)
	}

	return mongoDB
}

func New(cfg *Config) *TomoX {
	poCache, _ := lru.New(orderProcessedLimit)
	tokenDecimalCache, _ := lru.New(defaultCacheLimit)
	tomoX := &TomoX{
		Orderbooks:          make(map[string]*OrderBook),
		orderNonce:          make(map[common.Address]*big.Int),
		peers:               make(map[*Peer]struct{}),
		quit:                make(chan struct{}),
		envelopes:           make(map[common.Hash]*Envelope),
		syncAllowance:       DefaultSyncAllowance,
		expirations:         make(map[uint32]*set.SetNonTS),
		messageQueue:        make(chan *Envelope, messageQueueLimit),
		p2pMsgQueue:         make(chan *Envelope, messageQueueLimit),
		activePairs:         make(map[string]bool),
		processedOrderCache: poCache,
		tokenDecimalCache:   tokenDecimalCache,
	}
	switch cfg.DBEngine {
	case "leveldb":
		tomoX.db = NewLDBEngine(cfg)
		tomoX.sdkNode = false
	case "mongodb":
		tomoX.db = NewMongoDBEngine(cfg)
		tomoX.sdkNode = true
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
	if !tomoX.sdkNode {
		if err := tomoX.loadSnapshot(common.Hash{}); err != nil {
			log.Error("Failed to load tomox snapshot", "err", err)
		}
	}

	return tomoX
}

// Overflow returns an indication if the message queue is full.
func (tomox *TomoX) Overflow() bool {
	val, _ := tomox.settings.Load(overflowIdx)
	return val.(bool)
}

func (tomox *TomoX) IsSDKNode() bool {
	return tomox.sdkNode
}

func (tomox *TomoX) GetDB() OrderDao {
	return tomox.db
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

	if err := tomox.InsertOrder(order); err != nil {
		log.Error("Can't insert order", "order", order, "err", err)
		return nil
	}
	log.Debug("Inserted order to pending", "order", order)
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
func (tomox *TomoX) GetOrderBook(pairName string, dryrun bool, blockHash common.Hash) (*OrderBook, error) {
	return tomox.getAndCreateIfNotExisted(pairName, dryrun, blockHash)
}

func (tomox *TomoX) hasOrderBook(name string, dryrun bool, blockHash common.Hash) bool {
	key := crypto.Keccak256([]byte(name)) //name is already in lower format
	orderBookItemKey := append([]byte(orderbookItemPrefix), key...)
	val, err := tomox.db.Get(orderBookItemKey, &OrderBookItem{}, dryrun, blockHash)
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

func (tomox *TomoX) getAndCreateIfNotExisted(pairName string, dryrun bool, blockHash common.Hash) (*OrderBook, error) {

	name := strings.ToLower(pairName)

	if !tomox.hasOrderBook(name, dryrun, blockHash) {
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
		if err := ob.Restore(dryrun, blockHash); err != nil {
			log.Debug("Can't restore orderbook", "err", err)
			return nil, err
		}
		return ob, nil
	}
}

func (tomox *TomoX) InsertOrder(order *OrderItem) error {
	// TODO: after cache relayer information, please update full verify here
	if err := order.VerifyBasicOrderInfo(); err != nil {
		return err
	}
	if order.OrderID == 0 || order.Status == OrderStatusCancelled {
		if order.Status == OrderStatusCancelled {
			if err := tomox.addOrderToPending(order.Hash, true); err != nil {
				return err
			}
		} else {
			if err := tomox.verifyOrderNonce(order); err != nil {
				return err
			}
			if err := tomox.addOrderToPending(order.Hash, false); err != nil {
				return err
			}
			if err := tomox.UpdateOrderNonce(order.UserAddress, order.Nonce); err != nil {
				log.Error("Failed to update orderNonce", "err", err)
			}
		}
		if err := tomox.saveOrderPendingToDB(order, order.Status == OrderStatusCancelled); err != nil {
			return err
		}
	} else {
		log.Warn("Order has already processed", "orderhash", order.Hash)
	}

	return nil
}

func (tomox *TomoX) verifyOrderNonce(order *OrderItem) error {
	var (
		orderNonce *big.Int
		ok         bool
	)

	// in case of restarting nodes, data in memory has lost
	// should load from persistent storage
	if len(tomox.orderNonce) == 0 {
		if err := tomox.loadOrderNonce(); err != nil {
			// if a node has just started, its database doesn't have orderNonce information
			// Hence, we should not throw error here
			log.Debug("orderNonce is empty in leveldb", "err", err)
		}
	}
	if orderNonce, ok = tomox.orderNonce[order.UserAddress]; !ok {
		orderNonce = big.NewInt(-1)
	}

	if order.Nonce.Cmp(orderNonce) <= 0 {
		return ErrOrderNonceTooLow
	}
	distance := Sub(order.Nonce, orderNonce)
	if distance.Cmp(new(big.Int).SetUint64(LimitThresholdOrderNonceInQueue)) > 0 {
		return ErrOrderNonceTooHigh
	}
	return nil
}

func (tomox *TomoX) GetOrderNonce(address common.Address) (*big.Int, error) {
	if len(tomox.orderNonce) == 0 {
		if err := tomox.loadOrderNonce(); err != nil {
			return big.NewInt(0), nil
		}
	}
	orderNonce, ok := tomox.orderNonce[address]
	if !ok {
		return big.NewInt(0), nil
	}
	return orderNonce, nil
}

// load orderNonce from persistent storage
func (tomox *TomoX) loadOrderNonce() error {
	var (
		orderNonce map[common.Address]*big.Int
		err        error
		val        interface{}
	)
	val, err = tomox.db.Get([]byte(orderNonceKey), &[]byte{}, false, common.Hash{})
	if err != nil {
		return err
	}
	b := *val.(*[]byte)
	if err = json.Unmarshal(b, &orderNonce); err != nil {
		return err
	}
	tomox.orderNonce = orderNonce
	return nil
}

// update orderNonce to persistent storage
func (tomox *TomoX) UpdateOrderNonce(userAddress common.Address, newCount *big.Int) error {
	orderNonceList := tomox.orderNonce
	if orderNonce, ok := orderNonceList[userAddress]; !ok || newCount.Cmp(orderNonce) > 0 {
		orderNonceList[userAddress] = newCount
		blob, err := json.Marshal(orderNonceList)
		if err != nil {
			return err
		}
		if err := tomox.db.Put([]byte(orderNonceKey), &blob, false, common.Hash{}); err != nil {
			return err
		}
	}
	return nil
}

func (tomox *TomoX) GetBidsTree(pairName string, dryrun bool, blockHash common.Hash) (*OrderTree, error) {
	ob, err := tomox.GetOrderBook(pairName, dryrun, blockHash)
	if err != nil {
		return nil, err
	}
	return ob.Bids, nil
}

func (tomox *TomoX) GetAsksTree(pairName string, dryrun bool, blockHash common.Hash) (*OrderTree, error) {
	ob, err := tomox.GetOrderBook(pairName, dryrun, blockHash)
	if err != nil {
		return nil, err
	}
	return ob.Asks, nil
}

func (tomox *TomoX) ProcessOrderPending() []TxDataMatch {
	blockHash := common.StringToHash("COMMIT_NEW_WORK" + time.Now().String())
	txMatches := []TxDataMatch{}
	tomox.db.InitDryRunMode(blockHash)

	log.Debug("Get pending orders to process")
	pendingOrders := tomox.getPendingOrders()
	if len(pendingOrders) > 0 {
		for i, pendingOrder := range pendingOrders {
			if i < orderProcessLimit {
				order := tomox.getOrderPendingFromDB(pendingOrder.Hash, pendingOrder.Cancel)
				if order != nil {
					if pendingOrder.Cancel || !tomox.ExistProcessedOrderHash(pendingOrder.Hash) {
						var (
							ob  *OrderBook
							err error
						)

						// if orderbook has been processed before in this block, it should be in dry-run mode
						// otherwise it's in db
						ob, err = tomox.getAndCreateIfNotExisted(order.PairName, true, blockHash)
						if err != nil || ob == nil {
							log.Error("Fail to get/create orderbook", "order.PairName", order.PairName)
							continue
						}

						log.Info("Process order pending", "orderPending", order)
						obOld, err := ob.Hash()
						if err != nil {
							log.Error("Fail to get orderbook hash old", "err", err)
							continue
						}
						askOld, err := ob.Asks.Hash()
						if err != nil {
							log.Error("Fail to get ask tree hash old", "err", err)
							continue
						}
						bidOld, err := ob.Bids.Hash()
						if err != nil {
							log.Error("Fail to get bid tree hash old", "err", err)
							continue
						}
						originalOrder := &OrderItem{}
						*originalOrder = *order
						originalOrder.Quantity = CloneBigInt(order.Quantity)

						if pendingOrder.Cancel {
							order.Status = OrderStatusCancelled
						}

						trades, orderInBook, err := ob.ProcessOrder(order, true, true, blockHash)

						// remove order from pending list
						if err := tomox.RemoveOrderFromPending(order.Hash, order.Status == OrderStatusCancelled); err != nil {
							continue
						}

						// remove order pending
						if err := tomox.RemoveOrderPendingFromDB(order.Hash, order.Status == OrderStatusCancelled); err != nil {
							continue
						}

						if err != nil {
							log.Error("Can't process order", "order", order, "err", err)
							continue
						}
						obNew, err := ob.Hash()
						if err != nil {
							log.Error("Fail to get orderbook hash new", "err", err)
							continue
						}
						askNew, err := ob.Asks.Hash()
						if err != nil {
							log.Error("Fail to get ask tree hash new", "err", err)
							continue
						}
						bidNew, err := ob.Bids.Hash()
						if err != nil {
							log.Error("Fail to get bid tree hash new", "err", err)
							continue
						}

						// orderID has been updated
						originalOrder.OrderID = order.OrderID
						originalOrderValue, err := EncodeBytesItem(originalOrder)
						if err != nil {
							log.Error("Can't encode", "order", originalOrder, "err", err)
							continue
						}
						orderInBookValue := []byte{}
						if orderInBook != nil {
							orderInBookValue, err = EncodeBytesItem(orderInBook)
							if err != nil {
								log.Error("Can't encode orderInBook", "orderInBook", orderInBook, "err", err)
								continue
							}
						}
						log.Debug("Process OrderPending completed", "orderNonce", order.Nonce, "obNew", hex.EncodeToString(obNew.Bytes()), "bidNew", hex.EncodeToString(bidNew.Bytes()), "askNew", hex.EncodeToString(askNew.Bytes()))
						txMatch := TxDataMatch{
							Order:       originalOrderValue,
							Trades:      trades,
							OrderInBook: orderInBookValue,
							ObOld:       obOld,
							ObNew:       obNew,
							AskOld:      askOld,
							AskNew:      askNew,
							BidOld:      bidOld,
							BidNew:      bidNew,
						}
						txMatches = append(txMatches, txMatch)
					}
				} else {
					log.Error("Fail to get order pending from db", "hash", pendingOrder.Hash)
				}
			}
		}
	}

	return txMatches
}

func (tomox *TomoX) getOrderPendingFromDB(orderHash common.Hash, cancel bool) *OrderItem {
	var (
		val interface{}
		err error
	)
	prefix := []byte(pendingPrefix)
	if cancel {
		prefix = []byte(pendingCancelPrefix)
	}
	key := append(prefix, orderHash.Bytes()...)
	log.Debug("Get order pending", "order", orderHash, "key", hex.EncodeToString(key))
	if ok, _ := tomox.db.Has(key, false, common.Hash{}); ok {
		val, err = tomox.db.Get(key, &OrderItem{}, false, common.Hash{})
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

func (tomox *TomoX) saveOrderPendingToDB(order *OrderItem, cancel bool) error {
	prefix := []byte(pendingPrefix)
	if cancel {
		prefix = []byte(pendingCancelPrefix)
	}
	key := append(prefix, order.Hash.Bytes()...)
	// Insert new order pending.
	log.Debug("Add order pending", "order", order, "key", hex.EncodeToString(key))
	if err := tomox.db.Put(key, order, false, common.Hash{}); err != nil {
		log.Error("Fail to save order pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) RemoveOrderPendingFromDB(orderHash common.Hash, cancel bool) error {
	prefix := []byte(pendingPrefix)
	key := append(prefix, orderHash.Bytes()...)
	log.Debug("Remove order pending", "orderHash", orderHash, "key", hex.EncodeToString(key))
	if err := tomox.db.Delete(key, false, common.Hash{}); err != nil {
		log.Error("Fail to delete order pending", "with prefix", pendingPrefix, "err", err)
		return err
	}

	// cancel will remove both pendingprefix and pendingcancelprefix data.
	if cancel {
		prefix = []byte(pendingCancelPrefix)
		key := append(prefix, orderHash.Bytes()...)
		log.Debug("Remove order pending", "orderHash", orderHash, "key", hex.EncodeToString(key))
		if err := tomox.db.Delete(key, false, common.Hash{}); err != nil {
			log.Error("Fail to delete order pending", "with prefix", pendingCancelPrefix,"err", err)
			return err
		}
	}
	return nil
}

func (tomox *TomoX) addOrderToPending(orderHash common.Hash, cancel bool) error {
	log.Debug("Add order to pending", "orderHash", orderHash, "cancel", cancel)
	pendingOrders := tomox.getPendingOrders()
	if pendingOrders == nil {
		return nil
	}
	find := false
	for _, v := range pendingOrders {
		if v.Hash == orderHash && v.Cancel == cancel {
			find = true
		}
	}
	if !find {
		pendingOrders = append(pendingOrders, OrderPending{Hash: orderHash, Cancel: cancel})
	}
	// Store pending hash.
	key := []byte(pendingOrder)
	if err := tomox.db.Put(key, &pendingOrders, false, common.Hash{}); err != nil {
		log.Error("Fail to add order to pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) RemoveOrderFromPending(orderHash common.Hash, cancel bool) error {
	log.Debug("Remove pending hash", "orderHash", orderHash, "cancel", cancel)
	pendingOrders := tomox.getPendingOrders()
	if pendingOrders == nil {
		return nil
	}
	for i, v := range pendingOrders {
		if v.Hash == orderHash && v.Cancel == cancel {
			pendingOrders = append(pendingOrders[:i], pendingOrders[i+1:]...)
			break
		}
	}
	// Store pending hash.
	if err := tomox.db.Put([]byte(pendingOrder), &pendingOrders, false, common.Hash{}); err != nil {
		log.Error("Fail to delete order hash pending", "err", err)
		return err
	}

	return nil
}

func (tomox *TomoX) getPendingOrders() []OrderPending {
	var (
		val interface{}
		err error
	)
	key := []byte(pendingOrder)
	if ok, _ := tomox.db.Has(key, false, common.Hash{}); ok {
		if val, err = tomox.db.Get(key, &[]OrderPending{}, false, common.Hash{}); err != nil {
			log.Error("Fail to get pending hash", "err", err)
			return []OrderPending{}
		}
	}

	if val == nil {
		return []OrderPending{}
	}
	pendingOrders := *val.(*[]OrderPending)

	return pendingOrders
}

func (tomox *TomoX) addProcessedOrderHash(orderHash common.Hash, cancel bool) error {
	if !tomox.processedOrderCache.Add(orderHash, true) {
		// remove order from pending list
		if err := tomox.RemoveOrderFromPending(orderHash, cancel); err != nil {
			log.Warn("Double check pending order at addProcessedOrderHash. Failed to remove pending hash", "err", err, "orderHash", orderHash)
		}

		// remove order pending
		if err := tomox.RemoveOrderPendingFromDB(orderHash, cancel); err != nil {
			log.Warn("Double check pending order at addProcessedOrderHash. Failed to remove pending order", "err", err, "orderHash", orderHash)
		}
		return nil
	} else {
		return fmt.Errorf("Can't add processed order to cache: orderHash - %s", orderHash.Hex())
	}
}

func (tomox *TomoX) ExistProcessedOrderHash(orderHash common.Hash) bool {
	return tomox.processedOrderCache.Contains(orderHash)
}

func (tomox *TomoX) updatePairs(pairs map[string]bool) error {
	blob, err := json.Marshal(pairs)
	if err != nil {
		return err
	}
	if err := tomox.db.Put([]byte(activePairsKey), &blob, false, common.Hash{}); err != nil {
		return err
	}
	return nil
}

func (tomox *TomoX) loadPairs() (map[string]bool, error) {
	var (
		pairs map[string]bool
		val   interface{}
		err   error
	)
	val, err = tomox.db.Get([]byte(activePairsKey), &[]byte{}, false, common.Hash{})
	if err != nil {
		return map[string]bool{}, err
	}
	b := *val.(*[]byte)
	if err = json.Unmarshal(b, &pairs); err != nil {
		return map[string]bool{}, err
	}
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

func (tomox *TomoX) Snapshot(blockHash common.Hash) error {
	var (
		snap *Snapshot
		err  error
	)
	defer func(start time.Time) {
		if err != nil {
			log.Error("Failed to snapshot ", "err", err, "time", common.PrettyDuration(time.Since(start)), "hash", blockHash)
		} else {
			log.Debug("Snapshot process takes ", "time", common.PrettyDuration(time.Since(start)), "hash", blockHash)
		}
	}(time.Now())

	if snap, err = newSnapshot(tomox, blockHash); err != nil {
		return err
	}
	if err = snap.store(tomox.db); err != nil {
		return err
	}

	if err = tomox.db.Put([]byte(latestSnapshotKey), &blockHash, false, common.Hash{}); err != nil {
		return err
	}

	return nil
}

func (tomox *TomoX) loadSnapshot(hash common.Hash) error {
	// load orderbook from snapshot
	var (
		snap *Snapshot
		val  interface{}
		ob   *OrderBook
		err  error
	)

	defer func(start time.Time) {
		if err != nil {
			log.Error("Failed to load snapshot", "err", err, "time", common.PrettyDuration(time.Since(start)), "hash", hash)
		} else {
			log.Debug("Successfully load snapshot", "time", common.PrettyDuration(time.Since(start)), "hash", hash)
		}
	}(time.Now())

	if hash == (common.Hash{}) {
		if val, err = tomox.db.Get([]byte(latestSnapshotKey), &common.Hash{}, false, common.Hash{}); err != nil {
			// no snapshot found
			return err
		}
		hash = *val.(*common.Hash)
	}
	if snap, err = getSnapshot(tomox.db, hash); err != nil || len(snap.OrderBooks) == 0 {
		return err
	}
	for pair := range snap.OrderBooks {
		ob, err = snap.RestoreOrderBookFromSnapshot(tomox.db, pair)
		if err == nil {
			if err := ob.Save(false, common.Hash{}); err != nil {
				return err
			}
		}
	}
	return nil
}

// save orderbook after matching orders
// update order pending list, processed list
func (tomox *TomoX) ApplyTxMatches(orders []*OrderItem, blockHash common.Hash) error {
	if !tomox.IsSDKNode() {
		if err := tomox.db.SaveDryRunResult(blockHash); err != nil {
			log.Error("Failed to save dry-run result")
			return err
		}
	}

	for _, order := range orders {
		if err := tomox.addProcessedOrderHash(order.Hash, order.Status == OrderStatusCancelled); err != nil {
			log.Error("Failed to mark order as processed", "err", err)
		}
		log.Debug("Mark order as processed", "orderHash", hex.EncodeToString(order.Hash.Bytes()))
		if order.Status != OrderStatusCancelled {
			if err := tomox.UpdateOrderNonce(order.UserAddress, order.Nonce); err != nil {
				log.Error("Update orderNonce via ApplyTxMatches failed", "err", err)
			}
		}
	}
	tomox.db.InitDryRunMode(blockHash)
	return nil
}

// there are 3 tasks need to complete to update data in SDK nodes after matching
// 1. txMatchData.Order: order has been processed. This order should be put to `orders` collection with status sdktypes.OrderStatusOpen
// 2. txMatchData.Trades: includes information of matched orders.
// 		a. Put them to `trades` collection
// 		b. Update status of regrading orders to sdktypes.OrderStatusFilled
// 3. txMatchData.OrderInBook: remaining order after matching. If order has been matched but still remain in orderbook, update status to sdktypes.OrderStatusPartialFilled
func (tomox *TomoX) SyncDataToSDKNode(txDataMatch TxDataMatch, txHash common.Hash, statedb *state.StateDB) error {
	// apply for SDK nodes only
	if !tomox.IsSDKNode() {
		return nil
	}
	var (
		order, orderInBook *OrderItem
		err                error
	)
	db := tomox.GetDB()

	// 1. put processed order to db
	if order, err = txDataMatch.DecodeOrder(); err != nil {
		log.Error("SDK node decode order failed", "txDataMatch", txDataMatch)
		return fmt.Errorf("SDK node decode order failed")
	}

	if order.Status != OrderStatusCancelled {
		order.Status = OrderStatusOpen
	}

	log.Debug("Put processed order", "order", order)
	if err := db.Put(order.Hash.Bytes(), order, false, common.Hash{}); err != nil {
		return fmt.Errorf("SDKNode: failed to put processed order. Error: %s", err.Error())
	}
	if order.Status == OrderStatusCancelled {
		return nil
	}

	// 2. put trades to db and update status to FILLED
	trades := txDataMatch.GetTrades()
	log.Debug("Got trades", "number", len(trades), "trades", trades)
	for _, trade := range trades {

		// 2.a. put to trades
		tradeSDK := &Trade{}
		quantity := ToBigInt(trade[TradeQuantity])
		price := ToBigInt(trade[TradePrice])
		if price.Cmp(big.NewInt(0)) <= 0 || quantity.Cmp(big.NewInt(0)) <= 0 {
			return fmt.Errorf("trade misses important information. tradedPrice %v, tradedQuantity %v", price, quantity)
		}
		tradeSDK.Amount = quantity
		tradeSDK.PricePoint = price
		tradeSDK.PairName = order.PairName
		tradeSDK.BaseToken = order.BaseToken
		tradeSDK.QuoteToken = order.QuoteToken
		tradeSDK.Status = TradeStatusSuccess
		tradeSDK.Taker = order.UserAddress
		tradeSDK.Maker = common.HexToAddress(trade[TradeMaker])
		tradeSDK.TakerOrderHash = order.Hash
		tradeSDK.MakerOrderHash = common.HexToHash(trade[TradeMakerOrderHash])
		tradeSDK.TxHash = txHash
		tradeSDK.TakerOrderSide = order.Side
		tradeSDK.TakerExchange = order.ExchangeAddress
		tradeSDK.MakerExchange = common.HexToAddress(trade[TradeMakerExchange])

		// feeAmount: all fees are calculated in quoteToken
		quoteTokenQuantity := big.NewInt(0).Mul(quantity, price)
		quoteTokenQuantity = big.NewInt(0).Div(quoteTokenQuantity, common.BasePrice)
		takerFee := big.NewInt(0).Mul(quoteTokenQuantity, GetExRelayerFee(order.ExchangeAddress, statedb))
		takerFee = big.NewInt(0).Div(takerFee, common.TomoXBaseFee)
		tradeSDK.TakeFee = takerFee

		makerFee := big.NewInt(0).Mul(quoteTokenQuantity, GetExRelayerFee(common.HexToAddress(trade[TradeMakerExchange]), statedb))
		makerFee = big.NewInt(0).Div(makerFee, common.TomoXBaseFee)
		tradeSDK.MakeFee = makerFee

		tradeSDK.Hash = tradeSDK.ComputeHash()
		log.Debug("TRADE history", "order", order, "trade", tradeSDK)
		if err := db.Put(EmptyKey(), tradeSDK, false, common.Hash{}); err != nil {
			return fmt.Errorf("SDKNode: failed to store tradeSDK %s", err.Error())
		}

		// 2.b. update status and filledAmount
		filledAmount := quantity
		// update order status of relating orders
		if err := tomox.updateStatusOfMatchedOrder(trade[TradeMakerOrderHash], filledAmount); err != nil {
			return err
		}
		if err := tomox.updateStatusOfMatchedOrder(trade[TradeTakerOrderHash], filledAmount); err != nil {
			return err
		}
	}

	// 3. in case of remaining order in orderTree, the status of order should be partial_filled
	if orderInBook, err = txDataMatch.DecodeOrderInBook(); err != nil {
		log.Error("SDK node decode orderInBook failed", "txDataMatch", txDataMatch)
		return fmt.Errorf("SDK node decode orderInBook failed")
	}
	log.Debug("OrderInBook found in blockdata", "order", orderInBook)
	if orderInBook != nil {
		existingOrderInDB := &OrderItem{}
		key := orderInBook.Hash.Bytes()
		data, err := tomox.db.Get(key, &OrderItem{}, false, common.Hash{})
		if data != nil && err == nil {
			existingOrderInDB = data.(*OrderItem)
			log.Debug("OrderInBook found in database", "order", existingOrderInDB)
			if existingOrderInDB.Status == OrderStatusFilled {
				log.Debug("Update status to PARTIAL_FILLED", "order", existingOrderInDB)
				// if order already matched but still exist in orderbook, then update status to OrderStatusPartialFilled
				existingOrderInDB.Status = OrderStatusPartialFilled
				if err = db.Put(key, existingOrderInDB, false, common.Hash{}); err != nil {
					return fmt.Errorf("SDKNode: failed to update status to %s %s", existingOrderInDB.Status, err.Error())
				}
			}
		}
	}
	return nil
}

func (tomox *TomoX) updateStatusOfMatchedOrder(hashString string, filledAmount *big.Int) error {
	db := tomox.GetDB()
	orderHashBytes, err := hex.DecodeString(hashString)
	if err != nil {
		return fmt.Errorf("SDKNode: failed to decode orderKey. Key: %s", hashString)
	}
	val, err := db.Get(orderHashBytes, &OrderItem{}, false, common.Hash{})
	if err != nil || val == nil {
		return fmt.Errorf("SDKNode: failed to get order. Key: %s", hashString)
	}
	matchedOrder := val.(*OrderItem)
	matchedOrder.Status = OrderStatusFilled
	updatedFillAmount := new(big.Int)
	updatedFillAmount.Add(matchedOrder.FilledAmount, filledAmount)
	matchedOrder.FilledAmount = updatedFillAmount
	if err = db.Put(matchedOrder.Hash.Bytes(), matchedOrder, false, common.Hash{}); err != nil {
		return fmt.Errorf("SDKNode: failed to update matchedOrder to sdkNode %s", err.Error())
	}
	return nil
}
