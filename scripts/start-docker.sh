#!/bin/bash

KEY="royalkey"
CHAINID="decimal-localnet-1"
MONIKER="royalmoniker"
DATA_DIR=$(mktemp -d -t decimal-datadir.XXXXX)

echo "create and add new keys"
./dscd keys add $KEY --home "$DATA_DIR" --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init dsc with moniker=$MONIKER and chain-id=$CHAINID"
./dscd init $MONIKER --chain-id $CHAINID --home "$DATA_DIR"
echo "prepare genesis: Allocate genesis accounts"
./dscd add-genesis-account \
"$(./dscd keys show $KEY -a --home "$DATA_DIR" --keyring-backend test)" 1000000000000000000del,1000000000000000000stake \
--home "$DATA_DIR" --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./dscd gentx $KEY 1000000000000000000stake --keyring-backend test --home "$DATA_DIR" --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./dscd collect-gentxs --home "$DATA_DIR"
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./dscd validate-genesis --home "$DATA_DIR"

echo "starting dsc node $i in background ..."
./dscd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home "$DATA_DIR" \
>"$DATA_DIR"/node.log 2>&1 & disown

echo "started dsc node"
tail -f /dev/null