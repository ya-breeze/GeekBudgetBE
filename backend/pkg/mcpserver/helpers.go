package mcpserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// toJSON converts any value to indented JSON string
func toJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling to JSON: %v", err)
	}
	return string(data)
}

// textResult wraps plain text in an MCP result
func textResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, nil, nil
}

// jsonResult serializes data to JSON and wraps it in an MCP result
func jsonResult(data any) (*mcp.CallToolResult, any, error) {
	return textResult(toJSON(data))
}

// errorResult creates an MCP error result
func errorResult(err error) (*mcp.CallToolResult, any, error) {
	return nil, nil, err
}

// parseOptionalDate parses a date string in YYYY-MM-DD format, returns zero time if empty
func parseOptionalDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", s)
}

// defaultDateFrom parses a date string or returns N days ago if empty
func defaultDateFrom(s string, daysBack int) (time.Time, error) {
	if s == "" {
		return time.Now().AddDate(0, 0, -daysBack), nil
	}
	return time.Parse("2006-01-02", s)
}
