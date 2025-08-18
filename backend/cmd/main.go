package main

import (
	"log"
	"os"

	"github.com/efear/catdex/database"
	"github.com/efear/catdex/middleware"
	"github.com/efear/catdex/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDB()

	r := gin.Default()

	r.Use(middleware.RecoveryMiddleware())

	// Serve static files from the "uploads" directory
	r.Static("/static", "./uploads")

	// Define API routes
	r.GET("/api/cats/random", routes.GetRandomCat(database.DB))
	// r.GET("/api/cats", routes.GetAllCats)
	// r.GET("/api/cats/:id", routes.GetCatByID)
	// r.POST("/api/cats", routes.PostCat)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
