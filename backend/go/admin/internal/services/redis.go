package services

import (
	"admin/internal/config"
	"admin/internal/services/redisc"
	"context"
	"fmt"
)

type Redis struct {
	redisConf config.RedisConfig
	client    *redisc.RedisClient
}

func NewRedis(conf config.RedisConfig) *Redis {
	return &Redis{
		redisConf: conf,
		client:    redisc.New(conf),
	}
}

func (r *Redis) Start(ctx context.Context) error {
	return nil
}

func (r *Redis) String() string {
	return "redis"
}

func (r *Redis) State(ctx context.Context) (string, error) {
	if err := redisc.Client.Do(context.Background(), redisc.Client.B().Ping().Build()).Error(); err != nil {
		return "UNHEALTHY", fmt.Errorf("redis ping failed: %w", err)
	}
	return "HEALTHY", nil
}

func (r *Redis) Terminate(ctx context.Context) error {
	redisc.Client.Close()
	return nil
}
