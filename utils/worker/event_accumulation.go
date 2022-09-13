package worker

/*
Event accumulation:
Prepare data for backend.
Goal: good update set for decimal-models.

1) Balance changes
records like (address, address_type, coin_symbol, amount, operation_type)
address_type - 'single' - common from mnemonic, 'multisig'(?) - multisignature wallet...
operation_type: "+"/"-"/"fee"
!!! +fee deduction
from coin sends, multisig sends, nft reserving, swaps, etc...

2) Coin changes:
last event EditCoin group by coin

3) NFT owners changes ()
(denom,id,subtokens,old owner,new owner)

4) Replace owner (legacy events)

5) multisig events as is
*/
