package controllers

import (
	"example/event-app/config"
	"example/event-app/models"
	"net/http"
	"time"
	"os"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister
	// Validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash Password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed Encrypt Password",
		})
		return
	}

	// Save to DB
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Email maybe used",
		})
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Register Success",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"event": user.Event,
		},
	})

}

func LoginUser(context *gin.Context) {
	var input AuthInputLogin
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	userData:= config.DB.Where("email = ?", input.Email).First(&user)
	if userData.Error != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found",
		})
		return
	}

	errMatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatch != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password not match",
		})
		return
	}

	// make token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, errToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errToken != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user" : gin.H{
				"id":	user.ID,
				"name": user.Name,
				"email": user.Email,
				"event": user.Event,
			},
		})
}

func GetCurrentUser(context *gin.Context) {
	// Get userID from context
	userID, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to get user from context",
		})
		return
	}

	var user models.User
	userData := config.DB.Select("id", "name", "email").First(&user, userID).Error
	if userData != nil{
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"event": user.Event,
		},
	})	
}