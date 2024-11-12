package webapp

import (
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

//nolint:funlen,cyclop
func (r *WebAppRouter) unprocessedHandler(w http.ResponseWriter, req *http.Request) {
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

		currencies, err := r.db.GetCurrencies(userID)
		if err != nil {
			r.logger.Error("Failed to get currencies", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := req.URL.Query().Get("id")
		if id != "" {
			r.logger.Info("Skipping unprocessed transactions", "id", id)
		}
		s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
		unprocessed, err := s.PrepareUnprocessedTransactions(req.Context(), userID, true, id)
		if err != nil {
			r.logger.Error("Failed to get unprocessed", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(unprocessed) != 0 {
			u := unprocessed[0]
			web := WebUnprocessedTransaction{
				Transaction: transactionToWeb(u.Transaction, accounts, currencies),
			}
			for _, m := range u.Matched {
				web.Matched = append(web.Matched, WebMatcherAndTransaction{
					MatcherID: m.MatcherId,
					Transaction: transactionToWeb(
						transactionNoIDToTransaction(m.Transaction, u.Transaction.Id),
						accounts, currencies),
				})
			}
			for _, d := range u.Duplicates {
				web.Duplicates = append(web.Duplicates, transactionToWeb(d, accounts, currencies))
			}

			data["Unprocessed"] = &web
		}
	}

	if err := tmpl.ExecuteTemplate(w, "unprocessed.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
