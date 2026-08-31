---
name: add-api-endpoint
description: Checklist for adding a new Go API endpoint in envm-vault. Use when creating HTTP handlers, routes, or REST endpoints in the backend.
---

# Add API Endpoint

Use this checklist when adding a new endpoint to the Go API (Phase 2+).

## Before coding

1. Read [docs/architecture.md](../../docs/architecture.md) for domain context.
2. Confirm the endpoint fits the REST pattern: `/v1/{resource}`.
3. Identify required auth: access token, session, or public.

## Implementation checklist

```
- [ ] Define request/response types in internal/model/
- [ ] Add repository method in internal/repository/
- [ ] Add service method in internal/service/ (business logic)
- [ ] Add handler in internal/handler/ (input validation, HTTP response)
- [ ] Register route in cmd/server/main.go or router setup
- [ ] Add integration test
- [ ] Update docs if endpoint is user-facing
```

## Handler pattern

```go
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "project")
    if slug == "" {
        respondError(w, http.StatusBadRequest, "project slug required")
        return
    }

    project, err := h.service.GetProject(r.Context(), slug)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "failed to get project")
        return
    }

    respondJSON(w, http.StatusOK, project)
}
```

## Security checks

- [ ] Input validated and sanitized
- [ ] Auth middleware applied if required
- [ ] Token scope verified for secret-related endpoints
- [ ] No secret values in response unless explicitly requested (download endpoint)
- [ ] No secret values in error messages or logs
