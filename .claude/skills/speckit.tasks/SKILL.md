---
description: Generate a dependency-ordered tasks.md organized by user stories from spec.md. Produces a phased, immediately executable task plan from available design artifacts. Use this when design documents are ready and the user needs a structured implementation plan.
---

## User Input

```text
$ARGUMENTS
```

Consider the user input before proceeding (if not empty).

## Outline

1. **Setup**: Run `.specify/scripts/bash/check-prerequisites.sh --json` from repo root and parse FEATURE_DIR and AVAILABLE_DOCS list. All paths must be absolute. For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **Load design documents** from FEATURE_DIR:
   - Required: plan.md (tech stack, libraries, structure), spec.md (user stories with priorities)
   - Optional: data-model.md (entities), contracts/ (interface contracts), research.md (decisions), quickstart.md (test scenarios)
   - Not all projects have all documents. Generate tasks based on what is available.

3. **Execute task generation workflow**:
   - Extract tech stack, libraries, and project structure from plan.md
   - Extract user stories with their priorities (P1, P2, P3, etc.) from spec.md
   - If data-model.md exists: extract entities and map to user stories
   - If contracts/ exists: map interface contracts to user stories
   - If research.md exists: extract decisions for setup tasks
   - Generate tasks organized by user story (see Task Generation Rules)
   - Generate dependency graph showing user story completion order
   - Create parallel execution examples per user story
   - Validate task completeness: each user story should have all needed tasks and be independently testable

4. **Generate tasks.md** using `.specify/templates/tasks-template.md` as the structure. Fill with:
   - Feature name from plan.md
   - Phase 1: Setup tasks (project initialization)
   - Phase 2: Foundational tasks (blocking prerequisites for all user stories)
   - Phase 3+: One phase per user story in priority order from spec.md
   - Each phase includes: story goal, independent test criteria, tests (if requested), implementation tasks
   - Final Phase: Polish and cross-cutting concerns
   - All tasks follow the strict checklist format below
   - Clear file paths for each task
   - Dependencies section, parallel execution examples, and implementation strategy (MVP first, incremental delivery)

5. **Report**: Output the path to the generated tasks.md and a summary including:
   - Total task count and count per user story
   - Parallel opportunities identified
   - Independent test criteria for each story
   - Suggested MVP scope (typically User Story 1)
   - Format validation confirming all tasks follow the checklist format

Context for task generation: $ARGUMENTS

The tasks.md should be immediately executable -- each task must be specific enough that an LLM can complete it without additional context.

## Task Generation Rules

Tasks must be organized by user story so each story can be implemented and tested independently.

Tests are optional: only generate test tasks if explicitly requested in the feature specification or if the user requests a TDD approach.

### Checklist Format

Every task must follow this format exactly:

```text
- [ ] [TaskID] [P?] [Story?] Description with file path
```

**Format Components**:

1. **Checkbox**: Always start with `- [ ]` (markdown checkbox)
2. **Task ID**: Sequential number in execution order (T001, T002, T003...)
3. **[P] marker**: Include only if the task is parallelizable (different files, no dependencies on incomplete tasks)
4. **[Story] label**: Required for user story phase tasks only
   - Format: [US1], [US2], [US3], etc. (maps to user stories from spec.md)
   - Setup, Foundational, and Polish phases: no story label
   - User Story phases: must have story label
5. **Description**: Clear action with exact file path

**Correct examples**:

- `- [ ] T001 Create project structure per implementation plan`
- `- [ ] T005 [P] Implement authentication middleware in src/middleware/auth.py`
- `- [ ] T012 [P] [US1] Create User model in src/models/user.py`
- `- [ ] T014 [US1] Implement UserService in src/services/user_service.py`

**Wrong examples** (each violates a required component):

- `- [ ] Create User model` -- missing ID and Story label
- `T001 [US1] Create model` -- missing checkbox
- `- [ ] [US1] Create User model` -- missing Task ID
- `- [ ] T001 [US1] Create model` -- missing file path

### Task Organization

1. **From User Stories (spec.md)** -- this is the primary organizing principle:
   - Each user story (P1, P2, P3...) gets its own phase
   - Map all related components (models, services, interfaces, tests) to their story
   - Mark story dependencies; most stories should be independent

2. **From Contracts**: Map each interface contract to the user story it serves. If tests are requested, add contract test tasks [P] before implementation in that story's phase.

3. **From Data Model**: Map each entity to the user story that needs it. If an entity serves multiple stories, place it in the earliest story or the Setup phase. Relationships become service layer tasks in the appropriate story phase.

4. **From Setup/Infrastructure**: Shared infrastructure goes in Phase 1 (Setup). Blocking prerequisites go in Phase 2 (Foundational). Story-specific setup goes within that story's phase.

### Phase Structure

- **Phase 1**: Setup (project initialization)
- **Phase 2**: Foundational (blocking prerequisites that must complete before user stories)
- **Phase 3+**: User Stories in priority order (P1, P2, P3...)
  - Within each story: Tests (if requested) -> Models -> Services -> Endpoints -> Integration
  - Each phase should be a complete, independently testable increment
- **Final Phase**: Polish and Cross-Cutting Concerns
