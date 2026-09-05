---
name: add-api-endpoint
description: Checklist for adding a new Go API endpoint in envm-vault. Use when creating HTTP handlers, routes, or REST endpoints in the backend.
---

# Add API Endpoint

Use this checklist when adding a new endpoint to the Go API.

## Before coding

1. Read [docs/architecture.md](../../docs/architecture.md) for domain context.
2. Confirm the endpoint fits under `/v1/{context}/...`.
3. Identify required auth: access token, JWT session, or public.
4. Choose the route context: `health`, `auth`, `user`, `token`, or create a new `router/{context}.go`.

## Implementation checklist

```
- [ ] Add GORM model in internal/model/ (if new table)
- [ ] Add repository methods in internal/repository/
- [ ] Add service method in internal/service/{context}/
- [ ] Add handler method in internal/handler/{context}/
- [ ] Register route in internal/router/{context}.go
- [ ] Map errors via httperr.Write(c, err)
- [ ] Add tests under `api/tests/unit/` mirroring `internal/`
- [ ] Update docs if endpoint is user-facing
```

## Route registration pattern

```go
// internal/router/project.go
func registerProjectRoutes(v1 fiber.Router, deps Deps) {
    g := v1.Group("/projects", deps.JWT.Protect())
    g.Get("/", deps.Project.List)
    g.Post("/", deps.Project.Create)
}
```

Wire in `router.go`:

```go
registerProjectRoutes(v1, deps)
```

## Security checks

- [ ] Input validated and sanitized
- [ ] Auth middleware applied if required
- [ ] API key identifies the user for secret-related endpoints (membership checked when projects exist)
- [ ] No secret values in response unless explicitly requested
- [ ] No secret values in error messages or logs
