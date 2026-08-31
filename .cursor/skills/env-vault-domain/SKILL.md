---
name: env-vault-domain
description: Describes the envm-vault domain model including Project, Environment, Secret, User, and AccessToken entities. Use when implementing features, designing APIs, writing migrations, or discussing vault behavior.
---

# envm-vault Domain

Read [docs/architecture.md](../../docs/architecture.md) for the full specification.

## Entities

### Project
Groups environment variables under a slug (e.g. `my-app`). Owned by a User.

### Environment
A named deployment target within a Project (e.g. `dev`, `staging`, `prod`). Fully customizable — no fixed set.

### Secret (EnvVar)
A key-value pair scoped to Project + Environment. Values encrypted at rest (Phase 2).

### User
Account that owns projects and manages access tokens via the admin UI.

### AccessToken
Programmatic credential for downloading env vars. Scoped to specific project/environment pairs. Sent as `Authorization: Bearer <token>`.

## Primary Flow

```
GET /v1/projects/{project}/environments/{env}
Authorization: Bearer <access_token>
→ Returns env vars for that project+environment
```

## Key Rules

1. Tokens are scoped — a token for `my-app/dev` cannot access `my-app/prod`.
2. Secret values never appear in logs.
3. Environments are per-project — `dev` in project A is independent from `dev` in project B.

## Relationships

```
User → has many → Project
User → has many → AccessToken
Project → has many → Environment
Environment → has many → Secret
```
