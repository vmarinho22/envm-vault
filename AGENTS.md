# AGENTS.md

Instructions for AI coding agents working on **envm-vault**.

## Project Overview

envm-vault is an open-source environment variable vault. Users download environment variables for projects across customizable environments using a personal access token.

**Current status:** Phase 0 — documentation and agent layer only. No application code exists yet.

**Planned stack:**
- Backend: Go monolith API
- Frontend: Vue 3 admin UI
- Database: PostgreSQL

## Repository Map

```
envm-vault/
├── AGENTS.md              ← you are here
├── docs/
│   ├── architecture.md    ← domain model, data model, security principles
│   ├── roadmap.md         ← project phases
│   └── agent-guide.md     ← how to use agents in this repo
├── .cursor/
│   ├── rules/             ← Cursor-specific scoped rules
│   └── skills/            ← reusable workflow skills
├── backend/               ← Go API (Phase 1 — not created yet)
└── frontend/              ← Vue UI (Phase 3 — not created yet)
```

Always read `docs/architecture.md` before implementing domain logic.

## Security Rules (always apply)

1. **Never commit secrets** — no real tokens, passwords, or `.env` files with production values.
2. **Never log secret values** — redact env var values in all log output.
3. **Token scoping** — access tokens must be limited to explicit project/environment pairs.
4. **Minimal diffs** — do not refactor unrelated code in the same change.

## Conventions

### General

- **Code language:** English (identifiers, comments, commit messages).
- **Technical docs:** English (`docs/`, `AGENTS.md`).
- **README:** Portuguese (human-facing).
- **Commits:** Conventional commits — `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`.
- **Scope:** Keep changes focused. One concern per PR.

### Go Backend (planned — Phase 1+)

When `backend/` is created, follow these conventions:

- **Go version:** 1.22+
- **Layout:**
  ```
  backend/
  ├── cmd/server/main.go
  ├── internal/
  │   ├── handler/       # HTTP handlers
  │   ├── service/       # Business logic
  │   ├── repository/    # Database access
  │   └── model/         # Domain types
  └── pkg/               # Shared utilities (if needed)
  ```
- **Database:** PostgreSQL via `pgx`
- **HTTP router:** TBD (chi or echo — decide in Phase 1)
- **Error handling:** Wrap errors with context; never swallow errors silently
- **Tests:** Co-located `*_test.go` files; table-driven tests preferred

### Vue Frontend (planned — Phase 3+)

When `frontend/` is created, follow these conventions:

- **Vue 3** with Composition API and `<script setup lang="ts">`
- **Build:** Vite + TypeScript (strict mode)
- **State:** Pinia stores
- **Structure:**
  ```
  frontend/
  ├── src/
  │   ├── components/
  │   ├── composables/
  │   ├── views/
  │   ├── stores/
  │   └── api/
  ```
- **No secrets in client code** — the UI manages vault metadata; secret values are handled server-side only

## Skills

Project skills live in `.cursor/skills/`. Use them for domain knowledge and repeatable workflows:

| Skill | Use when |
|-------|----------|
| `env-vault-domain` | Working with Project, Environment, Token, Secret entities |
| `add-api-endpoint` | Adding a new Go API endpoint (Phase 2+) |
| `add-ui-screen` | Adding a new Vue screen (Phase 3+) |
| `security-review` | Reviewing auth, token, or secret-related changes |

## What NOT to Do in Phase 0

- Do not create `backend/`, `frontend/`, or application code without an explicit request.
- Do not add `go.mod`, `package.json`, or dependency files.
- Do not add CI workflows, Docker files, or Makefile until Phase 1.

## Testing (Phase 1+)

When code exists:

- Run tests before finishing any task.
- Add tests for new behavior; skip trivial assertions.
- Fix all test and lint failures before completing.

## Pull Request Guidelines

- Title: conventional commit format (e.g. `feat: add project CRUD endpoint`)
- Include a test plan describing what was verified.
- For security-related changes, confirm no secrets are included in the diff.

## Key References

- [Architecture](docs/architecture.md) — domain model and security principles
- [Roadmap](docs/roadmap.md) — project phases
- [Agent Guide](docs/agent-guide.md) — working with AI agents in this repo
- [Contributing](CONTRIBUTING.md) — contribution workflow
