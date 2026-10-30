package keeper

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegation"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationNft"
	"bitbucket.org/decimalteam/go-smart-node/contracts/delegationv7"
)

// Real mainnet logs emitted by the deployed (7-field Stake) delegation contract.
const mainnetDelegation = "0xa16c34ed1c0601c0e749e17ebef19752a15faa01"

func mainnetLog(t *testing.T, topic0, data string) *ethtypes.Log {
	t.Helper()
	return &ethtypes.Log{
		Address: common.HexToAddress(mainnetDelegation),
		Topics:  []common.Hash{common.HexToHash(topic0)},
		Data:    hexutil.MustDecode(data),
	}
}

var (
	// tx 0x57af0d5d… (block 33635482): 793 POLETELI hold stake by 0xd7da13… on 0x8ff2f2…
	logStakeUpdated = []string{"0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b",
		"0x0fd3239ecc6601c2e0f036869779f1323bc1b160f0a47fdfe203bd92eb15a8180000000000000000000000008ff2f220fb80b3f26cd96652aecbf3129983c115000000000000000000000000d7da13fb8ecf2e56019134027fb6aa5b4beafac1000000000000000000000000cf91ae318d8d48b0a8753b130e3e3eb2f7e3fd1900000000000000000000000000000000000000000000002afd1cbb9ca3e09f7700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000006c8fbb390000000000000000000000000000000000000000000000000000000000000001"}
	logStakeAmountUpdated = []string{"0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959",
		"0x0fd3239ecc6601c2e0f036869779f1323bc1b160f0a47fdfe203bd92eb15a81800000000000000000000000000000000000000000000002afd1cbb9ca3e09f77"}
	logWithdrawRequest = []string{"0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab",
		"0xfd43abc7fd14f5052b26a8a990838d54afc7072bca3ce590ab6e59bb508cb255000000000000000000000000000000000000000000000000000000000000449c0000000000000000000000009199bc585650bc60356f34d6e3cdc6aeeb4fa0c30000000000000000000000005479581e3004207a10bca590438294f3e7cb0b82000000000000000000000000ec2991de234a010fc5b58842d594fe9ae08d730400000000000000000000000000000000000000000000269fec7f0361d200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000006a75df77"}
	logStakeHolded = []string{"0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb",
		"0x8c0f8f1514e8f21b62b447693c1ad4ea142068e2f814db7bb23a82e4132bc09d0000000000000000000000000cd5f3b33e1617dddacef0c3be7559df7913f69c00000000000000000000000008b172e4fb254a4f1e6ce76762dde1b772a5cdf400000000000000000000000060cf963b60c8f61a45e1a1c2b6b50ba582c39880000000000000000000000000000000000000000000000a968163f0a57b40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000006c73c3e20000000000000000000000000000000000000000000000000000000000000001"}
	logTransferRequest = []string{"0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f",
		"0x22f643fecb41187c2ed61431d253bed5ef5707d1d7a6168c3097a560958076f50000000000000000000000000000000000000000000000000000000000004947000000000000000000000000cad64f9d9ebbd423a9e3d78c36afcf5cfa7fd565000000000000000000000000cb2ee62b7e316b0e0544658226f50a446125c398000000000000000000000000c5e341e2d1fd835e097e9433929e66cc7a37915b0000000000000000000000000000000000000000000003a99480e392c57a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000006a761665"}
	logStakeReset = []string{"0x0a825640f097d14edca3b473a8fa6e574cd5f56933823f4ba0114a2bce58c181",
		"0xd7ff2af0529a5cbc19f710294f803bce29110d9a8fb416699d83af2d6173b7510000000000000000000000008ff2f220fb80b3f26cd96652aecbf3129983c11500000000000000000000000017ebb1b504ecf2314d71a0c4d615b5b5793964b0000000000000000000000000a315eae6edac19c74f132683214176f0ab784ef3000000000000000000000000000000000000000000009a8fe2ae6f630800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000006a77ecb4"}
	logWithdrawCompleted = []string{"0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7",
		"0x00000000000000000000000000000000000000000000000000000000000041da"}
)

