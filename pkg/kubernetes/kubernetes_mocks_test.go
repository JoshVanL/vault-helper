// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/kubernetes/kubernetes.go

package kubernetes

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/vault/api"
	reflect "reflect"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return _m.recorder
}

// Ensure mocks base method
func (_m *MockBackend) Ensure() error {
	ret := _m.ctrl.Call(_m, "Ensure")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ensure indicates an expected call of Ensure
func (_mr *MockBackendMockRecorder) Ensure() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Ensure", reflect.TypeOf((*MockBackend)(nil).Ensure))
}

// Path mocks base method
func (_m *MockBackend) Path() string {
	ret := _m.ctrl.Call(_m, "Path")
	ret0, _ := ret[0].(string)
	return ret0
}

// Path indicates an expected call of Path
func (_mr *MockBackendMockRecorder) Path() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Path", reflect.TypeOf((*MockBackend)(nil).Path))
}

// MockVaultLogical is a mock of VaultLogical interface
type MockVaultLogical struct {
	ctrl     *gomock.Controller
	recorder *MockVaultLogicalMockRecorder
}

// MockVaultLogicalMockRecorder is the mock recorder for MockVaultLogical
type MockVaultLogicalMockRecorder struct {
	mock *MockVaultLogical
}

// NewMockVaultLogical creates a new mock instance
func NewMockVaultLogical(ctrl *gomock.Controller) *MockVaultLogical {
	mock := &MockVaultLogical{ctrl: ctrl}
	mock.recorder = &MockVaultLogicalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockVaultLogical) EXPECT() *MockVaultLogicalMockRecorder {
	return _m.recorder
}

