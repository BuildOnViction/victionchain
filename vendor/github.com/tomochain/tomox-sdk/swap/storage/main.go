package storage

import "github.com/tomochain/tomox-sdk/types"

// Storage is an interface that must be implemented by an object using
// persistent storage.
type Storage interface {
	// GetBitcoinBlockToProcess gets the number of Bitcoin block to process. `0` means the
	// processing should start from the current block.
	GetBlockToProcess(chain types.Chain) (uint64, error)
	// SaveLastProcessedBitcoinBlock should update the number of the last processed Bitcoin
	// block. It should only update the block if block > current block in atomic transaction.
	SaveLastProcessedBlock(chain types.Chain, block uint64) error
}
