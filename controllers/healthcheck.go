package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "pong")
}
