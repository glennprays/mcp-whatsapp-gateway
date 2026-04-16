# Build stage
FROM golang:1.25.0-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-whatsapp-gateway .

# Runtime stage using distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary from builder
COPY --from=builder /app/mcp-whatsapp-gateway /mcp-whatsapp-gateway

# Expose MCP port (default 8080)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/mcp-whatsapp-gateway"]
