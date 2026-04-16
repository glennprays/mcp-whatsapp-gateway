package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// CheckHealthInput represents the input for checking gateway health
type CheckHealthInput struct {
	// No input parameters required
}

// CheckHealthResult represents the result of checking gateway health
type CheckHealthResult struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// CheckHealth checks if the WhatsApp Gateway service is reachable
func CheckHealth(ctx context.Context, req *mcp.CallToolRequest, input CheckHealthInput) (
	*mcp.CallToolResult,
	CheckHealthResult,
	error,
) {
	// Get gateway client from context
	client, ok := ctx.Value("gateway").(gateway.GatewayClient)
	if !ok || client == nil {
		return nil, CheckHealthResult{}, fmt.Errorf("gateway client not available")
	}

	// Call the direct implementation
	result, err := CheckHealthDirect(client, input)
	if err != nil {
		return nil, CheckHealthResult{}, err
	}

	return nil, result, nil
}

// CheckHealthDirect checks if the WhatsApp Gateway service is reachable without using context
func CheckHealthDirect(client gateway.GatewayClient, input CheckHealthInput) (CheckHealthResult, error) {
	// Check health via gateway
	ctx := context.Background()
	health, err := client.Health(ctx)
	if err != nil {
		return CheckHealthResult{}, fmt.Errorf("check_health: %w", err)
	}

	result := CheckHealthResult{
		Status:    health.Status,
		Timestamp: health.Timestamp,
	}

	if health.Status == "ok" {
		result.Message = "Gateway service is reachable and healthy"
	} else {
		result.Message = fmt.Sprintf("Gateway service status: %s", health.Status)
	}

	return result, nil
}
