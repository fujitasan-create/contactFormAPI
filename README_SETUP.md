# セットアップガイド

## 前提条件

- Go 1.21以上
- PostgreSQL 12以上

## セットアップ手順

### 1. 依存関係のインストール

```bash
go mod download
```

### 2. データベースのセットアップ

PostgreSQLデータベースを作成します：

```bash
createdb contactform
```

または、psqlで：

```sql
CREATE DATABASE contactform;
```

### 3. マイグレーションの実行

マイグレーションファイルを実行してテーブルを作成します：

```bash
psql -d contactform -f migrations/001_create_contacts_table.up.sql
```

### 4. 環境変数の設定

必要な環境変数を設定します：

```bash
# データベース接続URL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/contactform?sslmode=disable"

# JWT秘密鍵（本番環境では強力なランダム文字列を使用）
export JWT_SECRET="your-secret-key-change-in-production"

# Adminユーザー名
export ADMIN_USERNAME="admin"

# Adminパスワードハッシュ（bcrypt）
# パスワードハッシュを生成するには：
go run scripts/generate_password_hash.go your_password
# 出力されたハッシュをコピーして使用
export ADMIN_PASSWORD_HASH="$2a$10$..."
```

### 5. アプリケーションの起動

```bash
go run cmd/api/main.go
```

サーバーは `http://localhost:8080` で起動します。

## APIエンドポイント

### 公開API

- `GET /health` - ヘルスチェック
- `POST /contact` - 問い合わせを登録

### 管理API

- `POST /admin/login` - Adminログイン（JWTトークン取得）
- `GET /admin/messages` - 問い合わせ一覧取得（JWT認証必須）

### Swagger UI

- `GET /swagger/index.html` - APIドキュメント

## 使用例

### 1. 問い合わせの送信

```bash
curl -X POST http://localhost:8080/contact \
  -H "Content-Type: application/json" \
  -d '{
    "contact": "example@example.com",
    "name": "山田太郎",
    "message": "お問い合わせ内容"
  }'
```

### 2. Adminログイン

```bash
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_password"
  }'
```

レスポンスから `access_token` を取得します。

### 3. 問い合わせ一覧の取得

```bash
curl -X GET http://localhost:8080/admin/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Swaggerドキュメントの生成

Swaggerドキュメントを再生成する場合：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go
```

