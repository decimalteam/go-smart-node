package upgrade

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

var (
	downloadStat = make(map[string]bool)
)

// BeginBlocker will check if there is a scheduled plan and if it is ready to be executed.
// If the current height is in the provided set of heights to skip, it will skip and clear the upgrade plan.
// If it is ready, it will execute it if the handler is installed, and panic/abort otherwise.
// If the plan is not ready, it will ensure the handler is not registered too early (and abort otherwise).
//
// The purpose is to ensure the binary is switched EXACTLY at the desired block, and to allow
// a migration to be executed if needed upon this switch (migration defined in the new binary)
// skipUpgradeHeightArray is a set of block heights for which the upgrade must be skipped
func BeginBlocker(k keeper.Keeper, ctx sdk.Context, _ abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	plan, found := k.GetUpgradePlan(ctx)

	if !k.DowngradeVerified() {
		k.SetDowngradeVerified(true)
		lastAppliedPlan, _ := k.GetLastCompletedUpgrade(ctx)
		// This check will make sure that we are using a valid binary.
		// It'll panic in these cases if there is no upgrade handler registered for the last applied upgrade.
		// 1. If there is no scheduled upgrade.
		// 2. If the plan is not ready.
		// 3. If the plan is ready and skip upgrade height is set for current height.
		if !found || !plan.ShouldExecute(ctx) || (plan.ShouldExecute(ctx) && k.IsSkipHeight(ctx.BlockHeight())) {
			if lastAppliedPlan != "" && !k.HasHandler(lastAppliedPlan) {
				var appVersion uint64
				cp := ctx.ConsensusParams()
				if cp != nil && cp.Version != nil {
					appVersion = cp.Version.AppVersion
				}
				panic(fmt.Sprintf("Wrong app version %d, upgrade handler is missing for %s upgrade plan", appVersion, lastAppliedPlan))
			}
		}
	}

	if !found {
		return
	}

	logger := ctx.Logger()

	allBlocks := config.UpdatesInfo.AllBlocks
	if _, ok := allBlocks[plan.Name]; ok {
		return
	}

	_, ok := downloadStat[plan.Name]

	// To make sure clear upgrade is executed at the same block
	if ctx.BlockHeight() < plan.Height && !ok {
		mapping := planMapping(plan)
		if mapping == nil {
			logger.Error("error: plan mapping decode")
			return
		}
		/*
			{
			  "binaries": {
			    "darwin/arm64": "http://127.0.0.1:8080/50/darwin/dscd?checksum=sha256:1ca5d7c987196f6d6fd9101c138dec210211b0fcbfacde26a5a97abd1668dd89"
			  }
			}
		*/
		system, ok := mapping[osArchForURL()]
		if !ok {
			logger.Error(fmt.Sprintf("error: plan mapping[os] for '%s' undefined", osArchForURL()))
			return
		}
		// example:
		// from http://127.0.0.1:8080/50/darwin/dscd?checksum=sha256:1ca5d7c987196f6d6fd9101c138dec210211b0fcbfacde26a5a97abd1668dd89
		// to "http://127.0.0.1:8080/50/darwin/dscd", 1ca5d7c987196f6d6fd9101c138dec210211b0fcbfacde26a5a97abd1668dd89
		newUrl, checksum := resolveDownloadURL(system)
		if newUrl == "" || checksum == "" {
			logger.Error("error: failed with generate url")
			return
		}

		if !doesPageExist(newUrl) {
			logger.Error("error: url page is not exists")
			return
		}

		downloadStat[plan.Name] = true
		downloadName := getDownloadFileName(config.AppBinName)

		if _, err := os.Stat(downloadName); os.IsNotExist(err) {
			go DownloadSecured(ctx, newUrl, downloadName, checksum)
		}
	}

	// To make sure clear upgrade is executed at the same block
	if plan.ShouldExecute(ctx) {
		// If skip upgrade has been set for current height, we clear the upgrade plan
		if k.IsSkipHeight(ctx.BlockHeight()) {
			skipUpgradeMsg := fmt.Sprintf("UPGRADE \"%s\" SKIPPED at %d: %s", plan.Name, plan.Height, plan.Info)
			logger.Info(skipUpgradeMsg)

			// Clear the upgrade plan at current height
			k.ClearUpgradePlan(ctx)
			return
		}

		// Prepare shutdown if we don't have an upgrade handler for this upgrade name (meaning this software is out of date)
		if !k.HasHandler(plan.Name) {
			if _, err := os.Stat(getDownloadFileName(config.AppBinName)); err == nil {
				err = changeBinary(plan)
				if err != nil {
					panic(fmt.Errorf("failed to change binaries err: %s", err.Error()))
				}
			}

			// Write the upgrade info to disk. The UpgradeStoreLoader uses this info to perform or skip
			// store migrations.
			err := k.DumpUpgradeInfoToDisk(ctx.BlockHeight(), plan)
			if err != nil {
				panic(fmt.Errorf("unable to write upgrade info to filesystem: %s", err.Error()))
			}

			upgradeMsg := BuildUpgradeNeededMsg(plan)
			logger.Error(upgradeMsg)
			panic(upgradeMsg)
		}

		// We have an upgrade handler for this upgrade name, so apply the upgrade
		logger.Info(fmt.Sprintf("applying upgrade \"%s\" at %s", plan.Name, plan.DueAt()))
		ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

		k.ApplyUpgrade(ctx, plan)

		config.UpdatesInfo.PushNewPlanHeight(plan.Height)
		err := config.UpdatesInfo.Save()
		if err != nil {
			logger.Error(fmt.Sprintf("push \"%s\" with error: %s", plan.Name, err.Error()))
			os.Exit(2)
		}

		config.UpdatesInfo.AddExecutedPlan(plan.Name, plan.Height)
		err = config.UpdatesInfo.Save()
		if err != nil {
			logger.Error(fmt.Sprintf("save \"%s\" with '%s'", plan.Name, err.Error()))
			os.Exit(3)
		}

		return
	}

	// if we have a pending upgrade, but it is not yet time, make sure we did not
	// set the handler already
	if k.HasHandler(plan.Name) {
		downgradeMsg := fmt.Sprintf("BINARY UPDATED BEFORE TRIGGER! UPGRADE \"%s\" - in binary but not executed on chain. Downgrade your binary", plan.Name)
		logger.Error(downgradeMsg)
		panic(downgradeMsg)
	}
}

// BuildUpgradeNeededMsg prints the message that notifies that an upgrade is needed.
func BuildUpgradeNeededMsg(plan types.Plan) string {
	return fmt.Sprintf("UPGRADE \"%s\" NEEDED at %s: %s", plan.Name, plan.DueAt(), plan.Info)
}

func loadVersion(urlPath string) string {
	const fileVersion = "version.txt"

	// example: "version.txt"
	u, err := url.Parse(fileVersion)
	if err != nil {
		log.Fatal(err)
	}

	// example: "https://testnet-repo.decimalchain.com/95000"
	base, err := u.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}

	// result: "https://testnet-repo.decimalchain.com/version.txt"
	resp, err := http.Get(base.ResolveReference(u).String())
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return strings.TrimSpace(string(body))
}
