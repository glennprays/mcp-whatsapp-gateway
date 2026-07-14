package tools

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// MockGatewayClient is a mock implementation of GatewayClient for testing
type MockGatewayClient struct {
	SendTextFunc            func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error)
	SendImageFunc           func(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error)
	SendAudioFunc           func(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool) (*gateway.SendMessageResponse, error)
	SendVideoFunc           func(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool) (*gateway.SendMessageResponse, error)
	SendDocumentFunc        func(ctx context.Context, msisdn string, document io.Reader, fileName, caption string) (*gateway.SendMessageResponse, error)
	SendLocationFunc        func(ctx context.Context, msisdn string, latitude, longitude float64, name, address string) (*gateway.SendMessageResponse, error)
	SendPollFunc            func(ctx context.Context, msisdn, question string, options []string, selectableCount int) (*gateway.SendMessageResponse, error)
	SendStickerFunc         func(ctx context.Context, msisdn string, sticker io.Reader) (*gateway.SendMessageResponse, error)
	EditMessageFunc         func(ctx context.Context, msisdn, messageID, newMessage string) error
	DeleteMessageFunc       func(ctx context.Context, msisdn, messageID string) error
	ReactToMessageFunc      func(ctx context.Context, msisdn, messageID, emoji string, senderMsisdn ...string) error
	GetIncomingMessagesFunc func(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error)
	GetJobStatusFunc        func(ctx context.Context, jobID string) (*gateway.JobStatusResponse, error)
	GetLoginStatusFunc      func(ctx context.Context) (*gateway.LoginStatus, error)
	CheckContactFunc        func(ctx context.Context, msisdn string) (*gateway.ContactCheckResponse, error)
	HealthFunc              func(ctx context.Context) (*gateway.HealthResponse, error)
	GetWebhookFunc          func(ctx context.Context) (*gateway.WebhookResponse, error)
	RegisterWebhookFunc     func(ctx context.Context, url, hmacSecret string) error
	DeleteWebhookFunc       func(ctx context.Context) error

	// CapturedSendOpts records the SendOptions passed to the most recent send,
	// so tests can assert chat/reply/mention options were built and forwarded.
	CapturedSendOpts []waga.SendOption
}

