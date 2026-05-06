package handlers

import (
	"net/http"
	"time"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/202102186ujmd/Minio_Api_Server/internal/domain/storage"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service *storage.Service
	cfg     config.Config
	start   time.Time
}

func NewHealthHandler(service *storage.Service, cfg config.Config) *HealthHandler {
	return &HealthHandler{service: service, cfg: cfg, start: time.Now()}
}

// Health
// @Summary Health check
// @Description Verifica que el servicio esté activo.
// @Tags Health
// @Success 200 {object} response.APIResponse "OK"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "OK", "Servicio activo", gin.H{
		"status": "up",
		"uptime": time.Since(h.start).String(),
	})
}

// Live
// @Summary Liveness
// @Description Verifica que el proceso esté vivo.
// @Tags Health
// @Success 200 {object} response.APIResponse "OK"
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response.Success(c, http.StatusOK, "OK", "Proceso vivo", gin.H{"status": "live"})
}

// Ready
// @Summary Readiness
// @Description Verifica conectividad con MinIO y bucket configurado.
// @Tags Health
// @Success 200 {object} response.APIResponse "Ready"
// @Failure 503 {object} response.APIResponse "No listo"
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	status := "ready"
	ready := true
	issues := make([]string, 0)

	if err := h.service.Ping(c.Request.Context()); err != nil {
		ready = false
		status = "degraded"
		issues = append(issues, "minio_unreachable")
	}

	if h.cfg.MinioBucket != "" {
		exists, err := h.service.BucketExists(c.Request.Context(), h.cfg.MinioBucket)
		if err != nil || !exists {
			ready = false
			status = "degraded"
			issues = append(issues, "bucket_not_ready")
		}
	}

	if !ready {
		response.Fail(c, http.StatusServiceUnavailable, "NOT_READY", "Servicio no listo", gin.H{
			"status": status,
			"issues": issues,
		})
		return
	}

	response.Success(c, http.StatusOK, "READY", "Servicio listo", gin.H{
		"status": status,
		"issues": issues,
	})
}
