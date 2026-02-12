package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) registerCurrencyTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_currencies",
		Description: "List all currencies defined in the system",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listCurrencies)
}

func (s *MCPServer) listCurrencies(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	currencies, err := s.storage.GetCurrencies(s.userID)
	if err != nil {
		s.logger.Error("Failed to get currencies", "error", err)
		return errorResult(err)
	}
	return jsonResult(currencies)
}
