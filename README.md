# Corporation API

Go言語で構築されたクリーンアーキテクチャベースのWebアプリケーション（APIサーバー）です。

## 📚 ドキュメント

- **[AI開発ガイドライン](AI_DEVELOPMENT_GUIDELINES.md)** - AI（GitHub Copilot等）を活用した開発のベストプラクティス
- **[クイックリファレンス](QUICK_REFERENCE.md)** - 日常開発で使用するコマンドとトラブルシューティング

## 技術スタック

- **言語**: Go 1.21
- **Webフレームワーク**: Chi router
- **データベース**: PostgreSQL
- **スキーマ管理**: sqldef (宣言的スキーマ管理)
- **SQLクエリ生成**: SQLC
- **API仕様**: OpenAPI 3.0 + Redocly
- **APIコード生成**: oapi-codegen
- **アーキテクチャ**: クリーンアーキテクチャ
- **実行環境**: Docker & Docker Compose

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
├── pkg/                  # 共有パッケージ
├── AI_DEVELOPMENT_GUIDELINES.md  # AI協働開発ガイド
└── QUICK_REFERENCE.md     # 開発クイックリファレンス
```

## セットアップ

### 必要な環境

- Docker & Docker Compose
- Make (オプション)

### アプリケーション起動

1. リポジトリをクローン
```bash
git clone <repository-url>
cd corporation-db
```

2. Docker環境でアプリケーションを起動
```bash
make docker-run
# または
docker-compose up --build
```

アプリケーションは http://localhost:8080 でアクセス可能になります。

## 📋 Docker-Compose 実行手順書

### 🚀 gBizINFO データ取得からAPI起動までの完全手順

Docker-Composeを使用して、データ取得・インポート・アプリケーション起動を順次実行する手順です。

#### 📚 前提条件
- Docker & Docker Compose がインストール済み
- プロジェクトのルートディレクトリにいること

---

### 🗂️ ステップ 1: データベース起動

```bash
# データベースを起動
docker-compose up -d db

# データベースの起動確認（約10-15秒待機）
docker-compose logs db

# データベース接続確認
docker-compose exec db psql -U postgres -d corporation_db -c "\dt"
```

**確認ポイント:**
- ✅ PostgreSQLコンテナが起動している
- ✅ データベースに接続できる
- ✅ テーブルが作成されている

---

### 📥 ステップ 2: データ取得バッチ実行

```bash
# ダウンロード用コンテナでデータ取得（推奨）
docker-compose run --rm download ./download \
  -output /data/gbiz_$(date +%Y%m%d_%H%M%S).zip

# または、固定ファイル名でダウンロード
docker-compose run --rm download ./download \
  -output /data/gbiz_latest.zip

# Makefileターゲット使用（推奨）
make download-data-docker
```

**確認ポイント:**
- ✅ ZIPファイルが `data/` ディレクトリに保存される
- ✅ ファイルサイズが適切（通常数百MB〜数GB）
- ✅ CSV headers が表示される

---

### 📊 ステップ 3: データインポートバッチ実行

```bash
# ダウンロードしたZIPファイルをインポート
docker-compose run --rm import go run ./cmd/import \
  -input /data/gbiz_latest.zip

# または、特定のタイムスタンプ付きファイルをインポート
docker-compose run --rm import go run ./cmd/import \
  -input /data/gbiz_20241231_120000.zip

# dry-runでインポート内容を事前確認（推奨）
docker-compose run --rm import go run ./cmd/import \
  -input /data/gbiz_latest.zip \
  -dry-run
```

**確認ポイント:**
- ✅ CSVファイルが正常に解析される
- ✅ 都道府県コードが自動抽出される
- ✅ データベースにレコードが保存される

---

### 🌐 ステップ 4: アプリケーション起動

```bash
# APIサーバーを起動
docker-compose up -d app

# ログ確認
docker-compose logs -f app

# ヘルスチェック
curl http://localhost:8080/health
```

**確認ポイント:**
- ✅ APIサーバーがポート8080で起動
- ✅ ヘルスチェックが `{"status":"ok"}` を返す
- ✅ データベース接続が正常

---

### 🔍 ステップ 5: 動作確認

#### API エンドポイントテスト

```bash
# 全法人情報取得
curl "http://localhost:8080/corporations?limit=5"

# 都道府県フィルタリング（東京都）
curl "http://localhost:8080/corporations?prefecture_code=13&limit=5"

# 都道府県フィルタリング（大阪府）
curl "http://localhost:8080/corporations?prefecture_code=27&limit=5"

