package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/cosmos/ibc-go/v3/testing/simapp"

	"github.com/tharsis/ethermint/encoding"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

// DefaultTestingAppInit defines the IBC application used for testing
var DefaultTestingAppInit func() (ibctesting.TestingApp, map[string]json.RawMessage) = SetupTestingApp

// DefaultConsensusParams defines the default Tendermint consensus params used in DSC testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func init() {
	cfg := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(cfg)
	cmdcfg.SetBip44CoinType(cfg)
}

// Setup initializes a new DSC. A Nop logger is set in DSC.
func Setup(
	isCheckTx bool,
	feemarketGenesis *feemarkettypes.GenesisState,
) *DSC {
	db := dbm.NewMemDB()
	app := NewDSC(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()

		// Verify feeMarket genesis
		if feemarketGenesis != nil {
			if err := feemarketGenesis.Validate(); err != nil {
				panic(err)
			}
			genesisState[feemarkettypes.ModuleName] = app.AppCodec().MustMarshalJSON(feemarketGenesis)
		}

		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				ChainId:         MainnetChainIDPrefix + "-1",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// SetupTestingApp initializes the IBC-go testing application
func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	cfg := encoding.MakeConfig(ModuleBasics)
	app := NewDSC(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, cfg, simapp.EmptyAppOptions{})
	return app, NewDefaultGenesisState()
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *DSC, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *DSC, ctx sdk.Context, accNum int, accAmt sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
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
		res, _ := sdk.AccAddressFromHex(buffer.String())
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
	res, err := sdk.AccAddressFromHex(addr)
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
