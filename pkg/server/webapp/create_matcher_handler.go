package webapp

import (
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (r *WebAppRouter) createMatcherHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactionID := req.Form.Get("transaction_id")
	matcherID := req.Form.Get("matcher_id")

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	accounts, err := r.db.GetAccounts(userID)
	if err != nil {
		r.logger.Error("Failed to get accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Accounts"] = accounts

	transaction := WebTransaction{}
	if transactionID != "" {
		t, err := r.db.GetTransaction(userID, transactionID)
		if err != nil {
			r.logger.Error("Failed to get transaction", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currencies, err := r.db.GetCurrencies(userID)
		if err != nil {
			r.logger.Error("Failed to get currencies", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transaction = transactionToWeb(t, accounts, currencies)
	}
	data["Transaction"] = transaction

	matcher := goserver.Matcher{}
	if matcherID != "" {
		m, err := r.db.GetMatcher(userID, matcherID)
		if err != nil {
			r.logger.Error("Failed to get matcher", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matcher = m
	} else {
		matcher = goserver.Matcher{
			Name:              transaction.Description,
			OutputDescription: transaction.Description,
			DescriptionRegExp: transaction.Description,
		}
	}
	data["Matcher"] = matcher

	if req.Method == "POST" {
		m := goserver.MatcherNoId{
			Name:              req.Form.Get("name"),
			OutputDescription: req.Form.Get("outputDescription"),
			DescriptionRegExp: req.Form.Get("descriptionRegExp"),
			OutputAccountId:   req.Form.Get("account"),
		}

		if m.OutputAccountId == "" {
			http.Error(w, "Account is required", http.StatusBadRequest)
			return
		}

		if matcher, err = r.db.CreateMatcher(userID, &m); err != nil {
			r.logger.Error("Failed to save matcher", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tmpl.ExecuteTemplate(w, "matcher_edit.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// http.Redirect(w, req, "/web/matchers", http.StatusFound)
}
