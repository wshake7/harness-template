package auth_router

import (
	"admin/internal/appsvc"
	"admin/internal/config"
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/query"

	"github.com/gofiber/fiber/v3"
)

func registerStorageRouters(router fiber.Router, conf *config.Config) {
	storageHandler := logic.NewStorageFileHandler(
		appsvc.NewFileStorage(query.Q, conf.Storage, objectstore.Current()),
		conf.Storage.MaxUploadBytes,
	)

	fileGroup := router.Group("/file")
	fileGroup.Post("/upload", handler.CtxFunc(storageHandler.Upload))
	fileGroup.Post("/prepareUpload", handler.CtxHandlerFunc(storageHandler.PrepareUpload))
	fileGroup.Post("/completeUpload", handler.CtxHandlerFunc(storageHandler.CompleteUpload))
	fileGroup.Post("/detail", handler.CtxHandlerFunc(storageHandler.Detail))
	fileGroup.Post("/presigned", handler.CtxHandlerFunc(storageHandler.Presigned))
	fileGroup.Post("/del", handler.CtxHandlerNilFunc(storageHandler.Del))
}
