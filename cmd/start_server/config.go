package start_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bqdanh/money_transfer/internal/adapters/distribute_lock"
	"github.com/bqdanh/money_transfer/internal/adapters/server/grpc_server"
	"github.com/bqdanh/money_transfer/internal/adapters/server/http_gateway"
	"github.com/bqdanh/money_transfer/internal/adapters/user_token"
	"github.com/bqdanh/money_transfer/internal/applications/accounts/link_account"
	"github.com/bqdanh/money_transfer/internal/applications/authenticate/generate_user_token"
	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/create_transaction"
	"github.com/bqdanh/money_transfer/internal/applications/transactions/deposit/process_transaction"
	implement_bank_account2 "github.com/bqdanh/money_transfer/internal/entities/account/sof/bank_account/implement_bank_account"
	"github.com/bqdanh/money_transfer/pkg/database"
	"github.com/bqdanh/money_transfer/pkg/logger"
	pkgredis "github.com/bqdanh/money_transfer/pkg/redis"
	"github.com/spf13/viper"
)

var (
	//load sof registry
	_ = implement_bank_account2.SourceOfFundCodeACB
	_ = implement_bank_account2.SourceOfFundCodeVCB
	_ = implement_bank_account2.SourceOfFundCodeVIB
)

type Config struct {
	Env                string                     `json:"env" mapstructure:"env"`
	GRPC               grpc_server.Config         `json:"grpc" mapstructure:"grpc"`
	HTTP               http_gateway.Config        `json:"http" mapstructure:"http"`
	Database           database.Config            `json:"database" mapstructure:"database"`
	Logger             logger.Config              `json:"logger" mapstructure:"logger"`
	JwtToken           user_token.Config          `json:"jwt_token" mapstructure:"jwt_token"`
	GenerateToken      generate_user_token.Config `json:"generate_token" mapstructure:"generate_token"`
	LinkAccount        link_account.Config        `json:"link_account" mapstructure:"link_account"`
	DistributeLock     distribute_lock.Config     `json:"distribute_lock" mapstructure:"distribute_lock"`
	RedisConnection    pkgredis.Config            `json:"redis_connection" mapstructure:"redis_connection"`
	CreateTransaction  create_transaction.Config  `json:"create_transaction" mapstructure:"create_transaction"`
	ProcessTransaction process_transaction.Config `json:"process_transaction" mapstructure:"process_transaction"`
}

func loadDefaultConfig() *Config {
	return &Config{
		Env: "local",
		GRPC: grpc_server.Config{
			Host: "0.0.0.0",
			Port: 9090,
		},
		HTTP: http_gateway.Config{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database:           database.Config{},
		Logger:             logger.Config{},
		JwtToken:           user_token.Config{},
		GenerateToken:      generate_user_token.Config{},
		LinkAccount:        link_account.DefaultConfig,
		DistributeLock:     distribute_lock.DefaultConfig,
		CreateTransaction:  create_transaction.DefaultConfig,
		ProcessTransaction: process_transaction.DefaultConfig,
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
