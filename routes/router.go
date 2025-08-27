package routes

import (
	"plant-reminder/container"
	"plant-reminder/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine, app *container.Application) {

	healthController := app.HealthController
	plantController := app.PlantController
	userController := app.UserController
	reminderController := app.ReminderController

	engine.GET("/ping", healthController.Ping)

	engine.POST("/login", userController.Login)
	engine.POST("/signup", userController.SignUp)
	engine.POST("/refresh", userController.RefreshToken)

	authGroup := engine.Group("/", middleware.VerifyAuth)

	authGroup.POST("/user/push_token", userController.SetPushToken)
	authGroup.DELETE("/user", userController.DeleteUser)
	authGroup.GET("/user/me", userController.GetMyProfile)

	authGroup.POST("/plant", plantController.AddPlant)
	authGroup.DELETE("/plant/:id", plantController.DeletePlant)
	authGroup.PUT("/plant/:id", plantController.UpdatePlant)
	authGroup.GET("/plant/:id", plantController.GetPlant)
	authGroup.GET("/plants", plantController.GetPlants)

	authGroup.POST("/plant/:id/reminder", reminderController.AddReminder)
	authGroup.DELETE("/plant/:id/reminder/:reminderId", reminderController.DeleteReminder)
	authGroup.PUT("/plant/:id/reminder", reminderController.UpdateReminder)
	authGroup.GET("/plant/:id/reminders", reminderController.GetPlantReminders)
	authGroup.GET("/plant/reminders", reminderController.GetAllReminders)
}
