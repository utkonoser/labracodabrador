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
- **Prometheus + Grafana** - мониторинг метрик и дашборды
- **Loki + Promtail** - централизованное логирование
- **Health checks** - автоматический перезапуск
- **GitHub Actions** - автоматический деплой

## 🏗️ Архитектура

```
                                    ┌─────────────────┐
                                    │   Web Explorer  │
                                    │   (Nginx:9080)  │
                                    └────────┬────────┘
                                             │
                   ┌─────────────────────────┼─────────────────────────┐
                   │                         │                         │
          ┌────────▼────────┐       ┌───────▼───────┐       ┌────────▼────────┐
          │   REST API      │       │   Grafana     │       │   Prometheus    │
          │   (Go:9081)     │       │   (4000)      │       │   (10090)       │
          └────────┬────────┘       └───────┬───────┘       └────────┬────────┘
                   │                        │                         │
          ┌────────▼────────┐               │                ┌────────▼────────┐
          │  Nginx Load     │               │                │   Loki          │
          │  Balancer:9545  │               │                │   (4100)        │
          └────────┬────────┘               │                └────────▲────────┘
                   │                        │                         │
        ┌──────────┴──────────┐             │                ┌────────┴────────┐
        │                     │             │                │   Promtail      │
┌───────▼───────┐     ┌──────▼──────┐      │                └─────────────────┘
│   RPC Node 1  │     │ RPC Node 2  │      │
│ (172.25.0.21) │     │(172.25.0.22)│      │
└───────┬───────┘     └──────┬──────┘      │
        │                    │              │
        └──────────┬─────────┘              │
                   │                        │
        ┌──────────┴──────────┐             │
        │                     │             │
┌───────▼───────┐  ┌─────────▼────────┐  ┌─▼──────────────┐
│ Signer Node 1 │  │  Signer Node 2   │  │ Signer Node 3  │
│(172.25.0.11)  │  │  (172.25.0.12)   │  │ (172.25.0.13)  │
│   Майнер #1   │  │    Майнер #2     │  │   Майнер #3    │
└───────────────┘  └──────────────────┘  └────────────────┘
```

### Компоненты

**Блокчейн уровень:**
- **3 Signer ноды** - валидаторы сети, создают блоки по алгоритму Clique PoA (Proof of Authority)
- **2 RPC ноды** - обрабатывают внешние запросы, не участвуют в консенсусе, синхронизируются с сетью

**Сетевой уровень:**
- **Nginx Load Balancer** - распределяет нагрузку между RPC нодами, обеспечивает отказоустойчивость
- **Docker Network** - изолированная сеть 172.25.0.0/16 для всех компонентов

**API уровень:**
- **REST API сервер (Go)** - предоставляет удобный HTTP API для работы с блокчейном
- **Web Explorer** - веб-интерфейс для просмотра блоков, транзакций и адресов

**Мониторинг и логирование:**
- **Prometheus** - сбор метрик со всех нод (block height, gas price, peer count и т.д.)
- **Grafana** - визуализация метрик и создание дашбордов
- **Loki** - централизованное хранилище логов
- **Promtail** - агент для сбора логов со всех контейнеров

### Особенности архитектуры

✅ **Высокая доступность** - отказ одной RPC ноды не влияет на доступность сети  
✅ **Отказоустойчивость** - минимум 2 из 3 signer нод должны работать для создания блоков  
✅ **Масштабируемость** - легко добавить больше RPC нод для обработки большей нагрузки  
✅ **Мониторинг** - полная наблюдаемость всех компонентов системы  
✅ **Безопасность** - изолированная сеть, контроль доступа через Nginx  

## 🛠️ Технологический стек

### Блокчейн
- **Geth v1.13.15** - официальная реализация Ethereum на Go
- **Clique PoA** - консенсус Proof of Authority для приватных сетей
- **Chain ID: 32382** - уникальный идентификатор сети

### Backend
- **Go 1.23** - язык программирования для API сервера
- **go-ethereum** - официальная библиотека Ethereum для Go
- **Gorilla Mux** - HTTP router для REST API

### Инфраструктура
- **Docker & Docker Compose** - контейнеризация и оркестрация
- **Nginx** - load balancing и reverse proxy
- **Linux** - операционная система (работает на любом дистрибутиве)

### Мониторинг и логирование
- **Prometheus** - time-series база данных для метрик
- **Grafana** - визуализация метрик и дашборды
- **Loki** - система логирования от Grafana Labs
- **Promtail** - агент для сбора логов

### DevOps
- **GitHub Actions** - CI/CD пайплайны для автоматического деплоя
- **Makefile** - автоматизация команд и упрощение управления
- **Bash scripts** - утилиты для бэкапа, восстановления и health checks

