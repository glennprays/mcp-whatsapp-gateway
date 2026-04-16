package tools

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
)

// MockGatewayClient is a mock implementation of GatewayClient for testing
type MockGatewayClient struct {
	SendTextFunc       func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error)
	SendImageFunc      func(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error)
	EditMessageFunc    func(ctx context.Context, msisdn, messageID, newMessage string) error
	DeleteMessageFunc  func(ctx context.Context, msisdn, messageID string) error
	ReactToMessageFunc func(ctx context.Context, msisdn, messageID, emoji string) error
	GetLoginStatusFunc func(ctx context.Context) (*gateway.LoginStatus, error)
	HealthFunc         func(ctx context.Context) (*gateway.HealthResponse, error)
	GetWebhookFunc     func(ctx context.Context) (*gateway.WebhookResponse, error)
	RegisterWebhookFunc func(ctx context.Context, url, hmacSecret string) error
	DeleteWebhookFunc  func(ctx context.Context) error
}

func (m *MockGatewayClient) SendText(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
	if m.SendTextFunc != nil {
		return m.SendTextFunc(ctx, msisdn, message)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_msg_id"}, nil
}

func (m *MockGatewayClient) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error) {
	if m.SendImageFunc != nil {
		return m.SendImageFunc(ctx, msisdn, image, caption, isViewOnce)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_img_msg_id"}, nil
}

func (m *MockGatewayClient) EditMessage(ctx context.Context, msisdn, messageID, newMessage string) error {
	if m.EditMessageFunc != nil {
		return m.EditMessageFunc(ctx, msisdn, messageID, newMessage)
	}
	return nil
}

func (m *MockGatewayClient) DeleteMessage(ctx context.Context, msisdn, messageID string) error {
	if m.DeleteMessageFunc != nil {
		return m.DeleteMessageFunc(ctx, msisdn, messageID)
	}
	return nil
}

func (m *MockGatewayClient) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string) error {
	if m.ReactToMessageFunc != nil {
		return m.ReactToMessageFunc(ctx, msisdn, messageID, emoji)
	}
	return nil
}

func (m *MockGatewayClient) GetLoginStatus(ctx context.Context) (*gateway.LoginStatus, error) {
	if m.GetLoginStatusFunc != nil {
		return m.GetLoginStatusFunc(ctx)
	}
	return &gateway.LoginStatus{Authenticated: true}, nil
}

func (m *MockGatewayClient) Health(ctx context.Context) (*gateway.HealthResponse, error) {
	if m.HealthFunc != nil {
		return m.HealthFunc(ctx)
	}
	return &gateway.HealthResponse{Status: "ok", Timestamp: "2024-01-01T00:00:00Z"}, nil
}

func (m *MockGatewayClient) GetWebhook(ctx context.Context) (*gateway.WebhookResponse, error) {
	if m.GetWebhookFunc != nil {
		return m.GetWebhookFunc(ctx)
	}
	return &gateway.WebhookResponse{URL: "https://example.com/webhook"}, nil
}

func (m *MockGatewayClient) RegisterWebhook(ctx context.Context, url, hmacSecret string) error {
	if m.RegisterWebhookFunc != nil {
		return m.RegisterWebhookFunc(ctx, url, hmacSecret)
	}
	return nil
}

func (m *MockGatewayClient) DeleteWebhook(ctx context.Context) error {
	if m.DeleteWebhookFunc != nil {
		return m.DeleteWebhookFunc(ctx)
	}
	return nil
}

// Test SendTextMessage

func TestSendTextMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "msg123"}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "Hello, World!",
	}

	_, result, err := SendTextMessage(ctx, nil, input)
	if err != nil {
		t.Fatalf("SendTextMessage() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.MessageID != "msg123" {
		t.Errorf("Expected message ID 'msg123', got '%s'", result.MessageID)
	}
}

func TestSendTextMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendMessageInput{
		To:      "",
		Message: "Hello, World!",
	}

	_, _, err := SendTextMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}

	if !errors.Is(err, errors.New("recipient address (to) is required")) && err.Error() != "recipient address (to) is required" {
		t.Errorf("Expected error about missing recipient, got: %v", err)
	}
}

func TestSendTextMessage_MissingMessage(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "",
	}

	_, _, err := SendTextMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing message")
	}

	if !errors.Is(err, errors.New("message content is required")) && err.Error() != "message content is required" {
		t.Errorf("Expected error about missing message, got: %v", err)
	}
}

func TestSendTextMessage_NoGatewayClient(t *testing.T) {
	ctx := context.Background()
	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "Hello, World!",
	}

	_, _, err := SendTextMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}

	if err.Error() != "gateway client not available" {
		t.Errorf("Expected error about missing gateway client, got: %v", err)
	}
}

func TestSendTextMessage_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "Hello, World!",
	}

	_, _, err := SendTextMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test SendImageMessage

func TestSendImageMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendImageFunc: func(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "img123"}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendImageInput{
		To:       "6281234567890@s.whatsapp.net",
		ImageURL: "https://example.com/image.jpg",
		Caption:  "Test image",
	}

	_, result, err := SendImageMessage(ctx, nil, input)
	if err != nil {
		t.Fatalf("SendImageMessage() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}
}

func TestSendImageMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendImageInput{
		To:       "",
		ImageURL: "https://example.com/image.jpg",
	}

	_, _, err := SendImageMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestSendImageMessage_MissingImageURL(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := SendImageInput{
		To:       "6281234567890@s.whatsapp.net",
		ImageURL: "",
	}

	_, _, err := SendImageMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing image URL")
	}
}

// Test EditMessage

func TestEditMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "msg123",
		NewMessage: "Edited message",
	}

	_, result, err := EditMessage(ctx, nil, input)
	if err != nil {
		t.Fatalf("EditMessage() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Status != "edited" {
		t.Errorf("Expected status 'edited', got '%s'", result.Status)
	}
}

func TestEditMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := EditMessageInput{
		To:         "",
		MessageID:  "msg123",
		NewMessage: "Edited message",
	}

	_, _, err := EditMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestEditMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "",
		NewMessage: "Edited message",
	}

	_, _, err := EditMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

func TestEditMessage_MissingNewMessage(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "msg123",
		NewMessage: "",
	}

	_, _, err := EditMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing new message")
	}
}

// Test DeleteMessage

func TestDeleteMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := DeleteMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
	}

	_, result, err := DeleteMessage(ctx, nil, input)
	if err != nil {
		t.Fatalf("DeleteMessage() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", result.Status)
	}
}

func TestDeleteMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := DeleteMessageInput{
		To:        "",
		MessageID: "msg123",
	}

	_, _, err := DeleteMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestDeleteMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := DeleteMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "",
	}

	_, _, err := DeleteMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

// Test ReactToMessage

func TestReactToMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
		Emoji:     "👍",
	}

	_, result, err := ReactToMessage(ctx, nil, input)
	if err != nil {
		t.Fatalf("ReactToMessage() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Status != "reacted" {
		t.Errorf("Expected status 'reacted', got '%s'", result.Status)
	}
}

func TestReactToMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := ReactToMessageInput{
		To:        "",
		MessageID: "msg123",
		Emoji:     "👍",
	}

	_, _, err := ReactToMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestReactToMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "",
		Emoji:     "👍",
	}

	_, _, err := ReactToMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

func TestReactToMessage_MissingEmoji(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
		Emoji:     "",
	}

	_, _, err := ReactToMessage(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing emoji")
	}
}

// Test CheckConnectionStatus

