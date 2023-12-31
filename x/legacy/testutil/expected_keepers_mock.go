// Code generated by MockGen. DO NOT EDIT.
// Source: bitbucket.org/decimalteam/go-smart-node/x/legacy/types (interfaces: BankKeeper,NftKeeper,MultisigKeeper,ValidatorKeeper)

// Package testutil is a generated GoMock package.
package testutil

import (
	reflect "reflect"

	types "bitbucket.org/decimalteam/go-smart-node/x/multisig/types"
	types0 "bitbucket.org/decimalteam/go-smart-node/x/nft/types"
	types1 "bitbucket.org/decimalteam/go-smart-node/x/validator/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockBankKeeper is a mock of BankKeeper interface.
type MockBankKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockBankKeeperMockRecorder
}

// MockBankKeeperMockRecorder is the mock recorder for MockBankKeeper.
type MockBankKeeperMockRecorder struct {
	mock *MockBankKeeper
}

// NewMockBankKeeper creates a new mock instance.
func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	mock := &MockBankKeeper{ctrl: ctrl}
	mock.recorder = &MockBankKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankKeeper) EXPECT() *MockBankKeeperMockRecorder {
	return m.recorder
}

// GetAllBalances mocks base method.
func (m *MockBankKeeper) GetAllBalances(arg0 types2.Context, arg1 types2.AccAddress) types2.Coins {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalances", arg0, arg1)
	ret0, _ := ret[0].(types2.Coins)
	return ret0
}

// GetAllBalances indicates an expected call of GetAllBalances.
func (mr *MockBankKeeperMockRecorder) GetAllBalances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalances", reflect.TypeOf((*MockBankKeeper)(nil).GetAllBalances), arg0, arg1)
}

// SendCoins mocks base method.
func (m *MockBankKeeper) SendCoins(arg0 types2.Context, arg1, arg2 types2.AccAddress, arg3 types2.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoins", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoins indicates an expected call of SendCoins.
func (mr *MockBankKeeperMockRecorder) SendCoins(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoins", reflect.TypeOf((*MockBankKeeper)(nil).SendCoins), arg0, arg1, arg2, arg3)
}

// SendCoinsFromModuleToAccount mocks base method.
func (m *MockBankKeeper) SendCoinsFromModuleToAccount(arg0 types2.Context, arg1 string, arg2 types2.AccAddress, arg3 types2.Coins) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCoinsFromModuleToAccount", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendCoinsFromModuleToAccount indicates an expected call of SendCoinsFromModuleToAccount.
func (mr *MockBankKeeperMockRecorder) SendCoinsFromModuleToAccount(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCoinsFromModuleToAccount", reflect.TypeOf((*MockBankKeeper)(nil).SendCoinsFromModuleToAccount), arg0, arg1, arg2, arg3)
}

// MockNftKeeper is a mock of NftKeeper interface.
type MockNftKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockNftKeeperMockRecorder
}

// MockNftKeeperMockRecorder is the mock recorder for MockNftKeeper.
type MockNftKeeperMockRecorder struct {
	mock *MockNftKeeper
}

// NewMockNftKeeper creates a new mock instance.
func NewMockNftKeeper(ctrl *gomock.Controller) *MockNftKeeper {
	mock := &MockNftKeeper{ctrl: ctrl}
	mock.recorder = &MockNftKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNftKeeper) EXPECT() *MockNftKeeperMockRecorder {
	return m.recorder
}

// GetSubTokens mocks base method.
func (m *MockNftKeeper) GetSubTokens(arg0 types2.Context, arg1 string) []types0.SubToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubTokens", arg0, arg1)
	ret0, _ := ret[0].([]types0.SubToken)
	return ret0
}

// GetSubTokens indicates an expected call of GetSubTokens.
func (mr *MockNftKeeperMockRecorder) GetSubTokens(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubTokens", reflect.TypeOf((*MockNftKeeper)(nil).GetSubTokens), arg0, arg1)
}

