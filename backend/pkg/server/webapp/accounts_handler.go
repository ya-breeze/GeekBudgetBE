package webapp

import (
	"net/http"
	"slices"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:dupl
func (r *WebAppRouter) accountsHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "accounts")

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
	data["Accounts"] = &accounts

	if err := tmpl.ExecuteTemplate(w, "accounts.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint:funlen,cyclop
func (r *WebAppRouter) accountsEditHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{}

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
	data["Accounts"] = &accounts

	if req.Method == http.MethodGet {
		id := req.URL.Query().Get("id")
		idx := slices.IndexFunc(accounts, func(a goserver.Account) bool {
			return a.Id == id
		})
		if idx != -1 {
			data["Id"] = accounts[idx].Id
			data["Name"] = accounts[idx].Name
			data["Type"] = accounts[idx].Type
			data["Description"] = accounts[idx].Description
		}
	} else {
		var acc goserver.Account
		id := req.FormValue("id")
		name := req.FormValue("name")

		if id == "" {
			if name == "" {
				r.logger.Info("name is empty")
				data["Error"] = "'Name' can't be empty"
			} else {
				if slices.ContainsFunc(accounts, func(a goserver.Account) bool {
					return a.Name == name
				}) {
					r.logger.Info("name already exists")
					data["Error"] = "Account with this name already exists"
				} else {
					r.logger.Info("creating account", "name", name)
					acc, err = r.db.CreateAccount(userID, &goserver.AccountNoId{
						Name:        name,
						Type:        req.FormValue("type"),
						Description: req.FormValue("description"),
					})
					if err != nil {
						r.logger.Error("Failed to create account", "error", err)
						r.RespondError(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		} else {
			r.logger.Info("updating account", "name", name)
			acc, err = r.db.UpdateAccount(userID, req.FormValue("id"), &goserver.AccountNoId{
				Name:        name,
				Type:        req.FormValue("type"),
				Description: req.FormValue("description"),
			})
			if err != nil {
				r.logger.Error("Failed to update account", "error", err)
				r.RespondError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data["Id"] = acc.Id
		data["Name"] = acc.Name
		data["Type"] = acc.Type
		data["Description"] = acc.Description
	}

	if err := tmpl.ExecuteTemplate(w, "accounts_edit.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		r.RespondError(w, err.Error(), http.StatusInternalServerError)
	}
}
