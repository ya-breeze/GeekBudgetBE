package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) registerAccountTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_accounts",
		Description: "List all accounts (asset, expense, income)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listAccounts)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_account",
		Description: "Get detailed information about a specific account by ID",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getAccount)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_account_balance",
		Description: "Get current balance for a specific account and currency",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getAccountBalance)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_account_history",
		Description: "Get all transactions for a specific account",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getAccountHistory)
}

func (s *MCPServer) listAccounts(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	accounts, err := s.storage.GetAccounts(s.userID)
	if err != nil {
		s.logger.Error("Failed to get accounts", "error", err)
		return errorResult(err)
	}
	return jsonResult(accounts)
}

type getAccountArgs struct {
	ID string `json:"id" jsonschema:"Account ID (UUID)"`
}

func (s *MCPServer) getAccount(ctx context.Context, req *mcp.CallToolRequest, args getAccountArgs) (*mcp.CallToolResult, any, error) {
	account, err := s.storage.GetAccount(s.userID, args.ID)
	if err != nil {
		s.logger.Error("Failed to get account", "error", err, "id", args.ID)
		return errorResult(err)
	}
	return jsonResult(account)
}

type getAccountBalanceArgs struct {
	AccountID  string `json:"accountId" jsonschema:"Account ID (UUID)"`
	CurrencyID string `json:"currencyId" jsonschema:"Currency ID (UUID)"`
}

func (s *MCPServer) getAccountBalance(ctx context.Context, req *mcp.CallToolRequest, args getAccountBalanceArgs) (*mcp.CallToolResult, any, error) {
	balance, err := s.storage.GetAccountBalance(s.userID, args.AccountID, args.CurrencyID)
	if err != nil {
		s.logger.Error("Failed to get account balance", "error", err, "accountId", args.AccountID, "currencyId", args.CurrencyID)
		return errorResult(err)
	}
	return textResult(fmt.Sprintf("Balance: %s", balance.String()))
}

type getAccountHistoryArgs struct {
	AccountID string `json:"accountId" jsonschema:"Account ID (UUID)"`
}

func (s *MCPServer) getAccountHistory(ctx context.Context, req *mcp.CallToolRequest, args getAccountHistoryArgs) (*mcp.CallToolResult, any, error) {
	transactions, err := s.storage.GetAccountHistory(s.userID, args.AccountID)
	if err != nil {
		s.logger.Error("Failed to get account history", "error", err, "accountId", args.AccountID)
		return errorResult(err)
	}
	return jsonResult(transactions)
}
