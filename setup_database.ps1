# PostgreSQLデータベースセットアップスクリプト
# このスクリプトはPostgreSQLがインストールされていることを前提としています

Write-Host "PostgreSQLデータベースのセットアップを開始します..." -ForegroundColor Green

# データベースが存在するか確認
$dbExists = psql -U postgres -lqt | Select-String -Pattern "contactform"

if ($dbExists) {
    Write-Host "データベース 'contactform' は既に存在します。" -ForegroundColor Yellow
    $response = Read-Host "削除して再作成しますか？ (y/n)"
    if ($response -eq "y") {
        Write-Host "データベースを削除しています..." -ForegroundColor Yellow
        psql -U postgres -c "DROP DATABASE IF EXISTS contactform;"
    } else {
        Write-Host "既存のデータベースを使用します。" -ForegroundColor Green
        exit 0
    }
}

# データベースを作成
Write-Host "データベース 'contactform' を作成しています..." -ForegroundColor Green
psql -U postgres -c "CREATE DATABASE contactform;"

if ($LASTEXITCODE -eq 0) {
    Write-Host "データベースの作成に成功しました。" -ForegroundColor Green
} else {
    Write-Host "データベースの作成に失敗しました。PostgreSQLがインストールされ、サービスが起動していることを確認してください。" -ForegroundColor Red
    exit 1
}

# マイグレーションを実行
Write-Host "マイグレーションを実行しています..." -ForegroundColor Green
psql -U postgres -d contactform -f migrations/001_create_contacts_table.up.sql

if ($LASTEXITCODE -eq 0) {
    Write-Host "マイグレーションの実行に成功しました。" -ForegroundColor Green
    Write-Host "セットアップが完了しました！" -ForegroundColor Green
} else {
    Write-Host "マイグレーションの実行に失敗しました。" -ForegroundColor Red
    exit 1
}

