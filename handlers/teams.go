package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"shirini-backend/database"
)

func currentUserID(c *gin.Context) (int, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return 0, false
	}

	userIDStr, ok := value.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return 0, false
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return 0, false
	}

	return userID, true
}

func GetTeams(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	teams, err := database.ListTeamsForUser(userID)
	if err != nil {
		log.Printf("GetTeams failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func CreateTeam(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var request struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Currency = strings.TrimSpace(request.Currency)

	if request.Name == "" || request.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and currency are required"})
		return
	}

	team, err := database.CreateTeam(request.Name, request.Currency, userID)
	if err != nil {
		log.Printf("CreateTeam failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create team"})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func GetTeamDetails(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	teamID := strings.TrimSpace(c.Param("team_id"))
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	team, err := database.GetTeamDetailsForUser(teamID, userID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrTeamNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		case errors.Is(err, database.ErrNotTeamMember):
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this team"})
		default:
			log.Printf("GetTeamDetails failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load team"})
		}
		return
	}

	c.JSON(http.StatusOK, team)
}
func GetTeamMembers(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	teamID := strings.TrimSpace(c.Param("team_id"))
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	members, err := database.GetTeamMembersForUser(teamID, userID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotTeamMember):
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not a member of this team"})
		default:
			log.Printf("GetTeamMembers failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load team members"})
		}
		return
	}

	c.JSON(http.StatusOK, members)
}