func (m *MockGatewayClient) SendText(ctx context.Context, msisdn, message string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendTextFunc != nil {
		return m.SendTextFunc(ctx, msisdn, message)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_msg_id"}, nil
}

func (m *MockGatewayClient) SendImage(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendImageFunc != nil {
		return m.SendImageFunc(ctx, msisdn, image, caption, isViewOnce)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_img_msg_id"}, nil
}

func (m *MockGatewayClient) SendAudio(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendAudioFunc != nil {
		return m.SendAudioFunc(ctx, msisdn, audio, isPTT, isViewOnce)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_audio_msg_id"}, nil
}

func (m *MockGatewayClient) SendVideo(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendVideoFunc != nil {
		return m.SendVideoFunc(ctx, msisdn, video, caption, isGif, isViewOnce)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_video_msg_id"}, nil
}

func (m *MockGatewayClient) SendDocument(ctx context.Context, msisdn string, document io.Reader, fileName, caption string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendDocumentFunc != nil {
		return m.SendDocumentFunc(ctx, msisdn, document, fileName, caption)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_document_msg_id"}, nil
}

func (m *MockGatewayClient) SendLocation(ctx context.Context, msisdn string, latitude, longitude float64, name, address string, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendLocationFunc != nil {
		return m.SendLocationFunc(ctx, msisdn, latitude, longitude, name, address)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_loc_msg_id"}, nil
}

func (m *MockGatewayClient) SendPoll(ctx context.Context, msisdn, question string, options []string, selectableCount int, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendPollFunc != nil {
		return m.SendPollFunc(ctx, msisdn, question, options, selectableCount)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_poll_msg_id"}, nil
}

func (m *MockGatewayClient) SendSticker(ctx context.Context, msisdn string, sticker io.Reader, opts ...waga.SendOption) (*gateway.SendMessageResponse, error) {
	m.CapturedSendOpts = opts
	if m.SendStickerFunc != nil {
		return m.SendStickerFunc(ctx, msisdn, sticker)
	}
	return &gateway.SendMessageResponse{Success: true, MessageID: "test_sticker_msg_id"}, nil
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

func (m *MockGatewayClient) ReactToMessage(ctx context.Context, msisdn, messageID, emoji string, senderMsisdn ...string) error {
	if m.ReactToMessageFunc != nil {
		return m.ReactToMessageFunc(ctx, msisdn, messageID, emoji, senderMsisdn...)
	}
	return nil
}

func (m *MockGatewayClient) GetLoginStatus(ctx context.Context) (*gateway.LoginStatus, error) {
	if m.GetLoginStatusFunc != nil {
		return m.GetLoginStatusFunc(ctx)
	}
	return &gateway.LoginStatus{Authenticated: true}, nil
}

func (m *MockGatewayClient) CheckContact(ctx context.Context, msisdn string) (*gateway.ContactCheckResponse, error) {
	if m.CheckContactFunc != nil {
		return m.CheckContactFunc(ctx, msisdn)
	}
	return &gateway.ContactCheckResponse{
		Query:        msisdn,
		JID:          msisdn,
		IsOnWhatsApp: true,
	}, nil
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

func (m *MockGatewayClient) GetIncomingMessages(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
	if m.GetIncomingMessagesFunc != nil {
		return m.GetIncomingMessagesFunc(ctx, limit)
	}
	return &waga.IncomingMessagesResponse{Success: true, Count: 0, Messages: []waga.IncomingMessage{}}, nil
}

func (m *MockGatewayClient) GetJobStatus(ctx context.Context, jobID string) (*gateway.JobStatusResponse, error) {
	if m.GetJobStatusFunc != nil {
		return m.GetJobStatusFunc(ctx, jobID)
	}
	return &gateway.JobStatusResponse{JobID: jobID, Status: "completed"}, nil
}

// Test SendTextMessageDirect

func TestSendTextMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "msg123"}, nil
		},
	}

	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "Hello, World!",
	}

	result, err := SendTextMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("SendTextMessageDirect() failed: %v", err)
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
	input := SendMessageInput{
		To:      "",
		Message: "Hello, World!",
	}

	_, err := SendTextMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}

	if err.Error() != "recipient address (chat or to) is required" {
		t.Errorf("Expected error about missing recipient, got: %v", err)
	}
}

// Test chat/to precedence and send-option forwarding (commit 2)

func TestResolveSendAddr_ChatWinsOverTo(t *testing.T) {
	addr, err := resolveSendAddr("123@g.us", "628@s.whatsapp.net")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "123@g.us" {
		t.Errorf("expected chat to win, got %q", addr)
	}
}

func TestResolveSendAddr_FallsBackToTo(t *testing.T) {
	addr, err := resolveSendAddr("", "628")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "628" {
		t.Errorf("expected fallback to to, got %q", addr)
	}
}

func TestResolveSendAddr_RequiresRecipient(t *testing.T) {
	if _, err := resolveSendAddr("", ""); err == nil {
		t.Error("expected error when neither chat nor to is set")
	}
}

func TestSendText_ChatReplyMentionsForwarded(t *testing.T) {
	mockClient := &MockGatewayClient{}
	_, err := SendTextMessageDirect(mockClient, SendMessageInput{
		Chat:          "123@g.us",
		To:            "628@s.whatsapp.net",
		Message:       "hi",
		ReplyToID:     "MSG1",
		ReplyToSender: "629@s.whatsapp.net",
		ReplyToText:   "original",
		Mentions:      []string{"629@s.whatsapp.net"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// WithChat + WithReply + WithMentions
	if len(mockClient.CapturedSendOpts) != 3 {
		t.Errorf("expected 3 send options forwarded, got %d", len(mockClient.CapturedSendOpts))
	}
}

func TestSendText_ChatOnlyForwardsSingleOption(t *testing.T) {
	mockClient := &MockGatewayClient{}
	_, err := SendTextMessageDirect(mockClient, SendMessageInput{Chat: "628", Message: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mockClient.CapturedSendOpts) != 1 {
		t.Errorf("expected 1 send option (WithChat), got %d", len(mockClient.CapturedSendOpts))
	}
}

func TestSendTextMessage_MissingMessage(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "",
	}

	_, err := SendTextMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing message")
	}

	if err.Error() != "message content is required" {
		t.Errorf("Expected error about missing message, got: %v", err)
	}
}

func TestSendTextMessage_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	input := SendMessageInput{
		To:      "6281234567890@s.whatsapp.net",
		Message: "Hello, World!",
	}

	_, err := SendTextMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test SendImageMessageDirect

func TestSendImageMessage_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("fake-image-data"))
	}))
	defer ts.Close()

	mockClient := &MockGatewayClient{
		SendImageFunc: func(ctx context.Context, msisdn string, image io.Reader, caption string, isViewOnce bool) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "img123"}, nil
		},
	}

	input := SendImageInput{
		To:       "6281234567890@s.whatsapp.net",
		ImageURL: ts.URL + "/image.jpg",
		Caption:  "Test image",
	}

	result, err := SendImageMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("SendImageMessageDirect() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.MessageID != "img123" {
		t.Errorf("Expected message ID 'img123', got '%s'", result.MessageID)
	}
}

