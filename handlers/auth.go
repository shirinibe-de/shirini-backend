package handlers

import (
	"crypto/rand"
	"encoding/base64"
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

		// The frontend URL to return to after auth; falls back to
		// the configured frontend callback URL.
		frontendRedirect := c.Query("redirect_uri")
		if frontendRedirect == "" {
			frontendRedirect = cfg.Google.FrontendCallbackURL
		}

		nonceBytes := make([]byte, 16)
		if _, err := rand.Read(nonceBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate state",
			})
			return
		}

		nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)

		state, err := signState(&oauthState{
			Nonce:       nonce,
			RedirectURI: frontendRedirect,
		}, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate state",
			})
			return
		}

		setStateCookie(c, nonce)

		authURL := "https://accounts.google.com/o/oauth2/v2/auth" +
			"?client_id=" + url.QueryEscape(clientID) +
			"&redirect_uri=" + url.QueryEscape(redirectURI) +
			"&response_type=code" +
			"&scope=openid%20email%20profile" +
			"&state=" + url.QueryEscape(state)

		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
