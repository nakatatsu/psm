# Feature Specification: Distinguish dry-run output from actual execution

**Feature Branch**: `008-dryrun-output`
**Created**: 2026-03-16
**Status**: Draft
**Input**: https://github.com/nakatatsu/psm/issues/10

## User Scenarios & Testing *(mandatory)*

### User Story 1 - dry-run 実行時に出力から即座に判別できる (Priority: P1)

ユーザーが `psm sync --dry-run` を実行したとき、出力を一目見るだけで「これは dry-run であり、実際には何も変更されていない」と判別できる。現状は dry-run と実行時の出力が同一で、誤認のリスクがある。

**Why this priority**: dry-run は「安全に差分を確認する」ためのもの。出力が本番実行と区別できなければ、その目的を果たせない。誤認による運用事故（適用済みと思い込む / 適用されたと焦る）を防ぐ。

**Independent Test**: `psm sync --store ssm --dry-run` を実行し、出力に dry-run であることを示す表記が含まれることを目視確認する。

**Acceptance Scenarios**:

1. **Given** 変更対象のパラメータが存在する状態, **When** `psm sync --dry-run` を実行する, **Then** 各アクション行とサマリー行に dry-run であることを示す表記が含まれる
2. **Given** 変更対象のパラメータが存在する状態, **When** `psm sync`（dry-run なし）を実行する, **Then** 出力に dry-run の表記が含まれない
3. **Given** 変更がない状態, **When** `psm sync --dry-run` を実行する, **Then** サマリー行に dry-run であることを示す表記が含まれる

### Edge Cases

- `--dry-run` なしの通常実行の出力フォーマットが変わっていないこと（後方互換性）
- `--dry-run` と `--prune` を組み合わせた場合も delete 行に dry-run 表記が出ること

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `--dry-run` 実行時、サマリー行に dry-run であることを示す表記を付与しなければならない
- **FR-002**: `--dry-run` 実行時、各アクション行（create / update / delete）にも dry-run であることを示す表記を付与しなければならない
- **FR-003**: `--dry-run` なしの通常実行時は、現行の出力フォーマットを維持しなければならない（後方互換性）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: dry-run の出力を見たユーザーが、実行済みの出力と混同しない（出力テキストに明示的な dry-run 表記がある）
- **SC-002**: 通常実行の出力が既存のフォーマットから変更されていない

## Assumptions

- 色付き出力（ANSI カラー）は使わない。psm は CI/CD パイプラインでも使われるため、プレーンテキストで判別できる必要がある
- 出力先が stderr / stdout の使い分けは現行のまま
