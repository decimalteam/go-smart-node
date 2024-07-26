// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validator

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

// IDecimalMasterValidatorStake is an auto generated low-level Go binding around an user-defined struct.
type IDecimalMasterValidatorStake struct {
	Token  common.Address
	Amount *big.Int
}

// IDecimalMasterValidatorValidator is an auto generated low-level Go binding around an user-defined struct.
type IDecimalMasterValidatorValidator struct {
	Status             uint8
	Paused             bool
	PenaltyPercantages []uint16
}

// IDecimalMasterValidatorValidatorMeta is an auto generated low-level Go binding around an user-defined struct.
type IDecimalMasterValidatorValidatorMeta struct {
	Validator common.Address
	Meta      string
}

// ValidatorMetaData contains all meta data concerning the Validator contract.
var ValidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyMember\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPenalty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStatus\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotMember\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"}],\"name\":\"ValidatorMetaUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"percentage\",\"type\":\"uint256\"}],\"name\":\"ValidatorPenalty\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIDecimalMasterValidator.ValidatorStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"ValidatorUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalMasterValidator.Stake\",\"name\":\"initialStake\",\"type\":\"tuple\"}],\"name\":\"addCandidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"}],\"name\":\"addCandidateDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"penaltyPercentage\",\"type\":\"uint16\"}],\"name\":\"addPenalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalMasterValidator.Stake\",\"name\":\"initialStake\",\"type\":\"tuple\"}],\"name\":\"addValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"}],\"name\":\"addValidatorDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"}],\"internalType\":\"structIDecimalMasterValidator.ValidatorMeta[]\",\"name\":\"metas\",\"type\":\"tuple[]\"}],\"name\":\"addValidators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"candidates\",\"type\":\"address[]\"}],\"name\":\"approveCandidates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImpl\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIDecimalMasterValidator.ValidatorStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"},{\"internalType\":\"uint16[]\",\"name\":\"penaltyPercantages\",\"type\":\"uint16[]\"}],\"internalType\":\"structIDecimalMasterValidator.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getValidatorPenalties\",\"outputs\":[{\"internalType\":\"uint16[]\",\"name\":\"\",\"type\":\"uint16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"getValidatorStatus\",\"outputs\":[{\"internalType\":\"enumIDecimalMasterValidator.ValidatorStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"isMember\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseSelf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpauseSelf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"unpauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"meta\",\"type\":\"string\"}],\"name\":\"updateValidatorMeta\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImpl\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidatorMetaData.ABI instead.
var ValidatorABI = ValidatorMetaData.ABI

// Validator is an auto generated Go binding around an Ethereum contract.
type Validator struct {
	ValidatorCaller     // Read-only binding to the contract
	ValidatorTransactor // Write-only binding to the contract
	ValidatorFilterer   // Log filterer for contract events
}

// ValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorSession struct {
	Contract     *Validator        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorCallerSession struct {
	Contract *ValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorTransactorSession struct {
	Contract     *ValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorRaw struct {
	Contract *Validator // Generic contract binding to access the raw methods on
}

// ValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorCallerRaw struct {
	Contract *ValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorTransactorRaw struct {
	Contract *ValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidator creates a new instance of Validator, bound to a specific deployed contract.
func NewValidator(address common.Address, backend bind.ContractBackend) (*Validator, error) {
	contract, err := bindValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validator{ValidatorCaller: ValidatorCaller{contract: contract}, ValidatorTransactor: ValidatorTransactor{contract: contract}, ValidatorFilterer: ValidatorFilterer{contract: contract}}, nil
}

// NewValidatorCaller creates a new read-only instance of Validator, bound to a specific deployed contract.
func NewValidatorCaller(address common.Address, caller bind.ContractCaller) (*ValidatorCaller, error) {
	contract, err := bindValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorCaller{contract: contract}, nil
}

// NewValidatorTransactor creates a new write-only instance of Validator, bound to a specific deployed contract.
func NewValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorTransactor, error) {
	contract, err := bindValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorTransactor{contract: contract}, nil
}

// NewValidatorFilterer creates a new log filterer instance of Validator, bound to a specific deployed contract.
func NewValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorFilterer, error) {
	contract, err := bindValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorFilterer{contract: contract}, nil
}

// bindValidator binds a generic wrapper to an already deployed contract.
func bindValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validator *ValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validator.Contract.ValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validator *ValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validator.Contract.ValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validator *ValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validator.Contract.ValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validator *ValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validator *ValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validator *ValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validator.Contract.contract.Transact(opts, method, params...)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Validator *ValidatorCaller) GetImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "getImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Validator *ValidatorSession) GetImpl() (common.Address, error) {
	return _Validator.Contract.GetImpl(&_Validator.CallOpts)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Validator *ValidatorCallerSession) GetImpl() (common.Address, error) {
	return _Validator.Contract.GetImpl(&_Validator.CallOpts)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address validator) view returns((uint8,bool,uint16[]))
