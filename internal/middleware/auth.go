package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/pkg/config"
	"github.com/halolight/halolight-api-go/pkg/utils"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			return
		}

		// Parse Bearer token
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, expected 'Bearer <token>'",
			})
			return
		}

		// Parse and validate JWT
		claims, err := utils.ValidateAccessToken(parts[1], cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}
