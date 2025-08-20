package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/efear/catdex/models"
	"github.com/gin-gonic/gin"
)

func GetRandomCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Getting random cat")
		cat, err := getRandomCatFromDB(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add a new field photo_url
		catWithURL := map[string]interface{}{
			"id":          cat.ID,
			"name":        cat.Name,
			"dateOfPhoto": cat.DateOfPhoto,
			"location":    cat.Location,
			"createdAt":   cat.CreatedAt,
			// Construct URL for frontend
			"photo_url": fmt.Sprintf("http://localhost:8080/images/%d", cat.ID),
		}

		c.JSON(http.StatusOK, catWithURL)
	}
}

func getRandomCatFromDB(db *sql.DB) (*models.Cat, error) {
	var cat models.Cat
	row := db.QueryRow(`
		SELECT *
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

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/api/cats/random", GetRandomCat(db))
	// Add other routes here
	// r.GET("/api/cats", GetAllCats(db))
	// r.GET("/api/cats/:id", GetCatByID(db))
	// r.POST("/api/cats", PostCat(db))
}
