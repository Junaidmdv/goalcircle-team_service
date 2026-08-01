package storage

import (
	"errors"

	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	minio_store "github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage/minio"
	"github.com/Junaidmdv/goalcircle-team_service/internal/infrastructure/storage/s3"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
)

type ObjectStorage interface {
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
