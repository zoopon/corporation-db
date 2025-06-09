.PHONY: help volume-create volume-status volume-remove volume-setup docker-up docker-down docker-rebuild backend-up frontend-up full-up

# デフォルトターゲット
help:
	@echo "Corporation DB - Docker Compose Management"
	@echo ""
	@echo "Available commands:"
	@echo "  full-up        - Start full application (nginx + backend + frontend + db)"
	@echo "  backend-up     - Start backend development environment only"
	@echo "  frontend-up    - Start frontend development environment only"
	@echo "  docker-up      - Start all services"
	@echo "  docker-down    - Stop all services"
	@echo "  docker-rebuild - Rebuild and restart all services"
	@echo ""
	@echo "Shared Volume operations:"
	@echo "  volume-create  - Create shared PostgreSQL volume"
	@echo "  volume-status  - Show shared volume status"
	@echo "  volume-remove  - Remove shared PostgreSQL volume"
	@echo "  volume-setup   - Setup shared volume configuration"
	@echo ""
	@echo "Database operations:"
	@echo "  db-connect     - Connect to PostgreSQL shell"
	@echo "  db-status      - Show database status"

# === Shared Volume Management ===
# 共有PostgreSQLボリューム名
SHARED_VOLUME_NAME := corporatioin-db-postgres-data

# 共有ボリュームを作成
volume-create:
	@echo "Creating shared PostgreSQL volume: $(SHARED_VOLUME_NAME)"
	@if docker volume inspect $(SHARED_VOLUME_NAME) >/dev/null 2>&1; then \
		echo "Volume $(SHARED_VOLUME_NAME) already exists"; \
	else \
		docker volume create $(SHARED_VOLUME_NAME); \
		echo "Volume $(SHARED_VOLUME_NAME) created successfully"; \
	fi

# 共有ボリュームの状態確認
volume-status:
	@echo "=== Shared Volume Status ==="
	@echo "Volume name: $(SHARED_VOLUME_NAME)"
	@if docker volume inspect $(SHARED_VOLUME_NAME) >/dev/null 2>&1; then \
		echo "Status: EXISTS"; \
		echo "Details:"; \
		docker volume inspect $(SHARED_VOLUME_NAME) --format "  Driver: {{.Driver}}"; \
		docker volume inspect $(SHARED_VOLUME_NAME) --format "  Mountpoint: {{.Mountpoint}}"; \
		docker volume inspect $(SHARED_VOLUME_NAME) --format "  Created: {{.CreatedAt}}"; \
	else \
		echo "Status: NOT EXISTS"; \
		echo "Run 'make volume-create' to create the shared volume"; \
	fi

# 共有ボリュームを削除（警告付き）
volume-remove:
	@echo "WARNING: This will permanently delete all PostgreSQL data!"
	@echo "Volume to remove: $(SHARED_VOLUME_NAME)"
	@read -p "Are you sure you want to continue? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		if docker volume inspect $(SHARED_VOLUME_NAME) >/dev/null 2>&1; then \
			docker volume rm $(SHARED_VOLUME_NAME); \
			echo "Volume $(SHARED_VOLUME_NAME) removed successfully"; \
		else \
			echo "Volume $(SHARED_VOLUME_NAME) does not exist"; \
		fi; \
	else \
		echo "Volume removal cancelled."; \
	fi

# 共有ボリューム設定
volume-setup:
	@echo "Setting up shared volume configuration..."
	@make volume-create
	@echo ""
	@echo "✅ Shared volume setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  make full-up     - Start full application stack"
	@echo "  make backend-up  - Start backend development only"

# === Docker Compose Operations ===

# フルスタックアプリケーション起動
full-up:
	@make volume-create
	@echo "Starting full application stack..."
	docker compose up -d
	@echo ""
	@echo "✅ Application started!"
	@echo "  Frontend: http://localhost"
	@echo "  Backend API: http://localhost/api"
	@echo "  Database: localhost:5432"

# バックエンド開発環境のみ起動
backend-up:
	@make volume-create
	@echo "Starting backend development environment..."
	cd backend && docker compose up -d
	@echo ""
	@echo "✅ Backend development environment started!"
	@echo "  API Server: http://localhost:8080"
	@echo "  Database: localhost:5432"

# フロントエンド開発環境のみ起動（バックエンドが必要）
frontend-up:
	@echo "Starting frontend development environment..."
	@echo "Note: Backend must be running for frontend to work properly"
	docker compose up frontend -d

# 全サービス起動
docker-up:
	@make volume-create
	docker compose up -d

# 全サービス停止
docker-down:
	docker compose down
	cd backend && docker compose down

# 全サービス再構築
docker-rebuild:
	@make volume-create
	docker compose down
	docker compose build --no-cache
	docker compose up -d

# === Database Operations ===

# データベース接続
db-connect:
	@echo "Connecting to PostgreSQL database..."
	@echo "Use \\q to exit the database shell"
	docker compose exec db psql -U postgres -d corporation_db

# データベース状態確認
db-status:
	@echo "=== Database Status ==="
	@if docker compose ps db | grep -q "Up"; then \
		echo "Database Status: Running"; \
		echo ""; \
		echo "Tables:"; \
		docker compose exec db psql -U postgres -d corporation_db -c "\dt"; \
		echo ""; \
		echo "Record counts:"; \
		docker compose exec db psql -U postgres -d corporation_db -c "SELECT 'corporations' as table_name, COUNT(*) as record_count FROM corporations UNION ALL SELECT 'bases' as table_name, COUNT(*) as record_count FROM bases;" 2>/dev/null || echo "Tables not yet created"; \
	else \
		echo "Database Status: Not running"; \
		echo "Run 'make full-up' or 'make backend-up' to start the database"; \
	fi
