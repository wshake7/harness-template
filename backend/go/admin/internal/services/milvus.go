package services

import (
	"admin/internal/config"
	"admin/internal/services/milvusc"
	"context"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type Milvus struct {
	conf   config.MilvusConfig
	client *milvusclient.Client
}

func NewMilvus(conf config.MilvusConfig) *Milvus {
	return &Milvus{conf: conf}
}

func (m *Milvus) Start(ctx context.Context) error {
	if !m.conf.Enabled {
		return nil
	}
	client, err := milvusc.New(ctx, m.conf)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *Milvus) String() string {
	return "milvus"
}

func (m *Milvus) State(ctx context.Context) (string, error) {
	if !m.conf.Enabled {
		return "DISABLED", nil
	}
	if m.client == nil {
		return "UNHEALTHY", fmt.Errorf("milvus client not initialized")
	}
	if _, err := m.client.ListDatabase(ctx, milvusclient.NewListDatabaseOption()); err != nil {
		return "UNHEALTHY", fmt.Errorf("milvus health check failed: %w", err)
	}
	return "HEALTHY", nil
}

func (m *Milvus) Terminate(ctx context.Context) error {
	if !m.conf.Enabled {
		return nil
	}
	return milvusc.Close(ctx)
}
