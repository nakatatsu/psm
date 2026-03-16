# Tasks: psm Example Project (Starter Template)

**Input**: Design documents from `/specs/007-docker/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Create directory structure and remove old artifacts

- [ ] T001 Create `example/.devcontainer/` directory structure
- [ ] T002 Remove old `docker/` directory (replaced by `example/`)

---

## Phase 2: User Story 1 - example/ をコピーして psm 管理リポジトリを立ち上げる (Priority: P1) 🎯 MVP

**Goal**: `example/` をコピーするだけで psm + SOPS + age + AWS CLI が揃った DevContainer 環境が手に入り、鍵生成から psm sync まで一連の流れが動く

**Independent Test**: example/ を別ディレクトリにコピー → DevContainer 起動 → 全ツール動作 → age キー生成 → sops 暗号化 → sops 復号 → AWS SSO → psm sync → SSM 確認

### Implementation for User Story 1

- [ ] T003 [P] [US1] Create DevContainer Dockerfile with psm, SOPS, age, AWS CLI v2 (all versions as ARG) in `example/.devcontainer/Dockerfile`
- [ ] T004 [P] [US1] Create devcontainer.json with host `~/.aws` mount and non-root user in `example/.devcontainer/devcontainer.json`
- [ ] T005 [P] [US1] Create sample secrets file with dummy Parameter Store paths in `example/secrets.yaml`
- [ ] T006 [P] [US1] Create SOPS configuration template with age public key placeholder in `example/.sops.yaml.example`
- [ ] T007 [US1] Create README with end-to-end instructions (prerequisites, key generation, encryption, SSO login, psm sync, verification) in `example/README.md`
- [ ] T008 [US1] Create manual verification command list covering full flow (tool checks → age keygen → sops encrypt/decrypt → SSO → psm sync → SSM verify → cleanup) in `example/test.sh`

**Checkpoint**: example/ を別ディレクトリにコピーし、DevContainer を開き、全ツールが動作する

---

## Phase 3: Polish & Cross-Cutting Concerns

**Purpose**: Independence verification and license compliance

- [ ] T009 Verify example/ has no references or dependencies to psm dev repository files (SC-003)
- [ ] T010 Add license information for all bundled tools to Dockerfile comments and README (FR-006)
- [ ] T011 HITL: Copy example/ to a separate directory, open DevContainer, run test.sh to verify full flow

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **User Story 1 (Phase 2)**: Depends on Phase 1 (directory exists)
- **Polish (Phase 3)**: Depends on Phase 2 completion

### Parallel Opportunities

- T003, T004, T005, T006 can all run in parallel (different files, no dependencies)
- T007 and T008 depend on T003-T006 (need to reference actual file contents/tool versions)

---

## Implementation Strategy

### MVP (Single Story)

1. Complete Phase 1: Setup directory, remove old docker/
2. Complete Phase 2: T003-T006 in parallel, then T007-T008
3. Complete Phase 3: Independence check, license info, HITL verification
4. Commit and open PR to merge `007-docker` → `main`

---

## Notes

- Single user story — no incremental delivery needed
- T011 is a HITL (Human-In-The-Loop) task requiring manual DevContainer verification
- No automated tests (Dockerfile/config files only, Test-First N/A)
- Total tasks: 11 (4 parallelizable)
