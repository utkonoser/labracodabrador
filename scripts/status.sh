#!/bin/bash

echo "==================================="
echo "Labracodabrador Status Check"
echo "==================================="

# Check if process is running
if pgrep -f labracodabrador > /dev/null; then
    echo "✓ Process is running"
    echo ""
    echo "PIDs:"
    pgrep -f labracodabrador
else
    echo "✗ Process is not running"
    exit 1
fi

echo ""
echo "Checking nodes..."

# Check node 1
echo ""
echo "Node 1 (8545):"
if curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' > /dev/null; then
    
    BLOCK=$(curl -s -X POST http://localhost:8545 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')
    
    PEERS=$(curl -s -X POST http://localhost:8545 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' | jq -r '.result')
    
    MINING=$(curl -s -X POST http://localhost:8545 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_mining","params":[],"id":1}' | jq -r '.result')
    
    echo "  Status:  ✓ Running"
    echo "  Block:   $BLOCK"
    echo "  Peers:   $PEERS"
    echo "  Mining:  $MINING"
else
    echo "  Status:  ✗ Not responding"
fi

# Check node 2
echo ""
echo "Node 2 (8547):"
if curl -s -X POST http://localhost:8547 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' > /dev/null; then
    
    BLOCK=$(curl -s -X POST http://localhost:8547 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')
    
    PEERS=$(curl -s -X POST http://localhost:8547 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' | jq -r '.result')
    
    MINING=$(curl -s -X POST http://localhost:8547 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_mining","params":[],"id":1}' | jq -r '.result')
    
    echo "  Status:  ✓ Running"
    echo "  Block:   $BLOCK"
    echo "  Peers:   $PEERS"
    echo "  Mining:  $MINING"
else
    echo "  Status:  ✗ Not responding"
fi

# Check node 3
echo ""
echo "Node 3 (8549):"
if curl -s -X POST http://localhost:8549 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' > /dev/null; then
    
    BLOCK=$(curl -s -X POST http://localhost:8549 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')
    
    PEERS=$(curl -s -X POST http://localhost:8549 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' | jq -r '.result')
    
    MINING=$(curl -s -X POST http://localhost:8549 \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_mining","params":[],"id":1}' | jq -r '.result')
    
    echo "  Status:  ✓ Running"
    echo "  Block:   $BLOCK"
    echo "  Peers:   $PEERS"
    echo "  Mining:  $MINING"
else
    echo "  Status:  ✗ Not responding"
fi

echo ""
echo "==================================="

