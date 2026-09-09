# envm-vault — Frontend

Interface web do [envm-vault](../README.md): vault open source para gerenciar e distribuir variáveis de ambiente entre projetos e ambientes customizáveis.

Nesta fase, a UI cobre autenticação por e-mail (OTP), chave de API, projetos, times, tags, ambientes e a matriz de secrets em `/project/:slug`.

## Stack


| Item      | Tecnologia                                                        |
| --------- | ----------------------------------------------------------------- |
| Framework | Vue 3 (Composition API) + TypeScript                              |
| Build     | Vite                                                              |
| Estilo    | Tailwind CSS v4 + daisyUI 5                                       |
| Estado    | Pinia (`sessionStorage` para JWT)                                 |
| i18n      | vue-i18n (`pt-BR` / `en-US`, detecção do navegador)               |
| Testes    | Vitest + Vue Test Utils                                           |
| Runtime   | Node 24 (via [mise](https://mise.jdx.dev/) — `mise.toml` na raiz) |


## Estrutura

Arquitetura orientada a módulos:

- `src/app/` — bootstrap, router, layouts, styles, `translations/`
- `src/app/translations/{pt-BR,en-US}/` — JSON por contexto (`login.json`, `nav.json`, …)
- `src/shared/` — `HttpClient`, `envs`, erros compartilhados
- `src/modules/` — domínios (`auth`, `user`, `token`, `health`, `home`)
- `tests/unit/` — unitários (services), espelhando `src/`
- `tests/component/` — componentes (Vue Test Utils), espelhando `src/`; helpers em `component/helpers/`

## Desenvolvimento

```bash
# Na raiz do repositório
mise install
make frontend-dev   # http://localhost:3100 (API em :8080 via make dev)
```

Ou direto neste diretório:

```bash
cd frontend
npm install
npm run dev         # http://localhost:3100
```

O Vite faz proxy de `/v1` e `/health` para a API.

## Variáveis de ambiente

O frontend **não** usa o `.env` da API. Configure só variáveis públicas aqui:

| Arquivo | Conteúdo |
|---------|----------|
| `../api/.env` | Secrets da API (nunca `VITE_*`) |
| `../.deploy/db/.env` | PostgreSQL / Docker Compose (`POSTGRES_*`) |
| `.env` | Só `VITE_*` — vão para o bundle do browser |

Leitura e validação ficam centralizadas em `src/shared/env/envs.ts`:

```ts
import { envs } from '@/shared/env/envs'

console.log(envs.api.baseUrl)
```

Não leia `import.meta.env` fora desse módulo.

```bash
# Na raiz: make setup  (cria api/.env e frontend/.env a partir dos .env.example)
```

- Em desenvolvimento local, deixe `VITE_API_BASE_URL` vazio (proxy do Vite).
- Em staging/prod, defina a URL absoluta `http(s)://…` da API (e alinhe `CORS_ORIGINS` em `api/.env`).
- Na boot, `assertEnvsValid()` falha se `VITE_API_BASE_URL` estiver setada e não for URL absoluta.

## Testes

```bash
npm test              # unit + component
npm run test:unit
npm run test:component
npm run type-check
```

Convenções:

- Testes de componente usam `mountComponent` de `tests/component/helpers/mount.ts`
- **Todo** teste de componente deve chamar `expectComponentSnapshot(wrapper)` e versionar os arquivos em `__snapshots__/`
- Quando a UI muda de propósito: `npm run test:component -- -u` para atualizar os snapshots

Documentação geral do projeto: [README.md](../README.md).