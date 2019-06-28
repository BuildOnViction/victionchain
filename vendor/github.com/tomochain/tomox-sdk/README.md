# DEX backend

Official decentralized exchange backend, forked from the Proof project.  
The matching-engine will be soon moved to blockchain services.

_**Note: This project is under active development. We are using Github Issues to keep track of milestones, tasks, bugs, ....
It will be officially released when the version is v1.0.0**_

# Getting Started

### Requirements

- **mongoDB** version 3.6 or newer
- **rabbitmq** version 3.7.7 or newer
- **Go Modules** latest

### Deployment guide step by step
```
https://github.com/tomochain/dex-smart-contract/blob/develop/Deployment_step_by_step.md
```

## REST API
```
https://dex.devnet.tomochain.com/api/swaggerui/#
```

## Websocket API

See [WEBSOCKET_API.md](WEBSOCKET_API.md)


## Types

### Orders

Orders contain the information that is required to register an order in the orderbook as a "Maker".

- **id** is the primary ID of the order (possibly deprecated)
- **orderType** is either BUY or SELL. It is currently not parsed by the server and compute directly from buyToken, sellToken, buyAmount, sellAmount
- **exchangeAddress** is the exchange smart contract address
- **maker** is the maker (usually sender) ethereum account address
- **buyToken** is the BUY token ethereum address
- **sellToken** is the SELL token ethereum address
- **buyAmount** is the BUY amount (in BUY_TOKEN units)
- **sellAmount** is the SELL amount (in SELL_TOKEN units)
- **expires** is the order expiration timestamp
- **nonce** is the nonce that corresponds to
- **feeMake** is the maker fee (not implemented yet)
- **feeTake** is the taker fee (not implemented yet)
- **pairID** is a hash of the corresponding
- **hash** is a hash of the order details (see details below)
- **signature** is a signature of the order hash. The signer must equal to the maker address for the order to be valid.
- **price** corresponds to the pricepoint computed by the matching engine (not parsed)
- **amount** corresponds to the amount computed by the matching engine (not parsed)

**Order Price and Amount**

There are two ways to describe the amount of tokens being bought/sold. The smart-contract requires (buyToken, sellToken, buyAmount, sellAmount) while the
orderbook requires (pairID, amount, price).

The conversion between both systems can be found in the engine.ComputeOrderPrice
function

**Order Hash**

The order hash is a sha-256 hash of the following elements:

- Exchange address
- Token Buy address
- Amount Buy
- Token Sell Address
- Amount Sell
- Expires
- Nonce
- Maker Address

### Trades

When an order matches another order in the orderbook, the "taker" is required
to sign a trade object that matches an order.

- **orderHash** is the hash of the matching order
- **amount** is the amount of tokens that will be traded
- **trade nonce** is a unique integer to distinguish successive but identical orders (note: can probably be renamed to nonce)
- **taker** is the taker ethereum account address
- **pairID** is a hash identifying the token pair that will be traded
- **hash** is a unique identifier hash of the trade details (see details below)
- **signature** is a signature of the trade hash

Trade Hash:

The trade hash is a sha-256 hash of the following elements:

- Order Hash
- Amount
- Taker Address
- Trade Nonce

The (Order, Trade) tuple can then be used to perform an on-chain transaction for this trade.

### Quote Tokens and Token Pairs

In the same way as traditional exchanges function with the idea of base
currencies and quote currencies, the decentralized exchange works with
base tokens and quote tokens under the following principles:

- Only the exchange operator can register a quote token
- Anybody can register a token pair (but the quote token needs to be registered)

Token pairs are identified by an ID (a hash of both token addresses)
