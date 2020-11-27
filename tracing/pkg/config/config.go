package config

import (
	"github.com/spf13/viper"
)

const (
	ConfigName = "config"
	ConfigType = "yml"
	ConfigPath = "."
)

type ServiceConfig struct {
	ServiceName   string
	Port          int
	Metrics       string
	Formatter     *Formatter
	Publisher     *Publisher
	FormatterHost string
	PublisherHost string
}

type Formatter struct {
	ServiceName string
	Host        string
	Port        int
}

type Publisher struct {
	ServiceName string
	Host        string
	Port        int
}

type Database struct {
	Host     string
	User     string
	Password string
	Scheme   string
}

type Authentication struct {
	Username string
	Password string
}

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
