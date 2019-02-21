package tomox

import (
	"time"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/math"
)

type Envelope struct {
	Expiry uint32
	TTL    uint32
	Data   []byte
	Nonce  uint64
	hash  common.Hash
}

// NewEnvelope wraps a Whisper message with expiration and destination data
// included into an envelope for network forwarding.
func NewEnvelope(ttl uint32, msg *sentMessage) *Envelope {
	env := Envelope{
		Expiry: uint32(time.Now().Add(time.Second * time.Duration(ttl)).Unix()),
		TTL:    ttl,
		Data:   msg.Raw,
		Nonce:  0,
	}

	return &env
}

// Hash returns the SHA3 hash of the envelope, calculating it if not yet done.
func (e *Envelope) Hash() common.Hash {
	if (e.hash == common.Hash{}) {
		encoded, _ := rlp.EncodeToBytes(e)
		e.hash = crypto.Keccak256Hash(encoded)
	}
	return e.hash
}

// rlpWithoutNonce returns the RLP encoded envelope contents, except the nonce.
func (e *Envelope) rlpWithoutNonce() []byte {
	res, _ := rlp.EncodeToBytes([]interface{}{e.Expiry, e.TTL, e.Data})
	return res
}

// DecodeRLP decodes an Envelope from an RLP data stream.
func (e *Envelope) DecodeRLP(s *rlp.Stream) error {
	raw, err := s.Raw()
	if err != nil {
		return err
	}
	// The decoding of Envelope uses the struct fields but also needs
	// to compute the hash of the whole RLP-encoded envelope. This
	// type has the same structure as Envelope but is not an
	// rlp.Decoder (does not implement DecodeRLP function).
	// Only public members will be encoded.
	type rlpenv Envelope
	if err := rlp.DecodeBytes(raw, (*rlpenv)(e)); err != nil {
		return err
	}
	e.hash = crypto.Keccak256Hash(raw)
	return nil
}

// GetEnvelope retrieves an envelope from the message queue by its hash.
// It returns nil if the envelope can not be found.
func (w *TomoX) GetEnvelope(hash common.Hash) *Envelope {
	w.poolMu.RLock()
	defer w.poolMu.RUnlock()
	return w.envelopes[hash]
}

// Seal closes the envelope by spending the requested amount of time as a proof
// of work on hashing the data.
func (e *Envelope) Seal(options *MessageParams) error {
	var target, bestBit int
	buf := make([]byte, 64)
	h := crypto.Keccak256(e.rlpWithoutNonce())
	copy(buf[:32], h)

	finish := time.Now().Add(time.Duration(options.WorkTime) * time.Second).UnixNano()
	for nonce := uint64(0); time.Now().UnixNano() < finish; {
		for i := 0; i < 1024; i++ {
			binary.BigEndian.PutUint64(buf[56:], nonce)
			d := new(big.Int).SetBytes(crypto.Keccak256(buf))
			firstBit := math.FirstBitSet(d)
			if firstBit > bestBit {
				e.Nonce, bestBit = nonce, firstBit
				if target > 0 && bestBit >= target {
					return nil
				}
			}
			nonce++
		}
	}

	return nil
}
