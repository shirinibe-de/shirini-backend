package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shirini-backend/config"
	"shirini-backend/database"
	"shirini-backend/services"
	"shirini-backend/utils"
)

func GoogleCallback(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		code := c.Query("code")

		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "authorization code is required",
			})
			return
		}

		// Verify the OAuth state: signature must be valid and the
		// nonce must match the browser cookie set at login start.
		stateParam := c.Query("state")

		state, err := verifyState(stateParam, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid oauth state",
			})
			return
		}

		cookieNonce, err := c.Cookie(stateCookieName)
		if err != nil || cookieNonce == "" || cookieNonce != state.Nonce {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid oauth state",
			})
			return
		}

		clearStateCookie(c)

		googleUser, err := services.ExchangeCode(code, cfg)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to authenticate with Google",
			})
			return
		}

		user, err := database.FindUserByGoogleID(googleUser.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to find user",
			})
			return
		}

		if user == nil {

			user, err = database.CreateUser(&database.User{
				GoogleID:  googleUser.ID,
				Email:     googleUser.Email,
				Name:      googleUser.Name,
				AvatarURL: googleUser.Picture,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to create user",
				})
				return
			}
		}

		token, err := utils.GenerateToken(
			strconv.Itoa(user.ID),
			cfg.JWT.Secret,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate token",
			})
			return
		}

		redirectURL := state.RedirectURI + "?token=" + token

		c.Redirect(
			http.StatusTemporaryRedirect,
			redirectURL,
		)
	}
}
