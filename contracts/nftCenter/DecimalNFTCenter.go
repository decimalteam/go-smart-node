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

// DecimalNFTCenterNFT is an auto generated low-level Go binding around an user-defined struct.
type DecimalNFTCenterNFT struct {
	Creator     common.Address
	Symbol      string
	Name        string
	ContractURI string
	Refundable  bool
}

// NFTState is an auto generated low-level Go binding around an user-defined struct.
type NFTState struct {
	Active  bool
	NftType uint8
}

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AlreadyDeployed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAllowMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialMint\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPayloadLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newContract\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"ContractUpgradedNFT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contractURI\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"nft\",\"type\":\"tuple\"}],\"name\":\"NFTCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"beacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"checkToken\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contractURI\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createERC1155\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"contractURI\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"refundable\",\"type\":\"bool\"}],\"internalType\":\"structDecimalNFTCenter.NFT\",\"name\":\"meta\",\"type\":\"tuple\"}],\"name\":\"createERC721\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"}],\"name\":\"getNftState\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"internalType\":\"structNFTState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressContractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newNFTImplementation\",\"type\":\"address\"},{\"internalType\":\"enumNFTType\",\"name\":\"nftType\",\"type\":\"uint8\"}],\"name\":\"upgradeNFT\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_Delegation *DelegationCaller) Beacon(opts *bind.CallOpts, nftType uint8) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "beacon", nftType)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_Delegation *DelegationSession) Beacon(nftType uint8) (common.Address, error) {
	return _Delegation.Contract.Beacon(&_Delegation.CallOpts, nftType)
}

// Beacon is a free data retrieval call binding the contract method 0x99f6c139.
//
// Solidity: function beacon(uint8 nftType) view returns(address)
func (_Delegation *DelegationCallerSession) Beacon(nftType uint8) (common.Address, error) {
	return _Delegation.Contract.Beacon(&_Delegation.CallOpts, nftType)
}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_Delegation *DelegationCaller) CheckToken(opts *bind.CallOpts, token common.Address) error {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "checkToken", token)

	if err != nil {
		return err
	}

	return err

}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_Delegation *DelegationSession) CheckToken(token common.Address) error {
	return _Delegation.Contract.CheckToken(&_Delegation.CallOpts, token)
}

// CheckToken is a free data retrieval call binding the contract method 0xf1880b24.
//
// Solidity: function checkToken(address token) view returns()
func (_Delegation *DelegationCallerSession) CheckToken(token common.Address) error {
	return _Delegation.Contract.CheckToken(&_Delegation.CallOpts, token)
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

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_Delegation *DelegationCaller) GetNftState(opts *bind.CallOpts, nft common.Address) (NFTState, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getNftState", nft)

	if err != nil {
		return *new(NFTState), err
	}

	out0 := *abi.ConvertType(out[0], new(NFTState)).(*NFTState)

	return out0, err

}

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_Delegation *DelegationSession) GetNftState(nft common.Address) (NFTState, error) {
	return _Delegation.Contract.GetNftState(&_Delegation.CallOpts, nft)
}

// GetNftState is a free data retrieval call binding the contract method 0xdd03ada8.
//
// Solidity: function getNftState(address nft) view returns((bool,uint8))
func (_Delegation *DelegationCallerSession) GetNftState(nft common.Address) (NFTState, error) {
	return _Delegation.Contract.GetNftState(&_Delegation.CallOpts, nft)
}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_Delegation *DelegationCaller) Implementation(opts *bind.CallOpts, nftType uint8) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "implementation", nftType)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_Delegation *DelegationSession) Implementation(nftType uint8) (common.Address, error) {
	return _Delegation.Contract.Implementation(&_Delegation.CallOpts, nftType)
}

