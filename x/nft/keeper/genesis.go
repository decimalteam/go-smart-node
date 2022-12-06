package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k Keeper, gs *types.GenesisState) {
	// Initialize params
	k.SetParams(ctx, gs.Params)

	if gs.Collections == nil || len(gs.Collections) < 1 {
		return
	}

	// iterate collections
	for _, collection := range gs.Collections {
		// validate collections
		if err := collection.Validate(); err != nil {
			panic(err)
		}
		// check for duplicate collections
		creator, err := sdk.AccAddressFromBech32(collection.Creator)
		if err != nil {
			panic(err)
		}
		if _, found := k.GetCollection(ctx, creator, collection.Denom); found {
			panic(errors.InvalidCollection) // TODO errors
		}

		// check for unvalid supply
		if int(collection.Supply) != len(collection.Tokens) {
			panic(errors.InvalidCollection)
		}

		// iterate tokens
		for _, token := range collection.Tokens {
			// validate token
			if err := token.Validate(); err != nil {
				panic(err)
			}
			// check for duplicate token IDs
			if _, found := k.GetToken(ctx, token.ID); found {
				panic(errors.NotUniqueTokenID)
			}
			// check for duplicate token URIs
			if exists := k.hasTokenURI(ctx, token.URI); exists {
				panic(errors.NotUniqueTokenURI)
			}

			// iterate subtokens
			for _, subToken := range token.SubTokens {
				// validate subtoken
				if err := subToken.Validate(); err != nil {
					panic(err)
				}
				// check for duplicate subtokens
				if _, found := k.GetSubToken(ctx, token.ID, subToken.ID); found {
					panic(errors.NotUniqueSubTokenIDs)
				}
				// reserve validity check
				// TODO: old subtokens can be slashed
				/*
					if subToken.Reserve.IsLT(token.Reserve) {
						panic(errors.InvalidReserve)
					}
				*/
				owner1, err1 := sdk.AccAddressFromBech32(subToken.Owner)
				owner2, err2 := sdk.GetFromBech32(subToken.Owner, "dx") // may be legacy address
				if err1 != nil && err2 != nil {
					panic(err)
				}
				var owner sdk.AccAddress
				if err1 == nil {
					owner = owner1
				}
				if err2 == nil {
					owner = sdk.AccAddress(owner2)
				}
				// write sub-token record
				k.SetSubToken(ctx, token.ID, *subToken)
				// write sub-token by owner index
				k.setSubTokenByOwner(ctx, owner, token.ID, subToken.ID)
			}

			k.CreateToken(ctx, collection, *token)
		}

		k.SetCollection(ctx, collection)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:      k.GetParams(ctx),
		Collections: k.GetCollections(ctx),
	}
}
