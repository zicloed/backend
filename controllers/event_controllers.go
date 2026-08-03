package controllers

import (
	"context"
	"example/event-app/config"
	"example/event-app/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvents(c *gin.Context) {
	userID, _ := c.Get("userID")

	// recieve from form-data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Picture is required",
		})
		return
	}

	defer file.Close()

	// upload image to imagekit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	// save to database
	event := models.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		Datetime:    parsedTime,
		Image:       uploadRes.URL,
		ImageID:     uploadRes.FileID,
		UserID:      uint(userID.(int)), // Set the UserID of the event to the authenticated user's ID
	}

	var user models.User
	if err := config.DB.First(&user, event.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
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
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Data Showing Detail Event",
		"event":   event,
	})
}

func UpdateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")
	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this event",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		ik := initImageKit()

		// upload image to imagekit
		fileName := header.Filename
		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: fileName,
		})

		if errUpload == nil {
			// // deleted old image from imagekit
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}
			// if event.ImageID != "" {
			// 	ik := initImageKit()
			// 	if _, errDelete := ik.Files.Delete(context.Background(), event.ImageID); errDelete != nil {
			// 		c.JSON(http.StatusInternalServerError, gin.H{
			// 			"error": "Failed to delete old image",
			// 		})
			// 		return
			// 	}
			// update event with new image
			event.Image = uploadRes.URL
			event.ImageID = uploadRes.FileID
		}
	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.Datetime = parseTime
		}
	}

	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data get update",
		"event":   event,
	})
}

func DeleteEvent(c *gin.Context) {
	userID, _ := c.Get("userID")
	var event models.Event
	paramId := c.Param("id")

	var eventData = config.DB.First(&event, paramId).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event.UserID != uint(userID.(int)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you are not authorized to delete this event",
		})
		return
	}

	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data success to Deleted",
		"event":   event,
	})
}
