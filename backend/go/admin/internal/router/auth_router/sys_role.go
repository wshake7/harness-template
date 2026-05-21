package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/service"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysRoleRouters(router fiber.Router) {
	sysRoleHandler := logic.NewSysRoleHandler(query.Q, service.NewCasbinService())
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("role"))
	createLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("role"), middleware.WithChangeQuery(sysRoleCreateChangeQuery))
	updateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("role"), middleware.WithChangeQuery(sysRoleUpdateChangeQuery))
	permissionLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("role"), middleware.WithChangeQuery(sysRolePermissionChangeQuery))

	router.Post("/list", handler.CtxHandlerFunc(sysRoleHandler.List))
	router.Get("/tree", handler.CtxFunc(sysRoleHandler.Tree))
	router.Get("/:id/permissions", handler.CtxHandlerFunc(sysRoleHandler.Permissions))
	router.Post("/create", createLogMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Create))
	router.Post("/update", updateLogMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Update))
	router.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.Del))
	router.Post("/permissions", permissionLogMiddleware, handler.CtxHandlerNilFunc(sysRoleHandler.SavePermissions))
}
