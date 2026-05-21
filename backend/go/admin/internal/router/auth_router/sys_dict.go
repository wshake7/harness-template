package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/logic"
	"admin/internal/service"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerSysDictRouters(router fiber.Router) {
	sysDictHandler := logic.NewSysDictHandler(query.Q, service.NewDataPermissionService())

	dictType := router.Group("/type")
	dictType.Post("/list", handler.CtxHandlerFunc(sysDictHandler.TypeList))
	logMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"))
	typeCreateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"), middleware.WithChangeQuery(dictTypeCreateChangeQuery))
	typeUpdateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"), middleware.WithChangeQuery(dictTypeUpdateChangeQuery))
	entryCreateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"), middleware.WithChangeQuery(dictEntryCreateChangeQuery))
	entryUpdateLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"), middleware.WithChangeQuery(dictEntryUpdateChangeQuery))
	entryBatchCopyLogMiddleware := middleware.ApiLogMiddleware(middleware.WithModule("dict"), middleware.WithChangeQuery(dictEntryBatchCopyChangeQuery))
	dictType.Post("/create", typeCreateLogMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeCreate))
	dictType.Post("/update", typeUpdateLogMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeUpdate))
	dictType.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.TypeDel))

	dictEntry := router.Group("/entry")
	dictEntry.Post("/list", handler.CtxHandlerFunc(sysDictHandler.EntryList))
	dictEntry.Post("/match", handler.CtxHandlerFunc(sysDictHandler.EntryMatch))
	dictEntry.Post("/create", entryCreateLogMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryCreate))
	dictEntry.Post("/update", entryUpdateLogMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryUpdate))
	dictEntry.Post("/del", logMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryDel))
	dictEntry.Post("/batch/copy", entryBatchCopyLogMiddleware, handler.CtxHandlerNilFunc(sysDictHandler.EntryBatchCopy))
}
