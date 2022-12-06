package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
	dscApi "bitbucket.org/decimalteam/go-smart-node/sdk/api"
)

func cmdVerify() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "verify-coins",
		Short: "Verify custom coins volume",
		Run: func(cmd *cobra.Command, args []string) {
			//
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			err = reactor.initApi(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			addresses, err := reactor.api.AllAccounts()
			if err != nil {
				fmt.Println(err)
				return
			}
			// coins info from coin keeper
			coinsInfo, err := reactor.coins()
			if err != nil {
				fmt.Println(err)
				return
			}
			balances := sdk.NewCoins()
			for _, adr := range addresses {
				coins, err := reactor.api.AddressBalance(adr)
				if err != nil {
					fmt.Println(err)
					return
				}
				balances = balances.Add(coins...)
			}
			for _, coinInfo := range coinsInfo {
				//if coinInfo.Denom == reactor.api.BaseCoin() {
				//	continue
				//}
				diff := "none"
				bal := balances.AmountOf(coinInfo.Denom)
				if !bal.Equal(coinInfo.Volume) {
					diff = fmt.Sprintf("volume=%s, balances=%s",
						coinInfo.Volume.String(), bal.String())
				}
				fmt.Printf("coin: (symbol: %s, volume: %s, reserve: %s), difference: %s\n",
					coinInfo.Denom, coinInfo.Volume, coinInfo.Reserve, diff)
			}
		},
	}

	return cmd
}

func cmdVerifyCoinsByBank() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "verify-coins-by-bank",
		Short: "Verify custom coins volume",
		Run: func(cmd *cobra.Command, args []string) {
			//
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			err = reactor.initApi(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			addresses, err := reactor.api.AllAccounts()
			if err != nil {
				fmt.Println(err)
				return
			}
			// coins info from coin keeper
			coinsInfo, err := reactor.coins()
			if err != nil {
				fmt.Println(err)
				return
			}
			balances := sdk.NewCoins()
			for _, adr := range addresses {
				coins, err := reactor.api.AddressBalance(adr)
				if err != nil {
					fmt.Println(err)
					return
				}
				balances = balances.Add(coins...)
			}
			for _, coinInfo := range coinsInfo {
				volume, err := reactor.api.GetSupply(coinInfo.Denom)
				if err != nil {
					fmt.Printf("GetSupply(%s) fail: %s\n", coinInfo.Denom, err.Error())
					continue
				}
				diff := "none"
				bal := balances.AmountOf(coinInfo.Denom)
				if !bal.Equal(volume) {
					diff = fmt.Sprintf("volume=%s, balances=%s",
						volume.String(), bal.String())
				}
				fmt.Printf("coin: (symbol: %s, bank volume: %s), difference: %s\n",
					coinInfo.Denom, volume, diff)
			}
		},
	}

	return cmd
}

func cmdValidators() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "validators",
		Short: "Show validators info",
		Run: func(cmd *cobra.Command, args []string) {
			//
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			err = reactor.initApi(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			// validators info
			vals, err := reactor.api.Validators()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, val := range vals {
				dels, err := reactor.api.ValidatorDelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("moniker: %s, address: %s, status: %d, online: %v, jailed: %v, stake: %d, rewards: %s, delegation: %d\n",
					val.Description.Moniker, val.OperatorAddress, val.Status, val.Online, val.Jailed, val.Stake, val.Rewards, len(dels))
			}
		},
	}

	return cmd
}

func cmdVerifyPools() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "verify-pools",
		Short: "Verify (un/re)delegations and validator pools",
		Run: func(cmd *cobra.Command, args []string) {
			//
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			reactor := stormReactor{}
			// init
			err = reactor.initApi(cmd.Flags())
			if err != nil {
				fmt.Println(err)
				return
			}
			// validators info
			vals, err := reactor.api.Validators()
			if err != nil {
				fmt.Println(err)
				return
			}
			ss := newStakeSummator(reactor.api)
			bondedAddress := moduleNameToAddress("bonded_tokens_pool")
			notBondedAddress := moduleNameToAddress("not_bonded_tokens_pool")
			for _, val := range vals {
				// delegations
				dels, err := reactor.api.ValidatorDelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, del := range dels {
					ss.addStake(del.Stake, val.Status)
				}
				// redelegations
				reds, err := reactor.api.ValidatorRedelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, red := range reds {
					for _, ent := range red.Entries {
						ss.addStake(ent.Stake, dscApi.BondStatus_Unbonded)
					}
				}
				// undelegations
				ubds, err := reactor.api.ValidatorUndelegations(val.OperatorAddress)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, ubd := range ubds {
					for _, ent := range ubd.Entries {
						ss.addStake(ent.Stake, dscApi.BondStatus_Unbonded)
					}
				}
			}

			// check pool
			balanceBondedPool, err := reactor.api.AddressBalance(bondedAddress)
			if err != nil {
				fmt.Println(err)
				return
			}
			balanceNotBondedPool, err := reactor.api.AddressBalance(notBondedAddress)
			if err != nil {
				fmt.Println(err)
				return
			}
			compareCoins("bonded_tokens_pool", balanceBondedPool, ss.bonded)
			compareCoins("not_bonded_tokens_pool", balanceNotBondedPool, ss.notBonded)
			if len(ss.nftDiff) > 0 {
				fmt.Printf("nft diffs:\n")
				for _, v := range ss.nftDiff {
					fmt.Printf("%s\n", v)
				}
			} else {
				fmt.Printf("nft owners is correct\n")
			}
			// check all nft
			nfts := make([]*dscApi.NFTToken, 0)
			colls, err := reactor.api.NFTCollections()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, coll := range colls {
				collWithTokens, err := reactor.api.NFTCollection(coll.Creator, coll.Denom)
				if err != nil {
					fmt.Println(err)
					continue
				}
				nfts = append(nfts, collWithTokens.Tokens...)
			}
			for _, nft := range nfts {
				tok, err := reactor.api.NFTToken(nft.ID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, sub := range tok.SubTokens {
					if !ss.nftDelegated[nftKey{nft.ID, sub.ID}] && (sub.Owner == bondedAddress || sub.Owner == notBondedAddress) {
						fmt.Printf("pool is owner of non-delegated token: %s, sub: %d\n", nft.ID, sub.ID)
					}
				}
			}
		},
	}

	return cmd
}

