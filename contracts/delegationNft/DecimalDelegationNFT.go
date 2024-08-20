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

// IDecimalDelegationCommonFrozenStake is an auto generated low-level Go binding around an user-defined struct.
type IDecimalDelegationCommonFrozenStake struct {
	Stake             IDecimalDelegationCommonStake
	FreezeStatus      uint8
	FreezeType        uint8
	UnfreezeTimestamp *big.Int
}

// IDecimalDelegationCommonStake is an auto generated low-level Go binding around an user-defined struct.
type IDecimalDelegationCommonStake struct {
	Validator     common.Address
	Delegator     common.Address
	Token         common.Address
	Amount        *big.Int
	TokenId       *big.Int
	TokenType     uint8
	HoldTimestamp *big.Int
}

// IDecimalDelegationCommonValidatorReserve is an auto generated low-level Go binding around an user-defined struct.
type IDecimalDelegationCommonValidatorReserve struct {
	PenaltyIndex *big.Int
	Reserve      *big.Int
}

// IDecimalDelegationCommonValidatorToken is an auto generated low-level Go binding around an user-defined struct.
type IDecimalDelegationCommonValidatorToken struct {
	Validator common.Address
	Token     common.Address
	TokenId   *big.Int
}

// DelegationNftMetaData contains all meta data concerning the DelegationNft contract.
var DelegationNftMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FrozenStakesQueueIsEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCollection\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFrozenType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTimestamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTokenType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValidator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCompletableFrozenStakes\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPenaltyToApply\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotDRC1155\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotDRC721\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeAlreadyUnfrozen\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeInactive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimestampError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawFreezeTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transferFreezeTime\",\"type\":\"uint256\"}],\"name\":\"FreezeTimeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newPenaltyIndex\",\"type\":\"uint256\"}],\"name\":\"PenaltyAppliedToStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"penaltyIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.ValidatorReserve\",\"name\":\"validatorReserve\",\"type\":\"tuple\"}],\"name\":\"PenaltyAppliedToValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"changedAmount\",\"type\":\"int256\"}],\"name\":\"StakeAmountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isNew\",\"type\":\"bool\"}],\"name\":\"StakeHolded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isNew\",\"type\":\"bool\"}],\"name\":\"StakeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"}],\"name\":\"TransferCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"frozenStake\",\"type\":\"tuple\"}],\"name\":\"TransferRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"}],\"name\":\"WithdrawCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"frozenStake\",\"type\":\"tuple\"}],\"name\":\"WithdrawRequest\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.ValidatorToken[]\",\"name\":\"validatorTokens\",\"type\":\"tuple[]\"}],\"name\":\"applyPenaltiesToValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"indexes\",\"type\":\"uint256[]\"}],\"name\":\"complete\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegateDRC1155\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateDRC1155ByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"delegateDRC721\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateDRC721ByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"delegateHoldDRC1155\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateHoldDRC1155ByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"delegateHoldDRC721\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateHoldDRC721ByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"}],\"name\":\"getFreezeTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getFrozenStake\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"stakeIndexes\",\"type\":\"uint256[]\"}],\"name\":\"getFrozenStakes\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.FrozenStake[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"getHoldStake\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"getHoldStakeId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImpl\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getStake\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getStakeId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"}],\"name\":\"getStakePenaltyIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stakeIds\",\"type\":\"bytes32[]\"}],\"name\":\"getStakes\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake[]\",\"name\":\"stakes\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getValidatorReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"penaltyIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.ValidatorReserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountToHold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldHoldTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newHoldTimestamp\",\"type\":\"uint256\"}],\"name\":\"hold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155BatchReceived\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"withdrawFreezeTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transferFreezeTime\",\"type\":\"uint256\"}],\"name\":\"setFreezeTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"transferHold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImpl\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nft\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"withdrawHold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_DelegationNft *DelegationNftCaller) GetContractCenter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getContractCenter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_DelegationNft *DelegationNftSession) GetContractCenter() (common.Address, error) {
	return _DelegationNft.Contract.GetContractCenter(&_DelegationNft.CallOpts)
}

// GetContractCenter is a free data retrieval call binding the contract method 0xba778bce.
//
// Solidity: function getContractCenter() view returns(address)
func (_DelegationNft *DelegationNftCallerSession) GetContractCenter() (common.Address, error) {
	return _DelegationNft.Contract.GetContractCenter(&_DelegationNft.CallOpts)
}

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) GetFreezeTime(opts *bind.CallOpts, freezeType uint8) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getFreezeTime", freezeType)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_DelegationNft *DelegationNftSession) GetFreezeTime(freezeType uint8) (*big.Int, error) {
	return _DelegationNft.Contract.GetFreezeTime(&_DelegationNft.CallOpts, freezeType)
}

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) GetFreezeTime(freezeType uint8) (*big.Int, error) {
	return _DelegationNft.Contract.GetFreezeTime(&_DelegationNft.CallOpts, freezeType)
}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_DelegationNft *DelegationNftCaller) GetFrozenStake(opts *bind.CallOpts, index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getFrozenStake", index)

	if err != nil {
		return *new(IDecimalDelegationCommonFrozenStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonFrozenStake)).(*IDecimalDelegationCommonFrozenStake)

	return out0, err

}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_DelegationNft *DelegationNftSession) GetFrozenStake(index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	return _DelegationNft.Contract.GetFrozenStake(&_DelegationNft.CallOpts, index)
}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_DelegationNft *DelegationNftCallerSession) GetFrozenStake(index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	return _DelegationNft.Contract.GetFrozenStake(&_DelegationNft.CallOpts, index)
}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_DelegationNft *DelegationNftCaller) GetFrozenStakes(opts *bind.CallOpts, stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getFrozenStakes", stakeIndexes)

	if err != nil {
		return *new([]IDecimalDelegationCommonFrozenStake), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDecimalDelegationCommonFrozenStake)).(*[]IDecimalDelegationCommonFrozenStake)

	return out0, err

}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_DelegationNft *DelegationNftSession) GetFrozenStakes(stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	return _DelegationNft.Contract.GetFrozenStakes(&_DelegationNft.CallOpts, stakeIndexes)
}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_DelegationNft *DelegationNftCallerSession) GetFrozenStakes(stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	return _DelegationNft.Contract.GetFrozenStakes(&_DelegationNft.CallOpts, stakeIndexes)
}

