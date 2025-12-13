# Docker実装ガイド

## 概要

このプロジェクトは**Docker Composeのみ**で開発・実行します。`go run cmd/api/main.go`は使用しません。

## ファイル構成

- `Dockerfile`: マルチステージビルドを使用したAPIアプリケーションのビルド定義
- `docker-compose.yml`: PostgreSQLとAPIアプリケーションのオーケストレーション
- `.dockerignore`: Dockerビルド時に除外するファイル
- `.env`: 環境変数設定ファイル（Docker Compose用）

## セットアップ

### 1. 環境変数の設定

プロジェクトルートに`.env`ファイルを作成し、以下の環境変数を設定してください：

```env
JWT_SECRET=okayu256
ADMIN_USERNAME=admin
ADMIN_PASSWORD_HASH=$$2a$$10$$F9R3SCNxCitnJotoN1.vMe9kbDOC/yKIVJMW87J1wBGpQYw3uBG2m
```

**重要**: `ADMIN_PASSWORD_HASH`の`$`文字は`$$`でエスケープしてください。  
docker-composeは`.env`ファイルの`$`文字を変数展開として解釈するため、bcryptハッシュの`$`文字は`$$`でエスケープする必要があります。

### 2. Dockerコンテナの起動

```bash
docker-compose up -d
```

これで、PostgreSQLとAPIアプリケーションの両方が起動します。

### 3. ログの確認

```bash
# APIコンテナのログ（リアルタイム）
docker-compose logs -f api

# PostgreSQLコンテナのログ
docker-compose logs -f postgres

# すべてのログ
docker-compose logs -f
```

### 4. コンテナの停止

```bash
docker-compose down
```

### 5. データの削除（ボリュームも含む）

```bash
docker-compose down -v
```

## 開発ワークフロー

### コード変更後の再ビルド

コードを変更した後、コンテナを再ビルドして再起動します：

```bash
# コンテナを停止
docker-compose down

# イメージを再ビルド
docker-compose build

# コンテナを起動
docker-compose up -d
```

または、一度に実行：

```bash
docker-compose up -d --build
```

### コンテナ内でのコマンド実行

```bash
# APIコンテナ内でシェルを実行
docker-compose exec api sh

# データベース接続テスト
docker-compose exec api ./main --help
```

## サービス

### PostgreSQL
- ポート: `5432`
- データベース名: `contactform`
- ユーザー名: `postgres`
- パスワード: `postgres`
- 接続文字列: `postgres://postgres:postgres@postgres:5432/contactform?sslmode=disable`

### API
- ポート: `8080`
- Swagger UI: http://localhost:8080/swagger/index.html
- ヘルスチェック: http://localhost:8080/health
- 公開API: http://localhost:8080/contact (POST)
- 管理API: http://localhost:8080/admin/login (POST), http://localhost:8080/admin/messages (GET)

## トラブルシューティング

### ポートが既に使用されている場合

ポート8080が既に使用されている場合、既存のプロセスを停止してください。

```powershell
# Windowsの場合
netstat -ano | findstr :8080
Stop-Process -Id <PID> -Force
```

### 環境変数が正しく読み込まれない場合

1. `.env`ファイルで`$`文字を`$$`でエスケープしているか確認してください
2. `.env`ファイルがプロジェクトルートにあるか確認してください
3. コンテナ内の環境変数を確認：

```bash
docker-compose exec api env | grep -E "JWT_SECRET|ADMIN_USERNAME|ADMIN_PASSWORD_HASH"
```

### データベース接続エラー

PostgreSQLコンテナが正常に起動しているか確認してください。

```bash
# コンテナの状態確認
docker-compose ps

# PostgreSQLのログ確認
docker-compose logs postgres

# PostgreSQLコンテナ内で接続テスト
docker-compose exec postgres psql -U postgres -d contactform -c "SELECT 1;"
```

### ビルドエラー

```bash
# キャッシュなしで再ビルド
docker-compose build --no-cache

# すべてのコンテナとイメージを削除して再ビルド
docker-compose down --rmi all
docker-compose build --no-cache
docker-compose up -d
```

## よく使うコマンド

```bash
# コンテナの状態確認
docker-compose ps

# ログの確認（最新50行）
docker-compose logs --tail=50 api

# コンテナの再起動
docker-compose restart api

# コンテナの停止と削除
docker-compose down

# ボリュームも含めて完全削除
docker-compose down -v

# イメージの再ビルドと起動
docker-compose up -d --build
```

## 注意事項

- **`go run cmd/api/main.go`は使用しません**。すべてDocker Composeで実行します
- `.env`ファイルは`.gitignore`に含まれているため、Gitにはコミットされません
- 本番環境では、環境変数は適切なシークレット管理システム（AWS Secrets Manager、HashiCorp Vaultなど）を使用してください
