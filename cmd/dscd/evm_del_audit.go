package main

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"text/tabwriter"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/decimalteam/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"bitbucket.org/decimalteam/go-smart-node/app"
	"bitbucket.org/decimalteam/go-smart-node/app/redenom"
	"bitbucket.org/decimalteam/go-smart-node/cmd/dscd/stakescan"
	"bitbucket.org/decimalteam/go-smart-node/contracts"
	nfttypes "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

var evmEmptyCodeHash = crypto.Keccak256(nil)

type unkContract struct {
	addr     common.Address
	bal      sdkmath.Int
	codeHash string
}

// EVMDelAuditCmd audits where native "del" reserves live in EVM: it classifies every
// CONTRACT holding a bank "del" balance and flags any that the redenomination doesn't
// already account for (WDEL, custom-coin DRC20 reserves, the DecimalChecks contract,
// system contracts, and x/nft-registered EVM NFT collections). Any UNKNOWN del-holding
// contract is a potential EVM NFT (or other) reserve whose internal ledger the upgrade
// would leave unscaled. For each x/nft EVM collection it also cross-checks
// Σ _reserve[tokenId].amount (DEL) against the contract's native "del" balance — a
// mismatch means the x/nft token mirror is missing tokenIds the EVM contract holds.
//
// Read-only; the node must be stopped.
func EVMDelAuditCmd(encCfg params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evm-del-audit",
		Short: "Audit where EVM contracts hold native DEL reserves (esp. NFT collections) before redenomination",
		Long: `Classifies every EVM contract holding a bank "del" balance and flags any not
accounted for by the redenomination, and cross-checks each x/nft EVM NFT collection's
on-chain _reserve[tokenId].amount sum against its native del backing. The node must be stopped.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := sdkserver.GetServerContextFromCmd(cmd)
			appOpts := serverCtx.Viper
			home := cast.ToString(appOpts.Get(flags.FlagHome))
			if home == "" {
				return fmt.Errorf("--home is required")
			}
			db, err := openRedenomAppDB(home, appOpts)
			if err != nil {
				return fmt.Errorf("open application db: %w", err)
			}
			defer db.Close()

			a := app.NewDSC(serverCtx.Logger, db, nil, true, map[int64]bool{}, home, uint(1), encCfg, appOpts)
			height := a.LastBlockHeight()
			chainID := readChainID(home)
			ctx := a.NewContext(true, tmproto.Header{Height: height, ChainID: chainID})
			base := a.ValidatorKeeper.BaseDenom(ctx)
			kv := ctx.KVStore(a.GetKey(evmtypes.StoreKey))

			// --inspect: dump one contract's storage + check whether any slot mirrors its
			// native del balance (the signature of an internal DEL ledger that must be scaled).
			if insp, _ := cmd.Flags().GetString("inspect"); insp != "" {
				addr := common.HexToAddress(insp)
				st := scanStorage(kv, addr)
				native := a.BankKeeper.GetBalance(ctx, sdk.AccAddress(addr.Bytes()), base).Amount
				codeLen := 0
				if acc := a.EvmKeeper.GetAccount(ctx, addr); acc != nil {
					codeLen = len(a.EvmKeeper.GetCode(ctx, common.BytesToHash(acc.CodeHash)))
				}
				fmt.Fprintf(cmd.OutOrStdout(), "inspect %s\n  native del = %s\n  code bytes = %d\n  storage slots = %d\n", addr.Hex(), native, codeLen, len(st))
				type sv struct {
					k common.Hash
					v *big.Int
				}
				var svs []sv
				mirrors := false
				for k, v := range st {
					val := new(big.Int).SetBytes(v.Bytes())
					svs = append(svs, sv{k, val})
					if val.Sign() != 0 && val.String() == native.String() {
						mirrors = true
					}
				}
				sort.Slice(svs, func(i, j int) bool { return svs[i].v.Cmp(svs[j].v) > 0 })
				for i, x := range svs {
					if i >= 25 {
						fmt.Fprintf(cmd.OutOrStdout(), "  ... (%d more slots)\n", len(svs)-25)
						break
					}
					fmt.Fprintf(cmd.OutOrStdout(), "  %s = %s\n", x.k.Hex(), x.v)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  >> a slot equals native del balance: %v (true ⇒ likely an internal DEL ledger needing scaling)\n", mirrors)
				return nil
			}

			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(out, "EVM DEL-reserve audit\tchain-id %s\theight %d\tbase %q\n", chainID, height, base)

			// Known contract universe.
			label := map[common.Address]string{}
			var nftCenterStore stakescan.Storage
			// DecimalNFTCenter._nfts mapping (field 0 of this ERC-7201 struct → field slot == base).
			nftCenterBase := common.HexToHash("0xc21943c83a3e2d63fb720a275e35cd8758d28f083d0aa8e8050823efbb0ffb00")
			if cc, ok := redenom.ContractCenterFor(chainID); ok {
				label[cc] = "contract-center"
				ccStore := scanStorage(kv, cc)
				for _, sl := range []string{
					contracts.NameOfSlugForGetAddressDelegation,
					contracts.NameOfSlugForGetAddressDelegationNft,
					contracts.NameOfSlugForGetAddressMasterValidator,
					contracts.NameOfSlugForGetAddressNftCenter,
					contracts.NameOfSlugForGetAddressTokenCenter,
					contracts.NameOfSlugForGetAddressWDEL,
				} {
					if addr := stakescan.ResolveAddressBySymbol(ccStore, sl); addr != (common.Address{}) {
						label[addr] = sl
						if sl == contracts.NameOfSlugForGetAddressNftCenter {
							nftCenterStore = scanStorage(kv, addr)
						}
					}
				}
			} else {
				fmt.Fprintf(out, "WARNING\tchain-id %q is an unknown network family (not in redenom.ContractCenterFor) — system contracts unlabeled\n", chainID)
			}
			_ = nftCenterStore
			_ = nftCenterBase
			// isNFTCollection reports whether addr is an EVM NFT collection: a DecimalNFTBeaconProxy
			// has the EIP-1967 beacon slot set. NFT collections inherit NFTReserve and hold DEL
			// reserves as native "del"; a del-holding beacon proxy therefore has a DEL _reserve
			// ledger. Robust and independent of x/nft (legacy) and the nft-center _nfts read.
			beaconSlot := common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")
			isNFTCollection := func(addr common.Address) bool {
				return a.EvmKeeper.GetState(ctx, addr, beaconSlot) != (common.Hash{})
			}
			// Custom-coin DRC20 reserve contracts.
			for _, c := range a.CoinKeeper.GetCoins(ctx) {
				if c.DRC20Contract != "" {
					label[common.HexToAddress(c.DRC20Contract)] = "coin-drc20:" + c.Denom
				}
			}
			// x/nft EVM NFT collections (AddressDRC mirror).
			type nftColl struct {
				creator sdk.AccAddress
				denom   string
			}
			nftColls := map[common.Address]nftColl{}
			for _, col := range a.NFTKeeper.GetCollections(ctx) {
				if col.AddressDRC == "" {
					continue
				}
				creator, err := sdk.AccAddressFromBech32(col.Creator)
				if err != nil {
					continue
				}
				addr := common.HexToAddress(col.AddressDRC)
				label[addr] = "nft-collection"
				nftColls[addr] = nftColl{creator: creator, denom: col.Denom}
			}

			// Classify every del-holding CONTRACT.
			type cat struct {
				count int
				total sdkmath.Int
			}
			cats := map[string]*cat{}
			eoaCount := 0
			eoaTotal := sdkmath.ZeroInt()
			var unknowns []unkContract
			add := func(k string, amt sdkmath.Int) {
				c := cats[k]
				if c == nil {
					c = &cat{total: sdkmath.ZeroInt()}
					cats[k] = c
				}
				c.count++
				c.total = c.total.Add(amt)
			}
			a.BankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, coin sdk.Coin) bool {
				if coin.Denom != base || len(addr.Bytes()) != 20 {
					return false
				}
				eth := common.BytesToAddress(addr.Bytes())
				acc := a.EvmKeeper.GetAccount(ctx, eth)
				if acc == nil || bytes.Equal(acc.CodeHash, evmEmptyCodeHash) {
					eoaCount++
					eoaTotal = eoaTotal.Add(coin.Amount)
					return false
				}
				lbl, ok := label[eth]
				if !ok {
					lbl = "UNKNOWN"
					unknowns = append(unknowns, unkContract{addr: eth, bal: coin.Amount, codeHash: common.Bytes2Hex(acc.CodeHash)})
				}
				// collapse coin-drc20:* into one bucket for the summary
				bucket := lbl
				if len(lbl) > 10 && lbl[:10] == "coin-drc20" {
					bucket = "coin-drc20"
				}
				add(bucket, coin.Amount)
				return false
			})

			// Summary by category.
			fmt.Fprintf(out, "\n-- del-holding contracts by category --\n")
			fmt.Fprintf(out, "category\tcontracts\ttotal del (wei)\n")
			keys := make([]string, 0, len(cats))
			for k := range cats {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(out, "%s\t%d\t%s\n", k, cats[k].count, cats[k].total)
			}
			fmt.Fprintf(out, "(EOAs)\t%d\t%s\n", eoaCount, eoaTotal)

			// Classify the unknown del-holding contracts: which are EVM NFT collections
			// registered in nft-center (the critical case — DEL reserves whose _reserve ledger
			// the upgrade would leave unscaled), and group the rest by bytecode.
			var nftUnknowns []unkContract
			nftTotal, otherTotal := sdkmath.ZeroInt(), sdkmath.ZeroInt()
			type codeGroup struct {
				n     int
				total sdkmath.Int
			}
			codeGroups := map[string]codeGroup{}
			for _, u := range unknowns {
				g := codeGroups[u.codeHash]
				if g.total.IsNil() {
					g.total = sdkmath.ZeroInt()
				}
				g.n++
				g.total = g.total.Add(u.bal)
				codeGroups[u.codeHash] = g
				if isNFTCollection(u.addr) {
					nftUnknowns = append(nftUnknowns, u)
					nftTotal = nftTotal.Add(u.bal)
				} else {
					otherTotal = otherTotal.Add(u.bal)
				}
			}

			fmt.Fprintf(out, "\n-- UNKNOWN del-holding contracts: classification --\n")
			fmt.Fprintf(out, "EVM NFT collections (beacon proxy, hold DEL reserve)\t%d contracts\tdel %s\n", len(nftUnknowns), nftTotal)
			fmt.Fprintf(out, "other unknown contracts (wallets / checks / etc.)\t%d contracts\tdel %s\n", len(unknowns)-len(nftUnknowns), otherTotal)

			// Distinct bytecodes among the unknowns (a single dominant code hash usually means
			// one contract type, e.g. a wallet template).
			type cg struct {
				hash  string
				n     int
				total sdkmath.Int
			}
			var cgs []cg
			for h, g := range codeGroups {
				cgs = append(cgs, cg{h, g.n, g.total})
			}
			sort.Slice(cgs, func(i, j int) bool { return cgs[i].total.GT(cgs[j].total) })
			fmt.Fprintf(out, "\n-- unknown contracts by bytecode (%d distinct) --\n", len(cgs))
			for i, g := range cgs {
				if i >= 12 {
					fmt.Fprintf(out, "... (%d more bytecode groups)\n", len(cgs)-12)
					break
				}
				short := g.hash
				if len(short) > 16 {
					short = short[:16]
				}
				fmt.Fprintf(out, "code %s…\t%d contracts\tdel %s\n", short, g.n, g.total)
			}

			if len(nftUnknowns) > 0 {
				fmt.Fprintf(out, "\n!! CRITICAL: %d EVM NFT collections hold DEL reserves but are NOT mirrored in x/nft — their _reserve ledgers would be left UNSCALED (insolvency). Sample:\n", len(nftUnknowns))
				sort.Slice(nftUnknowns, func(i, j int) bool { return nftUnknowns[i].bal.GT(nftUnknowns[j].bal) })
				for i, u := range nftUnknowns {
					if i >= 20 {
						fmt.Fprintf(out, "... (%d more)\n", len(nftUnknowns)-20)
						break
					}
					fmt.Fprintf(out, "%s\tdel=%s\n", u.addr.Hex(), u.bal)
				}
			}

			// Cosmos NFT reserve backing: the x/nft subtoken DEL reserves (which cosmos.go
			// scales) must be backed by the reserved_pool module account (bank "del", scaled by
			// scaleBankDel) — NOT by EVM collection contracts. reserved_pool ≈ Σ subtoken DEL
			// reserves ⇒ these NFTs are Cosmos-backed and fully handled. reserved_pool ≈ 0 while
			// Σ > 0 ⇒ the reserves are EVM-backed and the x/nft mirror is incomplete.
			reservedPool := a.AccountKeeper.GetModuleAddress(nfttypes.ReservedPool)
			poolDel := a.BankKeeper.GetBalance(ctx, reservedPool, base).Amount
			sumSub := sdkmath.ZeroInt()
			subCount := 0
			for _, col := range a.NFTKeeper.GetCollections(ctx) {
				creator, err := sdk.AccAddressFromBech32(col.Creator)
				if err != nil {
					continue
				}
				for _, tok := range a.NFTKeeper.GetTokens(ctx, creator, col.Denom) {
					for _, st := range a.NFTKeeper.GetSubTokens(ctx, tok.ID) {
						if st.Reserve != nil && st.Reserve.Denom == base {
							sumSub = sumSub.Add(st.Reserve.Amount)
							subCount++
						}
					}
				}
			}
			fmt.Fprintf(out, "\n-- Cosmos x/nft reserve backing (reserved_pool vs Σ subtoken DEL reserves) --\n")
			fmt.Fprintf(out, "reserved_pool del\t%s\n", poolDel)
			fmt.Fprintf(out, "Σ subtoken DEL reserves (%d subtokens)\t%s\n", subCount, sumSub)
			fmt.Fprintf(out, "match (Cosmos-backed, handled by cosmos.go + bank)\t%v\n", poolDel.Equal(sumSub))

			// Per-collection NFT reserve consistency: Σ on-chain _reserve[tokenId].amount (DEL)
			// vs the collection contract's native del balance.
			fmt.Fprintf(out, "\n-- x/nft EVM NFT collections (AddressDRC set): reserve consistency --\n")
			if len(nftColls) == 0 {
				fmt.Fprintf(out, "none — no x/nft collection has a non-empty AddressDRC (no EVM-bridged NFT collections)\n")
			} else {
				for addr, info := range nftColls {
					native := a.BankKeeper.GetBalance(ctx, sdk.AccAddress(addr.Bytes()), base).Amount
					st := scanStorage(kv, addr)
					sumDEL := sdkmath.ZeroInt()
					delTokens := 0
					for _, tok := range a.NFTKeeper.GetTokens(ctx, info.creator, info.denom) {
						tokenID, ok := new(big.Int).SetString(tok.ID, 10)
						if !ok {
							continue
						}
						rtype := new(big.Int).SetBytes(st[stakescan.NFTReserveTypeSlot(tokenID)].Bytes())
						if rtype.Uint64() != 1 { // 1 == DEL
							continue
						}
						amt := new(big.Int).SetBytes(st[stakescan.NFTReserveAmountSlot(tokenID)].Bytes())
						sumDEL = sumDEL.Add(sdkmath.NewIntFromBigInt(amt))
						delTokens++
					}
					status := "OK"
					if !sumDEL.Equal(native) {
						status = "MISMATCH — x/nft token mirror likely incomplete; some _reserve tokenIds unscaled by redenom"
					}
					fmt.Fprintf(out, "%s\tdelTokens=%d\tΣreserve=%s\tnative=%s\t%s\n",
						addr.Hex(), delTokens, sumDEL, native, status)
				}
			}

			return out.Flush()
		},
	}
	cmd.Flags().String(flags.FlagHome, app.DefaultNodeHome, "node home directory")
	cmd.Flags().String("inspect", "", "inspect a single contract's storage instead of auditing (hex address)")
	return cmd
}
