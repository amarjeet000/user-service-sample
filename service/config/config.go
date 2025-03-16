// Package config instantiates the service configuraion.
// It prioritizes ENV vars over the vars provided via the config.yml file.
package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	DefaultHost          = "0.0.0.0"
	DefaultPort          = "3030"
	ConfigFileName       = "config"
	ConfigFileDir        = "../service_config"
	DefaultKeyDir        = "../keys"
	DefaultSigningMethod = "rsa"
)

type Config struct {
	Host          string
	Port          string
	KeyDir        string
	SigningMethod string
}

// defaultConfig initializes config based on a config file.
func defaultConfig() {
	viper.SetConfigName(ConfigFileName)
	viper.AddConfigPath(ConfigFileDir)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("INFO: Could not read the config file, using defaults or ENV vars", err)
	}
}

// GetConfig instantiates a Config instance in a manner that prioritizes env over confg file, if same var is present in both env and the file.
func GetConfig() (*Config, error) {
	defaultConfig()
	viper.BindEnv("host", "HOST")
	viper.BindEnv("port", "PORT")
	viper.BindEnv("keydir", "KEYDIR")
	viper.BindEnv("signing-method", "SIGNING_METHOD")

	viper.SetDefault("host", DefaultHost)
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("keydir", DefaultKeyDir)
	viper.SetDefault("signing-method", DefaultSigningMethod)

	cfg := &Config{
		Host:          viper.GetString("host"),
		Port:          viper.GetString("port"),
		KeyDir:        viper.GetString("keydir"),
		SigningMethod: viper.GetString("signing-method"),
	}

	return cfg, nil
}
