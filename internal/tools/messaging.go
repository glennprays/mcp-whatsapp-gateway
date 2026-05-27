package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
)

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
	// Validate input
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.Message == "" {
		return SendMessageResult{}, fmt.Errorf("message content is required")
	}

	// Send message via gateway
	ctx := context.Background()
	resp, err := client.SendText(ctx, input.To, input.Message)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_text_message: %w", err)
	}

	result := SendMessageResult{
		Success:   resp.Success,
		MessageID: resp.MessageID,
		Status:    "sent",
	}

	return result, nil
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
	// Validate input
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.ImageURL == "" {
		return SendMessageResult{}, fmt.Errorf("image URL is required")
	}

	resp, err := http.Get(input.ImageURL)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SendMessageResult{}, fmt.Errorf("send_image_message: image download returned status %d", resp.StatusCode)
	}

	ctx := context.Background()
	sendResp, err := client.SendImage(ctx, input.To, resp.Body, input.Caption, input.ViewOnce)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: %w", err)
	}

	result := SendMessageResult{
		Success:   sendResp.Success,
		MessageID: sendResp.MessageID,
		Status:    "sent",
	}

	return result, nil
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
	// Validate input
	if input.To == "" {
		return EditMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return EditMessageResult{}, fmt.Errorf("message ID is required")
	}
	if input.NewMessage == "" {
		return EditMessageResult{}, fmt.Errorf("new message content is required")
	}

	// Edit message via gateway
	ctx := context.Background()
	err := client.EditMessage(ctx, input.To, input.MessageID, input.NewMessage)
	if err != nil {
		return EditMessageResult{}, fmt.Errorf("edit_message: %w", err)
	}

	result := EditMessageResult{
		Success: true,
		Status:  "edited",
	}

	return result, nil
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
	// Validate input
	if input.To == "" {
		return DeleteMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return DeleteMessageResult{}, fmt.Errorf("message ID is required")
	}

	// Delete message via gateway
	ctx := context.Background()
	err := client.DeleteMessage(ctx, input.To, input.MessageID)
	if err != nil {
		return DeleteMessageResult{}, fmt.Errorf("delete_message: %w", err)
	}

	result := DeleteMessageResult{
		Success: true,
		Status:  "deleted",
	}

	return result, nil
}

// ReactToMessageInput represents the input for reacting to a message
type ReactToMessageInput struct {
	To       string `json:"to"`
	MessageID string `json:"message_id"`
	Emoji    string `json:"emoji"`
}

// ReactToMessageResult represents the result of reacting to a message
type ReactToMessageResult struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

// ReactToMessageDirect reacts to a message with an emoji
func ReactToMessageDirect(client gateway.GatewayClient, input ReactToMessageInput) (ReactToMessageResult, error) {
	// Validate input
	if input.To == "" {
		return ReactToMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.MessageID == "" {
		return ReactToMessageResult{}, fmt.Errorf("message ID is required")
	}
	if input.Emoji == "" {
		return ReactToMessageResult{}, fmt.Errorf("emoji is required")
	}

	// React to message via gateway
	ctx := context.Background()
	err := client.ReactToMessage(ctx, input.To, input.MessageID, input.Emoji)
	if err != nil {
		return ReactToMessageResult{}, fmt.Errorf("react_to_message: %w", err)
	}

	result := ReactToMessageResult{
		Success: true,
		Status:  "reacted",
	}

	return result, nil
}
