## Prerequisites
0.1. Configure Golang development environment
```
https://golang.org/doc/install
```

0.2. Setup Node.js
```
https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-18-04#installing-using-nvm
```

0.3. Install `yarn`
```
npm i -g yarn
```

0.4. Setup Docker
```
https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-18-04
```
----------------

## tomox-launch-kit
1. Clone it: 
```
git clone git@github.com:tomochain/tomox-launch-kit.git
```
2. Checkout branch `features/tomox-integration`

3. Run `yarn install`

4. Update `deploy/.env`:
```
COINBASE_ADDRESS=0xF9D87abd60435b70415CcC1FAAcA4F8B91786eDb
DB_NAME=tomodex
MONGODB_URL=mongodb://localhost:27017
NETWORK_ID=89
REGISTER_PRIVATE_KEY=463D27C152040C4E49C5D9606BF3A27E7CE00ACBA25FF4E6A42DD486C27443DA
RELAYER_REGISTRATION_CONTRACT_ADDRESS=0x6214de5b30c872e09db48e88798476ecce8c8da2
RPC_URL=https://testnet.tomochain.com

```

5. Run `yarn reset-env`

5.5. Wait for about 30 seconds after finishing above command

6. Run `yarn seeds`

----------------
## tomox-sdk
1. Clone it:
```
git clone git@github.com:tomochain/tomox-sdk.git
```
2.  Checkout `features/tomox-integration` branch

3. Update `config/config.yaml` with the URL of TomoX:
```
  http_url: http://localhost:8501
  ws_url: ws://localhost:9501
```

4. Start the server
```
yarn start
```

----------------
## tomox-sdk-ui
1. Clone it:
```
git clone git@github.com:tomochain/tomox-sdk-ui.git
```

2.  Checkout `develop` branch

3. Install dependencies
```
yarn install
```
4. Install `sass`:
```
https://sass-lang.com/install
```
5. Run `yarn query-tokens`

6. Start the development server
```
yarn start
```
This command will also compile sass files

----------------
## DONE
