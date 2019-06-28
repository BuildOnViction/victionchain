package tomox

func (tx *TxDataMatch) DecodeOrder() (*OrderItem, error) {
	order := &OrderItem{}
	if err := DecodeBytesItem(tx.order, order); err != nil {
		return order, err
	}
	return order, nil
}

func (tx *TxDataMatch) GetTrades() []map[string]string {
	return tx.trades
}