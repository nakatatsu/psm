# Tasks: GitHub Releases Binary Distribution

**Input**: Design documents from `/specs/005-release/`

## Phase 1: Setup

- [x] T001 Create `.goreleaser.yaml` at repository root with build config (3 targets, CGO_ENABLED=0, ldflags)

## Phase 2: User Story 1 — Tag Push Release (P1)

**Goal**: `v*` tag push triggers binary build + GitHub Releases publish

- [x] T002 [US1] Create `.github/workflows/release.yml` with tag push trigger (`v*`)
- [x] T003 [US1] Add checkout and setup-go steps
- [x] T004 [US1] Add goreleaser-action step (v7, version '~> v2')
- [x] T005 [US1] Add `go install` to README.md

## Phase 3: Security

- [x] T006 Tag protection ruleset configured via Terraform (bypass: admin only, pattern: `v*`, creation/deletion/update restricted)

## Phase 4: Polish

- [ ] T007 Verify release workflow with test tag (`v0.0.1-test`)
- [ ] T008 Verify downloaded binary runs on target platform
- [ ] T009 Delete test release and tag after verification

## Dependencies

- T001 before T002-T004 (GoReleaser config must exist)
- T002-T004 are sequential (same file)
- T005 is independent
- T006-T008 after T004
