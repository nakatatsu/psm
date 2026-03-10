---
description: Execute a complete implementation from tasks.md -- builds the project phase by phase following the task plan. Use this after tasks have been generated and reviewed, when the user wants to start building.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

1. Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` from repo root and parse FEATURE_DIR and AVAILABLE_DOCS list. All paths must be absolute. For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **Check checklists status** (if FEATURE_DIR/checklists/ exists):
   - Scan all checklist files and count items matching `- [ ]`, `- [X]`, or `- [x]`
   - Display a status table:

     ```text
     | Checklist   | Total | Completed | Incomplete | Status |
     |-------------|-------|-----------|------------|--------|
     | ux.md       | 12    | 12        | 0          | PASS   |
     | test.md     | 8     | 5         | 3          | FAIL   |
     ```

   - If any checklist has incomplete items, STOP and ask: "Some checklists are incomplete. Do you want to proceed with implementation anyway? (yes/no)"
   - If all checklists pass, proceed automatically.

3. **Load implementation context** from FEATURE_DIR:
   - REQUIRED: tasks.md (complete task list and execution plan)
   - REQUIRED: plan.md (tech stack, architecture, file structure)
   - IF EXISTS: data-model.md, contracts/, research.md, quickstart.md

4. **Project Setup Verification** -- Detect project technologies from plan.md and the repo, then create or verify appropriate ignore files. The principle: each tool in the stack should have its ignore file with standard patterns for that technology.
   - Check `git rev-parse --git-dir 2>/dev/null` -- if git repo, verify .gitignore
   - Check for Dockerfile, .eslintrc*, .prettierrc*, package.json, *.tf, helm charts -- create corresponding ignore files (.dockerignore, .eslintignore, .prettierignore, etc.)
   - For example, a Node.js .gitignore should cover `node_modules/`, `dist/`, `.env*`; a Python .gitignore should cover `__pycache__/`, `.venv/`, `*.pyc`
   - If an ignore file already exists, append only missing critical patterns. If missing, create with the full standard pattern set for that technology.

5. **Parse tasks.md** and extract task phases, dependency ordering, file paths, and parallel markers [P].

6. **Execute implementation phase by phase**:
   - Complete each phase before moving to the next
   - Respect dependencies: run sequential tasks in order, parallel tasks [P] can run together
   - Follow TDD approach: execute test tasks before their corresponding implementation tasks
   - Tasks affecting the same files run sequentially
   - Validate each phase completion before proceeding

7. **Implementation order**:
   - Setup first: project structure, dependencies, configuration
   - Tests before code: contracts, entities, integration scenarios
   - Core development: models, services, CLI commands, endpoints
   - Integration: database connections, middleware, logging, external services
   - Polish: unit tests, performance optimization, documentation

8. **Progress tracking**:
   - Report progress after each completed task
   - Halt execution if any non-parallel task fails
   - For parallel tasks [P], continue with successful tasks and report failed ones
   - Provide clear error messages with context for debugging
   - Mark completed tasks as [X] in the tasks file

9. **Completion validation**:
   - Verify all required tasks are completed
   - Check that implemented features match the original specification
   - Validate that tests pass and coverage meets requirements
   - Report final status with summary of completed work

Note: This command assumes a complete task breakdown exists in tasks.md. If tasks are incomplete or missing, suggest running `/speckit.tasks` first to generate the task list.
