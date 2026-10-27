package handlers

// Mock implementations of the domain endpoints. These mirror the seed data
// and behavior of the frontend mock (shirini-webapp/lib/api/mock) so the
// webapp can run against the real backend before the domain logic is
// implemented. Responses are held in memory and reset on server restart.

import (
	"math/rand/v2"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"shirini-backend/database"

	"github.com/gin-gonic/gin"
)

// --- Types mirroring shirini-webapp/lib/api/types.ts ---

type mockUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}

type mockTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Currency    string `json:"currency"`
	InviteToken string `json:"invite_token"`
	CreatedAt   string `json:"created_at"`
}

type mockTeamMember struct {
	mockUser
	JoinedAt string `json:"joined_at"`
}

type mockAchievementType struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type mockVoteSummary struct {
	Approve int     `json:"approve"`
	Reject  int     `json:"reject"`
	MyVote  *string `json:"my_vote"`
}

type mockClaim struct {
	ID              string              `json:"id"`
	Claimer         mockUser            `json:"claimer"`
	ClaimedUser     mockUser            `json:"claimed_user"`
	AchievementType mockAchievementType `json:"achievement_type"`
	Quantity        int                 `json:"quantity"`
	Status          string              `json:"status"`
	Votes           mockVoteSummary     `json:"votes"`
	ExpiresAt       string              `json:"expires_at"`
	CreatedAt       string              `json:"created_at"`
}

type mockSettlement struct {
	DebtorID     string `json:"debtor_id"`
	PreviousDebt int    `json:"previous_debt"`
	SettledAt    string `json:"settled_at"`
}

// --- In-memory store mirroring mock/db.ts ---

// The frontend mock treats user-1 as the current user for every
// authenticated request.
const mockCurrentUserID = "user-1"

var mockStartTime = time.Now()

func mockHoursFromNow(h float64) string {
	return time.Now().Add(time.Duration(h * float64(time.Hour))).UTC().Format(time.RFC3339)
}

type mockStore struct {
	mu               sync.Mutex
	users            map[string]mockUser
	teams            map[string]mockTeam
	teamMembers      map[string][]string
	achievementTypes map[string][]mockAchievementType
	claims           []mockClaim
	votes            map[string]map[string]string
	settlements      map[string]map[string]string
}

func newMockStore() *mockStore {
	users := map[string]mockUser{
		"user-1": {
			ID:        "user-1",
			Email:     "ahmad@example.com",
			Name:      "Ahmad",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=Ahmad",
			CreatedAt: "2024-01-01T00:00:00Z",
		},
		"user-2": {
			ID:        "user-2",
			Email:     "sara@example.com",
			Name:      "Sara",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=Sara",
			CreatedAt: "2024-01-02T00:00:00Z",
		},
		"user-3": {
			ID:        "user-3",
			Email:     "reza@example.com",
			Name:      "Reza",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=Reza",
			CreatedAt: "2024-01-03T00:00:00Z",
		},
		"user-4": {
			ID:        "user-4",
			Email:     "neda@example.com",
			Name:      "Neda",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=Neda",
			CreatedAt: "2024-01-04T00:00:00Z",
		},
		"user-5": {
			ID:        "user-5",
			Email:     "ali@example.com",
			Name:      "Ali",
			AvatarURL: "https://api.dicebear.com/7.x/avataaars/svg?seed=Ali",
			CreatedAt: "2024-01-05T00:00:00Z",
		},
	}

	teams := map[string]mockTeam{
		"team-1": {
			ID:          "team-1",
			Name:        "Backend Squad",
			Currency:    "slice of pizza",
			InviteToken: "invite-backend-squad",
			CreatedAt:   "2024-01-10T00:00:00Z",
		},
		"team-2": {
			ID:          "team-2",
			Name:        "Frontend Ninjas",
			Currency:    "cup of coffee",
			InviteToken: "invite-frontend-ninjas",
			CreatedAt:   "2024-01-15T00:00:00Z",
		},
	}

	teamMembers := map[string][]string{
		"team-1": {"user-1", "user-2", "user-3", "user-5"},
		"team-2": {"user-1", "user-4"},
	}

	achievementTypes := map[string][]mockAchievementType{
		"team-1": {
			{ID: "ach-1", Name: "First deploy to production", CreatedAt: "2024-01-10T10:00:00Z"},
			{ID: "ach-2", Name: "Bought a new phone", CreatedAt: "2024-01-12T10:00:00Z"},
			{ID: "ach-3", Name: "Fixed a critical bug", CreatedAt: "2024-01-14T10:00:00Z"},
		},
		"team-2": {
			{ID: "ach-4", Name: "Shipped a new feature", CreatedAt: "2024-01-15T10:00:00Z"},
			{ID: "ach-5", Name: "Promoted to Senior", CreatedAt: "2024-01-16T10:00:00Z"},
		},
	}

	approved := "approve"
	claims := []mockClaim{
		{
			ID:              "claim-1",
			Claimer:         users["user-2"],
			ClaimedUser:     users["user-1"],
			AchievementType: achievementTypes["team-1"][0],
			Quantity:        2,
			Status:          "pending",
			Votes:           mockVoteSummary{Approve: 2, Reject: 0},
			ExpiresAt:       mockHoursFromNow(24),
			CreatedAt:       mockHoursFromNow(-2),
		},
		{
			ID:              "claim-2",
			Claimer:         users["user-1"],
			ClaimedUser:     users["user-3"],
			AchievementType: achievementTypes["team-1"][1],
			Quantity:        1,
			Status:          "pending",
			Votes:           mockVoteSummary{Approve: 1, Reject: 1},
			ExpiresAt:       mockHoursFromNow(20),
			CreatedAt:       mockHoursFromNow(-4),
		},
		{
			ID:              "claim-3",
			Claimer:         users["user-3"],
			ClaimedUser:     users["user-5"],
			AchievementType: achievementTypes["team-1"][2],
			Quantity:        3,
			Status:          "approved",
			Votes:           mockVoteSummary{Approve: 3, Reject: 0, MyVote: &approved},
			ExpiresAt:       mockHoursFromNow(22),
			CreatedAt:       mockHoursFromNow(-1),
		},
	}

	votes := map[string]map[string]string{
		"claim-1": {
			"user-2": "approve",
			"user-3": "approve",
		},
		"claim-2": {
			"user-2": "approve",
			"user-5": "reject",
		},
		"claim-3": {
			"user-1": "approve",
			"user-2": "approve",
			"user-3": "approve",
		},
	}

	return &mockStore{
		users:            users,
		teams:            teams,
		teamMembers:      teamMembers,
		achievementTypes: achievementTypes,
		claims:           claims,
		votes:            votes,
		settlements:      map[string]map[string]string{"team-1": {}},
	}
}

