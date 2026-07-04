package handlers

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Google OAuth is not configured",
		})
		return
	}

	authURL := "https://accounts.google.com/o/oauth2/v2/auth" +
		"?client_id=" + url.QueryEscape(clientID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code" +
		"&scope=openid%20email%20profile"

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}