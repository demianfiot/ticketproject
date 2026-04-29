package http

import (
	"fmt"

	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	fmt.Printf("DEBUG: Path: %s, Method: %s\n", c.Request.URL.Path, c.Request.Method) // ← Додайте
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		c.JSON(401, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(401, gin.H{"error": "invalid auth header"})
		return
	}
	ctx := c.Request.Context()
	userId, err := h.services.Authorization.ParseToken(ctx, headerParts[1])
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("DEBUG: UserID: %d", userId)
	c.Set(userCtx, userId)
	c.Next()
}
