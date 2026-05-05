package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbWriteUser := os.Getenv("DB_USER")
	dbWritePass := os.Getenv("DB_PASS")
	dbWriteHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbWriteUser == "" || dbWritePass == "" || dbWriteHost == "" || dbName == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	// DSN connection string
	writeDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		dbWriteUser, dbWritePass, dbWriteHost, dbName)

	// writeDSN := fmt.Sprintf("%s:%s@unix(%s)/%s", dbWriteUser, dbWritePass, dbWriteHost, dbName)

	db, err := sql.Open("mysql", writeDSN)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	if err := configureDB(db); err != nil {
		log.Fatal("Database configuration failed:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func configureDB(db *sql.DB) error {
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(500)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}
