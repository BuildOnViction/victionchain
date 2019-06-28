# Websocket API

**Websocket Endpoint**: `/socket`

There are 8 channels on the matching engine websocket API:

- orders
- ohlcv
- orderbook
- raw_orderbook
- trades
- price_board
- markets
- notification

To send a message to a specific channel, the channel the general format of a message is the following:

```json
{
  "channel": <channel_name>,
  "event": {
    "type": <event_type>,
    "payload": <payload>
  }
}
```

where

- \<channel_name> is either 'orders', 'ohlcv', 'orderbook', 'raw_orderbook', 'trades'
- \<event_type> is a string describing what type of message is being sent
- \<payload> is a JSON object

# Trades Channel

## Message:

- SUBSCRIBE_TRADES (client --> server)
- UNSUBSCRIBE_TRADES (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE_TRADES MESSAGE (client --> server)

```json
{
  "channel": "trades",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": <address>,
      "quoteToken": <address>,
      "name": <baseTokenSymbol>/<quoteTokenSymbol>,
    }
  }
}
```

The name parameter is optional but the API will return an error if the symbols do not correspond to the symbols registered
in the matching engine. This optional parameter can be used for verifying you are subscribing to the right channel.

###Example:

```json
{
  "channel": "trades",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E",
      "name": "FUN/WETH"
    }
  }
}
```

## UNSUBSCRIBE_TRADES MESSAGE (client --> server)

```json
{
  "channel": "trades",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the init message is the following:

```json
{
  "channel": "trades",
  "event": {
    "type": "INIT",
    "payload": [
      <trade>,
      <trade>,
      ...
    ]
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "trades",
  "event": {
    "type": "UPDATE",
    "payload": [
      <trade>,
      <trade>,
      ...
    ]
  }
}
```

The format of the message is identical to the INIT message.
This message differs from the INIT message in that the client is supposed
to updated trades with identical hashes in their data set. The INIT messages
supposes that there are no other existing trades for the currently subscribe pair.

# Orderbook Channel

## Message:

- SUBSCRIBE_ORDERBOOK (client --> server)
- UNSUBSCRIBE_ORDERBOOK (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE_ORDERBOOK MESSAGE (client --> server)

```json
{
  "channel": "orderbook",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": <address>,
      "quoteToken": <address>,
      "name": <baseTokenSymbol>/<quoteTokenSymbol>,
    }
  }
}
```

### Example:

```json
{
  "channel": "orderbook",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E"
      "name": "FUN/WETH",
    }
  }
}
```

## UNSUBSCRIBE_ORDERBOOK MESSAGE (client --> server)

```json
{
  "channel": "orderbook",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the INIT message is the following:

```json
{
  "channel": "trades",
  "event": {
    "type": "INIT",
    "payload": {
      "asks": [ <ask>, <ask>, ... ],
      "bids": [ <bid>, <bid>, ... ],
    }
  }
}
```

# Example:

```json
{
  "channel": "orderbook",
  "event": {
    "type": "INIT",
    "payload": {
      "asks": [
        { "amount": "10000", "pricepoint": "1000000" },
        { "amount": "10000", "pricepoint": "1000000" }
      ],
      "bids": [
        { "amount": "10000", "pricepoint": "1000000" },
        { "amount": "10000", "pricepoint": "1000000" }
      ]
    }
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "trades",
  "event": {
    "type": "UPDATE",
    "payload": {
      "asks": [ <ask>, <ask>, ... ],
      "bids": [ <bid>, <bid>, ... ],
    },
  },
}
```

The format of the message is identical to the INIT message.
This message differs from the INIT message in that the client is supposed
to updated trades with identical hashes in their data set. The INIT messages
supposes that there are no other existing trades for the currently subscribe pair.

# Example:

```json
{
  "channel": "orderbook",
  "event": {
    "type": "UDPATE",
    "payload": {
      "asks": [
        { "amount": "10000", "pricepoint": "1000000" },
        { "amount": "10000", "pricepoint": "1000000" }
      ],
      "bids": [
        { "amount": "10000", "pricepoint": "1000000" },
        { "amount": "10000", "pricepoint": "1000000" }
      ]
    }
  }
}
```

# OHLCV Channel

## Message:

- SUBSCRIBE_OHLCV (client --> server)
- UNSUBSCRIBE_OHLCV (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE_OHLCV MESSAGE (client --> server)

```json
{
  "channel": "ohlcv",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": <baseTokenAddress>,
      "quoteToken": <quoteTokenAddress>,
      "name": <baseTokenSymbol/quoteTokenSymbol>,
      "from": <from>,
      "to": <to>,
      "duration": <duration>,
      "units": <hour>
    }
  }
}
```

where

- \<baseTokenAddress> is the Ethereum address of the base token contract,
- \<quoteTokenAddress> is the Ethereum address of the quote token contract,
- \<duration> is the duration (in units, see param below) of each candlestick
- \<units> is the unit used to represent the above duration: "minute", "hour", "day", "week", "month"
- \<from> is the beginning timestamp from which ohlcv data has to be queried
- \<to> is the ending timestamp until which ohlcv data has to be queried

### Example:

```json
{
  "channel": "ohlcv",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E",
      "name": "FUN/WETH",
      "from": 1534746133,
      "to": 1540016533,
      "duration": 1,
      "units": "hour"
    }
  }
}
```

## UNSUBSCRIBE_OHLCV MESSAGE (client --> server)

```json
{
  "channel": "ohlcv",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

# Orders Channel

## Message:

- GET_TOKENS (client --> server)
- NEW_ORDER (client --> server)
- ORDER_ADDED (server --> client)
- CANCEL_ORDER (client --> server)
- ORDER_CANCELLED (server --> client) #CANCELLED with two L
- REQUEST_SIGNATURE (server --> client)
- SUBMIT_SIGNATURE (client --> server)
- ORDER_PENDING (server --> client)
- ORDER_SUCCESS (server --> client)
- ORDER_ERROR (server --> client)
- ERROR (server --> client)
- UPDATE (server --> client)

## GET_TOKENS (client --> server)

The general format of the GET_TOKENS message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "GET_TOKENS"
  }
}
```

## NEW_ORDER (client --> server)

The general format of the NEW_ORDER message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "NEW_ORDER",
    "hash": <orderhash>,
    "payload": <order>,
  }
}
```

where:

- \<orderhash> is the hash of the order in the payload (probably will get deprecated)
- \<order> is an order signed by the client

## Example:

```json
{
  "channel": "orders",
  "event": {
    "type": "NEW_ORDER",
    "hash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12",
    "payload": {
      "exchangeAddress": "0x44e5a8cc74C389e805DAa84993bacC2b833E13f0",
      "userAddress": "0xCf7389Dc6c63637598402907d5431160eC8972A5",
      "buyToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "sellToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E",
      "buyAmount": "1000000000000000000",
      "sellAmount": "1000000000000000000",
      "hash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12",
      "makeFee": "0",
      "takeFee": "0",
      "expires": "10000000000000",
      "nonce": "3426223084368067",
      "signature": {
        "r": "0x14f9f97a43df77b79e79fa0af2ccfc8798be6ab2fa1f4468f92f99ed3542c1b9",
        "s": "0x78b035c950f5cb54a58ca54961f82a8c301027693f50272ea0f178b55891f8fc",
        "v": 28
      }
    }
  }
}
```

Note: Take note that most values are strings (except for the V value in the signature).
Using numbers or floats instead of strings will fail. This is required

## ORDER_ADDED MESSAGE (server --> client)

The general format of the ORDER_ADDED message is the following:

```
{
  "channel": "orders",
  "event": {
    "type": "ORDER_ADDED",
    "hash": <orderhash>,
    "payload": <order>
  }
}
```

where:

- \<orderhash> is the hash of the order in the payload (probably will get deprecated)
- \<order> is the original order sent by the client. Additional fields are added such as the `status`

```json
{
  "channel": "orders",
  "event": {
    "userAddress": "0xcf7389dc6c63637598402907d5431160ec8972a5",
    "exchangeAddress": "0x44e5a8cc74c389e805daa84993bacc2b833e13f0",
    "buyToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
    "sellToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
    "buyAmount": "1000000000000000000",
    "sellAmount": "1000000000000000000",
    "hash": "0xb958a32836f4ca15c93c0e54a22e83b384dc7ec899c6b66952195f12b0ed5708",
    "makeFee": "0",
    "takeFee": "0",
    "expires": "10000000000000",
    "nonce": "7740163281176971",
    "amount": "1000000000000000000",
    "filledAmount": "0",
    "pricepoint": "1000000",
    "side": "SELL",
    "status": "OPEN",
    "pairName": "FUN/WETH",
    "signature": {
      "r": "0x1106adc21f1d85e37d245df4808c683b2dad7970f6e831c96d2c91c886281a9d",
      "s": "0x43dbfbb4c2c7cfcb2afadc3abd63854449eaad9e9d306981fad9bcd46dce0d0b",
      "v": 27},
    "quoteToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
    "baseToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
    "createdAt": "2018-10-20T15:21:55.119253+09:00"
    }
  }
}
```

## CANCEL_ORDER MESSAGE (client --> server)

The general format of the CANCEL_ORDER message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "CANCEL_ORDER",
    "hash": <hash>, # will get deprecated i think
    "payload": {
      "hash": <hash>,
      "orderHash": <orderHash>,
      "signature": <signature>,
    }
  }
}
```

