package webapp

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

func getTimeRange(req *http.Request, granularity utils.Granularity) (time.Time, time.Time, error) {
	var dateFrom, dateTo time.Time
	from := req.URL.Query().Get("from")
	if from != "" {
		var ts int64
		ts, err := strconv.ParseInt(req.URL.Query().Get("from"), 10, 64)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("failed to parse 'from' timestamp: %w", err)
		}
		//nolint:gosmopolitan // take TZ from user config eventually
		dateFrom = time.Unix(ts, 0).Local()
	} else {
		//nolint:gosmopolitan // take TZ from user config eventually
		dateFrom = utils.RoundToGranularity(time.Now().Local(), granularity, false)
	}
	dateTo = utils.RoundToGranularity(dateFrom, granularity, true)

	return dateFrom, dateTo, nil
}

//nolint:funlen,cyclop
func (r *WebAppRouter) transactionsHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "transactions")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	accountID := req.URL.Query().Get("accountID")
	data["AccountID"] = accountID

	accounts, err := r.db.GetAccounts(userID)
	if err != nil {
		r.logger.Error("Failed to get accounts", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Accounts"] = accounts

	currencies, err := r.db.GetCurrencies(userID)
	if err != nil {
		r.logger.Error("Failed to get currencies", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dateFrom, dateTo, err := getTimeRange(req, utils.GranularityMonth)
	if err != nil {
		r.logger.Error("Failed to get time range", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transactions, err := r.db.GetTransactions(userID, dateFrom, dateTo)
	if err != nil {
		r.logger.Error("Failed to get transactions", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []WebTransaction
	for _, t := range transactions {
		for _, m := range t.Movements {
			if accountID == "" || m.AccountId == accountID {
				result = append(result, transactionToWeb(t, accounts, currencies))
				break
			}
		}
	}

	data["From"] = dateFrom
	data["To"] = dateTo
	data["Current"] = dateFrom.Unix()
	data["Last"] = time.Date(
		dateFrom.Year(), dateFrom.Month(), 1, 0, 0, 0, 0, dateFrom.Location(),
	).AddDate(0, -1, 0).Unix()
	data["Next"] = dateTo.Unix()
	data["Transactions"] = &result

	if err := tmpl.ExecuteTemplate(w, "transactions.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint:funlen,cyclop
func (r *WebAppRouter) transactionsEditHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{}

	if err = req.ParseForm(); err != nil {
		r.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactionID := req.FormValue("id")

	// Validate session and obtain user ID using the centralized helper
	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}

	accounts, err := r.db.GetAccounts(userID)
	if err != nil {
		r.logger.Error("Failed to get accounts", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Accounts"] = accounts

	var currencies []goserver.Currency
	currencies, err = r.db.GetCurrencies(userID)
	if err != nil {
		r.logger.Error("Failed to get currencies", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var transaction goserver.Transaction
	if transactionID != "" {
		transaction, err = r.db.GetTransaction(userID, transactionID)
		if err != nil {
			r.logger.Error("Failed to get transaction", "error", err)
			r.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if req.Method == http.MethodPost {
		transaction.Description = req.FormValue("description")
		transaction.Tags = removeEmptyValues(strings.Split(req.FormValue("tags"), ","))
		for i := range transaction.Movements {
			transaction.Movements[i].AccountId = req.Form.Get(fmt.Sprintf("account_%d", i))
			transaction.Movements[i].CurrencyId = req.Form.Get(fmt.Sprintf("currency_%d", i))

			amountStr := req.Form.Get(fmt.Sprintf("amount_%d", i))
			transaction.Movements[i].Amount, err = strconv.ParseFloat(amountStr, 64)
			if err != nil {
				r.logger.Error("Failed to parse amount", "error", err)
				r.RespondError(w, "Invalid amount", http.StatusBadRequest)
				return
			}
		}

		if transactionID == "" {
			r.logger.Info("creating transaction", "transaction", transaction)
			if transaction, err = r.db.CreateTransaction(userID, &transaction); err != nil {
				r.logger.Error("Failed to create transaction", "error", err)
				r.RespondError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			r.logger.Info("updating transaction", "transaction", transaction)
			if transaction, err = r.db.UpdateTransaction(userID, transactionID, &transaction); err != nil {
				r.logger.Error("Failed to save transaction", "error", err)
				r.RespondError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	data["Transaction"] = transactionToWeb(transaction, accounts, currencies)

	if err := tmpl.ExecuteTemplate(w, "transaction_edit.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}
