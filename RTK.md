# RTK - Rust Token Killer

**Usage**: Token-optimized CLI proxy for shell commands in this repository.

## Rule

Always prefix shell commands with `rtk`. RTK has built-in filters for Go, Docker, Git, Make, and 80+ other tools. Unknown commands pass through unchanged.

```bash
# Good
rtk go test ./...
rtk go build ./...
rtk docker compose ps
rtk git status
rtk git diff

# Avoid (raw output wastes context)
go test ./...
docker compose logs
```

## envm-vault commands

Use RTK with project Makefile targets and API workflows:

```bash
rtk go test ./internal/...          # from api/
rtk go build ./cmd/server           # from api/
rtk docker compose up -d postgres   # from repo root
rtk docker compose logs postgres
```

Makefile shortcuts (already RTK-friendly where noted):

```bash
make test       # runs go test (use: cd api && rtk go test ./...)
make dev        # go run — keep raw for live server logs
```

## Meta commands

```bash
rtk gain            # token savings analytics
rtk gain --history  # recent command savings
rtk proxy <cmd>     # run without filtering (when you need full output)
rtk --version       # verify RTK is installed
```

## Custom project filters

Project-local filters live in [`.rtk/filters.toml`](.rtk/filters.toml). After editing, run:

```bash
rtk trust --yes
```

## Verification

```bash
rtk --version
which rtk
rtk gain
```

Install: https://www.rtk-ai.app/
