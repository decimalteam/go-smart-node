// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegationNft

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

// DelegationNftMetaData contains all meta data concerning the DelegationNft contract.
var DelegationNftMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserveType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MintNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNftDelegation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitUnauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReserveAlreadyInitialized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DisabledMint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"ReserveUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CREATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addedReserveAmount\",\"type\":\"uint256\"}],\"name\":\"addReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"addReserveByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"addReserveByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNftType\",\"outputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRefundable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initialSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialBaseURI\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialCreator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"initialRefundable\",\"type\":\"bool\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintByETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"mintByPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"percentage\",\"type\":\"uint16\"}],\"name\":\"penalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"setTokenURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DelegationNftABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegationNftMetaData.ABI instead.
var DelegationNftABI = DelegationNftMetaData.ABI

// DelegationNft is an auto generated Go binding around an Ethereum contract.
type DelegationNft struct {
	DelegationNftCaller     // Read-only binding to the contract
	DelegationNftTransactor // Write-only binding to the contract
	DelegationNftFilterer   // Log filterer for contract events
}

// DelegationNftCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegationNftCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationNftTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegationNftTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationNftFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegationNftFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationNftSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegationNftSession struct {
	Contract     *DelegationNft    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegationNftCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegationNftCallerSession struct {
	Contract *DelegationNftCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// DelegationNftTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegationNftTransactorSession struct {
	Contract     *DelegationNftTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DelegationNftRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegationNftRaw struct {
	Contract *DelegationNft // Generic contract binding to access the raw methods on
}

// DelegationNftCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegationNftCallerRaw struct {
	Contract *DelegationNftCaller // Generic read-only contract binding to access the raw methods on
}

// DelegationNftTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegationNftTransactorRaw struct {
	Contract *DelegationNftTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegationNft creates a new instance of DelegationNft, bound to a specific deployed contract.
func NewDelegationNft(address common.Address, backend bind.ContractBackend) (*DelegationNft, error) {
	contract, err := bindDelegationNft(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegationNft{DelegationNftCaller: DelegationNftCaller{contract: contract}, DelegationNftTransactor: DelegationNftTransactor{contract: contract}, DelegationNftFilterer: DelegationNftFilterer{contract: contract}}, nil
}

// NewDelegationNftCaller creates a new read-only instance of DelegationNft, bound to a specific deployed contract.
func NewDelegationNftCaller(address common.Address, caller bind.ContractCaller) (*DelegationNftCaller, error) {
	contract, err := bindDelegationNft(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationNftCaller{contract: contract}, nil
}

// NewDelegationNftTransactor creates a new write-only instance of DelegationNft, bound to a specific deployed contract.
func NewDelegationNftTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegationNftTransactor, error) {
	contract, err := bindDelegationNft(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransactor{contract: contract}, nil
}

// NewDelegationNftFilterer creates a new log filterer instance of DelegationNft, bound to a specific deployed contract.
func NewDelegationNftFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegationNftFilterer, error) {
	contract, err := bindDelegationNft(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegationNftFilterer{contract: contract}, nil
}

// bindDelegationNft binds a generic wrapper to an already deployed contract.
func bindDelegationNft(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegationNftMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegationNft *DelegationNftRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegationNft.Contract.DelegationNftCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegationNft *DelegationNftRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegationNftTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegationNft *DelegationNftRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegationNftTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegationNft *DelegationNftCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegationNft.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegationNft *DelegationNftTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegationNft.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegationNft *DelegationNftTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegationNft.Contract.contract.Transact(opts, method, params...)
}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftCaller) CREATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "CREATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftSession) CREATORROLE() ([32]byte, error) {
	return _DelegationNft.Contract.CREATORROLE(&_DelegationNft.CallOpts)
}

