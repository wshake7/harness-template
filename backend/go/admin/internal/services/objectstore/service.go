package objectstore

import (
	"admin/internal/config"
	"context"
	"fmt"
)

var current Engine = disabledEngine{}

type Service struct {
	conf   config.StorageConfig
	engine Engine
}

func NewService(conf config.StorageConfig) *Service {
	return &Service{conf: conf}
}

func Current() Engine {
	if current == nil {
		return disabledEngine{}
	}
	return current
}

func (s *Service) Start(ctx context.Context) error {
	engine, err := NewEngine(s.conf)
	if err != nil {
		return err
	}
	if err := engine.Health(ctx); err != nil {
		return err
	}
	s.engine = engine
	current = engine
	return nil
}

func (s *Service) String() string {
	return "objectstore"
}

func (s *Service) State(ctx context.Context) (string, error) {
	engine := s.engine
	if engine == nil {
		engine = Current()
	}
	if engine.Name() == "disabled" {
		return "DISABLED", nil
	}
	if err := engine.Health(ctx); err != nil {
		return "UNHEALTHY", fmt.Errorf("object storage health check failed: %w", err)
	}
	return "HEALTHY", nil
}

func (s *Service) Terminate(ctx context.Context) error {
	s.engine = disabledEngine{}
	current = disabledEngine{}
	return nil
}