// GetHoldStake is a free data retrieval call binding the contract method 0x8dcf9a3b.
//
// Solidity: function getHoldStake(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftCaller) GetHoldStake(opts *bind.CallOpts, validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getHoldStake", validator, delegator, nft, tokenId, holdTimestamp)

	if err != nil {
		return *new(IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonStake)).(*IDecimalDelegationCommonStake)

	return out0, err

}

// GetHoldStake is a free data retrieval call binding the contract method 0x8dcf9a3b.
//
// Solidity: function getHoldStake(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftSession) GetHoldStake(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetHoldStake(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId, holdTimestamp)
}

// GetHoldStake is a free data retrieval call binding the contract method 0x8dcf9a3b.
//
// Solidity: function getHoldStake(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftCallerSession) GetHoldStake(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetHoldStake(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId, holdTimestamp)
}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x304db495.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) pure returns(bytes32)
func (_DelegationNft *DelegationNftCaller) GetHoldStakeId(opts *bind.CallOpts, validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getHoldStakeId", validator, delegator, nft, tokenId, holdTimestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x304db495.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) pure returns(bytes32)
func (_DelegationNft *DelegationNftSession) GetHoldStakeId(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) ([32]byte, error) {
	return _DelegationNft.Contract.GetHoldStakeId(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId, holdTimestamp)
}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x304db495.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address nft, uint256 tokenId, uint256 holdTimestamp) pure returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) GetHoldStakeId(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) ([32]byte, error) {
	return _DelegationNft.Contract.GetHoldStakeId(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId, holdTimestamp)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_DelegationNft *DelegationNftCaller) GetImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_DelegationNft *DelegationNftSession) GetImpl() (common.Address, error) {
	return _DelegationNft.Contract.GetImpl(&_DelegationNft.CallOpts)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_DelegationNft *DelegationNftCallerSession) GetImpl() (common.Address, error) {
	return _DelegationNft.Contract.GetImpl(&_DelegationNft.CallOpts)
}

// GetStake is a free data retrieval call binding the contract method 0x347406a2.
//
// Solidity: function getStake(address validator, address delegator, address nft, uint256 tokenId) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftCaller) GetStake(opts *bind.CallOpts, validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) (IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getStake", validator, delegator, nft, tokenId)

	if err != nil {
		return *new(IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonStake)).(*IDecimalDelegationCommonStake)

	return out0, err

}

// GetStake is a free data retrieval call binding the contract method 0x347406a2.
//
// Solidity: function getStake(address validator, address delegator, address nft, uint256 tokenId) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftSession) GetStake(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) (IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetStake(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId)
}

// GetStake is a free data retrieval call binding the contract method 0x347406a2.
//
// Solidity: function getStake(address validator, address delegator, address nft, uint256 tokenId) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_DelegationNft *DelegationNftCallerSession) GetStake(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) (IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetStake(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId)
}

// GetStakeId is a free data retrieval call binding the contract method 0xe93f57c6.
//
// Solidity: function getStakeId(address validator, address delegator, address nft, uint256 tokenId) pure returns(bytes32)
func (_DelegationNft *DelegationNftCaller) GetStakeId(opts *bind.CallOpts, validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getStakeId", validator, delegator, nft, tokenId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetStakeId is a free data retrieval call binding the contract method 0xe93f57c6.
//
// Solidity: function getStakeId(address validator, address delegator, address nft, uint256 tokenId) pure returns(bytes32)
func (_DelegationNft *DelegationNftSession) GetStakeId(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) ([32]byte, error) {
	return _DelegationNft.Contract.GetStakeId(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId)
}

// GetStakeId is a free data retrieval call binding the contract method 0xe93f57c6.
//
// Solidity: function getStakeId(address validator, address delegator, address nft, uint256 tokenId) pure returns(bytes32)
func (_DelegationNft *DelegationNftCallerSession) GetStakeId(validator common.Address, delegator common.Address, nft common.Address, tokenId *big.Int) ([32]byte, error) {
	return _DelegationNft.Contract.GetStakeId(&_DelegationNft.CallOpts, validator, delegator, nft, tokenId)
}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_DelegationNft *DelegationNftCaller) GetStakePenaltyIndex(opts *bind.CallOpts, stakeId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getStakePenaltyIndex", stakeId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_DelegationNft *DelegationNftSession) GetStakePenaltyIndex(stakeId [32]byte) (*big.Int, error) {
	return _DelegationNft.Contract.GetStakePenaltyIndex(&_DelegationNft.CallOpts, stakeId)
}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_DelegationNft *DelegationNftCallerSession) GetStakePenaltyIndex(stakeId [32]byte) (*big.Int, error) {
	return _DelegationNft.Contract.GetStakePenaltyIndex(&_DelegationNft.CallOpts, stakeId)
}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_DelegationNft *DelegationNftCaller) GetStakes(opts *bind.CallOpts, stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getStakes", stakeIds)

	if err != nil {
		return *new([]IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDecimalDelegationCommonStake)).(*[]IDecimalDelegationCommonStake)

	return out0, err

}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_DelegationNft *DelegationNftSession) GetStakes(stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetStakes(&_DelegationNft.CallOpts, stakeIds)
}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_DelegationNft *DelegationNftCallerSession) GetStakes(stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	return _DelegationNft.Contract.GetStakes(&_DelegationNft.CallOpts, stakeIds)
}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x20af2fe5.
//
// Solidity: function getValidatorReserve(address validator, address token, uint256 tokenId) view returns((uint256,uint256))
func (_DelegationNft *DelegationNftCaller) GetValidatorReserve(opts *bind.CallOpts, validator common.Address, token common.Address, tokenId *big.Int) (IDecimalDelegationCommonValidatorReserve, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "getValidatorReserve", validator, token, tokenId)

	if err != nil {
		return *new(IDecimalDelegationCommonValidatorReserve), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonValidatorReserve)).(*IDecimalDelegationCommonValidatorReserve)

	return out0, err

}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x20af2fe5.
//
// Solidity: function getValidatorReserve(address validator, address token, uint256 tokenId) view returns((uint256,uint256))
func (_DelegationNft *DelegationNftSession) GetValidatorReserve(validator common.Address, token common.Address, tokenId *big.Int) (IDecimalDelegationCommonValidatorReserve, error) {
	return _DelegationNft.Contract.GetValidatorReserve(&_DelegationNft.CallOpts, validator, token, tokenId)
}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x20af2fe5.
//
// Solidity: function getValidatorReserve(address validator, address token, uint256 tokenId) view returns((uint256,uint256))
func (_DelegationNft *DelegationNftCallerSession) GetValidatorReserve(validator common.Address, token common.Address, tokenId *big.Int) (IDecimalDelegationCommonValidatorReserve, error) {
	return _DelegationNft.Contract.GetValidatorReserve(&_DelegationNft.CallOpts, validator, token, tokenId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DelegationNft *DelegationNftCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegationNft.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DelegationNft *DelegationNftSession) Owner() (common.Address, error) {
	return _DelegationNft.Contract.Owner(&_DelegationNft.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DelegationNft *DelegationNftCallerSession) Owner() (common.Address, error) {
	return _DelegationNft.Contract.Owner(&_DelegationNft.CallOpts)
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

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_DelegationNft *DelegationNftTransactor) ApplyPenaltiesToValidator(opts *bind.TransactOpts, validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "applyPenaltiesToValidator", validatorTokens)
}

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_DelegationNft *DelegationNftSession) ApplyPenaltiesToValidator(validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _DelegationNft.Contract.ApplyPenaltiesToValidator(&_DelegationNft.TransactOpts, validatorTokens)
}

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_DelegationNft *DelegationNftTransactorSession) ApplyPenaltiesToValidator(validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _DelegationNft.Contract.ApplyPenaltiesToValidator(&_DelegationNft.TransactOpts, validatorTokens)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_DelegationNft *DelegationNftTransactor) Complete(opts *bind.TransactOpts, indexes []*big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "complete", indexes)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_DelegationNft *DelegationNftSession) Complete(indexes []*big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Complete(&_DelegationNft.TransactOpts, indexes)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_DelegationNft *DelegationNftTransactorSession) Complete(indexes []*big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Complete(&_DelegationNft.TransactOpts, indexes)
}

