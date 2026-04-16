#!/bin/bash
set -e

VERSION=${1:-latest}
REGISTRY=${2:-ghcr.io}  # Options: ghcr.io or docker.io
USERNAME=${3:-glennprays}

echo "Building and publishing mcp-whatsapp-gateway..."
echo "Version: $VERSION"
echo "Registry: $REGISTRY"
echo "Username: $USERNAME"

# Build the image
echo "Building Docker image..."
docker build -t mcp-whatsapp-gateway .

# Tag for GHCR
if [ "$REGISTRY" == "ghcr.io" ]; then
    echo "Tagging for GitHub Container Registry..."
    docker tag mcp-whatsapp-gateway ghcr.io/$USERNAME/mcp-whatsapp-gateway:$VERSION
    docker tag mcp-whatsapp-gateway ghcr.io/$USERNAME/mcp-whatsapp-gateway:latest

    echo "Pushing to GHCR..."
    docker push ghcr.io/$USERNAME/mcp-whatsapp-gateway:$VERSION
    docker push ghcr.io/$USERNAME/mcp-whatsapp-gateway:latest

    echo "✅ Published to GHCR!"
    echo "Image: ghcr.io/$USERNAME/mcp-whatsapp-gateway:$VERSION"
    echo "Image: ghcr.io/$USERNAME/mcp-whatsapp-gateway:latest"
fi

# Tag for Docker Hub
if [ "$REGISTRY" == "docker.io" ]; then
    echo "Tagging for Docker Hub..."
    docker tag mcp-whatsapp-gateway $USERNAME/mcp-whatsapp-gateway:$VERSION
    docker tag mcp-whatsapp-gateway $USERNAME/mcp-whatsapp-gateway:latest

    echo "Pushing to Docker Hub..."
    docker push $USERNAME/mcp-whatsapp-gateway:$VERSION
    docker push $USERNAME/mcp-whatsapp-gateway:latest

    echo "✅ Published to Docker Hub!"
    echo "Image: $USERNAME/mcp-whatsapp-gateway:$VERSION"
    echo "Image: $USERNAME/mcp-whatsapp-gateway:latest"
fi

echo ""
echo "Users can now run:"
echo ""
if [ "$REGISTRY" == "ghcr.io" ]; then
    echo "docker pull ghcr.io/$USERNAME/mcp-whatsapp-gateway:$VERSION"
else
    echo "docker pull $USERNAME/mcp-whatsapp-gateway:$VERSION"
fi
