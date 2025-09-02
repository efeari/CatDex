package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/efeari/catdex/internal/db"
	"github.com/efeari/catdex/internal/store.go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)
	dbMaxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	dbMaxOpenConns, err := strconv.Atoi(dbMaxOpenConnsStr)
	if err != nil || dbMaxOpenConns <= 0 {
		dbMaxOpenConns = 30
	}

	dbMaxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	dbMaxIdleConns, err := strconv.Atoi(dbMaxIdleConnsStr)
	if err != nil || dbMaxIdleConns <= 0 {
		dbMaxIdleConns = 30
	}

	dbMaxIdleTime, found := os.LookupEnv("DB_MAX_IDLE_TIME")
	if !found {
		dbMaxIdleTime = "15m"
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
		db: dbConfig{
			addr:               connStr,
			maxOpenConnections: dbMaxOpenConns,
			maxIdleConnections: dbMaxIdleConns,
			maxIdleTime:        dbMaxIdleTime,
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConnections,
		cfg.db.maxIdleConnections,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Fatal(err)
	}

	//defer db.Close()

	log.Printf("database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
