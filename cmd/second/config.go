package second

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env string `json:"env" mapstructure:"env"`
}

func loadDefaultConfig() *Config {
	return &Config{
		Env: "local",
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
