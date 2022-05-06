#!/bin/bash

shopt -s extglob

echo "Current node is not a validator."
rm -r ~/.decimal/daemon/data
cd ~/.decimal/daemon/config && rm -- !(config.toml|genesis.json)

echo "Restating daemon."
sudo systemctl restart dscd
