#!/bin/bash
# Проверка здоровья блокчейна

set -e

RPC_URL="${RPC_URL:-http://localhost:8545}"
ALERT_WEBHOOK="${ALERT_WEBHOOK:-}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для отправки алертов
send_alert() {
    local message="$1"
    if [ -n "$ALERT_WEBHOOK" ]; then
        curl -s -X POST "$ALERT_WEBHOOK" \
            -H "Content-Type: application/json" \
            -d "{\"text\":\"⚠️ Blockchain Alert: $message\"}" > /dev/null
    fi
}

# Проверка RPC
check_rpc() {
    echo -n "🔍 Checking RPC connectivity... "
    RESPONSE=$(curl -s -X POST "$RPC_URL" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":1}' || echo "")
    
    if echo "$RESPONSE" | grep -q "Geth"; then
        echo -e "${GREEN}✅ OK${NC}"
        return 0
    else
        echo -e "${RED}❌ FAILED${NC}"
        send_alert "RPC not responding"
        return 1
    fi
}

# Проверка создания блоков
check_block_production() {
    echo -n "⛏️  Checking block production... "
    
    BLOCK1=$(curl -s -X POST "$RPC_URL" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r .result)
    
    sleep 10
    
    BLOCK2=$(curl -s -X POST "$RPC_URL" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r .result)
    
    BLOCK1_DEC=$((16#${BLOCK1#0x}))
    BLOCK2_DEC=$((16#${BLOCK2#0x}))
    
    if [ $BLOCK2_DEC -gt $BLOCK1_DEC ]; then
        echo -e "${GREEN}✅ OK${NC} (Block: $BLOCK1_DEC → $BLOCK2_DEC)"
        return 0
    else
        echo -e "${RED}❌ FAILED${NC} (Block stuck at $BLOCK1_DEC)"
        send_alert "Block production stopped at block $BLOCK1_DEC"
        return 1
    fi
}

# Проверка контейнеров
check_containers() {
    echo "🐳 Checking containers..."
    
    CONTAINERS=(
        "labracodabrador-signer1"
        "labracodabrador-signer2"
        "labracodabrador-signer3"
        "labracodabrador-rpc1"
        "labracodabrador-rpc2"
        "labracodabrador-nginx"
    )
    
    FAILED=0
    for CONTAINER in "${CONTAINERS[@]}"; do
        STATUS=$(docker inspect -f '{{.State.Status}}' "$CONTAINER" 2>/dev/null || echo "not found")
        HEALTH=$(docker inspect -f '{{.State.Health.Status}}' "$CONTAINER" 2>/dev/null || echo "unknown")
        
        echo -n "   $CONTAINER: "
        if [ "$STATUS" = "running" ]; then
            if [ "$HEALTH" = "healthy" ] || [ "$HEALTH" = "unknown" ]; then
                echo -e "${GREEN}✅ OK${NC} ($STATUS)"
            else
                echo -e "${YELLOW}⚠️  WARNING${NC} ($STATUS, health: $HEALTH)"
                FAILED=$((FAILED + 1))
            fi
        else
            echo -e "${RED}❌ FAILED${NC} ($STATUS)"
            send_alert "Container $CONTAINER is $STATUS"
            FAILED=$((FAILED + 1))
        fi
    done
    
    return $FAILED
}

# Проверка дискового пространства
check_disk_space() {
    echo -n "💾 Checking disk space... "
    
    USAGE=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    
    if [ $USAGE -lt 80 ]; then
        echo -e "${GREEN}✅ OK${NC} (${USAGE}% used)"
        return 0
    elif [ $USAGE -lt 90 ]; then
        echo -e "${YELLOW}⚠️  WARNING${NC} (${USAGE}% used)"
        return 0
    else
        echo -e "${RED}❌ CRITICAL${NC} (${USAGE}% used)"
        send_alert "Disk space critical: ${USAGE}% used"
        return 1
    fi
}

# Проверка peer connections
check_peers() {
    echo -n "👥 Checking peer connections... "
    
    PEERS=$(curl -s -X POST "$RPC_URL" \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' | jq -r .result)
    
    PEERS_DEC=$((16#${PEERS#0x}))
    
    if [ $PEERS_DEC -ge 2 ]; then
        echo -e "${GREEN}✅ OK${NC} ($PEERS_DEC peers)"
        return 0
    else
        echo -e "${YELLOW}⚠️  WARNING${NC} ($PEERS_DEC peers)"
        return 0
    fi
}

# Основная проверка
echo "================================"
echo "🏥 Blockchain Health Check"
echo "================================"
echo ""

FAILED=0

check_rpc || FAILED=$((FAILED + 1))
check_block_production || FAILED=$((FAILED + 1))
check_containers || FAILED=$((FAILED + 1))
check_peers || FAILED=$((FAILED + 1))
check_disk_space || FAILED=$((FAILED + 1))

echo ""
echo "================================"
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED check(s) failed!${NC}"
    exit 1
fi

