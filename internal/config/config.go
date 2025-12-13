package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	JWTSecret         string
	AdminUsername     string
	AdminPasswordHash string
	Port              string
}

func Load() (*Config, error) {
	// .envファイルを読み込む
	// プロジェクトルート（go.modがあるディレクトリ）を探す
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	projectRoot := wd

	// カレントディレクトリから親ディレクトリへ遡ってgo.modを探す
	maxDepth := 10 // 無限ループを防ぐ
	depth := 0
	for depth < maxDepth {
		goModPath := filepath.Join(projectRoot, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// go.modが見つかった = プロジェクトルート
			log.Printf("Found project root: %s", projectRoot)
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			// ルートディレクトリに到達したがgo.modが見つからない
			// カレントディレクトリを使用
			projectRoot = wd
			log.Printf("go.mod not found, using current directory: %s", projectRoot)
			break
		}
		projectRoot = parent
		depth++
	}

	// プロジェクトルートの.envファイルを読み込む
	envPath := filepath.Join(projectRoot, ".env")
	var loaded bool
	var adminPasswordHashFromFile string

	if _, err := os.Stat(envPath); err == nil {
		// まず通常の方法で読み込む
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env file from: %s", envPath)
			loaded = true
		}

		// ADMIN_PASSWORD_HASHだけ手動で読み込む（変数展開を避けるため）
		adminPasswordHashFromFile = readPasswordHashFromEnvFile(envPath)
	} else {
		log.Printf(".env file not found at: %s", envPath)
	}

	// プロジェクトルートで見つからない場合、カレントディレクトリの.envを試す
	if !loaded {
		currentEnv := ".env"
		if _, err := os.Stat(currentEnv); err == nil {
			if err := godotenv.Load(currentEnv); err == nil {
				log.Printf("Loaded .env file from: %s (current directory)", currentEnv)
				loaded = true
			}
			// ADMIN_PASSWORD_HASHだけ手動で読み込む
			if adminPasswordHashFromFile == "" {
				adminPasswordHashFromFile = readPasswordHashFromEnvFile(currentEnv)
			}
		}
	}

	if !loaded {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/contactform?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		AdminUsername:     getEnv("ADMIN_USERNAME", "admin"),
		AdminPasswordHash: getEnv("ADMIN_PASSWORD_HASH", ""),
		Port:              getEnv("PORT", "8080"),
	}

	// 手動で読み込んだハッシュがあればそれを使用
	if adminPasswordHashFromFile != "" {
		cfg.AdminPasswordHash = adminPasswordHashFromFile
		log.Printf("Loaded ADMIN_PASSWORD_HASH from .env file (length: %d)", len(cfg.AdminPasswordHash))
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if cfg.AdminPasswordHash == "" {
		return nil, fmt.Errorf("ADMIN_PASSWORD_HASH environment variable is required")
	}

	return cfg, nil
}

// readPasswordHashFromEnvFile は.envファイルからADMIN_PASSWORD_HASHを直接読み込む
func readPasswordHashFromEnvFile(envPath string) string {
	file, err := os.Open(envPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// コメント行や空行をスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// ADMIN_PASSWORD_HASH= で始まる行を探す
		if strings.HasPrefix(line, "ADMIN_PASSWORD_HASH=") {
			// = の後の部分を取得
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				value := parts[1]
				// 引用符を削除
				value = strings.Trim(value, `"'`)
				return value
			}
		}
	}
	return ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
