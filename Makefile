.PHONY: help setup db-up db-down db-reset migrate dev frontend-dev test test-unit test-frontend test-rtk rtk-trust hooks-install hooks-uninstall

help:
	@echo "Available commands:"
	@echo "  make setup             - Copy api/, frontend/, and .docker/db .env.example files if missing"
	@echo "  make db-up             - Start PostgreSQL via docker compose"
	@echo "  make db-down           - Stop PostgreSQL"
	@echo "  make db-reset          - Reset database volume and run AutoMigrate"
	@echo "  make migrate           - Run database migrations (GORM AutoMigrate)"
	@echo "  make dev               - Run API server"
	@echo "  make frontend-dev      - Run Vue admin UI (port 3100)"
	@echo "  make test              - Run API + frontend unit tests"
	@echo "  make test-unit         - Same as make test"
	@echo "  make test-frontend     - Run frontend unit tests (Vitest)"
	@echo "  make test-rtk          - Run tests via RTK (token-optimized output)"
	@echo "  make hooks-install     - Install git hooks via lefthook"
	@echo "  make hooks-uninstall   - Remove git hooks installed by lefthook"
	@echo "  make rtk-trust         - Trust project RTK filters (.rtk/filters.toml)"

COMPOSE_ENV := --env-file .docker/db/.env

setup:
	@test -f api/.env || cp api/.env.example api/.env
	@test -f frontend/.env || cp frontend/.env.example frontend/.env
	@test -f .docker/db/.env || cp .docker/db/.env.example .docker/db/.env
	@echo "Setup complete. Edit api/.env (API), frontend/.env (UI), and .docker/db/.env (Postgres) if needed."

db-up:
	docker compose $(COMPOSE_ENV) up -d postgres

db-down:
	docker compose $(COMPOSE_ENV) down

db-reset:
	docker compose $(COMPOSE_ENV) down -v && docker compose $(COMPOSE_ENV) up -d postgres
	@sleep 2
	@$(MAKE) migrate

migrate:
	cd api && go run ./cmd/migrate

dev:
	cd api && go run ./cmd/server

frontend-dev:
	cd frontend && npm run dev

test: test-unit

test-unit:
	cd api && go test ./...
	cd frontend && npm test

test-frontend:
	cd frontend && npm test

test-rtk:
	cd api && rtk go test ./...

hooks-install:
	lefthook install

hooks-uninstall:
	lefthook uninstall

rtk-trust:
	rtk trust --yes
