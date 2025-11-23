package tool

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/gomcpgo/api_wrapper/config"
)

// APIToolHandler implements the ToolHandler interface
type APIToolHandler struct {
	cfg *config.Config
}

// NewAPIToolHandler creates a new API tool handler
func NewAPIToolHandler(cfg *config.Config) *APIToolHandler {
	return &APIToolHandler{
		cfg: cfg,
	}
}

// HandleToolCall executes an API call
func (h *APIToolHandler) HandleToolCall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Find the tool configuration
	var toolCfg *config.ToolConfig
	for i := range h.cfg.Tools {
		if h.cfg.Tools[i].Name == req.Params.Name {
			toolCfg = &h.cfg.Tools[i]
			break
		}
	}

	if toolCfg == nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []interface{}{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Tool not found: %s", req.Params.Name),
				},
			},
		}, nil
	}

	// Execute the API call
	result, err := h.executeAPICall(ctx, toolCfg, req.Params.Arguments)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []interface{}{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("API call failed: %v", err),
				},
			},
		}, nil
	}

	// Return the result
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}
