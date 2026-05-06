package middleware

import (
	"net/http"
	"strings"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func APIKey(cfg config.Config) gin.HandlerFunc {
	excluded := []string{"/health", "/metrics", "/swagger"}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, prefix := range excluded {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

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
