package protocol

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tomox"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/protocols"
)

const (
	OrderbookName = "orderbook"
)

var (
	OrderbookProtocol = &protocols.Spec{
		Name:       OrderbookName,
		Version:    42,
		MaxMsgSize: 1024,
		Messages: []interface{}{
			&OrderbookHandshake{},
			&OrderbookMsg{},
			// &OrderbookCancelMsg{},
		},
	}
)

type OrderbookMsg struct {
	PairName  string
	OrderID   string
	Price     string
	Quantity  string
	Side      string
	Timestamp uint64
	TradeID   string
	Type      string
}

// type OrderbookCancelMsg struct {
// 	PairName  string
// 	OrderID   string
// 	Price     string
// 	Side      string
// 	Timestamp uint64
// }

func (msg *OrderbookMsg) ToQuote() map[string]string {
	quote := make(map[string]string)
	quote["timestamp"] = strconv.FormatUint(msg.Timestamp, 10)
	quote["type"] = msg.Type
	quote["side"] = msg.Side
	quote["quantity"] = msg.Quantity
	quote["price"] = msg.Price
	quote["trade_id"] = msg.TradeID
	quote["pair_name"] = msg.PairName
	// if insert id is not used, just for update
	quote["order_id"] = msg.OrderID
	return quote
}

// func (msg *OrderbookCancelMsg) ToQuote() map[string]string {
// 	quote := make(map[string]string)
// 	quote["timestamp"] = strconv.FormatUint(msg.Timestamp, 10)
// 	quote["side"] = msg.Side
// 	quote["price"] = msg.Price
// 	quote["pair_name"] = msg.PairName
// 	quote["order_id"] = msg.OrderID
// 	return quote
// }

func NewOrderbookMsg(quote map[string]string) (*OrderbookMsg, error) {
	timestamp, err := strconv.ParseUint(quote["timestamp"], 10, 64)
	return &OrderbookMsg{
		Timestamp: timestamp,
		Type:      quote["type"],
		Side:      quote["side"],
		Quantity:  quote["quantity"],
		Price:     quote["price"],
		TradeID:   quote["trade_id"],
		PairName:  quote["pair_name"],
		OrderID:   quote["order_id"],
	}, err
}

// func NewOrderbookCancelMsg(quote map[string]string) (*OrderbookCancelMsg, error) {
// 	timestamp, err := strconv.ParseUint(quote["timestamp"], 10, 64)
// 	return &OrderbookCancelMsg{
// 		Timestamp: timestamp,
// 		Side:      quote["side"],
// 		Price:     quote["price"],
// 		PairName:  quote["pair_name"],
// 		OrderID:   quote["order_id"],
// 	}, err
// }

type OrderbookHandshake struct {
	Nick string
	V    uint
}

// the protocols abstraction enables use of an external handler function
type OrderbookHandler struct {
	Engine *tomox.Engine
	Peer   *protocols.Peer
	InC    <-chan interface{}
	QuitC  <-chan struct{}
}

// checkProtoHandshake verifies local and remote protoHandshakes match
func checkProtoHandshake(testVersion uint) func(interface{}) error {
	return func(rhs interface{}) error {
		remote, ok := rhs.(*OrderbookHandshake)

		if ok && remote.V != testVersion {
			return fmt.Errorf("%d (!= %d)", remote.V, testVersion)
		}
		return nil
	}
}

func (orderbookHandler *OrderbookHandler) handleOrderbookMsg(message *OrderbookMsg) error {
	if message.TradeID == "" {
		return orderbookHandler.handleOrderbookCancelMsg(message)
	}
	log.Debug("Received order", "order", message, "peer", orderbookHandler.Peer)

	// add Order
	payload := message.ToQuote()
	log.Info("-> Add order", "payload", payload)

	trades, orderInBook := orderbookHandler.Engine.ProcessOrder(payload)
	log.Info("Orderbook result", "Trade", trades, "OrderInBook", orderInBook)
	return nil
}

