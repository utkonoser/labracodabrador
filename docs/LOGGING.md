# Логирование Labracodabrador Blockchain

## Обзор

Система логирования включает:
- **Loki** - централизованное хранение и индексация логов
- **Promtail** - сбор логов из контейнеров и файлов
- **Grafana** - визуализация и поиск логов

## Архитектура логирования

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Geth Nodes    │    │  Other Services │    │  Log Files      │
│   (Docker)      │    │  (Docker)       │    │  (File System)  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │       Promtail           │
                    │   (Log Collection)       │
                    └────────────┬─────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │        Loki             │
                    │   (Log Storage)         │
                    └────────────┬─────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │       Grafana           │
                    │   (Log Visualization)   │
                    └─────────────────────────┘
```

## Конфигурация

### Loki (`config/loki.yml`)
```yaml
auth_enabled: false
server:
  http_listen_port: 3100
  grpc_listen_port: 9095
common:
  path_prefix: /loki
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory
  storage:
    filesystem:
      directory: /loki/chunks
  wal:
    directory: /loki/wal
schema_config:
  configs:
    - from: 2020-10-27
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h
```

### Promtail (`config/promtail.yml`)
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://172.25.0.8:3100/loki/api/v1/push

scrape_configs:
  - job_name: containers
    static_configs:
      - targets:
          - localhost
        labels:
          job: containerlogs
          __path__: /var/log/containers/*.log

    pipeline_stages:
      - json:
          expressions:
            output: log
            stream: stream
            attrs:
      - json:
          expressions:
            tag:
          source: attrs
      - regex:
          expression: (?P<container_name>(?:[^|]*))\|
          source: tag
      - timestamp:
          format: RFC3339Nano
          source: time
      - labels:
          stream:
          container_name:
      - output:
          source: output

  - job_name: geth-logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: geth
          __path__: /var/log/geth/*.log

  - job_name: prometheus-logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: prometheus
          __path__: /var/log/prometheus/*.log

  - job_name: grafana-logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: grafana
          __path__: /var/log/grafana/*.log

  - job_name: nginx-logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: nginx
          __path__: /var/log/nginx/*.log
```

## Источники логов

### 1. Docker контейнеры
- **Источник**: Docker socket (`/var/run/docker.sock`)
- **Метки**: `container_name`, `stream` (stdout/stderr)
- **Формат**: JSON с полями `log`, `stream`, `time`, `attrs`

### 2. Файлы логов Geth
- **Путь**: `/var/log/geth/*.log`
- **Метки**: `job: geth`
- **Содержимое**: Логи Geth нод (signer1.log, signer2.log, signer3.log)

### 3. Файлы логов Prometheus
- **Путь**: `/var/log/prometheus/*.log`
- **Метки**: `job: prometheus`
- **Содержимое**: Логи Prometheus сервера

### 4. Файлы логов Grafana
- **Путь**: `/var/log/grafana/*.log`
- **Метки**: `job: grafana`
- **Содержимое**: Логи Grafana сервера

### 5. Файлы логов Nginx
- **Путь**: `/var/log/nginx/*.log`
- **Метки**: `job: nginx`
- **Содержимое**: Access и error логи Nginx

## Запросы логов (LogQL)

### Базовые запросы

```logql
# Все логи
{container_name=~".+"}

# Логи конкретного контейнера
{container_name="labracodabrador-signer1"}

# Логи по типу сервиса
{job="geth"}
{job="prometheus"}
{job="grafana"}
{job="nginx"}
```

### Фильтрация по содержимому

```logql
# Только ошибки
{container_name=~".+"} |~ "(?i)(error|fatal|panic)"

# Только предупреждения
{container_name=~".+"} |~ "(?i)(warn|warning)"

# Логи с определенным текстом
{container_name=~".+"} |= "block"
{container_name=~".+"} |= "transaction"
{container_name=~".+"} |= "peer"
```

### Временные фильтры

```logql
# Логи за последний час
{container_name=~".+"} [1h]

# Логи за последние 5 минут
{container_name=~".+"} [5m]

# Логи с определенного времени
{container_name=~".+"} |= "2024-01-01"
```

### Форматирование вывода

```logql
# Форматированный вывод с временными метками
{container_name=~".+"} | line_format "{{.timestamp}} [{{.level}}] {{.message}}"

# JSON формат
{container_name=~".+"} | json | line_format "{{.timestamp}} {{.level}} {{.message}}"
```

## Дашборды логов

