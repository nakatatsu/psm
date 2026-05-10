---
name: pr-lint
description: >
  Lint a PR for project-rule compliance immediately after it is created or updated, and
  auto-fix violations. Runs unconditionally right after `gh pr create` or
  `gh pr edit --body`, and always before reporting "PR created" / "PR updated" to the
  user. Current rules: (1) strip issue auto-close keywords (Closes / Fixes / Resolves
  etc.) from the PR body so the linked issue is not auto-closed on merge; (2) detect Test
  plan items that cannot be verified pre-merge (post-merge / post-release observations)
  and move them to the linked issue, since the PR Test plan is for approval-time
  verification only.
---

# pr-lint — PR Project-Rule Linter

This skill lints a PR that was just created or updated against project rules, auto-fixing
the two recurring mistakes
before the user sees the PR. The PR is the artifact a reviewer uses to decide "approve or
not", so the body must be free of (a) hooks that close the linked issue prematurely and
(b) checks that the reviewer cannot actually perform.

## When to Trigger

Run unconditionally immediately after `gh pr create` or `gh pr edit --body`, and always
before reporting "PR created" / "PR updated" to the user. Apply both to brand-new PRs and
to body edits on existing PRs.

## Procedure Overview

1. Fetch the PR body and the linked issue number.
2. Run Check A (auto-close keywords). Stage a rewrite if needed.
3. Run Check B (un-verifiable test plan items). Stage a rewrite + an issue-side append if
   needed.
4. If anything was staged, apply with a single `gh pr edit --body-file` (and a
   corresponding `gh issue edit` for items moved out).
5. Report only what changed; if nothing changed, stay silent.

---

## Check A — Strip Auto-Close Keywords

### Why

Issues in this project are closed manually by a human after verification is complete
(post-merge to `develop`, and post-release for release-related work). GitHub's auto-close
keywords close the linked issue the moment the PR is merged into the default branch —
before verification has happened — so verification is silently skipped.

GitHub treats these keywords (case-insensitive) as auto-close triggers when paired with an
issue reference:

```
close, closes, closed, fix, fixes, fixed, resolve, resolves, resolved
```

The trigger fires when the keyword is followed by `#<number>`, `<owner>/<repo>#<number>`,
or a full issue URL.

### Detect

```bash
gh pr view <PR-number> --json body --jq .body > /tmp/pr_body.txt
grep -inE '\b(close[sd]?|fix(e[sd])?|resolve[sd]?)\b[[:space:]]+(#[0-9]+|[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+#[0-9]+|https?://github\.com/[^/]+/[^/]+/issues/[0-9]+)' /tmp/pr_body.txt
```

If there are no matches, skip to Check B.

### Rewrite

Replace the keyword (and only the keyword) with `Refs`. The issue reference itself is
preserved so the PR-issue link is kept and `Refs #N` still renders as an auto-link.

- `Closes #70` → `Refs #70`
- `Fixes #70` → `Refs #70`
- `Resolves #70` → `Refs #70`
- Lowercase / past-tense variants (`close`, `closed`, `fix`, `fixed`, `resolve`,
  `resolved`) → `Refs #70`
- `owner/repo#70` and full-URL forms: replace only the keyword portion the same way

---

## Check B — Move Un-Verifiable Test Plan Items to the Issue

### Why

A PR Test plan exists so the reviewer can decide "approve / not approve". Items that the
reviewer cannot perform before merging create a deadlock: approval requires verification,
verification requires merge, merge requires approval. Such items belong in the linked
issue under a "post-merge verification" section, where they are picked up after merge and
gate the issue close instead of the PR approve.

The PR Test plan must contain only items verifiable before merge. Anything else moves to
the issue.

### Detect

Scan the Test plan section (a heading containing "Test plan" or "テスト" near a checklist
of `- [ ]` items) for phrases that signal "cannot be verified at PR-review time":

