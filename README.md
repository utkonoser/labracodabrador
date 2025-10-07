# 🚀 Labracodabrador - Production Ethereum PoA Blockchain

> Отказоустойчивый Ethereum PoA (Proof of Authority) блокчейн на Geth v1.13.15 с полной production инфраструктурой

[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)
[![Ethereum](https://img.shields.io/badge/Ethereum-3C3C3D?style=flat&logo=Ethereum&logoColor=white)](https://ethereum.org/)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)](https://golang.org/)

## ✨ Что это?

Готовая к production частная Ethereum сеть с:
- **3 Signer ноды** - создают блоки через Clique PoA
- **2 RPC ноды** - обрабатывают запросы
- **Nginx Load Balancer** - распределяет RPC запросы  
- **REST API** - удобный доступ к блокчейну
- **Web Explorer** - веб-интерфейс для просмотра блокчейна
- **Prometheus + Grafana** - мониторинг метрик
- **Health checks** - автоматический перезапуск

## 🚀 Быстрый старт

```bash
# 1. Инициализация
./scripts/setup-simple.sh
./scripts/create-accounts.sh

# 2. Запуск
make run

# 3. Проверка
make ps
make logs
```

**Готово!** 🎉

- 🌐 **Web Explorer:** http://localhost:8080
- 🔗 **REST API:** http://localhost:8081/api/v1
- ⚡ **RPC:** http://localhost:8545
- 📊 **Grafana:** http://localhost:3000
- 📈 **Prometheus:** http://localhost:9090

## 📚 Документация

Подробная документация находится в папке `docs/`:

- [**README.md**](docs/README.md) - Полная документация проекта
- [**QUICK-START-PRODUCTION.md**](docs/QUICK-START-PRODUCTION.md) - Production чеклист

## 🛠️ Команды Makefile

```bash
make help           # Показать все команды
make build          # Собрать API сервер
make run            # Запустить все сервисы
make stop           # Остановить сервисы
make restart        # Перезапустить
make logs           # Показать логи
make logs-api       # Логи API сервера
make logs-signer    # Логи signer нод
make ps             # Статус контейнеров
make backup         # Бэкап блокчейна
make restore        # Восстановление
make health         # Проверка здоровья
make clean          # Удалить все данные
```

## 📁 Структура проекта

```
labracodabrador/
├── config/                     # Конфигурации
│   ├── docker-compose.yml      # Docker orchestration
│   ├── genesis-poa.json        # Genesis block
│   ├── nginx.conf              # Load balancer
│   └── prometheus.yml          # Metrics
│
├── docker/                     # Docker файлы
│   └── Dockerfile.api          # API сервер
│
├── scripts/                    # Утилиты
│   ├── setup-simple.sh         # Инициализация
│   ├── create-accounts.sh      # Создание аккаунтов
│   ├── backup-blockchain.sh    # Бэкап
│   ├── restore-blockchain.sh   # Восстановление
│   └── health-check.sh         # Health check
│
├── cmd/                        # Go приложения
│   └── api/                    # REST API сервер
│
├── web-explorer/               # Web UI
│   └── index.html              # Веб интерфейс
│
├── docs/                       # Документация
│   ├── README.md               # Полная документация
│   └── QUICK-START-PRODUCTION.md
│
└── [Makefile, go.mod, go.sum, .gitignore]
```

## 🦊 Metamask

```
Network Name:  Labracodabrador PoA
RPC URL:       http://localhost:8545
Chain ID:      32382
Currency:      ETH
```

**Тестовые аккаунты:**
```
0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
0x70997970C51812dc3A010C7d01b50e0d17dc79C8
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
```

## 🔒 Security

⚠️ **Важно для production:**
1. Смените пароли в `password.txt`
2. Ограничьте CORS в конфигурации
3. Настройте SSL/TLS
4. Настройте firewall
5. Включите регулярные бэкапы

См. [Production Guide](docs/QUICK-START-PRODUCTION.md)

## 📊 Мониторинг

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

## 🤝 Полезные ссылки

- [Geth Documentation](https://geth.ethereum.org/docs)
- [Clique PoA Spec](https://eips.ethereum.org/EIPS/eip-225)
- [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)

## 📄 Лицензия

MIT License

