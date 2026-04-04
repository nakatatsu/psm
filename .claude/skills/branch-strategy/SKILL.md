---
name: branch-strategy
description: >
  GitFlow branching strategy guide. Use when creating branches, making commits,
  merging PRs, or when the user asks about branch naming, merge strategy, or release flow.
  Also use when the user says "new branch", "create feature", "release", or "hotfix".
---

# branch-strategy — GitFlow Branching Strategy

This project follows GitFlow with the following conventions.

## Branch Types

| Branch | Pattern | Purpose | Protection |
|--------|---------|---------|------------|
| Production | `main` | Production-ready code | PR required (1 approval), CI required |
| Development | `develop` | Integration branch | PR required (no approval), CI required |
| Feature | `feature/<issue-no>-<short-description>` | New features | No protection |
| Release | `release-<version>` | Release preparation | PR required (1 approval), CI required |
| Hotfix | `hotfix-<version>` | Urgent production fixes | PR required (1 approval), CI required |

## Branch Naming Rules

- **Feature**: `feature/<issue-no>-<short-description>` (e.g., `feature/45-branch-strategy`)
- **Release**: `release-<version>` (e.g., `release-1.2.0`)
- **Hotfix**: `hotfix-<version>` (e.g., `hotfix-1.2.1`)
- Use lowercase, hyphens for separators, no underscores

## Merge Strategy

- **Merge commit only** — squash and rebase are disabled
- All merges go through pull requests

## GitFlow Workflow

### Feature Development
1. Create `feature/<issue-no>-<short-description>` from `develop`
2. Work on the feature, commit with `#<issue-no>` in messages
3. Create PR to `develop`, CI must pass
4. Merge (merge commit)

### Release
1. Create `release-<version>` from `develop`
2. Final adjustments and QA on the release branch
3. Create PR to `main` (1 approval required, CI must pass)
4. Merge to `main`
5. Admin tags `main` with `v<version>` → triggers release workflow
6. Create PR from `release-<version>` back to `develop`

### Hotfix
1. Create `hotfix-<version>` from `main`
2. Fix the issue
3. Create PR to `main` (1 approval required, CI must pass)
4. Merge to `main`
5. Admin tags `main` with `v<version>` → triggers release workflow
6. Create PR from `hotfix-<version>` back to `develop`

## Commit Messages

- **Must include a GitHub Issue number** (e.g., `#45`)
- Enforced by `.githooks/commit-msg` hook
- Setup: `git config core.hooksPath .githooks`
- Merge commits (auto-generated) are exempt

## Important

- Never commit directly to `main` or `develop`
- Never delete `develop`, `release-*`, or `hotfix-*` branches (deletion protected)
- Never force-push to any protected branch
- Admin can bypass protection rules for emergencies only
