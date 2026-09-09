# AGENTS.md

@RTK.md

Instructions for AI coding agents working on **envm-vault**.

## Project Overview

envm-vault is an open-source environment variable vault. Users download environment variables for projects across customizable environments using a personal access token.

**Current status:** Phase 1 — API with auth, unique API key per user, and Vue login/onboarding UI.

**Stack:**
- API: Go + Fiber (`api/`)
- Frontend: Vue 3 + Vite + TypeScript (`frontend/`)
- Database: PostgreSQL via GORM + AutoMigrate

## Repository Map

```
envm-vault/
├── AGENTS.md              ← you are here
├── RTK.md                 ← token-optimized CLI usage for agents
├── .rtk/filters.toml      ← project-local RTK filters
├── mise.toml              ← Node + Go versions (mise)
├── api/                   ← Go API (Fiber)
│   ├── cmd/server/
│   ├── internal/
│   │   ├── env/           ← Envs singleton + Load()
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── service/
│   │   ├── email/
│   │   ├── database/      ← GORM connection + AutoMigrate
│   │   ├── model/         ← GORM models
│   │   ├── repository/    ← data access layer
│   └── tests/unit/        ← unit tests (mirrors internal/)
├── docs/
├── compose.yaml
├── Dockerfile
├── Makefile
└── frontend/              ← Vue 3 admin UI (Vite, login + API key)
```

Always read `docs/architecture.md` before implementing domain logic.

## Security Rules (always apply)

1. **Never commit secrets** — no real tokens, passwords, or `.env` files with production values.
2. **Env split** — Local: `api/.env`, `frontend/.env`. Containers: `.deploy/db/.env`, `.deploy/api/.env` (required at runtime); `.deploy/frontend/.env` only for image **build**. Image-only users create envs via `cp`/`curl` from `.env.example` (see README) — do not assume `make setup`. Never put `JWT_SECRET` / `DATABASE_URL` / `ENCRYPTION_KEY` in frontend env.
3. **Never log secret values** — OTP codes are logged only in `APP_ENV=development`.
4. **OTP lockout** — max 5 failed verify attempts per code; HTTP rate limit via `RATE_LIMIT_MAX` / `RATE_LIMIT_WINDOW` (API `/v1` only when SPA is served).
5. **API keys** — one active key per user, hashed at rest; plaintext returned only from `POST /v1/tokens`.
6. **Secret encryption** — `ENCRYPTION_KEY` (32-byte hex) encrypts secret values at rest (AES-256-GCM); separate from `JWT_SECRET`.
7. **Audit** — `audit_logs` via `service/audit`; never include OTP, JWT, or full personal token plaintext.
8. **Minimal diffs** — do not refactor unrelated code in the same change.
9. **Container** — Default Compose image is GHCR (`ghcr.io/envm-vault/envm-vault`). No-clone users: curl `compose.yaml` + `.deploy/{db,api}/.env`, then `docker compose up -d`. Local build: `make up` (sets `APP_IMAGE=envm-vault:local`).

## Conventions

### General

- **Code language:** English (identifiers, comments, commit messages).
- **Technical docs:** English (`docs/`, `AGENTS.md`).
- **README:** Portuguese (human-facing).
- **Commits:** Conventional commits — `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`.
- **Scope:** Keep changes focused. One concern per PR.

### Go API (`api/`)

- **Go version:** 1.25+ (pinned in root `mise.toml`)
- **Config / env:** `internal/env` — call `env.Load()` then read nested groups (`env.Envs.Database.URL`, `env.Envs.JWT.Secret`, …). Required: `DATABASE_URL`, `JWT_SECRET`. Do not read `os.Getenv` for app config outside this package.
- **HTTP:** Fiber v2
- **Database:** PostgreSQL via GORM + repository layer
- **Migrations:** GORM AutoMigrate (schema in `internal/model/`)
- **Layout:**
  ```
  api/
  ├── cmd/server/main.go       # wiring only — calls router.Register
  ├── cmd/migrate/main.go      # standalone AutoMigrate
  ├── internal/
  │   ├── env/               # Envs singleton + Load() validation
  │   ├── router/            # register routes by context (health, auth, user, token)
  │   ├── handler/           # subpackages: auth/, user/, token/, health/, httperr/
  │   ├── service/           # subpackages: auth/, user/, token/, audit/ + shared.go
  │   ├── middleware/        # JWT, auth guards, rate limit, audit origin
  │   ├── email/             # Email sender interface
  │   ├── database/          # Open, WithTx, Ping
  │   ├── model/             # GORM models (User, LoginCode, AccessToken, AuditLog)
  │   └── repository/        # data access (User, LoginCode, AccessToken, AuditLog)
  └── tests/unit/            # unit tests; mirrors internal/ (package xxx_test)
  ```
