package handlers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/202102186ujmd/Minio_Api_Server/internal/domain/storage"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/response"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/validator"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	service *storage.Service
	cfg     config.Config
}

type CopyMoveRequest struct {
	DestinationBucket string `json:"destination_bucket" validate:"required"`
	DestinationObject string `json:"destination_object" validate:"required"`
}

func NewStorageHandler(service *storage.Service, cfg config.Config) *StorageHandler {
	return &StorageHandler{service: service, cfg: cfg}
}

// ListBuckets
// @Summary Listar buckets
// @Description Devuelve todos los buckets disponibles.
// @Tags Buckets
// @Security ApiKeyAuth
// @Success 200 {object} response.APIResponse "Lista de buckets"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets [get]
func (h *StorageHandler) ListBuckets(c *gin.Context) {
	buckets, err := h.service.ListBuckets(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al listar buckets", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", "Buckets listados", buckets)
}

// ListObjects
// @Summary Listar objetos
// @Description Lista objetos y carpetas dentro de un bucket.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param prefix query string false "Prefijo para filtrar (simula carpetas)" example("videos/")
// @Param recursive query bool false "Recursivo" example(true)
// @Param limit query int false "Límite de resultados" example(100)
// @Success 200 {object} response.APIResponse "Lista de objetos"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects [get]
func (h *StorageHandler) ListObjects(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	if bucket == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket requerido", nil)
		return
	}

	prefix := c.DefaultQuery("prefix", "")
	recursive := c.DefaultQuery("recursive", "false") == "true"
	limit := c.DefaultQuery("limit", "0")

	objects, err := h.service.ListObjects(c.Request.Context(), bucket, prefix, recursive, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al listar objetos", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", "Objetos listados", objects)
}

// UploadObject
// @Summary Subir archivo
// @Description Sube un archivo al bucket indicado.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param file formData file true "Archivo a subir"
// @Param path formData string false "Ruta destino dentro del bucket" example("carpeta/subcarpeta")
// @Success 201 {object} response.APIResponse "Archivo subido"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 413 {object} response.APIResponse "Archivo demasiado grande"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects [post]
func (h *StorageHandler) UploadObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	if bucket == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket requerido", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Archivo requerido", err.Error())
		return
	}

	if err := h.validateFile(file); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, storage.ErrFileTooLarge) {
			status = http.StatusRequestEntityTooLarge
		}
		response.Fail(c, status, "VALIDATION", err.Error(), nil)
		return
	}

	path := strings.TrimSpace(c.PostForm("path"))
	objectName := h.buildObjectName(path, file.Filename)
	if objectName == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Nombre de archivo inválido", nil)
		return
	}

	opened, err := file.Open()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "FILE_OPEN_ERROR", "No se pudo abrir el archivo", err.Error())
		return
	}
	defer opened.Close()

	contentType := file.Header.Get("Content-Type")

	upload, err := h.service.UploadObject(c.Request.Context(), bucket, objectName, opened, file.Size, contentType)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al subir archivo", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "CREATED", "Archivo subido", upload)
}

// DownloadObject
// @Summary Descargar archivo
// @Description Descarga un archivo por nombre dentro del bucket.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param object path string true "Nombre del objeto"
// @Success 200 {file} file "Archivo"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 404 {object} response.APIResponse "No encontrado"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects/{object}/download [get]
func (h *StorageHandler) DownloadObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	object := strings.TrimSpace(c.Param("object"))
	if bucket == "" || object == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket y objeto requeridos", nil)
		return
	}

	obj, info, err := h.service.GetObject(c.Request.Context(), bucket, object)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al descargar archivo", err.Error())
		return
	}
	defer obj.Close()

	contentType := info.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Disposition", "attachment; filename=\""+filepath.Base(object)+"\"")
	c.DataFromReader(http.StatusOK, info.Size, contentType, obj, nil)
}

// PresignObject
// @Summary Generar URL firmada
// @Description Genera una URL firmada para GET o PUT.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param object path string true "Nombre del objeto"
// @Param method query string false "GET o PUT" example("GET")
// @Param expiry query int false "Expiración en segundos" example(3600)
// @Success 200 {object} response.APIResponse "URL firmada"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects/{object}/presign [get]
func (h *StorageHandler) PresignObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	object := strings.TrimSpace(c.Param("object"))
	if bucket == "" || object == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket y objeto requeridos", nil)
		return
	}

	method := strings.ToUpper(strings.TrimSpace(c.DefaultQuery("method", "GET")))
	if method != http.MethodGet && method != http.MethodPut {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Método inválido, usa GET o PUT", nil)
		return
	}

	expirySeconds := int64(3600)
	if q := strings.TrimSpace(c.Query("expiry")); q != "" {
		parsed, err := h.service.ParseExpiry(q)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "VALIDATION", "Expiración inválida", err.Error())
			return
		}
		expirySeconds = parsed
	}

	url, err := h.service.PresignURL(c.Request.Context(), bucket, object, method, time.Duration(expirySeconds)*time.Second)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al generar URL", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OK", "URL firmada generada", gin.H{
		"url":    url,
		"expiry": expirySeconds,
		"method": method,
	})
}

