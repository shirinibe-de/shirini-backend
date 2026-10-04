package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"shirini-backend/config"
)

func GoogleLogin(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		clientID := cfg.Google.ClientID
		redirectURI := cfg.Google.RedirectURI

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
}