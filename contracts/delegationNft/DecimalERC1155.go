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
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC1155InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"idsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valuesLength\",\"type\":\"uint256\"}],\"name\":\"ERC1155InvalidArrayLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC1155MissingApprovalForAll\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMinReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserve\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidReserveType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MintNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNftDelegation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitInvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PermitUnauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReserveAlreadyInitialized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DisabledMint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"indexed\":false,\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"reserve\",\"type\":\"tuple\"}],\"name\":\"ReserveUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CREATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_RESERVE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"addedReserveAmount\",\"type\":\"uint256\"}],\"name\":\"addReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"addReserveByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"addReserveByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowMint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNftType\",\"outputs\":[{\"internalType\":\"enumNFTType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRefundable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalNFTCommon.ReserveType\",\"name\":\"reserveType\",\"type\":\"uint8\"}],\"internalType\":\"structIDecimalNFTCommon.Reserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initialSymbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"initialBaseURI\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"initialCreator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"initialRefundable\",\"type\":\"bool\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"mintByETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToMint\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"reserveAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"reserveToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"mintByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToPenalty\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"percentage\",\"type\":\"uint16\"}],\"name\":\"penalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"rate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI\",\"type\":\"string\"}],\"name\":\"setTokenURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.BalanceOf(&_DelegationNft.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.BalanceOf(&_DelegationNft.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_DelegationNft *DelegationNftCaller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_DelegationNft *DelegationNftSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _DelegationNft.Contract.BalanceOfBatch(&_DelegationNft.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_DelegationNft *DelegationNftCallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _DelegationNft.Contract.BalanceOfBatch(&_DelegationNft.CallOpts, accounts, ids)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_DelegationNft *DelegationNftCaller) Exists(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "exists", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_DelegationNft *DelegationNftSession) Exists(id *big.Int) (bool, error) {
	return _DelegationNft.Contract.Exists(&_DelegationNft.CallOpts, id)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 id) view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) Exists(id *big.Int) (bool, error) {
	return _DelegationNft.Contract.Exists(&_DelegationNft.CallOpts, id)
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
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_DelegationNft *DelegationNftCaller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_DelegationNft *DelegationNftSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _DelegationNft.Contract.IsApprovedForAll(&_DelegationNft.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_DelegationNft *DelegationNftCallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _DelegationNft.Contract.IsApprovedForAll(&_DelegationNft.CallOpts, account, operator)
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

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftSession) Nonces(owner common.Address) (*big.Int, error) {
	return _DelegationNft.Contract.Nonces(&_DelegationNft.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _DelegationNft.Contract.Nonces(&_DelegationNft.CallOpts, owner)
}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_DelegationNft *DelegationNftCaller) Rate(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "rate", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_DelegationNft *DelegationNftSession) Rate(tokenId *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.Rate(&_DelegationNft.CallOpts, tokenId)
}

// Rate is a free data retrieval call binding the contract method 0xe7ee6ad6.
//
// Solidity: function rate(uint256 tokenId) view returns(uint256 amount)
func (_DelegationNft *DelegationNftCallerSession) Rate(tokenId *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.Rate(&_DelegationNft.CallOpts, tokenId)
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

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_DelegationNft *DelegationNftCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_DelegationNft *DelegationNftSession) TotalSupply() (*big.Int, error) {
	return _DelegationNft.Contract.TotalSupply(&_DelegationNft.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) TotalSupply() (*big.Int, error) {
	return _DelegationNft.Contract.TotalSupply(&_DelegationNft.CallOpts)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) TotalSupply0(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "totalSupply0", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.TotalSupply0(&_DelegationNft.CallOpts, id)
}

// TotalSupply0 is a free data retrieval call binding the contract method 0xbd85b039.
//
// Solidity: function totalSupply(uint256 id) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) TotalSupply0(id *big.Int) (*big.Int, error) {
	return _DelegationNft.Contract.TotalSupply0(&_DelegationNft.CallOpts, id)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftCaller) Uri(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "uri", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftSession) Uri(tokenId *big.Int) (string, error) {
	return _DelegationNft.Contract.Uri(&_DelegationNft.CallOpts, tokenId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 tokenId) view returns(string)
func (_DelegationNft *DelegationNftCallerSession) Uri(tokenId *big.Int) (string, error) {
	return _DelegationNft.Contract.Uri(&_DelegationNft.CallOpts, tokenId)
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

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "burn", tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftSession) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Burn(&_DelegationNft.TransactOpts, tokenId, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xb390c0ab.
//
// Solidity: function burn(uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactorSession) Burn(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Burn(&_DelegationNft.TransactOpts, tokenId, amount)
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

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_DelegationNft *DelegationNftTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mint", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_DelegationNft *DelegationNftSession) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Mint(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// Mint is a paid mutator transaction binding the contract method 0x89278622.
//
// Solidity: function mint(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken) returns()
func (_DelegationNft *DelegationNftTransactorSession) Mint(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Mint(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_DelegationNft *DelegationNftTransactor) MintByETH(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mintByETH", recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_DelegationNft *DelegationNftSession) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByETH(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByETH is a paid mutator transaction binding the contract method 0xa504cf16.
//
// Solidity: function mintByETH(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI) payable returns()
func (_DelegationNft *DelegationNftTransactorSession) MintByETH(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByETH(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) MintByPermit(opts *bind.TransactOpts, recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "mintByPermit", recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByPermit(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// MintByPermit is a paid mutator transaction binding the contract method 0x0f071717.
//
// Solidity: function mintByPermit(address recipient, uint256 tokenId, uint256 amountToMint, string tokenURI, uint256 reserveAmount, address reserveToken, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) MintByPermit(recipient common.Address, tokenId *big.Int, amountToMint *big.Int, tokenURI string, reserveAmount *big.Int, reserveToken common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.MintByPermit(&_DelegationNft.TransactOpts, recipient, tokenId, amountToMint, tokenURI, reserveAmount, reserveToken, deadline, v, r, s)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_DelegationNft *DelegationNftTransactor) Penalty(opts *bind.TransactOpts, tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "penalty", tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_DelegationNft *DelegationNftSession) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.Contract.Penalty(&_DelegationNft.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Penalty is a paid mutator transaction binding the contract method 0xc57ff62e.
//
// Solidity: function penalty(uint256 tokenId, uint256 amountToPenalty, uint16 percentage) returns()
func (_DelegationNft *DelegationNftTransactorSession) Penalty(tokenId *big.Int, amountToPenalty *big.Int, percentage uint16) (*types.Transaction, error) {
	return _DelegationNft.Contract.Penalty(&_DelegationNft.TransactOpts, tokenId, amountToPenalty, percentage)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "permit", owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Permit(&_DelegationNft.TransactOpts, owner, spender, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0x48613c28.
//
// Solidity: function permit(address owner, address spender, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) Permit(owner common.Address, spender common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Permit(&_DelegationNft.TransactOpts, owner, spender, deadline, v, r, s)
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

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_DelegationNft *DelegationNftTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_DelegationNft *DelegationNftSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeBatchTransferFrom(&_DelegationNft.TransactOpts, from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_DelegationNft *DelegationNftTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeBatchTransferFrom(&_DelegationNft.TransactOpts, from, to, ids, values, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_DelegationNft *DelegationNftTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "safeTransferFrom", from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_DelegationNft *DelegationNftSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom(&_DelegationNft.TransactOpts, from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_DelegationNft *DelegationNftTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.SafeTransferFrom(&_DelegationNft.TransactOpts, from, to, id, value, data)
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
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*DelegationNftApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftApprovalForAllIterator{contract: _DelegationNft.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *DelegationNftApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
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
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_DelegationNft *DelegationNftFilterer) ParseApprovalForAll(log types.Log) (*DelegationNftApprovalForAll, error) {
	event := new(DelegationNftApprovalForAll)
	if err := _DelegationNft.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
	TokenId     *big.Int
	TotalSupply *big.Int
	Reserve     IDecimalNFTCommonReserve
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReserveUpdated is a free log retrieval operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
func (_DelegationNft *DelegationNftFilterer) FilterReserveUpdated(opts *bind.FilterOpts) (*DelegationNftReserveUpdatedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "ReserveUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationNftReserveUpdatedIterator{contract: _DelegationNft.contract, event: "ReserveUpdated", logs: logs, sub: sub}, nil
}

// WatchReserveUpdated is a free log subscription operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
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

// ParseReserveUpdated is a log parse operation binding the contract event 0x416c94fe34624b6660ef8d22d507994befd5eee563a60424df0bc5a7e51262d7.
//
// Solidity: event ReserveUpdated(uint256 tokenId, uint256 totalSupply, (address,uint256,uint8) reserve)
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

// DelegationNftTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the DelegationNft contract.
type DelegationNftTransferBatchIterator struct {
	Event *DelegationNftTransferBatch // Event containing the contract specifics and raw log

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
func (it *DelegationNftTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftTransferBatch)
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
		it.Event = new(DelegationNftTransferBatch)
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
func (it *DelegationNftTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftTransferBatch represents a TransferBatch event raised by the DelegationNft contract.
type DelegationNftTransferBatch struct {
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
func (_DelegationNft *DelegationNftFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*DelegationNftTransferBatchIterator, error) {

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

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransferBatchIterator{contract: _DelegationNft.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_DelegationNft *DelegationNftFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *DelegationNftTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftTransferBatch)
				if err := _DelegationNft.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseTransferBatch(log types.Log) (*DelegationNftTransferBatch, error) {
	event := new(DelegationNftTransferBatch)
	if err := _DelegationNft.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the DelegationNft contract.
type DelegationNftTransferSingleIterator struct {
	Event *DelegationNftTransferSingle // Event containing the contract specifics and raw log

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
func (it *DelegationNftTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftTransferSingle)
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
		it.Event = new(DelegationNftTransferSingle)
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
func (it *DelegationNftTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftTransferSingle represents a TransferSingle event raised by the DelegationNft contract.
type DelegationNftTransferSingle struct {
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
func (_DelegationNft *DelegationNftFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*DelegationNftTransferSingleIterator, error) {

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

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransferSingleIterator{contract: _DelegationNft.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_DelegationNft *DelegationNftFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *DelegationNftTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftTransferSingle)
				if err := _DelegationNft.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseTransferSingle(log types.Log) (*DelegationNftTransferSingle, error) {
	event := new(DelegationNftTransferSingle)
	if err := _DelegationNft.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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