// CopyObject
// @Summary Copiar archivo
// @Description Copia un archivo a otro bucket u objeto.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param object path string true "Nombre del objeto"
// @Param payload body CopyMoveRequest true "Destino"
// @Success 200 {object} response.APIResponse "Archivo copiado"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects/{object}/copy [post]
func (h *StorageHandler) CopyObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	object := strings.TrimSpace(c.Param("object"))
	if bucket == "" || object == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket y objeto requeridos", nil)
		return
	}

	var payload CopyMoveRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Payload inválido", err.Error())
		return
	}
	if err := validator.ValidateStruct(payload); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Campos requeridos", err.Error())
		return
	}

	info, err := h.service.CopyObject(c.Request.Context(), bucket, object, payload.DestinationBucket, payload.DestinationObject)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al copiar archivo", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OK", "Archivo copiado", info)
}

// MoveObject
// @Summary Mover archivo
// @Description Mueve un archivo a otro bucket u objeto (copia y elimina el original).
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param object path string true "Nombre del objeto"
// @Param payload body CopyMoveRequest true "Destino"
// @Success 200 {object} response.APIResponse "Archivo movido"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects/{object}/move [post]
func (h *StorageHandler) MoveObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	object := strings.TrimSpace(c.Param("object"))
	if bucket == "" || object == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket y objeto requeridos", nil)
		return
	}

	var payload CopyMoveRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Payload inválido", err.Error())
		return
	}
	if err := validator.ValidateStruct(payload); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Campos requeridos", err.Error())
		return
	}

	info, err := h.service.MoveObject(c.Request.Context(), bucket, object, payload.DestinationBucket, payload.DestinationObject)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al mover archivo", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "OK", "Archivo movido", info)
}

// DeleteObject
// @Summary Eliminar archivo
// @Description Elimina un archivo por nombre dentro del bucket.
// @Tags Objects
// @Security ApiKeyAuth
// @Param bucket path string true "Nombre del bucket"
// @Param object path string true "Nombre del objeto"
// @Success 200 {object} response.APIResponse "Archivo eliminado"
// @Failure 400 {object} response.APIResponse "Validación"
// @Failure 401 {object} response.APIResponse "API key inválida"
// @Failure 404 {object} response.APIResponse "No encontrado"
// @Failure 500 {object} response.APIResponse "Error interno"
// @Router /buckets/{bucket}/objects/{object} [delete]
func (h *StorageHandler) DeleteObject(c *gin.Context) {
	bucket := strings.TrimSpace(c.Param("bucket"))
	object := strings.TrimSpace(c.Param("object"))
	if bucket == "" || object == "" {
		response.Fail(c, http.StatusBadRequest, "VALIDATION", "Bucket y objeto requeridos", nil)
		return
	}

	if err := h.service.DeleteObject(c.Request.Context(), bucket, object); err != nil {
		response.Fail(c, http.StatusInternalServerError, "MINIO_ERROR", "Error al eliminar archivo", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", "Archivo eliminado", gin.H{"object": object})
}

func (h *StorageHandler) validateFile(file *multipart.FileHeader) error {
	maxBytes := h.cfg.MaxUploadSizeMB * 1024 * 1024
	if maxBytes > 0 && file.Size > maxBytes {
		return storage.ErrFileTooLarge
	}

	contentType := file.Header.Get("Content-Type")
	if len(h.cfg.AllowedMIMEs) > 0 {
		allowed := false
		for _, mime := range h.cfg.AllowedMIMEs {
			if strings.EqualFold(strings.TrimSpace(mime), contentType) {
				allowed = true
				break
			}
		}
		if !allowed {
			return storage.ErrMimeNotAllowed
		}
	}

	name := filepath.Base(file.Filename)
	if name == "." || name == ".." || name == "" {
		return storage.ErrInvalidFileName
	}
	return nil
}

func (h *StorageHandler) buildObjectName(path, filename string) string {
	cleaned := strings.Trim(path, "/")
	filename = filepath.Base(filename)
	if filename == "" || filename == "." || filename == ".." {
		return ""
	}
	if cleaned == "" {
		return filename
	}
	return cleaned + "/" + filename
}
