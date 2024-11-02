package webapp

import "net/http"

func (r *WebAppRouter) bankImportersHandler(w http.ResponseWriter, req *http.Request) {
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

		bankimporters, err := r.db.GetBankImporters(userID)
		if err != nil {
			r.logger.Error("Failed to get bank importers", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.logger.Info("Bank importers", "bankimporters", bankimporters)

		data["BankImporters"] = &bankimporters
	}

	if err := tmpl.ExecuteTemplate(w, "bank_importers.html", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
