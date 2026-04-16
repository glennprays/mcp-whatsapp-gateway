package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// CheckConnectionStatusInput represents the input for checking connection status
type CheckConnectionStatusInput struct {
	// No input parameters required
}

// CheckConnectionStatusResult represents the result of checking connection status
type CheckConnectionStatusResult struct {
	Authenticated bool   `json:"authenticated"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// CheckConnectionStatus checks if the WhatsApp session is active and authenticated
func CheckConnectionStatus(ctx context.Context, req *mcp.CallToolRequest, input CheckConnectionStatusInput) (
	*mcp.CallToolResult,
	CheckConnectionStatusResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, CheckConnectionStatusResult{}, fmt.Errorf("gateway client not available")
	}

	// Check login status via gateway
	status, err := client.GetLoginStatus(ctx)
	if err != nil {
		return nil, CheckConnectionStatusResult{}, fmt.Errorf("check_connection_status: %w", err)
	}

	result := CheckConnectionStatusResult{
		Authenticated: status.Authenticated,
	}

	if status.Authenticated {
		result.Status = "connected"
		result.Message = "WhatsApp session is active and authenticated"
	} else {
		result.Status = "disconnected"
		result.Message = "WhatsApp session is not authenticated. Please reconnect to the gateway."
	}

	return nil, result, nil
}