// Implementation is a free data retrieval call binding the contract method 0xf19a74c5.
//
// Solidity: function implementation(uint8 nftType) view returns(address)
func (_Delegation *DelegationCallerSession) Implementation(nftType uint8) (common.Address, error) {
	return _Delegation.Contract.Implementation(&_Delegation.CallOpts, nftType)
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

// CreateERC1155 is a paid mutator transaction binding the contract method 0xdded1bd1.
//
// Solidity: function createERC1155((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationTransactor) CreateERC1155(opts *bind.TransactOpts, meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "createERC1155", meta)
}

// CreateERC1155 is a paid mutator transaction binding the contract method 0xdded1bd1.
//
// Solidity: function createERC1155((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationSession) CreateERC1155(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.Contract.CreateERC1155(&_Delegation.TransactOpts, meta)
}

// CreateERC1155 is a paid mutator transaction binding the contract method 0xdded1bd1.
//
// Solidity: function createERC1155((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationTransactorSession) CreateERC1155(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.Contract.CreateERC1155(&_Delegation.TransactOpts, meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x7ec4900d.
//
// Solidity: function createERC721((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationTransactor) CreateERC721(opts *bind.TransactOpts, meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "createERC721", meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x7ec4900d.
//
// Solidity: function createERC721((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationSession) CreateERC721(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.Contract.CreateERC721(&_Delegation.TransactOpts, meta)
}

// CreateERC721 is a paid mutator transaction binding the contract method 0x7ec4900d.
//
// Solidity: function createERC721((address,string,string,string,bool) meta) returns(address)
func (_Delegation *DelegationTransactorSession) CreateERC721(meta DecimalNFTCenterNFT) (*types.Transaction, error) {
	return _Delegation.Contract.CreateERC721(&_Delegation.TransactOpts, meta)
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

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_Delegation *DelegationTransactor) UpgradeNFT(opts *bind.TransactOpts, newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "upgradeNFT", newNFTImplementation, nftType)
}

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_Delegation *DelegationSession) UpgradeNFT(newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _Delegation.Contract.UpgradeNFT(&_Delegation.TransactOpts, newNFTImplementation, nftType)
}

