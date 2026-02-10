package middleware

import (
	"backend/internal/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			response.InternalServerError(c, err.Error())
		}
	}
}

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
