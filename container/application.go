package container

import (
	"plant-reminder/config"
	"plant-reminder/controllers"
	"plant-reminder/service"
)

type Application struct {
	PlantService    *service.PlantService
	UserService     *service.UserService
	ReminderService *service.ReminderService

	HealthController   *controllers.HealthController
	PlantController    *controllers.PlantController
	UserController     *controllers.UserController
	ReminderController *controllers.ReminderController
}

func NewApplication() *Application {
	db := config.DB
	plantService := service.NewPlantService(db)
	userService := service.NewUserService(db)
	reminderService := service.NewReminderService(plantService, db)

	healthController := controllers.NewHealthController()
	plantController := controllers.NewPlantController(plantService)
	userController := controllers.NewUserController(userService)
	reminderController := controllers.NewReminderController(reminderService)

	return &Application{
		PlantService:    plantService,
		UserService:     userService,
		ReminderService: reminderService,

		HealthController:   healthController,
		PlantController:    plantController,
		UserController:     userController,
		ReminderController: reminderController,
	}
}
