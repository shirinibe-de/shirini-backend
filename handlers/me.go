package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"shirini-backend/utils"
)

func Me(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(401, gin.H{"error": "no token"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": claims["user_id"],
	})
}