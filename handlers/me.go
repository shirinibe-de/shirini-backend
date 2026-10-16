package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shirini-backend/config"
	"shirini-backend/database"
)

func Me(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		userID, err := strconv.Atoi(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			return
		}

		user, err := database.FindUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to load user",
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         strconv.Itoa(user.ID),
			"email":      user.Email,
			"name":       user.Name,
			"avatar_url": user.AvatarURL,
		})
	}
}
