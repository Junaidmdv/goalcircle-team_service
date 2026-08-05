package minio

import (
	"context"
	"io"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minio sdk can store data  in locally and it's compatible wich s3.
type Minio struct {
	client *minio.Client
	Bucket string
	logger logger.Logger
}

func NewMinio(config config.ObjectStorageConfig, logger logger.Logger) (*Minio, error) {
	minio, err := minio.New(config.EndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			config.AccesskeyId, config.SecreteKey, "",
		),
		Secure: config.SSL,
	})

	if err != nil {
		return nil, apperror.NewInternalError("failed minio configration", err)
	}

	return &Minio{
		client: minio,
		logger: logger,
	}, nil
}

func (m *Minio) Upload(ctx context.Context, bucketName string, objectName string, data io.Reader, size int64, contentType string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	info, err := m.client.PutObject(
		ctx,
		bucketName,
		objectName,
		data,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: map[string]string{
				"uploaded-by": "goalcircle.team",
				"entity":      "team-logo",
			},
		})

	if err != nil {
		m.logger.Error("error minio store", "error", err, "method", "minio.Upload")
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return info.Key, nil
}

func (m *Minio) GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		m.logger.Error("failed get presigned url from minio", "error", err, "method", "minio.GetPresignedURL")
		return "", apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}

	return url.String(), nil
}

func (m *Minio) Delete(ctx context.Context, bucketName string, objectName string) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := m.client.RemoveObject(
		ctx,
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		m.logger.Error("failed delete operation", "error", err, "method", "minio.Delete")
		return apperror.NewInternalError(apperror.InternalErrorMsg, err)
	}
	return nil
}
