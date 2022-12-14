package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/encoding"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/client/cli"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

func TestCliCreateValidator(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewCreateValidatorCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())

	cmd.SetArgs([]string{
		fmt.Sprintf("--%s=%s", cli.FlagPubKey, consPubKeyBz),
		fmt.Sprintf("--%s=AFAF00C4", cli.FlagIdentity),
		fmt.Sprintf("--%s=moniker", cli.FlagMoniker),
		fmt.Sprintf("--%s=https://newvalidator.io", cli.FlagWebsite),
		fmt.Sprintf("--%s=contact@newvalidator.io", cli.FlagSecurityContact),
		fmt.Sprintf("--%s='Hey, I am a new validator. Please delegate!'", cli.FlagDetails),
		fmt.Sprintf("--%s=0.5", cli.FlagCommissionRate),
		fmt.Sprintf("--%s=100000del", cli.FlagAmount),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr0),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCreateValidator)
	require.True(t, ok)
	rate, err := sdk.NewDecFromStr("0.5")
	require.NoError(t, err)
	require.Equal(t, "100000del", msg.Stake.String())
	require.Equal(t, adr0.String(), msg.RewardAddress)
	require.Equal(t, rate, msg.Commission)
	require.Equal(t, "AFAF00C4", msg.Description.Identity)
	require.Equal(t, "'Hey, I am a new validator. Please delegate!'", msg.Description.Details)
	require.Equal(t, "moniker", msg.Description.Moniker)
	require.Equal(t, "contact@newvalidator.io", msg.Description.SecurityContact)
	require.Equal(t, "https://newvalidator.io", msg.Description.Website)
}

func TestCliEditValidator(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewEditValidatorCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", adr.String()),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=AFAF00C4", cli.FlagIdentity),
		fmt.Sprintf("--%s=moniker", cli.FlagEditMoniker),
		fmt.Sprintf("--%s=https://newvalidator.io", cli.FlagWebsite),
		fmt.Sprintf("--%s=contact@newvalidator.io", cli.FlagSecurityContact),
		fmt.Sprintf("--%s='Hey, I am a new validator. Please delegate!'", cli.FlagDetails),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgEditValidator)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.OperatorAddress)
	require.Equal(t, adr.String(), msg.RewardAddress)
	require.Equal(t, "AFAF00C4", msg.Description.Identity)
	require.Equal(t, "'Hey, I am a new validator. Please delegate!'", msg.Description.Details)
	require.Equal(t, "moniker", msg.Description.Moniker)
	require.Equal(t, "contact@newvalidator.io", msg.Description.SecurityContact)
	require.Equal(t, "https://newvalidator.io", msg.Description.Website)
}

func TestCliSetOnline(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewSetOnlineCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgSetOnline)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
}

func TestCliSetOffline(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewSetOfflineCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgSetOffline)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
}

func TestCliDelegate(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewDelegateCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "1000del"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgDelegate)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "1000del", msg.Coin.String())
}

func TestCliDelegateNFT(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewDelegateNFTCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "nft"),
		fmt.Sprintf("1,2,3"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgDelegateNFT)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "nft", msg.TokenID)
	require.Equal(t, []uint32{1, 2, 3}, msg.SubTokenIDs)
}

func TestCliUndelegate(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewUndelegateCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "1000del"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgUndelegate)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "1000del", msg.Coin.String())
}

func TestCliUndelegateNFT(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewUndelegateNFTCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "nft"),
		fmt.Sprintf("1,2,3"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgUndelegateNFT)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "nft", msg.TokenID)
	require.Equal(t, []uint32{1, 2, 3}, msg.SubTokenIDs)
}

func TestCliRedelegate(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)
	valSrcAddr := sdk.ValAddress(adr0)
	valDstAddr := sdk.ValAddress(adr1)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewRedelegateCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valSrcAddr.String()),
		fmt.Sprintf("%s", valDstAddr.String()),
		fmt.Sprintf("%s", "1000del"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr0),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgRedelegate)
	require.True(t, ok)
	require.Equal(t, valSrcAddr.String(), msg.ValidatorSrc)
	require.Equal(t, valDstAddr.String(), msg.ValidatorDst)
	require.Equal(t, adr0.String(), msg.Delegator)
	require.Equal(t, "1000del", msg.Coin.String())
}

