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

// Reserve is an auto generated low-level Go binding around an user-defined struct.
type Reserve struct {
	ContractAddress common.Address
	Amount          *big.Int
}

// CenterMetaData contains all meta data concerning the Center contract.
var CenterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllReservesAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AvailableOnlyDuringMigration\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DisabledDuringMigration\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FailedToMigrateReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStartIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStartOrEndIndex\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"provided\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"}],\"name\":\"InvalidValueToMigrateReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ListOfContractsIsEpmty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NothingToMigrate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PayloadMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"remainingMsgValue\",\"type\":\"uint256\"}],\"name\":\"SomethingGoesWrong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"ContractAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImpl\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"getMarkedReserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalMarkedReserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isMigrating\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structReserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"isReserveMigrated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structReserve[]\",\"name\":\"reserves\",\"type\":\"tuple[]\"}],\"name\":\"isReservesMigrated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structReserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"markReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structReserve[]\",\"name\":\"reserves\",\"type\":\"tuple[]\"}],\"name\":\"markReserves\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"migrateReserves\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"migrateReservesPage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"setAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"setContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"newIsMigrating\",\"type\":\"bool\"}],\"name\":\"setMigrating\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImpl\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string symbol) view returns(address)
func (_Center *CenterCaller) GetAddress(opts *bind.CallOpts, symbol string) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getAddress", symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string symbol) view returns(address)
func (_Center *CenterSession) GetAddress(symbol string) (common.Address, error) {
	return _Center.Contract.GetAddress(&_Center.CallOpts, symbol)
}

// GetAddress is a free data retrieval call binding the contract method 0xbf40fac1.
//
// Solidity: function getAddress(string symbol) view returns(address)
func (_Center *CenterCallerSession) GetAddress(symbol string) (common.Address, error) {
	return _Center.Contract.GetAddress(&_Center.CallOpts, symbol)
}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string symbol) view returns(address)
func (_Center *CenterCaller) GetContractAddress(opts *bind.CallOpts, symbol string) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getContractAddress", symbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string symbol) view returns(address)
func (_Center *CenterSession) GetContractAddress(symbol string) (common.Address, error) {
	return _Center.Contract.GetContractAddress(&_Center.CallOpts, symbol)
}

// GetContractAddress is a free data retrieval call binding the contract method 0x04433bbc.
//
// Solidity: function getContractAddress(string symbol) view returns(address)
func (_Center *CenterCallerSession) GetContractAddress(symbol string) (common.Address, error) {
	return _Center.Contract.GetContractAddress(&_Center.CallOpts, symbol)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Center *CenterCaller) GetImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Center *CenterSession) GetImpl() (common.Address, error) {
	return _Center.Contract.GetImpl(&_Center.CallOpts)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Center *CenterCallerSession) GetImpl() (common.Address, error) {
	return _Center.Contract.GetImpl(&_Center.CallOpts)
}

