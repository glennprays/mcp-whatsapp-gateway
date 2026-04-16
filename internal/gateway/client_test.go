package gateway

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
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

func TestNewWithHTTPClient(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		httpClient *http.Client
		wantErr    bool
	}{
		{
			name: "valid inputs",
			config: &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			},
			httpClient: &http.Client{},
			wantErr:    false,
		},
		{
			name:       "nil config",
			config:     nil,
			httpClient: &http.Client{},
			wantErr:    true,
		},
		{
			name: "nil http client",
			config: &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			},
			httpClient: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewWithHTTPClient(tt.config, tt.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithHTTPClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewWithHTTPClient() returned nil client for valid inputs")
			}
		})
	}
}

func TestClient_SendText(t *testing.T) {
	tests := []struct {
		name       string
		msisdn     string
		message    string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful send",
			msisdn:     "6281234567890@s.whatsapp.net",
			message:    "Hello, World!",
			response:   `{"success": true, "message_id": "msg123"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "gateway error",
			msisdn:     "6281234567890@s.whatsapp.net",
			message:    "Hello, World!",
			response:   `{"error": "Unauthorized"}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check authorization header
				auth := r.Header.Get("Authorization")
				if auth != "Bearer test-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Check method and path
				if r.Method != http.MethodPost || r.URL.Path != "/api/v1/message/text" {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			resp, err := client.SendText(context.Background(), tt.msisdn, tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("SendText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp != nil && !resp.Success {
				t.Error("SendText() returned unsuccessful response")
			}
		})
	}
}

func TestClient_GetLoginStatus(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantAuth   bool
	}{
		{
			name:       "authenticated",
			response:   `{"authenticated": true}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantAuth:   true,
		},
		{
			name:       "not authenticated",
			response:   `{"authenticated": false}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantAuth:   false,
		},
		{
			name:       "gateway error",
			response:   `{"error": "Internal server error"}`,
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/login/status" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			status, err := client.GetLoginStatus(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoginStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && status.Authenticated != tt.wantAuth {
				t.Errorf("GetLoginStatus() authenticated = %v, want %v", status.Authenticated, tt.wantAuth)
			}
		})
	}
}

func TestClient_Health(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "healthy",
			response:   `{"status": "ok", "timestamp": "2024-01-01T00:00:00Z"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "unhealthy",
			response:   `{"status": "error"}`,
			statusCode: http.StatusServiceUnavailable,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/health" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			health, err := client.Health(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && health.Status != "ok" {
				t.Errorf("Health() status = %v, want ok", health.Status)
			}
		})
	}
}

func TestClient_GetWebhook(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantURL    string
	}{
		{
			name:       "webhook registered",
			response:   `{"url": "https://example.com/webhook"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantURL:    "https://example.com/webhook",
		},
		{
			name:       "no webhook",
			response:   `{"url": ""}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/webhook" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			webhook, err := client.GetWebhook(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && webhook.URL != tt.wantURL {
				t.Errorf("GetWebhook() url = %v, want %v", webhook.URL, tt.wantURL)
			}
		})
	}
}

func TestClient_RegisterWebhook(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		hmacSecret string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful registration",
			url:        "https://example.com/webhook",
			hmacSecret: "secret123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid URL",
			url:        "invalid-url",
			hmacSecret: "",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/webhook" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.RegisterWebhook(context.Background(), tt.url, tt.hmacSecret)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteWebhook(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/webhook" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.DeleteWebhook(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_EditMessage(t *testing.T) {
	tests := []struct {
		name       string
		msisdn     string
		messageID  string
		newMessage string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful edit",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "msg123",
			newMessage: "Edited message",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "nonexistent",
			newMessage: "Edited message",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/message" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.EditMessage(context.Background(), tt.msisdn, tt.messageID, tt.newMessage)

			if (err != nil) != tt.wantErr {
				t.Errorf("EditMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteMessage(t *testing.T) {
	tests := []struct {
		name       string
		msisdn     string
		messageID  string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "msg123",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "nonexistent",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/message" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.DeleteMessage(context.Background(), tt.msisdn, tt.messageID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ReactToMessage(t *testing.T) {
	tests := []struct {
		name       string
		msisdn     string
		messageID  string
		emoji      string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful reaction",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "msg123",
			emoji:      "👍",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			msisdn:     "6281234567890@s.whatsapp.net",
			messageID:  "nonexistent",
			emoji:      "👍",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/message/react" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.ReactToMessage(context.Background(), tt.msisdn, tt.messageID, tt.emoji)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReactToMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CheckHealth(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantHealthy bool
	}{
		{
			name:        "healthy",
			response:    `{"status": "ok"}`,
			statusCode:  http.StatusOK,
			wantHealthy: true,
		},
		{
			name:        "unhealthy",
			response:    `{"status": "error"}`,
			statusCode:  http.StatusServiceUnavailable,
			wantHealthy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			err := client.CheckHealth(context.Background())
			isHealthy := err == nil

			if isHealthy != tt.wantHealthy {
				t.Errorf("CheckHealth() = %v, want %v", isHealthy, tt.wantHealthy)
			}
		})
	}
}

func TestClient_IsHealthy(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantHealthy bool
	}{
		{
			name:        "healthy",
			response:    `{"status": "ok"}`,
			statusCode:  http.StatusOK,
			wantHealthy: true,
		},
		{
			name:        "unhealthy",
			response:    `{"status": "error"}`,
			statusCode:  http.StatusServiceUnavailable,
			wantHealthy: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			cfg := &config.Config{
				WagaBaseURL:  server.URL + "/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			isHealthy := client.IsHealthy(context.Background())

			if isHealthy != tt.wantHealthy {
				t.Errorf("IsHealthy() = %v, want %v", isHealthy, tt.wantHealthy)
			}
		})
	}
}

func TestClient_SendImage(t *testing.T) {
	// Test the placeholder implementation
	tests := []struct {
		name      string
		msisdn    string
		caption   string
		isViewOnce bool
		wantErr   bool
	}{
		{
			name:      "successful send (placeholder)",
			msisdn:    "6281234567890@s.whatsapp.net",
			caption:   "Test image",
			isViewOnce: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				WagaBaseURL:  "http://localhost:3000/api/v1",
				WagaJWTToken: "test-token",
			}

			client, _ := New(cfg)
			imageData := []byte("fake image data")

			resp, err := client.SendImage(context.Background(), tt.msisdn, bytes.NewReader(imageData), tt.caption, tt.isViewOnce)

			if (err != nil) != tt.wantErr {
				t.Errorf("SendImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Error("SendImage() returned nil response")
			}
		})
	}
}

func TestGatewayClientInterface(t *testing.T) {
	// Ensure Client implements GatewayClient interface
	var _ GatewayClient = &Client{}
}

func TestClientTimeout(t *testing.T) {
	// Test that HTTP client timeout is respected
	cfg := &config.Config{
		WagaBaseURL:  "http://localhost:3000/api/v1",
		WagaJWTToken: "test-token",
	}

	client, _ := New(cfg)

	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("Client timeout = %v, want %v", client.httpClient.Timeout, DefaultTimeout)
	}
}
