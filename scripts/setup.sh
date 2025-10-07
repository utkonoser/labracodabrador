#!/bin/bash
set -e

echo "==================================="
echo "Labracodabrador Blockchain Setup"
echo "==================================="

# Check Go version
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go version: $GO_VERSION"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Please run this script from project root"
    exit 1
fi

echo ""
echo "Installing dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installed"

echo ""
echo "Building project..."
mkdir -p bin
go build -o bin/labracodabrador ./cmd/labracodabrador
echo "✓ Build complete"

echo ""
echo "Creating directories..."
mkdir -p data logs
echo "✓ Directories created"

echo ""
echo "==================================="
echo "Setup complete!"
echo "==================================="
echo ""
echo "To start the blockchain, run:"
echo "  ./bin/labracodabrador"
echo ""
echo "Or use make:"
echo "  make run"
echo ""

