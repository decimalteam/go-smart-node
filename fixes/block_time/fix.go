package block_time

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"net/http"
)

var fixes map[int64]int64

func init() {
	resp, err := http.Get("https://repo.decimalchain.com/block_time_fix.json")
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &fixes)
	if err != nil {
		panic(err)
	}
}

func GetFixedBlockTime(ctx sdk.Context) int64 {
	var fixedTime = ctx.BlockTime().Unix() + 5
	if ctx.BlockTime().Nanosecond() > 500000000 {
		fixedTime += 1
	}

	if fix, ok := fixes[ctx.BlockHeight()]; ok {
		fixedTime = fix
	}

	return fixedTime
}
