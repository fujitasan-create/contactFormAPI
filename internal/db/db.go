package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// テーブルの存在確認
	if err := checkTableExists(); err != nil {
		log.Printf("WARNING: Table check failed: %v", err)
		log.Println("NOTE: Make sure migrations have been run. The 'contacts' table may not exist.")
	}

	return nil
}

// checkTableExists はcontactsテーブルが存在するか確認する
func checkTableExists() error {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'contacts'
		)
	`
	err := DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("table 'contacts' does not exist - migrations may not have been run")
	}

	log.Println("Table 'contacts' exists")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
