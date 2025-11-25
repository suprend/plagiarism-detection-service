package s3

import (
	"bytes"
	"context"
	"io"
	"os"

	apperr "filestorage/internal/common/errors"
	"filestorage/internal/domain/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Repository struct {
	client *s3.Client
	bucket string
}

// NewS3Repository создает новый экземпляр S3 репозитория
func NewS3Repository(ctx context.Context, bucket, endpoint, region string) (repository.S3Repository, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	var cfg aws.Config
	var err error

	// Если указан endpoint (для S3-совместимых хранилищ типа MinIO, Yandex Object Storage)
	if endpoint != "" {
		// Используем credentials из переменных окружения для S3-совместимых хранилищ
		accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

		if accessKeyID != "" && secretAccessKey != "" {
			opts = append(opts, config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
			))
		}

		cfg, err = config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, err
		}

		cfg.BaseEndpoint = aws.String(endpoint)
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, err
		}
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.UsePathStyle = true // Для S3-совместимых хранилищ
		}
	})

	return &s3Repository{
		client: client,
		bucket: bucket,
	}, nil
}

// UploadFile загружает файл в S3
func (r *s3Repository) UploadFile(ctx context.Context, key string, data []byte, contentType string) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return apperr.Wrap(err, apperr.CodeStorage, "failed to upload object")
	}
	return nil
}

// GetFile получает файл из S3
func (r *s3Repository) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeStorage, "failed to get object")
	}

	return result.Body, nil
}

// DeleteFile удаляет файл из S3
func (r *s3Repository) DeleteFile(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return apperr.Wrap(err, apperr.CodeStorage, "failed to delete object")
	}
	return nil
}