// GetToken mocks base method.
func (m *MockNftKeeper) GetToken(arg0 types2.Context, arg1 string) (types0.Token, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken", arg0, arg1)
	ret0, _ := ret[0].(types0.Token)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetToken indicates an expected call of GetToken.
func (mr *MockNftKeeperMockRecorder) GetToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockNftKeeper)(nil).GetToken), arg0, arg1)
}

// ReplaceSubTokenOwner mocks base method.
func (m *MockNftKeeper) ReplaceSubTokenOwner(arg0 types2.Context, arg1 string, arg2 uint32, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceSubTokenOwner", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceSubTokenOwner indicates an expected call of ReplaceSubTokenOwner.
func (mr *MockNftKeeperMockRecorder) ReplaceSubTokenOwner(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceSubTokenOwner", reflect.TypeOf((*MockNftKeeper)(nil).ReplaceSubTokenOwner), arg0, arg1, arg2, arg3)
}

// MockMultisigKeeper is a mock of MultisigKeeper interface.
type MockMultisigKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockMultisigKeeperMockRecorder
}

// MockMultisigKeeperMockRecorder is the mock recorder for MockMultisigKeeper.
type MockMultisigKeeperMockRecorder struct {
	mock *MockMultisigKeeper
}

// NewMockMultisigKeeper creates a new mock instance.
func NewMockMultisigKeeper(ctrl *gomock.Controller) *MockMultisigKeeper {
	mock := &MockMultisigKeeper{ctrl: ctrl}
	mock.recorder = &MockMultisigKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMultisigKeeper) EXPECT() *MockMultisigKeeperMockRecorder {
	return m.recorder
}

// GetWallet mocks base method.
func (m *MockMultisigKeeper) GetWallet(arg0 types2.Context, arg1 string) (types.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWallet", arg0, arg1)
	ret0, _ := ret[0].(types.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWallet indicates an expected call of GetWallet.
func (mr *MockMultisigKeeperMockRecorder) GetWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWallet", reflect.TypeOf((*MockMultisigKeeper)(nil).GetWallet), arg0, arg1)
}

// SetWallet mocks base method.
func (m *MockMultisigKeeper) SetWallet(arg0 types2.Context, arg1 types.Wallet) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetWallet", arg0, arg1)
}

// SetWallet indicates an expected call of SetWallet.
func (mr *MockMultisigKeeperMockRecorder) SetWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWallet", reflect.TypeOf((*MockMultisigKeeper)(nil).SetWallet), arg0, arg1)
}

// MockValidatorKeeper is a mock of ValidatorKeeper interface.
type MockValidatorKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorKeeperMockRecorder
}

// MockValidatorKeeperMockRecorder is the mock recorder for MockValidatorKeeper.
type MockValidatorKeeperMockRecorder struct {
	mock *MockValidatorKeeper
}

// NewMockValidatorKeeper creates a new mock instance.
func NewMockValidatorKeeper(ctrl *gomock.Controller) *MockValidatorKeeper {
	mock := &MockValidatorKeeper{ctrl: ctrl}
	mock.recorder = &MockValidatorKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidatorKeeper) EXPECT() *MockValidatorKeeperMockRecorder {
	return m.recorder
}

// GetValidator mocks base method.
func (m *MockValidatorKeeper) GetValidator(arg0 types2.Context, arg1 types2.ValAddress) (types1.Validator, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", arg0, arg1)
	ret0, _ := ret[0].(types1.Validator)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator.
func (mr *MockValidatorKeeperMockRecorder) GetValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockValidatorKeeper)(nil).GetValidator), arg0, arg1)
}

// SetValidator mocks base method.
func (m *MockValidatorKeeper) SetValidator(arg0 types2.Context, arg1 types1.Validator) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetValidator", arg0, arg1)
}

// SetValidator indicates an expected call of SetValidator.
func (mr *MockValidatorKeeperMockRecorder) SetValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValidator", reflect.TypeOf((*MockValidatorKeeper)(nil).SetValidator), arg0, arg1)
}
