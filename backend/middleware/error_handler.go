package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://gin-gonic.com/en/docs/examples/error-handling-middleware/
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			c.JSON(http.StatusInternalServerError, map[string]any{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}
