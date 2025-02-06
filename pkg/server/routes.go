package server

import (
	"github.com/gin-gonic/gin"
	"github.com/loanengine/internal/handler/rest"
)

func (a *App) SetupRoutesAndMiddleware(router *gin.RouterGroup, restHandler rest.RestHandler) {
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	loan := v1.Group("/loan")
	loan.POST("/propose", restHandler.Propose)
	loan.POST("/invest", restHandler.Invest)
	loan.POST("/approve", restHandler.Approve)
	loan.POST("/disburse", restHandler.Disburse)

	agent := v1.Group("/agent")
	agent.GET("/loans", restHandler.GetProposedLoans)

	admin := v1.Group("/admin")
	admin.GET("/loan/status", restHandler.GetState)
}
