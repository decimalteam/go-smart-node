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

// MsgBurnNFT
type BurnNFTGenerator struct {
	knownNFT []*dscApi.NFTToken
	rnd      *rand.Rand
}

type BurnNFTAction struct {
	creator string // need for filter
	id      string
	denom   string
	subIds  []uint32
	nft     dscApi.NFTToken
}

func NewBurnNFTGenerator() *BurnNFTGenerator {
	return &BurnNFTGenerator{rnd: rand.New(rand.NewSource(time.Now().Unix()))}
}

func (gg *BurnNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFT = ui.NFTs
}

func (gg *BurnNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	// 10 attepmts to get ntf to burn
	for n := 0; n < 10; n++ {
		i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownNFT))))
		nftToBurn := gg.knownNFT[i]
		subTokenIDs := make([]uint32, 0)
		for _, sub := range nftToBurn.SubTokens {
			if sub.Owner == nftToBurn.Creator {
				subTokenIDs = append(subTokenIDs, sub.ID)
			}
		}
		if len(subTokenIDs) == 0 {
			continue
		}
		subToBurn := RandomSublist(gg.rnd, subTokenIDs)
		return &BurnNFTAction{
			creator: nftToBurn.Creator,
			id:      nftToBurn.ID,
			denom:   nftToBurn.Denom,
			subIds:  subToBurn,
			nft:     *nftToBurn,
		}
	}
	return &EmptyAction{}
}

func (aa *BurnNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
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

func (aa *BurnNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgBurnToken(
		sender,
		aa.id,
		aa.subIds,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *BurnNFTAction) String() string {
	return fmt.Sprintf("BurnNFT{ID: %s, Denom: %s, SubIds: %v}",
		aa.id,
		aa.denom,
		aa.subIds,
	)
}
