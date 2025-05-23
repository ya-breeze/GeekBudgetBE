// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Geek Budget - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

import (
	"context"
	"net/http"
	"os"
	"time"
)

// AccountsAPIRouter defines the required methods for binding the api requests to a responses for the AccountsAPI
// The AccountsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a AccountsAPIServicer to perform the required actions, then write the service results to the http response.
type AccountsAPIRouter interface {
	GetAccountHistory(http.ResponseWriter, *http.Request)
	GetAccounts(http.ResponseWriter, *http.Request)
	CreateAccount(http.ResponseWriter, *http.Request)
	GetAccount(http.ResponseWriter, *http.Request)
	UpdateAccount(http.ResponseWriter, *http.Request)
	DeleteAccount(http.ResponseWriter, *http.Request)
}

// AggregationsAPIRouter defines the required methods for binding the api requests to a responses for the AggregationsAPI
// The AggregationsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a AggregationsAPIServicer to perform the required actions, then write the service results to the http response.
type AggregationsAPIRouter interface {
	GetBalances(http.ResponseWriter, *http.Request)
	GetExpenses(http.ResponseWriter, *http.Request)
	GetIncomes(http.ResponseWriter, *http.Request)
}

// AuthAPIRouter defines the required methods for binding the api requests to a responses for the AuthAPI
// The AuthAPIRouter implementation should parse necessary information from the http request,
// pass the data to a AuthAPIServicer to perform the required actions, then write the service results to the http response.
type AuthAPIRouter interface {
	Authorize(http.ResponseWriter, *http.Request)
}

// BankImportersAPIRouter defines the required methods for binding the api requests to a responses for the BankImportersAPI
// The BankImportersAPIRouter implementation should parse necessary information from the http request,
// pass the data to a BankImportersAPIServicer to perform the required actions, then write the service results to the http response.
type BankImportersAPIRouter interface {
	GetBankImporters(http.ResponseWriter, *http.Request)
	CreateBankImporter(http.ResponseWriter, *http.Request)
	UpdateBankImporter(http.ResponseWriter, *http.Request)
	DeleteBankImporter(http.ResponseWriter, *http.Request)
	FetchBankImporter(http.ResponseWriter, *http.Request)
	UploadBankImporter(http.ResponseWriter, *http.Request)
}

// BudgetItemsAPIRouter defines the required methods for binding the api requests to a responses for the BudgetItemsAPI
// The BudgetItemsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a BudgetItemsAPIServicer to perform the required actions, then write the service results to the http response.
type BudgetItemsAPIRouter interface {
	GetBudgetItems(http.ResponseWriter, *http.Request)
	CreateBudgetItem(http.ResponseWriter, *http.Request)
	GetBudgetItem(http.ResponseWriter, *http.Request)
	UpdateBudgetItem(http.ResponseWriter, *http.Request)
	DeleteBudgetItem(http.ResponseWriter, *http.Request)
}

// CurrenciesAPIRouter defines the required methods for binding the api requests to a responses for the CurrenciesAPI
// The CurrenciesAPIRouter implementation should parse necessary information from the http request,
// pass the data to a CurrenciesAPIServicer to perform the required actions, then write the service results to the http response.
type CurrenciesAPIRouter interface {
	GetCurrencies(http.ResponseWriter, *http.Request)
	CreateCurrency(http.ResponseWriter, *http.Request)
	UpdateCurrency(http.ResponseWriter, *http.Request)
	DeleteCurrency(http.ResponseWriter, *http.Request)
}

// ExportAPIRouter defines the required methods for binding the api requests to a responses for the ExportAPI
// The ExportAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ExportAPIServicer to perform the required actions, then write the service results to the http response.
type ExportAPIRouter interface {
	Export(http.ResponseWriter, *http.Request)
}

// ImportAPIRouter defines the required methods for binding the api requests to a responses for the ImportAPI
// The ImportAPIRouter implementation should parse necessary information from the http request,
// pass the data to a ImportAPIServicer to perform the required actions, then write the service results to the http response.
type ImportAPIRouter interface {
	CallImport(http.ResponseWriter, *http.Request)
}

// MatchersAPIRouter defines the required methods for binding the api requests to a responses for the MatchersAPI
// The MatchersAPIRouter implementation should parse necessary information from the http request,
// pass the data to a MatchersAPIServicer to perform the required actions, then write the service results to the http response.
type MatchersAPIRouter interface {
	GetMatchers(http.ResponseWriter, *http.Request)
	CreateMatcher(http.ResponseWriter, *http.Request)
	UpdateMatcher(http.ResponseWriter, *http.Request)
	DeleteMatcher(http.ResponseWriter, *http.Request)
	CheckMatcher(http.ResponseWriter, *http.Request)
}

// NotificationsAPIRouter defines the required methods for binding the api requests to a responses for the NotificationsAPI
// The NotificationsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a NotificationsAPIServicer to perform the required actions, then write the service results to the http response.
type NotificationsAPIRouter interface {
	GetNotifications(http.ResponseWriter, *http.Request)
	DeleteNotification(http.ResponseWriter, *http.Request)
}

