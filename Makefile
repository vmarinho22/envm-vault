.PHONY: help setup db-up db-down db-reset migrate dev frontend-dev up down logs test test-unit test-frontend test-rtk rtk-trust hooks-install hooks-uninstall

help:
	@echo "Development (API + UI separately):"
	@echo "  make setup          - Copy .env examples if missing (optional; see README)"
	@echo "  make db-up          - Start Postgres only"
	@echo "  make migrate        - Run AutoMigrate"
	@echo "  make dev            - API on :8080"
	@echo "  make frontend-dev   - UI on :3100"
	@echo ""
	@echo "Container (API + UI in one image + Postgres):"
	@echo "  make up             - Build image, start stack → http://localhost:8080"
	@echo "  make logs           - Follow app logs (OTP when APP_ENV=development)"
	@echo "  make down           - Stop stack"
	@echo "  (published image: create .deploy/{db,api}/.env then APP_IMAGE=… docker compose up -d)"
	@echo ""
	@echo "Other: make test | hooks-install | rtk-trust"

COMPOSE := docker compose
COMPOSE_DB := $(COMPOSE) --env-file .deploy/db/.env

setup:
	@test -f api/.env || cp api/.env.example api/.env
	@test -f frontend/.env || cp frontend/.env.example frontend/.env
	@test -f .deploy/db/.env || cp .deploy/db/.env.example .deploy/db/.env
	@test -f .deploy/api/.env || cp .deploy/api/.env.example .deploy/api/.env
	@test -f .deploy/frontend/.env || cp .deploy/frontend/.env.example .deploy/frontend/.env
	@echo "Setup complete."
	@echo "  Dev:       api/.env + frontend/.env → make db-up && make migrate && make dev"
	@echo "  Container: .deploy/{db,api,frontend}/.env → make up && make logs"

db-up:
	@test -f .deploy/db/.env || cp .deploy/db/.env.example .deploy/db/.env
	$(COMPOSE_DB) up -d postgres

db-down:
	$(COMPOSE_DB) down

db-reset:
	$(COMPOSE_DB) down -v && $(COMPOSE_DB) up -d postgres
	@sleep 2
	@$(MAKE) migrate

migrate:
	cd api && go run ./cmd/migrate

dev:
	cd api && go run ./cmd/server

frontend-dev:
	cd frontend && npm run dev

up:
	@test -f .deploy/db/.env || cp .deploy/db/.env.example .deploy/db/.env
	@test -f .deploy/api/.env || cp .deploy/api/.env.example .deploy/api/.env
	@test -f .deploy/frontend/.env || cp .deploy/frontend/.env.example .deploy/frontend/.env
	@echo "Building and starting → http://localhost:8080"
	APP_IMAGE=envm-vault:local $(COMPOSE) up -d --build
	@echo "Ready. OTP: make logs | Stop: make down"

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f app

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
