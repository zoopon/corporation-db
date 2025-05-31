# AI実装者ガイドライン

## 概要

本ガイドラインは、AI（ChatGPT、Claude、GitHub Copilot等）が**corporatioin-db**プロジェクトで実装作業を行う際に従うべき原則とベストプラクティスをまとめたものです。

## 目次

1. [実装前の必須確認事項](#実装前の必須確認事項)
2. [プロジェクト理解の原則](#プロジェクト理解の原則)
3. [実装アプローチ](#実装アプローチ)
4. [コード生成・修正の原則](#コード生成修正の原則)
5. [品質保証チェックリスト](#品質保証チェックリスト)
6. [エラー対応プロトコル](#エラー対応プロトコル)
7. [コミュニケーション原則](#コミュニケーション原則)

---

## 実装前の必須確認事項

### 1. プロジェクト構造の把握

**必須アクション:**
```bash
# プロジェクト全体像を理解する
list_dir(/Users/zoo/projects/corporatioin-db)
read_file(README.md)
read_file(Makefile)
```

**確認すべき構造:**
```
corporatioin-db/
├── cmd/api/main.go          # エントリーポイント
├── internal/
│   ├── api/                 # OpenAPI生成コード（触らない）
│   ├── domain/              # ビジネスロジック・エンティティ
│   ├── infrastructure/      # DB、外部依存
│   ├── presentation/        # HTTP层・ルーティング
│   └── usecase/             # アプリケーション層
├── db/
│   ├── schema.sql           # DB スキーマ（唯一のソース）
│   ├── queries/             # SQLC用クエリ
│   └── migrations/          # マイグレーション
├── api/openapi.yaml         # API仕様（唯一のソース）
└── docker-compose.yml      # 開発環境
```

### 2. 技術スタックの理解

**必須知識:**
- **Go 1.21+**: メイン言語
- **PostgreSQL**: データベース
- **Clean Architecture**: 設計パターン
- **Docker-only**: ビルド・実行戦略
- **自動生成重視**: 手動コード最小化

**コード生成ツール:**
- **SQLC**: SQL → Go コード生成
- **oapi-codegen**: OpenAPI → Go コード生成
- **sqldef**: スキーママイグレーション

### 3. 開発フロー確認

**標準フロー:**
1. **仕様定義**: OpenAPI + DB schema
2. **コード生成**: `make generate`
3. **実装**: ビジネスロジック
4. **テスト**: `make test` または `make docker-test`
5. **動作確認**: `make docker-run`

---

## プロジェクト理解の原則

### 1. ファイル読取順序

**Phase 1: 設定・仕様理解**
1. `README.md` - プロジェクト概要
2. `Makefile` - 利用可能なコマンド
3. `docker-compose.yml` - 環境構成
4. `api/openapi.yaml` - API仕様
5. `db/schema.sql` - データベース設計

**Phase 2: 実装状況確認**
1. `cmd/api/main.go` - エントリーポイント
2. `internal/domain/` - ドメインモデル
3. `internal/infrastructure/` - データアクセス
4. `internal/usecase/` - ビジネスロジック
5. `internal/presentation/` - HTTP层

**Phase 3: 生成コード確認**
1. `internal/api/generated.go` - API実装
2. `internal/infrastructure/db/` - SQLC生成コード

### 2. 依存関係マップ

```
presentation → usecase → domain ← infrastructure
     ↓            ↓                      ↑
   HTTP       Business                Database
   Layer      Logic                   Access
```

**重要:** Clean Architectureの依存方向を守る

---

## 実装アプローチ

### 1. タスク分解原則

**大きなタスクの分解例:**
```
「ユーザー管理機能追加」
├── 1. DB設計（schema.sql更新）
├── 2. API設計（openapi.yaml更新）
├── 3. コード生成（make generate）
├── 4. ドメインモデル実装
├── 5. リポジトリ実装
├── 6. ユースケース実装
├── 7. ハンドラー実装
└── 8. 統合テスト
```

### 2. 実装順序

**必須順序:**
1. **設計変更** → 仕様ファイル更新
2. **コード生成** → `make generate`実行
3. **コンパイル確認** → `go build`成功確認
4. **ビジネスロジック** → domain, usecase実装
5. **統合** → presentation層実装
6. **テスト** → 動作確認

### 3. 変更影響範囲の確認

**変更前チェック:**
```bash
# 現在のエラー状況確認
get_errors([target_files])

# 既存機能の動作確認
run_in_terminal("make docker-test")

# 依存関係確認
list_code_usages("function_name")
```

---

## コード生成・修正の原則

### 1. 自動生成ファイルの扱い

**絶対に手動編集禁止:**
- `internal/api/generated.go`
- `internal/infrastructure/db/models.go`
- `internal/infrastructure/db/users.sql.go`
- `internal/infrastructure/db/db.go`

**編集する場合:**
1. ソースファイル（openapi.yaml, users.sql）を修正
2. `make generate`で再生成
3. 手動編集は一切行わない

### 2. SQLC操作原則

**問題発生時の対応:**
```bash
# 1. 生成ファイル完全削除
rm -rf internal/infrastructure/db/*.sql.go

# 2. クリーン再生成
sqlc generate

# 3. エラー確認
get_errors([generated_files])
```

**重複定義エラー対応:**
- キャッシュクリア: `go clean -cache`
- ファイル削除 → 再生成
- 手動編集の痕跡を完全除去

### 3. ファイル編集のベストプラクティス

**insert_edit_into_file使用時:**
```go
// 良い例
type User struct {
    // ...existing code...
    Phone   sql.NullString `json:"phone"`
    Address sql.NullString `json:"address"`
    // ...existing code...
}

// 悪い例 - 既存コードを重複記述
type User struct {
    ID        int32          `json:"id"`
    Name      string         `json:"name"`
    Email     string         `json:"email"`
    Phone     sql.NullString `json:"phone"`
    Address   sql.NullString `json:"address"`
    CreatedAt sql.NullTime   `json:"created_at"`
    UpdatedAt sql.NullTime   `json:"updated_at"`
}
```

**replace_string_in_file使用時:**
- 5行程度のコンテキストを含める
- 置換対象の一意性を保証
- 変更前後でファイル構造を破綻させない

---

## 品質保証チェックリスト

### 1. 実装完了前の必須チェック

**コンパイル・ビルド:**
- [ ] `go build cmd/api/main.go` 成功
- [ ] `make docker-build` 成功
- [ ] `get_errors()` でエラーなし

**機能性:**
- [ ] 新機能のAPI仕様準拠
- [ ] 既存機能への副作用なし
- [ ] データベース操作の正確性

**Clean Architecture準拠:**
- [ ] 依存方向の正確性
- [ ] 層間の適切な責務分離
- [ ] インターフェースを通じた依存注入

### 2. テスト実行

**必須テスト:**
```bash
# 単体テスト
go test ./...

# 統合テスト（Docker環境）
make docker-test

# API動作確認
make docker-run
curl http://localhost:8080/health
```

### 3. セキュリティチェック

**基本確認:**
- [ ] SQLインジェクション対策（SQLC使用）
- [ ] 入力バリデーション実装
- [ ] エラー情報の適切なマスキング
- [ ] CORS設定の確認

---

## エラー対応プロトコル

### 1. エラー分類と対応

**Type A: コンパイルエラー**
```
対応手順:
1. get_errors()で詳細確認
2. 文法・型エラーの修正
3. import文の調整
4. 再コンパイル確認
```

**Type B: SQLC生成エラー**
```
対応手順:
1. 生成ファイル完全削除
2. SQLクエリ構文確認
3. スキーマとの整合性確認
4. sqlc generate実行
```

**Type C: Docker関連エラー**
```
対応手順:
1. コンテナログ確認
2. docker-compose.yml検証
3. 環境変数設定確認
4. ネットワーク・ボリューム確認
```

### 2. エラー報告形式

**人間への報告時:**
```markdown
## エラー概要
[簡潔な説明]

## 発生状況
- ファイル: [対象ファイル]
- 操作: [実行していた操作]
- 環境: [Docker/ローカル]

## エラー詳細
```
[完全なエラーメッセージ]
```

## 実施した対応
1. [試行した解決策1]
2. [試行した解決策2]

## 追加調査が必要な点
[調査すべき箇所]
```

### 3. デバッグ戦略

**段階的アプローチ:**
1. **最小再現**: 最小コードでエラー再現
2. **分離確認**: 個別コンポーネントの動作確認
3. **統合テスト**: 段階的な統合
4. **環境確認**: Docker環境での最終確認

---

## コミュニケーション原則

### 1. 進捗報告

**実装開始時:**
```markdown
## 作業開始
- タスク: [具体的なタスク]
- アプローチ: [採用する手法]
- 影響範囲: [変更予定ファイル]
- 所要時間見積: [XX分程度]
```

**実装中:**
```markdown
## 進捗状況
- 完了: [完了した項目]
- 作業中: [現在の作業]
- 問題: [発生した問題]
- 次のステップ: [予定している作業]
```

**完了時:**
```markdown
## 実装完了
- 変更ファイル: [変更したファイル]
- テスト結果: [テスト実行結果]
- 動作確認: [確認した機能]
- 注意事項: [運用上の注意点]
```

### 2. 質問・確認事項

**技術的判断が必要な場合:**
```markdown
## 実装方針の確認

**状況:**
[現在の状況説明]

**選択肢:**
A. [選択肢A] - メリット: [...] デメリット: [...]
B. [選択肢B] - メリット: [...] デメリット: [...]

**推奨:**
[AI的推奨案と理由]

**判断をお願いします:**
[人間に決定してもらいたい事項]
```

### 3. 学習・改善提案

**プロジェクト改善案:**
```markdown
## 改善提案

**問題:**
[発見した問題・非効率]

**提案:**
[具体的な改善案]

**期待効果:**
[改善によるメリット]

**実装コスト:**
[必要な作業量]
```

---

## 付録

### A. 重要なファイル・ディレクトリ

**設定ファイル（要理解）:**
- `sqlc.yaml` - SQLC設定
- `oapi-codegen.yaml` - API生成設定
- `docker-compose.yml` - 環境設定
- `Makefile` - コマンド定義

**ソースファイル（編集対象）:**
- `api/openapi.yaml` - API仕様
- `db/schema.sql` - DB設計
- `db/queries/*.sql` - SQLクエリ
- `internal/domain/` - ドメインモデル
- `internal/usecase/` - ビジネスロジック
- `internal/presentation/` - HTTPハンドラー
- `internal/infrastructure/gbiz_client.go` - 外部データ取得・都道府県コード抽出

**生成ファイル（編集禁止）:**
- `internal/api/generated.go`
- `internal/infrastructure/db/`

### B. よく使用するコマンド

```bash
# コード生成
make generate

# テスト実行
make test
make docker-test

# 環境起動
make docker-run

# バッチ処理（gBizINFOデータインポート）
make batch-run

# クリーンアップ
make clean
go clean -cache

# SQLC単体実行
sqlc generate
sqlc vet
```

### C. 都道府県フィルタリング機能の実装パターン

**データベース設計:**
```sql
-- prefecture_codeカラム追加例
ALTER TABLE corporations ADD COLUMN prefecture_code VARCHAR(2);
CREATE INDEX idx_corporations_prefecture_code ON corporations(prefecture_code);
```

**SQL クエリパターン:**
```sql
-- 都道府県フィルタリング対応クエリ例
SELECT * FROM corporations 
WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
  AND ($2 = '' OR prefecture_code = $2)
ORDER BY id LIMIT $3 OFFSET $4;
```

**都道府県コード抽出ロジック:**
```go
// 所在地文字列から都道府県コードを抽出
func extractPrefectureCode(location string) string {
    prefectureMap := map[string]string{
        "北海道": "01", "青森": "02", "岩手": "03", // ...
    }
    
    for prefecture, code := range prefectureMap {
        if strings.Contains(location, prefecture) {
            return code
        }
    }
    return ""
}
```

---

## 実装例: 都道府県フィルタリング機能

### 完全実装例

都道府県コードによるフィルタリング機能の実装手順を示します。

#### 1. データベーススキーマ更新

```sql
-- db/schema.sql
CREATE TABLE IF NOT EXISTS corporations (
    id SERIAL PRIMARY KEY,
    corporate_number VARCHAR(13) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    location TEXT,
    prefecture_code VARCHAR(2), -- 追加
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- インデックス追加（フィルタリング高速化）
CREATE INDEX IF NOT EXISTS idx_corporations_prefecture_code 
ON corporations(prefecture_code);
```

#### 2. OpenAPI仕様更新

```yaml
# api/components/schemas/Corporation.yaml
Corporation:
  type: object
  properties:
    id:
      type: integer
    corporate_number:
      type: string
    name:
      type: string
    location:
      type: string
    prefecture_code:  # 追加
      type: string
      description: "JIS X 0401準拠の都道府県コード (01-47)"
      example: "13"

# api/paths/corporations.yaml
parameters:
  - name: prefecture_code
    in: query
    description: "都道府県コード (JIS X 0401準拠、01-47)"
    schema:
      type: string
      pattern: "^(0[1-9]|[1-4][0-9])$"
    example: "13"
```

#### 3. SQLクエリ更新

```sql
-- db/queries/corporations.sql
-- name: GetCorporationsWithFilter :many
SELECT * FROM corporations
WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
  AND ($2 = '' OR prefecture_code = $2)  -- 追加
ORDER BY id LIMIT $3 OFFSET $4;

-- name: CountCorporationsWithFilter :one
SELECT COUNT(*) FROM corporations
WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
  AND ($2 = '' OR prefecture_code = $2);  -- 追加
```

#### 4. ドメインモデル更新

```go
// internal/domain/corporation.go
type Corporation struct {
    ID               int32     `json:"id"`
    CorporateNumber  string    `json:"corporate_number"`
    Name            string    `json:"name"`
    Location        *string   `json:"location"`
    PrefectureCode  *string   `json:"prefecture_code"`  // 追加
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type CorporationFilter struct {
    Name           string `json:"name"`
    PrefectureCode string `json:"prefecture_code"`  // 追加
    Limit          int32  `json:"limit"`
    Offset         int32  `json:"offset"`
}
```

#### 5. 都道府県コード抽出ロジック

```go
// internal/infrastructure/gbiz_client.go
func extractPrefectureCode(location string) string {
    if location == "" {
        return ""
    }
    
    // 都道府県名マッピング (JIS X 0401準拠)
    prefectureMap := map[string]string{
        "北海道": "01", "青森": "02", "岩手": "03", "宮城": "04",
        "秋田": "05", "山形": "06", "福島": "07", "茨城": "08",
        // ... 全47都道府県
        "沖縄": "47",
    }
    
    // 主要都市による判定（フォールバック）
    cityToPrefecture := map[string]string{
        "札幌": "01", "青森": "02", "盛岡": "03", "仙台": "04",
        // ... 主要都市マッピング
        "那覇": "47",
    }
    
    // 都道府県名での完全一致チェック
    for prefecture, code := range prefectureMap {
        if strings.Contains(location, prefecture) {
            return code
        }
    }
    
    // 都市名でのフォールバック
    for city, code := range cityToPrefecture {
        if strings.Contains(location, city) {
            return code
        }
    }
    
    return ""
}
```

#### 6. APIハンドラー更新

```go
// internal/presentation/corporation_handler.go
func (h *CorporationHandler) GetCorporations(w http.ResponseWriter, r *http.Request) {
    nameFilter := r.URL.Query().Get("name")
    prefectureCode := r.URL.Query().Get("prefecture_code")  // 追加
    
    filter := domain.CorporationFilter{
        Name:           nameFilter,
        PrefectureCode: prefectureCode,  // 追加
        Limit:          parseLimit(r.URL.Query().Get("limit")),
        Offset:         parseOffset(r.URL.Query().Get("offset")),
    }
    
    corporations, err := h.usecase.GetCorporationsWithFilter(r.Context(), filter)
    // ... エラーハンドリング・レスポンス処理
}
```

#### 7. テスト実装

```go
// internal/infrastructure/gbiz_client_test.go
func TestExtractPrefectureCode(t *testing.T) {
    tests := []struct {
        name     string
        location string
        expected string
    }{
        {"Tokyo", "東京都新宿区", "13"},
        {"Hokkaido", "北海道札幌市", "01"},
        {"Okinawa", "沖縄県那覇市", "47"},
        {"Empty", "", ""},
        {"Unknown", "不明な地域", ""},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := extractPrefectureCode(tt.location)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 実装チェックリスト

#### ✅ データベース層
- [ ] スキーマにprefecture_codeカラム追加
- [ ] インデックス追加でフィルタリング高速化
- [ ] SQLクエリにWHERE条件追加

#### ✅ API層
- [ ] OpenAPI仕様にprefecture_codeパラメータ追加
- [ ] レスポンススキーマにprefecture_codeフィールド追加
- [ ] コード生成でAPI層更新

#### ✅ アプリケーション層
- [ ] ドメインモデルにPrefectureCodeフィールド追加
- [ ] フィルタ構造体にprefecture_code追加
- [ ] 都道府県コード抽出ロジック実装

#### ✅ インフラ層
- [ ] リポジトリでprefecture_codeフィルタリング対応
- [ ] CSVインポート時の都道府県コード自動抽出
- [ ] バッチ処理での一括更新対応

#### ✅ プレゼンテーション層
- [ ] HTTPハンドラーでクエリパラメータ受け取り
- [ ] レスポンスにprefecture_code含める
- [ ] エラーハンドリング適切に実装
