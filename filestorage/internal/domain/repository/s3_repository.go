package repository

import (
	"context"
	"io"
)

type S3Repository interface {
	UploadFile(ctx context.Context, key string, data []byte, contentType string) error

	GetFile(ctx context.Context, key string) (io.ReadCloser, error)

	DeleteFile(ctx context.Context, key string) error
}
