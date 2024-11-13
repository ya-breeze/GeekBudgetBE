package webapp

import (
	"errors"
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func (r *WebAppRouter) loginHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := req.Form.Get("username")
	password := req.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	userID, err := r.db.GetUserID(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		r.logger.Warn("failed to get user ID", "username", username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := r.cookies.Get(req, "session-name")
	session.Values["userID"] = userID
	err = session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "GeekBudget API",
		"UserID": userID,
	}

	if err := tmpl.ExecuteTemplate(w, "home.tpl", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
