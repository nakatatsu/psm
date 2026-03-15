---
name: speckit.survey
description: >
  spec の前提・方向性を広い視野で検証する事前調査スキル。
  /speckit.survey で呼び出す。
  speckit.clarify の後、speckit.plan の前に実行する。
  技術選定だけでなく、そもそもの方針の妥当性、コスト、リスク、代替アプローチ、
  既存コードや過去の決定との整合性など、多角的な観点から調査する。
  ユーザーが「調べて」「調査して」「検討して」と言った場合もこのスキルを使う。
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Goal

spec.md の前提・方向性・技術選択を広い視野で検証し、survey.md を生成する。plan フェーズが正しい方向に進むための土台を作る。plan の Phase 0 (research.md) が「具体的な未知事項の解決」であるのに対し、survey は「そもそもこの方向で良いのか」を問う上位の調査。

## survey と plan Phase 0 (research) の違い

| | survey | plan Phase 0 (research) |
|---|---|---|
| **問い** | そもそもこの方向で良いか？ | この方向で具体的に何を使うか？ |
| **視野** | 広い（代替アプローチ、コスト、リスク、過去の決定） | 狭い（NEEDS CLARIFICATION の解決） |
| **出力** | survey.md | research.md |
| **タイミング** | clarify の後、plan の前 | plan の Phase 0 |

## Execution Steps

### 1. Initialize

Run `.specify/scripts/bash/check-prerequisites.sh --json --paths-only` from repo root once. Parse the JSON payload for `FEATURE_DIR` and `FEATURE_SPEC`. If JSON parsing fails, abort and instruct user to re-run `/speckit.specify`.

### 2. Load Context

Read the following files:
- `FEATURE_SPEC` (spec.md)
- `.specify/memory/constitution.md`
- Clarifications section in spec.md (if exists)
- Assumptions section in spec.md

Also scan:
- Existing codebase (dependencies, architecture, patterns)
- Past feature research (e.g., `specs/001-psm/research.md`) for relevant prior decisions
- Constitution for potential conflicts

### 3. Survey Taxonomy

Systematically evaluate spec.md against the following categories. The most important question is always: **"Is the spec solving the right problem?"**

#### A. Problem Reframing (最重要)

Before evaluating the solution, question the problem itself:

1. **Problem Definition** — Is the spec solving the right problem? Or is it solving a symptom? Step back and ask what the *actual* goal is. Example: "IP が変わるドメインを IP ベースで許可する方法" が問題ではなく、"外向き通信のドメイン制御" が本当の問題 → フォワードプロキシで解決。
2. **Hidden Assumptions** — What assumptions are baked into the spec that nobody questioned? What if those assumptions are wrong? List them explicitly and challenge each one.
3. **Lateral Thinking** — If you had to solve this problem without the approach described in the spec, how would you do it? Think of at least 2 completely different strategies that solve the same underlying goal.

#### B. Solution Evaluation

Once the problem is validated, evaluate the proposed solution:

4. **Necessity** — Is this feature actually needed? Is there a simpler way to achieve the same goal without building this? What is the cost of NOT doing it?
5. **Approach Alternatives** — Are there fundamentally different approaches? Not just "which library" but "which strategy entirely"?
6. **Prior Decisions** — Do past decisions (in previous features, constitution, existing code) still hold? Have circumstances changed?
7. **Cost & Complexity** — What is the true cost? Not just implementation time, but ongoing maintenance, dependency burden, cognitive overhead, CI complexity.

#### C. Risk & Feasibility

8. **Risk & Failure Modes** — What can go wrong? What are the failure modes we haven't considered? What happens if a key assumption turns out to be wrong?
9. **External Dependencies** — What are we depending on externally? How stable/reliable are those dependencies? What happens if they break or change?
10. **Feasibility Verification** — Are there claims in the spec/assumptions that need hands-on verification (PoC/spike)?

#### D. Integration & Governance

11. **Integration Impact** — How does this interact with the rest of the system? What existing code/tests/workflows need to change? What are the ripple effects?
12. **Constitution Compliance** — Does this require amending the constitution? If so, is the amendment justified?
13. **Scope Boundaries** — Is the scope well-bounded? Are we trying to do too much? What should explicitly be out of scope?

### 4. Execute Survey

For each category where issues are found:

1. **Search and verify** — Use web search, documentation, codebase exploration, and past research to gather facts. Do NOT speculate or assume.
2. **Evaluate alternatives** — For significant decisions, identify at least 2 alternatives with trade-offs.
3. **Check against existing code** — Verify compatibility with what already exists.
4. **Flag items needing PoC** — If something cannot be verified without hands-on testing, say so explicitly.

### 5. Write survey.md

Create `FEATURE_DIR/survey.md` with the following structure:

```markdown
# Survey: [Feature Name]

**Date**: [DATE]
**Spec**: [link to spec.md]

## Summary

[1-2 paragraph overview of findings. Key risks, validated assumptions, open questions.]

## S[N]: [Topic]

**Category**: [from taxonomy above]
**Finding**: [What was discovered]
**Recommendation**: [What to do about it]
**Evidence**: [How this was verified — sources, code references, test results]
```

Include a final section:

```markdown
## Items Requiring PoC

[List any claims that could not be verified without hands-on experimentation]

## Constitution Impact

[Any amendments needed, or confirmation that no amendments are required]

## Recommendation

[Proceed to /speckit.plan, or address specific issues first]
```

### 6. Report

Output:
- Path to survey.md
- Number of survey items
- Critical findings (if any)
- Items requiring PoC
- Constitution impact
- Recommended next step

## Behavior Rules

- **Start with the problem, not the solution.** Always execute Section A (Problem Reframing) first. If the problem is wrong, the rest doesn't matter.
- **Apply lateral thinking.** For every spec assumption, ask: "What if we did the opposite?" or "What if this constraint didn't exist?" The Squid proxy example: everyone was trying to solve "how to track rotating IPs" when the real answer was "stop using IP-based filtering entirely."
- Never skip investigation by assuming. Verify claims with actual evidence.
- Think broadly first, then dig deep on what matters. Don't start with a conclusion.
- Past decisions are context, not gospel. Re-evaluate if circumstances have changed.
- If the survey reveals that the spec's fundamental approach is wrong, say so clearly. This is the most valuable output a survey can produce.
- If web search is needed, use it. Do not guess at compatibility, capabilities, or availability.
- If a finding cannot be verified without hands-on experimentation, flag it as needing a PoC.
