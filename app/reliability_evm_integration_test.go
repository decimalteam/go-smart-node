package app

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v6/testing/simapp"
	"github.com/decimalteam/ethermint/encoding"
	ethermint "github.com/decimalteam/ethermint/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	dsctest "bitbucket.org/decimalteam/go-smart-node/testutil/dsc"
	cointypes "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	validatorkeeper "bitbucket.org/decimalteam/go-smart-node/x/validator/keeper"
	validatortypes "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
)

const reliabilityChainID = "decimal_75-1"

var reliabilityTime = time.Unix(1800000000, 0).UTC()

type reliabilityFixture struct {
	ABI      json.RawMessage `json:"abi"`
	Bytecode string          `json:"bytecode"`
	Runtime  string          `json:"runtime"`
}

func loadReliabilityFixture(t *testing.T, name string) (reliabilityFixture, abi.ABI) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "reliability", name+".json"))
	require.NoError(t, err)
	var f reliabilityFixture
	require.NoError(t, json.Unmarshal(b, &f))
	a, err := abi.JSON(bytes.NewReader(f.ABI))
	require.NoError(t, err)
	return f, a
}
func reliabilityActivate(t *testing.T, enabled bool) {
	t.Helper()
	old := validatorkeeper.ReliabilityFixActivationHeights
	validatorkeeper.ReliabilityFixActivationHeights = map[string]int64{}
	if enabled {
		validatorkeeper.ReliabilityFixActivationHeights[reliabilityChainID] = 1
	}
	t.Cleanup(func() { validatorkeeper.ReliabilityFixActivationHeights = old })
}
func reliabilityOpen(t *testing.T, home string) (*DSC, dbm.DB) {
	t.Helper()
	db, err := dbm.NewGoLevelDB("application", filepath.Join(home, "data"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	dsc := NewDSC(log.NewNopLogger(), db, nil, true, map[int64]bool{}, home, 0, encoding.MakeConfig(ModuleBasics), simapp.EmptyAppOptions{})
	return dsc, db
}
func reliabilityHeader(t *testing.T, home string, height int64) tmproto.Header {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(home, "proposer"))
	require.NoError(t, err)
	return tmproto.Header{ChainID: reliabilityChainID, Height: height, Time: reliabilityTime.Add(time.Duration(height) * time.Second), ProposerAddress: b}
}
func reliabilityContext(t *testing.T, dsc *DSC, home string, height int64) sdk.Context {
	h := reliabilityHeader(t, home, height)
	dsc.BeginBlock(abci.RequestBeginBlock{Header: h})
	return dsc.BaseApp.NewContext(false, h).WithGasMeter(sdk.NewInfiniteGasMeter()).WithBlockGasMeter(sdk.NewInfiniteGasMeter())
}
func reliabilityNew(t *testing.T) (*DSC, sdk.Context, string, dbm.DB) {
	t.Helper()
	home := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(home, "data"), 0700))
	require.NoError(t, os.MkdirAll(filepath.Join(home, "config"), 0700))
	g := tmtypes.GenesisDoc{ChainID: reliabilityChainID, GenesisTime: reliabilityTime, ConsensusParams: tmtypes.DefaultConsensusParams(), AppState: json.RawMessage(`{}`)}
	require.NoError(t, g.SaveAs(filepath.Join(home, "config", "genesis.json")))
	dsc, db := reliabilityOpen(t, home)
	vals := dsctest.NewValidatorSet(t, 1)
	acc, balance := dsctest.NewAcc(t)
	require.NoError(t, os.WriteFile(filepath.Join(home, "proposer"), vals.Validators[0].Address, 0600))
	gs, err := dsctest.GenesisStateWithValSet(dsc.AppCodec(), NewDefaultGenesisState(dsc.AppCodec()), vals, []authtypes.GenesisAccount{acc}, balance)
	require.NoError(t, err)
	b, err := json.Marshal(gs)
	require.NoError(t, err)
	dsc.InitChain(abci.RequestInitChain{ChainId: reliabilityChainID, ConsensusParams: dsctest.DefaultConsensusParams, AppStateBytes: b})
	dsc.Commit()
	ctx := reliabilityContext(t, dsc, home, dsc.LastBlockHeight()+1)
	params := dsc.EvmKeeper.GetParams(ctx)
	params.EvmDenom = cmdcfg.BaseDenom
	dsc.EvmKeeper.SetParams(ctx, params)
	dsc.AccountKeeper.GetModuleAccount(ctx, cointypes.ModuleName)
	dsc.AccountKeeper.GetModuleAccount(ctx, validatortypes.ModuleName)
	zero := sdk.AccAddress(make([]byte, 20))
	a := dsc.AccountKeeper.NewAccountWithAddress(ctx, zero)
	dsc.AccountKeeper.SetAccount(ctx, a)
	return dsc, ctx, home, db
}
func reliabilityDeploy(t *testing.T, dsc *DSC, ctx sdk.Context, name string) (common.Address, abi.ABI) {
	t.Helper()
	f, a := loadReliabilityFixture(t, name)
	code, err := hex.DecodeString(f.Bytecode)
	require.NoError(t, err)
	from := common.Address{}
	nonce, err := dsc.AccountKeeper.GetSequence(ctx, sdk.AccAddress(from.Bytes()))
	require.NoError(t, err)
	addr := crypto.CreateAddress(from, nonce)
	res, err := dsc.EvmKeeper.CallEVMWithData(ctx, from, nil, code, true)
	require.NoError(t, err)
	require.False(t, res.Failed())
	require.Equal(t, crypto.Keccak256Hash(common.FromHex(f.Runtime)), common.BytesToHash(dsc.EvmKeeper.GetAccount(ctx, addr).CodeHash))
	return addr, a
}
func reliabilityRegistry(t *testing.T, dsc *DSC, ctx sdk.Context, auto, validator common.Address) {
	t.Helper()
	f, a := loadReliabilityFixture(t, "RegistryFixture")
	code := common.FromHex(f.Runtime)
	addr := common.HexToAddress(contracts.GetContractCenter(reliabilityChainID))
	hash := crypto.Keccak256Hash(code)
	acct := dsc.AccountKeeper.NewAccountWithAddress(ctx, sdk.AccAddress(addr.Bytes()))
	require.NoError(t, acct.(ethermint.EthAccountI).SetCodeHash(hash))
	dsc.AccountKeeper.SetAccount(ctx, acct)
	dsc.EvmKeeper.SetCode(ctx, hash.Bytes(), code)
	for slug, target := range map[string]common.Address{"delegation": auto, "master-validator": validator, "token-center": common.HexToAddress("0x1111"), "nft-center": common.HexToAddress("0x2222"), "wdel": common.HexToAddress("0x3333")} {
		_, err := dsc.EvmKeeper.CallEVM(ctx, a, common.Address{}, addr, true, "set", slug, target)
		require.NoError(t, err)
	}
}
func reliabilitySlot(dsc *DSC, ctx sdk.Context, addr common.Address, slot int64) int64 {
	return dsc.EvmKeeper.GetState(ctx, addr, common.BigToHash(big.NewInt(slot))).Big().Int64()
}

