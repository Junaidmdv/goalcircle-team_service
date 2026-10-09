package config

import (
	"errors"
	"os"
)

type NatsConfig struct {
	URL string
}

func (cb *configBuilder) WithNats() ConfigBuilder {

	url := os.Getenv("NATS_URL")
	if url == "" {
		cb.errors = append(cb.errors, errors.New("missing nats url"))
		return cb
	}
	cb.config.NatsConfig.URL = url

	return cb
}
