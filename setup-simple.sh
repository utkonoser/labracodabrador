#!/bin/bash
set -e

echo "🚀 Setting up Production Ethereum PoS Network"
echo "=============================================="
echo ""

# Создаем директории
mkdir -p jwt systemd logs

# Генерируем JWT secret
echo "🔐 Generating JWT secret..."
if [ ! -f "jwt/jwt.hex" ]; then
    openssl rand -hex 32 > jwt/jwt.hex
    echo "✓ JWT secret created"
else
    echo "✓ JWT secret exists"
fi

# Инициализируем Geth
echo "🔨 Initializing Geth..."
docker run --rm \
    -v $(pwd)/geth-data:/data \
    -v $(pwd)/genesis-pos.json:/genesis.json:ro \
    ethereum/client-go:v1.14.11 \
    init --datadir=/data /genesis.json

echo "✓ Geth initialized"

echo ""
echo "✅ Setup complete!"
echo ""
echo "🚀 Start the network:"
echo "   docker-compose up -d"
echo ""
echo "📊 Monitor:"
echo "   docker-compose ps"
echo "   docker-compose logs -f"
echo ""
echo "🦊 Metamask:"
echo "   RPC: http://localhost:8545"
echo "   Chain ID: 32382"
echo ""

