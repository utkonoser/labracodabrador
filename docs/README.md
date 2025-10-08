# 🚀 Labracodabrador - Production Ethereum PoA Blockchain

> Отказоустойчивый Ethereum PoA (Proof of Authority) блокчейн на Geth v1.13.15 с полной production инфраструктурой

## ✨ Что это?

Готовая к production частная Ethereum сеть с:
- **3 Signer ноды** - создают блоки через Clique PoA
- **2 RPC ноды** - обрабатывают запросы
- **Nginx Load Balancer** - распределяет RPC запросы  
- **Prometheus + Grafana** - мониторинг метрик и дашборды
- **Loki + Promtail** - централизованное логирование
- **Health checks** - автоматический перезапуск

## 🏗️ Архитектура

```
                    Metamask / DApp
                           │
                           ▼
            ┌──────────────────────────┐
            │  Nginx (Load Balancer)   │ :8545/:8546
            │    http://localhost      │
            └──────────┬───────────────┘
                       │
            ┌──────────┴──────────┐
            ▼          ▼          ▼
        ┌────────┬────────┬────────┐
        │Signer 1│Signer 2│Signer 3│  ← PoA Mining (блоки каждые 5с)
        └───┬────┴────┬───┴────┬───┘
            │         │        │
        ┌───┴─────────┴────────┴───┐
        │   Clique PoA Network     │
        │   Chain ID: 32382        │
        └──────────────────────────┘
                       │
            ┌──────────┴──────────┐
            ▼                     ▼
     ┌────────────┐         ┌─────────┐
     │ Prometheus │         │ Grafana │  ← Мониторинг
     └────────────┘         └─────────┘
            │                     │
            ▼                     ▼
     ┌────────────┐         ┌─────────┐
     │    Loki    │         │Promtail │  ← Логирование
     └────────────┘         └─────────┘
```

## 🚀 Быстрый старт

### 1. Первоначальная настройка

```bash
# Инициализация genesis блока и создание аккаунтов
./setup-simple.sh
./create-accounts.sh
```

### 2. Запуск сети

```bash
# Запустить все сервисы
docker-compose up -d

# Проверить статус
docker-compose ps

# Посмотреть логи
docker-compose logs -f
```

### 3. Проверка работы

```bash
# Текущий блок
curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Chain ID
curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'
```

## 🦊 Подключение Metamask

1. Откройте Metamask → **Добавить сеть**
2. Заполните:

```
Network Name:  Labracodabrador PoA
RPC URL:       http://localhost:8545
Chain ID:      32382
Currency:      ETH
```

3. Готово! ✅

### Тестовые аккаунты (с балансом)

Импортируйте в Metamask:

```
0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
0x70997970C51812dc3A010C7d01b50e0d17dc79C8
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
```

**Private keys** (только для testnet!):
```
Аккаунт 1: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
Аккаунт 2: 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
Аккаунт 3: 5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
```

## 🔍 Web Explorer & REST API

Полноценный blockchain explorer с REST API (работает на любой платформе):

```bash
# Запускается автоматически вместе с основной сетью
docker-compose up -d

# Web Explorer (UI)
open http://localhost:8080

# REST API
curl http://localhost:8081/api/v1/network
```

### 🌐 **Web Explorer (UI):**
- 🔍 **Поиск** - адреса, транзакции, блоки
- 💰 **Балансы** - проверка ETH балансов
- 📦 **Блоки** - последние блоки в реальном времени
- 📝 **Транзакции** - статус, gas, value
- ⚡ **Автообновление** - каждые 5-10 секунд

### 🔗 **REST API Endpoints:**
- `GET /api/v1/network` - Информация о сети (Chain ID, Latest Block, Gas Price)
- `GET /api/v1/blocks/latest` - Последние 6 блоков
- `GET /api/v1/blocks/{number}` - Детали конкретного блока
- `GET /api/v1/address/{address}` - Информация об адресе (баланс, tx count, тип)
- `GET /api/v1/transactions/{hash}` - Детали транзакции (статус, gas, block)
- `GET /health` - Health check API

