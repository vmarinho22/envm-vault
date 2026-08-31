# Contributing to envm-vault

Thank you for your interest in contributing! This guide covers the workflow and conventions for this project.

## Getting Started

1. Fork the repository.
2. Clone your fork: `git clone https://github.com/<you>/envm-vault.git`
3. Read [AGENTS.md](AGENTS.md) and [docs/architecture.md](docs/architecture.md) to understand the project.
4. Check [docs/roadmap.md](docs/roadmap.md) for the current phase and open work.

## Development Workflow

1. Create a branch from `main`: `git checkout -b feat/my-feature`
2. Make your changes following project conventions.
3. Commit using [conventional commits](https://www.conventionalcommits.org/):
   - `feat:` — new feature
   - `fix:` — bug fix
   - `docs:` — documentation only
   - `chore:` — maintenance, tooling
   - `refactor:` — code change that neither fixes a bug nor adds a feature
   - `test:` — adding or updating tests
4. Push and open a Pull Request against `main`.

## Pull Request Guidelines

- Use conventional commit format for the PR title.
- Fill in the PR template completely.
- For security-related changes, run through the security checklist in the template.
- Keep PRs focused — one concern per PR.

## Working with AI Agents

This repository is designed for AI-assisted development:

| Resource | Purpose |
|----------|---------|
| [AGENTS.md](AGENTS.md) | Project-wide agent instructions (works with any coding agent) |
| [.cursor/rules/](.cursor/rules/) | Cursor-specific scoped rules (Go, Vue) |
| [.cursor/skills/](.cursor/skills/) | Reusable workflow skills |
| [docs/agent-guide.md](docs/agent-guide.md) | Guide and suggested prompts |

If you use an agent, point it to `AGENTS.md` before asking it to implement anything.

## Code Conventions

- **Code language:** English (identifiers, comments, commit messages).
- **Technical docs:** English.
- **README:** Portuguese.
- **Scope:** Minimal diffs — do not refactor unrelated code.

## Security

- **Never commit secrets** — no real tokens, passwords, or production `.env` values.
- **Never include secrets in issues or PRs** — use placeholders.
- Report vulnerabilities via [SECURITY.md](SECURITY.md).

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). By participating, you agree to uphold its terms.

## Questions?

Open a [GitHub Discussion](https://github.com/envm-vault/envm-vault/discussions) or an issue with the `question` label.
