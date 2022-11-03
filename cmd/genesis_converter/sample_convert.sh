#!/bin/bash
go build

# testnet
./genesis_converter -basedenom tdel -injectlegacy \
    -decimal ~/smartchain/testnet-genesis/genesis.json \
    -source ~/go-smart-node/ci-scripts/ansible/roles/decimal-init/templates/genesis-testnet-params.json.j2 \
    -nftfix ~/37928/lost_nft.json \
    -nftdublicates "~/testnet-genesis/nftdublicates.json" \
    -nfturiprefix "https://testnet-nft.decimalchain.com/api/nfts/" \
    -result ~/go-smart-node/ci-scripts/ansible/roles/decimal-init/templates/genesis-testnet.json.j2 > conv.log
dscd validate-genesis ~/smartchain/testnet-genesis/dsc-genesis-testnet.json

#mainnet
./genesis_converter -basedenom del \
    -decimal ~/smartchain/mainnet-genesis/genesis.json \
    -source ~/go-smart-node/ci-scripts/ansible/roles/decimal-init/templates/genesis-testnet-params.json.j2 \
    -nftfix ~/37928/lost_nft.json \
    -nftdublicates ~/smartchain/mainnet-genesis/nftdublicates.json \
    -nfturiprefix "https://wherebuynft.com/api/nfts/" \
    -result ~/smartchain/mainnet-genesis/dsc-genesis-mainnet.json > conv-main.log
dscd validate-genesis ~/smartchain/mainnet-genesis/dsc-genesis-mainnet.json