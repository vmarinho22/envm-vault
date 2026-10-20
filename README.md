# envm-vault

![Screenshot of envm-vault](https://github.com/vmarinho22/envm-vault/blob/main/docs/envm-vault.png)



Vault open source para gerenciar e distribuir variáveis de ambiente.

Crie projetos, anexe secrets de um catálogo global, defina valores por ambiente e baixe um `.env` com um token pessoal (`ev_…`).

Detalhes técnicos: [docs/architecture.md](docs/architecture.md).

## Sumário

1. [Como usar](#como-usar)
2. [Stack](#stack)
3. [Arquivos `.env` (detalhes)](#arquivos-env-detalhes)
4. [Outras formas de rodar](#outras-formas-de-rodar)
5. [Download de env (API key)](#download-de-env-api-key)
6. [Testes e agentes](#testes-e-agentes)
7. [Documentação](#documentação)
8. [Licença](#licença)

## Como usar

Você **não precisa clonar** o repositório. Basta criar uma pasta com os `.env`, baixar o `compose.yaml` e subir a imagem pronta.

### 1. Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) (ou Podman) com Docker Compose

### 2. Criar a pasta e os arquivos `.env`

```bash
mkdir -p envm-vault/.deploy/db envm-vault/.deploy/api
cd envm-vault

# Postgres
curl -fsSL https://raw.githubusercontent.com/vmarinho22/envm-vault/main/.deploy/db/.env.example \
  -o .deploy/db/.env

# App (API + UI)
curl -fsSL https://raw.githubusercontent.com/vmarinho22/envm-vault/main/.deploy/api/.env.example \
  -o .deploy/api/.env
```

Os exemplos servem para testar em local. Em produção, edite `.deploy/api/.env` e troque pelo menos `JWT_SECRET` e `ENCRYPTION_KEY` (veja [Arquivos `.env`](#arquivos-env-detalhes)).

### 3. Baixar o Compose e subir

```bash
curl -fsSL https://raw.githubusercontent.com/vmarinho22/envm-vault/main/compose.yaml \
  -o compose.yaml

docker compose up -d
```

Isso puxa a imagem `ghcr.io/vmarinho22/envm-vault:latest`, sobe o Postgres e a app.

Abra: **http://localhost:8080**

### 4. Fazer login

1. Na UI, informe seu e-mail e peça o código OTP.
2. Com `APP_ENV=development` (padrão do exemplo), o código aparece no log do container:

```bash
docker compose logs -f app
# procure por: login code sent … code=……
```

3. Cole o código e entre. O **primeiro usuário** vira `super_admin`.

### 5. Usar o vault (resumo)

1. Crie ou use um **projeto** e vincule **ambientes** (ex.: `dev`).
2. Anexe **secrets** e preencha os valores por ambiente.
3. Em **Tokens**, gere uma chave pessoal (`ev_…`) — copie na hora; ela só aparece uma vez.
4. Baixe as variáveis:

```bash
export ENVM_TOKEN="ev_…"   # cole a chave gerada
curl -fsSL \
  -H "Authorization: Bearer ${ENVM_TOKEN}" \
  "http://localhost:8080/v1/projects/<slug-do-projeto>/env?environment=dev&format=env" \
  -o app.env
```

Troque `<slug-do-projeto>` pelo slug do seu projeto (ex.: `meu-app`).

### 6. Parar

```bash
docker compose down
```

---

Quer construir a imagem localmente, contribuir com código ou usar hot reload? Veja [Outras formas de rodar](#outras-formas-de-rodar).

## Stack

| Camada | Tecnologia |
|--------|------------|
| API | Go + Fiber |
| Banco | PostgreSQL (GORM + AutoMigrate) |
| Auth | OTP por e-mail + JWT · chave de API (`ev_…`) |
| UI | Vue 3 + Vite + TypeScript + Tailwind + daisyUI |

## Arquivos `.env` (detalhes)

Não commite valores reais. Em produção, gere chaves próprias:

```bash
openssl rand -hex 32   # cole o resultado em ENCRYPTION_KEY
```

| Arquivo | Precisa? | Para quê |
|---------|----------|----------|
| `.deploy/db/.env` | Sim (Compose com Postgres) | `POSTGRES_*` |
| `.deploy/api/.env` | Sim | Runtime da app |
| `.deploy/frontend/.env` | Só no **build** da imagem a partir do código | `VITE_*` (vazio = same-origin `/v1`) |

Variáveis **obrigatórias** em `.deploy/api/.env`:

| Variável | Notas |
|----------|--------|
| `DATABASE_URL` | No Compose o host é `postgres` (nome do serviço) |
| `JWT_SECRET` | Troque o valor de exemplo |
| `ENCRYPTION_KEY` | 64 caracteres hex (32 bytes). **Não mude** depois de gravar secrets — senão o download falha ao descriptografar |

Opcionais (já têm default): `APP_ENV`, `PORT`, `JWT_EXPIRY`, `CORS_ORIGINS`, `RATE_LIMIT_*`, `ENCRYPTION_KEY_VERSION`. Na imagem publicada, `FRONTEND_DIST=/app/ui` já vem definida.

## Outras formas de rodar

### Só a app (Postgres externo)

Se você já tem um PostgreSQL:

```bash
docker run --rm -p 8080:8080 \
  -e DATABASE_URL='postgres://USER:PASS@HOST:5432/DB?sslmode=disable' \
  -e JWT_SECRET='troque-isto' \
  -e ENCRYPTION_KEY="$(openssl rand -hex 32)" \
  -e APP_ENV=production \
  -e CORS_ORIGINS='["http://localhost:8080"]' \
  ghcr.io/vmarinho22/envm-vault:latest
```

### Clonar e construir localmente

```bash
git clone https://github.com/vmarinho22/envm-vault.git
cd envm-vault
make setup
make up          # build local → envm-vault:local
make logs
make down
```

### Desenvolvimento — hot reload

API e UI em processos separados (precisa de [mise](https://mise.jdx.dev/) ou Go/Node):

```bash
git clone https://github.com/vmarinho22/envm-vault.git
cd envm-vault
mise install
make setup
make db-up && make migrate
make dev              # API :8080
make frontend-dev     # UI  :3100
```

Arquivos nesse modo: `api/.env`, `frontend/.env`, `.deploy/db/.env` (ou `make setup`).

## Download de env (API key)

```bash
export ENVM_TOKEN="ev_…"
curl -fsSL \
  -H "Authorization: Bearer ${ENVM_TOKEN}" \
  "http://localhost:8080/v1/projects/<slug>/env?environment=dev&format=env" \
  -o app.env
```

Só o papel `super_admin` pode baixar o ambiente `production`.

## Testes e agentes

```bash
make test
brew install rtk && make rtk-trust && make test-rtk
```

- [AGENTS.md](AGENTS.md) · [RTK.md](RTK.md) · [docs/agent-guide.md](docs/agent-guide.md)

## Documentação

| Documento | Descrição |
|-----------|-----------|
| [Arquitetura](docs/architecture.md) | Domínio, API, segurança |
| [Contribuindo](CONTRIBUTING.md) | Como contribuir |
| [Segurança](SECURITY.md) | Vulnerabilidades |

## Licença

[MIT](LICENSE)

## English Summary

**envm-vault** — environment variable vault. No clone required: create `.deploy/db/.env` + `.deploy/api/.env`, download `compose.yaml`, run `docker compose up -d` (pulls `ghcr.io/vmarinho22/envm-vault`). Dev/contributors: `make setup` + `make up`. See [docs/architecture.md](docs/architecture.md).
