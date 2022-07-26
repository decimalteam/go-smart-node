* Sendstorm [IN PROGRESS]

Tool to test Decimal blockchain with randomized transactions

Usage:

0) Common flags for all commands:

`--mnemonics_file path_to_file, default: mnemonics.cfg`

`--node <node url>, default: http://localhost`

`--tport <tendermint port>, default: 26657`

`--rport <cosmos REST port>, default: 1317`

`--debug, default: false` -- turn on requests/resposes for resty library to sendstorm.log

1) Generate mnemonics for accounts

`sendstorm generate [--count count_of_mnemonics]`

2) Refill account with coins from faucet

`sendstorm faucet --faucet_mnemonic="..." [--only_empty]`

 `--only_empty false` refill without balance checking

3) Send bunch of transactions to blockchain

`sendstorm run [--actions actions_definition] [--tps transactions_per_second]`

Actions: `action_name1=weight1,action_name2=weight2...`, ex: `CreateCoin=1,SendCoin=99`