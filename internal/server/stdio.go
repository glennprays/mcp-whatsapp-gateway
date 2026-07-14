package server

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/tools"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServer wraps the MCP server instance
type MCPServer struct {
	server *mcp.Server
}

// NewStdioServer creates a new MCP server with stdio transport
func NewStdioServer(cfg *config.Config, gatewayClient gateway.GatewayClient) (*MCPServer, error) {
	// Verify gateway client is not nil
	if gatewayClient == nil {
		return nil, fmt.Errorf("gatewayClient is nil")
	}

	// Create MCP server instance
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "whatsapp-gateway",
		Version: "1.0.0",
	}, nil)

	// Register messaging tools using the SDK's typed handler system
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_text_message",
		Description: "Send a text message to a WhatsApp contact or group",
	}, createSendTextMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_image_message",
		Description: "Send an image message to a WhatsApp contact or group",
	}, createSendImageMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_audio_message",
		Description: "Send an audio message to a WhatsApp contact or group",
	}, createSendAudioMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_video_message",
		Description: "Send a video message to a WhatsApp contact or group",
	}, createSendVideoMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_document_message",
		Description: "Send a document message to a WhatsApp contact or group",
	}, createSendDocumentMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit_message",
		Description: "Edit a previously sent message",
	}, createEditMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_message",
		Description: "Delete a previously sent message",
	}, createDeleteMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "react_to_message",
		Description: "React to a message with an emoji",
	}, createReactToMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_location_message",
		Description: "Send a location with GPS coordinates to a WhatsApp number or group",
	}, createSendLocationMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_poll_message",
		Description: "Send a poll with question and options to a WhatsApp number or group",
	}, createSendPollMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_sticker_message",
		Description: "Send a sticker from a URL to a WhatsApp number or group",
	}, createSendStickerMessageHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_job_status",
		Description: "Poll the status of a queued message job by its job_id (returned by send_* tools when the gateway runs in queue mode). Returns status (queued/processing/completed/failed), the message_id once completed, and any error.",
	}, createGetJobStatusHandler(gatewayClient))

	// Register inbox tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_latest_incoming_messages",
		Description: "Fetch the most recent incoming WhatsApp messages. Useful for reading OTPs, verification codes, or recent conversation context. Returns newest first (text, image, video, audio, document, sticker, location, poll). Optional limit (default 10, max 50).",
	}, createGetLatestIncomingMessagesHandler(gatewayClient))

	// Register contact & group read tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_contacts",
		Description: "List the account's locally-synced WhatsApp contacts (paginated via limit/offset). Reads the local address book; an empty result is normal, never an error.",
	}, createListContactsHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_contact_info",
		Description: "Look up one contact's WhatsApp profile by chat (a number, user JID, or @lid): status text, current picture id, verified business name, linked-device count, and lid.",
	}, createGetContactInfoHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_avatar",
		Description: "Get a chat's profile picture URL (user or @g.us group). Set preview for the low-res thumbnail. Returns available=false with reason 'not_set' (no picture) or 'hidden' (privacy) instead of failing.",
	}, createGetAvatarHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_groups",
		Description: "List the account's joined WhatsApp groups as lightweight summaries (no participant roster; use get_group_info for one group's full detail).",
	}, createListGroupsHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_group_info",
		Description: "Get one group's full detail plus its participant roster by chat (a @g.us group JID). The account must be a member (403 if not, 404 if absent).",
	}, createGetGroupInfoHandler(gatewayClient))

	// Register connection tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_connection_status",
		Description: "Check if the WhatsApp session is active and authenticated",
	}, createCheckConnectionStatusHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_contact",
		Description: "Check whether a phone number is registered on WhatsApp",
	}, createCheckContactHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_health",
		Description: "Check if the WhatsApp Gateway service is reachable",
	}, createCheckHealthHandler(gatewayClient))

	// Register webhook tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_webhook",
		Description: "Get the current webhook configuration",
	}, createGetWebhookHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "register_webhook",
		Description: "Register a webhook URL to receive WhatsApp messages",
	}, createRegisterWebhookHandler(gatewayClient))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_webhook",
		Description: "Delete the currently registered webhook",
	}, createDeleteWebhookHandler(gatewayClient))

	return &MCPServer{server: server}, nil
}

// Handler creation functions - these bridge our existing tool handlers to MCP SDK format
// The gateway client is captured directly in the closure

func createSendTextMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendMessageInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendTextMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendImageMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendImageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendImageInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendImageMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendAudioMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendAudioInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendAudioInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendAudioMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendVideoMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendVideoInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendVideoInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendVideoMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendDocumentMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendDocumentInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendDocumentInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendDocumentMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createEditMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.EditMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.EditMessageInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.EditMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.DeleteMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteMessageInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.DeleteMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createReactToMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.ReactToMessageInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.ReactToMessageInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.ReactToMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendLocationMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendLocationInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendLocationInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendLocationMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendPollMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendPollInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendPollInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendPollMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createSendStickerMessageHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.SendStickerInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.SendStickerInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.SendStickerMessageDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetJobStatusHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetJobStatusInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetJobStatusInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetJobStatusDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetLatestIncomingMessagesHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetLatestIncomingMessagesInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetLatestIncomingMessagesInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetLatestIncomingMessagesDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createListContactsHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.ListContactsInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.ListContactsInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.ListContactsDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetContactInfoHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetContactInfoInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetContactInfoInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetContactInfoDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetAvatarHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetAvatarInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetAvatarInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetAvatarDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createListGroupsHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.ListGroupsInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.ListGroupsInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.ListGroupsDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetGroupInfoHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetGroupInfoInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetGroupInfoInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetGroupInfoDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckConnectionStatusHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.CheckConnectionStatusInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckConnectionStatusInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.CheckConnectionStatusDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckContactHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.CheckContactInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckContactInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.CheckContactDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createCheckHealthHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.CheckHealthInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.CheckHealthInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.CheckHealthDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createGetWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.GetWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetWebhookInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.GetWebhookDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createRegisterWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.RegisterWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.RegisterWebhookInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.RegisterWebhookDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

func createDeleteWebhookHandler(gatewayClient gateway.GatewayClient) mcp.ToolHandlerFor[tools.DeleteWebhookInput, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteWebhookInput) (*mcp.CallToolResult, any, error) {
		result, err := tools.DeleteWebhookDirect(gatewayClient, input)
		if err != nil {
			return nil, nil, err
		}
		return nil, result, nil
	}
}

// RunStdio starts the server with stdio transport
func (s *MCPServer) RunStdio(ctx context.Context) error {
	// Create stdio transport
	transport := &mcp.StdioTransport{}

	// Run the server
	if err := s.server.Run(ctx, transport); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

// Close shuts down the server
func (s *MCPServer) Close() error {
	// The MCP SDK handles cleanup automatically
	return nil
}
