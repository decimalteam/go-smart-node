package multisig_test

import (
	"regexp"
	"runtime"
	"testing"

	"bitbucket.org/decimalteam/go-smart-node/app"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/decimalteam/ethermint/encoding"
)

func FuzzPossibleProtoErrors(f *testing.F) {
	var skips = []*regexp.Regexp{
		regexp.MustCompile("unexpected EOF"),
		regexp.MustCompile(`proto: illegal wireType \d+`),
		regexp.MustCompile(`proto: \S+: wiretype end group for non-group`),
		regexp.MustCompile(`proto: \S+: illegal tag [-\d]+ \(wire type \d+\)`),
		regexp.MustCompile(`proto: wrong wireType = \d+ for field \S+`),
		regexp.MustCompile(`proto: integer overflow`),
		regexp.MustCompile(`proto: negative length found during unmarshaling`),
		regexp.MustCompile(`math/big: cannot unmarshal ".+?" into a \*big\.Int`),
		regexp.MustCompile(`no concrete type registered for type URL.*`),
	}

	var urls = []string{
		"/decimal.coin.v1.MsgCreateCoin",
		"/decimal.coin.v1.MsgUpdateCoin",
		"/decimal.coin.v1.MsgMultiSendCoin",
		"/decimal.coin.v1.MsgBuyCoin",
		"/decimal.coin.v1.MsgSellCoin",
		"/decimal.coin.v1.MsgSendCoin",
	}

	cfg := encoding.MakeConfig(app.ModuleBasics)
	f.Add([]byte{}, 0)
	f.Fuzz(func(t *testing.T, a []byte, b int) {
		if b < 0 || b >= len(urls) {
			return
		}
		runtime.GC()
		any := codectypes.Any{
			TypeUrl: urls[b],
			Value:   a,
		}
		var msg sdk.Msg
		err := cfg.Codec.UnpackAny(&any, &msg)
		if err != nil {
			s := err.Error()
			for _, skip := range skips {
				if skip.MatchString(s) {
					return
				}
			}
			t.Errorf("some errors '%s'", err.Error())
		}
	})

}