func TestCheckConnectionStatus_Success_Connected(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetLoginStatusFunc: func(ctx context.Context) (*gateway.LoginStatus, error) {
			return &gateway.LoginStatus{Authenticated: true}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := CheckConnectionStatus(ctx, nil, CheckConnectionStatusInput{})
	if err != nil {
		t.Fatalf("CheckConnectionStatus() failed: %v", err)
	}

	if !result.Authenticated {
		t.Error("Expected authenticated to be true")
	}

	if result.Status != "connected" {
		t.Errorf("Expected status 'connected', got '%s'", result.Status)
	}

	if result.Message != "WhatsApp session is active and authenticated" {
		t.Errorf("Expected message about active session, got '%s'", result.Message)
	}
}

func TestCheckConnectionStatus_Success_Disconnected(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetLoginStatusFunc: func(ctx context.Context) (*gateway.LoginStatus, error) {
			return &gateway.LoginStatus{Authenticated: false}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := CheckConnectionStatus(ctx, nil, CheckConnectionStatusInput{})
	if err != nil {
		t.Fatalf("CheckConnectionStatus() failed: %v", err)
	}

	if result.Authenticated {
		t.Error("Expected authenticated to be false")
	}

	if result.Status != "disconnected" {
		t.Errorf("Expected status 'disconnected', got '%s'", result.Status)
	}
}

func TestCheckConnectionStatus_NoGatewayClient(t *testing.T) {
	ctx := context.Background()

	_, _, err := CheckConnectionStatus(ctx, nil, CheckConnectionStatusInput{})
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}

	if err.Error() != "gateway client not available" {
		t.Errorf("Expected error about missing gateway client, got: %v", err)
	}
}

func TestCheckConnectionStatus_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetLoginStatusFunc: func(ctx context.Context) (*gateway.LoginStatus, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, _, err := CheckConnectionStatus(ctx, nil, CheckConnectionStatusInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test CheckHealth

func TestCheckHealth_Success_OK(t *testing.T) {
	mockClient := &MockGatewayClient{
		HealthFunc: func(ctx context.Context) (*gateway.HealthResponse, error) {
			return &gateway.HealthResponse{Status: "ok", Timestamp: "2024-01-01T12:00:00Z"}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := CheckHealth(ctx, nil, CheckHealthInput{})
	if err != nil {
		t.Fatalf("CheckHealth() failed: %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result.Status)
	}

	if result.Timestamp != "2024-01-01T12:00:00Z" {
		t.Errorf("Expected timestamp '2024-01-01T12:00:00Z', got '%s'", result.Timestamp)
	}

	if result.Message != "Gateway service is reachable and healthy" {
		t.Errorf("Expected message about healthy service, got '%s'", result.Message)
	}
}

func TestCheckHealth_Success_NotOK(t *testing.T) {
	mockClient := &MockGatewayClient{
		HealthFunc: func(ctx context.Context) (*gateway.HealthResponse, error) {
			return &gateway.HealthResponse{Status: "degraded", Timestamp: "2024-01-01T12:00:00Z"}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := CheckHealth(ctx, nil, CheckHealthInput{})
	if err != nil {
		t.Fatalf("CheckHealth() failed: %v", err)
	}

	if result.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", result.Status)
	}

	if result.Message != "Gateway service status: degraded" {
		t.Errorf("Expected message about service status, got '%s'", result.Message)
	}
}

func TestCheckHealth_NoGatewayClient(t *testing.T) {
	ctx := context.Background()

	_, _, err := CheckHealth(ctx, nil, CheckHealthInput{})
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}

	if err.Error() != "gateway client not available" {
		t.Errorf("Expected error about missing gateway client, got: %v", err)
	}
}

func TestCheckHealth_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		HealthFunc: func(ctx context.Context) (*gateway.HealthResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, _, err := CheckHealth(ctx, nil, CheckHealthInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test GetWebhook

func TestGetWebhook_Success_Registered(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetWebhookFunc: func(ctx context.Context) (*gateway.WebhookResponse, error) {
			return &gateway.WebhookResponse{
				URL: "https://example.com/webhook",
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := GetWebhook(ctx, nil, GetWebhookInput{})
	if err != nil {
		t.Fatalf("GetWebhook() failed: %v", err)
	}

	if result.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got '%s'", result.URL)
	}

	if result.Status != "registered" {
		t.Errorf("Expected status 'registered', got '%s'", result.Status)
	}
}

func TestGetWebhook_Success_NotRegistered(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetWebhookFunc: func(ctx context.Context) (*gateway.WebhookResponse, error) {
			return &gateway.WebhookResponse{
				URL: "",
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := GetWebhook(ctx, nil, GetWebhookInput{})
	if err != nil {
		t.Fatalf("GetWebhook() failed: %v", err)
	}

	if result.URL != "" {
		t.Errorf("Expected empty URL, got '%s'", result.URL)
	}

	if result.Status != "not_registered" {
		t.Errorf("Expected status 'not_registered', got '%s'", result.Status)
	}
}

func TestGetWebhook_NoGatewayClient(t *testing.T) {
	ctx := context.Background()

	_, _, err := GetWebhook(ctx, nil, GetWebhookInput{})
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}
}

func TestGetWebhook_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetWebhookFunc: func(ctx context.Context) (*gateway.WebhookResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, _, err := GetWebhook(ctx, nil, GetWebhookInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test RegisterWebhook

func TestRegisterWebhook_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := RegisterWebhookInput{
		URL:        "https://example.com/webhook",
		HMACSecret: "secret123",
	}

	_, result, err := RegisterWebhook(ctx, nil, input)
	if err != nil {
		t.Fatalf("RegisterWebhook() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.URL != "https://example.com/webhook" {
		t.Errorf("Expected URL 'https://example.com/webhook', got '%s'", result.URL)
	}

	if result.Status != "registered" {
		t.Errorf("Expected status 'registered', got '%s'", result.Status)
	}
}

func TestRegisterWebhook_MissingURL(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := RegisterWebhookInput{
		URL: "",
	}

	_, _, err := RegisterWebhook(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for missing URL")
	}

	if !errors.Is(err, errors.New("webhook URL is required")) && err.Error() != "webhook URL is required" {
		t.Errorf("Expected error about missing URL, got: %v", err)
	}
}

func TestRegisterWebhook_InvalidURL(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := RegisterWebhookInput{
		URL: "://invalid-url",
	}

	_, _, err := RegisterWebhook(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for invalid URL")
	}
}

func TestRegisterWebhook_InvalidURLScheme(t *testing.T) {
	mockClient := &MockGatewayClient{}
	ctx := context.WithValue(context.Background(), "gateway", mockClient)
	input := RegisterWebhookInput{
		URL: "ftp://example.com/webhook",
	}

	_, _, err := RegisterWebhook(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error for invalid URL scheme")
	}

	if !errors.Is(err, errors.New("webhook URL must use HTTP or HTTPS scheme")) && err.Error() != "webhook URL must use HTTP or HTTPS scheme" {
		t.Errorf("Expected error about URL scheme, got: %v", err)
	}
}

func TestRegisterWebhook_NoGatewayClient(t *testing.T) {
	ctx := context.Background()
	input := RegisterWebhookInput{
		URL: "https://example.com/webhook",
	}

	_, _, err := RegisterWebhook(ctx, nil, input)
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}
}

// Test DeleteWebhook

func TestDeleteWebhook_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, result, err := DeleteWebhook(ctx, nil, DeleteWebhookInput{})
	if err != nil {
		t.Fatalf("DeleteWebhook() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", result.Status)
	}
}

func TestDeleteWebhook_NoGatewayClient(t *testing.T) {
	ctx := context.Background()

	_, _, err := DeleteWebhook(ctx, nil, DeleteWebhookInput{})
	if err == nil {
		t.Fatal("Expected error when gateway client is not available")
	}
}

func TestDeleteWebhook_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		DeleteWebhookFunc: func(ctx context.Context) error {
			return errors.New("gateway connection failed")
		},
	}

	ctx := context.WithValue(context.Background(), "gateway", mockClient)

	_, _, err := DeleteWebhook(ctx, nil, DeleteWebhookInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}
