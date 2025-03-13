package webapp

import (
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:funlen,cyclop
func (r *WebAppRouter) unprocessedHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "unprocessed")

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

		id := req.URL.Query().Get("id")
		if id != "" {
			r.logger.Info("Skipping unprocessed transactions to specified ID", "id", id)
		}
		s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
		unprocessed, cnt, err := s.PrepareUnprocessedTransactions(req.Context(), userID, true, id)
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
			data["UnprocessedCount"] = cnt
		}
	}

	if err := tmpl.ExecuteTemplate(w, "unprocessed.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *WebAppRouter) unprocessedDeleteHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "no id", http.StatusBadRequest)
		return
	}
	duplicateOf := req.URL.Query().Get("duplicateOf")
	if duplicateOf == "" {
		http.Error(w, "no duplicateOf", http.StatusBadRequest)
		return
	}

	s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
	err := s.Delete(req.Context(), userID, id, duplicateOf)
	if err != nil {
		r.logger.Error("Failed to delete unprocessed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO redirect to the same page, not to the start of the unprocessed
	http.Redirect(w, req, "/web/unprocessed", http.StatusSeeOther)
}
