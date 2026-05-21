package httpc

import (
	"go.uber.org/zap"
	"resty.dev/v3"
)

type HttpcClient struct {
	*resty.Client
}

var Client *HttpcClient

func New(logger *zap.SugaredLogger) *HttpcClient {
	Client = &HttpcClient{resty.New().SetLogger(logger.With(zap.String("module", "httpc")))}
	return Client
}

func (h *HttpcClient) RWith(logger *zap.SugaredLogger) {
	h.R().SetLogger(logger.With(zap.String("module", "httpc")))
}
