# 問い合わせフォームAPI 設計書

## 概要
- 問い合わせフォーム用のREST APIを実装する
- 言語：Go
- Webフレームワーク：Gin
- コンテナ：Docker
- ホスティング：Fly.io
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
- 用途：Fly.io のヘルスチェック

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
├─ scripts/
├─ Dockerfile
├─ docker-compose.yml
├─ fly.toml
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

## Fly.io デプロイ設計

### インフラ構成
- **Fly.io**: アプリケーション実行（Dockerコンテナ）
- **Fly.io PostgreSQL**: データ保存（別アプリとして管理）
- **自動HTTPS**: `force_https = true` で自動的にHTTPS化
- **自動スケーリング**: リクエストに応じて自動起動・停止

### 設定ファイル（fly.toml）
- ポート：8080
- リージョン：nrt（東京）
- メモリ：512MB
- CPU：1コア

### 環境変数（Fly.io Secrets）
以下の環境変数を `fly secrets set` で設定：
- `DATABASE_URL`: Fly.io PostgreSQLの接続URL（自動設定される場合あり）
- `JWT_SECRET`: JWT署名用の秘密鍵
- `ADMIN_USERNAME`: Adminユーザー名
- `ADMIN_PASSWORD_HASH`: Adminパスワードのbcryptハッシュ
- `PORT`: アプリケーションのポート（デフォルト: 8080）

### デプロイ時の自動実行
- **マイグレーション**: `release_command`でデプロイ時に自動実行
  ```toml
  [deploy]
    release_command = "sh -c 'psql \"$DATABASE_URL\" -f migrations/001_create_contacts_table.up.sql || true'"
  ```

### ヘルスチェック
- エンドポイント：`GET /health`
- Fly.ioが自動的にヘルスチェックを実行

### デプロイコマンド
```bash
# デプロイ
fly deploy -a contactformapi

# ログ確認
fly logs -a contactformapi

# 環境変数設定
fly secrets set -a contactformapi JWT_SECRET=your-secret
```

## セキュリティ・運用上の考慮
- POST /contact にレート制限
- 入力文字数制限
- Admin API は JWT 必須
- 機密情報は Fly.io Secrets で管理（`fly secrets set`）
- HTTPS は自動的に有効化（`force_https = true`）

## Fly.io デプロイ手順

### 1. Fly.io CLI のインストール
```bash
# Windows (PowerShell)
iwr https://fly.io/install.ps1 -useb | iex
```

### 2. ログイン
```bash
fly auth login
```

### 3. アプリの作成（初回のみ）
```bash
fly launch
```

### 4. PostgreSQL の作成（別アプリとして）
```bash
fly postgres create --name contactform-db
```

### 5. 環境変数の設定
```bash
# パスワードハッシュを生成
go run scripts/generate_password_hash.go your_password

# 環境変数を設定
fly secrets set -a contactformapi \
  JWT_SECRET=your-secret \
  ADMIN_USERNAME=admin \
  ADMIN_PASSWORD_HASH='$2a$10$...'
```

### 6. デプロイ
```bash
fly deploy -a contactformapi
```

### 7. 動作確認
- Swagger UI: `https://contactformapi.fly.dev/swagger/index.html`
- ヘルスチェック: `https://contactformapi.fly.dev/health`

詳細なトラブルシューティングは [README_FLYIO_TROUBLESHOOTING.md](README_FLYIO_TROUBLESHOOTING.md) を参照してください。
