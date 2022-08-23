package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestEcrecover(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.Bech32PrefixAccAddr, config.Bech32PrefixAccPub)
	_config.SetBech32PrefixForValidator(config.Bech32PrefixValAddr, config.Bech32PrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.Bech32PrefixConsAddr, config.Bech32PrefixConsPub)

	_r, err := hex.DecodeString("d8c0c8ff4a9b168be168f480bae61ead0a7f2b973f983a038f867621451fa553")
	require.NoError(t, err)

	var r Hash
	copy(r[:], _r)

	_s, err := hex.DecodeString("641ba9f5749afbb425e83b69ecacb3a0c6e32e2431609d474d4300a7cce5eb41")
	require.NoError(t, err)

	var s Hash
	copy(s[:], _s)

	sender, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	recipient, err := sdk.AccAddressFromBech32("dx1tlhpwr6t9nnq95xjet3ap2lc9zlxyw9dhr9y0z")
	require.NoError(t, err)

	amount, ok := sdk.NewIntFromString("100000000000000000000")
	require.True(t, ok)

	msg := MsgSwapRedeem{
		sender.String(),
		sender.String(),
		recipient.String(),
		amount,
		"del",
		"123",
		2,
		1,
		27,
		&r,
		&s}

	transactionNumber, ok := sdk.NewIntFromString(msg.TransactionNumber)
	require.True(t, ok)

	hash, err := GetHash(transactionNumber, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	require.NoError(t, err)

	require.Equal(t, "b3d218b80efdaaac18e3df1647786f1200fb330cf90bfef72baa0073f6bf872b", hex.EncodeToString(hash[:]))

	R := big.NewInt(0)
	R.SetBytes(msg.R[:])

	S := big.NewInt(0)
	S.SetBytes(msg.S[:])

	type args struct {
		sighash [32]byte
		R       *big.Int
		S       *big.Int
		Vb      *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    ethcmn.Address
		wantErr bool
	}{
		{
			"Test1",
			args{
				sighash: hash,
				R:       R,
				S:       S,
				Vb:      sdk.NewInt(int64(msg.V)).BigInt(),
			},
			ethcmn.HexToAddress(DefaultCheckingAddress),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ecrecover(tt.args.sighash, tt.args.R, tt.args.S, tt.args.Vb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecrecover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if hex.EncodeToString(got.Bytes()) != DefaultCheckingAddress {
				t.Errorf("Ecrecover() got = %v, want %v", hex.EncodeToString(got.Bytes()), DefaultCheckingAddress)
			}
		})
	}
}
