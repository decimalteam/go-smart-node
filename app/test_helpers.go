package app

import (
	dsctestutil "bitbucket.org/decimalteam/go-smart-node/testutil/dsc"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"bytes"
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v5/testing/simapp"

	"github.com/evmos/ethermint/encoding"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

func init() {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)
	cmdcfg.SetBip44CoinType(cfg)
}

var chainID = MainnetChainIDPrefix + "-1"

func setup(withGenesis bool, invCheckPeriod uint) (*DSC, GenesisState) {
	db := dbm.NewMemDB()

	app := NewDSC(log.NewNopLogger(), db, nil, true,
		map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	if withGenesis {
		return app, NewDefaultGenesisState(app.AppCodec())
	}

	return app, GenesisState{}
}

// Setup initializes a new DSC. A Nop logger is set in DSC.
func Setup(
	t *testing.T,
	isCheckTx bool,
	feemarketGenesis *feemarkettypes.GenesisState,
) *DSC {
	t.Helper()

	// create validators set
	valSet := dsctestutil.NewValidatorSet(t, 1)

	// generate genesis account
	acc, balance := dsctestutil.NewAcc(t)

	return SetupWithGenesisValSet(t, valSet, feemarketGenesis, []authtypes.GenesisAccount{acc}, balance)
}

func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, feemarketGenesis *feemarkettypes.GenesisState, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *DSC {
	t.Helper()

	app, genesisState := setup(true, 5)
	// Verify feeMarket genesis
	if feemarketGenesis != nil {
		require.NoError(t, feemarketGenesis.Validate())
		genesisState[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(feemarketGenesis)
	}

	genesisState, err := dsctestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         chainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: dsctestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		ChainID:            chainID,
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app
}

func SetupWithGenesisAccounts(t *testing.T, feemarketGenesis *feemarkettypes.GenesisState, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *DSC {
	// create validators set
	valSet := dsctestutil.NewValidatorSet(t, 1)

	return SetupWithGenesisValSet(t, valSet, feemarketGenesis, genAccs, balances...)
}

// GenesisStateWithSingleValidator initializes GenesisState with a single validator and genesis accounts
// that also act as delegators.
func GenesisStateWithSingleValidator(t *testing.T, app *DSC) GenesisState {
	t.Helper()

	// create validators set
	valSet := dsctestutil.NewValidatorSet(t, 1)

	// generate genesis account
	acc, balance := dsctestutil.NewAcc(t)

	genesisState := NewDefaultGenesisState(app.appCodec)
	genesisState, err := dsctestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)
	require.NoError(t, err)

	return genesisState
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *DSC, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, dsctestutil.CreateRandomAccounts)
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *DSC, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, dsctestutil.CreateIncrementalAccounts)
}

func GetAddrs(dsc *DSC, ctx sdk.Context, number int) []sdk.AccAddress {
	addrs := AddTestAddrsIncremental(dsc, ctx, number, sdk.Coins{
		{
			Denom:  "del",
			Amount: helpers.EtherToWei(sdk.NewInt(1000000000000)),
		},
	})

	return addrs
}

func addTestAddrs(app *DSC, ctx sdk.Context, accNum int, accCoins sdk.Coins, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, accCoins)
	}

	return testAddrs
}

func initAccountWithCoins(app *DSC, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	ok, failedDenom := checkCoinsExist(ctx, &app.CoinKeeper, coins)
	if !ok {
		panic(fmt.Sprintf("coin with denom '%s' not exist", *failedDenom))
	}

	err := app.BankKeeper.MintCoins(ctx, cointypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
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

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
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

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
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

func CreateAddr() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address())
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(codec codec.Codec, genesisState map[string]json.RawMessage,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}
