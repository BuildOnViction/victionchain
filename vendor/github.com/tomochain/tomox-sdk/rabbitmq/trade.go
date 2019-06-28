package rabbitmq

import (
	"encoding/json"

	"github.com/tomochain/tomox-sdk/types"
)

func (c *Connection) SubscribeTrades(fn func(*types.OperatorMessage) error) error {
	ch := c.GetChannel("tradeSubscribe")
	q := c.GetQueue(ch, "trades")

	go func() {
		msgs, err := c.Consume(ch, q)
		if err != nil {
			logger.Error(err)
		}

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				msg := &types.OperatorMessage{}
				err := json.Unmarshal(d.Body, msg)
				if err != nil {
					logger.Error(err)
					continue
				}

				go fn(msg)
			}
		}()

		<-forever
	}()
	return nil
}

func (c *Connection) PublishTrades(matches *types.Matches) error {
	ch := c.GetChannel("tradePublish")
	q := c.GetQueue(ch, "trades")

	msg := &types.OperatorMessage{
		MessageType: "NEW_TRADE",
		Matches:     matches,
	}

	logger.Info("operator/", msg.String())

	b, err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = c.Publish(ch, q, b)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
