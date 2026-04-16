package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/glennprays/mcp-whatsapp-gateway/internal/config"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/gateway"
	"github.com/glennprays/mcp-whatsapp-gateway/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	log.Println("Starting WhatsApp Gateway MCP server")
	log.Printf("Configuration: env=%s, transport=%s, waga_base_url=%s",
		cfg.AppEnv, cfg.Transport, cfg.WagaBaseURL)

	// Initialize gateway client
	gatewayClient, err := gateway.New(cfg)
	if err != nil {
		log.Printf("Failed to initialize gateway client: %v\n", err)
		fmt.Fprintf(os.Stderr, "Gateway client initialization error: %v\n", err)
		os.Exit(1)
	}

	// Health check (warn on failure, continue)
	ctx := context.Background()
	health, err := gatewayClient.Health(ctx)
	if err != nil {
		log.Printf("Warning: Gateway health check failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Warning: Gateway health check failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Continuing anyway...\n")
	} else {
		log.Printf("Gateway health check passed: status=%s, timestamp=%s", health.Status, health.Timestamp)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start appropriate server based on transport
	if cfg.IsStdioTransport() {
		runStdioServer(ctx, cfg, gatewayClient)
	} else if cfg.IsHTTPTransport() {
		runHTTPServer(ctx, cfg, gatewayClient)
	} else {
		log.Printf("Invalid transport type: %s\n", cfg.Transport)
		fmt.Fprintf(os.Stderr, "Invalid transport type: %s\n", cfg.Transport)
		os.Exit(1)
	}

	log.Println("Server shutdown complete")
}

func runStdioServer(ctx context.Context, cfg *config.Config, gatewayClient gateway.GatewayClient) {
	log.Println("Starting stdio transport server")

	mcpServer, err := server.NewStdioServer(cfg, gatewayClient)
	if err != nil {
		log.Printf("Failed to create stdio server: %v\n", err)
		fmt.Fprintf(os.Stderr, "Failed to create stdio server: %v\n", err)
		os.Exit(1)
	}
	defer mcpServer.Close()

	log.Println("Stdio server ready, waiting for input...")
	if err := mcpServer.RunStdio(ctx); err != nil {
		log.Printf("Server error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func runHTTPServer(ctx context.Context, cfg *config.Config, gatewayClient gateway.GatewayClient) {
	log.Printf("Starting HTTP+SSE transport server on port %s\n", cfg.GetPort())

	if cfg.IsProduction() {
		log.Println("Production mode: Basic authentication enabled")
	}

	mcpServer, err := server.NewHTTPServer(cfg, gatewayClient)
	if err != nil {
		log.Printf("Failed to create HTTP server: %v\n", err)
		fmt.Fprintf(os.Stderr, "Failed to create HTTP server: %v\n", err)
		os.Exit(1)
	}
	defer mcpServer.Close()

	log.Println("HTTP+SSE server ready")
	if err := mcpServer.RunHTTP(ctx); err != nil {
		log.Printf("Server error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