// TransactionsAPIRouter defines the required methods for binding the api requests to a responses for the TransactionsAPI
// The TransactionsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a TransactionsAPIServicer to perform the required actions, then write the service results to the http response.
type TransactionsAPIRouter interface {
	GetTransactions(http.ResponseWriter, *http.Request)
	CreateTransaction(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	UpdateTransaction(http.ResponseWriter, *http.Request)
	DeleteTransaction(http.ResponseWriter, *http.Request)
}

// UnprocessedTransactionsAPIRouter defines the required methods for binding the api requests to a responses for the UnprocessedTransactionsAPI
// The UnprocessedTransactionsAPIRouter implementation should parse necessary information from the http request,
// pass the data to a UnprocessedTransactionsAPIServicer to perform the required actions, then write the service results to the http response.
type UnprocessedTransactionsAPIRouter interface {
	GetUnprocessedTransactions(http.ResponseWriter, *http.Request)
	ConvertUnprocessedTransaction(http.ResponseWriter, *http.Request)
	DeleteUnprocessedTransaction(http.ResponseWriter, *http.Request)
}

// UserAPIRouter defines the required methods for binding the api requests to a responses for the UserAPI
// The UserAPIRouter implementation should parse necessary information from the http request,
// pass the data to a UserAPIServicer to perform the required actions, then write the service results to the http response.
type UserAPIRouter interface {
	GetUser(http.ResponseWriter, *http.Request)
}

// AccountsAPIServicer defines the api actions for the AccountsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AccountsAPIServicer interface {
	GetAccountHistory(context.Context, string) (ImplResponse, error)
	GetAccounts(context.Context) (ImplResponse, error)
	CreateAccount(context.Context, AccountNoId) (ImplResponse, error)
	GetAccount(context.Context, string) (ImplResponse, error)
	UpdateAccount(context.Context, string, AccountNoId) (ImplResponse, error)
	DeleteAccount(context.Context, string) (ImplResponse, error)
}

// AggregationsAPIServicer defines the api actions for the AggregationsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AggregationsAPIServicer interface {
	GetBalances(context.Context, time.Time, time.Time, string) (ImplResponse, error)
	GetExpenses(context.Context, time.Time, time.Time, string) (ImplResponse, error)
	GetIncomes(context.Context, time.Time, time.Time, string) (ImplResponse, error)
}

// AuthAPIServicer defines the api actions for the AuthAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AuthAPIServicer interface {
	Authorize(context.Context, AuthData) (ImplResponse, error)
}

// BankImportersAPIServicer defines the api actions for the BankImportersAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type BankImportersAPIServicer interface {
	GetBankImporters(context.Context) (ImplResponse, error)
	CreateBankImporter(context.Context, BankImporterNoId) (ImplResponse, error)
	UpdateBankImporter(context.Context, string, BankImporterNoId) (ImplResponse, error)
	DeleteBankImporter(context.Context, string) (ImplResponse, error)
	FetchBankImporter(context.Context, string) (ImplResponse, error)
	UploadBankImporter(context.Context, string, string, *os.File) (ImplResponse, error)
}

// BudgetItemsAPIServicer defines the api actions for the BudgetItemsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type BudgetItemsAPIServicer interface {
	GetBudgetItems(context.Context) (ImplResponse, error)
	CreateBudgetItem(context.Context, BudgetItemNoId) (ImplResponse, error)
	GetBudgetItem(context.Context, string) (ImplResponse, error)
	UpdateBudgetItem(context.Context, string, BudgetItemNoId) (ImplResponse, error)
	DeleteBudgetItem(context.Context, string) (ImplResponse, error)
}

// CurrenciesAPIServicer defines the api actions for the CurrenciesAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type CurrenciesAPIServicer interface {
	GetCurrencies(context.Context) (ImplResponse, error)
	CreateCurrency(context.Context, CurrencyNoId) (ImplResponse, error)
	UpdateCurrency(context.Context, string, CurrencyNoId) (ImplResponse, error)
	DeleteCurrency(context.Context, string) (ImplResponse, error)
}

// ExportAPIServicer defines the api actions for the ExportAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ExportAPIServicer interface {
	Export(context.Context) (ImplResponse, error)
}

// ImportAPIServicer defines the api actions for the ImportAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ImportAPIServicer interface {
	CallImport(context.Context, WholeUserData) (ImplResponse, error)
}

// MatchersAPIServicer defines the api actions for the MatchersAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MatchersAPIServicer interface {
	GetMatchers(context.Context) (ImplResponse, error)
	CreateMatcher(context.Context, MatcherNoId) (ImplResponse, error)
	UpdateMatcher(context.Context, string, MatcherNoId) (ImplResponse, error)
	DeleteMatcher(context.Context, string) (ImplResponse, error)
	CheckMatcher(context.Context, CheckMatcherRequest) (ImplResponse, error)
}

// NotificationsAPIServicer defines the api actions for the NotificationsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type NotificationsAPIServicer interface {
	GetNotifications(context.Context) (ImplResponse, error)
	DeleteNotification(context.Context, string) (ImplResponse, error)
}

// TransactionsAPIServicer defines the api actions for the TransactionsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TransactionsAPIServicer interface {
	GetTransactions(context.Context, string, float64, float64, time.Time, time.Time) (ImplResponse, error)
	CreateTransaction(context.Context, TransactionNoId) (ImplResponse, error)
	GetTransaction(context.Context, string) (ImplResponse, error)
	UpdateTransaction(context.Context, string, TransactionNoId) (ImplResponse, error)
	DeleteTransaction(context.Context, string) (ImplResponse, error)
}

// UnprocessedTransactionsAPIServicer defines the api actions for the UnprocessedTransactionsAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UnprocessedTransactionsAPIServicer interface {
	GetUnprocessedTransactions(context.Context) (ImplResponse, error)
	ConvertUnprocessedTransaction(context.Context, string, TransactionNoId) (ImplResponse, error)
	DeleteUnprocessedTransaction(context.Context, string, string) (ImplResponse, error)
}

// UserAPIServicer defines the api actions for the UserAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type UserAPIServicer interface {
	GetUser(context.Context) (ImplResponse, error)
}
