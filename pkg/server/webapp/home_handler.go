package webapp

import (
	"net/http"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:funlen,cyclop
func (r *WebAppRouter) homeHandler(w http.ResponseWriter, req *http.Request) {
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

		a := api.NewAggregationsAPIServiceImpl(r.logger, r.db)
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
				CurrencyID:   currency.CurrencyId,
				CurrencyName: utils.GetCurrency(currency.CurrencyId, currencies).Name,
				Intervals:    expenses.Intervals,
			}
			if webCurrency.CurrencyName == "" {
				webCurrency.CurrencyName = "Unknown"
			}

			for _, account := range currency.Accounts {
				webAccount := AccountAggregation{
					AccountID:   account.AccountId,
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

	if err := tmpl.ExecuteTemplate(w, "home.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