// CREATORROLE is a free data retrieval call binding the contract method 0x8aeda25a.
//
// Solidity: function CREATOR_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) CREATORROLE() ([32]byte, error) {
	return _DelegationNft.Contract.CREATORROLE(&_DelegationNft.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DelegationNft.Contract.DEFAULTADMINROLE(&_DelegationNft.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DelegationNft.Contract.DEFAULTADMINROLE(&_DelegationNft.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_DelegationNft *DelegationNftCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_DelegationNft *DelegationNftSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _DelegationNft.Contract.DOMAINSEPARATOR(&_DelegationNft.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _DelegationNft.Contract.DOMAINSEPARATOR(&_DelegationNft.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_DelegationNft *DelegationNftCaller) MINRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "MIN_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_DelegationNft *DelegationNftSession) MINRESERVE() (*big.Int, error) {
	return _DelegationNft.Contract.MINRESERVE(&_DelegationNft.CallOpts)
}

// MINRESERVE is a free data retrieval call binding the contract method 0x09dfd0e6.
//
// Solidity: function MIN_RESERVE() view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) MINRESERVE() (*big.Int, error) {
	return _DelegationNft.Contract.MINRESERVE(&_DelegationNft.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _DelegationNft.Contract.BalanceOf(&_DelegationNft.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _DelegationNft.Contract.BalanceOf(&_DelegationNft.CallOpts, owner)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_DelegationNft *DelegationNftCaller) GetAllowMint(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getAllowMint")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_DelegationNft *DelegationNftSession) GetAllowMint() (bool, error) {
	return _DelegationNft.Contract.GetAllowMint(&_DelegationNft.CallOpts)
}

// GetAllowMint is a free data retrieval call binding the contract method 0xa6fde7ab.
//
// Solidity: function getAllowMint() view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) GetAllowMint() (bool, error) {
	return _DelegationNft.Contract.GetAllowMint(&_DelegationNft.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _DelegationNft.Contract.GetApproved(&_DelegationNft.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _DelegationNft.Contract.GetApproved(&_DelegationNft.CallOpts, tokenId)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_DelegationNft *DelegationNftCaller) GetNftType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getNftType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_DelegationNft *DelegationNftSession) GetNftType() (uint8, error) {
	return _DelegationNft.Contract.GetNftType(&_DelegationNft.CallOpts)
}

// GetNftType is a free data retrieval call binding the contract method 0x5bfb797f.
//
// Solidity: function getNftType() pure returns(uint8)
func (_DelegationNft *DelegationNftCallerSession) GetNftType() (uint8, error) {
	return _DelegationNft.Contract.GetNftType(&_DelegationNft.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_DelegationNft *DelegationNftCaller) GetRefundable(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getRefundable")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_DelegationNft *DelegationNftSession) GetRefundable() (bool, error) {
	return _DelegationNft.Contract.GetRefundable(&_DelegationNft.CallOpts)
}

// GetRefundable is a free data retrieval call binding the contract method 0x02912524.
//
// Solidity: function getRefundable() view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) GetRefundable() (bool, error) {
	return _DelegationNft.Contract.GetRefundable(&_DelegationNft.CallOpts)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_DelegationNft *DelegationNftCaller) GetReserve(opts *bind.CallOpts, tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getReserve", tokenId)

	if err != nil {
		return *new(IDecimalNFTCommonReserve), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalNFTCommonReserve)).(*IDecimalNFTCommonReserve)

	return out0, err

}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_DelegationNft *DelegationNftSession) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _DelegationNft.Contract.GetReserve(&_DelegationNft.CallOpts, tokenId)
}

// GetReserve is a free data retrieval call binding the contract method 0x77778db3.
//
// Solidity: function getReserve(uint256 tokenId) view returns((address,uint256,uint8))
func (_DelegationNft *DelegationNftCallerSession) GetReserve(tokenId *big.Int) (IDecimalNFTCommonReserve, error) {
	return _DelegationNft.Contract.GetReserve(&_DelegationNft.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegationNft *DelegationNftCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegationNft *DelegationNftSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DelegationNft.Contract.GetRoleAdmin(&_DelegationNft.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DelegationNft.Contract.GetRoleAdmin(&_DelegationNft.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegationNft *DelegationNftCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegationNft *DelegationNftSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DelegationNft.Contract.HasRole(&_DelegationNft.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DelegationNft.Contract.HasRole(&_DelegationNft.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_DelegationNft *DelegationNftCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_DelegationNft *DelegationNftSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _DelegationNft.Contract.IsApprovedForAll(&_DelegationNft.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _DelegationNft.Contract.IsApprovedForAll(&_DelegationNft.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_DelegationNft *DelegationNftCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_DelegationNft *DelegationNftSession) Name() (string, error) {
	return _DelegationNft.Contract.Name(&_DelegationNft.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_DelegationNft *DelegationNftCallerSession) Name() (string, error) {
	return _DelegationNft.Contract.Name(&_DelegationNft.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x141a468c.
//
// Solidity: function nonces(uint256 tokenId) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) Nonces(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "nonces", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x141a468c.
//
// Solidity: function nonces(uint256 tokenId) view returns(uint256)
func (_DelegationNft *DelegationNftSession) Nonces(tokenId *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.Nonces(&_DelegationNft.CallOpts, tokenId)
}

// Nonces is a free data retrieval call binding the contract method 0x141a468c.
//
// Solidity: function nonces(uint256 tokenId) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) Nonces(tokenId *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.Nonces(&_DelegationNft.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _DelegationNft.Contract.OwnerOf(&_DelegationNft.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_DelegationNft *DelegationNftCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _DelegationNft.Contract.OwnerOf(&_DelegationNft.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegationNft *DelegationNftCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegationNft *DelegationNftSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DelegationNft.Contract.SupportsInterface(&_DelegationNft.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DelegationNft.Contract.SupportsInterface(&_DelegationNft.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_DelegationNft *DelegationNftCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_DelegationNft *DelegationNftSession) Symbol() (string, error) {
	return _DelegationNft.Contract.Symbol(&_DelegationNft.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_DelegationNft *DelegationNftCallerSession) Symbol() (string, error) {
	return _DelegationNft.Contract.Symbol(&_DelegationNft.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftSession) TokenURI(tokenId *big.Int) (string, error) {
	return _DelegationNft.Contract.TokenURI(&_DelegationNft.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _DelegationNft.Contract.TokenURI(&_DelegationNft.CallOpts, tokenId)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_DelegationNft *DelegationNftTransactor) AddReserve(opts *bind.TransactOpts, tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "addReserve", tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_DelegationNft *DelegationNftSession) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserve(&_DelegationNft.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserve is a paid mutator transaction binding the contract method 0x726f77e3.
//
// Solidity: function addReserve(uint256 tokenId, uint256 addedReserveAmount) returns()
func (_DelegationNft *DelegationNftTransactorSession) AddReserve(tokenId *big.Int, addedReserveAmount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserve(&_DelegationNft.TransactOpts, tokenId, addedReserveAmount)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_DelegationNft *DelegationNftTransactor) AddReserveByETH(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "addReserveByETH", tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_DelegationNft *DelegationNftSession) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserveByETH(&_DelegationNft.TransactOpts, tokenId)
}

// AddReserveByETH is a paid mutator transaction binding the contract method 0x967c66a2.
//
// Solidity: function addReserveByETH(uint256 tokenId) payable returns()
func (_DelegationNft *DelegationNftTransactorSession) AddReserveByETH(tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserveByETH(&_DelegationNft.TransactOpts, tokenId)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) AddReserveByPermit(opts *bind.TransactOpts, tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "addReserveByPermit", tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserveByPermit(&_DelegationNft.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// AddReserveByPermit is a paid mutator transaction binding the contract method 0xdf5b24d4.
//
// Solidity: function addReserveByPermit(uint256 tokenId, uint256 reserveAmount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) AddReserveByPermit(tokenId *big.Int, reserveAmount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.AddReserveByPermit(&_DelegationNft.TransactOpts, tokenId, reserveAmount, deadline, v, r, s)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Approve(&_DelegationNft.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Approve(&_DelegationNft.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_DelegationNft *DelegationNftSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Burn(&_DelegationNft.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Burn(&_DelegationNft.TransactOpts, tokenId)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_DelegationNft *DelegationNftTransactor) DisableMint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "disableMint")
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_DelegationNft *DelegationNftSession) DisableMint() (*types.Transaction, error) {
	return _DelegationNft.Contract.DisableMint(&_DelegationNft.TransactOpts)
}

// DisableMint is a paid mutator transaction binding the contract method 0x34452f38.
//
// Solidity: function disableMint() returns()
func (_DelegationNft *DelegationNftTransactorSession) DisableMint() (*types.Transaction, error) {
	return _DelegationNft.Contract.DisableMint(&_DelegationNft.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.GrantRole(&_DelegationNft.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.GrantRole(&_DelegationNft.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialBaseURI, address initialCreator, bool initialRefundable) returns()
func (_DelegationNft *DelegationNftTransactor) Initialize(opts *bind.TransactOpts, initialSymbol string, initialName string, initialBaseURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "initialize", initialSymbol, initialName, initialBaseURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialBaseURI, address initialCreator, bool initialRefundable) returns()
func (_DelegationNft *DelegationNftSession) Initialize(initialSymbol string, initialName string, initialBaseURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _DelegationNft.Contract.Initialize(&_DelegationNft.TransactOpts, initialSymbol, initialName, initialBaseURI, initialCreator, initialRefundable)
}

// Initialize is a paid mutator transaction binding the contract method 0x28f7ed23.
//
// Solidity: function initialize(string initialSymbol, string initialName, string initialBaseURI, address initialCreator, bool initialRefundable) returns()
func (_DelegationNft *DelegationNftTransactorSession) Initialize(initialSymbol string, initialName string, initialBaseURI string, initialCreator common.Address, initialRefundable bool) (*types.Transaction, error) {
	return _DelegationNft.Contract.Initialize(&_DelegationNft.TransactOpts, initialSymbol, initialName, initialBaseURI, initialCreator, initialRefundable)
}

// Mint is a paid mutator transaction binding the contract method 0xedc6fb6d.
//
// Solidity: function mint(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken) returns(uint256)
func (_DelegationNft *DelegationNftTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mint", recipient, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0xedc6fb6d.
//
// Solidity: function mint(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken) returns(uint256)
func (_DelegationNft *DelegationNftSession) Mint(recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Mint(&_DelegationNft.TransactOpts, recipient, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0xedc6fb6d.
//
// Solidity: function mint(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken) returns(uint256)
func (_DelegationNft *DelegationNftTransactorSession) Mint(recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Mint(&_DelegationNft.TransactOpts, recipient, tokenURI, reserveAmount, reserveToken)
}

// MintByETH is a paid mutator transaction binding the contract method 0x054198c7.
//
// Solidity: function mintByETH(address recipient, string tokenURI) payable returns(uint256)
func (_DelegationNft *DelegationNftTransactor) MintByETH(opts *bind.TransactOpts, recipient common.Address, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mintByETH", recipient, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0x054198c7.
//
// Solidity: function mintByETH(address recipient, string tokenURI) payable returns(uint256)
func (_DelegationNft *DelegationNftSession) MintByETH(recipient common.Address, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByETH(&_DelegationNft.TransactOpts, recipient, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0x054198c7.
//
// Solidity: function mintByETH(address recipient, string tokenURI) payable returns(uint256)
func (_DelegationNft *DelegationNftTransactorSession) MintByETH(recipient common.Address, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByETH(&_DelegationNft.TransactOpts, recipient, tokenURI)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x29e76be6.
//
// Solidity: function mintByPermit(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_DelegationNft *DelegationNftTransactor) MintByPermit(opts *bind.TransactOpts, recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mintByPermit", recipient, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x29e76be6.
//
// Solidity: function mintByPermit(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_DelegationNft *DelegationNftSession) MintByPermit(recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByPermit(&_DelegationNft.TransactOpts, recipient, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x29e76be6.
//
// Solidity: function mintByPermit(address recipient, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_DelegationNft *DelegationNftTransactorSession) MintByPermit(recipient common.Address, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByPermit(&_DelegationNft.TransactOpts, recipient, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// Penalty is a paid mutator transaction binding the contract method 0x737c9cf6.
//
// Solidity: function penalty(uint256 tokenId, uint16 percentage) returns()
func (_DelegationNft *DelegationNftTransactor) Penalty(opts *bind.TransactOpts, tokenId *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "penalty", tokenId, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0x737c9cf6.
//
// Solidity: function penalty(uint256 tokenId, uint16 percentage) returns()
func (_DelegationNft *DelegationNftSession) Penalty(tokenId *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.Contract.Penalty(&_DelegationNft.TransactOpts, tokenId, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0x737c9cf6.
//
// Solidity: function penalty(uint256 tokenId, uint16 percentage) returns()
func (_DelegationNft *DelegationNftTransactorSession) Penalty(tokenId *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.Contract.Penalty(&_DelegationNft.TransactOpts, tokenId, percentage)
}

// Permit is a paid mutator transaction binding the contract method 0x7ac2ff7b.
//
// Solidity: function permit(address spender, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) Permit(opts *bind.TransactOpts, spender common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "permit", spender, tokenId, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x7ac2ff7b.
//
// Solidity: function permit(address spender, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) Permit(spender common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Permit(&_DelegationNft.TransactOpts, spender, tokenId, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x7ac2ff7b.
//
// Solidity: function permit(address spender, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) Permit(spender common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Permit(&_DelegationNft.TransactOpts, spender, tokenId, deadline, v, r, s)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DelegationNft *DelegationNftTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DelegationNft *DelegationNftSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.RenounceRole(&_DelegationNft.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DelegationNft *DelegationNftTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.RenounceRole(&_DelegationNft.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.RevokeRole(&_DelegationNft.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DelegationNft *DelegationNftTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.RevokeRole(&_DelegationNft.TransactOpts, role, account)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom(&_DelegationNft.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom(&_DelegationNft.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_DelegationNft *DelegationNftTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_DelegationNft *DelegationNftSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom0(&_DelegationNft.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_DelegationNft *DelegationNftTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom0(&_DelegationNft.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_DelegationNft *DelegationNftTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_DelegationNft *DelegationNftSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetApprovalForAll(&_DelegationNft.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_DelegationNft *DelegationNftTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetApprovalForAll(&_DelegationNft.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_DelegationNft *DelegationNftTransactor) SetBaseURI(opts *bind.TransactOpts, baseURI string) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "setBaseURI", baseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_DelegationNft *DelegationNftSession) SetBaseURI(baseURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetBaseURI(&_DelegationNft.TransactOpts, baseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string baseURI) returns()
func (_DelegationNft *DelegationNftTransactorSession) SetBaseURI(baseURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetBaseURI(&_DelegationNft.TransactOpts, baseURI)
}

// SetTokenURI is a paid mutator transaction binding the contract method 0x162094c4.
//
// Solidity: function setTokenURI(uint256 tokenId, string tokenURI) returns()
func (_DelegationNft *DelegationNftTransactor) SetTokenURI(opts *bind.TransactOpts, tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "setTokenURI", tokenId, tokenURI)
}

// SetTokenURI is a paid mutator transaction binding the contract method 0x162094c4.
//
// Solidity: function setTokenURI(uint256 tokenId, string tokenURI) returns()
func (_DelegationNft *DelegationNftSession) SetTokenURI(tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetTokenURI(&_DelegationNft.TransactOpts, tokenId, tokenURI)
}

// SetTokenURI is a paid mutator transaction binding the contract method 0x162094c4.
//
// Solidity: function setTokenURI(uint256 tokenId, string tokenURI) returns()
func (_DelegationNft *DelegationNftTransactorSession) SetTokenURI(tokenId *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetTokenURI(&_DelegationNft.TransactOpts, tokenId, tokenURI)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferFrom(&_DelegationNft.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferFrom(&_DelegationNft.TransactOpts, from, to, tokenId)
}

// DelegationNftApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the DelegationNft contract.
type DelegationNftApprovalIterator struct {
	Event *DelegationNftApproval // Event containing the contract specifics and raw log

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
func (it *DelegationNftApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftApproval)
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
		it.Event = new(DelegationNftApproval)
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
func (it *DelegationNftApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftApproval represents a Approval event raised by the DelegationNft contract.
type DelegationNftApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*DelegationNftApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftApprovalIterator{contract: _DelegationNft.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *DelegationNftApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftApproval)
				if err := _DelegationNft.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) ParseApproval(log types.Log) (*DelegationNftApproval, error) {
	event := new(DelegationNftApproval)
	if err := _DelegationNft.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the DelegationNft contract.
type DelegationNftApprovalForAllIterator struct {
	Event *DelegationNftApprovalForAll // Event containing the contract specifics and raw log

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
func (it *DelegationNftApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftApprovalForAll)
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
		it.Event = new(DelegationNftApprovalForAll)
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
func (it *DelegationNftApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftApprovalForAll represents a ApprovalForAll event raised by the DelegationNft contract.
type DelegationNftApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*DelegationNftApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftApprovalForAllIterator{contract: _DelegationNft.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *DelegationNftApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftApprovalForAll)
				if err := _DelegationNft.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) ParseApprovalForAll(log types.Log) (*DelegationNftApprovalForAll, error) {
	event := new(DelegationNftApprovalForAll)
	if err := _DelegationNft.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the DelegationNft contract.
type DelegationNftBatchMetadataUpdateIterator struct {
	Event *DelegationNftBatchMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *DelegationNftBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftBatchMetadataUpdate)
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
		it.Event = new(DelegationNftBatchMetadataUpdate)
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
func (it *DelegationNftBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the DelegationNft contract.
type DelegationNftBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_DelegationNft *DelegationNftFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*DelegationNftBatchMetadataUpdateIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &DelegationNftBatchMetadataUpdateIterator{contract: _DelegationNft.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_DelegationNft *DelegationNftFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *DelegationNftBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftBatchMetadataUpdate)
				if err := _DelegationNft.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
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

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_DelegationNft *DelegationNftFilterer) ParseBatchMetadataUpdate(log types.Log) (*DelegationNftBatchMetadataUpdate, error) {
	event := new(DelegationNftBatchMetadataUpdate)
	if err := _DelegationNft.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftDisabledMintIterator is returned from FilterDisabledMint and is used to iterate over the raw logs and unpacked data for DisabledMint events raised by the DelegationNft contract.
type DelegationNftDisabledMintIterator struct {
	Event *DelegationNftDisabledMint // Event containing the contract specifics and raw log

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
func (it *DelegationNftDisabledMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftDisabledMint)
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
		it.Event = new(DelegationNftDisabledMint)
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
func (it *DelegationNftDisabledMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftDisabledMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftDisabledMint represents a DisabledMint event raised by the DelegationNft contract.
type DelegationNftDisabledMint struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDisabledMint is a free log retrieval operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_DelegationNft *DelegationNftFilterer) FilterDisabledMint(opts *bind.FilterOpts) (*DelegationNftDisabledMintIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return &DelegationNftDisabledMintIterator{contract: _DelegationNft.contract, event: "DisabledMint", logs: logs, sub: sub}, nil
}

// WatchDisabledMint is a free log subscription operation binding the contract event 0x96786059fc12ef37dc62764d5fdd3131eeb87ad78f23b8476a8866eb7e6b57ce.
//
// Solidity: event DisabledMint()
func (_DelegationNft *DelegationNftFilterer) WatchDisabledMint(opts *bind.WatchOpts, sink chan<- *DelegationNftDisabledMint) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "DisabledMint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftDisabledMint)
				if err := _DelegationNft.contract.UnpackLog(event, "DisabledMint", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseDisabledMint(log types.Log) (*DelegationNftDisabledMint, error) {
	event := new(DelegationNftDisabledMint)
	if err := _DelegationNft.contract.UnpackLog(event, "DisabledMint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DelegationNft contract.
type DelegationNftInitializedIterator struct {
	Event *DelegationNftInitialized // Event containing the contract specifics and raw log

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
func (it *DelegationNftInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftInitialized)
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
		it.Event = new(DelegationNftInitialized)
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
func (it *DelegationNftInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftInitialized represents a Initialized event raised by the DelegationNft contract.
type DelegationNftInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DelegationNft *DelegationNftFilterer) FilterInitialized(opts *bind.FilterOpts) (*DelegationNftInitializedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DelegationNftInitializedIterator{contract: _DelegationNft.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DelegationNft *DelegationNftFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DelegationNftInitialized) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftInitialized)
				if err := _DelegationNft.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseInitialized(log types.Log) (*DelegationNftInitialized, error) {
	event := new(DelegationNftInitialized)
	if err := _DelegationNft.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the DelegationNft contract.
type DelegationNftMetadataUpdateIterator struct {
	Event *DelegationNftMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *DelegationNftMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftMetadataUpdate)
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
		it.Event = new(DelegationNftMetadataUpdate)
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
func (it *DelegationNftMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftMetadataUpdate represents a MetadataUpdate event raised by the DelegationNft contract.
type DelegationNftMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_DelegationNft *DelegationNftFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*DelegationNftMetadataUpdateIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &DelegationNftMetadataUpdateIterator{contract: _DelegationNft.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_DelegationNft *DelegationNftFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *DelegationNftMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftMetadataUpdate)
				if err := _DelegationNft.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_DelegationNft *DelegationNftFilterer) ParseMetadataUpdate(log types.Log) (*DelegationNftMetadataUpdate, error) {
	event := new(DelegationNftMetadataUpdate)
	if err := _DelegationNft.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftReserveUpdatedIterator is returned from FilterReserveUpdated and is used to iterate over the raw logs and unpacked data for ReserveUpdated events raised by the DelegationNft contract.
type DelegationNftReserveUpdatedIterator struct {
	Event *DelegationNftReserveUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationNftReserveUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftReserveUpdated)
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
		it.Event = new(DelegationNftReserveUpdated)
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
func (it *DelegationNftReserveUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftReserveUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftReserveUpdated represents a ReserveUpdated event raised by the DelegationNft contract.
type DelegationNftReserveUpdated struct {
	TokenId *big.Int
	Reserve IDecimalNFTCommonReserve
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterReserveUpdated is a free log retrieval operation binding the contract event 0x4d3d24953206aaccc71602e460c62760cc53db2d397ad36fdb2578f2eddd3f7a.
//
// Solidity: event ReserveUpdated(uint256 tokenId, (address,uint256,uint8) reserve)
func (_DelegationNft *DelegationNftFilterer) FilterReserveUpdated(opts *bind.FilterOpts) (*DelegationNftReserveUpdatedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationNftReserveUpdatedIterator{contract: _DelegationNft.contract, event: "ReserveUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveUpdated is a free log subscription operation binding the contract event 0x4d3d24953206aaccc71602e460c62760cc53db2d397ad36fdb2578f2eddd3f7a.
//
// Solidity: event ReserveUpdated(uint256 tokenId, (address,uint256,uint8) reserve)
func (_DelegationNft *DelegationNftFilterer) WatchReserveUpdated(opts *bind.WatchOpts, sink chan<- *DelegationNftReserveUpdated) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftReserveUpdated)
				if err := _DelegationNft.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
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

// ParseReserveUpdated is a log parse operation binding the contract event 0x4d3d24953206aaccc71602e460c62760cc53db2d397ad36fdb2578f2eddd3f7a.
//
// Solidity: event ReserveUpdated(uint256 tokenId, (address,uint256,uint8) reserve)
func (_DelegationNft *DelegationNftFilterer) ParseReserveUpdated(log types.Log) (*DelegationNftReserveUpdated, error) {
	event := new(DelegationNftReserveUpdated)
	if err := _DelegationNft.contract.UnpackLog(event, "ReserveUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the DelegationNft contract.
type DelegationNftRoleAdminChangedIterator struct {
	Event *DelegationNftRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *DelegationNftRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftRoleAdminChanged)
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
		it.Event = new(DelegationNftRoleAdminChanged)
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
func (it *DelegationNftRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftRoleAdminChanged represents a RoleAdminChanged event raised by the DelegationNft contract.
type DelegationNftRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DelegationNft *DelegationNftFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*DelegationNftRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftRoleAdminChangedIterator{contract: _DelegationNft.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DelegationNft *DelegationNftFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *DelegationNftRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftRoleAdminChanged)
				if err := _DelegationNft.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DelegationNft *DelegationNftFilterer) ParseRoleAdminChanged(log types.Log) (*DelegationNftRoleAdminChanged, error) {
	event := new(DelegationNftRoleAdminChanged)
	if err := _DelegationNft.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the DelegationNft contract.
type DelegationNftRoleGrantedIterator struct {
	Event *DelegationNftRoleGranted // Event containing the contract specifics and raw log

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
func (it *DelegationNftRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftRoleGranted)
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
		it.Event = new(DelegationNftRoleGranted)
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
func (it *DelegationNftRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftRoleGranted represents a RoleGranted event raised by the DelegationNft contract.
type DelegationNftRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DelegationNftRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftRoleGrantedIterator{contract: _DelegationNft.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *DelegationNftRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftRoleGranted)
				if err := _DelegationNft.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) ParseRoleGranted(log types.Log) (*DelegationNftRoleGranted, error) {
	event := new(DelegationNftRoleGranted)
	if err := _DelegationNft.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the DelegationNft contract.
type DelegationNftRoleRevokedIterator struct {
	Event *DelegationNftRoleRevoked // Event containing the contract specifics and raw log

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
func (it *DelegationNftRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftRoleRevoked)
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
		it.Event = new(DelegationNftRoleRevoked)
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
func (it *DelegationNftRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftRoleRevoked represents a RoleRevoked event raised by the DelegationNft contract.
type DelegationNftRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DelegationNftRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftRoleRevokedIterator{contract: _DelegationNft.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *DelegationNftRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftRoleRevoked)
				if err := _DelegationNft.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DelegationNft *DelegationNftFilterer) ParseRoleRevoked(log types.Log) (*DelegationNftRoleRevoked, error) {
	event := new(DelegationNftRoleRevoked)
	if err := _DelegationNft.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the DelegationNft contract.
type DelegationNftTransferIterator struct {
	Event *DelegationNftTransfer // Event containing the contract specifics and raw log

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
func (it *DelegationNftTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftTransfer)
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
		it.Event = new(DelegationNftTransfer)
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
func (it *DelegationNftTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftTransfer represents a Transfer event raised by the DelegationNft contract.
type DelegationNftTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*DelegationNftTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransferIterator{contract: _DelegationNft.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *DelegationNftTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftTransfer)
				if err := _DelegationNft.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_DelegationNft *DelegationNftFilterer) ParseTransfer(log types.Log) (*DelegationNftTransfer, error) {
	event := new(DelegationNftTransfer)
	if err := _DelegationNft.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the DelegationNft contract.
type DelegationNftURIIterator struct {
	Event *DelegationNftURI // Event containing the contract specifics and raw log

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
func (it *DelegationNftURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftURI)
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
		it.Event = new(DelegationNftURI)
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
func (it *DelegationNftURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftURI represents a URI event raised by the DelegationNft contract.
type DelegationNftURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_DelegationNft *DelegationNftFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*DelegationNftURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftURIIterator{contract: _DelegationNft.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_DelegationNft *DelegationNftFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *DelegationNftURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftURI)
				if err := _DelegationNft.contract.UnpackLog(event, "URI", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseURI(log types.Log) (*DelegationNftURI, error) {
	event := new(DelegationNftURI)
	if err := _DelegationNft.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
