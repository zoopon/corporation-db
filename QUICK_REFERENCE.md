# AI協働開発 クイックリファレンス

## 🚀 開発開始時のテンプレート

### AI への依頼テンプレート

```markdown
## 現在の状況
- **プロジェクト**: corporatioin-db (Go + PostgreSQL + Docker)
- **アーキテクチャ**: Clean Architecture
- **技術スタック**: 
  - API: OpenAPI + oapi-codegen + Chi
  - DB: SQLC + sqldef
  - 環境: Docker Compose
- **最後の作業**: [前回の作業内容を記載]

## 実装したい機能
[具体的な機能説明]

## 制約・要求事項
- [ ] 既存API仕様の維持
- [ ] Docker環境での動作
- [ ] 型安全性の保証
- [ ] [その他の制約]

## 期待する成果物
- [ ] [具体的な成果物1]
- [ ] [具体的な成果物2]
```

---

## 🔧 よく使う開発コマンド

### 基本セットアップ
```bash
# 依存関係インストール・ビルド
make docker-run

# 開発環境起動
make docker-run

# 全生成処理実行
make generate-all

# gBizINFOデータのインポート（バッチ処理）
make batch-run
```

### コード生成
```bash
# OpenAPI バンドル生成
make openapi-bundle

# API コード生成
make generate-api

# SQLC コード生成  
make sqlc-generate

# 全生成（推奨）
make generate-all
```

### トラブルシューティング
```bash
# Docker 完全クリーンアップ
docker compose down --volumes --rmi all
docker system prune -af

# SQLC クリーンアップ
rm -rf internal/infrastructure/db/*.sql.go
make sqlc-generate

# 依存関係更新
go mod tidy
go mod download
```

---

## 📋 品質チェックリスト

### コード生成後の確認項目

#### ✅ OpenAPI 関連
- [ ] `api/openapi-bundled.yaml` が最新
- [ ] `internal/api/generated.go` にエラーなし
- [ ] API エンドポイントが期待通り生成されている
- [ ] リクエスト・レスポンス型が正しい

#### ✅ SQLC 関連  
- [ ] `internal/infrastructure/db/models.go` が最新スキーマ反映
- [ ] `internal/infrastructure/db/users.sql.go` に重複関数なし
- [ ] 全クエリ関数が生成されている
- [ ] NULL可能フィールドが `sql.NullString` 等で適切に定義

#### ✅ Docker 関連
- [ ] `docker compose build` が成功
- [ ] `docker compose up` でサービス起動
- [ ] ヘルスチェックエンドポイント応答確認
- [ ] データベース接続確認
- [ ] 法人情報API（`/corporations`）の動作確認
- [ ] 都道府県フィルタリング（`?prefecture_code=13`等）の動作確認

---

## 🚨 エラー対処法

### 1. SQLC 重複定義エラー
```
duplicate function definition: GetUserByID
```

**対処法:**
```bash
cd /Users/zoo/projects/corporatioin-db
rm -rf internal/infrastructure/db/*.sql.go
docker run --rm -v $(pwd):/workspace -w /workspace sqlc/sqlc:1.29.0 generate
```

### 2. OpenAPI バンドルエラー
```
Error bundling OpenAPI: $ref resolution failed
```

**対処法:**
```bash
# パス確認
ls -la api/components/schemas/
ls -la api/paths/

# 手動バンドル
docker run --rm -v $(pwd):/workspace -w /workspace \
  redocly/cli:latest bundle api/openapi.yaml -o api/openapi-bundled.yaml
```

### 3. Docker ビルドエラー
```
no Go files in /workspace
```

**対処法:**
```bash
# Dockerfile 確認
cat docker/Dockerfile

# 権限確認
ls -la cmd/api/main.go

# キャッシュクリア再ビルド
docker compose build --no-cache
```

---

## 📚 ファイル構造クイックマップ

```
📦 corporatioin-db/
├── 🔧 設定ファイル群
│   ├── docker compose.yml     # 開発環境定義
│   ├── Makefile              # 開発コマンド
│   ├── sqlc.yaml             # DB コード生成設定
│   ├── oapi-codegen.yaml     # API コード生成設定
│   └── redocly.yaml          # OpenAPI 管理設定
│
├── 📋 API 仕様
│   ├── api/openapi.yaml      # メイン仕様（$ref使用）
│   ├── api/openapi-bundled.yaml  # バンドル済み仕様
│   ├── api/components/       # 再利用可能コンポーネント
│   └── api/paths/            # エンドポイント定義
│
├── 🗄️ データベース
│   ├── db/schema.sql         # テーブル定義（sqldef用）
│   └── db/queries/corporations.sql  # SQL クエリ（SQLC用）
│
├── 🏗️ アプリケーション
│   ├── cmd/api/main.go       # APIサーバーエントリーポイント
│   ├── cmd/download/main.go  # gBizINFOバッチエントリーポイント
│   ├── internal/api/         # 生成されたAPIコード
│   ├── internal/domain/      # ビジネスロジック
│   ├── internal/infrastructure/  # 外部依存・データアクセス
│   ├── internal/presentation/    # HTTP層
│   └── internal/usecase/     # アプリケーション層
│
└── 📖 ドキュメント
    ├── README.md                        # プロジェクト概要
    ├── AI_DEVELOPMENT_GUIDELINES.md     # AI開発ガイド
    └── QUICK_REFERENCE.md               # このファイル
```