where:

- \<orderhash> is the hash of the order that needs to be canceled
- \<hash> is a hash of the orderHash
- \<signature> is a signature of the previous \<hash> by the private key that was used to sign \<orderHash>

Example:

```json
{
  "channel": "orders",
  "event": {
    "type": "CANCEL_ORDER",
    "hash": "0xd3cad812e8a15d0efedb11187d14e82f4ec190df455844583b8844dbc2e068b2",
    "payload": {
      "hash": "0xd3cad812e8a15d0efedb11187d14e82f4ec190df455844583b8844dbc2e068b2",
      "orderHash": "0xb958a32836f4ca15c93c0e54a22e83b384dc7ec899c6b66952195f12b0ed5708",
      "signature": {
        "r": "0x584391c40d49fedd7760ba2d1becc960ceb8a220f1c93156198fe6d18d31db02",
        "s": "0x34a5775f8a0ef32ea0fea8158975203c74b0646e1f37b27fc640ec902290ceca",
        "v": 27
      }
    }
  }
}
```

## ORDER_CANCELLED_MESSAGE (server --> client)

The general format of the order cancelled message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "ORDER_CANCELLED",
    "hash": "0xd3cad812e8a15d0efedb11187d14e82f4ec190df455844583b8844dbc2e068b2",
    "payload": <order>,
  }
}
```

## Example:

```json
{
  "channel": "orders",
  "event": {
    "type": "ORDER_CANCELLED",
    "hash": "0xd3cad812e8a15d0efedb11187d14e82f4ec190df455844583b8844dbc2e068b2",
    "payload": {
      "exchangeAddress": "0x44e5a8cc74c389e805daa84993bacc2b833e13f0",
      "userAddress": "0xcf7389dc6c63637598402907d5431160ec8972a5",
      "buyToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
      "sellToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
      "buyAmount": "1000000000000000000",
      "sellAmount": "1000000000000000000",
      "makeFee": "0",
      "takeFee": "0",
      "expires": "10000000000000",
      "nonce": "7740163281176971",
      "side": "SELL",
      "status": "CANCELLED",
      "amount": "1000000000000000000",
      "filledAmount": "0",
      "hash": "0xb958a32836f4ca15c93c0e54a22e83b384dc7ec899c6b66952195f12b0ed5708",
      "pricepoint": "1000000",
      "pairName": "FUN/WETH",
      "quoteToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
      "baseToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
      "signature": {,
        "r": "0x1106adc21f1d85e37d245df4808c683b2dad7970f6e831c96d2c91c886281a9d",
        "s": "0x43dbfbb4c2c7cfcb2afadc3abd63854449eaad9e9d306981fad9bcd46dce0d0b",
        "v": 27,
      },
    }
}
```

## REQUEST_SIGNATURE MESSAGE (server --> client)

The general format of the request signature message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "SUBMIT_SIGNATURE",
    "hash": <orderHash>,
    "payload": {
      "order": <order>,
      "remainingOrder": <remainingOrder>,
      "matches": [
        {
          "order": <order>,
          "trade": <trade>
        },
        {
          "order": <order>,
          "trade": <trade>
        },
        ...
      ]
    }
  }
}
```

where:

- \<orderhash> is the order hash of the initially sent order (also the hash of in the REQUEST_SIGNATURE message)

- \<order> is the signed order given in the REQUEST_SIGNATURE message. It is the order that was initially made by the user (with additional fields such as `createdAt` or `filledAmount` added by the matching engine)

- \<remainingOrder> is a signed order that represents a new order that has to be signed by the user in case his initial order was a "taker" order and was a partial match. Indeed, in order for an order to be accepted by the exchange smart contract, the order data needs to be signed. In case of a partial match, the remaining order comes with a new amount meaning the order has to be resigned by the sender.

- `matches` contains an array of { order, trade } objects that have been matched.

## SUBMIT_SIGNATURE MESSAGE (client --> server)

The general format of the submit signature messages is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "SUBMIT_SIGNATURE",
    "hash": <orderHash>,
    "payload": {
      "order": <order>,
      "remainingOrder": <remainingOrder>,
      "matches": [
        {
          "order": <order>,
          "trade": <trade>
        },
        {
          "order": <order>,
          "trade": <trade>
        },
        ...
      ]
    }
  }
}
```

where:

- \<orderhash> is the order hash of the initially sent order (also the hash of in the REQUEST_SIGNATURE message)

- \<order> is the signed order given in the REQUEST_SIGNATURE message. It is the order that was initially made by the user (with additional fields such as `createdAt` or `filledAmount` added by the matching engine).

- \<remainingOrder> is a signed order that represents a new order that has to be signed by the user in case his initial order was a "taker" order and was a partial match. Indeed, in order for an order to be accepted by the exchange smart contract, the order data needs to be signed. In case of a partial match, the remaining order comes with a new amount meaning the order has to be resigned by the sender. The remaining order needs to be signed by the client

- `matches` contains an array of { order, trade } objects that have been matched. Each trade needs to be signed by the original client that made the request.

### Example:

```json
{
  "channel": "orders",
  "event": {
    "type": "SUBMIT_SIGNATURE",
    "hash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12",
    "payload": {
      "order": {
        "exchangeAddress": "0x44e5a8cc74c389e805daa84993bacc2b833e13f0",
        "userAddress": "0xcf7389dc6c63637598402907d5431160ec8972a5",
        "buyToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
        "sellToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
        "buyAmount": "1000000000000000000",
        "sellAmount": "1000000000000000000",
        "hash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12",
        "makeFee": "0",
        "takeFee": "0",
        "expires": "10000000000000",
        "nonce": "3426223084368067",
        "createdAt": "2018-10-20T14:48:27.619218+09:00",
        "amount": "1000000000000000000",
        "filledAmount": "1000000000000000000",
        "side": "BUY",
        "status": "FILLED",
        "pairName": "FUN/WETH",
        "pricepoint": "1000000",
        "signature": {
          "r": "0x14f9f97a43df77b79e79fa0af2ccfc8798be6ab2fa1f4468f92f99ed3542c1b9",
          "s": "0x78b035c950f5cb54a58ca54961f82a8c301027693f50272ea0f178b55891f8fc",
          "v": 28
        },
        "baseToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
        "quoteToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e"
      },
      "matches": [
        {
          "order": {
            "exchangeAddress": "0x44e5a8cc74c389e805daa84993bacc2b833e13f0",
            "userAddress": "0xcf7389dc6c63637598402907d5431160ec8972a5",
            "buyToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e",
            "sellToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
            "buyAmount": "1000000000000000000",
            "sellAmount": "1111111111111111111",
            "hash": "0x4c6b45c7b8d9df7a2f2f07521b545b22d075a6eb28cc2a849bb441ae5b91fe6c",
            "makeFee": "0",
            "takeFee": "0",
            "expires": "10000000000000",
            "nonce": "6155556364519545",
            "createdAt": "2018-10-20T02:57:51.287+09:00",
            "amount": "1000000000000000000",
            "filledAmount": "1000000000000000000",
            "side": "SELL",
            "status": "FILLED",
            "pairName": "FUN/WETH",
            "pricepoint": "900000",
            "signature": {
              "r": "0x0796313b59449ca056244d63089a6c7043bfe7c1eebd69f9a6374ae9975b206e",
              "s": "0x587b5f42ec0007ae07851c24d489cb9d5024a02d0c91a6066d135f451fa074da",
              "v": 28
            },
            "baseToken": "0x546d3b3d69e30859f4f3ba15f81809a2efce6e67",
            "quoteToken": "0x17b4e8b709ca82abf89e172366b151c72df9c62e"
          },
          "trade": {
            "maker": "0xcf7389dc6c63637598402907d5431160ec8972a5",
            "taker": "0xcf7389dc6c63637598402907d5431160ec8972a5",
            "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
            "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E",
            "amount": "1000000000000000000",
            "pricepoint": "1000000",
            "createdAt": "0001-01-01T00:00:00Z",
            "hash": "0xc7506b0c3316305cbdcf9d2135610550eff1a2c010f33a3a20e84fe5535c17a9",
            "orderHash": "0x4c6b45c7b8d9df7a2f2f07521b545b22d075a6eb28cc2a849bb441ae5b91fe6c",
            "takerOrderHash": "0x70034a326ab444412ae57b84dd30b0566c3a86e3cba15717d319847429444d12",
            "pairName": "FUN/WETH",
            "side": "BUY",
            "tradeNonce": "5596387790547337",
            "signature": {
              "r": "0xcc881329d253de13a9a7b9c7c48d381954038e9899099d4453a858d7799db7e1",
              "s": "0x55a78923f7d6d7cf4d0e06e2bfbc511d87dbfdfa58cad51de8882b99e5da0976",
              "v": 27
            }
          }
        }
      ]
    }
  }
}
```

## ORDER PENDING MESSAGE (server --> client)

The general format of the submit signature messages is the following:

```
{
  "channel": "orders",
  "event": {
    "type": "ORDER_PENDING",
    "hash": <order hash>
    "payload": {
      "order": <order>,
      "matches": [
        {
          "order": <order>,
          "trade": <trade>
        },
        {
          "order": <order>,
          "trade": <trade>
        },
        ...
      ]
    }
  }
}
```

where

- \<order> is the order that was initially sent by the user/client
- `matches` is an array of all orders that are being matched

This message can be used by the client to get updated information on each order such as the order status or simply
to confirm that this order was indeed been sent to the transaction queue

## ORDER SUCCESS MESSAGE (server --> client)

The order success message indicates that the order was successful and correctly executed on the Ethereum chain.

The general format of the order success message is the following:

```json
{
  "channel": "orders",
  "event": {
    "type": "ORDER_SUCCESS",
    "hash": <order hash>
    "payload": {
      "order": <order>,
      "matches": [
        {
          "order": <order>,
          "trade": <trade>
        },
        {
          "order": <order>,
          "trade": <trade>
        },
        ...
      ]
    }
  }
}
```

It is identical to the order success message except that order statuses are different.
This means that the trade transaction was successful.

## ORDER ERROR MESSAGE (server --> client)

The ORDER_ERROR message indicates that a trade transaction was sent to the blockchain but was rejected.

```json
{
  "channel": "orders",
  "event": {
    "type": "ORDER_ERROR",
    "hash": <orderhash>,
    "payload": {
      "order": <order>,
      "matches": [
        {
          "order": <order>,
          "trade": <trade>
        },
        {
          "order": <order>,
          "trade": <trade>
        },
        ...
      ]
    }
  }
}
```

It is identical to the order successs message except that order statuses are different.
The client should usually not receive this message and it can be interpreted as an 'internal server error' (bug in the system rather than a malformed payload or client error)

# Raw Orderbook Channel

## Message:

- SUBSCRIBE (client --> server)
- UNSUBSCRIBE (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE_RAW_ORDERBOOK MESSAGE (client --> server)

```json
{
  "channel": "raw_orderbook",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": <address>,
      "quoteToken": <address>,
      "name": <baseTokenSymbol>/<quoteTokenSymbol>,
    }
  }
}
```

### Example:

```json
{
  "channel": "raw_orderbook",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E"
      "name": "FUN/WETH",
    }
  }
}
```

## UNSUBSCRIBE MESSAGE (client --> server)

```json
{
  "channel": "raw_orderbook",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the INIT message is the following:

```json
{
  "channel": "raw_orderbook",
  "event": {
    "type": "INIT",
    "payload": [
      <order>,
      <order>,
      <order>
    ]
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "raw_orderbook",
  "event": {
    "type": "INIT",
    "payload": [
      <order>,
      <order>,
      <order>
    ]
  }
}
```

# Price Board Channel

## Message:

- SUBSCRIBE (client --> server)
- UNSUBSCRIBE (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE_PRICE_BOARD MESSAGE (client --> server)

```json
{
  "channel": "price_board",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": <address>,
      "quoteToken": <address>,
      "name": <baseTokenSymbol>/<quoteTokenSymbol>
    }
  }
}
```

### Example:

```json
{
  "channel": "price_board",
  "event": {
    "type": "SUBSCRIBE",
    "payload": {
      "baseToken": "0x546d3B3d69E30859f4F3bA15F81809a2efCE6e67",
      "quoteToken": "0x17b4E8B709ca82ABF89E172366b151c72DF9C62E",
      "name": "ETH/TOMO"
    }
  }
}
```

## UNSUBSCRIBE MESSAGE (client --> server)

```json
{
  "channel": "price_board",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the INIT message is the following:

```json
{
  "channel": "price_board",
  "event": {
    "type": "INIT",
    "payload": {
      "ticks": [
        {
          "open": "",
          "close": "",
          "high": "",
          "low": "",
          "volume": "",
          "pair": "",
          "count": "",
          "timestamp": ""
        }
      ],
      "usd": "",
      "last_trade_price": ""
    }
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "price_board",
  "event": {
    "type": "UPDATE",
    "payload": {
      "ticks": [
        {
          "open": "",
          "close": "",
          "high": "",
          "low": "",
          "volume": "",
          "pair": "",
          "count": "",
          "timestamp": ""
        }
      ],
      "usd": "",
      "last_trade_price": ""
    }
  }
}
```

# Markets Channel

## Message:

- SUBSCRIBE (client --> server)
- UNSUBSCRIBE (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE MESSAGE (client --> server)

```json
{
  "channel": "markets",
  "event": {
    "type": "SUBSCRIBE"
  }
}
```

## UNSUBSCRIBE MESSAGE (client --> server)

```json
{
  "channel": "markets",
  "event": {
    "type": "UNSUBSCRIBE"
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the INIT message is the following:

```json
{
  "channel": "markets",
  "event": {
    "type": "INIT",
    "payload": {
      "pairData": [
        {
          "askPrice": "251540000000000000000000000000000000000",
          "bidPrice": "251530000000000000000000000000000000000",
          "close": "251530000000000000000000000000000000000",
          "count": "4",
          "high": "251530000000000000000000000000000000000",
          "low": "251530000000000000000000000000000000000",
          "open": "251530000000000000000000000000000000000",
          "orderCount": "51",
          "orderVolume": "275000000000000000000",
          "pair": {
            "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
            "pairName": "ETH/TOMO",
            "quoteToken": "0x0000000000000000000000000000000000000001"
          },
          "price": "251535000000000000000000000000000000000",
          "rank": 0,
          "timestamp": 0,
          "volume": "11000000000000000000"
        },
        {
          "askPrice": "0",
          "bidPrice": "0",
          "close": "0",
          "count": "0",
          "high": "0",
          "low": "0",
          "open": "0",
          "orderCount": "0",
          "orderVolume": "0",
          "pair": {
            "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
            "pairName": "ETH/BTC",
            "quoteToken": "0x4Ce98e51bB710a9af0d0f417e4505bEb4e3b4632"
          },
          "price": "0",
          "rank": 0,
          "timestamp": 0,
          "volume": "0"
        },
        ...
      ],
      "smallChartsData": {
        "ETH/TOMO": [
          {
            "close": "251530000000000000000000000000000000000",
            "pair": {
              "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
              "pairName": "ETH/TOMO",
              "quoteToken": "0x0000000000000000000000000000000000000001"
            },
            "timestamp": 1556262000000
          }
        ]
      }
    }
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "markets",
  "event": {
    "type": "UPDATE",
    "payload": {
      "pairData": [
        {
          "askPrice": "251540000000000000000000000000000000000",
          "bidPrice": "251530000000000000000000000000000000000",
          "close": "251530000000000000000000000000000000000",
          "count": "4",
          "high": "251530000000000000000000000000000000000",
          "low": "251530000000000000000000000000000000000",
          "open": "251530000000000000000000000000000000000",
          "orderCount": "51",
          "orderVolume": "275000000000000000000",
          "pair": {
            "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
            "pairName": "ETH/TOMO",
            "quoteToken": "0x0000000000000000000000000000000000000001"
          },
          "price": "251535000000000000000000000000000000000",
          "rank": 0,
          "timestamp": 0,
          "volume": "11000000000000000000"
        },
        {
          "askPrice": "0",
          "bidPrice": "0",
          "close": "0",
          "count": "0",
          "high": "0",
          "low": "0",
          "open": "0",
          "orderCount": "0",
          "orderVolume": "0",
          "pair": {
            "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
            "pairName": "ETH/BTC",
            "quoteToken": "0x4Ce98e51bB710a9af0d0f417e4505bEb4e3b4632"
          },
          "price": "0",
          "rank": 0,
          "timestamp": 0,
          "volume": "0"
        },
        ...
      ],
      "smallChartsData": {
        "ETH/TOMO": [
          {
            "close": "251530000000000000000000000000000000000",
            "pair": {
              "baseToken": "0x260800BAb2E7a6C3BCB8501BeE79F8CFFE770d17",
              "pairName": "ETH/TOMO",
              "quoteToken": "0x0000000000000000000000000000000000000001"
            },
            "timestamp": 1556262000000
          }
        ]
      }
    }
  }
}
```

# Notification Channel

## Message:

- SUBSCRIBE (client --> server)
- INIT (server --> client)
- UPDATE (server --> client)

## SUBSCRIBE MESSAGE (client --> server)

```json
{
  "channel": "notification",
  "event": {
    "type": "SUBSCRIBE",
    "payload": "0x..." // User address
  }
}
```

## INIT MESSAGE (server --> client)

The general format of the INIT message is the following:

```json
{
  "channel": "notification",
  "event": {
    "type": "INIT",
    "payload": [
      {
          "_id" : "5cd145019eef1c0f04b280d8",
          "recipient" : "0xF069080F7acB9a6705b4a51F84d9aDc67b921bDF",
          "message" : "ORDER_ADDED - Order Hash: 0xc5fef1df4ee7a1a0aaf29aefbc022b0e4a32ada674f6bce18f47e59aa72e31f8",
          "type" : "LOG",
          "status" : "UNREAD",
          "createdAt" : "2019-05-07T08:42:41.961Z",
          "updatedAt" : "2019-05-07T08:42:41.961Z"
      },
      {
          "_id" : "5cd146769eef1c0f04b280d9",
          "recipient" : "0xF069080F7acB9a6705b4a51F84d9aDc67b921bDF",
          "message" : "ORDER_CANCELLED - Order Hash: 0xc5fef1df4ee7a1a0aaf29aefbc022b0e4a32ada674f6bce18f47e59aa72e31f8",
          "type" : "LOG",
          "status" : "UNREAD",
          "createdAt" : "2019-05-07T08:48:54.630Z",
          "updatedAt" : "2019-05-07T08:48:54.630Z"
      }
    ]
  }
}
```

## UPDATE MESSAGE (server --> client)

The general format of the update message is the following:

```json
{
  "channel": "notification",
  "event": {
    "type": "UPDATE",
    "payload": [
      {
          "_id" : "5cd145019eef1c0f04b280d8",
          "recipient" : "0xF069080F7acB9a6705b4a51F84d9aDc67b921bDF",
          "message" : "ORDER_ADDED - Order Hash: 0xc5fef1df4ee7a1a0aaf29aefbc022b0e4a32ada674f6bce18f47e59aa72e31f8",
          "type" : "LOG",
          "status" : "UNREAD",
          "createdAt" : "2019-05-07T08:42:41.961Z",
          "updatedAt" : "2019-05-07T08:42:41.961Z"
      },
      {
          "_id" : "5cd146769eef1c0f04b280d9",
          "recipient" : "0xF069080F7acB9a6705b4a51F84d9aDc67b921bDF",
          "message" : "ORDER_CANCELLED - Order Hash: 0xc5fef1df4ee7a1a0aaf29aefbc022b0e4a32ada674f6bce18f47e59aa72e31f8",
          "type" : "LOG",
          "status" : "UNREAD",
          "createdAt" : "2019-05-07T08:48:54.630Z",
          "updatedAt" : "2019-05-07T08:48:54.630Z"
      }
    ]
  }
}
```