package server

import (
	"context"
	"io"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// mockGatewayClient is a simple mock for testing
type mockGatewayClient struct{}

func (m *mockGatewayClient) SendText(ctx context.Context, msisdn, message string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_msg_id"}, nil
}

func (m *mockGatewayClient) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_img_msg_id"}, nil
}

func (m *mockGatewayClient) SendAudio(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_audio_msg_id"}, nil
}

func (m *mockGatewayClient) SendVideo(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_video_msg_id"}, nil
}

func (m *mockGatewayClient) SendDocument(ctx context.Context, msisdn string, document io.Reader, fileName, caption string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_doc_msg_id"}, nil
}

func (m *mockGatewayClient) SendLocation(ctx context.Context, msisdn string, latitude, longitude float64, name, address string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_loc_msg_id"}, nil
}

func (m *mockGatewayClient) SendPoll(ctx context.Context, msisdn, question string, options []string, selectableCount int, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_poll_msg_id"}, nil
}

func (m *mockGatewayClient) SendSticker(ctx context.Context, msisdn string, sticker io.Reader, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_sticker_msg_id"}, nil
}

func (m *mockGatewayClient) EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error {
	return nil
}

func (m *mockGatewayClient) DeleteMessage(ctx context.Context, msisdn, messageID string) error {
	return nil
}

func (m *mockGatewayClient) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string, senderMsisdn ...string) error {
	return nil
}

func (m *mockGatewayClient) GetLoginStatus(ctx context.Context) (*gateway.LoginStatus, error) {
	return &gateway.LoginStatus{Authenticated: true}, nil
}

func (m *mockGatewayClient) CheckContact(ctx context.Context, msisdn string) (*gateway.ContactCheckResponse, error) {
	return &gateway.ContactCheckResponse{Query: msisdn, JID: msisdn, IsOnWhatsApp: true}, nil
}

func (m *mockGatewayClient) Health(ctx context.Context) (*gateway.HealthResponse, error) {
	return &gateway.HealthResponse{Status: "ok", Timestamp: "2024-01-01T00:00:00Z"}, nil
}

func (m *mockGatewayClient) GetWebhook(ctx context.Context) (*gateway.WebhookResponse, error) {
	return &gateway.WebhookResponse{URL: "https://example.com/webhook"}, nil
}

func (m *mockGatewayClient) RegisterWebhook(ctx context.Context, url, hmacSecret string) error {
	return nil
}

func (m *mockGatewayClient) DeleteWebhook(ctx context.Context) error {
	return nil
}

func (m *mockGatewayClient) GetIncomingMessages(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
	return &waga.IncomingMessagesResponse{Success: true, Count: 0, Messages: []waga.IncomingMessage{}}, nil
}

func (m *mockGatewayClient) GetJobStatus(ctx context.Context, jobID string) (*gateway.JobStatusResponse, error) {
	return &gateway.JobStatusResponse{JobID: jobID, Status: "completed"}, nil
}

func (m *mockGatewayClient) ListContacts(ctx context.Context, limit, offset int) (*waga.ContactListResponse, error) {
	return &waga.ContactListResponse{Contacts: []waga.ContactListItem{}, Count: 0, Total: 0}, nil
}

func (m *mockGatewayClient) GetContactInfo(ctx context.Context, chat string) (*waga.ContactInfoResponse, error) {
	return &waga.ContactInfoResponse{JID: chat}, nil
}

func (m *mockGatewayClient) GetAvatar(ctx context.Context, chat string, preview bool) (*waga.AvatarResponse, error) {
	return &waga.AvatarResponse{JID: chat, URL: "https://cdn.example/avatar.jpg", ID: "abc", Type: "image"}, nil
}

func (m *mockGatewayClient) ListGroups(ctx context.Context) (*waga.GroupListResponse, error) {
	return &waga.GroupListResponse{Groups: []waga.GroupListItem{}, Count: 0}, nil
}

func (m *mockGatewayClient) GetGroupInfo(ctx context.Context, chat string) (*waga.GroupInfoResponse, error) {
	return &waga.GroupInfoResponse{JID: chat}, nil
}

func TestNewStdioServer_Success(t *testing.T) {
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
		AppEnv:       config.Dev,
		Transport:    "stdio",
	}

	mockClient := &mockGatewayClient{}

	server, err := NewStdioServer(cfg, mockClient)
	if err != nil {
		t.Fatalf("NewStdioServer() failed: %v", err)
	}

	if server == nil {
		t.Fatal("Expected server to be non-nil")
	}

	if server.server == nil {
		t.Error("Expected server.server to be non-nil")
	}
}

func TestMCPServer_Close(t *testing.T) {
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
		AppEnv:       config.Dev,
		Transport:    "stdio",
	}

	mockClient := &mockGatewayClient{}

	server, err := NewStdioServer(cfg, mockClient)
	if err != nil {
		t.Fatalf("NewStdioServer() failed: %v", err)
	}

	// Close should not error
	err = server.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}

func TestNewHTTPServer_Success(t *testing.T) {
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
		AppEnv:       config.Dev,
		Transport:    "http",
		Port:         "8080",
	}

	mockClient := &mockGatewayClient{}

	server, err := NewHTTPServer(cfg, mockClient)
	if err != nil {
		t.Fatalf("NewHTTPServer() failed: %v", err)
	}

	if server == nil {
		t.Fatal("Expected server to be non-nil")
	}

	if server.server == nil {
		t.Error("Expected server.server to be non-nil")
	}

	if server.port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", server.port)
	}
}

func TestHTTPServer_Close(t *testing.T) {
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
		AppEnv:       config.Dev,
		Transport:    "http",
		Port:         "8080",
	}

	mockClient := &mockGatewayClient{}

	server, err := NewHTTPServer(cfg, mockClient)
	if err != nil {
		t.Fatalf("NewHTTPServer() failed: %v", err)
	}

	// Close should not error
	err = server.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}
