// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package center

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

// CenterMetaData contains all meta data concerning the Center contract.
var CenterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedCommission\",\"type\":\"uint256\"}],\"name\":\"InvalidComission\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Log_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenSymbolExist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"TokenContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_INITIAL_MINT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateBuyInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateBuyOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateSellInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateSellOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"crr\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structDecimalTokenCenter.Token\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getCommissionSymbol\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isTokenExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressContractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"tokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTokenImplementation\",\"type\":\"address\"}],\"name\":\"upgradeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[],\"name\":\"isMigrating\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CenterABI is the input ABI used to generate the binding from.
// Deprecated: Use CenterMetaData.ABI instead.
var CenterABI = CenterMetaData.ABI

// Center is an auto generated Go binding around an Ethereum contract.
type Center struct {
	CenterCaller     // Read-only binding to the contract
	CenterTransactor // Write-only binding to the contract
	CenterFilterer   // Log filterer for contract events
}

// CenterCaller is an auto generated read-only Go binding around an Ethereum contract.
type CenterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CenterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CenterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CenterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CenterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CenterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CenterSession struct {
	Contract     *Center           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CenterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CenterCallerSession struct {
	Contract *CenterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CenterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CenterTransactorSession struct {
	Contract     *CenterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CenterRaw is an auto generated low-level Go binding around an Ethereum contract.
type CenterRaw struct {
	Contract *Center // Generic contract binding to access the raw methods on
}

// CenterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CenterCallerRaw struct {
	Contract *CenterCaller // Generic read-only contract binding to access the raw methods on
}

// CenterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CenterTransactorRaw struct {
	Contract *CenterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCenter creates a new instance of Center, bound to a specific deployed contract.
func NewCenter(address common.Address, backend bind.ContractBackend) (*Center, error) {
	contract, err := bindCenter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Center{CenterCaller: CenterCaller{contract: contract}, CenterTransactor: CenterTransactor{contract: contract}, CenterFilterer: CenterFilterer{contract: contract}}, nil
}

// NewCenterCaller creates a new read-only instance of Center, bound to a specific deployed contract.
func NewCenterCaller(address common.Address, caller bind.ContractCaller) (*CenterCaller, error) {
	contract, err := bindCenter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CenterCaller{contract: contract}, nil
}

// NewCenterTransactor creates a new write-only instance of Center, bound to a specific deployed contract.
func NewCenterTransactor(address common.Address, transactor bind.ContractTransactor) (*CenterTransactor, error) {
	contract, err := bindCenter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CenterTransactor{contract: contract}, nil
}

// NewCenterFilterer creates a new log filterer instance of Center, bound to a specific deployed contract.
func NewCenterFilterer(address common.Address, filterer bind.ContractFilterer) (*CenterFilterer, error) {
	contract, err := bindCenter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CenterFilterer{contract: contract}, nil
}

// bindCenter binds a generic wrapper to an already deployed contract.
func bindCenter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CenterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Center *CenterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Center.Contract.CenterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Center *CenterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.Contract.CenterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Center *CenterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Center.Contract.CenterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Center *CenterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Center.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Center *CenterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Center *CenterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Center.Contract.contract.Transact(opts, method, params...)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Center *CenterCaller) MAXCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MAX_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Center *CenterSession) MAXCRR() (*big.Int, error) {
	return _Center.Contract.MAXCRR(&_Center.CallOpts)
}

// MAXCRR is a free data retrieval call binding the contract method 0x090fab49.
//
// Solidity: function MAX_CRR() view returns(uint256)
func (_Center *CenterCallerSession) MAXCRR() (*big.Int, error) {
	return _Center.Contract.MAXCRR(&_Center.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterCaller) MAXTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MAX_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Center.Contract.MAXTOTALSUPPLY(&_Center.CallOpts)
}

// MAXTOTALSUPPLY is a free data retrieval call binding the contract method 0x33039d3d.
//
// Solidity: function MAX_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterCallerSession) MAXTOTALSUPPLY() (*big.Int, error) {
	return _Center.Contract.MAXTOTALSUPPLY(&_Center.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Center *CenterCaller) MINCRR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MIN_CRR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Center *CenterSession) MINCRR() (*big.Int, error) {
	return _Center.Contract.MINCRR(&_Center.CallOpts)
}

