package middleware

import (
	"net/http"

	response "github.com/efear/catdex/utils"
	"github.com/gin-gonic/gin"
)

// https://gin-gonic.com/en/docs/examples/error-handling-middleware/
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				OK:      false,
				Data:    nil,
				Reason:  "server_error",
				Message: err.Error(),
			})
		}
	}
}
