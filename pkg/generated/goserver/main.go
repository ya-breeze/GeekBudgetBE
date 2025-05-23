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
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

type CustomControllers struct {
	AccountsAPIService                AccountsAPIService
	AggregationsAPIService            AggregationsAPIService
	AuthAPIService                    AuthAPIService
	BankImportersAPIService           BankImportersAPIService
	BudgetItemsAPIService             BudgetItemsAPIService
	CurrenciesAPIService              CurrenciesAPIService
	ExportAPIService                  ExportAPIService
	ImportAPIService                  ImportAPIService
	MatchersAPIService                MatchersAPIService
	NotificationsAPIService           NotificationsAPIService
	TransactionsAPIService            TransactionsAPIService
	UnprocessedTransactionsAPIService UnprocessedTransactionsAPIService
	UserAPIService                    UserAPIService
}

func Serve(ctx context.Context, logger *slog.Logger, cfg *config.Config,
	controllers CustomControllers, extraRouters []Router, middlewares ...mux.MiddlewareFunc) (net.Addr, chan int, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to listen: %w", err)
	}
	logger.Info(fmt.Sprintf("Listening at port %d...", listener.Addr().(*net.TCPAddr).Port))

	AccountsAPIService := NewAccountsAPIService()
	if controllers.AccountsAPIService != nil {
		AccountsAPIService = controllers.AccountsAPIService
	}
	AccountsAPIController := NewAccountsAPIController(AccountsAPIService)

	AggregationsAPIService := NewAggregationsAPIService()
	if controllers.AggregationsAPIService != nil {
		AggregationsAPIService = controllers.AggregationsAPIService
	}
	AggregationsAPIController := NewAggregationsAPIController(AggregationsAPIService)

	AuthAPIService := NewAuthAPIService()
	if controllers.AuthAPIService != nil {
		AuthAPIService = controllers.AuthAPIService
	}
	AuthAPIController := NewAuthAPIController(AuthAPIService)

	BankImportersAPIService := NewBankImportersAPIService()
	if controllers.BankImportersAPIService != nil {
		BankImportersAPIService = controllers.BankImportersAPIService
	}
	BankImportersAPIController := NewBankImportersAPIController(BankImportersAPIService)

	BudgetItemsAPIService := NewBudgetItemsAPIService()
	if controllers.BudgetItemsAPIService != nil {
		BudgetItemsAPIService = controllers.BudgetItemsAPIService
	}
	BudgetItemsAPIController := NewBudgetItemsAPIController(BudgetItemsAPIService)

	CurrenciesAPIService := NewCurrenciesAPIService()
	if controllers.CurrenciesAPIService != nil {
		CurrenciesAPIService = controllers.CurrenciesAPIService
	}
	CurrenciesAPIController := NewCurrenciesAPIController(CurrenciesAPIService)

	ExportAPIService := NewExportAPIService()
	if controllers.ExportAPIService != nil {
		ExportAPIService = controllers.ExportAPIService
	}
	ExportAPIController := NewExportAPIController(ExportAPIService)

	ImportAPIService := NewImportAPIService()
	if controllers.ImportAPIService != nil {
		ImportAPIService = controllers.ImportAPIService
	}
	ImportAPIController := NewImportAPIController(ImportAPIService)

	MatchersAPIService := NewMatchersAPIService()
	if controllers.MatchersAPIService != nil {
		MatchersAPIService = controllers.MatchersAPIService
	}
	MatchersAPIController := NewMatchersAPIController(MatchersAPIService)

	NotificationsAPIService := NewNotificationsAPIService()
	if controllers.NotificationsAPIService != nil {
		NotificationsAPIService = controllers.NotificationsAPIService
	}
	NotificationsAPIController := NewNotificationsAPIController(NotificationsAPIService)

	TransactionsAPIService := NewTransactionsAPIService()
	if controllers.TransactionsAPIService != nil {
		TransactionsAPIService = controllers.TransactionsAPIService
	}
	TransactionsAPIController := NewTransactionsAPIController(TransactionsAPIService)

	UnprocessedTransactionsAPIService := NewUnprocessedTransactionsAPIService()
	if controllers.UnprocessedTransactionsAPIService != nil {
		UnprocessedTransactionsAPIService = controllers.UnprocessedTransactionsAPIService
	}
	UnprocessedTransactionsAPIController := NewUnprocessedTransactionsAPIController(UnprocessedTransactionsAPIService)

	UserAPIService := NewUserAPIService()
	if controllers.UserAPIService != nil {
		UserAPIService = controllers.UserAPIService
	}
	UserAPIController := NewUserAPIController(UserAPIService)

	routers := append(extraRouters, AccountsAPIController, AggregationsAPIController, AuthAPIController, BankImportersAPIController, BudgetItemsAPIController, CurrenciesAPIController, ExportAPIController, ImportAPIController, MatchersAPIController, NotificationsAPIController, TransactionsAPIController, UnprocessedTransactionsAPIController, UserAPIController)
	router := NewRouter(routers...)

	router.Use(middlewares...)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	server := &http.Server{
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(router),
	}

	go func() {
		server.Serve(listener)
	}()

	finishChan := make(chan int, 1)
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...")
		timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(timeout)
		finishChan <- 1
		logger.Info("Server stopped")
	}()

	return listener.Addr(), finishChan, nil
}
