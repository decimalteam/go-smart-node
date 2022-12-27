package types

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "validator"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// RouterKey defines module's messages routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)
)

// TODO: It will be better to prepare key bytes and then copy parts into it instead of multiple concatenations.

// Lets assume parts of keys are:
//   <validator>     = <val_address>           [20 bytes]
//   <val_src>       = <val_address>           [20 bytes]
//   <val_dst>       = <val_address>           [20 bytes]
//   <consensus>     = <cons_address>          [20 bytes]
//   <delegator>     = <acc_address>           [20 bytes]
//   <denom>         = <denom>                 [3-10 bytes]
//   <token_id>      = <hash(id)>              [32 bytes]
//   <stake_id>      = <denom> | <token_id>    [3-32 bytes]
//   <block_id>      = <block_number_dec>      [1-9 bytes]

// Core consensus related records (constant during a block):
//   - LastValidatorPowers:     0x11<validator>    : <int64>
//   - LastTotalPower:          0x12               : <int64>

// Validator related records:
//   - Validators:              0x21<validator>    : <Validator>
//   - ValidatorsByConsAddr:    0x22<cons_addr>    : <val_address>
//   - ValidatorsByPower:       0x23<int64>        : <val_address>

// Staking related records:
//   - Delegations:             0x31<delegator><validator><stake_id>          : <Stake>
//   - Redelegations:           0x32<delegator><val_src><val_dst>   		  : <Redelegation>
//   - RedelegationsByValSrc:   0x33<val_src><delegator><val_dst>             : []byte{}
//   - RedelegationsByValDst:   0x34<val_dst><delegator><val_src>             : []byte{}
//   - Undelegations:           0x35<delegator><validator><stake_id>          : <Undelegation>
//   - UndelegationsByValSrc:   0x36<validator><delegator><stake_id>          : []byte{}
//   - DelegationsByVal:        0x37<validator><delegator><stake_id>		  : []byte{}
//   - DelegationsCount:        0x38<validator>                               : <int32>

// Queues related records:
// TODO: Instead of storing array we need to store records separately and iterate over it when needed.
//   - ValidatorQueues:        0x41<time><height> : <ValAddresses>
//   - RedelegationQueues:     0x42<time>         : 0x32<delegator><val_src><val_dst><stake_id> (key of the redelegation record)
//   - UndelegationQueues:     0x43<time>         : 0x35<delegator><validator><stake_id>        (key of the undelegation record)

// ABCI related records:
//   - HistoricalInfo:          0x51<block_id>     : <HistoricalInfo>

// Missed blocks records:
//   - MissedBlock:          0x61<cons_addr><block_id> : []byte{}
//   - StartingBlock:        0x62<cons_addr>           : <block_id>

// Delegation related records:
//   - CustomCoinStaked:  0x71<denom>   :  sdkmath.Int
var (
	keyPrefixLastValidatorPowers        = []byte{0x11} // prefix for each key to a validator index (for bonded validators)
	keyPrefixLastTotalPower             = []byte{0x12} // prefix for the total power record
	keyPrefixValidators                 = []byte{0x21} // prefix for each key to a validator
	keyPrefixValidatorsByConsAddrIndex  = []byte{0x22} // prefix for each key to a validator index (by consensus address)
	keyPrefixValidatorsByPowerIndex     = []byte{0x23} // prefix for each key to a validator index (sorted by power)
	keyPrefixValidatorRewards           = []byte{0x24} // prefix for validator rewards
	keyPrefixDelegations                = []byte{0x31} // prefix for each key for a delegation
	keyPrefixRedelegations              = []byte{0x32} // prefix for each key for a redelegation
	keyPrefixRedelegationsByValSrcIndex = []byte{0x33} // prefix for each key for a redelegation index (by source validator address)
	keyPrefixRedelegationsByValDstIndex = []byte{0x34} // prefix for each key for a redelegation index (by destination validator address)
	keyPrefixUndelegations              = []byte{0x35} // prefix for each key for an undelegation
	keyPrefixUndelegationsByValIndex    = []byte{0x36} // prefix for each key for an undelegation index (by validator address)
	keyPrefixDelegationByValIndex       = []byte{0x37} // prefix for each key for a delegation key (by validator address)
	keyPrefixDelegationsCount           = []byte{0x38} // prefix for delegations count index (by validator address)
	keyPrefixValidatorQueue             = []byte{0x41} // prefix for the timestamps in validator queue
	keyPrefixRedelegationQueue          = []byte{0x42} // prefix for the timestamps in redelegations queue
	keyPrefixUndelegationQueue          = []byte{0x43} // prefix for the timestamps in unbonding queue
	keyPrefixHistoricalInfo             = []byte{0x51} // prefix for the historical info
	keyPrefixMissedBlock                = []byte{0x61} // prefix for missed blocks
	keyPrefixStartHeight                = []byte{0x62} // prefix for starting block
	keyPrefixCustomCoinStaked           = []byte{0x71} // prefix for custom coin total staked in delegations
)