func (orderbookHandler *OrderbookHandler) handleOrderbookCancelMsg(message *OrderbookMsg) error {
	log.Debug("Received cancel order", "cancel_order", message, "peer", orderbookHandler.Peer)

	// cancel Order
	payload := message.ToQuote()
	log.Info("-> Cancel order", "payload", payload)

	err := orderbookHandler.Engine.CancelOrder(payload)
	log.Info("Orderbook result", "err", err)
	return nil
}

func (orderbookHandler *OrderbookHandler) handleOrderbookHandshake(orderbookhs *OrderbookHandshake) error {
	log.Debug("Processing handshake", "from", orderbookhs.Nick, "version", orderbookhs.V)

	// now protocol is ok, we can inject channel to receive message
	go func() {
		for {
			select {
			case payload := <-orderbookHandler.InC:
				// log.Info("Internal received", "payload", payload)
				inmsg, ok := payload.(*OrderbookMsg)
				if ok {
					// maybe we have to use map[]chan
					// databytes, err := rlp.EncodeToBytes(inmsg)
					// databytes, err := json.Marshal(inmsg)

					log.Debug("Sending orderbook", "orderbook", inmsg)
					orderbookHandler.Peer.Send(inmsg)

				}

			// send quit command, break this loop
			case <-orderbookHandler.QuitC:
				break
			}
		}
	}()
	return nil
}

// we will receive message in handle
func (orderbookHandler *OrderbookHandler) handle(msg interface{}) error {

	// we got message or handshake

	// demo.LogWarn("Inbout", "inbout", orderbookHandler.Peer.Inbound())

	switch messageType := msg.(type) {
	case *OrderbookMsg:
		return orderbookHandler.handleOrderbookMsg(msg.(*OrderbookMsg))
	// case *OrderbookCancelMsg:
	// 	return orderbookHandler.handleOrderbookCancelMsg(msg.(*OrderbookCancelMsg))
	case *OrderbookHandshake:
		return orderbookHandler.handleOrderbookHandshake(msg.(*OrderbookHandshake))
	default:
		return fmt.Errorf("Unknown orderbook message type :%v", messageType)
	}

}

// create the protocol with the protocols extension
func NewProtocol(inC <-chan interface{}, quitC <-chan struct{}, orderbookEngine *tomox.Engine) *p2p.Protocol {
	return &p2p.Protocol{
		Name:    "Orderbook",
		Version: 42,
		// we may use more 1 custom message code
		Length: uint64(len(OrderbookProtocol.Messages)),
		// Length: 2,
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {

			// demo.LogWarn("running", "peer", p)

			var err error
			// create the enhanced peer, it will wrap p2p.Send with code from Message Spec
			pp := protocols.NewPeer(p, rw, OrderbookProtocol)

			// send the message, then handle it to make sure protocol success
			go func() {
				outmsg := &OrderbookHandshake{
					V: 42,
					// shortened hex string for terminal logging
					Nick: p.Name(),
				}

				// check handshake, should sleep a little bit before sending handshake so that
				// we can run handle first
				// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				// defer cancel()
				// hsCheck := checkProtoHandshake(outmsg.V)
				// _, err = pp.Handshake(ctx, outmsg, hsCheck)
				// if err != nil {
				// 	return
				// }

				// sleep a second
				time.Sleep(time.Second)

				err = pp.Send(outmsg)
				if err != nil {
					log.Error("Send p2p message fail", "err", err)
				}
				log.Info("Sending handshake", "peer", p, "handshake", outmsg)
			}()

			// protocols abstraction provides a separate blocking run loop for the peer
			// when this returns, the protocol will be terminated
			run := &OrderbookHandler{
				Engine: orderbookEngine,
				Peer:   pp,
				// assign channel
				InC:   inC,
				QuitC: quitC,
			}
			err = pp.Run(run.handle)
			return err
		},
	}
}
