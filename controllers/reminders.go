package controllers

import (
	"net/http"
	"plant-reminder/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddReminder(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reminder models.Reminder
	err = ctx.ShouldBindJSON(&reminder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminder.UserID = userID
	reminder.PlantID = plantID
	err = reminder.Save()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"reminder": reminder})
}

func GetReminders(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminders, err := models.GetReminders(userId, plantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminders": reminders})
}

func DeleteReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminderId, err := strconv.ParseInt(ctx.Param("reminderId"), 10, 64)
	if err != nil {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminderId": reminderId})
}

func UpdateReminder(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantIdStr := ctx.Param("id")
	plantID, err := strconv.ParseInt(plantIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var reminder models.Reminder
	err = ctx.ShouldBindJSON(&reminder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reminder.PlantID = plantID
	if reminder.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "reminderId must be valid"})
		return
	}

	err = reminder.Update(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reminder": reminder})
}
