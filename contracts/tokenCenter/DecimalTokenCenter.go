// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokenCenter

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

// DecimalTokenCenterToken is an auto generated low-level Go binding around an user-defined struct.
type DecimalTokenCenterToken struct {
	InitialMint    *big.Int
	MinTotalSupply *big.Int
	MaxTotalSupply *big.Int
	TokenOwner     common.Address
	Crr            uint8
	Identity       string
	Symbol         string
	Name           string
}

// TokenMetaData contains all meta data concerning the Token contract.
var TokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedCommission\",\"type\":\"uint256\"}],\"name\":\"InvalidComission\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Log_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenSymbolExist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"TokenContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"crr\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structDecimalTokenCenter.Token\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"TokenDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_INITIAL_MINT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateBuyInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateBuyOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateSellInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateSellOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"crr\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structDecimalTokenCenter.Token\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getCommissionSymbol\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isTokenExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressContractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"tokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTokenImplementation\",\"type\":\"address\"}],\"name\":\"upgradeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// TokenABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenMetaData.ABI instead.
var TokenABI = TokenMetaData.ABI

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
	TokenFilterer   // Log filterer for contract events
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewToken creates a new instance of Token, bound to a specific deployed contract.
func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenCaller: TokenCaller{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}, TokenFilterer: TokenFilterer{contract: contract}}, nil
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// NewTokenTransactor creates a new write-only instance of Token, bound to a specific deployed contract.
func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

