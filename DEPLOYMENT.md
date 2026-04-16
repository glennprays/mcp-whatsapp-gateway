# Deployment Guide

This guide covers testing and deploying the MCP WhatsApp Gateway server.

## 🧪 Local Development Testing

### Prerequisites

1. **Running WhatsApp Gateway**
   ```bash
   # Make sure your gateway is accessible
   curl http://localhost:3000/api/v1/health
   ```

2. **Valid JWT Token**
   - Register your phone with the gateway
   - Obtain a JWT token
   - Keep it secure!

### Development Testing

#### Option 1: Quick Test with `go run`

```bash
# Set environment variables
export WAGA_BASE_URL="http://localhost:3000/api/v1"
export WAGA_JWT_TOKEN="your_jwt_token"

# Test stdio transport (for Claude Desktop/Cursor)
go run main.go

# Test HTTP+SSE transport
export MCP_TRANSPORT="http"
export MCP_PORT="8080"
go run main.go
```

#### Option 2: Build and Test

```bash
# Build the binary
go build -o mcp-whatsapp-gateway

# Run with stdio transport
WAGA_BASE_URL="http://localhost:3000/api/v1" \
WAGA_JWT_TOKEN="your_jwt_token" \
./mcp-whatsapp-gateway
```

### Testing with MCP Clients

#### Claude Desktop Configuration

Edit your Claude Desktop config file:

**Mac**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "/path/to/mcp-whatsapp-gateway",
      "args": [],
      "env": {
        "WAGA_BASE_URL": "http://localhost:3000/api/v1",
        "WAGA_JWT_TOKEN": "your_jwt_token"
      }
    }
  }
}
```

#### Cursor IDE Configuration

```json
{
  "mcpServers": [
    {
      "name": "whatsapp-gateway",
      "command": "/path/to/mcp-whatsapp-gateway",
      "env": {
        "WAGA_BASE_URL": "http://localhost:3000/api/v1",
        "WAGA_JWT_TOKEN": "your_jwt_token"
      }
    }
  ]
}
```

### Verification

Test that tools are available:

```bash
# In Claude Desktop/Cursor, you should see these tools:
# - send_text_message
# - send_image_message
# - edit_message
# - delete_message
# - react_to_message
# - check_connection_status
# - check_health
# - get_webhook
# - register_webhook
# - delete_webhook
```

## 🐳 Docker Deployment

### Local Docker Testing

```bash
# Build the image
docker build -t mcp-whatsapp-gateway .

# Test with stdio transport (interactive)
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  mcp-whatsapp-gateway

# Test with HTTP+SSE transport
docker run -d --name whatsapp-gateway \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  mcp-whatsapp-gateway

# Test HTTP endpoint
curl http://localhost:8080/mcp

# Clean up
docker stop whatsapp-gateway
docker rm whatsapp-gateway
```

### Publishing to Container Registries

There's no specific "MCP hub" - MCP servers are distributed via regular container registries.

#### Option 1: GitHub Container Registry (GHCR) - Recommended

```bash
# Set version
VERSION=v1.0.0

# Tag for GHCR
docker tag mcp-whatsapp-gateway ghcr.io/glennprays/mcp-whatsapp-gateway:$VERSION
docker tag mcp-whatsapp-gateway ghcr.io/glennprays/mcp-whatsapp-gateway:latest

# Login to GHCR
echo "YOUR_GITHUB_TOKEN" | docker login ghcr.io -u glennprays --password-stdin

# Push to GHCR
docker push ghcr.io/glennprays/mcp-whatsapp-gateway:$VERSION
docker push ghcr.io/glennprays/mcp-whatsapp-gateway:latest
```

#### Option 2: Docker Hub

```bash
# Set version
VERSION=v1.0.0

# Tag for Docker Hub
docker tag mcp-whatsapp-gateway glennprays/mcp-whatsapp-gateway:$VERSION
docker tag mcp-whatsapp-gateway glennprays/mcp-whatsapp-gateway:latest