### ✨ **Преимущества:**
- ✅ **Кроссплатформенность** - работает на Mac, Linux, Windows
- ✅ **REST API** - легко интегрируется с любыми приложениями
- ✅ **Без БД** - не требует PostgreSQL или других баз данных
- ✅ **Легковесный** - Web UI ~100KB, API сервер ~10MB
- ✅ **Стабильный** - использует REST API вместо прямого RPC

💡 **Примечание:** Explorer работает через REST API, что обеспечивает стабильность и возможность интеграции с внешними приложениями.

---

## 📊 Мониторинг и логирование

- **Prometheus**: http://localhost:9090 - сбор метрик
- **Grafana**: http://localhost:3000 (admin/admin) - дашборды и визуализация
- **Loki**: http://localhost:3100 - централизованное логирование

### Доступные дашборды:
- **Blockchain Overview** - общий обзор сети
- **Node Details** - детали по узлам  
- **Logs Overview** - просмотр логов всех сервисов

**Метрики нод:**
- Signer1: http://172.25.0.11:6060/debug/metrics
- Signer2: http://172.25.0.12:6060/debug/metrics  
- Signer3: http://172.25.0.13:6060/debug/metrics

Подробнее: [MONITORING.md](MONITORING.md) | [LOGGING.md](LOGGING.md)

## 🔧 Управление

### Основные команды

```bash
# Статус
docker-compose ps

# Логи всех сервисов
docker-compose logs -f

# Логи конкретной ноды
docker-compose logs -f signer1

# Перезапуск
docker-compose restart

# Остановка
docker-compose down

# Полная очистка (удалит все данные блокчейна!)
docker-compose down -v
```

### API запросы

**Получить блок:**
```bash
curl http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

**Баланс аккаунта:**
```bash
curl http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","latest"],"id":1}'
```

**Отправить транзакцию:**
```bash
curl http://localhost:8545 -X POST -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"eth_sendTransaction",
    "params":[{
      "from":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      "to":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "value":"0xDE0B6B3A7640000",
      "gas":"0x5208"
    }],
    "id":1
  }'
```

## 💻 Использование в коде

### JavaScript (ethers.js)

```javascript
const { ethers } = require("ethers");

const provider = new ethers.JsonRpcProvider("http://localhost:8545");
const wallet = new ethers.Wallet("PRIVATE_KEY", provider);

// Получить блок
const blockNumber = await provider.getBlockNumber();
console.log("Block:", blockNumber);

// Отправить транзакцию
const tx = await wallet.sendTransaction({
  to: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
  value: ethers.parseEther("1.0")
});

const receipt = await tx.wait();
console.log("TX:", receipt.hash);
```

### Python (web3.py)

```python
from web3 import Web3

w3 = Web3(Web3.HTTPProvider('http://localhost:8545'))

# Проверить подключение
print(w3.is_connected())  # True

# Получить блок
block_number = w3.eth.block_number
print(f"Block: {block_number}")

# Отправить транзакцию
account = w3.eth.account.from_key('PRIVATE_KEY')
tx = {
    'from': account.address,
    'to': '0x70997970C51812dc3A010C7d01b50e0d17dc79C8',
    'value': w3.to_wei(1, 'ether'),
    'gas': 21000,
    'gasPrice': w3.eth.gas_price,
    'nonce': w3.eth.get_transaction_count(account.address),
    'chainId': 32382
}

signed = account.sign_transaction(tx)
tx_hash = w3.eth.send_raw_transaction(signed.raw_transaction)
print(f"TX: {tx_hash.hex()}")
```

### Go

```go
package main

