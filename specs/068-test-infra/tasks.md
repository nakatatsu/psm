# Tasks: テスト実行基盤の整備

**Input**: Design documents from `/specs/068-test-infra/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Organization**: このfeatureはUser Story形式ではなく要件ベースのため、FRに沿ってフェーズを構成する。

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: テストディレクトリの作成とテスト資材の移設

- [x] T001 Create `tests/integration/` directory
- [x] T002 [P] Copy `example/test.sh` to `tests/integration/test.sh`
- [x] T003 [P] Copy `example/secrets.example.yaml` to `tests/integration/secrets.example.yaml`
- [x] T004 Update `tests/integration/test.sh` — パスは既にSCRIPT_DIR相対のため変更不要。スキップ
- [ ] T005 Verify `tests/integration/test.sh` がローカルで実行可能であること（DevContainer環境で手動確認）

**Checkpoint**: テスト資材がtests/integration/に独立して配置されている

---

## Phase 2: CI結合テストワークフロー（FR-002, FR-004, FR-005）

**Purpose**: GitHub Actionsで結合テストを自動実行するワークフローの作成

- [x] T006 Create `.github/workflows/integration-test.yml` with the following configuration:
  - Trigger: `push` to `develop` and `release-*` branches
  - Permissions: `id-token: write`, `contents: read`
  - Job steps:
    1. `actions/checkout@v6`
    2. `actions/setup-go@v6` (go 1.26.1)
    3. `go build -o psm .`
    4. Install SOPS and age
    5. `aws-actions/configure-aws-credentials` with OIDC (role-to-assume from secrets)
    6. Run `bash tests/integration/test.sh` with `PSM_BIN` and `PSM_TEST_PROFILE` env vars
- [ ] T007 Add GitHub repository secrets/variables (manual setup required):
  - Secret `AWS_ROLE_ARN` — OIDC用IAMロールのARN
  - Variable `AWS_REGION` — テスト用リージョン
  - Variable `AWS_ACCOUNT_ID` — 期待するAWSアカウントID（test.sh の PSM_EXPECTED_ACCOUNT_ID に供給、defense-in-depth チェックに使用）

**Checkpoint**: developブランチへのpushで結合テストが自動実行される

---

## Phase 3: E2Eバイナリビルド（FR-006）

**Purpose**: リリースブランチからE2Eテスト用バイナリをCIで自動ビルド

- [x] T008 Add E2E binary build job to `.github/workflows/integration-test.yml`:
  - Trigger: `push` to `release-*` branches only
  - Steps:
    1. `actions/checkout@v6`
    2. `actions/setup-go@v6`
    3. `go build -ldflags "-X main.version=${GITHUB_SHA::7}" -o psm_${GITHUB_SHA::7}_linux_amd64 .`
    4. `actions/upload-artifact` to store the binary
- [ ] T009 Verify artifact is downloadable from GitHub Actions UI after a release branch push (manual verification after merge)

**Checkpoint**: リリースブランチへのpush時にE2Eバイナリがartifactとして作成される

---

## Phase 4: test.shのCI対応調整

**Purpose**: test.shをCI環境（非インタラクティブ）で動作させるための調整

- [x] T010 Review `tests/integration/test.sh` for CI compatibility
- [x] T011 Modify `tests/integration/test.sh` — CI=true時にsafety gateスキップ、PROFILE_FLAGをオプション化
- [ ] T012 Test the workflow end-to-end (manual verification after merge)

**Checkpoint**: test.shがCI環境で全シナリオ実行可能

---

## Phase 5: Polish

**Purpose**: ドキュメント更新と最終確認

- [x] T013 [P] Update `documents/backend/test-procedure.md` — 結合テストセクションをCI自動実行の記述に更新
- [x] T014 [P] Update `specs/068-test-infra/quickstart.md` — 既に最新、変更不要
- [ ] T015 Run quickstart.md validation（手順が正しく動作するか確認 — マージ後に実施）

**Checkpoint**: ドキュメントが最新状態

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (CI Workflow)**: Depends on Phase 1 completion
- **Phase 3 (E2E Binary)**: Depends on Phase 2 (same or related workflow file)
- **Phase 4 (CI Adjustment)**: Depends on Phase 2 (workflow must exist to test against)
- **Phase 5 (Polish)**: Depends on all previous phases

### Parallel Opportunities

- T002, T003 can run in parallel (different files)
- T013, T014 can run in parallel (different files)
- Phase 3 and Phase 4 can be worked in parallel after Phase 2

---

## Implementation Strategy

### MVP First

1. Phase 1: テスト資材の移設
2. Phase 2: CIワークフロー作成
3. Phase 4: test.shのCI対応
4. **VALIDATE**: developにpushして結合テストが動くことを確認

### Full Delivery

5. Phase 3: E2Eバイナリビルド追加
6. Phase 5: ドキュメント更新

---

## Notes

- AWS OIDCの前提条件（IAMロール、OIDCプロバイダ）はリポジトリ外の手動設定が必要
- FR-005（FAILでマージブロック）はGitHub branch protection rulesの手動設定が必要
- test.shのsafety gate（手動承認）はCI環境では自動スキップさせる必要がある
