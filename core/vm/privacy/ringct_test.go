package privacy

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/tomochain/tomochain/crypto"
)

//function returns(mutiple rings, private keys, message, error)
func GenerateMultiRingParams(numRing int, ringSize int, s int) (rings []Ring, privkeys []*ecdsa.PrivateKey, m [32]byte, err error) {
	for i := 0; i < numRing; i++ {
		privkey, err := crypto.GenerateKey()
		if err != nil {
			return nil, nil, [32]byte{}, err
		}
		privkeys = append(privkeys, privkey)

		ring, err := GenNewKeyRing(ringSize, privkey, s)
		if err != nil {
			return nil, nil, [32]byte{}, err
		}
		rings = append(rings, ring)
	}

	_, err = rand.Read(m[:])
	if err != nil {
		return nil, nil, [32]byte{}, err
	}
	return rings, privkeys, m, nil
}

func TestSign(t *testing.T) {
	/*for i := 14; i < 15; i++ {
	for j := 14; j < 15; j++ {
		for k := 0; k <= j; k++ {*/
	numRing := 5
	ringSize := 10
	s := 9
	fmt.Println("Generate random ring parameter ")
	rings, privkeys, m, err := GenerateMultiRingParams(numRing, ringSize, s)

	fmt.Println("numRing  ", numRing)
	fmt.Println("ringSize  ", ringSize)
	fmt.Println("index of real one  ", s)

	fmt.Println("Ring  ", rings)
	fmt.Println("privkeys  ", privkeys)
	fmt.Println("m  ", m)

	ringSignature, err := Sign(m, rings, privkeys, s)
	if err != nil {
		t.Error("Failed to create Ring signature")
	}

	sig, err := ringSignature.Serialize()
	if err != nil {
		t.Error("Failed to Serialize input Ring signature")
	}

	deserializedSig, err := Deserialize(sig)
	if err != nil {
		t.Error("Failed to Deserialize Ring signature")
	}
	verified := Verify(deserializedSig)

	if !verified {
		t.Error("Failed to verify Ring signature")
	}

}
