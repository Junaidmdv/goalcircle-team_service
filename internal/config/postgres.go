package config 

import (
	"errors"
	"os"
)

type PostgresConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func (cb *configBuilder) WithPostgres() ConfigBuilder {
	host := os.Getenv("DB_HOST")
	if host == "" {
		cb.errors = append(cb.errors, errors.New("DB_HOST is required"))

	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		cb.errors = append(cb.errors, errors.New("DB_PORT is required"))
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		cb.errors = append(cb.errors, errors.New("DB_NAME is required"))
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		cb.errors = append(cb.errors, errors.New("DB_USER is required"))
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		cb.errors = append(cb.errors, errors.New("DB_PASSWORD is required"))
	}

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.Postgres = &PostgresConfig{
		Host:     host,
		Port:     port,
		DBName:   dbName,
		User:     user,
		Password: password,
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	return cb
}
