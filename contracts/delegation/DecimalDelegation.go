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

// DelegationMetaData contains all meta data concerning the Delegation contract.
var DelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FrozenStakesQueueIsEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidFrozenType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTimestamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTokenType\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValidator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCompletableFrozenStakes\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoPenaltyToApply\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeAlreadyUnfrozen\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeInactive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimestampError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAmount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawFreezeTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transferFreezeTime\",\"type\":\"uint256\"}],\"name\":\"FreezeTimeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newPenaltyIndex\",\"type\":\"uint256\"}],\"name\":\"PenaltyAppliedToStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"penaltyIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.ValidatorReserve\",\"name\":\"validatorReserve\",\"type\":\"tuple\"}],\"name\":\"PenaltyAppliedToValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"changedAmount\",\"type\":\"int256\"}],\"name\":\"StakeAmountUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isNew\",\"type\":\"bool\"}],\"name\":\"StakeHolded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isNew\",\"type\":\"bool\"}],\"name\":\"StakeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"}],\"name\":\"TransferCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"frozenStake\",\"type\":\"tuple\"}],\"name\":\"TransferRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"}],\"name\":\"WithdrawCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakeIndex\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"frozenStake\",\"type\":\"tuple\"}],\"name\":\"WithdrawRequest\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"WETH\",\"outputs\":[{\"internalType\":\"contractIWETH\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.ValidatorToken[]\",\"name\":\"validatorTokens\",\"type\":\"tuple[]\"}],\"name\":\"applyPenaltiesToValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"indexes\",\"type\":\"uint256[]\"}],\"name\":\"complete\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"delegateDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"}],\"name\":\"delegateDELTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"delegateHold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"delegateHoldByPermit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"delegateHoldDEL\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegateTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContractCenter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"}],\"name\":\"getFreezeTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getFrozenStake\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.FrozenStake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"stakeIndexes\",\"type\":\"uint256[]\"}],\"name\":\"getFrozenStakes\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"stake\",\"type\":\"tuple\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeStatus\",\"name\":\"freezeStatus\",\"type\":\"uint8\"},{\"internalType\":\"enumIDecimalDelegationCommon.FreezeType\",\"name\":\"freezeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"unfreezeTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.FrozenStake[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"getHoldStake\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"getHoldStakeId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImpl\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getStake\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getStakeId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stakeId\",\"type\":\"bytes32\"}],\"name\":\"getStakePenaltyIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stakeIds\",\"type\":\"bytes32[]\"}],\"name\":\"getStakes\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"enumIDecimalDelegationCommon.TokenType\",\"name\":\"tokenType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.Stake[]\",\"name\":\"stakes\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getValidatorReserve\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"penaltyIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"internalType\":\"structIDecimalDelegationCommon.ValidatorReserve\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountToHold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"oldHoldTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newHoldTimestamp\",\"type\":\"uint256\"}],\"name\":\"hold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractCenter\",\"type\":\"address\"}],\"name\":\"setContractCenter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"withdrawFreezeTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"transferFreezeTime\",\"type\":\"uint256\"}],\"name\":\"setFreezeTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newValidator\",\"type\":\"address\"}],\"name\":\"transferHold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImpl\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"holdTimestamp\",\"type\":\"uint256\"}],\"name\":\"withdrawHold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Delegation *DelegationCaller) WETH(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "WETH")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Delegation *DelegationSession) WETH() (common.Address, error) {
	return _Delegation.Contract.WETH(&_Delegation.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Delegation *DelegationCallerSession) WETH() (common.Address, error) {
	return _Delegation.Contract.WETH(&_Delegation.CallOpts)
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

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_Delegation *DelegationCaller) GetFreezeTime(opts *bind.CallOpts, freezeType uint8) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getFreezeTime", freezeType)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_Delegation *DelegationSession) GetFreezeTime(freezeType uint8) (*big.Int, error) {
	return _Delegation.Contract.GetFreezeTime(&_Delegation.CallOpts, freezeType)
}

// GetFreezeTime is a free data retrieval call binding the contract method 0x0e08def3.
//
// Solidity: function getFreezeTime(uint8 freezeType) view returns(uint256)
func (_Delegation *DelegationCallerSession) GetFreezeTime(freezeType uint8) (*big.Int, error) {
	return _Delegation.Contract.GetFreezeTime(&_Delegation.CallOpts, freezeType)
}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_Delegation *DelegationCaller) GetFrozenStake(opts *bind.CallOpts, index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getFrozenStake", index)

	if err != nil {
		return *new(IDecimalDelegationCommonFrozenStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonFrozenStake)).(*IDecimalDelegationCommonFrozenStake)

	return out0, err

}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_Delegation *DelegationSession) GetFrozenStake(index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	return _Delegation.Contract.GetFrozenStake(&_Delegation.CallOpts, index)
}

// GetFrozenStake is a free data retrieval call binding the contract method 0xd8f06a8f.
//
// Solidity: function getFrozenStake(uint256 index) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256))
func (_Delegation *DelegationCallerSession) GetFrozenStake(index *big.Int) (IDecimalDelegationCommonFrozenStake, error) {
	return _Delegation.Contract.GetFrozenStake(&_Delegation.CallOpts, index)
}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_Delegation *DelegationCaller) GetFrozenStakes(opts *bind.CallOpts, stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getFrozenStakes", stakeIndexes)

	if err != nil {
		return *new([]IDecimalDelegationCommonFrozenStake), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDecimalDelegationCommonFrozenStake)).(*[]IDecimalDelegationCommonFrozenStake)

	return out0, err

}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_Delegation *DelegationSession) GetFrozenStakes(stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	return _Delegation.Contract.GetFrozenStakes(&_Delegation.CallOpts, stakeIndexes)
}

// GetFrozenStakes is a free data retrieval call binding the contract method 0x722c76f8.
//
// Solidity: function getFrozenStakes(uint256[] stakeIndexes) view returns(((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256)[])
func (_Delegation *DelegationCallerSession) GetFrozenStakes(stakeIndexes []*big.Int) ([]IDecimalDelegationCommonFrozenStake, error) {
	return _Delegation.Contract.GetFrozenStakes(&_Delegation.CallOpts, stakeIndexes)
}

