package drc20cosmos

import (
	"bytes"
	"embed"
	"encoding/json"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"

	"bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethante "github.com/evmos/ethermint/app/ante"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	evm "github.com/evmos/ethermint/x/evm/vm"
)

const (
	addressForContractOwner = "0x2941b073ad6b59b1de4fc70c69e39a9e130b25ce"
	firstReward             = 50

	firstIncrease = 5
)

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

type Drc20Cosmos struct {
	abi.ABI
	ctx        sdk.Context
	evmKeeper  ethante.EVMKeeper
	bankKeeper bankkeeper.Keeper
	stateDB    *statedb.StateDB
	evm        evm.EVM
	coin       types.Coin
}

// NewDrc20Cosmos create instance of contract
func NewDrc20Cosmos(ctx sdk.Context,
	evmKeeper ethante.EVMKeeper,
	bankKeeper bankkeeper.Keeper,
	msgEthTx *evmtypes.MsgEthereumTx,
	coinAction types.Coin,
) (*Drc20Cosmos, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, err
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return nil, err
	}

	params := evmKeeper.GetParams(ctx)
	ethCfg := params.ChainConfig.EthereumConfig(evmKeeper.ChainID())
	baseFee := evmKeeper.GetBaseFee(ctx, ethCfg)

	signer := ethtypes.MakeSigner(ethCfg, big.NewInt(ctx.BlockHeight()))

	coreMsg, err := msgEthTx.AsMessage(signer, baseFee)
	if err != nil {
		return nil, err
	}

	cfg := &evmtypes.EVMConfig{
		ChainConfig: ethCfg,
		Params:      params,
		CoinBase:    common.Address{},
		BaseFee:     baseFee,
	}

	stateNewDB := statedb.New(ctx, evmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(ctx.HeaderHash().Bytes())))
	evmNew := evmKeeper.NewEVM(ctx, coreMsg, cfg, evmtypes.NewNoOpTracer(), stateNewDB)

	// nonce := stateDB.GetNonce(common.HexToAddress("0x2941b073ad6b59b1de4fc70c69e39a9e130b25ce"))

	// stateDB.SetNonce(common.HexToAddress("0x2941b073ad6b59b1de4fc70c69e39a9e130b25ce"), nonce)
	// ret, _, leftoverGas, vmErr = evm.Create(sender, msg.Data(), leftoverGas, msg.Value())
	// stateDB.SetNonce(sender.Address(), msg.Nonce()+1)

	return &Drc20Cosmos{
		ABI:        newAbi,
		ctx:        ctx,
		evmKeeper:  evmKeeper,
		bankKeeper: bankKeeper,
		stateDB:    stateNewDB,
		evm:        evmNew,
		coin:       coinAction,
	}, nil
}

// CreateContractIfNotSet creation contract if not not to coin
func (drc Drc20Cosmos) CreateContractIfNotSet() (bool, error) {

	sender := vm.AccountRef(common.HexToAddress(addressForContractOwner))

	if drc.coin.Drc20Address == "" {
		drc.ctx.Logger().Info(drc.coin.Title)
	}
	drc.ctx.Logger().Info(drc.coin.Drc20Address)
	drc.ctx.Logger().Info(drc.coin.Denom)

	// receive nonce for owner address for new contract
	nonce := drc.stateDB.GetNonce(common.HexToAddress(addressForContractOwner))
	drc.stateDB.SetNonce(common.HexToAddress(addressForContractOwner), nonce)

	contractCode, err := f.ReadFile("creation.code")
	if err != nil {
		return false, err
	}

	ret, _, _, vmErr := drc.evm.Create(sender, contractCode, 10000, big.NewInt(100))
	drc.stateDB.SetNonce(sender.Address(), nonce+1)

	drc.ctx.Logger().With(ret).Info("Result create contract")

	if vmErr != nil {
		drc.ctx.Logger().Info(vmErr.Error())
		return false, sdkerrors.ErrUnknownRequest.Wrapf("failed to encode log %T", vmErr)
	}

	txLogAttrs := make([]sdk.Attribute, len(drc.stateDB.Logs()))
	for i, log := range drc.stateDB.Logs() {
		value, err := json.Marshal(log)
		if err != nil {
			drc.ctx.Logger().Info(drc.coin.Denom)
			return false, sdkerrors.ErrUnknownRequest.Wrapf("failed to encode log %T", err)
		}
		txLogAttrs[i] = sdk.NewAttribute(evmtypes.AttributeKeyTxLog, string(value))
	}

	// The dirty states in `StateDB` is either committed or discarded after return
	if err := drc.stateDB.Commit(); err != nil {
		return false, sdkerrors.ErrUnknownRequest.Wrapf("failed to encode log %T", err)
	}

	return true, nil
}
