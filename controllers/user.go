package controllers

import (
	"fmt"
	"log"
	"net/http"
	"plant-reminder/models"
	"plant-reminder/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	verifiedUser, err := models.VerifyUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.SignPayload(verifiedUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := utils.SignRefreshToken(verifiedUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          verifiedUser,
	})
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

	err := user.CreateUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.SignPayload(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := utils.SignRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	token, err := utils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID in token"})
		return
	}

	newAccessToken, err := utils.SignPayload(int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new access token"})
		return
	}

	newRefreshToken, err := utils.SignRefreshToken(int64(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
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

func DeleteUser(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	err := models.DeleteUser(userID)
	if err != nil {
		log.Printf("DeleteUser: failed to delete user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user and all associated data deleted successfully"})
}
