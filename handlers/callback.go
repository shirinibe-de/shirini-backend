package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shirini-backend/utils"
)

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	

	userID := "google-user-" + code

	token, _ := utils.GenerateToken(userID)

	redirectURL := "http://localhost:3000/auth/callback?token=" + token

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}