func TestIntegrationUpgradeCrashRealEVM(t *testing.T) {
	if os.Getenv("RELIABILITY_CRASH_CHILD") == "1" {
		home := os.Getenv("RELIABILITY_HOME")
		dsc, _ := reliabilityOpen(t, home)
		_, a := loadReliabilityFixture(t, "ValidatorFixture")
		b, err := os.ReadFile(filepath.Join(home, "target"))
		require.NoError(t, err)
		target := common.BytesToAddress(b)
		dsc.UpgradeKeeper.SetUpgradeHandler("reliability-crash", func(c sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			_, err := dsc.EvmKeeper.CallEVM(c, a, common.Address{}, target, true, "bump")
			return vm, err
		})
		ctx := reliabilityContext(t, dsc, home, 3)
		require.Equal(t, int64(1), reliabilitySlot(dsc, ctx, target, 0))
		require.Equal(t, int64(3), dsc.UpgradeKeeper.GetDoneHeight(ctx, "reliability-crash"))
		// Deliberately skip EndBlock, Commit, DB.Close and deferred cleanup.
		os.Exit(83)
	}
	reliabilityActivate(t, true)
	dsc, ctx, home, db := reliabilityNew(t)
	target, a := reliabilityDeploy(t, dsc, ctx, "ValidatorFixture")
	require.NoError(t, os.WriteFile(filepath.Join(home, "target"), target.Bytes(), 0600))
	require.NoError(t, dsc.UpgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{Name: "reliability-crash", Height: 3}))
	dsc.Commit()
	require.NoError(t, db.Close())
	child := exec.Command(os.Args[0], "-test.run=^TestIntegrationUpgradeCrashRealEVM$", "-test.v")
	child.Env = append(os.Environ(), "RELIABILITY_CRASH_CHILD=1", "RELIABILITY_HOME="+home)
	out, err := child.CombinedOutput()
	var exit *exec.ExitError
	require.ErrorAs(t, err, &exit, "child output: %s", out)
	require.Equal(t, 83, exit.ExitCode(), "child output: %s", out)
	recovered, recoveredDB := reliabilityOpen(t, home)
	defer recoveredDB.Close()
	recovered.UpgradeKeeper.SetUpgradeHandler("reliability-crash", func(c sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		_, err := recovered.EvmKeeper.CallEVM(c, a, common.Address{}, target, true, "bump")
		return vm, err
	})
	replay := reliabilityContext(t, recovered, home, 3)
	require.Equal(t, int64(1), reliabilitySlot(recovered, replay, target, 0), "EVM migration must execute again after real process exit before commit")
	require.Equal(t, int64(3), recovered.UpgradeKeeper.GetDoneHeight(replay, "reliability-crash"))
	recovered.Commit()
	require.NoError(t, recoveredDB.Close())
	again, againDB := reliabilityOpen(t, home)
	defer againDB.Close()
	again.UpgradeKeeper.SetUpgradeHandler("reliability-crash", func(c sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		_, err := again.EvmKeeper.CallEVM(c, a, common.Address{}, target, true, "bump")
		return vm, err
	})
	next := reliabilityContext(t, again, home, 4)
	require.Equal(t, int64(1), reliabilitySlot(again, next, target, 0), "committed EVM migration must not execute twice")
	t.Log("real process exit(83), on-disk LevelDB reopen, EVM migration replay and once-only committed state verified")
}

