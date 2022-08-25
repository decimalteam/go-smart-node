package actions

import (
	"bitbucket.org/decimalteam/go-smart-node/cmd/config"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"math/rand"
	"time"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgMintNFT
type MintNFTGenerator struct {
	// general values
	textLengthBottom, textLengthUp int64
	reserveBottom, reserveUp       int64 // in 10^18
	quantityBottom, quantityUp     int64
	knownAddresses                 []string
	rnd                            *rand.Rand
}

type MintNFTAction struct {
	recipient string
	id        string
	denom     string
	tokenURI  string
	quantity  sdkmath.Int
	reserve   sdk.Coin
	allowMint bool
}

func NewMintNFTGenerator(
	textLengthBottom, textLengthUp,
	reserveBottom, reserveUp,
	quantityBottom, quantityUp int64) *MintNFTGenerator {
	return &MintNFTGenerator{
		textLengthBottom: textLengthBottom,
		textLengthUp:     textLengthUp,
		reserveBottom:    reserveBottom,
		reserveUp:        reserveUp,
		quantityBottom:   quantityBottom,
		quantityUp:       quantityUp,
		rnd:              rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (gg *MintNFTGenerator) Update(ui UpdateInfo) {
	gg.knownAddresses = ui.Addresses
}

func (gg *MintNFTGenerator) Generate() Action {
	if len(gg.knownAddresses) == 0 {
		return &EmptyAction{}
	}
	return &MintNFTAction{
		recipient: RandomChoice(gg.rnd, gg.knownAddresses),
		id:        RandomString(gg.rnd, RandomRange(gg.rnd, gg.textLengthBottom, gg.textLengthUp), charsAll),
		denom:     RandomString(gg.rnd, RandomRange(gg.rnd, gg.textLengthBottom, gg.textLengthUp), charsAll),
		tokenURI:  RandomString(gg.rnd, RandomRange(gg.rnd, gg.textLengthBottom, gg.textLengthUp), charsAll),
		quantity:  sdk.NewInt(RandomRange(gg.rnd, gg.quantityBottom, gg.quantityUp)),
		reserve: sdk.NewCoin(
			config.BaseDenom,
			helpers.EtherToWei(sdk.NewInt(RandomRange(gg.rnd, gg.reserveBottom, gg.reserveUp))),
		),
		allowMint: gg.rnd.Int31n(2) == 0,
	}
}

func (aa *MintNFTAction) ChooseAccounts(saList []*stormTypes.StormAccount) []*stormTypes.StormAccount {
	var res []*stormTypes.StormAccount
	for i := range saList {
		if saList[i].IsDirty() {
			continue
		}
		if saList[i].BalanceForCoin(saList[i].FeeDenom()).LT(aa.reserve.Amount) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *MintNFTAction) GenerateTx(sa *stormTypes.StormAccount) ([]byte, error) {
	sender, err := sdk.AccAddressFromBech32(sa.Address())
	if err != nil {
		return nil, err
	}
	// NOTE: due NFT spec, only NFT creator can change NFT reserve
	// so i set recipient to sender for MsgUpdateReserve, MsgBurnNFT
	if rand.Int31n(2) == 0 {
		aa.recipient = sa.Address()
	}
	recipient, err := sdk.AccAddressFromBech32(aa.recipient)
	if err != nil {
		return nil, err
	}

	msg := dscTx.NewMsgMintNFT(
		sender,
		recipient,
		aa.id,
		aa.denom,
		aa.tokenURI,
		aa.quantity,
		aa.reserve,
		aa.allowMint,
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

func (aa *MintNFTAction) String() string {
	return fmt.Sprintf("MintNFT{ID: %s, Recipient: %s, Denom: %s, TokenURI: %s, Quantity: %s, Reserve: %s, AllowMint: %v}",
		aa.id,
		aa.recipient,
		aa.denom,
		aa.tokenURI,
		aa.quantity,
		aa.reserve,
		aa.allowMint,
	)
}
