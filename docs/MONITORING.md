# Мониторинг и логирование Labracodabrador Blockchain

## Обзор

Система мониторинга и логирования включает:
- **Prometheus** - сбор и хранение метрик
- **Grafana** - визуализация, дашборды и логи
- **Loki** - централизованное хранение логов
- **Promtail** - сбор логов из контейнеров и файлов
- **Alertmanager** (опционально) - отправка алертов

## Доступ к системам мониторинга

### Grafana
- **URL**: http://localhost:4000
- **Логин**: admin
- **Пароль**: admin

### Prometheus
- **URL**: http://localhost:10090
- **Targets**: http://localhost:10090/targets
- **Alerts**: http://localhost:10090/alerts
- **Rules**: http://localhost:10090/rules

### Loki
- **URL**: http://localhost:4100
- **API**: http://localhost:4100/loki/api/v1/query

## Доступные дашборды

### 1. Blockchain Overview
Главный дашборд для мониторинга всей сети:
- Высота блоков на всех узлах
- Количество подключенных пиров
- Скорость создания блоков (blocks/min)
- Пул транзакций (pending/queued)
- Использование газа
- Использование памяти и CPU
- Статистика по узлам (сколько узлов активны)
- Статистика RPC запросов
- Размер БД

### 2. Node Details
Детальный мониторинг отдельного узла:
- Статус узла (UP/DOWN)
- Текущий блок
- Количество пиров
- Pending транзакции
- Прогресс синхронизации
- P2P трафик (входящий/исходящий)
- Детальная статистика памяти
- Количество горутин
- Информация о подключенных пирах
- Активность пула транзакций
- Реорганизации цепи

### 3. Logs Overview
Централизованный просмотр логов всех сервисов:
- Все логи контейнеров
- Логи signer нод
- Логи RPC нод
- Фильтрация по ошибкам и предупреждениям
- Поиск по содержимому логов
- Временные диапазоны

## Система логирования

### Loki
Централизованное хранение логов:
- **Конфигурация**: `config/loki.yml`
- **Порт**: 3100
- **Хранение**: файловая система (`/loki/chunks`)
- **Retention**: 7 дней (настраивается)

### Promtail
Сбор логов из различных источников:
- **Конфигурация**: `config/promtail.yml`
- **Источники**:
  - Docker контейнеры (через Docker socket)
  - Файлы логов Geth (`/var/log/geth/*.log`)
  - Файлы логов Prometheus (`/var/log/prometheus/*.log`)
  - Файлы логов Grafana (`/var/log/grafana/*.log`)
  - Файлы логов Nginx (`/var/log/nginx/*.log`)

### Структура логов
Логи автоматически помечаются метками:
- `container_name` - имя контейнера
- `job` - тип сервиса (geth, prometheus, grafana, nginx)
- `level` - уровень лога (info, warn, error)
- `source` - источник лога

### Запросы логов в Grafana
```logql
# Все логи
{container_name=~".+"}

# Логи signer нод
{container_name=~".*signer.*"}

# Только ошибки
{container_name=~".+"} |~ "(?i)(error|fatal|panic)"

# Логи за последний час
{container_name=~".+"} |= "block" | line_format "{{.timestamp}} {{.message}}"
```

## Метрики Geth

Geth экспортирует метрики на эндпоинт: `http://<node-ip>:6060/debug/metrics/prometheus`

### Основные метрики

#### Блокчейн
- `chain_head_block` - текущая высота блока
- `chain_head_receipt_gasused` - использованный газ в последнем блоке
- `chain_reorg_add` - количество добавленных блоков при реорганизации
- `chain_reorg_drop` - количество удаленных блоков при реорганизации

#### P2P сеть
- `p2p_peers` - количество подключенных пиров
- `p2p_ingress` - входящий трафик (байты)
- `p2p_egress` - исходящий трафик (байты)

#### Пул транзакций
- `txpool_pending` - ожидающие транзакции
- `txpool_queued` - транзакции в очереди
- `txpool_valid` - валидные транзакции
- `txpool_invalid` - невалидные транзакции

#### Система
- `system_memory_allocs` - выделенная память
- `system_memory_used` - используемая память
- `system_memory_held` - удерживаемая память
- `system_goroutines` - количество горутин
- `system_cpu_sysload` - загрузка CPU

#### RPC
- `rpc_requests` - количество RPC запросов
- `rpc_duration_all` - время выполнения RPC запросов

#### База данных
- `eth_db_chaindata_disk_size` - размер БД на диске

## Настроенные алерты

### Критические (severity: critical)

#### NodeDown
- **Условие**: узел недоступен более 2 минут
- **Действие**: немедленно проверить статус узла

#### NoNewBlocks
- **Условие**: нет новых блоков более 3 минут
- **Действие**: проверить консенсус и подключение узлов

### Предупреждения (severity: warning)

#### LowPeerCount
- **Условие**: менее 2 подключенных пиров более 5 минут
- **Действие**: проверить сетевое подключение

#### HighCPUUsage
- **Условие**: использование CPU > 80% более 10 минут
- **Действие**: проверить нагрузку на узел

#### HighMemoryUsage
- **Условие**: использование памяти > 1GB более 10 минут
- **Действие**: мониторить утечки памяти

#### LargeTxPoolPending
- **Условие**: более 1000 ожидающих транзакций более 5 минут
- **Действие**: проверить скорость обработки транзакций