- **Error handling:** Wrap errors with context; use sentinel errors in service layer
- **Audit:** Security-relevant mutations call `audit.Logger.Record` after success. Middleware attaches IP/UA/Origin/channel; JWT Protect sets `logged_in_user`. Persist failures must not fail the HTTP response. Never put OTP, JWT, or full personal token plaintext in audit rows.
- **Tests:** `api/tests/unit/` mirroring `internal/` (`package xxx_test`); no co-located `*_test.go`

### Commands

```bash
make setup      # copy api/, frontend/, and .deploy/db .env.example → .env if missing
make db-up      # start PostgreSQL
make migrate    # run GORM AutoMigrate
make db-reset   # reset database volume + AutoMigrate
make dev        # run API server
make frontend-dev  # Vue UI on :3100
make test       # API + frontend unit tests
make test-unit  # same as make test (used by pre-commit hook)
make hooks-install  # install lefthook git hooks
```

**Use RTK** to reduce token usage when running shell commands. See [RTK.md](RTK.md).

```bash
cd api && rtk go test ./...
rtk docker compose ps
rtk gain        # check token savings
```

### Vue Frontend (`frontend/`)

- **Node:** version pinned in root `mise.toml` (`mise install`); do not use `.nvmrc`
- **Vue 3** with Composition API and `<script setup lang="ts">`
- **Build:** Vite + TypeScript (strict, `noUncheckedIndexedAccess`)
- **Architecture:** module-oriented (`app/` including `translations/`, `shared/`, `modules/{domain}/`)
- **Services:** typed interfaces + instances — `authService`, `userService`, `healthService`, `tokenService`, `httpClient`
- **State:** Pinia (`sessionStorage` for JWT)
- **i18n:** vue-i18n Composition API; locales `pt-BR` / `en-US` under `app/translations/{locale}/*.json` (one file per context); browser locale detection at boot
- **UI:** Tailwind CSS v4 + daisyUI 5 + Heroicons
- **Dev server:** port **3100** (`strictPort`); proxies `/v1` and `/health` to the API on `:8080`
- **Env:** `shared/env/envs.ts` — `import { envs } from '@/shared/env/envs'`; never read `import.meta.env` outside it
- **Tests:** Vitest under `frontend/tests/{unit,component}/` (`npm test` / `make test-frontend`)
- **Component snapshots:** every component test must call `expectComponentSnapshot(wrapper)` from `tests/component/helpers/mount` and commit the generated `__snapshots__/` files

## Skills

Project skills live in `.cursor/skills/`:

| Skill | Use when |
|-------|----------|
| `env-vault-domain` | Working with Project, Environment, Token, Secret entities |
| `add-api-endpoint` | Adding a new Go API endpoint |
| `add-ui-screen` | Adding a new Vue screen (Phase 3+) |
| `security-review` | Reviewing auth, token, or secret-related changes |

## Testing

- Run `make test` before finishing any task (API + frontend unit tests).
- API tests live in `api/tests/unit/` and mirror `internal/` (e.g. `internal/service/token/` → `tests/unit/service/token/`).
- Pre-commit hook runs API unit tests when Go files are staged, and frontend Vitest + `npm run type-check` when Vue/TS files are staged.
- Frontend component tests always include a HTML snapshot via `expectComponentSnapshot(wrapper)`.

## Pull Request Guidelines

- Title: conventional commit format (e.g. `feat: add project CRUD endpoint`)
- Include a test plan describing what was verified.
- For security-related changes, confirm no secrets are included in the diff.

## Key References

- [Architecture](docs/architecture.md) — domain model and security principles
- [Roadmap](docs/roadmap.md) — project phases
- [Agent Guide](docs/agent-guide.md) — working with AI agents in this repo
- [Contributing](CONTRIBUTING.md) — contribution workflow