// NewTokenFilterer creates a new log filterer instance of Token, bound to a specific deployed contract.
func NewTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenFilterer, error) {
	contract, err := bindToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenFilterer{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Token.Contract.TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Token *TokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Token *TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Token *TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Token.Contract.contract.Transact(opts, method, params...)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Token *TokenCaller) MAXCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MAX_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Token *TokenSession) MAXCRR() (*big.Int, error) {
	return _Token.Contract.MAXCRR(&_Token.CallOpts)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Token *TokenCallerSession) MAXCRR() (*big.Int, error) {
	return _Token.Contract.MAXCRR(&_Token.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenCaller) MAXTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MAX_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Token.Contract.MAXTOTALSUPPLY(&_Token.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenCallerSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Token.Contract.MAXTOTALSUPPLY(&_Token.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Token *TokenCaller) MINCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MIN_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Token *TokenSession) MINCRR() (*big.Int, error) {
	return _Token.Contract.MINCRR(&_Token.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Token *TokenCallerSession) MINCRR() (*big.Int, error) {
	return _Token.Contract.MINCRR(&_Token.CallOpts)
}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Token *TokenCaller) MININITIALMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MIN_INITIAL_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Token *TokenSession) MININITIALMINT() (*big.Int, error) {
	return _Token.Contract.MININITIALMINT(&_Token.CallOpts)
}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Token *TokenCallerSession) MININITIALMINT() (*big.Int, error) {
	return _Token.Contract.MININITIALMINT(&_Token.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Token *TokenCaller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Token *TokenSession) MINRESERVE() (*big.Int, error) {
	return _Token.Contract.MINRESERVE(&_Token.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Token *TokenCallerSession) MINRESERVE() (*big.Int, error) {
	return _Token.Contract.MINRESERVE(&_Token.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenCaller) MINTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "MIN_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Token.Contract.MINTOTALSUPPLY(&_Token.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Token *TokenCallerSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Token.Contract.MINTOTALSUPPLY(&_Token.CallOpts)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenCaller) CalculateBuyInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "calculateBuyInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateBuyInput(&_Token.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenCallerSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateBuyInput(&_Token.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenCaller) CalculateBuyOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "calculateBuyOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateBuyOutput(&_Token.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenCallerSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateBuyOutput(&_Token.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenCaller) CalculateSellInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "calculateSellInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateSellInput(&_Token.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Token *TokenCallerSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateSellInput(&_Token.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenCaller) CalculateSellOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "calculateSellOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateSellOutput(&_Token.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Token *TokenCallerSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Token.Contract.CalculateSellOutput(&_Token.CallOpts, supply, customReserve, customCrr, amountIn)
}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Token *TokenCaller) GetCommissionSymbol(opts *bind.CallOpts, symbol string) (*big.Int, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "getCommissionSymbol", symbol)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Token *TokenSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Token.Contract.GetCommissionSymbol(&_Token.CallOpts, symbol)
}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Token *TokenCallerSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Token.Contract.GetCommissionSymbol(&_Token.CallOpts, symbol)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Token *TokenCaller) GetContractCenter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "getContractCenter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Token *TokenSession) GetContractCenter() (common.Address, error) {
	return _Token.Contract.GetContractCenter(&_Token.CallOpts)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Token *TokenCallerSession) GetContractCenter() (common.Address, error) {
	return _Token.Contract.GetContractCenter(&_Token.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Token *TokenCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Token *TokenSession) GetImplementation() (common.Address, error) {
	return _Token.Contract.GetImplementation(&_Token.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Token *TokenCallerSession) GetImplementation() (common.Address, error) {
	return _Token.Contract.GetImplementation(&_Token.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Token *TokenCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Token *TokenSession) Implementation() (common.Address, error) {
	return _Token.Contract.Implementation(&_Token.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Token *TokenCallerSession) Implementation() (common.Address, error) {
	return _Token.Contract.Implementation(&_Token.CallOpts)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Token *TokenCaller) IsTokenExists(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "isTokenExists", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Token *TokenSession) IsTokenExists(token common.Address) (bool, error) {
	return _Token.Contract.IsTokenExists(&_Token.CallOpts, token)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Token *TokenCallerSession) IsTokenExists(token common.Address) (bool, error) {
	return _Token.Contract.IsTokenExists(&_Token.CallOpts, token)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Token *TokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Token *TokenSession) Owner() (common.Address, error) {
	return _Token.Contract.Owner(&_Token.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Token *TokenCallerSession) Owner() (common.Address, error) {
	return _Token.Contract.Owner(&_Token.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Token *TokenCaller) Tokens(opts *bind.CallOpts, symbol string) (common.Address, error) {
	var out []interface{}
	err := _Token.contract.Call(opts, &out, "tokens", symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Token *TokenSession) Tokens(symbol string) (common.Address, error) {
	return _Token.Contract.Tokens(&_Token.CallOpts, symbol)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Token *TokenCallerSession) Tokens(symbol string) (common.Address, error) {
	return _Token.Contract.Tokens(&_Token.CallOpts, symbol)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Token *TokenTransactor) Convert(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "convert", tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Token *TokenSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Token.Contract.Convert(&_Token.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Token *TokenTransactorSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Token.Contract.Convert(&_Token.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Token *TokenTransactor) Convert0(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "convert0", tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Token *TokenSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Token.Contract.Convert0(&_Token.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Token *TokenTransactorSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Token.Contract.Convert0(&_Token.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Token *TokenTransactor) CreateToken(opts *bind.TransactOpts, meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "createToken", meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Token *TokenSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Token.Contract.CreateToken(&_Token.TransactOpts, meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Token *TokenTransactorSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Token.Contract.CreateToken(&_Token.TransactOpts, meta)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Token *TokenTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Token *TokenSession) Initialize() (*types.Transaction, error) {
	return _Token.Contract.Initialize(&_Token.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Token *TokenTransactorSession) Initialize() (*types.Transaction, error) {
	return _Token.Contract.Initialize(&_Token.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _Token.Contract.RenounceOwnership(&_Token.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Token *TokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Token.Contract.RenounceOwnership(&_Token.TransactOpts)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Token *TokenTransactor) SetContractCenter(opts *bind.TransactOpts, addressContractCenter common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "setContractCenter", addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Token *TokenSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Token.Contract.SetContractCenter(&_Token.TransactOpts, addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Token *TokenTransactorSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Token.Contract.SetContractCenter(&_Token.TransactOpts, addressContractCenter)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Token *TokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Token *TokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Token.Contract.TransferOwnership(&_Token.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Token *TokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Token.Contract.TransferOwnership(&_Token.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Token *TokenTransactor) Upgrade(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "upgrade", newImplementation)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Token *TokenSession) Upgrade(newImplementation common.Address) (*types.Transaction, error) {
	return _Token.Contract.Upgrade(&_Token.TransactOpts, newImplementation)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Token *TokenTransactorSession) Upgrade(newImplementation common.Address) (*types.Transaction, error) {
	return _Token.Contract.Upgrade(&_Token.TransactOpts, newImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Token *TokenTransactor) UpgradeToken(opts *bind.TransactOpts, newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Token.contract.Transact(opts, "upgradeToken", newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Token *TokenSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Token.Contract.UpgradeToken(&_Token.TransactOpts, newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Token *TokenTransactorSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Token.Contract.UpgradeToken(&_Token.TransactOpts, newTokenImplementation)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Token.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenSession) Receive() (*types.Transaction, error) {
	return _Token.Contract.Receive(&_Token.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Token *TokenTransactorSession) Receive() (*types.Transaction, error) {
	return _Token.Contract.Receive(&_Token.TransactOpts)
}

// TokenContractUpgradedIterator is returned from FilterContractUpgraded and is used to iterate over the raw logs and unpacked data for ContractUpgraded events raised by the Token contract.
type TokenContractUpgradedIterator struct {
	Event *TokenContractUpgraded // Event containing the contract specifics and raw log

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
func (it *TokenContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenContractUpgraded)
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
		it.Event = new(TokenContractUpgraded)
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
func (it *TokenContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenContractUpgraded represents a ContractUpgraded event raised by the Token contract.
type TokenContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgraded is a free log retrieval operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) FilterContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*TokenContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &TokenContractUpgradedIterator{contract: _Token.contract, event: "ContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchContractUpgraded is a free log subscription operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) WatchContractUpgraded(opts *bind.WatchOpts, sink chan<- *TokenContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenContractUpgraded)
				if err := _Token.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
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

// ParseContractUpgraded is a log parse operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) ParseContractUpgraded(log types.Log) (*TokenContractUpgraded, error) {
	event := new(TokenContractUpgraded)
	if err := _Token.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Token contract.
type TokenInitializedIterator struct {
	Event *TokenInitialized // Event containing the contract specifics and raw log

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
func (it *TokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenInitialized)
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
		it.Event = new(TokenInitialized)
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
func (it *TokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenInitialized represents a Initialized event raised by the Token contract.
type TokenInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Token *TokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*TokenInitializedIterator, error) {

	logs, sub, err := _Token.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TokenInitializedIterator{contract: _Token.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Token *TokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TokenInitialized) (event.Subscription, error) {

	logs, sub, err := _Token.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenInitialized)
				if err := _Token.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Token *TokenFilterer) ParseInitialized(log types.Log) (*TokenInitialized, error) {
	event := new(TokenInitialized)
	if err := _Token.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Token contract.
type TokenOwnershipTransferredIterator struct {
	Event *TokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOwnershipTransferred)
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
		it.Event = new(TokenOwnershipTransferred)
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
func (it *TokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOwnershipTransferred represents a OwnershipTransferred event raised by the Token contract.
type TokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Token *TokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenOwnershipTransferredIterator{contract: _Token.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Token *TokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOwnershipTransferred)
				if err := _Token.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Token *TokenFilterer) ParseOwnershipTransferred(log types.Log) (*TokenOwnershipTransferred, error) {
	event := new(TokenOwnershipTransferred)
	if err := _Token.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenTokenContractUpgradedIterator is returned from FilterTokenContractUpgraded and is used to iterate over the raw logs and unpacked data for TokenContractUpgraded events raised by the Token contract.
type TokenTokenContractUpgradedIterator struct {
	Event *TokenTokenContractUpgraded // Event containing the contract specifics and raw log

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
func (it *TokenTokenContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenTokenContractUpgraded)
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
		it.Event = new(TokenTokenContractUpgraded)
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
func (it *TokenTokenContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenTokenContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenTokenContractUpgraded represents a TokenContractUpgraded event raised by the Token contract.
type TokenTokenContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenContractUpgraded is a free log retrieval operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) FilterTokenContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*TokenTokenContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &TokenTokenContractUpgradedIterator{contract: _Token.contract, event: "TokenContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchTokenContractUpgraded is a free log subscription operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) WatchTokenContractUpgraded(opts *bind.WatchOpts, sink chan<- *TokenTokenContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenTokenContractUpgraded)
				if err := _Token.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
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

// ParseTokenContractUpgraded is a log parse operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Token *TokenFilterer) ParseTokenContractUpgraded(log types.Log) (*TokenTokenContractUpgraded, error) {
	event := new(TokenTokenContractUpgraded)
	if err := _Token.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenTokenDeployedIterator is returned from FilterTokenDeployed and is used to iterate over the raw logs and unpacked data for TokenDeployed events raised by the Token contract.
type TokenTokenDeployedIterator struct {
	Event *TokenTokenDeployed // Event containing the contract specifics and raw log

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
func (it *TokenTokenDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenTokenDeployed)
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
		it.Event = new(TokenTokenDeployed)
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
func (it *TokenTokenDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenTokenDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenTokenDeployed represents a TokenDeployed event raised by the Token contract.
type TokenTokenDeployed struct {
	TokenAddress common.Address
	Meta         DecimalTokenCenterToken
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenDeployed is a free log retrieval operation binding the contract event 0xe17428da190a3db2fb16175ba372d6b241e61e2249e96a18ddac56ed4336aa19.
//
// Solidity: event TokenDeployed(address tokenAddress, (uint256,uint256,uint256,address,uint8,string,string,string) meta)
func (_Token *TokenFilterer) FilterTokenDeployed(opts *bind.FilterOpts) (*TokenTokenDeployedIterator, error) {

	logs, sub, err := _Token.contract.FilterLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return &TokenTokenDeployedIterator{contract: _Token.contract, event: "TokenDeployed", logs: logs, sub: sub}, nil
}

// WatchTokenDeployed is a free log subscription operation binding the contract event 0xe17428da190a3db2fb16175ba372d6b241e61e2249e96a18ddac56ed4336aa19.
//
// Solidity: event TokenDeployed(address tokenAddress, (uint256,uint256,uint256,address,uint8,string,string,string) meta)
func (_Token *TokenFilterer) WatchTokenDeployed(opts *bind.WatchOpts, sink chan<- *TokenTokenDeployed) (event.Subscription, error) {

	logs, sub, err := _Token.contract.WatchLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenTokenDeployed)
				if err := _Token.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
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

// ParseTokenDeployed is a log parse operation binding the contract event 0xe17428da190a3db2fb16175ba372d6b241e61e2249e96a18ddac56ed4336aa19.
//
// Solidity: event TokenDeployed(address tokenAddress, (uint256,uint256,uint256,address,uint8,string,string,string) meta)
func (_Token *TokenFilterer) ParseTokenDeployed(log types.Log) (*TokenTokenDeployed, error) {
	event := new(TokenTokenDeployed)
	if err := _Token.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Token contract.
type TokenUpgradedIterator struct {
	Event *TokenUpgraded // Event containing the contract specifics and raw log

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
func (it *TokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenUpgraded)
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
		it.Event = new(TokenUpgraded)
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
func (it *TokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenUpgraded represents a Upgraded event raised by the Token contract.
type TokenUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Token *TokenFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TokenUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Token.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TokenUpgradedIterator{contract: _Token.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Token *TokenFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TokenUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Token.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenUpgraded)
				if err := _Token.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Token *TokenFilterer) ParseUpgraded(log types.Log) (*TokenUpgraded, error) {
	event := new(TokenUpgraded)
	if err := _Token.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