func TestSendImageMessage_MissingRecipient(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := SendImageInput{
		To:       "",
		ImageURL: "https://example.com/image.jpg",
	}

	_, err := SendImageMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestSendImageMessage_MissingImageURL(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := SendImageInput{
		To:       "6281234567890@s.whatsapp.net",
		ImageURL: "",
	}

	_, err := SendImageMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing image URL")
	}
}

func TestSendImageMessage_DownloadFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	mockClient := &MockGatewayClient{}
	input := SendImageInput{
		To:       "6281234567890@s.whatsapp.net",
		ImageURL: ts.URL + "/missing.jpg",
	}

	_, err := SendImageMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for failed image download")
	}
}

func TestSendAudioMessage_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "audio/ogg")
		w.Write([]byte("fake-audio-data"))
	}))
	defer ts.Close()

	mockClient := &MockGatewayClient{
		SendAudioFunc: func(ctx context.Context, msisdn string, audio io.Reader, isPTT, isViewOnce bool) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "audio123"}, nil
		},
	}

	result, err := SendAudioMessageDirect(mockClient, SendAudioInput{
		To:       "6281234567890@s.whatsapp.net",
		AudioURL: ts.URL + "/audio.ogg",
		IsPTT:    true,
	})
	if err != nil {
		t.Fatalf("SendAudioMessageDirect() failed: %v", err)
	}
	if result.MessageID != "audio123" {
		t.Errorf("Expected message ID 'audio123', got '%s'", result.MessageID)
	}
}

func TestSendVideoMessage_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		w.Write([]byte("fake-video-data"))
	}))
	defer ts.Close()

	mockClient := &MockGatewayClient{
		SendVideoFunc: func(ctx context.Context, msisdn string, video io.Reader, caption string, isGif, isViewOnce bool) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "video123"}, nil
		},
	}

	result, err := SendVideoMessageDirect(mockClient, SendVideoInput{
		To:       "6281234567890@s.whatsapp.net",
		VideoURL: ts.URL + "/video.mp4",
		Caption:  "hello",
	})
	if err != nil {
		t.Fatalf("SendVideoMessageDirect() failed: %v", err)
	}
	if result.MessageID != "video123" {
		t.Errorf("Expected message ID 'video123', got '%s'", result.MessageID)
	}
}

func TestSendDocumentMessage_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake-document-data"))
	}))
	defer ts.Close()

	mockClient := &MockGatewayClient{
		SendDocumentFunc: func(ctx context.Context, msisdn string, document io.Reader, fileName, caption string) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "doc123"}, nil
		},
	}

	result, err := SendDocumentMessageDirect(mockClient, SendDocumentInput{
		To:          "6281234567890@s.whatsapp.net",
		DocumentURL: ts.URL + "/doc.pdf",
		FileName:    "doc.pdf",
	})
	if err != nil {
		t.Fatalf("SendDocumentMessageDirect() failed: %v", err)
	}
	if result.MessageID != "doc123" {
		t.Errorf("Expected message ID 'doc123', got '%s'", result.MessageID)
	}
}

// Test EditMessageDirect

func TestEditMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "msg123",
		NewMessage: "Edited message",
	}

	result, err := EditMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("EditMessageDirect() failed: %v", err)
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
	input := EditMessageInput{
		To:         "",
		MessageID:  "msg123",
		NewMessage: "Edited message",
	}

	_, err := EditMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestEditMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "",
		NewMessage: "Edited message",
	}

	_, err := EditMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

func TestEditMessage_MissingNewMessage(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := EditMessageInput{
		To:         "6281234567890@s.whatsapp.net",
		MessageID:  "msg123",
		NewMessage: "",
	}

	_, err := EditMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing new message")
	}
}

// Test DeleteMessageDirect

func TestDeleteMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := DeleteMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
	}

	result, err := DeleteMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("DeleteMessageDirect() failed: %v", err)
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
	input := DeleteMessageInput{
		To:        "",
		MessageID: "msg123",
	}

	_, err := DeleteMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestDeleteMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := DeleteMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "",
	}

	_, err := DeleteMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

