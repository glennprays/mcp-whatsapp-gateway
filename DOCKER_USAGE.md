# Docker Usage Guide

This guide provides comprehensive instructions for running the MCP WhatsApp Gateway server using Docker.

## Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Transport Options](#transport-options)
- [Running with stdio Transport](#running-with-stdio-transport)
- [Running with HTTP+SSE Transport](#running-with-httpsse-transport)
- [Configuring MCP Clients](#configuring-mcp-clients)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)

## Quick Start

```bash
# Pull the image
docker pull glennprays/mcp-whatsapp-gateway:latest

# Run with stdio transport (for Claude Desktop/Cursor/Claude Code)
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest

# Run with HTTP+SSE transport (for web-based clients)
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

## Prerequisites

### 1. WhatsApp Gateway Running

You need a running WhatsApp Gateway instance. For setup instructions, visit:
- **Documentation:** https://waga.glennprays.com
- **GitHub:** https://github.com/glennprays/whatsapp-gateway

### 2. JWT Token

After setting up the WhatsApp Gateway:
1. Register your phone number with the gateway
2. Obtain a JWT token from the gateway
3. Keep this token secure - you'll need it for configuration

### 3. Docker Installed

Ensure Docker is installed and running:
```bash
docker --version
docker info
```

## Transport Options

The MCP WhatsApp Gateway server supports two transport modes:

### stdio Transport
- **Use case:** Claude Desktop, Cursor, Claude Code CLI
- **Communication:** Standard input/output
- **Network:** None (direct communication)
- **Best for:** Local development, desktop applications

### HTTP+SSE Transport
- **Use case:** Web-based MCP clients, remote connections
- **Communication:** HTTP with Server-Sent Events
- **Network:** Requires port mapping
- **Best for:** Production deployments, web applications

## Running with stdio Transport

The stdio transport is ideal for use with Claude Desktop, Cursor IDE, and Claude Code CLI. It communicates via standard input/output, making it perfect for local development.

### Local Docker (stdio)

#### Development Mode

```bash
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

**Parameters explained:**
- `-i` - Keep stdin open (required for stdio transport)
- `--rm` - Remove container after exit (cleanup)
- `-e WAGA_BASE_URL` - Your WhatsApp Gateway URL
- `-e WAGA_JWT_TOKEN` - Your JWT authentication token

#### With Logging

```bash
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e LOG_LEVEL="debug" \
  glennprays/mcp-whatsapp-gateway:latest
```

#### Remote Gateway (stdio)

If your WhatsApp Gateway is running on a remote server:

```bash
docker run -i --rm \
  -e WAGA_BASE_URL="https://waga.example.com/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

### Container as Command (stdio)

For easier usage, you can create a shell alias or script:

#### Option 1: Shell Alias

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
alias whatsapp-mcp='docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest'
```

Usage:
```bash
whatsapp-mcp
```

#### Option 2: Shell Script

Create `~/whatsapp-mcp.sh`:

```bash
#!/bin/bash
docker run -i --rm \
  -e WAGA_BASE_URL="$WAGA_BASE_URL" \
  -e WAGA_JWT_TOKEN="$WAGA_JWT_TOKEN" \
  glennprays/mcp-whatsapp-gateway:latest "$@"
```

Make it executable:
```bash
chmod +x ~/whatsapp-mcp.sh
```

Usage:
```bash
export WAGA_BASE_URL="http://host.docker.internal:3000/api/v1"
export WAGA_JWT_TOKEN="your_jwt_token"
~/whatsapp-mcp.sh
```

## Running with HTTP+SSE Transport

The HTTP+SSE transport is ideal for production deployments and web-based MCP clients. It runs a web server that communicates via HTTP with Server-Sent Events.

### Local Docker (HTTP+SSE)

#### Development Mode (No Authentication)

```bash
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e MCP_PORT="8080" \
  -e APP_ENV="development" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

**Parameters explained:**
- `-d` - Run in detached mode (background)
- `--name` - Give the container a name
- `-p 8080:8080` - Map host port 8080 to container port 8080
- `-e MCP_TRANSPORT="http"` - Enable HTTP+SSE transport
- `-e MCP_PORT="8080"` - Port for HTTP server
- `-e APP_ENV="development"` - Development mode (no auth required)

#### Verify the Server is Running

```bash
# Check if container is running
docker ps | grep whatsapp-gateway-mcp

# Check logs
docker logs whatsapp-gateway-mcp

# Test the endpoint
curl http://localhost:8080/mcp
```

#### View Logs

```bash
# Follow logs in real-time
docker logs -f whatsapp-gateway-mcp

# View last 50 lines
docker logs --tail 50 whatsapp-gateway-mcp
```

#### Stop the Container

```bash
# Stop the container
docker stop whatsapp-gateway-mcp

# Remove the container
docker rm whatsapp-gateway-mcp
```

### Production Mode (With Authentication)

For production deployments, enable Basic authentication:

```bash
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e MCP_PORT="8080" \
  -e APP_ENV="production" \
  -e MCP_BASIC_AUTH_USER="admin" \
  -e MCP_BASIC_AUTH_PASSWORD="your_secure_password" \
  -e WAGA_BASE_URL="https://waga.example.com/api/v1" \
  -e WAGA_JWT_TOKEN="your_production_jwt_token" \
  --restart unless-stopped \
  glennprays/mcp-whatsapp-gateway:latest
```

**Additional parameters:**
- `-e MCP_BASIC_AUTH_USER` - Username for Basic authentication
- `-e MCP_BASIC_AUTH_PASSWORD` - Password for Basic authentication
- `--restart unless-stopped` - Auto-restart on failure or system reboot

#### Test Authentication

```bash
# Test with authentication
curl -u admin:your_secure_password http://localhost:8080/mcp

# Test without authentication (should return 401)
curl http://localhost:8080/mcp
```

### Docker Compose (HTTP+SSE)

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  whatsapp-gateway-mcp:
    image: glennprays/mcp-whatsapp-gateway:latest
    container_name: whatsapp-gateway-mcp
    ports:
      - "8080:8080"
    environment:
      - MCP_TRANSPORT=http
      - MCP_PORT=8080
      - APP_ENV=production
      - MCP_BASIC_AUTH_USER=${MCP_BASIC_AUTH_USER}
      - MCP_BASIC_AUTH_PASSWORD=${MCP_BASIC_AUTH_PASSWORD}
      - WAGA_BASE_URL=${WAGA_BASE_URL}
      - WAGA_JWT_TOKEN=${WAGA_JWT_TOKEN}
      - LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/mcp"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
```

Create a `.env` file:

```bash
MCP_BASIC_AUTH_USER=admin
MCP_BASIC_AUTH_PASSWORD=your_secure_password
WAGA_BASE_URL=https://waga.example.com/api/v1
WAGA_JWT_TOKEN=your_production_jwt_token
```

Start the service:

```bash
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

## Configuring MCP Clients

### Claude Code CLI

Claude Code CLI can use the MCP server via stdio transport.

#### Option 1: Direct Docker Command

Edit your Claude Code configuration:

```bash
# macOS
~/Library/Application\ Support/Claude/claude_desktop_config.json

# Windows
%APPDATA%\Claude\claude_desktop_config.json

# Linux
~/.config/Claude/claude_desktop_config.json
```

Add the MCP server configuration:

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

#### Option 2: Shell Script Wrapper

Create a wrapper script at `/usr/local/bin/whatsapp-mcp`:

```bash
#!/bin/bash
docker run -i --rm \
  -e WAGA_BASE_URL="$WAGA_BASE_URL" \
  -e WAGA_JWT_TOKEN="$WAGA_JWT_TOKEN" \
  glennprays/mcp-whatsapp-gateway:latest
```

Make it executable:
```bash
chmod +x /usr/local/bin/whatsapp-mcp
```

Configure Claude Code:

```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "/usr/local/bin/whatsapp-mcp",
      "env": {
        "WAGA_BASE_URL": "http://host.docker.internal:3000/api/v1",
        "WAGA_JWT_TOKEN": "your_jwt_token"
      }
    }
  }
}
```

#### Option 3: Built Binary (Alternative)

If you prefer to use a built binary instead of Docker:

```bash
# Download or build the binary
go install github.com/glennprays/mcp-whatsapp-gateway@latest

