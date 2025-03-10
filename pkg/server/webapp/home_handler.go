package webapp

import (
	"net/http"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:funlen,cyclop,gocognit
func (r *WebAppRouter) homeHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{}

	session, err := r.cookies.Get(req, "session-name")
	if err != nil {
		r.logger.Error("Failed to get session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		// dateFrom := utils.RoundToGranularity(time.Now(), utils.GranularityYear, false)
		// dateTo := utils.RoundToGranularity(time.Now(), utils.GranularityMonth, true)

		dateFrom, dateTo, err := getTimeRange(req, utils.GranularityYear)
		if err != nil {
			r.logger.Error("Failed to get time range", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
				Total:        make([]float64, len(expenses.Intervals)),
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
				for _, amount := range account.Amounts {
					webAccount.TotalForYear += amount
				}
				for i := range expenses.Intervals {
					webCurrency.Total[i] += account.Amounts[i]
				}

				webCurrency.Accounts = append(webCurrency.Accounts, webAccount)
			}

			webAggregation.Currencies = append(webAggregation.Currencies, webCurrency)
		}
		data["Expenses"] = &webAggregation
		data["From"] = dateFrom
		data["To"] = dateTo
		data["Current"] = dateFrom.Unix()
		data["Last"] = time.Date(
			dateFrom.Year(), 1, 1, 0, 0, 0, 0, dateFrom.Location(),
		).AddDate(-1, 0, 0).Unix()
		data["Next"] = dateTo.Unix()
	}

	if err := tmpl.ExecuteTemplate(w, "home.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
