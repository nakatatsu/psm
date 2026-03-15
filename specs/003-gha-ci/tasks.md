# Tasks: GitHub Actions CI

**Input**: Design documents from `/specs/003-gha-ci/`
**Prerequisites**: plan.md, spec.md, research.md, survey.md

**Tests**: Not applicable — this feature is a YAML workflow file, not Go code. Validation is done by creating a PR and observing CI behavior.

**Organization**: Single user story (P1: PR quality gate). No foundational phase needed — the project and branch protection already exist.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to

---

## Phase 1: Setup

**Purpose**: Prepare the CI configuration infrastructure

- [x] T001 Create `.github/workflows/` directory at repository root

---

## Phase 2: User Story 1 — PR Quality Gate (Priority: P1)

**Goal**: PR to main triggers automated checks; results visible on PR

**Independent Test**: Create a PR to main and confirm all checks run and results appear in Checks tab

### Implementation

- [x] T002 [US1] Update `.golangci.yml` to add `formatters:` section with goimports enabled
- [x] T003 [US1] Create `.github/workflows/ci.yml` with workflow triggers (pull_request + push to main) and job name `ci`
- [x] T004 [US1] Add checkout and setup-go steps (Go 1.26.1, module caching)
- [x] T005 [US1] Add gofumpt install step (pinned to v0.9.2) and format check step
- [x] T006 [US1] Add golangci-lint-action step (v2.9.0, uses `.golangci.yml`)
- [x] T007 [US1] Add `go test -race ./...` step
- [x] T008 [US1] Add govulncheck step with `continue-on-error: true`
- [x] T009 [US1] Add `go build ./...` step (place before lint/test steps)

**Checkpoint**: Push branch, create PR to main, verify CI runs all steps

---

## Phase 3: Polish

**Purpose**: Verify and clean up

- [ ] T010 Verify CI passes on a clean PR (all checks green)
- [ ] T011 Verify CI fails on a PR with intentional gofumpt violation
- [ ] T012 Verify CI fails on a PR with intentional test failure
- [x] T013 Remove `.tmp/check.md` if CI is confirmed working (already moved to .tmp/)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **User Story 1 (Phase 2)**: Depends on T001 (directory exists)
- **Polish (Phase 3)**: Depends on Phase 2 completion

### Within Phase 2

```
T002 (.golangci.yml update) ─── no dependency on other tasks
T003 (ci.yml skeleton)
  └─ T004 (checkout + setup-go)
       └─ T005 (gofumpt)      ┐
       └─ T006 (golangci-lint) ├─ all added to ci.yml sequentially
       └─ T007 (test)          │
       └─ T008 (govulncheck)   │
       └─ T009 (build)         ┘
```

T002 and T003 can be done in parallel (different files). T004-T009 are sequential additions to the same file (ci.yml).

---

## Implementation Strategy

### MVP

1. T001 → T002 + T003 (parallel) → T004-T009 (sequential) → push and create PR
2. Validate CI runs → T010-T012

### Estimated Effort

Total: 13 tasks. The core deliverable is 2 files (`.github/workflows/ci.yml` and `.golangci.yml` update). Implementation is straightforward YAML authoring.
