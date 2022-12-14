package dsc

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	//tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	coinkeeper "bitbucket.org/decimalteam/go-smart-node/x/coin/keeper"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
)

type GenerateAccountStrategy func(int) []sdk.AccAddress

func NewValidatorSet(t *testing.T, valCount int) *tmtypes.ValidatorSet {
	validators := make([]*tmtypes.Validator, valCount)

	for i := 0; i < valCount; i++ {
		privVal := mock.NewPV()
		pubKey, err := privVal.GetPubKey()
		require.NoError(t, err, "fail to get pubkey from privkey")

		// create validator set with single validator
		validators[i] = tmtypes.NewValidator(pubKey, 1)
	}

	valSet := tmtypes.NewValidatorSet(validators)

	return valSet
}

func NewAcc(t *testing.T) (*authtypes.BaseAccount, banktypes.Balance) {
	senderPrivKey, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err, "fail to create new genesis account")
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000000000000))),
	}

	return acc, balance
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(bankKeeper bankkeeper.Keeper, stakingKeeper *stakingkeeper.Keeper, coinKeeper *coinkeeper.Keeper, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(bankKeeper, stakingKeeper, coinKeeper, ctx, accNum, accAmt, CreateRandomAccounts)
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(bankKeeper bankkeeper.Keeper, stakingKeeper *stakingkeeper.Keeper, coinKeeper *coinkeeper.Keeper, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(bankKeeper, stakingKeeper, coinKeeper, ctx, accNum, accAmt, CreateIncrementalAccounts)
}

func GetAddrs(bankKeeper bankkeeper.Keeper, stakingKeeper *stakingkeeper.Keeper, coinKeeper *coinkeeper.Keeper, ctx sdk.Context, number int) []sdk.AccAddress {
	addrs := AddTestAddrsIncremental(bankKeeper, stakingKeeper, coinKeeper, ctx, number, sdk.Coins{
		{
			Denom:  "del",
			Amount: helpers.EtherToWei(sdk.NewInt(1000000000000)),
		},
	})

	return addrs
}

func addTestAddrs(bankKeeper bankkeeper.Keeper, stakingKeeper *stakingkeeper.Keeper, coinKeeper *coinkeeper.Keeper, ctx sdk.Context, accNum int, accCoins sdk.Coins, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	for _, addr := range testAddrs {
		initAccountWithCoins(bankKeeper, stakingKeeper, coinKeeper, ctx, addr, accCoins)
	}

	return testAddrs
}

func initAccountWithCoins(bankKeeper bankkeeper.Keeper, stakingKeeper *stakingkeeper.Keeper, coinKeeper *coinkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	ok, failedDenom := checkCoinsExist(ctx, coinKeeper, coins)
	if !ok {
		panic(fmt.Sprintf("coin with denom '%s' not exist", *failedDenom))
	}

	err := bankKeeper.MintCoins(ctx, cointypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

func checkCoinsExist(ctx sdk.Context, k cointypes.CoinKeeper, coins sdk.Coins) (bool, *string) {
	for _, coin := range coins {
		_, err := k.GetCoin(ctx, coin.Denom)
		if err != nil {
			return false, &coin.Denom
		}
	}

	return true, nil
}

// CreateRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func CreateRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// CreateIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func CreateIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

func CreateAddr() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address())
}
