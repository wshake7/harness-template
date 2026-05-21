package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysLoginLogRouters(router fiber.Router) {
	sysLoginLogHandler := logic.NewSysLoginLogHandler(query.Q)

	router.Post("/list", handler.CtxHandlerFunc(sysLoginLogHandler.List))
	router.Post("/detail", handler.CtxHandlerFunc(sysLoginLogHandler.Detail))
}
