// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ya-breeze/geekbudgetbe/pkg/database (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	database "github.com/ya-breeze/geekbudgetbe/pkg/database"
	models "github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	goserver "github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// CreateAccount mocks base method.
func (m *MockStorage) CreateAccount(arg0 string, arg1 *goserver.AccountNoId) (goserver.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(goserver.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockStorageMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockStorage)(nil).CreateAccount), arg0, arg1)
}

// CreateBankImporter mocks base method.
func (m *MockStorage) CreateBankImporter(arg0 string, arg1 *goserver.BankImporterNoId) (goserver.BankImporter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBankImporter", arg0, arg1)
	ret0, _ := ret[0].(goserver.BankImporter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBankImporter indicates an expected call of CreateBankImporter.
func (mr *MockStorageMockRecorder) CreateBankImporter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBankImporter", reflect.TypeOf((*MockStorage)(nil).CreateBankImporter), arg0, arg1)
}

// CreateCurrency mocks base method.
func (m *MockStorage) CreateCurrency(arg0 string, arg1 *goserver.CurrencyNoId) (goserver.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCurrency", arg0, arg1)
	ret0, _ := ret[0].(goserver.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCurrency indicates an expected call of CreateCurrency.
func (mr *MockStorageMockRecorder) CreateCurrency(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCurrency", reflect.TypeOf((*MockStorage)(nil).CreateCurrency), arg0, arg1)
}

// CreateMatcher mocks base method.
func (m *MockStorage) CreateMatcher(arg0 string, arg1 goserver.MatcherNoIdInterface) (goserver.Matcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMatcher", arg0, arg1)
	ret0, _ := ret[0].(goserver.Matcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMatcher indicates an expected call of CreateMatcher.
func (mr *MockStorageMockRecorder) CreateMatcher(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMatcher", reflect.TypeOf((*MockStorage)(nil).CreateMatcher), arg0, arg1)
}

// CreateTransaction mocks base method.
func (m *MockStorage) CreateTransaction(arg0 string, arg1 goserver.TransactionNoIdInterface) (goserver.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(goserver.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockStorageMockRecorder) CreateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockStorage)(nil).CreateTransaction), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStorage) CreateUser(arg0, arg1 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStorageMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorage)(nil).CreateUser), arg0, arg1)
}

// DeleteAccount mocks base method.
func (m *MockStorage) DeleteAccount(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockStorageMockRecorder) DeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockStorage)(nil).DeleteAccount), arg0, arg1)
}

// DeleteBankImporter mocks base method.
func (m *MockStorage) DeleteBankImporter(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBankImporter", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBankImporter indicates an expected call of DeleteBankImporter.
func (mr *MockStorageMockRecorder) DeleteBankImporter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBankImporter", reflect.TypeOf((*MockStorage)(nil).DeleteBankImporter), arg0, arg1)
}

