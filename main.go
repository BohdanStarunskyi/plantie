package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"plant-reminder/config"
	"plant-reminder/container"
	"plant-reminder/models"
	"plant-reminder/routes"
	"plant-reminder/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	app := initApp()

	server := setupServer(app)
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

func initApp() *container.Application {
	initDatabase()

	// Give the database a moment to settle before running migrations
	time.Sleep(1 * time.Second)

	runMigrations()
	initNotifier()

	app := container.NewApplication()

	setupCrons(app)

	return app
}

func setupServer(app *container.Application) *http.Server {
	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRouter(router, app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Addr:    ":" + port,
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
		log.Printf("Migration warning: %v", err)
	} else {
		log.Println("Database migration completed successfully")
	}
}

func setupCrons(app *container.Application) {
	if err := app.ReminderService.SetReminders(); err != nil {
		log.Fatalf("failed to start cron jobs: %v", err)
	}
}

func initNotifier() {
	if err := utils.InitNotifier(); err != nil {
		log.Fatalf("Failed to init notifier: %v", err)
	}
}
