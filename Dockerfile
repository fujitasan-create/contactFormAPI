# ビルドステージ
FROM golang:1.21-alpine AS builder

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係ファイルをコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Swaggerドキュメントを生成（必要に応じて）
# RUN go install github.com/swaggo/swag/cmd/swag@latest
# RUN swag init -g cmd/api/main.go

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 実行ステージ
FROM alpine:latest

# セキュリティのため、非rootユーザーを作成
# wgetをヘルスチェック用にインストール
# postgresql-clientをマイグレーション用にインストール
RUN apk --no-cache add ca-certificates wget postgresql-client && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 非rootユーザーに切り替え
USER appuser

# ポート8080を公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションを起動
CMD ["./main"]