func moduleNameToAddress(name string) string {
	address, err := bech32.ConvertAndEncode(cmdcfg.Bech32Prefix, cosmosAuthTypes.NewModuleAddress(name))
	if err != nil {
		panic(fmt.Sprintf("moduleNameToAddress(%s) = %s", name, err.Error()))
	}
	return address
}

type stakeSummator struct {
	bondedAddress    string
	notBondedAddress string
	bonded           sdk.Coins
	notBonded        sdk.Coins
	nftDiff          []string
	nftOwners        map[nftKey]string
	nftDelegated     map[nftKey]bool
	api              *dscApi.API
}

type nftKey struct {
	tokenID string
	subId   uint32
}

func newStakeSummator(api *dscApi.API) *stakeSummator {
	return &stakeSummator{
		bondedAddress:    moduleNameToAddress("bonded_tokens_pool"),
		notBondedAddress: moduleNameToAddress("not_bonded_tokens_pool"),
		bonded:           sdk.NewCoins(),
		notBonded:        sdk.NewCoins(),
		nftDiff:          []string{},
		nftOwners:        make(map[nftKey]string),
		nftDelegated:     make(map[nftKey]bool),
		api:              api,
	}
}

func (ss *stakeSummator) addNFT(tokenID string, subID uint32, owner string) {
	ss.nftOwners[nftKey{tokenID, subID}] = owner
}

func (ss *stakeSummator) addStake(stake dscApi.Stake, bondStatus dscApi.BondStatus) {
	switch stake.Type {
	case dscApi.StakeType_Coin:
		switch bondStatus {
		case dscApi.BondStatus_Bonded:
			ss.bonded = ss.bonded.Add(stake.Stake)
		case dscApi.BondStatus_Unbonded, dscApi.BondStatus_Unbonding:
			ss.notBonded = ss.notBonded.Add(stake.Stake)
		}
	case dscApi.StakeType_NFT:
		for _, subId := range stake.SubTokenIDs {
			ss.nftDelegated[nftKey{stake.ID, subId}] = true
			sub, err := ss.api.NFTSubToken(stake.ID, fmt.Sprintf("%d", subId))
			if err != nil {
				ss.nftDiff = append(ss.nftDiff, fmt.Sprintf("token: %s, sub id: %d, error: %s", stake.ID, subId, err.Error()))
			}
			switch bondStatus {
			case dscApi.BondStatus_Bonded:
				if sub.Owner != ss.bondedAddress {
					ss.nftDiff = append(ss.nftDiff, fmt.Sprintf("token: %s, sub id: %d, expect 'bonded' owner", stake.ID, subId))
				}
			case dscApi.BondStatus_Unbonded, dscApi.BondStatus_Unbonding:
				if sub.Owner != ss.notBondedAddress {
					ss.nftDiff = append(ss.nftDiff, fmt.Sprintf("token: %s, sub id: %d, expect 'not bonded' owner", stake.ID, subId))
				}
			}
		}

	}
}

func compareCoins(name string, coins1, coins2 sdk.Coins) {
	if coins1.IsEqual(coins2) {
		fmt.Printf("pool '%s' is correct\n", name)
	} else {
		fmt.Printf("pool '%s' differs:\n", name)
		denoms := make([]string, 0)
		for _, coin := range coins1.Add(coins2...) {
			denoms = append(denoms, coin.Denom)
		}
		sort.Strings(denoms)
		for _, denom := range denoms {
			if !coins1.AmountOf(denom).Equal(coins2.AmountOf(denom)) {
				fmt.Printf("denom '%s' %s != %s\n", denom, coins1.AmountOf(denom), coins2.AmountOf(denom))
			}
		}
	}
}