import (
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    client, _ := ethclient.Dial("http://localhost:8545")
    
    // Получить блок
    blockNumber, _ := client.BlockNumber(context.Background())
    fmt.Printf("Block: %d\n", blockNumber)
    
    // Баланс
    balance, _ := client.BalanceAt(
        context.Background(),
        common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
        nil,
    )
    fmt.Printf("Balance: %s ETH\n", balance)
}
```

## 🏗️ Конфигурация

### Основные параметры

- **Chain ID**: 32382
- **Network ID**: 32382
- **Consensus**: Clique PoA
- **Block time**: 5 секунд
- **Epoch**: 30000 блоков

### Порты

- `8545` - HTTP JSON-RPC (через Nginx)
- `8546` - WebSocket (через Nginx)
- `9090` - Prometheus
- `3000` - Grafana
- `30303` - P2P (каждая нода)

### Доступные API методы

**Через HTTP/WS:**
- `eth_*` - Ethereum methods
- `net_*` - Network methods
- `web3_*` - Web3 methods
- `txpool_*` - Transaction pool
- `debug_*` - Debug methods
- `admin_*` - Admin methods (только на нодах напрямую)

## 🛠️ Troubleshooting

### Порт 8545 занят

```bash
# Найти процесс
lsof -i :8545

# Или изменить порт в docker-compose.yml
ports:
  - "9545:8545"
```

### Блоки не создаются

```bash
# Проверить что signers работают
docker logs labracodabrador-signer1 | grep "Successfully sealed"

# Проверить количество peers
curl http://172.25.0.11:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}'
```

### Nginx не отвечает

```bash
# Проверить nginx
docker logs labracodabrador-nginx

# Проверить здоровье
curl http://localhost:8545/health
```

### Ноды не подключаются

```bash
# Проверить сеть
docker network inspect labracodabrador_blockchain-network

# Перезапустить сеть
docker-compose down
docker-compose up -d
```

## 🔒 Security (Production)

**⚠️ Перед выходом в production:**

### 1. Смените пароли

```bash
echo "YOUR_STRONG_PASSWORD" > password.txt
chmod 600 password.txt
```

### 2. Ограничьте доступ к RPC

В `docker-compose.yml`:
```yaml
- --http.corsdomain=https://yourdomain.com  # вместо *
- --http.vhosts=yourdomain.com              # вместо *
```

### 3. Используйте SSL/TLS

Настройте nginx с SSL сертификатом (Let's Encrypt).

### 4. Firewall

```bash
# Только нужные порты
ufw allow 8545/tcp  # RPC
ufw allow 8546/tcp  # WS
ufw allow 30303     # P2P (только между нодами)
```

### 5. Мониторинг и Alerts

Настройте Prometheus alerting rules для:
- Падение нод
- Отставание блоков
- Проблемы с consensus

### 6. Backups

```bash
# Регулярное резервное копирование
docker run --rm -v labracodabrador_signer1-data:/data \
  -v $(pwd)/backup:/backup alpine \
  tar czf /backup/signer1-$(date +%Y%m%d).tar.gz /data
```

## 📁 Структура проекта

```
labracodabrador/
├── docker-compose.yml      # Главная конфигурация
├── genesis-poa.json        # Genesis блок (Clique PoA)
├── nginx.conf              # Load balancer config
├── prometheus.yml          # Metrics config
├── password.txt            # Пароль для signer аккаунтов
├── setup-simple.sh         # Инициализация genesis
├── create-accounts.sh      # Создание signer аккаунтов
├── Makefile                # Make команды для Go проекта
├── cmd/                    # Go source code (опционально)
├── internal/               # Internal packages
└── README.md               # Документация
```

## 📚 Технологии

- **Geth**: v1.13.15 (последняя версия с PoA Clique)
- **Go**: 1.23
- **Docker**: Latest
- **Nginx**: Alpine
- **Prometheus**: Latest
- **Grafana**: Latest

## 🤝 Полезные ссылки

- [Geth Documentation](https://geth.ethereum.org/docs)
- [Clique PoA Spec](https://eips.ethereum.org/EIPS/eip-225)
- [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Metamask](https://metamask.io/)
- [Ethers.js](https://docs.ethers.org/)
- [Web3.py](https://web3py.readthedocs.io/)

## 📄 Лицензия

MIT License

