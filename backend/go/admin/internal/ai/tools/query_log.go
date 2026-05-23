package tools

import (
	"context"
	"fmt"

	e_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetLogMcpTool connects to a log-query MCP server at the given URL and returns its tools.
func GetLogMcpTool(mcpURL string) ([]tool.BaseTool, error) {
	if mcpURL == "" {
		return nil, fmt.Errorf("MCP_URL is empty")
	}

	ctx := context.Background()
	cli, err := client.NewSSEMCPClient(mcpURL)
	if err != nil {
		return nil, err
	}
	if err := cli.Start(ctx); err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "harness-template",
		Version: "1.0.0",
	}
	if _, err := cli.Initialize(ctx, initRequest); err != nil {
		return nil, err
	}

	return e_mcp.GetTools(ctx, &e_mcp.Config{Cli: cli})
}