func (_Validator *ValidatorCaller) GetValidator(opts *bind.CallOpts, validator common.Address) (IDecimalMasterValidatorValidator, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "getValidator", validator)

	if err != nil {
		return *new(IDecimalMasterValidatorValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalMasterValidatorValidator)).(*IDecimalMasterValidatorValidator)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address validator) view returns((uint8,bool,uint16[]))
func (_Validator *ValidatorSession) GetValidator(validator common.Address) (IDecimalMasterValidatorValidator, error) {
	return _Validator.Contract.GetValidator(&_Validator.CallOpts, validator)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address validator) view returns((uint8,bool,uint16[]))
func (_Validator *ValidatorCallerSession) GetValidator(validator common.Address) (IDecimalMasterValidatorValidator, error) {
	return _Validator.Contract.GetValidator(&_Validator.CallOpts, validator)
}

// GetValidatorPenalties is a free data retrieval call binding the contract method 0x6dd6f5d1.
//
// Solidity: function getValidatorPenalties(address validator) view returns(uint16[])
func (_Validator *ValidatorCaller) GetValidatorPenalties(opts *bind.CallOpts, validator common.Address) ([]uint16, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "getValidatorPenalties", validator)

	if err != nil {
		return *new([]uint16), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint16)).(*[]uint16)

	return out0, err

}

// GetValidatorPenalties is a free data retrieval call binding the contract method 0x6dd6f5d1.
//
// Solidity: function getValidatorPenalties(address validator) view returns(uint16[])
func (_Validator *ValidatorSession) GetValidatorPenalties(validator common.Address) ([]uint16, error) {
	return _Validator.Contract.GetValidatorPenalties(&_Validator.CallOpts, validator)
}

// GetValidatorPenalties is a free data retrieval call binding the contract method 0x6dd6f5d1.
//
// Solidity: function getValidatorPenalties(address validator) view returns(uint16[])
func (_Validator *ValidatorCallerSession) GetValidatorPenalties(validator common.Address) ([]uint16, error) {
	return _Validator.Contract.GetValidatorPenalties(&_Validator.CallOpts, validator)
}

// GetValidatorStatus is a free data retrieval call binding the contract method 0xa310624f.
//
// Solidity: function getValidatorStatus(address validator) view returns(uint8)
func (_Validator *ValidatorCaller) GetValidatorStatus(opts *bind.CallOpts, validator common.Address) (uint8, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "getValidatorStatus", validator)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetValidatorStatus is a free data retrieval call binding the contract method 0xa310624f.
//
// Solidity: function getValidatorStatus(address validator) view returns(uint8)
func (_Validator *ValidatorSession) GetValidatorStatus(validator common.Address) (uint8, error) {
	return _Validator.Contract.GetValidatorStatus(&_Validator.CallOpts, validator)
}

