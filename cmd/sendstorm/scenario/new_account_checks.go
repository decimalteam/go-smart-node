package scenario

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	dscWallet "bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

type RedeemChecksScenario struct {
	issuer *stormTypes.StormAccount
	api    *dscApi.API
	rnd    *rand.Rand
}

func NewRedeemChecksScenario(api *dscApi.API, issuer *stormTypes.StormAccount) *RedeemChecksScenario {
	return &RedeemChecksScenario{
		issuer: issuer,
		api:    api,
		rnd:    rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (rcs *RedeemChecksScenario) MakeCheck() {
	const password = "der_password"
	rcs.issuer.UpdateBalance()
	rcs.issuer.UpdateNumberSequence()

	nonce := sdkmath.NewInt(rand.Int63())
	dueBlock := rcs.api.GetLastHeight()
	if dueBlock == 0 {
		fmt.Printf("MakeCheck-block height is 0\n")
		return
	}

	dueBlock += rand.Int63n(100) + 100

	amount := helpers.EtherToWei(sdkmath.NewInt(rand.Int63n(10) + 1))

	checkBase58, err := dscTx.IssueCheck(rcs.issuer.Account(), "del", amount, nonce, uint64(dueBlock), password)
	if err != nil {
		fmt.Printf("MakeCheck-IssueCheck: %s\n", err.Error())
		return
	}

	mn, _ := dscWallet.NewMnemonic("")
	acc, err := dscWallet.NewAccountFromMnemonicWords(mn.Words(), "")
	if err != nil {
		fmt.Printf("MakeCheck-NewAccount: %s\n", err.Error())
		return
	}
	fmt.Printf("acc: %s\n", acc.Address())
	fmt.Printf("mnemonic: %s\n", mn.Words())

	acc.WithChainID(rcs.api.ChainID())

	msg, err := dscTx.CreateRedeemCheck(acc, checkBase58, password)
	if err != nil {
		fmt.Printf("MakeCheck-CreateRedeem: %s\n", err.Error())
		return
	}

	tx, err := dscTx.BuildTransaction(acc, []sdk.Msg{msg}, "", "del", rcs.api.GetFeeCalculationOptions())
	if err != nil {
		fmt.Printf("MakeCheck-BuildTransaction: %s\n", err.Error())
		return
	}
	err = tx.SignTransaction(acc)
	if err != nil {
		fmt.Printf("MakeCheck-SignTransaction: %s\n", err.Error())
		return
	}
	bz, err := tx.BytesToSend()
	if err != nil {
		fmt.Printf("MakeCheck-BytesToSend err: %s\n", err.Error())
		return
	}
	// check throught simulate
	simres, err := rcs.api.SimulateTx(bz)
	if err != nil {
		fmt.Printf("MakeCheck-Simulate err: %s\n", err.Error())
		return
	}
	fmt.Printf("simulate result: %d events\n", len(simres.Events))
	for _, ev := range simres.Events {
		fmt.Printf("event: %s\n", ev.Type)
	}

	resp, err := rcs.api.BroadcastTxCommit(bz)
	if err != nil {
		fmt.Printf("MakeCheck-BroadcastTxSync err: %s\n", err.Error())
	} else {
		fmt.Printf("result: %#v\n", resp)
	}
	fmt.Printf("Events:\n")
	for _, ev := range resp.Events {
		fmt.Printf("%s:\n", ev.Type)
		for _, kv := range ev.Attributes {
			fmt.Printf("\t%s = %s\n", string(kv.Key), string(kv.Value))
		}
	}
	fmt.Printf("\n\n\n")

}
