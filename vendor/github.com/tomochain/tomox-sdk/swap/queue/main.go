package queue

import "github.com/tomochain/tomox-sdk/types"

// Queue implements transactions queue.
// The queue must not allow duplicates (including history) or must implement deduplication
// interval so it should not allow duplicate entries for 5 minutes since the first
// entry with the same ID was added.
// This is a critical requirement! Otherwise ETH/BTC may be sent twice to Tomochain account.
// If you don't know what to do, use default AWS SQS FIFO queue or DB queue.
type Queue interface {
	// QueueAdd inserts the element to this queue. If element already exists in a queue, it should
	// return nil.
	QueueAdd(tx *types.DepositTransaction) error
	// QueuePool return the queue as a read-only channel.
	QueuePool() (<-chan *types.DepositTransaction, error)
}
