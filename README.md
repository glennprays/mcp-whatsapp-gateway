# mcp-whatsapp-gateway

A Model Context Protocol (MCP) server that exposes the [WhatsApp Gateway (waga)](https://waga.glennprays.com) as tools for AI agents. This allows Claude and other AI agents to send WhatsApp messages, manage webhooks, and check connection status through a pre-authenticated JWT session.

## What is WhatsApp Gateway?

[WhatsApp Gateway (waga)](https://github.com/glennprays/whatsapp-gateway) is a standalone service that provides a REST API for interacting with WhatsApp. It handles:

- **WhatsApp Integration**: Manages the connection to WhatsApp servers
- **Message Handling**: Sends/receives text and image messages
- **Session Management**: Maintains authenticated WhatsApp sessions
- **Webhook Support**: Delivers incoming messages via webhooks
- **Health Monitoring**: Provides status and health endpoints

## How This MCP Server Connects

This MCP server acts as a **bridge between AI agents and the WhatsApp Gateway**:

```
┌─────────────────┐     MCP Protocol      ┌──────────────────────┐     HTTP/JWT     ┌──────────────────┐
│  AI Agent       │ ←────────────────────→ │  MCP WhatsApp        │ ←────────────────→ │  WhatsApp        │
│  (Claude/Cursor │   (stdio or HTTP+SSE)   │  Gateway Server      │   (REST API)      │  Gateway         │
│   Claude Code)  │                         │                      │                   │  (waga)          │
└─────────────────┘                         └──────────────────────┘                   └──────────────────┘
```

**Data Flow:**
1. **AI Agent** calls MCP tools (e.g., `send_text_message`)
2. **MCP Server** receives the tool call and validates input
3. **Gateway Client** makes HTTP request to WhatsApp Gateway with JWT authentication
4. **WhatsApp Gateway** processes the request and interacts with WhatsApp
5. **Response** flows back through the chain to the AI agent

**Key Points:**
- The WhatsApp Gateway runs as a **separate service** that you need to deploy
- This MCP server connects to it via **HTTP using JWT authentication**
- The gateway handles all WhatsApp-specific logic and protocol details
- This MCP server only provides the MCP protocol interface for AI agents

## Prerequisites

Before using this MCP server, you need:

1. **Running WhatsApp Gateway instance**
   - Follow the setup guide at: https://waga.glennprays.com
   - Deploy the gateway service (Docker, binary, or cloud)
   - Ensure it's accessible via HTTP/HTTPS

2. **JWT Token from Gateway**
   - Register your phone number with the gateway
   - Obtain a JWT token for authentication
   - This token is used by the MCP server to authenticate with the gateway

3. **Gateway Configuration**
   - Set `WAGA_BASE_URL` to point to your gateway instance
   - Example: `http://localhost:3000/api/v1` (local development)
   - Example: `https://waga.example.com/api/v1` (production)

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

### 1. Pull the Docker Image

```bash
docker pull glennprays/mcp-whatsapp-gateway:latest
```

### 2. Run with Docker

**For Claude Desktop/Cursor/Claude Code (stdio transport):**
```bash
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

**For web-based clients (HTTP+SSE transport):**
```bash
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

### 3. Configure Your MCP Client

See [Configuring MCP Clients](#configuring-mcp-clients) below for detailed instructions.

---

## Installation

### Using Docker (Recommended)

Pull the pre-built image from Docker Hub:

```bash
docker pull glennprays/mcp-whatsapp-gateway:latest
```

### From Source (Development Only)

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

### Building Docker Image (Development)

```bash
# Build the Docker image
docker build -t mcp-whatsapp-gateway .

# Run the container
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  mcp-whatsapp-gateway
```

## Configuration

All configuration is done via environment variables:

### Required Variables

```bash
export WAGA_BASE_URL="http://localhost:3000/api/v1"  # Your gateway URL
export WAGA_JWT_TOKEN="your_jwt_token"              # From gateway registration
```

### Optional Variables

```bash
# Application Settings
export APP_ENV="development"         # development or production
export LOG_LEVEL="info"              # debug, info, warn, error

# Transport Settings
export MCP_TRANSPORT="stdio"         # stdio or http
export MCP_PORT="8080"               # HTTP+SSE port (when MCP_TRANSPORT=http)

# Production HTTP+SSE Only
export MCP_BASIC_AUTH_USER="admin"   # Required when APP_ENV=production and MCP_TRANSPORT=http
export MCP_BASIC_AUTH_PASSWORD="secure_password"  # Required when APP_ENV=production and MCP_TRANSPORT=http
```

### Running the Server

#### stdio Transport (Default)

For use with Claude Desktop, Cursor, or Claude Code:

**Using Docker:**
```bash
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

**Using binary:**
```bash
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

#### HTTP+SSE Transport

For web-based MCP clients:

**Development (no authentication):**
```bash
# Docker
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest

# Binary
MCP_TRANSPORT="http" MCP_PORT="8080" \
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

**Production (with authentication):**
```bash
# Docker
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="production" \
  -e MCP_BASIC_AUTH_USER="admin" \
  -e MCP_BASIC_AUTH_PASSWORD="secure_password" \
  -e WAGA_BASE_URL="https://your-gateway.com/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  --restart unless-stopped \
  glennprays/mcp-whatsapp-gateway:latest

# Binary
MCP_TRANSPORT="http" MCP_PORT="8080" \
APP_ENV="production" \
MCP_BASIC_AUTH_USER="admin" \
MCP_BASIC_AUTH_PASSWORD="secure_password" \
WAGA_BASE_URL="https://your-gateway.com/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

### Configuring MCP Clients

#### Claude Desktop

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%/Claude/claude_desktop_config.json`
**Linux:** `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "WAGA_BASE_URL=http://host.docker.internal:3000/api/v1",
        "-e", "WAGA_JWT_TOKEN=your_jwt_token",
        "glennprays/mcp-whatsapp-gateway:latest"
      ]
    }
  }
}
```

#### Cursor IDE

```json
{
  "mcpServers": [
    {
      "name": "whatsapp-gateway",
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "WAGA_BASE_URL=http://host.docker.internal:3000/api/v1",
        "-e", "WAGA_JWT_TOKEN=your_jwt_token",
        "glennprays/mcp-whatsapp-gateway:latest"
      ]
    }
  ]
}
```

**For comprehensive Docker deployment guide, see [DOCKER_USAGE.md](DOCKER_USAGE.md)**

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
APP_ENV="development" \
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="test_token" \
./mcp-whatsapp-gateway

# Test HTTP endpoint
curl http://localhost:8080/mcp
```

