---
description: Create or update the project constitution -- the non-negotiable principles, governance rules, and project identity that all specs and plans must follow. Use this when setting up a new project, adding or revising principles, or when governance needs to change.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

You are updating the project constitution at `.specify/memory/constitution.md`. This file is a template containing placeholder tokens in square brackets (e.g. `[PROJECT_NAME]`, `[PRINCIPLE_1_NAME]`). Your job is to collect or derive concrete values, fill the template precisely, and propagate amendments across dependent artifacts.

If `.specify/memory/constitution.md` does not exist, copy it from `.specify/templates/constitution-template.md` first.

Follow this execution flow:

1. Load the existing constitution at `.specify/memory/constitution.md`.
   - Identify every placeholder token of the form `[ALL_CAPS_IDENTIFIER]`.
   - The user may need fewer or more principles than the template provides. Adjust the document to match the requested number while following the template's general structure.

2. Collect or derive values for each placeholder:
   - If user input supplies a value, use it.
   - Otherwise infer from repo context (README, docs, prior constitution versions).
   - For governance dates: `RATIFICATION_DATE` is the original adoption date (ask or mark TODO if unknown). `LAST_AMENDED_DATE` is today if changes are made, otherwise keep previous.
   - Increment `CONSTITUTION_VERSION` using semantic versioning:
     - MAJOR: backward-incompatible governance or principle removals/redefinitions.
     - MINOR: new principle or section added, materially expanded guidance.
     - PATCH: clarifications, wording, typo fixes, non-semantic refinements.
   - If the bump type is ambiguous, state your reasoning before finalizing.

3. Draft the updated constitution:
   - Replace every placeholder with concrete text. No bracketed tokens should remain unless intentionally deferred (justify any that are).
   - Preserve heading hierarchy. Remove template comments once replaced, unless they still add guidance.
   - Each Principle section needs: a succinct name, a paragraph or bullet list of non-negotiable rules, and explicit rationale if not obvious.
   - The Governance section must list amendment procedure, versioning policy, and compliance review expectations.

4. Propagate consistency to dependent files:
   - Read `.specify/templates/plan-template.md` and verify any "Constitution Check" or rules align with updated principles.
   - Read `.specify/templates/spec-template.md` for scope/requirements alignment. Update if the constitution adds or removes mandatory sections or constraints.
   - Read `.specify/templates/tasks-template.md` and ensure task categorization reflects new or removed principle-driven task types (e.g., observability, versioning, testing discipline).
   - Read each command file in `.specify/templates/commands/*.md` to verify no outdated references remain.
   - Read any runtime guidance docs (e.g., `README.md`, `docs/quickstart.md`). Update references to changed principles.

5. Produce a Sync Impact Report (prepend as an HTML comment at top of the constitution file):
   - Version change: old -> new
   - Modified principles (old title -> new title if renamed)
   - Added and removed sections
   - Templates requiring updates (done / pending) with file paths
   - Follow-up TODOs if any placeholders were intentionally deferred

6. Validate before writing:
   - No remaining unexplained bracket tokens.
   - Version line matches the report.
   - Dates in ISO format (YYYY-MM-DD).
   - Principles are declarative, testable, and free of vague language.

7. Write the completed constitution back to `.specify/memory/constitution.md` (overwrite).

8. Output a final summary with:
   - New version and bump rationale.
   - Files flagged for manual follow-up.
   - Suggested commit message (e.g., `docs: amend constitution to vX.Y.Z (principle additions + governance update)`).

If the user supplies partial updates (e.g., only one principle revision), still run validation and version decision steps. If critical info is missing (e.g., ratification date truly unknown), insert `TODO(<FIELD_NAME>): explanation` and include it in the Sync Impact Report under deferred items.

Use Markdown headings exactly as in the template (do not change heading levels). Keep a single blank line between sections. Do not create a new template -- always operate on the existing `.specify/memory/constitution.md` file.
