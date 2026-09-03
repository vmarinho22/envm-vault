#!/usr/bin/env bash
set -euo pipefail

commit_msg_file="${1:?usage: validate-commit-msg.sh <commit-msg-file>}"

# Read first non-comment, non-empty line (subject)
subject=""
while IFS= read -r line || [[ -n "$line" ]]; do
  if [[ "$line" =~ ^[[:space:]]*# ]]; then
    continue
  fi
  if [[ -z "${line// /}" ]]; then
    continue
  fi
  subject="$line"
  break
done < "$commit_msg_file"

if [[ -z "$subject" ]]; then
  echo "commit-msg: empty commit message" >&2
  exit 1
fi

# Allow merge/revert commits from git
if [[ "$subject" =~ ^Merge\  ]] || [[ "$subject" =~ ^Revert\  ]]; then
  exit 0
fi

pattern='^(feat|fix|docs|chore|refactor|test)(\([a-z0-9._/-]+\))?!?: .+'

if [[ "$subject" =~ $pattern ]]; then
  exit 0
fi

cat >&2 <<EOF
commit-msg: message does not follow conventional commits

  Subject: $subject

Expected format:
  <type>(<optional-scope>): <description>

Allowed types: feat, fix, docs, chore, refactor, test

Examples:
  feat: add project CRUD endpoint
  fix(api): handle invalid OTP code
  docs: update architecture diagram
EOF

exit 1