////////////////////////////////////////////////////////////////////////////////
// Last total validator power //////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetLastValidatorPowersKey returns prefix key for bonded validators.
func GetLastValidatorPowersKey() []byte {
	return keyPrefixLastValidatorPowers
}

// GetLastValidatorPowerKey creates the bonded validator index key for the given operator address.
func GetLastValidatorPowerKey(validator sdk.ValAddress) []byte {
	return append(GetLastValidatorPowersKey(), validator.Bytes()...)
}

// GetLastTotalPowerKey returns key for the total power record.
func GetLastTotalPowerKey() []byte {
	return keyPrefixLastTotalPower
}

////////////////////////////////////////////////////////////////////////////////
// Validators //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetValidatorsKey returns prefix key for all validators.
func GetValidatorsKey() []byte {
	return keyPrefixValidators
}

// GetValidatorKey creates the key for the validator with the given operator address.
func GetValidatorKey(validator sdk.ValAddress) []byte {
	return append(GetValidatorsKey(), validator.Bytes()...)
}

// GetValidatorByConsAddrIndexKey creates the key for the validator with the given consensus address.
func GetValidatorByConsAddrIndexKey(addr sdk.ConsAddress) []byte {
	return append(keyPrefixValidatorsByConsAddrIndex, addr.Bytes()...)
}

// GetValidatorsByPowerIndexKey returns the prefix key for the validators sorted by voting power.
func GetValidatorsByPowerIndexKey() []byte {
	return keyPrefixValidatorsByPowerIndex
}

func GetValidatorsRewards() []byte {
	return keyPrefixValidatorRewards
}

func GetValidatorRewards(addr sdk.ValAddress) []byte {
	return append(GetValidatorsRewards(), addr.Bytes()...)
}

////////////////////////////////////////////////////////////////////////////////
// Delegations /////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetAllDelegationsKey returns the prefix key for all delegations.
func GetAllDelegationsKey() []byte {
	return keyPrefixDelegations
}

// GetDelegatorDelegationsKey creates the prefix for delegations from the given delegator.
func GetDelegatorDelegationsKey(delegator sdk.AccAddress) []byte {
	return append(GetAllDelegationsKey(), address.MustLengthPrefix(delegator)...)
}

// GetDelegationsKey creates the prefix key for delegations bond with validator.
func GetDelegationsKey(delegator sdk.AccAddress, validator sdk.ValAddress) []byte {
	return append(GetDelegatorDelegationsKey(delegator), address.MustLengthPrefix(validator)...)
}

// GetDelegationKey creates the key for the exact delegation in the specified coin.
func GetDelegationKey(delegator sdk.AccAddress, validator sdk.ValAddress, denom string) []byte {
	return append(GetDelegationsKey(delegator, validator), []byte(denom)...)
}

// GetValidatorAllDelegations returns the prefix for all validators delegations index
func GetValidatorAllDelegations() []byte {
	return keyPrefixDelegationByValIndex
}

// GetValidatorDelegationsKey create the key for validator all delegations
func GetValidatorDelegationsKey(val sdk.ValAddress) []byte {
	return append(GetValidatorAllDelegations(), address.MustLengthPrefix(val)...)
}

// GetValidatorDelegatorDelegationsKey create a key for all delegations between delegator and validator
func GetValidatorDelegatorDelegationsKey(val sdk.ValAddress, del sdk.AccAddress) []byte {
	return append(GetValidatorDelegationsKey(val), address.MustLengthPrefix(del)...)
}

// GetValidatorDelegatorDelegationKey create a key for validator-delegator-denom delegation
func GetValidatorDelegatorDelegationKey(val sdk.ValAddress, del sdk.AccAddress, denom string) []byte {
	return append(GetValidatorDelegatorDelegationsKey(val, del), []byte(denom)...)
}

