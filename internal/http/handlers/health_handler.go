package handlers

import (
	"net/http"

	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health
// @Summary Health check
// @Description Verifica que el servicio esté activo.
// @Tags Health
// @Success 200 {object} response.APIResponse "OK"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "OK", "Servicio activo", gin.H{"status": "up"})
}
