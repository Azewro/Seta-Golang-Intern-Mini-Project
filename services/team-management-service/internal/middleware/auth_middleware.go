package middleware

import (
	"net/http"
	"os"
	"strings"

	"team-management-service/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates JWT and sets user context.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			c.Abort()
			return
		}

		jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// ManagerOnly restricts access to manager users.
func ManagerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "manager" {
			c.JSON(http.StatusForbidden, gin.H{"error": "manager role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