// GetHoldStake is a free data retrieval call binding the contract method 0x3d249faa.
//
// Solidity: function getHoldStake(address validator, address delegator, address token, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationCaller) GetHoldStake(opts *bind.CallOpts, validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getHoldStake", validator, delegator, token, holdTimestamp)

	if err != nil {
		return *new(IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonStake)).(*IDecimalDelegationCommonStake)

	return out0, err

}

// GetHoldStake is a free data retrieval call binding the contract method 0x3d249faa.
//
// Solidity: function getHoldStake(address validator, address delegator, address token, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationSession) GetHoldStake(validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetHoldStake(&_Delegation.CallOpts, validator, delegator, token, holdTimestamp)
}

// GetHoldStake is a free data retrieval call binding the contract method 0x3d249faa.
//
// Solidity: function getHoldStake(address validator, address delegator, address token, uint256 holdTimestamp) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationCallerSession) GetHoldStake(validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) (IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetHoldStake(&_Delegation.CallOpts, validator, delegator, token, holdTimestamp)
}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x68b049d9.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address token, uint256 holdTimestamp) pure returns(bytes32)
func (_Delegation *DelegationCaller) GetHoldStakeId(opts *bind.CallOpts, validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getHoldStakeId", validator, delegator, token, holdTimestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x68b049d9.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address token, uint256 holdTimestamp) pure returns(bytes32)
func (_Delegation *DelegationSession) GetHoldStakeId(validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) ([32]byte, error) {
	return _Delegation.Contract.GetHoldStakeId(&_Delegation.CallOpts, validator, delegator, token, holdTimestamp)
}

// GetHoldStakeId is a free data retrieval call binding the contract method 0x68b049d9.
//
// Solidity: function getHoldStakeId(address validator, address delegator, address token, uint256 holdTimestamp) pure returns(bytes32)
func (_Delegation *DelegationCallerSession) GetHoldStakeId(validator common.Address, delegator common.Address, token common.Address, holdTimestamp *big.Int) ([32]byte, error) {
	return _Delegation.Contract.GetHoldStakeId(&_Delegation.CallOpts, validator, delegator, token, holdTimestamp)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Delegation *DelegationCaller) GetImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Delegation *DelegationSession) GetImpl() (common.Address, error) {
	return _Delegation.Contract.GetImpl(&_Delegation.CallOpts)
}

// GetImpl is a free data retrieval call binding the contract method 0xdfb80831.
//
// Solidity: function getImpl() view returns(address)
func (_Delegation *DelegationCallerSession) GetImpl() (common.Address, error) {
	return _Delegation.Contract.GetImpl(&_Delegation.CallOpts)
}

// GetStake is a free data retrieval call binding the contract method 0x5d518866.
//
// Solidity: function getStake(address validator, address delegator, address token) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationCaller) GetStake(opts *bind.CallOpts, validator common.Address, delegator common.Address, token common.Address) (IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getStake", validator, delegator, token)

	if err != nil {
		return *new(IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonStake)).(*IDecimalDelegationCommonStake)

	return out0, err

}

// GetStake is a free data retrieval call binding the contract method 0x5d518866.
//
// Solidity: function getStake(address validator, address delegator, address token) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationSession) GetStake(validator common.Address, delegator common.Address, token common.Address) (IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetStake(&_Delegation.CallOpts, validator, delegator, token)
}

// GetStake is a free data retrieval call binding the contract method 0x5d518866.
//
// Solidity: function getStake(address validator, address delegator, address token) view returns((address,address,address,uint256,uint256,uint8,uint256))
func (_Delegation *DelegationCallerSession) GetStake(validator common.Address, delegator common.Address, token common.Address) (IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetStake(&_Delegation.CallOpts, validator, delegator, token)
}

// GetStakeId is a free data retrieval call binding the contract method 0x1103c5dc.
//
// Solidity: function getStakeId(address validator, address delegator, address token) pure returns(bytes32)
func (_Delegation *DelegationCaller) GetStakeId(opts *bind.CallOpts, validator common.Address, delegator common.Address, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getStakeId", validator, delegator, token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetStakeId is a free data retrieval call binding the contract method 0x1103c5dc.
//
// Solidity: function getStakeId(address validator, address delegator, address token) pure returns(bytes32)
func (_Delegation *DelegationSession) GetStakeId(validator common.Address, delegator common.Address, token common.Address) ([32]byte, error) {
	return _Delegation.Contract.GetStakeId(&_Delegation.CallOpts, validator, delegator, token)
}

// GetStakeId is a free data retrieval call binding the contract method 0x1103c5dc.
//
// Solidity: function getStakeId(address validator, address delegator, address token) pure returns(bytes32)
func (_Delegation *DelegationCallerSession) GetStakeId(validator common.Address, delegator common.Address, token common.Address) ([32]byte, error) {
	return _Delegation.Contract.GetStakeId(&_Delegation.CallOpts, validator, delegator, token)
}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_Delegation *DelegationCaller) GetStakePenaltyIndex(opts *bind.CallOpts, stakeId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getStakePenaltyIndex", stakeId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_Delegation *DelegationSession) GetStakePenaltyIndex(stakeId [32]byte) (*big.Int, error) {
	return _Delegation.Contract.GetStakePenaltyIndex(&_Delegation.CallOpts, stakeId)
}

// GetStakePenaltyIndex is a free data retrieval call binding the contract method 0xe6376614.
//
// Solidity: function getStakePenaltyIndex(bytes32 stakeId) view returns(uint256)
func (_Delegation *DelegationCallerSession) GetStakePenaltyIndex(stakeId [32]byte) (*big.Int, error) {
	return _Delegation.Contract.GetStakePenaltyIndex(&_Delegation.CallOpts, stakeId)
}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_Delegation *DelegationCaller) GetStakes(opts *bind.CallOpts, stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getStakes", stakeIds)

	if err != nil {
		return *new([]IDecimalDelegationCommonStake), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDecimalDelegationCommonStake)).(*[]IDecimalDelegationCommonStake)

	return out0, err

}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_Delegation *DelegationSession) GetStakes(stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetStakes(&_Delegation.CallOpts, stakeIds)
}

