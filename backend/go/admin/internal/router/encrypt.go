package router

import (
	"admin/internal/appsvc"
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"

	"github.com/gofiber/fiber/v3"
)

func registerEncryptRouters(router fiber.Router) {
	encryptHandler := logic.NewEncryptHandler(appsvc.NewRedisCache())
	router.Get("/public/key", handler.CtxFunc(encryptHandler.PublicKey))
}
