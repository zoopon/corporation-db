# ビルドステージ
FROM golang:1.23.5-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git

WORKDIR /app

# Go modules関連のファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# API用のアプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# バッチ用のアプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o batch ./cmd/batch

# API実行用の軽量イメージ
FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# ビルドしたAPIバイナリをコピー
COPY --from=builder /app/api .

# ポート8080を公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./api"]

# バッチ実行用の軽量イメージ
FROM alpine:latest AS batch

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# ビルドしたバッチバイナリをコピー
COPY --from=builder /app/batch .

# バッチアプリケーションを実行
CMD ["./batch"]
