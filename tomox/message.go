package tomox

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// MessageParams specifies the exact way a message should be wrapped
// into an Envelope.
type MessageParams struct {
	TTL      uint32
	WorkTime uint32
	Payload  []byte
	Padding  []byte
	Topic    TopicType
}

// SentMessage represents an end-user data packet to transmit through the
// TomoX protocol. These are wrapped into Envelopes that need not be
// understood by intermediate nodes, just forwarded.
type sentMessage struct {
	Raw []byte
}

// ReceivedMessage represents a data packet to be received through the
// TomoX protocol and successfully decrypted.
type ReceivedMessage struct {
	Raw []byte

	Payload []byte
	Padding []byte
	Salt    []byte
	Topic   TopicType

	Sent uint32 // Time when the message was posted into the network
	TTL  uint32 // Maximum time to live allowed for the message

	EnvelopeHash common.Hash // Message envelope hash to act as a unique id
}

// NewSentMessage creates and initializes a non-signed, non-encrypted TomoX message.
func NewSentMessage(params *MessageParams) (*sentMessage, error) {
	const payloadSizeFieldMaxSize = 4
	msg := sentMessage{}
	msg.Raw = make([]byte, 1,
		flagsLength+payloadSizeFieldMaxSize+len(params.Payload)+len(params.Padding)+signatureLength+padSizeLimit)
	msg.Raw[0] = 0 // set all the flags to zero
	msg.addPayloadSizeField(params.Payload)
	msg.Raw = append(msg.Raw, params.Payload...)
	err := msg.appendPadding(params)
	return &msg, err
}

// addPayloadSizeField appends the auxiliary field containing the size of payload
func (msg *sentMessage) addPayloadSizeField(payload []byte) {
	fieldSize := getSizeOfPayloadSizeField(payload)
	field := make([]byte, 4)
	binary.LittleEndian.PutUint32(field, uint32(len(payload)))
	field = field[:fieldSize]
	msg.Raw = append(msg.Raw, field...)
	msg.Raw[0] |= byte(fieldSize)
}

// getSizeOfPayloadSizeField returns the number of bytes necessary to encode the size of payload
func getSizeOfPayloadSizeField(payload []byte) int {
	s := 1
	for i := len(payload); i >= 256; i /= 256 {
		s++
	}
	return s
}

// appendPadding appends the padding specified in params.
// If no padding is provided in params, then random padding is generated.
func (msg *sentMessage) appendPadding(params *MessageParams) error {
	if len(params.Padding) != 0 {
		// padding data was provided by the Dapp, just use it as is
		msg.Raw = append(msg.Raw, params.Padding...)
		return nil
	}

	rawSize := flagsLength + getSizeOfPayloadSizeField(params.Payload) + len(params.Payload)
	odd := rawSize % padSizeLimit
	paddingSize := padSizeLimit - odd
	pad := make([]byte, paddingSize)
	_, err := crand.Read(pad)
	if err != nil {
		return err
	}
	if !validateDataIntegrity(pad, paddingSize) {
		return errors.New("failed to generate random padding of size " + strconv.Itoa(paddingSize))
	}
	msg.Raw = append(msg.Raw, pad...)
	return nil
}

// Wrap bundles the message into an Envelope to transmit over the network.
func (msg *sentMessage) Wrap(options *MessageParams) (envelope *Envelope, err error) {
	if options.TTL == 0 {
		options.TTL = DefaultTTL
	}

	envelope = NewEnvelope(options.TTL, options.Topic, msg)
	if err = envelope.Seal(options); err != nil {
		return nil, err
	}
	return envelope, nil
}

// ValidateAndParse checks the message validity and extracts the fields in case of success.
func (msg *ReceivedMessage) ValidateAndParse() bool {
	end := len(msg.Raw)
	if end < 1 {
		return false
	}

	beg := 1
	payloadSize := 0
	sizeOfPayloadSizeField := int(msg.Raw[0] & SizeMask) // number of bytes indicating the size of payload
	if sizeOfPayloadSizeField != 0 {
		payloadSize = int(bytesToUintLittleEndian(msg.Raw[beg : beg+sizeOfPayloadSizeField]))
		if payloadSize+1 > end {
			return false
		}
		beg += sizeOfPayloadSizeField
		msg.Payload = msg.Raw[beg : beg+payloadSize]
	}

	beg += payloadSize
	msg.Padding = msg.Raw[beg:end]
	return true
}

// hash calculates the SHA3 checksum of the message flags, payload size field, payload and padding.
func (msg *ReceivedMessage) hash() []byte {
	return crypto.Keccak256(msg.Raw)
}
