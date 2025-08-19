package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/efear/catdex/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	if err := godotenv.Load("../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open test db: %v", err)
	}

	// Ensure schema exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS cats (
            id SERIAL PRIMARY KEY,
            name TEXT,
            date_of_photo DATE,
            location TEXT,
            photo_path TEXT,
            created_at TIMESTAMP DEFAULT now()
        )
    `)
	if err != nil {
		log.Fatalf("failed to create test cats table: %v", err)
	}
	return db
}

func TestGetRandomCatFromDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test cat
	_, err := db.Exec(
		`INSERT INTO cats 
	(name, date_of_photo, location, photo_path) 
	VALUES ('Sakiz',NOW(),'Istanbul','/fake/photo.jpg'
	)`)
	if err != nil {
		t.Fatal(err)
	}

	cat, err := getRandomCatFromDB(db)
	if err != nil {
		t.Fatal(err)
	}

	today := time.Now()
	expectedYear, expectedMonth, expectedDay := today.Year(), today.Month(), today.Day()

	catYear, catMonth, catDay := cat.DateOfPhoto.Year(), cat.DateOfPhoto.Month(), cat.DateOfPhoto.Day()

	if cat.Name != "Sakiz" || catYear != expectedYear || catMonth != expectedMonth || catDay != expectedDay {
		t.Errorf("Expected Sakiz on %04d-%02d-%02d, got %s on %04d-%02d-%02d",
			expectedYear, expectedMonth, expectedDay,
			cat.Name, catYear, catMonth, catDay)
	}
}

func TestGetRandomCatEndpoint(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test cat
	_, err := db.Exec(
		`INSERT INTO cats 
	(name, date_of_photo, location, photo_path) 
	VALUES ('Sakiz',CURRENT_DATE,'Istanbul','/fake/photo.jpg'
	)`)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.Default()
	router.GET("/api/cats/random", GetRandomCat(db))

	req := httptest.NewRequest("GET", "/api/cats/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var cat models.Cat
	if err := json.Unmarshal(w.Body.Bytes(), &cat); err != nil {
		t.Fatal(err)
	}

	if cat.Name != "Sakiz" {
		t.Errorf("Expected Sakiz, got %s", cat.Name)
	}
}
