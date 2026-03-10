---
description: Generate a custom checklist for the current feature based on user requirements. Checklists validate whether requirements are complete, clear, and consistent -- like unit tests for English prose. Use this whenever you need a quality gate for a spec before planning or implementation.
---

## Core Concept: Unit Tests for Requirements

Think of checklists as **unit tests for English**. If your spec is code written in prose, the checklist is its test suite. Each item tests whether a requirement is well-written, complete, unambiguous, and ready for implementation -- NOT whether the implementation works.

**FOR requirements quality validation:**

- "Are visual hierarchy requirements defined for all card types?" [Completeness]
- "Is 'prominent display' quantified with specific sizing/positioning?" [Clarity]
- "Are hover state requirements consistent across all interactive elements?" [Consistency]
- "Does the spec define what happens when logo image fails to load?" [Edge Cases]

**NOT for verification/testing:**

- WRONG: "Verify the button clicks correctly"
- WRONG: "Test error handling works"
- WRONG: "Confirm the API returns 200"

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Execution Steps

1. **Setup**: Run `.specify/scripts/bash/check-prerequisites.sh --json` from repo root and parse JSON for FEATURE_DIR and AVAILABLE_DOCS list.
   - All file paths must be absolute.
   - For single quotes in args like "I'm Groot", use escape syntax: e.g `'I'\''m Groot'` (or double-quote if possible: `"I'm Groot"`).

2. **Clarify intent**: Derive up to THREE contextual clarifying questions from the user's phrasing and signals in spec/plan/tasks. Only ask about information that materially changes checklist content. Skip any question already answered by `$ARGUMENTS`.

   Generation approach:
   1. Extract signals: domain keywords (auth, latency, UX, API), risk indicators ("critical", "must", "compliance"), stakeholder hints ("QA", "review"), explicit deliverables ("a11y", "rollback").
   2. Cluster into candidate focus areas (max 4) ranked by relevance.
   3. Identify probable audience and timing if not explicit.
   4. Detect missing dimensions: scope breadth, depth/rigor, risk emphasis, exclusion boundaries, measurable acceptance criteria.
   5. Formulate questions from these archetypes: scope refinement, risk prioritization, depth calibration, audience framing, boundary exclusion, scenario class gap.

   Present options as a compact table (Option | Candidate | Why It Matters) when useful; omit the table if free-form is clearer. Limit to A-E options. Never ask the user to restate what they already said.

   Defaults when interaction is impossible: Standard depth, Reviewer audience (Author if non-code), top 2 relevance clusters.

   Output questions as Q1/Q2/Q3. After answers, if >=2 scenario classes remain unclear, ask up to TWO more follow-ups (Q4/Q5) with one-line justification each. Do not exceed five total questions.

3. **Understand user request**: Combine `$ARGUMENTS` + clarifying answers to derive checklist theme, consolidate must-have items, map focus selections to categories, and infer missing context from spec/plan/tasks without hallucinating.

4. **Load feature context**: Read from FEATURE_DIR:
   - spec.md: Feature requirements and scope
   - plan.md (if exists): Technical details, dependencies
   - tasks.md (if exists): Implementation tasks

   Load only portions relevant to active focus areas. Prefer summarizing long sections into concise bullets. Use progressive disclosure -- add follow-on retrieval only if gaps are detected.

5. **Generate checklist**: Create requirement quality validation items.

   - Create `FEATURE_DIR/checklists/` directory if it doesn't exist.
   - Use a short, descriptive filename based on domain: `[domain].md` (e.g., `ux.md`, `api.md`, `security.md`).
   - If file does NOT exist: create new file, number items from CHK001.
   - If file exists: append new items, continuing from the last CHK ID.
   - Never delete or replace existing checklist content.

   **Writing Checklist Items**

   Every item evaluates the requirements themselves for: Completeness, Clarity, Consistency, Measurability, and Coverage.

   WRONG (tests implementation):
   - "Verify landing page displays 3 episode cards"
   - "Test hover states work on desktop"

   CORRECT (tests requirement quality):
   - "Are the exact number and layout of featured episodes specified?" [Completeness]
   - "Is 'prominent display' quantified with specific sizing/positioning?" [Clarity]
   - "Are hover state requirements consistent across all interactive elements?" [Consistency]
   - "Are keyboard navigation requirements defined for all interactive UI?" [Coverage]
   - "Is the fallback behavior specified when logo image fails to load?" [Edge Cases]

   Prohibited patterns -- these turn the checklist into an implementation test:
   - Items starting with "Verify", "Test", "Confirm", "Check" + implementation behavior
   - References to code execution, user actions, or system behavior
   - Phrases like "displays correctly", "works properly", "functions as expected"

   Required patterns -- these test requirement quality:
   - "Are [requirement type] defined/specified/documented for [scenario]?"
   - "Is [vague term] quantified/clarified with specific criteria?"
   - "Are requirements consistent between [section A] and [section B]?"
   - "Does the spec define [missing aspect]?"

   **Item Structure**: Question format about requirement quality, focused on what's written (or missing) in the spec/plan. Include quality dimension in brackets [Completeness/Clarity/Consistency/etc.]. Reference spec section `[Spec section-X.Y]` when checking existing requirements. Use `[Gap]` marker for missing requirements.

   **Category Structure** -- group items by requirement quality dimension:
   - Requirement Completeness
   - Requirement Clarity
   - Requirement Consistency
   - Acceptance Criteria Quality
   - Scenario Coverage (primary, alternate, exception, recovery, non-functional)
   - Edge Case Coverage
   - Non-Functional Requirements
   - Dependencies and Assumptions
   - Ambiguities and Conflicts

   **Scenario Coverage**: Check if requirements exist for primary, alternate, exception/error, recovery, and non-functional scenarios. For missing scenario classes, flag as: "Are [scenario type] requirements intentionally excluded or missing? [Gap]". Include resilience/rollback items when state mutation occurs.

   **Traceability**: At least 80% of items should include a traceability reference -- either a spec section `[Spec section-X.Y]` or a marker: `[Gap]`, `[Ambiguity]`, `[Conflict]`, `[Assumption]`.

   **Consolidation**: Soft cap at 40 items. Merge near-duplicates. If >5 low-impact edge cases remain, combine into one item: "Are edge cases X, Y, Z addressed in requirements? [Coverage]".

6. **Structure Reference**: Follow the canonical template in `.specify/templates/checklist-template.md` for title, meta section, category headings, and ID formatting. If template is unavailable, use: H1 title, purpose/created meta lines, `##` category sections containing `- [ ] CHK### <requirement item>` lines with globally incrementing IDs starting at CHK001.

7. **Report**: Output the full path to the checklist file, item count, and whether this created a new file or appended to an existing one. Summarize focus areas selected, depth level, actor/timing, and any user-specified must-have items incorporated.

## Example Checklist: UX Requirements Quality (`ux.md`)

Sample items showing the correct pattern:

- "Are visual hierarchy requirements defined with measurable criteria? [Clarity, Spec section-FR-1]"
- "Is the number and positioning of UI elements explicitly specified? [Completeness, Spec section-FR-1]"
- "Are interaction state requirements (hover, focus, active) consistently defined? [Consistency]"
- "Are accessibility requirements specified for all interactive elements? [Coverage, Gap]"
- "Is fallback behavior defined when images fail to load? [Edge Case, Gap]"
- "Can 'prominent display' be objectively measured? [Measurability, Spec section-FR-4]"
