# Corporation API

Go言語で構築されたクリーンアーキテクチャベースのWebアプリケーション（APIサーバー）です。

## 技術スタック

- **言語**: Go 1.21
- **Webフレームワーク**: Chi router
- **データベース**: PostgreSQL
- **スキーマ管理**: sqldef (宣言的スキーマ管理)
- **SQLクエリ生成**: SQLC
- **API仕様**: OpenAPI 3.0 + Redocly
- **APIコード生成**: oapi-codegen
- **アーキテクチャ**: クリーンアーキテクチャ
- **コンテナ**: Docker & Docker Compose

## プロジェクト構造

```
├── cmd/api/                 # メインアプリケーション
├── internal/
│   ├── api/                # 生成されたAPIコード（oapi-codegen）
│   ├── domain/             # ドメイン層（エンティティ、リポジトリインターフェース）
│   ├── usecase/            # ユースケース層（ビジネスロジック）
│   ├── infrastructure/     # インフラストラクチャ層（データベース、外部API）
│   └── presentation/       # プレゼンテーション層（HTTPハンドラー、ルーター）
├── db/
│   ├── schema.sql         # データベーススキーマ（sqldef用）
│   └── queries/           # SQLクエリ（SQLC用）
├── api/                   # OpenAPI定義（分割管理）
│   ├── openapi.yaml       # メインファイル
│   ├── components/        # 再利用可能なコンポーネント
│   │   └── schemas/       # スキーマ定義
│   └── paths/             # エンドポイント定義
├── docker/               # Docker関連ファイル
└── pkg/                  # 共有パッケージ
```

## セットアップ

### 必要な環境

- Go 1.21+
- Docker & Docker Compose
- Make (オプション)

### 開発環境の起動

1. リポジトリをクローン
```bash
git clone <repository-url>
cd corporation-db
```

2. Docker環境を起動
```bash
make docker-up
# または
docker-compose up -d
```

3. アプリケーションを起動
```bash
make run
# または
go run ./cmd/api
```

### Docker環境での完全な起動

```bash
make docker-run
# または
docker-compose up --build
```

## 開発状況

### ✅ 完了済み
1. **Atlas → sqldef マイグレーション**: 宣言的スキーマ管理への移行完了
2. **OpenAPI + oapi-codegen 実装**: 自動コード生成とChi統合完了
3. **Redocly分割ファイル管理**: モジュラーOpenAPI仕様管理完了
4. **データベース統合**: 実際のPostgreSQLクエリによるAPI実装完了

### 🚧 進行中/次のステップ
1. **APIドキュメント プレビューサーバー**: Redoclyプレビューの設定
2. **エラーハンドリングの改善**: より詳細なエラーレスポンス
3. **バリデーション強化**: リクエストデータの検証
4. **テストの追加**: ユニットテスト・統合テストの実装

### 📋 実装完了機能

#### API エンドポイント（実データベース統合済み）
- `GET /health` - ヘルスチェック ✅
- `GET /users` - ユーザー一覧取得 ✅
- `POST /users` - ユーザー作成 ✅
- `GET /users/{id}` - ユーザー詳細取得 ✅

#### データベース機能
- **スキーマ管理**: sqldefによる宣言的管理 ✅
- **クエリ生成**: SQLCによる型安全なクエリ ✅
- **フィールドサポート**: name, email, phone, address ✅
- **NULLable対応**: phone, addressのオプション扱い ✅

#### OpenAPI管理
- **モジュラー設計**: コンポーネント・パス分割 ✅
- **自動バンドリング**: Redoclyによる統合 ✅
- **リンティング**: 仕様品質チェック ✅
- **コード生成**: oapi-codegenによる自動化 ✅

## API開発ワークフロー

### OpenAPI仕様管理

このプロジェクトではRedoclyを使用してOpenAPI仕様を分割管理しています。

```bash
# OpenAPI仕様をlint
make openapi-lint

# 分割ファイルを単一ファイルにバンドル
make openapi-bundle

# APIドキュメントをプレビュー（http://localhost:8081）
make openapi-preview

# APIコードを再生成（OpenAPIから）
make generate-api
```

### スキーマ管理（sqldef）

sqldefによる宣言的スキーマ管理：

```bash
# スキーマ変更の確認（dry-run）
make schema-diff

# スキーマ変更を適用
make schema-apply
```

### コード生成

```bash
# SQLCでデータベースコードを生成
make sqlc-generate

# OpenAPIからAPIコードを生成
make generate-api
```

### エンドポイント

現在実装されているAPI：

- `GET /health` - ヘルスチェック
- `GET /users` - ユーザー一覧取得  
- `POST /users` - ユーザー作成
- `GET /users/{id}` - ユーザー詳細取得

### 使用例

```bash
# ヘルスチェック
curl http://localhost:8080/health

# ユーザー作成
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# ユーザー一覧取得
curl http://localhost:8080/users

# ユーザー詳細取得
curl http://localhost:8080/users/1
```

## 開発

### SQLCでコードを生成

SQLファイルを変更した後、以下のコマンドでコードを再生成します：

```bash
make sqlc-generate
# または
sqlc generate
```

### テスト実行

```bash
make test
# または
go test -v ./...
```

### ビルド

```bash
make build
# または
go build -o bin/api ./cmd/api
```

## 環境変数

以下の環境変数を設定できます（`.env`ファイルまたは環境変数として）：

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=corporation_db
PORT=8080
```

## データベース

PostgreSQL 15を使用しています。マイグレーションファイルは `db/migrations/` にあります。

### データベース接続

```bash
# Dockerコンテナのデータベースに接続
docker-compose exec db psql -U postgres -d corporation_db
```

## コントリビューション

1. フィーチャーブランチを作成
2. 変更を実装
3. テストを追加/更新
4. プルリクエストを作成

## ライセンス

MIT License