## Docker Usage

### Using Pre-built Images (Recommended)

```bash
# Pull the latest image
docker pull glennprays/mcp-whatsapp-gateway:latest

# Run in development mode (no auth)
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="development" \
  glennprays/mcp-whatsapp-gateway:latest

# Run in production mode (with auth)
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e WAGA_BASE_URL="https://your-gateway.com/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="production" \
  -e MCP_BASIC_AUTH_USER="admin" \
  -e MCP_BASIC_AUTH_PASSWORD="your_secure_password" \
  --restart unless-stopped \
  glennprays/mcp-whatsapp-gateway:latest
```

### Building from Source (Development Only)

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
  -e APP_ENV="development" \
  mcp-whatsapp-gateway

# Production mode (with auth)
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e MCP_TRANSPORT="http" \
  -e APP_ENV="production" \
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

### Relationship to WhatsApp Gateway

This MCP server is **designed to work exclusively with the WhatsApp Gateway**:

- **Separate Services**: The MCP server and WhatsApp Gateway are separate services
- **HTTP Communication**: MCP server communicates with gateway via HTTP/HTTPS
- **JWT Authentication**: Uses JWT tokens obtained from the gateway for authentication
- **API Endpoints**: Maps MCP tools to gateway REST API endpoints
  - `send_text_message` → POST /message/text
  - `check_health` → GET /health
  - `get_webhook` → GET /webhook
  - etc.

