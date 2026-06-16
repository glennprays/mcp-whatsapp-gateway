package tools

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// mediaDownloadTimeout matches the gateway client timeout so media
// downloads never fail earlier than the gateway request itself would.
const mediaDownloadTimeout = 30 * time.Second

// SendMessageInput represents the input for sending a text message
type SendMessageInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// SendMessageResult represents the result of sending a message.
//
// In queue mode the gateway returns Status "queued" and a JobID instead of a
// MessageID; poll get_job_status with the JobID to obtain the final outcome.
type SendMessageResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Status    string `json:"status"`
	JobID     string `json:"job_id,omitempty"`
}

// toSendResult maps a gateway send response into the tool result, deriving an
// honest status: the gateway's own status when present (e.g. "queued"),
// otherwise "sent" when a message ID came back, else "unknown".
func toSendResult(resp *gateway.SendMessageResponse) SendMessageResult {
	status := resp.Status
	if status == "" {
		if resp.MessageID != "" {
			status = "sent"
		} else {
			status = "unknown"
		}
	}
	return SendMessageResult{
		Success:   resp.Success,
		MessageID: resp.MessageID,
		Status:    status,
		JobID:     resp.JobID,
	}
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

	return toSendResult(resp), nil
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

	ctx, cancel := context.WithTimeout(context.Background(), mediaDownloadTimeout)
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

	return toSendResult(resp), nil
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

// SendLocationInput represents the input for sending a location message
type SendLocationInput struct {
	To        string  `json:"to"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// SendLocationMessageDirect sends a location message
func SendLocationMessageDirect(client gateway.GatewayClient, input SendLocationInput) (SendMessageResult, error) {
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.Latitude < -90 || input.Latitude > 90 {
		return SendMessageResult{}, fmt.Errorf("latitude must be between -90 and 90")
	}
	if input.Longitude < -180 || input.Longitude > 180 {
		return SendMessageResult{}, fmt.Errorf("longitude must be between -180 and 180")
	}

	ctx := context.Background()
	resp, err := client.SendLocation(ctx, waga.FormatMSISDN(input.To), input.Latitude, input.Longitude, input.Name, input.Address)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_location_message: %w", err)
	}

	return toSendResult(resp), nil
}

// SendPollInput represents the input for sending a poll message
type SendPollInput struct {
	To              string   `json:"to"`
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	SelectableCount int      `json:"selectable_count,omitempty"`
}

// SendPollMessageDirect sends a poll message
func SendPollMessageDirect(client gateway.GatewayClient, input SendPollInput) (SendMessageResult, error) {
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.Question == "" {
		return SendMessageResult{}, fmt.Errorf("question is required")
	}
	if len(input.Options) < 2 {
		return SendMessageResult{}, fmt.Errorf("at least 2 options are required")
	}

	ctx := context.Background()
	resp, err := client.SendPoll(ctx, waga.FormatMSISDN(input.To), input.Question, input.Options, input.SelectableCount)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_poll_message: %w", err)
	}

	return toSendResult(resp), nil
}

// SendStickerInput represents the input for sending a sticker message
type SendStickerInput struct {
	To         string `json:"to"`
	StickerURL string `json:"sticker_url"`
}

// SendStickerMessageDirect sends a sticker message
func SendStickerMessageDirect(client gateway.GatewayClient, input SendStickerInput) (SendMessageResult, error) {
	if input.To == "" {
		return SendMessageResult{}, fmt.Errorf("recipient address (to) is required")
	}
	if input.StickerURL == "" {
		return SendMessageResult{}, fmt.Errorf("sticker URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), mediaDownloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.StickerURL, nil)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_sticker_message: invalid sticker URL: %w", err)
	}
	dlResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_sticker_message: download failed: %w", err)
	}
	defer dlResp.Body.Close()
	if dlResp.StatusCode < 200 || dlResp.StatusCode >= 300 {
		return SendMessageResult{}, fmt.Errorf("send_sticker_message: sticker URL returned %d", dlResp.StatusCode)
	}

	resp, err := client.SendSticker(ctx, waga.FormatMSISDN(input.To), dlResp.Body)
	if err != nil {
		return SendMessageResult{}, fmt.Errorf("send_sticker_message: %w", err)
	}

	return toSendResult(resp), nil
}

// GetJobStatusInput represents the input for polling a queued message job.
type GetJobStatusInput struct {
	JobID string `json:"job_id"`
}

// JobStatusResult represents the status of a queued message job.
type JobStatusResult struct {
	JobID       string `json:"job_id"`
	Status      string `json:"status"`
	MessageID   string `json:"message_id,omitempty"`
	Error       string `json:"error,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// GetJobStatusDirect polls the status of a queued message job by its ID.
func GetJobStatusDirect(client gateway.GatewayClient, input GetJobStatusInput) (JobStatusResult, error) {
	if input.JobID == "" {
		return JobStatusResult{}, fmt.Errorf("job_id is required")
	}

	resp, err := client.GetJobStatus(context.Background(), input.JobID)
	if err != nil {
		return JobStatusResult{}, fmt.Errorf("get_job_status: %w", err)
	}

	result := JobStatusResult{
		JobID:     resp.JobID,
		Status:    resp.Status,
		CreatedAt: resp.CreatedAt,
	}
	if resp.MessageID != nil {
		result.MessageID = *resp.MessageID
	}
	if resp.Error != nil {
		result.Error = *resp.Error
	}
	if resp.CompletedAt != nil {
		result.CompletedAt = *resp.CompletedAt
	}
	return result, nil
}
