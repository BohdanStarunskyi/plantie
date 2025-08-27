package controllers

import (
	"log"
	"net/http"
	"plant-reminder/dto"
	"plant-reminder/service"
	"plant-reminder/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlantController struct {
	plantService service.PlantServiceInterface
}

func NewPlantController(plantService service.PlantServiceInterface) *PlantController {
	return &PlantController{
		plantService: plantService,
	}
}

func (pc *PlantController) AddPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")

	var plantRequest dto.PlantCreateRequest
	err := ctx.ShouldBindJSON(&plantRequest)
	if err != nil {
		log.Printf("AddPlant: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.Validate.Struct(plantRequest); err != nil {
		log.Printf("AddPlant: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plantResponse, err := pc.plantService.CreatePlant(&plantRequest, userId)
	if err != nil {
		log.Printf("AddPlant: failed to save plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"plant": plantResponse})
}

func (pc *PlantController) GetPlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("GetPlant: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plantResponse, err := pc.plantService.GetPlant(plantId, userId)
	if err != nil {
		log.Printf("GetPlant: failed to get plant: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plant": plantResponse})
}

func (pc *PlantController) GetPlants(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	plants, err := pc.plantService.GetPlants(userID)
	if err != nil {
		log.Printf("GetPlants: failed to get plants: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plants": plants})
}

func (pc *PlantController) UpdatePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	var plant dto.PlantUpdateRequest
	err := ctx.ShouldBindJSON(&plant)
	if err != nil {
		log.Printf("UpdatePlant: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("GetPlant: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.Validate.Struct(plant); err != nil {
		log.Printf("UpdatePlant: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = pc.plantService.UpdatePlant(&plant, plantId, userId)
	if err != nil {
		log.Printf("UpdatePlant: failed to update plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedPlant, err := pc.plantService.GetPlant(plantId, userId)
	if err != nil {
		log.Printf("UpdatePlant: failed to get updated plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"plant": updatedPlant})
}

func (pc *PlantController) DeletePlant(ctx *gin.Context) {
	userId := ctx.GetInt64("userID")
	plantId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Printf("DeletePlant: invalid plant id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = pc.plantService.DeletePlant(userId, plantId)
	if err != nil {
		log.Printf("DeletePlant: failed to delete plant: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