// UpgradeNFT is a paid mutator transaction binding the contract method 0x97451e91.
//
// Solidity: function upgradeNFT(address newNFTImplementation, uint8 nftType) returns()
func (_Delegation *DelegationTransactorSession) UpgradeNFT(newNFTImplementation common.Address, nftType uint8) (*types.Transaction, error) {
	return _Delegation.Contract.UpgradeNFT(&_Delegation.TransactOpts, newNFTImplementation, nftType)
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

// DelegationContractUpgradedNFTIterator is returned from FilterContractUpgradedNFT and is used to iterate over the raw logs and unpacked data for ContractUpgradedNFT events raised by the Delegation contract.
type DelegationContractUpgradedNFTIterator struct {
	Event *DelegationContractUpgradedNFT // Event containing the contract specifics and raw log

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
func (it *DelegationContractUpgradedNFTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationContractUpgradedNFT)
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
		it.Event = new(DelegationContractUpgradedNFT)
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
func (it *DelegationContractUpgradedNFTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationContractUpgradedNFTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationContractUpgradedNFT represents a ContractUpgradedNFT event raised by the Delegation contract.
type DelegationContractUpgradedNFT struct {
	OldContract common.Address
	NewContract common.Address
	NftType     uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgradedNFT is a free log retrieval operation binding the contract event 0xe4dfc1498d5921eb1933069a6657c79e0b3fa55db46bf016d331de23cee72a3a.
//
// Solidity: event ContractUpgradedNFT(address indexed oldContract, address indexed newContract, uint8 indexed nftType)
func (_Delegation *DelegationFilterer) FilterContractUpgradedNFT(opts *bind.FilterOpts, oldContract []common.Address, newContract []common.Address, nftType []uint8) (*DelegationContractUpgradedNFTIterator, error) {

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

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ContractUpgradedNFT", oldContractRule, newContractRule, nftTypeRule)
	if err != nil {
		return nil, err
	}
	return &DelegationContractUpgradedNFTIterator{contract: _Delegation.contract, event: "ContractUpgradedNFT", logs: logs, sub: sub}, nil
}

// WatchContractUpgradedNFT is a free log subscription operation binding the contract event 0xe4dfc1498d5921eb1933069a6657c79e0b3fa55db46bf016d331de23cee72a3a.
//
// Solidity: event ContractUpgradedNFT(address indexed oldContract, address indexed newContract, uint8 indexed nftType)
func (_Delegation *DelegationFilterer) WatchContractUpgradedNFT(opts *bind.WatchOpts, sink chan<- *DelegationContractUpgradedNFT, oldContract []common.Address, newContract []common.Address, nftType []uint8) (event.Subscription, error) {

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

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "ContractUpgradedNFT", oldContractRule, newContractRule, nftTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationContractUpgradedNFT)
				if err := _Delegation.contract.UnpackLog(event, "ContractUpgradedNFT", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseContractUpgradedNFT(log types.Log) (*DelegationContractUpgradedNFT, error) {
	event := new(DelegationContractUpgradedNFT)
	if err := _Delegation.contract.UnpackLog(event, "ContractUpgradedNFT", log); err != nil {
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

// DelegationNFTCreatedIterator is returned from FilterNFTCreated and is used to iterate over the raw logs and unpacked data for NFTCreated events raised by the Delegation contract.
type DelegationNFTCreatedIterator struct {
	Event *DelegationNFTCreated // Event containing the contract specifics and raw log

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
func (it *DelegationNFTCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNFTCreated)
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
		it.Event = new(DelegationNFTCreated)
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
func (it *DelegationNFTCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNFTCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNFTCreated represents a NFTCreated event raised by the Delegation contract.
type DelegationNFTCreated struct {
	TokenAddress common.Address
	NftType      uint8
	Nft          DecimalNFTCenterNFT
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNFTCreated is a free log retrieval operation binding the contract event 0x90f6044530e519155c7f7c84474c4a4fdfb6b203716628919bb7b48908bad236.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 indexed nftType, (address,string,string,string,bool) nft)
func (_Delegation *DelegationFilterer) FilterNFTCreated(opts *bind.FilterOpts, nftType []uint8) (*DelegationNFTCreatedIterator, error) {

	var nftTypeRule []interface{}
	for _, nftTypeItem := range nftType {
		nftTypeRule = append(nftTypeRule, nftTypeItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "NFTCreated", nftTypeRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNFTCreatedIterator{contract: _Delegation.contract, event: "NFTCreated", logs: logs, sub: sub}, nil
}

// WatchNFTCreated is a free log subscription operation binding the contract event 0x90f6044530e519155c7f7c84474c4a4fdfb6b203716628919bb7b48908bad236.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 indexed nftType, (address,string,string,string,bool) nft)
func (_Delegation *DelegationFilterer) WatchNFTCreated(opts *bind.WatchOpts, sink chan<- *DelegationNFTCreated, nftType []uint8) (event.Subscription, error) {

	var nftTypeRule []interface{}
	for _, nftTypeItem := range nftType {
		nftTypeRule = append(nftTypeRule, nftTypeItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "NFTCreated", nftTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNFTCreated)
				if err := _Delegation.contract.UnpackLog(event, "NFTCreated", log); err != nil {
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

// ParseNFTCreated is a log parse operation binding the contract event 0x90f6044530e519155c7f7c84474c4a4fdfb6b203716628919bb7b48908bad236.
//
// Solidity: event NFTCreated(address tokenAddress, uint8 indexed nftType, (address,string,string,string,bool) nft)
func (_Delegation *DelegationFilterer) ParseNFTCreated(log types.Log) (*DelegationNFTCreated, error) {
	event := new(DelegationNFTCreated)
	if err := _Delegation.contract.UnpackLog(event, "NFTCreated", log); err != nil {
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