// Test ReactToMessageDirect

func TestReactToMessage_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
		Emoji:     "\U0001f44d",
	}

	result, err := ReactToMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("ReactToMessageDirect() failed: %v", err)
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
	input := ReactToMessageInput{
		To:        "",
		MessageID: "msg123",
		Emoji:     "\U0001f44d",
	}

	_, err := ReactToMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing recipient")
	}
}

func TestReactToMessage_MissingMessageID(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "",
		Emoji:     "\U0001f44d",
	}

	_, err := ReactToMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing message ID")
	}
}

func TestReactToMessage_MissingEmoji(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := ReactToMessageInput{
		To:        "6281234567890@s.whatsapp.net",
		MessageID: "msg123",
		Emoji:     "",
	}

	_, err := ReactToMessageDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing emoji")
	}
}

// Test CheckConnectionStatusDirect

func TestCheckConnectionStatus_Success_Connected(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetLoginStatusFunc: func(ctx context.Context) (*gateway.LoginStatus, error) {
			return &gateway.LoginStatus{Authenticated: true}, nil
		},
	}

	result, err := CheckConnectionStatusDirect(mockClient, CheckConnectionStatusInput{})
	if err != nil {
		t.Fatalf("CheckConnectionStatusDirect() failed: %v", err)
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

	result, err := CheckConnectionStatusDirect(mockClient, CheckConnectionStatusInput{})
	if err != nil {
		t.Fatalf("CheckConnectionStatusDirect() failed: %v", err)
	}

	if result.Authenticated {
		t.Error("Expected authenticated to be false")
	}

	if result.Status != "disconnected" {
		t.Errorf("Expected status 'disconnected', got '%s'", result.Status)
	}
}

func TestCheckConnectionStatus_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetLoginStatusFunc: func(ctx context.Context) (*gateway.LoginStatus, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	_, err := CheckConnectionStatusDirect(mockClient, CheckConnectionStatusInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test CheckHealthDirect

func TestCheckHealth_Success_OK(t *testing.T) {
	mockClient := &MockGatewayClient{
		HealthFunc: func(ctx context.Context) (*gateway.HealthResponse, error) {
			return &gateway.HealthResponse{Status: "ok", Timestamp: "2024-01-01T12:00:00Z"}, nil
		},
	}

	result, err := CheckHealthDirect(mockClient, CheckHealthInput{})
	if err != nil {
		t.Fatalf("CheckHealthDirect() failed: %v", err)
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

	result, err := CheckHealthDirect(mockClient, CheckHealthInput{})
	if err != nil {
		t.Fatalf("CheckHealthDirect() failed: %v", err)
	}

	if result.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", result.Status)
	}

	if result.Message != "Gateway service status: degraded" {
		t.Errorf("Expected message about service status, got '%s'", result.Message)
	}
}

func TestCheckHealth_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		HealthFunc: func(ctx context.Context) (*gateway.HealthResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	_, err := CheckHealthDirect(mockClient, CheckHealthInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test GetWebhookDirect

func TestGetWebhook_Success_Registered(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetWebhookFunc: func(ctx context.Context) (*gateway.WebhookResponse, error) {
			return &gateway.WebhookResponse{
				URL: "https://example.com/webhook",
			}, nil
		},
	}

	result, err := GetWebhookDirect(mockClient, GetWebhookInput{})
	if err != nil {
		t.Fatalf("GetWebhookDirect() failed: %v", err)
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

	result, err := GetWebhookDirect(mockClient, GetWebhookInput{})
	if err != nil {
		t.Fatalf("GetWebhookDirect() failed: %v", err)
	}

	if result.URL != "" {
		t.Errorf("Expected empty URL, got '%s'", result.URL)
	}

	if result.Status != "not_registered" {
		t.Errorf("Expected status 'not_registered', got '%s'", result.Status)
	}
}

func TestGetWebhook_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetWebhookFunc: func(ctx context.Context) (*gateway.WebhookResponse, error) {
			return nil, errors.New("gateway connection failed")
		},
	}

	_, err := GetWebhookDirect(mockClient, GetWebhookInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test RegisterWebhookDirect

func TestRegisterWebhook_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := RegisterWebhookInput{
		URL:        "https://example.com/webhook",
		HMACSecret: "secret123",
	}

	result, err := RegisterWebhookDirect(mockClient, input)
	if err != nil {
		t.Fatalf("RegisterWebhookDirect() failed: %v", err)
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
	input := RegisterWebhookInput{
		URL: "",
	}

	_, err := RegisterWebhookDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for missing URL")
	}

	if err.Error() != "webhook URL is required" {
		t.Errorf("Expected error about missing URL, got: %v", err)
	}
}

func TestRegisterWebhook_InvalidURL(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := RegisterWebhookInput{
		URL: "://invalid-url",
	}

	_, err := RegisterWebhookDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for invalid URL")
	}
}

func TestRegisterWebhook_InvalidURLScheme(t *testing.T) {
	mockClient := &MockGatewayClient{}
	input := RegisterWebhookInput{
		URL: "ftp://example.com/webhook",
	}

	_, err := RegisterWebhookDirect(mockClient, input)
	if err == nil {
		t.Fatal("Expected error for invalid URL scheme")
	}

	if err.Error() != "webhook URL must use HTTP or HTTPS scheme" {
		t.Errorf("Expected error about URL scheme, got: %v", err)
	}
}

// Test DeleteWebhookDirect

func TestDeleteWebhook_Success(t *testing.T) {
	mockClient := &MockGatewayClient{}

	result, err := DeleteWebhookDirect(mockClient, DeleteWebhookInput{})
	if err != nil {
		t.Fatalf("DeleteWebhookDirect() failed: %v", err)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", result.Status)
	}
}

func TestDeleteWebhook_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		DeleteWebhookFunc: func(ctx context.Context) error {
			return errors.New("gateway connection failed")
		},
	}

	_, err := DeleteWebhookDirect(mockClient, DeleteWebhookInput{})
	if err == nil {
		t.Fatal("Expected error from gateway")
	}
}

