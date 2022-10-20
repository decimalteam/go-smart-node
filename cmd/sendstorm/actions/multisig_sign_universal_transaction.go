package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/strings"
)

type SignMultisigUniversalTransactionGenerator struct {
	knownWallets      []dscApi.MultisigWallet
	knownTransactions []dscApi.MultisigUniversalTransactionResponse
	rnd               *rand.Rand
}

type SignMultisigUniversalTransactionAction struct {
	txID            string
	possibleSigners []string
}

func NewSignMultisigUniversalTransactionGenerator() *SignMultisigUniversalTransactionGenerator {
	return &SignMultisigUniversalTransactionGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *SignMultisigUniversalTransactionGenerator) Update(ui UpdateInfo) {
	gg.knownWallets = ui.MultisigWallets
	gg.knownTransactions = ui.MultisigUniversalTransactions
}

func extractPossibleSignersFromUniversal(wallet dscApi.MultisigWallet, signedBy []string) []string {
	var result []string
	for _, owner := range wallet.Owners {
		if !strings.StringInSlice(owner, signedBy) {
			result = append(result, owner)
		}
	}
	return result
}

func (gg *SignMultisigUniversalTransactionGenerator) Generate() Action {
	if len(gg.knownWallets) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownTransactions) == 0 {
		return &EmptyAction{}
	}
	// try to find unexecuted transaction
	for n := 0; n < 10; n++ {
		tx := RandomChoice(gg.rnd, gg.knownTransactions)
		if tx.Completed {
			continue
		}
		wallet := dscApi.MultisigWallet{}
		for _, w := range gg.knownWallets {
			if w.Address == tx.Transaction.Wallet {
				wallet = w
				break
			}
		}
		signers := extractPossibleSignersFromUniversal(wallet, tx.Signers)
		if len(signers) == 0 {
			return &EmptyAction{}
		}
		return &SignMultisigUniversalTransactionAction{
			txID:            tx.Transaction.Id,
			possibleSigners: signers,
		}
	}
	return &EmptyAction{}
}

func (aa *SignMultisigUniversalTransactionAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if !strings.StringInSlice(saList[i].Address(), aa.possibleSigners) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *SignMultisigUniversalTransactionAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSignUniversalTransaction(sender, aa.txID)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *SignMultisigUniversalTransactionAction) String() string {
	return fmt.Sprintf("SignMultisigUniversalTransaction{txID: %s, signers: %v}",
		aa.txID, aa.possibleSigners)
}
