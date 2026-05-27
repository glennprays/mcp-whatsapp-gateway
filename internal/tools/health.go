package tools

import (
	"context"
	"fmt"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
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

// CheckHealthDirect checks if the WhatsApp Gateway service is reachable
func CheckHealthDirect(client gateway.GatewayClient, input CheckHealthInput) (CheckHealthResult, error) {
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
