package bitcoin

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/swap/storage"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tyler-smith/go-bip32"
)

var logger = utils.Logger

var (
	eight = big.NewInt(8)
	ten   = big.NewInt(10)
	// satInBtc = 10^8
	satInBtc = new(big.Rat).SetInt(new(big.Int).Exp(ten, eight, nil))
)

// Listener listens for transactions using bitcoin-core RPC. It calls TransactionHandler for each new
// transactions. It will reprocess the block if TransactionHandler returns error. It will
// start from the block number returned from Storage.GetBlockToProcess("bitcoin") or the latest block
// if it returned 0. Transactions can be processed more than once, it's TransactionHandler
// responsibility to ignore duplicates.
// Listener tracks only P2PKH payments.
// You can run multiple Listeners if Storage is implemented correctly.
type Listener struct {
	Enabled              bool
	Client               Client          `inject:""`
	Storage              storage.Storage `inject:""`
	TransactionHandler   TransactionHandler
	Testnet              bool
	ConfirmedBlockNumber uint64
	chainParams          *chaincfg.Params
}

type Client interface {
	GetBlockCount() (int64, error)
	GetBlockHash(blockHeight int64) (*chainhash.Hash, error)
	GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error)
}

type TransactionHandler func(transaction Transaction) error

type Transaction struct {
	Hash       string
	TxOutIndex int
	// Value in sats
	ValueSat int64
	To       string
}

type AddressGenerator struct {
	masterPublicKey *bip32.Key
	chainParams     *chaincfg.Params
}

func BtcToSat(btc string) (int64, error) {

	valueRat := new(big.Rat)
	_, ok := valueRat.SetString(btc)
	if !ok {
		return 0, errors.New("Could not convert to *big.Rat")
	}

	// Calculate value in satoshi
	valueRat.Mul(valueRat, satInBtc)

	// Ensure denominator is equal `1`
	if valueRat.Denom().Cmp(big.NewInt(1)) != 0 {
		return 0, errors.New("Invalid precision, is value smaller than 1 satoshi?")
	}

	return valueRat.Num().Int64(), nil
}
