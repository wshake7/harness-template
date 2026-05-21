package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/service"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysResourceApiRouters(router fiber.Router) {
	sysResourceApiHandler := logic.NewSysResourceApiHandler(query.Q, service.NewCasbinService())
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_api"))
	createLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_api"), middleware.WithChangeQuery(resourceApiCreateChangeQuery))
	updateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_api"), middleware.WithChangeQuery(resourceApiUpdateChangeQuery))

	router.Post("/list", handler.CtxHandlerFunc(sysResourceApiHandler.List))
	router.Post("/create", createLogMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Create))
	router.Post("/update", updateLogMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysResourceApiHandler.Del))
}