// GetStakes is a free data retrieval call binding the contract method 0x226f6ea2.
//
// Solidity: function getStakes(bytes32[] stakeIds) view returns((address,address,address,uint256,uint256,uint8,uint256)[] stakes)
func (_Delegation *DelegationCallerSession) GetStakes(stakeIds [][32]byte) ([]IDecimalDelegationCommonStake, error) {
	return _Delegation.Contract.GetStakes(&_Delegation.CallOpts, stakeIds)
}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x330b5c45.
//
// Solidity: function getValidatorReserve(address validator, address token) view returns((uint256,uint256))
func (_Delegation *DelegationCaller) GetValidatorReserve(opts *bind.CallOpts, validator common.Address, token common.Address) (IDecimalDelegationCommonValidatorReserve, error) {
	var out []interface{}
	err := _Delegation.contract.Call(opts, &out, "getValidatorReserve", validator, token)

	if err != nil {
		return *new(IDecimalDelegationCommonValidatorReserve), err
	}

	out0 := *abi.ConvertType(out[0], new(IDecimalDelegationCommonValidatorReserve)).(*IDecimalDelegationCommonValidatorReserve)

	return out0, err

}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x330b5c45.
//
// Solidity: function getValidatorReserve(address validator, address token) view returns((uint256,uint256))
func (_Delegation *DelegationSession) GetValidatorReserve(validator common.Address, token common.Address) (IDecimalDelegationCommonValidatorReserve, error) {
	return _Delegation.Contract.GetValidatorReserve(&_Delegation.CallOpts, validator, token)
}

// GetValidatorReserve is a free data retrieval call binding the contract method 0x330b5c45.
//
// Solidity: function getValidatorReserve(address validator, address token) view returns((uint256,uint256))
func (_Delegation *DelegationCallerSession) GetValidatorReserve(validator common.Address, token common.Address) (IDecimalDelegationCommonValidatorReserve, error) {
	return _Delegation.Contract.GetValidatorReserve(&_Delegation.CallOpts, validator, token)
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

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_Delegation *DelegationTransactor) ApplyPenaltiesToValidator(opts *bind.TransactOpts, validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "applyPenaltiesToValidator", validatorTokens)
}

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_Delegation *DelegationSession) ApplyPenaltiesToValidator(validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _Delegation.Contract.ApplyPenaltiesToValidator(&_Delegation.TransactOpts, validatorTokens)
}

// ApplyPenaltiesToValidator is a paid mutator transaction binding the contract method 0x1996c40a.
//
// Solidity: function applyPenaltiesToValidator((address,address,uint256)[] validatorTokens) returns()
func (_Delegation *DelegationTransactorSession) ApplyPenaltiesToValidator(validatorTokens []IDecimalDelegationCommonValidatorToken) (*types.Transaction, error) {
	return _Delegation.Contract.ApplyPenaltiesToValidator(&_Delegation.TransactOpts, validatorTokens)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_Delegation *DelegationTransactor) Complete(opts *bind.TransactOpts, indexes []*big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "complete", indexes)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_Delegation *DelegationSession) Complete(indexes []*big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Complete(&_Delegation.TransactOpts, indexes)
}

// Complete is a paid mutator transaction binding the contract method 0x5d95f94f.
//
// Solidity: function complete(uint256[] indexes) returns()
func (_Delegation *DelegationTransactorSession) Complete(indexes []*big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Complete(&_Delegation.TransactOpts, indexes)
}

// Delegate is a paid mutator transaction binding the contract method 0xc26f3acf.
//
// Solidity: function delegate(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactor) Delegate(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegate", validator, token, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0xc26f3acf.
//
// Solidity: function delegate(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationSession) Delegate(validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Delegate(&_Delegation.TransactOpts, validator, token, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0xc26f3acf.
//
// Solidity: function delegate(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactorSession) Delegate(validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Delegate(&_Delegation.TransactOpts, validator, token, amount)
}

// DelegateByPermit is a paid mutator transaction binding the contract method 0x1418cb93.
//
// Solidity: function delegateByPermit(address validator, address token, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) DelegateByPermit(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateByPermit", validator, token, amount, deadline, v, r, s)
}

// DelegateByPermit is a paid mutator transaction binding the contract method 0x1418cb93.
//
// Solidity: function delegateByPermit(address validator, address token, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) DelegateByPermit(validator common.Address, token common.Address, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateByPermit(&_Delegation.TransactOpts, validator, token, amount, deadline, v, r, s)
}

// DelegateByPermit is a paid mutator transaction binding the contract method 0x1418cb93.
//
// Solidity: function delegateByPermit(address validator, address token, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) DelegateByPermit(validator common.Address, token common.Address, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateByPermit(&_Delegation.TransactOpts, validator, token, amount, deadline, v, r, s)
}

// DelegateDEL is a paid mutator transaction binding the contract method 0xcc5e5500.
//
// Solidity: function delegateDEL(address validator) payable returns()
func (_Delegation *DelegationTransactor) DelegateDEL(opts *bind.TransactOpts, validator common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateDEL", validator)
}

// DelegateDEL is a paid mutator transaction binding the contract method 0xcc5e5500.
//
// Solidity: function delegateDEL(address validator) payable returns()
func (_Delegation *DelegationSession) DelegateDEL(validator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateDEL(&_Delegation.TransactOpts, validator)
}

// DelegateDEL is a paid mutator transaction binding the contract method 0xcc5e5500.
//
// Solidity: function delegateDEL(address validator) payable returns()
func (_Delegation *DelegationTransactorSession) DelegateDEL(validator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateDEL(&_Delegation.TransactOpts, validator)
}

// DelegateDELTo is a paid mutator transaction binding the contract method 0x2d9473ea.
//
// Solidity: function delegateDELTo(address delegator, address validator) payable returns()
func (_Delegation *DelegationTransactor) DelegateDELTo(opts *bind.TransactOpts, delegator common.Address, validator common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateDELTo", delegator, validator)
}

// DelegateDELTo is a paid mutator transaction binding the contract method 0x2d9473ea.
//
// Solidity: function delegateDELTo(address delegator, address validator) payable returns()
func (_Delegation *DelegationSession) DelegateDELTo(delegator common.Address, validator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateDELTo(&_Delegation.TransactOpts, delegator, validator)
}