# 法人名検索 + 都道府県フィルタ
curl "http://localhost:8080/corporations?name=株式会社&prefecture_code=13&limit=5"

# 特定法人番号で詳細取得
curl "http://localhost:8080/corporations/1234567890123"
```

---

### 📊 全サービス管理コマンド

#### 全サービス起動
```bash
# 全サービス起動（バックグラウンド）
docker-compose up -d

# 全サービス起動（フォアグラウンド）
docker-compose up
```

#### サービス状態確認
```bash
# サービス状態確認
docker-compose ps

# ログ確認
docker-compose logs

# 特定サービスのログ確認
docker-compose logs app
docker-compose logs db
```

#### サービス停止・クリーンアップ
```bash
# 全サービス停止
docker-compose down

# データ保持して停止
docker-compose down

# データも削除して停止（注意！）
docker-compose down -v

# 不要なイメージも削除
docker-compose down --rmi all
```

---

### 🚨 トラブルシューティング

#### 🗃️ データベースリセット

データベースを完全にリセットする場合の手順：

```bash
# 方法1: ボリュームを削除してデータベースを完全リセット
docker-compose down -v
docker-compose up -d db

# 方法2: 既存のデータベースコンテナを削除して再構築
docker-compose down
docker volume rm corporatioin-db_postgres_data
docker-compose up -d db

# 方法3: テーブルのデータのみを削除（テーブル構造は保持）
docker-compose exec db psql -U postgres -d corporation_db -c "TRUNCATE TABLE corporations, users RESTART IDENTITY CASCADE;"

# 方法4: 特定のテーブルのみリセット（corporationsテーブルのみ）
docker-compose exec db psql -U postgres -d corporation_db -c "TRUNCATE TABLE corporations RESTART IDENTITY CASCADE;"

# リセット後の確認
docker-compose exec db psql -U postgres -d corporation_db -c "SELECT COUNT(*) FROM corporations;"
docker-compose exec db psql -U postgres -d corporation_db -c "SELECT COUNT(*) FROM users;"
```

**⚠️ 注意事項:**
- `docker-compose down -v` は **すべてのデータが永久に削除** されます
- `TRUNCATE` コマンドは **データのみ削除**、テーブル構造は保持されます
- 本番環境では絶対に実行しないでください

#### 📊 データベース状態確認

```bash
# データベース内のテーブル一覧
docker-compose exec db psql -U postgres -d corporation_db -c "\dt"

# 各テーブルのレコード数確認
docker-compose exec db psql -U postgres -d corporation_db -c "
  SELECT 'corporations' as table_name, COUNT(*) as record_count FROM corporations
  UNION ALL
  SELECT 'users' as table_name, COUNT(*) as record_count FROM users;"

# 都道府県別の法人数統計
docker-compose exec db psql -U postgres -d corporation_db -c "
  SELECT prefecture_code, COUNT(*) as corporation_count 
  FROM corporations 
  WHERE prefecture_code IS NOT NULL 
  GROUP BY prefecture_code 
  ORDER BY prefecture_code;"

# データベースサイズ確認
docker-compose exec db psql -U postgres -d corporation_db -c "
  SELECT pg_size_pretty(pg_database_size('corporation_db')) as database_size;"
```

#### データベース接続エラー
```bash
# データベースコンテナの状態確認
docker-compose logs db

# データベース再起動
docker-compose restart db

# データベース接続確認
docker-compose exec db psql -U postgres -d corporation_db -c "SELECT COUNT(*) FROM corporations;"
```

#### インポートエラー
```bash
# インポートログ確認
docker-compose logs import

# dry-runで事前確認
docker-compose run --rm import go run ./cmd/import -input /data/gbiz_latest.zip -dry-run

# データベースの現在の状態確認
docker-compose exec db psql -U postgres -d corporation_db -c "SELECT COUNT(*) FROM corporations WHERE prefecture_code IS NOT NULL;"
```

#### APIサーバーエラー
```bash
# APIサーバーログ確認
docker-compose logs app

# ヘルスチェック
curl http://localhost:8080/health

# APIサーバー再起動
docker-compose restart app
```

---

### 📝 定期実行用自動化スクリプト

#### データ更新スクリプト例
```bash
#!/bin/bash
# update_gbiz_data.sh

set -e

# 現在日時のファイル名
FILENAME="gbiz_$(date +%Y%m%d_%H%M%S).zip"

echo "Starting gBizINFO data update..."

# 1. データベース起動確認
docker-compose up -d db
sleep 10

