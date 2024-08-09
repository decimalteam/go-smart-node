// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nft1155

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

// Nft1155MetaData contains all meta data concerning the Nft1155 contract.
var Nft1155MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC1155InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"idsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valuesLength\",\"type\":\"uint256\"}],\"name\":\"ERC1155InvalidArrayLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC1155MissingApprovalForAll\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserveType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MintNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCreator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNftDelegation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitUnauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReserveAlreadyInitialized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ContractURIUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DisabledMint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"ReserveUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addedReserveAmount\",\"type\":\"uint256\"}],\"name\":\"addReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"addReserveByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"addReserveByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contractURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"creator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNftType\",\"outputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRefundable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initialSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialContractURI\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialCreator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"initialRefundable\",\"type\":\"bool\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"mintByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"percentage\",\"type\":\"uint16\"}],\"name\":\"penalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"rate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Nft1155ABI is the input ABI used to generate the binding from.
// Deprecated: Use Nft1155MetaData.ABI instead.
var Nft1155ABI = Nft1155MetaData.ABI

// Nft1155 is an auto generated Go binding around an Ethereum contract.
type Nft1155 struct {
	Nft1155Caller     // Read-only binding to the contract
	Nft1155Transactor // Write-only binding to the contract
	Nft1155Filterer   // Log filterer for contract events
}

