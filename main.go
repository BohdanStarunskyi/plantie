package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"plant-reminder/config"
	"plant-reminder/models"
	"plant-reminder/routes"
	"plant-reminder/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	initApp()

	server := setupServer()
	startServer(server)
	gracefulShutdown(server)
}

func loadEnv() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func initApp() {
	initDatabase()
	runMigrations()
	setupCrons()
	initNotifier()
}

func setupServer() *http.Server {
	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRouter(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return &http.Server{
		Addr:    port,
		Handler: router,
	}
}

func startServer(server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Printf("Server is running on %s", server.Addr)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}

func initDatabase() {
	config.InitDb()
}

func runMigrations() {
	err := config.DB.AutoMigrate(&models.User{}, &models.Plant{}, &models.Reminder{})
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
}

func setupCrons() {
	if err := models.SetReminders(); err != nil {
		log.Fatalf("Failed to start cron jobs: %v", err)
	}
}

func initNotifier() {
	if err := utils.InitNotifier(); err != nil {
		log.Fatalf("Failed to init notifier: %v", err)
	}
}
