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
	ErrInvalidExpiry   = errors.New("expiración fuera de rango")
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

type CopyMoveResult struct {
	Bucket         string `json:"bucket"`
	ObjectName     string `json:"object_name"`
	Destination    string `json:"destination"`
	DestinationKey string `json:"destination_key"`
	ETag           string `json:"etag"`
}

func NewService(client *minio.Client, cfg config.Config) *Service {
	return &Service{client: client, cfg: cfg}
}

func (s *Service) Ping(ctx context.Context) error {
	_, err := s.client.ListBuckets(ctx)
	return err
}

func (s *Service) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return s.client.BucketExists(ctx, bucket)
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

func (s *Service) GetObject(ctx context.Context, bucket, objectName string) (*minio.Object, minio.ObjectInfo, error) {
	obj, err := s.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}
	info, err := obj.Stat()
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}
	return obj, info, nil
}

func (s *Service) PresignURL(ctx context.Context, bucket, objectName, method string, expiry time.Duration) (string, error) {
	if method == "GET" {
		url, err := s.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
		if err != nil {
			return "", err
		}
		return url.String(), nil
	}

	url, err := s.client.PresignedPutObject(ctx, bucket, objectName, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *Service) ParseExpiry(value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	if parsed < 60 || parsed > 604800 {
		return 0, ErrInvalidExpiry
	}
	return parsed, nil
}

func (s *Service) CopyObject(ctx context.Context, bucket, objectName, destBucket, destObject string) (*CopyMoveResult, error) {
	src := minio.CopySrcOptions{Bucket: bucket, Object: objectName}
	dst := minio.CopyDestOptions{Bucket: destBucket, Object: destObject}
	info, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return nil, err
	}

	return &CopyMoveResult{
		Bucket:         bucket,
		ObjectName:     objectName,
		Destination:    destBucket,
		DestinationKey: destObject,
		ETag:           info.ETag,
	}, nil
}

func (s *Service) MoveObject(ctx context.Context, bucket, objectName, destBucket, destObject string) (*CopyMoveResult, error) {
	info, err := s.CopyObject(ctx, bucket, objectName, destBucket, destObject)
	if err != nil {
		return nil, err
	}

	if err := s.DeleteObject(ctx, bucket, objectName); err != nil {
		return nil, err
	}
	return info, nil
}

func (s *Service) DeleteObject(ctx context.Context, bucket, objectName string) error {
	return s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}
