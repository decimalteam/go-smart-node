<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [decimal/coin/v1/coin.proto](#decimal/coin/v1/coin.proto)
    - [Coin](#decimal.coin.v1.Coin)
  
- [decimal/coin/v1/genesis.proto](#decimal/coin/v1/genesis.proto)
    - [GenesisState](#decimal.coin.v1.GenesisState)
  
- [decimal/coin/v1/query.proto](#decimal/coin/v1/query.proto)
    - [QueryCoinRequest](#decimal.coin.v1.QueryCoinRequest)
    - [QueryCoinResponse](#decimal.coin.v1.QueryCoinResponse)
    - [QueryCoinsRequest](#decimal.coin.v1.QueryCoinsRequest)
    - [QueryCoinsResponse](#decimal.coin.v1.QueryCoinsResponse)
  
    - [Query](#decimal.coin.v1.Query)
  
- [decimal/coin/v1/tx.proto](#decimal/coin/v1/tx.proto)
    - [MsgBuyCoin](#decimal.coin.v1.MsgBuyCoin)
    - [MsgBuyCoinResponse](#decimal.coin.v1.MsgBuyCoinResponse)
    - [MsgCreateCoin](#decimal.coin.v1.MsgCreateCoin)
    - [MsgCreateCoinResponse](#decimal.coin.v1.MsgCreateCoinResponse)
    - [MsgMultiSendCoin](#decimal.coin.v1.MsgMultiSendCoin)
    - [MsgMultisendCoinResponse](#decimal.coin.v1.MsgMultisendCoinResponse)
    - [MsgRedeemCheck](#decimal.coin.v1.MsgRedeemCheck)
    - [MsgRedeemCheckResponse](#decimal.coin.v1.MsgRedeemCheckResponse)
    - [MsgSellAllCoin](#decimal.coin.v1.MsgSellAllCoin)
    - [MsgSellAllCoinResponse](#decimal.coin.v1.MsgSellAllCoinResponse)
    - [MsgSellCoin](#decimal.coin.v1.MsgSellCoin)
    - [MsgSellCoinReponse](#decimal.coin.v1.MsgSellCoinReponse)
    - [MsgSendCoin](#decimal.coin.v1.MsgSendCoin)
    - [MsgSendCoinResponse](#decimal.coin.v1.MsgSendCoinResponse)
    - [MsgUpdateCoin](#decimal.coin.v1.MsgUpdateCoin)
    - [MsgUpdateCoinResponse](#decimal.coin.v1.MsgUpdateCoinResponse)
    - [Send](#decimal.coin.v1.Send)
  
    - [Msg](#decimal.coin.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="decimal/coin/v1/coin.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/coin/v1/coin.proto



<a name="decimal.coin.v1.Coin"></a>

### Coin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `constant_reserve_ratio` | [uint64](#uint64) |  |  |
| `symbol` | [string](#string) |  |  |
| `reserve` | [bytes](#bytes) |  |  |
| `limit_volume` | [bytes](#bytes) |  |  |
| `volume` | [bytes](#bytes) |  |  |
| `creator` | [string](#string) |  |  |
| `identity` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="decimal/coin/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/coin/v1/genesis.proto



<a name="decimal.coin.v1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `limit_volume` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="decimal/coin/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/coin/v1/query.proto



<a name="decimal.coin.v1.QueryCoinRequest"></a>

### QueryCoinRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `symbol` | [string](#string) |  |  |






<a name="decimal.coin.v1.QueryCoinResponse"></a>

### QueryCoinResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="decimal.coin.v1.QueryCoinsRequest"></a>

### QueryCoinsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="decimal.coin.v1.QueryCoinsResponse"></a>

### QueryCoinsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="decimal.coin.v1.Query"></a>

### Query
Query defines the gRPC querier service for coin module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Coin` | [QueryCoinRequest](#decimal.coin.v1.QueryCoinRequest) | [QueryCoinResponse](#decimal.coin.v1.QueryCoinResponse) | Coin queries existing coin by specific symbol. | GET|/coin/v1/coin/{symbol}|
| `Coins` | [QueryCoinsRequest](#decimal.coin.v1.QueryCoinsRequest) | [QueryCoinsResponse](#decimal.coin.v1.QueryCoinsResponse) | Coins queries list of all existing coins. | GET|/coin/v1/coins|

 <!-- end services -->



<a name="decimal/coin/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/coin/v1/tx.proto



<a name="decimal.coin.v1.MsgBuyCoin"></a>

### MsgBuyCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `coin_to_buy` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `max_coin_to_sell` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="decimal.coin.v1.MsgBuyCoinResponse"></a>

### MsgBuyCoinResponse







<a name="decimal.coin.v1.MsgCreateCoin"></a>

### MsgCreateCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `title` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `constant_reserve_ration` | [uint64](#uint64) |  |  |
| `initial_volume` | [bytes](#bytes) |  |  |
| `initial_reserve` | [bytes](#bytes) |  |  |
| `limit_volume` | [bytes](#bytes) |  |  |
| `identity` | [string](#string) |  |  |






<a name="decimal.coin.v1.MsgCreateCoinResponse"></a>

### MsgCreateCoinResponse







<a name="decimal.coin.v1.MsgMultiSendCoin"></a>

### MsgMultiSendCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `sends` | [Send](#decimal.coin.v1.Send) | repeated |  |






<a name="decimal.coin.v1.MsgMultisendCoinResponse"></a>

### MsgMultisendCoinResponse







<a name="decimal.coin.v1.MsgRedeemCheck"></a>

### MsgRedeemCheck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `check` | [bytes](#bytes) |  |  |
| `proof` | [bytes](#bytes) |  |  |






<a name="decimal.coin.v1.MsgRedeemCheckResponse"></a>

### MsgRedeemCheckResponse







<a name="decimal.coin.v1.MsgSellAllCoin"></a>

### MsgSellAllCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `coin_to_sell` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `min_coin_to_buy` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="decimal.coin.v1.MsgSellAllCoinResponse"></a>

### MsgSellAllCoinResponse







<a name="decimal.coin.v1.MsgSellCoin"></a>

### MsgSellCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `coin_to_sell` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `min_coin_to_buy` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="decimal.coin.v1.MsgSellCoinReponse"></a>

### MsgSellCoinReponse







<a name="decimal.coin.v1.MsgSendCoin"></a>

### MsgSendCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `receiver` | [string](#string) |  |  |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="decimal.coin.v1.MsgSendCoinResponse"></a>

### MsgSendCoinResponse







<a name="decimal.coin.v1.MsgUpdateCoin"></a>

### MsgUpdateCoin



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `limit_volume` | [bytes](#bytes) |  |  |
| `identity` | [string](#string) |  |  |






<a name="decimal.coin.v1.MsgUpdateCoinResponse"></a>

### MsgUpdateCoinResponse







<a name="decimal.coin.v1.Send"></a>

### Send



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `receiver` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="decimal.coin.v1.Msg"></a>

### Msg
Msg defines the coin Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateCoin` | [MsgCreateCoin](#decimal.coin.v1.MsgCreateCoin) | [MsgCreateCoinResponse](#decimal.coin.v1.MsgCreateCoinResponse) | CreateCoin defines message for new coin creation. | |
| `UpdateCoin` | [MsgUpdateCoin](#decimal.coin.v1.MsgUpdateCoin) | [MsgUpdateCoinResponse](#decimal.coin.v1.MsgUpdateCoinResponse) | UpdateCoin defines message for modifying existing coin. | |
| `SendCoin` | [MsgSendCoin](#decimal.coin.v1.MsgSendCoin) | [MsgSendCoinResponse](#decimal.coin.v1.MsgSendCoinResponse) | SendCoin defines message for transfering specific coin. | |
| `MultiSendCoin` | [MsgMultiSendCoin](#decimal.coin.v1.MsgMultiSendCoin) | [MsgMultisendCoinResponse](#decimal.coin.v1.MsgMultisendCoinResponse) | MultiSendCoin defines message for transfering specific coins as a batch. | |
| `BuyCoin` | [MsgBuyCoin](#decimal.coin.v1.MsgBuyCoin) | [MsgBuyCoinResponse](#decimal.coin.v1.MsgBuyCoinResponse) | BuyCoin defines message for buying specific coin. | |
| `SellCoin` | [MsgSellCoin](#decimal.coin.v1.MsgSellCoin) | [MsgSellCoinReponse](#decimal.coin.v1.MsgSellCoinReponse) | SellCoin defines message for selling specific coin. | |
| `SellAllCoin` | [MsgSellAllCoin](#decimal.coin.v1.MsgSellAllCoin) | [MsgSellAllCoinResponse](#decimal.coin.v1.MsgSellAllCoinResponse) | SellAllCoin defines message for selling all specific coin. | |
| `RedeemCheck` | [MsgRedeemCheck](#decimal.coin.v1.MsgRedeemCheck) | [MsgRedeemCheckResponse](#decimal.coin.v1.MsgRedeemCheckResponse) | RedeemCheck defines message for redeeming checks. | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

