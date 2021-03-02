package configs

import (
	"time"

	"github.com/spf13/viper"
)

const (
	ConfigName = "config"
	ConfigType = "yml"
	ConfigPath = "."
)

type ServiceConfig struct {
	ManagerClient  *ManagerClient
	JWT            *JWT
	Database       *Database
	Authentication *Authentication
	// server using
	AccessibleRoles map[string][]string
	// client using
	AuthMethods map[string]bool
}

// json web token
type JWT struct {
	SecretKey string
	Duration  time.Duration
	Issuer    string
}

type Database struct {
	Host     string
	User     string
	Password string
	Scheme   string
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

// type AccessibleRoles map[string][]string

func LoadConfig(cfgFile string, cfg interface{}) error {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.SetConfigName(ConfigName)
		viper.SetConfigType(ConfigType)
		viper.AddConfigPath(ConfigPath)
	}
	viper.AutomaticEnv()
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