// GetDelegationKeyFromValIndexKey rearranges the ValIndexKey to get the DelegationKey
func GetDelegationKeyFromValIndexKey(indexKey []byte) []byte {
	kv.AssertKeyAtLeastLength(indexKey, 2)
	addrs := indexKey[1:] // remove prefix bytes

	// get validator
	valAddrLen := addrs[0]
	kv.AssertKeyAtLeastLength(addrs, int(valAddrLen)+2)
	validator := addrs[1 : valAddrLen+1]

	// get delegator
	delAddrLen := addrs[valAddrLen+1]
	kv.AssertKeyAtLeastLength(addrs, int(valAddrLen)+int(delAddrLen)+3)
	delegator := addrs[valAddrLen+2 : valAddrLen+2+delAddrLen]

	// get denom
	denom := string(addrs[valAddrLen+delAddrLen+2:])

	return GetDelegationKey(delegator, validator, denom)
}

////////////////////////////////////////////////////////////////////////////////
// Redelegations ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetAllREDsKey returns a key prefix for indexing all redelegations.
func GetAllREDsKey() []byte {
	return keyPrefixRedelegations
}

// GetREDsKey returns a key prefix for indexing a redelegation from a delegator address.
func GetREDsKey(delegator sdk.AccAddress) []byte {
	return append(keyPrefixRedelegations, address.MustLengthPrefix(delegator)...)
}

// GetREDKey returns a key prefix for indexing a redelegation from a delegator and source validator to a destination validator.
func GetREDKey(delegator sdk.AccAddress, validatorSrc sdk.ValAddress, validatorDst sdk.ValAddress) []byte {
	// key is of the form GetREDsKey || valSrcAddrLen (1 byte) || validatorSrc || valDstAddrLen (1 byte) || validatorDst
	key := make([]byte, 1+3+len(delegator)+len(validatorSrc)+len(validatorDst))
	copy(key[0:2+len(delegator)], GetREDsKey(delegator.Bytes()))
	key[2+len(delegator)] = byte(len(validatorSrc))
	copy(key[3+len(delegator):3+len(delegator)+len(validatorSrc)], validatorSrc.Bytes())
	key[3+len(delegator)+len(validatorSrc)] = byte(len(validatorDst))
	copy(key[4+len(delegator)+len(validatorSrc):], validatorDst.Bytes())

	return key
}

// GetREDsFromValSrcIndexKey returns a key prefix for indexing a redelegation to a source validator.
func GetREDsFromValSrcIndexKey(validatorSrc sdk.ValAddress) []byte {
	return append(keyPrefixRedelegationsByValSrcIndex, address.MustLengthPrefix(validatorSrc)...)
}

// GetREDByValSrcIndexKey creates the index-key for a redelegation, stored by source-validator-index.
func GetREDByValSrcIndexKey(delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress) []byte {
	REDSFromValsSrcKey := GetREDsFromValSrcIndexKey(validatorSrc)
	offset := len(REDSFromValsSrcKey)

	// key is of the form REDSFromValsSrcKey || delAddrLen (1 byte) || delegator || valDstAddrLen (1 byte) || validatorDst
	key := make([]byte, offset+2+len(delegator)+len(validatorDst))
	copy(key[0:offset], REDSFromValsSrcKey)
	key[offset] = byte(len(delegator))
	copy(key[offset+1:offset+1+len(delegator)], delegator.Bytes())
	key[offset+1+len(delegator)] = byte(len(validatorDst))
	copy(key[offset+2+len(delegator):], validatorDst.Bytes())

	return key
}

// GetREDsToValDstIndexKey returns a key prefix for indexing a redelegation to a destination (target) validator.
func GetREDsToValDstIndexKey(validatorDst sdk.ValAddress) []byte {
	return append(keyPrefixRedelegationsByValDstIndex, address.MustLengthPrefix(validatorDst)...)
}

// GetREDByValDstIndexKey creates the index-key for a redelegation, stored by destination-validator-index.
func GetREDByValDstIndexKey(delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress) []byte {
	REDSToValsDstKey := GetREDsToValDstIndexKey(validatorDst)
	offset := len(REDSToValsDstKey)

	// key is of the form REDSToValsDstKey || delAddrLen (1 byte) || delegator || valSrcAddrLen (1 byte) || validatorSrc
	key := make([]byte, offset+2+len(delegator)+len(validatorSrc))
	copy(key[0:offset], REDSToValsDstKey)
	key[offset] = byte(len(delegator))
	copy(key[offset+1:offset+1+len(delegator)], delegator.Bytes())
	key[offset+1+len(delegator)] = byte(len(validatorSrc))
	copy(key[offset+2+len(delegator):], validatorSrc.Bytes())

	return key
}

