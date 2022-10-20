package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// MsgUpdateReserveNFT
type UpdateReserveNFTGenerator struct {
	increaseBottom, increaseUp int64 // value in 10^18 (del)
	knownNFT                   []*dscApi.NFTToken
	knownSubtokenReserves      map[NFTSubTokenKey]sdk.Coin
	rnd                        *rand.Rand
}

type UpdateReserveNFTAction struct {
	creator    string // need for filter
	id         string
	denom      string
	newReserve sdk.Coin
	subIds     []uint32
}

func NewUpdateReserveNFTGenerator(increaseBottom, increaseUp int64) *UpdateReserveNFTGenerator {
	return &UpdateReserveNFTGenerator{
		increaseBottom: increaseBottom,
		increaseUp:     increaseUp,
		rnd:            rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *UpdateReserveNFTGenerator) Update(ui UpdateInfo) {
	gg.knownNFT = ui.NFTs
	gg.knownSubtokenReserves = ui.NFTSubTokenReserves
}

func (gg *UpdateReserveNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	// 10 attepmts to get ntf to burn
	for n := 0; n < 10; n++ {
		i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownNFT))))
		nftToUpdateReserve := gg.knownNFT[i]
		subTokenIDs := make([]uint32, 0)
		for _, sub := range nftToUpdateReserve.SubTokens {
			if sub.Owner == nftToUpdateReserve.Creator {
				subTokenIDs = append(subTokenIDs, sub.ID)
			}
		}
		// creator not in owners
		if len(subTokenIDs) == 0 {
			continue
		}
		increase := RandomRange(gg.rnd, gg.increaseBottom, gg.increaseUp)
		subToUpdate := RandomSublist(gg.rnd, subTokenIDs)
		// get max reserve of subtokens
		newReserve := sdk.ZeroInt()
		for _, s := range subToUpdate {
			key := NFTSubTokenKey{TokenID: nftToUpdateReserve.ID, ID: s}
			reserve, ok := gg.knownSubtokenReserves[key]
			if !ok {
				continue
			}
			if newReserve.LT(reserve.Amount) {
				newReserve = reserve.Amount
			}
		}
		newReserve = newReserve.Add(helpers.EtherToWei(sdkmath.NewInt(increase)))

		return &UpdateReserveNFTAction{
			creator:    nftToUpdateReserve.Creator,
			id:         nftToUpdateReserve.ID,
			denom:      nftToUpdateReserve.Denom,
			subIds:     subToUpdate,
			newReserve: sdk.NewCoin(cmdcfg.BaseDenom, newReserve),
		}
	}
	return &EmptyAction{}
}

func (aa *UpdateReserveNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != aa.creator {
			continue
		}
		// TODO
		//if saList[i].BalanceForCoin(saList[i].FeeDenom()).LT(helpers.EtherToWei(sdk.NewInt(aa.increase))) {
		//	continue
		//}
		res = append(res, saList[i])
	}
	return res
}

func (aa *UpdateReserveNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgUpdateReserve(
		sender,
		aa.id,
		aa.subIds,
		aa.newReserve,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *UpdateReserveNFTAction) String() string {
	return fmt.Sprintf("UpdateReserveNFT{ID: %s, Denom: %s, SubIds: %v, newReserve: %s}",
		aa.id,
		aa.denom,
		aa.subIds,
		aa.newReserve.String(),
	)
}
