package main

import (
	"log"
	"os"
	"plant-reminder/config"
	"plant-reminder/models"
	"plant-reminder/routes"
	"plant-reminder/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitDb()
	migrateDb()
	setCrons()
	initNotifier()
	server := gin.Default()
	routes.SetupRouter(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	server.Run(port)
}

func setCrons() {
	err := models.SetReminders()
	if err != nil {
		panic("Failed to start cron jobs: " + err.Error())
	}
}

func migrateDb() {
	err := config.DB.AutoMigrate(&models.User{}, &models.Plant{}, &models.Reminder{})
	if err != nil {
		panic("Failed to migrate models: " + err.Error())
	}
}

func initNotifier() {
	err := utils.InitNotifier()
	if err != nil {
		panic("Failed to init notifier: " + err.Error())
	}
}
