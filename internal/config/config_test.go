package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Set required environment variables
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-jwt-token")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.WagaBaseURL != "http://localhost:3000/api/v1" {
		t.Errorf("Expected WagaBaseURL 'http://localhost:3000/api/v1', got '%s'", cfg.WagaBaseURL)
	}

	if cfg.WagaJWTToken != "test-jwt-token" {
		t.Errorf("Expected WagaJWTToken 'test-jwt-token', got '%s'", cfg.WagaJWTToken)
	}

	// Check defaults
	if cfg.AppEnv != defaultAppEnv {
		t.Errorf("Expected default AppEnv '%s', got '%s'", defaultAppEnv, cfg.AppEnv)
	}

	if cfg.LogLevel != defaultLogLevel {
		t.Errorf("Expected default LogLevel '%s', got '%s'", defaultLogLevel, cfg.LogLevel)
	}

	if cfg.Transport != defaultTransport {
		t.Errorf("Expected default Transport '%s', got '%s'", defaultTransport, cfg.Transport)
	}

	if cfg.Port != defaultPort {
		t.Errorf("Expected default Port '%s', got '%s'", defaultPort, cfg.Port)
	}
}

func TestLoad_MissingWagaBaseURL(t *testing.T) {
	os.Unsetenv("WAGA_BASE_URL")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when WAGA_BASE_URL is missing, got nil")
	}

	if !containsString(err.Error(), "WAGA_BASE_URL is required") {
		t.Errorf("Expected error about missing WAGA_BASE_URL, got: %v", err)
	}
}

func TestLoad_MissingJWTToken(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Unsetenv("WAGA_JWT_TOKEN")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when WAGA_JWT_TOKEN is missing, got nil")
	}

	if !containsString(err.Error(), "WAGA_JWT_TOKEN is required") {
		t.Errorf("Expected error about missing WAGA_JWT_TOKEN, got: %v", err)
	}
}

func TestLoad_InvalidAppEnv(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	os.Setenv("APP_ENV", "invalid")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when APP_ENV is invalid, got nil")
	}

	if !containsString(err.Error(), "invalid APP_ENV value") {
		t.Errorf("Expected error about invalid APP_ENV, got: %v", err)
	}
}

func TestLoad_InvalidTransport(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	os.Setenv("MCP_TRANSPORT", "invalid")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when MCP_TRANSPORT is invalid, got nil")
	}

	if !containsString(err.Error(), "invalid MCP_TRANSPORT value") {
		t.Errorf("Expected error about invalid MCP_TRANSPORT, got: %v", err)
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	os.Setenv("LOG_LEVEL", "invalid")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when LOG_LEVEL is invalid, got nil")
	}

	if !containsString(err.Error(), "invalid LOG_LEVEL value") {
		t.Errorf("Expected error about invalid LOG_LEVEL, got: %v", err)
	}
}

func TestLoad_ProdHTTP_MissingBasicAuth(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	os.Setenv("APP_ENV", "prod")
	os.Setenv("MCP_TRANSPORT", "http")
	os.Unsetenv("MCP_BASIC_AUTH_USER")
	os.Unsetenv("MCP_BASIC_AUTH_PASSWORD")
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error when Basic Auth is missing in prod+http mode, got nil")
	}

	if !containsString(err.Error(), "missing Basic Auth credentials") {
		t.Errorf("Expected error about missing Basic Auth, got: %v", err)
	}
}

func TestLoad_ProdHTTP_WithBasicAuth(t *testing.T) {
	os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
	os.Setenv("WAGA_JWT_TOKEN", "test-token")
	os.Setenv("APP_ENV", "prod")
	os.Setenv("MCP_TRANSPORT", "http")
	os.Setenv("MCP_BASIC_AUTH_USER", "admin")
	os.Setenv("MCP_BASIC_AUTH_PASSWORD", "secret")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.BasicAuthUser != "admin" {
		t.Errorf("Expected BasicAuthUser 'admin', got '%s'", cfg.BasicAuthUser)
	}

	if cfg.BasicAuthPassword != "secret" {
		t.Errorf("Expected BasicAuthPassword 'secret', got '%s'", cfg.BasicAuthPassword)
	}
}

func TestLoad_AllValidLogLevels(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}

	for _, level := range validLevels {
		t.Run(level, func(t *testing.T) {
			os.Setenv("WAGA_BASE_URL", "http://localhost:3000/api/v1")
			os.Setenv("WAGA_JWT_TOKEN", "test-token")
			os.Setenv("LOG_LEVEL", level)
			defer cleanupEnv()

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() failed for log level %s: %v", level, err)
			}

			if cfg.LogLevel != level {
				t.Errorf("Expected LogLevel '%s', got '%s'", level, cfg.LogLevel)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name     string
		appEnv   string
		expected bool
	}{
		{"production", "prod", true},
		{"development", "dev", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{AppEnv: tt.appEnv}
			if got := cfg.IsProduction(); got != tt.expected {
				t.Errorf("IsProduction() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name     string
		appEnv   string
		expected bool
	}{
		{"development", "dev", true},
		{"production", "prod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{AppEnv: tt.appEnv}
			if got := cfg.IsDevelopment(); got != tt.expected {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_TransportChecks(t *testing.T) {
	tests := []struct {
		name         string
		transport    string
		isStdio      bool
		isHTTP       bool
	}{
		{"stdio transport", "stdio", true, false},
		{"http transport", "http", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Transport: tt.transport}
			if got := cfg.IsStdioTransport(); got != tt.isStdio {
				t.Errorf("IsStdioTransport() = %v, want %v", got, tt.isStdio)
			}
			if got := cfg.IsHTTPTransport(); got != tt.isHTTP {
				t.Errorf("IsHTTPTransport() = %v, want %v", got, tt.isHTTP)
			}
		})
	}
}

func TestConfig_Getters(t *testing.T) {
	cfg := &Config{
		LogLevel: "debug",
		Port:     "9090",
	}

	if got := cfg.GetLogLevel(); got != "debug" {
		t.Errorf("GetLogLevel() = %v, want debug", got)
	}

	if got := cfg.GetPort(); got != "9090" {
		t.Errorf("GetPort() = %v, want 9090", got)
	}
}

// Helper functions

func cleanupEnv() {
	os.Unsetenv("WAGA_BASE_URL")
	os.Unsetenv("WAGA_JWT_TOKEN")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("MCP_TRANSPORT")
	os.Unsetenv("MCP_PORT")
	os.Unsetenv("MCP_BASIC_AUTH_USER")
	os.Unsetenv("MCP_BASIC_AUTH_PASSWORD")
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsInString(s, substr))
}

func containsInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