// GetValidatorStatus is a free data retrieval call binding the contract method 0xa310624f.
//
// Solidity: function getValidatorStatus(address validator) view returns(uint8)
func (_Validator *ValidatorCallerSession) GetValidatorStatus(validator common.Address) (uint8, error) {
	return _Validator.Contract.GetValidatorStatus(&_Validator.CallOpts, validator)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address validator) view returns(bool)
func (_Validator *ValidatorCaller) IsActive(opts *bind.CallOpts, validator common.Address) (bool, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "isActive", validator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address validator) view returns(bool)
func (_Validator *ValidatorSession) IsActive(validator common.Address) (bool, error) {
	return _Validator.Contract.IsActive(&_Validator.CallOpts, validator)
}

// IsActive is a free data retrieval call binding the contract method 0x9f8a13d7.
//
// Solidity: function isActive(address validator) view returns(bool)
func (_Validator *ValidatorCallerSession) IsActive(validator common.Address) (bool, error) {
	return _Validator.Contract.IsActive(&_Validator.CallOpts, validator)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address validator) view returns(bool)
func (_Validator *ValidatorCaller) IsMember(opts *bind.CallOpts, validator common.Address) (bool, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "isMember", validator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address validator) view returns(bool)
func (_Validator *ValidatorSession) IsMember(validator common.Address) (bool, error) {
	return _Validator.Contract.IsMember(&_Validator.CallOpts, validator)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address validator) view returns(bool)
func (_Validator *ValidatorCallerSession) IsMember(validator common.Address) (bool, error) {
	return _Validator.Contract.IsMember(&_Validator.CallOpts, validator)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validator *ValidatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validator *ValidatorSession) Owner() (common.Address, error) {
	return _Validator.Contract.Owner(&_Validator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validator *ValidatorCallerSession) Owner() (common.Address, error) {
	return _Validator.Contract.Owner(&_Validator.CallOpts)
}

// AddCandidate is a paid mutator transaction binding the contract method 0xce6416e7.
//
// Solidity: function addCandidate(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorTransactor) AddCandidate(opts *bind.TransactOpts, validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addCandidate", validator, meta, initialStake)
}

// AddCandidate is a paid mutator transaction binding the contract method 0xce6416e7.
//
// Solidity: function addCandidate(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorSession) AddCandidate(validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.Contract.AddCandidate(&_Validator.TransactOpts, validator, meta, initialStake)
}

// AddCandidate is a paid mutator transaction binding the contract method 0xce6416e7.
//
// Solidity: function addCandidate(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorTransactorSession) AddCandidate(validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.Contract.AddCandidate(&_Validator.TransactOpts, validator, meta, initialStake)
}

// AddCandidateDEL is a paid mutator transaction binding the contract method 0xda4d56f0.
//
// Solidity: function addCandidateDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorTransactor) AddCandidateDEL(opts *bind.TransactOpts, validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addCandidateDEL", validator, meta)
}

// AddCandidateDEL is a paid mutator transaction binding the contract method 0xda4d56f0.
//
// Solidity: function addCandidateDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorSession) AddCandidateDEL(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.AddCandidateDEL(&_Validator.TransactOpts, validator, meta)
}

// AddCandidateDEL is a paid mutator transaction binding the contract method 0xda4d56f0.
//
// Solidity: function addCandidateDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorTransactorSession) AddCandidateDEL(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.AddCandidateDEL(&_Validator.TransactOpts, validator, meta)
}

// AddPenalty is a paid mutator transaction binding the contract method 0x0af66c5a.
//
// Solidity: function addPenalty(address validator, uint16 penaltyPercentage) returns()
func (_Validator *ValidatorTransactor) AddPenalty(opts *bind.TransactOpts, validator common.Address, penaltyPercentage uint16) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addPenalty", validator, penaltyPercentage)
}

