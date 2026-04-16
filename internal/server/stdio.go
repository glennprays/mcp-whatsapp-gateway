package server

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/tools"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServer wraps the MCP server instance
type MCPServer struct {
	server *mcp.Server
}

// NewStdioServer creates a new MCP server with stdio transport
func NewStdioServer(cfg *config.Config, gatewayClient gateway.GatewayClient) (*MCPServer, error) {
	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whatsapp-gateway",
		Version: "1.0.0",
	}, nil)

	// Create context with gateway client for tool handlers
	ctx := context.WithValue(context.Background(), "gateway", gatewayClient)

	// Register messaging tools using the SDK's typed handler system
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_text_message",
		Description: "Send a text message to a WhatsApp contact or group",
	}, createSendTextMessageHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_image_message",
		Description: "Send an image message to a WhatsApp contact or group",
	}, createSendImageMessageHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_message",
		Description: "Edit a previously sent message",
	}, createEditMessageHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_message",
		Description: "Delete a previously sent message",
	}, createDeleteMessageHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "react_to_message",
		Description: "React to a message with an emoji",
	}, createReactToMessageHandler(ctx))

	// Register connection tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_connection_status",
		Description: "Check if the WhatsApp session is active and authenticated",
	}, createCheckConnectionStatusHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_health",
		Description: "Check if the WhatsApp Gateway service is reachable",
	}, createCheckHealthHandler(ctx))

	// Register webhook tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_webhook",
		Description: "Get the current webhook configuration",
	}, createGetWebhookHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "register_webhook",
		Description: "Register a webhook URL to receive WhatsApp messages",
	}, createRegisterWebhookHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete the currently registered webhook",
	}, createDeleteWebhookHandler(ctx))

	return &MCPServer{server: server}, nil
}

// Handler creation functions - these bridge our existing tool handlers to MCP SDK format

func createSendTextMessageHandler(ctx context.Context) mcp.ToolHandlerFor[tools.SendMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendMessageInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.SendTextMessage(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendImageMessageHandler(ctx context.Context) mcp.ToolHandlerFor[tools.SendImageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendImageInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.SendImageMessage(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createEditMessageHandler(ctx context.Context) mcp.ToolHandlerFor[tools.EditMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.EditMessageInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.EditMessage(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteMessageHandler(ctx context.Context) mcp.ToolHandlerFor[tools.DeleteMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteMessageInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.DeleteMessage(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createReactToMessageHandler(ctx context.Context) mcp.ToolHandlerFor[tools.ReactToMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.ReactToMessageInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.ReactToMessage(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckConnectionStatusHandler(ctx context.Context) mcp.ToolHandlerFor[tools.CheckConnectionStatusInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckConnectionStatusInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.CheckConnectionStatus(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckHealthHandler(ctx context.Context) mcp.ToolHandlerFor[tools.CheckHealthInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckHealthInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.CheckHealth(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetWebhookHandler(ctx context.Context) mcp.ToolHandlerFor[tools.GetWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetWebhookInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.GetWebhook(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createRegisterWebhookHandler(ctx context.Context) mcp.ToolHandlerFor[tools.RegisterWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.RegisterWebhookInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.RegisterWebhook(ctx, req, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteWebhookHandler(ctx context.Context) mcp.ToolHandlerFor[tools.DeleteWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteWebhookInput) (*mcp.CallToolResult, any, error) {
		_, result, err := tools.DeleteWebhook(ctx, req, input)
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
