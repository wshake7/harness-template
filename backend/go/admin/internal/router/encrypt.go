package router

import (
	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/service"

	"github.com/gofiber/fiber/v3"
)

func registerEncryptRouters(router fiber.Router) {
	encryptHandler := logic.NewEncryptHandler(service.NewRedisCache())
	router.Get("/public/key", handler.CtxFunc(encryptHandler.PublicKey))
}
