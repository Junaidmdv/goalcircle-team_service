package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type ObjectStorageConfig struct {
	StorageProvider    string
	AccesskeyId        string
	SecreteKey         string
	EndPoint           string
	Region             string
	Bucket             string
	SSL                bool
	PresignedURLExpiry time.Duration
}

func (cb *configBuilder) WithObjectStorage() ConfigBuilder {

	storageProvider := os.Getenv("STORAGE_PROVIDER")

	var endpoint, region string
	switch storageProvider {
	case "minio":
		endpoint = os.Getenv("STORAGE_ENDPOINT")
		if endpoint == "" {
			cb.errors = append(cb.errors, errors.New("missing storage endpoint"))
		}

	case "s3":
		region = os.Getenv("STORAGE_REGION")
		if region == "" {
			cb.errors = append(cb.errors, errors.New("missing storage region"))
		}

	default:
		cb.errors = append(cb.errors, errors.New("invalid STORAGE_PROVIDER"))
	}

	accessKey := os.Getenv("STORAGE_ACCESSKEY")
	if accessKey == "" {
		cb.errors = append(cb.errors, errors.New("missing stroage access key id"))
	}

	secreteKey := os.Getenv("STORAGE_SECRETEKEY")
	if secreteKey == "" {
		cb.errors = append(cb.errors, errors.New("missing storage secretekey"))
	}

	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		cb.errors = append(cb.errors, errors.New("failed to add bucket name"))
	}

	useSSL, err := strconv.ParseBool(os.Getenv("STORAGE_SSL"))
	if err != nil {
		cb.errors = append(cb.errors, errors.New("failed to add bucket name"))
	}

	duration, err := time.ParseDuration(os.Getenv("PRESIGNED_URL_EXPIRY"))
	if err != nil {
		cb.errors = append(cb.errors, errors.New("failed to add presigned url expiry"))
	}

	if len(cb.errors) > 0 {
		return cb
	}
	

	cb.config.StorageConfig = &ObjectStorageConfig{
		StorageProvider:    storageProvider,
		AccesskeyId:        accessKey,
		SecreteKey:         secreteKey,
		EndPoint:           endpoint,
		Bucket:             bucket,
		Region:             region,
		SSL:                useSSL,
		PresignedURLExpiry: duration,
	}

	return cb
}
