package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type UserService struct {
	Host    string
	Port    int
	TimeOut time.Duration
}

func (cb *configBuilder) WithUserService() ConfigBuilder {
	host := os.Getenv("USER_SERVICE_HOST")
	if host == "" {
		cb.errors = append(cb.errors, fmt.Errorf("USER_SERVICE_HOST is required"))
		return cb
	}

	portStr := os.Getenv("USER_SERVICE_PORT")
	if portStr == "" {
		cb.errors = append(cb.errors, fmt.Errorf("USER_SERVICE_PORT is required"))
		return cb
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("USER_SERVICE_PORT must be a number"))
		return cb
	}

	timeoutStr := os.Getenv("USER_SERVICE_TIMEOUT")
	if timeoutStr == "" {
		cb.errors = append(cb.errors, fmt.Errorf("USER_SERVICE_TIMEOUT is required"))
		return cb
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		cb.errors = append(cb.errors, fmt.Errorf("USER_SERVICE_TIMEOUT must be a valid duration"))
		return cb
	}

	// UserServiceClient: Service{
	// 	Host:    host,
	// 	Port:    port,
	// 	TimeOut: timeout,
	// },

	cb.config.UserSrv = &UserService{
		Host:    host,
		Port:    port,
		TimeOut: timeout,
	}

	return cb
}
