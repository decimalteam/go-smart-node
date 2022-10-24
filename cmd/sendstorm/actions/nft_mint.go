package actions

import (
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stormTypes "bitbucket.org/decimalteam/go-smart-node/cmd/sendstorm/types"
	dscTx "bitbucket.org/decimalteam/go-smart-node/sdk/tx"
	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
)

// MsgMintNFT
type MintNFTGenerator struct {
	// general values
	textLengthBottom, textLengthUp int64
	reserveBottom, reserveUp       int64 // in 10^18
	quantityBottom, quantityUp     int64
	knownAddresses                 []string
	knownCoins                     []string
	rnd                            *rand.Rand
}

type MintNFTAction struct {
	recipient string
	id        string
	denom     string
	tokenURI  string
	quantity  uint32
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
	gg.knownCoins = ui.Coins
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
		quantity:  uint32(RandomRange(gg.rnd, gg.quantityBottom, gg.quantityUp)),
		reserve: sdk.NewCoin(
			RandomChoice(gg.rnd, gg.knownCoins),
			helpers.EtherToWei(sdkmath.NewInt(RandomRange(gg.rnd, gg.reserveBottom, gg.reserveUp))),
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
		if saList[i].BalanceForCoin(aa.reserve.Denom).LT(aa.reserve.Amount.Mul(sdkmath.NewInt(int64(aa.quantity)))) {
			continue
		}
		res = append(res, saList[i])
	}
	return res
}

func (aa *MintNFTAction) GenerateTx(sa *stormTypes.StormAccount, feeConfig *stormTypes.FeeConfiguration) ([]byte, error) {
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

	msg := dscTx.NewMsgMintToken(
		sender,
		aa.denom,
		aa.id,
		aa.tokenURI,
		aa.allowMint,
		recipient,
		aa.quantity,
		aa.reserve,
	)

	return feeConfig.MakeTransaction(sa, msg)
}

func (aa *MintNFTAction) String() string {
	return fmt.Sprintf("MintNFT{ID: %s, Recipient: %s, Denom: %s, TokenURI: %s, Quantity: %d, Reserve: %s, AllowMint: %v}",
		aa.id,
		aa.recipient,
		aa.denom,
		aa.tokenURI,
		aa.quantity,
		aa.reserve,
		aa.allowMint,
	)
}
