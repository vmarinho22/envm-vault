# envm-vault — Frontend

Interface web do [envm-vault](../README.md): vault open source para gerenciar e distribuir variáveis de ambiente entre projetos e ambientes customizáveis.

Nesta fase, a UI cobre autenticação por e-mail (OTP) e o shell autenticado (menu + início). Telas de projetos, secrets e tokens virão nas próximas fases.

## Stack


| Item      | Tecnologia                                                        |
| --------- | ----------------------------------------------------------------- |
| Framework | Vue 3 (Composition API) + TypeScript                              |
| Build     | Vite                                                              |
| Estilo    | Tailwind CSS v4 + daisyUI 5                                       |
| Estado    | Pinia (`sessionStorage` para JWT)                                 |
| Testes    | Vitest + Vue Test Utils                                           |
| Runtime   | Node 24 (via [mise](https://mise.jdx.dev/) — `mise.toml` na raiz) |


## Estrutura

Arquitetura orientada a módulos:

- `src/app/` — bootstrap, router, layouts, styles
- `src/shared/` — `HttpClient`, erros compartilhados
- `src/modules/` — domínios (`auth`, `user`, `health`, `home`)
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