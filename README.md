# envm-vault

Vault open source para gerenciar e distribuir variáveis de ambiente.

Crie projetos, anexe secrets de um catálogo global, defina valores por ambiente e baixe um `.env` com um token pessoal (`ev_…`) — no CI, no terminal ou no painel admin.

## O que você pode fazer

- **Projetos** com times, tags e ambientes (`dev`, `staging`, `production`, …)
- **Catálogo global de secrets** compartilhados entre projetos (valores criptografados em repouso)
- **Prefixos por secret × ambiente** no download (ex.: `API_URL` → `VITE_API_URL`)
- **UI admin** (Vue) para login OTP, chave de API, projetos e secrets
- **Download programático** via API key para pipelines e desenvolvedores

Detalhes de domínio e segurança: [docs/architecture.md](docs/architecture.md).

## Stack

| Camada | Tecnologia |
|--------|------------|
| API | Go + Fiber |
| Banco | PostgreSQL (GORM + AutoMigrate) |
| Auth | OTP por e-mail + JWT · chave de API (`ev_…`) |
| UI | Vue 3 + Vite + TypeScript + Tailwind + daisyUI |

## Quick Start

Pré-requisitos: [mise](https://mise.jdx.dev/) (Node + Go), Docker, Make.

```bash
git clone https://github.com/envm-vault/envm-vault.git
cd envm-vault

mise install        # versões de Node e Go (mise.toml)
make setup          # copia .env.example → .env se faltarem
make db-up          # PostgreSQL
make migrate        # AutoMigrate
make dev            # API em :8080
make frontend-dev   # UI em :3100
make hooks-install  # git hooks (testes + conventional commits)
```

| Arquivo | Uso |
|---------|-----|
| `api/.env` | Secrets da API (`DATABASE_URL`, `JWT_*`, `ENCRYPTION_KEY`, CORS, rate limit) |
| `.docker/db/.env` | PostgreSQL / Docker Compose (`POSTGRES_*`) |
| `frontend/.env` | Só `VITE_*` públicas (expostas no browser) |

A UI em `http://localhost:3100` faz proxy de `/v1` para a API. Em desenvolvimento, o OTP aparece no log da API.

Gere `ENCRYPTION_KEY` (32 bytes hex) se ainda não existir:

```bash
openssl rand -hex 32
```

### Primeiro acesso

1. Abra `http://localhost:3100/login` e entre com um e-mail (o primeiro usuário vira `super_admin`).
2. Gere sua chave de API em **Chave de API** (mostrada uma vez).
3. Crie projetos, ambientes e secrets no painel.

Via curl (bootstrap):

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com"}'

# Código OTP no log da API (APP_ENV=development)
curl -X POST http://localhost:8080/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "code": "12345678"}'
```

### Baixar variáveis (API key)

No detalhe do projeto a UI mostra o snippet pronto. Exemplo:

```bash
export ENVM_TOKEN="ev_…"

curl -fsSL \
  -H "Authorization: Bearer ${ENVM_TOKEN}" \
  "http://localhost:8080/v1/projects/<slug>/env?environment=dev&format=env" \
  -o .env
```

Só `super_admin` pode baixar o ambiente `production`.

## Testes

```bash
make test           # API + frontend
make test-frontend  # só Vitest
```

Em testes de componente Vue, use `expectComponentSnapshot(wrapper)` e versione `__snapshots__/`. Ver [frontend/README.md](frontend/README.md).

## Agentes de IA

Este repositório é otimizado para contribuições assistidas por agentes:

- **[AGENTS.md](AGENTS.md)** — instruções para coding agents
- **[RTK.md](RTK.md)** — prefixe comandos shell com `rtk` para economizar tokens
- **[docs/agent-guide.md](docs/agent-guide.md)** — guia completo

```bash
brew install rtk          # ou: https://www.rtk-ai.app/
make rtk-trust
make test-rtk
```

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [Arquitetura](docs/architecture.md) | Domínio, API, segurança |
| [Guia de agentes](docs/agent-guide.md) | Como usar agentes de IA |
| [AGENTS.md](AGENTS.md) | Convenções para agents |
| [RTK.md](RTK.md) | CLI otimizada para tokens |
| [Contribuindo](CONTRIBUTING.md) | Como contribuir |
| [Segurança](SECURITY.md) | Reporte de vulnerabilidades |

## Contribuindo

Contribuições são bem-vindas! Leia [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença

[MIT](LICENSE)

## English Summary

**envm-vault** is an open-source environment variable vault: shared secret catalog, per-project environments, optional download-time name prefixes, Go (Fiber) API, Vue 3 admin UI, and personal API keys for CI/developer downloads. See [docs/architecture.md](docs/architecture.md).
