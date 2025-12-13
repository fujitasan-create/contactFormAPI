# Fly.io トラブルシューティングガイド

## 現在のエラー状況

- ✅ ヘルスチェック (`/health`) - 成功
- ❌ `POST /contact` - 500エラー
- ❌ `POST /admin/login` - 401エラー

## 1. ログの確認

まず、Fly.ioのログを確認して実際のエラーを特定します：

```bash
# リアルタイムでログを確認
fly logs

# 最新の100行を表示
fly logs -n 100
```

## 2. 500エラーの原因と解決方法

### 原因の可能性
1. **データベーステーブルが存在しない**（マイグレーション未実行）
2. **データベース接続エラー**

### 解決手順

#### ステップ1: データベースの状態を確認

```bash
# Fly.ioのPostgreSQLに接続
fly postgres connect -a <your-postgres-app-name>

# または、DATABASE_URLを使って直接接続
fly ssh console -a contactformapi
# その後、psqlコマンドで接続
```

#### ステップ2: テーブルの存在確認

PostgreSQLに接続後、以下を実行：

```sql
-- テーブル一覧を確認
\dt

-- contactsテーブルが存在するか確認
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'contacts'
);
```

#### ステップ3: マイグレーションの実行

テーブルが存在しない場合、マイグレーションを実行します：

**方法1: Fly.ioのSSHコンソールから実行**

```bash
# SSHコンソールに接続
fly ssh console -a contactformapi

# マイグレーションファイルを確認
cat migrations/001_create_contacts_table.up.sql

# psqlでマイグレーションを実行（DATABASE_URL環境変数を使用）
psql $DATABASE_URL -f migrations/001_create_contacts_table.up.sql
```

**方法2: ローカルから実行（推奨）**

```bash
# DATABASE_URLを取得
fly secrets list -a contactformapi | grep DATABASE_URL

# または、Fly.ioのPostgreSQLアプリから取得
fly postgres connect -a <your-postgres-app-name>

# ローカルからマイグレーションを実行
psql <DATABASE_URL> -f migrations/001_create_contacts_table.up.sql
```

**方法3: fly.tomlにinitコマンドを追加（自動実行）**

`fly.toml`に以下を追加：

```toml
[deploy]
  release_command = "psql $DATABASE_URL -f migrations/001_create_contacts_table.up.sql || true"
```

## 3. 401エラーの原因と解決方法

### 原因の可能性
1. **ADMIN_PASSWORD_HASHが正しく設定されていない**
2. **$文字のエスケープ問題**
3. **ADMIN_USERNAMEが一致しない**

### 解決手順

#### ステップ1: 環境変数の確認

```bash
# 現在の環境変数を確認
fly secrets list -a contactformapi
```

#### ステップ2: パスワードハッシュの再生成

ローカルで正しいパスワードハッシュを生成：

```bash
go run scripts/generate_password_hash.go okayu256
```

出力例：
```
Password hash: $2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

#### ステップ3: 環境変数の設定

**重要**: Fly.ioでは`$`文字をエスケープする必要はありません。そのまま設定してください。

```bash
# 環境変数を設定
fly secrets set -a contactformapi \
  JWT_SECRET=okayu256 \
  ADMIN_USERNAME=admin \
  ADMIN_PASSWORD_HASH='$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
```

**注意**: `ADMIN_PASSWORD_HASH`はシングルクォートで囲むことを推奨します。

#### ステップ4: アプリの再起動

環境変数を変更した後、アプリを再起動：

```bash
fly apps restart -a contactformapi
```

## 4. データベース接続の確認

### DATABASE_URLの確認

```bash
# DATABASE_URLが設定されているか確認
fly secrets list -a contactformapi | grep DATABASE_URL

# もし設定されていない場合、Fly.ioのPostgreSQLアプリから取得
fly postgres connect -a <your-postgres-app-name>
# または
fly postgres db list -a <your-postgres-app-name>
```

### DATABASE_URLの設定

Fly.ioのPostgreSQLを使用している場合：

```bash
# PostgreSQLアプリの名前を確認
fly apps list

# DATABASE_URLを設定（Fly.ioが自動生成する場合もある）
fly secrets set -a contactformapi DATABASE_URL='postgres://user:password@host:5432/dbname?sslmode=require'
```

**Fly.ioのPostgreSQLを使用している場合、通常は自動的に`DATABASE_URL`が設定されます。**

## 5. デバッグ用のログ出力

コードにデバッグログが追加されているので、以下の情報がログに出力されます：

- データベース接続状態
- テーブル存在確認
- 認証試行の詳細（username, password hashの一部）

ログを確認：

```bash
fly logs -a contactformapi
```

## 6. よくある問題と解決方法

### 問題: "table 'contacts' does not exist"

**解決方法**: マイグレーションを実行（上記のステップ2を参照）

### 問題: "invalid credentials" (401)

**解決方法**: 
1. `ADMIN_PASSWORD_HASH`が正しく設定されているか確認
2. パスワードハッシュを再生成して再設定
3. アプリを再起動

### 問題: "failed to ping database"

**解決方法**:
1. `DATABASE_URL`が正しく設定されているか確認
2. Fly.ioのPostgreSQLアプリが起動しているか確認
3. ネットワーク設定を確認

## 7. 完全な再デプロイ手順

問題が解決しない場合、完全に再デプロイ：

```bash
# 1. 現在のアプリを停止（必要に応じて）
fly apps destroy -a contactformapi

# 2. 新しくアプリを作成
fly launch

# 3. 環境変数を設定
fly secrets set -a contactformapi \
  JWT_SECRET=okayu256 \
  ADMIN_USERNAME=admin \
  ADMIN_PASSWORD_HASH='$2a$10$...' \
  DATABASE_URL='postgres://...'

# 4. マイグレーションを実行
fly ssh console -a contactformapi
psql $DATABASE_URL -f migrations/001_create_contacts_table.up.sql

# 5. アプリをデプロイ
fly deploy
```

## 8. 次のステップ

1. ログを確認してエラーの詳細を特定
2. マイグレーションを実行
3. 環境変数を正しく設定
4. アプリを再起動
5. 再度テスト

問題が解決しない場合は、ログの内容を共有してください。

