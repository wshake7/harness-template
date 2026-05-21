package middleware

import (
	"admin/internal/domains"
	"admin/internal/fiberc/handler"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDKey string = "trace_id"
)

func TraceMiddleware() fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		traceId := ctx.RequestID()
		if traceId == "" {
			traceId = uuid.NewString()
		}
		ctx.Set(domains.HeaderXRequestID, traceId)
		ctx.TraceId = traceId
		ctx.AddLogFields(zap.String(TraceIDKey, traceId), zap.String("method", ctx.Method()), zap.String("path", ctx.Path()))
		ctx.AddResLogFields(zap.String(TraceIDKey, traceId))
		return ctx.Next()
	})
}
