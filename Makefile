.PHONY: help build run test clean docker-up docker-down sqlc-generate

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

# データベースのマイグレーション（開発用）
db-migrate:
	docker-compose exec db psql -U postgres -d corporation_db -f /docker-entrypoint-initdb.d/001_create_users_table.sql
