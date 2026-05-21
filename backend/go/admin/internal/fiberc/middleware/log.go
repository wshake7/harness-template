package middleware

import (
	"admin/internal/fiberc/handler"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func LogMiddleware(logger *zap.Logger) fiber.Handler {
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		ctx.SetLogger(logger)
		return ctx.Next()
	})
}
