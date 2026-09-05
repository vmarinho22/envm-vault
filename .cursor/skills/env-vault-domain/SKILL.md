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
Account that owns projects and manages a unique API key via the admin UI.

### AccessToken
Programmatic credential for downloading env vars. **One active key per user.** Stored hashed (SHA-256); plaintext shown only at generate/rotate. Sent as `Authorization: Bearer <token>`. Access to a given project/environment is determined by the user's membership, not by scopes on the token.

## Primary Flow

```
GET /v1/projects/{project}/environments/{env}
Authorization: Bearer <access_token>
→ Identifies the user, then returns env vars for that project+environment
```

## Key Rules

1. One active API key per user — generate rotates and invalidates the previous hash.
2. Secret values and API key plaintext never appear in logs.
3. Environments are per-project — `dev` in project A is independent from `dev` in project B.

## Relationships

```
User → has many → Project
User → has one → AccessToken
Project → has many → Environment
Environment → has many → Secret
```
