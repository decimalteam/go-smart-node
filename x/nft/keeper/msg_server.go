package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-smart-node/utils/events"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/errors"
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
)

var _ types.MsgServer = &Keeper{}

////////////////////////////////////////////////////////////////
// MintToken
////////////////////////////////////////////////////////////////

func (k Keeper) MintToken(c context.Context, msg *types.MsgMintToken) (*types.MsgMintTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	recipient := sdk.MustAccAddressFromBech32(msg.Recipient)

	// retrieve NFT collection
	collection, collectionExists := k.GetCollection(ctx, sender, msg.Denom)
	if !collectionExists {
		// create NFT collection
		collection = types.Collection{
			Creator: msg.Sender,
			Denom:   msg.Denom,
		}
	} else {
		// ensure creator address
		if collection.Creator != msg.Sender {
			return nil, errors.NotCreatorMint
		}
	}

	// retrieve NFT token
	token, tokenExists := k.GetToken(ctx, msg.TokenID)
	if !tokenExists {
		// ensure new NFT token is valid
		if k.hasTokenURI(ctx, msg.TokenURI) {
			return nil, errors.NotUniqueTokenURI
		}
		if msg.Reserve.Amount.LT(k.GetMinReserve(ctx)) {
			return nil, errors.InvalidReserve
		}
		// create new NFT token
		token = types.Token{
			Creator:   msg.Sender,
			ID:        msg.TokenID,
			URI:       msg.TokenURI,
			Reserve:   msg.Reserve,
			AllowMint: msg.AllowMint,
		}
	} else {
		// ensure additional minting is allowed
		if token.Creator != msg.Sender {
			return nil, errors.NotCreatorMint
		}
		if !token.AllowMint {
			return nil, errors.NotAllowedMint
		}

		// make sure the same denom is used
		if token.Reserve.Denom != msg.Reserve.Denom {
			return nil, errors.WrongReserveCoinDenom
		}
	}

	if !collectionExists && tokenExists {
		return nil, errors.UnknownCollection
	}

	// prepare new NFT sub-tokens
	subTokenIDs := make([]uint32, msg.Quantity)
	subTokens := make([]types.SubToken, msg.Quantity)
	for i, c, o := uint32(0), msg.Quantity, token.Minted+1; i < c; i++ {
		subTokenIDs[i] = o + i
		subTokens[i] = types.SubToken{
			ID:    o + i,
			Owner: recipient.String(),
		}
	}

	// update NFT collection in the store
	if !collectionExists {
		collection.Supply = 1
		// write collection with it's counter
		k.SetCollection(ctx, collection)
	} else if !tokenExists {
		collection.Supply++
		// write collection counter separately
		k.setCollectionCounter(ctx, sender, collection.Denom, types.CollectionCounter{
			Supply: collection.Supply,
		})
	}

	// update NFT token in the store
	token.Minted += msg.Quantity
	if !tokenExists {
		// write token with it's counter and indexes
		k.CreateToken(ctx, collection, token)
	} else {
		// write token counter separately
		k.setTokenCounter(ctx, token.ID, types.TokenCounter{
			Minted: token.Minted,
			Burnt:  token.Burnt,
		})
	}

	// write new NFT sub-tokens to the store
	for _, subToken := range subTokens {
		// write sub-token record
		k.SetSubToken(ctx, token.ID, subToken)
		// write sub-token by owner index
		k.setSubTokenByOwner(ctx, recipient, token.ID, subToken.ID)
	}

	// calculate needed amount of coins to reserve
	reserveAmount := msg.Reserve.Amount
	if tokenExists {
		reserveAmount = token.Reserve.Amount
	}
	reserveCoin := sdk.NewCoin(msg.Reserve.Denom, reserveAmount.Mul(sdkmath.NewInt(int64(msg.Quantity))))

	// reserve needed amount of coins
	err := k.ReserveTokens(ctx, sdk.NewCoins(reserveCoin), sender)
	if err != nil {
		return nil, errors.InsufficientFunds
	}

	// emit NFT collection events
	if !collectionExists {
		// emit create NFT collection event
		err = events.EmitTypedEvent(ctx, &types.EventCreateCollection{
			Creator: collection.Creator,
			Denom:   collection.Denom,
			Supply:  collection.Supply,
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("err: %s", err.Error())
		}
	} else {
		// emit update NFT collection event
		err = events.EmitTypedEvent(ctx, &types.EventUpdateCollection{
			Creator: collection.Creator,
			Denom:   collection.Denom,
			Supply:  collection.Supply,
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("err: %s", err.Error())
		}
	}

	// emit NFT token events
	if !tokenExists {
		// emit create NFT token event
		err = events.EmitTypedEvent(ctx, &types.EventCreateToken{
			Creator:     collection.Creator,
			Denom:       collection.Denom,
			ID:          msg.TokenID,
			URI:         msg.TokenURI,
			AllowMint:   msg.AllowMint,
			Reserve:     msg.Reserve.String(),
			Recipient:   msg.Recipient,
			SubTokenIDs: subTokenIDs,
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("err: %s", err.Error())
		}
	} else {
		// emit mint NFT token event
		err = events.EmitTypedEvent(ctx, &types.EventMintToken{
			Creator:     collection.Creator,
			Denom:       collection.Denom,
			ID:          msg.TokenID,
			Reserve:     token.Reserve.String(),
			Recipient:   msg.Recipient,
			SubTokenIDs: subTokenIDs,
		})
		if err != nil {
			return nil, errors.Internal.Wrapf("err: %s", err.Error())
		}
	}

	return &types.MsgMintTokenResponse{}, nil
}

////////////////////////////////////////////////////////////////
// UpdateToken
////////////////////////////////////////////////////////////////

func (k Keeper) UpdateToken(c context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// retrieve NFT token
	token, tokenExists := k.GetToken(ctx, msg.TokenID)
	if !tokenExists {
		return nil, errors.UnknownNFT
	}

	// ensure creator address
	if token.Creator != msg.Sender {
		return nil, errors.NotCreatorUpdate
	}
	// ensure token URI differs and unique
	if token.URI == msg.TokenURI {
		return nil, errors.SameTokenURI
	}
	if k.hasTokenURI(ctx, msg.TokenURI) {
		return nil, errors.NotUniqueTokenURI
	}

	// update token URI in token
	token.URI = msg.TokenURI
	k.setToken(ctx, token)
	// update token URI indexes in the KVStore
	k.updateTokenURI(ctx, token.URI, msg.TokenURI)

	// emit NFT token update event
	err := events.EmitTypedEvent(ctx, &types.EventUpdateToken{
		Sender: msg.Sender,
		ID:     msg.TokenID,
		URI:    msg.TokenURI,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUpdateTokenResponse{}, nil
}

////////////////////////////////////////////////////////////////
// UpdateReserve
////////////////////////////////////////////////////////////////

func (k Keeper) UpdateReserve(c context.Context, msg *types.MsgUpdateReserve) (*types.MsgUpdateReserveResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	// retrieve NFT token
	token, tokenExists := k.GetToken(ctx, msg.TokenID)
	if !tokenExists {
		return nil, errors.UnknownNFT
	}

	// ensure creator address
	if token.Creator != msg.Sender {
		return nil, errors.NotCreatorUpdateReserve
	}
	// ensure reserve coin denom
	if token.Reserve.Denom != msg.Reserve.Denom {
		return nil, errors.WrongReserveCoinDenom
	}

	// update NFT sub-token to refill in total
	refillAmount := sdkmath.ZeroInt()
	for _, subTokenID := range msg.SubTokenIDs {
		// retrieve NFT sub-token
		subToken, found := k.GetSubToken(ctx, token.ID, subTokenID)
		if !found {
			return nil, errors.UnknownSubTokenForNFT
		}
		// ensure NFT sub-token is owned by the sender
		if subToken.Owner != msg.Sender {
			return nil, errors.NotCreatorUpdateReserve
		}
		// ensure new reserve is valid
		reserve, newReserve := token.Reserve.Amount, msg.Reserve.Amount
		if subToken.Reserve != nil {
			reserve = subToken.Reserve.Amount
		}
		if newReserve.LTE(reserve) {
			return nil, errors.NotSetValueLowerNow
		}
		// update NFT sub-token
		refillAmount = refillAmount.Add(newReserve.Sub(reserve))
		k.SetSubToken(ctx, token.ID, types.SubToken{
			ID:      subToken.ID,
			Owner:   subToken.Owner,
			Reserve: &msg.Reserve,
		})
	}

	// reserve needed amount of coins
	refillCoin := sdk.NewCoin(token.Reserve.Denom, refillAmount)
	err := k.ReserveTokens(ctx, sdk.NewCoins(refillCoin), sender)
	if err != nil {
		return nil, errors.InsufficientFunds
	}

	// emit NFT token reserve update event
	err = events.EmitTypedEvent(ctx, &types.EventUpdateReserve{
		Sender:      msg.Sender,
		ID:          msg.TokenID,
		Reserve:     msg.Reserve.String(),
		Refill:      refillCoin.String(),
		SubTokenIDs: msg.SubTokenIDs,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgUpdateReserveResponse{}, nil
}

////////////////////////////////////////////////////////////////
// SendToken
////////////////////////////////////////////////////////////////

func (k Keeper) SendToken(c context.Context, msg *types.MsgSendToken) (*types.MsgSendTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	recipient := sdk.MustAccAddressFromBech32(msg.Recipient)

	err := k.TransferSubTokens(ctx, sender, recipient, msg.TokenID, msg.GetSubTokenIDs())
	if err != nil {
		return nil, err
	}

	return &types.MsgSendTokenResponse{}, nil
}

////////////////////////////////////////////////////////////////
// BurnToken
////////////////////////////////////////////////////////////////

func (k Keeper) BurnToken(c context.Context, msg *types.MsgBurnToken) (*types.MsgBurnTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	subTokenCount := uint32(len(msg.SubTokenIDs))

	// retrieve NFT token
	token, tokenExists := k.GetToken(ctx, msg.TokenID)
	if !tokenExists {
		return nil, errors.UnknownNFT
	}

	// ensure creator address
	// TODO: only creator? can owner to burn NFT own sub-tokens?
	// TODO: why creator can burn NFT sub-tokens that are not owned by creator?
	if token.Creator != msg.Sender {
		return nil, errors.NotCreatorBurn
	}

	// retrieve NFT sub-tokens
	subTokens := make([]types.SubToken, subTokenCount)
	for i, subTokenID := range msg.SubTokenIDs {
		subToken, subTokenExists := k.GetSubToken(ctx, token.ID, subTokenID)
		if !subTokenExists {
			return nil, errors.SubTokenDoesNotExists
		}
		// ensure NFT sub-token is owned by the sender
		if subToken.Owner != msg.Sender {
			return nil, errors.OwnerDoesNotOwnSubTokenID
		}
		if subToken.Reserve == nil {
			subToken.Reserve = &token.Reserve
		}
		subTokens[i] = subToken
	}

	// update NFT token counter
	token.Burnt += subTokenCount
	k.setTokenCounter(ctx, token.ID, types.TokenCounter{
		Minted: token.Minted,
		Burnt:  token.Burnt,
	})

	// burn NFT sub-tokens
	returnAmount := sdkmath.ZeroInt()
	for _, subToken := range subTokens {
		returnAmount = returnAmount.Add(subToken.Reserve.Amount)
		k.removeSubToken(ctx, token.ID, subToken.ID)
		k.removeSubTokenByOwner(ctx, sender, token.ID, subToken.ID)
	}
	returnCoin := sdk.NewCoin(token.Reserve.Denom, returnAmount)

	// return unlocked amount of coins
	err := k.ReturnTokensTo(ctx, sdk.NewCoins(returnCoin), sender)
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	// emit NFT sub-tokens burn event
	err = events.EmitTypedEvent(ctx, &types.EventBurnToken{
		Sender:      msg.Sender,
		ID:          msg.TokenID,
		Return:      returnCoin.String(),
		SubTokenIDs: msg.SubTokenIDs,
	})
	if err != nil {
		return nil, errors.Internal.Wrapf("err: %s", err.Error())
	}

	return &types.MsgBurnTokenResponse{}, nil
}

func (k Keeper) TransferSubTokens(ctx sdk.Context, sender, recipient sdk.AccAddress, tokenID string, subTokenIDs []uint32) error {
	// retrieve NFT token
	token, tokenExists := k.GetToken(ctx, tokenID)
	if !tokenExists {
		return errors.UnknownNFT
	}

	// retrieve NFT sub-tokens
	subTokens := make([]types.SubToken, len(subTokenIDs))
	for i, subTokenID := range subTokenIDs {
		subToken, subTokenExists := k.GetSubToken(ctx, token.ID, subTokenID)
		if !subTokenExists {
			return errors.SubTokenDoesNotExists
		}
		// ensure NFT sub-token is owned by the sender
		if subToken.Owner != sender.String() {
			return errors.OwnerDoesNotOwnSubTokenID
		}
		subTokens[i] = subToken
	}

	// transfer NFT sub-tokens
	for _, subToken := range subTokens {
		subToken.Owner = recipient.String()
		k.SetSubToken(ctx, token.ID, subToken)
		k.transferSubToken(ctx, sender, recipient, token.ID, subToken.ID)
	}

	// emit NFT sub-tokens transfer event
	err := events.EmitTypedEvent(ctx, &types.EventSendToken{
		Sender:      sender.String(),
		ID:          tokenID,
		Recipient:   recipient.String(),
		SubTokenIDs: subTokenIDs,
	})
	if err != nil {
		return errors.Internal.Wrapf("err: %s", err.Error())
	}

	return nil
}
