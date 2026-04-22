# Feature Specification: テスト実行基盤の整備

**Feature Branch**: `068-test-infra`
**Created**: 2026-04-21
**Status**: Draft
**Input**: https://github.com/nakatatsu/psm/issues/68

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: 結合テスト資材（テストスクリプト、テストデータ等）は専用のテストディレクトリに配置されなければならない。ただし、example/以下は、ユーザー向け環境として今回は触らないこと。
- **FR-002**: テストディレクトリの結合テストは、リリースブランチおよびdevelopブランチにマージされた時に結合テストを実行されること
- **FR-004**: CIからAWS開発アカウントへの認証はGitHub Actions OIDCを使用する。IAMロールの信頼ポリシーでリリースブランチおよびdevelopブランチのみにassumeを制限し、他のブランチやforkからの認証を拒否しなければならない
- **FR-005**: 結合テストのFAILはPRのマージをブロックしなければならない
- **FR-006**: CIにより、E2Eテスト用にリリースブランチからビルドしたバイナリが作成されること。またこの際、名前がかぶらないようにすること。

## Constraints _(mandatory)_

- 特記事項無し

## Success Criteria _(mandatory)_

- **SC-001**: developブランチへのpush時に結合テストがCI上で自動実行され、結果が確認できる
- **SC-002**: 結合テストの全14シナリオがCI上でPASS/FAIL/SKIPのいずれかの結果を返す
- **SC-004**: E2Eテスト用バイナリをリリースブランチから準備する手順がCIにより自動化されている

## Edge Cases

- CIで結合テスト実行中にAWSクレデンシャルが失効した場合、失敗でよい
- 結合テストが並行実行された場合、テスト用パラメータ（/myapp/配下）が競合するかもしれないが、したら失敗でよい
- example/のDevContainerのDockerfileに記載されているpsmバージョン（PSM_VERSION）と、テスト対象バイナリのバージョンが異なる場合は問題なので、当然ビルドしたばかりのバイナリを使わねばならない

## Assumptions

- テスト用AWSアカウントでOIDCの設定が行われていること
