#!/bin/bash
set -e

echo "🔑 Creating signer accounts..."

# Известные тестовые приватные ключи (из Hardhat/Ganache)
# Signer 1: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
PRIV_KEY_1="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

# Signer 2: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8  
PRIV_KEY_2="59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

# Signer 3: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
PRIV_KEY_3="5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

# Импортируем аккаунты
echo "Importing signer1 account..."
echo "$PRIV_KEY_1" > /tmp/key1.txt
docker run --rm \
    -v $(pwd)/signer1-init:/data \
    -v $(pwd)/password.txt:/password.txt \
    -v /tmp/key1.txt:/key.txt \
    ethereum/client-go:v1.14.11 \
    account import --datadir=/data --password=/password.txt /key.txt

echo "Importing signer2 account..."
echo "$PRIV_KEY_2" > /tmp/key2.txt
docker run --rm \
    -v $(pwd)/signer2-init:/data \
    -v $(pwd)/password.txt:/password.txt \
    -v /tmp/key2.txt:/key.txt \
    ethereum/client-go:v1.14.11 \
    account import --datadir=/data --password=/password.txt /key.txt

echo "Importing signer3 account..."
echo "$PRIV_KEY_3" > /tmp/key3.txt
docker run --rm \
    -v $(pwd)/signer3-init:/data \
    -v $(pwd)/password.txt:/password.txt \
    -v /tmp/key3.txt:/key.txt \
    ethereum/client-go:v1.14.11 \
    account import --datadir=/data --password=/password.txt /key.txt

# Очищаем временные файлы
rm -f /tmp/key*.txt

echo "✅ All accounts imported!"

