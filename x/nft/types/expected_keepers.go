package types

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	//GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account //TODO
}
