package webapp

import "net/http"

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
		r.logger.Info("Matchers", "matchers", matchers)

		data["Matchers"] = &matchers
	}

	if err := tmpl.ExecuteTemplate(w, "matchers.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