# Configure Claude Code
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "/Users/your-username/go/bin/mcp-whatsapp-gateway",
      "env": {
        "WAGA_BASE_URL": "http://localhost:3000/api/v1",
        "WAGA_JWT_TOKEN": "your_jwt_token"
      }
    }
  }
}
```

### Claude Desktop

For Claude Desktop, use the stdio transport with Docker:

#### macOS Configuration

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

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

#### Windows Configuration

Edit `%APPDATA%\Claude\claude_desktop_config.json`:

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

**Note for Windows:** Use `host.docker.internal` to access host services from Docker.

#### Linux Configuration

Edit `~/.config/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "whatsapp-gateway": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "WAGA_BASE_URL=http://172.17.0.1:3000/api/v1",
        "-e", "WAGA_JWT_TOKEN=your_jwt_token",
        "glennprays/mcp-whatsapp-gateway:latest"
      ]
    }
  }
}
```

**Note for Linux:** Use `172.17.0.1` (default Docker bridge) or your host's IP address.

### Cursor IDE

Cursor IDE uses a similar configuration format:

#### macOS/Linux

Edit `~/.cursor/mcp.json` or use Cursor's settings UI:

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

#### Windows

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

### HTTP+SSE with Web Clients

For web-based MCP clients using HTTP+SSE transport:

```bash
# Start the MCP server with HTTP+SSE
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  glennprays/mcp-whatsapp-gateway:latest
```

The MCP endpoint will be available at:
- **Local:** `http://localhost:8080/mcp`
- **Production:** `https://your-server.com:8080/mcp` (with authentication)