func TestIntegrationAutoUnbondRealEVM(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		name := "legacy"
		if enabled {
			name = "activated"
		}
		t.Run(name, func(t *testing.T) {
			reliabilityActivate(t, enabled)
			dsc, ctx, _, db := reliabilityNew(t)
			defer db.Close()
			target, a := reliabilityDeploy(t, dsc, ctx, "AutoUnbondFixture")
			master, _ := reliabilityDeploy(t, dsc, ctx, "ValidatorFixture")
			reliabilityRegistry(t, dsc, ctx, target, master)
			coin, err := dsc.CoinKeeper.GetCoin(ctx, cmdcfg.BaseDenom)
			require.NoError(t, err)
			coin.DRC20Contract = common.HexToAddress("0x4444").Hex()
			dsc.CoinKeeper.SetCoin(ctx, coin)
			val := sdk.ValAddress(bytes.Repeat([]byte{0x41}, 20))
			del := sdk.AccAddress(bytes.Repeat([]byte{0x42}, 20))
			stake := validatortypes.NewStakeCoin(sdk.NewInt64Coin(cmdcfg.BaseDenom, 100))
			stake.Holds = []*validatortypes.StakeHold{{Amount: sdk.NewInt(10), HoldEndTime: 1900000000}, {Amount: sdk.NewInt(20), HoldEndTime: 2000000000}}
			dsc.ValidatorKeeper.SetDelegation(ctx, validatortypes.NewDelegation(del, val, stake))
			params := dsc.ValidatorKeeper.GetParams(ctx)
			params.AutoUnbondTimeout = time.Hour
			dsc.ValidatorKeeper.SetParams(ctx, params)
			dsc.ValidatorKeeper.SetValidatorOfflineSince(ctx, val, ctx.BlockTime().Add(-2*time.Hour))
			_, err = dsc.EvmKeeper.CallEVM(ctx, a, common.Address{}, target, true, "setFailHold", big.NewInt(2000000000))
			require.NoError(t, err)
			probe, _ := ctx.CacheContext()
			probeErr := dsc.ValidatorKeeper.ExecuteAutoUnbondEnqueue(probe, validatortypes.NewDelegation(del, val, stake))
			t.Logf("enqueue probe: err=%v prefix=%d", probeErr, reliabilitySlot(dsc, probe, target, 0))
			dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)
			_, tracked := dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, val)
			if !enabled {
				require.False(t, tracked)
				require.Equal(t, int64(1), reliabilitySlot(dsc, ctx, target, 0))
				t.Log("legacy control: first real EVM enqueue committed; second reverted; retry marker deleted")
				return
			}
			require.True(t, tracked, "failed job must remain retryable")
			require.Equal(t, int64(0), reliabilitySlot(dsc, ctx, target, 0), "first successful EVM call must roll back with failed group")
			_, err = dsc.EvmKeeper.CallEVM(ctx, a, common.Address{}, target, true, "setFailHold", big.NewInt(0))
			require.NoError(t, err)
			dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)
			require.Equal(t, int64(3), reliabilitySlot(dsc, ctx, target, 0))
			_, tracked = dsc.ValidatorKeeper.GetValidatorOfflineSince(ctx, val)
			require.False(t, tracked)
			dsc.ValidatorKeeper.ProcessAutoUnbond(ctx)
			require.Equal(t, int64(3), reliabilitySlot(dsc, ctx, target, 0), "successful job must not enqueue twice")
			t.Log("real EVM revert on second enqueue: atomic rollback, retained marker, 3 successful calls on retry, no second enqueue verified")
		})
	}
}

