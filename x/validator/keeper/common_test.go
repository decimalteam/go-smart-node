package keeper_test

import (
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/decimalteam/ethermint/server/config"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	feemarkettypes "github.com/decimalteam/ethermint/x/feemarket/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	"bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

var PKs = simapp.CreateTestPubKeys(500)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// createTestInput Returns a app.DSC with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput(t *testing.T) (*codec.LegacyAmino, *app.DSC, sdk.Context) {
	dsc := app.Setup(t, false, feemarkettypes.DefaultGenesisState())
	ctx := dsc.BaseApp.NewContext(false, tmproto.Header{})

	dsc.ValidatorKeeper = keeper.NewKeeper(
		dsc.AppCodec(),
		dsc.GetKey(types.StoreKey),
		dsc.GetSubspace(types.ModuleName),
		dsc.AccountKeeper,
		dsc.BankKeeper,
		&dsc.NFTKeeper,
		&dsc.CoinKeeper,
		&dsc.MultisigKeeper,
		&dsc.EvmKeeper,
	)
	return dsc.LegacyAmino(), dsc, ctx
}

// intended to be used with require/assert:  require.True(ValEq(...))
func ValEq(t *testing.T, exp, got types.Validator) (*testing.T, bool, string, types.Validator, types.Validator) {
	return t, exp.MinEqual(&got), "expected:\n%v\ngot:\n%v", exp, got
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(dsc *app.DSC, ctx sdk.Context, numAddrs int, accCoins sdk.Coins) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := app.AddTestAddrsIncremental(dsc, ctx, numAddrs, accCoins)
	addrVals := app.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

// DeployTestContract deploy a test erc20 contract and returns the contract address
func DeployTestContractCenter(t *testing.T, app *app.DSC, ctx sdk.Context, addressOwner common.Address, keyringSigner keyring.Signer) common.Address {
	//ctx := sdk.WrapSDKContext(ctx)
	chainID := app.EvmKeeper.ChainID()

	//abiContract, _ := contracts.ContractCenterMetaData.GetAbi()

	nonce := app.EvmKeeper.GetNonce(ctx, addressOwner)

	data := []byte(contracts.ContractCenterBin)
	args, err := json.Marshal(&evmtypes.TransactionArgs{
		From: &addressOwner,
		Data: (*hexutil.Bytes)(&data),
	})
	require.NoError(t, err)
	res, err := app.EvmKeeper.EstimateGas(ctx, &evmtypes.EthCallRequest{
		Args:            args,
		GasCap:          uint64(config.DefaultGasCap),
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	})
	require.NoError(t, err)

	var erc20DeployTx *evmtypes.MsgEthereumTx
	erc20DeployTx = evmtypes.NewTxContract(
		chainID,
		nonce,
		nil,     // amount
		res.Gas, // gasLimit
		nil,     // gasPrice
		app.FeeKeeper.GetBaseFee(ctx),
		big.NewInt(1),
		data,                   // input
		&ethtypes.AccessList{}, // accesses
	)

	erc20DeployTx.From = addressOwner.Hex()
	err = erc20DeployTx.Sign(ethtypes.LatestSignerForChainID(chainID), keyringSigner)
	require.NoError(t, err)
	rsp, err := app.EvmKeeper.EthereumTx(ctx, erc20DeployTx)
	require.NoError(t, err)
	require.Empty(t, rsp.VmError)
	return crypto.CreateAddress(addressOwner, nonce)
}
