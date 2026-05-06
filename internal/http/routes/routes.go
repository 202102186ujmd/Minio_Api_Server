package routes

import (
	"github.com/202102186ujmd/Minio_Api_Server/internal/http/handlers"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(r *gin.Engine, health *handlers.HealthHandler, storage *handlers.StorageHandler) {
	r.GET("/health", health.Health)
	r.GET("/health/live", health.Live)
	r.GET("/health/ready", health.Ready)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := r.Group("/v1")
	{
		v1.GET("/buckets", storage.ListBuckets)
		v1.GET("/buckets/:bucket/objects", storage.ListObjects)
		v1.POST("/buckets/:bucket/objects", storage.UploadObject)
		v1.GET("/buckets/:bucket/objects/:object/download", storage.DownloadObject)
		v1.GET("/buckets/:bucket/objects/:object/presign", storage.PresignObject)
		v1.POST("/buckets/:bucket/objects/:object/copy", storage.CopyObject)
		v1.POST("/buckets/:bucket/objects/:object/move", storage.MoveObject)
		v1.DELETE("/buckets/:bucket/objects/:object", storage.DeleteObject)
	}
}