// GetMarkedReserve is a free data retrieval call binding the contract method 0x95d9d943.
//
// Solidity: function getMarkedReserve(address contractAddress) view returns(uint256)
func (_Center *CenterCaller) GetMarkedReserve(opts *bind.CallOpts, contractAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getMarkedReserve", contractAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMarkedReserve is a free data retrieval call binding the contract method 0x95d9d943.
//
// Solidity: function getMarkedReserve(address contractAddress) view returns(uint256)
func (_Center *CenterSession) GetMarkedReserve(contractAddress common.Address) (*big.Int, error) {
	return _Center.Contract.GetMarkedReserve(&_Center.CallOpts, contractAddress)
}

// GetMarkedReserve is a free data retrieval call binding the contract method 0x95d9d943.
//
// Solidity: function getMarkedReserve(address contractAddress) view returns(uint256)
func (_Center *CenterCallerSession) GetMarkedReserve(contractAddress common.Address) (*big.Int, error) {
	return _Center.Contract.GetMarkedReserve(&_Center.CallOpts, contractAddress)
}

// GetTotalMarkedReserve is a free data retrieval call binding the contract method 0x9c6da6ef.
//
// Solidity: function getTotalMarkedReserve() view returns(uint256)
func (_Center *CenterCaller) GetTotalMarkedReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "getTotalMarkedReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalMarkedReserve is a free data retrieval call binding the contract method 0x9c6da6ef.
//
// Solidity: function getTotalMarkedReserve() view returns(uint256)
func (_Center *CenterSession) GetTotalMarkedReserve() (*big.Int, error) {
	return _Center.Contract.GetTotalMarkedReserve(&_Center.CallOpts)
}

// GetTotalMarkedReserve is a free data retrieval call binding the contract method 0x9c6da6ef.
//
// Solidity: function getTotalMarkedReserve() view returns(uint256)
func (_Center *CenterCallerSession) GetTotalMarkedReserve() (*big.Int, error) {
	return _Center.Contract.GetTotalMarkedReserve(&_Center.CallOpts)
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

// IsReserveMigrated is a free data retrieval call binding the contract method 0xf144f040.
//
// Solidity: function isReserveMigrated((address,uint256) reserve) view returns(bool)
func (_Center *CenterCaller) IsReserveMigrated(opts *bind.CallOpts, reserve Reserve) (bool, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "isReserveMigrated", reserve)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsReserveMigrated is a free data retrieval call binding the contract method 0xf144f040.
//
// Solidity: function isReserveMigrated((address,uint256) reserve) view returns(bool)
func (_Center *CenterSession) IsReserveMigrated(reserve Reserve) (bool, error) {
	return _Center.Contract.IsReserveMigrated(&_Center.CallOpts, reserve)
}

// IsReserveMigrated is a free data retrieval call binding the contract method 0xf144f040.
//
// Solidity: function isReserveMigrated((address,uint256) reserve) view returns(bool)
func (_Center *CenterCallerSession) IsReserveMigrated(reserve Reserve) (bool, error) {
	return _Center.Contract.IsReserveMigrated(&_Center.CallOpts, reserve)
}

// IsReservesMigrated is a free data retrieval call binding the contract method 0xbfb25601.
//
// Solidity: function isReservesMigrated((address,uint256)[] reserves) view returns(bool)
func (_Center *CenterCaller) IsReservesMigrated(opts *bind.CallOpts, reserves []Reserve) (bool, error) {
	var out []interface{}
	err := _Center.contract.Call(opts, &out, "isReservesMigrated", reserves)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsReservesMigrated is a free data retrieval call binding the contract method 0xbfb25601.
//
// Solidity: function isReservesMigrated((address,uint256)[] reserves) view returns(bool)
func (_Center *CenterSession) IsReservesMigrated(reserves []Reserve) (bool, error) {
	return _Center.Contract.IsReservesMigrated(&_Center.CallOpts, reserves)
}

// IsReservesMigrated is a free data retrieval call binding the contract method 0xbfb25601.
//
// Solidity: function isReservesMigrated((address,uint256)[] reserves) view returns(bool)
func (_Center *CenterCallerSession) IsReservesMigrated(reserves []Reserve) (bool, error) {
	return _Center.Contract.IsReservesMigrated(&_Center.CallOpts, reserves)
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

// MarkReserve is a paid mutator transaction binding the contract method 0x73274bce.
//
// Solidity: function markReserve((address,uint256) reserve) returns()
func (_Center *CenterTransactor) MarkReserve(opts *bind.TransactOpts, reserve Reserve) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "markReserve", reserve)
}

// MarkReserve is a paid mutator transaction binding the contract method 0x73274bce.
//
// Solidity: function markReserve((address,uint256) reserve) returns()
func (_Center *CenterSession) MarkReserve(reserve Reserve) (*types.Transaction, error) {
	return _Center.Contract.MarkReserve(&_Center.TransactOpts, reserve)
}

// MarkReserve is a paid mutator transaction binding the contract method 0x73274bce.
//
// Solidity: function markReserve((address,uint256) reserve) returns()
func (_Center *CenterTransactorSession) MarkReserve(reserve Reserve) (*types.Transaction, error) {
	return _Center.Contract.MarkReserve(&_Center.TransactOpts, reserve)
}

// MarkReserves is a paid mutator transaction binding the contract method 0x55c82602.
//
// Solidity: function markReserves((address,uint256)[] reserves) returns()
func (_Center *CenterTransactor) MarkReserves(opts *bind.TransactOpts, reserves []Reserve) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "markReserves", reserves)
}

