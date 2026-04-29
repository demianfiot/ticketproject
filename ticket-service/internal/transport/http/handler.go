package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/demianfiot/ticketproject/ticket-service/internal/metrics"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) CreateTicket(c *gin.Context) {
	start := time.Now()

	metrics.HTTPRequests.WithLabelValues("POST", "/tickets").Inc()

	defer func() {
		metrics.HTTPDuration.
			WithLabelValues("/tickets").
			Observe(time.Since(start).Seconds())
	}()
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	resp, err := h.services.Ticket.CreateTicket(c.Request.Context(), service.CreateTicketInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetTicketByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ticket id",
		})
		return
	}

	resp, err := h.services.Ticket.GetTicketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get ticket",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
