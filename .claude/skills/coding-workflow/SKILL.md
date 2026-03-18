---
name: coding-workflow
description: >
  Project coding workflow guide. Use when deciding what to do next, how to proceed,
  or when the user asks about the workflow. Also use when the user says
  "next step", "what now", "how do we proceed", or needs guidance on the development process.
---

# coding-workflow — Project Development Workflow

This project follows a structured workflow with SpecKit for planning and implementation.

## Development Flow

Feature development follows these steps in order. Each step uses the corresponding `/speckit.*` skill.

1. `/speckit.specify` → Create `specs/<NNN>-<name>/spec.md` (feature specification)
2. `/speckit.clarify` → Identify and resolve ambiguities in the spec
3. `/speckit.survey` → Validate direction from multiple angles (`survey.md`)
4. `/speckit.plan` → Implementation plan (`research.md`, `data-model.md`, `contracts/`, `quickstart.md`)
5. `/speckit.checklist` → Generate requirements checklist
6. `/speckit.tasks` → Generate dependency-ordered `tasks.md`
7. `/speckit.analyze` → Cross-check consistency across spec/plan/tasks
8. `/speckit.implement` → Implement based on `tasks.md`
9. `/simplify` → Review changed code for reuse, quality, and efficiency (required if code was changed)
10. `/check` → Run mechanical checks (formatter, linter, tests) after coding

### As Needed

- `/speckit.taskstoissues` → Create GitHub Issues from `tasks.md`
- `/adr` → Record architectural decisions in `.specify/decisions/`

## Directory Structure

| Path                  | Description                                                                                                          |
| --------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `specs/<NNN>-<name>/` | Planning artifacts (spec, survey, plan, tasks, etc.). Point-in-time snapshots, not the source of truth after implementation |
| `.specify/decisions/` | ADRs (Architecture Decision Records). Record of technical decisions                                                  |

## Important Notes

- **Do not skip steps**: Wait for user instruction before proceeding to the next step
- **specs/ is history**: Planning docs are snapshots. The source of truth lives elsewhere (e.g., README.md, the code itself)