| Pattern (case-insensitive) | Reason it cannot be verified pre-merge |
|----------------------------|----------------------------------------|
| `マージ後` / `merge 後` / `after merge` / `post-merge` / `post merge` | Requires the merge to happen first |
| `develop merge` / `default branch` / `default ブランチ` | Requires landing on the default / target branch |
| `次回リリース` / `next release` / `on release` | Requires a future release event |
| `Dependabot` (in a check item)| Dependabot only reads config from default branch |
| Mention of a workflow trigger that does not fire on PR (e.g. `pull_request: closed`, `push: develop`, `release-*` push) used as a verification step | The trigger does not fire during PR review |

If there are no such items, skip to Apply step.

### Reclassify

Split the Test plan into two parts:

- **Pre-merge verifiable** (stays in PR): items confirmable from the PR's CI runs, the PR
  diff, or by observation on the feature branch.
- **Post-merge / post-release** (moves out): everything matched by the patterns above.

### Find the linked issue

The issue number is taken from the PR body. Look for `#<number>` references; if multiple,
ask the user. If none, skip Check B and report that the items were detected but no issue
to move them to.

### Rewrite

- PR body: keep only the pre-merge items in the Test plan. If a Test plan section becomes
  empty, replace it with a single line noting that all verification is post-merge and
  pointing to the issue (`See #N for verification checklist`).
- Issue body: append (or update if present) a `## 動作確認` / `## Verification` section
  containing the moved items. If the section already exists, merge intelligently —
  deduplicate, keep existing checklist state.

Apply with:

```bash
gh pr edit <PR-number> --body-file /tmp/pr_body_fixed.txt
gh issue edit <issue-number> --body-file /tmp/issue_body_fixed.txt
```

**Always check the exit code and verify the change took effect**, because `gh pr edit` can
exit non-zero on this repo due to a deprecated-Projects-classic GraphQL warning even when
the user thinks it succeeded — and in that failure mode the body is **not** updated. If
`gh pr edit` returns non-zero, fall back to the REST API which bypasses the Projects
field:

```bash
# Fallback for gh pr edit
jq -Rs '{body: .}' < /tmp/pr_body_fixed.txt > /tmp/payload.json
gh api -X PATCH repos/<owner>/<repo>/pulls/<PR-number> --input /tmp/payload.json

# Fallback for gh issue edit
jq -Rs '{body: .}' < /tmp/issue_body_fixed.txt > /tmp/payload.json
gh api -X PATCH repos/<owner>/<repo>/issues/<issue-number> --input /tmp/payload.json
```

After applying (whichever path), **always re-fetch the body and confirm the rewrite is
actually visible** before reporting success. Do not trust the apply command's output
alone:

```bash
gh pr view <PR-number> --json body --jq .body | grep -F "<a phrase that should now be present>"
```

If the rewrite is not visible, treat the lint as failed and surface the real status to
the user — never report a successful rewrite that didn't land.

---

## Reporting

Report only the checks that changed something. Skip silent checks.

- Check A changed something:
  > PR 本文の auto-close キーワード（Closes/Fixes/Resolves 等）を `Refs` に置換しました。Issue は動作確認完了後に手動でクローズしてください。
- Check B moved items:
  > PR Test plan から merge 前に検証不可能な N 件を Issue #M の動作確認セクションへ移しました。
- Both:
  > PR 本文を修正しました: auto-close キーワード除去、merge 前検証不可な N 件を Issue #M へ移動。

If nothing changed, stay silent and let the original PR-created / PR-updated report
continue.

## Notes

- Commit messages can also trigger auto-close. The convention is "do not write `Closes`
  etc. in commit messages in the first place"; this skill's scope is the PR body only.
- `Refs #N` does not break GitHub's auto-link.
- Requires a token with PR-edit and issue-edit permission. If `gh pr edit` or
  `gh issue edit` returns an auth error, refresh via `/gh-token` and retry.
- Be conservative on Check B — when in doubt about whether an item is pre-merge
  verifiable, leave it in the PR and tell the user. Better to over-keep in the PR than
  to silently drop a real check.
