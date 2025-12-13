# データベースセットアップガイド

## PostgreSQLのインストール

PostgreSQLがインストールされていない場合は、以下のいずれかの方法でインストールしてください。

### Windowsでのインストール方法

1. **公式インストーラーを使用**
   - https://www.postgresql.org/download/windows/ からダウンロード
   - インストール時にパスワードを設定（デフォルトユーザー: postgres）

2. **Chocolateyを使用（推奨）**
   ```powershell
   choco install postgresql
   ```

3. **Dockerを使用（開発環境推奨）**
   ```powershell
   docker run --name postgres-contactform -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=contactform -p 5432:5432 -d postgres:15
   ```

## データベースのセットアップ

### 方法1: PowerShellスクリプトを使用（推奨）

```powershell
.\setup_database.ps1
```

### 方法2: 手動でセットアップ

1. **データベースの作成**
   ```powershell
   psql -U postgres -c "CREATE DATABASE contactform;"
   ```

2. **マイグレーションの実行**
   ```powershell
   psql -U postgres -d contactform -f migrations/001_create_contacts_table.up.sql
   ```

### 方法3: Dockerを使用

既にDockerでPostgreSQLを起動している場合：

```powershell
# マイグレーションの実行
docker exec -i postgres-contactform psql -U postgres -d contactform < migrations/001_create_contacts_table.up.sql
```

または、PowerShellから：

```powershell
Get-Content migrations/001_create_contacts_table.up.sql | docker exec -i postgres-contactform psql -U postgres -d contactform
```

## データベース接続の確認

```powershell
psql -U postgres -d contactform -c "\dt"
```

`contacts` テーブルが表示されれば成功です。

## トラブルシューティング

### PostgreSQLに接続できない場合

1. PostgreSQLサービスが起動しているか確認
   ```powershell
   Get-Service -Name postgresql*
   ```

2. サービスを起動
   ```powershell
   Start-Service postgresql-x64-15  # バージョンに応じて変更
   ```

### パスワード認証エラーの場合

`.env`ファイルの`DATABASE_URL`を確認し、正しいパスワードを設定してください：

```
DATABASE_URL=postgres://postgres:あなたのパスワード@localhost:5432/contactform?sslmode=disable
```

