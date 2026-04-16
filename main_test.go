package main

import (
	"os"
	"testing"
)

func TestMain_CanSetEnvironment(t *testing.T) {
	// Test that we can set required environment variables
	requiredVars := []struct {
		key   string
		value string
	}{
		{"WAGA_BASE_URL", "http://localhost:3000/api/v1"},
		{"WAGA_JWT_TOKEN", "test-token"},
		{"APP_ENV", "development"},
		{"MCP_TRANSPORT", "stdio"},
	}

	for _, v := range requiredVars {
		if err := os.Setenv(v.key, v.value); err != nil {
			t.Errorf("Failed to set environment variable %s: %v", v.key, err)
		}
		defer os.Unsetenv(v.key)
	}
}
