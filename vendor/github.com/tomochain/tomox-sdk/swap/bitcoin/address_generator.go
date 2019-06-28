package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tyler-smith/go-bip32"
)

// NewAddressGenerator return AddressGenerator from the master public key
func NewAddressGenerator(masterPublicKeyString string, chainParams *chaincfg.Params) (*AddressGenerator, error) {
	deserializedMasterPublicKey, err := bip32.B58Deserialize(masterPublicKeyString)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing master public key")
	}

	if deserializedMasterPublicKey.IsPrivate {
		return nil, errors.New("Key is not a master public key")
	}

	return &AddressGenerator{deserializedMasterPublicKey, chainParams}, nil
}

// Generate return public key corresponding to the index
func (g *AddressGenerator) Generate(index uint32) (string, error) {
	if g.masterPublicKey == nil {
		return "", errors.New("No master public key set")
	}

	accountKey, err := g.masterPublicKey.NewChildKey(index)
	if err != nil {
		return "", errors.Wrap(err, "Error creating new child key")
	}

	address, err := btcutil.NewAddressPubKey(accountKey.Key, g.chainParams)
	if err != nil {
		return "", errors.Wrap(err, "Error creating address for new child key")
	}

	return address.AddressPubKeyHash().EncodeAddress(), nil
}
