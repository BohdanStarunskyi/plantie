package controllers

import (
	"net/http"
	"plant-reminder/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")

	var plant models.Plant
	err := ctx.ShouldBindJSON(&plant)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plant.UserID = userId

	err = plant.Save()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plant": plant})
}

func GetPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plant, err := models.GetPlant(plantId, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plant": plant})
}

func GetPlants(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plants, err := models.GetPlants(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plants": plants})
}

func UpdatePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")

	var plant models.Plant
	err := ctx.ShouldBindJSON(&plant)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if plant.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "plantId must be valid"})
		return
	}

	err = plant.Update(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plant": plant})
}

func DeletePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DeletePlant(userId, plantId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plantId": plantId})
}
