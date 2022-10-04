package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ethereumCommon "github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

func init() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("dx", "dxpub")
}

func TestIsSupportedKeys(t *testing.T) {
	testCases := []struct {
		name        string
		pk          cryptotypes.PubKey
		isSupported bool
	}{
		{
			"nil key",
			nil,
			false,
		},
		{
			"ed25519 key",
			&ed25519.PubKey{},
			true,
		},
		{
			"secp256k1 key",
			&secp256k1.PubKey{},
			false,
		},
		{
			"ethsecp256k1 key",
			&ethsecp256k1.PubKey{},
			true,
		},
		{
			"multisig key - no pubkeys",
			&multisig.LegacyAminoPubKey{},
			false,
		},
		{
			"multisig key - valid pubkeys",
			multisig.NewLegacyAminoPubKey(2, []cryptotypes.PubKey{&ed25519.PubKey{}, &ed25519.PubKey{}, &ed25519.PubKey{}}),
			true,
		},
		{
			"multisig key - nested multisig",
			multisig.NewLegacyAminoPubKey(2, []cryptotypes.PubKey{&ed25519.PubKey{}, &ed25519.PubKey{}, &multisig.LegacyAminoPubKey{}}),
			false,
		},
		{
			"multisig key - invalid pubkey",
			multisig.NewLegacyAminoPubKey(2, []cryptotypes.PubKey{&ed25519.PubKey{}, &ed25519.PubKey{}, &secp256k1.PubKey{}}),
			false,
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.isSupported, IsSupportedKey(tc.pk), tc.name)
	}
}

func TestGetDecimalAddressFromBech32(t *testing.T) {
	testCases := []struct {
		name       string
		address    string
		expAddress string
		expError   bool
	}{
		{
			"blank bech32 address",
			" ",
			"",
			true,
		},
		{
			"invalid bech32 address",
			"dx",
			"",
			true,
		},
		{
			"invalid address bytes",
			"dx1123",
			"",
			true,
		},
		{
			"decimal address",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			false,
		},
		{
			"cosmos address",
			"cosmos1qql8ag4cluz6r4dz28p3w00dnc9w8ueulg2gmc",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			false,
		},
		{
			"osmosis address",
			"osmo1qql8ag4cluz6r4dz28p3w00dnc9w8ueuhnecd2",
			"dx1qql8ag4cluz6r4dz28p3w00dnc9w8ueue3x6fd",
			false,
		},
	}

	for _, tc := range testCases {
		addr, err := GetDecimalAddressFromBech32(tc.address)
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.expAddress, addr.String(), tc.name)
		}
	}
}

func TestOldNewAddresForPubKey(t *testing.T) {
	// Generate addresses from cli to check test data
	// deccli --keyring-backend test keys add 111 --dry-run -i --recover
	// dscd --keyring-backend test keys add 111 --dry-run --recover
	// pubKeyBytes taken by using decimal-go-sdk:
	/*
		acc1, _ := wallet.NewAccountFromMnemonicWords(mnemonic, "")
		bt := acc1.PrivateKey().PublicKey().BytesCompressed()
		fmt.Printf("mnemonic pub key = %#v\n", bt)
	*/
	var testData = []struct {
		pubKeyBytes []byte
		oldAddress  string
		newAddress  string
	}{
		{
			// gasp history river forget aware wide dance velvet weather rain rail dry cliff assault coach jelly choose spirit shoulder isolate kidney outer trust message
			[]byte{0x3, 0x44, 0x8e, 0x6b, 0x3d, 0x50, 0xd6, 0xa3, 0x9c, 0xab, 0x3b, 0xab, 0xaa, 0x4a, 0xa2, 0xb0, 0x88, 0x5f, 0x55, 0x6f, 0xe0, 0x5d, 0x71, 0x49, 0x88, 0x5a, 0x5, 0xa0, 0xe7, 0x94, 0xa, 0x7e, 0x4f},
			"dx1w98j4vk6dkpyndjnv5dn2eemesq6a2c2j9depy",
			"dx1xp6aqad49te7vsfga6str8hrdeh24r9jnxuadn",
		},
		{
			// section jeans evoke hockey result spell dish zero merge actress pink resource loan afford fitness install purity duck cannon ugly session stereo pattern spawn
			[]byte{0x3, 0x16, 0x18, 0x96, 0x7d, 0x77, 0xf1, 0xe4, 0x90, 0xd4, 0x1f, 0xc0, 0xe0, 0xc0, 0xc8, 0xb4, 0xb0, 0x47, 0x99, 0xe1, 0x16, 0x82, 0x4b, 0xea, 0x8, 0x47, 0x33, 0xe4, 0x63, 0x43, 0x20, 0xca, 0x3},
			"dx1m3eg7v6pu0dga2knj9zm4683dk9c8800j9nfw0",
			"dx18c8mer8lq2y8yw8cq8f4c6fdqfa8xcjg3pv33f",
		},
		{
			// citizen marine borrow just apology mistake trumpet border sauce drip smile current excuse sing shove puppy dial ticket margin fabric afraid identify rookie elite
			[]byte{0x3, 0x9a, 0x1, 0x12, 0x1d, 0x8f, 0xc3, 0xc5, 0x9e, 0xe7, 0xcb, 0x7, 0xbe, 0x27, 0x68, 0x8d, 0x75, 0x23, 0xd4, 0xb2, 0xb0, 0xbc, 0xf2, 0x4e, 0x83, 0x14, 0xbb, 0x27, 0xc8, 0xe7, 0xa8, 0x6f, 0x96},
			"dx1lw2q66zph22x3hzmc527em25kd4zfydnx7arw7",
			"dx164ea54aqgsmp7dp6wzs0y8n6vjehudnkgcn4fe",
		},
	}
	for _, tc := range testData {
		newPubKey := ethsecp256k1.PubKey{Key: tc.pubKeyBytes}
		oldPubKey := secp256k1.PubKey{Key: tc.pubKeyBytes}
		newAddress, err := bech32.ConvertAndEncode("dx", newPubKey.Address())
		require.NoError(t, err)
		require.Equal(t, tc.newAddress, newAddress)
		oldAddress, err := GetLegacyAddressFromPubKey(oldPubKey.Key)
		require.NoError(t, err)
		require.Equal(t, tc.oldAddress, oldAddress)
	}
}

// this test check conversion from ethereum address representation to bech32 representation
func TestEth2dx(t *testing.T) {
	ethAdr := "0xa618f8e2b953593c1f08f2b3dce2a963ce130916"
	dxAdr := "dx15cv03c4e2dvnc8cg72eaec4fv08pxzgkmr255d"
	adr := ethereumCommon.HexToAddress(ethAdr)
	bch, err := sdk.Bech32ifyAddressBytes("dx", adr.Bytes())
	require.NoError(t, err)
	require.Equal(t, dxAdr, bch)
}