// DelegateDRC1155 is a paid mutator transaction binding the contract method 0x1caa077d.
//
// Solidity: function delegateDRC1155(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateDRC1155(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateDRC1155", validator, nft, tokenId, amount)
}

// DelegateDRC1155 is a paid mutator transaction binding the contract method 0x1caa077d.
//
// Solidity: function delegateDRC1155(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftSession) DelegateDRC1155(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC1155(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount)
}

// DelegateDRC1155 is a paid mutator transaction binding the contract method 0x1caa077d.
//
// Solidity: function delegateDRC1155(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateDRC1155(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC1155(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount)
}

// DelegateDRC1155ByPermit is a paid mutator transaction binding the contract method 0x019a2159.
//
// Solidity: function delegateDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateDRC1155ByPermit(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateDRC1155ByPermit", validator, nft, tokenId, amount, deadline, v, r, s)
}

// DelegateDRC1155ByPermit is a paid mutator transaction binding the contract method 0x019a2159.
//
// Solidity: function delegateDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) DelegateDRC1155ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC1155ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, deadline, v, r, s)
}

// DelegateDRC1155ByPermit is a paid mutator transaction binding the contract method 0x019a2159.
//
// Solidity: function delegateDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateDRC1155ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC1155ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, deadline, v, r, s)
}

// DelegateDRC721 is a paid mutator transaction binding the contract method 0x556e055f.
//
// Solidity: function delegateDRC721(address validator, address nft, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateDRC721(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateDRC721", validator, nft, tokenId)
}

// DelegateDRC721 is a paid mutator transaction binding the contract method 0x556e055f.
//
// Solidity: function delegateDRC721(address validator, address nft, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftSession) DelegateDRC721(validator common.Address, nft common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC721(&_DelegationNft.TransactOpts, validator, nft, tokenId)
}

// DelegateDRC721 is a paid mutator transaction binding the contract method 0x556e055f.
//
// Solidity: function delegateDRC721(address validator, address nft, uint256 tokenId) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateDRC721(validator common.Address, nft common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC721(&_DelegationNft.TransactOpts, validator, nft, tokenId)
}

// DelegateDRC721ByPermit is a paid mutator transaction binding the contract method 0x8bbf0a25.
//
// Solidity: function delegateDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateDRC721ByPermit(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateDRC721ByPermit", validator, nft, tokenId, deadline, v, r, s)
}

// DelegateDRC721ByPermit is a paid mutator transaction binding the contract method 0x8bbf0a25.
//
// Solidity: function delegateDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) DelegateDRC721ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC721ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, deadline, v, r, s)
}

// DelegateDRC721ByPermit is a paid mutator transaction binding the contract method 0x8bbf0a25.
//
// Solidity: function delegateDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateDRC721ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateDRC721ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, deadline, v, r, s)
}

// DelegateHoldDRC1155 is a paid mutator transaction binding the contract method 0xd55c3e08.
//
// Solidity: function delegateHoldDRC1155(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateHoldDRC1155(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateHoldDRC1155", validator, nft, tokenId, amount, holdTimestamp)
}

// DelegateHoldDRC1155 is a paid mutator transaction binding the contract method 0xd55c3e08.
//
// Solidity: function delegateHoldDRC1155(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftSession) DelegateHoldDRC1155(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC1155(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp)
}

// DelegateHoldDRC1155 is a paid mutator transaction binding the contract method 0xd55c3e08.
//
// Solidity: function delegateHoldDRC1155(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateHoldDRC1155(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC1155(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp)
}

// DelegateHoldDRC1155ByPermit is a paid mutator transaction binding the contract method 0x57420901.
//
// Solidity: function delegateHoldDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateHoldDRC1155ByPermit(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateHoldDRC1155ByPermit", validator, nft, tokenId, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDRC1155ByPermit is a paid mutator transaction binding the contract method 0x57420901.
//
// Solidity: function delegateHoldDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) DelegateHoldDRC1155ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC1155ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDRC1155ByPermit is a paid mutator transaction binding the contract method 0x57420901.
//
// Solidity: function delegateHoldDRC1155ByPermit(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateHoldDRC1155ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC1155ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDRC721 is a paid mutator transaction binding the contract method 0x3ae02675.
//
// Solidity: function delegateHoldDRC721(address validator, address nft, uint256 tokenId, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateHoldDRC721(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateHoldDRC721", validator, nft, tokenId, holdTimestamp)
}

// DelegateHoldDRC721 is a paid mutator transaction binding the contract method 0x3ae02675.
//
// Solidity: function delegateHoldDRC721(address validator, address nft, uint256 tokenId, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftSession) DelegateHoldDRC721(validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC721(&_DelegationNft.TransactOpts, validator, nft, tokenId, holdTimestamp)
}

