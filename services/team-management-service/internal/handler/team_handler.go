package handler

import (
	"errors"
	"net/http"
	"strconv"

	"team-management-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

// TeamHandler binds HTTP requests to team usecases.
type TeamHandler struct {
	teamUsecase usecase.TeamUsecase
}

func NewTeamHandler(teamUsecase usecase.TeamUsecase) *TeamHandler {
	return &TeamHandler{teamUsecase: teamUsecase}
}

// CreateTeam godoc
// @Summary Create team
// @Description Global manager only. Creator becomes main manager.
// @Tags teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body usecase.CreateTeamRequest true "Team name"
// @Success 201 {object} usecase.TeamResponse
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, role, _, ok := getActorContext(c)
	if !ok {
		return
	}

	var req usecase.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamUsecase.CreateTeam(userID, role, &req)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, team)
}

// ListMyTeams godoc
// @Summary List my teams
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Success 200 {object} teamsListBody
// @Failure 401 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/my [get]
func (h *TeamHandler) ListMyTeams(c *gin.Context) {
	userID, _, _, ok := getActorContext(c)
	if !ok {
		return
	}

	teams, err := h.teamUsecase.ListMyTeams(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch teams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// GetTeam godoc
// @Summary Get team
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Param teamId path int true "Team ID"
// @Success 200 {object} usecase.TeamResponse
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 404 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/{teamId} [get]
func (h *TeamHandler) GetTeam(c *gin.Context) {
	userID, _, _, ok := getActorContext(c)
	if !ok {
		return
	}

	teamID, parseErr := parseUintParam(c, "teamId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	team, err := h.teamUsecase.GetTeam(userID, teamID)
	if err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, team)
}

// AddMember godoc
// @Summary Add member
// @Tags teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param teamId path int true "Team ID"
// @Param body body usecase.TeamActionRequest true "Target user"
// @Success 200 {object} messageBody
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 404 {object} errorBody
// @Failure 409 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/{teamId}/members [post]
func (h *TeamHandler) AddMember(c *gin.Context) {
	h.handleTeamAction(c, h.teamUsecase.AddMember, "Member added")
}

// RemoveMember godoc
// @Summary Remove member
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Param teamId path int true "Team ID"
// @Param userId path int true "User ID"
// @Success 200 {object} messageBody
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 404 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/{teamId}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	userID, role, _, ok := getActorContext(c)
	if !ok {
		return
	}

	teamID, parseErr := parseUintParam(c, "teamId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	memberUserID, parseErr := parseUintParam(c, "userId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	if err := h.teamUsecase.RemoveMember(userID, role, teamID, memberUserID); err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// AddManager godoc
// @Summary Add manager
// @Tags teams
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param teamId path int true "Team ID"
// @Param body body usecase.TeamActionRequest true "Target user (must be global manager)"
// @Success 200 {object} messageBody
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 404 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/{teamId}/managers [post]
func (h *TeamHandler) AddManager(c *gin.Context) {
	h.handleTeamAction(c, h.teamUsecase.AddManager, "Manager added")
}

// RemoveManager godoc
// @Summary Remove manager
// @Tags teams
// @Security BearerAuth
// @Produce json
// @Param teamId path int true "Team ID"
// @Param userId path int true "User ID"
// @Success 200 {object} messageBody
// @Failure 400 {object} errorBody
// @Failure 403 {object} errorBody
// @Failure 404 {object} errorBody
// @Failure 500 {object} errorBody
// @Router /api/v1/teams/{teamId}/managers/{userId} [delete]
func (h *TeamHandler) RemoveManager(c *gin.Context) {
	userID, role, _, ok := getActorContext(c)
	if !ok {
		return
	}

	teamID, parseErr := parseUintParam(c, "teamId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}
	managerUserID, parseErr := parseUintParam(c, "userId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	if err := h.teamUsecase.RemoveManager(userID, role, teamID, managerUserID); err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed"})
}

func (h *TeamHandler) handleTeamAction(c *gin.Context, fn func(actorID uint, actorRole string, token string, teamID uint, targetUserID uint) error, successMessage string) {
	userID, role, token, ok := getActorContext(c)
	if !ok {
		return
	}

	teamID, parseErr := parseUintParam(c, "teamId")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	var req usecase.TeamActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := fn(userID, role, token, teamID, req.UserID); err != nil {
		h.handleUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": successMessage})
}

func (h *TeamHandler) handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrTeamNameRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrTeamNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrAlreadyInTeam):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrNotTeamMember):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrMainManagerRequired), errors.Is(err, usecase.ErrForbidden), errors.Is(err, usecase.ErrNotTeamManager):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrGlobalManagerRequired), errors.Is(err, usecase.ErrUseManagerEndpoint), errors.Is(err, usecase.ErrCannotRemoveMain):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	raw := c.Param(key)
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid path parameter")
	}
	return uint(parsed), nil
}

func getActorContext(c *gin.Context) (uint, string, string, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return 0, "", "", false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return 0, "", "", false
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role context"})
		return 0, "", "", false
	}

	role, ok := roleValue.(string)
	if !ok || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role context"})
		return 0, "", "", false
	}

	tokenValue, exists := c.Get("token")
	token := ""
	if exists {
		token, _ = tokenValue.(string)
	}

	return userID, role, token, true
}

// Swagger-only DTOs
type teamsListBody struct {
	Data []usecase.TeamResponse `json:"data"`
}

type messageBody struct {
	Message string `json:"message"`
}

type errorBody struct {
	Error string `json:"error"`
}
