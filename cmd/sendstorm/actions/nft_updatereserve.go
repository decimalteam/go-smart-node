package actions

import (
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgUpdateReserveNFT
type UpdateReserveNFTGenerator struct {
	increaseBottom, increaseUp int64 // value in 10^18 (del)
	knownNFT                   []dscApi.NFT
	rnd                        *rand.Rand
}

type UpdateReserveNFTAction struct {
	creator  string // need for filter
	id       string
	denom    string
	increase int64 // value in 10^18 (del)
	subIds   []uint64
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
}

func (gg *UpdateReserveNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownNFT))))
	nftToUpdateReserve := gg.knownNFT[i]
	subTokenIDs := make([]uint64, 0)
	for _, o := range nftToUpdateReserve.Owners {
		if o.Address == nftToUpdateReserve.Creator {
			subTokenIDs = append(subTokenIDs, o.SubTokenIDs...)
			break
		}

	}
	// creator not in owners
	if len(subTokenIDs) == 0 {
		return &EmptyAction{}
	}
	increase := RandomRange(gg.rnd, gg.increaseBottom, gg.increaseUp)
	subToUpdate := RandomSublist(gg.rnd, subTokenIDs)

	return &UpdateReserveNFTAction{
		creator:  nftToUpdateReserve.Creator,
		id:       nftToUpdateReserve.ID,
		denom:    nftToUpdateReserve.Denom,
		subIds:   subToUpdate,
		increase: increase,
	}
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
		if saList[i].BalanceForCoin(saList[i].FeeDenom()).LT(helpers.EtherToWei(sdk.NewInt(aa.increase))) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *UpdateReserveNFTAction) GenerateTx(sa *stormTypes.StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgUpdateReserveNFT(
		sender,
		aa.id,
		aa.denom,
		aa.subIds,
		helpers.EtherToWei(sdk.NewInt(aa.increase)),
	)
	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", sa.FeeDenom())
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (aa *UpdateReserveNFTAction) String() string {
	return fmt.Sprintf("UpdateReserveNFT{ID: %s, Denom: %s, SubIds: %v, Increase: %d}",
		aa.id,
		aa.denom,
		aa.subIds,
		aa.increase,
	)
}
