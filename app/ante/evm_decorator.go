package ante

import (
	"encoding/hex"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// EVMDecorator execute evm action to cosmos
type EVMDecorator struct {
	evmKeeper  ethante.EVMKeeper
	bankKeeper BankKeeper
}

// NewEVMDecorator creates a new EVMDecorator
func NewEVMDecorator(
	evmKeeper ethante.EVMKeeper,
	bankKeeper BankKeeper,
) EVMDecorator {
	return EVMDecorator{
		evmKeeper,
		bankKeeper,
	}
}

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"name\":\"name\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"name\":\"approve\",\"type\":\"function\",\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"constant\":false,\"stateMutability\":\"nonpayable\"},{\"name\":\"totalSupply\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"name\":\"transferFrom\",\"type\":\"function\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"constant\":false,\"stateMutability\":\"nonpayable\"},{\"name\":\"withdraw\",\"type\":\"function\",\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"outputs\":[],\"payable\":false,\"constant\":false,\"stateMutability\":\"nonpayable\"},{\"name\":\"decimals\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"name\":\"balanceOf\",\"type\":\"function\",\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"name\":\"symbol\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"name\":\"transfer\",\"type\":\"function\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"constant\":false,\"stateMutability\":\"nonpayable\"},{\"name\":\"deposit\",\"type\":\"function\",\"inputs\":[],\"outputs\":[],\"payable\":true,\"constant\":false,\"stateMutability\":\"payable\"},{\"name\":\"allowance\",\"type\":\"function\",\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"constant\":true,\"stateMutability\":\"view\"},{\"type\":\"fallback\",\"payable\":true,\"stateMutability\":\"payable\"},{\"name\":\"Approval\",\"type\":\"event\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true},{\"name\":\"guy\",\"type\":\"address\",\"indexed\":true},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false},{\"name\":\"Transfer\",\"type\":\"event\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true},{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false},{\"name\":\"Deposit\",\"type\":\"event\",\"inputs\":[{\"name\":\"dst\",\"type\":\"address\",\"indexed\":true},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false},{\"name\":\"Withdrawal\",\"type\":\"event\",\"inputs\":[{\"name\":\"src\",\"type\":\"address\",\"indexed\":true},{\"name\":\"wad\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false}]"

// AnteHandle validates that the Ethereum tx message has enough to cover intrinsic gas
// (during CheckTx only) and that the sender has enough balance to pay for the gas cost.
//
// Intrinsic gas for a transaction is the amount of gas that the transaction uses before the
// transaction is executed. The gas is a constant value plus any cost inccured by additional bytes
// of data supplied with the transaction.
//
// This AnteHandler decorator will fail if:
// - the message is not a MsgEthereumTx
// - sender account cannot be found
// - transaction's gas limit is lower than the intrinsic gas
// - user doesn't have enough balance to deduct the transaction fees (gas_limit * gas_price)
// - transaction or block gas meter runs out of gas
// - sets the gas meter limit
func (ed EVMDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, sdkerrors.ErrUnknownRequest.Wrapf("invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		ctx.Logger().Info(hex.EncodeToString(msgEthTx.AsTransaction().Data()[:]))

		ctx.Logger().Info(hex.EncodeToString(msgEthTx.AsTransaction().To()[:]))

		parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
		if err != nil {
			ctx.Logger().Info(err.Error())
		}

		method, err := parsed.MethodById(msgEthTx.AsTransaction().Data())
		if err != nil {
			ctx.Logger().Info(err.Error())
		}

		ctx.Logger().Info(method.Name)
		// ctx.Logger().Info(parsed.abi.Get)

		// bind.NewBoundContract("0x588ae821faf69761598d3b3b689672d9fbe91d36", parsed, caller, transactor, filterer)
	}
	// we know that we have enough gas on the pool to cover the intrinsic gas
	return next(ctx, tx, simulate)
}
