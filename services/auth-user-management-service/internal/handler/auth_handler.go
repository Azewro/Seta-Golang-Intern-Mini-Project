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

// Register (POST /api/v1/auth/register)
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

// VerifyEmail (GET /api/v1/auth/verify-email?token=...)
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

// ResendVerification (POST /api/v1/auth/resend-verification)
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

// Login (POST /api/v1/auth/login)
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

// Logout (POST /api/v1/auth/logout)
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

// ListUsers (GET /api/v1/users) manager only
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

// Me (GET /api/v1/users/me)
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

	me, err := h.authUsecase.GetMyProfile(userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, me)
}
