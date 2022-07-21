package keeper

import (
	"bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) MintNFT(c context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err == nil {
		if nft.GetCreator() != msg.Sender || !nft.GetAllowMint() {
			return nil, types.ErrNotAllowedMint()
		}
	} else {
		if k.ExistTokenURI(ctx, msg.TokenURI) {
			return nil, types.ErrNotUniqueTokenURI()
		}
		if k.ExistTokenID(ctx, msg.ID) {
			return nil, types.ErrNotUniqueTokenID()
		}
		if msg.Reserve.LT(types.NewMinReserve2) {
			return nil, types.ErrInvalidReserve(msg.Reserve.String())
		}
	}

	lastSubTokenID, err := k.Mint(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyNFTTokenURI, msg.TokenURI),
			sdk.NewAttribute(types.AttributeKeySubTokenIDStartRange, strconv.FormatUint(lastSubTokenID, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgMintNFTResponse{}, nil
}

func (k Keeper) TransferNFT(c context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, err := k.Transfer(ctx, msg.Denom, msg.ID, msg.Sender, msg.Recipient, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgTransferNFTResponse{}, nil
}

func (k Keeper) EditNFTMetadata(c context.Context, msg *types.MsgEditNFTMetadata) (*types.MsgEditNFTMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if nft.GetCreator() != msg.Sender {
		return nil, types.ErrNotAllowedMint()
	}

	// update NFT
	err = k.EditNFT(ctx, msg.Denom, msg.ID, msg.TokenURI)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditNFTMetadata,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyNFTTokenURI, msg.TokenURI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgEditNFTMetadataResponse{}, nil
}

func (k Keeper) BurnNFT(c context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if nft.GetCreator() != msg.Sender {
		return nil, types.ErrNotAllowedBurn()
	}

	// remove NFT
	err = k.DeleteNFTSubTokens(ctx, msg.Denom, msg.ID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnNFTResponse{}, nil
}

func (k Keeper) UpdateReserveNFT(c context.Context, msg *types.MsgUpdateReserveNFT) (*types.MsgUpdateReserveNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if nft.GetCreator() != msg.Sender {
		return nil, types.ErrNotAllowedUpdateReserve()
	}

	// update reserve nft
	err = k.UpdateNFTReserve(ctx, msg.Denom, msg.ID, msg.SubTokenIDs, msg.NewReserveNFT)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateReserveNFT,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgUpdateReserveNFTResponse{}, nil
}
