package webapp

import (
	"fmt"
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

func (r *WebAppRouter) unprocessedConvertHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactionID := req.Form.Get("transaction_id")

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	t, err := r.db.GetTransaction(userID, transactionID)
	if err != nil {
		r.logger.Error("Failed to get transaction", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range t.Movements {
		t.Movements[i].AccountId = req.Form.Get(fmt.Sprintf("account_%d", i))
	}

	s := api.NewUnprocessedTransactionsAPIServiceImpl(r.logger, r.db)
	_, err = s.Convert(req.Context(), userID, transactionID, &t)
	if err != nil {
		r.logger.Error("Failed to convert unprocessed transaction", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/web/unprocessed?id="+t.Id, http.StatusSeeOther)
}
