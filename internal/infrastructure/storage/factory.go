package storage

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	minio_store "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage/minio"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage/s3"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type ObjectStorage interface {
	Upload(ctx context.Context, bucketName string, objectName string, imageData io.Reader, size int64, contentType string) (string, error)
	GetPresignedURL(ctx context.Context, bucketName string, objectName string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, bucketName string, objectName string) error
}

func ObjectStorageFactoryMethod(ob config.ObjectStorageConfig, logger logger.Logger) (ObjectStorage, error) {
	switch ob.StorageProvider {
	case "minio":
		return minio_store.NewMinio(ob, logger)
	case "sw":
		return &s3.S3Storage{}, nil
	default:

		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, errors.New("invalid object storage method"))
	}

}