// GetREDKeyFromValSrcIndexKey rearranges the ValSrcIndexKey to get the REDKey
func GetREDKeyFromValSrcIndexKey(indexKey []byte) []byte {
	// note that first byte is prefix byte, which we remove
	kv.AssertKeyAtLeastLength(indexKey, 2)
	addrs := indexKey[1:]

	valSrcAddrLen := addrs[0]
	kv.AssertKeyAtLeastLength(addrs, int(valSrcAddrLen)+2)
	validatorSrc := addrs[1 : valSrcAddrLen+1]
	delAddrLen := addrs[valSrcAddrLen+1]
	kv.AssertKeyAtLeastLength(addrs, int(valSrcAddrLen)+int(delAddrLen)+2)
	delegator := addrs[valSrcAddrLen+2 : valSrcAddrLen+2+delAddrLen]
	kv.AssertKeyAtLeastLength(addrs, int(valSrcAddrLen)+int(delAddrLen)+4)
	validatorDst := addrs[valSrcAddrLen+delAddrLen+3:]

	return GetREDKey(delegator, validatorSrc, validatorDst)
}

// GetREDKeyFromValDstIndexKey rearranges the ValDstIndexKey to get the REDKey
func GetREDKeyFromValDstIndexKey(indexKey []byte) []byte {
	// note that first byte is prefix byte, which we remove
	kv.AssertKeyAtLeastLength(indexKey, 2)
	addrs := indexKey[1:]

	valDstAddrLen := addrs[0]
	kv.AssertKeyAtLeastLength(addrs, int(valDstAddrLen)+2)
	validatorDst := addrs[1 : valDstAddrLen+1]
	delAddrLen := addrs[valDstAddrLen+1]
	kv.AssertKeyAtLeastLength(addrs, int(valDstAddrLen)+int(delAddrLen)+3)
	delegator := addrs[valDstAddrLen+2 : valDstAddrLen+2+delAddrLen]
	kv.AssertKeyAtLeastLength(addrs, int(valDstAddrLen)+int(delAddrLen)+4)
	validatorSrc := addrs[valDstAddrLen+delAddrLen+3:]

	return GetREDKey(delegator, validatorSrc, validatorDst)
}

// GetREDsByDelToValDstIndexKey returns a key prefix for indexing a redelegation from an address to a source validator.
func GetREDsByDelToValDstIndexKey(delegator sdk.AccAddress, validatorDst sdk.ValAddress) []byte {
	return append(GetREDsToValDstIndexKey(validatorDst), address.MustLengthPrefix(delegator)...)
}

////////////////////////////////////////////////////////////////////////////////
// Undelegations ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetAllUBDsKey returns a key prefix for indexing all undelegations.
func GetAllUBDsKey() []byte {
	return keyPrefixUndelegations
}

// GetUBDsKey creates the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delegator sdk.AccAddress) []byte {
	return append(GetAllUBDsKey(), address.MustLengthPrefix(delegator)...)
}

// GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
func GetUBDKey(delegator sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsKey(delegator.Bytes()), address.MustLengthPrefix(valAddr)...)
}

// GetUBDsByValIndexKey creates the prefix keyspace for the indexes of unbonding delegations for a validator
func GetUBDsByValIndexKey(valAddr sdk.ValAddress) []byte {
	return append(keyPrefixUndelegationsByValIndex, address.MustLengthPrefix(valAddr)...)
}

// GetUBDByValIndexKey creates the index-key for an unbonding delegation, stored by validator-index
func GetUBDByValIndexKey(delegator sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsByValIndexKey(valAddr), address.MustLengthPrefix(delegator)...)
}

// GetUBDKeyFromValIndexKey rearranges the ValIndexKey to get the UBDKey
func GetUBDKeyFromValIndexKey(indexKey []byte) []byte {
	kv.AssertKeyAtLeastLength(indexKey, 2)
	addrs := indexKey[1:] // remove prefix bytes

	valAddrLen := addrs[0]
	kv.AssertKeyAtLeastLength(addrs, 2+int(valAddrLen))
	valAddr := addrs[1 : 1+valAddrLen]
	kv.AssertKeyAtLeastLength(addrs, 3+int(valAddrLen))
	delegator := addrs[valAddrLen+2:]

	return GetUBDKey(delegator, valAddr)
}

////////////////////////////////////////////////////////////////////////////////
// Delegations count ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetAllDelegationsCount returns the prefix key for all counts of delegations.
func GetAllDelegationsCount() []byte {
	return keyPrefixDelegationsCount
}

