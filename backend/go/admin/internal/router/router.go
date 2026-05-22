package router

import (
	"admin/internal/config"
	"admin/internal/fiberc/middleware"
	"admin/internal/router/auth_router"
	"github.com/gofiber/fiber/v3"
)

type Router struct {
	Conf *config.Config
}

func (r *Router) RegisterRouters(group fiber.Router) {
	group = group.Group("/api")
	defaultGroup := group.Use(
		middleware.TimestampMiddleware(),
		//middleware.NonceMiddleware(),
		middleware.LanguageMiddleware(r.Conf.DefaultLanguage),
	)
	registerAccountRouters(defaultGroup.Group("/account"), r.Conf)
	registerEncryptRouters(defaultGroup.Group("/encrypt"))
	auth_router.RegisterRouters(defaultGroup, r.Conf)
}