// DeleteCurrency mocks base method.
func (m *MockStorage) DeleteCurrency(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCurrency", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCurrency indicates an expected call of DeleteCurrency.
func (mr *MockStorageMockRecorder) DeleteCurrency(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCurrency", reflect.TypeOf((*MockStorage)(nil).DeleteCurrency), arg0, arg1)
}

// DeleteDuplicateTransaction mocks base method.
func (m *MockStorage) DeleteDuplicateTransaction(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDuplicateTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDuplicateTransaction indicates an expected call of DeleteDuplicateTransaction.
func (mr *MockStorageMockRecorder) DeleteDuplicateTransaction(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDuplicateTransaction", reflect.TypeOf((*MockStorage)(nil).DeleteDuplicateTransaction), arg0, arg1, arg2)
}

// DeleteMatcher mocks base method.
func (m *MockStorage) DeleteMatcher(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMatcher", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMatcher indicates an expected call of DeleteMatcher.
func (mr *MockStorageMockRecorder) DeleteMatcher(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMatcher", reflect.TypeOf((*MockStorage)(nil).DeleteMatcher), arg0, arg1)
}

// DeleteTransaction mocks base method.
func (m *MockStorage) DeleteTransaction(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTransaction indicates an expected call of DeleteTransaction.
func (mr *MockStorageMockRecorder) DeleteTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTransaction", reflect.TypeOf((*MockStorage)(nil).DeleteTransaction), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockStorage) GetAccount(arg0, arg1 string) (goserver.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(goserver.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStorageMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStorage)(nil).GetAccount), arg0, arg1)
}

// GetAccountHistory mocks base method.
func (m *MockStorage) GetAccountHistory(arg0, arg1 string) ([]goserver.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountHistory", arg0, arg1)
	ret0, _ := ret[0].([]goserver.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountHistory indicates an expected call of GetAccountHistory.
func (mr *MockStorageMockRecorder) GetAccountHistory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountHistory", reflect.TypeOf((*MockStorage)(nil).GetAccountHistory), arg0, arg1)
}

// GetAccounts mocks base method.
func (m *MockStorage) GetAccounts(arg0 string) ([]goserver.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", arg0)
	ret0, _ := ret[0].([]goserver.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockStorageMockRecorder) GetAccounts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockStorage)(nil).GetAccounts), arg0)
}

// GetAllBankImporters mocks base method.
func (m *MockStorage) GetAllBankImporters() ([]database.ImportInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBankImporters")
	ret0, _ := ret[0].([]database.ImportInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBankImporters indicates an expected call of GetAllBankImporters.
func (mr *MockStorageMockRecorder) GetAllBankImporters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBankImporters", reflect.TypeOf((*MockStorage)(nil).GetAllBankImporters))
}

// GetBankImporter mocks base method.
func (m *MockStorage) GetBankImporter(arg0, arg1 string) (goserver.BankImporter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankImporter", arg0, arg1)
	ret0, _ := ret[0].(goserver.BankImporter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBankImporter indicates an expected call of GetBankImporter.
func (mr *MockStorageMockRecorder) GetBankImporter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankImporter", reflect.TypeOf((*MockStorage)(nil).GetBankImporter), arg0, arg1)
}

// GetBankImporters mocks base method.
func (m *MockStorage) GetBankImporters(arg0 string) ([]goserver.BankImporter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankImporters", arg0)
	ret0, _ := ret[0].([]goserver.BankImporter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBankImporters indicates an expected call of GetBankImporters.
func (mr *MockStorageMockRecorder) GetBankImporters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankImporters", reflect.TypeOf((*MockStorage)(nil).GetBankImporters), arg0)
}

// GetCNBRates mocks base method.
func (m *MockStorage) GetCNBRates(arg0 time.Time) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCNBRates", arg0)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCNBRates indicates an expected call of GetCNBRates.
func (mr *MockStorageMockRecorder) GetCNBRates(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCNBRates", reflect.TypeOf((*MockStorage)(nil).GetCNBRates), arg0)
}

// GetCurrencies mocks base method.
func (m *MockStorage) GetCurrencies(arg0 string) ([]goserver.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrencies", arg0)
	ret0, _ := ret[0].([]goserver.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrencies indicates an expected call of GetCurrencies.
func (mr *MockStorageMockRecorder) GetCurrencies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrencies", reflect.TypeOf((*MockStorage)(nil).GetCurrencies), arg0)
}

// GetCurrency mocks base method.
func (m *MockStorage) GetCurrency(arg0, arg1 string) (goserver.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrency", arg0, arg1)
	ret0, _ := ret[0].(goserver.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrency indicates an expected call of GetCurrency.
func (mr *MockStorageMockRecorder) GetCurrency(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrency", reflect.TypeOf((*MockStorage)(nil).GetCurrency), arg0, arg1)
}

// GetMatcher mocks base method.
func (m *MockStorage) GetMatcher(arg0, arg1 string) (goserver.Matcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatcher", arg0, arg1)
	ret0, _ := ret[0].(goserver.Matcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatcher indicates an expected call of GetMatcher.
func (mr *MockStorageMockRecorder) GetMatcher(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatcher", reflect.TypeOf((*MockStorage)(nil).GetMatcher), arg0, arg1)
}

// GetMatcherRuntime mocks base method.
func (m *MockStorage) GetMatcherRuntime(arg0, arg1 string) (database.MatcherRuntime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatcherRuntime", arg0, arg1)
	ret0, _ := ret[0].(database.MatcherRuntime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatcherRuntime indicates an expected call of GetMatcherRuntime.
func (mr *MockStorageMockRecorder) GetMatcherRuntime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatcherRuntime", reflect.TypeOf((*MockStorage)(nil).GetMatcherRuntime), arg0, arg1)
}

// GetMatchers mocks base method.
func (m *MockStorage) GetMatchers(arg0 string) ([]goserver.Matcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchers", arg0)
	ret0, _ := ret[0].([]goserver.Matcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchers indicates an expected call of GetMatchers.
func (mr *MockStorageMockRecorder) GetMatchers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchers", reflect.TypeOf((*MockStorage)(nil).GetMatchers), arg0)
}

// GetMatchersRuntime mocks base method.
func (m *MockStorage) GetMatchersRuntime(arg0 string) ([]database.MatcherRuntime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchersRuntime", arg0)
	ret0, _ := ret[0].([]database.MatcherRuntime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchersRuntime indicates an expected call of GetMatchersRuntime.
func (mr *MockStorageMockRecorder) GetMatchersRuntime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchersRuntime", reflect.TypeOf((*MockStorage)(nil).GetMatchersRuntime), arg0)
}

// GetTransaction mocks base method.
func (m *MockStorage) GetTransaction(arg0, arg1 string) (goserver.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", arg0, arg1)
	ret0, _ := ret[0].(goserver.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockStorageMockRecorder) GetTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockStorage)(nil).GetTransaction), arg0, arg1)
}

// GetTransactions mocks base method.
func (m *MockStorage) GetTransactions(arg0 string, arg1, arg2 time.Time) ([]goserver.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions", arg0, arg1, arg2)
	ret0, _ := ret[0].([]goserver.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactions indicates an expected call of GetTransactions.
func (mr *MockStorageMockRecorder) GetTransactions(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockStorage)(nil).GetTransactions), arg0, arg1, arg2)
}

// GetUser mocks base method.
func (m *MockStorage) GetUser(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStorageMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorage)(nil).GetUser), arg0)
}

// GetUserID mocks base method.
func (m *MockStorage) GetUserID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserID indicates an expected call of GetUserID.
func (mr *MockStorageMockRecorder) GetUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserID", reflect.TypeOf((*MockStorage)(nil).GetUserID), arg0)
}

// Open mocks base method.
func (m *MockStorage) Open() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockStorageMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockStorage)(nil).Open))
}

