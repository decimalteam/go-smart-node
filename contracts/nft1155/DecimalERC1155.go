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

// IDecimalNFTCommonReserve is an auto generated low-level Go binding around an user-defined struct.
type IDecimalNFTCommonReserve struct {
	Token       common.Address
	Amount      *big.Int
	ReserveType uint8
}

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC1155InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"idsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valuesLength\",\"type\":\"uint256\"}],\"name\":\"ERC1155InvalidArrayLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC1155MissingApprovalForAll\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserveType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MintNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCreator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNftDelegation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitUnauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReserveAlreadyInitialized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ContractURIUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DisabledMint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"ReserveUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addedReserveAmount\",\"type\":\"uint256\"}],\"name\":\"addReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"addReserveByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"addReserveByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contractURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"creator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNftType\",\"outputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRefundable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initialSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialContractURI\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialCreator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"initialRefundable\",\"type\":\"bool\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"mintByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"percentage\",\"type\":\"uint16\"}],\"name\":\"penalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"rate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Delegation *DelegationCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Delegation *DelegationSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Delegation.Contract.BalanceOf(&_Delegation.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Delegation *DelegationCallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Delegation.Contract.BalanceOf(&_Delegation.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Delegation *DelegationCaller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Delegation *DelegationSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Delegation.Contract.BalanceOfBatch(&_Delegation.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Delegation *DelegationCallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Delegation.Contract.BalanceOfBatch(&_Delegation.CallOpts, accounts, ids)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Delegation *DelegationCaller) ContractURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "contractURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Delegation *DelegationSession) ContractURI() (string, error) {
	return _Delegation.Contract.ContractURI(&_Delegation.CallOpts)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Delegation *DelegationCallerSession) ContractURI() (string, error) {
	return _Delegation.Contract.ContractURI(&_Delegation.CallOpts)
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

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Delegation *DelegationCaller) Exists(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "exists", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Delegation *DelegationSession) Exists(id *big.Int) (bool, error) {
	return _Delegation.Contract.Exists(&_Delegation.CallOpts, id)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Delegation *DelegationCallerSession) Exists(id *big.Int) (bool, error) {
	return _Delegation.Contract.Exists(&_Delegation.CallOpts, id)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Delegation *DelegationCaller) GetAllowMint(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getAllowMint")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Delegation *DelegationSession) GetAllowMint() (bool, error) {
	return _Delegation.Contract.GetAllowMint(&_Delegation.CallOpts)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Delegation *DelegationCallerSession) GetAllowMint() (bool, error) {
	return _Delegation.Contract.GetAllowMint(&_Delegation.CallOpts)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Delegation *DelegationCaller) GetNftType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getNftType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Delegation *DelegationSession) GetNftType() (uint8, error) {
	return _Delegation.Contract.GetNftType(&_Delegation.CallOpts)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Delegation *DelegationCallerSession) GetNftType() (uint8, error) {
	return _Delegation.Contract.GetNftType(&_Delegation.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Delegation *DelegationCaller) GetRefundable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getRefundable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Delegation *DelegationSession) GetRefundable() (bool, error) {
	return _Delegation.Contract.GetRefundable(&_Delegation.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Delegation *DelegationCallerSession) GetRefundable() (bool, error) {
	return _Delegation.Contract.GetRefundable(&_Delegation.CallOpts)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Delegation *DelegationCaller) GetReserve(opts *bind.CallOpts, tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getReserve", tokenId)

	if err != nil {
		return *new(IDecimalNFTCommonReserve), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalNFTCommonReserve)).(*IDecimalNFTCommonReserve)

	return out0, err

}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Delegation *DelegationSession) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _Delegation.Contract.GetReserve(&_Delegation.CallOpts, tokenId)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Delegation *DelegationCallerSession) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _Delegation.Contract.GetReserve(&_Delegation.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Delegation *DelegationCaller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Delegation *DelegationSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Delegation.Contract.IsApprovedForAll(&_Delegation.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Delegation *DelegationCallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Delegation.Contract.IsApprovedForAll(&_Delegation.CallOpts, account, operator)
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

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Delegation *DelegationCaller) Rate(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "rate", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Delegation *DelegationSession) Rate(tokenId *big.Int) (*big.Int, error) {
	return _Delegation.Contract.Rate(&_Delegation.CallOpts, tokenId)
}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Delegation *DelegationCallerSession) Rate(tokenId *big.Int) (*big.Int, error) {
	return _Delegation.Contract.Rate(&_Delegation.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegation *DelegationCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegation *DelegationSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegation.Contract.SupportsInterface(&_Delegation.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegation *DelegationCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegation.Contract.SupportsInterface(&_Delegation.CallOpts, interfaceId)
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

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Delegation *DelegationCaller) TotalSupply0(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "totalSupply0", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Delegation *DelegationSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _Delegation.Contract.TotalSupply0(&_Delegation.CallOpts, id)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Delegation *DelegationCallerSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _Delegation.Contract.TotalSupply0(&_Delegation.CallOpts, id)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Delegation *DelegationCaller) Uri(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "uri", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Delegation *DelegationSession) Uri(tokenId *big.Int) (string, error) {
	return _Delegation.Contract.Uri(&_Delegation.CallOpts, tokenId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Delegation *DelegationCallerSession) Uri(tokenId *big.Int) (string, error) {
	return _Delegation.Contract.Uri(&_Delegation.CallOpts, tokenId)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Delegation *DelegationTransactor) AddReserve(opts *bind.TransactOpts, tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "addReserve", tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Delegation *DelegationSession) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserve(&_Delegation.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Delegation *DelegationTransactorSession) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserve(&_Delegation.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Delegation *DelegationTransactor) AddReserveByETH(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "addReserveByETH", tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Delegation *DelegationSession) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserveByETH(&_Delegation.TransactOpts, tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Delegation *DelegationTransactorSession) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserveByETH(&_Delegation.TransactOpts, tokenId)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) AddReserveByPermit(opts *bind.TransactOpts, tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "addReserveByPermit", tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserveByPermit(&_Delegation.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.AddReserveByPermit(&_Delegation.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Delegation *DelegationTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "burn", tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Delegation *DelegationSession) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Burn(&_Delegation.TransactOpts, tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Delegation *DelegationTransactorSession) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Burn(&_Delegation.TransactOpts, tokenId, amount)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Delegation *DelegationTransactor) DisableMint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "disableMint")
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Delegation *DelegationSession) DisableMint() (*types.Transaction, error) {
	return _Delegation.Contract.DisableMint(&_Delegation.TransactOpts)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Delegation *DelegationTransactorSession) DisableMint() (*types.Transaction, error) {
	return _Delegation.Contract.DisableMint(&_Delegation.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Delegation *DelegationTransactor) Initialize(opts *bind.TransactOpts, initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "initialize", initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Delegation *DelegationSession) Initialize(initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts, initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Delegation *DelegationTransactorSession) Initialize(initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Delegation.Contract.Initialize(&_Delegation.TransactOpts, initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Delegation *DelegationTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "mint", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Delegation *DelegationSession) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Mint(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Delegation *DelegationTransactorSession) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Mint(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Delegation *DelegationTransactor) MintByETH(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "mintByETH", recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Delegation *DelegationSession) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Delegation.Contract.MintByETH(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Delegation *DelegationTransactorSession) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Delegation.Contract.MintByETH(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) MintByPermit(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "mintByPermit", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.MintByPermit(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.MintByPermit(&_Delegation.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Delegation *DelegationTransactor) Penalty(opts *bind.TransactOpts, tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "penalty", tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Delegation *DelegationSession) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Delegation.Contract.Penalty(&_Delegation.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Delegation *DelegationTransactorSession) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Delegation.Contract.Penalty(&_Delegation.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "permit", owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Permit(&_Delegation.TransactOpts, owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.Permit(&_Delegation.TransactOpts, owner, spender, deadline, v, r, s)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Delegation *DelegationTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Delegation *DelegationSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.SafeBatchTransferFrom(&_Delegation.TransactOpts, from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Delegation *DelegationTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.SafeBatchTransferFrom(&_Delegation.TransactOpts, from, to, ids, values, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Delegation *DelegationTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "safeTransferFrom", from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Delegation *DelegationSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.SafeTransferFrom(&_Delegation.TransactOpts, from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Delegation *DelegationTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.SafeTransferFrom(&_Delegation.TransactOpts, from, to, id, value, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Delegation *DelegationTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Delegation *DelegationSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Delegation.Contract.SetApprovalForAll(&_Delegation.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Delegation *DelegationTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Delegation.Contract.SetApprovalForAll(&_Delegation.TransactOpts, operator, approved)
}

// DelegationApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Delegation contract.
type DelegationApprovalForAllIterator struct {
	Event *DelegationApprovalForAll // Event containing the contract specifics and raw log

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
func (it *DelegationApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationApprovalForAll)
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
		it.Event = new(DelegationApprovalForAll)
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
func (it *DelegationApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationApprovalForAll represents a ApprovalForAll event raised by the Delegation contract.
type DelegationApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Delegation *DelegationFilterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*DelegationApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationApprovalForAllIterator{contract: _Delegation.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Delegation *DelegationFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *DelegationApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationApprovalForAll)
				if err := _Delegation.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Delegation *DelegationFilterer) ParseApprovalForAll(log types.Log) (*DelegationApprovalForAll, error) {
	event := new(DelegationApprovalForAll)
	if err := _Delegation.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationContractURIUpdatedIterator is returned from FilterContractURIUpdated and is used to iterate over the raw logs and unpacked data for ContractURIUpdated events raised by the Delegation contract.
type DelegationContractURIUpdatedIterator struct {
	Event *DelegationContractURIUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationContractURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationContractURIUpdated)
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
		it.Event = new(DelegationContractURIUpdated)
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
func (it *DelegationContractURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationContractURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationContractURIUpdated represents a ContractURIUpdated event raised by the Delegation contract.
type DelegationContractURIUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterContractURIUpdated is a free log retrieval operation binding the contract event 0xa5d4097edda6d87cb9329af83fb3712ef77eeb13738ffe43cc35a4ce305ad962.
//
// Solidity: event ContractURIUpdated()
func (_Delegation *DelegationFilterer) FilterContractURIUpdated(opts *bind.FilterOpts) (*DelegationContractURIUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationContractURIUpdatedIterator{contract: _Delegation.contract, event: "ContractURIUpdated", logs: logs, sub: sub}, nil
}

// WatchContractURIUpdated is a free log subscription operation binding the contract event 0xa5d4097edda6d87cb9329af83fb3712ef77eeb13738ffe43cc35a4ce305ad962.
//
// Solidity: event ContractURIUpdated()
func (_Delegation *DelegationFilterer) WatchContractURIUpdated(opts *bind.WatchOpts, sink chan<- *DelegationContractURIUpdated) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationContractURIUpdated)
				if err := _Delegation.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
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

// ParseContractURIUpdated is a log parse operation binding the contract event 0xa5d4097edda6d87cb9329af83fb3712ef77eeb13738ffe43cc35a4ce305ad962.
//
// Solidity: event ContractURIUpdated()
func (_Delegation *DelegationFilterer) ParseContractURIUpdated(log types.Log) (*DelegationContractURIUpdated, error) {
	event := new(DelegationContractURIUpdated)
	if err := _Delegation.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationDisabledMintIterator is returned from FilterDisabledMint and is used to iterate over the raw logs and unpacked data for DisabledMint events raised by the Delegation contract.
type DelegationDisabledMintIterator struct {
	Event *DelegationDisabledMint // Event containing the contract specifics and raw log

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
func (it *DelegationDisabledMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationDisabledMint)
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
		it.Event = new(DelegationDisabledMint)
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
func (it *DelegationDisabledMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationDisabledMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationDisabledMint represents a DisabledMint event raised by the Delegation contract.
type DelegationDisabledMint struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDisabledMint is a free log retrieval operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_Delegation *DelegationFilterer) FilterDisabledMint(opts *bind.FilterOpts) (*DelegationDisabledMintIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return &DelegationDisabledMintIterator{contract: _Delegation.contract, event: "DisabledMint", logs: logs, sub: sub}, nil
}

// WatchDisabledMint is a free log subscription operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_Delegation *DelegationFilterer) WatchDisabledMint(opts *bind.WatchOpts, sink chan<- *DelegationDisabledMint) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationDisabledMint)
				if err := _Delegation.contract.UnpackLog(event, "DisabledMint", log); err != nil {
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

// ParseDisabledMint is a log parse operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_Delegation *DelegationFilterer) ParseDisabledMint(log types.Log) (*DelegationDisabledMint, error) {
	event := new(DelegationDisabledMint)
	if err := _Delegation.contract.UnpackLog(event, "DisabledMint", log); err != nil {
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
	TokenId     *big.Int
	TotalSupply *big.Int
	Reserve     IDecimalNFTCommonReserve
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReserveUpdated is a free log retrieval operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
func (_Delegation *DelegationFilterer) FilterReserveUpdated(opts *bind.FilterOpts) (*DelegationReserveUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationReserveUpdatedIterator{contract: _Delegation.contract, event: "ReserveUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveUpdated is a free log subscription operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
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

// ParseReserveUpdated is a log parse operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
func (_Delegation *DelegationFilterer) ParseReserveUpdated(log types.Log) (*DelegationReserveUpdated, error) {
	event := new(DelegationReserveUpdated)
	if err := _Delegation.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Delegation contract.
type DelegationTransferBatchIterator struct {
	Event *DelegationTransferBatch // Event containing the contract specifics and raw log

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
func (it *DelegationTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTransferBatch)
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
		it.Event = new(DelegationTransferBatch)
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
func (it *DelegationTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTransferBatch represents a TransferBatch event raised by the Delegation contract.
type DelegationTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Delegation *DelegationFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*DelegationTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationTransferBatchIterator{contract: _Delegation.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Delegation *DelegationFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *DelegationTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTransferBatch)
				if err := _Delegation.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Delegation *DelegationFilterer) ParseTransferBatch(log types.Log) (*DelegationTransferBatch, error) {
	event := new(DelegationTransferBatch)
	if err := _Delegation.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Delegation contract.
type DelegationTransferSingleIterator struct {
	Event *DelegationTransferSingle // Event containing the contract specifics and raw log

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
func (it *DelegationTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTransferSingle)
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
		it.Event = new(DelegationTransferSingle)
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
func (it *DelegationTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTransferSingle represents a TransferSingle event raised by the Delegation contract.
type DelegationTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Delegation *DelegationFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*DelegationTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationTransferSingleIterator{contract: _Delegation.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Delegation *DelegationFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *DelegationTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTransferSingle)
				if err := _Delegation.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Delegation *DelegationFilterer) ParseTransferSingle(log types.Log) (*DelegationTransferSingle, error) {
	event := new(DelegationTransferSingle)
	if err := _Delegation.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the Delegation contract.
type DelegationURIIterator struct {
	Event *DelegationURI // Event containing the contract specifics and raw log

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
func (it *DelegationURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationURI)
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
		it.Event = new(DelegationURI)
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
func (it *DelegationURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationURI represents a URI event raised by the Delegation contract.
type DelegationURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Delegation *DelegationFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*DelegationURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &DelegationURIIterator{contract: _Delegation.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Delegation *DelegationFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *DelegationURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationURI)
				if err := _Delegation.contract.UnpackLog(event, "URI", log); err != nil {
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

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Delegation *DelegationFilterer) ParseURI(log types.Log) (*DelegationURI, error) {
	event := new(DelegationURI)
	if err := _Delegation.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
