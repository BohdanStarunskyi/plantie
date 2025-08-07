package routes

import (
	"plant-reminder/controllers"
	"plant-reminder/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(engine *gin.Engine) {
	engine.GET("/ping", controllers.Ping)

	engine.POST("/login", controllers.Login)
	engine.POST("/signup", controllers.SignUp)
	engine.POST("/refresh", controllers.RefreshToken)

	authGroup := engine.Group("/", middleware.VerifyAuth)

	authGroup.POST("/user/push_token", controllers.SetPushToken)

	authGroup.POST("/plant", controllers.AddPlant)
	authGroup.DELETE("/plant/:id", controllers.DeletePlant)
	authGroup.PUT("/plant", controllers.UpdatePlant)
	authGroup.GET("/plant/:id", controllers.GetPlant)
	authGroup.GET("/plants", controllers.GetPlants)

	authGroup.POST("/plant/:id/reminder", controllers.AddReminder)
	authGroup.DELETE("plant/:id/reminder/:reminderId", controllers.DeleteReminder)
	authGroup.PUT("/plant/:id/reminder", controllers.UpdateReminder)
	authGroup.GET("/plant/:id/reminders", controllers.GetPlantReminders)
	authGroup.GET("/plant/reminders", controllers.GetAllReminders)
}
