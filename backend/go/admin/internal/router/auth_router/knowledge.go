package auth_router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerKnowledgeRouters(router fiber.Router) {
	collectionHandler := logic.NewKnowledgeCollectionHandler(query.Q)
	documentHandler := logic.NewKnowledgeDocumentHandler(query.Q)

	collectionGroup := router.Group("/collection")
	collectionGroup.Post("/list", handler.CtxHandlerFunc(collectionHandler.List))
	collectionGroup.Post("/detail", handler.CtxHandlerFunc(collectionHandler.Detail))
	collectionGroup.Post("/create", handler.CtxHandlerNilFunc(collectionHandler.Create))
	collectionGroup.Post("/update", handler.CtxHandlerNilFunc(collectionHandler.Update))
	collectionGroup.Post("/del", handler.CtxHandlerNilFunc(collectionHandler.Del))
	documentGroup := router.Group("/document")
	documentGroup.Post("/list", handler.CtxHandlerFunc(documentHandler.List))
	documentGroup.Post("/listByCollection", handler.CtxHandlerFunc(documentHandler.ListByCollection))
	documentGroup.Post("/detail", handler.CtxHandlerFunc(documentHandler.Detail))
	documentGroup.Post("/create", handler.CtxHandlerNilFunc(documentHandler.Create))
	documentGroup.Post("/update", handler.CtxHandlerNilFunc(documentHandler.Update))
	documentGroup.Post("/del", handler.CtxHandlerNilFunc(documentHandler.Del))
}
