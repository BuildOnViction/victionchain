package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/tomochain/tomox-sdk/errors"

	"github.com/tomochain/tomox-sdk/types"
)

// by default the exchange type is fanout, which means it broadcast all the queues
// topic is just a list of words, direct is not filterable

// QueuePoolDepositTransactions : return a queue as a channel
func (c *Connection) QueuePoolDepositTransactions() (<-chan *types.DepositTransaction, error) {
	ch := c.GetChannel("depositSubscribe")
	q := c.GetQueue(ch, "deposit")

	// get channel
	msgs, err := c.Consume(ch, q)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	out := make(chan *types.DepositTransaction)

	go func() {
		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &types.DepositTransaction{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logger.Error(err)
					continue
				}
				// feed this queue forever
				out <- msg
			}
		}()

		// wait for background feeding
		<-forever
		// if this program is terminate
		close(out)
	}()

	return out, nil
}

func (c *Connection) PublishDepositTransaction(transaction *types.DepositTransaction) error {
	ch := c.GetChannel("depositPublish")
	q := c.GetQueue(ch, "deposit")

	bytes, err := json.Marshal(transaction)
	if err != nil {
		log.Fatal("Failed to marshal deposit: ", err)
		return errors.New("Failed to marshal deposit: " + err.Error())
	}

	err = c.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
