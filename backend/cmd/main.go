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
	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(middleware.RecoveryMiddleware())

	// Serve static files from the "uploads" directory
	r.Static("/static", "./uploads")
	r.Static("/photos", "C:/Users/efear/Documents/VS Code Projects/CatDex/photos")

	// Define API routes
	routes.RegisterRoutes(r, database.DB)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
