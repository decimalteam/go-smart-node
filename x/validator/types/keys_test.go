package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var (
	keysPK1   = ed25519.GenPrivKeyFromSecret([]byte{1}).PubKey()
	keysPK2   = ed25519.GenPrivKeyFromSecret([]byte{2}).PubKey()
	keysPK3   = ed25519.GenPrivKeyFromSecret([]byte{3}).PubKey()
	keysAddr1 = keysPK1.Address()
	keysAddr2 = keysPK2.Address()
	keysAddr3 = keysPK3.Address()
)

func TestDelegations(t *testing.T) {
	val := sdk.ValAddress(keysAddr1)
	del := sdk.AccAddress(keysAddr2)
	denom := cmdcfg.BaseDenom

	key := append(append(append(types.GetAllDelegationsKey(), address.MustLengthPrefix(del)...), address.MustLengthPrefix(val)...), []byte(denom)...)
	// <delegator> <validator> <denom>
	delegationKey := types.GetDelegationKey(del, val, denom)
	require.Equal(t, key, delegationKey)

	// <delegator> <validator>
	delegatorDelegationsKey := types.GetDelegationsKey(del, val)
	require.Equal(t, delegationKey[:len(delegationKey)-len(denom)], delegatorDelegationsKey)

	// <delegator>
	delegatorAllDelegations := types.GetDelegatorDelegationsKey(del)
	require.Equal(t, delegationKey[:len(delegationKey)-len(del)-1-len(denom)], delegatorAllDelegations)

	// <validator> <delegator> <denom>
	valDelKey := types.GetValidatorDelegatorDelegationKey(val, del, denom)

	// <validator> <delegator> <denom> ------> <delegator> <validator> <denom>
	converterdKey := types.GetDelegationKeyFromValIndexKey(valDelKey)
	require.Equal(t, delegationKey, converterdKey)
}

func TestRedelegations(t *testing.T) {
	del := sdk.AccAddress(keysAddr2)
	val_src := sdk.ValAddress(keysAddr3)
	val_dst := sdk.ValAddress(keysAddr1)

	key := append(append(append(types.GetAllREDsKey(), address.MustLengthPrefix(del)...), address.MustLengthPrefix(val_src)...), address.MustLengthPrefix(val_dst)...)
	// <delegator> <src_validator> <dst_validator>
	redelegationKey := types.GetREDKey(del, val_src, val_dst)
	require.Equal(t, key, redelegationKey)

	// <src_validator> <delegator> <dst_validator> ----> <delegator> <src_validator> <dst_validator>
	src := types.GetREDByValSrcIndexKey(del, val_src, val_dst)
	redKey := types.GetREDKeyFromValSrcIndexKey(src)
	require.Equal(t, redelegationKey, redKey)

	// <dst_validator> <delegator> <src_validator> -> <delegator> <src_validator> <dst_validator>
	dst := types.GetREDByValDstIndexKey(del, val_src, val_dst)
	redKey = types.GetREDKeyFromValDstIndexKey(dst)
	require.Equal(t, redelegationKey, redKey)

	// <dst_validator> <delegator>
	dstDelReds := types.GetREDsByDelToValDstIndexKey(del, val_dst)
	t.Log(len(dst) - len(val_src))
	require.Equal(t, dst[:len(dst)-len(val_src)-1], dstDelReds)
}

func TestUndelegations(t *testing.T) {
	del := sdk.AccAddress(keysAddr2)
	val := sdk.ValAddress(keysAddr3)

	key := append(append(append(types.GetAllUBDsKey(), address.MustLengthPrefix(del)...), address.MustLengthPrefix(val)...))

	// <delegator> <validator>
	undelegationKey := types.GetUBDKey(del, val)
	require.Equal(t, key, undelegationKey)

	// <delegator>
	ubds := types.GetUBDsKey(del)
	require.Equal(t, key[:len(key)-len(val)-1], ubds)

	//<validator> <delegator>
	validatorUbd := types.GetUBDByValIndexKey(del, val)
	recoverKey := types.GetUBDKeyFromValIndexKey(validatorUbd)
	require.Equal(t, key, recoverKey)

	// <validator>
	validatorUbds := types.GetUBDsByValIndexKey(val)
	require.Equal(t, validatorUbd[:len(validatorUbd)-len(del)-1], validatorUbds)
}
