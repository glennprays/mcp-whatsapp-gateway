package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
)

const (
	// DefaultTimeout is the default timeout for gateway requests
	DefaultTimeout = 30 * time.Second
)

// GatewayClient defines the interface for gateway operations
type GatewayClient interface {
	// Message operations
	SendText(ctx context.Context, msisdn, message string) (*SendMessageResponse, error)
	SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*SendMessageResponse, error)
	EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error
	DeleteMessage(ctx context.Context, msisdn, messageID string) error
	ReactToMessage(ctx context.Context, msisdn, messageID, emoji string) error

	// Connection operations
	GetLoginStatus(ctx context.Context) (*LoginStatus, error)
	Health(ctx context.Context) (*HealthResponse, error)

	// Webhook operations
	GetWebhook(ctx context.Context) (*WebhookResponse, error)
	RegisterWebhook(ctx context.Context, url, hmacSecret string) error
	DeleteWebhook(ctx context.Context) error
}

// Client implements GatewayClient interface
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// New creates a new gateway client from the provided configuration
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}, nil
}

// NewWithHTTPClient creates a new gateway client with a custom HTTP client
func NewWithHTTPClient(cfg *config.Config, httpClient *http.Client) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if httpClient == nil {
		return nil, fmt.Errorf("http client cannot be nil")
	}

	return &Client{
		config:     cfg,
		httpClient: httpClient,
	}, nil
}

// makeRequest makes an HTTP request to the gateway with proper authentication
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader, response interface{}) error {
	url := c.config.WagaBaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+c.config.WagaJWTToken)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(respBody))
	}

	if response != nil {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// SendMessage sends a text message to the specified recipient
func (c *Client) SendText(ctx context.Context, msisdn, message string) (*SendMessageResponse, error) {
	reqBody := map[string]interface{}{
		"msisdn":  msisdn,
		"message": message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var response SendMessageResponse
	if err := c.makeRequest(ctx, http.MethodPost, "/message/text", bytes.NewReader(jsonBody), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendImage sends an image message to the specified recipient
func (c *Client) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*SendMessageResponse, error) {
	// For image uploads, we need to use multipart/form-data
	// This is a simplified placeholder implementation
	// In production, you would use multipart.Writer to properly handle image uploads
	_ = msisdn  // Will be used in multipart form data
	_ = image   // Will be used in multipart form data
	_ = caption // Will be used in multipart form data
	_ = isViewOnce // Will be used in multipart form data

	// TODO: Implement proper multipart/form-data upload
	return &SendMessageResponse{
		Success:   true,
		MessageID: "placeholder_message_id",
	}, nil
}

// EditMessage edits a previously sent message
func (c *Client) EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error {
	reqBody := map[string]interface{}{
		"msisdn":     msisdn,
		"message_id": messageID,
		"new_message": newMessage,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return c.makeRequest(ctx, http.MethodPut, "/message", bytes.NewReader(jsonBody), nil)
}

// DeleteMessage deletes a previously sent message
func (c *Client) DeleteMessage(ctx context.Context, msisdn, messageID string) error {
	reqBody := map[string]interface{}{
		"msisdn":     msisdn,
		"message_id": messageID,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return c.makeRequest(ctx, http.MethodDelete, "/message", bytes.NewReader(jsonBody), nil)
}

// ReactToMessage reacts to a message with an emoji
func (c *Client) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string) error {
	reqBody := map[string]interface{}{
		"msisdn":     msisdn,
		"message_id": messageID,
		"emoji":      emoji,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return c.makeRequest(ctx, http.MethodPost, "/message/react", bytes.NewReader(jsonBody), nil)
}

// GetLoginStatus checks if the WhatsApp session is authenticated
func (c *Client) GetLoginStatus(ctx context.Context) (*LoginStatus, error) {
	var status LoginStatus
	if err := c.makeRequest(ctx, http.MethodGet, "/login/status", nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// Health checks if the gateway is reachable
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	url := c.config.WagaBaseURL + "/health"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Health endpoint doesn't require authentication
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var health HealthResponse
	if err := json.Unmarshal(respBody, &health); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &health, nil
}

// GetWebhook retrieves the currently registered webhook
func (c *Client) GetWebhook(ctx context.Context) (*WebhookResponse, error) {
	var webhook WebhookResponse
	if err := c.makeRequest(ctx, http.MethodGet, "/webhook", nil, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// RegisterWebhook registers a webhook URL for incoming message notifications
func (c *Client) RegisterWebhook(ctx context.Context, url, hmacSecret string) error {
	reqBody := map[string]interface{}{
		"url": url,
	}

	if hmacSecret != "" {
		reqBody["hmac_secret"] = hmacSecret
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	return c.makeRequest(ctx, http.MethodPost, "/webhook", bytes.NewReader(jsonBody), nil)
}

// DeleteWebhook removes the currently registered webhook
func (c *Client) DeleteWebhook(ctx context.Context) error {
	return c.makeRequest(ctx, http.MethodDelete, "/webhook", nil, nil)
}

// CheckHealth is a convenience method that returns nil if healthy, error otherwise
func (c *Client) CheckHealth(ctx context.Context) error {
	_, err := c.Health(ctx)
	return err
}

// IsHealthy returns true if the gateway is reachable
func (c *Client) IsHealthy(ctx context.Context) bool {
	return c.CheckHealth(ctx) == nil
}

// Response types

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
}

// LoginStatus represents the WhatsApp session status
type LoginStatus struct {
	Authenticated bool `json:"authenticated"`
}

// HealthResponse represents the gateway health status
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// WebhookResponse represents the webhook configuration
type WebhookResponse struct {
	URL string `json:"url"`
}
