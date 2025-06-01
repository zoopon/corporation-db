.PHONY: help docker-run docker-up docker-down docker-rebuild sqlc-generate schema-apply schema-diff schema-dry-run generate-api openapi-lint openapi-bundle openapi-preview schema-export schema-check docs-serve docs-build generate-all db-reset db-reset-data db-status db-connect

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  docker-run     - Start development environment with hot reload"
	@echo "  docker-up      - Start Docker containers (background)"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-rebuild - Force rebuild Docker images (no cache)" docker-run docker-up docker-down docker-dev docker-dev-up docker-dev-down docker-rebuild sqlc-generate schema-apply schema-diff schema-dry-run generate-api openapi-lint openapi-bundle openapi-preview schema-export schema-check docs-serve docs-build generate-all db-reset db-reset-data db-status db-connect

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  docker-run      - Start application with Docker (production)"
	@echo "  docker-up       - Start Docker containers (production)"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  docker-dev      - Start development environment with hot reload"
	@echo "  docker-dev-up   - Start development containers"
	@echo "  docker-dev-down - Stop development containers"
	@echo "  docker-rebuild  - Force rebuild Docker images (no cache)"
	@echo "  docker-rebuild-dev - Force rebuild development Docker images"
	@echo "  generate-all   - Generate all code (API + SQLC)"
	@echo ""
	@echo "Code generation:"
	@echo "  sqlc-generate  - Generate code from SQL"
	@echo "  generate-api   - Generate API code from OpenAPI spec"
	@echo ""
	@echo "OpenAPI/Redocly commands:"
	@echo "  openapi-lint   - Lint OpenAPI specification"
	@echo "  openapi-bundle - Bundle split OpenAPI files into single file"
	@echo "  openapi-preview - Preview API documentation"
	@echo ""
	@echo "sqldef commands:"
	@echo "  schema-apply   - Apply schema changes to database"
	@echo "  schema-diff    - Show schema differences (dry-run)"
	@echo "  schema-dry-run - Dry run schema changes"
	@echo "  schema-export  - Export current database schema"
	@echo "  schema-check   - Check database connection and status"
	@echo ""
	@echo "Database operations:"
	@echo "  db-reset       - Reset database (remove all data and volumes)"
	@echo "  db-reset-data  - Reset only data (keep table structure)"
	@echo "  db-status      - Show database status and record counts"
	@echo "  db-connect     - Connect to database shell"
	@echo ""
	@echo "Documentation:"
	@echo "  docs-serve     - Serve API documentation locally"
	@echo "  docs-build     - Build static documentation"
	@echo ""
	@echo "Batch operations:"
	@echo "  batch-build    - Build batch command"
	@echo "  batch-run      - Run batch import (requires database)"
	@echo "  batch-dry-run  - Run batch import in dry-run mode"
	@echo "  batch-docker   - Run batch import in Docker"
	@echo ""
	@echo "Download and Import operations:"
	@echo "  download-build       - Build download command"
	@echo "  import-build         - Build import command"
	@echo "  build-all            - Build all commands"
	@echo "  download-data        - Download data from gBizINFO (local)"
	@echo "  download-data-docker - Download data from gBizINFO (Docker)"
	@echo "  import-data          - Import data from ZIP file (local)"
	@echo "  import-data-docker   - Import data from ZIP file (Docker)"

# Docker環境でアプリケーションを起動（開発環境・ホットリロード対応）
docker-run:
	docker-compose up --build

# Dockerコンテナを起動（バックグラウンド）
docker-up:
	docker-compose up -d

# Dockerコンテナを停止
docker-down:
	docker-compose down

# Dockerイメージを強制的に再ビルド（キャッシュなし）
docker-rebuild:
	docker-compose build --no-cache

# SQLCでコードを生成
sqlc-generate:
	sqlc generate

# 開発環境でSQLCコードを生成（Docker内で実行）
sqlc-generate-docker:
	docker-compose exec app sqlc generate

# oapi-codegenでAPIコードを生成
generate-api:
	make openapi-bundle
	oapi-codegen --config oapi-codegen.yaml api/openapi-bundled.yaml

# Redocly: OpenAPI仕様をlint
openapi-lint:
	docker-compose --profile docs run --rm redocly lint api/openapi.yaml

# Redocly: 分割されたOpenAPIファイルを単一ファイルにバンドル
openapi-bundle:
	docker-compose --profile docs run --rm redocly bundle api/openapi.yaml --output api/openapi-bundled.yaml

# Redocly: API ドキュメントをプレビュー
openapi-preview:
	docker-compose --profile docs run --rm -p 8081:8080 redocly preview-docs api/openapi.yaml --host 0.0.0.0

# Redocly: API ドキュメントをプレビュー（バックグラウンド）
openapi-preview-bg:
	@echo "Starting API documentation preview at http://localhost:8081"
	docker-compose --profile docs up -d redocly
	docker-compose --profile docs exec redocly redocly preview-docs api/openapi.yaml --host 0.0.0.0 --port 8080 &

