package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// GetLatestIncomingMessagesInput is the request payload for the
// get_latest_incoming_messages tool.
type GetLatestIncomingMessagesInput struct {
	// Limit caps the number of messages returned. Optional; the gateway
	// substitutes 10 when zero and clamps anything above 50.
	Limit int `json:"limit,omitempty"`
}

// GetLatestIncomingMessagesResult is the structured response delivered to
// the calling agent. Messages are sorted newest first.
type GetLatestIncomingMessagesResult struct {
	Success   bool                   `json:"success"`
	Timestamp int64                  `json:"timestamp"`
	Count     int                    `json:"count"`
	Messages  []waga.IncomingMessage `json:"messages"`
}

// GetLatestIncomingMessagesDirect fetches the latest incoming messages via
// the gateway client. No client-side clamping — the gateway is the trust
// boundary and enforces the limit range.
func GetLatestIncomingMessagesDirect(
	client gateway.GatewayClient,
	input GetLatestIncomingMessagesInput,
) (GetLatestIncomingMessagesResult, error) {
	resp, err := client.GetIncomingMessages(context.Background(), input.Limit)
	if err != nil {
		return GetLatestIncomingMessagesResult{}, fmt.Errorf("get_latest_incoming_messages: %w", err)
	}

	return GetLatestIncomingMessagesResult{
		Success:   resp.Success,
		Timestamp: resp.Timestamp,
		Count:     resp.Count,
		Messages:  resp.Messages,
	}, nil
}
