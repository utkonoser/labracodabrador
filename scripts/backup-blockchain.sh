#!/bin/bash
# Скрипт для бэкапа блокчейна

set -e

BACKUP_ROOT="/backup/ethereum"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="$BACKUP_ROOT/$TIMESTAMP"
RETENTION_DAYS=7

echo "🔄 Starting blockchain backup at $(date)"

# Создать директорию для бэкапа
mkdir -p "$BACKUP_DIR"

# Бэкап каждой signer ноды
for NODE in signer1 signer2 signer3; do
    echo "📦 Backing up $NODE..."
    
    # Бэкап chaindata через docker volume
    docker run --rm \
        -v labracodabrador_${NODE}-data:/data:ro \
        -v "$BACKUP_DIR:/backup" \
        alpine \
        tar -czf /backup/${NODE}-data.tar.gz -C /data .
    
    echo "✅ $NODE backup completed"
done

# Бэкап keystore (если есть)
if [ -d "/secure/ethereum/keystore" ]; then
    echo "🔑 Backing up keystore..."
    tar -czf "$BACKUP_DIR/keystore.tar.gz" -C /secure/ethereum keystore
fi

# Бэкап конфигурации
echo "⚙️  Backing up configuration..."
cp docker-compose.yml "$BACKUP_DIR/"
cp genesis-poa.json "$BACKUP_DIR/"
cp nginx.conf "$BACKUP_DIR/"

# Создать единый архив
echo "📚 Creating final archive..."
tar -czf "$BACKUP_ROOT/blockchain-backup-$TIMESTAMP.tar.gz" -C "$BACKUP_DIR" .
rm -rf "$BACKUP_DIR"

# Размер бэкапа
BACKUP_SIZE=$(du -h "$BACKUP_ROOT/blockchain-backup-$TIMESTAMP.tar.gz" | cut -f1)
echo "💾 Backup size: $BACKUP_SIZE"

# Очистка старых бэкапов
echo "🧹 Cleaning old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_ROOT" -name "blockchain-backup-*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "✅ Backup completed: blockchain-backup-$TIMESTAMP.tar.gz"
echo "📊 Available backups:"
ls -lh "$BACKUP_ROOT"/blockchain-backup-*.tar.gz 2>/dev/null || echo "No backups found"

