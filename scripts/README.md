# 🛠️ Утилиты для управления блокчейном

Этот каталог содержит утилитные скрипты для управления production блокчейном.

## 📦 Доступные скрипты

### 1. backup-blockchain.sh

**Назначение:** Создание бэкапа всех Docker volumes блокчейна.

**Использование:**
```bash
./scripts/backup-blockchain.sh
```

**Что делает:**
- Создает архив каждой signer ноды (signer1-data, signer2-data, signer3-data)
- Сохраняет конфигурацию (docker-compose.yml, genesis-poa.json, nginx.conf)
- Создает единый tar.gz архив
- Удаляет бэкапы старше 7 дней

**Бэкапы сохраняются в:** `/backup/ethereum/blockchain-backup-TIMESTAMP.tar.gz`

---

### 2. restore-blockchain.sh

**Назначение:** Восстановление блокчейна из бэкапа.

**Использование:**
```bash
./scripts/restore-blockchain.sh /backup/ethereum/blockchain-backup-20251008_123456.tar.gz
```

**Что делает:**
- Останавливает все контейнеры
- Удаляет существующие Docker volumes
- Восстанавливает данные из бэкапа
- Запускает контейнеры заново

**⚠️ Внимание:** Это удалит все текущие данные блокчейна!

---

### 3. health-check.sh

**Назначение:** Проверка здоровья блокчейна и всех компонентов.

**Использование:**
```bash
./scripts/health-check.sh
```

**Что проверяет:**
- ✅ RPC connectivity
- ✅ Создание блоков
- ✅ Статус Docker контейнеров
- ✅ Peer connections
- ✅ Свободное место на диске

**Алерты:** Может отправлять webhook алерты при проблемах (настраивается через `ALERT_WEBHOOK`)

**Пример с алертами:**
```bash
ALERT_WEBHOOK="https://hooks.slack.com/..." ./scripts/health-check.sh
```

---

## 🚀 Использование через Makefile

Все скрипты доступны через удобные Make команды:

```bash
make backup          # Бэкап блокчейна
make restore BACKUP=/path/to/backup.tar.gz  # Восстановление
make health          # Проверка здоровья
```

## 📅 Автоматизация

### Автоматический бэкап (cron)

Добавьте в crontab для ежедневного бэкапа в 2:00:

```bash
0 2 * * * /path/to/labracodabrador/scripts/backup-blockchain.sh >> /var/log/blockchain-backup.log 2>&1
```

### Автоматический health check

Проверка каждые 5 минут:

```bash
*/5 * * * * /path/to/labracodabrador/scripts/health-check.sh >> /var/log/blockchain-health.log 2>&1
```

## 🔧 Переменные окружения

### backup-blockchain.sh
- `BACKUP_ROOT` - директория для бэкапов (по умолчанию: `/backup/ethereum`)
- `RETENTION_DAYS` - количество дней хранения (по умолчанию: `7`)

### health-check.sh
- `RPC_URL` - URL RPC endpoint (по умолчанию: `http://localhost:8545`)
- `ALERT_WEBHOOK` - webhook для алертов (опционально)

### restore-blockchain.sh
- Не требует переменных окружения

---

**💡 Совет:** Используйте эти скрипты в production для обеспечения надежности и быстрого восстановления!
