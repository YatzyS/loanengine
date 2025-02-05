package server

import (
	"github.com/gin-gonic/gin"
	"github.com/loanengine/internal/handler/rest"
)

func (a *App) SetupRoutesAndMiddleware(router *gin.RouterGroup, restHandler rest.RestHandler) {
	router.Use(gin.Recovery())
	// TODO: Setup routes here
}
