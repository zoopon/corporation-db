.PHONY: help docker-run docker-up docker-down sqlc-generate schema-apply schema-diff schema-dry-run generate-api openapi-lint openapi-bundle openapi-preview schema-export schema-check docs-serve docs-build generate-all

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  docker-run     - Start application with Docker (recommended)"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
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
	@echo "  download-build - Build download command"
	@echo "  import-build   - Build import command"
	@echo "  build-all      - Build all commands"
	@echo "  download-data  - Download data from gBizINFO"
	@echo "  import-data    - Import data from ZIP file"

# Docker環境でアプリケーションを起動（推奨）
docker-run:
	docker-compose up --build

# Dockerコンテナを起動（バックグラウンド）
docker-up:
	docker-compose up -d

# Dockerコンテナを停止
docker-down:
	docker-compose down

# SQLCでコードを生成
sqlc-generate:
	sqlc generate

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

import-data: import-build
	@echo "Importing data from ZIP file..."
	@read -p "Enter ZIP file path: " zipfile; \
	./bin/import -input "$$zipfile"
