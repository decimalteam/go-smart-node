#!/bin/bash

echo "Current node is a validator."
rm -r ~/.decimal/daemon/data/*.db
rm -r ~/.decimal/daemon/data/*.wal
rm -r ~/.decimal/daemon/config/gentx
rm ~/.decimal/daemon/config/write-file-atomic-*

echo "Wipe priv_validator_state."
PRIV_VAL_STATE_PATH=~/.decimal/daemon/data/priv_validator_state.json
PRIV_VAL_STATE_ZERO=$(jq <$PRIV_VAL_STATE_PATH '.height="0"|.step=0')
echo "$PRIV_VAL_STATE_ZERO" >$PRIV_VAL_STATE_PATH

echo "Restating daemon."
sudo systemctl restart dscd
