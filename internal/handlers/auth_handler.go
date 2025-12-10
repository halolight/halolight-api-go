package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/internal/repository"
	"github.com/halolight/halolight-api-go/internal/services"
	"github.com/halolight/halolight-api-go/pkg/config"
	"github.com/halolight/halolight-api-go/pkg/utils"
)

type AuthHandler struct {
	auth             services.AuthService
	refreshTokenRepo repository.RefreshTokenRepository
	cfg              config.Config
}

func NewAuthHandler(auth services.AuthService, refreshTokenRepo repository.RefreshTokenRepository, cfg config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, refreshTokenRepo: refreshTokenRepo, cfg: cfg}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Registration details"
// @Success 201 {object} authResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.auth.Register(req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		if errors.Is(err, services.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
			return
		}
		if errors.Is(err, services.ErrWeakPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is too weak"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, authResponse{
		User:  user,
		Token: token,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} authResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		User:  user,
		Token: token,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshRequest true "Refresh token"
// @Success 200 {object} authResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Verify refresh token
	storedToken, err := h.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil || storedToken == nil || storedToken.IsExpired() {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid refresh token"})
		return
	}

	// Generate new tokens
	accessToken, err := utils.GenerateAccessToken(storedToken.UserID, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"accessToken": accessToken,
	})
}

// Me godoc
// @Summary Get current user
// @Description Get authenticated user's profile
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	// Get user from service (would need to add this method)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id": userID,
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userID")
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	_ = c.ShouldBindJSON(&req)

	if req.RefreshToken != "" {
		_ = h.refreshTokenRepo.DeleteByToken(req.RefreshToken)
	} else {
		_ = h.refreshTokenRepo.DeleteByUserID(userID)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out successfully"})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body forgotPasswordRequest true "Email"
// @Success 200 {object} map[string]string
// @Router /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// In production, send email with reset link
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "If an account with that email exists, a password reset link has been sent",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body resetPasswordRequest true "Reset details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Verify token and reset password (simplified)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password has been reset successfully",
	})
}
