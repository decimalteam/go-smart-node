<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [decimal/swap/v1/tx.proto](#decimal/swap/v1/tx.proto)
    - [MsgHTLT](#decimal.swap.v1.MsgHTLT)
    - [MsgHTLTResponse](#decimal.swap.v1.MsgHTLTResponse)
    - [MsgRedeem](#decimal.swap.v1.MsgRedeem)
    - [MsgRedeemResponse](#decimal.swap.v1.MsgRedeemResponse)
    - [MsgRefund](#decimal.swap.v1.MsgRefund)
    - [MsgRefundResponse](#decimal.swap.v1.MsgRefundResponse)
  
    - [TransferType](#decimal.swap.v1.TransferType)
  
    - [Msg](#decimal.swap.v1.Msg)
  
- [decimal/swap/v1/swap.proto](#decimal/swap/v1/swap.proto)
    - [Params](#decimal.swap.v1.Params)
    - [Swap](#decimal.swap.v1.Swap)
  
- [decimal/swap/v1/genesis.proto](#decimal/swap/v1/genesis.proto)
    - [GenesisState](#decimal.swap.v1.GenesisState)
  
- [Scalar Value Types](#scalar-value-types)



<a name="decimal/swap/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/swap/v1/tx.proto



<a name="decimal.swap.v1.MsgHTLT"></a>

### MsgHTLT



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfer_type` | [TransferType](#decimal.swap.v1.TransferType) |  |  |
| `from` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `hashed_secret` | [bytes](#bytes) |  |  |
| `secret` | [bytes](#bytes) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="decimal.swap.v1.MsgHTLTResponse"></a>

### MsgHTLTResponse







<a name="decimal.swap.v1.MsgRedeem"></a>

### MsgRedeem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `secret` | [bytes](#bytes) |  |  |






<a name="decimal.swap.v1.MsgRedeemResponse"></a>

### MsgRedeemResponse







<a name="decimal.swap.v1.MsgRefund"></a>

### MsgRefund



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `hashed_secret` | [bytes](#bytes) |  |  |






<a name="decimal.swap.v1.MsgRefundResponse"></a>

### MsgRefundResponse






 <!-- end messages -->


<a name="decimal.swap.v1.TransferType"></a>

### TransferType


| Name | Number | Description |
| ---- | ------ | ----------- |
| TransferTypeOut | 0 |  |
| TransferTypeIn | 1 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="decimal.swap.v1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `HTLT` | [MsgHTLT](#decimal.swap.v1.MsgHTLT) | [MsgHTLTResponse](#decimal.swap.v1.MsgHTLTResponse) |  | |
| `Redeem` | [MsgRedeem](#decimal.swap.v1.MsgRedeem) | [MsgRedeemResponse](#decimal.swap.v1.MsgRedeemResponse) |  | |
| `Refund` | [MsgRefund](#decimal.swap.v1.MsgRefund) | [MsgRefundResponse](#decimal.swap.v1.MsgRefundResponse) |  | |

 <!-- end services -->



<a name="decimal/swap/v1/swap.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/swap/v1/swap.proto



<a name="decimal.swap.v1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `locked_time_out` | [int64](#int64) |  |  |
| `locked_time_in` | [int64](#int64) |  |  |






<a name="decimal.swap.v1.Swap"></a>

### Swap



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfer_type` | [TransferType](#decimal.swap.v1.TransferType) |  |  |
| `hashed_secret` | [string](#string) |  |  |
| `from` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `timestamp` | [uint64](#uint64) |  |  |
| `redeemed` | [bool](#bool) |  |  |
| `refunded` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="decimal/swap/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## decimal/swap/v1/genesis.proto



<a name="decimal.swap.v1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#decimal.swap.v1.Params) |  | params defines all the paramaters of related to deposit. |
| `swaps` | [Swap](#decimal.swap.v1.Swap) | repeated | swaps defines the swaps active at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

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

