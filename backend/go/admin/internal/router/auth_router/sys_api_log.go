package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysApiLogRouters(router fiber.Router) {
	sysApiLogHandler := logic.NewSysApiLogHandler(query.Q)

	router.Post("/list", handler.CtxHandlerFunc(sysApiLogHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(sysApiLogHandler.Detail))
}
