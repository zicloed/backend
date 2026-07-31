package config

import (
	"example/event-app/models"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(){
	dsn:= os.Getenv("DATABASE_URL")
	if dsn == ""{
		log.Fatal("Environment variabel belum di isi")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal("Failed Connection to Database", err)
	}

	err = database.AutoMigrate(&models.User{}, models.Event{})
	if err != nil{
		log.Fatal("Failed Migration Database", err)
	}

	DB= database
	log.Println("Database Connection Success")
}