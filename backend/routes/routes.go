package routes

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/efear/catdex/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetRandomCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Getting random cat")
		cat, err := getRandomCatFromDB(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, buildCatResponse(cat, db))
	}
}

func getRandomCatFromDB(db *sql.DB) (*models.Cat, error) {
	var cat models.Cat
	row := db.QueryRow(`
		SELECT id, name, date_of_photo, location, photo_path, created_at
		FROM cats
		ORDER BY RANDOM()
		LIMIT 1;
	`)
	err := row.Scan(&cat.ID, &cat.Name, &cat.DateOfPhoto, &cat.Location, &cat.PhotoPath, &cat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func buildCatResponse(cat *models.Cat, db *sql.DB) map[string]interface{} {
	photoFilename := filepath.Base(cat.PhotoPath)

	// determine availability of neighbors
	nextCat, _ := getNextCatFromDB(db, cat.CreatedAt)
	prevCat, _ := getPreviousCatFromDB(db, cat.CreatedAt)

	return map[string]interface{}{
		"id":            cat.ID,
		"name":          cat.Name,
		"date_of_photo": cat.DateOfPhoto,
		"location":      cat.Location,
		"created_at":    cat.CreatedAt,
		"photo_url":     fmt.Sprintf("http://localhost:8080/photos/%s", photoFilename),
		"has_next":      nextCat != nil,
		"has_previous":  prevCat != nil,
	}
}

func GetCatByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var cat models.Cat
		row := db.QueryRow(`
			SELECT id, name, date_of_photo, location, photo_path, created_at
			FROM cats
			WHERE id = $1;
		`, id)
		if err := row.Scan(&cat.ID, &cat.Name, &cat.DateOfPhoto, &cat.Location, &cat.PhotoPath, &cat.CreatedAt); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
			return
		}

		c.JSON(http.StatusOK, buildCatResponse(&cat, db))
	}
}

func GetNextCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		afterStr := c.Query("after")
		after, err := time.Parse(time.RFC3339, afterStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		cat, err := getNextCatFromDB(db, after)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if cat == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "no next cat"})
			return
		}

		c.JSON(http.StatusOK, buildCatResponse(cat, db))
	}
}

func GetPreviousCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		beforeStr := c.Query("before")
		before, err := time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		cat, err := getPreviousCatFromDB(db, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if cat == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "no previous cat"})
			return
		}

		c.JSON(http.StatusOK, buildCatResponse(cat, db))
	}
}

func getNextCatFromDB(db *sql.DB, after time.Time) (*models.Cat, error) {
	query := (`
		SELECT id, name, date_of_photo, location, photo_path, created_at
        FROM cats
        WHERE created_at > $1
        ORDER BY created_at ASC
        LIMIT 1
	`)

	cat := &models.Cat{}
	err := db.QueryRow(query, after).Scan(
		&cat.ID,
		&cat.Name,
		&cat.DateOfPhoto,
		&cat.Location,
		&cat.PhotoPath,
		&cat.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // no next cat
	}
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func getPreviousCatFromDB(db *sql.DB, before time.Time) (*models.Cat, error) {
	query := (`
        SELECT id, name, date_of_photo, location, photo_path, created_at
        FROM cats
        WHERE created_at < $1
        ORDER BY created_at DESC
        LIMIT 1
    `)

	cat := &models.Cat{}
	err := db.QueryRow(query, before).Scan(
		&cat.ID,
		&cat.Name,
		&cat.DateOfPhoto,
		&cat.Location,
		&cat.PhotoPath,
		&cat.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func PostCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		dateOfPhoto := c.PostForm("date_of_photo")
		location := c.PostForm("location")

		file, err := c.FormFile("photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Photo upload failed"})
			return
		}

		catID := uuid.New().String()
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		filename := fmt.Sprintf("%s%s", catID, ext)

		// Ensure folder exists
		err = os.MkdirAll("../../photos", os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photos folder"})
			return
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer src.Close()

		// Create destination file
		out, err := os.Create(filepath.Join("../../photos", filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
		defer out.Close()

		// Copy content
		_, err = io.Copy(out, src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
			return
		}

		// Insert into DB
		_, err = db.Exec(`
			INSERT INTO cats (id, name, date_of_photo, location, photo_path, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`, catID, name, dateOfPhoto, location, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cat data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cat uploaded successfully", "id": catID})
	}
}

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/api/cats/random", GetRandomCat(db))
	r.GET("/api/cats/:id", GetCatByID(db))
	r.GET("/api/cats/previous", GetPreviousCat(db))
	r.GET("/api/cats/next", GetNextCat(db))
	r.POST("/api/cats", PostCat(db))
}
