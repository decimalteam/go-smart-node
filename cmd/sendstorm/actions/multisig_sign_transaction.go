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

type SignMultisigTransactionGenerator struct {
	knownWallets      []dscApi.MultisigWallet
	knownTransactions []dscApi.MultisigTransactionInfo
	rnd               *rand.Rand
}

type SignMultisigTransactionAction struct {
	txID            string
	possibleSigners []string
}

func NewSignMultisigTransactionGenerator() *SignMultisigTransactionGenerator {
	return &SignMultisigTransactionGenerator{
		rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *SignMultisigTransactionGenerator) Update(ui UpdateInfo) {
	gg.knownWallets = ui.MultisigWallets
	gg.knownTransactions = ui.MultisigTransactions
}

func isExecuted(tx dscApi.MultisigTransactionInfo) bool {
	return tx.Completed
}

func extractPossibleSigners(wallet dscApi.MultisigWallet, tx dscApi.MultisigTransactionInfo) []string {
	var result []string
	for i := range wallet.Owners {
		if !strings.StringInSlice(wallet.Owners[i], tx.Signers) {
			result = append(result, wallet.Owners[i])
		}
	}
	return result
}

func (gg *SignMultisigTransactionGenerator) Generate() Action {
	if len(gg.knownWallets) == 0 {
		return &EmptyAction{}
	}
	if len(gg.knownTransactions) == 0 {
		return &EmptyAction{}
	}
	// try to find unexecuted transaction
	for n := 0; n < 10; n++ {
		i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownTransactions))))
		tx := gg.knownTransactions[i]
		wallet := dscApi.MultisigWallet{}
		for _, w := range gg.knownWallets {
			if w.Address == tx.Transaction.Wallet {
				wallet = w
				break
			}
		}
		if !isExecuted(tx) {
			signers := extractPossibleSigners(wallet, tx)
			if len(signers) == 0 {
				return &EmptyAction{}
			}
			return &SignMultisigTransactionAction{
				txID:            tx.Transaction.Id,
				possibleSigners: signers,
			}
		}
	}
	return &EmptyAction{}
}

func (aa *SignMultisigTransactionAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (aa *SignMultisigTransactionAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgSignTransaction(sender, aa.txID)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *SignMultisigTransactionAction) String() string {
	return fmt.Sprintf("SignMultisigTransaction{txID: %s, signers: %v}",
		aa.txID, aa.possibleSigners)
}