// MINCRR is a free data retrieval call binding the contract method 0x4fdcc43c.
//
// Solidity: function MIN_CRR() view returns(uint256)
func (_Center *CenterCallerSession) MINCRR() (*big.Int, error) {
	return _Center.Contract.MINCRR(&_Center.CallOpts)
}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Center *CenterCaller) MININITIALMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MIN_INITIAL_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Center *CenterSession) MININITIALMINT() (*big.Int, error) {
	return _Center.Contract.MININITIALMINT(&_Center.CallOpts)
}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Center *CenterCallerSession) MININITIALMINT() (*big.Int, error) {
	return _Center.Contract.MININITIALMINT(&_Center.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Center *CenterCaller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Center *CenterSession) MINRESERVE() (*big.Int, error) {
	return _Center.Contract.MINRESERVE(&_Center.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Center *CenterCallerSession) MINRESERVE() (*big.Int, error) {
	return _Center.Contract.MINRESERVE(&_Center.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterCaller) MINTOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "MIN_TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Center.Contract.MINTOTALSUPPLY(&_Center.CallOpts)
}

// MINTOTALSUPPLY is a free data retrieval call binding the contract method 0x5122c409.
//
// Solidity: function MIN_TOTAL_SUPPLY() view returns(uint256)
func (_Center *CenterCallerSession) MINTOTALSUPPLY() (*big.Int, error) {
	return _Center.Contract.MINTOTALSUPPLY(&_Center.CallOpts)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterCaller) CalculateBuyInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "calculateBuyInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateBuyInput(&_Center.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyInput is a free data retrieval call binding the contract method 0xb8d23ea0.
//
// Solidity: function calculateBuyInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterCallerSession) CalculateBuyInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateBuyInput(&_Center.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterCaller) CalculateBuyOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "calculateBuyOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateBuyOutput(&_Center.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateBuyOutput is a free data retrieval call binding the contract method 0x15380182.
//
// Solidity: function calculateBuyOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterCallerSession) CalculateBuyOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateBuyOutput(&_Center.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterCaller) CalculateSellInput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "calculateSellInput", supply, customReserve, customCrr, amountOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateSellInput(&_Center.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellInput is a free data retrieval call binding the contract method 0x2e391d02.
//
// Solidity: function calculateSellInput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountOut) pure returns(uint256)
func (_Center *CenterCallerSession) CalculateSellInput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountOut *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateSellInput(&_Center.CallOpts, supply, customReserve, customCrr, amountOut)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterCaller) CalculateSellOutput(opts *bind.CallOpts, supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "calculateSellOutput", supply, customReserve, customCrr, amountIn)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateSellOutput(&_Center.CallOpts, supply, customReserve, customCrr, amountIn)
}

