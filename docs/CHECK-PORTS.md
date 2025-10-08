# Проверка портов на сервере

## Быстрая проверка

### Linux
```bash
# Проверить все наши порты
sudo netstat -tuln | grep -E ':(4000|4100|9080|9081|9545|9546|10090)'

# Проверить конкретный порт
sudo netstat -tuln | grep :9545
sudo lsof -i :9545
```

### macOS
```bash
# Проверить все наши порты
lsof -i :4000 -i :4100 -i :9080 -i :9081 -i :9545 -i :9546 -i :10090

# Проверить конкретный порт
lsof -i :9545
```

## Используемые порты

| Порт | Сервис | Описание |
|------|--------|----------|
| 9080 | Web Explorer | Веб-интерфейс для просмотра блокчейна |
| 9081 | REST API | REST API для доступа к блокчейну |
| 9545 | RPC HTTP | JSON-RPC через HTTP |
| 9546 | RPC WebSocket | JSON-RPC через WebSocket |
| 4000 | Grafana | Мониторинг дашборды |
| 10090 | Prometheus | Сбор метрик |
| 4100 | Loki | Сбор и хранение логов |

## Что делать если порт занят?

1. **Узнать что использует порт:**
   ```bash
   sudo lsof -i :9545
   ```

2. **Остановить процесс:**
   ```bash
   # Найти PID
   sudo lsof -i :9545 | grep LISTEN
   
   # Остановить процесс
   sudo kill <PID>
   ```

3. **Или изменить порт в docker-compose.yml:**
   ```yaml
   nginx:
     ports:
       - "НОВЫЙ_ПОРТ:8545"  # Измените на свободный порт
   ```

## Проверка после запуска

```bash
# Все контейнеры запущены?
docker ps | grep labracodabrador

# Проверка доступности сервисов
curl http://localhost:9545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

curl http://localhost:9081/api/v1/network
curl http://localhost:4000/api/health
curl http://localhost:10090/-/healthy
curl http://localhost:4100/ready
```
