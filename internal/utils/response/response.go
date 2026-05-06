package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      any         `json:"data,omitempty"`
	Errors    any         `json:"errors,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func JSON(c *gin.Context, status int, code, message string, data any, errors any) {
	c.JSON(status, APIResponse{
		Success:   status < http.StatusBadRequest,
		Code:      code,
		Message:   message,
		Data:      data,
		Errors:    errors,
		RequestID: RequestIDFromContext(c),
	})
}

func Success(c *gin.Context, status int, code, message string, data any) {
	JSON(c, status, code, message, data, nil)
}

func Fail(c *gin.Context, status int, code, message string, errors any) {
	JSON(c, status, code, message, nil, errors)
}

func RequestIDFromContext(c *gin.Context) string {
	if val, ok := c.Get("request_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