func TestCliRedelegateNFT(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)
	valSrcAddr := sdk.ValAddress(adr0)
	valDstAddr := sdk.ValAddress(adr1)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewRedelegateNFTCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valSrcAddr.String()),
		fmt.Sprintf("%s", valDstAddr.String()),
		fmt.Sprintf("%s", "nft"),
		fmt.Sprintf("1,2,3"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr0),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgRedelegateNFT)
	require.True(t, ok)
	require.Equal(t, valSrcAddr.String(), msg.ValidatorSrc)
	require.Equal(t, valDstAddr.String(), msg.ValidatorDst)
	require.Equal(t, adr0.String(), msg.Delegator)
	require.Equal(t, "nft", msg.TokenID)
	require.Equal(t, []uint32{1, 2, 3}, msg.SubTokenIDs)
}

func TestCliCancelUndelegate(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewCancelUndelegateCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "1000del"),
		fmt.Sprintf("%s", "1000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCancelUndelegation)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "1000del", msg.Coin.String())
	require.Equal(t, int64(1000), msg.CreationHeight)
}

func TestCliCancelUndelegateNFT(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr, err := accs[0].GetAddress()
	require.NoError(t, err)
	valAddr := sdk.ValAddress(adr)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewCancelUndelegateNFTCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valAddr.String()),
		fmt.Sprintf("%s", "nft"),
		fmt.Sprintf("1,2,3"),
		fmt.Sprintf("%s", "1000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCancelUndelegationNFT)
	require.True(t, ok)
	require.Equal(t, valAddr.String(), msg.Validator)
	require.Equal(t, adr.String(), msg.Delegator)
	require.Equal(t, "nft", msg.TokenID)
	require.Equal(t, []uint32{1, 2, 3}, msg.SubTokenIDs)
	require.Equal(t, int64(1000), msg.CreationHeight)
}

func TestCliCancelRedelegate(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)
	valSrcAddr := sdk.ValAddress(adr0)
	valDstAddr := sdk.ValAddress(adr1)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewCancelRedelegateCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valSrcAddr.String()),
		fmt.Sprintf("%s", valDstAddr.String()),
		fmt.Sprintf("%s", "1000del"),
		fmt.Sprintf("%s", "1000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr0),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCancelRedelegation)
	require.True(t, ok)
	require.Equal(t, valSrcAddr.String(), msg.ValidatorSrc)
	require.Equal(t, valDstAddr.String(), msg.ValidatorDst)
	require.Equal(t, adr0.String(), msg.Delegator)
	require.Equal(t, "1000del", msg.Coin.String())
	require.Equal(t, int64(1000), msg.CreationHeight)
}

func TestCliCancelRedelegateNFT(t *testing.T) {
	clientCtx, accs, result := setUpCliTest(t, 2)

	enc := encoding.MakeConfig(app.ModuleBasics)
	adr0, err := accs[0].GetAddress()
	require.NoError(t, err)
	adr1, err := accs[1].GetAddress()
	require.NoError(t, err)
	valSrcAddr := sdk.ValAddress(adr0)
	valDstAddr := sdk.ValAddress(adr1)

	consPrivKey := ed25519.GenPrivKey()
	consPubKeyBz, err := enc.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
	require.NoError(t, err)
	require.NotNil(t, consPubKeyBz)

	cmd := cli.NewCancelRedelegateNFTCmd()
	ctx := setUpCmd(t, cmd, clientCtx, adr0.String())
	cmd.SetArgs([]string{
		fmt.Sprintf("%s", valSrcAddr.String()),
		fmt.Sprintf("%s", valDstAddr.String()),
		fmt.Sprintf("%s", "nft"),
		fmt.Sprintf("1,2,3"),
		fmt.Sprintf("%s", "1000"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, adr0),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	})
	err = cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// check
	require.Equal(t, 1, len(result.msgs))
	msg, ok := result.msgs[0].(*types.MsgCancelRedelegationNFT)
	require.True(t, ok)
	require.Equal(t, valSrcAddr.String(), msg.ValidatorSrc)
	require.Equal(t, valDstAddr.String(), msg.ValidatorDst)
	require.Equal(t, adr0.String(), msg.Delegator)
	require.Equal(t, "nft", msg.TokenID)
	require.Equal(t, []uint32{1, 2, 3}, msg.SubTokenIDs)
	require.Equal(t, int64(1000), msg.CreationHeight)
}
