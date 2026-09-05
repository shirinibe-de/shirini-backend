package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"shirini-backend/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, secret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)

		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
