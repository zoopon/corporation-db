# Corporation API

Go言語で構築されたクリーンアーキテクチャベースのWebアプリケーション（APIサーバー）です。

## 技術スタック

- **言語**: Go 1.21
- **Webフレームワーク**: Chi router
- **データベース**: PostgreSQL
- **ORM**: SQLC
- **API仕様**: OpenAPI 3.0
- **アーキテクチャ**: クリーンアーキテクチャ
- **コンテナ**: Docker & Docker Compose

## プロジェクト構造

```
├── cmd/api/                 # メインアプリケーション
├── internal/
│   ├── domain/             # ドメイン層（エンティティ、リポジトリインターフェース）
│   ├── usecase/            # ユースケース層（ビジネスロジック）
│   ├── infrastructure/     # インフラストラクチャ層（データベース、外部API）
│   └── presentation/       # プレゼンテーション層（HTTPハンドラー、ルーター）
├── db/
│   ├── migrations/         # データベースマイグレーション
│   └── queries/           # SQLクエリ（SQLC用）
├── api/                   # OpenAPI定義
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

## API仕様

OpenAPI 3.0仕様は `api/openapi.yaml` にあります。

### エンドポイント

- `GET /health` - ヘルスチェック
- `GET /users` - ユーザー一覧取得
- `POST /users` - ユーザー作成
- `GET /users/{id}` - ユーザー詳細取得
- `PUT /users/{id}` - ユーザー更新
- `DELETE /users/{id}` - ユーザー削除

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
