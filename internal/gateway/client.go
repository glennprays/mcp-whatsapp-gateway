package gateway

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

const (
	// DefaultTimeout is the default timeout for gateway requests
	DefaultTimeout = 30 * time.Second
)

// GatewayClient defines the interface for gateway operations
type GatewayClient interface {
	// Message operations. Sends take a trailing variadic of SDK SendOptions
	// (WithChat / WithReply / WithMentions) so callers can address by canonical
	// chat and thread replies/mentions without changing the positional args.
	SendText(ctx context.Context, msisdn, message string, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendAudio(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendVideo(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendDocument(ctx context.Context, msisdn string, document io.Reader, fileName, caption string, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendLocation(ctx context.Context, msisdn string, latitude, longitude float64, name, address string, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendPoll(ctx context.Context, msisdn, question string, options []string, selectableCount int, opts ...waga.SendOption) (*SendMessageResponse, error)
	SendSticker(ctx context.Context, msisdn string, sticker io.Reader, opts ...waga.SendOption) (*SendMessageResponse, error)
	EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error
	DeleteMessage(ctx context.Context, msisdn, messageID string) error
	ReactToMessage(ctx context.Context, msisdn, messageID, emoji string, senderMsisdn ...string) error

	// Read operations
	GetIncomingMessages(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error)
	GetJobStatus(ctx context.Context, jobID string) (*JobStatusResponse, error)

	// Contact & group reads. These return the SDK response types directly; the
	// error chains preserve the SDK sentinels (waga.IsNotFound / waga.IsForbidden)
	// so tools can surface soft failures as results.
	ListContacts(ctx context.Context, limit, offset int) (*waga.ContactListResponse, error)
	GetContactInfo(ctx context.Context, chat string) (*waga.ContactInfoResponse, error)
	GetAvatar(ctx context.Context, chat string, preview bool) (*waga.AvatarResponse, error)
	ListGroups(ctx context.Context) (*waga.GroupListResponse, error)
	GetGroupInfo(ctx context.Context, chat string) (*waga.GroupInfoResponse, error)

	// Two-way conversation actions. These are outbound, conversation-affecting
	// calls governed server-side by the gateway's outbound pacer (per-account
	// pace + per-recipient cap); over-budget calls are paced or 429'd.
	MarkRead(ctx context.Context, chat string, messageIDs []string, sender string) error
	SendChatPresence(ctx context.Context, chat, state string) error

	// Connection operations
	GetLoginStatus(ctx context.Context) (*LoginStatus, error)
	CheckContact(ctx context.Context, msisdn string) (*ContactCheckResponse, error)
	Health(ctx context.Context) (*HealthResponse, error)

	// Webhook operations
	GetWebhook(ctx context.Context) (*WebhookResponse, error)
	RegisterWebhook(ctx context.Context, url, hmacSecret string) error
	DeleteWebhook(ctx context.Context) error
}

// Client wraps the WhatsApp Gateway SDK client
type Client struct {
	client *waga.Client
	config *config.Config
}

// New creates a new gateway client from the provided configuration
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Initialize the gateway SDK client with base URL and JWT token
	client := waga.NewClient(
		waga.WithBaseURL(cfg.WagaBaseURL),
		waga.WithToken(cfg.WagaJWTToken),
		waga.WithTimeout(DefaultTimeout),
	)

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// NewWithClient creates a new gateway client with a custom SDK client
func NewWithClient(cfg *config.Config, sdkClient *waga.Client) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if sdkClient == nil {
		return nil, fmt.Errorf("SDK client cannot be nil")
	}

	return &Client{
		client: sdkClient,
		config: cfg,
	}, nil
}

// GetClient returns the underlying WhatsApp Gateway SDK client
func (c *Client) GetClient() *waga.Client {
	return c.client
}

// SendText sends a text message to the specified recipient
func (c *Client) SendText(ctx context.Context, msisdn, message string, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendText(ctx, msisdn, message, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send text message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendImage sends an image message to the specified recipient
func (c *Client) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendImage(ctx, msisdn, image, caption, isViewOnce, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send image message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendAudio sends an audio message to the specified recipient
func (c *Client) SendAudio(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendAudio(ctx, msisdn, audio, isPTT, isViewOnce, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send audio message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendVideo sends a video message to the specified recipient
func (c *Client) SendVideo(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendVideo(ctx, msisdn, video, caption, isGif, isViewOnce, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send video message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendDocument sends a document message to the specified recipient
func (c *Client) SendDocument(ctx context.Context, msisdn string, document io.Reader, fileName, caption string, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendDocument(ctx, msisdn, document, fileName, caption, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send document message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendLocation sends a location message to the specified recipient
func (c *Client) SendLocation(ctx context.Context, msisdn string, latitude, longitude float64, name, address string, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendLocation(ctx, msisdn, latitude, longitude, name, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send location message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendPoll sends a poll message to the specified recipient
func (c *Client) SendPoll(ctx context.Context, msisdn, question string, options []string, selectableCount int, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendPoll(ctx, msisdn, question, options, selectableCount, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send poll message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// SendSticker sends a sticker message to the specified recipient
func (c *Client) SendSticker(ctx context.Context, msisdn string, sticker io.Reader, opts ...waga.SendOption) (*SendMessageResponse, error) {
	resp, err := c.client.SendSticker(ctx, msisdn, sticker, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to send sticker message: %w", err)
	}

	return &SendMessageResponse{
		Success:   resp.Success,
		MessageID: resp.MessageId,
		Status:    resp.Status,
		JobID:     resp.JobID,
	}, nil
}

// EditMessage edits a previously sent message
func (c *Client) EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error {
	err := c.client.EditMessage(ctx, msisdn, messageID, newMessage)
	if err != nil {
		return fmt.Errorf("failed to edit message: %w", err)
	}
	return nil
}

// DeleteMessage deletes a previously sent message
func (c *Client) DeleteMessage(ctx context.Context, msisdn, messageID string) error {
	err := c.client.DeleteMessage(ctx, msisdn, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// ReactToMessage reacts to a message with an emoji
func (c *Client) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string, senderMsisdn ...string) error {
	err := c.client.React(ctx, msisdn, messageID, emoji, senderMsisdn...)
	if err != nil {
		return fmt.Errorf("failed to react to message: %w", err)
	}
	return nil
}

// GetIncomingMessages fetches the most recent incoming WhatsApp messages
// buffered by the gateway for the authenticated session, newest first.
//
// The limit caps the number of messages returned. If limit <= 0 the gateway
// substitutes its default (10); values above 50 are clamped server-side.
func (c *Client) GetIncomingMessages(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
	resp, err := c.client.GetIncomingMessages(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get incoming messages: %w", err)
	}
	return resp, nil
}

// GetJobStatus polls the status of a queued message job by its ID.
func (c *Client) GetJobStatus(ctx context.Context, jobID string) (*JobStatusResponse, error) {
	resp, err := c.client.GetJobStatus(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}
	return &JobStatusResponse{
		JobID:       resp.JobID,
		Status:      resp.Status,
		MessageID:   resp.MessageID,
		Error:       resp.Error,
		CreatedAt:   resp.CreatedAt,
		CompletedAt: resp.CompletedAt,
	}, nil
}

// ListContacts returns a page of the account's locally-synced contacts. An empty
// address book is not an error (the gateway never 404s on empty).
func (c *Client) ListContacts(ctx context.Context, limit, offset int) (*waga.ContactListResponse, error) {
	resp, err := c.client.ListContacts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list contacts: %w", err)
	}
	return resp, nil
}

// GetContactInfo returns a server-side profile lookup for one user.
func (c *Client) GetContactInfo(ctx context.Context, chat string) (*waga.ContactInfoResponse, error) {
	resp, err := c.client.GetContactInfo(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact info: %w", err)
	}
	return resp, nil
}

// GetAvatar returns a chat's profile picture. preview requests the low-res
// thumbnail. The returned error preserves the SDK sentinels: a chat with no
// picture yields waga.ErrNotFound (404), a hidden picture waga.ErrForbidden (403).
func (c *Client) GetAvatar(ctx context.Context, chat string, preview bool) (*waga.AvatarResponse, error) {
	resp, err := c.client.GetAvatar(ctx, chat, preview)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}
	return resp, nil
}

// ListGroups returns the account's joined groups as lightweight summaries.
func (c *Client) ListGroups(ctx context.Context) (*waga.GroupListResponse, error) {
	resp, err := c.client.ListGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	return resp, nil
}

// GetGroupInfo returns a single group's full detail plus its participant roster.
func (c *Client) GetGroupInfo(ctx context.Context, chat string) (*waga.GroupInfoResponse, error) {
	resp, err := c.client.GetGroupInfo(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}
	return resp, nil
}

// MarkRead marks one or more messages in a chat as read (blue ticks). sender is
// the message author's JID/number and is required for group chats.
func (c *Client) MarkRead(ctx context.Context, chat string, messageIDs []string, sender string) error {
	if err := c.client.MarkRead(ctx, chat, messageIDs, sender); err != nil {
		return fmt.Errorf("failed to mark messages read: %w", err)
	}
	return nil
}

// SendChatPresence sets the typing indicator in a chat (composing/recording/paused).
func (c *Client) SendChatPresence(ctx context.Context, chat, state string) error {
	if err := c.client.SendChatPresence(ctx, chat, state); err != nil {
		return fmt.Errorf("failed to send chat presence: %w", err)
	}
	return nil
}

// GetLoginStatus checks if the WhatsApp session is authenticated
func (c *Client) GetLoginStatus(ctx context.Context) (*LoginStatus, error) {
	status, err := c.client.GetLoginStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get login status: %w", err)
	}

	return &LoginStatus{
		Authenticated: status.Authenticated,
	}, nil
}

// CheckContact validates if a recipient number is on WhatsApp.
func (c *Client) CheckContact(ctx context.Context, msisdn string) (*ContactCheckResponse, error) {
	contact, err := c.client.CheckContact(ctx, msisdn)
	if err != nil {
		return nil, fmt.Errorf("failed to check contact: %w", err)
	}

	return &ContactCheckResponse{
		Query:        contact.Query,
		JID:          contact.JID,
		IsOnWhatsApp: contact.IsOnWhatsApp,
		VerifiedName: contact.VerifiedName,
	}, nil
}

// Health checks if the gateway is reachable
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	health, err := c.client.Health(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check gateway health: %w", err)
	}

	return &HealthResponse{
		Status:    health.Status,
		Timestamp: health.Timestamp,
	}, nil
}

// GetWebhook retrieves the currently registered webhook
func (c *Client) GetWebhook(ctx context.Context) (*WebhookResponse, error) {
	webhook, err := c.client.GetWebhook(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}

	return &WebhookResponse{
		URL: webhook.URL,
	}, nil
}

// RegisterWebhook registers a webhook URL for incoming message notifications
func (c *Client) RegisterWebhook(ctx context.Context, url, hmacSecret string) error {
	err := c.client.RegisterWebhook(ctx, url, hmacSecret)
	if err != nil {
		return fmt.Errorf("failed to register webhook: %w", err)
	}
	return nil
}

// DeleteWebhook removes the currently registered webhook
func (c *Client) DeleteWebhook(ctx context.Context) error {
	err := c.client.UnregisterWebhook(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
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

// Response types - these wrap the SDK types for consistency

// SendMessageResponse represents the response from sending a message.
//
// In direct mode the gateway returns MessageID. In queue mode it returns
// Status ("queued") and JobID instead, with an empty MessageID — poll the
// job with GetJobStatus to obtain the final message ID/outcome.
type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Status    string `json:"status,omitempty"`
	JobID     string `json:"job_id,omitempty"`
}

// JobStatusResponse represents the status of a queued message job.
type JobStatusResponse struct {
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	MessageID   *string `json:"message_id,omitempty"`
	Error       *string `json:"error,omitempty"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

// LoginStatus represents the WhatsApp session status
type LoginStatus struct {
	Authenticated bool `json:"authenticated"`
}

// ContactCheckResponse represents a contact validation result.
type ContactCheckResponse struct {
	Query        string  `json:"query"`
	JID          string  `json:"jid"`
	IsOnWhatsApp bool    `json:"is_on_whatsapp"`
	VerifiedName *string `json:"verified_name,omitempty"`
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
