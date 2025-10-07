#!/bin/bash
# Скрипт для восстановления блокчейна из бэкапа

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <backup-file.tar.gz>"
    echo ""
    echo "Available backups:"
    ls -lh /backup/ethereum/blockchain-backup-*.tar.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"
RESTORE_DIR="/tmp/blockchain-restore-$$"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "⚠️  WARNING: This will stop all blockchain containers and restore from backup!"
echo "Backup file: $BACKUP_FILE"
read -p "Continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

# Остановить контейнеры
echo "🛑 Stopping containers..."
cd "$(dirname "$0")/.." && docker-compose -f config/docker-compose.yml down

# Распаковать бэкап
echo "📦 Extracting backup..."
mkdir -p "$RESTORE_DIR"
tar -xzf "$BACKUP_FILE" -C "$RESTORE_DIR"

# Восстановить каждую ноду
for NODE in signer1 signer2 signer3; do
    if [ -f "$RESTORE_DIR/${NODE}-data.tar.gz" ]; then
        echo "🔄 Restoring $NODE..."
        
        # Удалить старый volume
        docker volume rm labracodabrador_${NODE}-data 2>/dev/null || true
        
        # Создать новый volume
        docker volume create labracodabrador_${NODE}-data
        
        # Восстановить данные
        docker run --rm \
            -v labracodabrador_${NODE}-data:/data \
            -v "$RESTORE_DIR:/backup" \
            alpine \
            sh -c "cd /data && tar -xzf /backup/${NODE}-data.tar.gz"
        
        echo "✅ $NODE restored"
    fi
done

# Восстановить keystore
if [ -f "$RESTORE_DIR/keystore.tar.gz" ]; then
    echo "🔑 Restoring keystore..."
    mkdir -p /secure/ethereum
    tar -xzf "$RESTORE_DIR/keystore.tar.gz" -C /secure/ethereum
fi

# Очистка
rm -rf "$RESTORE_DIR"

# Запустить контейнеры
echo "🚀 Starting containers..."
cd "$(dirname "$0")/.." && docker-compose -f config/docker-compose.yml up -d

echo "✅ Restore completed!"
echo "📊 Check status with: make ps"
echo "📋 Check logs with: make logs"

