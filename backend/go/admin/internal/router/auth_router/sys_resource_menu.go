package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysResourceMenuRouters(router fiber.Router) {
	sysResourceMenuHandler := logic.NewSysResourceMenuHandler(query.Q)
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_menu"))
	createLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_menu"), middleware.WithChangeQuery(resourceMenuCreateChangeQuery))
	updateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("resource_menu"), middleware.WithChangeQuery(resourceMenuUpdateChangeQuery))

	router.Post("/list", handler.CtxHandlerFunc(sysResourceMenuHandler.List))
	router.Get("/tree", handler.CtxFunc(sysResourceMenuHandler.Tree))
	router.Post("/create", createLogMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Create))
	router.Post("/update", updateLogMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysResourceMenuHandler.Del))
}
