.PHONY: help build run test clean docker-up docker-down sqlc-generate schema-apply schema-diff schema-dry-run generate-api schema-export schema-check

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  sqlc-generate  - Generate code from SQL"
	@echo "  generate-api   - Generate API code from OpenAPI spec"
	@echo ""
	@echo "sqldef commands:"
	@echo "  schema-apply   - Apply schema changes to database"
	@echo "  schema-diff    - Show schema differences (dry-run)"
	@echo "  schema-dry-run - Dry run schema changes"
	@echo "  schema-export  - Export current database schema"
	@echo "  schema-check   - Check database connection and status"

# アプリケーションをビルド
build:
	go build -o bin/api ./cmd/api

# アプリケーションをローカルで実行
run:
	go run ./cmd/api

# テストを実行
test:
	go test -v ./...

# ビルド成果物をクリーンアップ
clean:
	rm -rf bin/

# Dockerコンテナを起動
docker-up:
	docker-compose up -d

# Dockerコンテナを停止
docker-down:
	docker-compose down

# Docker環境でアプリケーションを起動（DBも含む）
docker-run:
	docker-compose up --build

# SQLCでコードを生成
sqlc-generate:
	sqlc generate

# oapi-codegenでAPIコードを生成
generate-api:
	oapi-codegen --config oapi-codegen.yaml api/openapi.yaml

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
