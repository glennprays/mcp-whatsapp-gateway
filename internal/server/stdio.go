package server

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/tools"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Package-level gateway client for tool handlers
var globalGatewayClient gateway.GatewayClient

// MCPServer wraps the MCP server instance
type MCPServer struct {
	server *mcp.Server
}

// NewStdioServer creates a new MCP server with stdio transport
func NewStdioServer(cfg *config.Config, gatewayClient gateway.GatewayClient) (*MCPServer, error) {
	// Verify gateway client is not nil
	if gatewayClient == nil {
		return nil, fmt.Errorf("gatewayClient is nil")
	}

	// Store gateway client globally for tool handlers
	globalGatewayClient = gatewayClient

	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whatsapp-gateway",
		Version: "1.0.0",
	}, nil)

	// Register messaging tools using the SDK's typed handler system
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_text_message",
		Description: "Send a text message to a WhatsApp contact or group",
	}, createSendTextMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_image_message",
		Description: "Send an image message to a WhatsApp contact or group",
	}, createSendImageMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_message",
		Description: "Edit a previously sent message",
	}, createEditMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_message",
		Description: "Delete a previously sent message",
	}, createDeleteMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "react_to_message",
		Description: "React to a message with an emoji",
	}, createReactToMessageHandler(gatewayClient))

	// Register connection tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_connection_status",
		Description: "Check if the WhatsApp session is active and authenticated",
	}, createCheckConnectionStatusHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_health",
		Description: "Check if the WhatsApp Gateway service is reachable",
	}, createCheckHealthHandler(gatewayClient))

	// Register webhook tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_webhook",
		Description: "Get the current webhook configuration",
	}, createGetWebhookHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "register_webhook",
		Description: "Register a webhook URL to receive WhatsApp messages",
	}, createRegisterWebhookHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete the currently registered webhook",
	}, createDeleteWebhookHandler(gatewayClient))

	return &MCPServer{server: server}, nil
}

// Handler creation functions - these bridge our existing tool handlers to MCP SDK format
// The gateway client is captured directly in the closure

func createSendTextMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendMessageInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}

		// Call the tool function directly with the gateway client
		result, err := tools.SendTextMessageDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendImageMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendImageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendImageInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.SendImageMessageDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createEditMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.EditMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.EditMessageInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.EditMessageDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.DeleteMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteMessageInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.DeleteMessageDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createReactToMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.ReactToMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.ReactToMessageInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.ReactToMessageDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckConnectionStatusHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.CheckConnectionStatusInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckConnectionStatusInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.CheckConnectionStatusDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckHealthHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.CheckHealthInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckHealthInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.CheckHealthDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetWebhookInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.GetWebhookDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createRegisterWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.RegisterWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.RegisterWebhookInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.RegisterWebhookDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.DeleteWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteWebhookInput) (*mcp.CallToolResult, any, error) {
		// Use global gateway client instead of parameter
		client := globalGatewayClient
		if client == nil {
			return nil, nil, fmt.Errorf("gateway client not available (global is nil)")
		}
		result, err := tools.DeleteWebhookDirect(client, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

// RunStdio starts the server with stdio transport
func (s *MCPServer) RunStdio(ctx context.Context) error {
	// Create stdio transport
	transport := &mcp.StdioTransport{}

	// Run the server
	if err := s.server.Run(ctx, transport); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

// Close shuts down the server
func (s *MCPServer) Close() error {
	// The MCP SDK handles cleanup automatically
	return nil
}
