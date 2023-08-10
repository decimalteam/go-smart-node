package drc20cosmos

import (
	"bytes"
	"embed"
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

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

type Drc20Cosmos struct {
	abi.ABI
	ctx        sdk.Context
	evmKeeper  ethante.EVMKeeper
	bankKeeper bankkeeper.Keeper
	stateDB    vm.StateDB
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

	if drc.coin.Erc20Address == "" {
		drc.ctx.Logger().Info(drc.coin.Title)
	}
	drc.ctx.Logger().Info(drc.coin.Erc20Address)
	drc.ctx.Logger().Info(drc.coin.Denom)

	return true, nil
}
