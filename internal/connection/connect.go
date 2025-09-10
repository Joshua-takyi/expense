package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitPsql() (*sql.DB, error) {
	var connStr string

	// Check if DATABASE_URL is set (for production/cloud deployment)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// Load environment variables from .env.local file
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatalf("failed to read data from the env file: %v", err)
		}

		// Now read the environment variables after loading the .env file
		dbname := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("PASSWORD")

		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			return nil, fmt.Errorf("missing required database configuration: host=%s, port=%s, user=%s, password=%s, dbname=%s",
				host, port, user, password, dbname)
		}

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("failed to close the database connection: %v", err)
	}

}
