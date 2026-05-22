package milvusc

import (
	"admin/internal/config"
	"context"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var Client *milvusclient.Client

func New(ctx context.Context, conf config.MilvusConfig) (*milvusclient.Client, error) {
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: conf.Address,
		APIKey:  conf.APIKey,
		DBName:  conf.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("create milvus client: %w", err)
	}
	Client = client
	return client, nil
}

func Close(ctx context.Context) error {
	if Client == nil {
		return nil
	}
	if err := Client.Close(ctx); err != nil {
		return fmt.Errorf("close milvus client: %w", err)
	}
	Client = nil
	return nil
}
