package routes

import (
	"github.com/202102186ujmd/Minio_Api_Server/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, health *handlers.HealthHandler, storage *handlers.StorageHandler) {
	r.GET("/health", health.Health)

	v1 := r.Group("/v1")
	{
		v1.GET("/buckets", storage.ListBuckets)
		v1.GET("/buckets/:bucket/objects", storage.ListObjects)
		v1.POST("/buckets/:bucket/objects", storage.UploadObject)
		v1.DELETE("/buckets/:bucket/objects/:object", storage.DeleteObject)
	}
}
