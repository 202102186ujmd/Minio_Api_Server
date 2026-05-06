package middleware

import (
	"net/http"

	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/logger"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(logg *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logg.Error("panic recovered", logger.String("request_id", requestIDFromContext(c)), zap.Any("error", err))
				response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Ocurrió un error interno", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
