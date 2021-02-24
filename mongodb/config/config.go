package config

type DatabaseConfig struct {
	ConnectionURI string
	Database      string
	PoolSize      uint64
}
