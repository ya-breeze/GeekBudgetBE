package webapp

import (
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (r *WebAppRouter) matchersHandler(w http.ResponseWriter, req *http.Request) {
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

		// accounts, err := r.db.GetAccounts(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get accounts", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// currencies, err := r.db.GetCurrencies(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get currencies", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		if req.Method == http.MethodDelete {
			id := req.URL.Query().Get("id")
			if id != "" {
				if err := r.db.DeleteMatcher(userID, id); err != nil {
					r.logger.Error("Failed to delete matcher", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		matchers, err := r.db.GetMatchers(userID)
		if err != nil {
			r.logger.Error("Failed to get matchers", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["Matchers"] = &matchers
	}

	if err := tmpl.ExecuteTemplate(w, "matchers.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint:funlen,cyclop
func (r *WebAppRouter) matcherEditHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Title": "GeekBudget API",
	}

	if err = req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matcherID := req.FormValue("id")
	transactionID := req.FormValue("transaction_id")

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
		var t goserver.Transaction
		t, err = r.db.GetTransaction(userID, transactionID)
		if err != nil {
			r.logger.Error("Failed to get transaction", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var currencies []goserver.Currency
		currencies, err = r.db.GetCurrencies(userID)
		if err != nil {
			r.logger.Error("Failed to get currencies", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transaction = transactionToWeb(t, accounts, currencies)
	}
	data["Transaction"] = transaction

	var matcher goserver.Matcher
	if matcherID != "" {
		var m goserver.Matcher
		m, err = r.db.GetMatcher(userID, matcherID)
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

	if req.Method == http.MethodPost {
		m := goserver.MatcherNoId{
			Name:              req.FormValue("name"),
			OutputDescription: req.FormValue("outputDescription"),
			DescriptionRegExp: req.FormValue("descriptionRegExp"),
			OutputAccountId:   req.FormValue("account"),
		}

		if m.OutputAccountId == "" {
			http.Error(w, "Account is required", http.StatusBadRequest)
			return
		}

		if matcherID == "" {
			r.logger.Info("creating matcher", "name", m.Name)
			if matcher, err = r.db.CreateMatcher(userID, &m); err != nil {
				r.logger.Error("Failed to create matcher", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			r.logger.Info("updating matcher", "name", m.Name)
			if matcher, err = r.db.UpdateMatcher(userID, matcherID, &m); err != nil {
				r.logger.Error("Failed to save matcher", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	data["Matcher"] = matcher

	if err := tmpl.ExecuteTemplate(w, "matcher_edit.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// http.Redirect(w, req, "/web/matchers", http.StatusFound)
}

func (r *WebAppRouter) matcherDeleteHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id := req.FormValue("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := r.db.DeleteMatcher(userID, id); err != nil {
		r.logger.Error("Failed to delete matcher", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/web/matchers", http.StatusFound)
}
