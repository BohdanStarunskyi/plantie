package controllers

import (
	"log"
	"net/http"
	"plant-reminder/models"
	"plant-reminder/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")

	var plant models.Plant
	err := ctx.ShouldBindJSON(&plant)
	if err != nil {
		log.Printf("AddPlant: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.Validate.Struct(plant); err != nil {
		log.Printf("AddPlant: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plant.UserID = userId
	err = plant.Save()
	if err != nil {
		log.Printf("AddPlant: failed to save plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plant": plant})
}

func GetPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("GetPlant: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plant, err := models.GetPlant(plantId, userId)
	if err != nil {
		log.Printf("GetPlant: failed to get plant: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plant": plant})
}

func GetPlants(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plants, err := models.GetPlants(userID)
	if err != nil {
		log.Printf("GetPlants: failed to get plants: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plants": plants})
}

func UpdatePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")

	var plant models.Plant
	err := ctx.ShouldBindJSON(&plant)
	if err != nil {
		log.Printf("UpdatePlant: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if plant.ID == 0 {
		log.Printf("UpdatePlant: invalid plant ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "plantId must be valid"})
		return
	}

	if err := utils.Validate.Struct(plant); err != nil {
		log.Printf("UpdatePlant: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = plant.Update(userId)
	if err != nil {
		log.Printf("UpdatePlant: failed to update plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plant": plant})
}

func DeletePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("DeletePlant: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DeletePlant(userId, plantId)
	if err != nil {
		log.Printf("DeletePlant: failed to delete plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