func GetValidatorDelegationsCount(valAddr sdk.ValAddress) []byte {
	return append(GetAllDelegationsCount(), valAddr...)
}

func ParseValidatorDelegationsCountKey(key []byte) sdk.ValAddress {
	return key[1:]
}

////////////////////////////////////////////////////////////////////////////////
// Queues //////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func GetAllValidatorQueueKey() []byte {
	return keyPrefixValidatorQueue
}

// GetValidatorQueueKey returns the prefix key used for getting a set of unbonding
// validators whose unbonding completion occurs at the given time and height.
func GetValidatorQueueKey(timestamp time.Time, height int64) []byte {
	heightBz := sdk.Uint64ToBigEndian(uint64(height))
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len(keyPrefixValidatorQueue)

	bz := make([]byte, prefixL+8+timeBzL+8)
	// copy the prefix
	copy(bz[:prefixL], keyPrefixValidatorQueue)
	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))
	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	// copy the encoded height
	copy(bz[prefixL+8+timeBzL:], heightBz)

	return bz
}

// GetAllRedelegationsTimeKey creates the prefix for undelegations.
func GetAllRedelegationsTimeKey() []byte {
	return keyPrefixRedelegationQueue
}

// GetRedelegationsTimeKey returns a key prefix for indexing redelegations based on a completion time.
func GetRedelegationsTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(keyPrefixRedelegationQueue, bz...)
}

// GetAllUndelegationsTimeKey creates the prefix for undelegations.
func GetAllUndelegationsTimeKey() []byte {
	return keyPrefixUndelegationQueue
}

// GetUndelegationsTimeKey creates the prefix for undelegations based on a completion time.
func GetUndelegationsTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(keyPrefixUndelegationQueue, bz...)
}

////////////////////////////////////////////////////////////////////////////////
// Historical Info /////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetHistoricalInfosKey returns a key prefix for all HistoricalInfo objects.
func GetHistoricalInfosKey() []byte {
	return keyPrefixHistoricalInfo
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(GetHistoricalInfosKey(), []byte(strconv.FormatInt(height, 10))...)
}

////////////////////////////////////////////////////////////////////////////////
// Power Key ///////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AddressFromLastValidatorPowerKey creates the validator operator address from LastValidatorPowerKey.
func AddressFromLastValidatorPowerKey(key []byte) []byte {
	kv.AssertKeyAtLeastLength(key, 2)
	return key[1:] // remove prefix bytes and address length
}

// ParseValidatorPowerKey parses the validators operator address from voting power key.
func ParseValidatorPowerKey(key []byte) (operAddr []byte) {
	powerBytesLen := 8

	// key is of format prefix (1 byte) || powerbytes || addrLen (1byte) || addrBytes
	operAddr = sdk.CopyBytes(key[powerBytesLen+2:])

	for i, b := range operAddr {
		operAddr[i] = ^b
	}

	return operAddr
}

// ParseValidatorQueueKey returns the encoded time and height from a key created from GetValidatorQueueKey.
func ParseValidatorQueueKey(bz []byte) (time.Time, int64, error) {
	prefixL := len(keyPrefixValidatorQueue)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, keyPrefixValidatorQueue) {
		return time.Time{}, 0, fmt.Errorf("invalid prefix; expected: %X, got: %X", keyPrefixValidatorQueue, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	ts, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, 0, err
	}

	height := sdk.BigEndianToUint64(bz[prefixL+8+int(timeBzL):])

	return ts, int64(height), nil
}

////////////////////////////////////////////////////////////////////////////////
// Missed Blocks ///////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// GetMissedBlockKey returns key for save validator missed blocks
func GetMissedBlockKey(addr sdk.ConsAddress, height int64) []byte {
	// key format: prefix (1 byte) || consensus address || height (8 bytes)
	key := append(keyPrefixMissedBlock, addr.Bytes()...)
	key = append(key, sdk.Uint64ToBigEndian(uint64(height))...)
	return key
}

// GetStartHeightKey returns key for save validator first block stake
func GetStartHeightKey(addr sdk.ConsAddress) []byte {
	// key format: prefix (1 byte) || consensus address
	return append(keyPrefixStartHeight, addr.Bytes()...)
}

////////////////////////////////////////////////////////////////////////////////
// Custom Coins Staked /////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func GetAllCustomCoinsStaked() []byte {
	return keyPrefixCustomCoinStaked
}

func GetCustomCoinStaked(denom string) []byte {
	return append(keyPrefixCustomCoinStaked, []byte(denom)...)
}