## Production Deployment

### Production Checklist

- [ ] Use `APP_ENV=production`
- [ ] Enable Basic authentication
- [ ] Use HTTPS for WAGA_BASE_URL
- [ ] Set strong passwords
- [ ] Use `--restart unless-stopped`
- [ ] Monitor container health
- [ ] Set up log aggregation
- [ ] Configure resource limits

### Production Docker Run Command

```bash
docker run -d --name whatsapp-gateway-mcp \
  -p 8080:8080 \
  -e MCP_TRANSPORT="http" \
  -e MCP_PORT="8080" \
  -e APP_ENV="production" \
  -e MCP_BASIC_AUTH_USER="admin" \
  -e MCP_BASIC_AUTH_PASSWORD="your_strong_password_here" \
  -e WAGA_BASE_URL="https://waga.example.com/api/v1" \
  -e WAGA_JWT_TOKEN="your_production_jwt_token" \
  -e LOG_LEVEL="info" \
  --restart unless-stopped \
  --memory="512m" \
  --cpus="1.0" \
  --health-cmd="wget --spider -q http://localhost:8080/mcp || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  glennprays/mcp-whatsapp-gateway:latest
```

### Production Docker Compose

```yaml
version: '3.8'

services:
  whatsapp-gateway-mcp:
    image: glennprays/mcp-whatsapp-gateway:latest
    container_name: whatsapp-gateway-mcp
    ports:
      - "8080:8080"
    environment:
      - MCP_TRANSPORT=http
      - MCP_PORT=8080
      - APP_ENV=production
      - MCP_BASIC_AUTH_USER=${MCP_BASIC_AUTH_USER}
      - MCP_BASIC_AUTH_PASSWORD=${MCP_BASIC_AUTH_PASSWORD}
      - WAGA_BASE_URL=${WAGA_BASE_URL}
      - WAGA_JWT_TOKEN=${WAGA_JWT_TOKEN}
      - LOG_LEVEL=info
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/mcp"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Monitoring and Logs

```bash
# View logs
docker logs -f whatsapp-gateway-mcp

# Check container health
docker inspect --format='{{.State.Health.Status}}' whatsapp-gateway-mcp

# View resource usage
docker stats whatsapp-gateway-mcp

# Check recent logs
docker logs --since 1h whatsapp-gateway-mcp
```

## Troubleshooting

### Container Issues

**Problem:** Container exits immediately

```bash
# Check logs
docker logs whatsapp-gateway-mcp

