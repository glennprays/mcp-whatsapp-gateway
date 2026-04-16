# mcp-whatsapp-gateway — Project Brief

## Overview

`mcp-whatsapp-gateway` is an MCP (Model Context Protocol) server that exposes the
[WhatsApp Gateway (waga)](https://waga.glennprays.com) as tools for AI agents. It allows
agents like Claude to send WhatsApp messages, manage webhooks, and check connection
status — all through a pre-authenticated JWT session.

## Repository

**Name:** `mcp-whatsapp-gateway`
**Language:** Go
**MCP Library:** `github.com/modelcontextprotocol/go-sdk` (official SDK)
**MCP Transport:** Both `stdio` and `HTTP+SSE` supported in the same binary, selected via config

---

## Auth Strategy: Pre-Login JWT

This phase uses a **pre-authenticated JWT only**. There is no login flow inside the MCP server.

- The operator registers a phone number against the gateway and obtains a JWT **outside**
  of this MCP server (e.g. via `curl` or a setup script)
- The operator logs into WhatsApp once (QR or pair code) via the gateway directly
- The JWT is injected into the MCP server at startup via environment variable
- The MCP server assumes the session is alive and simply uses the JWT on every request

**No login, logout, or QR/pair code tools are implemented in this phase.**

If the session drops, the operator must reconnect manually via the gateway. The
`check_connection_status` tool is provided so agents can detect this and inform the user.

---

## References

### Gateway API
- **llms.txt:** https://waga.glennprays.com/llms.txt — read this for gateway capabilities,
  endpoint structure, auth flow, webhook events, and recipient address formats
- **OpenAPI spec:** https://waga.glennprays.com/openapi.yaml — canonical reference for all
  request/response schemas

### Go SDK
- **SDK:** `github.com/glennprays/whatsapp-gateway-sdk-go`
- Use this SDK for all HTTP calls to the gateway. Do not hand-roll HTTP clients.
- Check the SDK source for available methods before implementing any tool.

---

## Transport

Both transports are supported in the same binary, controlled by the `MCP_TRANSPORT` env var.

| `MCP_TRANSPORT` | Behavior |
|---|---|
| `stdio` | Runs as stdio MCP server (for Claude Desktop, Cursor, Claude Code) |
| `http` | Runs HTTP+SSE MCP server on `MCP_PORT` |

Default: `stdio`

### HTTP+SSE Auth

Auth is **environment-dependent**:

- `APP_ENV=dev` — no auth, open access
- `APP_ENV=prod` — **Basic auth** required on all HTTP+SSE endpoints

Basic auth credentials are set via `MCP_BASIC_AUTH_USER` and `MCP_BASIC_AUTH_PASSWORD`.
Requests without valid credentials in prod return `401`. Fail fast on startup if these are
missing when `APP_ENV=prod` and `MCP_TRANSPORT=http`.

---

## Configuration

All config loaded from environment variables at startup. No config files.

| Variable | Required | Default | Description |
|---|---|---|---|
| `WAGA_BASE_URL` | Yes | — | Gateway base URL e.g. `http://localhost:3000/api/v1` |
| `WAGA_JWT_TOKEN` | Yes | — | Pre-obtained JWT for the WhatsApp account |
| `APP_ENV` | No | `dev` | `dev` or `prod` |
| `LOG_LEVEL` | No | `info` | `debug`, `info`, `warn`, `error` |
| `MCP_TRANSPORT` | No | `stdio` | `stdio` or `http` |
| `MCP_PORT` | No | `8080` | HTTP+SSE server port (only used when `MCP_TRANSPORT=http`) |
| `MCP_BASIC_AUTH_USER` | prod+http only | — | Basic auth username for HTTP+SSE |
| `MCP_BASIC_AUTH_PASSWORD` | prod+http only | — | Basic auth password for HTTP+SSE |

### Startup Validation

On startup, in order:

1. Fail fast if `WAGA_BASE_URL` or `WAGA_JWT_TOKEN` are missing or empty
2. Fail fast if `APP_ENV=prod` + `MCP_TRANSPORT=http` and Basic auth credentials are missing
3. Ping gateway `GET /health` — if unreachable, **log a warning and continue** (do not panic)
4. Log the resolved config (transport, env, log level) — never log secrets

---

## MCP Tools

### Recipient Address Format
- Individual: `{phone}@s.whatsapp.net` e.g. `6281234567890@s.whatsapp.net`
- Group: `{group_id}@g.us` e.g. `120363xxxxx@g.us`
- Always document this in each tool's input schema description.

### Messaging

| Tool | Gateway Endpoint | Description |
|---|---|---|
| `send_text_message` | POST `/message/text` | Send a text message to a contact or group |
| `send_image_message` | POST `/message/image` | Send an image with optional caption and view-once flag |
| `edit_message` | PUT `/message` | Edit a previously sent message by message ID |
| `delete_message` | DELETE `/message` | Delete a sent message by message ID |
| `react_to_message` | POST `/message/react` | React to a message with an emoji |

### Connection

| Tool | Gateway Endpoint | Description |
|---|---|---|
| `check_connection_status` | GET `/login/status` | Check if the WhatsApp session is active |

### Health

| Tool | Gateway Endpoint | Description |
|---|---|---|
| `check_health` | GET `/health` | Check if the gateway service is reachable |

### Webhook Management

| Tool | Gateway Endpoint | Description |
|---|---|---|
| `get_webhook` | GET `/webhook` | Get the currently registered webhook URL |
| `register_webhook` | POST `/webhook` | Register a webhook URL with optional HMAC secret |
| `delete_webhook` | DELETE `/webhook` | Remove the registered webhook |

---

## Project Structure

```
mcp-whatsapp-gateway/
├── main.go                  # Entrypoint: init logger, load config, validate, start server
├── config/
│   └── config.go            # Load and validate env vars, fail fast logic
├── server/
│   ├── stdio.go             # stdio transport setup
│   └── http.go              # HTTP+SSE transport setup + Basic auth middleware
├── tools/
│   ├── messaging.go         # send_text, send_image, edit, delete, react
│   ├── connection.go        # check_connection_status
│   ├── health.go            # check_health
│   └── webhook.go           # get_webhook, register_webhook, delete_webhook
├── gateway/
│   └── client.go            # Initializes the waga Go SDK client from config
├── Dockerfile
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

---

## Logging

Use `github.com/glennprays/log` for all logging. Do not use `fmt.Println`, `log.Printf`,
or any other logging package.

### Initialization

```go
logger, err := log.New(log.Config{
    Service:      "mcp-whatsapp-gateway",
    Env:          cfg.AppEnv,
    Level:        resolvedLevel,        // mapped from LOG_LEVEL env var
    Output:       log.OutputStdout,
    EnableCaller: cfg.AppEnv == "dev",  // caller info in dev only
})
if err != nil {
    panic(err)
}
defer logger.Sync()
```

### Log Level by Environment

| `APP_ENV` | Default Level | `EnableCaller` |
|---|---|---|
| `dev` | `debug` | `true` |
| `prod` | `info` | `false` |

`LOG_LEVEL` env var overrides the default level in both environments.

### Usage Pattern

- Every log call **must** include a non-empty `traceId`. Empty string will panic.
- For MCP tool calls, use the MCP request ID if available, otherwise generate a UUID.
- Use `nil` metadata for simple internal logs.
- Use `metadata` map for request-scoped context (tool name, recipient, transport, etc).
- **Never log** JWT tokens, Basic auth credentials, HMAC secrets, or message content at
  `Info` level or above. Debug level only if truly necessary.

```go
// Tool invocation
logger.Info(traceID, "tool invoked", map[string]any{
    "tool":      "send_text_message",
    "transport": cfg.Transport,
}, log.String("env", cfg.AppEnv))

// Gateway warning (startup health check)
logger.Warn(traceID, "gateway health check failed, continuing", map[string]any{
    "base_url": cfg.WagaBaseURL,
}, log.Error(err))

// Tool error
logger.Error(traceID, "tool execution failed", map[string]any{
    "tool": "send_text_message",
}, log.Error(err))
```

---

## Error Handling

- All gateway errors must be returned as MCP tool errors, not panics.
- Wrap SDK errors with context: `fmt.Errorf("send_text_message: %w", err)`
- On `401` from gateway: return — "JWT token is invalid or expired. Re-register the phone
  number against the gateway to obtain a new token."
- On `403` from gateway: return — "Session may be disconnected. Run
  `check_connection_status` to verify."
- On `500` from gateway: return the error with a suggestion to check gateway logs.

---

## Docker

Provide a single `Dockerfile` using a multi-stage build:

- **Stage 1:** Build the Go binary
- **Stage 2:** Minimal runtime image (distroless or Alpine)
- Expose `MCP_PORT` (default `8080`) — only relevant when `MCP_TRANSPORT=http`
- All config via environment variables, no config files baked into the image

Provide a `.env.example` with all variables documented and safe placeholder values.

---

## Documentation
Create proper documentation on README.md, and Docker as default running process

---

## Out of Scope (Phase 1 / MVP)

- Login / logout / QR code / pair code tools
- Incoming message handling or webhook receiving
- Multi-account support
- Message history, read receipts, or presence
- Incoming message store or polling tool
- docker-compose (Dockerfile only for now)
- Any background goroutines or persistent state inside the MCP process