// CalculateSellOutput is a free data retrieval call binding the contract method 0x200ff9f4.
//
// Solidity: function calculateSellOutput(uint256 supply, uint256 customReserve, uint256 customCrr, uint256 amountIn) pure returns(uint256)
func (_Center *CenterCallerSession) CalculateSellOutput(supply *big.Int, customReserve *big.Int, customCrr *big.Int, amountIn *big.Int) (*big.Int, error) {
	return _Center.Contract.CalculateSellOutput(&_Center.CallOpts, supply, customReserve, customCrr, amountIn)
}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Center *CenterCaller) GetCommissionSymbol(opts *bind.CallOpts, symbol string) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getCommissionSymbol", symbol)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Center *CenterSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Center.Contract.GetCommissionSymbol(&_Center.CallOpts, symbol)
}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Center *CenterCallerSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Center.Contract.GetCommissionSymbol(&_Center.CallOpts, symbol)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Center *CenterCaller) GetContractCenter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getContractCenter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Center *CenterSession) GetContractCenter() (common.Address, error) {
	return _Center.Contract.GetContractCenter(&_Center.CallOpts)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Center *CenterCallerSession) GetContractCenter() (common.Address, error) {
	return _Center.Contract.GetContractCenter(&_Center.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Center *CenterCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Center *CenterSession) GetImplementation() (common.Address, error) {
	return _Center.Contract.GetImplementation(&_Center.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Center *CenterCallerSession) GetImplementation() (common.Address, error) {
	return _Center.Contract.GetImplementation(&_Center.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Center *CenterCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Center *CenterSession) Implementation() (common.Address, error) {
	return _Center.Contract.Implementation(&_Center.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Center *CenterCallerSession) Implementation() (common.Address, error) {
	return _Center.Contract.Implementation(&_Center.CallOpts)
}

// IsMigrating is a free data retrieval call binding the contract method 0xf05e777d.
//
// Solidity: function isMigrating() view returns(bool)
func (_Center *CenterCaller) IsMigrating(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "isMigrating")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMigrating is a free data retrieval call binding the contract method 0xf05e777d.
//
// Solidity: function isMigrating() view returns(bool)
func (_Center *CenterSession) IsMigrating() (bool, error) {
	return _Center.Contract.IsMigrating(&_Center.CallOpts)
}

// IsMigrating is a free data retrieval call binding the contract method 0xf05e777d.
//
// Solidity: function isMigrating() view returns(bool)
func (_Center *CenterCallerSession) IsMigrating() (bool, error) {
	return _Center.Contract.IsMigrating(&_Center.CallOpts)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Center *CenterCaller) IsTokenExists(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "isTokenExists", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Center *CenterSession) IsTokenExists(token common.Address) (bool, error) {
	return _Center.Contract.IsTokenExists(&_Center.CallOpts, token)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Center *CenterCallerSession) IsTokenExists(token common.Address) (bool, error) {
	return _Center.Contract.IsTokenExists(&_Center.CallOpts, token)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Center *CenterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Center *CenterSession) Owner() (common.Address, error) {
	return _Center.Contract.Owner(&_Center.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Center *CenterCallerSession) Owner() (common.Address, error) {
	return _Center.Contract.Owner(&_Center.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Center *CenterCaller) Tokens(opts *bind.CallOpts, symbol string) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "tokens", symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Center *CenterSession) Tokens(symbol string) (common.Address, error) {
	return _Center.Contract.Tokens(&_Center.CallOpts, symbol)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Center *CenterCallerSession) Tokens(symbol string) (common.Address, error) {
	return _Center.Contract.Tokens(&_Center.CallOpts, symbol)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Center *CenterTransactor) Convert(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "convert", tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Center *CenterSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Center.Contract.Convert(&_Center.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Center *CenterTransactorSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Center.Contract.Convert(&_Center.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Center *CenterTransactor) Convert0(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "convert0", tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Center *CenterSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Center.Contract.Convert0(&_Center.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Center *CenterTransactorSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Center.Contract.Convert0(&_Center.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Center *CenterTransactor) CreateToken(opts *bind.TransactOpts, meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "createToken", meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Center *CenterSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Center.Contract.CreateToken(&_Center.TransactOpts, meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Center *CenterTransactorSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Center.Contract.CreateToken(&_Center.TransactOpts, meta)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Center *CenterTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Center *CenterSession) Initialize() (*types.Transaction, error) {
	return _Center.Contract.Initialize(&_Center.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Center *CenterTransactorSession) Initialize() (*types.Transaction, error) {
	return _Center.Contract.Initialize(&_Center.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Center *CenterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Center *CenterSession) RenounceOwnership() (*types.Transaction, error) {
	return _Center.Contract.RenounceOwnership(&_Center.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Center *CenterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Center.Contract.RenounceOwnership(&_Center.TransactOpts)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Center *CenterTransactor) SetContractCenter(opts *bind.TransactOpts, addressContractCenter common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "setContractCenter", addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Center *CenterSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetContractCenter(&_Center.TransactOpts, addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Center *CenterTransactorSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetContractCenter(&_Center.TransactOpts, addressContractCenter)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Center *CenterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Center *CenterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Center.Contract.TransferOwnership(&_Center.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Center *CenterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Center.Contract.TransferOwnership(&_Center.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Center *CenterTransactor) Upgrade(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "upgrade", newImplementation)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Center *CenterSession) Upgrade(newImplementation common.Address) (*types.Transaction, error) {
	return _Center.Contract.Upgrade(&_Center.TransactOpts, newImplementation)
}

// Upgrade is a paid mutator transaction binding the contract method 0x0900f010.
//
// Solidity: function upgrade(address newImplementation) returns()
func (_Center *CenterTransactorSession) Upgrade(newImplementation common.Address) (*types.Transaction, error) {
	return _Center.Contract.Upgrade(&_Center.TransactOpts, newImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Center *CenterTransactor) UpgradeToken(opts *bind.TransactOpts, newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "upgradeToken", newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Center *CenterSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Center.Contract.UpgradeToken(&_Center.TransactOpts, newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Center *CenterTransactorSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Center.Contract.UpgradeToken(&_Center.TransactOpts, newTokenImplementation)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Center *CenterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Center *CenterSession) Receive() (*types.Transaction, error) {
	return _Center.Contract.Receive(&_Center.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Center *CenterTransactorSession) Receive() (*types.Transaction, error) {
	return _Center.Contract.Receive(&_Center.TransactOpts)
}

// CenterContractUpgradedIterator is returned from FilterContractUpgraded and is used to iterate over the raw logs and unpacked data for ContractUpgraded events raised by the Center contract.
type CenterContractUpgradedIterator struct {
	Event *CenterContractUpgraded // Event containing the contract specifics and raw log

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
func (it *CenterContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterContractUpgraded)
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
		it.Event = new(CenterContractUpgraded)
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
func (it *CenterContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterContractUpgraded represents a ContractUpgraded event raised by the Center contract.
type CenterContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgraded is a free log retrieval operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Center *CenterFilterer) FilterContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*CenterContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Center.contract.FilterLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &CenterContractUpgradedIterator{contract: _Center.contract, event: "ContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchContractUpgraded is a free log subscription operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Center *CenterFilterer) WatchContractUpgraded(opts *bind.WatchOpts, sink chan<- *CenterContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Center.contract.WatchLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterContractUpgraded)
				if err := _Center.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
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
func (_Center *CenterFilterer) ParseContractUpgraded(log types.Log) (*CenterContractUpgraded, error) {
	event := new(CenterContractUpgraded)
	if err := _Center.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CenterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Center contract.
type CenterInitializedIterator struct {
	Event *CenterInitialized // Event containing the contract specifics and raw log

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
func (it *CenterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterInitialized)
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
		it.Event = new(CenterInitialized)
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
func (it *CenterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterInitialized represents a Initialized event raised by the Center contract.
type CenterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Center *CenterFilterer) FilterInitialized(opts *bind.FilterOpts) (*CenterInitializedIterator, error) {

	logs, sub, err := _Center.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &CenterInitializedIterator{contract: _Center.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Center *CenterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *CenterInitialized) (event.Subscription, error) {

	logs, sub, err := _Center.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterInitialized)
				if err := _Center.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Center *CenterFilterer) ParseInitialized(log types.Log) (*CenterInitialized, error) {
	event := new(CenterInitialized)
	if err := _Center.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CenterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Center contract.
type CenterOwnershipTransferredIterator struct {
	Event *CenterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CenterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterOwnershipTransferred)
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
		it.Event = new(CenterOwnershipTransferred)
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
func (it *CenterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterOwnershipTransferred represents a OwnershipTransferred event raised by the Center contract.
type CenterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Center *CenterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CenterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Center.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CenterOwnershipTransferredIterator{contract: _Center.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Center *CenterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CenterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Center.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterOwnershipTransferred)
				if err := _Center.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Center *CenterFilterer) ParseOwnershipTransferred(log types.Log) (*CenterOwnershipTransferred, error) {
	event := new(CenterOwnershipTransferred)
	if err := _Center.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CenterTokenContractUpgradedIterator is returned from FilterTokenContractUpgraded and is used to iterate over the raw logs and unpacked data for TokenContractUpgraded events raised by the Center contract.
type CenterTokenContractUpgradedIterator struct {
	Event *CenterTokenContractUpgraded // Event containing the contract specifics and raw log

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
func (it *CenterTokenContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterTokenContractUpgraded)
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
		it.Event = new(CenterTokenContractUpgraded)
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
func (it *CenterTokenContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterTokenContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterTokenContractUpgraded represents a TokenContractUpgraded event raised by the Center contract.
type CenterTokenContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenContractUpgraded is a free log retrieval operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Center *CenterFilterer) FilterTokenContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*CenterTokenContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Center.contract.FilterLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &CenterTokenContractUpgradedIterator{contract: _Center.contract, event: "TokenContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchTokenContractUpgraded is a free log subscription operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Center *CenterFilterer) WatchTokenContractUpgraded(opts *bind.WatchOpts, sink chan<- *CenterTokenContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Center.contract.WatchLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterTokenContractUpgraded)
				if err := _Center.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
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
func (_Center *CenterFilterer) ParseTokenContractUpgraded(log types.Log) (*CenterTokenContractUpgraded, error) {
	event := new(CenterTokenContractUpgraded)
	if err := _Center.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CenterTokenDeployedIterator is returned from FilterTokenDeployed and is used to iterate over the raw logs and unpacked data for TokenDeployed events raised by the Center contract.
type CenterTokenDeployedIterator struct {
	Event *CenterTokenDeployed // Event containing the contract specifics and raw log

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
func (it *CenterTokenDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterTokenDeployed)
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
		it.Event = new(CenterTokenDeployed)
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
func (it *CenterTokenDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterTokenDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterTokenDeployed represents a TokenDeployed event raised by the Center contract.
type CenterTokenDeployed struct {
	TokenAddress common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenDeployed is a free log retrieval operation binding the contract event 0x91d24864a084ab70b268a1f865e757ca12006cf298d763b6be697302ef86498c.
//
// Solidity: event TokenDeployed(address tokenAddress)
func (_Center *CenterFilterer) FilterTokenDeployed(opts *bind.FilterOpts) (*CenterTokenDeployedIterator, error) {

	logs, sub, err := _Center.contract.FilterLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return &CenterTokenDeployedIterator{contract: _Center.contract, event: "TokenDeployed", logs: logs, sub: sub}, nil
}

// WatchTokenDeployed is a free log subscription operation binding the contract event 0x91d24864a084ab70b268a1f865e757ca12006cf298d763b6be697302ef86498c.
//
// Solidity: event TokenDeployed(address tokenAddress)
func (_Center *CenterFilterer) WatchTokenDeployed(opts *bind.WatchOpts, sink chan<- *CenterTokenDeployed) (event.Subscription, error) {

	logs, sub, err := _Center.contract.WatchLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterTokenDeployed)
				if err := _Center.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
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

// ParseTokenDeployed is a log parse operation binding the contract event 0x91d24864a084ab70b268a1f865e757ca12006cf298d763b6be697302ef86498c.
//
// Solidity: event TokenDeployed(address tokenAddress)
func (_Center *CenterFilterer) ParseTokenDeployed(log types.Log) (*CenterTokenDeployed, error) {
	event := new(CenterTokenDeployed)
	if err := _Center.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CenterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Center contract.
type CenterUpgradedIterator struct {
	Event *CenterUpgraded // Event containing the contract specifics and raw log

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
func (it *CenterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterUpgraded)
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
		it.Event = new(CenterUpgraded)
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
func (it *CenterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterUpgraded represents a Upgraded event raised by the Center contract.
type CenterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Center *CenterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*CenterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Center.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &CenterUpgradedIterator{contract: _Center.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Center *CenterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *CenterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Center.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterUpgraded)
				if err := _Center.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Center *CenterFilterer) ParseUpgraded(log types.Log) (*CenterUpgraded, error) {
	event := new(CenterUpgraded)
	if err := _Center.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