var store = newMockStore()

// mockID generates IDs like the frontend mock:
// "<prefix>-" + Math.random().toString(36).substring(7)
func mockID(prefix string) string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 7)
	for i := range b {
		b[i] = chars[rand.IntN(len(chars))]
	}
	return prefix + "-" + string(b)
}

// calculateDebt mirrors mock/db.ts calculateDebt: sum of quantities of the
// user's approved claims created after their last settlement.
func (s *mockStore) calculateDebt(teamID, userID string) int {
	lastSettlement := s.settlements[teamID][userID]
	debt := 0
	for _, claim := range s.claims {
		if claim.ClaimedUser.ID != userID || claim.Status != "approved" {
			continue
		}
		if lastSettlement == "" || claim.CreatedAt > lastSettlement {
			debt += claim.Quantity
		}
	}
	return debt
}

// --- Handlers mirroring mock/handlers.ts ---

func mockError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// GET /auth/me (mock): the frontend mock ignores auth and returns user-1.
func MockMe(c *gin.Context) {
	store.mu.Lock()
	defer store.mu.Unlock()
	c.JSON(http.StatusOK, store.users[mockCurrentUserID])
}

// GET /teams
func MockGetTeams(c *gin.Context) {
	store.mu.Lock()
	defer store.mu.Unlock()

	teamIDs := make([]string, 0)
	for id, members := range store.teamMembers {
		for _, m := range members {
			if m == mockCurrentUserID {
				teamIDs = append(teamIDs, id)
				break
			}
		}
	}
	sort.Strings(teamIDs)

	teams := make([]mockTeam, 0, len(teamIDs))
	for _, id := range teamIDs {
		teams = append(teams, store.teams[id])
	}
	c.JSON(http.StatusOK, teams)
}

type mockCreateTeamBody struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

// POST /teams

