package tools

import (
	"context"
	"fmt"
	"io"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetGatewayClient retrieves the gateway client from context or returns an error
func GetGatewayClient(ctx context.Context) (gateway.GatewayClient, error) {
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, fmt.Errorf("gateway client not available")
	}
	return client, nil
}

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

// SendTextMessage sends a text message to a WhatsApp contact or group
func SendTextMessage(ctx context.Context, req *mcp.CallToolRequest, input SendMessageInput) (
	*mcp.CallToolResult,
	SendMessageResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, SendMessageResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := SendTextMessageDirect(client, input)
	if err != nil {
		return nil, SendMessageResult{}, err
	}

	return nil, result, nil
}

// SendTextMessageDirect sends a text message without using context
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

// SendImageMessage sends an image message to a WhatsApp contact or group
func SendImageMessage(ctx context.Context, req *mcp.CallToolRequest, input SendImageInput) (
	*mcp.CallToolResult,
	SendMessageResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, SendMessageResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := SendImageMessageDirect(client, input)
	if err != nil {
		return nil, SendMessageResult{}, err
	}

	return nil, result, nil
}

// SendImageMessageDirect sends an image message without using context
func SendImageMessageDirect(client gateway.GatewayClient, input SendImageInput) (SendMessageResult, error) {
	// Validate input
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.ImageURL == "" {
		return SendMessageResult{}, fmt.Errorf("image URL is required")
	}

	// For now, use a placeholder reader
	// In a real implementation, you would download the image from the URL
	// This is a placeholder - image URL handling would go here
	imageReader := io.NopCloser(nil)
	_ = input.ImageURL // Will be used to download the image

	// Send image via gateway
	ctx := context.Background()
	resp, err := client.SendImage(ctx, input.To, imageReader, input.Caption, input.ViewOnce)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_image_message: %w", err)
	}

	result := SendMessageResult{
		Success:   resp.Success,
		MessageID: resp.MessageID,
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

// EditMessage edits a previously sent message
func EditMessage(ctx context.Context, req *mcp.CallToolRequest, input EditMessageInput) (
	*mcp.CallToolResult,
	EditMessageResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, EditMessageResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := EditMessageDirect(client, input)
	if err != nil {
		return nil, EditMessageResult{}, err
	}

	return nil, result, nil
}

// EditMessageDirect edits a previously sent message without using context
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

// DeleteMessage deletes a previously sent message
func DeleteMessage(ctx context.Context, req *mcp.CallToolRequest, input DeleteMessageInput) (
	*mcp.CallToolResult,
	DeleteMessageResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, DeleteMessageResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := DeleteMessageDirect(client, input)
	if err != nil {
		return nil, DeleteMessageResult{}, err
	}

	return nil, result, nil
}

// DeleteMessageDirect deletes a previously sent message without using context
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

// ReactToMessage reacts to a message with an emoji
func ReactToMessage(ctx context.Context, req *mcp.CallToolRequest, input ReactToMessageInput) (
	*mcp.CallToolResult,
	ReactToMessageResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, ReactToMessageResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := ReactToMessageDirect(client, input)
	if err != nil {
		return nil, ReactToMessageResult{}, err
	}

	return nil, result, nil
}

// ReactToMessageDirect reacts to a message with an emoji without using context
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
