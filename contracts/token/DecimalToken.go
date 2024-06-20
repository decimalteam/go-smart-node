// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegation

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"ERC2612ExpiredSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC2612InvalidSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAmountInput\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientAmountOutput\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"currentNonce\",\"type\":\"uint256\"}],\"name\":\"InvalidAccountNonce\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCrr\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMaxTotalSupply\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinTotalSupply\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSymbol\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewMaxTotalSupplyTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewReserveTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewSupplyTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewSupplyTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCreator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Log_InputTooSmall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newReserve\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSupply\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newPrice\",\"type\":\"uint256\"}],\"name\":\"ReserveUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CREATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"buy\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"buyExactTokenForDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"buyTokenForExactDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateBuyInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateBuyInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateBuyOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateBuyOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateSellInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateSellInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateSellOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateSellOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"creator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crr\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"identity\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initialName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialSymbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialCreator\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"initialCrr\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialMinTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialMaxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"initialIdentity\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"permitHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sell\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sellExactTokensForDEL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"sellTokensForExactDEL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newIdentity\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"newMaxTotalSupply\",\"type\":\"uint256\"}],\"name\":\"updateDetails\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// DelegationABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegationMetaData.ABI instead.
var DelegationABI = DelegationMetaData.ABI

// Delegation is an auto generated Go binding around an Ethereum contract.
type Delegation struct {
	DelegationCaller     // Read-only binding to the contract
	DelegationTransactor // Write-only binding to the contract
	DelegationFilterer   // Log filterer for contract events
}

// DelegationCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegationSession struct {
	Contract     *Delegation       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegationCallerSession struct {
	Contract *DelegationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// DelegationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegationTransactorSession struct {
	Contract     *DelegationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DelegationRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegationRaw struct {
	Contract *Delegation // Generic contract binding to access the raw methods on
}

// DelegationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegationCallerRaw struct {
	Contract *DelegationCaller // Generic read-only contract binding to access the raw methods on
}

// DelegationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegationTransactorRaw struct {
	Contract *DelegationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegation creates a new instance of Delegation, bound to a specific deployed contract.
func NewDelegation(address common.Address, backend bind.ContractBackend) (*Delegation, error) {
	contract, err := bindDelegation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegation{DelegationCaller: DelegationCaller{contract: contract}, DelegationTransactor: DelegationTransactor{contract: contract}, DelegationFilterer: DelegationFilterer{contract: contract}}, nil
}

// NewDelegationCaller creates a new read-only instance of Delegation, bound to a specific deployed contract.
func NewDelegationCaller(address common.Address, caller bind.ContractCaller) (*DelegationCaller, error) {
	contract, err := bindDelegation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationCaller{contract: contract}, nil
}

// NewDelegationTransactor creates a new write-only instance of Delegation, bound to a specific deployed contract.
func NewDelegationTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegationTransactor, error) {
	contract, err := bindDelegation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationTransactor{contract: contract}, nil
}

// NewDelegationFilterer creates a new log filterer instance of Delegation, bound to a specific deployed contract.
func NewDelegationFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegationFilterer, error) {
	contract, err := bindDelegation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegationFilterer{contract: contract}, nil
}

