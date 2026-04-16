package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Environment represents the application environment
type Environment string

const (
	// Dev is the development environment
	Dev Environment = "development"
	// Prod is the production environment
	Prod Environment = "production"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}

// Config holds all configuration loaded from environment variables
type Config struct {
	// WhatsApp Gateway Configuration
	WagaBaseURL  string `mapstructure:"WAGA_BASE_URL"`
	WagaJWTToken string `mapstructure:"WAGA_JWT_TOKEN"`

	// Application Configuration
	AppEnv   Environment `mapstructure:"APP_ENV" default:"development"`
	LogLevel string      `mapstructure:"LOG_LEVEL" default:"info"`

	// MCP Server Configuration
	Transport string `mapstructure:"MCP_TRANSPORT" default:"stdio"`
	Port      string `mapstructure:"MCP_PORT" default:"8080"`

	// HTTP+SSE Authentication (prod only)
	BasicAuthUser     string `mapstructure:"MCP_BASIC_AUTH_USER" default:""`
	BasicAuthPassword string `mapstructure:"MCP_BASIC_AUTH_PASSWORD" default:""`
}

// Load reads configuration from environment variables and validates it
func Load() (*Config, error) {
	// Create config instance
	cfg := &Config{}

	// Apply defaults from struct tags
	if err := defaults.Set(cfg); err != nil {
		return nil, fmt.Errorf("failed to set defaults: %w", err)
	}

	// Load .env file if it exists
	_ = godotenv.Load()

	// Get environment from ENV variable
	envStr := strings.ToLower(os.Getenv("APP_ENV"))
	env := Environment(envStr)
	if env == "" {
		env = Dev
	}

	// Configure Viper to read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Auto-bind each struct field by key
	t := reflect.TypeOf(cfg).Elem()
	for i := range t.NumField() {
		field := t.Field(i)
		key := field.Tag.Get("mapstructure")
		if key != "" {
			err := viper.BindEnv(key)
			if err != nil {
				return nil, fmt.Errorf("failed to bind env key %s: %w", key, err)
			}
		}
	}

	// Unmarshal environment variables into config
	// This will override defaults with actual env values
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that all required configuration is present and valid
func (c *Config) validate() error {
	// Check required configuration
	if c.WagaBaseURL == "" {
		return fmt.Errorf("missing required configuration: WAGA_BASE_URL is required")
	}

	if c.WagaJWTToken == "" {
		return fmt.Errorf("missing required configuration: WAGA_JWT_TOKEN is required")
	}

	// Validate APP_ENV
	if !isValidAppEnv(string(c.AppEnv)) {
		return fmt.Errorf("invalid APP_ENV value: %s (must be 'development' or 'production')", c.AppEnv)
	}

	// Validate Transport
	if !isValidTransport(c.Transport) {
		return fmt.Errorf("invalid MCP_TRANSPORT value: %s (must be 'stdio' or 'http')", c.Transport)
	}

	// Validate LogLevel
	if !isValidLogLevel(c.LogLevel) {
		return fmt.Errorf("invalid LOG_LEVEL value: %s (must be 'debug', 'info', 'warn', or 'error')", c.LogLevel)
	}

	// Validate Basic Auth for prod+http
	if c.AppEnv == Prod && c.Transport == "http" {
		if c.BasicAuthUser == "" || c.BasicAuthPassword == "" {
			return fmt.Errorf("missing required configuration: MCP_BASIC_AUTH_USER and MCP_BASIC_AUTH_PASSWORD are required when APP_ENV=production and MCP_TRANSPORT=http")
		}
	}

	return nil
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == Prod
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == Dev
}

// IsStdioTransport returns true if the transport is stdio
func (c *Config) IsStdioTransport() bool {
	return c.Transport == "stdio"
}

// IsHTTPTransport returns true if the transport is HTTP+SSE
func (c *Config) IsHTTPTransport() bool {
	return c.Transport == "http"
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

// isValidAppEnv checks if the app environment value is valid
func isValidAppEnv(env string) bool {
	return env == string(Dev) || env == string(Prod)
}

// isValidTransport checks if the transport value is valid
func isValidTransport(transport string) bool {
	return transport == "stdio" || transport == "http"
}

// isValidLogLevel checks if the log level value is valid
func isValidLogLevel(level string) bool {
	switch strings.ToLower(level) {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}
