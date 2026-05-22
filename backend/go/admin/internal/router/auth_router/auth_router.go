package auth_router

import (
	"admin/internal/config"
	"admin/internal/fiberc/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRouters(router fiber.Router, conf *config.Config) {
	eventGroup := router.Use(middleware.AuthMiddleware(conf.Auth.TokenName), middleware.CasbinAPIMiddleware(), middleware.LanguageMiddleware(conf.DefaultLanguage))
	registerEventRouters(eventGroup)

	group := router.Use(middleware.AuthMiddleware(conf.Auth.TokenName), middleware.CasbinAPIMiddleware(), middleware.EncryptMiddleware(), middleware.LanguageMiddleware(conf.DefaultLanguage))
	registerSysRoleRouters(group.Group("/sys/role"))
	registerSysUserRouters(group.Group("/sys/user"))
	registerSysDictRouters(group.Group("/sys/dict"))
	registerSysLanguageRouters(group.Group("/sys/language"))
	registerSysApiLogRouters(group.Group("/sys/api/log"))
	registerSysLoginLogRouters(group.Group("/sys/login/log"))
	registerSysResourceMenuRouters(group.Group("/sys/resource/menu"))
	registerSysResourceApiRouters(group.Group("/sys/resource/api"))
	registerJobScheduleRouters(group.Group("/sys/job/schedule"), conf)
	registerJobExecutionRouters(group.Group("/sys/job/execution"))
}
