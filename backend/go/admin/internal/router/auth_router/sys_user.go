package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysUserRouters(router fiber.Router) {
	sysUserHandler := logic.NewSysUserHandler(query.Q)
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("user"))
	createLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("user"), middleware.WithChangeQuery(sysUserCreateChangeQuery))
	updateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("user"), middleware.WithChangeQuery(sysUserUpdateChangeQuery))

	router.Post("/list", handler.CtxHandlerFunc(sysUserHandler.List))
	router.Post("/create", createLogMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Create))
	router.Post("/update", updateLogMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysUserHandler.Del))
}