func TestNormalizeDelegationLogs_TranslatesEveryLegacyStakeEvent(t *testing.T) {
	current, err := delegation.DelegationMetaData.GetAbi()
	require.NoError(t, err)
	legacy, err := delegationv7.DelegationMetaData.GetAbi()
	require.NoError(t, err)

	for name, raw := range map[string][]string{
		"StakeUpdated": logStakeUpdated, "StakeHolded": logStakeHolded, "StakeReset": logStakeReset,
		"WithdrawRequest": logWithdrawRequest, "TransferRequest": logTransferRequest,
	} {
		t.Run(name, func(t *testing.T) {
			in := mainnetLog(t, raw[0], raw[1])
			_, err := current.EventByID(in.Topics[0])
			require.Error(t, err, "the current ABI must not know the legacy topic (that is the bug)")

			out, err := normalizeDelegationLogs([]*ethtypes.Log{in}, mainnetDelegation, current, legacy)
			require.NoError(t, err)
			require.Len(t, out, 1)
			ev, err := current.EventByID(out[0].Topics[0])
			require.NoError(t, err)
			require.Equal(t, name, ev.Name)

			// The translated stake must equal the legacy one field by field.
			var want delegationv7.IDecimalDelegationCommonStake
			var got delegation.IDecimalDelegationCommonStake
			switch name {
			case "WithdrawRequest", "TransferRequest":
				var l delegationv7.DelegationWithdrawRequest
				var c delegation.DelegationWithdrawRequest
				require.NoError(t, contracts.UnpackLog(legacy, &l, name, in))
				require.NoError(t, contracts.UnpackLog(current, &c, name, out[0]))
				require.Equal(t, l.StakeId, c.StakeId)
				require.Equal(t, 0, l.StakeIndex.Cmp(c.StakeIndex))
				require.Equal(t, l.FrozenStake.FreezeStatus, c.FrozenStake.FreezeStatus)
				require.Equal(t, l.FrozenStake.FreezeType, c.FrozenStake.FreezeType)
				require.Equal(t, 0, l.FrozenStake.UnfreezeTimestamp.Cmp(c.FrozenStake.UnfreezeTimestamp))
				want, got = l.FrozenStake.Stake, c.FrozenStake.Stake
			case "StakeReset":
				var l delegationv7.DelegationStakeReset
				var c delegation.DelegationStakeReset
				require.NoError(t, contracts.UnpackLog(legacy, &l, name, in))
				require.NoError(t, contracts.UnpackLog(current, &c, name, out[0]))
				require.Equal(t, l.StakeId, c.StakeId)
				want, got = l.Stake, c.Stake
			default:
				var l delegationv7.DelegationStakeUpdated
				var c delegation.DelegationStakeUpdated
				require.NoError(t, contracts.UnpackLog(legacy, &l, name, in))
				require.NoError(t, contracts.UnpackLog(current, &c, name, out[0]))
				require.Equal(t, l.StakeId, c.StakeId)
				require.Equal(t, l.IsNew, c.IsNew)
				want, got = l.Stake, c.Stake
			}
			require.Equal(t, want.Validator, got.Validator)
			require.Equal(t, want.Delegator, got.Delegator)
			require.Equal(t, want.Token, got.Token)
			require.Equal(t, 0, want.Amount.Cmp(got.Amount))
			require.Equal(t, 0, want.TokenId.Cmp(got.TokenId))
			require.Equal(t, want.TokenType, got.TokenType)
			require.Equal(t, 0, want.HoldTimestamp.Cmp(got.HoldTimestamp))
			require.Equal(t, 0, got.HoldStartTime.Sign(), "legacy events carry no hold start; 0 makes the node use block time")
		})
	}
}

func TestNormalizeDelegationLogs_PoletelStake(t *testing.T) {
	current, _ := delegation.DelegationMetaData.GetAbi()
	legacy, _ := delegationv7.DelegationMetaData.GetAbi()
	out, err := normalizeDelegationLogs([]*ethtypes.Log{
		mainnetLog(t, logStakeAmountUpdated[0], logStakeAmountUpdated[1]),
		mainnetLog(t, logStakeUpdated[0], logStakeUpdated[1]),
	}, mainnetDelegation, current, legacy)
	require.NoError(t, err)
	var su delegation.DelegationStakeUpdated
	require.NoError(t, contracts.UnpackLog(current, &su, "StakeUpdated", out[1]))
	require.Equal(t, common.HexToAddress("0x8ff2f220fb80b3f26cd96652aecbf3129983c115"), su.Stake.Validator)
	require.Equal(t, common.HexToAddress("0xd7da13Fb8ECF2E56019134027fB6Aa5b4beAFAC1"), su.Stake.Delegator)
	require.Equal(t, common.HexToAddress("0xcf91ae318d8d48b0a8753b130e3e3eb2f7e3fd19"), su.Stake.Token)
	amount, _ := new(big.Int).SetString("793001909968183533431", 10)
	require.Equal(t, 0, amount.Cmp(su.Stake.Amount))
	require.Equal(t, int64(1821358905), su.Stake.HoldTimestamp.Int64())
	require.True(t, su.IsNew)
}

