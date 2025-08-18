package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS cats (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		date_of_photo DATE,
		location TEXT,
		photo_path TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	)`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized and table created.")
}
