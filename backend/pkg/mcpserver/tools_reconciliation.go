package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (s *MCPServer) registerReconciliationTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_reconciliation_status",
		Description: "Get reconciliation status for all asset accounts (compares app balance vs bank balance)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getReconciliationStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_reconciliation_history",
		Description: "Get reconciliation history for a specific account",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getReconciliationHistory)
}

type reconciliationStatusResponse struct {
	AccountID                       string     `json:"accountId"`
	AccountName                     string     `json:"accountName"`
	CurrencyID                      string     `json:"currencyId"`
	CurrencySymbol                  string     `json:"currencySymbol"`
	BankBalance                     string     `json:"bankBalance"`
	AppBalance                      string     `json:"appBalance"`
	Delta                           string     `json:"delta"`
	HasUnprocessedTransactions      bool       `json:"hasUnprocessedTransactions"`
	HasBankImporter                 bool       `json:"hasBankImporter"`
	BankBalanceAt                   *time.Time `json:"bankBalanceAt,omitempty"`
	HasTransactionsAfterBankBalance bool       `json:"hasTransactionsAfterBankBalance"`
	LastReconciledAt                *time.Time `json:"lastReconciledAt,omitempty"`
	LastReconciledBalance           string     `json:"lastReconciledBalance,omitempty"`
	IsManualReconciliationEnabled   bool       `json:"isManualReconciliationEnabled"`
}

func (s *MCPServer) getReconciliationStatus(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	accounts, err := s.storage.GetAccounts(s.userID)
	if err != nil {
		s.logger.Error("Failed to get accounts", "error", err)
		return errorResult(err)
	}

	bankImporters, err := s.storage.GetBankImporters(s.userID)
	if err != nil {
		s.logger.Error("Failed to get bank importers", "error", err)
		return errorResult(err)
	}

	currencies, err := s.storage.GetCurrencies(s.userID)
	if err != nil {
		s.logger.Error("Failed to get currencies", "error", err)
		return errorResult(err)
	}

	currencyMap := make(map[string]string)
	for _, c := range currencies {
		currencyMap[c.Id] = c.Name
	}

	accountsWithImporter := make(map[string]bool)
	for _, bi := range bankImporters {
		accountsWithImporter[bi.AccountId] = true
	}

	bulkData, err := s.storage.GetBulkReconciliationData(s.userID)
	if err != nil {
		s.logger.Error("Failed to get bulk reconciliation data", "error", err)
		return errorResult(err)
	}

	var statuses []reconciliationStatusResponse
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
			unprocessedCount := bulkData.UnprocessedCounts[acc.Id]

			bankBalance := b.ClosingBalance
			var bankBalanceAt *time.Time = b.LastUpdatedAt

			if !hasImporter && lastRec != nil {
				bankBalance = lastRec.ReconciledBalance
				bankBalanceAt = &lastRec.ReconciledAt
			}

			status := reconciliationStatusResponse{
				AccountID:                  acc.Id,
				AccountName:                acc.Name,
				CurrencyID:                 b.CurrencyId,
				CurrencySymbol:             currencyMap[b.CurrencyId],
				BankBalance:                bankBalance.String(),
				AppBalance:                 appBalance.String(),
				Delta:                      appBalance.Sub(bankBalance).String(),
				HasUnprocessedTransactions: unprocessedCount > 0,
				HasBankImporter:            hasImporter,
				BankBalanceAt:              bankBalanceAt,
			}

			if bankBalanceAt != nil {
				maxDate := bulkData.MaxTransactionDates[acc.Id][b.CurrencyId]
				status.HasTransactionsAfterBankBalance = maxDate.After(*bankBalanceAt)
			}

			if lastRec != nil {
				status.LastReconciledAt = &lastRec.ReconciledAt
				status.LastReconciledBalance = lastRec.ReconciledBalance.String()
				status.IsManualReconciliationEnabled = lastRec.IsManual
			}

			statuses = append(statuses, status)
		}
	}

	return jsonResult(statuses)
}

type getReconciliationHistoryArgs struct {
	AccountID string `json:"accountId" jsonschema:"Account ID (UUID)"`
}

func (s *MCPServer) getReconciliationHistory(ctx context.Context, req *mcp.CallToolRequest, args getReconciliationHistoryArgs) (*mcp.CallToolResult, any, error) {
	reconciliations, err := s.storage.GetReconciliationsForAccount(s.userID, args.AccountID)
	if err != nil {
		s.logger.Error("Failed to get reconciliation history", "error", err, "accountId", args.AccountID)
		return errorResult(err)
	}
	return jsonResult(reconciliations)
}
