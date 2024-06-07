// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nftCenter

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

// DecimalNFTCenterNFT is an auto generated low-level Go binding around an user-defined struct.
type DecimalNFTCenterNFT struct {
	TokenOwner common.Address
	Symbol     string
	Name       string
	Refundable bool
}

// NFTState is an auto generated low-level Go binding around an user-defined struct.
type NFTState struct {
	Active  bool
	NftType uint8
}

// NftCenterMetaData contains all meta data concerning the NftCenter contract.
var NftCenterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDeployed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAllowMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayloadLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"ContractUpgradedNFT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"nft\",\"type\":\"tuple\"}],\"name\":\"NFTCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"beacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"checkToken\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createERC1155\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createERC721\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"}],\"name\":\"getNftState\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"internalType\":\"structNFTState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressContractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newNFTImplementation\",\"type\":\"address\"},{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"upgradeNFT\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// NftCenterABI is the input ABI used to generate the binding from.
// Deprecated: Use NftCenterMetaData.ABI instead.
var NftCenterABI = NftCenterMetaData.ABI

// NftCenter is an auto generated Go binding around an Ethereum contract.
type NftCenter struct {
	NftCenterCaller     // Read-only binding to the contract
	NftCenterTransactor // Write-only binding to the contract
	NftCenterFilterer   // Log filterer for contract events
}

