---
description: Perform cross-artifact consistency and quality analysis across spec.md, plan.md, and tasks.md. Use this after task generation to catch inconsistencies, coverage gaps, ambiguities, and constitution violations before implementation begins.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Goal

Identify inconsistencies, duplications, ambiguities, and underspecified items across the three core artifacts (`spec.md`, `plan.md`, `tasks.md`) before implementation. This command should run only after `/speckit.tasks` has produced a complete `tasks.md`.

## Operating Constraints

This analysis is read-only -- modifying files during analysis would defeat the purpose of an independent quality check. Output a structured analysis report. Offer an optional remediation plan (user must explicitly approve before any edits are applied).

The project constitution (`.specify/memory/constitution.md`) is the authority within this analysis scope. Constitution conflicts are automatically CRITICAL severity and require adjustment of the spec, plan, or tasks -- not dilution or reinterpretation. If a principle itself needs to change, that must occur in a separate constitution update outside `/speckit.analyze`.

Keep analysis token-efficient: focus on actionable findings over exhaustive documentation. Load artifacts incrementally rather than dumping all content. Limit findings to 50 rows and summarize overflow. Rerunning without changes should produce consistent IDs and counts.

If a section is absent, report it accurately -- do not hallucinate missing content. Report zero issues gracefully with coverage statistics.

## Execution Steps

### 1. Initialize Analysis Context

Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` once from repo root and parse JSON for FEATURE_DIR and AVAILABLE_DOCS. Derive absolute paths:

- SPEC = FEATURE_DIR/spec.md
- PLAN = FEATURE_DIR/plan.md
- TASKS = FEATURE_DIR/tasks.md

Abort with an error message if any required file is missing (instruct the user to run the missing prerequisite command).
For single quotes in args like "I'm Groot", use escape syntax: e.g `'I'\''m Groot'` (or double-quote if possible: `"I'm Groot"`).

### 2. Load Artifacts

Load only the portions relevant to analysis from each artifact:

**From spec.md:** Overview/Context, Functional Requirements, Non-Functional Requirements, User Stories, Edge Cases (if present).

**From plan.md:** Architecture/stack choices, Data Model references, Phases, Technical constraints.

**From tasks.md:** Task IDs, Descriptions, Phase grouping, Parallel markers [P], Referenced file paths.

**From constitution:** Load `.specify/memory/constitution.md` for principle validation.

### 3. Build Semantic Models

Create internal representations (do not include raw artifacts in output):

- **Requirements inventory**: Each functional + non-functional requirement with a stable slug key (e.g., "User can upload file" -> `user-can-upload-file`)
- **User story/action inventory**: Discrete user actions with acceptance criteria
- **Task coverage mapping**: Map each task to requirements or stories by keyword/explicit reference
- **Constitution rule set**: Extract principle names and MUST/SHOULD normative statements

### 4. Detection Passes

Run these passes, limiting to 50 findings total with an overflow summary.

**A. Duplication** -- Identify near-duplicate requirements; mark lower-quality phrasing for consolidation.

**B. Ambiguity** -- Flag vague adjectives (fast, scalable, secure, intuitive, robust) lacking measurable criteria. Flag unresolved placeholders (TODO, TKTK, ???, `<placeholder>`, etc.).

**C. Underspecification** -- Requirements with verbs but missing object or measurable outcome. User stories missing acceptance criteria alignment. Tasks referencing files or components not defined in spec/plan.

**D. Constitution Alignment** -- Any requirement or plan element conflicting with a MUST principle. Missing mandated sections or quality gates from constitution.

**E. Coverage Gaps** -- Requirements with zero associated tasks. Tasks with no mapped requirement/story. Non-functional requirements not reflected in tasks (e.g., performance, security).

**F. Inconsistency** -- Terminology drift (same concept named differently across files). Data entities referenced in plan but absent in spec (or vice versa). Task ordering contradictions without dependency notes. Conflicting requirements (e.g., one requires Next.js while another specifies Vue).

### 5. Severity Assignment

- **CRITICAL**: Violates constitution MUST, missing core spec artifact, or requirement with zero coverage that blocks baseline functionality
- **HIGH**: Duplicate or conflicting requirement, ambiguous security/performance attribute, untestable acceptance criterion
- **MEDIUM**: Terminology drift, missing non-functional task coverage, underspecified edge case
- **LOW**: Style/wording improvements, minor redundancy not affecting execution order

### 6. Produce Compact Analysis Report

Output a Markdown report (no file writes) with this structure:

## Specification Analysis Report

| ID | Category | Severity | Location(s) | Summary | Recommendation |
|----|----------|----------|-------------|---------|----------------|
| A1 | Duplication | HIGH | spec.md:L120-134 | Two similar requirements ... | Merge phrasing; keep clearer version |

(One row per finding; stable IDs prefixed by category initial.)

**Coverage Summary Table:**

| Requirement Key | Has Task? | Task IDs | Notes |
|-----------------|-----------|----------|-------|

**Constitution Alignment Issues:** (if any)

**Unmapped Tasks:** (if any)

**Metrics:**

- Total Requirements
- Total Tasks
- Coverage % (requirements with >=1 task)
- Ambiguity Count
- Duplication Count
- Critical Issues Count

### 7. Provide Next Actions

At the end of the report, output a concise Next Actions block:

- If CRITICAL issues exist: recommend resolving before `/speckit.implement`
- If only LOW/MEDIUM: user may proceed, with improvement suggestions
- Provide explicit command suggestions, e.g., "Run /speckit.specify with refinement", "Run /speckit.plan to adjust architecture", "Manually edit tasks.md to add coverage for 'performance-metrics'"

### 8. Offer Remediation

Ask the user: "Would you like me to suggest concrete remediation edits for the top N issues?" Do not apply them automatically.

## Context

$ARGUMENTS
