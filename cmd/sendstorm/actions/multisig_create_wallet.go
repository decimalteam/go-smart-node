package actions

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateWallet
type CreateMultisigWalletGenerator struct {
	knownAddresses []string
	rnd            *rand.Rand
}

type CreateMultisigWalletAction struct {
	owners    []string
	weights   []uint32
	threshold uint32
}

func NewCreateMultisigWalletGenerator() *CreateMultisigWalletGenerator {
	return &CreateMultisigWalletGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *CreateMultisigWalletGenerator) Update(ui UpdateInfo) {
	gg.knownAddresses = ui.Addresses
}

func (gg *CreateMultisigWalletGenerator) Generate() Action {
	ownersCount := int(RandomRange(gg.rnd, 2, 16+1)) // x/multisig/types/config.go
	if len(gg.knownAddresses) < ownersCount {
		return &EmptyAction{}
	}
	var action CreateMultisigWalletAction
	var sumOfWeights uint32
	action.owners = make([]string, ownersCount)
	action.weights = make([]uint32, ownersCount)

	for i := 0; i < ownersCount; i++ {
		var owner string
		var keepFind = true
		for keepFind {
			owner = RandomChoice(gg.rnd, gg.knownAddresses)
			keepFind = false
			for j := 0; j < ownersCount; j++ {
				if action.owners[j] == owner {
					keepFind = true
				}
			}
		}
		action.owners[i] = owner
		action.weights[i] = uint32(RandomRange(gg.rnd, 1, 1024))
		sumOfWeights += action.weights[i]
	}
	action.threshold = uint32(RandomRange(gg.rnd, 1, int64(sumOfWeights)))

	return &action
}

func (aa *CreateMultisigWalletAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *CreateMultisigWalletAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgCreateWallet(
		sender,
		aa.owners,
		aa.weights,
		aa.threshold,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *CreateMultisigWalletAction) String() string {
	return fmt.Sprintf("CreateMultisigWallet{Owners: [%s], Weights: %v, Threshols: %d}",
		strings.Join(aa.owners, ", "),
		aa.weights, aa.threshold)
}