// AddPenalty is a paid mutator transaction binding the contract method 0x0af66c5a.
//
// Solidity: function addPenalty(address validator, uint16 penaltyPercentage) returns()
func (_Validator *ValidatorSession) AddPenalty(validator common.Address, penaltyPercentage uint16) (*types.Transaction, error) {
	return _Validator.Contract.AddPenalty(&_Validator.TransactOpts, validator, penaltyPercentage)
}

// AddPenalty is a paid mutator transaction binding the contract method 0x0af66c5a.
//
// Solidity: function addPenalty(address validator, uint16 penaltyPercentage) returns()
func (_Validator *ValidatorTransactorSession) AddPenalty(validator common.Address, penaltyPercentage uint16) (*types.Transaction, error) {
	return _Validator.Contract.AddPenalty(&_Validator.TransactOpts, validator, penaltyPercentage)
}

// AddValidator is a paid mutator transaction binding the contract method 0x8e7d0f53.
//
// Solidity: function addValidator(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorTransactor) AddValidator(opts *bind.TransactOpts, validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addValidator", validator, meta, initialStake)
}

// AddValidator is a paid mutator transaction binding the contract method 0x8e7d0f53.
//
// Solidity: function addValidator(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorSession) AddValidator(validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.Contract.AddValidator(&_Validator.TransactOpts, validator, meta, initialStake)
}

// AddValidator is a paid mutator transaction binding the contract method 0x8e7d0f53.
//
// Solidity: function addValidator(address validator, string meta, (address,uint256) initialStake) returns()
func (_Validator *ValidatorTransactorSession) AddValidator(validator common.Address, meta string, initialStake IDecimalMasterValidatorStake) (*types.Transaction, error) {
	return _Validator.Contract.AddValidator(&_Validator.TransactOpts, validator, meta, initialStake)
}

// AddValidatorDEL is a paid mutator transaction binding the contract method 0x91234204.
//
// Solidity: function addValidatorDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorTransactor) AddValidatorDEL(opts *bind.TransactOpts, validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addValidatorDEL", validator, meta)
}

// AddValidatorDEL is a paid mutator transaction binding the contract method 0x91234204.
//
// Solidity: function addValidatorDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorSession) AddValidatorDEL(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.AddValidatorDEL(&_Validator.TransactOpts, validator, meta)
}

// AddValidatorDEL is a paid mutator transaction binding the contract method 0x91234204.
//
// Solidity: function addValidatorDEL(address validator, string meta) payable returns()
func (_Validator *ValidatorTransactorSession) AddValidatorDEL(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.AddValidatorDEL(&_Validator.TransactOpts, validator, meta)
}

// AddValidators is a paid mutator transaction binding the contract method 0xcdba1d92.
//
// Solidity: function addValidators((address,string)[] metas) returns()
func (_Validator *ValidatorTransactor) AddValidators(opts *bind.TransactOpts, metas []IDecimalMasterValidatorValidatorMeta) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "addValidators", metas)
}

// AddValidators is a paid mutator transaction binding the contract method 0xcdba1d92.
//
// Solidity: function addValidators((address,string)[] metas) returns()
func (_Validator *ValidatorSession) AddValidators(metas []IDecimalMasterValidatorValidatorMeta) (*types.Transaction, error) {
	return _Validator.Contract.AddValidators(&_Validator.TransactOpts, metas)
}

// AddValidators is a paid mutator transaction binding the contract method 0xcdba1d92.
//
// Solidity: function addValidators((address,string)[] metas) returns()
func (_Validator *ValidatorTransactorSession) AddValidators(metas []IDecimalMasterValidatorValidatorMeta) (*types.Transaction, error) {
	return _Validator.Contract.AddValidators(&_Validator.TransactOpts, metas)
}

// ApproveCandidates is a paid mutator transaction binding the contract method 0x8c9ad016.
//
// Solidity: function approveCandidates(address[] candidates) returns()
func (_Validator *ValidatorTransactor) ApproveCandidates(opts *bind.TransactOpts, candidates []common.Address) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "approveCandidates", candidates)
}

