package tools

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

const imageDownloadTimeout = 20 * time.Second

// SendMessageInput represents the input for sending a text message
type SendMessageInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// SendMessageResult represents the result of sending a message
type SendMessageResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// SendTextMessageDirect sends a text message
func SendTextMessageDirect(client gateway.GatewayClient, input SendMessageInput) (SendMessageResult, error) {
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.Message == "" {
		return SendMessageResult{}, fmt.Errorf("message content is required")
	}

	ctx := context.Background()
	resp, err := client.SendText(ctx, waga.FormatMSISDN(input.To), input.Message)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_text_message: %w", err)
	}

	return SendMessageResult{
		Success:   resp.Success,
		MessageID: resp.MessageID,
		Status:    "sent",
	}, nil
}

// SendImageInput represents the input for sending an image message
type SendImageInput struct {
	To       string `json:"to"`
	ImageURL string `json:"image_url"`
	Caption  string `json:"caption"`
	ViewOnce bool   `json:"view_once"`
}

// SendImageMessageDirect sends an image message
func SendImageMessageDirect(client gateway.GatewayClient, input SendImageInput) (SendMessageResult, error) {
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.ImageURL == "" {
		return SendMessageResult{}, fmt.Errorf("image URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), imageDownloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.ImageURL, nil)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: invalid image URL: %w", err)
	}
	dlResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: download failed: %w", err)
	}
	defer dlResp.Body.Close()
	if dlResp.StatusCode < 200 || dlResp.StatusCode >= 300 {
		return SendMessageResult{}, fmt.Errorf("send_image_message: image URL returned %d", dlResp.StatusCode)
	}

	resp, err := client.SendImage(ctx, waga.FormatMSISDN(input.To), dlResp.Body, input.Caption, input.ViewOnce)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: %w", err)
	}

	return SendMessageResult{
		Success:   resp.Success,
		MessageID: resp.MessageID,
		Status:    "sent",
	}, nil
}

// EditMessageInput represents the input for editing a message
type EditMessageInput struct {
	To         string `json:"to"`
	MessageID  string `json:"message_id"`
	NewMessage string `json:"new_message"`
}

// EditMessageResult represents the result of editing a message
type EditMessageResult struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

// EditMessageDirect edits a previously sent message
func EditMessageDirect(client gateway.GatewayClient, input EditMessageInput) (EditMessageResult, error) {
	if input.To == "" {
		return EditMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return EditMessageResult{}, fmt.Errorf("message ID is required")
	}
	if input.NewMessage == "" {
		return EditMessageResult{}, fmt.Errorf("new message content is required")
	}

	ctx := context.Background()
	err := client.EditMessage(ctx, waga.FormatMSISDN(input.To), input.MessageID, input.NewMessage)
	if err != nil {
		return EditMessageResult{}, fmt.Errorf("edit_message: %w", err)
	}

	return EditMessageResult{
		Success: true,
		Status:  "edited",
	}, nil
}

// DeleteMessageInput represents the input for deleting a message
type DeleteMessageInput struct {
	To        string `json:"to"`
	MessageID string `json:"message_id"`
}

// DeleteMessageResult represents the result of deleting a message
type DeleteMessageResult struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

// DeleteMessageDirect deletes a previously sent message
func DeleteMessageDirect(client gateway.GatewayClient, input DeleteMessageInput) (DeleteMessageResult, error) {
	if input.To == "" {
		return DeleteMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return DeleteMessageResult{}, fmt.Errorf("message ID is required")
	}

	ctx := context.Background()
	err := client.DeleteMessage(ctx, waga.FormatMSISDN(input.To), input.MessageID)
	if err != nil {
		return DeleteMessageResult{}, fmt.Errorf("delete_message: %w", err)
	}

	return DeleteMessageResult{
		Success: true,
		Status:  "deleted",
	}, nil
}

// ReactToMessageInput represents the input for reacting to a message
type ReactToMessageInput struct {
	To        string `json:"to"`
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

// ReactToMessageResult represents the result of reacting to a message
type ReactToMessageResult struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

// ReactToMessageDirect reacts to a message with an emoji
func ReactToMessageDirect(client gateway.GatewayClient, input ReactToMessageInput) (ReactToMessageResult, error) {
	if input.To == "" {
		return ReactToMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return ReactToMessageResult{}, fmt.Errorf("message ID is required")
	}
	if input.Emoji == "" {
		return ReactToMessageResult{}, fmt.Errorf("emoji is required")
	}

	ctx := context.Background()
	err := client.ReactToMessage(ctx, waga.FormatMSISDN(input.To), input.MessageID, input.Emoji)
	if err != nil {
		return ReactToMessageResult{}, fmt.Errorf("react_to_message: %w", err)
	}

	return ReactToMessageResult{
		Success: true,
		Status:  "reacted",
	}, nil
}
