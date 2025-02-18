package main

import (
	"log"

	"github.com/kensamaa/blockchain-medical-records/rest-api/routes" // adjust the module path if necessary

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router instance
	router := gin.Default()

	// Register all API routes
	routes.RegisterRoutes(router)

	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
