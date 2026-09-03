# AGENTS.md

@RTK.md

Instructions for AI coding agents working on **envm-vault**.

## Project Overview

envm-vault is an open-source environment variable vault. Users download environment variables for projects across customizable environments using a personal access token.

**Current status:** Phase 1 — API with auth (health, login, verify, register).

**Stack:**
- API: Go + Fiber (`api/`)
- Frontend: Vue 3 admin UI (Phase 3)
- Database: PostgreSQL via GORM + AutoMigrate

## Repository Map

```
envm-vault/
├── AGENTS.md              ← you are here
├── RTK.md                 ← token-optimized CLI usage for agents
├── .rtk/filters.toml      ← project-local RTK filters
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
└── frontend/              ← Vue UI (Phase 3 — not created yet)
```

Always read `docs/architecture.md` before implementing domain logic.

## Security Rules (always apply)

1. **Never commit secrets** — no real tokens, passwords, or `.env` files with production values.
2. **Never log secret values** — OTP codes are logged only in `APP_ENV=development`.
3. **Token scoping** — access tokens must be limited to explicit project/environment pairs.
4. **Minimal diffs** — do not refactor unrelated code in the same change.

## Conventions

### General

- **Code language:** English (identifiers, comments, commit messages).
- **Technical docs:** English (`docs/`, `AGENTS.md`).
- **README:** Portuguese (human-facing).
- **Commits:** Conventional commits — `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`.
- **Scope:** Keep changes focused. One concern per PR.

### Go API (`api/`)

- **Go version:** 1.25+
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
- **Tests:** Co-located `*_test.go`; integration tests skip without `DATABASE_URL`

### Commands

```bash
make setup      # copy .env.example → .env
make db-up      # start PostgreSQL
make migrate    # run GORM AutoMigrate
make db-reset   # reset database volume + AutoMigrate
make dev        # run API server
make test       # run tests
```

**Use RTK** to reduce token usage when running shell commands. See [RTK.md](RTK.md).

```bash
cd api && rtk go test ./...
rtk docker compose ps
rtk gain        # check token savings
```

### Vue Frontend (planned — Phase 3+)

- **Vue 3** with Composition API and `<script setup lang="ts">`
- **Build:** Vite + TypeScript (strict mode)
- **State:** Pinia stores

## Skills

Project skills live in `.cursor/skills/`:

| Skill | Use when |
|-------|----------|
| `env-vault-domain` | Working with Project, Environment, Token, Secret entities |
| `add-api-endpoint` | Adding a new Go API endpoint |
| `add-ui-screen` | Adding a new Vue screen (Phase 3+) |
| `security-review` | Reviewing auth, token, or secret-related changes |

## Testing

- Run `make test` before finishing any task.
- Add tests for new behavior; skip trivial assertions.
- Integration tests require `DATABASE_URL` and running PostgreSQL.

## Pull Request Guidelines

- Title: conventional commit format (e.g. `feat: add project CRUD endpoint`)
- Include a test plan describing what was verified.
- For security-related changes, confirm no secrets are included in the diff.

## Key References

- [Architecture](docs/architecture.md) — domain model and security principles
- [Roadmap](docs/roadmap.md) — project phases
- [Agent Guide](docs/agent-guide.md) — working with AI agents in this repo
- [Contributing](CONTRIBUTING.md) — contribution workflow
