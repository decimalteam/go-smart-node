#!/bin/bash

KEY="royalkey"
CHAINID="decimal_202020-1"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# Decimal paths
DECIMAL_CONFIG="$HOME/.decimal/config/config.toml"
DECIMAL_GENESIS="$HOME/.decimal/config/genesis.json"
DECIMAL_GENESIS_TMP="$HOME/.decimal/config/tmp_genesis.json"

# Calidate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# Reinstall daemon
rm -rf ~/.decimal*
make install

# Set client config
dscd config keyring-backend $KEYRING
dscd config chain-id $CHAINID

# if $KEY exists it should be deleted
dscd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO

# Set moniker and chain-id for DSC (Moniker can be anything, chain-id must be an integer)
dscd init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to del
cat $DECIMAL_GENESIS | jq '.app_state["staking"]["params"]["bond_denom"]="del"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS
cat $DECIMAL_GENESIS | jq '.app_state["crisis"]["constant_fee"]["denom"]="del"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS
cat $DECIMAL_GENESIS | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="del"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS
cat $DECIMAL_GENESIS | jq '.app_state["evm"]["params"]["evm_denom"]="del"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS
cat $DECIMAL_GENESIS | jq '.app_state["inflation"]["params"]["mint_denom"]="del"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Set gas limit in genesis
cat $DECIMAL_GENESIS | jq '.consensus_params["block"]["max_gas"]="10000000"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Set claims start time
node_address=$(dscd keys list | grep  "address: " | cut -c12-)
current_date=$(date -u +"%Y-%m-%dT%TZ")
cat $DECIMAL_GENESIS | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["airdrop_start_time"]=$current_date' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Set claims records for validator account
amount_to_claim=10000
cat $DECIMAL_GENESIS | jq -r --arg node_address "$node_address" --arg amount_to_claim "$amount_to_claim" '.app_state["claims"]["claims_records"]=[{"initial_claimable_amount":$amount_to_claim, "actions_completed":[false, false, false, false],"address":$node_address}]' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Set claims decay
cat $DECIMAL_GENESIS | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["duration_of_decay"]="1000000s"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS
cat $DECIMAL_GENESIS | jq -r --arg current_date "$current_date" '.app_state["claims"]["params"]["duration_until_decay"]="100000s"' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Claim module account:
# 0xA61808Fe40fEb8B3433778BBC2ecECCAA47c8c47 || dx15cvq3ljql6utxseh0zau9m8ve2j8erz8prfj7l
cat $DECIMAL_GENESIS | jq -r --arg amount_to_claim "$amount_to_claim" '.app_state["bank"]["balances"] += [{"address":"dx15cvq3ljql6utxseh0zau9m8ve2j8erz8prfj7l","coins":[{"denom":"del", "amount":$amount_to_claim}]}]' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# disable produce empty block
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $DECIMAL_CONFIG
  else
    sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $DECIMAL_CONFIG
fi

if [[ $1 == "pending" ]]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $DECIMAL_CONFIG
      sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $DECIMAL_CONFIG
  else
      sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $DECIMAL_CONFIG
      sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $DECIMAL_CONFIG
  fi
fi

# Allocate genesis accounts (cosmos formatted addresses)
dscd add-genesis-account $KEY 100000000000000000000000000del --keyring-backend $KEYRING

# Update total supply with claim values
# validators_supply=$(cat $DECIMAL_GENESIS | jq -r '.app_state["bank"]["supply"][0]["amount"]')
# Bc is required to add this big numbers
# total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
total_supply=100000000000000000000010000
cat $DECIMAL_GENESIS | jq -r --arg total_supply "$total_supply" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' > $DECIMAL_GENESIS_TMP && mv $DECIMAL_GENESIS_TMP $DECIMAL_GENESIS

# Sign genesis transaction
dscd gentx $KEY 1000000000000000000000del --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
dscd collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
dscd validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
dscd start --pruning=nothing $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001del --json-rpc.api eth,txpool,personal,net,debug,web3