// ApproveCandidates is a paid mutator transaction binding the contract method 0x8c9ad016.
//
// Solidity: function approveCandidates(address[] candidates) returns()
func (_Validator *ValidatorSession) ApproveCandidates(candidates []common.Address) (*types.Transaction, error) {
	return _Validator.Contract.ApproveCandidates(&_Validator.TransactOpts, candidates)
}

// ApproveCandidates is a paid mutator transaction binding the contract method 0x8c9ad016.
//
// Solidity: function approveCandidates(address[] candidates) returns()
func (_Validator *ValidatorTransactorSession) ApproveCandidates(candidates []common.Address) (*types.Transaction, error) {
	return _Validator.Contract.ApproveCandidates(&_Validator.TransactOpts, candidates)
}

// PauseSelf is a paid mutator transaction binding the contract method 0xaf69ca76.
//
// Solidity: function pauseSelf() returns()
func (_Validator *ValidatorTransactor) PauseSelf(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "pauseSelf")
}

// PauseSelf is a paid mutator transaction binding the contract method 0xaf69ca76.
//
// Solidity: function pauseSelf() returns()
func (_Validator *ValidatorSession) PauseSelf() (*types.Transaction, error) {
	return _Validator.Contract.PauseSelf(&_Validator.TransactOpts)
}

// PauseSelf is a paid mutator transaction binding the contract method 0xaf69ca76.
//
// Solidity: function pauseSelf() returns()
func (_Validator *ValidatorTransactorSession) PauseSelf() (*types.Transaction, error) {
	return _Validator.Contract.PauseSelf(&_Validator.TransactOpts)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address validator) returns()
func (_Validator *ValidatorTransactor) PauseValidator(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "pauseValidator", validator)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address validator) returns()
func (_Validator *ValidatorSession) PauseValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.PauseValidator(&_Validator.TransactOpts, validator)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address validator) returns()
func (_Validator *ValidatorTransactorSession) PauseValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.PauseValidator(&_Validator.TransactOpts, validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_Validator *ValidatorTransactor) RemoveValidator(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "removeValidator", validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_Validator *ValidatorSession) RemoveValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.RemoveValidator(&_Validator.TransactOpts, validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address validator) returns()
func (_Validator *ValidatorTransactorSession) RemoveValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.RemoveValidator(&_Validator.TransactOpts, validator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validator *ValidatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validator *ValidatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validator.Contract.RenounceOwnership(&_Validator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validator *ValidatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validator.Contract.RenounceOwnership(&_Validator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validator *ValidatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validator *ValidatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validator.Contract.TransferOwnership(&_Validator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validator *ValidatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validator.Contract.TransferOwnership(&_Validator.TransactOpts, newOwner)
}

// UnpauseSelf is a paid mutator transaction binding the contract method 0x9391d5e8.
//
// Solidity: function unpauseSelf() returns()
func (_Validator *ValidatorTransactor) UnpauseSelf(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "unpauseSelf")
}

// UnpauseSelf is a paid mutator transaction binding the contract method 0x9391d5e8.
//
// Solidity: function unpauseSelf() returns()
func (_Validator *ValidatorSession) UnpauseSelf() (*types.Transaction, error) {
	return _Validator.Contract.UnpauseSelf(&_Validator.TransactOpts)
}

// UnpauseSelf is a paid mutator transaction binding the contract method 0x9391d5e8.
//
// Solidity: function unpauseSelf() returns()
func (_Validator *ValidatorTransactorSession) UnpauseSelf() (*types.Transaction, error) {
	return _Validator.Contract.UnpauseSelf(&_Validator.TransactOpts)
}

// UnpauseValidator is a paid mutator transaction binding the contract method 0x0437e4fd.
//
// Solidity: function unpauseValidator(address validator) returns()
func (_Validator *ValidatorTransactor) UnpauseValidator(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "unpauseValidator", validator)
}

// UnpauseValidator is a paid mutator transaction binding the contract method 0x0437e4fd.
//
// Solidity: function unpauseValidator(address validator) returns()
func (_Validator *ValidatorSession) UnpauseValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.UnpauseValidator(&_Validator.TransactOpts, validator)
}

