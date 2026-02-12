package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) registerNotificationTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_notifications",
		Description: "List all active notifications (warnings, errors, information messages)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listNotifications)
}

func (s *MCPServer) listNotifications(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	notifications, err := s.storage.GetNotifications(s.userID)
	if err != nil {
		s.logger.Error("Failed to get notifications", "error", err)
		return errorResult(err)
	}
	return jsonResult(notifications)
}
