package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"auth-user-management-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

// AuthHandler holds the dependencies for HTTP routing
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler injects usecase layer into handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register creates a new user account.
// @Summary Register
// @Description Create account (role manager or member; default member). Email verification flow depends on SMTP configuration.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body usecase.RegisterRequest true "Credentials"
// @Success 201 {object} registerCreatedResponse
// @Failure 400 {object} errorJSON
// @Failure 409 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.authUsecase.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrSMTPNotConfigured):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email service is not configured"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please verify your email within 5 minutes.",
		"user":    newUser,
	})
}

// VerifyEmail marks a user verified using a one-time token from email.
// @Summary Verify email
// @Tags auth
// @Produce json
// @Param token query string true "Verification token from email link"
// @Success 200 {object} messageJSON
// @Failure 400 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	rawToken := strings.TrimSpace(c.Query("token"))
	if rawToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token is required"})
		return
	}

	err := h.authUsecase.VerifyEmail(rawToken)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidVerifyToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrExpiredVerifyToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully. You can now log in."})
}

// ResendVerification emails a new verification link if the account exists and is unverified.
// @Summary Resend verification email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body usecase.ResendVerificationRequest true "Email"
// @Success 200 {object} messageJSON
// @Failure 400 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req usecase.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authUsecase.ResendVerification(&req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrSMTPNotConfigured):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email service is not configured"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resend verification email"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the account exists and is unverified, a new verification email has been sent."})
}

// Login returns a JWT and user profile.
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body usecase.LoginRequest true "Credentials"
// @Success 200 {object} usecase.LoginResponse
// @Failure 400 {object} errorJSON
// @Failure 401 {object} errorJSON
// @Failure 403 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := h.authUsecase.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrJWTSecretMissing):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

// Logout revokes the current session (JWT must still be valid).
// @Summary Logout
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} messageJSON
// @Failure 401 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenIDValue, exists := c.Get("tokenID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token context"})
		return
	}

	tokenID, ok := tokenIDValue.(string)
	if !ok || tokenID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token context"})
		return
	}

	if err := h.authUsecase.Logout(tokenID); err != nil {
		if errors.Is(err, usecase.ErrSessionUnavailable) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me returns the authenticated user profile.
// @Summary Current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} usecase.UserResponse
// @Failure 401 {object} errorJSON
// @Failure 404 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/users/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return
	}

	profile, err := h.authUsecase.GetMyProfile(userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// BulkGetUsers resolves user summaries for a list of IDs.
// @Summary Bulk get users
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body usecase.BulkGetUsersRequest true "User IDs"
// @Success 200 {object} bulkUsersResponse
// @Failure 400 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/users/bulk [post]
func (h *AuthHandler) BulkGetUsers(c *gin.Context) {
	var req usecase.BulkGetUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.authUsecase.GetUsersByIDs(req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// ListUsers returns a paginated list (global manager only).
// @Summary List users
// @Description Requires global role manager.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Page size" default(20)
// @Success 200 {object} listUsersResponse
// @Failure 403 {object} errorJSON
// @Failure 500 {object} errorJSON
// @Router /api/v1/users [get]
func (h *AuthHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, err := h.authUsecase.ListUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  users,
	})
}

// Swagger response helpers (referenced by comments only).
type registerCreatedResponse struct {
	Message string               `json:"message"`
	User    usecase.UserResponse `json:"user"`
}

type messageJSON struct {
	Message string `json:"message"`
}

type errorJSON struct {
	Error string `json:"error"`
}

type bulkUsersResponse struct {
	Data []usecase.UserResponse `json:"data"`
}

type listUsersResponse struct {
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
	Data  []usecase.UserResponse `json:"data"`
}
