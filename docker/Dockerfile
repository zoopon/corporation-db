FROM golang:1.23.5-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git

WORKDIR /app

# Go modules関連のファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 実行用の軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# ポート8080を公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]
