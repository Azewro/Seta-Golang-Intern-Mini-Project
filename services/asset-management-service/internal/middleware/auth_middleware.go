package middleware

import (
	"net/http"
	"strings"

	"asset-management-service/pkg/client"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates token via auth service and sets context.
func AuthRequired(authClient client.AuthClient) gin.HandlerFunc {
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

		user, err := authClient.VerifyToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Set("role", user.Role)
		c.Set("token", parts[1])
		c.Next()
	}
}