# 2. 新しいデータをダウンロード
echo "Downloading latest data..."
docker-compose run --rm download -output "/data/$FILENAME"

# 3. データをインポート
echo "Importing data..."
docker-compose run --rm import go run ./cmd/import -input "/data/$FILENAME"

# 4. APIサーバー起動
echo "Starting API server..."
docker-compose up -d app

# 5. 動作確認
sleep 5
curl -f http://localhost:8080/health

echo "Data update completed successfully!"
```

#### 使用方法
```bash
# スクリプトに実行権限を付与
chmod +x update_gbiz_data.sh

# 実行
./update_gbiz_data.sh
```

---

### 💾 データ永続化について

現在の設定では、PostgreSQLデータベースのデータは **永続化されています**。

- **ボリューム**: `postgres_data` 名前付きボリュームでデータ保存
- **再起動後も保持**: Docker/Docker Composeを再起動してもデータは削除されません
- **完全削除**: `docker-compose down -v` でボリュームも削除される

```bash
# データ保持して再起動
docker-compose down && docker-compose up -d

# データも含めて完全削除（注意！）
docker-compose down -v
```

## 開発状況

### ✅ 完了済み
1. **Atlas → sqldef マイグレーション**: 宣言的スキーマ管理への移行完了
2. **OpenAPI + oapi-codegen 実装**: 自動コード生成とChi統合完了
3. **Redocly分割ファイル管理**: モジュラーOpenAPI仕様管理完了
4. **データベース統合**: 実際のPostgreSQLクエリによるAPI実装完了
5. **Docker実行環境修正**: Docker download実行時の引数渡し問題とファイルシステム問題の修正完了

### 🔧 最近の修正 (2025年5月31日)

#### Docker Download実行時の問題修正

**問題**: 
- Makefileの`docker-compose run`コマンドでdownloadバイナリに引数が正しく渡されない
- Dockerボリュームマウント環境でのファイル移動（`os.Rename`）がクロスファイルシステムエラーで失敗

**修正内容**:

1. **Makefile修正**: Docker引数渡し問題の解決
   ```makefile
   # 修正前（引数が正しく渡されない）
   docker-compose run --rm download -output /data/gbiz_$(shell date +%Y%m%d_%H%M%S).zip

   # 修正後（バイナリパスを明示的に指定）
   docker-compose run --rm download ./download -output /data/gbiz_$(shell date +%Y%m%d_%H%M%S).zip
   ```

2. **ファイルシステム修正**: `cmd/download/main.go`
   ```go
   // 修正前: クロスファイルシステムでエラー
   err = os.Rename(zipPath, *outputPath)

   // 修正後: copyFile関数による安全なファイル移動
   err = copyFile(zipPath, *outputPath)
   
   // copyFile関数を追加
   func copyFile(src, dst string) error {
       // io.Copy使用による安全なファイルコピー実装
   }
   ```

**修正結果**:
- ✅ Docker環境でのダウンロードコマンドが正常動作
- ✅ ファイル移動処理がクロスファイルシステム環境で安定動作
- ✅ データ取得からインポートまでの完全なパイプラインが機能

### 🚧 進行中/次のステップ
1. **APIドキュメント プレビューサーバー**: Redoclyプレビューの設定
2. **エラーハンドリングの改善**: より詳細なエラーレスポンス
3. **バリデーション強化**: リクエストデータの検証
4. **テストの追加**: ユニットテスト・統合テストの実装
5. **都道府県フィルタリングAPI**: 法人情報の効率的な地域別検索 ✅

### 📋 実装完了機能

#### gBizINFOバッチシステム
- **データダウンロード**: gBizINFO基本情報CSV自動取得 ✅
- **CSVパース**: ZIP展開・CSV解析・エラーハンドリング ✅
- **大量データ処理**: バッチサイズ制御・プログレス表示 ✅
- **UPSERTサポート**: 既存データ更新・新規追加対応 ✅
- **バッチコマンド**: スタンドアロン実行・Docker対応 ✅
- **都道府県コード自動抽出**: 所在地から JIS X 0401 準拠のコード生成 ✅

#### API エンドポイント（実データベース統合済み）
- `GET /health` - ヘルスチェック ✅
- `GET /users` - ユーザー一覧取得 ✅
- `POST /users` - ユーザー作成 ✅
- `GET /users/{id}` - ユーザー詳細取得 ✅
- `GET /corporations` - 法人情報一覧取得（都道府県フィルタリング対応） ✅

#### データベース機能
- **スキーマ管理**: sqldefによる宣言的管理 ✅
- **クエリ生成**: SQLCによる型安全なクエリ ✅
- **フィールドサポート**: name, email, phone, address ✅
- **NULLable対応**: phone, addressのオプション扱い ✅
- **都道府県コード**: prefecture_code (JIS X 0401準拠) による高速フィルタリング ✅

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

## gBizINFOバッチシステム

このプロジェクトには、gBizINFO（政府法人情報）から基本法人情報をダウンロード・インポートするバッチシステムが含まれています。

### バッチ操作

```bash
# バッチコマンドをビルド
make batch-build

