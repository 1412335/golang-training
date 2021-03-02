package config

import (
	"fw/pkg/config"
	"time"
)

type ServiceConfig struct {
	JWT      *JWT
	Database *config.Database
}

// json web token
type JWT struct {
	SecretKey string
	Duration  time.Duration
	Issuer    string
}

// manager grpc-pool
type ManagerClient struct {
	MaxPoolSize int
	TimeOut     int
	// method need to request with authentication
	AuthMethods map[string]bool
	// credentials authentication
	Authentication *Authentication
	// jwt token
	RefreshDuration time.Duration
}

type Authentication struct {
	Username string
	Password string
}
