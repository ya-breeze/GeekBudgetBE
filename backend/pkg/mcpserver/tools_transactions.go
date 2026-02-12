package mcpserver

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (s *MCPServer) registerTransactionTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_transactions",
		Description: "List transactions with optional date range and suspicious flag filters. Defaults to last 30 days if no dates provided.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listTransactions)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_transaction",
		Description: "Get detailed information about a specific transaction by ID",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getTransaction)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_transactions",
		Description: "Search transactions by text query (searches description, partner name, place, and tags)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.searchTransactions)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_duplicate_transactions",
		Description: "Get all transactions that are marked as duplicates of a given transaction",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getDuplicateTransactions)
}

type listTransactionsArgs struct {
	DateFrom       string `json:"dateFrom,omitempty" jsonschema:"Start date YYYY-MM-DD, defaults to 30 days ago"`
	DateTo         string `json:"dateTo,omitempty" jsonschema:"End date YYYY-MM-DD, defaults to today"`
	OnlySuspicious bool   `json:"onlySuspicious,omitempty" jsonschema:"If true, only return suspicious transactions"`
}

func (s *MCPServer) listTransactions(ctx context.Context, req *mcp.CallToolRequest, args listTransactionsArgs) (*mcp.CallToolResult, any, error) {
	dateFrom, err := defaultDateFrom(args.DateFrom, 30)
	if err != nil {
		return errorResult(err)
	}

	dateTo, err := parseOptionalDate(args.DateTo)
	if err != nil {
		return errorResult(err)
	}

	transactions, err := s.storage.GetTransactions(s.userID, dateFrom, dateTo, args.OnlySuspicious)
	if err != nil {
		s.logger.Error("Failed to get transactions", "error", err)
		return errorResult(err)
	}

	return jsonResult(transactions)
}

type getTransactionArgs struct {
	ID string `json:"id" jsonschema:"Transaction ID (UUID)"`
}

func (s *MCPServer) getTransaction(ctx context.Context, req *mcp.CallToolRequest, args getTransactionArgs) (*mcp.CallToolResult, any, error) {
	transaction, err := s.storage.GetTransaction(s.userID, args.ID)
	if err != nil {
		s.logger.Error("Failed to get transaction", "error", err, "id", args.ID)
		return errorResult(err)
	}
	return jsonResult(transaction)
}

type searchTransactionsArgs struct {
	Query    string `json:"query" jsonschema:"Search query text"`
	DateFrom string `json:"dateFrom,omitempty" jsonschema:"Start date YYYY-MM-DD, defaults to 30 days ago"`
	DateTo   string `json:"dateTo,omitempty" jsonschema:"End date YYYY-MM-DD, defaults to today"`
}

func (s *MCPServer) searchTransactions(ctx context.Context, req *mcp.CallToolRequest, args searchTransactionsArgs) (*mcp.CallToolResult, any, error) {
	dateFrom, err := defaultDateFrom(args.DateFrom, 30)
	if err != nil {
		return errorResult(err)
	}

	dateTo, err := parseOptionalDate(args.DateTo)
	if err != nil {
		return errorResult(err)
	}

	transactions, err := s.storage.GetTransactions(s.userID, dateFrom, dateTo, false)
	if err != nil {
		s.logger.Error("Failed to get transactions for search", "error", err)
		return errorResult(err)
	}

	query := strings.ToLower(args.Query)
	var filtered []goserver.Transaction
	for _, txn := range transactions {
		if matchesQuery(txn, query) {
			filtered = append(filtered, txn)
		}
	}

	return jsonResult(filtered)
}

func matchesQuery(txn goserver.Transaction, query string) bool {
	if strings.Contains(strings.ToLower(txn.Description), query) {
		return true
	}
	if txn.PartnerName != "" && strings.Contains(strings.ToLower(txn.PartnerName), query) {
		return true
	}
	if txn.Place != "" && strings.Contains(strings.ToLower(txn.Place), query) {
		return true
	}
	for _, tag := range txn.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

type getDuplicateTransactionsArgs struct {
	TransactionID string `json:"transactionId" jsonschema:"Transaction ID (UUID)"`
}

func (s *MCPServer) getDuplicateTransactions(ctx context.Context, req *mcp.CallToolRequest, args getDuplicateTransactionsArgs) (*mcp.CallToolResult, any, error) {
	duplicateIDs, err := s.storage.GetDuplicateTransactionIDs(s.userID, args.TransactionID)
	if err != nil {
		s.logger.Error("Failed to get duplicate transaction IDs", "error", err, "transactionId", args.TransactionID)
		return errorResult(err)
	}

	var duplicates []goserver.Transaction
	for _, id := range duplicateIDs {
		txn, err := s.storage.GetTransaction(s.userID, id)
		if err != nil {
			s.logger.Warn("Failed to get duplicate transaction", "error", err, "id", id)
			continue
		}
		duplicates = append(duplicates, txn)
	}

	return jsonResult(duplicates)
}
