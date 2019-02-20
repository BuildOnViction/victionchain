package tomox

import (
	"sync"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/fatih/set.v0"
	"golang.org/x/sync/syncmap"
	"crypto/ecdsa"
)

const (
	ProtocolName = "tomoX"
	ProtocolVersion = uint64(1)
	ProtocolVersionStr = "1.0"
	expirationCycle   = time.Second
	transmissionCycle = 300 * time.Millisecond
	statusCode           = 0   // used by TomoX protocol
	messagesCode         = 1   // normal TomoX message
	p2pMessageCode       = 127 // peer-to-peer message (to be consumed by the peer, but not forwarded any further)
	DefaultTTL           = 50 // seconds
	DefaultSyncAllowance = 10 // seconds
	messageQueueLimit = 1024
	overflowIdx                    // Indicator of message queue overflow
	signatureLength = 65 // in bytes
	padSizeLimit      = 256 // just an arbitrary number, could be changed without breaking the protocol
	flagsLength     = 1
	signatureFlag = byte(4)
	SizeMask      = byte(3) // mask used to extract the size of payload size field from the flags
)

type Config struct {
	DataDir string `toml:",omitempty"`
}

var AllowedPairs = map[string]*big.Int{
	"TOMO/WETH": big.NewInt(10e9),
}

type TomoX struct {
	// Order related
	Orderbooks map[string]*OrderBook
	db         *BatchDatabase
	// pair and max volume ...
	allowedPairs map[string]*big.Int

	// P2P messaging related
	protocol p2p.Protocol
	quit chan struct{}
	peers  map[*Peer]struct{} // Set of currently active peers
	peerMu sync.RWMutex       // Mutex to sync the active peer set

	messageQueue chan *Envelope // Message queue for normal TomoX messages
	p2pMsgQueue  chan *Envelope // Message queue for peer-to-peer messages (not to be forwarded any further)

	envelopes   map[common.Hash]*Envelope
	expirations map[uint32]*set.SetNonTS  // Message expiration pool
	poolMu      sync.RWMutex  // Mutex to sync the message and expiration pools

	syncAllowance int // maximum time in seconds allowed to process the tomoX-related messages

	lightClient bool // indicates is this node is pure light client (does not forward any messages)

	statsMu sync.Mutex // guard stats

	settings syncmap.Map // holds configuration settings that can be dynamically changed

}

func New(cfg *Config) *TomoX {
	datadir := cfg.DataDir
	batchDB := NewBatchDatabaseWithEncode(datadir, 0, 0,
		EncodeBytesItem, DecodeBytesItem)

	fixAllowedPairs := make(map[string]*big.Int)
	for key, value := range AllowedPairs {
		fixAllowedPairs[strings.ToLower(key)] = value
	}

	tomoX := &TomoX{
		Orderbooks:   make(map[string]*OrderBook),
		db:           batchDB,
		allowedPairs: fixAllowedPairs,
		peers:         make(map[*Peer]struct{}),
		quit:          make(chan struct{}),
		envelopes:     make(map[common.Hash]*Envelope),
		syncAllowance: DefaultSyncAllowance,
		expirations:   make(map[uint32]*set.SetNonTS),
		messageQueue:  make(chan *Envelope, messageQueueLimit),
		p2pMsgQueue:   make(chan *Envelope, messageQueueLimit),
	}

	tomoX.settings.Store(overflowIdx, false)

	// p2p tomoX sub protocol handler
	tomoX.protocol = p2p.Protocol{
		Name: ProtocolName,
		Version: uint(ProtocolVersion),
		Run: tomoX.HandlePeer,
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version":        ProtocolVersionStr,
			}
		},
	}

	return tomoX
}

// HandlePeer is called by the underlying P2P layer when the TomoX sub-protocol
// connection is negotiated.
func (tomox *TomoX) HandlePeer(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
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
		return err
	}
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
		case  p2pMessageCode:
			// peer-to-peer message, sent directly to peer.
			// this message is not supposed to be forwarded to other peers.
			// these messages are only accepted from the trusted peer.
			if p.trusted {
				var envelope Envelope
				if err := packet.Decode(&envelope); err != nil {
					log.Warn("failed to decode direct message, peer will be disconnected", "peer", p.peer.ID(), "err", err)
					return errors.New("invalid direct message")
				}
				tomox.postEvent(&envelope, true)
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
		tomox.postEvent(envelope, isP2P) // notify the local node about the new message
	}
	return true, nil
}

// postEvent queues the message for further processing.
func (tomox *TomoX) postEvent(envelope *Envelope, isP2P bool) {
	if isP2P {
		tomox.p2pMsgQueue <- envelope
	} else {
		tomox.checkOverflow()
		tomox.messageQueue <- envelope
	}
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

// Overflow returns an indication if the message queue is full.
func (tomox *TomoX) Overflow() bool {
	val, _ := tomox.settings.Load(overflowIdx)
	return val.(bool)
}

// Protocols returns the whisper sub-protocols ran by this particular client.
func (tomox *TomoX) Protocols() []p2p.Protocol {
	return []p2p.Protocol{tomox.protocol}
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

// Stop implements node.Service, stopping the background data propagation thread
// of the TomoX protocol.
func (tomox *TomoX) Stop() error {
	close(tomox.quit)
	log.Info("tomoX stopped")
	return nil
}

// Start implements node.Service, starting the background data propagation thread
// of the TomoX protocol.
func (tomox *TomoX) Start(*p2p.Server) error {
	//log.Info("started tomoX v." + ProtocolVersionStr)
	//go tomox.update()
	//
	//numCPU := runtime.NumCPU()
	//for i := 0; i < numCPU; i++ {
	//	go tomox.processQueue()
	//}

	return nil
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

// ValidatePublicKey checks the format of the given public key.
func ValidatePublicKey(k *ecdsa.PublicKey) bool {
	return k != nil && k.X != nil && k.Y != nil && k.X.Sign() != 0 && k.Y.Sign() != 0
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
/*
=====================================================================================
*/

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
		// check allow pair
		if _, ok := tomox.allowedPairs[name]; !ok {
			return nil, fmt.Errorf("Orderbook not found for pair :%s", pairName)
		}

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

func (tomox *TomoX) ProcessOrder(quote map[string]string) ([]map[string]string, map[string]string) {

	ob, _ := tomox.getAndCreateIfNotExisted(quote["pair_name"])
	var trades []map[string]string
	var orderInBook map[string]string

	if ob != nil {
		// get map as general input, we can set format later to make sure there is no problem
		orderID, err := strconv.ParseUint(quote["order_id"], 10, 64)
		if err == nil {
			// insert
			if orderID == 0 {
				log.Info("Process order")
				trades, orderInBook = ob.ProcessOrder(quote, true)
			} else {
				log.Info("Update order")
				err = ob.UpdateOrder(quote)
				if err != nil {
					log.Info("Update order failed", "quote", quote, "err", err)
				}
			}
		}

	}

	return trades, orderInBook

}

func (tomox *TomoX) CancelOrder(quote map[string]string) error {
	ob, err := tomox.getAndCreateIfNotExisted(quote["pair_name"])
	if ob != nil {
		orderID, err := strconv.ParseUint(quote["order_id"], 10, 64)
		if err == nil {

			price, ok := new(big.Int).SetString(quote["price"], 10)
			if !ok {
				return fmt.Errorf("Price is not correct :%s", quote["price"])
			}

			return ob.CancelOrder(quote["side"], orderID, price)
		}
	}

	return err
}
