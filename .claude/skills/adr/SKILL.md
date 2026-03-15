---
name: adr
description: >
  Create and manage Architecture Decision Records (ADRs). Records technical decisions
  in a standard format under .specify/decisions/. Use this for design choices, technology
  selections, policy changes, rejected proposals — anything where "why did we decide this?"
  matters for the future. Also use when the user says "write an ADR", "record this decision",
  "document why we chose this", or after a speckit.survey leads to a conclusion.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Goal

Record technical decisions in a standard format so that future team members — or your future self — can understand why a particular choice was made. An ADR captures the context, the decision, and its consequences in a single, lightweight, searchable document.

## When to Write an ADR

- Multiple technical options exist and you want to record why one was chosen
- A policy or principle was changed (e.g., constitution amendment)
- A proposal was evaluated and rejected (the rejection rationale is valuable)
- A decision will likely prompt "why is it done this way?" in the future

When in doubt, write one. You'll never regret having an ADR; you'll often regret not having one.

## ADR Template

```markdown
# ADR-[NNN]: [Title]

**Date**: [YYYY-MM-DD]
**Status**: [Proposed | Accepted | Rejected | Deprecated | Superseded by ADR-XXX]

## Context

What is the issue? Why is a decision needed now? What forces are at play
(technical, business, organizational)? Include enough background that a reader
unfamiliar with the project can understand the situation.

## Decision

State the decision clearly and concisely. One or two sentences.

## Consequences

What changes as a result of this decision? Include:

- What becomes easier or possible
- What becomes harder or impossible
- What follow-up actions are needed
- What to revisit if circumstances change
```

## Execution Steps

### 1. Determine ADR Number

Scan `/.specify/decisions/` for existing ADR files. Extract the highest number and increment by 1. If the decision relates to a feature branch (e.g., `004-aws-test-stub`), use the feature number as the ADR number for traceability. If the decision is not tied to a feature, use the next available number.

### 2. Gather Context

From the conversation, spec, survey, or user input, extract:

- What problem or question prompted this decision
- What options were considered
- What was chosen and why
- What the impact is

If a survey.md exists for this feature, reference its findings as evidence.

### 3. Write the ADR

Use the template above. Follow these principles:

- **Be concise but complete.** A good ADR is one page. If it's getting long, the decision might need splitting into multiple ADRs.
- **Write for future readers.** Someone joining the project in 6 months should understand the context without needing to ask anyone.
- **Record what was rejected and why.** The rejected alternatives are often more valuable than the chosen one — they prevent the team from re-debating the same options.
- **Include a clear Status.** This makes it easy to scan which decisions are current vs. superseded.
- **Link to evidence.** Reference survey.md, research.md, spec.md, or other artifacts rather than repeating their content.

### 4. Save

Write to `.specify/decisions/[NNN]-[short-name].md`.

Naming convention:

- `004-aws-test-stub-rejected.md` — feature-linked, rejected
- `005-ci-tool-versions.md` — feature-linked, accepted
- `010-go-error-handling.md` — standalone decision

### 5. Report

Output:

- Path to the ADR file
- ADR number, title, and status
- Summary of the decision in one sentence

## Superseding an ADR

When a previous decision is reversed or replaced:

1. Create a new ADR with the new decision
2. Update the old ADR's status to `Superseded by ADR-XXX`
3. In the new ADR's Context, reference the old ADR and explain what changed

This preserves the full history of reasoning without cluttering the current state.
