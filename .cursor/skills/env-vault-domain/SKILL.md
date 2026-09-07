---
name: env-vault-domain
description: Describes the envm-vault domain model including Project, Environment, Secret, User, AccessToken, Team, and Tag entities. Use when implementing features, designing APIs, writing migrations, or discussing vault behavior.
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
Account that owns projects and manages a unique API key via the admin UI. May be blocked (`blocked_at`) by admin/super_admin — denies login, JWT, and API key create/rotate; blocking revokes the active key.

### AccessToken
Programmatic credential for downloading env vars. **One active key per user.** Stored hashed (SHA-256); plaintext shown only at generate/rotate. Sent as `Authorization: Bearer <token>`. Access to a given project/environment is determined by the user's membership, not by scopes on the token.

### Team
Organizational unit with `name`, unique `slug` (from name via shared slug helper), and one accountable `User` (for future notifications). No memberships yet; not linked to Project yet. Any authenticated user may list/search; only `admin` / `super_admin` mutate.

### Tag
Global catalog label with unique `title` and hex `color` (`#RRGGBB`) for identifying projects later. Not linked to Project yet. Any authenticated user may list, create, update, or delete.

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
4. Slugs are generated with `api/internal/slug.FromName` (mirrored in `frontend/src/shared/slug`).

## Relationships

```
User → has many → Project
User → has one → AccessToken
User → accountable for many → Team
Project → has many → Environment
Environment → has many → Secret
Tag (catalog; Project link TBD)
```
