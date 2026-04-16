# mcp-whatsapp-gateway

A Model Context Protocol (MCP) server that exposes the [WhatsApp Gateway (waga)](https://waga.glennprays.com) as tools for AI agents. This allows Claude and other AI agents to send WhatsApp messages, manage webhooks, and check connection status through a pre-authenticated JWT session.

## Overview

`mcp-whatsapp-gateway` is a Go-based MCP server that provides a clean interface between AI agents and the WhatsApp Gateway. It supports both stdio and HTTP+SSE transports, making it suitable for use with Claude Desktop, Cursor, Claude Code, or any MCP-compatible client.

**Transport Support:**
- **stdio**: For Claude Desktop, Cursor, and Claude Code (default)
- **HTTP+SSE**: For web-based MCP clients with configurable Basic authentication

**Key Features:**
- Send text and image messages to WhatsApp contacts and groups
- Edit, delete, and react to sent messages
- Check WhatsApp connection status
- Manage webhook URLs for incoming message notifications
- Health monitoring for the gateway service
- Pre-authenticated JWT session (no login flow required)
- Comprehensive logging with trace IDs

## Quick Start

### Prerequisites

- Go 1.25 or later
- A running instance of [WhatsApp Gateway](https://github.com/glennprays/whatsapp-gateway)
- A JWT token for your WhatsApp account (obtained from the gateway)

### Installation

#### From Source

```bash
# Clone the repository
git clone https://github.com/glennprays/mcp-whatsapp-gateway.git
cd mcp-whatsapp-gateway

# Install dependencies
go mod download

# Build the binary
go build -o mcp-whatsapp-gateway

# Or install directly
go install github.com/glennprays/mcp-whatsapp-gateway@latest
```

#### Using Docker

```bash
# Build the Docker image
docker build -t mcp-whatsapp-gateway .

# Run the container (example for HTTP+SSE transport)
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="dev" \
  mcp-whatsapp-gateway
```

### Configuration

All configuration is done via environment variables. Create a `.env` file or set them directly:

```bash
# Required Configuration
export WAGA_BASE_URL="http://localhost:3000/api/v1"
export WAGA_JWT_TOKEN="your_jwt_token"

# Optional Configuration (with defaults)
export APP_ENV="dev"                    # dev or prod
export LOG_LEVEL="info"                 # debug, info, warn, error
export MCP_TRANSPORT="stdio"            # stdio or http
export MCP_PORT="8080"                  # HTTP+SSE port (when MCP_TRANSPORT=http)

# Production HTTP+SSE Only
export MCP_BASIC_AUTH_USER="admin"      # Required when APP_ENV=prod and MCP_TRANSPORT=http
export MCP_BASIC_AUTH_PASSWORD="secure_password"  # Required when APP_ENV=prod and MCP_TRANSPORT=http
```

### Usage

#### stdio Transport (Default)

For use with Claude Desktop, Cursor, or Claude Code:

```bash
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

#### HTTP+SSE Transport

For web-based MCP clients:

```bash
# Development (no authentication)
MCP_TRANSPORT="http" \
MCP_PORT="8080" \
APP_ENV="dev" \
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway

# Production (with Basic authentication)
MCP_TRANSPORT="http" \
MCP_PORT="8080" \
APP_ENV="prod" \
MCP_BASIC_AUTH_USER="admin" \
MCP_BASIC_AUTH_PASSWORD="secure_password" \
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

## MCP Tools

The following tools are exposed by this MCP server:

### Messaging Tools

#### send_text_message
Send a text message to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
  - Individual: `{phone}@s.whatsapp.net` (e.g., `6281234567890@s.whatsapp.net`)
  - Group: `{group_id}@g.us` (e.g., `120363xxxxx@g.us`)
- `message` (string, required): Text message content

**Returns:** Message ID and status

#### send_image_message
Send an image message to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `image_url` (string, required): URL of the image to send
- `caption` (string, optional): Image caption
- `view_once` (boolean, optional): Whether the image should be view-once (default: false)

**Returns:** Message ID and status

#### edit_message
Edit a previously sent message.

**Input:**
- `to` (string, required): Recipient address in JID format
- `message_id` (string, required): ID of the message to edit
- `new_message` (string, required): New message content

**Returns:** Success status

#### delete_message
Delete a previously sent message.

**Input:**
- `to` (string, required): Recipient address in JID format
- `message_id` (string, required): ID of the message to delete

**Returns:** Success status

#### react_to_message
React to a message with an emoji.

**Input:**
- `to` (string, required): Recipient address in JID format
- `message_id` (string, required): ID of the message to react to
- `emoji` (string, required): Emoji reaction

**Returns:** Success status

### Connection Tools

#### check_connection_status
Check if the WhatsApp session is active.

**Input:** None

**Returns:** Authentication status and connection information

#### check_health
Check if the WhatsApp Gateway service is reachable.

**Input:** None

**Returns:** Gateway health status

### Webhook Tools

#### get_webhook
Get the currently registered webhook URL.

**Input:** None

**Returns:** Webhook URL and configuration

#### register_webhook
Register a webhook URL for incoming message notifications.

**Input:**
- `url` (string, required): Webhook URL to register
- `hmac_secret` (string, optional): HMAC secret for webhook signature verification

**Returns:** Success status

#### delete_webhook
Remove the registered webhook.

**Input:** None

**Returns:** Success status

## Error Handling

The server provides specific error messages for common scenarios:

- **401 Unauthorized**: "JWT token is invalid or expired. Re-register the phone number against the gateway to obtain a new token."
- **403 Forbidden**: "Session may be disconnected. Run `check_connection_status` to verify."
- **500 Internal Server Error**: Suggestion to check gateway logs

All errors include trace IDs for debugging and monitoring.

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests (requires a running gateway)
go test -tags=integration ./...
```

### Manual Testing

```bash
# Test stdio transport
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="test_token" \
./mcp-whatsapp-gateway

# Test HTTP+SSE transport
MCP_TRANSPORT="http" \
MCP_PORT="8080" \
APP_ENV="dev" \
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="test_token" \
./mcp-whatsapp-gateway

# Test HTTP endpoint
curl http://localhost:8080/mcp
```

## Docker Usage

### Building

```bash
docker build -t mcp-whatsapp-gateway .
```

### Running

```bash
# Development mode (no auth)
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="dev" \
  mcp-whatsapp-gateway

# Production mode (with auth)
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="prod" \
  -e MCP_BASIC_AUTH_USER="admin" \
  -e MCP_BASIC_AUTH_PASSWORD="secure_password" \
  mcp-whatsapp-gateway
```

## Architecture

This project follows Go best practices with an internal package structure:

```
mcp-whatsapp-gateway/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration loading and validation
│   ├── gateway/           # WhatsApp Gateway client initialization
│   ├── server/            # MCP server transports (stdio, HTTP+SSE)
│   └── tools/             # MCP tool implementations
└── Dockerfile             # Multi-stage Docker build
```

### Key Components

- **Config**: Environment-based configuration with validation
- **Gateway Client**: WhatsApp Gateway SDK integration
- **Server**: MCP protocol support with stdio and HTTP+SSE transports
- **Tools**: Message handling, connection management, and webhook operations

## Security Considerations

- **JWT Token**: Store securely and never log at Info level or above
- **Basic Auth**: Required for production HTTP+SSE deployments
- **Environment**: Use `APP_ENV=prod` in production for security features
- **Docker**: Uses distroless base image for minimal attack surface
- **Logging**: Sensitive data is never logged at Info level or above

## Dependencies

- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [whatsapp-gateway-sdk-go](https://github.com/glennprays/whatsapp-gateway-sdk-go) - WhatsApp Gateway Go SDK
- [log](https://github.com/glennprays/log) - Structured logging package

## Related Projects

- [WhatsApp Gateway](https://github.com/glennprays/whatsapp-gateway) - The underlying gateway service
- [WhatsApp Gateway SDK Go](https://github.com/glennprays/whatsapp-gateway-sdk-go) - Go SDK for the gateway
- [Model Context Protocol](https://modelcontextprotocol.io) - The protocol specification

## License

[Specify your license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions:
- GitHub Issues: [Create an issue](https://github.com/glennprays/mcp-whatsapp-gateway/issues)
- WhatsApp Gateway Documentation: https://waga.glennprays.com
