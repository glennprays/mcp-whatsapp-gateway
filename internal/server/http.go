package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPServer wraps the HTTP MCP server instance
type HTTPServer struct {
	server  *mcp.Server
	port    string
	handler http.Handler
}

// NewHTTPServer creates a new MCP server with HTTP+SSE transport
func NewHTTPServer(cfg *config.Config, gatewayClient gateway.GatewayClient) (*HTTPServer, error) {
	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whatsapp-gateway",
		Version: "1.0.0",
	}, nil)

	// Register all tools (same as stdio)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_text_message",
		Description: "Send a text message to a WhatsApp contact or group",
	}, createSendTextMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_image_message",
		Description: "Send an image message to a WhatsApp contact or group",
	}, createSendImageMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_audio_message",
		Description: "Send an audio message to a WhatsApp contact or group",
	}, createSendAudioMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_video_message",
		Description: "Send a video message to a WhatsApp contact or group",
	}, createSendVideoMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_document_message",
		Description: "Send a document message to a WhatsApp contact or group",
	}, createSendDocumentMessageHandler(gatewayClient))

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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_location_message",
		Description: "Send a location with GPS coordinates to a WhatsApp number or group",
	}, createSendLocationMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_poll_message",
		Description: "Send a poll with question and options to a WhatsApp number or group",
	}, createSendPollMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_sticker_message",
		Description: "Send a sticker from a URL to a WhatsApp number or group",
	}, createSendStickerMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_job_status",
		Description: "Poll the status of a queued message job by its job_id (returned by send_* tools when the gateway runs in queue mode). Returns status (queued/processing/completed/failed), the message_id once completed, and any error.",
	}, createGetJobStatusHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_latest_incoming_messages",
		Description: "Fetch the most recent incoming WhatsApp messages. Useful for reading OTPs, verification codes, or recent conversation context. Returns newest first (text, image, video, audio, document, sticker, location, poll). Optional limit (default 10, max 50).",
	}, createGetLatestIncomingMessagesHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_connection_status",
		Description: "Check if the WhatsApp session is active and authenticated",
	}, createCheckConnectionStatusHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_contact",
		Description: "Check whether a phone number is registered on WhatsApp",
	}, createCheckContactHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_health",
		Description: "Check if the WhatsApp Gateway service is reachable",
	}, createCheckHealthHandler(gatewayClient))

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
		server:  server,
		port:    cfg.GetPort(),
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

	// Start HTTP server; surface startup failures (e.g. port in use)
	// instead of only printing them inside the goroutine
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("HTTP+SSE server listening on %s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for context cancellation or server failure
	select {
	case err := <-errCh:
		return fmt.Errorf("HTTP server error: %w", err)
	case <-ctx.Done():
	}

	// Graceful shutdown, bounded so a hung connection cannot block exit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server shutdown: %w", err)
	}

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
