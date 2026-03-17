# Feature Specification: README.md に Usage セクションを追加する

**Feature Branch**: `012-readme-usage`
**Created**: 2026-03-16
**Status**: Draft
**Input**: https://github.com/nakatatsu/psm/issues/12

## User Scenarios & Testing _(mandatory)_

### User Story 1 - psm を初めて使う開発者が基本操作を理解する (Priority: P1)

新しく psm を導入した開発者が README の Usage セクションを読み、`psm sync` コマンドの基本的な使い方（対象ストア指定、プロファイル指定、ドライラン、プルーン）を理解できる。

**Why this priority**: Usage がないと開発者はソースコードやヘルプを読むしかなく、最も基本的なオンボーディング障壁になっている。

**Independent Test**: README の Usage セクションだけを読み、`psm sync` の主要オプションを正しく説明できるか第三者に確認してもらう。

**Acceptance Scenarios**:

1. **Given** README を開いた開発者, **When** Usage セクションを読む, **Then** `psm sync --store ssm secrets.yaml` の基本パターンが理解できる
2. **Given** README を開いた開発者, **When** Usage セクションを読む, **Then** `--store`, `--profile`, `--dry-run`, `--prune` の各オプションの役割が理解できる

---

### User Story 2 - export コマンドの使い方を知る (Priority: P2)

開発者が `psm export` で既存のパラメータを YAML にエクスポートする方法を README から把握できる。

**Why this priority**: export は sync と対になる重要な機能だが、sync より利用頻度は低い。

**Independent Test**: README の export セクションだけを読み、エクスポートコマンドを正しく実行できるか確認する。

**Acceptance Scenarios**:

1. **Given** README を開いた開発者, **When** export の説明を読む, **Then** `psm export --store ssm output.yaml` の使い方が理解できる

---

### User Story 3 - SOPS との組み合わせパターンを知る (Priority: P3)

暗号化された secrets を扱う開発者が、SOPS で復号しながら psm にパイプで渡すパターンを README から把握できる。

**Why this priority**: SOPS 連携は psm の実運用で頻出するパターンだが、psm 単体の使い方を先に理解している前提。

**Independent Test**: README の SOPS セクションだけを読み、パイプパターンのコマンドを正しく組み立てられるか確認する。

**Acceptance Scenarios**:

1. **Given** README を開いた開発者, **When** SOPS 連携の説明を読む, **Then** `sops -d secrets.enc.yaml | psm sync --store ssm /dev/stdin` のパターンが理解できる

---

### User Story 4 - secrets.yaml のフォーマットを確認する (Priority: P2)

開発者が psm に渡す YAML ファイルの構造（キーと値の対応）を README のサンプルから理解できる。

**Why this priority**: 入力ファイルのフォーマットがわからないとコマンドを実行しても期待通りに動かない。export と同等の優先度。

**Independent Test**: README のフォーマット例を見て、正しい secrets.yaml を新規作成できるか確認する。

**Acceptance Scenarios**:

1. **Given** README を開いた開発者, **When** YAML フォーマット例を見る, **Then** psm が期待するキー・値の構造を理解し、自分の secrets.yaml を作成できる

---

### Edge Cases

- README の Usage セクションが既存セクション（Install, Access to AWS）と重複・矛盾しないこと
- コマンド例が実際の psm CLI のオプション体系と一致していること

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: README に Usage セクションが存在し、`psm sync` の基本的な使い方（`--store`, `--profile`, `--dry-run`, `--prune`）を説明していること
- **FR-002**: README に `psm export` の使い方が記載されていること
- **FR-003**: README に SOPS との組み合わせ例（パイプで渡すパターン）が記載されていること
- **FR-004**: README に secrets.yaml のフォーマット例が記載されていること
- **FR-005**: 記載されたコマンド例が psm の実際の CLI オプション体系と一致していること
- **FR-006**: Usage セクションが既存の Install セクションの直後に配置されていること

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: psm を初めて使う開発者が README だけを読んで、5分以内に最初の `psm sync --dry-run` コマンドを正しく組み立てられる
- **SC-002**: README に記載された全コマンド例が、実際に psm の `--help` 出力と矛盾なく対応している
- **SC-003**: example/README.md と本体 README の間でコマンド記法に不整合がない

## Addition

- Write /workspace/README.md in English. Translate it and write README.ja.md in Japanese
