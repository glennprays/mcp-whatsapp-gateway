package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	waga "github.com/glennprays/whatsapp-gateway-sdk-go"
)

// MarkReadInput is the request payload for the mark_read tool.
type MarkReadInput struct {
	// Chat is the canonical recipient (a bare number, a user JID, or a "@g.us" group).
	Chat string `json:"chat"`
	// MessageIDs are the message IDs to mark read.
	MessageIDs []string `json:"message_ids"`
	// Sender is the message author's JID/number; required for group chats.
	Sender string `json:"sender,omitempty"`
}

// MarkReadResult is the result of a mark_read call.
type MarkReadResult struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
}

// MarkReadDirect marks one or more messages in a chat as read (blue ticks).
func MarkReadDirect(client gateway.GatewayClient, input MarkReadInput) (MarkReadResult, error) {
	if input.Chat == "" {
		return MarkReadResult{}, fmt.Errorf("chat is required")
	}
	if len(input.MessageIDs) == 0 {
		return MarkReadResult{}, fmt.Errorf("at least one message_id is required")
	}
	if err := client.MarkRead(context.Background(), input.Chat, input.MessageIDs, input.Sender); err != nil {
		return MarkReadResult{}, fmt.Errorf("mark_read: %w", err)
	}
	return MarkReadResult{Success: true, Status: "read"}, nil
}

// SendTypingInput is the request payload for the send_typing tool.
type SendTypingInput struct {
	// Chat is the canonical recipient (a bare number, a user JID, or a "@g.us" group).
	Chat string `json:"chat"`
	// State is one of composing (typing…), recording (recording audio…), or paused (cleared).
	State string `json:"state"`
}

// SendTypingResult is the result of a send_typing call.
type SendTypingResult struct {
	Success bool   `json:"success"`
	State   string `json:"state"`
}

// SendTypingDirect sets the typing indicator in a chat.
func SendTypingDirect(client gateway.GatewayClient, input SendTypingInput) (SendTypingResult, error) {
	if input.Chat == "" {
		return SendTypingResult{}, fmt.Errorf("chat is required")
	}
	switch input.State {
	case waga.PresenceComposing, waga.PresenceRecording, waga.PresencePaused:
	default:
		return SendTypingResult{}, fmt.Errorf("state must be one of composing, recording, paused")
	}
	if err := client.SendChatPresence(context.Background(), input.Chat, input.State); err != nil {
		return SendTypingResult{}, fmt.Errorf("send_typing: %w", err)
	}
	return SendTypingResult{Success: true, State: input.State}, nil
}
