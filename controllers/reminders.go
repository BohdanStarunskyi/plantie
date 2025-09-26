package controllers

import (
	"net/http"
	"plant-reminder/dto"
	"plant-reminder/service"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
)

type ReminderController struct {
	reminderService service.ReminderServiceInterface
}

func NewReminderController(reminderService service.ReminderServiceInterface) *ReminderController {
	return &ReminderController{
		reminderService: reminderService,
	}
}

func (rc *ReminderController) AddReminder(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plantID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("AddReminder: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	var reminderRequest dto.ReminderCreateRequest
	if err := ctx.ShouldBindJSON(&reminderRequest); err != nil {
		log.Printf("AddReminder: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if err := reminderRequest.Validate(); err != nil {
		log.Printf("AddReminder: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminderResponse, err := rc.reminderService.CreateReminder(&reminderRequest, plantID, userID)
	if err != nil {
		log.Printf("AddReminder: failed to save reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"reminder": reminderResponse})
}

func (rc *ReminderController) GetPlantReminders(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		log.Printf("GetReminders: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminders, err := rc.reminderService.GetPlantReminders(plantID, userId)
	if err != nil {
		log.Printf("GetReminders: failed to get reminders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminders": reminders})
}

func (rc *ReminderController) GetAllReminders(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	reminders, err := rc.reminderService.GetUserReminders(userId)
	if err != nil {
		log.Printf("GetReminders: failed to get reminders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminders": reminders})
}

func (rc *ReminderController) DeleteReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	reminderId, err := strconv.ParseInt(ctx.Param("reminderId"), 10, 64)
	if err != nil {
		log.Printf("DeleteReminder: invalid reminder id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rc.reminderService.DeleteReminder(reminderId, userId)
	if err != nil {
		log.Printf("DeleteReminder: failed to delete reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (rc *ReminderController) TestReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	err := rc.reminderService.TestReminder(userId)
	if err != nil {
		log.Printf("DeleteReminder: failed to test reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (rc *ReminderController) UpdateReminder(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plantID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("UpdateReminder: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	var reminderRequest dto.ReminderUpdateRequest
	if err := ctx.ShouldBindJSON(&reminderRequest); err != nil {
		log.Printf("UpdateReminder: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	if reminderRequest.ID == 0 {
		log.Printf("UpdateReminder: invalid reminder ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "reminder ID must be set"})
		return
	}

	if err := reminderRequest.Validate(); err != nil {
		log.Printf("UpdateReminder: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rc.reminderService.UpdateReminder(&reminderRequest, userID, plantID); err != nil {
		log.Printf("UpdateReminder: failed to update reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedReminder, err := rc.reminderService.GetReminder(reminderRequest.ID, userID)
	if err != nil {
		log.Printf("UpdateReminder: failed to get updated reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"reminder": updatedReminder})
}
