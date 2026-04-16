package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPServer wraps the HTTP MCP server instance
type HTTPServer struct {
	server *mcp.Server
	port   string
	handler http.Handler
}

// NewHTTPServer creates a new MCP server with HTTP+SSE transport
func NewHTTPServer(cfg *config.Config, gatewayClient gateway.GatewayClient) (*HTTPServer, error) {
	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whatsapp-gateway",
		Version: "1.0.0",
	}, nil)

	// Create context with gateway client for tool handlers
	ctx := context.WithValue(context.Background(), "gateway", gatewayClient)

	// Register all tools (same as stdio)
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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_connection_status",
		Description: "Check if the WhatsApp session is active and authenticated",
	}, createCheckConnectionStatusHandler(ctx))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_health",
		Description: "Check if the WhatsApp Gateway service is reachable",
	}, createCheckHealthHandler(ctx))

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

	// Create SSE handler
	handler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		// Apply Basic auth middleware in production
		if cfg.IsProduction() && cfg.IsHTTPTransport() {
			if !checkBasicAuth(r, cfg.BasicAuthUser, cfg.BasicAuthPassword) {
				return nil // Return nil to trigger 400 Bad Request
			}
		}
		return server
	}, nil)

	return &HTTPServer{
		server: server,
		port:   cfg.GetPort(),
		handler: handler,
	}, nil
}

// RunHTTP starts the server with HTTP+SSE transport
func (s *HTTPServer) RunHTTP(ctx context.Context) error {
	// Configure HTTP server
	addr := fmt.Sprintf(":%s", s.port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}

	// Start HTTP server
	go func() {
		fmt.Printf("HTTP+SSE server listening on %s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx := context.Background()
	httpServer.Shutdown(shutdownCtx)

	return nil
}

// checkBasicAuth verifies Basic authentication credentials
func checkBasicAuth(r *http.Request, username, password string) bool {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return user == username && pass == password
}

// Close shuts down the server
func (s *HTTPServer) Close() error {
	// The MCP SDK handles cleanup automatically
	return nil
}
