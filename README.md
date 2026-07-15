# mcp-whatsapp-gateway

A Model Context Protocol (MCP) server that exposes the [WhatsApp Gateway (waga)](https://waga.glennprays.com) as tools for AI agents. This allows Claude and other AI agents to send WhatsApp messages, manage webhooks, and check connection status through a pre-authenticated JWT session.

** Documentation:**
- **[WhatsApp Gateway MCP Documentation](https://waga.glennprays.com/#docs:mcp/introduction)** - Main comprehensive documentation
- **[DOCKER_USAGE.md](DOCKER_USAGE.md)** - Docker deployment guide
- **[README.md](#quick-start)** - Quick start and overview (this file)

## What is WhatsApp Gateway?

[WhatsApp Gateway (waga)](https://github.com/glennprays/whatsapp-gateway) is a standalone service that provides a REST API for interacting with WhatsApp. It handles:

- **WhatsApp Integration**: Manages the connection to WhatsApp servers
- **Message Handling**: Sends/receives text, image, location, poll, and sticker messages
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

> ** For detailed setup instructions, see the [comprehensive documentation](https://waga.glennprays.com/#docs:mcp/introduction)**

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
- Send text, image, audio, video, document, location, poll, and sticker messages to WhatsApp contacts and groups
- Edit, delete, and react to sent messages
- Check whether a number is registered on WhatsApp
- Check WhatsApp connection status
- Manage webhook URLs for incoming message notifications
- Health monitoring for the gateway service
- Pre-authenticated JWT session (no login flow required)
- Comprehensive logging with trace IDs

## Quick Start

> ** For comprehensive setup and configuration guide, see the [main documentation](https://waga.glennprays.com/#docs:mcp/introduction)**

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

> ** For comprehensive Docker deployment guide, transport options, and production configuration, see [DOCKER_USAGE.md](DOCKER_USAGE.md)**

### 3. Configure Your MCP Client

See [Configuring MCP Clients](#configuring-mcp-clients) below for detailed instructions.

---

## Installation

### Using Docker (Recommended)

Pull the pre-built image from Docker Hub:

```bash
docker pull glennprays/mcp-whatsapp-gateway:latest
```

> ** For detailed Docker setup, deployment options, and production configuration, see [DOCKER_USAGE.md](DOCKER_USAGE.md)**

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

> ** For Docker deployment options, transport configuration, and production setup, see [DOCKER_USAGE.md](DOCKER_USAGE.md)**

#### Claude Code CLI

For Claude Code CLI, configure MCP servers using the command line:

```bash
# Add the MCP server with environment variables
claude mcp add -e WAGA_BASE_URL=http://host.docker.internal:3000/api/v1 \
  -e WAGA_JWT_TOKEN=your_jwt_token \
  whatsapp-gateway -- /path/to/mcp-whatsapp-gateway

# Or use a project-specific .mcp.json file
# See .mcp.json.example in this repository
```

## MCP Tools

The following tools are exposed by this MCP server:

> **Curated subset — read & messaging only.** This MCP deliberately exposes only
> read and messaging capabilities. It does **not** expose group/community
> **mutations** (create, leave, participants, settings, name, topic, photo,
> invite, join, requests), **community** operations, the operator **admin** plane
> (sessions), or **metrics** — even though the underlying SDK and gateway support
> them. This keeps an autonomous agent from performing destructive or
> account-wide actions. A guard test (`internal/server/manifest_test.go`)
> enforces the exclusion. To perform those operations, call the gateway's REST
> API directly from a trusted backend.

### Messaging Tools

**Common send arguments.** Every send tool below accepts canonical chat
addressing plus optional reply/mention threading, in addition to its
tool-specific fields:

- `chat` (string): Canonical recipient — a bare number, a user JID
  (`@s.whatsapp.net`), a group JID (`@g.us`), or a `@lid`. Preferred; wins over
  `to` when both are set.
- `to` (string): Back-compat alias for the recipient. Either `chat` or `to` is required.
- `reply_to_id` (string, optional): Quote an existing message by id.
- `reply_to_sender` (string, optional): Author JID/number of the quoted message.
- `reply_to_text` (string, optional): Caller-supplied preview of the quoted text (the gateway is storeless and does not look it up).
- `mentions` (array of strings, optional): Numbers/JIDs to @-tag in the message.

#### send_text_message
Send a text message to a WhatsApp contact or group.

**Input:**
- `chat` (string): Canonical recipient (preferred) — see Common send arguments above.
- `to` (string): Back-compat recipient alias. Either `chat` or `to` is required.
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

#### send_audio_message
Send an audio message to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `audio_url` (string, required): URL of the audio file to send
- `is_ptt` (boolean, optional): Send as voice note bubble (default: false)
- `view_once` (boolean, optional): Send as view-once media (default: false)

**Returns:** Message ID and status

#### send_video_message
Send a video message to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `video_url` (string, required): URL of the video file to send
- `caption` (string, optional): Video caption
- `is_gif` (boolean, optional): Render as GIF-style playback (default: false)
- `view_once` (boolean, optional): Send as view-once media (default: false)

**Returns:** Message ID and status

#### send_document_message
Send a document message to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `document_url` (string, required): URL of the file to send
- `file_name` (string, optional): Override filename shown in WhatsApp
- `caption` (string, optional): Document caption

**Returns:** Message ID and status

#### send_location_message
Send a location message with GPS coordinates to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
  - Individual: `{phone}@s.whatsapp.net` (e.g., `6281234567890@s.whatsapp.net`)
  - Group: `{group_id}@g.us` (e.g., `120363xxxxx@g.us`)
- `latitude` (number, required): Geographic latitude of the location
- `longitude` (number, required): Geographic longitude of the location
- `name` (string, optional): Name of the location
- `address` (string, optional): Address of the location

**Returns:** Message ID and status

#### send_poll_message
Send a poll message with a question and options to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `question` (string, required): Poll question text
- `options` (array of strings, required): List of poll options (minimum 2)
- `selectable_count` (integer, optional): Maximum number of options a user can select (0 = no limit)

**Returns:** Message ID and status

#### send_sticker_message
Send a sticker message from a URL to a WhatsApp contact or group.

**Input:**
- `to` (string, required): Recipient address in JID format
- `sticker_url` (string, required): URL of the sticker to send (WebP format)

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

### Inbox Tools

#### get_latest_incoming_messages
Fetch the most recent incoming WhatsApp messages buffered by the gateway. Useful for reading OTPs, verification codes, or recent conversation context. Returns newest first.

**Input:**
- `limit` (integer, optional): Maximum number of messages to return. Defaults to 10. Clamped server-side to [1, 50].

**Returns:**
- `success` (boolean): True on a successful read.
- `timestamp` (integer): Unix milliseconds when the response was generated.
- `count` (integer): Number of messages returned.
- `messages` (array): Each entry mirrors the webhook payload vocabulary:
  - `message_id`, `from`, `chat`, `is_group`, `push_name`, `timestamp`
  - `type`: one of `text` / `image` / `video` / `audio` / `document` / `sticker` / `contact` / `location` / `unknown`
  - `text`: present for text messages
  - `media`: object with `type`, `mime_type`, `size`, `file_name`, `caption` (metadata only — media URLs are not exposed by this tool; subscribe to webhooks if you need fetchable links)

**Notes:**
- The gateway keeps the most recent 100 incoming messages per session in memory. Buffer is lost on gateway restart — acceptable for short-lived OTP retrieval.
- Messages from the user (`is_from_me`) are excluded.
- `group_name` resolution is deferred — group messages return `is_group: true` and the group JID in `chat`.

### Contact & Group Read Tools

#### list_contacts
List the account's locally-synced WhatsApp contacts. Reads the local address book; an empty result is normal, never an error.

**Input:**
- `limit` (integer, optional): Page size. Gateway default 100, max 500.
- `offset` (integer, optional): Pagination offset.

**Returns:** `contacts[]` (`jid`, `push_name`, `full_name`, `first_name`, `business_name`), `count`, `total`, and an optional `note`.

#### get_contact_info
Look up one contact's WhatsApp profile.

**Input:**
- `chat` (string, required): Canonical recipient — a number, user JID, or `@lid`.

**Returns:** `jid`, `status`, `picture_id`, `verified_name`, `device_count`, `lid`.

#### get_avatar
Get a chat's profile picture URL (user or group). Soft failures are surfaced as results, not errors.

**Input:**
- `chat` (string, required): Canonical recipient (user or `@g.us` group).
- `preview` (boolean, optional): Request the low-res thumbnail instead of full-resolution.

**Returns:** `jid`, `available` (boolean). When available: `url` (time-limited CDN link), `id` (ETag), `type` (`image`/`preview`). When not: `available=false` with `reason` `not_set` (404, no picture) or `hidden` (403, privacy).

#### list_groups
List the account's joined groups as lightweight summaries (no participant roster; use `get_group_info` for one group's full detail).

**Input:** None

**Returns:** `groups[]` (`jid`, `name`, `topic`, `owner_jid`, `participant_count`, `is_announce`, `is_locked`, `is_community`) and `count`. Server-hitting read subject to a per-account budget (`429` when exhausted).

#### get_group_info
Get one group's full detail plus its participant roster.

**Input:**
- `chat` (string, required): A group JID (`@g.us`). The account must be a member.

**Returns:** Group detail plus `participants[]` (`jid`, `phone_number`, `lid`, `is_admin`, `is_super_admin`). `403` if not a member, `404` if absent.

### Conversation Tools

> These are conversation-affecting **outbound** actions. They are governed by the
> gateway's outbound pacer (per-account pace + per-recipient cap); over-budget
> calls are paced or rejected with `429`. This is real pacing, not an interim cap.

#### mark_read
Mark one or more messages in a chat as read (blue ticks).

**Input:**
- `chat` (string, required): Canonical recipient.
- `message_ids` (array of strings, required): Message IDs to mark read.
- `sender` (string, optional): Message author's JID/number — required for group chats.

**Returns:** Success status.

#### send_typing
Set the typing indicator in a chat.

**Input:**
- `chat` (string, required): Canonical recipient.
- `state` (string, required): One of `composing` (typing…), `recording` (recording audio…), or `paused` (cleared).

**Returns:** Success status and the applied state.

### Connection Tools

#### check_connection_status
Check if the WhatsApp session is active.

**Input:** None

**Returns:** Authentication status and connection information

#### check_contact
Check whether a phone number is registered on WhatsApp.

**Input:**
- `msisdn` (string, required): Phone number or JID to validate

**Returns:** Query result with canonical JID, `is_on_whatsapp`, and optional `verified_name`

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

For comprehensive Docker deployment instructions including:
- Transport options (stdio vs HTTP+SSE)
- Production deployment with authentication
- Docker networking configuration
- Environment variable reference
- Troubleshooting Docker-specific issues

** See the complete [DOCKER_USAGE.md](DOCKER_USAGE.md) guide**

### Quick Docker Commands

```bash
# Pull the latest image
docker pull glennprays/mcp-whatsapp-gateway:latest

# Build from source (development only)
docker build -t mcp-whatsapp-gateway .
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

This MCP server requires a running WhatsApp Gateway instance.

** For comprehensive WhatsApp Gateway MCP setup, configuration, and usage documentation:**
- **[Main Documentation](https://waga.glennprays.com/#docs:mcp/introduction)** - Complete guide for MCP integration
- **[Docker Deployment Guide](DOCKER_USAGE.md)** - Docker-specific deployment instructions

**Additional Resources:**
- **WhatsApp Gateway Documentation:** https://waga.glennprays.com
- **GitHub Repository:** https://github.com/glennprays/whatsapp-gateway

Once your gateway is running and you have a JWT token, configure this MCP server:

```bash
export WAGA_BASE_URL="http://localhost:3000/api/v1"  # Your gateway URL
export WAGA_JWT_TOKEN="your_jwt_token_here"          # From gateway registration
```

Then run the MCP server using Docker (see [DOCKER_USAGE.md](DOCKER_USAGE.md)) or from source (see Installation section above).

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

** For comprehensive documentation, guides, and troubleshooting:**
- **[Main Documentation](https://waga.glennprays.com/#docs:mcp/introduction)** - Complete MCP integration guide
- **[Docker Deployment Guide](DOCKER_USAGE.md)** - Docker-specific instructions

**For issues and questions:**
- GitHub Issues: [Create an issue](https://github.com/glennprays/mcp-whatsapp-gateway/issues)
- WhatsApp Gateway Documentation: https://waga.glennprays.com
