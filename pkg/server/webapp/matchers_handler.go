package webapp

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

//nolint:dupl
func (r *WebAppRouter) matchersHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "matchers")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	matchers, err := r.db.GetMatchers(userID)
	if err != nil {
		r.logger.Error("Failed to get matchers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Matchers"] = &matchers

	if err := tmpl.ExecuteTemplate(w, "matchers.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint:dupl
func (r *WebAppRouter) matchersDeleteHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

//nolint:funlen,cyclop
func (r *WebAppRouter) matcherEditHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{}

	if err = req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matcherID := req.FormValue("id")
	transactionID := req.FormValue("transaction_id")

	userID, _, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		matcher, err = r.db.GetMatcher(userID, matcherID)
		if err != nil {
			r.logger.Error("Failed to get matcher", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		matcher = goserver.Matcher{
			Name:                       transaction.Description,
			OutputDescription:          transaction.Description,
			DescriptionRegExp:          transaction.Description,
			PartnerAccountNumberRegExp: transaction.PartnerAccount,
			PartnerNameRegExp:          transaction.PartnerName,
		}
	}

	if req.Method == http.MethodPost {
		matcher = goserver.Matcher{
			Name:                       req.FormValue("name"),
			OutputDescription:          req.FormValue("outputDescription"),
			DescriptionRegExp:          req.FormValue("descriptionRegExp"),
			OutputAccountId:            req.FormValue("account"),
			PartnerAccountNumberRegExp: req.FormValue("partnerAccountNumberRegExp"),
			OutputTags:                 removeEmptyValues(strings.Split(req.FormValue("outputTags"), ",")),
		}

		if matcher.OutputAccountId == "" {
			http.Error(w, "Account is required", http.StatusBadRequest)
			return
		}

		if matcherID == "" {
			r.logger.Info("creating matcher", "matcher", matcher)
			if matcher, err = r.db.CreateMatcher(userID, &matcher); err != nil {
				r.logger.Error("Failed to create matcher", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			r.logger.Info("updating matcher", "matcher", matcher)
			if matcher, err = r.db.UpdateMatcher(userID, matcherID, &matcher); err != nil {
				r.logger.Error("Failed to save matcher", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	data["Matcher"] = matcher

	if err := tmpl.ExecuteTemplate(w, "matcher_edit.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// http.Redirect(w, req, "/web/matchers", http.StatusFound)
}

func removeEmptyValues(arr []string) []string {
	var result []string
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

// matcherCheckHandler handles POST requests to /web/matchers/check
// This is a web wrapper around the API endpoint that handles authentication via session cookies
func (r *WebAppRouter) matcherCheckHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, code, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		http.Error(w, "Unauthorized", code)
		return
	}

	// Parse the request body
	var checkRequest goserver.CheckMatcherRequest
	if decodeErr := json.NewDecoder(req.Body).Decode(&checkRequest); decodeErr != nil {
		r.logger.Error("Failed to decode request body", "error", decodeErr)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a new context with the userID
	ctx := context.WithValue(req.Context(), common.UserIDKey, userID)

	// Call the API service directly
	matchersService := api.NewMatchersAPIServiceImpl(r.logger, r.db)
	result, err := matchersService.CheckMatcher(ctx, checkRequest)
	if err != nil {
		r.logger.Error("Failed to check matcher", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Code)
	if result.Body != nil {
		if err := json.NewEncoder(w).Encode(result.Body); err != nil {
			r.logger.Error("Failed to encode response", "error", err)
		}
	}
}

//nolint:dupl
func (r *WebAppRouter) matcherDeleteHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, _, err := r.GetUserIDFromSession(req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
