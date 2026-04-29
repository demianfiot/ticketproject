package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/tickets/", handler.CreateTicket)
	router.GET("/tickets/:id", handler.GetTicketByID)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
