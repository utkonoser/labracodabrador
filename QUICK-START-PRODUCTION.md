# 🚀 Quick Start: Production Deployment

## Краткая шпаргалка для вывода в продакшн

### 📋 Что уже работает (не требует изменений)

✅ **Хранение транзакций** - Geth автоматически сохраняет все транзакции и блоки в базу данных `/data/geth/chaindata`

✅ **Отказоустойчивость** - 3 signer ноды с Clique PoA консенсусом

✅ **Load Balancing** - Nginx распределяет RPC запросы

✅ **Мониторинг** - Prometheus + Grafana готовы

### 🔥 Критические изменения для продакшна

#### 1. Безопасность (ОБЯЗАТЕЛЬНО!)

```bash
# Переместите keystore за пределы репозитория
mkdir -p /secure/ethereum/keystore
mv signer*-init/keystore/* /secure/ethereum/keystore/

# Обновите docker-compose.yml:
# volumes:
#   - /secure/ethereum/keystore/signer1:/data/keystore

# Добавьте в .gitignore (уже добавлено):
# password.txt
# signer*-init/keystore/*
```

#### 2. Настройте бэкапы

```bash
# Добавьте в crontab:
0 3 * * * /opt/labracodabrador/scripts/backup-blockchain.sh >> /var/log/blockchain-backup.log 2>&1

# Создайте директорию для бэкапов:
mkdir -p /backup/ethereum
```

#### 3. Health Check

```bash
# Добавьте в crontab для регулярной проверки:
*/5 * * * * /opt/labracodabrador/scripts/health-check.sh || echo "Health check failed"
```

### 🛠️ Полезные команды

```bash
# Проверка здоровья
./scripts/health-check.sh

# Бэкап
./scripts/backup-blockchain.sh

# Восстановление
./scripts/restore-blockchain.sh /backup/ethereum/blockchain-backup-TIMESTAMP.tar.gz

# Просмотр логов
docker-compose logs -f signer1

# Статус всех нод
docker-compose ps

# Текущий блок
curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq

# Количество пиров
curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' | jq
```

### 🌐 Production конфигурация

#### Минимальные требования к серверу

- **CPU:** 4 cores (рекомендуется 8)
- **RAM:** 8 GB (рекомендуется 16 GB)
- **Disk:** 100 GB SSD (рекомендуется 500 GB+)
- **Network:** 1 Gbps

#### Настройка firewall

```bash
# UFW (Ubuntu)
ufw allow 22/tcp    # SSH
ufw allow 8545/tcp  # RPC (ограничьте IP)
ufw allow 30303     # P2P
ufw enable

# iptables (альтернатива)
iptables -A INPUT -p tcp --dport 8545 -s YOUR_TRUSTED_IP -j ACCEPT
iptables -A INPUT -p tcp --dport 8545 -j DROP
```

#### SSL/TLS для RPC (рекомендуется)

```bash
# Получите Let's Encrypt сертификат
certbot certonly --standalone -d your-domain.com

# Обновите nginx.conf для HTTPS
# (см. PRODUCTION.md)
```

### 📊 Мониторинг

- **Grafana:** http://your-domain:3000 (admin/admin)
- **Prometheus:** http://your-domain:9090
- **Metrics endpoint:** http://localhost:6060/debug/metrics

### 🔍 Troubleshooting

#### Проблема: Ноды не создают блоки

```bash
# Проверьте логи
docker-compose logs -f signer1 | grep -i error

# Проверьте пиры
curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"admin_peers","params":[],"id":1}' | jq
```

#### Проблема: RPC не отвечает

```bash
# Проверьте статус контейнеров
docker-compose ps

# Перезапустите RPC ноды
docker-compose restart rpc1 rpc2

# Проверьте nginx
docker-compose logs nginx
```

#### Проблема: Нехватка места на диске

```bash
# Проверьте использование
df -h

# Очистите старые логи
docker-compose logs --tail=1000 > /dev/null

# Удалите старые бэкапы
find /backup/ethereum -name "*.tar.gz" -mtime +30 -delete
```

### 📝 Чек-лист перед запуском в продакшн

- [ ] Keystore файлы защищены (не в Git!)
- [ ] Пароли безопасные и уникальные
- [ ] Firewall настроен
- [ ] SSL/TLS для RPC (если доступ извне)
- [ ] Бэкапы настроены и протестированы
- [ ] Мониторинг работает
- [ ] Health checks настроены
- [ ] Документация обновлена
- [ ] Disaster recovery план готов
- [ ] Алерты настроены (Slack/Email)
- [ ] Тестирование под нагрузкой выполнено

### 🎯 Следующие шаги

1. Прочитайте полную документацию: **[PRODUCTION.md](./PRODUCTION.md)**
2. Настройте безопасность согласно чек-листу
3. Запустите health-check для проверки
4. Настройте бэкапы
5. Настройте мониторинг и алерты
6. Проведите нагрузочное тестирование
7. Подготовьте disaster recovery план

### 📚 Дополнительные ресурсы

- [PRODUCTION.md](./PRODUCTION.md) - Полное руководство по продакшну
- [README.md](./README.md) - Основная документация
- [Geth Documentation](https://geth.ethereum.org/docs)
- [Clique PoA](https://geth.ethereum.org/docs/fundamentals/consensus)

---

## 💡 Важные моменты

### О хранении транзакций

**Geth автоматически хранит:**
- ✅ Все блоки в `/data/geth/chaindata/ancient/chain`
- ✅ Все транзакции в каждом блоке
- ✅ State database (балансы аккаунтов)
- ✅ Transaction index для быстрого поиска

**Вам НЕ нужно:**
- ❌ Создавать отдельную базу данных для транзакций
- ❌ Писать код для сохранения транзакций
- ❌ Управлять хранением вручную

**Что вы МОЖЕТЕ сделать дополнительно:**
- ✅ Настроить BlockScout для web explorer
- ✅ Экспортировать транзакции в аналитику
- ✅ Создать индексатор для специфичных данных

### Параметры хранения

```yaml
# В docker-compose.yml можно настроить:
--gcmode=archive      # Хранить всю историю (рекомендуется)
--txlookuplimit=0     # Индексировать все транзакции
--cache=2048          # Увеличить кеш для производительности
```

### Размер базы данных

Примерные размеры для планирования:
- **Новая сеть:** ~100 MB
- **Месяц работы (5 TPS):** ~1-5 GB
- **Год работы (10 TPS):** ~50-100 GB

Регулярно мониторьте: `df -h /var/lib/docker/volumes/`