// MarkReserves is a paid mutator transaction binding the contract method 0x55c82602.
//
// Solidity: function markReserves((address,uint256)[] reserves) returns()
func (_Center *CenterSession) MarkReserves(reserves []Reserve) (*types.Transaction, error) {
	return _Center.Contract.MarkReserves(&_Center.TransactOpts, reserves)
}

// MarkReserves is a paid mutator transaction binding the contract method 0x55c82602.
//
// Solidity: function markReserves((address,uint256)[] reserves) returns()
func (_Center *CenterTransactorSession) MarkReserves(reserves []Reserve) (*types.Transaction, error) {
	return _Center.Contract.MarkReserves(&_Center.TransactOpts, reserves)
}

// MigrateReserves is a paid mutator transaction binding the contract method 0x41b63c6f.
//
// Solidity: function migrateReserves() payable returns()
func (_Center *CenterTransactor) MigrateReserves(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "migrateReserves")
}

// MigrateReserves is a paid mutator transaction binding the contract method 0x41b63c6f.
//
// Solidity: function migrateReserves() payable returns()
func (_Center *CenterSession) MigrateReserves() (*types.Transaction, error) {
	return _Center.Contract.MigrateReserves(&_Center.TransactOpts)
}

// MigrateReserves is a paid mutator transaction binding the contract method 0x41b63c6f.
//
// Solidity: function migrateReserves() payable returns()
func (_Center *CenterTransactorSession) MigrateReserves() (*types.Transaction, error) {
	return _Center.Contract.MigrateReserves(&_Center.TransactOpts)
}

// MigrateReservesPage is a paid mutator transaction binding the contract method 0xc4d678af.
//
// Solidity: function migrateReservesPage(uint256 start, uint256 end) payable returns()
func (_Center *CenterTransactor) MigrateReservesPage(opts *bind.TransactOpts, start *big.Int, end *big.Int) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "migrateReservesPage", start, end)
}

// MigrateReservesPage is a paid mutator transaction binding the contract method 0xc4d678af.
//
// Solidity: function migrateReservesPage(uint256 start, uint256 end) payable returns()
func (_Center *CenterSession) MigrateReservesPage(start *big.Int, end *big.Int) (*types.Transaction, error) {
	return _Center.Contract.MigrateReservesPage(&_Center.TransactOpts, start, end)
}

// MigrateReservesPage is a paid mutator transaction binding the contract method 0xc4d678af.
//
// Solidity: function migrateReservesPage(uint256 start, uint256 end) payable returns()
func (_Center *CenterTransactorSession) MigrateReservesPage(start *big.Int, end *big.Int) (*types.Transaction, error) {
	return _Center.Contract.MigrateReservesPage(&_Center.TransactOpts, start, end)
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

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string symbol, address contractAddress) returns()
func (_Center *CenterTransactor) SetAddress(opts *bind.TransactOpts, symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "setAddress", symbol, contractAddress)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string symbol, address contractAddress) returns()
func (_Center *CenterSession) SetAddress(symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetAddress(&_Center.TransactOpts, symbol, contractAddress)
}

// SetAddress is a paid mutator transaction binding the contract method 0x9b2ea4bd.
//
// Solidity: function setAddress(string symbol, address contractAddress) returns()
func (_Center *CenterTransactorSession) SetAddress(symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetAddress(&_Center.TransactOpts, symbol, contractAddress)
}

// SetContractAddress is a paid mutator transaction binding the contract method 0x534e785c.
//
// Solidity: function setContractAddress(string symbol, address contractAddress) returns()
func (_Center *CenterTransactor) SetContractAddress(opts *bind.TransactOpts, symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "setContractAddress", symbol, contractAddress)
}

// SetContractAddress is a paid mutator transaction binding the contract method 0x534e785c.
//
// Solidity: function setContractAddress(string symbol, address contractAddress) returns()
func (_Center *CenterSession) SetContractAddress(symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetContractAddress(&_Center.TransactOpts, symbol, contractAddress)
}