# Redocly: プレビューサーバーを停止
openapi-preview-stop:
	docker-compose --profile docs down

# sqldef: スキーマを適用
schema-apply:
	docker-compose --profile migration run --rm -T -e PGPASSWORD=password sqldef \
		-h db -p 5432 -U postgres \
		corporation_db < db/schema.sql

# sqldef: スキーマの差分を確認（dry-run）
schema-diff:
	docker-compose --profile migration run --rm -T -e PGPASSWORD=password sqldef \
		-h db -p 5432 -U postgres --dry-run \
		corporation_db < db/schema.sql

# sqldef: ドライラン（変更内容を表示のみ）
schema-dry-run: schema-diff

# sqldef: 現在のスキーマをエクスポート
schema-export:
	docker-compose --profile migration run --rm -T -e PGPASSWORD=password sqldef \
		-h db -p 5432 -U postgres --export \
		corporation_db

# sqldef: データベース接続確認
schema-check:
	docker-compose exec db psql -U postgres -d corporation_db -c "SELECT current_database(), current_user;"

# 全コード生成（推奨）
generate-all: openapi-bundle generate-api sqlc-generate
	@echo "All code generation completed"

# API documentation commands
docs-serve:
	@echo "Starting API documentation server..."
	docker-compose --profile docs up redocly

docs-build:
	@echo "Building static API documentation..."
	docker-compose --profile docs run --rm redocly build-docs api/openapi.yaml --output docs/

# Batch operation commands
batch-build:
	@echo "Building batch command..."
	go build -o bin/batch cmd/batch/main.go

batch-run: batch-build
	@echo "Running batch import..."
	./bin/batch

batch-dry-run: batch-build
	@echo "Running batch import (dry-run)..."
	./bin/batch -dry-run

batch-docker:
	@echo "Running batch import in Docker..."
	docker-compose --profile batch up --build batch

# Download and Import commands
download-build:
	@echo "Building download command..."
	go build -o bin/download cmd/download/main.go

import-build:
	@echo "Building import command..."
	go build -o bin/import cmd/import/main.go

build-all: batch-build download-build import-build
	@echo "All commands built successfully"

download-data: download-build
	@echo "Downloading data from gBizINFO..."
	./bin/download -output ./data/gbiz_$(shell date +%Y%m%d_%H%M%S).zip

download-data-docker:
	@echo "Downloading data from gBizINFO using Docker..."
	docker-compose run --rm download ./download -output /data/gbiz_$(shell date +%Y%m%d_%H%M%S).zip

import-data: import-build
	@echo "Importing data from ZIP file..."
	@read -p "Enter ZIP file path: " zipfile; \
	./bin/import -input "$$zipfile"

import-data-docker:
	@echo "Importing data from ZIP file using Docker..."
	@read -p "Enter ZIP file path (relative to ./data/): " zipfile; \
	docker-compose run --rm import ./import -input "/data/$$zipfile"

# Database operations
db-reset:
	@echo "⚠️  WARNING: This will permanently delete ALL database data!"
	@echo "This action cannot be undone."
	@read -p "Are you sure you want to reset the database? (yes/NO): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "Stopping containers and removing volumes..."; \
		docker-compose down -v; \
		echo "Starting fresh database..."; \
		docker-compose up -d db; \
		echo "Waiting for database to be ready..."; \
		sleep 10; \
		echo "Database reset completed successfully!"; \
	else \
		echo "Database reset cancelled."; \
	fi

db-reset-data:
	@echo "⚠️  WARNING: This will delete all data from corporations and users tables!"
	@echo "Table structure will be preserved."
	@read -p "Are you sure you want to reset table data? (yes/NO): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "Truncating tables..."; \
		docker-compose exec db psql -U postgres -d corporation_db -c "TRUNCATE TABLE corporations, users RESTART IDENTITY CASCADE;"; \
		echo "Data reset completed successfully!"; \
		make db-status; \
	else \
		echo "Data reset cancelled."; \
	fi

db-status:
	@echo "=== Database Status ==="
	@echo "Tables:"
	@docker-compose exec db psql -U postgres -d corporation_db -c "\dt"
	@echo ""
	@echo "Record counts:"
	@docker-compose exec db psql -U postgres -d corporation_db -c "SELECT 'corporations' as table_name, COUNT(*) as record_count FROM corporations UNION ALL SELECT 'users' as table_name, COUNT(*) as record_count FROM users;"
	@echo ""
	@echo "Database size:"
	@docker-compose exec db psql -U postgres -d corporation_db -c "SELECT pg_size_pretty(pg_database_size('corporation_db')) as database_size;"

db-connect:
	@echo "Connecting to database shell..."
	@echo "Use \\q to exit the database shell"
	docker-compose exec db psql -U postgres -d corporation_db