func MockCreateTeam(c *gin.Context) {
	var body mockCreateTeamBody
	if err := c.ShouldBindJSON(&body); err != nil {
		mockError(c, http.StatusBadRequest, "Name and currency are required")
		return
	}

	rawUserID, exists := c.Get("user_id")
	if !exists {
		mockError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	userIDString, ok := rawUserID.(string)
	if !ok || userIDString == "" {
		mockError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil || userID <= 0 {
		mockError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	inviteToken := mockID("invite")

	dbTeam, err := database.CreateTeam(
		body.Name,
		body.Currency,
		inviteToken,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create team",
		})
		return
	}

	team := mockTeam{
		ID:          dbTeam.ID,
		Name:        dbTeam.Name,
		Currency:    body.Currency,
		InviteToken: inviteToken,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Temporary compatibility with the remaining mock handlers.
	store.mu.Lock()
	store.teams[team.ID] = team
	store.teamMembers[team.ID] = []string{mockCurrentUserID}
	store.achievementTypes[team.ID] = []mockAchievementType{}
	store.settlements[team.ID] = map[string]string{}
	store.mu.Unlock()

	c.JSON(http.StatusOK, team)
}

// GET /teams/:team_id
func MockGetTeamDetails(c *gin.Context) {
	teamID := c.Param("team_id")

	store.mu.Lock()
	defer store.mu.Unlock()

	team, ok := store.teams[teamID]
	if !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           team.ID,
		"name":         team.Name,
		"currency":     team.Currency,
		"invite_token": team.InviteToken,
		"created_at":   team.CreatedAt,
		"member_count": len(store.teamMembers[teamID]),
	})
}

// GET /teams/:team_id/members
func MockGetTeamMembers(c *gin.Context) {
	teamID := c.Param("team_id")

	store.mu.Lock()
	defer store.mu.Unlock()

	memberIDs, ok := store.teamMembers[teamID]
	if !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}

	members := make([]mockTeamMember, 0, len(memberIDs))
	for _, id := range memberIDs {
		members = append(members, mockTeamMember{
			mockUser: store.users[id],
			JoinedAt: "2024-01-01T00:00:00Z",
		})
	}
	c.JSON(http.StatusOK, members)
}

// GET /invites/:token (no auth, per api.md)
func MockResolveInvite(c *gin.Context) {
	token := c.Param("token")

	store.mu.Lock()
	defer store.mu.Unlock()

	for _, team := range store.teams {
		if team.InviteToken == token {
			c.JSON(http.StatusOK, gin.H{
				"team_id":   team.ID,
				"team_name": team.Name,
				"currency":  team.Currency,
			})
			return
		}
	}
	mockError(c, http.StatusNotFound, "Invite not found")
}

// POST /invites/:token/accept
func MockAcceptInvite(c *gin.Context) {
	token := c.Param("token")

	store.mu.Lock()
	defer store.mu.Unlock()

	var team *mockTeam
	for id := range store.teams {
		if store.teams[id].InviteToken == token {
			t := store.teams[id]
			team = &t
			break
		}
	}
	if team == nil {
		mockError(c, http.StatusNotFound, "Invite not found")
		return
	}

	members := store.teamMembers[team.ID]
	isMember := false
	for _, m := range members {
		if m == mockCurrentUserID {
			isMember = true
			break
		}
	}
	if !isMember {
		store.teamMembers[team.ID] = append(members, mockCurrentUserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"team_id":   team.ID,
		"team_name": team.Name,
	})
}

// GET /teams/:team_id/achievement-types
func MockGetAchievementTypes(c *gin.Context) {
	teamID := c.Param("team_id")

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.teams[teamID]; !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}
	types := store.achievementTypes[teamID]
	if types == nil {
		types = []mockAchievementType{}
	}
	c.JSON(http.StatusOK, types)
}

type mockCreateAchievementTypeBody struct {
	Name string `json:"name" binding:"required"`
}

// POST /teams/:team_id/achievement-types
func MockCreateAchievementType(c *gin.Context) {
	teamID := c.Param("team_id")

	var body mockCreateAchievementTypeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		mockError(c, http.StatusBadRequest, "Name is required")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.teams[teamID]; !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}

	ach := mockAchievementType{
		ID:        mockID("ach"),
		Name:      body.Name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	store.achievementTypes[teamID] = append(store.achievementTypes[teamID], ach)
	c.JSON(http.StatusOK, ach)
}

// GET /teams/:team_id/claims?status=pending
func MockGetClaims(c *gin.Context) {
	teamID := c.Param("team_id")
	status := c.DefaultQuery("status", "pending")

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.teams[teamID]; !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}

	claims := make([]mockClaim, 0)
	for _, claim := range store.claims {
		if claim.Status == status {
			claims = append(claims, claim)
		}
	}
	c.JSON(http.StatusOK, claims)
}

type mockCreateClaimBody struct {
	ClaimedUserID     string `json:"claimed_user_id"`
	AchievementTypeID string `json:"achievement_type_id"`
	Quantity          int    `json:"quantity"`
}

