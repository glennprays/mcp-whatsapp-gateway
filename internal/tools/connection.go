package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
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

// CheckConnectionStatusDirect checks if the WhatsApp session is active and authenticated
func CheckConnectionStatusDirect(client gateway.GatewayClient, input CheckConnectionStatusInput) (CheckConnectionStatusResult, error) {
	ctx := context.Background()
	status, err := client.GetLoginStatus(ctx)
	if err != nil {
		return CheckConnectionStatusResult{}, fmt.Errorf("check_connection_status: %w", err)
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

	return result, nil
}
