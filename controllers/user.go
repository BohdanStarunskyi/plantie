package controllers

import (
	"fmt"
	"net/http"
	"plant-reminder/models"
	"log"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user *models.User
	ctx.ShouldBindJSON(&user)
	token, user, err := models.VerifyUser(user.Email, user.Password)
	if err != nil {
		log.Printf("Login: failed to verify user: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func SignUp(ctx *gin.Context) {
	var user *models.User
	ctx.ShouldBindJSON(&user)
	token, err := user.CreateUser()
	if err != nil {
		log.Printf("SignUp: failed to create user: %v", err)
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
		log.Printf("SetPushToken: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid push token"})
		return
	}

	err := models.SetPushToken(fmt.Sprintf("%d", userID), req.Token)
	if err != nil {
		log.Printf("SetPushToken: failed to set push token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "push token set successfully"})
}
