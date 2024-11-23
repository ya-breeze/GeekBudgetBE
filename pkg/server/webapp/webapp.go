package webapp

import (
	"html/template"
	"log/slog"
	"math"
	"path/filepath"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

type WebAppRouter struct {
	commit  string
	logger  *slog.Logger
	cfg     *config.Config
	db      database.Storage
	cookies *sessions.CookieStore
}

func NewWebAppRouter(
	commit string, logger *slog.Logger, cfg *config.Config, db database.Storage,
) *WebAppRouter {
	return &WebAppRouter{
		commit:  commit,
		logger:  logger,
		cfg:     cfg,
		db:      db,
		cookies: sessions.NewCookieStore([]byte("SESSION_KEY")),
	}
}

//nolint:funlen // This is a webapp router, it's supposed to have many routes.
func (r *WebAppRouter) Routes() goserver.Routes {
	return goserver.Routes{
		"RootPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: r.homeHandler,
		},
		"Login": goserver.Route{
			Method:      "POST",
			Pattern:     "/",
			HandlerFunc: r.loginHandler,
		},
		"AboutPath": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/about",
			HandlerFunc: r.aboutHandler,
		},
		"BankImporters": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/bank-importers",
			HandlerFunc: r.bankImportersHandler,
		},
		"Matchers": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/matchers",
			HandlerFunc: r.matchersHandler,
		},
		"MatchersDelete": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/matchers/delete",
			HandlerFunc: r.matchersDeleteHandler,
		},
		"MatcherEditGET": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/matchers/edit",
			HandlerFunc: r.matcherEditHandler,
		},
		"MatcherEditPOST": goserver.Route{
			Method:      "POST",
			Pattern:     "/web/matchers/edit",
			HandlerFunc: r.matcherEditHandler,
		},
		"MatcherDelete": goserver.Route{
			Method:      "DELETE",
			Pattern:     "/web/matchers",
			HandlerFunc: r.matcherDeleteHandler,
		},
		"Unprocessed": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/unprocessed",
			HandlerFunc: r.unprocessedHandler,
		},
		"UnprocessedConvert": goserver.Route{
			Method:      "POST",
			Pattern:     "/web/unprocessed/convert",
			HandlerFunc: r.unprocessedConvertHandler,
		},
		"UnprocessedDelete": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/unprocessed/delete",
			HandlerFunc: r.unprocessedDeleteHandler,
		},
		"Accounts": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/accounts",
			HandlerFunc: r.accountsHandler,
		},
		"AccountEditGet": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/accounts/edit",
			HandlerFunc: r.accountsEditHandler,
		},
		"AccountEditPost": goserver.Route{
			Method:      "POST",
			Pattern:     "/web/accounts/edit",
			HandlerFunc: r.accountsEditHandler,
		},

		"Transactions": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/transactions",
			HandlerFunc: r.transactionsHandler,
		},
		"TransactionEditGet": goserver.Route{
			Method:      "GET",
			Pattern:     "/web/transactions/edit",
			HandlerFunc: r.transactionsEditHandler,
		},
		"TransactionEditPost": goserver.Route{
			Method:      "POST",
			Pattern:     "/web/transactions/edit",
			HandlerFunc: r.transactionsEditHandler,
		},
	}
}

func (r *WebAppRouter) loadTemplates() (*template.Template, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"formatTime": utils.FormatTime,
		"decrease": func(i int) int {
			return i - 1
		},
		"money": func(num float64) float64 {
			return math.Round(num*100) / 100
		},
		"timestamp": func(t time.Time) int64 {
			return t.Unix()
		},
	}).ParseGlob(filepath.Join("webapp", "templates", "*.tpl"))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func transactionToWeb(
	t goserver.Transaction, accounts []goserver.Account, currencies []goserver.Currency,
) WebTransaction {
	res := WebTransaction{
		ID:             t.Id,
		Date:           t.Date,
		Description:    t.Description,
		Place:          t.Place,
		Tags:           t.Tags,
		PartnerName:    t.PartnerName,
		PartnerAccount: t.PartnerAccount,
		Movements:      make([]WebMovement, 0, len(t.Movements)),
	}

	for _, m := range t.Movements {
		res.Movements = append(res.Movements, WebMovement{
			Amount:       m.Amount,
			AccountID:    m.AccountId,
			AccountName:  utils.GetAccount(m.AccountId, accounts).Name,
			CurrencyID:   m.CurrencyId,
			CurrencyName: utils.GetCurrency(m.CurrencyId, currencies).Name,
		})
	}

	return res
}

func transactionNoIDToTransaction(t goserver.TransactionNoId, id string) goserver.Transaction {
	res := goserver.Transaction{
		Id:             id,
		Date:           t.Date,
		Description:    t.Description,
		Place:          t.Place,
		Tags:           t.Tags,
		PartnerName:    t.PartnerName,
		PartnerAccount: t.PartnerAccount,
		Movements:      make([]goserver.Movement, 0, len(t.Movements)),
	}

	for _, m := range t.Movements {
		res.Movements = append(res.Movements, goserver.Movement{
			Amount:     m.Amount,
			AccountId:  m.AccountId,
			CurrencyId: m.CurrencyId,
		})
	}

	return res
}
