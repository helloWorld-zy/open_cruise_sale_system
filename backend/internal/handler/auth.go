package handler

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/response"

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
	jwtConfig *config.JWTConfig
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtConfig *config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		jwtConfig: jwtConfig,
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

	// TODO: Validate credentials against database
	// For now, return mock response

	// Generate tokens
	token, err := middleware.GenerateToken("user-id", req.Username, "operations", h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken("user-id", h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	response.Success(c, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:       "user-id",
			Username: req.Username,
			Name:     "Admin User",
			Role:     "operations",
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

	// TODO: Validate refresh token and get user info

	token, err := middleware.GenerateToken("user-id", "username", "operations", h.jwtConfig)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	response.Success(c, LoginResponse{
		Token:        token,
		RefreshToken: "new-refresh-token",
		User: UserInfo{
			ID:       "user-id",
			Username: "username",
			Name:     "Admin User",
			Role:     "operations",
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
	// TODO: Add token to blacklist
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
		Name:     "Admin User",
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

	// TODO: Validate old password and update to new password
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

	// TODO: Exchange WeChat code for user info
	// Generate tokens
	token, _ := middleware.GenerateToken("wechat-user-id", "wechat-user", "customer", h.jwtConfig)
	refreshToken, _ := middleware.GenerateRefreshToken("wechat-user-id", h.jwtConfig)

	response.Success(c, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:       "wechat-user-id",
			Username: "wechat-user",
			Name:     "WeChat User",
			Role:     "customer",
		},
	})
}

// WeChatLoginRequest represents WeChat login request
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}
