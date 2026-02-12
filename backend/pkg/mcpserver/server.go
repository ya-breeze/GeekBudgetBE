package mcpserver

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

// MCPServer wraps the MCP server with GeekBudget-specific context
type MCPServer struct {
	logger  *slog.Logger
	storage database.Storage
	userID  string
}

// Run starts the MCP stdio server
func Run(ctx context.Context, logger *slog.Logger, storage database.Storage, userID string) error {
	s := &MCPServer{
		logger:  logger,
		storage: storage,
		userID:  userID,
	}

	// Create MCP server with implementation info
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "geekbudget",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		Instructions: instructions,
	})

	// Register all tool groups
	s.registerAccountTools(server)
	s.registerCurrencyTools(server)
	s.registerTransactionTools(server)
	s.registerMatcherTools(server)
	s.registerBudgetTools(server)
	s.registerReconciliationTools(server)
	s.registerNotificationTools(server)
	s.registerAnalysisTools(server)

	logger.Info("MCP server registered all tools")

	// Run stdio transport
	transport := &mcp.StdioTransport{}
	return server.Run(ctx, transport)
}
