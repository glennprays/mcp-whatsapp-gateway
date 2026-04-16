package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetWebhookInput represents the input for getting webhook configuration
type GetWebhookInput struct {
	// No input parameters required
}

// GetWebhookResult represents the result of getting webhook configuration
type GetWebhookResult struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// GetWebhook retrieves the current webhook configuration
func GetWebhook(ctx context.Context, req *mcp.CallToolRequest, input GetWebhookInput) (
	*mcp.CallToolResult,
	GetWebhookResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, GetWebhookResult{}, fmt.Errorf("gateway client not available")
	}

	// Get webhook via gateway
	webhook, err := client.GetWebhook(ctx)
	if err != nil {
		return nil, GetWebhookResult{}, fmt.Errorf("get_webhook: %w", err)
	}

	result := GetWebhookResult{
		URL: webhook.URL,
	}

	if webhook.URL != "" {
		result.Status = "registered"
		result.Description = "Webhook is registered and active"
	} else {
		result.Status = "not_registered"
		result.Description = "No webhook is currently registered"
	}

	return nil, result, nil
}

// RegisterWebhookInput represents the input for registering a webhook
type RegisterWebhookInput struct {
	URL        string `json:"url"`
	HMACSecret string `json:"hmac_secret"`
}

// RegisterWebhookResult represents the result of registering a webhook
type RegisterWebhookResult struct {
	Success     bool   `json:"success"`
	URL         string `json:"url"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// RegisterWebhook registers a webhook URL to receive WhatsApp messages
func RegisterWebhook(ctx context.Context, req *mcp.CallToolRequest, input RegisterWebhookInput) (
	*mcp.CallToolResult,
	RegisterWebhookResult,
	error,
) {
	// Validate input
	if input.URL == "" {
		return nil, RegisterWebhookResult{}, fmt.Errorf("webhook URL is required")
	}

	// Validate URL format
	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return nil, RegisterWebhookResult{}, fmt.Errorf("invalid webhook URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, RegisterWebhookResult{}, fmt.Errorf("webhook URL must use HTTP or HTTPS scheme")
	}

	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, RegisterWebhookResult{}, fmt.Errorf("gateway client not available")
	}

	// Register webhook via gateway
	err = client.RegisterWebhook(ctx, input.URL, input.HMACSecret)
	if err != nil {
		return nil, RegisterWebhookResult{}, fmt.Errorf("register_webhook: %w", err)
	}

	result := RegisterWebhookResult{
		Success: true,
		URL:     input.URL,
		Status:  "registered",
		Description: "Webhook registered successfully. Incoming WhatsApp messages will be sent to this URL.",
	}

	return nil, result, nil
}

// DeleteWebhookInput represents the input for deleting a webhook
type DeleteWebhookInput struct {
	// No input parameters required
}

// DeleteWebhookResult represents the result of deleting a webhook
type DeleteWebhookResult struct {
	Success     bool   `json:"success"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// DeleteWebhook removes the currently registered webhook
func DeleteWebhook(ctx context.Context, req *mcp.CallToolRequest, input DeleteWebhookInput) (
	*mcp.CallToolResult,
	DeleteWebhookResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, DeleteWebhookResult{}, fmt.Errorf("gateway client not available")
	}

	// Delete webhook via gateway
	err := client.DeleteWebhook(ctx)
	if err != nil {
		return nil, DeleteWebhookResult{}, fmt.Errorf("delete_webhook: %w", err)
	}

	result := DeleteWebhookResult{
		Success: true,
		Status:  "deleted",
		Description: "Webhook deleted successfully. Incoming WhatsApp messages will no longer be sent.",
	}

	return nil, result, nil
}
