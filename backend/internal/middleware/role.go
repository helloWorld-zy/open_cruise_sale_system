package middleware

import (
	"backend/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		if userRole == "" {
			response.Unauthorized(c, "Role not found in context")
			c.Abort()
			return
		}

		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Insufficient permissions")
		c.Abort()
	}
}
