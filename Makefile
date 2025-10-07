.PHONY: help build run stop logs clean backup restore health

# Help command
help:
	@echo "Labracodabrador - Ethereum PoA Blockchain"
	@echo ""
	@echo "Available commands:"
	@echo "  make build        - Build API server Docker image"
	@echo "  make run          - Start all services (blockchain + API + explorer)"
	@echo "  make stop         - Stop all services"
	@echo "  make restart      - Restart all services"
	@echo "  make logs         - Show logs from all services"
	@echo "  make logs-api     - Show API server logs"
	@echo "  make logs-signer  - Show signer nodes logs"
	@echo "  make clean        - Stop and remove all containers and volumes"
	@echo "  make backup       - Backup blockchain data"
	@echo "  make restore      - Restore blockchain from backup"
	@echo "  make health       - Check blockchain health"
	@echo "  make ps           - Show container status"

# Build API server
build:
	@echo "Building API server..."
	cd $(PWD) && docker-compose -f config/docker-compose.yml build api-server
	@echo "✅ Build complete"

# Start all services
run:
	@echo "Starting Labracodabrador blockchain..."
	cd $(PWD) && docker-compose -f config/docker-compose.yml up -d
	@echo "✅ Services started"
	@echo ""
	@echo "🌐 Web Explorer: http://localhost:8080"
	@echo "🔗 REST API:     http://localhost:8081/api/v1"
	@echo "⚡ RPC:          http://localhost:8545"
	@echo "📊 Grafana:      http://localhost:3000"
	@echo "📈 Prometheus:   http://localhost:9090"

# Stop all services
stop:
	@echo "Stopping services..."
	cd $(PWD) && docker-compose -f config/docker-compose.yml down
	@echo "✅ Services stopped"

# Restart all services
restart:
	@echo "Restarting services..."
	cd $(PWD) && docker-compose -f config/docker-compose.yml restart
	@echo "✅ Services restarted"

# Show logs
logs:
	cd $(PWD) && docker-compose -f config/docker-compose.yml logs -f

logs-api:
	cd $(PWD) && docker-compose -f config/docker-compose.yml logs -f api-server

logs-signer:
	cd $(PWD) && docker-compose -f config/docker-compose.yml logs -f signer1 signer2 signer3

# Container status
ps:
	cd $(PWD) && docker-compose -f config/docker-compose.yml ps

# Clean all data
clean:
	@echo "⚠️  WARNING: This will remove all blockchain data!"
	@read -p "Continue? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		cd $(PWD) && docker-compose -f config/docker-compose.yml down -v; \
		echo "✅ All data removed"; \
	else \
		echo "Aborted."; \
	fi

# Backup blockchain
backup:
	@echo "Running blockchain backup..."
	./scripts/backup-blockchain.sh
	@echo "✅ Backup complete"

# Restore blockchain
restore:
	@if [ -z "$(BACKUP)" ]; then \
		echo "Usage: make restore BACKUP=/path/to/backup.tar.gz"; \
		echo ""; \
		echo "Available backups:"; \
		ls -lh /backup/ethereum/blockchain-backup-*.tar.gz 2>/dev/null || echo "No backups found"; \
	else \
		./scripts/restore-blockchain.sh $(BACKUP); \
	fi

# Health check
health:
	@echo "Running health check..."
	./scripts/health-check.sh
