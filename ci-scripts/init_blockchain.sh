#!/bin/bash

DECIMAL_VALIDATOR_01_TX_PATH="./build/gentx/validator-01.tmp.json"
DECIMAL_VALIDATOR_02_TX_PATH="./build/gentx/validator-02.tmp.json"
DECIMAL_VALIDATOR_03_TX_PATH="./build/gentx/validator-03.tmp.json"
DECIMAL_INITIAL_STAKE_AMOUNT=30000000000000000000000000del

# Ensure tmp directory exists
mkdir "./build/gentx" -p

################################################################
# Validator 1
################################################################

# Generate and sign initial declare candidate transaction
dscd gentx dev-dsc-validator-0 $DECIMAL_INITIAL_STAKE_AMOUNT \
    --home "$HOME/.dscd/daemon" \
    --chain-id "decimal_2020-22093001" \
    --details "Initial validator Alpha (dev-node-01)" \
    --moniker "Validator Alpha" \
    --website "decimalchain.com" \
    --security-contact "security@decimalchain.com" \
    --node-id "ab5c39fc39819c7cbfef3a9346f3f5913bf6b5f9" \
    --ip "185.242.122.118" \
    --pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"gpLQaLp96T5XmC49PwcS0zP+CzD9j1VcODCspcFhgiE="}' \
    --offline --sequence 0 --keyring-backend test --output-document "$DECIMAL_VALIDATOR_01_TX_PATH" 2>/dev/null

# Print signature only and remove temporary file
jq <"$DECIMAL_VALIDATOR_01_TX_PATH" '.signatures[0]'

################################################################
# Validator 2
################################################################

# Generate and sign initial declare candidate transaction
dscd gentx dev-dsc-validator-1 $DECIMAL_INITIAL_STAKE_AMOUNT \
    --home "$HOME/.dscd/daemon" \
    --chain-id "decimal_2020-22093001" \
    --details "Initial validator Beta (dev-node-02)" \
    --moniker "Validator Beta" \
    --website "decimalchain.com" \
    --security-contact "security@decimalchain.com" \
    --node-id "7a71ef63609058a99573e2150b1b12b5a56d0c69" \
    --ip "185.242.122.119" \
    --pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"0WHq1uc29s1VEZoUy0hT90JPGnWmVSlniiBXIpvMhdo="}' \
    --offline --sequence 0 --keyring-backend test --output-document "$DECIMAL_VALIDATOR_02_TX_PATH" 2>/dev/null

# Print signature only and remove temporary file
jq <"$DECIMAL_VALIDATOR_02_TX_PATH" '.signatures[0]'

################################################################
# Validator 3
################################################################

# Generate and sign initial declare candidate transaction
dscd gentx dev-dsc-validator-2 $DECIMAL_INITIAL_STAKE_AMOUNT \
    --home "$HOME/.dscd/daemon" \
    --chain-id "decimal_2020-22093001" \
    --details "Initial validator Gamma (dev-node-03)" \
    --moniker "Validator Gamma" \
    --website "decimalchain.com" \
    --security-contact "security@decimalchain.com" \
    --node-id "8930e434c77caab942e67b066c6b43f53a15b3ad" \
    --ip "185.242.122.124" \
    --pubkey '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"gBn6rZrJPvpcqbgSg1f616Czd3ixIFyCLD5+mCpBD90="}' \
    --offline --sequence 0 --keyring-backend test --output-document "$DECIMAL_VALIDATOR_03_TX_PATH" 2>/dev/null

# Print signature only and remove temporary file
jq <"$DECIMAL_VALIDATOR_03_TX_PATH" '.signatures[0]'

# Also remove `build/gentx` folder
rm -rf "./build/gentx"
