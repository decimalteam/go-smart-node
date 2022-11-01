#!/bin/bash

KEY="royalkey"
FAUCET_KEY="faucet"
CHAINID="decimal_202020-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# Decimal paths
DECIMAL_CONFIG="$HOME/.decimal/daemon/config/config.toml"
DECIMAL_GENESIS="$HOME/.decimal/daemon/config/genesis.json"
DECIMAL_GENESIS_TMP="$HOME/.decimal/daemon/config/tmp_genesis.json"

# Calidate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Reinstall daemon
rm -rf "$HOME/.decimal/"
make install

# Set client config
dscd config keyring-backend $KEYRING
dscd config chain-id $CHAINID

echo "Keys and mnemonics" > keys-and-mnemonics.txt
echo "Validator:" >> keys-and-mnemonics.txt
# if $KEY exists it should be deleted
dscd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO >> keys-and-mnemonics.txt 2>&1

echo "--------------" >> keys-and-mnemonics.txt
echo "" >> keys-and-mnemonics.txt
echo "Faucet:" >> keys-and-mnemonics.txt
dscd keys add $FAUCET_KEY --keyring-backend $KEYRING --algo $KEYALGO >> keys-and-mnemonics.txt 2>&1

# Set moniker and chain-id for DSC (Moniker can be anything, chain-id must be an integer)
dscd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to del
jq <"$DECIMAL_GENESIS" '.app_state["validator"]["params"]["base_denom"]="del"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"
jq <"$DECIMAL_GENESIS" '.app_state["crisis"]["constant_fee"]["denom"]="del"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"
jq <"$DECIMAL_GENESIS" '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="del"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"
jq <"$DECIMAL_GENESIS" '.app_state["evm"]["params"]["evm_denom"]="del"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"
jq <"$DECIMAL_GENESIS" '.app_state["inflation"]["params"]["mint_denom"]="del"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"

# Set gas limit in genesis
jq <"$DECIMAL_GENESIS" '.consensus_params["block"]["max_gas"]="10000000"' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"

# Set claims start time
node_address=$(dscd keys show $KEY | grep  "address: " | cut -c12-)
current_date=$(date -u +"%Y-%m-%dT%TZ")

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = false/create_empty_blocks = true/g' "$DECIMAL_CONFIG"
  else
    sed -i 's/create_empty_blocks = false/create_empty_blocks = true/g' "$DECIMAL_CONFIG"
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' "$DECIMAL_CONFIG"
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' "$DECIMAL_CONFIG"
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' "$DECIMAL_CONFIG"
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' "$DECIMAL_CONFIG"
  fi
fi

# Allocate genesis accounts (cosmos formatted addresses)
dscd add-genesis-account $KEY 100000000000000000000000000del --keyring-backend $KEYRING
dscd add-genesis-account $FAUCET_KEY 1000000000000000000000000000000000000del --keyring-backend $KEYRING

dscd selfdelegation 100000000000000000000000del --keyring-backend $KEYRING --from $KEY

jq < "$DECIMAL_GENESIS" 'del(.app_state["bank"]["supply"])' > "$DECIMAL_GENESIS_TMP" && mv "$DECIMAL_GENESIS_TMP" "$DECIMAL_GENESIS"

# Run this to ensure everything worked and that the genesis file is setup correctly
echo "### Validate genesis"
dscd validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
#dscd start "$TRACE" --pruning "nothing" --log_level "$LOGLEVEL" --minimum-gas-prices "0.0001del" --json-rpc.api "eth,txpool,personal,net,debug,web3"
