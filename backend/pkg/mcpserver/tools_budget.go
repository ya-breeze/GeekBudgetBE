package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) registerBudgetTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_budget_items",
		Description: "List all budget items (expected spending per account/category)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listBudgetItems)
}

func (s *MCPServer) listBudgetItems(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	budgetItems, err := s.storage.GetBudgetItems(s.userID)
	if err != nil {
		s.logger.Error("Failed to get budget items", "error", err)
		return errorResult(err)
	}
	return jsonResult(budgetItems)
}
