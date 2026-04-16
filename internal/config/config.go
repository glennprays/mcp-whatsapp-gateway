package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration loaded from environment variables
type Config struct {
	// WhatsApp Gateway Configuration
	WagaBaseURL  string
	WagaJWTToken string

	// Application Configuration
	AppEnv   string
	LogLevel string

	// MCP Server Configuration
	Transport string
	Port      string

	// HTTP+SSE Authentication (prod only)
	BasicAuthUser     string
	BasicAuthPassword string
}

const (
	// Environment Constants
	defaultAppEnv   = "dev"
	defaultLogLevel = "info"
	defaultTransport = "stdio"
	defaultPort      = "8080"

	// Valid Environment Values
	envDev  = "dev"
	envProd = "prod"

	// Valid Transport Values
	transportStdio = "stdio"
	transportHTTP  = "http"

	// Valid Log Levels
	logDebug = "debug"
	logInfo  = "info"
	logWarn  = "warn"
	logError = "error"
)

var (
	// ErrMissingRequiredConfig is returned when required configuration is missing
	ErrMissingRequiredConfig = errors.New("missing required configuration")

	// ErrInvalidAppEnv is returned when APP_ENV has an invalid value
	ErrInvalidAppEnv = errors.New("invalid APP_ENV value")

	// ErrInvalidTransport is returned when MCP_TRANSPORT has an invalid value
	ErrInvalidTransport = errors.New("invalid MCP_TRANSPORT value")

	// ErrInvalidLogLevel is returned when LOG_LEVEL has an invalid value
	ErrInvalidLogLevel = errors.New("invalid LOG_LEVEL value")

	// ErrMissingBasicAuth is returned when Basic Auth credentials are missing in prod+http mode
	ErrMissingBasicAuth = errors.New("missing Basic Auth credentials for production HTTP+SSE mode")
)

// Load reads configuration from environment variables and validates it
func Load() (*Config, error) {
	cfg := &Config{
		WagaBaseURL:        os.Getenv("WAGA_BASE_URL"),
		WagaJWTToken:       os.Getenv("WAGA_JWT_TOKEN"),
		AppEnv:            getEnvWithDefault("APP_ENV", defaultAppEnv),
		LogLevel:          getEnvWithDefault("LOG_LEVEL", defaultLogLevel),
		Transport:         getEnvWithDefault("MCP_TRANSPORT", defaultTransport),
		Port:              getEnvWithDefault("MCP_PORT", defaultPort),
		BasicAuthUser:     os.Getenv("MCP_BASIC_AUTH_USER"),
		BasicAuthPassword: os.Getenv("MCP_BASIC_AUTH_PASSWORD"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that all required configuration is present and valid
func (c *Config) validate() error {
	// Check required configuration
	if c.WagaBaseURL == "" {
		return fmt.Errorf("%w: WAGA_BASE_URL is required", ErrMissingRequiredConfig)
	}

	if c.WagaJWTToken == "" {
		return fmt.Errorf("%w: WAGA_JWT_TOKEN is required", ErrMissingRequiredConfig)
	}

	// Validate APP_ENV
	if !isValidAppEnv(c.AppEnv) {
		return fmt.Errorf("%w: %s (must be 'dev' or 'prod')", ErrInvalidAppEnv, c.AppEnv)
	}

	// Validate Transport
	if !isValidTransport(c.Transport) {
		return fmt.Errorf("%w: %s (must be 'stdio' or 'http')", ErrInvalidTransport, c.Transport)
	}

	// Validate LogLevel
	if !isValidLogLevel(c.LogLevel) {
		return fmt.Errorf("%w: %s (must be 'debug', 'info', 'warn', or 'error')", ErrInvalidLogLevel, c.LogLevel)
	}

	// Validate Basic Auth for prod+http
	if c.AppEnv == envProd && c.Transport == transportHTTP {
		if c.BasicAuthUser == "" || c.BasicAuthPassword == "" {
			return fmt.Errorf("%w: MCP_BASIC_AUTH_USER and MCP_BASIC_AUTH_PASSWORD are required when APP_ENV=prod and MCP_TRANSPORT=http", ErrMissingBasicAuth)
		}
	}

	return nil
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == envProd
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == envDev
}

// IsStdioTransport returns true if the transport is stdio
func (c *Config) IsStdioTransport() bool {
	return c.Transport == transportStdio
}

// IsHTTPTransport returns true if the transport is HTTP+SSE
func (c *Config) IsHTTPTransport() bool {
	return c.Transport == transportHTTP
}

// GetLogLevel returns the log level as a string
func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

// GetPort returns the port for HTTP+SSE server
func (c *Config) GetPort() string {
	return c.Port
}

// Helper functions

// getEnvWithDefault returns the environment variable value or a default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isValidAppEnv checks if the app environment value is valid
func isValidAppEnv(env string) bool {
	return env == envDev || env == envProd
}

// isValidTransport checks if the transport value is valid
func isValidTransport(transport string) bool {
	return transport == transportStdio || transport == transportHTTP
}

// isValidLogLevel checks if the log level value is valid
func isValidLogLevel(level string) bool {
	switch strings.ToLower(level) {
	case logDebug, logInfo, logWarn, logError:
		return true
	default:
		return false
	}
}