// Write mocks base method
func (_m *MockVaultLogical) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	ret := _m.ctrl.Call(_m, "Write", path, data)
	ret0, _ := ret[0].(*api.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write
func (_mr *MockVaultLogicalMockRecorder) Write(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Write", reflect.TypeOf((*MockVaultLogical)(nil).Write), arg0, arg1)
}

// Read mocks base method
func (_m *MockVaultLogical) Read(path string) (*api.Secret, error) {
	ret := _m.ctrl.Call(_m, "Read", path)
	ret0, _ := ret[0].(*api.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (_mr *MockVaultLogicalMockRecorder) Read(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Read", reflect.TypeOf((*MockVaultLogical)(nil).Read), arg0)
}

// MockVaultSys is a mock of VaultSys interface
type MockVaultSys struct {
	ctrl     *gomock.Controller
	recorder *MockVaultSysMockRecorder
}

// MockVaultSysMockRecorder is the mock recorder for MockVaultSys
type MockVaultSysMockRecorder struct {
	mock *MockVaultSys
}

// NewMockVaultSys creates a new mock instance
func NewMockVaultSys(ctrl *gomock.Controller) *MockVaultSys {
	mock := &MockVaultSys{ctrl: ctrl}
	mock.recorder = &MockVaultSysMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockVaultSys) EXPECT() *MockVaultSysMockRecorder {
	return _m.recorder
}

// ListMounts mocks base method
func (_m *MockVaultSys) ListMounts() (map[string]*api.MountOutput, error) {
	ret := _m.ctrl.Call(_m, "ListMounts")
	ret0, _ := ret[0].(map[string]*api.MountOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMounts indicates an expected call of ListMounts
func (_mr *MockVaultSysMockRecorder) ListMounts() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ListMounts", reflect.TypeOf((*MockVaultSys)(nil).ListMounts))
}

// ListPolicies mocks base method
func (_m *MockVaultSys) ListPolicies() ([]string, error) {
	ret := _m.ctrl.Call(_m, "ListPolicies")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPolicies indicates an expected call of ListPolicies
func (_mr *MockVaultSysMockRecorder) ListPolicies() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "ListPolicies", reflect.TypeOf((*MockVaultSys)(nil).ListPolicies))
}

// Mount mocks base method
func (_m *MockVaultSys) Mount(path string, mountInfo *api.MountInput) error {
	ret := _m.ctrl.Call(_m, "Mount", path, mountInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mount indicates an expected call of Mount
func (_mr *MockVaultSysMockRecorder) Mount(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Mount", reflect.TypeOf((*MockVaultSys)(nil).Mount), arg0, arg1)
}

// PutPolicy mocks base method
func (_m *MockVaultSys) PutPolicy(name string, rules string) error {
	ret := _m.ctrl.Call(_m, "PutPolicy", name, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutPolicy indicates an expected call of PutPolicy
func (_mr *MockVaultSysMockRecorder) PutPolicy(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutPolicy", arg0, arg1)
}

// TuneMount mocks base method
func (_m *MockVaultSys) TuneMount(path string, config api.MountConfigInput) error {
	ret := _m.ctrl.Call(_m, "TuneMount", path, config)
	ret0, _ := ret[0].(error)
	return ret0
}

// TuneMount indicates an expected call of TuneMount
func (_mr *MockVaultSysMockRecorder) TuneMount(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TuneMount", arg0, arg1)
}

// GetPolicy mocks base method
func (_m *MockVaultSys) GetPolicy(name string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetPolicy", name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolicy indicates an expected call of GetPolicy
func (_mr *MockVaultSysMockRecorder) GetPolicy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPolicy", arg0)
}

// MockVaultAuth is a mock of VaultAuth interface
type MockVaultAuth struct {
	ctrl     *gomock.Controller
	recorder *MockVaultAuthMockRecorder
}

// MockVaultAuthMockRecorder is the mock recorder for MockVaultAuth
type MockVaultAuthMockRecorder struct {
	mock *MockVaultAuth
}

// NewMockVaultAuth creates a new mock instance
func NewMockVaultAuth(ctrl *gomock.Controller) *MockVaultAuth {
	mock := &MockVaultAuth{ctrl: ctrl}
	mock.recorder = &MockVaultAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockVaultAuth) EXPECT() *MockVaultAuthMockRecorder {
	return _m.recorder
}

// Token mocks base method
func (_m *MockVaultAuth) Token() VaultToken {
	ret := _m.ctrl.Call(_m, "Token")
	ret0, _ := ret[0].(VaultToken)
	return ret0
}

// Token indicates an expected call of Token
func (_mr *MockVaultAuthMockRecorder) Token() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Token")
}

// MockVaultToken is a mock of VaultToken interface
type MockVaultToken struct {
	ctrl     *gomock.Controller
	recorder *MockVaultTokenMockRecorder
}

// MockVaultTokenMockRecorder is the mock recorder for MockVaultToken
type MockVaultTokenMockRecorder struct {
	mock *MockVaultToken
}

// NewMockVaultToken creates a new mock instance
func NewMockVaultToken(ctrl *gomock.Controller) *MockVaultToken {
	mock := &MockVaultToken{ctrl: ctrl}
	mock.recorder = &MockVaultTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockVaultToken) EXPECT() *MockVaultTokenMockRecorder {
	return _m.recorder
}

// CreateOrphan mocks base method
func (_m *MockVaultToken) CreateOrphan(opts *api.TokenCreateRequest) (*api.Secret, error) {
	ret := _m.ctrl.Call(_m, "CreateOrphan", opts)
	ret0, _ := ret[0].(*api.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrphan indicates an expected call of CreateOrphan
func (_mr *MockVaultTokenMockRecorder) CreateOrphan(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateOrphan", arg0)
}

// RevokeOrphan mocks base method
func (_m *MockVaultToken) RevokeOrphan(token string) error {
	ret := _m.ctrl.Call(_m, "RevokeOrphan", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeOrphan indicates an expected call of RevokeOrphan
func (_mr *MockVaultTokenMockRecorder) RevokeOrphan(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RevokeOrphan", arg0)
}

// Lookup mocks base method
func (_m *MockVaultToken) Lookup(token string) (*api.Secret, error) {
	ret := _m.ctrl.Call(_m, "Lookup", token)
	ret0, _ := ret[0].(*api.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lookup indicates an expected call of Lookup
func (_mr *MockVaultTokenMockRecorder) Lookup(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Lookup", arg0)
}

// MockVault is a mock of Vault interface
type MockVault struct {
	ctrl     *gomock.Controller
	recorder *MockVaultMockRecorder
}

// MockVaultMockRecorder is the mock recorder for MockVault
type MockVaultMockRecorder struct {
	mock *MockVault
}

// NewMockVault creates a new mock instance
func NewMockVault(ctrl *gomock.Controller) *MockVault {
	mock := &MockVault{ctrl: ctrl}
	mock.recorder = &MockVaultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockVault) EXPECT() *MockVaultMockRecorder {
	return _m.recorder
}

// Logical mocks base method
func (_m *MockVault) Logical() VaultLogical {
	ret := _m.ctrl.Call(_m, "Logical")
	ret0, _ := ret[0].(VaultLogical)
	return ret0
}

// Logical indicates an expected call of Logical
func (_mr *MockVaultMockRecorder) Logical() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Logical")
}

// Sys mocks base method
func (_m *MockVault) Sys() VaultSys {
	ret := _m.ctrl.Call(_m, "Sys")
	ret0, _ := ret[0].(VaultSys)
	return ret0
}

// Sys indicates an expected call of Sys
func (_mr *MockVaultMockRecorder) Sys() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Sys")
}

// Auth mocks base method
func (_m *MockVault) Auth() VaultAuth {
	ret := _m.ctrl.Call(_m, "Auth")
	ret0, _ := ret[0].(VaultAuth)
	return ret0
}

// Auth indicates an expected call of Auth
func (_mr *MockVaultMockRecorder) Auth() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Auth")
}
