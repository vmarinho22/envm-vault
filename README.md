# envm-vault

Vault open source para gerenciamento e distribuição de variáveis de ambiente.

Com um token de acesso pessoal, desenvolvedores e pipelines de CI baixam variáveis de ambiente de projetos em ambientes customizáveis (`dev`, `tst`, `staging`, `prod`, ou qualquer nome que você definir).

## Status

**Fase 1** — API com autenticação por e-mail (OTP). UI admin em breve.

## Stack

| Camada | Tecnologia |
|--------|------------|
| API | Go + Fiber |
| Banco | PostgreSQL (GORM + AutoMigrate) |
| Auth | OTP por e-mail + JWT |
| UI Admin | Vue 3 + TypeScript + Vite (Fase 3) |

## Quick Start

Pré-requisitos: Go 1.25+, Docker, Make.

```bash
git clone https://github.com/envm-vault/envm-vault.git
cd envm-vault

make setup          # copia .env.example → .env
make db-up          # sobe PostgreSQL
make migrate        # roda AutoMigrate (GORM)
make dev            # inicia API em :8080
make hooks-install  # instala git hooks (testes unitários + conventional commits)
```

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

## Endpoints

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| `GET` | `/v1/health` | — | Health check |
| `GET` | `/health` | — | Health check (alias infra) |
| `POST` | `/v1/auth/login` | — | Solicita OTP por e-mail |
| `POST` | `/v1/auth/verify` | — | Valida OTP, retorna JWT |
| `POST` | `/v1/auth/register` | JWT (admin+) | Cadastra usuário |

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

**envm-vault** is an open-source environment variable vault. Built with Go (Fiber API) and PostgreSQL. Currently in Phase 1 with email OTP authentication. See [docs/architecture.md](docs/architecture.md) for details.
