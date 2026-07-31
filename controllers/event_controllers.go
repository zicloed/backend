package controllers

import (
	"example/event-app/config"
	"example/event-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEvents(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, event.UserID).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := config.DB.Create(&event).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Data Created",
		"event":   event,
	})
}

func GetEvents(context *gin.Context) {
	var event []models.Event

	config.DB.Find(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data Showing All",
		"event":   event,
	})
}

func GetEventsbyId(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Event not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Data Showing Detail Event",
		"event":   event,
	})
}

func UpdateEvent(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	var input models.Event
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(input)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data get update",
		"event":   event,
	})
}

func DeleteEvent(context *gin.Context) {
	var event models.Event
	paramId := context.Param("id")

	var eventData = config.DB.First(&event, paramId).Error
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data success to Deleted",
		"event":   event,
	})
}
