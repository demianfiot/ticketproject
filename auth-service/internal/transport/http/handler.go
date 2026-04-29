package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/demianfiot/ticketproject/auth-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var input RegisterUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input body",
		})
		return
	}

	id, err := h.services.Authorization.CreateUser(c.Request.Context(), service.CreateUserInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, RegisterUserResponse{
		ID: id,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var input LoginUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input body",
		})
		return
	}

	token, err := h.services.Authorization.GenerateToken(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to login user",
		})
		return
	}

	c.JSON(http.StatusOK, LoginUserResponse{
		Token: token,
	})
}
