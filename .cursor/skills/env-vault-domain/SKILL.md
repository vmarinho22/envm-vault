---
name: env-vault-domain
description: Describes the envm-vault domain model including Project, Environment, Secret, User, AccessToken, Team, and Tag entities. Use when implementing features, designing APIs, writing migrations, or discussing vault behavior.
---

# envm-vault Domain

Read [docs/architecture.md](../../docs/architecture.md) for the full specification.

## Entities

### Project
Central hub with unique `title` and unique `slug` (business key for `/project/<slug>`). Many-to-many with Teams, Environments, and Tags. Any authenticated user may create/update; only `admin` / `super_admin` may delete. Owns Secrets.

### Environment
Global catalog of deployment targets with unique `title` (e.g. `dev`, `staging`, `production`). Linked to projects via M2M. Any authenticated user may list/search; only `admin` / `super_admin` mutate. Seeded defaults: `dev`, `test`, `staging`, `production`.

### Secret (EnvVar)
A named key scoped to a **Project** (`secrets` table). Values live in `secret_values` per Environment (sparse). Values encrypted at rest (AES-256-GCM, `ENCRYPTION_KEY`). Unique `(project_id, name)`.

### User
Account that owns projects and manages a unique API key via the admin UI. May be blocked (`blocked_at`) by admin/super_admin — denies login, JWT, and API key create/rotate; blocking revokes the active key.

### AccessToken
Programmatic credential for downloading env vars. **One active key per user.** Stored hashed (SHA-256); plaintext shown only at generate/rotate. Sent as `Authorization: Bearer <token>`. Access to a given project/environment will be determined by membership (TBD); currently any valid key may download.

### Team
Organizational unit with `name`, unique `slug` (from name via shared slug helper), and one accountable `User`. Linked to projects via M2M. Any authenticated user may list/search; only `admin` / `super_admin` mutate.

### Tag
Global catalog label with unique `title` and hex `color` (`#RRGGBB`). Linked to projects via M2M. Any authenticated user may list, create, update, or delete. Seeded defaults: `backend` `#000000`, `frontend` `#FFFFFF`, `infra` `#F97316`.

## Primary Flow

```
GET /v1/projects/{projectSlug}/env?environment={envTitle}
Authorization: Bearer <access_token>
→ Identifies the user, then returns env vars for that project+environment
→ environment=production requires role super_admin
```

## Key Rules

1. One active API key per user — generate rotates and invalidates the previous hash.
2. Secret values and API key plaintext never appear in logs.
3. Project `slug` is the unique natural key for UI routes (`/project/<slug>`); UUID is the PK for mutations and junctions.
4. Slugs are generated with `api/internal/slug.FromName` (mirrored in `frontend/src/shared/slug`).
5. Secret names match `^[A-Za-z_][A-Za-z0-9_]*$` (max 128).
6. Only `super_admin` may download the `production` environment (API key or JWT export).

## Relationships

```
User → has one → AccessToken
User → accountable for many → Team
Project ↔ Team (M2M)
Project ↔ Environment (M2M)
Project ↔ Tag (M2M)
Project → has many → Secret
Secret → has many → SecretValue (per Environment)
```
