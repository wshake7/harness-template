package services

import (
	"admin/internal/services/casbin"
	"context"
	"gorm.io/gorm"
)

type Casbin struct {
	db *gorm.DB
}

func NewCasbin(db *gorm.DB) *Casbin {
	return &Casbin{db: db}
}

func (c *Casbin) Start(ctx context.Context) error {
	casbin.New(c.db)
	return nil
}

func (c *Casbin) String() string {
	return "casbin"
}

func (c *Casbin) State(ctx context.Context) (string, error) {
	return "HEALTHY", nil
}

func (c *Casbin) Terminate(ctx context.Context) error {
	return casbin.Adapter.Close()
}
