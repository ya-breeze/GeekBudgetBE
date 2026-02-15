package webapp

import (
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:funlen,cyclop,gocognit
func (r *WebAppRouter) unprocessedHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "unprocessed")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

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

	id := req.URL.Query().Get("id")
	if id != "" {
		r.logger.Info("Skipping unprocessed transactions to specified ID", "id", id)
	}
	s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
	unprocessed, cnt, err := s.PrepareUnprocessedTransactions(req.Context(), userID, true, id)
	if err != nil {
		r.logger.Error("Failed to get unprocessed", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(unprocessed) != 0 {
		u := unprocessed[0]
		web := WebUnprocessedTransaction{
			Transaction: transactionToWeb(u.Transaction, accounts, currencies),
		}

		allMatcherIDs := make([]string, 0, len(u.Matched))
		for _, m := range u.Matched {
			allMatcherIDs = append(allMatcherIDs, m.MatcherId)
		}
		for _, m := range u.Matched {
			others := make([]string, 0, len(allMatcherIDs)-1)
			for _, id := range allMatcherIDs {
				if id != m.MatcherId {
					others = append(others, id)
				}
			}

			// Fetch matcher to obtain confirmation history (reverted to DB fetch)
			confirmationsOK := 0
			confirmationsTotal := 0
			if matcher, err := r.db.GetMatcher(userID, m.MatcherId); err == nil {
				// Use the new fields on the Matcher object or calculate from history
				// faster to just use the fields if they are populated, or recalculate
				// The retrieved matcher is models.Matcher
				history := matcher.GetConfirmationHistory()
				if history != nil {
					for _, v := range history {
						if v {
							confirmationsOK++
						}
					}
					confirmationsTotal = len(history)
				}
			} else {
				// Log the error but continue; leave counts at 0
				r.logger.Warn("failed to load matcher for unprocessed", "matcherId", m.MatcherId, "error", err)
			}

			web.Matched = append(web.Matched, WebMatcherAndTransaction{
				MatcherID:       m.MatcherId,
				OtherMatcherIDs: others,
				Transaction: transactionToWeb(
					transactionNoIDToTransaction(m.Transaction, u.Transaction.Id),
					accounts, currencies),
				ConfirmationsOK:    confirmationsOK,
				ConfirmationsTotal: confirmationsTotal,
				ConfidenceClass: func() string {
					if confirmationsTotal == 0 {
						return "bg-secondary"
					}
					ratio := float64(confirmationsOK) / float64(confirmationsTotal)
					switch {
					case ratio >= 0.7:
						return "bg-success"
					case ratio >= 0.4:
						return "bg-warning text-dark"
					default:
						return "bg-danger"
					}
				}(),
			})
		}
		for _, d := range u.Duplicates {
			web.Duplicates = append(web.Duplicates, transactionToWeb(d, accounts, currencies))
		}

		data["Unprocessed"] = &web
		data["UnprocessedCount"] = cnt
	}

	if err := tmpl.ExecuteTemplate(w, "unprocessed.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *WebAppRouter) unprocessedDeleteHandler(w http.ResponseWriter, req *http.Request) {
	userID, code, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.RespondError(w, http.StatusText(code), code)
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		r.RespondError(w, "no id", http.StatusBadRequest)
		return
	}
	duplicateOf := req.URL.Query().Get("duplicateOf")
	if duplicateOf == "" {
		r.RespondError(w, "no duplicateOf", http.StatusBadRequest)
		return
	}

	_, err = r.db.MergeTransactions(userID, duplicateOf, id)
	if err != nil {
		r.logger.Error("Failed to merge duplicate", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO redirect to the same page, not to the start of the unprocessed
	http.Redirect(w, req, "/web/unprocessed", http.StatusSeeOther)
}
