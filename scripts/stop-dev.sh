#!/bin/bash

# 開発サーバー停止・クリーンアップスクリプト

set -e

echo "🛑 Corporation DB 開発環境を停止しています..."

# Docker Compose停止
docker-compose -f docker-compose.full.yml down

# 必要に応じてボリュームもクリーンアップ
if [ "$1" = "--clean" ]; then
    echo "🧹 データとイメージをクリーンアップしています..."
    docker-compose -f docker-compose.full.yml down -v
    docker system prune -f
fi

echo "✅ 停止完了"
