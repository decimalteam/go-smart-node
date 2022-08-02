#!/bin/bash

# this script prepares testnet files to launch decimal blockchain by 'docker-compose' with COUNT validators
# usage:
# ./testnetgenerate.sh [validators_count [path_for_files]]
# ./testnetgenerate.sh           : 3 validators, files in ./testnet
# ./testnetgenerate.sh 5         : 5 validators, files in ./testnet
# ./testnetgenerate.sh 5 ./test  : 5 validators, files in ./test

DOCKERCOMPOSEFILE=docker-compose.yml
# subnet range, must be X.Y.Z
IPPREFIX=172.18.2

# need build for 'testnet' command
make install

# build docker image 'dscd-testnet' and add network
docker build -t dscd-testnet .
docker container prune -f
docker network rm smartnode_net
docker network prune -f
docker network create --gateway $IPPREFIX.1 --subnet $IPPREFIX.0/24 smartnode_net

### input parameters
COUNT=$1
if [[ -z "$1" ]]; then
   COUNT=3
fi

ROOTHOME=$2
if [[ -z "$2" ]]; then
   ROOTHOME=./testnet
fi

# sudo: workaround for 'root' UID from docker
sudo rm -rf $ROOTHOME
mkdir $ROOTHOME
CHAINID="decimal_202020-1"
KEYALGO="eth_secp256k1"

# generate testnet files
# default amount of coins is 5000 for every account
# default stake is 100 for every validator
# mnemonics will be save in key_seed.json for every node
dscd --keyring-backend test --algo $KEYALGO --chain-id $CHAINID testnet init-files --v $COUNT --output-dir "$ROOTHOME" --node-dir-prefix "node" --starting-ip-address $IPPREFIX.100 --node-daemon-home ".decimal/daemon"

# make 'docker-compose.yml'
echo "version: '3'" > $DOCKERCOMPOSEFILE
echo "networks:" >> $DOCKERCOMPOSEFILE
echo "  smartnode_net:" >> $DOCKERCOMPOSEFILE
echo "    external: true" >> $DOCKERCOMPOSEFILE

fullprefix=$(realpath "$ROOTHOME")/node

echo "services:" >> $DOCKERCOMPOSEFILE
for ((i = 0 ; i < $COUNT ; i++))
do
ip=$((100+i))
echo "  node$i:" >> $DOCKERCOMPOSEFILE
echo "    image: dscd-testnet" >> $DOCKERCOMPOSEFILE
echo "    container_name: node$i" >> $DOCKERCOMPOSEFILE
echo "    volumes:" >> $DOCKERCOMPOSEFILE
echo "      - $fullprefix$i:/root" >> $DOCKERCOMPOSEFILE
echo "    networks:" >> $DOCKERCOMPOSEFILE
echo "      smartnode_net:" >> $DOCKERCOMPOSEFILE
echo "        ipv4_address: $IPPREFIX.$ip" >> $DOCKERCOMPOSEFILE
done
