package asynq

import (
	"admin/internal/config"
	"github.com/hibiken/asynq"
)

type Asynq struct {
	*asynq.Client
}

var Client *Asynq

func New(conf config.RedisConfig) *Asynq {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Network:      "tcp",
		Addr:         conf.Addr[0],
		Username:     conf.Username,
		Password:     conf.Password,
		DB:           conf.SelectDB,
		DialTimeout:  conf.ConnDialTimeout,
		ReadTimeout:  conf.ConnReadTimeout,
		WriteTimeout: conf.ConnWriteTimeout,
		PoolSize:     conf.BlockingPoolSize,
		TLSConfig:    nil,
	})
	Client = &Asynq{client}
	return Client
}