// NftCenterCaller is an auto generated read-only Go binding around an Ethereum contract.
type NftCenterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftCenterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NftCenterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftCenterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NftCenterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NftCenterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NftCenterSession struct {
	Contract     *NftCenter        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NftCenterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NftCenterCallerSession struct {
	Contract *NftCenterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// NftCenterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NftCenterTransactorSession struct {
	Contract     *NftCenterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// NftCenterRaw is an auto generated low-level Go binding around an Ethereum contract.
type NftCenterRaw struct {
	Contract *NftCenter // Generic contract binding to access the raw methods on
}

// NftCenterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NftCenterCallerRaw struct {
	Contract *NftCenterCaller // Generic read-only contract binding to access the raw methods on
}

// NftCenterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NftCenterTransactorRaw struct {
	Contract *NftCenterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNftCenter creates a new instance of NftCenter, bound to a specific deployed contract.
func NewNftCenter(address common.Address, backend bind.ContractBackend) (*NftCenter, error) {
	contract, err := bindNftCenter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NftCenter{NftCenterCaller: NftCenterCaller{contract: contract}, NftCenterTransactor: NftCenterTransactor{contract: contract}, NftCenterFilterer: NftCenterFilterer{contract: contract}}, nil
}

// NewNftCenterCaller creates a new read-only instance of NftCenter, bound to a specific deployed contract.
func NewNftCenterCaller(address common.Address, caller bind.ContractCaller) (*NftCenterCaller, error) {
	contract, err := bindNftCenter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NftCenterCaller{contract: contract}, nil
}

// NewNftCenterTransactor creates a new write-only instance of NftCenter, bound to a specific deployed contract.
func NewNftCenterTransactor(address common.Address, transactor bind.ContractTransactor) (*NftCenterTransactor, error) {
	contract, err := bindNftCenter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NftCenterTransactor{contract: contract}, nil
}

// NewNftCenterFilterer creates a new log filterer instance of NftCenter, bound to a specific deployed contract.
func NewNftCenterFilterer(address common.Address, filterer bind.ContractFilterer) (*NftCenterFilterer, error) {
	contract, err := bindNftCenter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NftCenterFilterer{contract: contract}, nil
}

// bindNftCenter binds a generic wrapper to an already deployed contract.
func bindNftCenter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NftCenterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NftCenter *NftCenterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NftCenter.Contract.NftCenterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NftCenter *NftCenterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftCenter.Contract.NftCenterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NftCenter *NftCenterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NftCenter.Contract.NftCenterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NftCenter *NftCenterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NftCenter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NftCenter *NftCenterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftCenter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NftCenter *NftCenterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NftCenter.Contract.contract.Transact(opts, method, params...)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_NftCenter *NftCenterCaller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_NftCenter *NftCenterSession) MINRESERVE() (*big.Int, error) {
	return _NftCenter.Contract.MINRESERVE(&_NftCenter.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_NftCenter *NftCenterCallerSession) MINRESERVE() (*big.Int, error) {
	return _NftCenter.Contract.MINRESERVE(&_NftCenter.CallOpts)
}

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterCaller) Beacon(opts *bind.CallOpts, nftType uint8) (common.Address, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "beacon", nftType)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterSession) Beacon(nftType uint8) (common.Address, error) {
	return _NftCenter.Contract.Beacon(&_NftCenter.CallOpts, nftType)
}

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterCallerSession) Beacon(nftType uint8) (common.Address, error) {
	return _NftCenter.Contract.Beacon(&_NftCenter.CallOpts, nftType)
}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_NftCenter *NftCenterCaller) CheckToken(opts *bind.CallOpts, token common.Address) error {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "checkToken", token)

	if err != nil {
		return err
	}

	return err

}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_NftCenter *NftCenterSession) CheckToken(token common.Address) error {
	return _NftCenter.Contract.CheckToken(&_NftCenter.CallOpts, token)
}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_NftCenter *NftCenterCallerSession) CheckToken(token common.Address) error {
	return _NftCenter.Contract.CheckToken(&_NftCenter.CallOpts, token)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_NftCenter *NftCenterCaller) GetContractCenter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "getContractCenter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_NftCenter *NftCenterSession) GetContractCenter() (common.Address, error) {
	return _NftCenter.Contract.GetContractCenter(&_NftCenter.CallOpts)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_NftCenter *NftCenterCallerSession) GetContractCenter() (common.Address, error) {
	return _NftCenter.Contract.GetContractCenter(&_NftCenter.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_NftCenter *NftCenterCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_NftCenter *NftCenterSession) GetImplementation() (common.Address, error) {
	return _NftCenter.Contract.GetImplementation(&_NftCenter.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_NftCenter *NftCenterCallerSession) GetImplementation() (common.Address, error) {
	return _NftCenter.Contract.GetImplementation(&_NftCenter.CallOpts)
}

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_NftCenter *NftCenterCaller) GetNftState(opts *bind.CallOpts, nft common.Address) (NFTState, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "getNftState", nft)

	if err != nil {
		return *new(NFTState), err
	}

	out0 := *abi.ConvertType(out[0], new(NFTState)).(*NFTState)

	return out0, err

}

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_NftCenter *NftCenterSession) GetNftState(nft common.Address) (NFTState, error) {
	return _NftCenter.Contract.GetNftState(&_NftCenter.CallOpts, nft)
}

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_NftCenter *NftCenterCallerSession) GetNftState(nft common.Address) (NFTState, error) {
	return _NftCenter.Contract.GetNftState(&_NftCenter.CallOpts, nft)
}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterCaller) Implementation(opts *bind.CallOpts, nftType uint8) (common.Address, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "implementation", nftType)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterSession) Implementation(nftType uint8) (common.Address, error) {
	return _NftCenter.Contract.Implementation(&_NftCenter.CallOpts, nftType)
}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_NftCenter *NftCenterCallerSession) Implementation(nftType uint8) (common.Address, error) {
	return _NftCenter.Contract.Implementation(&_NftCenter.CallOpts, nftType)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NftCenter *NftCenterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NftCenter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NftCenter *NftCenterSession) Owner() (common.Address, error) {
	return _NftCenter.Contract.Owner(&_NftCenter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NftCenter *NftCenterCallerSession) Owner() (common.Address, error) {
	return _NftCenter.Contract.Owner(&_NftCenter.CallOpts)
}

// CreateERC1155 is a paid mutator transaction binding the contract method 0x969a269d.
//
// Solidity: function createERC1155((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterTransactor) CreateERC1155(opts *bind.TransactOpts, meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "createERC1155", meta)
}

// CreateERC1155 is a paid mutator transaction binding the contract method 0x969a269d.
//
// Solidity: function createERC1155((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterSession) CreateERC1155(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.Contract.CreateERC1155(&_NftCenter.TransactOpts, meta)
}

