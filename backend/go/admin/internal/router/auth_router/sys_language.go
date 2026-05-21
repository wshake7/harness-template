package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysLanguageRouters(router fiber.Router) {
	sysLanguageHandler := logic.NewSysLanguageHandler(query.Q)

	langType := router.Group("/type")
	langType.Post("/list", handler.CtxHandlerFunc(sysLanguageHandler.TypeList))
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"))
	typeCreateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"), middleware.WithChangeQuery(langTypeCreateChangeQuery))
	typeUpdateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"), middleware.WithChangeQuery(langTypeUpdateChangeQuery))
	entryCreateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"), middleware.WithChangeQuery(langEntryCreateChangeQuery))
	entryUpdateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"), middleware.WithChangeQuery(langEntryUpdateChangeQuery))
	entryBatchCreateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("language"), middleware.WithChangeQuery(langEntryBatchCreateChangeQuery))
	langType.Post("/create", typeCreateLogMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeCreate))
	langType.Post("/update", typeUpdateLogMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeUpdate))
	langType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.TypeDel))

	langEntry := router.Group("/entry")
	langEntry.Post("/list", handler.CtxHandlerFunc(sysLanguageHandler.EntryList))
	langEntry.Post("/create", entryCreateLogMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryCreate))
	langEntry.Post("/update", entryUpdateLogMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryUpdate))
	langEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryDel))
	langEntry.Post("/batch/create", entryBatchCreateLogMiddleware, handler.CtxHandlerNilFunc(sysLanguageHandler.EntryBatchCreate))
}
