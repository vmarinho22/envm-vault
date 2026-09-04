.PHONY: help setup db-up db-down db-reset migrate dev frontend-dev test test-unit test-integration test-frontend test-rtk rtk-trust hooks-install hooks-uninstall

help:
	@echo "Available commands:"
	@echo "  make setup             - Copy root and frontend .env.example files if missing"
	@echo "  make db-up             - Start PostgreSQL via docker compose"
	@echo "  make db-down           - Stop PostgreSQL"
	@echo "  make db-reset          - Reset database volume and run AutoMigrate"
	@echo "  make migrate           - Run database migrations (GORM AutoMigrate)"
	@echo "  make dev               - Run API server"
	@echo "  make frontend-dev      - Run Vue admin UI (port 3100)"
	@echo "  make test              - Run unit + integration tests"
	@echo "  make test-unit         - Run API + frontend unit tests (no Postgres)"
	@echo "  make test-frontend     - Run frontend unit tests (Vitest)"
	@echo "  make test-integration  - Run API integration tests (requires Postgres)"
	@echo "  make test-rtk          - Run tests via RTK (token-optimized output)"
	@echo "  make hooks-install     - Install git hooks via lefthook"
	@echo "  make hooks-uninstall   - Remove git hooks installed by lefthook"
	@echo "  make rtk-trust         - Trust project RTK filters (.rtk/filters.toml)"

setup:
	@test -f .env || cp .env.example .env
	@test -f frontend/.env || cp frontend/.env.example frontend/.env
	@echo "Setup complete. Edit .env (API) and frontend/.env (UI) if needed."

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

frontend-dev:
	cd frontend && npm run dev

test: test-unit test-integration

test-unit:
	cd api && go test ./...
	cd frontend && npm test

test-frontend:
	cd frontend && npm test

test-integration:
	@set -a && . ./.env 2>/dev/null; set +a && cd api && go test -tags=integration ./...

test-rtk:
	cd api && rtk go test ./...

hooks-install:
	lefthook install

hooks-uninstall:
	lefthook uninstall

rtk-trust:
	rtk trust --yes
