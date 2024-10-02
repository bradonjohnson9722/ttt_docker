package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func main() {

	// Initialize MongoDB connection
	InitDB()

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// Define API endpoints
	router.POST("/start-game", StartGame)
	router.POST("/make-move", MakeMove)

	// Start the server
	log.Println("Server starting on port 1972...")
	router.Run(":1972") // Server will run on port 8080
}
