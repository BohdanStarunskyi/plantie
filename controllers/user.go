package controllers

import (
	"fmt"
	"log"
	"net/http"
	"plant-reminder/models"
	"plant-reminder/utils"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := utils.Validate.Var(user.Email, "required,email"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	if err := utils.Validate.Var(user.Password, "required,min=6"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token, verifiedUser, err := models.VerifyUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": verifiedUser})
}

func SignUp(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