// bindDelegation binds a generic wrapper to an already deployed contract.
func bindDelegation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegation *DelegationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegation.Contract.DelegationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegation *DelegationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.Contract.DelegationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegation *DelegationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegation.Contract.DelegationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegation *DelegationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegation *DelegationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegation *DelegationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegation.Contract.contract.Transact(opts, method, params...)
}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_Delegation *DelegationCaller) CREATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "CREATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_Delegation *DelegationSession) CREATORROLE() ([32]byte, error) {
	return _Delegation.Contract.CREATORROLE(&_Delegation.CallOpts)
}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_Delegation *DelegationCallerSession) CREATORROLE() ([32]byte, error) {
	return _Delegation.Contract.CREATORROLE(&_Delegation.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Delegation *DelegationCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Delegation *DelegationSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Delegation.Contract.DOMAINSEPARATOR(&_Delegation.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Delegation *DelegationCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Delegation.Contract.DOMAINSEPARATOR(&_Delegation.CallOpts)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Delegation *DelegationCaller) MAXCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MAX_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Delegation *DelegationSession) MAXCRR() (*big.Int, error) {
	return _Delegation.Contract.MAXCRR(&_Delegation.CallOpts)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Delegation *DelegationCallerSession) MAXCRR() (*big.Int, error) {
	return _Delegation.Contract.MAXCRR(&_Delegation.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationCaller) MAXTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MAX_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Delegation.Contract.MAXTOTALSUPPLY(&_Delegation.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationCallerSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Delegation.Contract.MAXTOTALSUPPLY(&_Delegation.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Delegation *DelegationCaller) MINCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MIN_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Delegation *DelegationSession) MINCRR() (*big.Int, error) {
	return _Delegation.Contract.MINCRR(&_Delegation.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Delegation *DelegationCallerSession) MINCRR() (*big.Int, error) {
	return _Delegation.Contract.MINCRR(&_Delegation.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Delegation *DelegationCaller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Delegation *DelegationSession) MINRESERVE() (*big.Int, error) {
	return _Delegation.Contract.MINRESERVE(&_Delegation.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Delegation *DelegationCallerSession) MINRESERVE() (*big.Int, error) {
	return _Delegation.Contract.MINRESERVE(&_Delegation.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationCaller) MINTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MIN_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Delegation.Contract.MINTOTALSUPPLY(&_Delegation.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Delegation *DelegationCallerSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Delegation.Contract.MINTOTALSUPPLY(&_Delegation.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Delegation *DelegationCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "PERMIT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Delegation *DelegationSession) PERMITTYPEHASH() ([32]byte, error) {
	return _Delegation.Contract.PERMITTYPEHASH(&_Delegation.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_Delegation *DelegationCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _Delegation.Contract.PERMITTYPEHASH(&_Delegation.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Delegation *DelegationCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Delegation *DelegationSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Delegation.Contract.Allowance(&_Delegation.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Delegation *DelegationCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Delegation.Contract.Allowance(&_Delegation.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Delegation *DelegationCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Delegation *DelegationSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Delegation.Contract.BalanceOf(&_Delegation.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Delegation *DelegationCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Delegation.Contract.BalanceOf(&_Delegation.CallOpts, account)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationCaller) CalculateBuyInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateBuyInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyInput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyInput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyInput0 is a free data retrieval call binding the contract method 0xf804cd49.
//
// Solidity: function calculateBuyInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCaller) CalculateBuyInput0(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateBuyInput0", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyInput0 is a free data retrieval call binding the contract method 0xf804cd49.
//
// Solidity: function calculateBuyInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationSession) CalculateBuyInput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyInput0(&_Delegation.CallOpts, amount)
}

// CalculateBuyInput0 is a free data retrieval call binding the contract method 0xf804cd49.
//
// Solidity: function calculateBuyInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateBuyInput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyInput0(&_Delegation.CallOpts, amount)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationCaller) CalculateBuyOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateBuyOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyOutput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyOutput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateBuyOutput0 is a free data retrieval call binding the contract method 0x6cc94a2b.
//
// Solidity: function calculateBuyOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCaller) CalculateBuyOutput0(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateBuyOutput0", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyOutput0 is a free data retrieval call binding the contract method 0x6cc94a2b.
//
// Solidity: function calculateBuyOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationSession) CalculateBuyOutput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyOutput0(&_Delegation.CallOpts, amount)
}

// CalculateBuyOutput0 is a free data retrieval call binding the contract method 0x6cc94a2b.
//
// Solidity: function calculateBuyOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateBuyOutput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateBuyOutput0(&_Delegation.CallOpts, amount)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationCaller) CalculateSellInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateSellInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellInput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellInput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellInput0 is a free data retrieval call binding the contract method 0x7058c1aa.
//
// Solidity: function calculateSellInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCaller) CalculateSellInput0(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateSellInput0", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellInput0 is a free data retrieval call binding the contract method 0x7058c1aa.
//
// Solidity: function calculateSellInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationSession) CalculateSellInput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellInput0(&_Delegation.CallOpts, amount)
}

// CalculateSellInput0 is a free data retrieval call binding the contract method 0x7058c1aa.
//
// Solidity: function calculateSellInput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateSellInput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellInput0(&_Delegation.CallOpts, amount)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationCaller) CalculateSellOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateSellOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellOutput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellOutput(&_Delegation.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellOutput0 is a free data retrieval call binding the contract method 0xc193d62e.
//
// Solidity: function calculateSellOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCaller) CalculateSellOutput0(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "calculateSellOutput0", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellOutput0 is a free data retrieval call binding the contract method 0xc193d62e.
//
// Solidity: function calculateSellOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationSession) CalculateSellOutput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellOutput0(&_Delegation.CallOpts, amount)
}

// CalculateSellOutput0 is a free data retrieval call binding the contract method 0xc193d62e.
//
// Solidity: function calculateSellOutput(uint256 amount) view returns(uint256)
func (_Delegation *DelegationCallerSession) CalculateSellOutput0(amount *big.Int) (*big.Int, error) {
	return _Delegation.Contract.CalculateSellOutput0(&_Delegation.CallOpts, amount)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Delegation *DelegationCaller) Creator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "creator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Delegation *DelegationSession) Creator() (common.Address, error) {
	return _Delegation.Contract.Creator(&_Delegation.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Delegation *DelegationCallerSession) Creator() (common.Address, error) {
	return _Delegation.Contract.Creator(&_Delegation.CallOpts)
}

// Crr is a free data retrieval call binding the contract method 0x68213256.
//
// Solidity: function crr() view returns(uint256)
func (_Delegation *DelegationCaller) Crr(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "crr")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Crr is a free data retrieval call binding the contract method 0x68213256.
//
// Solidity: function crr() view returns(uint256)
func (_Delegation *DelegationSession) Crr() (*big.Int, error) {
	return _Delegation.Contract.Crr(&_Delegation.CallOpts)
}

// Crr is a free data retrieval call binding the contract method 0x68213256.
//
// Solidity: function crr() view returns(uint256)
func (_Delegation *DelegationCallerSession) Crr() (*big.Int, error) {
	return _Delegation.Contract.Crr(&_Delegation.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Delegation *DelegationCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Delegation *DelegationSession) Decimals() (uint8, error) {
	return _Delegation.Contract.Decimals(&_Delegation.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Delegation *DelegationCallerSession) Decimals() (uint8, error) {
	return _Delegation.Contract.Decimals(&_Delegation.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_Delegation *DelegationCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_Delegation *DelegationSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _Delegation.Contract.Eip712Domain(&_Delegation.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_Delegation *DelegationCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _Delegation.Contract.Eip712Domain(&_Delegation.CallOpts)
}

// GetCurrentPrice is a free data retrieval call binding the contract method 0xeb91d37e.
//
// Solidity: function getCurrentPrice() view returns(uint256)
func (_Delegation *DelegationCaller) GetCurrentPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getCurrentPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentPrice is a free data retrieval call binding the contract method 0xeb91d37e.
//
// Solidity: function getCurrentPrice() view returns(uint256)
func (_Delegation *DelegationSession) GetCurrentPrice() (*big.Int, error) {
	return _Delegation.Contract.GetCurrentPrice(&_Delegation.CallOpts)
}

// GetCurrentPrice is a free data retrieval call binding the contract method 0xeb91d37e.
//
// Solidity: function getCurrentPrice() view returns(uint256)
func (_Delegation *DelegationCallerSession) GetCurrentPrice() (*big.Int, error) {
	return _Delegation.Contract.GetCurrentPrice(&_Delegation.CallOpts)
}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Delegation *DelegationCaller) Identity(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "identity")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Delegation *DelegationSession) Identity() (string, error) {
	return _Delegation.Contract.Identity(&_Delegation.CallOpts)
}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Delegation *DelegationCallerSession) Identity() (string, error) {
	return _Delegation.Contract.Identity(&_Delegation.CallOpts)
}

// MaxTotalSupply is a free data retrieval call binding the contract method 0x2ab4d052.
//
// Solidity: function maxTotalSupply() view returns(uint256)
func (_Delegation *DelegationCaller) MaxTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "maxTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxTotalSupply is a free data retrieval call binding the contract method 0x2ab4d052.
//
// Solidity: function maxTotalSupply() view returns(uint256)
func (_Delegation *DelegationSession) MaxTotalSupply() (*big.Int, error) {
	return _Delegation.Contract.MaxTotalSupply(&_Delegation.CallOpts)
}

// MaxTotalSupply is a free data retrieval call binding the contract method 0x2ab4d052.
//
// Solidity: function maxTotalSupply() view returns(uint256)
func (_Delegation *DelegationCallerSession) MaxTotalSupply() (*big.Int, error) {
	return _Delegation.Contract.MaxTotalSupply(&_Delegation.CallOpts)
}

// MinTotalSupply is a free data retrieval call binding the contract method 0x79db6346.
//
// Solidity: function minTotalSupply() view returns(uint256)
func (_Delegation *DelegationCaller) MinTotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "minTotalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinTotalSupply is a free data retrieval call binding the contract method 0x79db6346.
//
// Solidity: function minTotalSupply() view returns(uint256)
func (_Delegation *DelegationSession) MinTotalSupply() (*big.Int, error) {
	return _Delegation.Contract.MinTotalSupply(&_Delegation.CallOpts)
}

// MinTotalSupply is a free data retrieval call binding the contract method 0x79db6346.
//
// Solidity: function minTotalSupply() view returns(uint256)
func (_Delegation *DelegationCallerSession) MinTotalSupply() (*big.Int, error) {
	return _Delegation.Contract.MinTotalSupply(&_Delegation.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Delegation *DelegationCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Delegation *DelegationSession) Name() (string, error) {
	return _Delegation.Contract.Name(&_Delegation.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Delegation *DelegationCallerSession) Name() (string, error) {
	return _Delegation.Contract.Name(&_Delegation.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Delegation *DelegationCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Delegation *DelegationSession) Nonces(owner common.Address) (*big.Int, error) {
	return _Delegation.Contract.Nonces(&_Delegation.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Delegation *DelegationCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _Delegation.Contract.Nonces(&_Delegation.CallOpts, owner)
}

// PermitHash is a free data retrieval call binding the contract method 0x02bfc0f9.
//
// Solidity: function permitHash(address owner, address spender, uint256 value, uint256 nonce, uint256 deadline) view returns(bytes32)
func (_Delegation *DelegationCaller) PermitHash(opts *bind.CallOpts, owner common.Address, spender common.Address, value *big.Int, nonce *big.Int, deadline *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "permitHash", owner, spender, value, nonce, deadline)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PermitHash is a free data retrieval call binding the contract method 0x02bfc0f9.
//
// Solidity: function permitHash(address owner, address spender, uint256 value, uint256 nonce, uint256 deadline) view returns(bytes32)
func (_Delegation *DelegationSession) PermitHash(owner common.Address, spender common.Address, value *big.Int, nonce *big.Int, deadline *big.Int) ([32]byte, error) {
	return _Delegation.Contract.PermitHash(&_Delegation.CallOpts, owner, spender, value, nonce, deadline)
}

// PermitHash is a free data retrieval call binding the contract method 0x02bfc0f9.
//
// Solidity: function permitHash(address owner, address spender, uint256 value, uint256 nonce, uint256 deadline) view returns(bytes32)
func (_Delegation *DelegationCallerSession) PermitHash(owner common.Address, spender common.Address, value *big.Int, nonce *big.Int, deadline *big.Int) ([32]byte, error) {
	return _Delegation.Contract.PermitHash(&_Delegation.CallOpts, owner, spender, value, nonce, deadline)
}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() view returns(uint256)
func (_Delegation *DelegationCaller) Reserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "reserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() view returns(uint256)
func (_Delegation *DelegationSession) Reserve() (*big.Int, error) {
	return _Delegation.Contract.Reserve(&_Delegation.CallOpts)
}

// Reserve is a free data retrieval call binding the contract method 0xcd3293de.
//
// Solidity: function reserve() view returns(uint256)
func (_Delegation *DelegationCallerSession) Reserve() (*big.Int, error) {
	return _Delegation.Contract.Reserve(&_Delegation.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Delegation *DelegationCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Delegation *DelegationSession) Symbol() (string, error) {
	return _Delegation.Contract.Symbol(&_Delegation.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Delegation *DelegationCallerSession) Symbol() (string, error) {
	return _Delegation.Contract.Symbol(&_Delegation.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Delegation *DelegationCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Delegation *DelegationSession) TotalSupply() (*big.Int, error) {
	return _Delegation.Contract.TotalSupply(&_Delegation.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Delegation *DelegationCallerSession) TotalSupply() (*big.Int, error) {
	return _Delegation.Contract.TotalSupply(&_Delegation.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Delegation *DelegationTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Delegation *DelegationSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Approve(&_Delegation.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Delegation *DelegationTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Approve(&_Delegation.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Delegation *DelegationTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Delegation *DelegationSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Burn(&_Delegation.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_Delegation *DelegationTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Burn(&_Delegation.TransactOpts, value)
}

// Buy is a paid mutator transaction binding the contract method 0x7deb6025.
//
// Solidity: function buy(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationTransactor) Buy(opts *bind.TransactOpts, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "buy", amountOutMin, recipient)
}

// Buy is a paid mutator transaction binding the contract method 0x7deb6025.
//
// Solidity: function buy(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationSession) Buy(amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Buy(&_Delegation.TransactOpts, amountOutMin, recipient)
}

// Buy is a paid mutator transaction binding the contract method 0x7deb6025.
//
// Solidity: function buy(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationTransactorSession) Buy(amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Buy(&_Delegation.TransactOpts, amountOutMin, recipient)
}

// BuyExactTokenForDEL is a paid mutator transaction binding the contract method 0xf6228a48.
//
// Solidity: function buyExactTokenForDEL(uint256 amountOut, address recipient) payable returns()
func (_Delegation *DelegationTransactor) BuyExactTokenForDEL(opts *bind.TransactOpts, amountOut *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "buyExactTokenForDEL", amountOut, recipient)
}

// BuyExactTokenForDEL is a paid mutator transaction binding the contract method 0xf6228a48.
//
// Solidity: function buyExactTokenForDEL(uint256 amountOut, address recipient) payable returns()
func (_Delegation *DelegationSession) BuyExactTokenForDEL(amountOut *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.BuyExactTokenForDEL(&_Delegation.TransactOpts, amountOut, recipient)
}

// BuyExactTokenForDEL is a paid mutator transaction binding the contract method 0xf6228a48.
//
// Solidity: function buyExactTokenForDEL(uint256 amountOut, address recipient) payable returns()
func (_Delegation *DelegationTransactorSession) BuyExactTokenForDEL(amountOut *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.BuyExactTokenForDEL(&_Delegation.TransactOpts, amountOut, recipient)
}

// BuyTokenForExactDEL is a paid mutator transaction binding the contract method 0x98c9c2a9.
//
// Solidity: function buyTokenForExactDEL(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationTransactor) BuyTokenForExactDEL(opts *bind.TransactOpts, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "buyTokenForExactDEL", amountOutMin, recipient)
}

// BuyTokenForExactDEL is a paid mutator transaction binding the contract method 0x98c9c2a9.
//
// Solidity: function buyTokenForExactDEL(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationSession) BuyTokenForExactDEL(amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.BuyTokenForExactDEL(&_Delegation.TransactOpts, amountOutMin, recipient)
}

// BuyTokenForExactDEL is a paid mutator transaction binding the contract method 0x98c9c2a9.
//
// Solidity: function buyTokenForExactDEL(uint256 amountOutMin, address recipient) payable returns()
func (_Delegation *DelegationTransactorSession) BuyTokenForExactDEL(amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.BuyTokenForExactDEL(&_Delegation.TransactOpts, amountOutMin, recipient)
}

// Initialize is a paid mutator transaction binding the contract method 0x9e4c800c.
//
// Solidity: function initialize(string initialName, string initialSymbol, address initialCreator, uint8 initialCrr, uint256 initialMint, uint256 initialMinTotalSupply, uint256 initialMaxTotalSupply, string initialIdentity) payable returns()
func (_Delegation *DelegationTransactor) Initialize(opts *bind.TransactOpts, initialName string, initialSymbol string, initialCreator common.Address, initialCrr uint8, initialMint *big.Int, initialMinTotalSupply *big.Int, initialMaxTotalSupply *big.Int, initialIdentity string) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "initialize", initialName, initialSymbol, initialCreator, initialCrr, initialMint, initialMinTotalSupply, initialMaxTotalSupply, initialIdentity)
}

// Initialize is a paid mutator transaction binding the contract method 0x9e4c800c.
//
// Solidity: function initialize(string initialName, string initialSymbol, address initialCreator, uint8 initialCrr, uint256 initialMint, uint256 initialMinTotalSupply, uint256 initialMaxTotalSupply, string initialIdentity) payable returns()
func (_Delegation *DelegationSession) Initialize(initialName string, initialSymbol string, initialCreator common.Address, initialCrr uint8, initialMint *big.Int, initialMinTotalSupply *big.Int, initialMaxTotalSupply *big.Int, initialIdentity string) (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts, initialName, initialSymbol, initialCreator, initialCrr, initialMint, initialMinTotalSupply, initialMaxTotalSupply, initialIdentity)
}

// Initialize is a paid mutator transaction binding the contract method 0x9e4c800c.
//
// Solidity: function initialize(string initialName, string initialSymbol, address initialCreator, uint8 initialCrr, uint256 initialMint, uint256 initialMinTotalSupply, uint256 initialMaxTotalSupply, string initialIdentity) payable returns()
func (_Delegation *DelegationTransactorSession) Initialize(initialName string, initialSymbol string, initialCreator common.Address, initialCrr uint8, initialMint *big.Int, initialMinTotalSupply *big.Int, initialMaxTotalSupply *big.Int, initialIdentity string) (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts, initialName, initialSymbol, initialCreator, initialCrr, initialMint, initialMinTotalSupply, initialMaxTotalSupply, initialIdentity)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Permit(&_Delegation.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Permit(&_Delegation.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Sell is a paid mutator transaction binding the contract method 0xd04c6983.
//
// Solidity: function sell(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactor) Sell(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "sell", amountIn, amountOutMin, recipient)
}

// Sell is a paid mutator transaction binding the contract method 0xd04c6983.
//
// Solidity: function sell(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationSession) Sell(amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Sell(&_Delegation.TransactOpts, amountIn, amountOutMin, recipient)
}

// Sell is a paid mutator transaction binding the contract method 0xd04c6983.
//
// Solidity: function sell(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactorSession) Sell(amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Sell(&_Delegation.TransactOpts, amountIn, amountOutMin, recipient)
}

// SellExactTokensForDEL is a paid mutator transaction binding the contract method 0x7627a7af.
//
// Solidity: function sellExactTokensForDEL(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactor) SellExactTokensForDEL(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "sellExactTokensForDEL", amountIn, amountOutMin, recipient)
}

// SellExactTokensForDEL is a paid mutator transaction binding the contract method 0x7627a7af.
//
// Solidity: function sellExactTokensForDEL(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationSession) SellExactTokensForDEL(amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SellExactTokensForDEL(&_Delegation.TransactOpts, amountIn, amountOutMin, recipient)
}

// SellExactTokensForDEL is a paid mutator transaction binding the contract method 0x7627a7af.
//
// Solidity: function sellExactTokensForDEL(uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactorSession) SellExactTokensForDEL(amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SellExactTokensForDEL(&_Delegation.TransactOpts, amountIn, amountOutMin, recipient)
}

// SellTokensForExactDEL is a paid mutator transaction binding the contract method 0xca017dcd.
//
// Solidity: function sellTokensForExactDEL(uint256 amountOut, uint256 amountInMax, address recipient) returns()
func (_Delegation *DelegationTransactor) SellTokensForExactDEL(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "sellTokensForExactDEL", amountOut, amountInMax, recipient)
}

// SellTokensForExactDEL is a paid mutator transaction binding the contract method 0xca017dcd.
//
// Solidity: function sellTokensForExactDEL(uint256 amountOut, uint256 amountInMax, address recipient) returns()
func (_Delegation *DelegationSession) SellTokensForExactDEL(amountOut *big.Int, amountInMax *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SellTokensForExactDEL(&_Delegation.TransactOpts, amountOut, amountInMax, recipient)
}

// SellTokensForExactDEL is a paid mutator transaction binding the contract method 0xca017dcd.
//
// Solidity: function sellTokensForExactDEL(uint256 amountOut, uint256 amountInMax, address recipient) returns()
func (_Delegation *DelegationTransactorSession) SellTokensForExactDEL(amountOut *big.Int, amountInMax *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SellTokensForExactDEL(&_Delegation.TransactOpts, amountOut, amountInMax, recipient)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Delegation *DelegationTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Delegation *DelegationSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Transfer(&_Delegation.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_Delegation *DelegationTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Transfer(&_Delegation.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Delegation *DelegationTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Delegation *DelegationSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.TransferFrom(&_Delegation.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_Delegation *DelegationTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.TransferFrom(&_Delegation.TransactOpts, from, to, value)
}

// UpdateDetails is a paid mutator transaction binding the contract method 0xd4a2949e.
//
// Solidity: function updateDetails(string newIdentity, uint256 newMaxTotalSupply) returns()
func (_Delegation *DelegationTransactor) UpdateDetails(opts *bind.TransactOpts, newIdentity string, newMaxTotalSupply *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "updateDetails", newIdentity, newMaxTotalSupply)
}

// UpdateDetails is a paid mutator transaction binding the contract method 0xd4a2949e.
//
// Solidity: function updateDetails(string newIdentity, uint256 newMaxTotalSupply) returns()
func (_Delegation *DelegationSession) UpdateDetails(newIdentity string, newMaxTotalSupply *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.UpdateDetails(&_Delegation.TransactOpts, newIdentity, newMaxTotalSupply)
}

// UpdateDetails is a paid mutator transaction binding the contract method 0xd4a2949e.
//
// Solidity: function updateDetails(string newIdentity, uint256 newMaxTotalSupply) returns()
func (_Delegation *DelegationTransactorSession) UpdateDetails(newIdentity string, newMaxTotalSupply *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.UpdateDetails(&_Delegation.TransactOpts, newIdentity, newMaxTotalSupply)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegation *DelegationTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegation *DelegationSession) Receive() (*types.Transaction, error) {
	return _Delegation.Contract.Receive(&_Delegation.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegation *DelegationTransactorSession) Receive() (*types.Transaction, error) {
	return _Delegation.Contract.Receive(&_Delegation.TransactOpts)
}

// DelegationApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Delegation contract.
type DelegationApprovalIterator struct {
	Event *DelegationApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationApproval represents a Approval event raised by the Delegation contract.
type DelegationApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Delegation *DelegationFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*DelegationApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &DelegationApprovalIterator{contract: _Delegation.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Delegation *DelegationFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *DelegationApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationApproval)
				if err := _Delegation.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Delegation *DelegationFilterer) ParseApproval(log types.Log) (*DelegationApproval, error) {
	event := new(DelegationApproval)
	if err := _Delegation.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the Delegation contract.
type DelegationEIP712DomainChangedIterator struct {
	Event *DelegationEIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationEIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationEIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationEIP712DomainChanged represents a EIP712DomainChanged event raised by the Delegation contract.
type DelegationEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_Delegation *DelegationFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*DelegationEIP712DomainChangedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &DelegationEIP712DomainChangedIterator{contract: _Delegation.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_Delegation *DelegationFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *DelegationEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationEIP712DomainChanged)
				if err := _Delegation.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_Delegation *DelegationFilterer) ParseEIP712DomainChanged(log types.Log) (*DelegationEIP712DomainChanged, error) {
	event := new(DelegationEIP712DomainChanged)
	if err := _Delegation.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Delegation contract.
type DelegationInitializedIterator struct {
	Event *DelegationInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationInitialized represents a Initialized event raised by the Delegation contract.
type DelegationInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Delegation *DelegationFilterer) FilterInitialized(opts *bind.FilterOpts) (*DelegationInitializedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DelegationInitializedIterator{contract: _Delegation.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Delegation *DelegationFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DelegationInitialized) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationInitialized)
				if err := _Delegation.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Delegation *DelegationFilterer) ParseInitialized(log types.Log) (*DelegationInitialized, error) {
	event := new(DelegationInitialized)
	if err := _Delegation.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationReserveUpdatedIterator is returned from FilterReserveUpdated and is used to iterate over the raw logs and unpacked data for ReserveUpdated events raised by the Delegation contract.
type DelegationReserveUpdatedIterator struct {
	Event *DelegationReserveUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationReserveUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationReserveUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationReserveUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationReserveUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationReserveUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationReserveUpdated represents a ReserveUpdated event raised by the Delegation contract.
type DelegationReserveUpdated struct {
	NewReserve *big.Int
	NewSupply  *big.Int
	NewPrice   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReserveUpdated is a free log retrieval operation binding the contract event 0x736a4a5812ced57865d349f18ffc358079c6b479326c0dfd1dae30c465b1daf2.
//
// Solidity: event ReserveUpdated(uint256 newReserve, uint256 newSupply, uint256 newPrice)
func (_Delegation *DelegationFilterer) FilterReserveUpdated(opts *bind.FilterOpts) (*DelegationReserveUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationReserveUpdatedIterator{contract: _Delegation.contract, event: "ReserveUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveUpdated is a free log subscription operation binding the contract event 0x736a4a5812ced57865d349f18ffc358079c6b479326c0dfd1dae30c465b1daf2.
//
// Solidity: event ReserveUpdated(uint256 newReserve, uint256 newSupply, uint256 newPrice)
func (_Delegation *DelegationFilterer) WatchReserveUpdated(opts *bind.WatchOpts, sink chan<- *DelegationReserveUpdated) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationReserveUpdated)
				if err := _Delegation.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReserveUpdated is a log parse operation binding the contract event 0x736a4a5812ced57865d349f18ffc358079c6b479326c0dfd1dae30c465b1daf2.
//
// Solidity: event ReserveUpdated(uint256 newReserve, uint256 newSupply, uint256 newPrice)
func (_Delegation *DelegationFilterer) ParseReserveUpdated(log types.Log) (*DelegationReserveUpdated, error) {
	event := new(DelegationReserveUpdated)
	if err := _Delegation.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Delegation contract.
type DelegationTransferIterator struct {
	Event *DelegationTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DelegationTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DelegationTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DelegationTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTransfer represents a Transfer event raised by the Delegation contract.
type DelegationTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Delegation *DelegationFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*DelegationTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationTransferIterator{contract: _Delegation.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Delegation *DelegationFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *DelegationTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTransfer)
				if err := _Delegation.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Delegation *DelegationFilterer) ParseTransfer(log types.Log) (*DelegationTransfer, error) {
	event := new(DelegationTransfer)
	if err := _Delegation.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
