package controllers

import (
	"net/http"
	"plant-reminder/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"log"
)

func AddReminder(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		log.Printf("AddReminder: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reminder models.Reminder
	err = ctx.ShouldBindJSON(&reminder)
	if err != nil {
		log.Printf("AddReminder: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminder.UserID = userID
	reminder.PlantID = plantID
	err = reminder.Save()
	if err != nil {
		log.Printf("AddReminder: failed to save reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"reminder": reminder})
}

func GetReminders(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		log.Printf("GetReminders: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminders, err := models.GetReminders(userId, plantID)
	if err != nil {
		log.Printf("GetReminders: failed to get reminders: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminders": reminders})
}

func DeleteReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		log.Printf("DeleteReminder: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminderId, err := strconv.ParseInt(ctx.Param("reminderId"), 10, 64)
	if err != nil {
		log.Printf("DeleteReminder: invalid reminder id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reminder := models.Reminder{
		ID:      reminderId,
		PlantID: plantID,
		UserID:  userId,
	}
	err = reminder.Delete()
	if err != nil {
		log.Printf("DeleteReminder: failed to delete reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func UpdateReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		log.Printf("UpdateReminder: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reminder models.Reminder
	err = ctx.ShouldBindJSON(&reminder)
	if err != nil {
		log.Printf("UpdateReminder: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminder.PlantID = plantID
	if reminder.ID == 0 {
		log.Printf("UpdateReminder: invalid reminder ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "reminderId must be valid"})
		return
	}

	err = reminder.Update(userId)
	if err != nil {
		log.Printf("UpdateReminder: failed to update reminder: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminder": reminder})
}
