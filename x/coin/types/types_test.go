package types

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cosmos/cosmos-sdk/codec"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"bitbucket.org/decimalteam/go-smart-node/utils/helpers"
// )

// func TestChangeCoinStruct(t *testing.T) {
// 	volume := helpers.BipToPip(sdk.NewInt(100000))
// 	reserve := helpers.BipToPip(sdk.NewInt(100000))

// 	type Coin struct {
// 		Title       string  `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
// 		CRR         uint    `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
// 		Symbol      string  `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
// 		Reserve     sdk.Int `json:"reserve" yaml:"reserve"`
// 		LimitVolume sdk.Int `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
// 		Volume      sdk.Int `json:"volume" yaml:"volume"`
// 	}

// 	coin := Coin{
// 		Title:       "TEST COIN",
// 		CRR:         10,
// 		Symbol:      "test",
// 		Reserve:     reserve,
// 		LimitVolume: volume.Mul(sdk.NewInt(10)),
// 		Volume:      volume,
// 	}

// 	cdc := codec.New()
// 	value := cdc.MustMarshalBinaryLengthPrefixed(coin)
// 	fmt.Println(value)

// 	type NewCoin struct {
// 		Title       string         `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
// 		CRR         uint           `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
// 		Symbol      string         `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
// 		Reserve     sdk.Int        `json:"reserve" yaml:"reserve"`
// 		LimitVolume sdk.Int        `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
// 		Volume      sdk.Int        `json:"volume" yaml:"volume"`
// 		Creator     sdk.AccAddress `json:"creator" yaml:"creator"`
// 		Icon        string         `json:"icon" yaml:"icon"`
// 	}

// 	newCoin := NewCoin{}
// 	cdc.MustUnmarshalBinaryLengthPrefixed(value, &newCoin)
// 	value = cdc.MustMarshalBinaryLengthPrefixed(newCoin)
// 	fmt.Println(value)
// }
