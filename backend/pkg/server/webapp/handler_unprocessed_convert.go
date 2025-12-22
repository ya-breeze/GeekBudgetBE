package webapp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

func (r *WebAppRouter) unprocessedConvertHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactionID := req.Form.Get("transaction_id")

	userID, code, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		r.RespondError(w, err.Error(), code)
		return
	}

	t, err := r.db.GetTransaction(userID, transactionID)
	if err != nil {
		r.logger.Error("Failed to get transaction", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range t.Movements {
		t.Movements[i].AccountId = req.Form.Get(fmt.Sprintf("account_%d", i))
	}

	// Update matcher ID if provided
	matcherID := req.Form.Get("matcher_id")
	if matcherID != "" {
		t.MatcherId = matcherID
		t.IsAuto = false
	}

	s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
	_, err = s.Convert(req.Context(), userID, transactionID, &t)
	if err != nil {
		r.logger.Error("Failed to convert unprocessed transaction", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.handleMatcherConfirmations(req, userID)

	http.Redirect(w, req, "/web/unprocessed?id="+t.Id, http.StatusSeeOther)
}

// handleMatcherConfirmations updates matcher confirmation history based on form fields.
// Extracted to reduce cyclomatic complexity of the main handler.
func (r *WebAppRouter) handleMatcherConfirmations(req *http.Request, userID string) {
	matcherID := req.Form.Get("matcher_id")
	otherMatchers := req.Form.Get("other_matchers")

	if matcherID != "" {
		if err := r.db.AddMatcherConfirmation(userID, matcherID, true); err != nil {
			r.logger.Warn("Failed to add confirmation to matcher", "matcher_id", matcherID, "error", err)
		}
	}

	if otherMatchers != "" {
		parts := strings.Split(otherMatchers, ",")
		for _, id := range parts {
			id = strings.TrimSpace(id)
			if id == "" || id == matcherID {
				continue
			}

			if err := r.db.AddMatcherConfirmation(userID, id, false); err != nil {
				r.logger.Warn("Failed to add confirmation to other matcher", "matcher_id", id, "error", err)
			}
		}
	}
}