// UnpauseValidator is a paid mutator transaction binding the contract method 0x0437e4fd.
//
// Solidity: function unpauseValidator(address validator) returns()
func (_Validator *ValidatorTransactorSession) UnpauseValidator(validator common.Address) (*types.Transaction, error) {
	return _Validator.Contract.UnpauseValidator(&_Validator.TransactOpts, validator)
}

// UpdateValidatorMeta is a paid mutator transaction binding the contract method 0xf2a19f61.
//
// Solidity: function updateValidatorMeta(address validator, string meta) returns()
func (_Validator *ValidatorTransactor) UpdateValidatorMeta(opts *bind.TransactOpts, validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "updateValidatorMeta", validator, meta)
}

// UpdateValidatorMeta is a paid mutator transaction binding the contract method 0xf2a19f61.
//
// Solidity: function updateValidatorMeta(address validator, string meta) returns()
func (_Validator *ValidatorSession) UpdateValidatorMeta(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.UpdateValidatorMeta(&_Validator.TransactOpts, validator, meta)
}

// UpdateValidatorMeta is a paid mutator transaction binding the contract method 0xf2a19f61.
//
// Solidity: function updateValidatorMeta(address validator, string meta) returns()
func (_Validator *ValidatorTransactorSession) UpdateValidatorMeta(validator common.Address, meta string) (*types.Transaction, error) {
	return _Validator.Contract.UpdateValidatorMeta(&_Validator.TransactOpts, validator, meta)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Validator *ValidatorTransactor) Upgrade(opts *bind.TransactOpts, newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Validator.contract.Transact(opts, "upgrade", newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Validator *ValidatorSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Validator.Contract.Upgrade(&_Validator.TransactOpts, newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Validator *ValidatorTransactorSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Validator.Contract.Upgrade(&_Validator.TransactOpts, newImpl, data)
}

// ValidatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validator contract.
type ValidatorInitializedIterator struct {
	Event *ValidatorInitialized // Event containing the contract specifics and raw log

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
func (it *ValidatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorInitialized)
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
		it.Event = new(ValidatorInitialized)
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
func (it *ValidatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorInitialized represents a Initialized event raised by the Validator contract.
type ValidatorInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validator *ValidatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*ValidatorInitializedIterator, error) {

	logs, sub, err := _Validator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ValidatorInitializedIterator{contract: _Validator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validator *ValidatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ValidatorInitialized) (event.Subscription, error) {

	logs, sub, err := _Validator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorInitialized)
				if err := _Validator.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validator *ValidatorFilterer) ParseInitialized(log types.Log) (*ValidatorInitialized, error) {
	event := new(ValidatorInitialized)
	if err := _Validator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validator contract.
type ValidatorOwnershipTransferredIterator struct {
	Event *ValidatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorOwnershipTransferred)
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
		it.Event = new(ValidatorOwnershipTransferred)
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
func (it *ValidatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorOwnershipTransferred represents a OwnershipTransferred event raised by the Validator contract.
type ValidatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validator *ValidatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorOwnershipTransferredIterator{contract: _Validator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validator *ValidatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorOwnershipTransferred)
				if err := _Validator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validator *ValidatorFilterer) ParseOwnershipTransferred(log types.Log) (*ValidatorOwnershipTransferred, error) {
	event := new(ValidatorOwnershipTransferred)
	if err := _Validator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Validator contract.
type ValidatorUpgradedIterator struct {
	Event *ValidatorUpgraded // Event containing the contract specifics and raw log

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
func (it *ValidatorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorUpgraded)
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
		it.Event = new(ValidatorUpgraded)
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
func (it *ValidatorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorUpgraded represents a Upgraded event raised by the Validator contract.
type ValidatorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validator *ValidatorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ValidatorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validator.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorUpgradedIterator{contract: _Validator.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validator *ValidatorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ValidatorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validator.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorUpgraded)
				if err := _Validator.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Validator *ValidatorFilterer) ParseUpgraded(log types.Log) (*ValidatorUpgraded, error) {
	event := new(ValidatorUpgraded)
	if err := _Validator.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorValidatorMetaUpdatedIterator is returned from FilterValidatorMetaUpdated and is used to iterate over the raw logs and unpacked data for ValidatorMetaUpdated events raised by the Validator contract.
type ValidatorValidatorMetaUpdatedIterator struct {
	Event *ValidatorValidatorMetaUpdated // Event containing the contract specifics and raw log

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
func (it *ValidatorValidatorMetaUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorValidatorMetaUpdated)
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
		it.Event = new(ValidatorValidatorMetaUpdated)
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
func (it *ValidatorValidatorMetaUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorValidatorMetaUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorValidatorMetaUpdated represents a ValidatorMetaUpdated event raised by the Validator contract.
type ValidatorValidatorMetaUpdated struct {
	Validator common.Address
	Meta      string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorMetaUpdated is a free log retrieval operation binding the contract event 0x8ca43af3725e78a8c75f7147c291edb09a28164c8917ee3029afdde5f55826c4.
//
// Solidity: event ValidatorMetaUpdated(address indexed validator, string meta)
func (_Validator *ValidatorFilterer) FilterValidatorMetaUpdated(opts *bind.FilterOpts, validator []common.Address) (*ValidatorValidatorMetaUpdatedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.FilterLogs(opts, "ValidatorMetaUpdated", validatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorValidatorMetaUpdatedIterator{contract: _Validator.contract, event: "ValidatorMetaUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorMetaUpdated is a free log subscription operation binding the contract event 0x8ca43af3725e78a8c75f7147c291edb09a28164c8917ee3029afdde5f55826c4.
//
// Solidity: event ValidatorMetaUpdated(address indexed validator, string meta)
func (_Validator *ValidatorFilterer) WatchValidatorMetaUpdated(opts *bind.WatchOpts, sink chan<- *ValidatorValidatorMetaUpdated, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.WatchLogs(opts, "ValidatorMetaUpdated", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorValidatorMetaUpdated)
				if err := _Validator.contract.UnpackLog(event, "ValidatorMetaUpdated", log); err != nil {
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

// ParseValidatorMetaUpdated is a log parse operation binding the contract event 0x8ca43af3725e78a8c75f7147c291edb09a28164c8917ee3029afdde5f55826c4.
//
// Solidity: event ValidatorMetaUpdated(address indexed validator, string meta)
func (_Validator *ValidatorFilterer) ParseValidatorMetaUpdated(log types.Log) (*ValidatorValidatorMetaUpdated, error) {
	event := new(ValidatorValidatorMetaUpdated)
	if err := _Validator.contract.UnpackLog(event, "ValidatorMetaUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorValidatorPenaltyIterator is returned from FilterValidatorPenalty and is used to iterate over the raw logs and unpacked data for ValidatorPenalty events raised by the Validator contract.
type ValidatorValidatorPenaltyIterator struct {
	Event *ValidatorValidatorPenalty // Event containing the contract specifics and raw log

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
func (it *ValidatorValidatorPenaltyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorValidatorPenalty)
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
		it.Event = new(ValidatorValidatorPenalty)
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
func (it *ValidatorValidatorPenaltyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorValidatorPenaltyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorValidatorPenalty represents a ValidatorPenalty event raised by the Validator contract.
type ValidatorValidatorPenalty struct {
	Validator  common.Address
	Percentage *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValidatorPenalty is a free log retrieval operation binding the contract event 0x19ff0fc7a4a07f7c07236f0b2fe800e25ec7beae3194fe7a2bcb823057b2b145.
//
// Solidity: event ValidatorPenalty(address indexed validator, uint256 percentage)
func (_Validator *ValidatorFilterer) FilterValidatorPenalty(opts *bind.FilterOpts, validator []common.Address) (*ValidatorValidatorPenaltyIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.FilterLogs(opts, "ValidatorPenalty", validatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorValidatorPenaltyIterator{contract: _Validator.contract, event: "ValidatorPenalty", logs: logs, sub: sub}, nil
}

// WatchValidatorPenalty is a free log subscription operation binding the contract event 0x19ff0fc7a4a07f7c07236f0b2fe800e25ec7beae3194fe7a2bcb823057b2b145.
//
// Solidity: event ValidatorPenalty(address indexed validator, uint256 percentage)
func (_Validator *ValidatorFilterer) WatchValidatorPenalty(opts *bind.WatchOpts, sink chan<- *ValidatorValidatorPenalty, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.WatchLogs(opts, "ValidatorPenalty", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorValidatorPenalty)
				if err := _Validator.contract.UnpackLog(event, "ValidatorPenalty", log); err != nil {
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

// ParseValidatorPenalty is a log parse operation binding the contract event 0x19ff0fc7a4a07f7c07236f0b2fe800e25ec7beae3194fe7a2bcb823057b2b145.
//
// Solidity: event ValidatorPenalty(address indexed validator, uint256 percentage)
func (_Validator *ValidatorFilterer) ParseValidatorPenalty(log types.Log) (*ValidatorValidatorPenalty, error) {
	event := new(ValidatorValidatorPenalty)
	if err := _Validator.contract.UnpackLog(event, "ValidatorPenalty", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorValidatorUpdatedIterator is returned from FilterValidatorUpdated and is used to iterate over the raw logs and unpacked data for ValidatorUpdated events raised by the Validator contract.
type ValidatorValidatorUpdatedIterator struct {
	Event *ValidatorValidatorUpdated // Event containing the contract specifics and raw log

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
func (it *ValidatorValidatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorValidatorUpdated)
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
		it.Event = new(ValidatorValidatorUpdated)
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
func (it *ValidatorValidatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorValidatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorValidatorUpdated represents a ValidatorUpdated event raised by the Validator contract.
type ValidatorValidatorUpdated struct {
	Validator common.Address
	Status    uint8
	Paused    bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorUpdated is a free log retrieval operation binding the contract event 0x522fda59479e70d03bc6366f2541872b0144f9099ef5d2260c2faf0360979a1d.
//
// Solidity: event ValidatorUpdated(address indexed validator, uint8 status, bool paused)
func (_Validator *ValidatorFilterer) FilterValidatorUpdated(opts *bind.FilterOpts, validator []common.Address) (*ValidatorValidatorUpdatedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.FilterLogs(opts, "ValidatorUpdated", validatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorValidatorUpdatedIterator{contract: _Validator.contract, event: "ValidatorUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorUpdated is a free log subscription operation binding the contract event 0x522fda59479e70d03bc6366f2541872b0144f9099ef5d2260c2faf0360979a1d.
//
// Solidity: event ValidatorUpdated(address indexed validator, uint8 status, bool paused)
func (_Validator *ValidatorFilterer) WatchValidatorUpdated(opts *bind.WatchOpts, sink chan<- *ValidatorValidatorUpdated, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Validator.contract.WatchLogs(opts, "ValidatorUpdated", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorValidatorUpdated)
				if err := _Validator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
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

// ParseValidatorUpdated is a log parse operation binding the contract event 0x522fda59479e70d03bc6366f2541872b0144f9099ef5d2260c2faf0360979a1d.
//
// Solidity: event ValidatorUpdated(address indexed validator, uint8 status, bool paused)
func (_Validator *ValidatorFilterer) ParseValidatorUpdated(log types.Log) (*ValidatorValidatorUpdated, error) {
	event := new(ValidatorValidatorUpdated)
	if err := _Validator.contract.UnpackLog(event, "ValidatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
