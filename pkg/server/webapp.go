package server

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

type RootRouter struct {
	commit  string
	logger  *slog.Logger
	cfg     *config.Config
	db      database.Storage
	cookies *sessions.CookieStore
}

func NewRootRouter(
	commit string, logger *slog.Logger, cfg *config.Config, db database.Storage,
) *RootRouter {
	return &RootRouter{
		commit:  commit,
		logger:  logger,
		cfg:     cfg,
		db:      db,
		cookies: sessions.NewCookieStore([]byte("SESSION_KEY")),
	}
}

func (r *RootRouter) Routes() goserver.Routes {
	return goserver.Routes{
		"RootPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: r.homeHandler,
		},
		"Login": goserver.Route{
			Method:      "POST",
			Pattern:     "/",
			HandlerFunc: r.loginHandler,
		},
		"AboutPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/about",
			HandlerFunc: r.aboutHandler,
		},
		"BankImporters": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/bank-importers",
			HandlerFunc: r.bankImportersHandler,
		},
		"Matchers": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/matchers",
			HandlerFunc: r.matchersHandler,
		},
		"Unprocessed": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/unprocessed",
			HandlerFunc: r.unprocessedHandler,
		},
	}
}

func (r *RootRouter) loadTemplates() (*template.Template, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"formatTime": utils.FormatTime,
	}).ParseGlob(filepath.Join("webapp", "templates", "*.html"))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (r *RootRouter) homeHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

		accounts, err := r.db.GetAccounts(userID)
		if err != nil {
			r.logger.Error("Failed to get accounts", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["Accounts"] = accounts

		currencies, err := r.db.GetCurrencies(userID)
		if err != nil {
			r.logger.Error("Failed to get currencies", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		a := NewAggregationsAPIServiceImpl(r.logger, r.db)
		dateFrom := utils.RoundToGranularity(time.Now(), utils.GranularityYear, false)
		dateTo := utils.RoundToGranularity(time.Now(), utils.GranularityMonth, true)

		expenses, err := a.GetAggregatedExpenses(req.Context(), userID, dateFrom, dateTo, "")
		if err != nil {
			r.logger.Error("Failed to get aggregated expenses", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		webAggregation := WebAggregation{
			From:        expenses.From,
			To:          expenses.To,
			Granularity: expenses.Granularity,
			Intervals:   expenses.Intervals,
			Currencies:  make([]WebCurrencyAggregation, 0, len(expenses.Currencies)),
		}
		for _, currency := range expenses.Currencies {
			webCurrency := WebCurrencyAggregation{
				CurrencyId:   currency.CurrencyId,
				CurrencyName: utils.GetCurrency(currency.CurrencyId, currencies).Name,
				Intervals:    expenses.Intervals,
			}
			if webCurrency.CurrencyName == "" {
				webCurrency.CurrencyName = "Unknown"
			}

			for _, account := range currency.Accounts {
				webAccount := AccountAggregation{
					AccountId:   account.AccountId,
					AccountName: utils.GetAccount(account.AccountId, accounts).Name,
					Amounts:     account.Amounts,
				}
				if webAccount.AccountName == "" {
					webAccount.AccountName = "Unknown"
				}
				webCurrency.Accounts = append(webCurrency.Accounts, webAccount)
			}

			webAggregation.Currencies = append(webAggregation.Currencies, webCurrency)
		}

		data["Expenses"] = &webAggregation
	}

	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RootRouter) aboutHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "about.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RootRouter) loginHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := req.Form.Get("username")
	password := req.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	userID, err := r.db.GetUserID(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		r.logger.Warn("failed to get user ID", "username", username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := r.cookies.Get(req, "session-name")
	session.Values["userID"] = userID
	err = session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "GeekBudget API",
		"UserID": userID,
	}

	if err := tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RootRouter) bankImportersHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

		// accounts, err := r.db.GetAccounts(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get accounts", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// currencies, err := r.db.GetCurrencies(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get currencies", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		bankimporters, err := r.db.GetBankImporters(userID)
		if err != nil {
			r.logger.Error("Failed to get bank importers", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.logger.Info("Bank importers", "bankimporters", bankimporters)

		data["BankImporters"] = &bankimporters
	}

	if err := tmpl.ExecuteTemplate(w, "bank_importers.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RootRouter) matchersHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

		// accounts, err := r.db.GetAccounts(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get accounts", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// currencies, err := r.db.GetCurrencies(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get currencies", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		matchers, err := r.db.GetMatchers(userID)
		if err != nil {
			r.logger.Error("Failed to get matchers", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.logger.Info("Matchers", "matchers", matchers)

		data["Matchers"] = &matchers
	}

	if err := tmpl.ExecuteTemplate(w, "matchers.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RootRouter) unprocessedHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

		// accounts, err := r.db.GetAccounts(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get accounts", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// currencies, err := r.db.GetCurrencies(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get currencies", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		unprocessed, err := r.db.GetTransactions(userID, time.Time{}, time.Time{})
		if err != nil {
			r.logger.Error("Failed to get unprocessed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data["Unprocessed"] = &unprocessed
	}

	if err := tmpl.ExecuteTemplate(w, "unprocessed.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