# Common causes:
# 1. Missing required environment variables
# 2. Invalid JWT token
# 3. Gateway not reachable
```

**Problem:** Cannot connect to WhatsApp Gateway

```bash
# Test gateway connectivity from container
docker run --rm \
  glennprays/mcp-whatsapp-gateway:latest \
  curl http://host.docker.internal:3000/api/v1/health

# For remote gateways, test from host:
curl https://waga.example.com/api/v1/health
```

**Problem:** Permission errors

```bash
# Check if port is already in use
lsof -i :8080

# Use a different port
docker run -p 9090:8080 \
  -e MCP_PORT="8080" \
  ...
```

### Docker Networking

**Problem:** Container cannot reach host services

**Solutions:**

1. **macOS/Windows:** Use `host.docker.internal`
   ```bash
   -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1"
   ```

2. **Linux:** Use Docker bridge IP
   ```bash
   -e WAGA_BASE_URL="http://172.17.0.1:3000/api/v1"
   ```

3. **Linux:** Use host network (not recommended for production)
   ```bash
   docker run --network host ...
   -e WAGA_BASE_URL="http://localhost:3000/api/v1"
   ```

### MCP Client Issues

**Problem:** MCP tools not showing in Claude Desktop

**Solutions:**

1. **Check configuration file syntax:**
   ```bash
   # Validate JSON
   cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .
   ```

2. **Check Claude Desktop logs:**
   ```bash
   # macOS
   ~/Library/Logs/Claude/

   # Windows
   %APPDATA%\Claude\logs\
   ```

3. **Test MCP server manually:**
   ```bash
   docker run -i --rm \
     -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
     -e WAGA_JWT_TOKEN="your_token" \
     glennprays/mcp-whatsapp-gateway:latest

   # Should see MCP initialization message
   ```

**Problem:** 401 Unauthorized in production

**Solution:**
```bash
# Ensure Basic auth credentials are set
-e MCP_BASIC_AUTH_USER="admin"
-e MCP_BASIC_AUTH_PASSWORD="your_password"

# Test with curl
curl -u admin:your_password http://localhost:8080/mcp
```

### Authentication Issues

**Problem:** JWT token errors

**Solutions:**

1. **Verify token is valid:**
   ```bash
   curl -H "Authorization: Bearer your_token" \
        http://localhost:3000/api/v1/health
   ```

2. **Check token hasn't expired**
3. **Re-register phone with gateway** to get new token

**Problem:** Session disconnected (403 Forbidden)

**Solution:**
1. Check WhatsApp Gateway connection status
2. Reconnect gateway if needed
3. Verify WhatsApp session is active

### Debug Mode

Enable debug logging for troubleshooting:

```bash
docker run -i --rm \
  -e WAGA_BASE_URL="http://host.docker.internal:3000/api/v1" \
  -e WAGA_JWT_TOKEN="your_jwt_token" \
  -e LOG_LEVEL="debug" \
  glennprays/mcp-whatsapp-gateway:latest
```

## Environment Variables Reference

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `WAGA_BASE_URL` | WhatsApp Gateway API base URL | `http://localhost:3000/api/v1` |
| `WAGA_JWT_TOKEN` | JWT authentication token | `eyJhbGciOi...` |

### Optional Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `APP_ENV` | Application environment | `development` | `production` |
| `LOG_LEVEL` | Logging level | `info` | `debug`, `info`, `warn`, `error` |
| `MCP_TRANSPORT` | Transport type | `stdio` | `stdio`, `http` |
| `MCP_PORT` | HTTP server port | `8080` | `8080`, `9090` |
| `MCP_BASIC_AUTH_USER` | Basic auth username (prod only) | - | `admin` |
| `MCP_BASIC_AUTH_PASSWORD` | Basic auth password (prod only) | - | `secure_password` |

## Support

For issues and questions:
- **GitHub Issues:** https://github.com/glennprays/mcp-whatsapp-gateway/issues
- **WhatsApp Gateway Documentation:** https://waga.glennprays.com

## Additional Resources

- [README.md](README.md) - Main project documentation
- [WhatsApp Gateway](https://github.com/glennprays/whatsapp-gateway) - The underlying gateway service
- [MCP Protocol](https://modelcontextprotocol.io) - Protocol specification