### Сеть
- **JSON-RPC** - стандартный протокол Ethereum для взаимодействия с нодами
- **WebSocket** - для real-time подписок на события блокчейна
- **HTTP REST API** - удобный API для веб-приложений

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

- 🌐 **Web Explorer:** http://localhost:9080
- 🔗 **REST API:** http://localhost:9081/api/v1
- ⚡ **RPC:** http://localhost:9545
- 📊 **Grafana:** http://localhost:4000 (admin/admin)
- 📈 **Prometheus:** http://localhost:10090
- 📝 **Loki:** http://localhost:4100

> 💡 **Порты:** Все порты подняты на +1000 для избежания конфликтов с другими приложениями на сервере

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
│   ├── prometheus.yml          # Metrics
│   ├── prometheus-alerts.yml   # Alert rules
│   ├── loki.yml                # Log aggregation
│   └── promtail.yml            # Log collection
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
├── grafana/                    # Дашборды
│   └── dashboards/             # Grafana дашборды
│
├── .github/                    # GitHub Actions
│   └── workflows/              # CI/CD пайплайны
│
└── [Makefile, go.mod, go.sum, .gitignore]
```

## 🦊 MetaMask

```
Network Name:  Labracodabrador PoA
RPC URL:       http://your-server-ip:9545
Chain ID:      32382
Currency:      ETH
Block Explorer: http://your-server-ip:9080
```

**Тестовые аккаунты:**
```
0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
0x70997970C51812dc3A010C7d01b50e0d17dc79C8
0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
```

> ⚠️ **Важно:** Отключите Smart Transactions в MetaMask для корректной работы с приватными сетями!

## 🔒 Security

⚠️ **Важно для production:**

1. **Смените пароли** в `password.txt`
2. **Ограничьте CORS** в конфигурации
3. **Настройте SSL/TLS** для HTTPS
4. **Настройте firewall** - откройте только нужные порты
5. **Включите регулярные бэкапы**
6. **Используйте сильные пароли** для keystore файлов
7. **Ограничьте доступ** к RPC endpoints

### Порты и безопасность:

| Порт | Сервис | Доступ | Рекомендация |
|------|--------|--------|--------------|
| 9080 | Web Explorer | Публичный | Ограничить по IP |
| 9081 | REST API | Публичный | Ограничить по IP |
| 9545 | RPC | Публичный | **Ограничить обязательно!** |
| 4000 | Grafana | Админ | VPN или localhost |
| 10090 | Prometheus | Админ | VPN или localhost |

## 📊 Мониторинг и логирование

### Prometheus (http://localhost:10090)
- Сбор метрик с всех узлов
- Алерты при проблемах
- Исторические данные

### Grafana (http://localhost:4000)
- **Логин:** admin / admin
- **Blockchain Overview** - общий обзор сети
- **Node Details** - детали по узлам  
- **Logs Overview** - просмотр логов всех сервисов

### Loki (http://localhost:4100)
- Централизованное логирование
- Поиск по логам
- Интеграция с Grafana

## 🚀 Автоматический деплой

Настроен GitHub Actions для автоматического деплоя при пуше в `main`:

1. **Тесты** - запуск тестов Go
2. **Сборка** - сборка проекта
3. **Деплой** - автоматический деплой на сервер

### Настройка деплоя:

1. Добавьте Secrets в GitHub:
   - `SERVER_HOST` - IP сервера
   - `SERVER_USER` - пользователь SSH
   - `SSH_PRIVATE_KEY` - приватный SSH ключ

2. На сервере клонируйте репозиторий:
   ```bash
   cd /opt
   git clone <your-repo-url> labracodabrador
   ```

3. Готово! При пуше в `main` код автоматически задеплоится.

## 🔧 Troubleshooting

### Проблемы с MetaMask:

1. **Отключите Smart Transactions** в настройках MetaMask
2. **Проверьте Chain ID** - должен быть `32382`
3. **Увеличьте Gas Price** до 20 Gwei
4. **Увеличьте Gas Limit** до 300000

### Проблемы с подтверждениями:

- В PoA сети блоки создаются каждые 60 секунд
- Подождите 1-2 минуты для подтверждения
- Проверьте что signer ноды работают: `make ps`

### Проблемы с мониторингом:

```bash
# Проверьте статус
make health

# Перезапустите Grafana
docker-compose -f config/docker-compose.yml restart grafana

# Очистите кэш Grafana
docker volume rm labracodabrador_grafana-data
```

## 🤝 Полезные ссылки

- [Geth Documentation](https://geth.ethereum.org/docs)
- [Clique PoA Spec](https://eips.ethereum.org/EIPS/eip-225)
- [JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Prometheus](https://prometheus.io/docs/)
- [Grafana](https://grafana.com/docs/)

## 📄 Лицензия

MIT License

---