### Logs Overview
Главный дашборд для просмотра логов:
- **Все логи контейнеров** - общий поток логов
- **Логи signer нод** - логи узлов, создающих блоки
- **Логи RPC нод** - логи узлов, обрабатывающих запросы
- **Ошибки** - фильтрация по ошибкам и критическим сообщениям
- **Предупреждения** - фильтрация по предупреждениям

### Настройка дашборда
1. Откройте Grafana: http://localhost:3000
2. Перейдите в **Dashboards** → **Logs Overview**
3. Используйте фильтры для поиска нужных логов
4. Настройте временной диапазон
5. Используйте LogQL для сложных запросов

## Мониторинг системы логирования

### Проверка статуса сервисов

```bash
# Статус Loki
curl http://localhost:4100/ready

# Статус Promtail
curl http://localhost:10080/ready

# Логи Promtail
docker logs labracodabrador-promtail

# Логи Loki
docker logs labracodabrador-loki
```

### Проверка сбора логов

```bash
# Проверка позиций Promtail
docker exec labracodabrador-promtail cat /tmp/positions.yaml

# Проверка конфигурации Promtail
docker exec labracodabrador-promtail cat /etc/promtail/config.yml

# Проверка доступности Loki API
curl "http://localhost:4100/loki/api/v1/query?query={container_name=~".+"}&limit=10"
```

### Проверка файлов логов

```bash
# Проверка логов Geth
ls -la logs/geth/
tail -f logs/geth/signer1.log

# Проверка логов Prometheus
ls -la logs/prometheus/
tail -f logs/prometheus/prometheus.log

# Проверка логов Grafana
ls -la logs/grafana/
tail -f logs/grafana/grafana.log
```

## Troubleshooting

### Логи не собираются

1. **Проверьте статус Promtail**:
   ```bash
   docker logs labracodabrador-promtail
   ```

2. **Проверьте доступность Loki**:
   ```bash
   curl http://localhost:4100/ready
   ```

3. **Проверьте права доступа к Docker socket**:
   ```bash
   ls -la /var/run/docker.sock
   ```

4. **Проверьте конфигурацию Promtail**:
   ```bash
   docker exec labracodabrador-promtail cat /etc/promtail/config.yml
   ```

### Логи не отображаются в Grafana

1. **Проверьте datasource Loki**:
   - Configuration → Data Sources
   - Убедитесь, что Loki доступен

2. **Проверьте запросы в дашборде**:
   - Используйте простые запросы: `{container_name=~".+"}`
   - Проверьте временной диапазон

3. **Проверьте доступность Loki из Grafana**:
   ```bash
   docker exec labracodabrador-grafana curl http://loki:3100/ready
   ```

### Медленная работа Loki

1. **Проверьте использование диска**:
   ```bash
   docker exec labracodabrador-loki df -h /loki
   ```

2. **Очистите старые данные**:
   ```bash
   # Удаление данных старше 7 дней
   docker exec labracodabrador-loki find /loki/chunks -type f -mtime +7 -delete
   ```

3. **Увеличьте retention в конфигурации**:
   ```yaml
   limits_config:
     retention_period: 168h  # 7 дней
   ```

## Настройка алертов на логи

### Создание алертов в Grafana

1. **Alert на ошибки**:
   - Query: `{container_name=~".+"} |~ "(?i)(error|fatal|panic)"`
   - Condition: `IS ABOVE 0`
   - For: `1m`

2. **Alert на отсутствие логов**:
   - Query: `{container_name="labracodabrador-signer1"}`
   - Condition: `IS BELOW 1`
   - For: `5m`

3. **Alert на высокую частоту ошибок**:
   - Query: `rate({container_name=~".+"} |~ "(?i)(error|fatal|panic)"[5m])`
   - Condition: `IS ABOVE 0.1`
   - For: `2m`

## Рекомендации

1. **Регулярно мониторьте дашборд Logs Overview** для выявления проблем
2. **Настройте алерты** на критические ошибки
3. **Используйте фильтры** для быстрого поиска нужной информации
4. **Настройте retention** в зависимости от потребностей
5. **Мониторьте производительность Loki** при большом объеме логов
6. **Регулярно очищайте старые логи** для экономии места
7. **Используйте LogQL** для сложных запросов и анализа
8. **Сохраняйте важные логи** при инцидентах

## Полезные ссылки

- [Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/)
- [Grafana Logs Panel](https://grafana.com/docs/grafana/latest/panels/visualizations/logs-panel/)
