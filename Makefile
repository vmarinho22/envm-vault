.PHONY: help setup db-up db-down db-reset migrate dev test test-rtk rtk-trust

help:
	@echo "Available commands:"
	@echo "  make setup      - Copy .env.example to .env if missing"
	@echo "  make db-up      - Start PostgreSQL via docker compose"
	@echo "  make db-down    - Stop PostgreSQL"
	@echo "  make db-reset   - Reset database volume and run AutoMigrate"
	@echo "  make migrate    - Run database migrations (GORM AutoMigrate)"
	@echo "  make dev        - Run API server"
	@echo "  make test       - Run tests"
	@echo "  make test-rtk   - Run tests via RTK (token-optimized output)"
	@echo "  make rtk-trust  - Trust project RTK filters (.rtk/filters.toml)"

setup:
	@test -f .env || cp .env.example .env
	@echo "Setup complete. Edit .env if needed."

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-reset:
	docker compose down -v && docker compose up -d postgres
	@sleep 2
	@$(MAKE) migrate

migrate:
	@set -a && . ./.env 2>/dev/null; set +a && cd api && go run ./cmd/migrate

dev:
	@set -a && . ./.env 2>/dev/null; set +a && cd api && go run ./cmd/server

test:
	@set -a && . ./.env 2>/dev/null; set +a && cd api && go test ./...

test-rtk:
	cd api && rtk go test ./...

rtk-trust:
	rtk trust --yes
