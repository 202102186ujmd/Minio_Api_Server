package middleware

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func requestIDFromContext(c *gin.Context) string {
	if val, ok := c.Get("request_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
