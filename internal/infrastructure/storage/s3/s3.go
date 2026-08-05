package s3

import (
	"context"
	"io"
	"time"
)

type S3Storage struct {
}

func (m *S3Storage) Upload(ctx context.Context, bucketName string, objectName string, data io.Reader, size int64, contentType string) (string, error) {
	return "", nil
}

func (m *S3Storage) GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error) {
	return "", nil
}

func (m *S3Storage) Delete(ctx context.Context, bucketName string, objectName string) error {
	return nil
}
