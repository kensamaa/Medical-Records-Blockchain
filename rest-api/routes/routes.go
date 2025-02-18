package routes

import (
	"github.com/kensamaa/blockchain-medical-records/rest-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Public routes (no auth middleware)
	public := router.Group("/api")
	{
		public.POST("/login", controllers.Login)
	}

	// Protected routes - here you would attach your JWT middleware
	protected := router.Group("/api")
	// e.g., protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/records", controllers.CreateRecord)
		protected.GET("/records/:id", controllers.GetRecord)
		protected.PUT("/records/:id", controllers.UpdateRecord)
	}
}
