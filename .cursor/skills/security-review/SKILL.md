---
name: security-review
description: Reviews code changes for security issues in envm-vault, focusing on authentication, token handling, secret storage, and log redaction. Use when reviewing PRs or changes involving auth, tokens, or secrets.
---

# Security Review

Use this checklist when reviewing changes that touch authentication, tokens, or secrets.

## Authentication & Tokens

- [ ] Access tokens are hashed before storage (never store plaintext tokens)
- [ ] Token scope is enforced — request for `project/env` rejects out-of-scope tokens
- [ ] Expired and revoked tokens are rejected
- [ ] Token values are never returned in API responses after creation (show once pattern)

## Secret Handling

- [ ] Secret values are encrypted at rest (Phase 2+)
- [ ] Secret values never appear in logs, error messages, or stack traces
- [ ] List endpoints do not return decrypted secret values
- [ ] Download endpoint requires valid scoped token

## Input Validation

- [ ] All user input is validated (slug format, key names, lengths)
- [ ] SQL injection prevented (parameterized queries via pgx)
- [ ] Path traversal prevented in project/environment slugs

## Repository & Config

- [ ] No real secrets, tokens, or credentials in the diff
- [ ] No `.env` files with production values committed
- [ ] Environment-specific config loaded from env vars, not hardcoded

## Reporting

Format findings as:
- **Critical** — must fix before merge (secret exposure, missing auth)
- **Warning** — should fix (missing validation, weak scoping)
- **Note** — consider improving (logging, error messages)

Reference: [docs/architecture.md — Security Principles](../../docs/architecture.md#security-principles)
