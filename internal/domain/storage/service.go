package storage

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/minio/minio-go/v7"
)

var (
	ErrFileTooLarge    = errors.New("el archivo excede el tamaño permitido")
	ErrMimeNotAllowed  = errors.New("tipo MIME no permitido")
	ErrInvalidFileName = errors.New("nombre de archivo inválido")
)

type Service struct {
	client *minio.Client
	cfg    config.Config
}

type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type,omitempty"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	IsDir        bool      `json:"is_dir"`
}

type UploadResult struct {
	Bucket      string `json:"bucket"`
	ObjectName  string `json:"object_name"`
	Size        int64  `json:"size"`
	ETag        string `json:"etag"`
	VersionID   string `json:"version_id,omitempty"`
	Location    string `json:"location"`
	ContentType string `json:"content_type,omitempty"`
}

func NewService(client *minio.Client, cfg config.Config) *Service {
	return &Service{client: client, cfg: cfg}
}

func (s *Service) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return s.client.ListBuckets(ctx)
}

func (s *Service) ListObjects(ctx context.Context, bucket, prefix string, recursive bool, limitStr string) ([]ObjectInfo, error) {
	options := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}

	limit := int64(0)
	if limitStr != "" {
		if parsed, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = parsed
		}
	}

	objects := make([]ObjectInfo, 0)
	count := int64(0)

	for object := range s.client.ListObjects(ctx, bucket, options) {
		if object.Err != nil {
			return nil, object.Err
		}

		objects = append(objects, ObjectInfo{
			Key:          object.Key,
			Size:         object.Size,
			ContentType:  object.ContentType,
			LastModified: object.LastModified,
			ETag:         object.ETag,
			IsDir:        object.IsDir,
		})

		count++
		if limit > 0 && count >= limit {
			break
		}
	}

	return objects, nil
}

func (s *Service) UploadObject(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	info, err := s.client.PutObject(ctx, bucket, objectName, reader, size, opts)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Bucket:      bucket,
		ObjectName:  info.Key,
		Size:        info.Size,
		ETag:        info.ETag,
		VersionID:   info.VersionID,
		Location:    info.Location,
		ContentType: contentType,
	}, nil
}

func (s *Service) DeleteObject(ctx context.Context, bucket, objectName string) error {
	return s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}