#### LowTxProcessingRate
- **Условие**: низкая скорость обработки при наличии ожидающих транзакций
- **Действие**: проверить загрузку узла

#### HighRPCLatency
- **Условие**: среднее время ответа RPC > 1 секунды
- **Действие**: оптимизировать нагрузку на RPC узлы

## Запросы для анализа

### Prometheus Query Examples

#### Текущая высота блока
```promql
chain_head_block
```

#### Скорость создания блоков (blocks/min)
```promql
rate(chain_head_block[5m]) * 60
```

#### Общее количество пиров в сети
```promql
sum(p2p_peers)
```

#### Среднее количество пиров на узел
```promql
avg(p2p_peers)
```

#### Узлы с низким количеством пиров
```promql
p2p_peers < 2
```

#### Общий размер txpool в сети
```promql
sum(txpool_pending)
```

#### Использование памяти по узлам
```promql
system_memory_used{node_type="signer"}
```

#### RPC запросы в секунду
```promql
rate(rpc_requests[5m])
```

#### Топ-5 самых частых RPC методов
```promql
topk(5, rate(rpc_requests[5m]))
```

## Настройка алертов через Alertmanager

Для отправки алертов на почту, Slack, Telegram и т.д., настройте Alertmanager:

1. Создайте `alertmanager.yml`:
```yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'cluster']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'email'

receivers:
  - name: 'email'
    email_configs:
      - to: 'your-email@example.com'
        from: 'alertmanager@labracodabrador.local'
        smarthost: 'smtp.gmail.com:587'
        auth_username: 'your-email@gmail.com'
        auth_password: 'your-app-password'
```

2. Добавьте Alertmanager в docker-compose.yml
3. Обновите Prometheus для отправки алертов в Alertmanager

## Полезные команды

### Проверка метрик напрямую
```bash
# Проверка метрик с узла
curl http://172.25.0.11:6060/debug/metrics/prometheus | grep chain_head_block

# Проверка через Prometheus API
curl 'http://localhost:10090/api/v1/query?query=chain_head_block'
```

### Проверка статуса targets
```bash
curl -s http://localhost:10090/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health, error: .lastError}'
```

### Проверка активных алертов
```bash
curl -s http://localhost:10090/api/v1/alerts | jq '.data.alerts[] | {alert: .labels.alertname, state: .state}'
```

## Расширенная настройка

### Увеличение retention Prometheus
В docker-compose.yml измените параметр:
```yaml
--storage.tsdb.retention.time=90d  # вместо 30d
```

### Добавление новых дашбордов
1. Создайте дашборд в Grafana UI
2. Экспортируйте JSON через Settings -> JSON Model
3. Сохраните в `grafana/dashboards/`
4. Перезапустите Grafana

### Добавление новых алертов
1. Отредактируйте `config/prometheus-alerts.yml`
2. Добавьте новое правило в группу
3. Перезапустите Prometheus: `make restart`

## Troubleshooting

### Метрики не собираются
1. Проверьте targets: http://localhost:10090/targets
2. Убедитесь, что узлы доступны: `curl http://172.25.0.11:6060/debug/metrics/prometheus`
3. Проверьте логи Prometheus: `docker logs labracodabrador-prometheus`

### Дашборды не загружаются
1. Проверьте datasource в Grafana: Configuration -> Data Sources
2. Убедитесь, что Prometheus доступен из Grafana
3. Проверьте логи Grafana: `docker logs labracodabrador-grafana`

### Алерты не работают
1. Проверьте правила: http://localhost:10090/rules
2. Убедитесь, что файл алертов смонтирован: `docker exec labracodabrador-prometheus cat /etc/prometheus/alerts/prometheus-alerts.yml`
3. Проверьте синтаксис правил в Prometheus UI

### Логи не собираются
1. Проверьте статус Promtail: `docker logs labracodabrador-promtail`
2. Убедитесь, что Loki доступен: `curl http://localhost:4100/ready`
3. Проверьте конфигурацию Promtail: `docker exec labracodabrador-promtail cat /etc/promtail/config.yml`
4. Проверьте права доступа к Docker socket: `ls -la /var/run/docker.sock`

### Логи не отображаются в Grafana
1. Проверьте datasource Loki в Grafana: Configuration -> Data Sources
2. Убедитесь, что Loki доступен из Grafana
3. Проверьте запросы в дашборде Logs Overview
4. Проверьте временной диапазон в дашборде

## Рекомендации

1. **Регулярно мониторьте дашборд Blockchain Overview** для общего состояния сети
2. **Настройте алерты** для критических метрик
3. **Используйте дашборд Logs Overview** для анализа проблем
4. **Сохраняйте снапшоты** важных графиков при инцидентах
5. **Регулярно проверяйте размер БД** и планируйте дисковое пространство
6. **Мониторьте количество пиров** - если их мало, сеть может быть нестабильной
7. **Следите за txpool** - большое количество pending транзакций может указывать на проблемы
8. **Настройте retention для логов** - по умолчанию 7 дней, увеличьте при необходимости
9. **Используйте фильтры логов** для быстрого поиска ошибок и предупреждений
10. **Мониторьте производительность Loki** - при большом объеме логов может потребоваться оптимизация

## Полезные ссылки

- [Geth Metrics Documentation](https://geth.ethereum.org/docs/monitoring/metrics)
- [Prometheus Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/)
- [LOGGING.md](LOGGING.md) - Подробное руководство по системе логирования

