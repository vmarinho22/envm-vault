# AGENTS.md

@RTK.md

Instructions for AI coding agents working on **envm-vault**.

## Project Overview

envm-vault is an open-source environment variable vault. Users download environment variables for projects across customizable environments using a personal access token.

**Current status:** Phase 1 — API with auth + Vue login UI.

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
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── service/
│   │   ├── email/
│   │   ├── database/      ← GORM connection + AutoMigrate
│   │   ├── model/         ← GORM models
│   │   ├── repository/    ← data access layer
├── docs/
├── docker-compose.yml
├── Makefile
└── frontend/              ← Vue 3 admin UI (Vite, login OTP)
```

Always read `docs/architecture.md` before implementing domain logic.

## Security Rules (always apply)

1. **Never commit secrets** — no real tokens, passwords, or `.env` files with production values.
2. **Env split** — root `.env` is API + Docker only; `frontend/.env` is `VITE_*` only (public, baked into the browser bundle). Never put `JWT_SECRET` / `DATABASE_URL` in the frontend env.
3. **Never log secret values** — OTP codes are logged only in `APP_ENV=development`.
4. **OTP lockout** — max 5 failed verify attempts per code; HTTP rate limit via `RATE_LIMIT_MAX` / `RATE_LIMIT_WINDOW`.
5. **Token scoping** — access tokens must be limited to explicit project/environment pairs.
6. **Minimal diffs** — do not refactor unrelated code in the same change.

## Conventions

### General

- **Code language:** English (identifiers, comments, commit messages).
- **Technical docs:** English (`docs/`, `AGENTS.md`).
- **README:** Portuguese (human-facing).
- **Commits:** Conventional commits — `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`.
- **Scope:** Keep changes focused. One concern per PR.

### Go API (`api/`)

- **Go version:** 1.25+ (pinned in root `mise.toml`)
- **HTTP:** Fiber v2
- **Database:** PostgreSQL via GORM + repository layer
- **Migrations:** GORM AutoMigrate (schema in `internal/model/`)
- **Layout:**
  ```
  api/
  ├── cmd/server/main.go       # wiring only — calls router.Register
  ├── cmd/migrate/main.go      # standalone AutoMigrate
  ├── internal/
  │   ├── router/              # register routes by context (health, auth, user)
  │   ├── handler/           # subpackages: auth/, user/, health/, httperr/
  │   ├── service/           # subpackages: auth/, user/ + shared.go
  │   ├── middleware/        # JWT, auth guards
  │   ├── email/             # Email sender interface
  │   ├── database/          # Open, WithTx, Ping
  │   ├── model/             # GORM models (User, LoginCode)
  │   └── repository/        # data access (UserRepository, LoginCodeRepository)
  ```
- **Error handling:** Wrap errors with context; use sentinel errors in service layer
- **Tests:** Co-located `*_test.go`; integration tests use `//go:build integration` and require Postgres

### Commands

```bash
make setup      # copy root + frontend .env.example → .env if missing
make db-up      # start PostgreSQL
make migrate    # run GORM AutoMigrate
make db-reset   # reset database volume + AutoMigrate
make dev        # run API server
make frontend-dev  # Vue UI on :3100
make test       # unit + integration tests
make test-unit  # API + frontend unit tests (used by pre-commit hook)
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
- **Architecture:** module-oriented (`app/`, `shared/`, `modules/{domain}/`)
- **Services:** typed interfaces + instances — `authService`, `userService`, `healthService`, `httpClient`
- **State:** Pinia (`sessionStorage` for JWT)
- **UI:** Tailwind CSS v4 + daisyUI 5 + Heroicons
- **Dev server:** port **3100** (`strictPort`); proxies `/v1` and `/health` to the API on `:8080`
- **Env:** `frontend/.env` (`VITE_*` only); API secrets stay in root `.env`
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

- Run `make test` before finishing any task (unit + integration).
- Run `make test-unit` for fast feedback without Postgres (API + frontend).
- Integration tests use build tag `integration` and require `DATABASE_URL` + running PostgreSQL.
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
