#!/bin/bash

# フルスタック開発環境セットアップスクリプト

set -e

echo "🚀 Corporation DB フルスタック開発環境セットアップ"

# 現在のディレクトリを確認
if [ ! -f "docker-compose.full.yml" ]; then
    echo "❌ エラー: プロジェクトルートディレクトリで実行してください"
    exit 1
fi

# 環境変数チェック
if [ ! -f "backend/.env" ]; then
    echo "📝 バックエンド環境変数ファイルを作成しています..."
    cp backend/.env.example backend/.env
    echo "⚠️  backend/.env ファイルにOPENAI_API_KEYを設定してください"
fi

# フロントエンド依存関係インストール
echo "📦 フロントエンド依存関係をインストールしています..."
cd frontend
npm install
cd ..

# OpenAPI型生成
echo "🔧 OpenAPI型を生成しています..."
cd frontend
npm run generate-types
cd ..

# Dockerイメージビルド
echo "🐳 Dockerイメージをビルドしています..."
docker-compose -f docker-compose.full.yml build

echo "✅ セットアップ完了!"
echo ""
echo "🎯 次のステップ:"
echo "1. backend/.env ファイルでOPENAI_API_KEYを設定"
echo "2. 開発サーバー起動: ./scripts/start-dev.sh"
echo "3. ブラウザで http://localhost:3000 を開く"
