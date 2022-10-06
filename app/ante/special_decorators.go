package ante

import (
	coin "bitbucket.org/decimalteam/go-smart-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// Out limits: only 1 message in transaction
type CountMsgDecorator struct {
}

func NewCountMsgDecorator() CountMsgDecorator {
	return CountMsgDecorator{}
}

func (cd CountMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) > 1 {
		return ctx, CountOfMsgsMustBeOne
	}
	return next(ctx, tx, simulate)
}

/*
PreCreateAccountDecorator + PostCreateAccountDecorator are workaround for case of RedeemCheck and new account
New account has no account number i.e. his account number is 0, and account number is 0 in signature info
But after account creating account number is not 0, so we temporary set it to 0 for signature verification.
*/

// PreCreateAccountDecorator creates account in case of check redeeming from account unknown in the blockchain.
// Such accounts sign their first transaction with account number equal to 0. This is the reason why
// creating account is divided in two parts (PreCreateAccountDecorator and PostCreateAccountDecorator):
// it is necessary to create account in the beginning of the Ante chain with account number 0, but just
// before the end of the Ante chain it is necessary to restore correct account number.

type PreCreateAccountDecorator struct {
	ak evmtypes.AccountKeeper
}

// NewPreCreateAccountDecorator creates new PreCreateAccountDecorator.
func NewPreCreateAccountDecorator(ak evmtypes.AccountKeeper) PreCreateAccountDecorator {
	return PreCreateAccountDecorator{ak: ak}
}

// AnteHandle implements sdk.AnteHandler function.
func (cad PreCreateAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	if len(msgs) > 0 {
		switch msgs[0].(type) {
		case *coin.MsgRedeemCheck:
			signers := msgs[0].GetSigners()
			if len(signers) == 1 {
				acc := cad.ak.GetAccount(ctx, signers[0])
				if acc == nil {
					acc = cad.ak.NewAccountWithAddress(ctx, signers[0])
					ctx = ctx.WithValue("created_account_address", signers[0].String())
					ctx = ctx.WithValue("created_account_number", acc.GetAccountNumber())
					acc.SetAccountNumber(0) // necessary to validate signature
					cad.ak.SetAccount(ctx, acc)
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

// PostCreateAccountDecorator restores account number in case of check redeeming from account unknown for the blockchain.
type PostCreateAccountDecorator struct {
	ak evmtypes.AccountKeeper
}

// NewPostCreateAccountDecorator creates new PostCreateAccountDecorator.
func NewPostCreateAccountDecorator(ak evmtypes.AccountKeeper) PostCreateAccountDecorator {
	return PostCreateAccountDecorator{ak: ak}
}

// AnteHandle implements sdk.AnteHandler function.
func (cad PostCreateAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	accAddress := ctx.Value("created_account_address")
	accNumber := ctx.Value("created_account_number")
	if accAddress != nil && accNumber != nil {
		ctx = ctx.WithValue("created_account_address", nil)
		ctx = ctx.WithValue("created_account_number", nil)
		accAddr, err := sdk.AccAddressFromBech32(accAddress.(string))
		if err != nil {
			return ctx, InvalidAddressOfCreatedAccount
		}
		acc := cad.ak.GetAccount(ctx, accAddr)
		if acc == nil {
			return ctx, UnableToFindCreatedAccount
		}
		acc.SetAccountNumber(accNumber.(uint64))
		cad.ak.SetAccount(ctx, acc)
	}

	return next(ctx, tx, simulate)
}