func TestNormalizeDelegationLogs_PassThroughAndStrictness(t *testing.T) {
	current, _ := delegation.DelegationMetaData.GetAbi()
	legacy, _ := delegationv7.DelegationMetaData.GetAbi()

	amountLog := mainnetLog(t, logStakeAmountUpdated[0], logStakeAmountUpdated[1])
	completed := mainnetLog(t, logWithdrawCompleted[0], logWithdrawCompleted[1])
	foreign := &ethtypes.Log{Address: common.HexToAddress("0xcf91ae318d8d48b0a8753b130e3e3eb2f7e3fd19"),
		Topics: []common.Hash{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}, Data: []byte{1}}
	legacyLog := mainnetLog(t, logStakeUpdated[0], logStakeUpdated[1])
	origData := append([]byte(nil), legacyLog.Data...)
	origTopic := legacyLog.Topics[0]

	in := []*ethtypes.Log{foreign, amountLog, completed, legacyLog}
	out, err := normalizeDelegationLogs(in, mainnetDelegation, current, legacy)
	require.NoError(t, err)
	require.Len(t, out, 4)
	require.Same(t, foreign, out[0], "other contracts' logs are untouched")
	require.Same(t, amountLog, out[1], "layout-independent events are untouched")
	require.Same(t, completed, out[2])
	require.NotSame(t, legacyLog, out[3])
	// the receipt's own log is never mutated (it is part of the tx result)
	require.True(t, bytes.Equal(origData, legacyLog.Data))
	require.Equal(t, origTopic, legacyLog.Topics[0])

	// An event from the delegation contract that neither layout knows must fail.
	unknown := mainnetLog(t, "0x1111111111111111111111111111111111111111111111111111111111111111", "0x")
	_, err = normalizeDelegationLogs([]*ethtypes.Log{unknown}, mainnetDelegation, current, legacy)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "unknown to the node"))

	// The same unknown topic from another contract is not the node's business.
	unknown.Address = common.HexToAddress("0x0000000000000000000000000000000000000001")
	_, err = normalizeDelegationLogs([]*ethtypes.Log{unknown}, mainnetDelegation, current, legacy)
	require.NoError(t, err)
}

func TestIsLegacyNftDeposit(t *testing.T) {
	nftAbi, err := delegationNft.DelegationNftMetaData.GetAbi()
	require.NoError(t, err)
	const legacyNft = "0xe45adfcc739a0d10ce9462b58866c9a1a06035e2"
	at := func(addr string, name string) *ethtypes.Log {
		return &ethtypes.Log{Address: common.HexToAddress(addr), Topics: []common.Hash{nftAbi.Events[name].ID}}
	}
	for _, name := range []string{"StakeUpdated", "StakeHolded"} {
		require.True(t, isLegacyNftDeposit([]*ethtypes.Log{at(legacyNft, name)}, legacyNft, nftAbi), name)
		require.False(t, isLegacyNftDeposit([]*ethtypes.Log{at(mainnetDelegation, name)}, legacyNft, nftAbi), "other address")
	}
	for _, name := range []string{"WithdrawRequest", "TransferRequest", "WithdrawCompleted", "TransferCompleted"} {
		require.False(t, isLegacyNftDeposit([]*ethtypes.Log{at(legacyNft, name)}, legacyNft, nftAbi), "exits stay allowed: "+name)
		// as emitted: the exit's own stake update comes with it (before or after)
		require.False(t, isLegacyNftDeposit([]*ethtypes.Log{at(legacyNft, "StakeAmountUpdated"), at(legacyNft, "StakeUpdated"), at(legacyNft, name)}, legacyNft, nftAbi), "exit with its stake update: "+name)
		require.False(t, isLegacyNftDeposit([]*ethtypes.Log{at(legacyNft, name), at(legacyNft, "StakeUpdated")}, legacyNft, nftAbi), "completion with re-stake: "+name)
	}
	// an exit of another contract does not excuse a legacy deposit
	require.True(t, isLegacyNftDeposit([]*ethtypes.Log{at(mainnetDelegation, "WithdrawRequest"), at(legacyNft, "StakeUpdated")}, legacyNft, nftAbi))
}
