package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/efear/catdex/models"
	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func GetRandomCat(c *gin.Context) {
	fmt.Println("Getting random cat")
	cat, err := getRandomCatFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func getRandomCatFromDB() (*models.Cat, error) {
	var cat models.Cat
	row := DB.QueryRow("SELECT * FROM cats ORDER BY RANDOM() LIMIT 1")
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