// DelegateHoldDRC721 is a paid mutator transaction binding the contract method 0x3ae02675.
//
// Solidity: function delegateHoldDRC721(address validator, address nft, uint256 tokenId, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateHoldDRC721(validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC721(&_DelegationNft.TransactOpts, validator, nft, tokenId, holdTimestamp)
}

// DelegateHoldDRC721ByPermit is a paid mutator transaction binding the contract method 0xac029ddb.
//
// Solidity: function delegateHoldDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactor) DelegateHoldDRC721ByPermit(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "delegateHoldDRC721ByPermit", validator, nft, tokenId, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDRC721ByPermit is a paid mutator transaction binding the contract method 0xac029ddb.
//
// Solidity: function delegateHoldDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftSession) DelegateHoldDRC721ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC721ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDRC721ByPermit is a paid mutator transaction binding the contract method 0xac029ddb.
//
// Solidity: function delegateHoldDRC721ByPermit(address validator, address nft, uint256 tokenId, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_DelegationNft *DelegationNftTransactorSession) DelegateHoldDRC721ByPermit(validator common.Address, nft common.Address, tokenId *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.DelegateHoldDRC721ByPermit(&_DelegationNft.TransactOpts, validator, nft, tokenId, holdTimestamp, deadline, v, r, s)
}

// Hold is a paid mutator transaction binding the contract method 0x0029ae70.
//
// Solidity: function hold(address validator, address token, uint256 tokenId, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_DelegationNft *DelegationNftTransactor) Hold(opts *bind.TransactOpts, validator common.Address, token common.Address, tokenId *big.Int, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "hold", validator, token, tokenId, amountToHold, oldHoldTimestamp, newHoldTimestamp)
}

// Hold is a paid mutator transaction binding the contract method 0x0029ae70.
//
// Solidity: function hold(address validator, address token, uint256 tokenId, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_DelegationNft *DelegationNftSession) Hold(validator common.Address, token common.Address, tokenId *big.Int, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Hold(&_DelegationNft.TransactOpts, validator, token, tokenId, amountToHold, oldHoldTimestamp, newHoldTimestamp)
}

// Hold is a paid mutator transaction binding the contract method 0x0029ae70.
//
// Solidity: function hold(address validator, address token, uint256 tokenId, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_DelegationNft *DelegationNftTransactorSession) Hold(validator common.Address, token common.Address, tokenId *big.Int, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Hold(&_DelegationNft.TransactOpts, validator, token, tokenId, amountToHold, oldHoldTimestamp, newHoldTimestamp)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC1155BatchReceived(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC1155BatchReceived(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC1155Received(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC1155Received(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC721Received(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_DelegationNft *DelegationNftTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.OnERC721Received(&_DelegationNft.TransactOpts, arg0, arg1, arg2, arg3)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelegationNft *DelegationNftTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelegationNft *DelegationNftSession) RenounceOwnership() (*types.Transaction, error) {
	return _DelegationNft.Contract.RenounceOwnership(&_DelegationNft.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DelegationNft *DelegationNftTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DelegationNft.Contract.RenounceOwnership(&_DelegationNft.TransactOpts)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address contractCenter) returns()
func (_DelegationNft *DelegationNftTransactor) SetContractCenter(opts *bind.TransactOpts, contractCenter common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "setContractCenter", contractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address contractCenter) returns()
func (_DelegationNft *DelegationNftSession) SetContractCenter(contractCenter common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetContractCenter(&_DelegationNft.TransactOpts, contractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address contractCenter) returns()
func (_DelegationNft *DelegationNftTransactorSession) SetContractCenter(contractCenter common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetContractCenter(&_DelegationNft.TransactOpts, contractCenter)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_DelegationNft *DelegationNftTransactor) SetFreezeTime(opts *bind.TransactOpts, withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "setFreezeTime", withdrawFreezeTime, transferFreezeTime)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_DelegationNft *DelegationNftSession) SetFreezeTime(withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetFreezeTime(&_DelegationNft.TransactOpts, withdrawFreezeTime, transferFreezeTime)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_DelegationNft *DelegationNftTransactorSession) SetFreezeTime(withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.SetFreezeTime(&_DelegationNft.TransactOpts, withdrawFreezeTime, transferFreezeTime)
}

// Transfer is a paid mutator transaction binding the contract method 0x47338bc3.
//
// Solidity: function transfer(address validator, address nft, uint256 tokenId, uint256 amount, address newValidator) returns()
func (_DelegationNft *DelegationNftTransactor) Transfer(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "transfer", validator, nft, tokenId, amount, newValidator)
}

// Transfer is a paid mutator transaction binding the contract method 0x47338bc3.
//
// Solidity: function transfer(address validator, address nft, uint256 tokenId, uint256 amount, address newValidator) returns()
func (_DelegationNft *DelegationNftSession) Transfer(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Transfer(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, newValidator)
}

