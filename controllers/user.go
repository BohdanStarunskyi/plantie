package controllers

import (
	"fmt"
	"net/http"
	"plant-reminder/models"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user *models.User
	ctx.ShouldBindJSON(&user)
	token, user, err := models.VerifyUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func SignUp(ctx *gin.Context) {
	var user *models.User
	ctx.ShouldBindJSON(&user)
	token, err := user.CreateUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func SetPushToken(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid push token"})
		return
	}

	err := models.SetPushToken(fmt.Sprintf("%d", userID), req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "push token set successfully"})
}