// DelegateDELTo is a paid mutator transaction binding the contract method 0x2d9473ea.
//
// Solidity: function delegateDELTo(address delegator, address validator) payable returns()
func (_Delegation *DelegationTransactorSession) DelegateDELTo(delegator common.Address, validator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateDELTo(&_Delegation.TransactOpts, delegator, validator)
}

// DelegateHold is a paid mutator transaction binding the contract method 0x40bbbd6d.
//
// Solidity: function delegateHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationTransactor) DelegateHold(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateHold", validator, token, amount, holdTimestamp)
}

// DelegateHold is a paid mutator transaction binding the contract method 0x40bbbd6d.
//
// Solidity: function delegateHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationSession) DelegateHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp)
}

// DelegateHold is a paid mutator transaction binding the contract method 0x40bbbd6d.
//
// Solidity: function delegateHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationTransactorSession) DelegateHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp)
}

// DelegateHoldByPermit is a paid mutator transaction binding the contract method 0x713ec4ac.
//
// Solidity: function delegateHoldByPermit(address validator, address token, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactor) DelegateHoldByPermit(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateHoldByPermit", validator, token, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldByPermit is a paid mutator transaction binding the contract method 0x713ec4ac.
//
// Solidity: function delegateHoldByPermit(address validator, address token, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationSession) DelegateHoldByPermit(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHoldByPermit(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldByPermit is a paid mutator transaction binding the contract method 0x713ec4ac.
//
// Solidity: function delegateHoldByPermit(address validator, address token, uint256 amount, uint256 holdTimestamp, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_Delegation *DelegationTransactorSession) DelegateHoldByPermit(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHoldByPermit(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp, deadline, v, r, s)
}

// DelegateHoldDEL is a paid mutator transaction binding the contract method 0x9e0f4c7c.
//
// Solidity: function delegateHoldDEL(address validator, uint256 holdTimestamp) payable returns()
func (_Delegation *DelegationTransactor) DelegateHoldDEL(opts *bind.TransactOpts, validator common.Address, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateHoldDEL", validator, holdTimestamp)
}

// DelegateHoldDEL is a paid mutator transaction binding the contract method 0x9e0f4c7c.
//
// Solidity: function delegateHoldDEL(address validator, uint256 holdTimestamp) payable returns()
func (_Delegation *DelegationSession) DelegateHoldDEL(validator common.Address, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHoldDEL(&_Delegation.TransactOpts, validator, holdTimestamp)
}

// DelegateHoldDEL is a paid mutator transaction binding the contract method 0x9e0f4c7c.
//
// Solidity: function delegateHoldDEL(address validator, uint256 holdTimestamp) payable returns()
func (_Delegation *DelegationTransactorSession) DelegateHoldDEL(validator common.Address, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateHoldDEL(&_Delegation.TransactOpts, validator, holdTimestamp)
}

// DelegateTo is a paid mutator transaction binding the contract method 0x2c1b15f4.
//
// Solidity: function delegateTo(address delegator, address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactor) DelegateTo(opts *bind.TransactOpts, delegator common.Address, validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "delegateTo", delegator, validator, token, amount)
}

// DelegateTo is a paid mutator transaction binding the contract method 0x2c1b15f4.
//
// Solidity: function delegateTo(address delegator, address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationSession) DelegateTo(delegator common.Address, validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateTo(&_Delegation.TransactOpts, delegator, validator, token, amount)
}

// DelegateTo is a paid mutator transaction binding the contract method 0x2c1b15f4.
//
// Solidity: function delegateTo(address delegator, address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactorSession) DelegateTo(delegator common.Address, validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.DelegateTo(&_Delegation.TransactOpts, delegator, validator, token, amount)
}

// Hold is a paid mutator transaction binding the contract method 0x2dbfa4bb.
//
// Solidity: function hold(address validator, address token, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_Delegation *DelegationTransactor) Hold(opts *bind.TransactOpts, validator common.Address, token common.Address, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "hold", validator, token, amountToHold, oldHoldTimestamp, newHoldTimestamp)
}

// Hold is a paid mutator transaction binding the contract method 0x2dbfa4bb.
//
// Solidity: function hold(address validator, address token, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_Delegation *DelegationSession) Hold(validator common.Address, token common.Address, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Hold(&_Delegation.TransactOpts, validator, token, amountToHold, oldHoldTimestamp, newHoldTimestamp)
}

// Hold is a paid mutator transaction binding the contract method 0x2dbfa4bb.
//
// Solidity: function hold(address validator, address token, uint256 amountToHold, uint256 oldHoldTimestamp, uint256 newHoldTimestamp) returns()
func (_Delegation *DelegationTransactorSession) Hold(validator common.Address, token common.Address, amountToHold *big.Int, oldHoldTimestamp *big.Int, newHoldTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Hold(&_Delegation.TransactOpts, validator, token, amountToHold, oldHoldTimestamp, newHoldTimestamp)
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
// Solidity: function setContractCenter(address contractCenter) returns()
func (_Delegation *DelegationTransactor) SetContractCenter(opts *bind.TransactOpts, contractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "setContractCenter", contractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address contractCenter) returns()
func (_Delegation *DelegationSession) SetContractCenter(contractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SetContractCenter(&_Delegation.TransactOpts, contractCenter)
}

// SetContractCenter is a paid mutator transaction binding the contract method 0x5fb599d4.
//
// Solidity: function setContractCenter(address contractCenter) returns()
func (_Delegation *DelegationTransactorSession) SetContractCenter(contractCenter common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.SetContractCenter(&_Delegation.TransactOpts, contractCenter)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_Delegation *DelegationTransactor) SetFreezeTime(opts *bind.TransactOpts, withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "setFreezeTime", withdrawFreezeTime, transferFreezeTime)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_Delegation *DelegationSession) SetFreezeTime(withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.SetFreezeTime(&_Delegation.TransactOpts, withdrawFreezeTime, transferFreezeTime)
}

