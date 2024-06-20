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

// DecimalTokenCenterToken is an auto generated low-level Go binding around an user-defined struct.
type DecimalTokenCenterToken struct {
	InitialMint    *big.Int
	MinTotalSupply *big.Int
	MaxTotalSupply *big.Int
	Creator        common.Address
	Crr            uint8
	Identity       string
	Symbol         string
	Name           string
}

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expectedCommission\",\"type\":\"uint256\"}],\"name\":\"InvalidComission\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Log_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenSymbolExist\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"TokenContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"crr\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structDecimalTokenCenter.Token\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"TokenDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_CRR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_INITIAL_MINT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateBuyInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateBuyOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"calculateSellInput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"customCrr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"name\":\"calculateSellOutput\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"convert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"initialMint\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"crr\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"internalType\":\"structDecimalTokenCenter.Token\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getCommissionSymbol\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isTokenExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressContractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"tokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTokenImplementation\",\"type\":\"address\"}],\"name\":\"upgradeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Delegation *DelegationCaller) MININITIALMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "MIN_INITIAL_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Delegation *DelegationSession) MININITIALMINT() (*big.Int, error) {
	return _Delegation.Contract.MININITIALMINT(&_Delegation.CallOpts)
}

// MININITIALMINT is a free data retrieval call binding the contract method 0xc265bbd1.
//
// Solidity: function MIN_INITIAL_MINT() view returns(uint256)
func (_Delegation *DelegationCallerSession) MININITIALMINT() (*big.Int, error) {
	return _Delegation.Contract.MININITIALMINT(&_Delegation.CallOpts)
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

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Delegation *DelegationCaller) GetCommissionSymbol(opts *bind.CallOpts, symbol string) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getCommissionSymbol", symbol)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Delegation *DelegationSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Delegation.Contract.GetCommissionSymbol(&_Delegation.CallOpts, symbol)
}

