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

// MsgTransferNFT
type TransferNFTGenerator struct {
	knownAddresses []string
	knownNFT       []*dscApi.NFTToken
	rnd            *rand.Rand
}

type TransferNFTAction struct {
	owner     string // need for filter
	recipient string
	id        string
	denom     string
	subIds    []uint32
}

func NewTransferNFTGenerator() *TransferNFTGenerator {
	return &TransferNFTGenerator{rnd: rand.New(rand.NewSource(time.Now().Unix()))}
}

func (gg *TransferNFTGenerator) Update(ui UpdateInfo) {
	gg.knownAddresses = ui.Addresses
	gg.knownNFT = ui.NFTs
}

func (gg *TransferNFTGenerator) Generate() Action {
	if len(gg.knownNFT) == 0 {
		return &EmptyAction{}
	}
	i := int(RandomRange(gg.rnd, 0, int64(len(gg.knownNFT))))
	nftToTransfer := gg.knownNFT[i]
	i = int(RandomRange(gg.rnd, 0, int64(len(nftToTransfer.SubTokens))))
	tokenOwner := nftToTransfer.SubTokens[i].Owner
	subIds := make([]uint32, 0)
	for _, sub := range nftToTransfer.SubTokens {
		if sub.Owner == tokenOwner {
			subIds = append(subIds, sub.ID)
		}
	}
	subIds = RandomSublist32(gg.rnd, subIds)
	var recipient = RandomChoice(gg.rnd, gg.knownAddresses)
	for j := 0; j < 10; j++ {
		// 10 attempts to get recipient != owner
		if j == 9 {
			return &EmptyAction{}
		}
		if recipient != tokenOwner {
			break
		}
		recipient = RandomChoice(gg.rnd, gg.knownAddresses)
	}
	return &TransferNFTAction{
		owner:     tokenOwner,
		recipient: RandomChoice(gg.rnd, gg.knownAddresses),
		id:        nftToTransfer.ID,
		denom:     nftToTransfer.Denom,
		subIds:    subIds,
	}
}

func (aa *TransferNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].Address() != aa.owner {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *TransferNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	recipient, err := sdk.AccAddressFromBech32(aa.recipient)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgTransferNFT(
		sender,
		recipient,
		aa.denom,
		aa.id,
		aa.subIds,
	)
	tx, err := dscTx.BuildTransaction(sa.Account(), []sdk.Msg{msg}, "", sa.FeeDenom(), feeConfig.DelPrice, feeConfig.Params)
	if err != nil {
		return nil, err
	}
	err = tx.SignTransaction(sa.Account())
	if err != nil {
		return nil, err
	}
	return tx.BytesToSend()
}

func (aa *TransferNFTAction) String() string {
	return fmt.Sprintf("TransferNFT{ID: %s, Recipient: %s, Denom: %s, SubIds: %v}",
		aa.id,
		aa.recipient,
		aa.denom,
		aa.subIds,
	)
}
