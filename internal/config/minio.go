package config

import (
	"errors"
	"log"
	"os"
	"strconv"
)

type MinioConfig struct {
	AccesskeyId string
	SecreteKey  string
	EndPoint    string
	SSL         bool
}

func (cb *configBuilder) WithMinio() ConfigBuilder {
	accessKey := os.Getenv("MINIO_ROOT_USER")
	if accessKey == "" {
		cb.errors = append(cb.errors, errors.New("missing minio access key id"))
	}

	secreteKey := os.Getenv("MINIO_ROOT_PASSWORD")
	if secreteKey == "" {
		cb.errors = append(cb.errors, errors.New("missing minio secretekey"))
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		cb.errors = append(cb.errors, errors.New("missing minio endpoint "))
	}

	bucket := os.Getenv("MINIO_BUCKET")

	if bucket == "" {
		cb.errors = append(cb.errors, errors.New("failed to add bucket name"))
	}

	isSSL, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))

	if !isSSL {
		log.Printf("minio ssl is not set, running without ssl")
	}

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.MinioConfig = &MinioConfig{
		AccesskeyId: accessKey,
		SecreteKey:  secreteKey,
		EndPoint:    endpoint,
		SSL:         isSSL,
	}

	return cb
}
