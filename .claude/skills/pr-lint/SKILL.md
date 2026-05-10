---
name: pr-lint
description: >
  Lint a PR for project-rule compliance immediately after it is created or updated, and
  auto-fix violations. Runs unconditionally right after `gh pr create` or
  `gh pr edit --body`, and always before reporting "PR created" / "PR updated" to the
  user. Current rules: (1) strip issue auto-close keywords (Closes / Fixes / Resolves
  etc.) from the PR body so the linked issue is not auto-closed on merge; (2) keep only
  items that gate the merge decision in the PR Test plan, moving everything else
  (post-merge observations, follow-up tests, regression watchlists, doc tasks) to the
  linked Issue — creating an Issue if none exists — so reviewers can decide approve /
  not-approve cleanly.
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

1. Fetch the PR body and the linked Issue number.
2. Run Check A (auto-close keywords). Stage a rewrite if needed.
3. Run Check B (non-merge-decision items). Stage a PR rewrite + an Issue-side append if
   needed. If no linked Issue exists, create one and use it.
4. If anything was staged, apply: `gh pr edit --body-file` and (if items moved)
   `gh issue edit --body-file` or `gh issue create`.
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

## Check B — Keep Only Merge-Decision Tests in the PR

### Why

A PR Test plan exists for one purpose only: **to let a reviewer decide approve / not
approve**. Anything that does not feed that decision does not belong in the PR.

Two categories must be moved out:

1. **Post-merge / post-release observations** (Dependabot scan after merge, next-release
   `release.yml` behavior, etc.) — the reviewer cannot perform them, and putting them in
   the PR creates a deadlock: approval requires verification, verification requires
   merge, merge requires approval.
2. **Follow-up tests, regression watchlists, ops checklists, future rebase reminders,
   doc updates** — these may be testable in principle, but they do not gate this PR's
   merge. They belong with the Issue, not the PR.

Both categories move to the linked Issue under a single `## 動作確認` / `## Verification`
section. This way the PR stays sharp ("can I approve?") and the Issue stays the durable
home for everything else ("are we done with this work?"). The Issue is closed only when
all of its verification items are checked.

If no linked Issue exists, **create one** for this work and move the items there. Never
discard a real test just because the PR cannot host it; the test belongs somewhere
durable.

### Detect

Two passes against the Test plan section (a heading containing "Test plan" or "テスト"
near a checklist of `- [ ]` items).

**Pass 1 — keyword auto-detect.** Items matching any pattern below are unambiguously
non-merge-decision and move out without further judgment:

| Pattern (case-insensitive) | Reason |
|----------------------------|--------|
| `マージ後` / `merge 後` / `after merge` / `post-merge` / `post merge` | Requires the merge to happen first |
| `develop merge` / `default branch` / `default ブランチ` | Requires landing on the default / target branch |
| `次回リリース` / `next release` / `on release` | Requires a future release event |
| `Dependabot` (in a check item) | Dependabot only reads config from default branch |
| Workflow trigger that does not fire on PR (`pull_request: closed`, `push: develop`, `release-*` push, etc.) used as a verification step | The trigger does not fire during PR review |
| `follow-up` / `フォローアップ` / `後追い` / `regression watch` / `回帰チェック` / `次回 PR で` | Stated as future work, not this PR's gate |

**Pass 2 — judgment.** For each remaining item, ask: *"can a reviewer perform this check
right now, and does its outcome change whether this PR should be approved?"*. If the
answer to either part is no, the item is not a merge-decision test. Examples:

- "Update the README in a follow-up PR" → cannot be verified here, doesn't gate this PR
- "Run nightly load test against staging next week" → not testable now, not a merge gate
- "After merge, monitor error rate for 24h" → post-merge observation

When Pass 2 is ambiguous (the line could plausibly gate the merge), **leave it in the PR
and surface it to the user** instead of guessing.

If neither pass produces matches, skip to Apply step.

### Reclassify

Split the Test plan into two parts:

- **Merge-decision tests** (stays in PR): items the reviewer can perform now and whose
  outcome would change the approve / not-approve decision.
- **Everything else** (moves out): all items matched by Pass 1, plus Pass-2 items
  judged non-gating with confidence.

### Find or create the linked Issue

1. Read the PR body for `#<number>` references and pick the Issue most plausibly tied to
   this PR's work (usually the one mentioned in `Refs #N` or in the Summary).
2. If multiple plausible candidates exist, ask the user which is correct.
3. If none exists, **create one** with `gh issue create`. Title: mirror the PR title (drop
   the trailing `#N` if present). Body: brief restatement of what this PR delivers, plus
   the moved verification items under `## 動作確認`. Then add a `Refs #<new-issue>` line
   to the PR Summary so the connection is visible.

### Rewrite

- PR body: keep only the merge-decision items in the Test plan. If the Test plan section
  becomes empty, replace it with a single line: `See #N for the verification checklist
  (this PR has no reviewer-side test required)`.
- Issue body: append (or update if present) a `## 動作確認` / `## Verification` section
  containing the moved items. If the section already exists, merge intelligently —
  deduplicate, keep existing checklist state. Note that the Issue closes only when all
  verification items are checked, not on PR merge.

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
  > PR Test plan から merge 判定外の N 件を Issue #M の動作確認セクションへ移しました。
- Check B created an Issue then moved items:
  > 紐づく Issue が無かったため Issue #M を作成し、merge 判定外の N 件をそこへ移しました。
- Both:
  > PR 本文を修正しました: auto-close キーワード除去、merge 判定外の N 件を Issue #M へ移動。

If nothing changed, stay silent and let the original PR-created / PR-updated report
continue.

## Notes

- Commit messages can also trigger auto-close. The convention is "do not write `Closes`
  etc. in commit messages in the first place"; this skill's scope is the PR body only.
- `Refs #N` does not break GitHub's auto-link.
- Requires a token with PR-edit and issue-edit permission. If `gh pr edit` or
  `gh issue edit` returns an auth error, refresh via `/gh-token` and retry.
- Be conservative on Check B — when in doubt about whether an item gates this PR's merge,
  leave it in the PR and tell the user. Better to over-keep in the PR than to silently
  drop a real check or move it to the wrong place.
- Never delete a verification item just because it does not belong in the PR. It belongs
  somewhere durable; either the existing linked Issue or a freshly created one.