# バッチインポートを実行（データベースが必要）
make batch-run

# ドライラン（実際のインポートなし）
make batch-dry-run

# Dockerでバッチ実行
make batch-docker
```

### バッチコマンドオプション

```bash
# ヘルプを表示
./bin/batch -help

# カスタムデータベース設定
./bin/batch -db-host mydb.com -db-port 5433 -db-user admin

# 環境変数での設定
DB_HOST=production.db.com DB_USER=admin ./bin/batch
```

### バッチ処理内容

1. **ダウンロード**: gBizINFOから最新の基本法人情報ZIPファイルを取得
2. **解析**: ZIP展開・CSV解析・データ変換
3. **インポート**: 大量データの効率的なUPSERT処理
4. **進捗表示**: リアルタイム進捗・統計情報・エラーハンドリング

### 法人データベース

バッチで取得される法人情報：

- **基本情報**: 法人番号（13桁）、法人名、所在地
- **詳細情報**: 代表者、資本金、従業員数、業種
- **連絡先**: 電話番号、メールアドレス、ウェブサイト
- **メタデータ**: 設立年月日、法人状態、最終更新日
- **地域情報**: 都道府県コード (JIS X 0401準拠、01-47)による高速検索対応

## エンドポイント

現在実装されているAPI：

- `GET /health` - ヘルスチェック
- `GET /users` - ユーザー一覧取得  
- `POST /users` - ユーザー作成
- `GET /users/{id}` - ユーザー詳細取得
- `GET /corporations` - 法人情報一覧取得（都道府県フィルタリング対応）

### 法人情報 API の使用例

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

### 都道府県コード一覧 (JIS X 0401準拠)

| コード | 都道府県 | コード | 都道府県 | コード | 都道府県 | コード | 都道府県 |
|-------|----------|-------|----------|-------|----------|-------|----------|
| 01 | 北海道 | 13 | 東京都 | 25 | 滋賀県 | 37 | 香川県 |
| 02 | 青森県 | 14 | 神奈川県 | 26 | 京都府 | 38 | 愛媛県 |
| 03 | 岩手県 | 15 | 新潟県 | 27 | 大阪府 | 39 | 高知県 |
| 04 | 宮城県 | 16 | 富山県 | 28 | 兵庫県 | 40 | 福岡県 |
| 05 | 秋田県 | 17 | 石川県 | 29 | 奈良県 | 41 | 佐賀県 |
| 06 | 山形県 | 18 | 福井県 | 30 | 和歌山県 | 42 | 長崎県 |
| 07 | 福島県 | 19 | 山梨県 | 31 | 鳥取県 | 43 | 熊本県 |
| 08 | 茨城県 | 20 | 長野県 | 32 | 島根県 | 44 | 大分県 |
| 09 | 栃木県 | 21 | 岐阜県 | 33 | 岡山県 | 45 | 宮崎県 |
| 10 | 群馬県 | 22 | 静岡県 | 34 | 広島県 | 46 | 鹿児島県 |
| 11 | 埼玉県 | 23 | 愛知県 | 35 | 山口県 | 47 | 沖縄県 |
| 12 | 千葉県 | 24 | 三重県 | 36 | 徳島県 | | |

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

## 開発ワークフロー

### コード生成とツール

このプロジェクトではDockerベースのツールチェーンを使用しています。

```bash
# SQLCでデータベースコードを生成
make sqlc-generate

# OpenAPIからAPIコードを生成
make generate-api
```

## 環境設定

### Docker Compose環境変数

アプリケーションの設定は`docker-compose.yml`で管理されています：

```yaml
environment:
  - DB_HOST=db
  - DB_PORT=5432
  - DB_USER=postgres
  - DB_PASSWORD=password
  - DB_NAME=corporation_db
  - PORT=8080
```

カスタマイズが必要な場合は、`.env`ファイルを作成するか`docker-compose.override.yml`を使用してください。

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
