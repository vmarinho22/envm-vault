# Agent Guide

This guide helps contributors — human or AI — work effectively in the envm-vault repository.

## Three Layers of Agent Context

| Layer | Location | Scope | Tool support |
|-------|----------|-------|--------------|
| **AGENTS.md** | Repository root | Project-wide instructions | Cursor, Copilot, Codex, Aider, and 20+ other agents |
| **RTK.md** | Repository root | Token-optimized CLI commands | Any agent running shell commands |
| **Cursor Rules** | `.cursor/rules/*.mdc` | File-pattern scoped rules | Cursor only |
| **Skills** | `.cursor/skills/*/SKILL.md` | Reusable workflows | Cursor (project skills) |

### When to use each

- **AGENTS.md** — Always the starting point. Portable markdown read by any coding agent.
- **RTK.md** — Prefix shell commands with `rtk` to reduce token usage (`rtk go test`, `rtk docker compose ps`).
- **Cursor Rules** — Automatic context when editing matching files (e.g. `api/**` triggers Go rules).
- **Skills** — Invoke explicitly or let the agent discover them by description. Use for domain knowledge and repeatable checklists.

## RTK Setup

[RTK](https://www.rtk-ai.app/) (Rust Token Killer) filters command output before it reaches the agent context.

```bash
brew install rtk              # or see rtk-ai.app
make rtk-trust                # trust project filters in .rtk/filters.toml
rtk gain                      # view token savings
```

Project filters cover: make targets and go test in `api/`.

## Getting Started with an Agent

Before asking an agent to implement anything, provide this context:

```
Read AGENTS.md, RTK.md, and docs/architecture.md before making any changes.
Use rtk prefix for shell commands (e.g. rtk go test ./...).
```

## Suggested Prompts

### Understanding the project

```
Summarize the envm-vault domain model based on docs/architecture.md.
What entities exist and how do they relate?
```

### Planning a feature (Phase 1+)

```
I want to implement [feature]. Read AGENTS.md, docs/architecture.md,
and the env-vault-domain skill. Propose an implementation plan
following project conventions before writing code.
```

### Security-sensitive work (Phase 2+)

```
Review this change for security issues. Use the security-review skill.
Focus on token handling, secret storage, and log redaction.
```

### Adding an API endpoint (Phase 2+)

```
Add a new API endpoint for [description]. Follow the add-api-endpoint skill
and Go conventions in AGENTS.md.
```

## Working Without Cursor

If you use Copilot, Codex, Aider, or another agent:

1. Point it to `AGENTS.md` — it works with any agent that reads markdown context files.
2. Reference `docs/architecture.md` for domain knowledge.
3. Skills in `.cursor/skills/` are plain markdown — copy relevant sections into your prompt if your tool doesn't auto-discover them.

## Conventions for Agent-Generated Code

When code exists (Phase 1+):

- **Language:** English for code, comments, and technical docs.
- **Commits:** Conventional commits (`feat:`, `fix:`, `docs:`, `chore:`).
- **Scope:** Minimal diffs — do not refactor unrelated code.
- **Security:** Never commit secrets, tokens, or real `.env` values.
- **Tests:** Add tests for new behavior; do not add trivial tests.

## Reporting Issues with Agent Output

If an agent produces incorrect or insecure code:

1. Do not merge it.
2. Open an issue describing what went wrong.
3. Consider updating `AGENTS.md`, a rule, or a skill to prevent recurrence.
