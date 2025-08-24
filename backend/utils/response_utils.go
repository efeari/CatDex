package response

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	OK      bool   `json:"ok"`
	Data    any    `json:"data,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Success with data
func Success(c *gin.Context, data any) {
	c.JSON(200, APIResponse{
		OK:   true,
		Data: data,
	})
}

// Failure with reason/message
func Fail(c *gin.Context, reason, message string) {
	c.JSON(200, APIResponse{
		OK:      false,
		Reason:  reason,
		Message: message,
	})
}
