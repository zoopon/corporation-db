# OpenAI API Integration Implementation Summary

## 実装概要

法人番号を使用してOpenAI APIから拠点/支店情報を取得し、データベースに保存するエンドポイントを実装しました。

## 実装したコンポーネント

### 1. OpenAPI仕様
- **ファイル**: `api/paths/corporations_{corporate_number}_fetch-bases.yaml`
- **エンドポイント**: `POST /corporations/{corporate_number}/fetch-bases`
- **機能**: 法人の拠点情報をOpenAI APIで取得・保存

### 2. ドメイン層
- **ファイル**: `internal/domain/openai.go`
- **定義**: OpenAI API用のリクエスト/レスポンス構造体、サービスインターフェース
- **主要な型**:
  - `OpenAIRequest`, `OpenAIResponse`: API通信用
  - `URLDiscoveryResult`: URL発見結果
  - `BaseExtractionResult`: 拠点情報抽出結果
  - `OpenAIService`: サービスインターフェース

### 3. インフラストラクチャ層
- **ファイル**: `internal/infrastructure/openai_client.go`
- **機能**: OpenAI API クライアント実装（github.com/sashabaranov/go-openaiを使用せず独自実装）
- **主要メソッド**:
  - `DiscoverURLs()`: 企業の公式サイトから拠点情報を含むURLを発見
  - `ExtractBasesFromURL()`: URLから拠点情報を抽出
  - `chatCompletion()`: OpenAI APIとの通信

### 4. ユースケース層
- **ファイル**: `internal/usecase/fetch_bases.go`
- **機能**: OpenAI APIを使用した拠点情報取得のビジネスロジック
- **主要機能**:
  - URL発見 → 拠点情報抽出 → データベース保存
  - 重複拠点の検出・除外
  - エラーハンドリング

### 5. プレゼンテーション層
- **ファイル**: `internal/presentation/corporation_handler.go`
- **機能**: HTTPリクエスト処理
- **メソッド**: `FetchCorporationBases()` - 新しいエンドポイントのハンドラ

### 6. ルーティング更新
- **ファイル**: `internal/presentation/router.go`
- **変更**: oapi-codegenで生成されたハンドラーを使用するよう修正
- **追加**: POSTメソッドをCORSで許可

### 7. メイン関数更新
- **ファイル**: `cmd/api/main.go`
- **追加**: OpenAIクライアントとFetchBasesUseCaseの初期化

## 環境変数

新しく追加された環境変数：
```
OPENAI_API_KEY=your_openai_api_key_here
```

## API仕様

### エンドポイント
```
POST /corporations/{corporate_number}/fetch-bases
```

### パラメータ
- `corporate_number`: 13桁の法人番号

### レスポンス例
```json
{
  "message": "Successfully fetched and saved 3 base/branch offices",
  "bases_count": 3,
  "urls_found": [
    "https://example.com/company/offices",
    "https://example.com/about/locations"
  ]
}
```

### エラーレスポンス
- `400`: 不正な法人番号形式
- `404`: 法人が見つからない
- `500`: OpenAI API エラーまたは内部エラー

## 動作確認

✅ **エンドポイント認識**: POST /corporations/{corporate_number}/fetch-bases が正常に認識される
✅ **法人検索**: 指定された法人番号で法人情報を取得
✅ **OpenAI API呼び出し**: DiscoverURLs メソッドが実行される
⚠️ **API Key**: 実際のOpenAI API Keyが必要（現在は500エラーで期待通りの動作）

## テスト例
```bash
curl -X POST -i "http://localhost:8080/corporations/8030001013271/fetch-bases"
```

現在の応答（OPENAI_API_KEYが未設定のため）:
```
HTTP/1.1 500 Internal Server Error
{"error":"Failed to fetch bases information"}
```

## 次のステップ

1. **実際のOpenAI API Key設定**: `.env`ファイルに有効なAPIキーを設定
2. **プロンプト調整**: URL発見と拠点情報抽出のプロンプトを実際のケースでテスト・調整
3. **エラーハンドリング改善**: より詳細なエラーメッセージとログ
4. **レート制限対応**: OpenAI APIのレート制限を考慮した実装
5. **テストケース追加**: ユニットテストとインテグレーションテスト

## アーキテクチャ遵守

✅ **Clean Architecture**: ドメイン → ユースケース → インフラ → プレゼンテーションの層分離
✅ **依存性注入**: インターフェースによる疎結合
✅ **OpenAPI Code Generation**: oapi-codegenによる自動生成コード使用
✅ **既存パターン踏襲**: 既存のリポジトリパターンとエラーハンドリングに準拠
