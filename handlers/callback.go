package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shirini-backend/utils"
)

func GoogleCallback(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")

		userID := "google-user-" + code

		token, err := utils.GenerateToken(userID, secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate token",
			})
			return
		}

		redirectURL := "http://localhost:3000/auth/callback?token=" + token

		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}