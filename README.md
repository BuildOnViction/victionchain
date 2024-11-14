# Viction Blockchain

Viction Blockchain (or Viction for short) is an innovative solution to the scalability problem with the Ethereum blockchain.
Our mission is to be a leading force in building the Internet of Value, and its infrastructure.
We are working to create an alternative, scalable financial system which is more secure, transparent, efficient, inclusive, and equitable for everyone.

Viction relies on a system of 150 Masternodes with a Proof of Stake Voting consensus that can support near-zero fee, and 2-second transaction confirmation times. Security, stability, and chain finality are guaranteed via novel techniques such as double validation, staking via smart-contracts, and "true" randomization processes.

Viction supports all EVM-compatible smart-contracts, protocols, and atomic cross-chain token transfers.
New scaling techniques such as sharding, private-chain generation, and hardware integration will be continuously researched and incorporated into Viction Blockchain's masternode architecture. This architecture will be an ideal scalable smart-contract public blockchain for decentralized apps, token issuances, and token integrations for small and big businesses.

More details can be found at our [white papers](https://docs.viction.xyz/whitepaper-and-research).

Read more about us on:

- our website: [https://viction.xyz](https://viction.xyz)
- our documentation portal: [https://docs.viction.xyz](https://docs.viction.xyz)
- our blockchain explorer: [https://vicscan.xyz](https://vicscan.xyz)

## How To Build

Viction supports both binaries build and Docker build for your convenience. Building Viction requires Go lang and C compiler on all platforms (Linux, Windows, MacOS).

### Binary Build

#### Install Go1.18.10

Due to some changes in Go-lang, Viction need to be build specificly with Go `1.18`. The following command will install Go 1.18.10 side-by-side with current system wide Go version:

```bash
go install golang.org/dl/go1.18.10@latest
go1.18.10 download
```

#### Build `tomo`

Clone this repository and change working directory to where you clone it, then run the following commands:

```bash
go1.18.10 run build/ci.go install
```

### Docker Build

Clone this repository and change working directory to where you clone it, then run the following commands:

```bash
docker build --file Dockerfile.node -t "viction:2.4.3" .
```

### Pre-built Bianries

Alternatively, you could quickly download our pre-complied binary from our [github release page](https://github.com/BuildOnViction/victionchain/releases)

## Running `tomo`

Going through all the possible command line flags is out of scope here, but we've enumerated a few common parameter combos to get you up to speed quickly on how you can run your own `tomo` instance.

### Initialize / Import accounts for the node keystore

Viction requires an account when running the node, even it's a full node. If you already had an existing account, import it. Otherwise, please initialize new accounts.

Initialize new account

```bash
tomo account new --password /path/to/password_file --keystore /path/to/keystore_dir
```

Import account

```bash
tomo account import /path/to/privatekey_file --password /path/to/password_file --keystore /path/to/keystore_dir
```

### Run a full node on mainnet

To run full node with default settings, simply run this command.

```bash
tomo --datadir /path/to/data_dir --keystore /path/to/keystore_dir --password /path/to/password_file --unlock 0
```

The following also run a full node on mainnet, with more customization.

```bash
tomo --datadir /path/to/data_dir \
  --keystore /path/to/keystore_dir \
  --password /path/to/password_file --unlock 0 \
  --identity my-full-node \
  --networkid 88 \
  --gasprice 250000000 \
  --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpcvhosts "*" --rpccorsdomain "*" \
  --rpcapi "eth,debug,net,db,personal,web3" \
  --ws --wsaddr 0.0.0.0 --wsport 8546 --wsorigins "*" \
  --mine \
  --bootnodes "enode://fd3da177f9492a39d1e7ce036b05745512894df251399cb3ec565081cb8c6dfa1092af8fac27991e66b6af47e9cb42e02420cc89f8549de0ce513ee25ebffc3a@3.212.20.0:30303,enode://97f0ca95a653e3c44d5df2674e19e9324ea4bf4d47a46b1d8560f3ed4ea328f725acec3fcfcb37eb11706cf07da669e9688b091f1543f89b2425700a68bc8876@104.248.98.78:30301,enode://b72927f349f3a27b789d0ca615ffe3526f361665b496c80e7cc19dace78bd94785fdadc270054ab727dbb172d9e3113694600dd31b2558dd77ad85a869032dea@188.166.207.189:30301,enode://c8f2f0643527d4efffb8cb10ef9b6da4310c5ac9f2e988a7f85363e81d42f1793f64a9aa127dbaff56b1e8011f90fe9ff57fa02a36f73220da5ff81d8b8df351@104.248.98.60:30301" \
  --port 30303 \
  --syncmode "full" --gcmode "full" \
  --ethstats my-full-node:getty-site-pablo-auger-room-sos-blair-shin-whiz-delhi@stats.viction.xyz \
  --verbosity 3
```

Brief explainations on the used flags:

```text
--datadir: path to your data directory created above.
--keystore: path to your account's keystore created above.
--password: your account's password.
--identity: your full node's name.
--networkid: our network ID.
--tomo-testnet: required when running a network other than Mainnet
--gasprice: Minimal gas price to accept for mining a transaction.
--rpc, --rpcaddr, --rpcport, --rpcvhosts, --rpccorsdomain: configure HTTP-RPC.
--ws, --wsaddr, --wsport, --wsorigins: configure Websocket.
--mine: your full-node wants to register to be a candidate for masternode selection.
--bootnodes: list of enodes of other peers that your full-node will try to connect at startup
--port: your full-node's listening port (default to 30303)
--nat NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>) to let other peer connect to your node easier
--synmode: blockchain sync mode ("fast", "full", or "light". More detail: https://github.com/BuildOnViction/victionchain/blob/master/eth/downloader/modes.go#L24)
--gcmode: blockchain garbage collection mode ("full", "archive")
--store-reward: store reward report. must be used in conjuction with --gcmode archive for archive node
--ethstats: send data to stats website
--verbosity: log level from 1 to 5. Here we're using 4 for debug messages
```

### Run Docker

```bash
docker run --name viction \
  -v "/path/to/data_dir:/tomochain/data" \
  -v "/path/to/keystore_dir:/tomochain/keystore" \
  -v "/path/to/password_file:/tomochain/password" \
  -p 8545:8545 \
  -p 8546:8546 \
  -p 30303:30303 \
  -e IDENTITY=my-full-node \
  -e NETWORK_ID=88
  -e BOOTNODES=enode://fd3da177f9492a39d1e7ce036b05745512894df251399cb3ec565081cb8c6dfa1092af8fac27991e66b6af47e9cb42e02420cc89f8549de0ce513ee25ebffc3a@3.212.20.0:30303,enode://97f0ca95a653e3c44d5df2674e19e9324ea4bf4d47a46b1d8560f3ed4ea328f725acec3fcfcb37eb11706cf07da669e9688b091f1543f89b2425700a68bc8876@104.248.98.78:30301,enode://b72927f349f3a27b789d0ca615ffe3526f361665b496c80e7cc19dace78bd94785fdadc270054ab727dbb172d9e3113694600dd31b2558dd77ad85a869032dea@188.166.207.189:30301,enode://c8f2f0643527d4efffb8cb10ef9b6da4310c5ac9f2e988a7f85363e81d42f1793f64a9aa127dbaff56b1e8011f90fe9ff57fa02a36f73220da5ff81d8b8df351@104.248.98.60:30301 \
  -e NETSTATS_HOST=stats.viction.xyz \
  -e NETSTATS_PORT=443 \
  -e WS_SECRET=getty-site-pablo-auger-room-sos-blair-shin-whiz-delhi \
  -e VERBOSITY=3 \
  buildonviction/node:2.4.3
```

Brief explainations on the supported variables:

```text
EXTIP: Your IP on the internet to let other peers connect to your nodes. Only use this if you have trouble connecting to peers.
IDENTITY: your full node's name.
PRIVATE_KEY: Private key of node's account in plain text.
PASSWORD: Password to encrypt/decrypt node's account in plain text.
NETWORK_ID: our network ID.
BOOTNODES: list of enodes of other peers that your full-node will try to connect at startup.
NETSTATS_HOST: Hostname of Ethstats service.
NETSTATS_PORT: Port of Ethstats service.
WS_SECRET: Secret of Ethstats service.
DEBUG_MODE: Enable archive mode.
VERBOSITY: log level from 1 to 5. Here we're using 4 for debug messages.
```

### Other usecases

For full featured guide. Please check our docs: [https://docs.viction.xyz/masternode](https://docs.viction.xyz/masternode)

## Contribution

Thank you for considering to try out our network and/or help out with the source code.
We would love to get your help; feel free to lend a hand.
Even the smallest bit of code, bug reporting, or just discussing ideas are highly appreciated.

If you would like to contribute to the tomochain source code, please refer to our Developer Guide for details on configuring development environment, managing dependencies, compiling, testing and submitting your code changes to our repo.

Please also make sure your contributions adhere to the base coding guidelines:

- Code must adhere to official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e uses [gofmt](https://golang.org/cmd/gofmt/)).
- Code comments must adhere to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
- Pull requests need to be based on and opened against the `master` branch.
- Any code you are trying to contribute must be well-explained as an issue on our [github issue page](https://github.com/BuildOnViction/victionchain/issues)
- Commit messages should be short but clear enough and should refer to the corresponding pre-logged issue mentioned above.
