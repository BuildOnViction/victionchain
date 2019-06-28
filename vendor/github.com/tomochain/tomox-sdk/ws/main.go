package ws

func isClientConnected(clients []*Client, client *Client) bool {
	for _, c := range clients {
		if c == client {
			logger.Info("Client is connected")
			return true
		}
	}

	logger.Info("Client is not connected")
	return false
}
