package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        *GRPCServerConfig
	Postgres      *PostgresConfig
	UserSrv       *UserService
	StorageConfig *ObjectStorageConfig
}

type ConfigBuilder interface {
	WithGrpcServer() ConfigBuilder
	WithPostgres() ConfigBuilder
	WithUserService() ConfigBuilder
	WithObjectStorage() ConfigBuilder
	Build() (*Config, []error)
}

type configBuilder struct {
	config *Config
	errors []error
}

func LoadConfig() ConfigBuilder {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, reading from environment")
		return nil
	}

	return &configBuilder{
		config: &Config{},
	}
}

func (cb *configBuilder) Build() (*Config, []error) {
	if len(cb.errors) > 0 {
		return nil, cb.errors
	}

	return cb.config, nil
}
