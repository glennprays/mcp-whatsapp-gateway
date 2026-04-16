package server

import (
	"context"
	"io"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
)

// mockGatewayClient is a simple mock for testing
type mockGatewayClient struct{}

func (m *mockGatewayClient) SendText(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_msg_id"}, nil
}

func (m *mockGatewayClient) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error) {
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_img_msg_id"}, nil
}

func (m *mockGatewayClient) EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error {
	return nil
}

func (m *mockGatewayClient) DeleteMessage(ctx context.Context, msisdn, messageID string) error {
	return nil
}

func (m *mockGatewayClient) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string) error {
	return nil
}

func (m *mockGatewayClient) GetLoginStatus(ctx context.Context) (*gateway.LoginStatus, error) {
	return &gateway.LoginStatus{Authenticated: true}, nil
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