**Why This Architecture?**
- **Separation of Concerns**: WhatsApp Gateway handles WhatsApp protocol, MCP server handles AI agent integration
- **Scalability**: Each service can be deployed and scaled independently
- **Flexibility**: Multiple MCP servers can connect to the same gateway
- **Security**: JWT tokens provide secure authentication without exposing WhatsApp credentials

## Troubleshooting

### Gateway Connection Issues

**Problem**: "Failed to initialize gateway client" or "Gateway health check failed"

**Solutions**:
1. **Verify Gateway is Running**: Ensure the WhatsApp Gateway service is running
   ```bash
   curl http://localhost:3000/api/v1/health
   ```

2. **Check WAGA_BASE_URL**: Verify the URL is correct and accessible
   - Local: `http://localhost:3000/api/v1`
   - Docker: Use `host.docker.internal` on Mac/Windows
   - Remote: Ensure firewall allows connections

3. **Validate JWT Token**: Ensure your JWT token is valid and not expired
   ```bash
   # Test gateway connection with curl
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
        http://localhost:3000/api/v1/health
   ```

4. **Check Gateway Logs**: Look for errors in the WhatsApp Gateway logs

**Problem**: "401 Unauthorized" when sending messages

**Solutions**:
1. **JWT Token Expired**: Re-register your phone number with the gateway to get a new token
2. **Invalid Token**: Verify the token is correctly set in `WAGA_JWT_TOKEN`

**Problem**: "403 Forbidden" when sending messages

**Solutions**:
1. **Session Disconnected**: Your WhatsApp session may have disconnected
2. **Run Connection Check**: Use `check_connection_status` tool to verify
3. **Reconnect Gateway**: Restart the WhatsApp Gateway and re-authenticate

### MCP Server Issues

**Problem**: "Invalid transport type"

**Solution**: Set `MCP_TRANSPORT` to either `stdio` or `http`

**Problem**: HTTP+SSE returns 401 in production

**Solution**: Set `MCP_BASIC_AUTH_USER` and `MCP_BASIC_AUTH_PASSWORD` when `APP_ENV=production`

### Docker Networking

**Problem**: Container cannot reach host services

**Solution**: Use the appropriate host address:
- macOS/Windows: `host.docker.internal`
- Linux: Use host's IP address or `172.17.0.1` (default Docker bridge)

```bash
# Example for macOS/Windows
docker run -p 8080:8080 \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_token" \
  mcp-whatsapp-gateway
```

## Security Considerations

- **JWT Token**: Store securely and never log at Info level or above
- **Basic Auth**: Required for production HTTP+SSE deployments
- **Environment**: Use `APP_ENV=production` in production for security features
- **Docker**: Uses distroless base image for minimal attack surface
- **Logging**: Sensitive data is never logged at Info level or above

## Getting Started with WhatsApp Gateway

This MCP server requires a running WhatsApp Gateway instance. For detailed setup instructions, deployment options, and configuration guidance, visit:

**WhatsApp Gateway Documentation:** https://waga.glennprays.com

Once your gateway is running and you have a JWT token, configure this MCP server:

```bash
export WAGA_BASE_URL="http://localhost:3000/api/v1"  # Your gateway URL
export WAGA_JWT_TOKEN="your_jwt_token_here"          # From gateway registration
```

Then run the MCP server using Docker or from source (see Installation section above).

## Dependencies

- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official MCP Go SDK
- [whatsapp-gateway-sdk-go](https://github.com/glennprays/whatsapp-gateway-sdk-go) - WhatsApp Gateway Go SDK
- [log](https://github.com/glennprays/log) - Structured logging package

## Related Projects

- [WhatsApp Gateway](https://github.com/glennprays/whatsapp-gateway) - The underlying gateway service
- [WhatsApp Gateway SDK Go](https://github.com/glennprays/whatsapp-gateway-sdk-go) - Go SDK for the gateway
- [Model Context Protocol](https://modelcontextprotocol.io) - The protocol specification

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**Copyright (c) 2024 Glenn Prays**

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, USE OR OTHER DEALINGS IN THE SOFTWARE.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions:
- GitHub Issues: [Create an issue](https://github.com/glennprays/mcp-whatsapp-gateway/issues)
- WhatsApp Gateway Documentation: https://waga.glennprays.com