// SetContractAddress is a paid mutator transaction binding the contract method 0x534e785c.
//
// Solidity: function setContractAddress(string symbol, address contractAddress) returns()
func (_Center *CenterTransactorSession) SetContractAddress(symbol string, contractAddress common.Address) (*types.Transaction, error) {
	return _Center.Contract.SetContractAddress(&_Center.TransactOpts, symbol, contractAddress)
}

// SetMigrating is a paid mutator transaction binding the contract method 0xf785f03d.
//
// Solidity: function setMigrating(bool newIsMigrating) returns()
func (_Center *CenterTransactor) SetMigrating(opts *bind.TransactOpts, newIsMigrating bool) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "setMigrating", newIsMigrating)
}

// SetMigrating is a paid mutator transaction binding the contract method 0xf785f03d.
//
// Solidity: function setMigrating(bool newIsMigrating) returns()
func (_Center *CenterSession) SetMigrating(newIsMigrating bool) (*types.Transaction, error) {
	return _Center.Contract.SetMigrating(&_Center.TransactOpts, newIsMigrating)
}

// SetMigrating is a paid mutator transaction binding the contract method 0xf785f03d.
//
// Solidity: function setMigrating(bool newIsMigrating) returns()
func (_Center *CenterTransactorSession) SetMigrating(newIsMigrating bool) (*types.Transaction, error) {
	return _Center.Contract.SetMigrating(&_Center.TransactOpts, newIsMigrating)
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

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Center *CenterTransactor) Upgrade(opts *bind.TransactOpts, newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Center.contract.Transact(opts, "upgrade", newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Center *CenterSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Center.Contract.Upgrade(&_Center.TransactOpts, newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Center *CenterTransactorSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Center.Contract.Upgrade(&_Center.TransactOpts, newImpl, data)
}

// CenterContractAddedIterator is returned from FilterContractAdded and is used to iterate over the raw logs and unpacked data for ContractAdded events raised by the Center contract.
type CenterContractAddedIterator struct {
	Event *CenterContractAdded // Event containing the contract specifics and raw log

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
func (it *CenterContractAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CenterContractAdded)
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
		it.Event = new(CenterContractAdded)
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
func (it *CenterContractAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CenterContractAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CenterContractAdded represents a ContractAdded event raised by the Center contract.
type CenterContractAdded struct {
	Symbol          string
	ContractAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterContractAdded is a free log retrieval operation binding the contract event 0x8b4ef7d4e5bc8f098e6f637ac0acf4aee47b3f027efea6307264b06b4bc9d298.
//
// Solidity: event ContractAdded(string symbol, address indexed contractAddress)
func (_Center *CenterFilterer) FilterContractAdded(opts *bind.FilterOpts, contractAddress []common.Address) (*CenterContractAddedIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Center.contract.FilterLogs(opts, "ContractAdded", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &CenterContractAddedIterator{contract: _Center.contract, event: "ContractAdded", logs: logs, sub: sub}, nil
}

// WatchContractAdded is a free log subscription operation binding the contract event 0x8b4ef7d4e5bc8f098e6f637ac0acf4aee47b3f027efea6307264b06b4bc9d298.
//
// Solidity: event ContractAdded(string symbol, address indexed contractAddress)
func (_Center *CenterFilterer) WatchContractAdded(opts *bind.WatchOpts, sink chan<- *CenterContractAdded, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Center.contract.WatchLogs(opts, "ContractAdded", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CenterContractAdded)
				if err := _Center.contract.UnpackLog(event, "ContractAdded", log); err != nil {
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

// ParseContractAdded is a log parse operation binding the contract event 0x8b4ef7d4e5bc8f098e6f637ac0acf4aee47b3f027efea6307264b06b4bc9d298.
//
// Solidity: event ContractAdded(string symbol, address indexed contractAddress)
func (_Center *CenterFilterer) ParseContractAdded(log types.Log) (*CenterContractAdded, error) {
	event := new(CenterContractAdded)
	if err := _Center.contract.UnpackLog(event, "ContractAdded", log); err != nil {
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
