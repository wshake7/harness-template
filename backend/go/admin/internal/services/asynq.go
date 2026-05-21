package services

import (
	"admin/internal/config"
	"admin/internal/services/asynq"
	"context"
)

type Asynq struct {
	conf   config.RedisConfig
	client *asynq.Asynq
}

func NewAsynq(conf config.RedisConfig) *Asynq {
	return &Asynq{conf: conf}
}

func (a *Asynq) Start(ctx context.Context) error {
	a.client = asynq.New(a.conf)
	return a.client.Ping()
}

func (a *Asynq) String() string {
	return "asynq"
}

func (a *Asynq) State(ctx context.Context) (string, error) {
	err := a.client.Ping()
	if err != nil {
		return "FAIL", err
	}
	return "HEALTHY", nil
}

func (a *Asynq) Terminate(ctx context.Context) error {
	return a.client.Close()
}
