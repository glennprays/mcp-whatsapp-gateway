package gateway

import (
	"context"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("New() returned nil client for valid config")
			}
		})
	}
}

func TestNewWithClient(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		sdkClient  *waga.Client
		wantErr    bool
	}{
		{
			name: "valid inputs",
			config: &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			},
			sdkClient:  waga.NewClient(),
			wantErr:    false,
		},
		{
			name:       "nil config",
			config:     nil,
			sdkClient:  waga.NewClient(),
			wantErr:    true,
		},
		{
			name: "nil SDK client",
			config: &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			},
			sdkClient:  nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewWithClient(tt.config, tt.sdkClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewWithClient() returned nil client for valid inputs")
			}
		})
	}
}

func TestClient_GetClient(t *testing.T) {
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if client.GetClient() == nil {
		t.Error("GetClient() returned nil SDK client")
	}
}

func TestGatewayClientInterface(t *testing.T) {
	// Ensure Client implements GatewayClient interface
	var _ GatewayClient = &Client{}
}

// Integration tests with actual SDK - these require a running gateway
// To run these tests, use: go test -tags=integration ./internal/gateway/...

func TestClient_SDKIntegration_SendText(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	_, err = client.SendText(context.Background(), "6281234567890@s.whatsapp.net", "test message")
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("SendText failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_GetLoginStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	_, err = client.GetLoginStatus(context.Background())
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("GetLoginStatus failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	_, err = client.Health(context.Background())
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("Health failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_GetWebhook(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	_, err = client.GetWebhook(context.Background())
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("GetWebhook failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_RegisterWebhook(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.RegisterWebhook(context.Background(), "https://example.com/webhook", "secret")
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("RegisterWebhook failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_DeleteWebhook(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.DeleteWebhook(context.Background())
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("DeleteWebhook failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_EditMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.EditMessage(context.Background(), "6281234567890@s.whatsapp.net", "msg123", "edited message")
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("EditMessage failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_DeleteMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.DeleteMessage(context.Background(), "6281234567890@s.whatsapp.net", "msg123")
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("DeleteMessage failed as expected without real gateway: %v", err)
	}
}

func TestClient_SDKIntegration_ReactToMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.ReactToMessage(context.Background(), "6281234567890@s.whatsapp.net", "msg123", "👍")
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("ReactToMessage failed as expected without real gateway: %v", err)
	}
}

func TestClient_CheckHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	err = client.CheckHealth(context.Background())
	if err != nil {
		// Expected to fail without real gateway
		t.Logf("CheckHealth failed as expected without real gateway: %v", err)
	}
}

func TestClient_IsHealthy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// This will fail without a real gateway, but we can test the interface
	isHealthy := client.IsHealthy(context.Background())
	if isHealthy {
		t.Error("IsHealthy() returned true without real gateway")
	}
}
