
package middleware

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"shirini-backend/database"
)

var uuidPattern = regexp.MustCompile(
	`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
)

func RequireTeamMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		teamID := c.Param("team_id")

		// Team IDs in the database are UUIDs.
		if !uuidPattern.MatchString(teamID) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "team not found",
			})
			c.Abort()
			return
		}

		rawUserID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			c.Abort()
			return
		}

		userIDString, ok := rawUserID.(string)
		if !ok || userIDString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDString)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user",
			})
			c.Abort()
			return
		}

		team, err := database.FindTeamByID(teamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check team",
			})
			c.Abort()
			return
		}

		if team == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "team not found",
			})
			c.Abort()
			return
		}

		membership, err := database.FindTeamMembership(teamID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to check team membership",
			})
			c.Abort()
			return
		}

		if membership == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "not a member of this team",
			})
			c.Abort()
			return
		}

		c.Set("team", team)
		c.Set("team_membership", membership)

		c.Next()
	}
}
