package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
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

// GetWebhookDirect retrieves the current webhook configuration
func GetWebhookDirect(client gateway.GatewayClient, input GetWebhookInput) (GetWebhookResult, error) {
	ctx := context.Background()
	webhook, err := client.GetWebhook(ctx)
	if err != nil {
		return GetWebhookResult{}, fmt.Errorf("get_webhook: %w", err)
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

	return result, nil
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

// RegisterWebhookDirect registers a webhook URL to receive WhatsApp messages
func RegisterWebhookDirect(client gateway.GatewayClient, input RegisterWebhookInput) (RegisterWebhookResult, error) {
	if input.URL == "" {
		return RegisterWebhookResult{}, fmt.Errorf("webhook URL is required")
	}

	parsedURL, err := url.Parse(input.URL)
	if err != nil {
		return RegisterWebhookResult{}, fmt.Errorf("invalid webhook URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return RegisterWebhookResult{}, fmt.Errorf("webhook URL must use HTTP or HTTPS scheme")
	}

	ctx := context.Background()
	err = client.RegisterWebhook(ctx, input.URL, input.HMACSecret)
	if err != nil {
		return RegisterWebhookResult{}, fmt.Errorf("register_webhook: %w", err)
	}

	result := RegisterWebhookResult{
		Success:     true,
		URL:         input.URL,
		Status:      "registered",
		Description: "Webhook registered successfully. Incoming WhatsApp messages will be sent to this URL.",
	}

	return result, nil
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

// DeleteWebhookDirect removes the currently registered webhook
func DeleteWebhookDirect(client gateway.GatewayClient, input DeleteWebhookInput) (DeleteWebhookResult, error) {
	ctx := context.Background()
	err := client.DeleteWebhook(ctx)
	if err != nil {
		return DeleteWebhookResult{}, fmt.Errorf("delete_webhook: %w", err)
	}

	result := DeleteWebhookResult{
		Success:     true,
		Status:      "deleted",
		Description: "Webhook deleted successfully. Incoming WhatsApp messages will no longer be sent.",
	}

	return result, nil
}
