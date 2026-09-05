# envm-vault

Vault open source para gerenciamento e distribuição de variáveis de ambiente.

Com um token de acesso pessoal, desenvolvedores e pipelines de CI baixam variáveis de ambiente de projetos em ambientes customizáveis (`dev`, `tst`, `staging`, `prod`, ou qualquer nome que você definir).

## Status

**Fase 1** — API com autenticação por e-mail (OTP), chave de API única por usuário e UI de login/onboarding.

## Stack

| Camada | Tecnologia |
|--------|------------|
| API | Go + Fiber |
| Banco | PostgreSQL (GORM + AutoMigrate) |
| Auth | OTP por e-mail + JWT (máx. 5 tentativas por código) |
| UI Admin | Vue 3 + Vite + TypeScript + Tailwind + daisyUI |

Testes da UI: Vitest (`make test-frontend`). Em testes de componente, sempre use `expectComponentSnapshot(wrapper)` e versione os `__snapshots__/`. Detalhes em [frontend/README.md](frontend/README.md).

Rate limit HTTP global por IP: `RATE_LIMIT_MAX` / `RATE_LIMIT_WINDOW` (padrão `100` / `1m`; `/health` isento).

## Quick Start

Pré-requisitos: [mise](https://mise.jdx.dev/) (Node 24 + Go 1.25), Docker, Make.

```bash
git clone https://github.com/envm-vault/envm-vault.git
cd envm-vault

mise install        # instala as versões de Node e Go definidas em mise.toml
make setup          # copia .env.example → .env e frontend/.env.example → frontend/.env
make db-up          # sobe PostgreSQL
make migrate        # roda AutoMigrate (GORM)
make dev            # inicia API em :8080
make frontend-dev   # inicia UI em :3100
make hooks-install  # instala git hooks (testes unitários + conventional commits)
```

Ambiente separado:

| Arquivo | Uso |
|---------|-----|
| `.env` (raiz) | API + Docker Compose (secrets: DB, JWT, CORS) |
| `frontend/.env` | Só variáveis públicas `VITE_*` (expostas no browser) |

A UI em `http://localhost:3100` faz proxy de `/v1` para a API. Em desenvolvimento, o código OTP aparece no log da API.

### Primeiro acesso (bootstrap)

```bash
# 1. Solicitar código OTP
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com"}'

# 2. Verificar código (veja o código no log do terminal em dev)
curl -X POST http://localhost:8080/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "code": "12345678"}'

# Resposta: JWT + user com role super_admin
```

### Cadastrar usuário (requer JWT de admin)

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt>" \
  -d '{"email": "user@example.com", "role": "user"}'
```

### Gerar chave de API (requer JWT)

```bash
curl -X POST http://localhost:8080/v1/tokens \
  -H "Authorization: Bearer <jwt>"
# Resposta: { "token": "ev_…", "prefix": "ev_a1b2c3d4", "created_at": "…" }
# O valor completo só aparece nesta resposta.
```

## Endpoints

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| `GET` | `/v1/health` | — | Health check |
| `GET` | `/health` | — | Health check (alias infra) |
| `POST` | `/v1/auth/login` | — | Solicita OTP por e-mail |
| `POST` | `/v1/auth/verify` | — | Valida OTP, retorna JWT |
| `POST` | `/v1/auth/register` | JWT (admin+) | Cadastra usuário |
| `GET` | `/v1/tokens/me` | JWT | Status da chave de API (sem o secret) |
| `POST` | `/v1/tokens` | JWT | Gera ou rotaciona a chave (plaintext uma vez) |
| `DELETE` | `/v1/tokens` | JWT | Revoga a própria chave |
| `GET` | `/v1/users` | JWT (admin+) | Lista usuários e status da chave |
| `DELETE` | `/v1/users/:id/token` | JWT (admin+) | Revoga a chave de outro usuário |

## Agentes de IA e RTK

Este repositório é otimizado para contribuições assistidas por agentes:

- **[AGENTS.md](AGENTS.md)** — instruções portáveis para coding agents
- **[RTK.md](RTK.md)** — prefixe comandos shell com `rtk` para economizar tokens
- **[`.rtk/filters.toml`](.rtk/filters.toml)** — filtros customizados (make, go test)

Setup RTK (recomendado para quem usa Cursor/Copilot/Codex):

```bash
brew install rtk          # ou: https://www.rtk-ai.app/
make rtk-trust            # confia nos filtros do projeto
make test-rtk             # testes com output compacto
```

Veja [docs/agent-guide.md](docs/agent-guide.md) para mais detalhes.

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [Arquitetura](docs/architecture.md) | Modelo de domínio, fluxos e segurança |
| [Roadmap](docs/roadmap.md) | Fases do projeto |
| [Guia de agentes](docs/agent-guide.md) | Como usar agentes de IA |
| [AGENTS.md](AGENTS.md) | Instruções para coding agents |
| [RTK.md](RTK.md) | Comandos shell otimizados para agentes |
| [Contribuindo](CONTRIBUTING.md) | Como contribuir |
| [Segurança](SECURITY.md) | Reporte de vulnerabilidades |

## Contribuindo

Contribuições são bem-vindas! Leia [CONTRIBUTING.md](CONTRIBUTING.md) para começar.

## Licença

[MIT](LICENSE)

## English Summary

**envm-vault** is an open-source environment variable vault. Built with Go (Fiber API), Vue 3, and PostgreSQL. Currently in Phase 1 with email OTP authentication and a unique API key per user. See [docs/architecture.md](docs/architecture.md) for details.
