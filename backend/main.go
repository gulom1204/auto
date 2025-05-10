package main

import (
	"autoparts/middleware"
	"autoparts/views"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()

	// Basic route
	r.GET("/", views.GetHome)

	r.POST("/sign_up", views.Registration)
	r.POST("/login", views.Login)
	adminGroup := r.Group("/Admin")
	ApplyAdminMiddleware(adminGroup) 
	
	{
		adminGroup.POST("/create_items", views.CreateItems)
	}
	
	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 

func ApplyAdminMiddleware(group *gin.RouterGroup) {
	group.Use(middleware.IsAdmin())
}