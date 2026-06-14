package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type GRPCServerConfig struct {
	Port    int
	TimeOut time.Duration
}

func (cb *configBuilder) WithGrpcServer() ConfigBuilder {
	port, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		cb.errors = append(cb.errors, errors.New("port is missing"))
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		cb.errors = append(cb.errors, errors.New("invlid or missing timout data"))
	}

	if len(cb.errors) > 0 {
		return cb
	}

	cb.config.Server = &GRPCServerConfig{
		Port:    port,
		TimeOut: timeout,
	}
	return cb
}
