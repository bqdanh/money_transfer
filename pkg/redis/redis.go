package redis

import (
	"context"
	"fmt"
	"time"

	goredislib "github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
	Username string `json:"username" mapstructure:"username"`
	PoolSize int    `json:"pool_size" mapstructure:"pool_size"`
	DB       int    `json:"db" mapstructure:"db"`
}

func NewRedisClient(cfg Config) (*goredislib.Client, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		Username: cfg.Username,
		PoolSize: cfg.PoolSize,
		DB:       cfg.DB,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	result := client.Ping(ctx)
	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return client, nil
}
