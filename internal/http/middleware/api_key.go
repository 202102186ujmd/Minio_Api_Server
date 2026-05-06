package middleware

import (
	"net/http"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func APIKey(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || apiKey != cfg.APIKey {
			response.Fail(c, http.StatusUnauthorized, "AUTH_INVALID", "API key inválida", gin.H{
				"hint": "Incluye el header X-API-Key con tu clave válida",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