// GetCommissionSymbol is a free data retrieval call binding the contract method 0xc73636c1.
//
// Solidity: function getCommissionSymbol(string symbol) pure returns(uint256)
func (_Delegation *DelegationCallerSession) GetCommissionSymbol(symbol string) (*big.Int, error) {
	return _Delegation.Contract.GetCommissionSymbol(&_Delegation.CallOpts, symbol)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Delegation *DelegationCaller) GetContractCenter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getContractCenter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Delegation *DelegationSession) GetContractCenter() (common.Address, error) {
	return _Delegation.Contract.GetContractCenter(&_Delegation.CallOpts)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_Delegation *DelegationCallerSession) GetContractCenter() (common.Address, error) {
	return _Delegation.Contract.GetContractCenter(&_Delegation.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Delegation *DelegationCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Delegation *DelegationSession) GetImplementation() (common.Address, error) {
	return _Delegation.Contract.GetImplementation(&_Delegation.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_Delegation *DelegationCallerSession) GetImplementation() (common.Address, error) {
	return _Delegation.Contract.GetImplementation(&_Delegation.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Delegation *DelegationCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Delegation *DelegationSession) Implementation() (common.Address, error) {
	return _Delegation.Contract.Implementation(&_Delegation.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_Delegation *DelegationCallerSession) Implementation() (common.Address, error) {
	return _Delegation.Contract.Implementation(&_Delegation.CallOpts)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Delegation *DelegationCaller) IsTokenExists(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "isTokenExists", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Delegation *DelegationSession) IsTokenExists(token common.Address) (bool, error) {
	return _Delegation.Contract.IsTokenExists(&_Delegation.CallOpts, token)
}

// IsTokenExists is a free data retrieval call binding the contract method 0x9ed4fa5a.
//
// Solidity: function isTokenExists(address token) view returns(bool)
func (_Delegation *DelegationCallerSession) IsTokenExists(token common.Address) (bool, error) {
	return _Delegation.Contract.IsTokenExists(&_Delegation.CallOpts, token)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegation *DelegationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegation *DelegationSession) Owner() (common.Address, error) {
	return _Delegation.Contract.Owner(&_Delegation.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegation *DelegationCallerSession) Owner() (common.Address, error) {
	return _Delegation.Contract.Owner(&_Delegation.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Delegation *DelegationCaller) Tokens(opts *bind.CallOpts, symbol string) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "tokens", symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Delegation *DelegationSession) Tokens(symbol string) (common.Address, error) {
	return _Delegation.Contract.Tokens(&_Delegation.CallOpts, symbol)
}

// Tokens is a free data retrieval call binding the contract method 0x04c2320b.
//
// Solidity: function tokens(string symbol) view returns(address)
func (_Delegation *DelegationCallerSession) Tokens(symbol string) (common.Address, error) {
	return _Delegation.Contract.Tokens(&_Delegation.CallOpts, symbol)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactor) Convert(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "convert", tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Convert(&_Delegation.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert is a paid mutator transaction binding the contract method 0x069ffbe5.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient) returns()
func (_Delegation *DelegationTransactorSession) Convert(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Convert(&_Delegation.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) Convert0(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "convert0", tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Convert0(&_Delegation.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// Convert0 is a paid mutator transaction binding the contract method 0xb7a7a048.
//
// Solidity: function convert(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) Convert0(tokenIn common.Address, tokenOut common.Address, amountIn *big.Int, amountOutMin *big.Int, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Convert0(&_Delegation.TransactOpts, tokenIn, tokenOut, amountIn, amountOutMin, recipient, deadline, v, r, s)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Delegation *DelegationTransactor) CreateToken(opts *bind.TransactOpts, meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "createToken", meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Delegation *DelegationSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Delegation.Contract.CreateToken(&_Delegation.TransactOpts, meta)
}

// CreateToken is a paid mutator transaction binding the contract method 0x4dfb5142.
//
// Solidity: function createToken((uint256,uint256,uint256,address,uint8,string,string,string) meta) payable returns()
func (_Delegation *DelegationTransactorSession) CreateToken(meta DecimalTokenCenterToken) (*types.Transaction, error) {
	return _Delegation.Contract.CreateToken(&_Delegation.TransactOpts, meta)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Delegation *DelegationTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Delegation *DelegationSession) Initialize() (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Delegation *DelegationTransactorSession) Initialize() (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegation *DelegationTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegation *DelegationSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegation.Contract.RenounceOwnership(&_Delegation.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegation *DelegationTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegation.Contract.RenounceOwnership(&_Delegation.TransactOpts)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Delegation *DelegationTransactor) SetContractCenter(opts *bind.TransactOpts, addressContractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "setContractCenter", addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Delegation *DelegationSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SetContractCenter(&_Delegation.TransactOpts, addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_Delegation *DelegationTransactorSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SetContractCenter(&_Delegation.TransactOpts, addressContractCenter)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegation *DelegationTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegation *DelegationSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.TransferOwnership(&_Delegation.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegation *DelegationTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.TransferOwnership(&_Delegation.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_Delegation *DelegationTransactor) Upgrade(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "upgrade", newImplementation, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_Delegation *DelegationSession) Upgrade(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.Upgrade(&_Delegation.TransactOpts, newImplementation, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_Delegation *DelegationTransactorSession) Upgrade(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.Upgrade(&_Delegation.TransactOpts, newImplementation, data)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Delegation *DelegationTransactor) UpgradeToken(opts *bind.TransactOpts, newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "upgradeToken", newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Delegation *DelegationSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.UpgradeToken(&_Delegation.TransactOpts, newTokenImplementation)
}

// UpgradeToken is a paid mutator transaction binding the contract method 0x6ee31a18.
//
// Solidity: function upgradeToken(address newTokenImplementation) returns()
func (_Delegation *DelegationTransactorSession) UpgradeToken(newTokenImplementation common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.UpgradeToken(&_Delegation.TransactOpts, newTokenImplementation)
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

// DelegationContractUpgradedIterator is returned from FilterContractUpgraded and is used to iterate over the raw logs and unpacked data for ContractUpgraded events raised by the Delegation contract.
type DelegationContractUpgradedIterator struct {
	Event *DelegationContractUpgraded // Event containing the contract specifics and raw log

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
func (it *DelegationContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationContractUpgraded)
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
		it.Event = new(DelegationContractUpgraded)
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
func (it *DelegationContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationContractUpgraded represents a ContractUpgraded event raised by the Delegation contract.
type DelegationContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgraded is a free log retrieval operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Delegation *DelegationFilterer) FilterContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*DelegationContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &DelegationContractUpgradedIterator{contract: _Delegation.contract, event: "ContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchContractUpgraded is a free log subscription operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Delegation *DelegationFilterer) WatchContractUpgraded(opts *bind.WatchOpts, sink chan<- *DelegationContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationContractUpgraded)
				if err := _Delegation.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseContractUpgraded(log types.Log) (*DelegationContractUpgraded, error) {
	event := new(DelegationContractUpgraded)
	if err := _Delegation.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
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

// DelegationOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Delegation contract.
type DelegationOwnershipTransferredIterator struct {
	Event *DelegationOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DelegationOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationOwnershipTransferred)
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
		it.Event = new(DelegationOwnershipTransferred)
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
func (it *DelegationOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationOwnershipTransferred represents a OwnershipTransferred event raised by the Delegation contract.
type DelegationOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegation *DelegationFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DelegationOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DelegationOwnershipTransferredIterator{contract: _Delegation.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegation *DelegationFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DelegationOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationOwnershipTransferred)
				if err := _Delegation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseOwnershipTransferred(log types.Log) (*DelegationOwnershipTransferred, error) {
	event := new(DelegationOwnershipTransferred)
	if err := _Delegation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTokenContractUpgradedIterator is returned from FilterTokenContractUpgraded and is used to iterate over the raw logs and unpacked data for TokenContractUpgraded events raised by the Delegation contract.
type DelegationTokenContractUpgradedIterator struct {
	Event *DelegationTokenContractUpgraded // Event containing the contract specifics and raw log

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
func (it *DelegationTokenContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTokenContractUpgraded)
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
		it.Event = new(DelegationTokenContractUpgraded)
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
func (it *DelegationTokenContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTokenContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTokenContractUpgraded represents a TokenContractUpgraded event raised by the Delegation contract.
type DelegationTokenContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenContractUpgraded is a free log retrieval operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Delegation *DelegationFilterer) FilterTokenContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*DelegationTokenContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &DelegationTokenContractUpgradedIterator{contract: _Delegation.contract, event: "TokenContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchTokenContractUpgraded is a free log subscription operation binding the contract event 0x27bf8a17dff3ae6812ef6a2059d654c298fd3a87c570f2bab5c34b166dd868aa.
//
// Solidity: event TokenContractUpgraded(address indexed oldContract, address indexed newContract)
func (_Delegation *DelegationFilterer) WatchTokenContractUpgraded(opts *bind.WatchOpts, sink chan<- *DelegationTokenContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TokenContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTokenContractUpgraded)
				if err := _Delegation.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseTokenContractUpgraded(log types.Log) (*DelegationTokenContractUpgraded, error) {
	event := new(DelegationTokenContractUpgraded)
	if err := _Delegation.contract.UnpackLog(event, "TokenContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTokenDeployedIterator is returned from FilterTokenDeployed and is used to iterate over the raw logs and unpacked data for TokenDeployed events raised by the Delegation contract.
type DelegationTokenDeployedIterator struct {
	Event *DelegationTokenDeployed // Event containing the contract specifics and raw log

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
func (it *DelegationTokenDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTokenDeployed)
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
		it.Event = new(DelegationTokenDeployed)
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
func (it *DelegationTokenDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTokenDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTokenDeployed represents a TokenDeployed event raised by the Delegation contract.
type DelegationTokenDeployed struct {
	TokenAddress common.Address
	Meta         DecimalTokenCenterToken
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenDeployed is a free log retrieval operation binding the contract event 0xe17428da190a3db2fb16175ba372d6b241e61e2249e96a18ddac56ed4336aa19.
//
// Solidity: event TokenDeployed(address tokenAddress, (uint256,uint256,uint256,address,uint8,string,string,string) meta)
func (_Delegation *DelegationFilterer) FilterTokenDeployed(opts *bind.FilterOpts) (*DelegationTokenDeployedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return &DelegationTokenDeployedIterator{contract: _Delegation.contract, event: "TokenDeployed", logs: logs, sub: sub}, nil
}

// WatchTokenDeployed is a free log subscription operation binding the contract event 0xe17428da190a3db2fb16175ba372d6b241e61e2249e96a18ddac56ed4336aa19.
//
// Solidity: event TokenDeployed(address tokenAddress, (uint256,uint256,uint256,address,uint8,string,string,string) meta)
func (_Delegation *DelegationFilterer) WatchTokenDeployed(opts *bind.WatchOpts, sink chan<- *DelegationTokenDeployed) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TokenDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTokenDeployed)
				if err := _Delegation.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseTokenDeployed(log types.Log) (*DelegationTokenDeployed, error) {
	event := new(DelegationTokenDeployed)
	if err := _Delegation.contract.UnpackLog(event, "TokenDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Delegation contract.
type DelegationUpgradedIterator struct {
	Event *DelegationUpgraded // Event containing the contract specifics and raw log

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
func (it *DelegationUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationUpgraded)
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
		it.Event = new(DelegationUpgraded)
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
func (it *DelegationUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationUpgraded represents a Upgraded event raised by the Delegation contract.
type DelegationUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Delegation *DelegationFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DelegationUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DelegationUpgradedIterator{contract: _Delegation.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Delegation *DelegationFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DelegationUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationUpgraded)
				if err := _Delegation.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseUpgraded(log types.Log) (*DelegationUpgraded, error) {
	event := new(DelegationUpgraded)
	if err := _Delegation.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