# Login to Docker Hub
docker login

# Push to Docker Hub
docker push glennprays/mcp-whatsapp-gateway:$VERSION
docker push glennprays/mcp-whatsapp-gateway:latest
```

### Using the Publishing Script

```bash
# Publish to GHCR (default)
./publish-docker.sh v1.0.0 ghcr.io glennprays

# Publish to Docker Hub
./publish-docker.sh v1.0.0 docker.io glennprays
```

## 🌐 Production Deployment

### Production Configuration

```bash
docker run -d --name whatsapp-gateway \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e MCP_TRANSPORT=http \
  -e MCP_PORT=8080 \
  -e MCP_BASIC_AUTH_USER=admin \
  -e MCP_BASIC_AUTH_PASSWORD=your_secure_password \
  -e WAGA_BASE_URL=https://your-gateway.com/api/v1 \
  -e WAGA_JWT_TOKEN=your_production_jwt_token \
  -e LOG_LEVEL=info \
  --restart unless-stopped \
  ghcr.io/glennprays/mcp-whatsapp-gateway:latest
```

### Docker Compose Deployment

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  mcp-whatsapp-gateway:
    image: ghcr.io/glennprays/mcp-whatsapp-gateway:latest
    container_name: whatsapp-gateway-mcp
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - MCP_TRANSPORT=http
      - MCP_PORT=8080
      - MCP_BASIC_AUTH_USER=${MCP_BASIC_AUTH_USER}
      - MCP_BASIC_AUTH_PASSWORD=${MCP_BASIC_AUTH_PASSWORD}
      - WAGA_BASE_URL=${WAGA_BASE_URL}
      - WAGA_JWT_TOKEN=${WAGA_JWT_TOKEN}
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/mcp"]
      interval: 30s
      timeout: 10s
      retries: 3
```

```bash
# Create .env file
cat > .env << EOF
MCP_BASIC_AUTH_USER=admin
MCP_BASIC_AUTH_PASSWORD=your_secure_password
WAGA_BASE_URL=https://your-gateway.com/api/v1
WAGA_JWT_TOKEN=your_production_jwt_token
EOF

# Deploy
docker-compose up -d
```

## 🔍 MCP Server Discovery

### Adding to Claude Desktop

After publishing, users can add your MCP server:

**Option 1: Local Binary**
```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "/path/to/mcp-whatsapp-gateway",
      "env": {
        "WAGA_BASE_URL": "https://your-gateway.com/api/v1",
        "WAGA_JWT_TOKEN": "user_jwt_token"
      }
    }
  }
}
```

**Option 2: Docker Container**
```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "WAGA_BASE_URL=https://your-gateway.com/api/v1",
        "-e", "WAGA_JWT_TOKEN=user_jwt_token",
        "ghcr.io/glennprays/mcp-whatsapp-gateway:latest"
      ]
    }
  }
}
```

### Publishing to MCP Marketplace

Currently, there's no centralized "MCP hub" or marketplace. MCP servers are distributed through:

1. **GitHub**: Source code and documentation
2. **Container Registries**: Docker Hub, GHCR, etc.
3. **Community Discovery**: Through blog posts, social media, and word-of-mouth

**Best practices for discovery:**
- Clear documentation with examples
- Proper tagging in Docker registries (mcp, whatsapp, claude)
- GitHub issues/discussions for support
- Blog posts or tutorials
- Social media announcements

## 📋 Verification Checklist

Before publishing:

- [ ] All tests pass locally (`go test ./...`)
- [ ] Docker image builds successfully
- [ ] Test with real WhatsApp Gateway
- [ ] Verify all 11 tools work correctly
- [ ] Test with Claude Desktop/Cursor
- [ ] Documentation is complete
- [ ] Security review completed
- [ ] Version tags are semantic
- [ ] Environment variables documented

After publishing:

- [ ] Test pulling the published image
- [ ] Verify from clean environment
- [ ] Update README with published image locations
- [ ] Create GitHub release
- [ ] Announce to community
