package middleware

import (
	"time"

	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logg *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)

		if raw != "" {
			path = path + "?" + raw
		}

		logg.Info("request",
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.Int("status", c.Writer.Status()),
			logger.String("client_ip", c.ClientIP()),
			logger.String("request_id", requestIDFromContext(c)),
			logger.Int64("latency_ms", latency.Milliseconds()),
		)
	}
}
