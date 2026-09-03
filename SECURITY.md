# Security Policy

## Reporting a Vulnerability

envm-vault handles sensitive data (environment variables and access tokens). We take security seriously.

If you discover a security vulnerability, please report it responsibly:

1. **Do NOT open a public GitHub issue** for security vulnerabilities.
2. Email the maintainers at **[INSERT SECURITY EMAIL]** with:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)
3. We will acknowledge receipt within 48 hours and provide an estimated timeline for a fix.

## What to Report

- Authentication or authorization bypass
- Token scope enforcement failures
- Secret exposure (in logs, responses, or storage)
- Encryption weaknesses
- Injection vulnerabilities (SQL, path traversal)

## What NOT to Include in Reports

- Real access tokens or credentials
- Production environment variable values
- Personal data of other users

Use placeholders and describe the issue conceptually.

## Supported Versions

| Version | Supported |
|---------|-----------|
| main (latest) | Yes |

## Security Best Practices for Contributors

- Never commit real secrets, tokens, or `.env` files with production values.
- Never log secret values — redact in all output.
- **OTP codes are logged to stdout only when `APP_ENV=development`.** Never enable this in production.
- Use the [security-review skill](.cursor/skills/security-review/SKILL.md) when reviewing auth or secret-related changes.
- Follow the security principles in [docs/architecture.md](docs/architecture.md#security-principles).

## Public Issues

For non-sensitive security improvements (e.g. adding security headers, improving input validation), feel free to open a regular GitHub issue or PR.
