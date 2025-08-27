package controllers

import (
	"fmt"
	"log"
	"net/http"
	"plant-reminder/dto"
	"plant-reminder/service"
	"plant-reminder/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	userService service.UserServiceInterface
}

func NewUserController(userService service.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Login(ctx *gin.Context) {
	var loginRequest dto.UserLoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := utils.Validate.Struct(loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := uc.userService.VerifyUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authResponse)
}

func (uc *UserController) SignUp(ctx *gin.Context) {
	var userRequest dto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := utils.Validate.Struct(userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := uc.userService.CreateUser(&userRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, authResponse)
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
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

func (uc *UserController) SetPushToken(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	var req dto.PushTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("SetPushToken: failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid push token"})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		log.Printf("SetPushToken: validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userService.SetPushToken(fmt.Sprintf("%d", userID), req.Token)
	if err != nil {
		log.Printf("SetPushToken: failed to set push token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "push token set successfully"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	err := uc.userService.DeleteUser(userID)
	if err != nil {
		log.Printf("DeleteUser: failed to delete user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user and all associated data deleted successfully"})
}

func (uc *UserController) GetMyProfile(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")
	userResponse, err := uc.userService.GetUser(userID)
	if err != nil {
		log.Printf("GetUser: failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userResponse == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": userResponse})
}
