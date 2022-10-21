package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgEditNFT
type EditNFTGenerator struct {
	textLengthBottom, textLengthUp int64
	knownNFT                       []*dscApi.NFTToken
	rnd                            *rand.Rand
}

type EditNFTAction struct {
	creator     string // need for filter
	id          string
	denom       string
	newTokenUri string
}

func NewEditNFTGenerator(textLengthBottom, textLengthUp int64) *EditNFTGenerator {
	return &EditNFTGenerator{
		textLengthBottom: textLengthBottom,
		textLengthUp:     textLengthUp,
		rnd:              rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *EditNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFT = ui.NFTs
}

func (gg *EditNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownNFT))))
	nftToEdit := gg.knownNFT[i]

	return &EditNFTAction{
		creator:     nftToEdit.Creator,
		id:          nftToEdit.ID,
		denom:       nftToEdit.Denom,
		newTokenUri: RandomString(gg.rnd, RandomRange(gg.rnd, gg.textLengthBottom, gg.textLengthUp), charsAll),
	}
}

func (aa *EditNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != aa.creator {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *EditNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgUpdateToken(
		sender,
		aa.id,
		aa.newTokenUri,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *EditNFTAction) String() string {
	return fmt.Sprintf("EditNFT{ID: %s, Denom: %s, TokenURI: %s}",
		aa.id,
		aa.denom,
		aa.newTokenUri,
	)
}
