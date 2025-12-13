# 問い合わせフォームAPI 設計書

## 概要
- 問い合わせフォーム用のREST APIを実装する
- 言語：Go
- Webフレームワーク：Gin
- コンテナ：Docker
- ホスティング：AWS ECS (Fargate)
- 管理画面はUIを作らず、Swagger UI からJSONを確認する運用とする

## 開発方法

**このプロジェクトはDocker Composeのみで開発・実行します。**

```bash
# コンテナの起動
docker-compose up -d

# ログの確認
docker-compose logs -f api

# コンテナの停止
docker-compose down
```

詳細は [README_DOCKER.md](README_DOCKER.md) を参照してください。

**注意**: `go run cmd/api/main.go`は使用しません。すべてDocker Composeで実行します。

## 要件

### 公開API（一般ユーザー向け）
- 問い合わせをPOSTで受け取る
- すべて必須フィールドとする
  - 連絡先（メールアドレス・電話番号など）
    - 型：string
    - フォーマットは厳密に制限しない
  - 名前
    - 型：string
  - メッセージ
    - 型：string

### 管理API（Admin専用）
- Adminのみが問い合わせ一覧を取得できる
- UIは作成せず、Swagger UI 経由でJSONを閲覧する
- JWTによる認証を行う

## API仕様

### Public API

### POST /contact
- 概要：問い合わせを登録する
- Request Body（JSON）
{
  "contact": "example@example.com",
  "name": "山田太郎",
  "message": "お問い合わせ内容"
}

- Validation
  - contact, name, message はすべて必須
  - 空文字は不可
  - message は最大長制限を設ける（例：2000文字）

- Response
  - 201 Created
{
  "status": "created"
}

### GET /health
- 概要：ヘルスチェック用
- 用途：ECS / ALB のヘルスチェック

- Response
  - 200 OK
{
  "status": "ok"
}

## Admin API

### POST /admin/login
- 概要：Adminログイン
- Request Body
{
  "username": "admin",
  "password": "password"
}

- 処理内容
  - username を環境変数と比較
  - password を bcrypt で検証
  - 正常時に JWT を発行

- Response
{
  "access_token": "JWT_TOKEN"
}

### GET /admin/messages
- 概要：問い合わせ一覧取得
- 認証：JWT必須
- Header
Authorization: Bearer <JWT_TOKEN>

- Response
[
  {
    "id": 1,
    "contact": "example@example.com",
    "name": "山田太郎",
    "message": "お問い合わせ内容",
    "created_at": "2025-01-01T12:00:00Z"
  }
]

## 認証設計
- 認証方式：JWT
- Admin情報はDBに持たず、環境変数で管理する

### 使用する環境変数
- ADMIN_USERNAME
- ADMIN_PASSWORD_HASH（bcrypt）
- JWT_SECRET

### JWT仕様
- 有効期限：例 12時間
- 署名方式：HMAC（HS256）

## データベース設計

### contacts テーブル
- id (bigserial, primary key)
- contact (text, not null)
- name (text, not null)
- message (text, not null)
- created_at (timestamptz, not null, default now())
- ip (text, null)
- user_agent (text, null)

## ディレクトリ構成

.
├─ cmd/
│  └─ api/
│     └─ main.go
├─ internal/
│  ├─ config/
│  ├─ http/
│  ├─ auth/
│  ├─ db/
│  └─ repository/
├─ migrations/
├─ docs/
├─ Dockerfile
├─ docker-compose.yml
└─ go.mod

## Swagger
- swaggo/swag
- swaggo/gin-swagger
- アクセス例
  /swagger/index.html

## Docker設計
- マルチステージビルド
- ポート：8080
- 環境変数で設定を注入
- /health エンドポイントを実装

## AWS ECS (Fargate) 設計
- ECR：Dockerイメージ格納
- ECS Fargate：API実行
- ALB：外部公開
- RDS(PostgreSQL)：データ保存
- CloudWatch Logs：ログ管理

### Task Definition
- ポート：8080
- 環境変数
  - DATABASE_URL
  - JWT_SECRET
  - ADMIN_USERNAME
  - ADMIN_PASSWORD_HASH
- ヘルスチェック
  - GET /health

## セキュリティ・運用上の考慮
- POST /contact にレート制限
- 入力文字数制限
- Admin API は JWT 必須
- 機密情報は Secrets Manager 管理を推奨