// POST /teams/:team_id/claims
func MockCreateClaim(c *gin.Context) {
	teamID := c.Param("team_id")

	var body mockCreateClaimBody
	if err := c.ShouldBindJSON(&body); err != nil {
		mockError(c, http.StatusBadRequest, "Invalid claim body")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.teams[teamID]; !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}
	if body.ClaimedUserID == mockCurrentUserID {
		mockError(c, http.StatusBadRequest, "Cannot claim yourself")
		return
	}

	claimer, claimerOK := store.users[mockCurrentUserID]
	claimedUser, claimedOK := store.users[body.ClaimedUserID]
	var achievementType *mockAchievementType
	for i := range store.achievementTypes[teamID] {
		if store.achievementTypes[teamID][i].ID == body.AchievementTypeID {
			achievementType = &store.achievementTypes[teamID][i]
			break
		}
	}
	if !claimerOK || !claimedOK || achievementType == nil {
		mockError(c, http.StatusNotFound, "User or achievement type not found")
		return
	}

	claim := mockClaim{
		ID:              mockID("claim"),
		Claimer:         claimer,
		ClaimedUser:     claimedUser,
		AchievementType: *achievementType,
		Quantity:        body.Quantity,
		Status:          "pending",
		Votes:           mockVoteSummary{},
		ExpiresAt:       mockHoursFromNow(24),
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
	}
	store.claims = append(store.claims, claim)
	c.JSON(http.StatusOK, claim)
}

type mockCastVoteBody struct {
	Value string `json:"value" binding:"required"`
}

// POST /claims/:claim_id/votes
func MockCastVote(c *gin.Context) {
	claimID := c.Param("claim_id")

	var body mockCastVoteBody
	if err := c.ShouldBindJSON(&body); err != nil || (body.Value != "approve" && body.Value != "reject") {
		mockError(c, http.StatusBadRequest, "Value must be approve or reject")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	var claim *mockClaim
	for i := range store.claims {
		if store.claims[i].ID == claimID {
			claim = &store.claims[i]
			break
		}
	}
	if claim == nil {
		mockError(c, http.StatusNotFound, "Claim not found")
		return
	}
	if claim.Status != "pending" {
		mockError(c, http.StatusBadRequest, "Claim is not pending")
		return
	}

	if store.votes[claimID] == nil {
		store.votes[claimID] = map[string]string{}
	}
	store.votes[claimID][mockCurrentUserID] = body.Value

	approveCount, rejectCount := 0, 0
	for _, v := range store.votes[claimID] {
		if v == "approve" {
			approveCount++
		} else if v == "reject" {
			rejectCount++
		}
	}

	myVote := body.Value
	claim.Votes = mockVoteSummary{Approve: approveCount, Reject: rejectCount, MyVote: &myVote}

	// Frontend mock: claim resolves to approved once approve > reject.
	if approveCount > rejectCount {
		claim.Status = "approved"
	}

	c.JSON(http.StatusOK, gin.H{
		"claim_id":     claimID,
		"value":        body.Value,
		"claim_status": claim.Status,
	})
}

// GET /teams/:team_id/debts
func MockGetDebts(c *gin.Context) {
	teamID := c.Param("team_id")

	store.mu.Lock()
	defer store.mu.Unlock()

	team, ok := store.teams[teamID]
	if !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}

	memberIDs := store.teamMembers[teamID]
	type debtEntry struct {
		User      mockUser `json:"user"`
		TotalDebt int      `json:"total_debt"`
		Currency  string   `json:"currency"`
	}
	debts := make([]debtEntry, 0, len(memberIDs))
	for _, userID := range memberIDs {
		debts = append(debts, debtEntry{
			User:      store.users[userID],
			TotalDebt: store.calculateDebt(teamID, userID),
			Currency:  team.Currency,
		})
	}
	sort.SliceStable(debts, func(i, j int) bool {
		return debts[i].TotalDebt > debts[j].TotalDebt
	})
	c.JSON(http.StatusOK, debts)
}

type mockCreateSettlementBody struct {
	DebtorID string `json:"debtor_id" binding:"required"`
}

// POST /teams/:team_id/settlements
func MockCreateSettlement(c *gin.Context) {
	teamID := c.Param("team_id")

	var body mockCreateSettlementBody
	if err := c.ShouldBindJSON(&body); err != nil {
		mockError(c, http.StatusBadRequest, "debtor_id is required")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.teams[teamID]; !ok {
		mockError(c, http.StatusNotFound, "Team not found")
		return
	}

	previousDebt := store.calculateDebt(teamID, body.DebtorID)
	settledAt := time.Now().UTC().Format(time.RFC3339)
	if store.settlements[teamID] == nil {
		store.settlements[teamID] = map[string]string{}
	}
	store.settlements[teamID][body.DebtorID] = settledAt

	c.JSON(http.StatusOK, mockSettlement{
		DebtorID:     body.DebtorID,
		PreviousDebt: previousDebt,
		SettledAt:    settledAt,
	})
}
