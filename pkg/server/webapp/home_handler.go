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
	data := utils.CreateTemplateData(req, "home")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
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

	dateFrom, dateTo, err := getTimeRange(req, utils.GranularityMonth)
	if err != nil {
		r.logger.Error("Failed to get time range", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dateFrom = dateFrom.AddDate(0, -12, 0)

	outputCurrencyName := req.URL.Query().Get("currency")
	expenses, err := a.GetAggregatedExpenses(req.Context(), userID, dateFrom, dateTo, outputCurrencyName)
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

	if utils.IsMobile(req.Header.Get("User-Agent")) {
		data["Template"] = "home_mobile.tpl"
	} else {
		data["Template"] = "home.tpl"
	}

	templateName, ok := data["Template"].(string)
	if !ok {
		r.logger.Error("Failed to assert template name")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		r.logger.Warn("failed to execute template", "error", err, "template", templateName)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