// Test GetLatestIncomingMessagesDirect

func TestGetLatestIncomingMessages_Success(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetIncomingMessagesFunc: func(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
			if limit != 5 {
				t.Errorf("expected limit=5, got %d", limit)
			}
			return &waga.IncomingMessagesResponse{
				Success:   true,
				Timestamp: 1715900000000,
				Count:     2,
				Messages: []waga.IncomingMessage{
					{MessageId: "NEW", From: "6281@s.whatsapp.net", PushName: "Alice", Type: waga.IncomingMessageTypeText, Text: "Your code is 123456", Timestamp: 2},
					{MessageId: "OLD", From: "6282@s.whatsapp.net", PushName: "Bob", Type: waga.IncomingMessageTypeText, Text: "hi", Timestamp: 1},
				},
			}, nil
		},
	}

	result, err := GetLatestIncomingMessagesDirect(mockClient, GetLatestIncomingMessagesInput{Limit: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Error("expected Success=true")
	}
	if result.Count != 2 {
		t.Errorf("expected Count=2, got %d", result.Count)
	}
	if len(result.Messages) != 2 || result.Messages[0].MessageId != "NEW" {
		t.Errorf("expected newest-first ordering, got %+v", result.Messages)
	}
}

func TestGetLatestIncomingMessages_DefaultLimitPassthrough(t *testing.T) {
	var capturedLimit int
	mockClient := &MockGatewayClient{
		GetIncomingMessagesFunc: func(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
			capturedLimit = limit
			return &waga.IncomingMessagesResponse{Success: true, Count: 0, Messages: []waga.IncomingMessage{}}, nil
		},
	}

	_, err := GetLatestIncomingMessagesDirect(mockClient, GetLatestIncomingMessagesInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedLimit != 0 {
		t.Errorf("expected limit=0 (gateway substitutes default), got %d", capturedLimit)
	}
}

func TestGetLatestIncomingMessages_Empty(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetIncomingMessagesFunc: func(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
			return &waga.IncomingMessagesResponse{Success: true, Count: 0, Messages: []waga.IncomingMessage{}}, nil
		},
	}

	result, err := GetLatestIncomingMessagesDirect(mockClient, GetLatestIncomingMessagesInput{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 0 || len(result.Messages) != 0 {
		t.Errorf("expected empty response, got count=%d len=%d", result.Count, len(result.Messages))
	}
}

func TestGetLatestIncomingMessages_GatewayError(t *testing.T) {
	mockClient := &MockGatewayClient{
		GetIncomingMessagesFunc: func(ctx context.Context, limit int) (*waga.IncomingMessagesResponse, error) {
			return nil, errors.New("gateway exploded")
		},
	}

	_, err := GetLatestIncomingMessagesDirect(mockClient, GetLatestIncomingMessagesInput{Limit: 10})
	if err == nil {
		t.Fatal("expected error from gateway")
	}
	if got := err.Error(); got == "" || got[:len("get_latest_incoming_messages")] != "get_latest_incoming_messages" {
		t.Errorf("expected error prefixed with tool name, got %q", got)
	}
}

func TestCheckContact_Success(t *testing.T) {
	mockClient := &MockGatewayClient{
		CheckContactFunc: func(ctx context.Context, msisdn string) (*gateway.ContactCheckResponse, error) {
			return &gateway.ContactCheckResponse{
				Query:        "6281234567890",
				JID:          "6281234567890@s.whatsapp.net",
				IsOnWhatsApp: true,
			}, nil
		},
	}

	res, err := CheckContactDirect(mockClient, CheckContactInput{Msisdn: "6281234567890"})
	if err != nil {
		t.Fatalf("CheckContactDirect failed: %v", err)
	}
	if !res.IsOnWhatsApp {
		t.Fatal("expected IsOnWhatsApp true")
	}
}

// Test SendLocationMessageDirect coordinate validation

func TestSendLocationMessage_NullIslandIsValid(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendLocationFunc: func(ctx context.Context, msisdn string, latitude, longitude float64, name, address string) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "msg123"}, nil
		},
	}

	input := SendLocationInput{To: "6281234567890@s.whatsapp.net", Latitude: 0, Longitude: 0}
	result, err := SendLocationMessageDirect(mockClient, input)
	if err != nil {
		t.Fatalf("coordinates (0,0) must be valid, got error: %v", err)
	}
	if !result.Success {
		t.Error("Expected success to be true")
	}
}

