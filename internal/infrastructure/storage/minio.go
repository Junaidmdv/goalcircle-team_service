package storage

import (
	"github.com/Junaidmdv/goalcircle-team_service/internal/config"
	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minio sdk can store data  in locally and it's compatible wich s3.
type Minio struct {
	Client *minio.Client
	Bucket string
}

func NewMinio(config config.MinioConfig) (*Minio, error) {
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
		Client: minio,
	}, nil
}
