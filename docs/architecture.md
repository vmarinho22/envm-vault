# Architecture

This document describes the architecture of **envm-vault**.

## Overview

envm-vault is an open-source environment variable vault. Users download environment variables for projects across customizable environments using a personal access token.

```mermaid
flowchart LR
  Client[CI_or_Developer] -->|"Bearer user_token"| API[Go_API_Fiber]
  Admin[Admin_UI_Vue] -->|JWT| API
  API --> DB[(PostgreSQL)]
  API -->|"GET env vars"| Client
```

## Current Implementation (Phase 1)

The API lives in `api/` and provides authentication via email OTP plus a unique API key per user.

### Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.25+ |
| HTTP | Fiber v2 |
| Database | PostgreSQL 16 |
| DB driver | GORM (postgres driver) |
| Schema | GORM AutoMigrate |
| Data access | repository layer |
| Auth | Email OTP (8 digits, 10 min, max 5 verify attempts) + JWT |
| Admin UI | Vue 3 + Vite + TypeScript |
| UI styling | Tailwind CSS v4 + daisyUI 5 |

### API Endpoints

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| `GET` | `/v1/health` | — | Health check (API + DB) |
| `GET` | `/health` | — | Health check alias (infra probes) |
| `POST` | `/v1/auth/login` | — | Request OTP via email |
| `POST` | `/v1/auth/verify` | — | Verify OTP, returns JWT |
| `POST` | `/v1/auth/register` | JWT (admin+) | Register user (email + role) |
| `GET` | `/v1/tokens/me` | JWT | Current user's API key status (no secret) |
| `POST` | `/v1/tokens` | JWT | Generate or rotate the user's API key (plaintext once) |
| `DELETE` | `/v1/tokens` | JWT | Revoke the current user's API key |
| `GET` | `/v1/users` | JWT (admin+) | List users with API key status |
| `DELETE` | `/v1/users/:id/token` | JWT (admin+) | Revoke another user's API key |

### Auth Flow

1. **Bootstrap:** When no users exist, the first successful OTP verification creates a `super_admin`.
2. **Login:** Registered users request OTP via `/v1/auth/login`, verify via `/v1/auth/verify`.
3. **Register:** `super_admin` or `admin` creates users via `/v1/auth/register` with role `admin` or `user`.
4. **API key:** After login, the UI requires generating a unique API key (`POST /v1/tokens`) before other screens. The plaintext is shown once. `POST` again rotates (invalidates the previous hash). The user can revoke their key; `admin` can revoke `user`/`admin` keys; only `super_admin` can revoke a `super_admin` key.

In development (`APP_ENV=development`), OTP codes are logged to stdout instead of sent via email.

The admin UI (`frontend/`) implements the same OTP flow on `/login`, then `/setup-token` when no active key exists. Dev server: `http://localhost:3100` (proxies `/v1` to the API). CORS origins are configured via `CORS_ORIGINS` (JSON array).

JWT is the admin UI session. The API key (`Authorization: Bearer ev_…`) is the programmatic credential for future env downloads. Which projects/environments a key can access will be decided by user membership, not by scopes on the token.

## Core Concepts (planned)

### Project

A logical grouping of environment variables (e.g. `my-app`, `payment-service`).

### Environment

A named deployment target within a project. Environments are fully customizable per project.

### Secret / EnvVar

A key-value pair scoped to a **Project + Environment** combination.

### User

An account with role `super_admin`, `admin`, or `user`. Authenticates via email OTP.

### AccessToken

A unique programmatic credential per user for downloading env vars. One active key at a time. Stored as a SHA-256 hash; only a display prefix is returned after creation. Sent as `Authorization: Bearer <token>` on env download (endpoint not implemented yet).

## Repository Structure

```
envm-vault/
├── api/                  # Go API (Fiber)
│   ├── cmd/server/
│   ├── tests/unit/       # unit tests (mirrors internal/)
│   ├── internal/
│   │   ├── env/          # Envs singleton + Load()
│   │   ├── router/       # Route registration by context (health, auth, user, token)
│   │   ├── handler/      # Handlers per domain (auth/, user/, token/, health/)
│   │   ├── service/      # Business logic per domain (auth/, user/, token/)
│   │   ├── middleware/
│   │   ├── database/     # GORM connection + AutoMigrate
│   │   ├── model/        # GORM models
│   │   ├── repository/   # data access layer
├── frontend/             # Vue 3 admin UI (Vite :3100)
├── docs/
├── docker-compose.yml
└── Makefile
```

## Data Model

```mermaid
erDiagram
  User ||--o{ LoginCode : requests
  User ||--o| AccessToken : has
  User ||--o{ Project : owns
  Project ||--o{ Environment : contains
  Environment ||--o{ Secret : stores

  User {
    uuid id PK
    string email UK
    string role
    timestamp created_at
    timestamp updated_at
  }

  LoginCode {
    uuid id PK
    string email
    string code_hash
    timestamp expires_at
    timestamp used_at
  }

  AccessToken {
    uuid id PK
    uuid user_id UK
    string token_hash UK
    string prefix
    timestamp revoked_at
    uuid revoked_by
    timestamp created_at
  }
```

### User Roles

| Role | Description |
|------|-------------|
| `super_admin` | First user (bootstrap). Full access. Can revoke any API key. |
| `admin` | Can register users and revoke `user`/`admin` API keys. |
| `user` | Standard user. Manages their own API key. |

## Security Principles

1. **Never store plaintext OTP codes** — only SHA-256 hashes in `login_codes`.
2. **Never store plaintext API keys** — only SHA-256 hashes in `access_tokens`. The plaintext is returned only from `POST /v1/tokens`.
3. **Never log secret values** — OTP codes are logged only in development. API keys are never logged.
4. **OTP verify lockout** — max 5 failed attempts per code (`attempt_count`); then the code is invalidated and the user must request a new OTP.
5. **Global HTTP rate limit** — per-IP limiter (`RATE_LIMIT_MAX` / `RATE_LIMIT_WINDOW`, defaults 100/1m). `/health` and `/v1/health` are exempt. In-memory store (single instance); shared store needed for multi-replica later.
6. **Env authorization** — future env download identifies the user via the API key; access to a project/environment is based on membership, not token scopes.
7. **No secrets in the repository** — never commit real tokens or `.env` files.
8. **JWT expiry** — configurable via `JWT_EXPIRY` (default 24h).

## Open Decisions

| Decision | Status |
|----------|--------|
| Response format (`.env` vs JSON) | TBD (Phase 2) |
| SMTP provider for production email | TBD |
| Encryption at rest for secrets | TBD (Phase 2) |
