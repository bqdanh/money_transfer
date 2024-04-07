package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bqdanh/money_transfer/internal/adapters/grpc"
	"github.com/bqdanh/money_transfer/internal/adapters/http_gateway"
	"github.com/bqdanh/money_transfer/pkg/database"
	"github.com/bqdanh/money_transfer/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Env      string              `json:"env" mapstructure:"env"`
	GRPC     grpc.Config         `json:"grpc" mapstructure:"grpc"`
	HTTP     http_gateway.Config `json:"http" mapstructure:"http"`
	Database database.Config     `json:"database" mapstructure:"database"`
	Logger   logger.Config       `json:"logger" mapstructure:"logger"`
}

func loadDefaultConfig() *Config {
	return &Config{
		Env: "local",
		GRPC: grpc.Config{
			Host: "0.0.0.0",
			Port: 9090,
		},
		HTTP: http_gateway.Config{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: database.Config{},
		Logger:   logger.Config{},
	}
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	/**
	|-------------------------------------------------------------------------
	| You should set default config value here
	| 1. Populate the default value in (Source code)
	| 2. Then merge from config (YAML) and OS environment
	|-----------------------------------------------------------------------*/
	// Load default config
	c := loadDefaultConfig()
	configBuffer, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	//1. Populate the default value in (Source code)
	if err := viper.ReadConfig(bytes.NewBuffer(configBuffer)); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	//2. Then merge from config (YAML) and OS environment
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge in config: %w", err)
	}
	// Populate all config again
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return c, err
}
