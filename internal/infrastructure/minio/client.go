package minio

import (
	"context"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(cfg config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
		Region: cfg.MinioRegion,
	})
	if err != nil {
		return nil, err
	}

	if cfg.MinioBucket != "" {
		exists, err := client.BucketExists(context.Background(), cfg.MinioBucket)
		if err != nil {
			return nil, err
		}
		if !exists {
			if err := client.MakeBucket(context.Background(), cfg.MinioBucket, minio.MakeBucketOptions{Region: cfg.MinioRegion}); err != nil {
				return nil, err
			}
		}
	}

	return client, nil
}
