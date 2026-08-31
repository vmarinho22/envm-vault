# Roadmap

## Phase 0 — Documentation & Agent Layer (current)

**Goal:** Establish the open-source foundation so humans and AI agents can contribute effectively.

Deliverables:
- [x] Project documentation (`docs/`)
- [x] Agent instructions (`AGENTS.md`, `.cursor/rules/`, `.cursor/skills/`)
- [x] OSS files (README, CONTRIBUTING, LICENSE, SECURITY)
- [x] GitHub templates (PR, issues)

No application code in this phase.

---

## Phase 1 — Scaffold & API Skeleton

**Goal:** Runnable monorepo with health check and basic project structure.

Deliverables:
- `backend/` — Go module, `cmd/server`, `internal/` layout
- `frontend/` — Vue 3 + Vite + TypeScript scaffold
- `docker-compose.yml` — PostgreSQL for local dev
- `Makefile` — common dev commands (`setup`, `db-up`, `dev`)
- Database migrations (initial schema)
- `GET /health` endpoint
- Basic CI (lint + test on PR)

---

## Phase 2 — Core API

**Goal:** Full CRUD and env var download via access token.

Deliverables:
- User registration and authentication
- Project, Environment, Secret CRUD
- Access token generation with scoping
- `GET /v1/projects/{project}/environments/{env}` — download env vars
- Encryption at rest for secret values
- API integration tests

---

## Phase 3 — Admin UI

**Goal:** Vue admin interface for managing vault contents.

Deliverables:
- Login / session management
- Project and environment management
- Secret editor (create, update, delete)
- Access token management UI
- Responsive layout

---

## Phase 4 — Production Readiness

**Goal:** Deployable, maintainable open-source release.

Deliverables:
- Docker image for API
- GitHub Actions release pipeline
- Helm chart or docker-compose for self-hosting
- Documentation site or expanded README
- Security audit checklist

---

## How to Contribute to a Phase

1. Check open issues labeled with the phase (e.g. `phase-1`).
2. Read [AGENTS.md](../AGENTS.md) and [architecture.md](architecture.md) before starting.
3. Open a draft PR early for feedback on approach.
4. Follow [CONTRIBUTING.md](../CONTRIBUTING.md) conventions.