---

## 🤖 AI との効果的な対話例

### 新機能追加時

**👤 開発者:**
```markdown
現在のcorporatioin-dbプロジェクトに、ユーザーの所属部署管理機能を追加したいです。

**要件:**
- ユーザーは複数の部署に所属可能
- 部署には名前、説明、作成日時が必要
- 既存のユーザーAPIに部署情報を含める

**制約:**
- 既存のAPI仕様を破壊しない
- NULLフィールドの適切な処理
- Clean Architectureの維持

どのような手順で実装すべきでしょうか？
```

**🤖 AI の期待回答構造:**
1. データベーススキーマ設計提案
2. OpenAPI仕様更新案
3. 実装順序の提案
4. 影響範囲の分析

### エラー解決時

**👤 開発者:**
```markdown
Docker ビルド時に以下のエラーが発生しています：

```
[ERROR] /workspace/internal/infrastructure/db/users.sql.go:76:1: 
func (q *Queries) GetUserByID redeclared in this block
```

**現在の状況:**
- SQLC で users.sql から生成したコード
- 最近 phone, address フィールドを追加
- sqlc.yaml は確認済みで正しい設定

**試行済み:**
- docker compose build --no-cache
- rm -rf internal/infrastructure/db/*.sql.go

どのような原因と解決策が考えられますか？
```

**🤖 AI の期待回答構造:**
1. 問題の原因分析
2. 段階的な解決手順
3. 今後の予防策
4. 関連する設定確認ポイント

---

## 📈 進捗トラッキング

### 機能実装の進捗管理

```markdown
## [機能名] 実装進捗

### ✅ 完了
- [ ] データベーススキーマ設計
- [ ] OpenAPI仕様更新  
- [ ] SQLC クエリ作成
- [ ] ビジネスロジック実装
- [ ] API ハンドラー実装
- [ ] 統合テスト
- [ ] ドキュメント更新

### 🚧 進行中
- [ ] [現在作業中の項目]

### ⏳ 待機中  
- [ ] [今後の予定項目]

### ❌ ブロック中
- [ ] [問題により停止中の項目] - [問題詳細]
```

---

## 🔗 API使用例

### 法人情報API

```bash
# 全法人情報を取得
curl http://localhost:8080/corporations

# 東京都の法人のみ取得 (prefecture_code=13)
curl "http://localhost:8080/corporations?prefecture_code=13"

# 北海道の法人のみ取得 (prefecture_code=01) 
curl "http://localhost:8080/corporations?prefecture_code=01"

# 沖縄県の法人のみ取得 (prefecture_code=47)
curl "http://localhost:8080/corporations?prefecture_code=47"

# 複数条件での検索例
curl "http://localhost:8080/corporations?name=株式会社&prefecture_code=13"
```

### データインポート確認

```bash
# gBizINFOデータをダウンロード・インポート
make batch-run

# インポート後の確認
curl "http://localhost:8080/corporations?limit=5"

# 特定都道府県のデータ確認
curl "http://localhost:8080/corporations?prefecture_code=13&limit=10"
```

### 都道府県コード参照表

| コード | 都道府県 | 使用例 |
|-------|----------|--------|
| 01 | 北海道 | `?prefecture_code=01` |
| 13 | 東京都 | `?prefecture_code=13` |
| 27 | 大阪府 | `?prefecture_code=27` |
| 40 | 福岡県 | `?prefecture_code=40` |
| 47 | 沖縄県 | `?prefecture_code=47` |

*全47都道府県のコード一覧は[README.md](README.md)を参照*

---

## 🐛 トラブルシューティング: 都道府県フィルタリング

### 都道府県コードが正しく抽出されない場合

```bash
# CSVデータの都道府県フィールド確認
head -5 data/gbiz_*.csv | cut -d',' -f4

# データベースの都道府県コード分布確認
docker compose exec db psql -U postgres -d corporation_db -c \
  "SELECT prefecture_code, COUNT(*) FROM corporations GROUP BY prefecture_code ORDER BY prefecture_code;"
```

### APIフィルタリングが効かない場合

```bash
# エンドポイント仕様確認
curl -v "http://localhost:8080/corporations?prefecture_code=13"

# データベース直接確認
docker compose exec db psql -U postgres -d corporation_db -c \
  "SELECT name, location, prefecture_code FROM corporations WHERE prefecture_code = '13' LIMIT 5;"
```

---

**最終更新**: 2025年5月31日  
**バージョン**: 1.1 - Prefecture Code Filtering対応
