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
		c.JSON(http.StatusOK, cat)
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

// func getAllCats(c *gin.Context) {
// 	// Implementation here
// }

// func getCatByID(c *gin.Context) {
// 	// Implementation here
// }

// func postCat(c *gin.Context) {
// 	// Implementation here
// }
