package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *MCPServer) registerMatcherTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_matchers",
		Description: "List all matchers (regex-based rules for auto-categorizing transactions)",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.listMatchers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_matcher",
		Description: "Get detailed information about a specific matcher by ID",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}, s.getMatcher)
}

func (s *MCPServer) listMatchers(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	matchers, err := s.storage.GetMatchers(s.userID)
	if err != nil {
		s.logger.Error("Failed to get matchers", "error", err)
		return errorResult(err)
	}
	return jsonResult(matchers)
}

type getMatcherArgs struct {
	ID string `json:"id" jsonschema:"Matcher ID (UUID)"`
}

func (s *MCPServer) getMatcher(ctx context.Context, req *mcp.CallToolRequest, args getMatcherArgs) (*mcp.CallToolResult, any, error) {
	matcher, err := s.storage.GetMatcher(s.userID, args.ID)
	if err != nil {
		s.logger.Error("Failed to get matcher", "error", err, "id", args.ID)
		return errorResult(err)
	}
	return jsonResult(matcher)
}
