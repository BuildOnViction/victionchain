#!/bin/sh

# constants
DATA_DIR="data"
KEYSTORE_DIR="keystore"

# variables
genesisPath=""
params=""
accountsCount=$(
  tomo account list --datadir ${DATA_DIR}  --keystore ${KEYSTORE_DIR} \
  2> /dev/null \
  | wc -l
)

# file to env
for env in IDENTITY NETWORK_ID SYNC_MODE \
           BOOTNODES EXTIP P2P_PORT NETSTATS_HOST NETSTATS_PORT WS_SECRET \
           PASSWORD PRIVATE_KEY ANNOUNCE_TXS VERBOSITY ''; do
  file=$(eval echo "\$${env}_FILE")
  if [[ -f $file ]] && [[ ! -z $file ]]; then
    echo "Replacing $env by $file"
    export $env=$(cat $file)
  elif [[ "$env" == "BOOTNODES" ]] && [[ ! -z $file ]]; then
    echo "Bootnodes file is not available. Waiting for it to be provisioned..."
    while true ; do
      if [[ -f $file ]] && [[ $(grep -e enode $file) ]]; then
        echo "Fount bootnode file."
        break
      fi
      echo "Still no bootnodes file, sleeping..."
      sleep 5
    done
    export $env=$(cat $file)
  fi
done

# networkid
if [[ ! -z ${NETWORK_ID} ]]; then
  case ${NETWORK_ID} in
    88 )
      genesisPath="mainnet.json"
      ;;
    89 )
      genesisPath="testnet.json"
      params="$params --tomo-testnet"
      ;;
    * )
      params="$params --tomo-testnet"
      ;;
  esac
fi

# custom genesis path
if [[ ! -z ${GENESIS_PATH} ]]; then
  genesisPath="${GENESIS_PATH}"
fi

# data dir
if [[ ! -d ${DATA_DIR}/tomo ]]; then
  echo "No blockchain data, creating genesis block."
  tomo init $genesisPath --datadir ${DATA_DIR} 2> /dev/null
fi

# identity
if [[ -z ${IDENTITY} ]]; then
  IDENTITY="unnamed_$(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c6)"
fi

# bootnodes
if [[ ! -z ${BOOTNODES} ]]; then
  params="$params --bootnodes ${BOOTNODES}"
fi

# extip
if [[ ! -z ${EXTIP} ]]; then
  params="$params --nat extip:${EXTIP}"
fi

# netstats
if [[ ! -z ${WS_SECRET} ]]; then
  echo "Will report to netstats server ${NETSTATS_HOST}:${NETSTATS_PORT}"
  params="$params --ethstats ${IDENTITY}:${WS_SECRET}@${NETSTATS_HOST}:${NETSTATS_PORT}"
else
  echo "WS_SECRET not set, will not report to netstats server."
fi

# password file
if [[ ! -f ./password ]]; then
  if [[ ! -z ${PASSWORD} ]]; then
    echo "Password env is set. Writing into file."
    echo "${PASSWORD}" > ./password
  else
    echo "No password set (or empty), generating a new one"
    $(< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c${1:-32} > password)
  fi
fi

# private key
if [[ $accountsCount -le 0 ]]; then
  echo "No accounts found"
  if [[ ! -z ${PRIVATE_KEY} ]]; then
    echo "Creating account from private key"
    echo "${PRIVATE_KEY}" > ./private_key
    tomo  account import ./private_key \
      --datadir ${DATA_DIR} \
      --keystore ${KEYSTORE_DIR} \
      --password ./password
    rm ./private_key
  else
    echo "Creating new account"
    tomo account new \
      --datadir ${DATA_DIR} \
      --keystore ${KEYSTORE_DIR} \
      --password ./password
  fi
fi
account=$(
  tomo account list --datadir ${DATA_DIR}  --keystore ${KEYSTORE_DIR} \
  2> /dev/null \
  | head -n 1 \
  | cut -d"{" -f 2 | cut -d"}" -f 1
)
echo "Using account $account"
params="$params --unlock $account"

# annonce txs
if [[ ! -z ${ANNOUNCE_TXS} ]]; then
  params="$params --announce-txs"
fi

# debug mode
if [[ ! -z ${DEBUG_MODE} ]]; then
  params="$params --gcmode archive --rpcapi db,eth,net,web3,debug,posv"
fi

# store reward
if [[ ! -z ${STORE_REWARD} ]]; then
  params="$params --store-reward"
fi

# dump
echo "dump: ${IDENTITY} $account ${BOOTNODES}"

set -x

exec tomo --identity ${IDENTITY} \
  --datadir ${DATA_DIR} \
  --networkid ${NETWORK_ID} \
  --syncmode ${SYNC_MODE} \
  --rpc \
  --rpccorsdomain "*" \
  --rpcaddr 0.0.0.0 \
  --rpcport 8545 \
  --rpcvhosts "*" \
  --ws \
  --wsaddr 0.0.0.0 \
  --wsport 8546 \
  --wsorigins "*" \
  --port ${P2P_PORT} \
  --maxpeers ${MAX_PEERS} \
  --txpool.globalqueue 5000 \
  --txpool.globalslots 5000 \
  --mine \
  --keystore ${KEYSTORE_DIR} \
  --password ./password \
  --gasprice "250000000" \
  --targetgaslimit "84000000" \
  --verbosity ${VERBOSITY} \
  ${params} \
  "$@"