// Transfer is a paid mutator transaction binding the contract method 0x47338bc3.
//
// Solidity: function transfer(address validator, address nft, uint256 tokenId, uint256 amount, address newValidator) returns()
func (_DelegationNft *DelegationNftTransactorSession) Transfer(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.Transfer(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x3547805a.
//
// Solidity: function transferHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_DelegationNft *DelegationNftTransactor) TransferHold(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "transferHold", validator, nft, tokenId, amount, holdTimestamp, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x3547805a.
//
// Solidity: function transferHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_DelegationNft *DelegationNftSession) TransferHold(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferHold(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x3547805a.
//
// Solidity: function transferHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_DelegationNft *DelegationNftTransactorSession) TransferHold(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferHold(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp, newValidator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelegationNft *DelegationNftTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelegationNft *DelegationNftSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferOwnership(&_DelegationNft.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DelegationNft *DelegationNftTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DelegationNft.Contract.TransferOwnership(&_DelegationNft.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_DelegationNft *DelegationNftTransactor) Upgrade(opts *bind.TransactOpts, newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "upgrade", newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_DelegationNft *DelegationNftSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Upgrade(&_DelegationNft.TransactOpts, newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_DelegationNft *DelegationNftTransactorSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _DelegationNft.Contract.Upgrade(&_DelegationNft.TransactOpts, newImpl, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactor) Withdraw(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "withdraw", validator, nft, tokenId, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftSession) Withdraw(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Withdraw(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7bfe950c.
//
// Solidity: function withdraw(address validator, address nft, uint256 tokenId, uint256 amount) returns()
func (_DelegationNft *DelegationNftTransactorSession) Withdraw(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.Withdraw(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0xc52c4025.
//
// Solidity: function withdrawHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactor) WithdrawHold(opts *bind.TransactOpts, validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.contract.Transact(opts, "withdrawHold", validator, nft, tokenId, amount, holdTimestamp)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0xc52c4025.
//
// Solidity: function withdrawHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftSession) WithdrawHold(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.WithdrawHold(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0xc52c4025.
//
// Solidity: function withdrawHold(address validator, address nft, uint256 tokenId, uint256 amount, uint256 holdTimestamp) returns()
func (_DelegationNft *DelegationNftTransactorSession) WithdrawHold(validator common.Address, nft common.Address, tokenId *big.Int, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _DelegationNft.Contract.WithdrawHold(&_DelegationNft.TransactOpts, validator, nft, tokenId, amount, holdTimestamp)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DelegationNft *DelegationNftTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegationNft.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DelegationNft *DelegationNftSession) Receive() (*types.Transaction, error) {
	return _DelegationNft.Contract.Receive(&_DelegationNft.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DelegationNft *DelegationNftTransactorSession) Receive() (*types.Transaction, error) {
	return _DelegationNft.Contract.Receive(&_DelegationNft.TransactOpts)
}

// DelegationNftFreezeTimeUpdatedIterator is returned from FilterFreezeTimeUpdated and is used to iterate over the raw logs and unpacked data for FreezeTimeUpdated events raised by the DelegationNft contract.
type DelegationNftFreezeTimeUpdatedIterator struct {
	Event *DelegationNftFreezeTimeUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationNftFreezeTimeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftFreezeTimeUpdated)
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
		it.Event = new(DelegationNftFreezeTimeUpdated)
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
func (it *DelegationNftFreezeTimeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftFreezeTimeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftFreezeTimeUpdated represents a FreezeTimeUpdated event raised by the DelegationNft contract.
type DelegationNftFreezeTimeUpdated struct {
	WithdrawFreezeTime *big.Int
	TransferFreezeTime *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterFreezeTimeUpdated is a free log retrieval operation binding the contract event 0x575c3ee3963e6c410284c700524d556e59f04881e7ff2702126adcfc33e2e22e.
//
// Solidity: event FreezeTimeUpdated(uint256 withdrawFreezeTime, uint256 transferFreezeTime)
func (_DelegationNft *DelegationNftFilterer) FilterFreezeTimeUpdated(opts *bind.FilterOpts) (*DelegationNftFreezeTimeUpdatedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "FreezeTimeUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationNftFreezeTimeUpdatedIterator{contract: _DelegationNft.contract, event: "FreezeTimeUpdated", logs: logs, sub: sub}, nil
}

// WatchFreezeTimeUpdated is a free log subscription operation binding the contract event 0x575c3ee3963e6c410284c700524d556e59f04881e7ff2702126adcfc33e2e22e.
//
// Solidity: event FreezeTimeUpdated(uint256 withdrawFreezeTime, uint256 transferFreezeTime)
func (_DelegationNft *DelegationNftFilterer) WatchFreezeTimeUpdated(opts *bind.WatchOpts, sink chan<- *DelegationNftFreezeTimeUpdated) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "FreezeTimeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftFreezeTimeUpdated)
				if err := _DelegationNft.contract.UnpackLog(event, "FreezeTimeUpdated", log); err != nil {
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

// ParseFreezeTimeUpdated is a log parse operation binding the contract event 0x575c3ee3963e6c410284c700524d556e59f04881e7ff2702126adcfc33e2e22e.
//
// Solidity: event FreezeTimeUpdated(uint256 withdrawFreezeTime, uint256 transferFreezeTime)
func (_DelegationNft *DelegationNftFilterer) ParseFreezeTimeUpdated(log types.Log) (*DelegationNftFreezeTimeUpdated, error) {
	event := new(DelegationNftFreezeTimeUpdated)
	if err := _DelegationNft.contract.UnpackLog(event, "FreezeTimeUpdated", log); err != nil {
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

// DelegationNftOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DelegationNft contract.
type DelegationNftOwnershipTransferredIterator struct {
	Event *DelegationNftOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DelegationNftOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftOwnershipTransferred)
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
		it.Event = new(DelegationNftOwnershipTransferred)
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
func (it *DelegationNftOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftOwnershipTransferred represents a OwnershipTransferred event raised by the DelegationNft contract.
type DelegationNftOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DelegationNft *DelegationNftFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DelegationNftOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftOwnershipTransferredIterator{contract: _DelegationNft.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DelegationNft *DelegationNftFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DelegationNftOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftOwnershipTransferred)
				if err := _DelegationNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseOwnershipTransferred(log types.Log) (*DelegationNftOwnershipTransferred, error) {
	event := new(DelegationNftOwnershipTransferred)
	if err := _DelegationNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftPenaltyAppliedToStakeIterator is returned from FilterPenaltyAppliedToStake and is used to iterate over the raw logs and unpacked data for PenaltyAppliedToStake events raised by the DelegationNft contract.
type DelegationNftPenaltyAppliedToStakeIterator struct {
	Event *DelegationNftPenaltyAppliedToStake // Event containing the contract specifics and raw log

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
func (it *DelegationNftPenaltyAppliedToStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftPenaltyAppliedToStake)
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
		it.Event = new(DelegationNftPenaltyAppliedToStake)
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
func (it *DelegationNftPenaltyAppliedToStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftPenaltyAppliedToStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftPenaltyAppliedToStake represents a PenaltyAppliedToStake event raised by the DelegationNft contract.
type DelegationNftPenaltyAppliedToStake struct {
	StakeId         [32]byte
	NewPenaltyIndex *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPenaltyAppliedToStake is a free log retrieval operation binding the contract event 0x22599c23e6ea88f2526338256e99d3dda3a4734dd2f2cf04b452c215101568e2.
//
// Solidity: event PenaltyAppliedToStake(bytes32 indexed stakeId, uint256 newPenaltyIndex)
func (_DelegationNft *DelegationNftFilterer) FilterPenaltyAppliedToStake(opts *bind.FilterOpts, stakeId [][32]byte) (*DelegationNftPenaltyAppliedToStakeIterator, error) {

	var stakeIdRule []interface{}
	for _, stakeIdItem := range stakeId {
		stakeIdRule = append(stakeIdRule, stakeIdItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "PenaltyAppliedToStake", stakeIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftPenaltyAppliedToStakeIterator{contract: _DelegationNft.contract, event: "PenaltyAppliedToStake", logs: logs, sub: sub}, nil
}

// WatchPenaltyAppliedToStake is a free log subscription operation binding the contract event 0x22599c23e6ea88f2526338256e99d3dda3a4734dd2f2cf04b452c215101568e2.
//
// Solidity: event PenaltyAppliedToStake(bytes32 indexed stakeId, uint256 newPenaltyIndex)
func (_DelegationNft *DelegationNftFilterer) WatchPenaltyAppliedToStake(opts *bind.WatchOpts, sink chan<- *DelegationNftPenaltyAppliedToStake, stakeId [][32]byte) (event.Subscription, error) {

	var stakeIdRule []interface{}
	for _, stakeIdItem := range stakeId {
		stakeIdRule = append(stakeIdRule, stakeIdItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "PenaltyAppliedToStake", stakeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftPenaltyAppliedToStake)
				if err := _DelegationNft.contract.UnpackLog(event, "PenaltyAppliedToStake", log); err != nil {
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

// ParsePenaltyAppliedToStake is a log parse operation binding the contract event 0x22599c23e6ea88f2526338256e99d3dda3a4734dd2f2cf04b452c215101568e2.
//
// Solidity: event PenaltyAppliedToStake(bytes32 indexed stakeId, uint256 newPenaltyIndex)
func (_DelegationNft *DelegationNftFilterer) ParsePenaltyAppliedToStake(log types.Log) (*DelegationNftPenaltyAppliedToStake, error) {
	event := new(DelegationNftPenaltyAppliedToStake)
	if err := _DelegationNft.contract.UnpackLog(event, "PenaltyAppliedToStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftPenaltyAppliedToValidatorIterator is returned from FilterPenaltyAppliedToValidator and is used to iterate over the raw logs and unpacked data for PenaltyAppliedToValidator events raised by the DelegationNft contract.
type DelegationNftPenaltyAppliedToValidatorIterator struct {
	Event *DelegationNftPenaltyAppliedToValidator // Event containing the contract specifics and raw log

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
func (it *DelegationNftPenaltyAppliedToValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftPenaltyAppliedToValidator)
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
		it.Event = new(DelegationNftPenaltyAppliedToValidator)
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
func (it *DelegationNftPenaltyAppliedToValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftPenaltyAppliedToValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftPenaltyAppliedToValidator represents a PenaltyAppliedToValidator event raised by the DelegationNft contract.
type DelegationNftPenaltyAppliedToValidator struct {
	Validator        common.Address
	Token            common.Address
	TokenType        uint8
	TokenId          *big.Int
	ValidatorReserve IDecimalDelegationCommonValidatorReserve
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPenaltyAppliedToValidator is a free log retrieval operation binding the contract event 0x3b900bc82221c5a44e5d8e57d6f36ba5c3c9e10bcfdbef3441dd877a5e75981c.
//
// Solidity: event PenaltyAppliedToValidator(address indexed validator, address token, uint8 tokenType, uint256 tokenId, (uint256,uint256) validatorReserve)
func (_DelegationNft *DelegationNftFilterer) FilterPenaltyAppliedToValidator(opts *bind.FilterOpts, validator []common.Address) (*DelegationNftPenaltyAppliedToValidatorIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "PenaltyAppliedToValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftPenaltyAppliedToValidatorIterator{contract: _DelegationNft.contract, event: "PenaltyAppliedToValidator", logs: logs, sub: sub}, nil
}

// WatchPenaltyAppliedToValidator is a free log subscription operation binding the contract event 0x3b900bc82221c5a44e5d8e57d6f36ba5c3c9e10bcfdbef3441dd877a5e75981c.
//
// Solidity: event PenaltyAppliedToValidator(address indexed validator, address token, uint8 tokenType, uint256 tokenId, (uint256,uint256) validatorReserve)
func (_DelegationNft *DelegationNftFilterer) WatchPenaltyAppliedToValidator(opts *bind.WatchOpts, sink chan<- *DelegationNftPenaltyAppliedToValidator, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "PenaltyAppliedToValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftPenaltyAppliedToValidator)
				if err := _DelegationNft.contract.UnpackLog(event, "PenaltyAppliedToValidator", log); err != nil {
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

// ParsePenaltyAppliedToValidator is a log parse operation binding the contract event 0x3b900bc82221c5a44e5d8e57d6f36ba5c3c9e10bcfdbef3441dd877a5e75981c.
//
// Solidity: event PenaltyAppliedToValidator(address indexed validator, address token, uint8 tokenType, uint256 tokenId, (uint256,uint256) validatorReserve)
func (_DelegationNft *DelegationNftFilterer) ParsePenaltyAppliedToValidator(log types.Log) (*DelegationNftPenaltyAppliedToValidator, error) {
	event := new(DelegationNftPenaltyAppliedToValidator)
	if err := _DelegationNft.contract.UnpackLog(event, "PenaltyAppliedToValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftStakeAmountUpdatedIterator is returned from FilterStakeAmountUpdated and is used to iterate over the raw logs and unpacked data for StakeAmountUpdated events raised by the DelegationNft contract.
type DelegationNftStakeAmountUpdatedIterator struct {
	Event *DelegationNftStakeAmountUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationNftStakeAmountUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftStakeAmountUpdated)
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
		it.Event = new(DelegationNftStakeAmountUpdated)
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
func (it *DelegationNftStakeAmountUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftStakeAmountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftStakeAmountUpdated represents a StakeAmountUpdated event raised by the DelegationNft contract.
type DelegationNftStakeAmountUpdated struct {
	StakeId       [32]byte
	ChangedAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeAmountUpdated is a free log retrieval operation binding the contract event 0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959.
//
// Solidity: event StakeAmountUpdated(bytes32 stakeId, int256 changedAmount)
func (_DelegationNft *DelegationNftFilterer) FilterStakeAmountUpdated(opts *bind.FilterOpts) (*DelegationNftStakeAmountUpdatedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "StakeAmountUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationNftStakeAmountUpdatedIterator{contract: _DelegationNft.contract, event: "StakeAmountUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeAmountUpdated is a free log subscription operation binding the contract event 0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959.
//
// Solidity: event StakeAmountUpdated(bytes32 stakeId, int256 changedAmount)
func (_DelegationNft *DelegationNftFilterer) WatchStakeAmountUpdated(opts *bind.WatchOpts, sink chan<- *DelegationNftStakeAmountUpdated) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "StakeAmountUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftStakeAmountUpdated)
				if err := _DelegationNft.contract.UnpackLog(event, "StakeAmountUpdated", log); err != nil {
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

// ParseStakeAmountUpdated is a log parse operation binding the contract event 0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959.
//
// Solidity: event StakeAmountUpdated(bytes32 stakeId, int256 changedAmount)
func (_DelegationNft *DelegationNftFilterer) ParseStakeAmountUpdated(log types.Log) (*DelegationNftStakeAmountUpdated, error) {
	event := new(DelegationNftStakeAmountUpdated)
	if err := _DelegationNft.contract.UnpackLog(event, "StakeAmountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftStakeHoldedIterator is returned from FilterStakeHolded and is used to iterate over the raw logs and unpacked data for StakeHolded events raised by the DelegationNft contract.
type DelegationNftStakeHoldedIterator struct {
	Event *DelegationNftStakeHolded // Event containing the contract specifics and raw log

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
func (it *DelegationNftStakeHoldedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftStakeHolded)
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
		it.Event = new(DelegationNftStakeHolded)
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
func (it *DelegationNftStakeHoldedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftStakeHoldedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftStakeHolded represents a StakeHolded event raised by the DelegationNft contract.
type DelegationNftStakeHolded struct {
	StakeId [32]byte
	Stake   IDecimalDelegationCommonStake
	IsNew   bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeHolded is a free log retrieval operation binding the contract event 0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb.
//
// Solidity: event StakeHolded(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) FilterStakeHolded(opts *bind.FilterOpts) (*DelegationNftStakeHoldedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "StakeHolded")
	if err != nil {
		return nil, err
	}
	return &DelegationNftStakeHoldedIterator{contract: _DelegationNft.contract, event: "StakeHolded", logs: logs, sub: sub}, nil
}

// WatchStakeHolded is a free log subscription operation binding the contract event 0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb.
//
// Solidity: event StakeHolded(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) WatchStakeHolded(opts *bind.WatchOpts, sink chan<- *DelegationNftStakeHolded) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "StakeHolded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftStakeHolded)
				if err := _DelegationNft.contract.UnpackLog(event, "StakeHolded", log); err != nil {
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

// ParseStakeHolded is a log parse operation binding the contract event 0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb.
//
// Solidity: event StakeHolded(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) ParseStakeHolded(log types.Log) (*DelegationNftStakeHolded, error) {
	event := new(DelegationNftStakeHolded)
	if err := _DelegationNft.contract.UnpackLog(event, "StakeHolded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftStakeUpdatedIterator is returned from FilterStakeUpdated and is used to iterate over the raw logs and unpacked data for StakeUpdated events raised by the DelegationNft contract.
type DelegationNftStakeUpdatedIterator struct {
	Event *DelegationNftStakeUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationNftStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftStakeUpdated)
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
		it.Event = new(DelegationNftStakeUpdated)
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
func (it *DelegationNftStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftStakeUpdated represents a StakeUpdated event raised by the DelegationNft contract.
type DelegationNftStakeUpdated struct {
	StakeId [32]byte
	Stake   IDecimalDelegationCommonStake
	IsNew   bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdated is a free log retrieval operation binding the contract event 0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b.
//
// Solidity: event StakeUpdated(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) FilterStakeUpdated(opts *bind.FilterOpts) (*DelegationNftStakeUpdatedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationNftStakeUpdatedIterator{contract: _DelegationNft.contract, event: "StakeUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeUpdated is a free log subscription operation binding the contract event 0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b.
//
// Solidity: event StakeUpdated(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) WatchStakeUpdated(opts *bind.WatchOpts, sink chan<- *DelegationNftStakeUpdated) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftStakeUpdated)
				if err := _DelegationNft.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
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

// ParseStakeUpdated is a log parse operation binding the contract event 0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b.
//
// Solidity: event StakeUpdated(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_DelegationNft *DelegationNftFilterer) ParseStakeUpdated(log types.Log) (*DelegationNftStakeUpdated, error) {
	event := new(DelegationNftStakeUpdated)
	if err := _DelegationNft.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftTransferCompletedIterator is returned from FilterTransferCompleted and is used to iterate over the raw logs and unpacked data for TransferCompleted events raised by the DelegationNft contract.
type DelegationNftTransferCompletedIterator struct {
	Event *DelegationNftTransferCompleted // Event containing the contract specifics and raw log

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
func (it *DelegationNftTransferCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftTransferCompleted)
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
		it.Event = new(DelegationNftTransferCompleted)
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
func (it *DelegationNftTransferCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftTransferCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftTransferCompleted represents a TransferCompleted event raised by the DelegationNft contract.
type DelegationNftTransferCompleted struct {
	StakeIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferCompleted is a free log retrieval operation binding the contract event 0xfd987fd5c7de5139db194fdde15ddabcec1b78c3bfc832ad563ac57a9bfa9b36.
//
// Solidity: event TransferCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) FilterTransferCompleted(opts *bind.FilterOpts) (*DelegationNftTransferCompletedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "TransferCompleted")
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransferCompletedIterator{contract: _DelegationNft.contract, event: "TransferCompleted", logs: logs, sub: sub}, nil
}

// WatchTransferCompleted is a free log subscription operation binding the contract event 0xfd987fd5c7de5139db194fdde15ddabcec1b78c3bfc832ad563ac57a9bfa9b36.
//
// Solidity: event TransferCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) WatchTransferCompleted(opts *bind.WatchOpts, sink chan<- *DelegationNftTransferCompleted) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "TransferCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftTransferCompleted)
				if err := _DelegationNft.contract.UnpackLog(event, "TransferCompleted", log); err != nil {
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

// ParseTransferCompleted is a log parse operation binding the contract event 0xfd987fd5c7de5139db194fdde15ddabcec1b78c3bfc832ad563ac57a9bfa9b36.
//
// Solidity: event TransferCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) ParseTransferCompleted(log types.Log) (*DelegationNftTransferCompleted, error) {
	event := new(DelegationNftTransferCompleted)
	if err := _DelegationNft.contract.UnpackLog(event, "TransferCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftTransferRequestIterator is returned from FilterTransferRequest and is used to iterate over the raw logs and unpacked data for TransferRequest events raised by the DelegationNft contract.
type DelegationNftTransferRequestIterator struct {
	Event *DelegationNftTransferRequest // Event containing the contract specifics and raw log

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
func (it *DelegationNftTransferRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftTransferRequest)
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
		it.Event = new(DelegationNftTransferRequest)
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
func (it *DelegationNftTransferRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftTransferRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftTransferRequest represents a TransferRequest event raised by the DelegationNft contract.
type DelegationNftTransferRequest struct {
	StakeId     [32]byte
	StakeIndex  *big.Int
	FrozenStake IDecimalDelegationCommonFrozenStake
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTransferRequest is a free log retrieval operation binding the contract event 0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f.
//
// Solidity: event TransferRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) FilterTransferRequest(opts *bind.FilterOpts) (*DelegationNftTransferRequestIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "TransferRequest")
	if err != nil {
		return nil, err
	}
	return &DelegationNftTransferRequestIterator{contract: _DelegationNft.contract, event: "TransferRequest", logs: logs, sub: sub}, nil
}

// WatchTransferRequest is a free log subscription operation binding the contract event 0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f.
//
// Solidity: event TransferRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) WatchTransferRequest(opts *bind.WatchOpts, sink chan<- *DelegationNftTransferRequest) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "TransferRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftTransferRequest)
				if err := _DelegationNft.contract.UnpackLog(event, "TransferRequest", log); err != nil {
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

// ParseTransferRequest is a log parse operation binding the contract event 0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f.
//
// Solidity: event TransferRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) ParseTransferRequest(log types.Log) (*DelegationNftTransferRequest, error) {
	event := new(DelegationNftTransferRequest)
	if err := _DelegationNft.contract.UnpackLog(event, "TransferRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DelegationNft contract.
type DelegationNftUpgradedIterator struct {
	Event *DelegationNftUpgraded // Event containing the contract specifics and raw log

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
func (it *DelegationNftUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftUpgraded)
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
		it.Event = new(DelegationNftUpgraded)
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
func (it *DelegationNftUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftUpgraded represents a Upgraded event raised by the DelegationNft contract.
type DelegationNftUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DelegationNft *DelegationNftFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DelegationNftUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DelegationNftUpgradedIterator{contract: _DelegationNft.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DelegationNft *DelegationNftFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DelegationNftUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftUpgraded)
				if err := _DelegationNft.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_DelegationNft *DelegationNftFilterer) ParseUpgraded(log types.Log) (*DelegationNftUpgraded, error) {
	event := new(DelegationNftUpgraded)
	if err := _DelegationNft.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftWithdrawCompletedIterator is returned from FilterWithdrawCompleted and is used to iterate over the raw logs and unpacked data for WithdrawCompleted events raised by the DelegationNft contract.
type DelegationNftWithdrawCompletedIterator struct {
	Event *DelegationNftWithdrawCompleted // Event containing the contract specifics and raw log

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
func (it *DelegationNftWithdrawCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftWithdrawCompleted)
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
		it.Event = new(DelegationNftWithdrawCompleted)
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
func (it *DelegationNftWithdrawCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftWithdrawCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftWithdrawCompleted represents a WithdrawCompleted event raised by the DelegationNft contract.
type DelegationNftWithdrawCompleted struct {
	StakeIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCompleted is a free log retrieval operation binding the contract event 0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7.
//
// Solidity: event WithdrawCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) FilterWithdrawCompleted(opts *bind.FilterOpts) (*DelegationNftWithdrawCompletedIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "WithdrawCompleted")
	if err != nil {
		return nil, err
	}
	return &DelegationNftWithdrawCompletedIterator{contract: _DelegationNft.contract, event: "WithdrawCompleted", logs: logs, sub: sub}, nil
}

// WatchWithdrawCompleted is a free log subscription operation binding the contract event 0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7.
//
// Solidity: event WithdrawCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) WatchWithdrawCompleted(opts *bind.WatchOpts, sink chan<- *DelegationNftWithdrawCompleted) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "WithdrawCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftWithdrawCompleted)
				if err := _DelegationNft.contract.UnpackLog(event, "WithdrawCompleted", log); err != nil {
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

// ParseWithdrawCompleted is a log parse operation binding the contract event 0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7.
//
// Solidity: event WithdrawCompleted(uint256 stakeIndex)
func (_DelegationNft *DelegationNftFilterer) ParseWithdrawCompleted(log types.Log) (*DelegationNftWithdrawCompleted, error) {
	event := new(DelegationNftWithdrawCompleted)
	if err := _DelegationNft.contract.UnpackLog(event, "WithdrawCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationNftWithdrawRequestIterator is returned from FilterWithdrawRequest and is used to iterate over the raw logs and unpacked data for WithdrawRequest events raised by the DelegationNft contract.
type DelegationNftWithdrawRequestIterator struct {
	Event *DelegationNftWithdrawRequest // Event containing the contract specifics and raw log

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
func (it *DelegationNftWithdrawRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationNftWithdrawRequest)
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
		it.Event = new(DelegationNftWithdrawRequest)
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
func (it *DelegationNftWithdrawRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationNftWithdrawRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationNftWithdrawRequest represents a WithdrawRequest event raised by the DelegationNft contract.
type DelegationNftWithdrawRequest struct {
	StakeId     [32]byte
	StakeIndex  *big.Int
	FrozenStake IDecimalDelegationCommonFrozenStake
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequest is a free log retrieval operation binding the contract event 0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab.
//
// Solidity: event WithdrawRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) FilterWithdrawRequest(opts *bind.FilterOpts) (*DelegationNftWithdrawRequestIterator, error) {

	logs, sub, err := _DelegationNft.contract.FilterLogs(opts, "WithdrawRequest")
	if err != nil {
		return nil, err
	}
	return &DelegationNftWithdrawRequestIterator{contract: _DelegationNft.contract, event: "WithdrawRequest", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequest is a free log subscription operation binding the contract event 0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab.
//
// Solidity: event WithdrawRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) WatchWithdrawRequest(opts *bind.WatchOpts, sink chan<- *DelegationNftWithdrawRequest) (event.Subscription, error) {

	logs, sub, err := _DelegationNft.contract.WatchLogs(opts, "WithdrawRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationNftWithdrawRequest)
				if err := _DelegationNft.contract.UnpackLog(event, "WithdrawRequest", log); err != nil {
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

// ParseWithdrawRequest is a log parse operation binding the contract event 0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab.
//
// Solidity: event WithdrawRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_DelegationNft *DelegationNftFilterer) ParseWithdrawRequest(log types.Log) (*DelegationNftWithdrawRequest, error) {
	event := new(DelegationNftWithdrawRequest)
	if err := _DelegationNft.contract.UnpackLog(event, "WithdrawRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
