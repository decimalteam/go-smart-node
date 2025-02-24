package block_time

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"net/http"
	"strconv"
	"time"
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
	if fix, ok := fixes[ctx.BlockHeight()]; ok {
		return fix
	}

	return getBlockTime(ctx.BlockHeight() + 1)
}

func getBlockTime(height int64) int64 {
	resp, err := http.Get("https://node.decimalchain.com/rpc/block?height=" + strconv.Itoa(int(height)))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	block := blockResponse{}
	err = json.Unmarshal(body, &block)
	if err != nil {
		panic(err)
	}

	return block.Result.Block.Header.Time.Unix()
}

type blockResponse struct {
	Result struct {
		Block struct {
			Header struct {
				Time time.Time `json:"time"`
			} `json:"header"`
		} `json:"block"`
	} `json:"result"`
}
