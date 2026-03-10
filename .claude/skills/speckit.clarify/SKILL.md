---
description: Identify underspecified areas in the current feature spec by asking up to 5 highly targeted clarification questions, then encoding answers directly into the spec. Use this after writing or updating a spec to reduce ambiguity and fill gaps before planning begins.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Goal

Detect and reduce ambiguity or missing decision points in the active feature specification, then record clarifications directly in the spec file.

This workflow should run before invoking `/speckit.plan`. If the user explicitly states they are skipping clarification (e.g., exploratory spike), proceed but warn that downstream rework risk increases.

## Execution Steps

### 1. Initialize

Run `.specify/scripts/bash/check-prerequisites.sh --json --paths-only` from repo root once. Parse the JSON payload for `FEATURE_DIR` and `FEATURE_SPEC`. If JSON parsing fails, abort and instruct user to re-run `/speckit.specify`.
For single quotes in args like "I'm Groot", use escape syntax: e.g `'I'\''m Groot'` (or double-quote if possible: `"I'm Groot"`).

### 2. Scan for Ambiguity and Gaps

Load the current spec file. For each category below, mark status as Clear / Partial / Missing. Produce an internal coverage map for prioritization (do not output the raw map unless no questions will be asked).

**Functional Scope**: Core user goals and success criteria, explicit out-of-scope declarations, user roles/personas.

**Domain and Data Model**: Entities, attributes, relationships, identity/uniqueness rules, lifecycle/state transitions, data volume/scale assumptions.

**Interaction and UX Flow**: Critical user journeys, error/empty/loading states, accessibility or localization notes.

**Non-Functional Quality Attributes**: Performance targets, scalability limits, reliability/availability expectations, observability signals, security/privacy posture, compliance constraints.

**Integration and External Dependencies**: External services/APIs and failure modes, data import/export formats, protocol/versioning assumptions.

**Edge Cases and Failure Handling**: Negative scenarios, rate limiting/throttling, conflict resolution (e.g., concurrent edits).

**Constraints and Tradeoffs**: Technical constraints (language, storage, hosting), explicit tradeoffs or rejected alternatives.

**Terminology and Consistency**: Canonical glossary terms, avoided synonyms or deprecated terms.

**Completion Signals**: Acceptance criteria testability, measurable definition-of-done indicators.

**Placeholders**: TODO markers, unresolved decisions, ambiguous adjectives ("robust", "intuitive") lacking quantification.

For each Partial or Missing category, add a candidate question unless clarification would not materially change implementation or the information is better deferred to planning.

### 3. Build Question Queue

Generate up to 5 prioritized candidate questions internally. Do not output them all at once. Each question must be answerable with either:
- A short multiple-choice selection (2-5 options), OR
- A short-phrase answer (constrain: "Answer in <=5 words")

Only include questions whose answers materially impact architecture, data modeling, task decomposition, test design, UX behavior, or compliance. Cover the highest-impact unresolved categories first. If more than 5 categories remain unresolved, select top 5 by (Impact x Uncertainty).

### 4. Sequential Questioning Loop

Present ONE question at a time. Maximum 5 questions total.

**For multiple-choice questions:**
- Analyze all options and determine the most suitable based on best practices, common patterns, risk reduction, and alignment with project goals.
- Present your recommendation prominently: `**Recommended:** Option [X] - <reasoning>`
- Render all options as a Markdown table (Option | Description), including a "Short" row for free-form alternatives when appropriate.
- Add: `Reply with the option letter, accept the recommendation by saying "yes", or provide your own short answer.`

**For short-answer questions:**
- Provide a suggested answer: `**Suggested:** <answer> - <brief reasoning>`
- Add: `Format: Short answer (<=5 words). Say "yes" to accept, or provide your own.`

**After each answer:**
- If user replies "yes", "recommended", or "suggested", use the stated recommendation.
- If ambiguous, ask for quick disambiguation (does not count as a new question).
- Record the answer in working memory and proceed to the next question.

**Stop when:** all critical ambiguities are resolved early, user signals completion ("done", "good", "no more"), or 5 questions have been asked. Never reveal future queued questions.

### 5. Integrate Each Answer

Maintain an in-memory representation of the spec loaded at start.

For the first answer in this session:
- Ensure a `## Clarifications` section exists (create it after the highest-level overview section if missing).
- Under it, create a `### Session YYYY-MM-DD` subheading for today if not present.

After each accepted answer:
- Append a bullet: `- Q: <question> -> A: <final answer>`
- Apply the clarification to the most appropriate section:
  - Functional ambiguity -> update Functional Requirements
  - Actor/role distinction -> update User Stories or Actors subsection
  - Data shape/entities -> update Data Model with fields, types, relationships
  - Non-functional constraint -> add/modify measurable criteria in Quality Attributes (convert vague adjectives to metrics)
  - Edge case/negative flow -> add bullet under Edge Cases / Error Handling
  - Terminology conflict -> normalize term across spec; retain original only with `(formerly referred to as "X")` once
- If the clarification invalidates an earlier ambiguous statement, replace it rather than duplicating -- leave no obsolete contradictory text.
- Save the spec file after each integration to minimize context loss risk.
- Preserve formatting: do not reorder unrelated sections or break heading hierarchy.
- Keep each inserted clarification minimal and testable.

### 6. Validate

After each write and as a final pass, verify: one bullet per accepted answer in the Clarifications section, no duplicates, total questions <=5, updated sections contain no lingering vague placeholders the answer was meant to resolve, no contradictory earlier statements remain, markdown structure is valid, and terminology is consistent across all updated sections.

### 7. Save

Write the updated spec back to `FEATURE_SPEC`.

### 8. Report

Output:
- Number of questions asked and answered
- Path to updated spec
- Sections touched
- Coverage summary table: each taxonomy category with status (Resolved, Deferred, Clear, or Outstanding)
- If Outstanding or Deferred items remain, recommend whether to proceed to `/speckit.plan` or run `/speckit.clarify` again later
- Suggested next command

## Behavior Rules

- If no meaningful ambiguities found, respond: "No critical ambiguities detected worth formal clarification." and suggest proceeding.
- If spec file is missing, instruct user to run `/speckit.specify` first.
- Never exceed 5 total questions (disambiguation retries for the same question do not count as new questions).
- Avoid speculative tech stack questions unless the absence blocks functional clarity.
- Respect user early termination signals ("stop", "done", "proceed").
- If quota is reached with unresolved high-impact categories, flag them under Deferred with rationale.

Context for prioritization: $ARGUMENTS
