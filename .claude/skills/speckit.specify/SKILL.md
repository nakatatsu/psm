---
description: Create a feature specification from a GitHub Issue URL. Extracts the Issue number as the SpecKit feature number, retrieves Issue details via GitHub API, generates a numbered feature branch, writes a structured spec, validates quality, and handles clarifications. Use this when the user has a GitHub Issue and wants to create a spec. Requires a GitHub Issue URL as input.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

The text the user typed after the command is a GitHub Issue URL. A valid GitHub Issue URL is **required**. Do not ask the user to repeat it unless they provided an empty command.

Given the user input:

1. **Detect and validate the GitHub Issue URL**:

   a. Scan the user input for a URL matching `https://github.com/{owner}/{repo}/issues/{number}`.
      - If no Issue URL is found: **ERROR** "GitHub Issue URL が必要です。先に GitHub Issue を作成し、その URL を渡してください。例: `/speckit.specify https://github.com/owner/repo/issues/123`"
      - If the URL does not match the expected pattern (e.g., it's a PR URL like `/pull/10`): **ERROR** "不正な Issue URL です。`https://github.com/{owner}/{repo}/issues/{number}` 形式で指定してください"
   b. Extract the Issue number from the URL.
   c. Extract the `{owner}/{repo}` from the URL.
   d. Compare `{owner}/{repo}` against the current repository's remote (via `git remote get-url origin`). If they differ, display a warning but continue.

2. **Verify the Issue exists and retrieve its details**:

   Run `gh issue view {number} --repo {owner}/{repo} --json title,body,labels` to:
   - Confirm the Issue exists. If the command fails: **ERROR** "Issue #{number} が見つかりません。URL を確認してください"
   - Retrieve the Issue title and body for use in subsequent steps.

3. **Generate a concise short name** (2-4 words) for the branch:
   - Extract meaningful keywords from the **Issue title and body** (title is the primary source; body provides supplementary context when the title alone is insufficient).
   - Use action-noun format when possible (e.g., "add-user-auth", "fix-payment-timeout").
   - Preserve technical terms and acronyms (OAuth2, API, JWT, etc.).
   - Example: Issue title "Add user authentication" -> "user-auth"

4. **Check for branch collisions and create the feature branch**:

   a. Fetch all remote branches:
      ```bash
      git fetch --all --prune
      ```

   b. Check if a branch with the Issue number already exists:
      - Remote branches: `git ls-remote --heads origin | grep -E 'refs/heads/0*{number}-'`
      - Local branches: `git branch | grep -E '^[* ]*0*{number}-'`
      - Specs directories: Check for directories matching `specs/0*{number}-*`
      - If any match is found: **ERROR** "ブランチ番号 {number} は既に使用されています（既存: {matched branch/dir}）。対応する Issue を確認してください"

   c. Run `.specify/scripts/bash/create-new-feature.sh --json --number {issue-number} --short-name "your-short-name" "{Issue title}"`.
      - The `--number` is always the Issue number extracted in step 1.
      - For single quotes in args, use: `"I'm Groot"` (double-quote) or `'I'\''m Groot'` (escape).

   The JSON output contains BRANCH_NAME and SPEC_FILE paths -- always refer to it. Run this script only once per feature.

5. Load `.specify/templates/spec-template.md` to understand required sections.

6. **Build the specification** following this flow:
   1. Parse the Issue title and body retrieved in step 2. If both are empty: ERROR "Issue has no content".
   2. Extract key concepts: actors, actions, data, constraints.
   3. For unclear aspects, make informed guesses based on context and industry standards. Only use `[NEEDS CLARIFICATION: specific question]` markers (maximum 3 total) when the choice significantly impacts scope or UX, multiple reasonable interpretations exist, and no reasonable default works. Prioritize by: scope > security/privacy > UX > technical details.
   4. Fill User Scenarios and Testing section. If no clear user flow can be determined: ERROR "Cannot determine user scenarios".
   5. Generate Functional Requirements -- each must be testable. Use reasonable defaults for unspecified details and document assumptions in the Assumptions section.
   6. Define Success Criteria: measurable, technology-agnostic, verifiable outcomes. Include both quantitative metrics (time, performance) and qualitative measures (task completion, satisfaction).
   7. Identify Key Entities if data is involved.
   8. Return: SUCCESS (spec ready for planning).

7. Write the specification to SPEC_FILE using the template structure, replacing placeholders with concrete details while preserving section order and headings.

8. **Specification Quality Validation**: After writing the spec, validate it against quality criteria.

   a. **Create a checklist** at `FEATURE_DIR/checklists/requirements.md` covering:
      - Content Quality: no implementation details, focused on user value, written for non-technical stakeholders, all mandatory sections completed.
      - Requirement Completeness: no unresolved NEEDS CLARIFICATION markers, testable and unambiguous requirements, measurable and technology-agnostic success criteria, edge cases identified, scope bounded, dependencies listed.
      - Feature Readiness: functional requirements have acceptance criteria, user scenarios cover primary flows, no implementation details leak into the spec.

   b. **Run validation**: Review the spec against each checklist item and document specific issues found.

   c. **Handle results**:
      - All items pass: mark checklist complete and proceed.
      - Items fail (excluding NEEDS CLARIFICATION): list failing items, update the spec, and re-validate (max 3 iterations). If still failing, document remaining issues and warn the user.
      - NEEDS CLARIFICATION markers remain: extract them (keep only the 3 most critical if more exist, make informed guesses for the rest). Present each to the user:

        ```markdown
        ## Question [N]: [Topic]

        **Context**: [Quote relevant spec section]
        **What we need to know**: [Specific question]

        | Option | Answer | Implications |
        |--------|--------|--------------|
        | A      | [First answer] | [Impact on feature] |
        | B      | [Second answer] | [Impact on feature] |
        | C      | [Third answer] | [Impact on feature] |
        | Custom | Your own answer | [How to provide it] |
        ```

        Present all questions together, numbered Q1-Q3. Wait for the user to respond (e.g., "Q1: A, Q2: Custom - [details], Q3: B"). Update the spec with answers and re-validate.

   d. **Update the checklist** with current pass/fail status after each iteration.

9. Report completion with branch name, spec file path, checklist results, and readiness for `/speckit.clarify` or `/speckit.plan`.

The script creates and checks out the new branch and initializes the spec file before writing.

## Guidelines

- Focus on WHAT users need and WHY -- avoid HOW (no tech stack, APIs, code structure).
- Write for business stakeholders, not developers.
- Do not embed checklists inside the spec itself; those are separate artifacts.
- Mandatory sections must be completed for every feature. Optional sections should be removed entirely when not relevant (do not leave as "N/A").
- Make informed guesses using context and common patterns. Document assumptions in the Assumptions section.
- Think like a tester: every vague requirement should fail the "testable and unambiguous" check.

Reasonable defaults that do not need clarification: industry-standard data retention, standard web/mobile performance expectations, user-friendly error handling, standard authentication methods (session-based or OAuth2), and project-appropriate integration patterns.

### Success Criteria

Success criteria must be measurable, technology-agnostic, user-focused, and verifiable without knowing implementation details.

Good: "Users complete checkout in under 3 minutes", "System supports 10,000 concurrent users", "95% of searches return results in under 1 second".

Bad: "API response under 200ms" (use a user-facing metric instead), "Redis cache hit rate above 80%" (technology-specific).