func TestSendLocationMessage_OutOfRangeLatitude(t *testing.T) {
	input := SendLocationInput{To: "6281234567890@s.whatsapp.net", Latitude: 91, Longitude: 0}
	if _, err := SendLocationMessageDirect(&MockGatewayClient{}, input); err == nil {
		t.Error("expected error for latitude > 90")
	}
}

func TestSendLocationMessage_OutOfRangeLongitude(t *testing.T) {
	input := SendLocationInput{To: "6281234567890@s.whatsapp.net", Latitude: 0, Longitude: -181}
	if _, err := SendLocationMessageDirect(&MockGatewayClient{}, input); err == nil {
		t.Error("expected error for longitude < -180")
	}
}

func TestSendTextMessage_QueueModeSurfacesJobID(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			// Queue mode: 202 with status+job_id, empty message_id
			return &gateway.SendMessageResponse{Success: true, Status: "queued", JobID: "job-77"}, nil
		},
	}
	res, err := SendTextMessageDirect(mockClient, SendMessageInput{To: "6281234567890", Message: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "queued" {
		t.Errorf("expected status queued, got %q", res.Status)
	}
	if res.JobID != "job-77" {
		t.Errorf("expected job id job-77, got %q", res.JobID)
	}
	if res.MessageID != "" {
		t.Errorf("expected empty message id in queue mode, got %q", res.MessageID)
	}
}

func TestSendTextMessage_DirectModeStatusSent(t *testing.T) {
	mockClient := &MockGatewayClient{
		SendTextFunc: func(ctx context.Context, msisdn, message string) (*gateway.SendMessageResponse, error) {
			return &gateway.SendMessageResponse{Success: true, MessageID: "msg-1"}, nil
		},
	}
	res, _ := SendTextMessageDirect(mockClient, SendMessageInput{To: "6281234567890", Message: "hi"})
	if res.Status != "sent" || res.MessageID != "msg-1" {
		t.Errorf("expected sent/msg-1, got %q/%q", res.Status, res.MessageID)
	}
}

func TestGetJobStatus(t *testing.T) {
	msgID := "3EB0ABC"
	mockClient := &MockGatewayClient{
		GetJobStatusFunc: func(ctx context.Context, jobID string) (*gateway.JobStatusResponse, error) {
			return &gateway.JobStatusResponse{JobID: jobID, Status: "completed", MessageID: &msgID}, nil
		},
	}
	res, err := GetJobStatusDirect(mockClient, GetJobStatusInput{JobID: "job-77"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "completed" || res.MessageID != "3EB0ABC" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestGetJobStatus_RequiresJobID(t *testing.T) {
	if _, err := GetJobStatusDirect(&MockGatewayClient{}, GetJobStatusInput{}); err == nil {
		t.Error("expected error for empty job_id")
	}
}
