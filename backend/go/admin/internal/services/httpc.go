package services

import (
	"admin/internal/services/httpc"
	"context"
	"go.uber.org/zap"
)

type Httpc struct {
}

func NewHttpc() *Httpc {
	return &Httpc{}
}

func (h *Httpc) Start(ctx context.Context) error {
	httpc.New(zap.S())
	return nil
}

func (h *Httpc) String() string {
	return "httpc"
}

func (h *Httpc) State(ctx context.Context) (string, error) {
	return "HEALTHY", nil
}

func (h *Httpc) Terminate(ctx context.Context) error {
	return httpc.Client.Close()
}
