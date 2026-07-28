package main

import (
	"example/event-app/config"
	"example/event-app/models"
	"log"
	"net/http"

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
	// server port
	server.Run(":8080")
	// router
	api := server.Group("/api")
	{
		api.POST("/events", createEvents)
		api.GET("/events", getEvents)
	}
}

// Function Handler
// createEvents
func createEvents(context *gin.Context) {
	var events models.Event
	err := context.ShouldBindJSON(&events)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	events.UserId = 1

	// save input
	events.Save()
	context.JSON(http.StatusCreated, gin.H{
		"message": "Event Created",
		"events":  events,
	})
}

// getEvents
func getEvents(context *gin.Context) {
	events := models.GetAllEvent()
	context.JSON(http.StatusOK, events)
}