func TestIntegrationEventRollbackRealEVM(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, empty := range []bool{false, true} {
			name := "legacy"
			if enabled {
				name = "activated"
			}
			if empty {
				name += "_empty"
			}
			t.Run(name, func(t *testing.T) {
				reliabilityActivate(t, enabled)
				dsc, ctx, _, db := reliabilityNew(t)
				defer db.Close()
				auto, _ := reliabilityDeploy(t, dsc, ctx, "AutoUnbondFixture")
				target, a := reliabilityDeploy(t, dsc, ctx, "ValidatorFixture")
				reliabilityRegistry(t, dsc, ctx, auto, target)
				dsc.ValidatorKeeper.SetRewardPerBlockOverride(ctx, sdk.NewInt(7))
				key, err := crypto.GenerateKey()
				require.NoError(t, err)
				sender := crypto.PubkeyToAddress(key.PublicKey)
				acct := dsc.AccountKeeper.NewAccountWithAddress(ctx, sdk.AccAddress(sender.Bytes()))
				dsc.AccountKeeper.SetAccount(ctx, acct)
				coins := sdk.NewCoins(sdk.NewInt64Coin(cmdcfg.BaseDenom, 1000000))
				require.NoError(t, dsc.BankKeeper.MintCoins(ctx, cointypes.ModuleName, coins))
				require.NoError(t, dsc.BankKeeper.SendCoinsFromModuleToAccount(ctx, cointypes.ModuleName, sdk.AccAddress(sender.Bytes()), coins))
				payload, err := a.Pack("emitBad", empty)
				require.NoError(t, err)
				tx := ethtypes.NewTx(&ethtypes.LegacyTx{Nonce: dsc.EvmKeeper.GetNonce(ctx, sender), To: &target, Gas: 500000, GasPrice: big.NewInt(0), Value: big.NewInt(123), Data: payload})
				signed, err := ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(big.NewInt(75)), key)
				require.NoError(t, err)
				msg := &evmtypes.MsgEthereumTx{}
				require.NoError(t, msg.FromEthereumTx(signed))
				msg.From = sender.Hex()
				res, err := dsc.EvmKeeper.ApplyTransaction(ctx, msg)
				require.NoError(t, err)
				t.Logf("real transaction: vmError=%q postError=%q logs=%d gas=%d", res.VmError, res.PostTxError, len(res.Logs), res.GasUsed)
				_, override := dsc.ValidatorKeeper.GetRewardPerBlockOverride(ctx)
				if !enabled {
					require.False(t, res.Failed())
					require.Equal(t, int64(1), reliabilitySlot(dsc, ctx, target, 0))
					require.False(t, override)
					require.Equal(t, int64(123), dsc.EvmKeeper.GetBalance(ctx, target).Int64())
					t.Log("legacy control: malformed event accepted; override cleared; storage and value transfer committed")
					return
				}
				require.True(t, res.Failed(), "hook error must mark transaction failed")
				require.Equal(t, evmtypes.ErrPostTxProcessing.Error(), res.VmError, "failure must come from the hook, not an unrelated VM failure")
				require.Empty(t, res.Logs)
				require.Equal(t, int64(0), reliabilitySlot(dsc, ctx, target, 0))
				require.True(t, override)
				require.Zero(t, dsc.EvmKeeper.GetBalance(ctx, target).Sign(), "value transfer must roll back")
				payload, err = a.Pack("emitGood", big.NewInt(42), true)
				require.NoError(t, err)
				tx = ethtypes.NewTx(&ethtypes.LegacyTx{Nonce: dsc.EvmKeeper.GetNonce(ctx, sender), To: &target, Gas: 500000, GasPrice: big.NewInt(0), Data: payload})
				signed, err = ethtypes.SignTx(tx, ethtypes.LatestSignerForChainID(big.NewInt(75)), key)
				require.NoError(t, err)
				msg = &evmtypes.MsgEthereumTx{}
				require.NoError(t, msg.FromEthereumTx(signed))
				msg.From = sender.Hex()
				res, err = dsc.EvmKeeper.ApplyTransaction(ctx, msg)
				require.NoError(t, err)
				require.False(t, res.Failed())
				amount, ok := dsc.ValidatorKeeper.GetRewardPerBlockOverride(ctx)
				require.True(t, ok)
				require.Equal(t, "42", amount.String())
				require.Equal(t, int64(1), reliabilitySlot(dsc, ctx, target, 0))
				t.Log("signed Ethereum transaction: malformed event reverted storage + native value transfer + reward state; valid event committed")
			})
		}
	}
}