// Nft1155Caller is an auto generated read-only Go binding around an Ethereum contract.
type Nft1155Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Nft1155Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Nft1155Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Nft1155Session struct {
	Contract     *Nft1155          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Nft1155CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Nft1155CallerSession struct {
	Contract *Nft1155Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Nft1155TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Nft1155TransactorSession struct {
	Contract     *Nft1155Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Nft1155Raw is an auto generated low-level Go binding around an Ethereum contract.
type Nft1155Raw struct {
	Contract *Nft1155 // Generic contract binding to access the raw methods on
}

// Nft1155CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Nft1155CallerRaw struct {
	Contract *Nft1155Caller // Generic read-only contract binding to access the raw methods on
}

// Nft1155TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Nft1155TransactorRaw struct {
	Contract *Nft1155Transactor // Generic write-only contract binding to access the raw methods on
}

// NewNft1155 creates a new instance of Nft1155, bound to a specific deployed contract.
func NewNft1155(address common.Address, backend bind.ContractBackend) (*Nft1155, error) {
	contract, err := bindNft1155(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nft1155{Nft1155Caller: Nft1155Caller{contract: contract}, Nft1155Transactor: Nft1155Transactor{contract: contract}, Nft1155Filterer: Nft1155Filterer{contract: contract}}, nil
}

// NewNft1155Caller creates a new read-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Caller(address common.Address, caller bind.ContractCaller) (*Nft1155Caller, error) {
	contract, err := bindNft1155(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Caller{contract: contract}, nil
}

// NewNft1155Transactor creates a new write-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Transactor(address common.Address, transactor bind.ContractTransactor) (*Nft1155Transactor, error) {
	contract, err := bindNft1155(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Transactor{contract: contract}, nil
}

// NewNft1155Filterer creates a new log filterer instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Filterer(address common.Address, filterer bind.ContractFilterer) (*Nft1155Filterer, error) {
	contract, err := bindNft1155(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Nft1155Filterer{contract: contract}, nil
}

// bindNft1155 binds a generic wrapper to an already deployed contract.
func bindNft1155(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Nft1155MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.Nft1155Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Nft1155 *Nft1155Caller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Nft1155 *Nft1155Session) DOMAINSEPARATOR() ([32]byte, error) {
	return _Nft1155.Contract.DOMAINSEPARATOR(&_Nft1155.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_Nft1155 *Nft1155CallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _Nft1155.Contract.DOMAINSEPARATOR(&_Nft1155.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Nft1155 *Nft1155Caller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Nft1155 *Nft1155Session) MINRESERVE() (*big.Int, error) {
	return _Nft1155.Contract.MINRESERVE(&_Nft1155.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) MINRESERVE() (*big.Int, error) {
	return _Nft1155.Contract.MINRESERVE(&_Nft1155.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Caller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Session) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Caller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Session) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155CallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Nft1155 *Nft1155Caller) ContractURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "contractURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Nft1155 *Nft1155Session) ContractURI() (string, error) {
	return _Nft1155.Contract.ContractURI(&_Nft1155.CallOpts)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_Nft1155 *Nft1155CallerSession) ContractURI() (string, error) {
	return _Nft1155.Contract.ContractURI(&_Nft1155.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Nft1155 *Nft1155Caller) Creator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "creator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Nft1155 *Nft1155Session) Creator() (common.Address, error) {
	return _Nft1155.Contract.Creator(&_Nft1155.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_Nft1155 *Nft1155CallerSession) Creator() (common.Address, error) {
	return _Nft1155.Contract.Creator(&_Nft1155.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Nft1155 *Nft1155Caller) Exists(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "exists", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Nft1155 *Nft1155Session) Exists(id *big.Int) (bool, error) {
	return _Nft1155.Contract.Exists(&_Nft1155.CallOpts, id)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) Exists(id *big.Int) (bool, error) {
	return _Nft1155.Contract.Exists(&_Nft1155.CallOpts, id)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Nft1155 *Nft1155Caller) GetAllowMint(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "getAllowMint")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Nft1155 *Nft1155Session) GetAllowMint() (bool, error) {
	return _Nft1155.Contract.GetAllowMint(&_Nft1155.CallOpts)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_Nft1155 *Nft1155CallerSession) GetAllowMint() (bool, error) {
	return _Nft1155.Contract.GetAllowMint(&_Nft1155.CallOpts)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Nft1155 *Nft1155Caller) GetNftType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "getNftType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Nft1155 *Nft1155Session) GetNftType() (uint8, error) {
	return _Nft1155.Contract.GetNftType(&_Nft1155.CallOpts)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_Nft1155 *Nft1155CallerSession) GetNftType() (uint8, error) {
	return _Nft1155.Contract.GetNftType(&_Nft1155.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Nft1155 *Nft1155Caller) GetRefundable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "getRefundable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Nft1155 *Nft1155Session) GetRefundable() (bool, error) {
	return _Nft1155.Contract.GetRefundable(&_Nft1155.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_Nft1155 *Nft1155CallerSession) GetRefundable() (bool, error) {
	return _Nft1155.Contract.GetRefundable(&_Nft1155.CallOpts)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Nft1155 *Nft1155Caller) GetReserve(opts *bind.CallOpts, tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "getReserve", tokenId)

	if err != nil {
		return *new(IDecimalNFTCommonReserve), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalNFTCommonReserve)).(*IDecimalNFTCommonReserve)

	return out0, err

}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Nft1155 *Nft1155Session) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _Nft1155.Contract.GetReserve(&_Nft1155.CallOpts, tokenId)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_Nft1155 *Nft1155CallerSession) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _Nft1155.Contract.GetReserve(&_Nft1155.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155Caller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155Session) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, account, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nft1155 *Nft1155Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nft1155 *Nft1155Session) Name() (string, error) {
	return _Nft1155.Contract.Name(&_Nft1155.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Nft1155 *Nft1155CallerSession) Name() (string, error) {
	return _Nft1155.Contract.Name(&_Nft1155.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Nft1155 *Nft1155Caller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Nft1155 *Nft1155Session) Nonces(owner common.Address) (*big.Int, error) {
	return _Nft1155.Contract.Nonces(&_Nft1155.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _Nft1155.Contract.Nonces(&_Nft1155.CallOpts, owner)
}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Nft1155 *Nft1155Caller) Rate(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "rate", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Nft1155 *Nft1155Session) Rate(tokenId *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.Rate(&_Nft1155.CallOpts, tokenId)
}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_Nft1155 *Nft1155CallerSession) Rate(tokenId *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.Rate(&_Nft1155.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nft1155 *Nft1155Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nft1155 *Nft1155Session) Symbol() (string, error) {
	return _Nft1155.Contract.Symbol(&_Nft1155.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Nft1155 *Nft1155CallerSession) Symbol() (string, error) {
	return _Nft1155.Contract.Symbol(&_Nft1155.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Nft1155 *Nft1155Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Nft1155 *Nft1155Session) TotalSupply() (*big.Int, error) {
	return _Nft1155.Contract.TotalSupply(&_Nft1155.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) TotalSupply() (*big.Int, error) {
	return _Nft1155.Contract.TotalSupply(&_Nft1155.CallOpts)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Caller) TotalSupply0(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "totalSupply0", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Session) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.TotalSupply0(&_Nft1155.CallOpts, id)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.TotalSupply0(&_Nft1155.CallOpts, id)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Nft1155 *Nft1155Caller) Uri(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "uri", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Nft1155 *Nft1155Session) Uri(tokenId *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, tokenId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_Nft1155 *Nft1155CallerSession) Uri(tokenId *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, tokenId)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Nft1155 *Nft1155Transactor) AddReserve(opts *bind.TransactOpts, tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "addReserve", tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Nft1155 *Nft1155Session) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserve(&_Nft1155.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_Nft1155 *Nft1155TransactorSession) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserve(&_Nft1155.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Nft1155 *Nft1155Transactor) AddReserveByETH(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "addReserveByETH", tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Nft1155 *Nft1155Session) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserveByETH(&_Nft1155.TransactOpts, tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_Nft1155 *Nft1155TransactorSession) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserveByETH(&_Nft1155.TransactOpts, tokenId)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Transactor) AddReserveByPermit(opts *bind.TransactOpts, tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "addReserveByPermit", tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Session) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserveByPermit(&_Nft1155.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155TransactorSession) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.AddReserveByPermit(&_Nft1155.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Nft1155 *Nft1155Transactor) Burn(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "burn", tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Nft1155 *Nft1155Session) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Burn(&_Nft1155.TransactOpts, tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_Nft1155 *Nft1155TransactorSession) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Burn(&_Nft1155.TransactOpts, tokenId, amount)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Nft1155 *Nft1155Transactor) DisableMint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "disableMint")
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Nft1155 *Nft1155Session) DisableMint() (*types.Transaction, error) {
	return _Nft1155.Contract.DisableMint(&_Nft1155.TransactOpts)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_Nft1155 *Nft1155TransactorSession) DisableMint() (*types.Transaction, error) {
	return _Nft1155.Contract.DisableMint(&_Nft1155.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Nft1155 *Nft1155Transactor) Initialize(opts *bind.TransactOpts, initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "initialize", initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Nft1155 *Nft1155Session) Initialize(initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Nft1155.Contract.Initialize(&_Nft1155.TransactOpts, initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialContractURI, address initialCreator, bool initialRefundable) returns()
func (_Nft1155 *Nft1155TransactorSession) Initialize(initialSymbol string, initialName string, initialContractURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _Nft1155.Contract.Initialize(&_Nft1155.TransactOpts, initialSymbol, initialName, initialContractURI, initialCreator, initialRefundable)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Nft1155 *Nft1155Transactor) Mint(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mint", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Nft1155 *Nft1155Session) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_Nft1155 *Nft1155TransactorSession) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Nft1155 *Nft1155Transactor) MintByETH(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mintByETH", recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Nft1155 *Nft1155Session) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Nft1155.Contract.MintByETH(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_Nft1155 *Nft1155TransactorSession) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _Nft1155.Contract.MintByETH(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Transactor) MintByPermit(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mintByPermit", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Session) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.MintByPermit(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155TransactorSession) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.MintByPermit(&_Nft1155.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Nft1155 *Nft1155Transactor) Penalty(opts *bind.TransactOpts, tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "penalty", tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Nft1155 *Nft1155Session) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.Penalty(&_Nft1155.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_Nft1155 *Nft1155TransactorSession) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.Penalty(&_Nft1155.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Transactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "permit", owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155Session) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.Permit(&_Nft1155.TransactOpts, owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Nft1155 *Nft1155TransactorSession) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Nft1155.Contract.Permit(&_Nft1155.TransactOpts, owner, spender, deadline, v, r, s)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, values, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeTransferFrom", from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, value, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// Nft1155ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Nft1155 contract.
type Nft1155ApprovalForAllIterator struct {
	Event *Nft1155ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *Nft1155ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ApprovalForAll)
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
		it.Event = new(Nft1155ApprovalForAll)
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
func (it *Nft1155ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ApprovalForAll represents a ApprovalForAll event raised by the Nft1155 contract.
type Nft1155ApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*Nft1155ApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155ApprovalForAllIterator{contract: _Nft1155.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Nft1155ApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ApprovalForAll)
				if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseApprovalForAll(log types.Log) (*Nft1155ApprovalForAll, error) {
	event := new(Nft1155ApprovalForAll)
	if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155ContractURIUpdatedIterator is returned from FilterContractURIUpdated and is used to iterate over the raw logs and unpacked data for ContractURIUpdated events raised by the Nft1155 contract.
type Nft1155ContractURIUpdatedIterator struct {
	Event *Nft1155ContractURIUpdated // Event containing the contract specifics and raw log

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
func (it *Nft1155ContractURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ContractURIUpdated)
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
		it.Event = new(Nft1155ContractURIUpdated)
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
func (it *Nft1155ContractURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ContractURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ContractURIUpdated represents a ContractURIUpdated event raised by the Nft1155 contract.
type Nft1155ContractURIUpdated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterContractURIUpdated is a free log retrieval operation binding the contract event 0xa5d4097edda6d87cb9329af83fb3712ef77eeb13738ffe43cc35a4ce305ad962.
//
// Solidity: event ContractURIUpdated()
func (_Nft1155 *Nft1155Filterer) FilterContractURIUpdated(opts *bind.FilterOpts) (*Nft1155ContractURIUpdatedIterator, error) {

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return &Nft1155ContractURIUpdatedIterator{contract: _Nft1155.contract, event: "ContractURIUpdated", logs: logs, sub: sub}, nil
}

// WatchContractURIUpdated is a free log subscription operation binding the contract event 0xa5d4097edda6d87cb9329af83fb3712ef77eeb13738ffe43cc35a4ce305ad962.
//
// Solidity: event ContractURIUpdated()
func (_Nft1155 *Nft1155Filterer) WatchContractURIUpdated(opts *bind.WatchOpts, sink chan<- *Nft1155ContractURIUpdated) (event.Subscription, error) {

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ContractURIUpdated)
				if err := _Nft1155.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseContractURIUpdated(log types.Log) (*Nft1155ContractURIUpdated, error) {
	event := new(Nft1155ContractURIUpdated)
	if err := _Nft1155.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155DisabledMintIterator is returned from FilterDisabledMint and is used to iterate over the raw logs and unpacked data for DisabledMint events raised by the Nft1155 contract.
type Nft1155DisabledMintIterator struct {
	Event *Nft1155DisabledMint // Event containing the contract specifics and raw log

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
func (it *Nft1155DisabledMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155DisabledMint)
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
		it.Event = new(Nft1155DisabledMint)
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
func (it *Nft1155DisabledMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155DisabledMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155DisabledMint represents a DisabledMint event raised by the Nft1155 contract.
type Nft1155DisabledMint struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDisabledMint is a free log retrieval operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_Nft1155 *Nft1155Filterer) FilterDisabledMint(opts *bind.FilterOpts) (*Nft1155DisabledMintIterator, error) {

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return &Nft1155DisabledMintIterator{contract: _Nft1155.contract, event: "DisabledMint", logs: logs, sub: sub}, nil
}

// WatchDisabledMint is a free log subscription operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_Nft1155 *Nft1155Filterer) WatchDisabledMint(opts *bind.WatchOpts, sink chan<- *Nft1155DisabledMint) (event.Subscription, error) {

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155DisabledMint)
				if err := _Nft1155.contract.UnpackLog(event, "DisabledMint", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseDisabledMint(log types.Log) (*Nft1155DisabledMint, error) {
	event := new(Nft1155DisabledMint)
	if err := _Nft1155.contract.UnpackLog(event, "DisabledMint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Nft1155 contract.
type Nft1155InitializedIterator struct {
	Event *Nft1155Initialized // Event containing the contract specifics and raw log

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
func (it *Nft1155InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155Initialized)
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
		it.Event = new(Nft1155Initialized)
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
func (it *Nft1155InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155Initialized represents a Initialized event raised by the Nft1155 contract.
type Nft1155Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Nft1155 *Nft1155Filterer) FilterInitialized(opts *bind.FilterOpts) (*Nft1155InitializedIterator, error) {

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Nft1155InitializedIterator{contract: _Nft1155.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Nft1155 *Nft1155Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Nft1155Initialized) (event.Subscription, error) {

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155Initialized)
				if err := _Nft1155.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseInitialized(log types.Log) (*Nft1155Initialized, error) {
	event := new(Nft1155Initialized)
	if err := _Nft1155.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155ReserveUpdatedIterator is returned from FilterReserveUpdated and is used to iterate over the raw logs and unpacked data for ReserveUpdated events raised by the Nft1155 contract.
type Nft1155ReserveUpdatedIterator struct {
	Event *Nft1155ReserveUpdated // Event containing the contract specifics and raw log

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
func (it *Nft1155ReserveUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ReserveUpdated)
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
		it.Event = new(Nft1155ReserveUpdated)
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
func (it *Nft1155ReserveUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ReserveUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ReserveUpdated represents a ReserveUpdated event raised by the Nft1155 contract.
type Nft1155ReserveUpdated struct {
	TokenId     *big.Int
	TotalSupply *big.Int
	Reserve     IDecimalNFTCommonReserve
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReserveUpdated is a free log retrieval operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
func (_Nft1155 *Nft1155Filterer) FilterReserveUpdated(opts *bind.FilterOpts) (*Nft1155ReserveUpdatedIterator, error) {

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return &Nft1155ReserveUpdatedIterator{contract: _Nft1155.contract, event: "ReserveUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveUpdated is a free log subscription operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
func (_Nft1155 *Nft1155Filterer) WatchReserveUpdated(opts *bind.WatchOpts, sink chan<- *Nft1155ReserveUpdated) (event.Subscription, error) {

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ReserveUpdated)
				if err := _Nft1155.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseReserveUpdated(log types.Log) (*Nft1155ReserveUpdated, error) {
	event := new(Nft1155ReserveUpdated)
	if err := _Nft1155.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Nft1155 contract.
type Nft1155TransferBatchIterator struct {
	Event *Nft1155TransferBatch // Event containing the contract specifics and raw log

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
func (it *Nft1155TransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferBatch)
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
		it.Event = new(Nft1155TransferBatch)
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
func (it *Nft1155TransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferBatch represents a TransferBatch event raised by the Nft1155 contract.
type Nft1155TransferBatch struct {
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
func (_Nft1155 *Nft1155Filterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferBatchIterator, error) {

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

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferBatchIterator{contract: _Nft1155.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *Nft1155TransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferBatch)
				if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseTransferBatch(log types.Log) (*Nft1155TransferBatch, error) {
	event := new(Nft1155TransferBatch)
	if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Nft1155 contract.
type Nft1155TransferSingleIterator struct {
	Event *Nft1155TransferSingle // Event containing the contract specifics and raw log

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
func (it *Nft1155TransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferSingle)
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
		it.Event = new(Nft1155TransferSingle)
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
func (it *Nft1155TransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferSingle represents a TransferSingle event raised by the Nft1155 contract.
type Nft1155TransferSingle struct {
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
func (_Nft1155 *Nft1155Filterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferSingleIterator, error) {

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

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferSingleIterator{contract: _Nft1155.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *Nft1155TransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferSingle)
				if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseTransferSingle(log types.Log) (*Nft1155TransferSingle, error) {
	event := new(Nft1155TransferSingle)
	if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155URIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the Nft1155 contract.
type Nft1155URIIterator struct {
	Event *Nft1155URI // Event containing the contract specifics and raw log

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
func (it *Nft1155URIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155URI)
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
		it.Event = new(Nft1155URI)
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
func (it *Nft1155URIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155URIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155URI represents a URI event raised by the Nft1155 contract.
type Nft1155URI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*Nft1155URIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155URIIterator{contract: _Nft1155.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) WatchURI(opts *bind.WatchOpts, sink chan<- *Nft1155URI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155URI)
				if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseURI(log types.Log) (*Nft1155URI, error) {
	event := new(Nft1155URI)
	if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
