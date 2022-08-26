package types

import (
	"time"
)

const (
	DefaultLockedTimeOut = time.Hour * 24
	DefaultLockedTimeIn  = time.Hour * 12
)

const DefaultCheckingAddress = "18fa71ffcf736d5ec0d06f2330a33b4f85a6d69f"

//TODO: actualize for new addresses
//const ChainActivatorAddress = "dx16aeq4ypsx5ar4076v507ch5z8ryd6usx32tnru"
const DefaultSwapServiceAddress = "dx1jqx7chw0faswfmw78cdejzzery5akzmk5zc5x5"

const SwapPool = "atomic_swap_pool"
