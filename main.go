package main

import (
	"example/event-app/config"
	"example/event-app/controllers"
	"example/event-app/middlewares"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()
	// router
	api := server.Group("/api")
	{
		api.POST("/events", controllers.CreateEvents)
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventsbyId)
		api.PUT("/events/:id", controllers.UpdateEvent)
		api.DELETE("/events/:id", controllers.DeleteEvent)

		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)
		
		protected := api.Group("/")
		protected.Use(middlewares.RequiredAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)
		}
	}

	// server port
	server.Run(":8080")
}