// PutUser mocks base method.
func (m *MockStorage) PutUser(arg0 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutUser indicates an expected call of PutUser.
func (mr *MockStorageMockRecorder) PutUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutUser", reflect.TypeOf((*MockStorage)(nil).PutUser), arg0)
}

// SaveCNBRates mocks base method.
func (m *MockStorage) SaveCNBRates(arg0 map[string]float64, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCNBRates", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCNBRates indicates an expected call of SaveCNBRates.
func (mr *MockStorageMockRecorder) SaveCNBRates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCNBRates", reflect.TypeOf((*MockStorage)(nil).SaveCNBRates), arg0, arg1)
}

// UpdateAccount mocks base method.
func (m *MockStorage) UpdateAccount(arg0, arg1 string, arg2 *goserver.AccountNoId) (goserver.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(goserver.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockStorageMockRecorder) UpdateAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockStorage)(nil).UpdateAccount), arg0, arg1, arg2)
}

// UpdateBankImporter mocks base method.
func (m *MockStorage) UpdateBankImporter(arg0, arg1 string, arg2 goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBankImporter", arg0, arg1, arg2)
	ret0, _ := ret[0].(goserver.BankImporter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBankImporter indicates an expected call of UpdateBankImporter.
func (mr *MockStorageMockRecorder) UpdateBankImporter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBankImporter", reflect.TypeOf((*MockStorage)(nil).UpdateBankImporter), arg0, arg1, arg2)
}

// UpdateCurrency mocks base method.
func (m *MockStorage) UpdateCurrency(arg0, arg1 string, arg2 *goserver.CurrencyNoId) (goserver.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCurrency", arg0, arg1, arg2)
	ret0, _ := ret[0].(goserver.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCurrency indicates an expected call of UpdateCurrency.
func (mr *MockStorageMockRecorder) UpdateCurrency(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCurrency", reflect.TypeOf((*MockStorage)(nil).UpdateCurrency), arg0, arg1, arg2)
}

// UpdateMatcher mocks base method.
func (m *MockStorage) UpdateMatcher(arg0, arg1 string, arg2 goserver.MatcherNoIdInterface) (goserver.Matcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMatcher", arg0, arg1, arg2)
	ret0, _ := ret[0].(goserver.Matcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMatcher indicates an expected call of UpdateMatcher.
func (mr *MockStorageMockRecorder) UpdateMatcher(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMatcher", reflect.TypeOf((*MockStorage)(nil).UpdateMatcher), arg0, arg1, arg2)
}

// UpdateTransaction mocks base method.
func (m *MockStorage) UpdateTransaction(arg0, arg1 string, arg2 goserver.TransactionNoIdInterface) (goserver.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(goserver.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTransaction indicates an expected call of UpdateTransaction.
func (mr *MockStorageMockRecorder) UpdateTransaction(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTransaction", reflect.TypeOf((*MockStorage)(nil).UpdateTransaction), arg0, arg1, arg2)
}