// SetFreezeTime is a paid mutator transaction binding the contract method 0xb71c52bd.
//
// Solidity: function setFreezeTime(uint256 withdrawFreezeTime, uint256 transferFreezeTime) returns()
func (_Delegation *DelegationTransactorSession) SetFreezeTime(withdrawFreezeTime *big.Int, transferFreezeTime *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.SetFreezeTime(&_Delegation.TransactOpts, withdrawFreezeTime, transferFreezeTime)
}

// Transfer is a paid mutator transaction binding the contract method 0xf9ce7813.
//
// Solidity: function transfer(address validator, address token, uint256 amount, address newValidator) returns()
func (_Delegation *DelegationTransactor) Transfer(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "transfer", validator, token, amount, newValidator)
}

// Transfer is a paid mutator transaction binding the contract method 0xf9ce7813.
//
// Solidity: function transfer(address validator, address token, uint256 amount, address newValidator) returns()
func (_Delegation *DelegationSession) Transfer(validator common.Address, token common.Address, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Transfer(&_Delegation.TransactOpts, validator, token, amount, newValidator)
}

// Transfer is a paid mutator transaction binding the contract method 0xf9ce7813.
//
// Solidity: function transfer(address validator, address token, uint256 amount, address newValidator) returns()
func (_Delegation *DelegationTransactorSession) Transfer(validator common.Address, token common.Address, amount *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.Transfer(&_Delegation.TransactOpts, validator, token, amount, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x108b2ccf.
//
// Solidity: function transferHold(address validator, address token, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_Delegation *DelegationTransactor) TransferHold(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "transferHold", validator, token, amount, holdTimestamp, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x108b2ccf.
//
// Solidity: function transferHold(address validator, address token, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_Delegation *DelegationSession) TransferHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.TransferHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp, newValidator)
}

// TransferHold is a paid mutator transaction binding the contract method 0x108b2ccf.
//
// Solidity: function transferHold(address validator, address token, uint256 amount, uint256 holdTimestamp, address newValidator) returns()
func (_Delegation *DelegationTransactorSession) TransferHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int, newValidator common.Address) (*types.Transaction, error) {
	return _Delegation.Contract.TransferHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp, newValidator)
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
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Delegation *DelegationTransactor) Upgrade(opts *bind.TransactOpts, newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "upgrade", newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Delegation *DelegationSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.Upgrade(&_Delegation.TransactOpts, newImpl, data)
}

