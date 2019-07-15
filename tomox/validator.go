package tomox

import "github.com/ethereum/go-ethereum/log"

func (tx *TxDataMatch) DecodeOrder() (*OrderItem, error) {
	order := &OrderItem{}
	log.Debug("tx.order", "tx.order", tx.Order)
	if err := DecodeBytesItem(tx.Order, order); err != nil {
		return order, err
	}
	return order, nil
}

func (tx *TxDataMatch) GetTrades() []map[string]string {
	return tx.Trades
}
