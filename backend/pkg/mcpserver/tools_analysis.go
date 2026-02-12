package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (s *MCPServer) registerAnalysisTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "financial_summary",
		Description: "Get a comprehensive financial overview including all accounts, balances, notifications, and reconciliation status. Use this first to understand the user's financial state.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.financialSummary)
}

type financialSummaryResponse struct {
	Accounts      []accountSummary      `json:"accounts"`
	Currencies    []currencySummary     `json:"currencies"`
	Balances      []balanceSummary      `json:"balances"`
	Notifications []notificationSummary `json:"notifications"`
}

type accountSummary struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Type                    string `json:"type"`
	ShowInReconciliation    bool   `json:"showInReconciliation"`
	ManualReconciliation    bool   `json:"manualReconciliation"`
	UnprocessedTransactions int    `json:"unprocessedTransactions"`
}

type currencySummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type balanceSummary struct {
	AccountID      string `json:"accountId"`
	AccountName    string `json:"accountName"`
	CurrencyID     string `json:"currencyId"`
	CurrencySymbol string `json:"currencySymbol"`
	AppBalance     string `json:"appBalance"`
	BankBalance    string `json:"bankBalance,omitempty"`
	Delta          string `json:"delta,omitempty"`
	NeedsAttention bool   `json:"needsAttention"`
}

type notificationSummary struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

func (s *MCPServer) financialSummary(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	accounts, err := s.storage.GetAccounts(s.userID)
	if err != nil {
		s.logger.Error("Failed to get accounts", "error", err)
		return errorResult(err)
	}

	currencies, err := s.storage.GetCurrencies(s.userID)
	if err != nil {
		s.logger.Error("Failed to get currencies", "error", err)
		return errorResult(err)
	}

	bulkData, err := s.storage.GetBulkReconciliationData(s.userID)
	if err != nil {
		s.logger.Error("Failed to get bulk reconciliation data", "error", err)
		return errorResult(err)
	}

	notifications, err := s.storage.GetNotifications(s.userID)
	if err != nil {
		s.logger.Error("Failed to get notifications", "error", err)
		return errorResult(err)
	}

	currencyMap := make(map[string]string)
	for _, c := range currencies {
		currencyMap[c.Id] = c.Name
	}

	accountMap := make(map[string]string)
	for _, a := range accounts {
		accountMap[a.Id] = a.Name
	}

	bankImporters, err := s.storage.GetBankImporters(s.userID)
	if err != nil {
		s.logger.Error("Failed to get bank importers", "error", err)
		return errorResult(err)
	}
	accountsWithImporter := make(map[string]bool)
	for _, bi := range bankImporters {
		accountsWithImporter[bi.AccountId] = true
	}

	resp := financialSummaryResponse{
		Accounts:      make([]accountSummary, 0, len(accounts)),
		Currencies:    make([]currencySummary, 0, len(currencies)),
		Balances:      make([]balanceSummary, 0),
		Notifications: make([]notificationSummary, 0, len(notifications)),
	}

	for _, acc := range accounts {
		resp.Accounts = append(resp.Accounts, accountSummary{
			ID:                      acc.Id,
			Name:                    acc.Name,
			Type:                    acc.Type,
			ShowInReconciliation:    acc.ShowInReconciliation,
			ManualReconciliation:    acc.ShowInReconciliation,
			UnprocessedTransactions: bulkData.UnprocessedCounts[acc.Id],
		})
	}

	for _, cur := range currencies {
		resp.Currencies = append(resp.Currencies, currencySummary{
			ID:     cur.Id,
			Name:   cur.Name,
			Symbol: cur.Name,
		})
	}

	for _, acc := range accounts {
		if acc.Type != "asset" {
			continue
		}

		for _, b := range acc.BankInfo.Balances {
			var lastRec *goserver.Reconciliation
			if accLatest, ok := bulkData.LatestReconciliations[acc.Id]; ok {
				lastRec = accLatest[b.CurrencyId]
			}

			hasImporter := accountsWithImporter[acc.Id]
			if !hasImporter && lastRec == nil && !acc.ShowInReconciliation {
				continue
			}

			appBalance := bulkData.Balances[acc.Id][b.CurrencyId]
			bankBalance := b.ClosingBalance

			if !hasImporter && lastRec != nil {
				bankBalance = lastRec.ReconciledBalance
			}

			delta := appBalance.Sub(bankBalance)
			needsAttention := delta.Abs().GreaterThan(decimal.NewFromFloat(0.01)) ||
				bulkData.UnprocessedCounts[acc.Id] > 0

			balance := balanceSummary{
				AccountID:      acc.Id,
				AccountName:    acc.Name,
				CurrencyID:     b.CurrencyId,
				CurrencySymbol: currencyMap[b.CurrencyId],
				AppBalance:     appBalance.String(),
				BankBalance:    bankBalance.String(),
				Delta:          delta.String(),
				NeedsAttention: needsAttention,
			}

			resp.Balances = append(resp.Balances, balance)
		}
	}

	for _, acc := range accounts {
		if acc.Type == "asset" {
			continue
		}

		if balances, ok := bulkData.Balances[acc.Id]; ok {
			for currencyID, balance := range balances {
				resp.Balances = append(resp.Balances, balanceSummary{
					AccountID:      acc.Id,
					AccountName:    acc.Name,
					CurrencyID:     currencyID,
					CurrencySymbol: currencyMap[currencyID],
					AppBalance:     balance.String(),
					NeedsAttention: false,
				})
			}
		}
	}

	for _, notif := range notifications {
		resp.Notifications = append(resp.Notifications, notificationSummary{
			ID:      notif.Id,
			Message: notif.Description,
			Level:   notif.Type,
		})
	}

	return jsonResult(resp)
}
