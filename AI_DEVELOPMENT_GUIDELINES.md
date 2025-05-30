# AI駆動開発ガイドライン

## 概要

本ガイドラインは、**corporatioin-db**プロジェクトの開発を通じて得られた知見をもとに、AI（GitHub Copilot等）を活用した効率的な開発手法をまとめたものです。

## 目次

1. [プロジェクト構成の原則](#プロジェクト構成の原則)
2. [技術スタック選定基準](#技術スタック選定基準)
3. [AI協働開発のベストプラクティス](#ai協働開発のベストプラクティス)
4. [自動化戦略](#自動化戦略)
5. [品質保証](#品質保証)
6. [トラブルシューティング](#トラブルシューティング)
7. [継続的改善](#継続的改善)

---

## プロジェクト構成の原則

### 1. Clean Architecture の採用

```
internal/
├── api/           # OpenAPI生成コード
├── domain/        # ビジネスロジック・エンティティ
├── infrastructure/ # 外部依存（DB、外部API等）
├── presentation/  # HTTP層・ルーティング
└── usecase/       # アプリケーション層
```

**AIとの協働メリット:**
- 構造が明確でAIが理解しやすい
- 各層の責任範囲が明確で的確な提案を得られる
- テストコード生成時に依存関係を正しく認識できる

### 2. 設定ファイルの集約

プロジェクトルートに設定ファイルを配置し、AIが全体像を把握しやすくする:

```
project-root/
├── docker-compose.yml    # 開発環境定義
├── sqlc.yaml            # データベースコード生成
├── oapi-codegen.yaml    # API生成設定
├── redocly.yaml         # OpenAPI管理
└── Makefile             # 開発コマンド集約
```

### 3. ドキュメント駆動開発

- **README.md**: プロジェクト概要、セットアップ、使用方法
- **API仕様**: OpenAPIでAPI-First設計
- **スキーマ定義**: 宣言的スキーマ管理（sqldef）

---

## 技術スタック選定基準

### 1. AI協働に適したツール選択

| 分野 | 選択技術 | AI協働の利点 |
|------|----------|-------------|
| API設計 | OpenAPI 3.0 | 仕様から自動コード生成、AIが仕様を理解しやすい |
| データベース | SQLC + sqldef | 宣言的管理、型安全、AIが構造を理解しやすい |
| Web框架 | Chi Router | シンプルで予測可能、生成コードとの親和性が高い |
| 開発環境 | Docker Compose | 環境統一、AIが環境差異に悩まされない |
| 文档管理 | Redocly | モジュラー管理、大規模APIでもAIが処理しやすい |

### 2. 自動生成重視の選択

**採用理由:**
- 手動コードを最小化し、AIによる生成・更新を容易にする
- 設定ファイルベースで変更意図をAIが理解しやすい
- 型安全性により、AIが生成したコードの品質を保証

---

## AI協働開発のベストプラクティス

### 1. コンテキスト提供の技法

#### ✅ 良い依頼方法

```markdown
**現在の状況:**
- Atlas から sqldef に移行済み
- OpenAPI + oapi-codegen で API 生成済み
- Docker-only 戦略採用

**目標:**
- SQLC の重複関数定義エラーを解決
- phone、address フィールドに対応したクエリ生成

**制約:**
- 既存の API 仕様は変更不可
- Docker 環境での動作必須
```

#### ❌ 避けるべき依頼

```markdown
エラーが出ています。直してください。
```

### 2. 段階的な開発アプローチ

1. **設計フェーズ**: 全体構造をAIと議論
2. **実装フェーズ**: 小さな単位で逐次実装
3. **統合フェーズ**: AIにテスト・検証支援を依頼
4. **改善フェーズ**: AIと共に問題点を特定・解決

### 3. ファイル変更の原則

- **1つの責務**: 1回の変更で1つの機能・修正に集中
- **完全なコンテキスト**: 変更対象ファイルの全体像を提供
- **依存関係の明示**: 影響範囲をAIに伝える

---

## 自動化戦略

### 1. 開発フロー自動化

```makefile
# 開発の基本フロー
.PHONY: dev-setup
dev-setup: ## 開発環境初期化
	docker-compose build
	docker-compose run --rm api sqldef --dry-run

.PHONY: generate
generate: ## コード生成
	docker run --rm -v $(PWD):/workspace -w /workspace \
		redocly/cli:latest bundle api/openapi.yaml -o api/openapi-bundled.yaml
	go generate ./...
	docker run --rm -v $(PWD):/workspace -w /workspace \
		sqlc/sqlc:latest generate

.PHONY: test
test: ## テスト実行
	docker-compose run --rm api go test ./...
```

### 2. AI支援のための自動化

- **依存関係可視化**: 定期的に依存グラフを生成
- **コード品質チェック**: linter、formatter の自動実行
- **ドキュメント生成**: API文档の自動更新

---

## 品質保証

### 1. AI生成コードの検証

#### 必須チェックポイント

- [ ] **型安全性**: コンパイルエラーがないか
- [ ] **API契約**: OpenAPI仕様との整合性
- [ ] **データベース整合性**: スキーマとクエリの一致
- [ ] **エラーハンドリング**: 適切な例外処理

#### 自動検証の設定

```yaml
# .github/workflows/ai-validation.yml
name: AI Generated Code Validation
on: [push, pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate OpenAPI
        run: redocly lint api/openapi.yaml
      - name: Validate SQLC
        run: sqlc vet
      - name: Build Check
        run: docker-compose build
```

### 2. 段階的デプロイメント

1. **開発環境**: AI生成コードの基本動作確認
2. **ステージング環境**: 統合テスト
3. **本番環境**: 段階的ロールアウト

---

## トラブルシューティング

### 1. よくある問題と対処法

#### SQLC重複定義エラー

**症状:** 同じ関数が複数回定義される
**原因:** キャッシュや古いファイルの残存
**対処法:**
```bash
# 完全クリーンアップ
rm -rf internal/infrastructure/db/*.sql.go
docker run --rm -v $(PWD):/workspace -w /workspace sqlc/sqlc:latest generate
```

#### OpenAPI生成コードの型不一致

**症状:** 生成されたGoコードが期待通りでない
**原因:** OpenAPI仕様の曖昧性
**対処法:**
```yaml
# より厳密な型定義
components:
  schemas:
    User:
      type: object
      required: [id, name, email]  # 必須フィールドを明示
      properties:
        phone:
          type: string
          nullable: true  # NULL許可を明示
```

#### Docker ビルドエラー

**症状:** コンテナビルド失敗
**原因:** 依存関係の不整合、ファイル権限
**対処法:**
```bash
# キャッシュクリア
docker-compose build --no-cache
# 権限問題の解決
chmod +x scripts/*.sh
```

### 2. AI協働時のデバッグ手法

1. **エラー共有**: 完全なエラーメッセージとスタックトレースを提供
2. **環境情報**: Go版本、Docker版本、OS情報を明示
3. **再現手順**: AIが問題を再現できる手順を提供

---

## 継続的改善

### 1. 定期的な技術更新

- **月次**: 依存関係の更新確認
- **四半期**: 新しいAI開発ツールの評価
- **半年**: アーキテクチャの見直し

### 2. AI協働の改善

#### 効果測定指標

- **開発速度**: 機能実装までの時間
- **品質**: バグ発生率、テストカバレッジ
- **保守性**: コード変更の影響範囲

#### 改善サイクル

1. **振り返り**: AI協働での問題点抽出
2. **実験**: 新しいアプローチの試行
3. **評価**: 効果測定と判断
4. **標準化**: 有効な手法のガイドライン化

---

## 付録

### A. 推奨開発環境

```json
// .vscode/settings.json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "files.associations": {
    "*.yaml": "yaml",
    "docker-compose*.yml": "dockercompose"
  }
}
```

### B. 有用なVSCode拡張機能

- **Go**: Go言語サポート
- **Docker**: コンテナ管理
- **OpenAPI (Swagger) Editor**: API仕様編集
- **SQLTools**: データベース管理
- **GitHub Copilot**: AI支援

### C. AI協働チェックリスト

開発開始時:
- [ ] プロジェクト構造をAIに説明
- [ ] 現在の技術スタックを明示
- [ ] 制約条件・要件を整理

実装中:
- [ ] 小さな単位で逐次確認
- [ ] エラーは即座に共有
- [ ] 生成コードの検証を実施

完了時:
- [ ] 統合テスト実行
- [ ] ドキュメント更新
- [ ] 次回開発のメモ記録

---

**作成日**: 2025年5月30日  
**バージョン**: 1.0  
**プロジェクト**: corporatioin-db  
**更新履歴**: 初版作成
