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

		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventsbyId)

		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.RequiredAuth())
		{
			protected.GET("/auth/me", controllers.GetCurrentUser)
			protected.POST("/events", controllers.CreateEvents)
			protected.PUT("/events/:id", controllers.UpdateEvent)
			protected.DELETE("/events/:id", controllers.DeleteEvent)
		}
	}

	// server port
	server.Run(":8080")
}