// Upgrade is a paid mutator transaction binding the contract method 0xc987336c.
//
// Solidity: function upgrade(address newImpl, bytes data) returns()
func (_Delegation *DelegationTransactorSession) Upgrade(newImpl common.Address, data []byte) (*types.Transaction, error) {
	return _Delegation.Contract.Upgrade(&_Delegation.TransactOpts, newImpl, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactor) Withdraw(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "withdraw", validator, token, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationSession) Withdraw(validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Withdraw(&_Delegation.TransactOpts, validator, token, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address validator, address token, uint256 amount) returns()
func (_Delegation *DelegationTransactorSession) Withdraw(validator common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.Withdraw(&_Delegation.TransactOpts, validator, token, amount)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0x124ffa81.
//
// Solidity: function withdrawHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationTransactor) WithdrawHold(opts *bind.TransactOpts, validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.contract.Transact(opts, "withdrawHold", validator, token, amount, holdTimestamp)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0x124ffa81.
//
// Solidity: function withdrawHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationSession) WithdrawHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.WithdrawHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp)
}

// WithdrawHold is a paid mutator transaction binding the contract method 0x124ffa81.
//
// Solidity: function withdrawHold(address validator, address token, uint256 amount, uint256 holdTimestamp) returns()
func (_Delegation *DelegationTransactorSession) WithdrawHold(validator common.Address, token common.Address, amount *big.Int, holdTimestamp *big.Int) (*types.Transaction, error) {
	return _Delegation.Contract.WithdrawHold(&_Delegation.TransactOpts, validator, token, amount, holdTimestamp)
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

// DelegationFreezeTimeUpdatedIterator is returned from FilterFreezeTimeUpdated and is used to iterate over the raw logs and unpacked data for FreezeTimeUpdated events raised by the Delegation contract.
type DelegationFreezeTimeUpdatedIterator struct {
	Event *DelegationFreezeTimeUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationFreezeTimeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationFreezeTimeUpdated)
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
		it.Event = new(DelegationFreezeTimeUpdated)
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
func (it *DelegationFreezeTimeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationFreezeTimeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationFreezeTimeUpdated represents a FreezeTimeUpdated event raised by the Delegation contract.
type DelegationFreezeTimeUpdated struct {
	WithdrawFreezeTime *big.Int
	TransferFreezeTime *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterFreezeTimeUpdated is a free log retrieval operation binding the contract event 0x575c3ee3963e6c410284c700524d556e59f04881e7ff2702126adcfc33e2e22e.
//
// Solidity: event FreezeTimeUpdated(uint256 withdrawFreezeTime, uint256 transferFreezeTime)
func (_Delegation *DelegationFilterer) FilterFreezeTimeUpdated(opts *bind.FilterOpts) (*DelegationFreezeTimeUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "FreezeTimeUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationFreezeTimeUpdatedIterator{contract: _Delegation.contract, event: "FreezeTimeUpdated", logs: logs, sub: sub}, nil
}

// WatchFreezeTimeUpdated is a free log subscription operation binding the contract event 0x575c3ee3963e6c410284c700524d556e59f04881e7ff2702126adcfc33e2e22e.
//
// Solidity: event FreezeTimeUpdated(uint256 withdrawFreezeTime, uint256 transferFreezeTime)
func (_Delegation *DelegationFilterer) WatchFreezeTimeUpdated(opts *bind.WatchOpts, sink chan<- *DelegationFreezeTimeUpdated) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "FreezeTimeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationFreezeTimeUpdated)
				if err := _Delegation.contract.UnpackLog(event, "FreezeTimeUpdated", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseFreezeTimeUpdated(log types.Log) (*DelegationFreezeTimeUpdated, error) {
	event := new(DelegationFreezeTimeUpdated)
	if err := _Delegation.contract.UnpackLog(event, "FreezeTimeUpdated", log); err != nil {
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

// DelegationPenaltyAppliedToStakeIterator is returned from FilterPenaltyAppliedToStake and is used to iterate over the raw logs and unpacked data for PenaltyAppliedToStake events raised by the Delegation contract.
type DelegationPenaltyAppliedToStakeIterator struct {
	Event *DelegationPenaltyAppliedToStake // Event containing the contract specifics and raw log

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
func (it *DelegationPenaltyAppliedToStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationPenaltyAppliedToStake)
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
		it.Event = new(DelegationPenaltyAppliedToStake)
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
func (it *DelegationPenaltyAppliedToStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationPenaltyAppliedToStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationPenaltyAppliedToStake represents a PenaltyAppliedToStake event raised by the Delegation contract.
type DelegationPenaltyAppliedToStake struct {
	StakeId         [32]byte
	NewPenaltyIndex *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPenaltyAppliedToStake is a free log retrieval operation binding the contract event 0x22599c23e6ea88f2526338256e99d3dda3a4734dd2f2cf04b452c215101568e2.
//
// Solidity: event PenaltyAppliedToStake(bytes32 indexed stakeId, uint256 newPenaltyIndex)
func (_Delegation *DelegationFilterer) FilterPenaltyAppliedToStake(opts *bind.FilterOpts, stakeId [][32]byte) (*DelegationPenaltyAppliedToStakeIterator, error) {

	var stakeIdRule []interface{}
	for _, stakeIdItem := range stakeId {
		stakeIdRule = append(stakeIdRule, stakeIdItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "PenaltyAppliedToStake", stakeIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationPenaltyAppliedToStakeIterator{contract: _Delegation.contract, event: "PenaltyAppliedToStake", logs: logs, sub: sub}, nil
}

// WatchPenaltyAppliedToStake is a free log subscription operation binding the contract event 0x22599c23e6ea88f2526338256e99d3dda3a4734dd2f2cf04b452c215101568e2.
//
// Solidity: event PenaltyAppliedToStake(bytes32 indexed stakeId, uint256 newPenaltyIndex)
func (_Delegation *DelegationFilterer) WatchPenaltyAppliedToStake(opts *bind.WatchOpts, sink chan<- *DelegationPenaltyAppliedToStake, stakeId [][32]byte) (event.Subscription, error) {

	var stakeIdRule []interface{}
	for _, stakeIdItem := range stakeId {
		stakeIdRule = append(stakeIdRule, stakeIdItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "PenaltyAppliedToStake", stakeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationPenaltyAppliedToStake)
				if err := _Delegation.contract.UnpackLog(event, "PenaltyAppliedToStake", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParsePenaltyAppliedToStake(log types.Log) (*DelegationPenaltyAppliedToStake, error) {
	event := new(DelegationPenaltyAppliedToStake)
	if err := _Delegation.contract.UnpackLog(event, "PenaltyAppliedToStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationPenaltyAppliedToValidatorIterator is returned from FilterPenaltyAppliedToValidator and is used to iterate over the raw logs and unpacked data for PenaltyAppliedToValidator events raised by the Delegation contract.
type DelegationPenaltyAppliedToValidatorIterator struct {
	Event *DelegationPenaltyAppliedToValidator // Event containing the contract specifics and raw log

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
func (it *DelegationPenaltyAppliedToValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationPenaltyAppliedToValidator)
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
		it.Event = new(DelegationPenaltyAppliedToValidator)
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
func (it *DelegationPenaltyAppliedToValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationPenaltyAppliedToValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationPenaltyAppliedToValidator represents a PenaltyAppliedToValidator event raised by the Delegation contract.
type DelegationPenaltyAppliedToValidator struct {
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
func (_Delegation *DelegationFilterer) FilterPenaltyAppliedToValidator(opts *bind.FilterOpts, validator []common.Address) (*DelegationPenaltyAppliedToValidatorIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "PenaltyAppliedToValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationPenaltyAppliedToValidatorIterator{contract: _Delegation.contract, event: "PenaltyAppliedToValidator", logs: logs, sub: sub}, nil
}

// WatchPenaltyAppliedToValidator is a free log subscription operation binding the contract event 0x3b900bc82221c5a44e5d8e57d6f36ba5c3c9e10bcfdbef3441dd877a5e75981c.
//
// Solidity: event PenaltyAppliedToValidator(address indexed validator, address token, uint8 tokenType, uint256 tokenId, (uint256,uint256) validatorReserve)
func (_Delegation *DelegationFilterer) WatchPenaltyAppliedToValidator(opts *bind.WatchOpts, sink chan<- *DelegationPenaltyAppliedToValidator, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "PenaltyAppliedToValidator", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationPenaltyAppliedToValidator)
				if err := _Delegation.contract.UnpackLog(event, "PenaltyAppliedToValidator", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParsePenaltyAppliedToValidator(log types.Log) (*DelegationPenaltyAppliedToValidator, error) {
	event := new(DelegationPenaltyAppliedToValidator)
	if err := _Delegation.contract.UnpackLog(event, "PenaltyAppliedToValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationStakeAmountUpdatedIterator is returned from FilterStakeAmountUpdated and is used to iterate over the raw logs and unpacked data for StakeAmountUpdated events raised by the Delegation contract.
type DelegationStakeAmountUpdatedIterator struct {
	Event *DelegationStakeAmountUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationStakeAmountUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationStakeAmountUpdated)
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
		it.Event = new(DelegationStakeAmountUpdated)
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
func (it *DelegationStakeAmountUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationStakeAmountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationStakeAmountUpdated represents a StakeAmountUpdated event raised by the Delegation contract.
type DelegationStakeAmountUpdated struct {
	StakeId       [32]byte
	ChangedAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeAmountUpdated is a free log retrieval operation binding the contract event 0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959.
//
// Solidity: event StakeAmountUpdated(bytes32 stakeId, int256 changedAmount)
func (_Delegation *DelegationFilterer) FilterStakeAmountUpdated(opts *bind.FilterOpts) (*DelegationStakeAmountUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "StakeAmountUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationStakeAmountUpdatedIterator{contract: _Delegation.contract, event: "StakeAmountUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeAmountUpdated is a free log subscription operation binding the contract event 0x58d3d64d5ec2e281761bdebe1b59491d61dd8ff0d2fc459c5ee3b128d29f0959.
//
// Solidity: event StakeAmountUpdated(bytes32 stakeId, int256 changedAmount)
func (_Delegation *DelegationFilterer) WatchStakeAmountUpdated(opts *bind.WatchOpts, sink chan<- *DelegationStakeAmountUpdated) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "StakeAmountUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationStakeAmountUpdated)
				if err := _Delegation.contract.UnpackLog(event, "StakeAmountUpdated", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseStakeAmountUpdated(log types.Log) (*DelegationStakeAmountUpdated, error) {
	event := new(DelegationStakeAmountUpdated)
	if err := _Delegation.contract.UnpackLog(event, "StakeAmountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationStakeHoldedIterator is returned from FilterStakeHolded and is used to iterate over the raw logs and unpacked data for StakeHolded events raised by the Delegation contract.
type DelegationStakeHoldedIterator struct {
	Event *DelegationStakeHolded // Event containing the contract specifics and raw log

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
func (it *DelegationStakeHoldedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationStakeHolded)
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
		it.Event = new(DelegationStakeHolded)
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
func (it *DelegationStakeHoldedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationStakeHoldedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationStakeHolded represents a StakeHolded event raised by the Delegation contract.
type DelegationStakeHolded struct {
	StakeId [32]byte
	Stake   IDecimalDelegationCommonStake
	IsNew   bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeHolded is a free log retrieval operation binding the contract event 0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb.
//
// Solidity: event StakeHolded(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_Delegation *DelegationFilterer) FilterStakeHolded(opts *bind.FilterOpts) (*DelegationStakeHoldedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "StakeHolded")
	if err != nil {
		return nil, err
	}
	return &DelegationStakeHoldedIterator{contract: _Delegation.contract, event: "StakeHolded", logs: logs, sub: sub}, nil
}

// WatchStakeHolded is a free log subscription operation binding the contract event 0xa0a8b22bc7aca2e71ba792f9390bbc1875d1fa8b0d9a82c0158f96c7f4b89cdb.
//
// Solidity: event StakeHolded(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_Delegation *DelegationFilterer) WatchStakeHolded(opts *bind.WatchOpts, sink chan<- *DelegationStakeHolded) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "StakeHolded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationStakeHolded)
				if err := _Delegation.contract.UnpackLog(event, "StakeHolded", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseStakeHolded(log types.Log) (*DelegationStakeHolded, error) {
	event := new(DelegationStakeHolded)
	if err := _Delegation.contract.UnpackLog(event, "StakeHolded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationStakeUpdatedIterator is returned from FilterStakeUpdated and is used to iterate over the raw logs and unpacked data for StakeUpdated events raised by the Delegation contract.
type DelegationStakeUpdatedIterator struct {
	Event *DelegationStakeUpdated // Event containing the contract specifics and raw log

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
func (it *DelegationStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationStakeUpdated)
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
		it.Event = new(DelegationStakeUpdated)
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
func (it *DelegationStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationStakeUpdated represents a StakeUpdated event raised by the Delegation contract.
type DelegationStakeUpdated struct {
	StakeId [32]byte
	Stake   IDecimalDelegationCommonStake
	IsNew   bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdated is a free log retrieval operation binding the contract event 0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b.
//
// Solidity: event StakeUpdated(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_Delegation *DelegationFilterer) FilterStakeUpdated(opts *bind.FilterOpts) (*DelegationStakeUpdatedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return &DelegationStakeUpdatedIterator{contract: _Delegation.contract, event: "StakeUpdated", logs: logs, sub: sub}, nil
}

// WatchStakeUpdated is a free log subscription operation binding the contract event 0x71a822c8b2dd1c5369373dd93ec6a6b04cf7d41eb154314433c73f4f8856c03b.
//
// Solidity: event StakeUpdated(bytes32 stakeId, (address,address,address,uint256,uint256,uint8,uint256) stake, bool isNew)
func (_Delegation *DelegationFilterer) WatchStakeUpdated(opts *bind.WatchOpts, sink chan<- *DelegationStakeUpdated) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "StakeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationStakeUpdated)
				if err := _Delegation.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseStakeUpdated(log types.Log) (*DelegationStakeUpdated, error) {
	event := new(DelegationStakeUpdated)
	if err := _Delegation.contract.UnpackLog(event, "StakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTransferCompletedIterator is returned from FilterTransferCompleted and is used to iterate over the raw logs and unpacked data for TransferCompleted events raised by the Delegation contract.
type DelegationTransferCompletedIterator struct {
	Event *DelegationTransferCompleted // Event containing the contract specifics and raw log

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
func (it *DelegationTransferCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTransferCompleted)
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
		it.Event = new(DelegationTransferCompleted)
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
func (it *DelegationTransferCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTransferCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTransferCompleted represents a TransferCompleted event raised by the Delegation contract.
type DelegationTransferCompleted struct {
	StakeIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransferCompleted is a free log retrieval operation binding the contract event 0xfd987fd5c7de5139db194fdde15ddabcec1b78c3bfc832ad563ac57a9bfa9b36.
//
// Solidity: event TransferCompleted(uint256 stakeIndex)
func (_Delegation *DelegationFilterer) FilterTransferCompleted(opts *bind.FilterOpts) (*DelegationTransferCompletedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TransferCompleted")
	if err != nil {
		return nil, err
	}
	return &DelegationTransferCompletedIterator{contract: _Delegation.contract, event: "TransferCompleted", logs: logs, sub: sub}, nil
}

// WatchTransferCompleted is a free log subscription operation binding the contract event 0xfd987fd5c7de5139db194fdde15ddabcec1b78c3bfc832ad563ac57a9bfa9b36.
//
// Solidity: event TransferCompleted(uint256 stakeIndex)
func (_Delegation *DelegationFilterer) WatchTransferCompleted(opts *bind.WatchOpts, sink chan<- *DelegationTransferCompleted) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TransferCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTransferCompleted)
				if err := _Delegation.contract.UnpackLog(event, "TransferCompleted", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseTransferCompleted(log types.Log) (*DelegationTransferCompleted, error) {
	event := new(DelegationTransferCompleted)
	if err := _Delegation.contract.UnpackLog(event, "TransferCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationTransferRequestIterator is returned from FilterTransferRequest and is used to iterate over the raw logs and unpacked data for TransferRequest events raised by the Delegation contract.
type DelegationTransferRequestIterator struct {
	Event *DelegationTransferRequest // Event containing the contract specifics and raw log

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
func (it *DelegationTransferRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationTransferRequest)
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
		it.Event = new(DelegationTransferRequest)
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
func (it *DelegationTransferRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationTransferRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationTransferRequest represents a TransferRequest event raised by the Delegation contract.
type DelegationTransferRequest struct {
	StakeId     [32]byte
	StakeIndex  *big.Int
	FrozenStake IDecimalDelegationCommonFrozenStake
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTransferRequest is a free log retrieval operation binding the contract event 0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f.
//
// Solidity: event TransferRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_Delegation *DelegationFilterer) FilterTransferRequest(opts *bind.FilterOpts) (*DelegationTransferRequestIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "TransferRequest")
	if err != nil {
		return nil, err
	}
	return &DelegationTransferRequestIterator{contract: _Delegation.contract, event: "TransferRequest", logs: logs, sub: sub}, nil
}

// WatchTransferRequest is a free log subscription operation binding the contract event 0x82483dd4ee8f4c0625f25624d8973d4bc4d9bca5110e002425f5ffe8aed3f44f.
//
// Solidity: event TransferRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_Delegation *DelegationFilterer) WatchTransferRequest(opts *bind.WatchOpts, sink chan<- *DelegationTransferRequest) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "TransferRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationTransferRequest)
				if err := _Delegation.contract.UnpackLog(event, "TransferRequest", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseTransferRequest(log types.Log) (*DelegationTransferRequest, error) {
	event := new(DelegationTransferRequest)
	if err := _Delegation.contract.UnpackLog(event, "TransferRequest", log); err != nil {
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

// DelegationWithdrawCompletedIterator is returned from FilterWithdrawCompleted and is used to iterate over the raw logs and unpacked data for WithdrawCompleted events raised by the Delegation contract.
type DelegationWithdrawCompletedIterator struct {
	Event *DelegationWithdrawCompleted // Event containing the contract specifics and raw log

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
func (it *DelegationWithdrawCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationWithdrawCompleted)
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
		it.Event = new(DelegationWithdrawCompleted)
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
func (it *DelegationWithdrawCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationWithdrawCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationWithdrawCompleted represents a WithdrawCompleted event raised by the Delegation contract.
type DelegationWithdrawCompleted struct {
	StakeIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCompleted is a free log retrieval operation binding the contract event 0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7.
//
// Solidity: event WithdrawCompleted(uint256 stakeIndex)
func (_Delegation *DelegationFilterer) FilterWithdrawCompleted(opts *bind.FilterOpts) (*DelegationWithdrawCompletedIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "WithdrawCompleted")
	if err != nil {
		return nil, err
	}
	return &DelegationWithdrawCompletedIterator{contract: _Delegation.contract, event: "WithdrawCompleted", logs: logs, sub: sub}, nil
}

// WatchWithdrawCompleted is a free log subscription operation binding the contract event 0xf5c5432de34b3d8107ca8b5be5379cdfbb805a494446d6f3a0b6d1714a3e2dc7.
//
// Solidity: event WithdrawCompleted(uint256 stakeIndex)
func (_Delegation *DelegationFilterer) WatchWithdrawCompleted(opts *bind.WatchOpts, sink chan<- *DelegationWithdrawCompleted) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "WithdrawCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationWithdrawCompleted)
				if err := _Delegation.contract.UnpackLog(event, "WithdrawCompleted", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseWithdrawCompleted(log types.Log) (*DelegationWithdrawCompleted, error) {
	event := new(DelegationWithdrawCompleted)
	if err := _Delegation.contract.UnpackLog(event, "WithdrawCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationWithdrawRequestIterator is returned from FilterWithdrawRequest and is used to iterate over the raw logs and unpacked data for WithdrawRequest events raised by the Delegation contract.
type DelegationWithdrawRequestIterator struct {
	Event *DelegationWithdrawRequest // Event containing the contract specifics and raw log

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
func (it *DelegationWithdrawRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationWithdrawRequest)
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
		it.Event = new(DelegationWithdrawRequest)
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
func (it *DelegationWithdrawRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationWithdrawRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationWithdrawRequest represents a WithdrawRequest event raised by the Delegation contract.
type DelegationWithdrawRequest struct {
	StakeId     [32]byte
	StakeIndex  *big.Int
	FrozenStake IDecimalDelegationCommonFrozenStake
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequest is a free log retrieval operation binding the contract event 0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab.
//
// Solidity: event WithdrawRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_Delegation *DelegationFilterer) FilterWithdrawRequest(opts *bind.FilterOpts) (*DelegationWithdrawRequestIterator, error) {

	logs, sub, err := _Delegation.contract.FilterLogs(opts, "WithdrawRequest")
	if err != nil {
		return nil, err
	}
	return &DelegationWithdrawRequestIterator{contract: _Delegation.contract, event: "WithdrawRequest", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequest is a free log subscription operation binding the contract event 0xa7d58a1e25d2fea28ca76b72e922e949d2592d46feda6df44d1f9f801071c6ab.
//
// Solidity: event WithdrawRequest(bytes32 stakeId, uint256 stakeIndex, ((address,address,address,uint256,uint256,uint8,uint256),uint8,uint8,uint256) frozenStake)
func (_Delegation *DelegationFilterer) WatchWithdrawRequest(opts *bind.WatchOpts, sink chan<- *DelegationWithdrawRequest) (event.Subscription, error) {

	logs, sub, err := _Delegation.contract.WatchLogs(opts, "WithdrawRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationWithdrawRequest)
				if err := _Delegation.contract.UnpackLog(event, "WithdrawRequest", log); err != nil {
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
func (_Delegation *DelegationFilterer) ParseWithdrawRequest(log types.Log) (*DelegationWithdrawRequest, error) {
	event := new(DelegationWithdrawRequest)
	if err := _Delegation.contract.UnpackLog(event, "WithdrawRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