// CreateERC1155 is a paid mutator transaction binding the contract method 0x969a269d.
//
// Solidity: function createERC1155((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterTransactorSession) CreateERC1155(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.Contract.CreateERC1155(&_NftCenter.TransactOpts, meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x976d659c.
//
// Solidity: function createERC721((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterTransactor) CreateERC721(opts *bind.TransactOpts, meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "createERC721", meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x976d659c.
//
// Solidity: function createERC721((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterSession) CreateERC721(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.Contract.CreateERC721(&_NftCenter.TransactOpts, meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x976d659c.
//
// Solidity: function createERC721((address,string,string,bool) meta) returns(address)
func (_NftCenter *NftCenterTransactorSession) CreateERC721(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _NftCenter.Contract.CreateERC721(&_NftCenter.TransactOpts, meta)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NftCenter *NftCenterTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NftCenter *NftCenterSession) Initialize() (*types.Transaction, error) {
	return _NftCenter.Contract.Initialize(&_NftCenter.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_NftCenter *NftCenterTransactorSession) Initialize() (*types.Transaction, error) {
	return _NftCenter.Contract.Initialize(&_NftCenter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NftCenter *NftCenterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NftCenter *NftCenterSession) RenounceOwnership() (*types.Transaction, error) {
	return _NftCenter.Contract.RenounceOwnership(&_NftCenter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NftCenter *NftCenterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NftCenter.Contract.RenounceOwnership(&_NftCenter.TransactOpts)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_NftCenter *NftCenterTransactor) SetContractCenter(opts *bind.TransactOpts, addressContractCenter common.Address) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "setContractCenter", addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_NftCenter *NftCenterSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _NftCenter.Contract.SetContractCenter(&_NftCenter.TransactOpts, addressContractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address addressContractCenter) returns()
func (_NftCenter *NftCenterTransactorSession) SetContractCenter(addressContractCenter common.Address) (*types.Transaction, error) {
	return _NftCenter.Contract.SetContractCenter(&_NftCenter.TransactOpts, addressContractCenter)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NftCenter *NftCenterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NftCenter *NftCenterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NftCenter.Contract.TransferOwnership(&_NftCenter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NftCenter *NftCenterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NftCenter.Contract.TransferOwnership(&_NftCenter.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_NftCenter *NftCenterTransactor) Upgrade(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "upgrade", newImplementation, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_NftCenter *NftCenterSession) Upgrade(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NftCenter.Contract.Upgrade(&_NftCenter.TransactOpts, newImplementation, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImplementation, bytes data) returns()
func (_NftCenter *NftCenterTransactorSession) Upgrade(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NftCenter.Contract.Upgrade(&_NftCenter.TransactOpts, newImplementation, data)
}

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_NftCenter *NftCenterTransactor) UpgradeNFT(opts *bind.TransactOpts, newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _NftCenter.contract.Transact(opts, "upgradeNFT", newNFTImplementation, nftType)
}

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_NftCenter *NftCenterSession) UpgradeNFT(newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _NftCenter.Contract.UpgradeNFT(&_NftCenter.TransactOpts, newNFTImplementation, nftType)
}

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_NftCenter *NftCenterTransactorSession) UpgradeNFT(newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _NftCenter.Contract.UpgradeNFT(&_NftCenter.TransactOpts, newNFTImplementation, nftType)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NftCenter *NftCenterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NftCenter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NftCenter *NftCenterSession) Receive() (*types.Transaction, error) {
	return _NftCenter.Contract.Receive(&_NftCenter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_NftCenter *NftCenterTransactorSession) Receive() (*types.Transaction, error) {
	return _NftCenter.Contract.Receive(&_NftCenter.TransactOpts)
}

// NftCenterContractUpgradedIterator is returned from FilterContractUpgraded and is used to iterate over the raw logs and unpacked data for ContractUpgraded events raised by the NftCenter contract.
type NftCenterContractUpgradedIterator struct {
	Event *NftCenterContractUpgraded // Event containing the contract specifics and raw log

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
func (it *NftCenterContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterContractUpgraded)
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
		it.Event = new(NftCenterContractUpgraded)
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
func (it *NftCenterContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterContractUpgraded represents a ContractUpgraded event raised by the NftCenter contract.
type NftCenterContractUpgraded struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgraded is a free log retrieval operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_NftCenter *NftCenterFilterer) FilterContractUpgraded(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address) (*NftCenterContractUpgradedIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return &NftCenterContractUpgradedIterator{contract: _NftCenter.contract, event: "ContractUpgraded", logs: logs, sub: sub}, nil
}

// WatchContractUpgraded is a free log subscription operation binding the contract event 0x2e4cc16c100f0b55e2df82ab0b1a7e294aa9cbd01b48fbaf622683fbc0507a49.
//
// Solidity: event ContractUpgraded(address indexed oldContract, address indexed newContract)
func (_NftCenter *NftCenterFilterer) WatchContractUpgraded(opts *bind.WatchOpts, sink chan<- *NftCenterContractUpgraded, oldContract []common.Address, newContract []common.Address) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "ContractUpgraded", oldContractRule, newContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterContractUpgraded)
				if err := _NftCenter.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
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
func (_NftCenter *NftCenterFilterer) ParseContractUpgraded(log types.Log) (*NftCenterContractUpgraded, error) {
	event := new(NftCenterContractUpgraded)
	if err := _NftCenter.contract.UnpackLog(event, "ContractUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftCenterContractUpgradedNFTIterator is returned from FilterContractUpgradedNFT and is used to iterate over the raw logs and unpacked data for ContractUpgradedNFT events raised by the NftCenter contract.
type NftCenterContractUpgradedNFTIterator struct {
	Event *NftCenterContractUpgradedNFT // Event containing the contract specifics and raw log

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
func (it *NftCenterContractUpgradedNFTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterContractUpgradedNFT)
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
		it.Event = new(NftCenterContractUpgradedNFT)
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
func (it *NftCenterContractUpgradedNFTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterContractUpgradedNFTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterContractUpgradedNFT represents a ContractUpgradedNFT event raised by the NftCenter contract.
type NftCenterContractUpgradedNFT struct {
	OldContract common.Address
	NewContract common.Address
	NftType     uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgradedNFT is a free log retrieval operation binding the contract event 0xe4dfc1498d5921eb1933069a6657c79e0b3fa55db46bf016d331de23cee72a3a.
//
// Solidity: event ContractUpgradedNFT(address indexed oldContract, address indexed newContract, uint8 indexed nftType)
func (_NftCenter *NftCenterFilterer) FilterContractUpgradedNFT(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address, nftType []uint8) (*NftCenterContractUpgradedNFTIterator, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}
	var nftTypeRule []interface{}
	for _, nftTypeItem := range nftType {
		nftTypeRule = append(nftTypeRule, nftTypeItem)
	}

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "ContractUpgradedNFT", oldContractRule, newContractRule, nftTypeRule)
	if err != nil {
		return nil, err
	}
	return &NftCenterContractUpgradedNFTIterator{contract: _NftCenter.contract, event: "ContractUpgradedNFT", logs: logs, sub: sub}, nil
}

// WatchContractUpgradedNFT is a free log subscription operation binding the contract event 0xe4dfc1498d5921eb1933069a6657c79e0b3fa55db46bf016d331de23cee72a3a.
//
// Solidity: event ContractUpgradedNFT(address indexed oldContract, address indexed newContract, uint8 indexed nftType)
func (_NftCenter *NftCenterFilterer) WatchContractUpgradedNFT(opts *bind.WatchOpts, sink chan<- *NftCenterContractUpgradedNFT, oldContract []common.Address, newContract []common.Address, nftType []uint8) (event.Subscription, error) {

	var oldContractRule []interface{}
	for _, oldContractItem := range oldContract {
		oldContractRule = append(oldContractRule, oldContractItem)
	}
	var newContractRule []interface{}
	for _, newContractItem := range newContract {
		newContractRule = append(newContractRule, newContractItem)
	}
	var nftTypeRule []interface{}
	for _, nftTypeItem := range nftType {
		nftTypeRule = append(nftTypeRule, nftTypeItem)
	}

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "ContractUpgradedNFT", oldContractRule, newContractRule, nftTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterContractUpgradedNFT)
				if err := _NftCenter.contract.UnpackLog(event, "ContractUpgradedNFT", log); err != nil {
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

// ParseContractUpgradedNFT is a log parse operation binding the contract event 0xe4dfc1498d5921eb1933069a6657c79e0b3fa55db46bf016d331de23cee72a3a.
//
// Solidity: event ContractUpgradedNFT(address indexed oldContract, address indexed newContract, uint8 indexed nftType)
func (_NftCenter *NftCenterFilterer) ParseContractUpgradedNFT(log types.Log) (*NftCenterContractUpgradedNFT, error) {
	event := new(NftCenterContractUpgradedNFT)
	if err := _NftCenter.contract.UnpackLog(event, "ContractUpgradedNFT", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftCenterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NftCenter contract.
type NftCenterInitializedIterator struct {
	Event *NftCenterInitialized // Event containing the contract specifics and raw log

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
func (it *NftCenterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterInitialized)
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
		it.Event = new(NftCenterInitialized)
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
func (it *NftCenterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterInitialized represents a Initialized event raised by the NftCenter contract.
type NftCenterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NftCenter *NftCenterFilterer) FilterInitialized(opts *bind.FilterOpts) (*NftCenterInitializedIterator, error) {

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NftCenterInitializedIterator{contract: _NftCenter.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NftCenter *NftCenterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NftCenterInitialized) (event.Subscription, error) {

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterInitialized)
				if err := _NftCenter.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NftCenter *NftCenterFilterer) ParseInitialized(log types.Log) (*NftCenterInitialized, error) {
	event := new(NftCenterInitialized)
	if err := _NftCenter.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftCenterNFTCreatedIterator is returned from FilterNFTCreated and is used to iterate over the raw logs and unpacked data for NFTCreated events raised by the NftCenter contract.
type NftCenterNFTCreatedIterator struct {
	Event *NftCenterNFTCreated // Event containing the contract specifics and raw log

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
func (it *NftCenterNFTCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterNFTCreated)
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
		it.Event = new(NftCenterNFTCreated)
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
func (it *NftCenterNFTCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterNFTCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterNFTCreated represents a NFTCreated event raised by the NftCenter contract.
type NftCenterNFTCreated struct {
	TokenAddress common.Address
	NftType      uint8
	Nft          DecimalNFTCenterNFT
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNFTCreated is a free log retrieval operation binding the contract event 0x80e7e91eab84672429971f005c9d800477ffaaa3de26a1a85269fbdbc45dc465.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 nftType, (address,string,string,bool) nft)
func (_NftCenter *NftCenterFilterer) FilterNFTCreated(opts *bind.FilterOpts) (*NftCenterNFTCreatedIterator, error) {

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "NFTCreated")
	if err != nil {
		return nil, err
	}
	return &NftCenterNFTCreatedIterator{contract: _NftCenter.contract, event: "NFTCreated", logs: logs, sub: sub}, nil
}

// WatchNFTCreated is a free log subscription operation binding the contract event 0x80e7e91eab84672429971f005c9d800477ffaaa3de26a1a85269fbdbc45dc465.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 nftType, (address,string,string,bool) nft)
func (_NftCenter *NftCenterFilterer) WatchNFTCreated(opts *bind.WatchOpts, sink chan<- *NftCenterNFTCreated) (event.Subscription, error) {

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "NFTCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterNFTCreated)
				if err := _NftCenter.contract.UnpackLog(event, "NFTCreated", log); err != nil {
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

// ParseNFTCreated is a log parse operation binding the contract event 0x80e7e91eab84672429971f005c9d800477ffaaa3de26a1a85269fbdbc45dc465.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 nftType, (address,string,string,bool) nft)
func (_NftCenter *NftCenterFilterer) ParseNFTCreated(log types.Log) (*NftCenterNFTCreated, error) {
	event := new(NftCenterNFTCreated)
	if err := _NftCenter.contract.UnpackLog(event, "NFTCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftCenterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NftCenter contract.
type NftCenterOwnershipTransferredIterator struct {
	Event *NftCenterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NftCenterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterOwnershipTransferred)
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
		it.Event = new(NftCenterOwnershipTransferred)
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
func (it *NftCenterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterOwnershipTransferred represents a OwnershipTransferred event raised by the NftCenter contract.
type NftCenterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NftCenter *NftCenterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NftCenterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NftCenterOwnershipTransferredIterator{contract: _NftCenter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NftCenter *NftCenterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NftCenterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterOwnershipTransferred)
				if err := _NftCenter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NftCenter *NftCenterFilterer) ParseOwnershipTransferred(log types.Log) (*NftCenterOwnershipTransferred, error) {
	event := new(NftCenterOwnershipTransferred)
	if err := _NftCenter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NftCenterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NftCenter contract.
type NftCenterUpgradedIterator struct {
	Event *NftCenterUpgraded // Event containing the contract specifics and raw log

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
func (it *NftCenterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NftCenterUpgraded)
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
		it.Event = new(NftCenterUpgraded)
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
func (it *NftCenterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NftCenterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NftCenterUpgraded represents a Upgraded event raised by the NftCenter contract.
type NftCenterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NftCenter *NftCenterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*NftCenterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NftCenter.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &NftCenterUpgradedIterator{contract: _NftCenter.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NftCenter *NftCenterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NftCenterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NftCenter.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NftCenterUpgraded)
				if err := _NftCenter.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_NftCenter *NftCenterFilterer) ParseUpgraded(log types.Log) (*NftCenterUpgraded, error) {
	event := new(NftCenterUpgraded)
	if err := _NftCenter.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
