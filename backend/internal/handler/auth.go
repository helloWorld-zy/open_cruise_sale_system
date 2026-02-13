package handler

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// AuthHandler handles authentication
type AuthHandler struct {
	jwtConfig      *config.JWTConfig
	userRepo       repository.UserRepository
	rbac           *auth.RBAC
	tokenBlacklist middleware.TokenBlacklist
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	jwtConfig *config.JWTConfig,
	userRepo repository.UserRepository,
	rbac *auth.RBAC,
	tokenBlacklist middleware.TokenBlacklist,
) *AuthHandler {
	return &AuthHandler{
		jwtConfig:      jwtConfig,
		userRepo:       userRepo,
		rbac:           rbac,
		tokenBlacklist: tokenBlacklist,
	}
}

// Login godoc
// @Summary Staff login
// @Description Authenticate staff and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// SEC-001: Look up user by username (phone/email/nickname)
	user, err := h.userRepo.GetByUsername(c.Request.Context(), req.Username)
	if err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// Verify password using bcrypt
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// Determine user role (default to customer if not set)
	role := "customer"
	if user.Nickname != "" {
		// In a real system, roles would be stored in a separate table
		// For now, use a simplified role lookup
		role = "operations"
	}

	// Generate tokens
	token, err := middleware.GenerateToken(user.ID.String(), user.Nickname, role, h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID.String(), h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	response.Success(c, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:       user.ID.String(),
			Username: user.Nickname,
			Name:     user.Nickname,
			Role:     role,
		},
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer refresh_token"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(c, "Authorization header required")
		return
	}

	// Extract Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		response.Unauthorized(c, "Invalid authorization header format")
		return
	}

	// SEC-001: Parse and validate the refresh token
	claims, err := middleware.ParseToken(parts[1], h.jwtConfig)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	// Verify it's a refresh token
	if claims.TokenType != "refresh" {
		response.Unauthorized(c, "Invalid token type, expected refresh token")
		return
	}

	// Check if token is blacklisted
	if h.tokenBlacklist != nil && h.tokenBlacklist.IsRevoked(c.Request.Context(), claims.ID) {
		response.Unauthorized(c, "Token has been revoked")
		return
	}

	// Look up user to verify they still exist and are active
	user, err := h.userRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Unauthorized(c, "User not found")
		return
	}

	// Generate new token pair
	token, err := middleware.GenerateToken(user.ID.String(), claims.Username, claims.Role, h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID.String(), h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	response.Success(c, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:       user.ID.String(),
			Username: user.Nickname,
			Name:     user.Nickname,
			Role:     claims.Role,
		},
	})
}

// Logout godoc
// @Summary Logout
// @Description Invalidate current token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// SEC-002: Add token to blacklist
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := middleware.ParseToken(parts[1], h.jwtConfig)
			if err == nil && h.tokenBlacklist != nil {
				// Add token to blacklist with remaining TTL
				remaining := time.Until(claims.ExpiresAt.Time)
				if remaining > 0 {
					h.tokenBlacklist.Revoke(c.Request.Context(), claims.ID, remaining)
				}
			}
		}
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

// Me godoc
// @Summary Get current user
// @Description Get information about currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=UserInfo}
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	response.Success(c, UserInfo{
		ID:       userID.(string),
		Username: username.(string),
		Name:     username.(string),
		Role:     role.(string),
	})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change current user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 200 {object} response.Response
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// SEC-001: Get current user
	userID, _ := c.Get("userID")
	user, err := h.userRepo.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify old password
	if err := auth.CheckPassword(req.OldPassword, user.PasswordHash); err != nil {
		response.Unauthorized(c, "Old password is incorrect")
		return
	}

	// Hash and save new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		response.InternalServerError(c, "Failed to hash password")
		return
	}

	user.PasswordHash = hashedPassword
	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		response.InternalServerError(c, "Failed to update password")
		return
	}

	response.Success(c, gin.H{"message": "Password changed successfully"})
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// WeChatLogin godoc
// @Summary WeChat Mini Program login
// @Description Login via WeChat code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body WeChatLoginRequest true "WeChat login request"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Router /auth/wechat-login [post]
func (h *AuthHandler) WeChatLogin(c *gin.Context) {
	var req WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// TODO: Exchange WeChat code for user info via WeChat SDK
	// This requires WeChat Mini Program AppID and AppSecret
	// 1. Call wx.login() API to get session_key and openid
	// 2. Create or find user by openid
	// 3. Generate JWT tokens
	response.Error(c, http.StatusNotImplemented, "WeChat login requires WeChat SDK integration - configure AppID and AppSecret first")
}

// WeChatLoginRequest represents WeChat login request
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}
