#!/bin/bash

# フルスタック開発サーバー起動スクリプト

set -e

echo "🚀 Corporation DB 開発サーバーを起動しています..."

# 環境変数チェック
if [ ! -f "backend/.env" ]; then
    echo "❌ backend/.env ファイルが見つかりません"
    echo "💡 ./scripts/setup-dev.sh を先に実行してください"
    exit 1
fi

# OpenAI API Key チェック
if ! grep -q "OPENAI_API_KEY=" backend/.env || grep -q "OPENAI_API_KEY=$" backend/.env; then
    echo "⚠️  OPENAI_API_KEYが設定されていません"
    echo "💡 backend/.env ファイルを確認してください"
fi

# Docker Compose起動
echo "🐳 フルスタック環境を起動しています..."
docker-compose -f docker-compose.full.yml up --build

echo "✅ 開発サーバーが起動しました!"
echo ""
echo "🌐 アクセス先:"
echo "  フロントエンド: http://localhost:3000"
echo "  バックエンドAPI: http://localhost:8080"
echo ""
echo "🛑 停止するには Ctrl+C を押してください"
