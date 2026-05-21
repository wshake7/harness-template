package router

import (
	"admin/internal/fiberc/middleware"
	"admin/internal/router/auth_router"
	"github.com/gofiber/fiber/v3"
)

type Router struct{}

func (r *Router) RegisterRouters(group fiber.Router) {
	group = group.Group("/api")
	defaultGroup := group.Use(
		middleware.TimestampMiddleware(),
		//middleware.NonceMiddleware(),
		middleware.LanguageMiddleware(),
	)
	registerAccountRouters(defaultGroup.Group("/account"))
	registerEncryptRouters(defaultGroup.Group("/encrypt"))
	auth_router.RegisterRouters(defaultGroup)
}
